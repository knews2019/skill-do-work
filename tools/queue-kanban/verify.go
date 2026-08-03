package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Verify probe categories. Each finding carries one so callers (and tests) can
// name a probe instead of matching report prose.
const (
	verifyCategoryVersionChangelogMismatch     = "version-changelog-mismatch"
	verifyCategoryChangelogVersionNotAhead     = "changelog-version-not-ahead"
	verifyCategoryReusedChangelogTitle         = "reused-changelog-title"
	verifyCategoryDuplicateRequestId           = "duplicate-req-id"
	verifyCategoryMergedWorktreeLeftover       = "merged-worktree-leftover"
	verifyCategoryUnmergedWorktreeLeftover     = "unmerged-worktree-leftover"
	verifyCategoryUndeterminedWorktreeLeftover = "worktree-merge-state-undetermined"
	verifyCategoryWorktreeWroteQueueState      = "worktree-wrote-queue-state"
	verifyCategoryCheckpointGhostRequest       = "checkpoint-names-missing-req"
	verifyCategoryClaimNeedsAttention          = "claim-needs-attention"
	verifyCategoryStrandedFinishedRequest      = "stranded-finished-req"
)

// staleClaimThreshold is how long a `claimed` REQ may sit before verify reports
// it. It mirrors the takeover threshold in actions/work-reference.md → Crash
// Recovery (Step 1), and carries the same meaning there: it bounds how long a
// dead claim goes unnoticed. It is NOT a liveness test — a Route C REQ with a
// remediation loop can legitimately run longer — so this probe reports and never
// authorizes anything.
const staleClaimThreshold = 3 * time.Hour

// worktreeAgentNamePrefix is the naming convention worktree dispatch mode uses
// for a builder's worktree directory and branch (both share the string). Only
// names carrying it are in scope; a developer's own worktrees are never touched
// or reported.
const worktreeAgentNamePrefix = "worktree-agent-"

// requestIdMentionPattern finds every REQ-NNN token in free text — used to read
// the ids a CHECKPOINT.md names.
var requestIdMentionPattern = regexp.MustCompile(`REQ-\d+`)

// duplicateRequestIdWarningPrefix is how model.go's duplicate-id resolution
// phrases its warning. Verify reads those warnings rather than re-deriving the
// duplicate set — one definition, in the parser that already owns it.
const duplicateRequestIdWarningPrefix = "duplicate REQ id "

// VerifyFinding is one thing verify found wrong.
//
// Fixable means specifically: `do-work cleanup` can mechanically resolve it.
// Everything else is a human decision, and must not be advertised otherwise —
// an inflated fixable count sends the user to a command that will not help.
type VerifyFinding struct {
	Category string
	Detail   string
	Fixable  bool
	Remedy   string
}

// VerifyReport is the whole read-only result: what was wrong, and which probes
// could not run. A skipped probe is reported, never silently dropped — silence
// would read as "checked and clean."
type VerifyReport struct {
	RepoRoot      string
	Findings      []VerifyFinding
	SkippedProbes []string
}

// FixableCount is how many findings `do-work cleanup` can resolve mechanically.
func (report VerifyReport) FixableCount() int {
	fixableCount := 0
	for _, finding := range report.Findings {
		if finding.Fixable {
			fixableCount++
		}
	}
	return fixableCount
}

// ExitCode is 1 when there are findings, 0 otherwise. Skipped probes do not fail
// the run — a missing git binary is not a repo defect.
func (report VerifyReport) ExitCode() int {
	if len(report.Findings) > 0 {
		return 1
	}
	return 0
}

// runVerifyProbes checks the release-ritual and queue invariants that are
// otherwise verified by hand, and reports. It is strictly READ-ONLY: it opens
// files and shells out to git for read commands, and writes nothing anywhere —
// repairs belong to actions/cleanup.md, which asks before it acts.
//
// `now` is injected so claim-age assertions are deterministic in tests.
func runVerifyProbes(repoRootOverride string, now time.Time) (VerifyReport, error) {
	repoRoot, resolveError := resolveRepoRootOrDefault(repoRootOverride)
	if resolveError != nil {
		return VerifyReport{}, resolveError
	}
	report := VerifyReport{RepoRoot: repoRoot}

	board, buildError := buildBoard(repoRoot, now, defaultRecentWindow, lookupGitCommitDate)
	if buildError != nil {
		return VerifyReport{}, buildError
	}

	appendReleaseFindings(&report, repoRoot)
	appendDuplicateRequestIdFindings(&report, board)
	appendCheckpointGhostFindings(&report, repoRoot, board)
	appendClaimFindings(&report, board, now)
	appendStrandedFinishedFindings(&report, board)
	appendWorktreeFindings(&report, repoRoot)

	return report, nil
}

// appendReleaseFindings covers CLAUDE.md § Before Every Commit items 1 and 2 —
// the two cross-file checks a human is asked to perform every commit, and the two
// that have already been gotten wrong.
func appendReleaseFindings(report *VerifyReport, repoRoot string) {
	versionFilePath := resolveVersionFilePath(repoRoot, "")
	changelogPath := filepath.Join(repoRoot, "CHANGELOG.md")

	currentVersion, versionError := readCurrentVersion(versionFilePath)
	if versionError != nil {
		report.SkippedProbes = append(report.SkippedProbes,
			fmt.Sprintf("version-vs-changelog probes: %v", versionError))
		return
	}

	changelogEntries, changelogError := readChangelogEntries(changelogPath)
	if changelogError != nil {
		report.SkippedProbes = append(report.SkippedProbes,
			fmt.Sprintf("version-vs-changelog probes: no readable CHANGELOG.md at %s", changelogPath))
		return
	}
	if len(changelogEntries) == 0 {
		// The repo either has no entries yet or uses a different convention.
		// actions/work-reference.md → Changelog Entry Procedure says match the
		// existing format rather than impose the house one, so these two probes
		// have nothing to assert against and say so.
		report.SkippedProbes = append(report.SkippedProbes,
			fmt.Sprintf("version-vs-changelog probes: %s has no `## X.Y.Z — Title (YYYY-MM-DD)` entries (empty, or a different changelog convention)", changelogPath))
		return
	}

	newestEntry := changelogEntries[0]

	// The version file and the newest entry must AGREE. "Strictly greater" is the
	// rule while you are composing a release (CLAUDE.md § Before Every Commit item
	// 1: the new version must beat the newest entry that already exists) — the
	// moment the entry is written, the two are equal, and equality is the steady
	// state a check running at an arbitrary time should see. So the finding is a
	// mismatch, and its direction names the cause: ahead means a bump landed
	// without its entry (the pre-commit half-done state), behind means an entry
	// landed with a version the version file never received.
	comparison, compareError := compareSemanticVersions(currentVersion, newestEntry.Version)
	switch {
	case compareError != nil:
		report.SkippedProbes = append(report.SkippedProbes,
			fmt.Sprintf("version-vs-changelog probe: %v", compareError))
	case comparison > 0:
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryVersionChangelogMismatch,
			Detail: fmt.Sprintf("version %s is ahead of the newest CHANGELOG.md entry %s (%s) — a bump without its changelog entry",
				currentVersion, newestEntry.Version, newestEntry.Title),
			Remedy: "write the entry for " + currentVersion + ", or revert the bump (expected mid-release; a finding only if the release is done)",
		})
	case comparison < 0:
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryVersionChangelogMismatch,
			Detail: fmt.Sprintf("version %s is behind the newest CHANGELOG.md entry %s (%s) — an entry whose version was never written to the version file",
				currentVersion, newestEntry.Version, newestEntry.Title),
			Remedy: "set the version to " + newestEntry.Version + ", or correct the entry heading",
		})
	}

	// The "strictly greater" ordering that actually survives into the committed
	// state is WITHIN the changelog: the newest entry must beat every earlier one.
	// This is the check that catches the duplicate version numbers CLAUDE.md
	// records as having already happened.
	for _, earlierEntry := range changelogEntries[1:] {
		entryComparison, entryCompareError := compareSemanticVersions(newestEntry.Version, earlierEntry.Version)
		if entryCompareError != nil || entryComparison > 0 {
			continue
		}
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryChangelogVersionNotAhead,
			Detail: fmt.Sprintf("newest CHANGELOG.md entry %s is not strictly greater than the earlier entry %s (%s)",
				newestEntry.Version, earlierEntry.Version, earlierEntry.Date),
			Remedy: "renumber the newest entry — a duplicate or out-of-order version breaks every version-ordered reader",
		})
		break
	}

	for _, earlierEntry := range changelogEntries[1:] {
		if strings.EqualFold(earlierEntry.Title, newestEntry.Title) {
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryReusedChangelogTitle,
				Detail: fmt.Sprintf("newest entry title %q is already used by %s (%s)",
					newestEntry.Title, earlierEntry.Version, earlierEntry.Date),
				Remedy: "retitle the newest entry so a reader scanning headings can tell them apart",
			})
			break
		}
	}
}

// appendDuplicateRequestIdFindings lifts the parser's own duplicate-id warnings
// into findings. Two REQ files claiming one number is cheap to fix after the fact
// (the title is in the filename), which is why detection here is enough.
func appendDuplicateRequestIdFindings(report *VerifyReport, board *Board) {
	for _, warningText := range board.Warnings {
		if strings.HasPrefix(warningText, duplicateRequestIdWarningPrefix) {
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryDuplicateRequestId,
				Detail:   warningText,
				Remedy:   "renumber one of the two files (its title is in the filename, so the pair is easy to tell apart)",
			})
		}
	}
}

// appendCheckpointGhostFindings flags any REQ id do-work/CHECKPOINT.md names that
// no longer exists anywhere in the tree. The checkpoint is read by the next
// session — and since it is now Crash Recovery's input (actions/work-reference.md
// → Crash Recovery (Step 1)), a dangling id there is a stale premise, not just
// cosmetic staleness.
func appendCheckpointGhostFindings(report *VerifyReport, repoRoot string, board *Board) {
	checkpointPath := filepath.Join(repoRoot, "do-work", "CHECKPOINT.md")
	checkpointBytes, readError := os.ReadFile(checkpointPath)
	if readError != nil {
		return // no checkpoint is normal, not a finding
	}

	alreadyReported := map[string]bool{}
	for _, mentionedId := range requestIdMentionPattern.FindAllString(string(checkpointBytes), -1) {
		if alreadyReported[mentionedId] {
			continue
		}
		if _, exists := board.RequestsById[mentionedId]; exists {
			continue
		}
		alreadyReported[mentionedId] = true
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryCheckpointGhostRequest,
			Detail:   fmt.Sprintf("do-work/CHECKPOINT.md names %s, which does not exist in queue/, working/, or archive/", mentionedId),
			Remedy:   "edit or delete the checkpoint — a later session reads it, and crash recovery classifies against it",
		})
	}
}

// appendClaimFindings reports a claim that has sat past the threshold, and any
// claimed REQ whose claimed_at cannot be trusted. A bad stamp is reported rather
// than ignored for the same reason actions/work-reference.md's takeover guard
// treats it as immediately eligible: a meaningless age must push toward a human
// looking, not toward silence.
func appendClaimFindings(report *VerifyReport, board *Board, now time.Time) {
	skewHorizon := now.Add(futureTimestampSkewAllowance)
	for _, claimedTicket := range board.Columns.Claimed {
		rawClaimStamp := strings.TrimSpace(claimedTicket.ClaimedAt)
		if rawClaimStamp == "" {
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryClaimNeedsAttention,
				Detail:   fmt.Sprintf("%s is claimed but carries no claimed_at — its age cannot be known", claimedTicket.RequestId),
				Remedy:   "see actions/forensics.md Check 1 for the manual reset",
			})
			continue
		}
		claimInstant, parsedOk := parseTimestamp(rawClaimStamp)
		if !parsedOk {
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryClaimNeedsAttention,
				Detail:   fmt.Sprintf("%s has an unparseable claimed_at (%q)", claimedTicket.RequestId, rawClaimStamp),
				Remedy:   "fix the stamp to a UTC ISO-8601 instant, or reset the REQ (actions/forensics.md Check 1)",
			})
			continue
		}
		if claimInstant.After(skewHorizon) {
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryClaimNeedsAttention,
				Detail: fmt.Sprintf("%s has a future-dated claimed_at (%s) — usually local wall-clock time written with a Z suffix",
					claimedTicket.RequestId, rawClaimStamp),
				Remedy: "re-stamp it with `date -u +%Y-%m-%dT%H:%M:%SZ` (the Timestamp rule in actions/work-reference.md)",
			})
			continue
		}
		if claimAge := now.Sub(claimInstant); claimAge >= staleClaimThreshold {
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryClaimNeedsAttention,
				Detail: fmt.Sprintf("%s has been claimed for %s (threshold %s) — reported, not judged dead",
					claimedTicket.RequestId, formatApproximateDuration(claimAge), formatApproximateDuration(staleClaimThreshold)),
				Remedy: "a long build can legitimately exceed this; take it over only if you know the run is gone (actions/forensics.md Check 1)",
			})
		}
	}
}

// appendStrandedFinishedFindings reports a REQ with a terminal status still
// sitting in queue/ or working/ — actions/forensics.md Check 9's definition, and
// the one finding class besides orphan worktrees that cleanup resolves
// mechanically (Pass 0 sweeps it into the archive).
func appendStrandedFinishedFindings(report *VerifyReport, board *Board) {
	for _, ticket := range board.AllRequests {
		if !isTerminalResolvedStatus(ticket.Status) {
			continue
		}
		treeSection := ticket.TreeSection
		if treeSection != "queue" && treeSection != "working" {
			continue
		}
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryStrandedFinishedRequest,
			Detail:   fmt.Sprintf("%s has terminal status %q but still sits in do-work/%s/", ticket.RequestId, ticket.Status, treeSection),
			Fixable:  true,
			Remedy:   "cleanup Pass 0 moves it into do-work/archive/",
		})
	}
}

// worktreeMergeState is what verify can honestly say about a worktree-agent-*
// leftover: its branch is already contained in the integration branch, it is
// not, or git could not answer. It says nothing about whether a builder is still
// running — see classifyWorktreeMergeState.
type worktreeMergeState int

const (
	worktreeMergeStateMerged worktreeMergeState = iota
	worktreeMergeStateUnmerged
	worktreeMergeStateUndetermined
)

// classifyWorktreeMergeState answers whether leftoverName is already contained in
// the integration branch, which it reads as the repo-root checkout's HEAD.
//
// Merged-ness is HEAD-relative — `git branch -d`'s own trap, documented at
// actions/work-reference.md → Worktree Dispatch Mode, "Cleanup — happy path":
// asked from an unrelated checkout, a perfectly merged branch reads unmerged and
// an unmerged one can read merged. `git -C repoRoot` pins the question to the
// main checkout the orchestrator merges into, never to a builder's worktree.
//
// It CANNOT tell a builder that is still running from one that died and left this
// behind. There is no lock, heartbeat, or claim registry to ask, and REQ-073
// forbids adding one; a time threshold is not a stand-in either (see
// staleClaimThreshold's own doc comment). So the unmerged case names the
// still-in-flight possibility in its remedy instead of guessing between them.
func classifyWorktreeMergeState(repoRoot string, leftoverName string) worktreeMergeState {
	command := exec.Command("git", "-C", repoRoot, "merge-base", "--is-ancestor", leftoverName, "HEAD")
	runError := command.Run()
	if runError == nil {
		return worktreeMergeStateMerged
	}
	// Exit 1 is git's answer "not an ancestor". Anything else — most often exit
	// 128 for a worktree whose branch is gone — is git declining to answer, which
	// is not the same claim and must not be reported as one.
	if exitError, isExitError := runError.(*exec.ExitError); isExitError && exitError.ExitCode() == 1 {
		return worktreeMergeStateUnmerged
	}
	return worktreeMergeStateUndetermined
}

// routeWorktreeLeftover maps a merge state onto the finding it produces.
//
// Fixable is true for merged residue and nothing else, because that is the only
// state actions/cleanup.md → Pass 5 resolves mechanically: `git worktree remove`
// plus `git branch -d`, neither forcing. Every other state lands on Pass 5's
// consent-gated path, where the pass "stops being mechanical" and asks — a human
// decision, which VerifyFinding.Fixable's doc comment says must not be advertised
// otherwise.
//
// An unmerged leftover stays a reported finding during a live run rather than
// being suppressed while builders are in flight. VerifyReport's doc comment is
// explicit that silence reads as "checked and clean," and verify has no way to
// know a run is active (see classifyWorktreeMergeState) — so suppression would
// have to guess, and would hide genuinely stranded work whenever it guessed
// wrong. This mirrors how version-changelog-mismatch handles its own expected
// mid-release state: reported, with the transient case named in the remedy, and
// not fixable.
func routeWorktreeLeftover(mergeState worktreeMergeState) (category string, fixable bool, remedy string) {
	switch mergeState {
	case worktreeMergeStateMerged:
		return verifyCategoryMergedWorktreeLeftover, true,
			"cleanup Pass 5 removes it mechanically — the branch is already contained in HEAD, so nothing is lost"
	case worktreeMergeStateUnmerged:
		return verifyCategoryUnmergedWorktreeLeftover, false,
			"this is either a builder still in flight or work that outlived a dead run — verify cannot tell those apart. Leave it alone during a run; otherwise cleanup Pass 5 asks before discarding it, because the branch may hold the only copy"
	default:
		return verifyCategoryUndeterminedWorktreeLeftover, false,
			"git could not say whether this is merged (typically a worktree whose branch is gone) — inspect it by hand; cleanup Pass 5 deletes nothing it cannot establish a merge target for"
	}
}

// appendWorktreeFindings covers the two worktree-dispatch invariants: no
// `worktree-agent-*` leftovers should outlive their run, and a builder must never
// write queue state (actions/work-reference.md → Worktree Dispatch Mode, "state
// stays home" and "sole integrator").
//
// Leftovers are classified by merge state rather than reported as one kind of
// thing, so the report routes only the mechanically-resolvable ones to cleanup.
func appendWorktreeFindings(report *VerifyReport, repoRoot string) {
	if !gitBinaryAvailable() {
		report.SkippedProbes = append(report.SkippedProbes, "worktree probes: git is not on PATH")
		return
	}

	worktreePathsByBranch, listError := listWorktreeAgentWorktrees(repoRoot)
	if listError != nil {
		report.SkippedProbes = append(report.SkippedProbes,
			fmt.Sprintf("worktree probes: %v", listError))
		return
	}
	agentBranches := listWorktreeAgentBranches(repoRoot)

	reportedNames := map[string]bool{}
	orderedNames := make([]string, 0, len(worktreePathsByBranch)+len(agentBranches))
	for name := range worktreePathsByBranch {
		orderedNames = append(orderedNames, name)
	}
	for _, branchName := range agentBranches {
		if _, alreadyListed := worktreePathsByBranch[branchName]; !alreadyListed {
			orderedNames = append(orderedNames, branchName)
		}
	}
	sort.Strings(orderedNames)

	for _, leftoverName := range orderedNames {
		if reportedNames[leftoverName] {
			continue
		}
		reportedNames[leftoverName] = true
		worktreePath, hasWorktree := worktreePathsByBranch[leftoverName]
		locationDetail := "branch only (no worktree)"
		if hasWorktree {
			locationDetail = worktreePath
		}
		category, fixable, remedy := routeWorktreeLeftover(classifyWorktreeMergeState(repoRoot, leftoverName))
		report.Findings = append(report.Findings, VerifyFinding{
			Category: category,
			Detail:   fmt.Sprintf("%s%s exists — %s", worktreeAgentNamePrefix, strings.TrimPrefix(leftoverName, worktreeAgentNamePrefix), locationDetail),
			Fixable:  fixable,
			Remedy:   remedy,
		})

		if !hasWorktree {
			continue
		}
		if dirtyQueuePaths := worktreeDirtyQueueState(worktreePath); len(dirtyQueuePaths) > 0 {
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryWorktreeWroteQueueState,
				Detail: fmt.Sprintf("%s has uncommitted changes under do-work/ (%s) — a builder wrote queue state the orchestrator alone owns",
					worktreePath, strings.Join(dirtyQueuePaths, ", ")),
				Remedy: "discard those changes in the worktree; every claim, status flip, and archive move belongs to the main tree",
			})
		}
	}
}

// listWorktreeAgentWorktrees maps each worktree-agent-* name to its worktree
// path, read from `git worktree list --porcelain`. The name is the directory's
// basename, which worktree dispatch mode keeps identical to the branch name.
func listWorktreeAgentWorktrees(repoRoot string) (map[string]string, error) {
	command := exec.Command("git", "-C", repoRoot, "worktree", "list", "--porcelain")
	output, runError := command.Output()
	if runError != nil {
		return nil, fmt.Errorf("`git worktree list` failed (no worktree support, or not a git repo)")
	}
	pathsByName := map[string]string{}
	for _, line := range strings.Split(string(output), "\n") {
		if !strings.HasPrefix(line, "worktree ") {
			continue
		}
		worktreePath := strings.TrimSpace(strings.TrimPrefix(line, "worktree "))
		worktreeName := filepath.Base(worktreePath)
		if strings.HasPrefix(worktreeName, worktreeAgentNamePrefix) {
			pathsByName[worktreeName] = worktreePath
		}
	}
	return pathsByName, nil
}

// listWorktreeAgentBranches returns every local worktree-agent-* branch name. A
// branch can outlive its worktree (the REQ archived, the branch did not), which
// is why branches are enumerated separately from worktrees.
func listWorktreeAgentBranches(repoRoot string) []string {
	command := exec.Command("git", "-C", repoRoot, "branch", "--list", worktreeAgentNamePrefix+"*", "--format=%(refname:short)")
	output, runError := command.Output()
	if runError != nil {
		return nil
	}
	var branchNames []string
	for _, line := range strings.Split(string(output), "\n") {
		branchName := strings.TrimSpace(line)
		if branchName != "" {
			branchNames = append(branchNames, branchName)
		}
	}
	return branchNames
}

// worktreeDirtyQueueState returns the do-work/ paths a builder's worktree has
// modified. Uncommitted changes there are the detectable signature of the "state
// stays home" rule being broken; a stale committed snapshot (which a worktree
// legitimately carries where the consumer commits do-work/) is not.
func worktreeDirtyQueueState(worktreePath string) []string {
	command := exec.Command("git", "-C", worktreePath, "status", "--porcelain", "--untracked-files=all", "--", "do-work/")
	output, runError := command.Output()
	if runError != nil {
		return nil
	}
	var dirtyPaths []string
	for _, line := range strings.Split(string(output), "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		statusFields := strings.Fields(trimmedLine)
		dirtyPaths = append(dirtyPaths, statusFields[len(statusFields)-1])
	}
	return dirtyPaths
}

// formatApproximateDuration renders a duration as "4h20m" / "35m" — short enough
// for a report line, precise enough to judge a claim by.
func formatApproximateDuration(duration time.Duration) string {
	wholeHours := int(duration.Hours())
	remainingMinutes := int(duration.Minutes()) - wholeHours*60
	if wholeHours == 0 {
		return fmt.Sprintf("%dm", remainingMinutes)
	}
	return fmt.Sprintf("%dh%dm", wholeHours, remainingMinutes)
}

// renderVerifyReport formats the report for a terminal. The fixable count and its
// pointer at `do-work cleanup` are the last line, so a single cheap invocation
// tells the user what to run next without this tool ever repairing anything.
func renderVerifyReport(report VerifyReport) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "queue-kanban verify — %s\n", report.RepoRoot)

	if len(report.Findings) == 0 {
		fmt.Fprintf(&builder, "  OK: no findings\n")
	}
	for _, finding := range report.Findings {
		fixableMarker := ""
		if finding.Fixable {
			fixableMarker = " [fixable]"
		}
		fmt.Fprintf(&builder, "  ! %s%s: %s\n", finding.Category, fixableMarker, finding.Detail)
		if finding.Remedy != "" {
			fmt.Fprintf(&builder, "      → %s\n", finding.Remedy)
		}
	}
	for _, skippedProbe := range report.SkippedProbes {
		fmt.Fprintf(&builder, "  - skipped %s\n", skippedProbe)
	}
	if fixableCount := report.FixableCount(); fixableCount > 0 {
		fmt.Fprintf(&builder, "  %d fixable: run do-work cleanup\n", fixableCount)
	}
	return builder.String()
}
