package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// verifyFixtureFile is one file to seed into a synthetic repo: a path relative to
// the repo root plus its literal content.
type verifyFixtureFile struct {
	RelativePath string
	Content      string
}

// writeVerifyFixture seeds a synthetic repo root. Every fixture gets do-work/queue/
// so the tree walk has something to walk.
func writeVerifyFixture(t *testing.T, files []verifyFixtureFile) string {
	t.Helper()
	repoRoot := t.TempDir()
	if mkdirError := os.MkdirAll(filepath.Join(repoRoot, "do-work", "queue"), 0o755); mkdirError != nil {
		t.Fatalf("mkdir do-work/queue: %v", mkdirError)
	}
	for _, fixtureFile := range files {
		absolutePath := filepath.Join(repoRoot, fixtureFile.RelativePath)
		if mkdirError := os.MkdirAll(filepath.Dir(absolutePath), 0o755); mkdirError != nil {
			t.Fatalf("mkdir for %s: %v", fixtureFile.RelativePath, mkdirError)
		}
		if writeError := os.WriteFile(absolutePath, []byte(fixtureFile.Content), 0o644); writeError != nil {
			t.Fatalf("write %s: %v", fixtureFile.RelativePath, writeError)
		}
	}
	return repoRoot
}

// findingsMentioning returns the findings whose Category matches, so assertions
// name the probe rather than counting anonymous strings.
func findingsMentioning(report VerifyReport, category string) []VerifyFinding {
	var matched []VerifyFinding
	for _, finding := range report.Findings {
		if finding.Category == category {
			matched = append(matched, finding)
		}
	}
	return matched
}

// A healthy repo's version file AGREES with the newest changelog entry — the
// steady state after a release lands. See verify.go's release probes on why
// "strictly greater" is a mid-release condition, not a checkable invariant.
const cleanVersionFile = "# Version Action\n\n**Current version**: 0.163.3\n"

const cleanChangelog = `# Changelog

---

## 0.163.3 — Board Copy Includes the REQ Id (2026-08-02)

Lead.

## 0.163.2 — Audit Sweep (2026-08-01)

Lead.
`

func TestVerifyPassesOnACleanTree(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/queue/REQ-071-only-one.md", "---\nid: REQ-071\nstatus: pending\ntitle: fixture\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	if len(report.Findings) != 0 {
		t.Errorf("clean tree produced %d findings, want 0:\n%s", len(report.Findings), renderVerifyReport(report))
	}
	if report.ExitCode() != 0 {
		t.Errorf("clean tree exit code = %d, want 0", report.ExitCode())
	}
}

// The version file and the newest changelog entry must agree. A mismatch is the
// finding, and each direction has its own cause — CLAUDE.md § Before Every Commit
// item 1 is the mid-release form of this rule ("strictly greater than the entry
// that already exists"), which becomes equality the moment the entry is written.
func TestVerifyFlagsVersionAheadOfNewestChangelogEntry(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", "**Current version**: 0.164.0\n"}, // bumped, entry not written yet
		{"CHANGELOG.md", cleanChangelog},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	versionFindings := findingsMentioning(report, verifyCategoryVersionChangelogMismatch)
	if len(versionFindings) != 1 {
		t.Fatalf("got %d version-mismatch findings, want 1:\n%s", len(versionFindings), renderVerifyReport(report))
	}
	if !strings.Contains(versionFindings[0].Detail, "0.164.0") || !strings.Contains(versionFindings[0].Detail, "0.163.3") {
		t.Errorf("finding must name both versions, got: %s", versionFindings[0].Detail)
	}
	if !strings.Contains(versionFindings[0].Detail, "ahead") {
		t.Errorf("finding must name the direction so the cause is legible, got: %s", versionFindings[0].Detail)
	}
	if versionFindings[0].Fixable {
		t.Error("a version mismatch is a human decision — it must not be advertised as cleanup-fixable")
	}
	if report.ExitCode() == 0 {
		t.Error("verify must exit non-zero when it has findings")
	}
}

func TestVerifyFlagsVersionBehindNewestChangelogEntry(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", "**Current version**: 0.163.2\n"}, // entry written, version file never bumped
		{"CHANGELOG.md", cleanChangelog},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	versionFindings := findingsMentioning(report, verifyCategoryVersionChangelogMismatch)
	if len(versionFindings) != 1 {
		t.Fatalf("got %d version-mismatch findings, want 1:\n%s", len(versionFindings), renderVerifyReport(report))
	}
	if !strings.Contains(versionFindings[0].Detail, "behind") {
		t.Errorf("finding must name the direction so the cause is legible, got: %s", versionFindings[0].Detail)
	}
}

// The "strictly greater" ordering that survives into the committed state is
// within the changelog itself — this is the duplicate-version-number failure
// CLAUDE.md records as having already happened.
func TestVerifyFlagsDuplicateChangelogVersion(t *testing.T) {
	duplicateVersionChangelog := `# Changelog

## 0.163.3 — A Newer Thing (2026-08-02)

Lead.

## 0.163.3 — The Same Number Reused (2026-08-01)

Lead.
`
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", "**Current version**: 0.163.3\n"},
		{"CHANGELOG.md", duplicateVersionChangelog},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	orderingFindings := findingsMentioning(report, verifyCategoryChangelogVersionNotAhead)
	if len(orderingFindings) != 1 {
		t.Fatalf("got %d changelog-ordering findings, want 1:\n%s", len(orderingFindings), renderVerifyReport(report))
	}
	if !strings.Contains(orderingFindings[0].Detail, "0.163.3") {
		t.Errorf("finding must name the duplicated version, got: %s", orderingFindings[0].Detail)
	}
}

// CLAUDE.md § Before Every Commit, item 2: the newest entry's title must not
// already be in use.
func TestVerifyFlagsReusedChangelogTitle(t *testing.T) {
	reusedTitleChangelog := `# Changelog

## 0.163.3 — Audit Sweep (2026-08-02)

Lead.

## 0.120.0 — Audit Sweep (2026-07-01)

The earlier entry that already used this title.
`
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", "**Current version**: 0.163.3\n"},
		{"CHANGELOG.md", reusedTitleChangelog},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	titleFindings := findingsMentioning(report, verifyCategoryReusedChangelogTitle)
	if len(titleFindings) != 1 {
		t.Fatalf("got %d reused-title findings, want 1:\n%s", len(titleFindings), renderVerifyReport(report))
	}
	if !strings.Contains(titleFindings[0].Detail, "0.120.0") {
		t.Errorf("finding must name the earlier entry that already used the title, got: %s", titleFindings[0].Detail)
	}
}

func TestVerifyFlagsDuplicateRequestIds(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/queue/REQ-071-first-slug.md", "---\nid: REQ-071\nstatus: pending\ntitle: first\n---\n"},
		{"do-work/queue/REQ-071-second-slug.md", "---\nid: REQ-071\nstatus: pending\ntitle: second\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	duplicateFindings := findingsMentioning(report, verifyCategoryDuplicateRequestId)
	if len(duplicateFindings) == 0 {
		t.Fatalf("two REQ-071 files under different slugs produced no duplicate finding:\n%s", renderVerifyReport(report))
	}
	if !strings.Contains(duplicateFindings[0].Detail, "REQ-071") {
		t.Errorf("finding must name the duplicated id, got: %s", duplicateFindings[0].Detail)
	}
}

func TestVerifyFlagsCheckpointNamingAMissingRequest(t *testing.T) {
	checkpointBody := `---
session_ended: 2026-08-03T09:00:00Z
last_completed: REQ-070
---

# Session Checkpoint

## In Progress (interrupted)
- REQ-999: a REQ that does not exist anywhere — stopped at implementing
`
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/CHECKPOINT.md", checkpointBody},
		{"do-work/archive/REQ-070-real.md", "---\nid: REQ-070\nstatus: completed\ntitle: real\ncompleted_at: 2026-08-02T10:00:00Z\ncommit: abc1234\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	checkpointFindings := findingsMentioning(report, verifyCategoryCheckpointGhostRequest)
	if len(checkpointFindings) != 1 {
		t.Fatalf("got %d checkpoint-ghost findings, want 1:\n%s", len(checkpointFindings), renderVerifyReport(report))
	}
	if !strings.Contains(checkpointFindings[0].Detail, "REQ-999") {
		t.Errorf("finding must name the ghost id, got: %s", checkpointFindings[0].Detail)
	}
	// REQ-070 exists, so it must not be flagged.
	for _, finding := range checkpointFindings {
		if strings.Contains(finding.Detail, "REQ-070") {
			t.Errorf("REQ-070 exists in the archive and must not be flagged: %s", finding.Detail)
		}
	}
}

func TestVerifyFlagsStaleAndMalformedClaims(t *testing.T) {
	fixedNow := time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)
	staleStamp := fixedNow.Add(-5 * time.Hour).Format(time.RFC3339)
	freshStamp := fixedNow.Add(-20 * time.Minute).Format(time.RFC3339)
	futureStamp := fixedNow.Add(3 * time.Hour).Format(time.RFC3339)

	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/working/REQ-101-stale-claim.md", "---\nid: REQ-101\nstatus: claimed\ntitle: stale\nclaimed_at: " + staleStamp + "\n---\n"},
		{"do-work/working/REQ-102-fresh-claim.md", "---\nid: REQ-102\nstatus: claimed\ntitle: fresh\nclaimed_at: " + freshStamp + "\n---\n"},
		{"do-work/working/REQ-103-future-claim.md", "---\nid: REQ-103\nstatus: claimed\ntitle: future\nclaimed_at: " + futureStamp + "\n---\n"},
		{"do-work/working/REQ-104-unparseable-claim.md", "---\nid: REQ-104\nstatus: claimed\ntitle: unparseable\nclaimed_at: not-a-timestamp\n---\n"},
		{"do-work/working/REQ-105-missing-claim-stamp.md", "---\nid: REQ-105\nstatus: claimed\ntitle: no stamp\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, fixedNow)
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	claimFindings := findingsMentioning(report, verifyCategoryClaimNeedsAttention)
	flaggedIds := map[string]bool{}
	for _, finding := range claimFindings {
		for _, requestId := range []string{"REQ-101", "REQ-102", "REQ-103", "REQ-104", "REQ-105"} {
			if strings.Contains(finding.Detail, requestId) {
				flaggedIds[requestId] = true
			}
		}
	}
	for _, shouldFlag := range []string{"REQ-101", "REQ-103", "REQ-104", "REQ-105"} {
		if !flaggedIds[shouldFlag] {
			t.Errorf("%s should be flagged (stale, future-dated, unparseable, or absent claimed_at):\n%s", shouldFlag, renderVerifyReport(report))
		}
	}
	if flaggedIds["REQ-102"] {
		t.Error("REQ-102 was claimed 20 minutes ago and is inside the threshold — it must not be flagged")
	}
}

func TestVerifyFlagsStrandedFinishedRequests(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/queue/REQ-060-finished-but-queued.md", "---\nid: REQ-060\nstatus: completed\ntitle: stranded\ncompleted_at: 2026-08-01T10:00:00Z\ncommit: abc1234\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	strandedFindings := findingsMentioning(report, verifyCategoryStrandedFinishedRequest)
	if len(strandedFindings) != 1 {
		t.Fatalf("got %d stranded-finished findings, want 1:\n%s", len(strandedFindings), renderVerifyReport(report))
	}
	if !strandedFindings[0].Fixable {
		t.Error("a stranded finished REQ is exactly what cleanup Pass 0 sweeps — it must be marked fixable")
	}
}

// Requirement 5: the report names how many findings cleanup can mechanically
// resolve and points at it. It must not claim a human-decision finding is fixable.
func TestVerifyReportRoutesFixableFindingsToCleanup(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", "**Current version**: 0.164.0\n"}, // mismatch — not fixable
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/queue/REQ-060-finished-but-queued.md", "---\nid: REQ-060\nstatus: completed\ntitle: stranded\ncompleted_at: 2026-08-01T10:00:00Z\ncommit: abc1234\n---\n"}, // fixable
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	if report.FixableCount() != 1 {
		t.Errorf("FixableCount = %d, want 1 (the stranded REQ only)", report.FixableCount())
	}
	renderedReport := renderVerifyReport(report)
	if !strings.Contains(renderedReport, "1 fixable: run do-work cleanup") {
		t.Errorf("report must route fixable findings to cleanup, got:\n%s", renderedReport)
	}
}

// Requirement 4 / the Constraints: verify is read-only. Nothing it touches may
// change on disk — least of all CHANGELOG.md.
func TestVerifyWritesNothing(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", "**Current version**: 0.164.0\n"},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/queue/REQ-060-finished-but-queued.md", "---\nid: REQ-060\nstatus: completed\ntitle: stranded\ncompleted_at: 2026-08-01T10:00:00Z\n---\n"},
	})

	snapshotBefore := snapshotTreeContents(t, repoRoot)
	if _, verifyError := runVerifyProbes(repoRoot, time.Now()); verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	snapshotAfter := snapshotTreeContents(t, repoRoot)

	if len(snapshotBefore) != len(snapshotAfter) {
		t.Fatalf("file count changed from %d to %d — verify must be read-only", len(snapshotBefore), len(snapshotAfter))
	}
	for relativePath, beforeContent := range snapshotBefore {
		if snapshotAfter[relativePath] != beforeContent {
			t.Errorf("verify modified %s — it must be read-only", relativePath)
		}
	}
}

// snapshotTreeContents maps every file under repoRoot to its content, so a
// read-only assertion can compare the whole tree rather than a sampled file.
func snapshotTreeContents(t *testing.T, repoRoot string) map[string]string {
	t.Helper()
	contentsByPath := map[string]string{}
	walkError := filepath.Walk(repoRoot, func(path string, info os.FileInfo, walkError error) error {
		if walkError != nil || info.IsDir() {
			return nil
		}
		relativePath, relativeError := filepath.Rel(repoRoot, path)
		if relativeError != nil {
			return nil
		}
		fileBytes, readError := os.ReadFile(path)
		if readError != nil {
			return nil
		}
		contentsByPath[relativePath] = string(fileBytes)
		return nil
	})
	if walkError != nil {
		t.Fatalf("walk: %v", walkError)
	}
	return contentsByPath
}

// runGitInFixture runs one git command while building a worktree fixture. Setup
// only — the probes under test do their own shelling out.
func runGitInFixture(t *testing.T, workingDirectory string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", arguments...)
	command.Dir = workingDirectory
	if output, runError := command.CombinedOutput(); runError != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(arguments, " "), runError, output)
	}
}

// newWorktreeFixtureRepo seeds a real single-commit git repo whose release files
// are already clean, so the only findings a worktree test sees are the ones it
// planted. Skips when git is unavailable, matching how the probe itself degrades.
func newWorktreeFixtureRepo(t *testing.T) string {
	t.Helper()
	if _, lookupError := exec.LookPath("git"); lookupError != nil {
		t.Skip("git is not on PATH — the worktree probes skip for the same reason")
	}
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/queue/REQ-071-only-one.md", "---\nid: REQ-071\nstatus: pending\ntitle: fixture\n---\n"},
	})
	runGitInFixture(t, repoRoot, "init", "--quiet")
	runGitInFixture(t, repoRoot, "config", "user.email", "fixture@example.test")
	runGitInFixture(t, repoRoot, "config", "user.name", "Verify Fixture")
	runGitInFixture(t, repoRoot, "add", "-A")
	runGitInFixture(t, repoRoot, "commit", "--quiet", "-m", "fixture base")
	return repoRoot
}

// addFixtureWorktree creates a branch and its worktree under a parent directory
// OUTSIDE the repo, the way worktree dispatch mode requires.
func addFixtureWorktree(t *testing.T, repoRoot string, worktreeParent string, name string) string {
	t.Helper()
	worktreePath := filepath.Join(worktreeParent, name)
	runGitInFixture(t, repoRoot, "worktree", "add", "--quiet", "-b", name, worktreePath)
	return worktreePath
}

// commitFixtureWork puts a commit on whatever branch workingDirectory has checked
// out — what makes a builder branch genuinely unmerged.
func commitFixtureWork(t *testing.T, workingDirectory string, fileName string) {
	t.Helper()
	if writeError := os.WriteFile(filepath.Join(workingDirectory, fileName), []byte("builder work\n"), 0o644); writeError != nil {
		t.Fatalf("write %s: %v", fileName, writeError)
	}
	runGitInFixture(t, workingDirectory, "add", fileName)
	runGitInFixture(t, workingDirectory, "commit", "--quiet", "-m", "builder work")
}

// The contract assertion, not just the finding: an unmerged builder branch is
// cleanup Pass 5's consent-gated path, so it must never be counted as something
// `do-work cleanup` will mechanically resolve. This is the assertion that stops
// the Fixable contract regressing — the finding's presence is the easy half.
func TestVerifyDoesNotAdvertiseAnUnmergedWorktreeAsFixable(t *testing.T) {
	repoRoot := newWorktreeFixtureRepo(t)
	worktreeParent := t.TempDir()

	unmergedWorktree := addFixtureWorktree(t, repoRoot, worktreeParent, "worktree-agent-REQ-001-unmerged")
	commitFixtureWork(t, unmergedWorktree, "builder-output.txt")

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	if report.FixableCount() != 0 {
		t.Errorf("FixableCount = %d, want 0 — unmerged builder work is a human decision:\n%s",
			report.FixableCount(), renderVerifyReport(report))
	}
	unmergedFindings := findingsMentioning(report, verifyCategoryUnmergedWorktreeLeftover)
	if len(unmergedFindings) != 1 {
		t.Fatalf("got %d unmerged-leftover findings, want 1:\n%s", len(unmergedFindings), renderVerifyReport(report))
	}
	if unmergedFindings[0].Fixable {
		t.Error("cleanup Pass 5 asks before discarding unmerged work — it must not be marked fixable")
	}
	// verify cannot tell a live builder from a dead one, so the remedy has to say
	// so rather than let the reader assume the leftover is abandoned.
	if !strings.Contains(unmergedFindings[0].Remedy, "in flight") {
		t.Errorf("the remedy must name the still-running case verify cannot rule out, got: %s", unmergedFindings[0].Remedy)
	}
	if strings.Contains(renderVerifyReport(report), "fixable: run do-work cleanup") {
		t.Errorf("a lone unmerged leftover must not route the user to cleanup:\n%s", renderVerifyReport(report))
	}
}

// The full state table REQ-072 deferred. One fixture, five leftovers, so the
// merged/unmerged split is proved to be a classification rather than a blanket
// demotion — and a developer's own worktree is proved to stay invisible.
func TestVerifyClassifiesWorktreeLeftoversByMergeState(t *testing.T) {
	repoRoot := newWorktreeFixtureRepo(t)
	worktreeParent := t.TempDir()

	// 1. Unmerged, worktree still present — the in-flight-or-dead case.
	unmergedWithWorktree := addFixtureWorktree(t, repoRoot, worktreeParent, "worktree-agent-REQ-001-unmerged-live")
	commitFixtureWork(t, unmergedWithWorktree, "one.txt")

	// 2. Unmerged, worktree already gone — the branch outlived its worktree.
	unmergedBranchOnly := addFixtureWorktree(t, repoRoot, worktreeParent, "worktree-agent-REQ-002-unmerged-branch")
	commitFixtureWork(t, unmergedBranchOnly, "two.txt")
	runGitInFixture(t, repoRoot, "worktree", "remove", unmergedBranchOnly)

	// 3. Merged, worktree still present — pure residue Pass 5 removes mechanically.
	mergedWithWorktree := addFixtureWorktree(t, repoRoot, worktreeParent, "worktree-agent-REQ-003-merged-live")
	commitFixtureWork(t, mergedWithWorktree, "three.txt")
	runGitInFixture(t, repoRoot, "merge", "--no-ff", "--no-edit", "--quiet", "worktree-agent-REQ-003-merged-live")

	// 4. Merged, branch only — points at HEAD, so it is trivially contained.
	runGitInFixture(t, repoRoot, "branch", "worktree-agent-REQ-004-merged-branch")

	// 5. A worktree that is not a builder's. It must be ignored entirely.
	runGitInFixture(t, repoRoot, "worktree", "add", "--quiet", "-b", "my-own-feature",
		filepath.Join(worktreeParent, "my-own-feature"))

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	renderedReport := renderVerifyReport(report)

	categoryByLeftover := map[string]string{}
	for _, category := range []string{verifyCategoryMergedWorktreeLeftover, verifyCategoryUnmergedWorktreeLeftover, verifyCategoryUndeterminedWorktreeLeftover} {
		for _, finding := range findingsMentioning(report, category) {
			for _, leftoverName := range []string{
				"worktree-agent-REQ-001-unmerged-live",
				"worktree-agent-REQ-002-unmerged-branch",
				"worktree-agent-REQ-003-merged-live",
				"worktree-agent-REQ-004-merged-branch",
			} {
				if strings.Contains(finding.Detail, leftoverName) {
					categoryByLeftover[leftoverName] = category
				}
			}
		}
	}

	expectedCategoryByLeftover := map[string]string{
		"worktree-agent-REQ-001-unmerged-live":   verifyCategoryUnmergedWorktreeLeftover,
		"worktree-agent-REQ-002-unmerged-branch": verifyCategoryUnmergedWorktreeLeftover,
		"worktree-agent-REQ-003-merged-live":     verifyCategoryMergedWorktreeLeftover,
		"worktree-agent-REQ-004-merged-branch":   verifyCategoryMergedWorktreeLeftover,
	}
	for leftoverName, wantCategory := range expectedCategoryByLeftover {
		if gotCategory := categoryByLeftover[leftoverName]; gotCategory != wantCategory {
			t.Errorf("%s classified as %q, want %q:\n%s", leftoverName, gotCategory, wantCategory, renderedReport)
		}
	}

	if strings.Contains(renderedReport, "my-own-feature") {
		t.Errorf("a developer's own worktree must never be reported:\n%s", renderedReport)
	}
	// Exactly the two merged leftovers are fixable — proof this is a
	// classification, not a blanket demotion of the whole probe.
	if report.FixableCount() != 2 {
		t.Errorf("FixableCount = %d, want 2 (the merged leftovers only):\n%s", report.FixableCount(), renderedReport)
	}
}

// A worktree whose name matches the convention but has no branch behind it (a
// detached checkout) is git declining to answer, not an answer. Reporting an
// unanswered question as "unmerged" would be the same class of defect this probe
// was fixed for, so it gets its own state — and, like unmerged, it is never
// advertised as something cleanup can mechanically resolve.
func TestVerifyReportsAnUndeterminedMergeStateSeparately(t *testing.T) {
	repoRoot := newWorktreeFixtureRepo(t)
	worktreeParent := t.TempDir()

	runGitInFixture(t, repoRoot, "worktree", "add", "--quiet", "--detach",
		filepath.Join(worktreeParent, "worktree-agent-REQ-005-detached"))

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	undeterminedFindings := findingsMentioning(report, verifyCategoryUndeterminedWorktreeLeftover)
	if len(undeterminedFindings) != 1 {
		t.Fatalf("got %d undetermined-state findings, want 1:\n%s", len(undeterminedFindings), renderVerifyReport(report))
	}
	if undeterminedFindings[0].Fixable {
		t.Error("verify could not establish a merge target — it must not be advertised as cleanup-fixable")
	}
	if report.FixableCount() != 0 {
		t.Errorf("FixableCount = %d, want 0:\n%s", report.FixableCount(), renderVerifyReport(report))
	}
}

// A repo whose changelog does not follow the house entry key must not produce
// bogus release findings — the Changelog Entry Procedure's precedence rule says
// match the existing format, so these two probes skip with a note instead.
func TestVerifySkipsReleaseProbesOnAForeignChangelogFormat(t *testing.T) {
	keepAChangelogBody := `# Changelog

## [Unreleased]
### Added
- Something.
`
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", keepAChangelogBody},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	if len(findingsMentioning(report, verifyCategoryVersionChangelogMismatch)) != 0 {
		t.Error("a foreign changelog format must not yield a version-ordering finding")
	}
	if len(findingsMentioning(report, verifyCategoryReusedChangelogTitle)) != 0 {
		t.Error("a foreign changelog format must not yield a reused-title finding")
	}
	// The skip must be reported and must name the changelog — silence would read
	// as "checked and clean," and a generic skip line would not tell the user
	// which invariant went unverified.
	skippedChangelogProbe := false
	for _, skippedProbe := range report.SkippedProbes {
		if strings.Contains(skippedProbe, "CHANGELOG.md") {
			skippedChangelogProbe = true
		}
	}
	if !skippedChangelogProbe {
		t.Errorf("the changelog-format skip must be reported by name, got: %v", report.SkippedProbes)
	}
}
