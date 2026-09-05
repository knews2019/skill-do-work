package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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
	verifyCategoryStructurallyDamagedRequest   = "structurally-damaged-req"
	verifyCategoryUnrecognizedRequestStatus    = "unrecognized-req-status"
	verifyCategoryMergedWorktreeLeftover       = "merged-worktree-leftover"
	verifyCategoryUnmergedWorktreeLeftover     = "unmerged-worktree-leftover"
	verifyCategoryUndeterminedWorktreeLeftover = "worktree-merge-state-undetermined"
	// The three states in which a worktree-agent-* worktree is PRESENT rather
	// than leftover: its branch is merged, but something the repository records
	// says the worktree itself is not disposable. None of them is fixable.
	verifyCategoryWorktreePresentUncommittedWork = "worktree-present-uncommitted-work"
	verifyCategoryWorktreePresentRunInFlight     = "worktree-present-run-in-flight"
	verifyCategoryWorktreePresentStateUnknown    = "worktree-present-state-unknown"
	verifyCategoryWorktreeWroteQueueState        = "worktree-wrote-queue-state"
	verifyCategoryWorktreeCommittedQueueState    = "worktree-committed-queue-state"
	verifyCategoryCheckpointGhostRequest         = "checkpoint-names-missing-req"
	verifyCategoryClaimNeedsAttention            = "claim-needs-attention"
	verifyCategoryStrandedFinishedRequest        = "stranded-finished-req"
	verifyCategoryStrayRequestFile               = "stray-req-file"
	verifyCategoryAssignedElsewhereClaimedHere   = "assigned-elsewhere-claimed-here"
	verifyCategoryArchivedUserRequestLiveMember  = "ur-archived-with-live-member"
	verifyCategoryCompletionAnomaly              = "completion-anomaly"
	verifyCategoryTimestampOrdering              = "timestamp-ordering"
	verifyCategoryCalibrationLogMismatch         = "calibration-log-mismatch"
	verifyCategoryCalibrationRowUnreconcilable   = "calibration-row-unreconcilable"
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
// Subject is the thing the finding is about — a worktree name, a REQ id, a file
// path — and it exists so a reader can see several findings about one thing as
// one block. It is set by the probe that knows, and stays empty when a probe has
// no natural subject; the board groups on an exact string match and never parses
// it back out of Detail, which is prose and free to change.
type VerifyFinding struct {
	Category string
	Detail   string
	Subject  string
	Fixable  bool
	Remedy   string
}

// VerifyReport is the whole read-only result: what was wrong, which probes could not run,
// and which do not apply to this repo at all. A skipped probe is reported, never silently
// dropped — silence would read as "checked and clean."
//
// SkippedProbes and NotApplicableProbes are different claims and are rendered separately.
// Skipped means "this invariant is real here and went unverified" — someone should act.
// Not applicable means "this invariant is not this repo's" — there is nothing to act on,
// and filing it as a skip trains readers to scroll past the skipped section, including on
// the day a genuinely skipped probe means something (REQ-282).
type VerifyReport struct {
	RepoRoot            string
	Findings            []VerifyFinding
	SkippedProbes       []string
	NotApplicableProbes []string
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

// ExitCode is 1 when there are findings, 0 otherwise. Neither skipped nor not-applicable
// probes fail the run — a missing git binary is not a repo defect, and neither is a repo
// that simply is not the suite.
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

	board, buildError := buildBoard(repoRoot, now, defaultRecentWindow, lookupGitCommitDate)
	if buildError != nil {
		return VerifyReport{}, buildError
	}

	return collectVerifyFindings(repoRoot, board, now), nil
}

// collectVerifyFindings is runVerifyProbes' body over a board the caller already
// built. The split exists so the board's own producer can emit these findings into
// the page without building the board a second time per request (REQ-284) — the
// board build is the expensive half, the probes are the cheap half.
//
// `now` is a parameter rather than read here for the same reason it always was:
// claim-age findings must be deterministic in tests, and — since serve calls this
// outside its mtime cache — a claim must be able to cross the staleness threshold
// on a tree where no file has changed. Passing a stale `now` would silently restore
// the blind spot the split was made to remove.
func collectVerifyFindings(repoRoot string, board *Board, now time.Time) VerifyReport {
	report := VerifyReport{RepoRoot: repoRoot}

	appendReleaseFindings(&report, repoRoot)
	appendDuplicateRequestIdFindings(&report, board)
	appendStructuralDamageFindings(&report, board)
	appendUnrecognizedStatusFindings(&report, board)
	appendCompletionAnomalyFindings(&report, repoRoot, board)
	appendTimestampOrderingFindings(&report, board)
	appendCalibrationLogFindings(&report, repoRoot, board)
	appendCheckpointGhostFindings(&report, repoRoot, board)
	appendClaimFindings(&report, board, now)
	appendStrandedFinishedFindings(&report, board)
	appendStrayRequestFileFindings(&report, board)
	appendAssignedElsewhereFindings(&report, board)
	appendArchivedUserRequestLiveMemberFindings(&report, board)
	appendWorktreeFindings(&report, repoRoot, board)

	return report
}

// appendStrayRequestFileFindings forwards the board walker's structured
// evidence for REQ files outside queue/working/archive. Strays remain outside
// AllRequests and every normal request probe; verify neither parses warning prose
// nor performs a second filesystem walk.
func appendStrayRequestFileFindings(report *VerifyReport, board *Board) {
	for _, strayRequest := range board.StrayRequestFiles {
		relativePath := filepath.ToSlash(strayRequest.RelativePath)
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryStrayRequestFile,
			Subject:  "do-work/" + relativePath,
			Detail: fmt.Sprintf("REQ file at do-work/%s is outside queue/, working/, and archive/ and is invisible to normal request/card probes",
				relativePath),
			Remedy: "inspect the REQ, then move it to do-work/archive/ if resolved or do-work/queue/ if it still needs work; verify is read-only and relocation is a human decision",
		})
	}
}

// releaseFindingSubject is the one thing all three release probes are about, so
// the board prints that heading once instead of repeating it per row.
const releaseFindingSubject = "CHANGELOG.md"

// appendReleaseFindings covers the two release invariants the commit ritual asks a
// human to verify by eye every time: the version file must agree with the newest
// CHANGELOG.md entry, and that entry's title must not already be in use. They are
// the two cross-file checks, and the two that have already been gotten wrong.
func appendReleaseFindings(report *VerifyReport, repoRoot string) {
	// Not the writer's resolver: this one knows the modular suite layout, and it separates
	// "no version file anywhere" (this is not a suite checkout, so the probes do not apply)
	// from every failure reachable after one resolves (a real unverified invariant, still
	// skipped). Collapsing those two is what would let the not-applicable path silence the
	// probes in the very repo that owns them — the failure REQ-282 exists to end.
	versionFilePath, isSuiteCheckout := resolveReleaseProbeVersionFilePath(repoRoot)
	if !isSuiteCheckout {
		report.NotApplicableProbes = append(report.NotApplicableProbes,
			"release probes: they verify the suite's own release ritual, and this repo root is not a suite checkout")
		return
	}
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
	// rule while you are composing a release (the new version must beat the newest
	// entry that already exists) — the
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
			Subject:  releaseFindingSubject,
			Detail: fmt.Sprintf("version %s is ahead of the newest CHANGELOG.md entry %s (%s) — a bump without its changelog entry",
				currentVersion, newestEntry.Version, newestEntry.Title),
			Remedy: "write the entry for " + currentVersion + ", or revert the bump (expected mid-release; a finding only if the release is done)",
		})
	case comparison < 0:
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryVersionChangelogMismatch,
			Subject:  releaseFindingSubject,
			Detail: fmt.Sprintf("version %s is behind the newest CHANGELOG.md entry %s (%s) — an entry whose version was never written to the version file",
				currentVersion, newestEntry.Version, newestEntry.Title),
			Remedy: "set the version to " + newestEntry.Version + ", or correct the entry heading",
		})
	}

	// The "strictly greater" ordering that actually survives into the committed
	// state is WITHIN the changelog: the newest entry must beat every earlier one.
	// This is the check that catches duplicate version numbers, which have already
	// happened in this repo more than once.
	for _, earlierEntry := range changelogEntries[1:] {
		entryComparison, entryCompareError := compareSemanticVersions(newestEntry.Version, earlierEntry.Version)
		if entryCompareError != nil || entryComparison > 0 {
			continue
		}
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryChangelogVersionNotAhead,
			Subject:  releaseFindingSubject,
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
				Subject:  releaseFindingSubject,
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

// appendStructuralDamageFindings reports a REQ file whose frontmatter no longer
// holds the fields the whole pipeline reads off it. The parser recovers from
// every shape below rather than erroring (splitFrontmatter and
// parseFrontmatterFields both document recovery as their contract), so the damage
// surfaces as EMPTY FIELDS on a card that still renders — which is why verify
// printed `OK: no findings` on a queue where six of seven files were broken.
// Detection here changes nothing about that leniency: a damaged REQ still parses
// and still appears on the board, it just stops passing the mechanical check.
//
// Like appendCompletionAnomalyFindings and the stray-file probe, this forwards
// what buildBoard already produced — the ticket's own parsed fields and its
// retained frontmatter bytes. No warning prose is matched and the tree is not
// walked a second time.
//
// A file with no leading fence reports ONE finding and stops. Its id, status and
// user_request are all empty as a CONSEQUENCE of the missing fence, and the
// remedy for the fence repairs all of them; reporting each separately would be
// four findings for one defect, the same double-report that
// appendTimestampOrderingFindings' outer-pair carve-out avoids.
func appendStructuralDamageFindings(report *VerifyReport, board *Board) {
	for _, ticket := range board.AllRequests {
		if ticket.FrontmatterMarkdown == "" {
			missingFenceDetail := fmt.Sprintf("%s has no leading frontmatter fence, so id, status, user_request and every other field parsed empty (the id named here was recovered from the filename)",
				ticket.RequestId)
			missingFenceRemedy := "restore the opening `---` as the file's very first line and the closing `---` after the last field, then re-check the fields it was hiding"
			if bodyStartsWithOpeningFrontmatterFence(ticket.BodyMarkdown) {
				missingFenceDetail = fmt.Sprintf("%s has an opening frontmatter fence but no closing fence, so id, status, user_request and every other field parsed empty (the id named here was recovered from the filename)",
					ticket.RequestId)
				missingFenceRemedy = "restore the closing `---` after the last frontmatter field, then re-check the fields it was hiding"
			}
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryStructurallyDamagedRequest,
				Subject:  ticket.RequestId,
				Detail:   missingFenceDetail,
				Remedy:   missingFenceRemedy,
			})
			continue
		}

		frontmatterFields := requestFrontmatterFields(ticket)
		if coerceScalarToString(frontmatterFields["id"]) == "" {
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryStructurallyDamagedRequest,
				Subject:  ticket.RequestId,
				Detail: fmt.Sprintf("%s has an empty or absent id: field — caution: its id was recovered from the filename, so renaming the file silently renumbers the REQ",
					ticket.RequestId),
				Remedy: fmt.Sprintf("write `id: %s` into the frontmatter", ticket.RequestId),
			})
		}

		if ticket.UserRequestId != "" {
			continue
		}
		if userRequestMayBeAbsent(frontmatterFields) {
			continue
		}
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryStructurallyDamagedRequest,
			Subject:  ticket.RequestId,
			Detail: fmt.Sprintf("%s carries no user_request: pointer, so it belongs to no UR — nothing links it to the input that asked for it, and UR closure cannot see it",
				ticket.RequestId),
			Remedy: "write `user_request: UR-NNN` naming the UR this REQ was captured under (do-work/user-requests/); documented stakeholder, code-review, scoped review-generated, and context_ref shapes are exempt",
		})
	}
}

// bodyStartsWithOpeningFrontmatterFence recognizes the exact first physical line
// that splitFrontmatter accepts as an opening fence, including BOM and CRLF forms.
// When no closing fence exists splitFrontmatter leaves the whole raw file in
// BodyMarkdown, so this can name the missing closing delimiter without a second
// filesystem walk. A hand-built ticket with an empty body remains a missing-opening
// case, preserving the structured-evidence test seam.
func bodyStartsWithOpeningFrontmatterFence(bodyMarkdown string) bool {
	if strings.HasPrefix(bodyMarkdown, "\ufeff") {
		bodyMarkdown = strings.TrimPrefix(bodyMarkdown, "\ufeff")
	}
	return strings.HasPrefix(bodyMarkdown, "---\n") || strings.HasPrefix(bodyMarkdown, "---\r\n")
}

// stakeholderMarkerFieldName is the frontmatter key that marks a stakeholder REQ —
// the fold discriminator in actions/work-reference.md → Stakeholder REQ Template,
// and the only positive evidence that a missing user_request is deliberate rather
// than damage.
const stakeholderMarkerFieldName = "stakeholder"

// userRequestMayBeAbsent identifies the documented REQ schemas that deliberately
// lack UR membership. Each exemption requires affirmative frontmatter evidence,
// never a directory name: a damaged ordinary REQ cannot cheaply look like one of
// these shapes by merely being archived.
func userRequestMayBeAbsent(frontmatterFields map[string]any) bool {
	if coerceScalarToString(frontmatterFields[stakeholderMarkerFieldName]) != "" {
		return true
	}
	if coerceScalarToString(frontmatterFields["source"]) == "code-review" {
		return true
	}
	if coerceScalarToString(frontmatterFields["review_generated"]) == "true" &&
		coerceScalarToString(frontmatterFields["scope"]) != "" {
		return true
	}
	return coerceScalarToString(frontmatterFields["context_ref"]) != ""
}

// requestFrontmatterFields re-reads one ticket's OWN retained frontmatter bytes
// through the same two parsers buildBoard used. This is neither a second tree walk
// nor a second parser: RequestTicket.FrontmatterMarkdown is the original fence
// bytes buildBoard already sliced off the file and kept on the ticket, and both
// functions called here are the production ones.
//
// It exists because two of the fields this probe needs are the ones the ticket
// cannot answer for. RequestTicket.RequestId falls back to the filename, so an
// empty `id:` is indistinguishable from a present one by the time it reaches a
// ticket; `stakeholder:` is not parsed onto the ticket at all. Promoting either to
// a RequestTicket field is model.go's call.
//
// parseFrontmatterFields never returns an error (recovery is its contract), so a
// block strict YAML rejects still yields whatever the lenient line scan recovered —
// the leniency this REQ was told to keep.
func requestFrontmatterFields(ticket *RequestTicket) map[string]any {
	yamlText, _, _, hasFrontmatter := splitFrontmatter(ticket.FrontmatterMarkdown)
	if !hasFrontmatter {
		return nil
	}
	frontmatterFields, _ := parseFrontmatterFields(yamlText)
	return frontmatterFields
}

// appendUnrecognizedStatusFindings lifts the parser's own unrecognized-status
// verdict into a finding, the way appendDuplicateRequestIdFindings lifts the
// duplicate-id one — but off RequestTicket.StatusUnrecognized, the structured form
// of that verdict, rather than off the warning sentence bucketColumns writes beside
// it. Status is the field the entire pipeline routes on: a REQ whose status did not
// survive parsing is parked in Needs input / Blocked, where it is indistinguishable
// from a REQ that is genuinely blocked, and nothing will ever claim it.
//
// A file with no frontmatter fence at all is left to
// appendStructuralDamageFindings: its status is empty because the fence is gone,
// and the fence finding already carries the repair.
func appendUnrecognizedStatusFindings(report *VerifyReport, board *Board) {
	for _, ticket := range board.AllRequests {
		if !ticket.StatusUnrecognized || ticket.FrontmatterMarkdown == "" {
			continue
		}
		detail := fmt.Sprintf("%s has an unrecognized status: value %q", ticket.RequestId, ticket.OriginalStatus)
		if strings.TrimSpace(ticket.OriginalStatus) == "" {
			detail = fmt.Sprintf("%s has an empty or absent status: field", ticket.RequestId)
		}
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryUnrecognizedRequestStatus,
			Subject:  ticket.RequestId,
			Detail:   detail + " — it is parked under Needs input / Blocked, indistinguishable from a genuinely blocked REQ, and no run will claim it",
			Remedy:   "edit `status:` to a Schema Read Contract value (actions/work-reference.md), or run do-work forensics to route it",
		})
	}
}

// appendCompletionAnomalyFindings lifts the board's completion-anomaly tickets
// into findings, so a broken terminal record fails the mechanical check instead
// of passing an `OK: no findings` while the summary shows a flagged strip
// (REQ-214 — verify was blind to every anomaly class until then). It forwards
// buildBoard's structured evidence, exactly like the stray-file probe: no
// warning-prose parsing, no second walk. The per-ticket reason already names
// the broken field(s) and the fix, so the remedy stays generic routing.
//
// One class is environment-sensitive: a ticket whose ONLY problem is a commit
// hash git could not date (completed_at absent) is indistinguishable from a
// healthy record when git itself is missing or the tree is not a repository —
// the same dating probe would fail for every valid hash. Per the ExitCode
// contract, availability gaps are skipped probes, never findings; every other
// class (unparseable completed_at, reversed span, neither field) is a genuine
// on-disk defect regardless of git and still fails the check.
func appendCompletionAnomalyFindings(report *VerifyReport, repoRoot string, board *Board) {
	commitDatingUsable := gitBinaryAvailable()
	if commitDatingUsable {
		if _, gitDirProbeError := exec.Command("git", "-C", repoRoot, "rev-parse", "--git-dir").Output(); gitDirProbeError != nil {
			commitDatingUsable = false
		}
	}
	for _, anomalousTicket := range board.Columns.CompletionAnomalies {
		if !commitDatingUsable && anomalousTicket.CompletedAt == "" && anomalousTicket.CommitHash != "" {
			report.SkippedProbes = append(report.SkippedProbes, fmt.Sprintf(
				"completion-anomaly probe: %s carries only a commit hash and git/the repository is unavailable — cannot distinguish a bad hash from an undatable valid one",
				anomalousTicket.RequestId))
			continue
		}
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryCompletionAnomaly,
			Subject:  anomalousTicket.RequestId,
			Detail: fmt.Sprintf("%s (status %s): %s",
				anomalousTicket.RequestId, anomalousTicket.Status, anomalousTicket.CompletionAnomalyReason),
			Remedy: "repair the named frontmatter field(s) in the archived REQ — the reason states which stamp or hash is wrong and what to write instead",
		})
	}
}

// appendTimestampOrderingFindings reports a REQ whose stamps cannot describe a
// real sequence of events: created_at after claimed_at, or claimed_at after
// completed_at. detectCompletionAnomaly (model.go) already covers the single
// completed_at < claimed_at case for terminal tickets, but only there — nothing
// checked created_at at all, so a claimed_at fabricated to any instant before
// completion passed every check the suite had. That is the class this probe
// closes, and it deliberately covers queue, working AND archive: the archive is
// where nothing auto-repairs, because the SessionStart repairer never touches it.
//
// The predicate is NOT a fourth spelling of the ordering rule. It is the same
// created_at <= claimed_at <= completed_at that the registered Go SessionStart
// repair owner enforces. Two boundary decisions are shared by both read and repair sides:
// the comparison is strict (equal stamps are legal — Step 2's claim and Step 3.6's
// estimate can read the same instant), and an absent or unparseable stamp is other
// checks' territory rather than a violation here, matching detectCompletionAnomaly's
// carve-out. If either spelling changes, change both.
//
// Each violated pair is its own finding: a REQ can be wrong in both pairs, and
// collapsing them would hide half the repair. Not Fixable — `do-work cleanup` does
// not rewrite stamps; the remedy names the tool that does, which differs by where
// the file lives.
func appendTimestampOrderingFindings(report *VerifyReport, board *Board) {
	for _, ticket := range board.AllRequests {
		createdInstant, createdParsed := parseTimestamp(ticket.CreatedAt)
		claimedInstant, claimedParsed := parseTimestamp(ticket.ClaimedAt)
		completedInstant, completedParsed := parseTimestamp(ticket.CompletedAt)

		if createdParsed && claimedParsed && claimedInstant.Before(createdInstant) {
			report.Findings = append(report.Findings, timestampOrderingFinding(
				ticket, "created_at", ticket.CreatedAt, "claimed_at", ticket.ClaimedAt,
				"claimed before it was created"))
		}
		if claimedParsed && completedParsed && completedInstant.Before(claimedInstant) {
			report.Findings = append(report.Findings, timestampOrderingFinding(
				ticket, "claimed_at", ticket.ClaimedAt, "completed_at", ticket.CompletedAt,
				"completed before it was claimed"))
		}
		// The outer pair, checked ONLY when claimed_at cannot carry the comparison.
		// created_at <= completed_at is implied transitively whenever claimed_at
		// parses, so checking it there would report a third finding for the same
		// defect; with claimed_at absent or unparseable the implication has nothing
		// to travel through, and an impossible ordering would otherwise pass. That
		// gap is the same shape as the one this whole probe closes — a real
		// violation invisible because the checked pairs did not span it.
		if !claimedParsed && createdParsed && completedParsed &&
			completedInstant.Before(createdInstant) {
			report.Findings = append(report.Findings, timestampOrderingFinding(
				ticket, "created_at", ticket.CreatedAt, "completed_at", ticket.CompletedAt,
				"completed before it was created, and no parseable claimed_at sits between them"))
		}
	}
}

// timestampOrderingFinding names both fields and both raw values, and routes the
// remedy by where the file actually lives — a reader must never be sent to do by
// hand what an installed script already does. The archive repair is deliberately a
// conscious invocation and is never hook-wired, which is why the two halves of this
// remedy differ.
func timestampOrderingFinding(
	ticket *RequestTicket, earlierField string, earlierValue string,
	laterField string, laterValue string, plainSummary string,
) VerifyFinding {
	remedy := "queue/ and working/ stamps are repaired mechanically by the registered Go SessionStart hook owner on the next session — no hand edit needed; run `skills/do-work/tools/do-work-cli.sh --format text repair-req-timestamps` for immediate recovery"
	if ticket.TreeSection == "archive" {
		remedy = "run `skills/do-work/tools/do-work-cli.sh --format text audit-archive-timestamps` to report, and add --fix to " +
			"repair from the stamp's own git history — the SessionStart repairer deliberately never " +
			"touches the archive, so this one is a conscious invocation"
	}
	return VerifyFinding{
		Category: verifyCategoryTimestampOrdering,
		Subject:  ticket.RequestId,
		Detail: fmt.Sprintf("%s (status %s): %s %q is later than %s %q — %s",
			ticket.RequestId, ticket.Status, earlierField, earlierValue, laterField, laterValue, plainSummary),
		Fixable: false,
		Remedy:  remedy,
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

	for _, mentionedId := range checkpointMentionedRequestIds(string(checkpointBytes)) {
		if _, exists := board.RequestsById[mentionedId]; exists {
			continue
		}
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryCheckpointGhostRequest,
			Subject:  mentionedId,
			Detail:   fmt.Sprintf("do-work/CHECKPOINT.md names %s, which does not exist in queue/, working/, or archive/", mentionedId),
			Remedy:   "edit that id's own entry out of the checkpoint rather than deleting the file — the file may also hold entries under another checkout's writer: label, which are live claims there; a later session reads it, and crash recovery classifies against it",
		})
	}
}

// checkpointMentionedRequestIds returns the distinct REQ ids a checkpoint's
// free text names, in first-mention order. A match whose digit run continues
// straight into `[` is skipped: that is a quoted shell glob
// (`REQ-0[0-9][0-9]-*.md` in session notes), not a REQ id — its `REQ-0` prefix
// used to be reported as a ghost. Go's RE2 has no lookahead, so the boundary
// is checked here rather than in the pattern; only `[` is treated as a glob
// continuation, because `*` and `?` legitimately follow real ids in prose
// (`**REQ-093**` emphasis, a sentence ending "…REQ-093?").
func checkpointMentionedRequestIds(checkpointText string) []string {
	var mentionedIds []string
	seenIds := map[string]bool{}
	for _, matchSpan := range requestIdMentionPattern.FindAllStringIndex(checkpointText, -1) {
		if matchSpan[1] < len(checkpointText) && checkpointText[matchSpan[1]] == '[' {
			continue
		}
		mentionedId := checkpointText[matchSpan[0]:matchSpan[1]]
		if seenIds[mentionedId] {
			continue
		}
		seenIds[mentionedId] = true
		mentionedIds = append(mentionedIds, mentionedId)
	}
	return mentionedIds
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
				Subject:  claimedTicket.RequestId,
				Detail:   fmt.Sprintf("%s is claimed but carries no claimed_at — its age cannot be known", claimedTicket.RequestId),
				Remedy:   "run `do-work run`; actions/work-reference.md -> Crash Recovery (Step 1) owns the human reset decision",
			})
			continue
		}
		claimInstant, parsedOk := parseTimestamp(rawClaimStamp)
		if !parsedOk {
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryClaimNeedsAttention,
				Subject:  claimedTicket.RequestId,
				Detail:   fmt.Sprintf("%s has an unparseable claimed_at (%q)", claimedTicket.RequestId, rawClaimStamp),
				Remedy:   "fix the stamp to a UTC ISO-8601 instant, or run `do-work run` and use actions/work-reference.md -> Crash Recovery (Step 1) for the reset decision",
			})
			continue
		}
		if claimInstant.After(skewHorizon) {
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryClaimNeedsAttention,
				Subject:  claimedTicket.RequestId,
				Detail: fmt.Sprintf("%s has a future-dated claimed_at (%s) — usually %s",
					claimedTicket.RequestId, rawClaimStamp, futureStampCauseClause),
				Remedy: "re-stamp it with the current UTC instant — `skills/do-work/tools/do-work-cli.sh --format text now` prints exactly that shape on any platform (the Timestamp rule in actions/work-reference.md)",
			})
			continue
		}
		if claimAge := now.Sub(claimInstant); claimAge >= staleClaimThreshold {
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryClaimNeedsAttention,
				Subject:  claimedTicket.RequestId,
				Detail: fmt.Sprintf("%s has been claimed for %s (threshold %s) — reported, not judged dead",
					claimedTicket.RequestId, formatApproximateDuration(claimAge), formatApproximateDuration(staleClaimThreshold)),
				Remedy: "a long build can legitimately exceed this; run `do-work run` and use actions/work-reference.md -> Crash Recovery (Step 1) to decide whether to take it over",
			})
		}
	}
}

// appendStrandedFinishedFindings reports a REQ with a terminal status still
// sitting in queue/ or working/ — doctor's STRANDED-TERMINAL-REQUEST definition, and
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
			Subject:  ticket.RequestId,
			Detail:   fmt.Sprintf("%s has terminal status %q but still sits in do-work/%s/", ticket.RequestId, ticket.Status, treeSection),
			Fixable:  true,
			Remedy:   "cleanup Pass 0 moves it into do-work/archive/",
		})
	}
}

// appendAssignedElsewhereFindings flags a REQ that reached do-work/working/ while
// still carrying assigned_to. Clearing that field is part of Step 2's claim
// (actions/work.md), so its survival past the move means either the claim skipped
// the clear or a session claimed work earmarked for another one. The marker is now
// actively wrong: it tells every other checkout to skip a REQ this one is building.
//
// A terminally-resolved REQ stranded in working/ is deliberately NOT reported here:
// appendStrandedFinishedFindings already owns that state, it is mechanically
// fixable, and this probe's remedy would tell the user to clear or release a claim
// on work that is already done.
//
// Read-only and not fixable — whose claim wins is a human call, and cleanup asks.
func appendAssignedElsewhereFindings(report *VerifyReport, board *Board) {
	for _, ticket := range board.AllRequests {
		if ticket.AssignedTo == "" || ticket.TreeSection != "working" {
			continue
		}
		// Stranded, not being built — see the carve-out in this function's doc comment.
		if isTerminalResolvedStatus(ticket.Status) {
			continue
		}
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryAssignedElsewhereClaimedHere,
			Subject:  ticket.RequestId,
			Detail: fmt.Sprintf("%s sits in do-work/working/ but is still assigned to %q — the claim did not clear the marker",
				ticket.RequestId, ticket.AssignedTo),
			Remedy: "cleanup asks before touching it: clear assigned_to if this checkout is the one building it, or release the claim if the earmark should stand",
		})
	}
}

// appendArchivedUserRequestLiveMemberFindings flags an archived UR that still has a
// member REQ sitting in do-work/queue/ or do-work/working/. A UR reaches
// do-work/archive/ only once every member is terminally resolved (actions/work.md
// Step 8), so this state means the closure check ran on stale information or a
// folder was moved by hand — and the live REQ is now orphaned from the input.md
// that records why it exists.
//
// Membership is the UR's RequestIds, which linkRequestsToUserRequests fills from
// each REQ's `user_request:` frontmatter — deliberately NOT the UR's own
// `requests:` array. That array is written once at capture time and can be wrong in
// both directions (a follow-up REQ captured later is absent from it, a split or
// abandoned REQ lingers in it), which is exactly why Step 8's closure predicate is
// a frontmatter scan. A probe reading the array would go silent on the follow-up
// case, which is the common one.
//
// A terminally-resolved member stranded in queue/ or working/ is deliberately NOT
// reported here: appendStrandedFinishedFindings already owns that state, it is
// mechanically fixable, and this probe's remedy would tell the user to run or
// abandon a REQ that is already resolved.
//
// After that ownership carve-out, an exact review_generated: true marker exempts
// a non-terminal queue/working member: REQ-193 keeps that follow-up on the same
// already-closed UR without reopening it. Ordinary siblings remain anomalies.
//
// Read-only and not fixable: the archived UR stays closed while a human resolves
// or abandons the ordinary live REQ, or corrects its user_request association.
func appendArchivedUserRequestLiveMemberFindings(report *VerifyReport, board *Board) {
	for _, userRequestTicket := range board.UserRequests {
		if !isArchivedUserRequestPath(userRequestTicket.FilePath) {
			continue
		}
		var liveMemberIds []string
		for _, memberRequestId := range userRequestTicket.RequestIds {
			memberTicket, memberFound := board.RequestsById[memberRequestId]
			if !memberFound {
				continue
			}
			// Stranded, not live — see the carve-out in this function's doc comment.
			if isTerminalResolvedStatus(memberTicket.Status) {
				continue
			}
			if memberTicket.TreeSection != "queue" && memberTicket.TreeSection != "working" {
				continue
			}
			if memberTicket.ReviewGenerated {
				continue
			}
			liveMemberIds = append(liveMemberIds, memberRequestId)
		}
		if len(liveMemberIds) == 0 {
			continue
		}
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryArchivedUserRequestLiveMember,
			Subject:  userRequestTicket.UserRequestId,
			Detail: fmt.Sprintf("%s is archived but still has live member(s) %s in do-work/queue/ or do-work/working/",
				userRequestTicket.UserRequestId, strings.Join(liveMemberIds, ", ")),
			Remedy: "keep the UR archived; resolve or abandon each ordinary live member, or correct its user_request association if it does not belong to this UR",
		})
	}
}

// isArchivedUserRequestPath reports whether a UR's input.md lives under
// do-work/archive/. The UR ticket carries no TreeSection of its own (unlike a REQ),
// so the path is the only evidence.
//
// The match must NOT require a leading separator. resolveRepoRootOrDefault returns
// an explicit --repo-root override verbatim, so `verify --repo-root .` yields ticket
// paths like "do-work/archive/UR-090/input.md" with no leading slash at all — and a
// leading-slash pattern silently recognized zero archived URs in that supported CLI
// mode, which is exactly the mode actions/forensics.md Check 14 can invoke. Anchor
// on the trailing separator instead: "do-work/archive/" still cannot match a
// directory merely named "archive", and it matches under both absolute and relative
// roots. Guard the prefix case explicitly rather than trusting Contains, so a path
// like "my-do-work/archive/" cannot satisfy it either.
func isArchivedUserRequestPath(userRequestFilePath string) bool {
	normalizedPath := filepath.ToSlash(userRequestFilePath)
	const archiveSegment = "do-work/archive/"
	return strings.HasPrefix(normalizedPath, archiveSegment) ||
		strings.Contains(normalizedPath, "/"+archiveSegment)
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

// worktreeLeftoverDisposition is the verdict for one worktree-agent-* name once
// every probe has run. A merged branch tip is NECESSARY for mechanical removal
// and never sufficient: `git merge-base --is-ancestor` proves the COMMITS are
// contained in the integration branch, and says nothing about the worktree
// holding them or about the run that created it (REQ-458). Two repository facts
// finish the question, and neither is a liveness signal.
type worktreeLeftoverDisposition int

const (
	// Merged, clean, and its REQ has left do-work/working/ — the one state
	// cleanup Pass 5 resolves mechanically.
	worktreeLeftoverFinishedResidue worktreeLeftoverDisposition = iota
	// Merged, but the worktree holds uncommitted work. Pass 5's `git worktree
	// remove` runs without --force and refuses exactly this, so calling it
	// fixable would contradict the command the remedy names.
	worktreeLeftoverUncommittedWork
	// Merged and clean, but the REQ the name carries is still in
	// do-work/working/ — the run owns the worktree until its REQ reaches its
	// final path, which is after review and any remediation.
	worktreeLeftoverRunInFlight
	// A probe could not answer, so neither cleanliness nor finishedness is
	// established. Fail safe: unestablished is never advertised as fixable.
	worktreeLeftoverStateUnknown
	worktreeLeftoverUnmergedBranch
	worktreeLeftoverMergeStateUndetermined
)

// worktreeAgentRequestIdPattern reads a leftover's own REQ id out of the
// worktree-agent-REQ-NNN-<suffix> convention (actions/work-reference.md →
// Worktree Dispatch Mode, "Naming"). Anchored at the prefix on purpose: a suffix
// derived from the REQ's title slug can mention another REQ id, and that mention
// is not this worktree's owner.
var worktreeAgentRequestIdPattern = regexp.MustCompile(`^` + worktreeAgentNamePrefix + `(REQ-\d+)`)

// requestIdFromWorktreeName returns the REQ id the name carries, or "" when the
// name does not follow the convention closely enough to name one.
func requestIdFromWorktreeName(leftoverName string) string {
	match := worktreeAgentRequestIdPattern.FindStringSubmatch(leftoverName)
	if match == nil {
		return ""
	}
	return match[1]
}

// requestPipelineState is what the board can honestly say about the REQ a
// worktree name carries. Three answers rather than two, because a boolean folds
// "the board has never seen this id" into "this REQ has moved on" and only the
// second of those is something verify read (REQ-458).
//
// Absence has causes verify cannot tell apart — a REQ file never created,
// deleted, renamed, or parked somewhere the walk reports as a stray (see
// appendStrayRequestFileFindings) — and none of them is evidence that the run
// which created the worktree finished.
type requestPipelineState int

const (
	// The REQ sits in do-work/working/, where the pipeline keeps a request from
	// its claim until it reaches its final path.
	requestPipelineStateInFlight requestPipelineState = iota
	// The REQ is on the board and has left do-work/working/ — the only state in
	// which "the run that created this worktree has finished" is a read rather
	// than an assumption.
	requestPipelineStateSettled
	// The board's request index carries no REQ with this id at all.
	requestPipelineStateAbsentFromBoard
)

// classifyRequestPipelineState reads the board's already-parsed request index
// rather than walking the tree a second time — the same reuse
// appendStrandedFinishedFindings makes of TreeSection, and one reader means one
// answer. This is a question about where a file is, not about whether a builder
// process is running: REQ-073 forbids the latter, and nothing here approaches it.
func classifyRequestPipelineState(board *Board, requestId string) requestPipelineState {
	ticket, exists := board.RequestsById[requestId]
	switch {
	case !exists:
		return requestPipelineStateAbsentFromBoard
	case ticket.TreeSection == "working":
		return requestPipelineStateInFlight
	default:
		return requestPipelineStateSettled
	}
}

// noOptionalLocksFlag keeps `git status` from refreshing the index it read.
// Without it, a status run rewrites .git/worktrees/<name>/index whenever a stat
// looks stale, so verify would leave a changed file behind on every run — which
// contradicts runVerifyProbes' contract that it writes nothing anywhere, and the
// board prime's "`frontmatter` and `verify` write nothing at all".
//
// It is a top-level git option, so it must come BEFORE -C and before the
// subcommand; git rejects it after the subcommand name.
const noOptionalLocksFlag = "--no-optional-locks"

// worktreeHasUncommittedWork reports whether a builder's worktree holds any
// uncommitted change at all — modified tracked files, staged edits, or untracked
// files.
//
// Deliberately the WHOLE worktree, not the do-work/ subset worktreeDirtyQueueState
// asks about. The two answer different questions: that one asks whether a builder
// broke the "state stays home" rule, this one asks whether Pass 5's non-forced
// `git worktree remove` would refuse, and it refuses for dirt anywhere.
//
// A non-nil error means the question went unanswered. Callers must not read that
// as clean — VerifyReport's contract is that silence reads as "checked and clean,"
// which is the one thing an unanswered probe must never say.
func worktreeHasUncommittedWork(worktreePath string) (bool, error) {
	command := exec.Command("git", noOptionalLocksFlag, "-C", worktreePath, "status", "--porcelain", "--untracked-files=all")
	output, runError := command.Output()
	if runError != nil {
		return false, fmt.Errorf("`git status` in %s failed, so whether it holds uncommitted work is unknown", worktreePath)
	}
	return strings.TrimSpace(string(output)) != "", nil
}

// classifyWorktreeLeftover decides one leftover's disposition from repository
// reads only: the merge probe, `git status` inside the worktree, and where the
// REQ named by the worktree sits in the do-work tree.
//
// None of the three asks whether a builder is alive, and none may be replaced by
// something that does — REQ-073 rules out a lock, heartbeat, PID check, mtime
// heuristic and time threshold alike (see staleClaimThreshold's own comment). The
// two added facts hold regardless of any process: a dirty worktree is dirty
// whether or not anyone still holds it, and a REQ in working/ is unfinished
// whether or not its builder still runs.
//
// requestId is the REQ id the caller already read out of leftoverName, passed in
// rather than re-derived so one name yields one id.
//
// The returned error is non-nil exactly when a fact could not be established; the
// disposition is then worktreeLeftoverStateUnknown, never a silent "clean". A REQ
// id the board does not carry is one of those unestablished facts, not a quiet
// "finished": absence proves nothing about a run.
//
// Precedence among the merged sub-states is deliberate. An in-flight run is
// reported first because it is the fact that decides what to do — a builder's
// worktree being dirty during its own run is expected rather than notable, while
// "this belongs to a run that has not finished" is the reason to leave it alone
// either way. Cleanliness is probed last, and only once the board has placed the
// REQ outside working/: a worktree whose finishedness verify could not establish
// gets inspected by hand whether or not it is also dirty, so the earlier question
// is the one that decides.
func classifyWorktreeLeftover(repoRoot string, board *Board, leftoverName string, requestId string, worktreePath string, hasWorktree bool) (worktreeLeftoverDisposition, error) {
	switch classifyWorktreeMergeState(repoRoot, leftoverName) {
	case worktreeMergeStateUnmerged:
		return worktreeLeftoverUnmergedBranch, nil
	case worktreeMergeStateUndetermined:
		return worktreeLeftoverMergeStateUndetermined, nil
	}

	if requestId == "" {
		return worktreeLeftoverStateUnknown, fmt.Errorf(
			"%s carries no REQ-NNN id, so whether the run that created it has finished cannot be read from do-work/working/", leftoverName)
	}
	switch classifyRequestPipelineState(board, requestId) {
	case requestPipelineStateInFlight:
		return worktreeLeftoverRunInFlight, nil
	case requestPipelineStateAbsentFromBoard:
		return worktreeLeftoverStateUnknown, fmt.Errorf(
			"%s names %s, which the board carries no REQ file for in do-work/queue/, working/ or archive/, so whether the run that created it has finished cannot be read from do-work/working/", leftoverName, requestId)
	}

	// A branch whose worktree is already gone has no working tree to be dirty:
	// there is nothing left to hold uncommitted work, and nothing to remove but
	// the branch. Answered explicitly rather than probed, because probing a path
	// that does not exist would fail and be reported as an unanswered question.
	if !hasWorktree {
		return worktreeLeftoverFinishedResidue, nil
	}
	hasUncommittedWork, statusError := worktreeHasUncommittedWork(worktreePath)
	if statusError != nil {
		return worktreeLeftoverStateUnknown, statusError
	}
	if hasUncommittedWork {
		return worktreeLeftoverUncommittedWork, nil
	}
	return worktreeLeftoverFinishedResidue, nil
}

// routeWorktreeLeftover maps a disposition onto the finding it produces.
// requestId is the leftover's own REQ id, used to name the evidence in the
// in-flight remedy; it may be empty for the states that do not cite it.
//
// Fixable is true for finished residue and nothing else, because that is the only
// state actions/cleanup.md → Pass 5 resolves mechanically: `git worktree remove`
// plus `git branch -d`, neither forcing. Every other state lands on Pass 5's
// consent-gated path, where the pass "stops being mechanical" and asks — a human
// decision, which VerifyFinding.Fixable's doc comment says must not be advertised
// otherwise. The two present-and-non-fixable states earn separate categories and
// separate remedies because their remedies genuinely differ: one asks the reader
// to resolve uncommitted edits, the other asks them to leave the worktree alone.
//
// An unmerged leftover stays a reported finding during a live run rather than
// being suppressed while builders are in flight. VerifyReport's doc comment is
// explicit that silence reads as "checked and clean," and for an unmerged branch
// verify has no way to know a run is active (see classifyWorktreeMergeState) — so
// suppression would have to guess, and would hide genuinely stranded work whenever
// it guessed wrong. This mirrors how version-changelog-mismatch handles its own
// expected mid-release state: reported, with the transient case named in the
// remedy, and not fixable.
func routeWorktreeLeftover(disposition worktreeLeftoverDisposition, requestId string) (category string, fixable bool, remedy string) {
	switch disposition {
	case worktreeLeftoverFinishedResidue:
		return verifyCategoryMergedWorktreeLeftover, true,
			"cleanup Pass 5 removes it mechanically — the branch is already contained in HEAD, the worktree is clean, and its REQ has left do-work/working/, so nothing is lost"
	case worktreeLeftoverUncommittedWork:
		return verifyCategoryWorktreePresentUncommittedWork, false,
			"this worktree is present, not leftover: its branch is merged, but the worktree holds uncommitted changes that are in no commit. cleanup Pass 5 removes nothing here — its `git worktree remove` runs without --force and refuses a dirty worktree. Inspect the changes and commit or discard them in the worktree before anything is removed"
	case worktreeLeftoverRunInFlight:
		return verifyCategoryWorktreePresentRunInFlight, false,
			"this worktree is present, not leftover: " + requestId + " is still in do-work/working/, so the run that owns it has not reached review, remediation, or its final path. Leave it in place — worktree cleanup happens after the REQ is archived, not before"
	case worktreeLeftoverStateUnknown:
		return verifyCategoryWorktreePresentStateUnknown, false,
			"this worktree is present and verify could not establish that it is both clean and finished — the read that failed is named in this report's own skipped-probes list, which the board shows as a `not checked` row in the same list as this finding. Inspect it by hand; an unestablished state is never advertised as mechanically removable"
	case worktreeLeftoverUnmergedBranch:
		return verifyCategoryUnmergedWorktreeLeftover, false,
			"this is either a builder still in flight or work that outlived a dead run — verify cannot tell those apart. Leave it alone during a run; otherwise cleanup Pass 5 asks before discarding it, because the branch may hold the only copy"
	default:
		return verifyCategoryUndeterminedWorktreeLeftover, false,
			"git could not say whether this is merged (typically a worktree whose branch is gone) — inspect it by hand; cleanup Pass 5 deletes nothing it cannot establish a merge target for. The same unresolved branch stopped the committed-queue-state check, so whether a builder committed queue state under do-work/ on it is unknown, not clean"
	}
}

// appendWorktreeFindings covers the two worktree-dispatch invariants: no
// `worktree-agent-*` leftovers should outlive their run, and a builder must never
// write queue state (actions/work-reference.md → Worktree Dispatch Mode, "state
// stays home" and "sole integrator").
//
// Each name is classified rather than reported as one kind of thing, so the report
// routes only the mechanically-resolvable ones to cleanup. The board answers where
// a leftover's REQ sits — in do-work/working/, past it, or nowhere the index knows
// — and it is the board this run already built, never a second walk of the tree.
func appendWorktreeFindings(report *VerifyReport, repoRoot string, board *Board) {
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

	// Resolved once: every builder branch is compared against the same integration
	// point. An unresolvable one disables only the committed-state half, which is
	// reported rather than passed over — silence would read as "checked and clean."
	integrationRef, integrationRefError := resolveIntegrationBranchRef(repoRoot)
	if integrationRefError != nil {
		report.SkippedProbes = append(report.SkippedProbes,
			fmt.Sprintf("committed-queue-state probe: %v", integrationRefError))
	}

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
		requestId := requestIdFromWorktreeName(leftoverName)
		disposition, probeError := classifyWorktreeLeftover(repoRoot, board, leftoverName, requestId, worktreePath, hasWorktree)
		if probeError != nil {
			report.SkippedProbes = append(report.SkippedProbes,
				fmt.Sprintf("worktree removability probe for %s: %v", leftoverName, probeError))
		}
		category, fixable, remedy := routeWorktreeLeftover(disposition, requestId)
		report.Findings = append(report.Findings, VerifyFinding{
			Category: category,
			Subject:  leftoverName,
			Detail:   fmt.Sprintf("%s%s exists — %s", worktreeAgentNamePrefix, strings.TrimPrefix(leftoverName, worktreeAgentNamePrefix), locationDetail),
			Fixable:  fixable,
			Remedy:   remedy,
		})

		// Checked before the worktree guard below: a builder's committed queue edits
		// live in the branch, so they are detectable — and just as wrong — after the
		// worktree itself is gone.
		//
		// Not checked at all when the merge state came back undetermined: the branch
		// git could not resolve for `merge-base` is the same one `git diff
		// <ref>...<name>` needs, so the probe would fail for the reason already
		// reported and print a second row for one fact. The undetermined remedy
		// carries "this went unchecked" instead — moved, not dropped, because
		// silence in this report reads as checked-and-clean.
		if integrationRefError == nil && disposition != worktreeLeftoverMergeStateUndetermined {
			committedQueuePaths, committedError := worktreeCommittedQueueState(repoRoot, integrationRef, leftoverName)
			switch {
			case committedError != nil:
				report.SkippedProbes = append(report.SkippedProbes,
					fmt.Sprintf("committed-queue-state probe for %s: %v", leftoverName, committedError))
			case len(committedQueuePaths) > 0:
				report.Findings = append(report.Findings, VerifyFinding{
					Category: verifyCategoryWorktreeCommittedQueueState,
					Subject:  leftoverName,
					Detail: fmt.Sprintf("%s has committed changes under do-work/ on its branch (%s) — a builder wrote queue state the orchestrator alone owns",
						leftoverName, strings.Join(committedQueuePaths, ", ")),
					Remedy: "those commits are in the branch about to be merged — drop or revert them there before integrating; every claim, status flip, and archive move belongs to the main tree",
				})
			}
		}

		if !hasWorktree {
			continue
		}
		if dirtyQueuePaths := worktreeDirtyQueueState(worktreePath); len(dirtyQueuePaths) > 0 {
			report.Findings = append(report.Findings, VerifyFinding{
				Category: verifyCategoryWorktreeWroteQueueState,
				Subject:  leftoverName,
				Detail: fmt.Sprintf("%s has uncommitted changes under do-work/ (%s) — a builder wrote queue state the orchestrator alone owns",
					worktreePath, strings.Join(dirtyQueuePaths, ", ")),
				Remedy: "discard those changes in the worktree; every claim, status flip, and archive move belongs to the main tree",
			})
		}
	}
}

// resolveIntegrationBranchRef names the point a builder branch is compared
// against: the repo-root checkout's own branch.
//
// Named explicitly rather than passed as `HEAD`, because `HEAD` means whatever
// checkout the command runs in — inside a builder's worktree it names the
// builder's own branch, and the comparison would silently become branch-against-
// itself. That is the same class of error as `git branch -d` testing merged-ness
// against whatever HEAD happens to be (actions/work-reference.md → Worktree
// Dispatch Mode, "Cleanup — happy path"). A detached repo-root checkout has no
// branch name, so the commit id is returned instead — it names the same point
// just as explicitly.
func resolveIntegrationBranchRef(repoRoot string) (string, error) {
	branchCommand := exec.Command("git", "-C", repoRoot, "rev-parse", "--abbrev-ref", "HEAD")
	branchOutput, branchError := branchCommand.Output()
	if branchError != nil {
		return "", fmt.Errorf("cannot resolve the integration branch at %s", repoRoot)
	}
	integrationRef := strings.TrimSpace(string(branchOutput))
	if integrationRef != "" && integrationRef != "HEAD" {
		return integrationRef, nil
	}
	commitCommand := exec.Command("git", "-C", repoRoot, "rev-parse", "HEAD")
	commitOutput, commitError := commitCommand.Output()
	if commitError != nil {
		return "", fmt.Errorf("cannot resolve the integration commit at %s (detached checkout with no commits?)", repoRoot)
	}
	return strings.TrimSpace(string(commitOutput)), nil
}

// worktreeCommittedQueueState returns the do-work/ paths a builder's BRANCH has
// added or changed relative to the integration branch — the owner impersonation
// that survives a commit, and the shape REQ-072 requirement 3 actually asked for.
//
// The three dots are load-bearing. `A...B` diffs from merge-base(A,B) to B, so it
// reports what the builder's branch changed and stays blind to how far the
// integration branch has moved since the branch point. That distinction is what
// worktreeDirtyQueueState's doc comment (below) is protecting: where a consumer
// commits do-work/, the worktree legitimately carries a stale snapshot while the
// main tree moves constantly as the orchestrator claims and archives. A two-tree
// comparison would fire on nearly every run; this one fires only on the builder's
// own writes.
//
// Where do-work/ is untracked — the common install — there is nothing committed
// to diff and this simply returns nothing, leaving the porcelain check as the
// only probe that can see anything. That is correct, not a gap.
func worktreeCommittedQueueState(repoRoot string, integrationRef string, branchName string) ([]string, error) {
	command := exec.Command("git", "-C", repoRoot, "diff", "--name-only",
		integrationRef+"..."+branchName, "--", "do-work/")
	output, runError := command.Output()
	if runError != nil {
		return nil, fmt.Errorf("`git diff %s...%s -- do-work/` failed (no such branch, or unrelated histories)", integrationRef, branchName)
	}
	var committedPaths []string
	for _, line := range strings.Split(string(output), "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			committedPaths = append(committedPaths, trimmedLine)
		}
	}
	return committedPaths, nil
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
	command := exec.Command("git", noOptionalLocksFlag, "-C", worktreePath, "status", "--porcelain", "--untracked-files=all", "--", "do-work/")
	output, runError := command.Output()
	if runError != nil {
		return nil
	}
	return parsePorcelainStatusPaths(string(output))
}

// parsePorcelainStatusPaths extracts the path from each `git status --porcelain`
// v1 line: a fixed two-character status, one space, then the path — which may
// itself contain spaces, so the path is everything after the 3-byte prefix, not
// the last whitespace-separated field. Rename lines carry `old -> new`; the
// destination side is kept. Git double-quotes paths with special characters;
// the quotes are stripped for display but inner escapes are left as git printed
// them (these paths feed a report line, not file operations).
func parsePorcelainStatusPaths(statusOutput string) []string {
	var dirtyPaths []string
	for _, line := range strings.Split(statusOutput, "\n") {
		if len(line) < 4 {
			continue
		}
		pathText := line[3:]
		if renameArrowIndex := strings.LastIndex(pathText, " -> "); renameArrowIndex >= 0 {
			pathText = pathText[renameArrowIndex+len(" -> "):]
		}
		pathText = strings.Trim(pathText, `"`)
		if pathText != "" {
			dirtyPaths = append(dirtyPaths, pathText)
		}
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
	// Rendered with its own marker and word so a reader can tell the two apart at a glance:
	// a skip is work left undone, a not-applicable is a question this repo was never asked.
	for _, notApplicableProbe := range report.NotApplicableProbes {
		fmt.Fprintf(&builder, "  ~ not applicable: %s\n", notApplicableProbe)
	}
	if fixableCount := report.FixableCount(); fixableCount > 0 {
		fmt.Fprintf(&builder, "  %d fixable: run do-work cleanup\n", fixableCount)
	}
	return builder.String()
}

// calibrationToleranceMinutes absorbs truncation-versus-rounding between the
// writer and this reader, which is noise. Anything past it is a real divergence
// between two records that were supposed to describe the same span.
const calibrationToleranceMinutes = 1

// appendCalibrationLogFindings reconciles do-work/calibration-log.tsv against the
// frontmatter each row was derived from. The log is an independent third record —
// written once by actions/work.md Step 8 substep 7.5 as completed_at − claimed_at,
// never revised, and read back by actions/estimate-reference.md as the corpus the
// scoring table is fit from. Nothing compared the two, so a corpus that feeds every
// future estimate was decaying unaudited: measured on this repo, 10 of 72 rows
// disagree and eight of those materially.
//
// It reports and never repairs, and it deliberately does NOT pick a winner. Either
// record can legitimately be the wrong one: the log line is written once, while the
// frontmatter can be rewritten afterwards by the registered Go SessionStart repair
// owner, by the canonical audit command's --fix mode, or by a crash-recovery
// pass that cleared and re-stamped a claim. Resolving the disagreement needs a human
// who knows which happened, which is also why nothing here is Fixable.
//
// A row this probe cannot reconcile at all — malformed, naming a REQ that exists
// nowhere, or naming one whose stamps are absent or unparseable — is its own finding
// and never a disagreement. Reporting it as a mismatch would put a number next to a
// value that was never computed.
func appendCalibrationLogFindings(report *VerifyReport, repoRoot string, board *Board) {
	logPath := filepath.Join(repoRoot, "do-work", "calibration-log.tsv")
	logBytes, readError := os.ReadFile(logPath)
	if readError != nil {
		// A repo that has archived nothing yet has no log. That is not a defect, but
		// it is also not a verified invariant, so it is reported rather than silent.
		report.SkippedProbes = append(report.SkippedProbes, fmt.Sprintf(
			"calibration-log probe: do-work/calibration-log.tsv is unreadable or absent (%v) — no rows to reconcile", readError))
		return
	}

	for lineNumber, rawLine := range strings.Split(string(logBytes), "\n") {
		trimmedLine := strings.TrimSpace(rawLine)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "req_id\t") {
			continue // blank line, or the header
		}
		humanLineNumber := lineNumber + 1
		columns := strings.Split(trimmedLine, "\t")
		if len(columns) < 4 {
			report.Findings = append(report.Findings, calibrationRowFinding(humanLineNumber,
				fmt.Sprintf("has %d tab-separated column(s), want at least 4 (req_id, route, estimated_p50_minutes, wall_minutes)", len(columns))))
			continue
		}
		requestId := strings.TrimSpace(columns[0])
		loggedMinutes, parseError := strconv.Atoi(strings.TrimSpace(columns[3]))
		if parseError != nil {
			report.Findings = append(report.Findings, calibrationRowFinding(humanLineNumber,
				fmt.Sprintf("%s: wall_minutes %q is not an integer", requestId, columns[3])))
			continue
		}

		ticket, ticketFound := board.RequestsById[requestId]
		if !ticketFound {
			report.Findings = append(report.Findings, calibrationRowFinding(humanLineNumber,
				fmt.Sprintf("%s: the log has a row for it, but no such REQ exists in queue/, working/, or archive/", requestId)))
			continue
		}
		claimedInstant, claimedParsed := parseTimestamp(ticket.ClaimedAt)
		completedInstant, completedParsed := parseTimestamp(ticket.CompletedAt)
		if !claimedParsed || !completedParsed {
			report.Findings = append(report.Findings, calibrationRowFinding(humanLineNumber,
				fmt.Sprintf("%s: cannot recompute the span — claimed_at %q and completed_at %q do not both parse",
					requestId, ticket.ClaimedAt, ticket.CompletedAt)))
			continue
		}

		// Integer minutes, truncated, exactly as the write site computes it.
		recomputedMinutes := int(completedInstant.Sub(claimedInstant).Minutes())
		differenceMinutes := recomputedMinutes - loggedMinutes
		if differenceMinutes < 0 {
			differenceMinutes = -differenceMinutes
		}
		if differenceMinutes <= calibrationToleranceMinutes {
			continue
		}
		report.Findings = append(report.Findings, VerifyFinding{
			Category: verifyCategoryCalibrationLogMismatch,
			Subject:  requestId,
			Detail: fmt.Sprintf(
				"do-work/calibration-log.tsv line %d: %s logs wall_minutes %d, but its frontmatter recomputes to %d (claimed_at %q → completed_at %q)",
				humanLineNumber, requestId, loggedMinutes, recomputedMinutes, ticket.ClaimedAt, ticket.CompletedAt),
			Fixable: false,
			Remedy: "either record may be the correct one — the log row is written once and never revised, " +
				"while the frontmatter can have been rewritten since by the SessionStart repairer, by " +
				"audit-archive-timestamps.sh --fix, or by a crash-recovery re-stamp. Decide which describes " +
				"what actually happened before changing anything; the estimator fits its scoring table from these rows",
		})
	}
}

// calibrationRowFinding is the not-a-disagreement case: the row could not be
// reconciled at all, so no recomputed number exists to compare against. Kept as one
// category with the specific reason in the detail, because the remedy is the same
// question in every case — which record is wrong, the log or the tree.
func calibrationRowFinding(lineNumber int, reason string) VerifyFinding {
	return VerifyFinding{
		Category: verifyCategoryCalibrationRowUnreconcilable,
		Detail:   fmt.Sprintf("do-work/calibration-log.tsv line %d %s", lineNumber, reason),
		Fixable:  false,
		Remedy: "the row cannot be checked against any frontmatter — find the REQ it names, repair its stamps, " +
			"or correct the row. The estimator reads this file as its corpus, so an unreconcilable row is an " +
			"unverifiable input rather than a harmless one",
	}
}
