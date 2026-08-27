package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
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
		// A clean REQ carries its user_request upward pointer. Both clean-base
		// fixtures omitted it until the structural probe landed, which made
		// "clean" mean something no captured REQ actually looks like.
		{"do-work/queue/REQ-071-only-one.md", "---\nid: REQ-071\nstatus: pending\ntitle: fixture\nuser_request: UR-071\n---\n"},
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
// finding, and each direction has its own cause. The mid-release form of the rule is
// "strictly greater than the entry that already exists", which becomes equality the
// moment the entry is written.
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
// within the changelog itself — this is the duplicate-version-number failure that
// has already happened in this repo more than once.
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

// The newest changelog entry's title must not already be in use — the second of the
// two release invariants the commit ritual verifies by eye.
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

func TestVerifyFlagsStrayRequestFiles(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/user-requests/UR-095/input.md", "---\nid: UR-095\ntitle: open UR with misplaced REQs\nrequests: []\n---\n"},
		{"do-work/user-requests/UR-095/REQ-090-pending-stray.md",
			"---\nid: REQ-090\nstatus: pending\ntitle: pending stray\nuser_request: UR-095\n---\n"},
		{"do-work/user-requests/UR-095/REQ-091-completed-stray.md",
			"---\nid: REQ-091\nstatus: completed\ntitle: completed stray\nuser_request: UR-095\ncompleted_at: 2026-08-01T10:00:00Z\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	strayFindings := findingsMentioning(report, "stray-req-file")
	if len(strayFindings) != 2 {
		t.Fatalf("got %d stray-REQ findings, want 2 (pending and completed):\n%s",
			len(strayFindings), renderVerifyReport(report))
	}
	for _, expectedPath := range []string{
		"do-work/user-requests/UR-095/REQ-090-pending-stray.md",
		"do-work/user-requests/UR-095/REQ-091-completed-stray.md",
	} {
		matchingCount := 0
		for _, finding := range strayFindings {
			if strings.Contains(finding.Detail, expectedPath) {
				matchingCount++
			}
			if finding.Fixable {
				t.Errorf("stray relocation is a human decision and must not be cleanup-fixable: %+v", finding)
			}
		}
		if matchingCount != 1 {
			t.Errorf("path %s appeared in %d stray findings, want exactly 1", expectedPath, matchingCount)
		}
	}
	if strandedFindings := findingsMentioning(report, verifyCategoryStrandedFinishedRequest); len(strandedFindings) != 0 {
		t.Fatalf("a completed stray must remain outside normal request probes, got %d stranded findings:\n%s",
			len(strandedFindings), renderVerifyReport(report))
	}
}

func TestAppendStrayRequestFileFindingsUsesStructuredEvidence(t *testing.T) {
	structuredPath := "user-requests/UR-096/REQ-095-structured-only.md"
	board := &Board{
		RepoRoot: filepath.Join(t.TempDir(), "missing-repo"),
		StrayRequestFiles: []strayRequestFile{
			{RelativePath: structuredPath},
		},
		Warnings: nil,
	}
	var report VerifyReport
	if _, statError := os.Stat(board.RepoRoot); !os.IsNotExist(statError) {
		t.Fatalf("fixture repo root must remain absent so a filesystem re-walk has no evidence; stat error = %v", statError)
	}

	appendStrayRequestFileFindings(&report, board)

	if len(report.Findings) != 1 {
		t.Fatalf("got %d findings from structured stray evidence, want exactly 1: %+v",
			len(report.Findings), report.Findings)
	}
	finding := report.Findings[0]
	if finding.Category != verifyCategoryStrayRequestFile {
		t.Errorf("category = %q, want %q", finding.Category, verifyCategoryStrayRequestFile)
	}
	if !strings.Contains(finding.Detail, structuredPath) {
		t.Errorf("detail %q does not name structured path %q", finding.Detail, structuredPath)
	}
	if finding.Fixable {
		t.Errorf("structured stray finding must remain non-fixable: %+v", finding)
	}
}

// A terminally-resolved member stranded in queue/ under an ARCHIVED UR used to
// trip both probes, and the archived-UR remedy then told the user to run or
// abandon a REQ that was already resolved. The stranded-finished probe owns that
// state (it is fixable, and cleanup Pass 0 sweeps it), so the archived-UR probe
// must stay silent on it. Both halves of the terminal set are covered here —
// `completed` and `cancelled` — because the carve-out keys on
// isTerminalResolvedStatus, not on completion alone. The converse case (an
// archived UR whose queued member is genuinely unresolved must still be flagged)
// is TestVerifyFlagsArchivedUserRequestWithALiveMember, below.
func TestVerifyStrandedFinishedMemberOfAnArchivedUserRequestFiresOnlyTheStrandedProbe(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/UR-094/input.md", "---\nid: UR-094\ntitle: archived with stranded members\nrequests: [REQ-086, REQ-088]\n---\n"},
		{"do-work/queue/REQ-086-finished-but-queued.md",
			"---\nid: REQ-086\nstatus: completed\ntitle: stranded completed member\nuser_request: UR-094\ncompleted_at: 2026-08-01T10:00:00Z\ncommit: abc1234\n---\n"},
		{"do-work/queue/REQ-088-cancelled-but-queued.md",
			"---\nid: REQ-088\nstatus: cancelled\ntitle: stranded cancelled member\nuser_request: UR-094\ncompleted_at: 2026-08-01T11:00:00Z\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	strandedFindings := findingsMentioning(report, verifyCategoryStrandedFinishedRequest)
	if len(strandedFindings) != 2 {
		t.Fatalf("got %d stranded-finished findings, want 2 (one completed, one cancelled):\n%s",
			len(strandedFindings), renderVerifyReport(report))
	}
	if found := findingsMentioning(report, verifyCategoryArchivedUserRequestLiveMember); len(found) != 0 {
		t.Fatalf("a terminally-resolved member is not live — the archived-UR probe must not double-fire, got %d findings:\n%s",
			len(found), renderVerifyReport(report))
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
		// A clean REQ carries its user_request upward pointer. Both clean-base
		// fixtures omitted it until the structural probe landed, which made
		// "clean" mean something no captured REQ actually looks like.
		{"do-work/queue/REQ-071-only-one.md", "---\nid: REQ-071\nstatus: pending\ntitle: fixture\nuser_request: UR-071\n---\n"},
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

// writeFixtureFile writes one file under a checkout, creating parents. Used to
// plant queue state where it does or does not belong.
func writeFixtureFile(t *testing.T, checkoutPath string, relativePath string, content string) {
	t.Helper()
	absolutePath := filepath.Join(checkoutPath, relativePath)
	if mkdirError := os.MkdirAll(filepath.Dir(absolutePath), 0o755); mkdirError != nil {
		t.Fatalf("mkdir for %s: %v", relativePath, mkdirError)
	}
	if writeError := os.WriteFile(absolutePath, []byte(content), 0o644); writeError != nil {
		t.Fatalf("write %s: %v", relativePath, writeError)
	}
}

const forgedQueueRequest = "---\nid: REQ-999\nstatus: pending\ntitle: forged by a builder\n---\n"

// The gap REQ-072 requirement 3 actually asked about: a builder that wrote queue
// state and then COMMITTED it on its own branch. That is the likely shape, since
// a builder commits its work by design — and it leaves the worktree clean, so a
// porcelain-only probe sees nothing.
func TestVerifyFlagsQueueStateCommittedOnABuilderBranch(t *testing.T) {
	repoRoot := newWorktreeFixtureRepo(t)
	worktreeParent := t.TempDir()

	builderWorktree := addFixtureWorktree(t, repoRoot, worktreeParent, "worktree-agent-REQ-010-impersonator")
	writeFixtureFile(t, builderWorktree, "do-work/queue/REQ-999-forged.md", forgedQueueRequest)
	runGitInFixture(t, builderWorktree, "add", "do-work/queue/REQ-999-forged.md")
	runGitInFixture(t, builderWorktree, "commit", "--quiet", "-m", "builder wrote queue state")

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	committedFindings := findingsMentioning(report, verifyCategoryWorktreeCommittedQueueState)
	if len(committedFindings) != 1 {
		t.Fatalf("got %d committed-queue-state findings, want 1:\n%s", len(committedFindings), renderVerifyReport(report))
	}
	if !strings.Contains(committedFindings[0].Detail, "REQ-999-forged.md") {
		t.Errorf("finding must name the path the builder wrote, got: %s", committedFindings[0].Detail)
	}
	// The two states need different remedies — one is inside the branch about to
	// be merged, the other is loose in a working tree — so the detail must say
	// which was found rather than leaving the reader to guess.
	if !strings.Contains(committedFindings[0].Detail, "committed") {
		t.Errorf("finding must say the edits are committed on the branch, got: %s", committedFindings[0].Detail)
	}
	if committedFindings[0].Fixable {
		t.Error("discarding a builder's commits is never mechanical — this must not be marked fixable")
	}
}

// Requirement 1: the merge-base comparison is added ALONGSIDE the porcelain
// check, not in place of it. A merge-base diff cannot see uncommitted edits.
func TestVerifyStillFlagsUncommittedQueueStateInAWorktree(t *testing.T) {
	repoRoot := newWorktreeFixtureRepo(t)
	worktreeParent := t.TempDir()

	builderWorktree := addFixtureWorktree(t, repoRoot, worktreeParent, "worktree-agent-REQ-011-loose-edit")
	writeFixtureFile(t, builderWorktree, "do-work/queue/REQ-999-forged.md", forgedQueueRequest)

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	dirtyFindings := findingsMentioning(report, verifyCategoryWorktreeWroteQueueState)
	if len(dirtyFindings) != 1 {
		t.Fatalf("got %d uncommitted-queue-state findings, want 1:\n%s", len(dirtyFindings), renderVerifyReport(report))
	}
	if dirtyFindings[0].Fixable {
		t.Error("uncommitted builder queue edits are a human decision — not fixable")
	}
}

// Requirement 4, and the regression test that makes widening the probe safe: a
// worktree whose do-work/ is merely BEHIND the main tree is the legitimate stale
// snapshot the original narrowing was protecting. Three-dot diff semantics are
// what keep it silent — the orchestrator's later queue commits are on the
// integration branch's side of the merge base, not the builder's.
//
// It also pins requirement 6: REQ-082 grants a builder exactly one main-tree
// write, its own REQ-NNN-handback.md, absolute and never committed. Planted here
// in the main tree, it must be invisible to both states.
func TestVerifyIgnoresAStaleQueueSnapshotAndTheHandBackFile(t *testing.T) {
	repoRoot := newWorktreeFixtureRepo(t)
	worktreeParent := t.TempDir()

	// Branch point: the builder's worktree is created here and never touches do-work/.
	addFixtureWorktree(t, repoRoot, worktreeParent, "worktree-agent-REQ-012-innocent")

	// The orchestrator then moves the main tree on, exactly as a live run does.
	writeFixtureFile(t, repoRoot, "do-work/archive/REQ-071-archived-later.md",
		"---\nid: REQ-071\nstatus: completed\ntitle: archived after the branch point\ncompleted_at: 2026-08-03T10:00:00Z\ncommit: abc1234\n---\n")
	runGitInFixture(t, repoRoot, "add", "-A")
	runGitInFixture(t, repoRoot, "commit", "--quiet", "-m", "orchestrator archived a REQ")

	// REQ-082's one permitted main-tree write by the builder: absolute path, never
	// committed. It must not read as an owner impersonation.
	writeFixtureFile(t, repoRoot, "do-work/runs/work-2026-08-03-120000/REQ-012-handback.md",
		"# Hand-back\n\nbranch: worktree-agent-REQ-012-innocent\n")

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	if committedFindings := findingsMentioning(report, verifyCategoryWorktreeCommittedQueueState); len(committedFindings) != 0 {
		t.Errorf("a stale snapshot is not an impersonation — got %d findings:\n%s", len(committedFindings), renderVerifyReport(report))
	}
	if dirtyFindings := findingsMentioning(report, verifyCategoryWorktreeWroteQueueState); len(dirtyFindings) != 0 {
		t.Errorf("the main tree's uncommitted hand-back file must not be read as a worktree write — got %d findings:\n%s", len(dirtyFindings), renderVerifyReport(report))
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

// TestParsePorcelainStatusPathsKeepsSpacesAndRenames pins the porcelain-v1
// parsing contract: the path is everything after the fixed 3-byte prefix, so a
// path containing spaces survives whole (the last-whitespace-field parse used
// to truncate "REQ-12 draft copy.md" to "copy.md" in the dirty-worktree report
// line), a rename keeps its destination side, and git's double-quoting of
// special-character paths is stripped for display.
func TestParsePorcelainStatusPathsKeepsSpacesAndRenames(t *testing.T) {
	statusOutput := strings.Join([]string{
		` M do-work/queue/REQ-12 draft copy.md`,
		`?? do-work/queue/REQ-13-normal.md`,
		`R  do-work/queue/REQ-14-old.md -> do-work/queue/REQ-14-new.md`,
		`?? "do-work/queue/REQ-15 quoted name.md"`,
		``,
	}, "\n")
	got := parsePorcelainStatusPaths(statusOutput)
	want := []string{
		"do-work/queue/REQ-12 draft copy.md",
		"do-work/queue/REQ-13-normal.md",
		"do-work/queue/REQ-14-new.md",
		"do-work/queue/REQ-15 quoted name.md",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parsePorcelainStatusPaths mismatch:\ngot  %q\nwant %q", got, want)
	}
}

// TestCheckpointMentionedRequestIdsSkipsGlobPatterns pins the false positive a
// real CHECKPOINT.md produced: session notes quoting the shell glob
// `REQ-0[0-9][0-9]-*.md` made the scanner report a ghost "REQ-0". A digit run
// continuing into `[` is a pattern, not an id; ids wrapped in markdown
// emphasis or ending a question must still be found.
func TestCheckpointMentionedRequestIdsSkipsGlobPatterns(t *testing.T) {
	checkpointText := strings.Join([]string{
		"last_completed: REQ-093",
		"cleanup ran `rm -f kb/raw/inbox/REQ-0[0-9][0-9]-*.md` over the inbox",
		"**REQ-088** closed; did REQ-074? REQ-093 again.",
	}, "\n")
	got := checkpointMentionedRequestIds(checkpointText)
	want := []string{"REQ-093", "REQ-088", "REQ-074"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("checkpointMentionedRequestIds mismatch:\ngot  %q\nwant %q", got, want)
	}
}

// assigned_to is advisory, and clearing it is part of Step 2's claim
// (actions/work.md). A REQ that reached do-work/working/ still carrying the field
// means either the claim skipped the clear, or a session claimed work earmarked
// for somebody else without meaning to. Either way the marker is now lying to
// every other checkout: it says "skip me", about a REQ already being built here.
func TestVerifyFlagsAssignedRequestClaimedHere(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/working/REQ-070-claimed-but-assigned.md",
			"---\nid: REQ-070\nstatus: claimed\ntitle: claimed but still earmarked\nassigned_to: \"cloud-alpha\"\nclaimed_at: " +
				time.Now().UTC().Format(time.RFC3339) + "\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	assignedFindings := findingsMentioning(report, verifyCategoryAssignedElsewhereClaimedHere)
	if len(assignedFindings) != 1 {
		t.Fatalf("got %d assigned-elsewhere findings, want 1:\n%s", len(assignedFindings), renderVerifyReport(report))
	}
	if !strings.Contains(assignedFindings[0].Detail, "cloud-alpha") {
		t.Errorf("finding must name the assignee verbatim so the reader knows whose marker it is, got %q", assignedFindings[0].Detail)
	}
	if assignedFindings[0].Fixable {
		t.Error("clearing somebody else's claim marker is a human decision — cleanup must not advertise it as mechanical")
	}
}

func TestVerifyIgnoresAnAssignedRequestStillInTheQueue(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		// Earmarked and NOT claimed here: this is the field working exactly as designed.
		{"do-work/queue/REQ-071-earmarked.md",
			"---\nid: REQ-071\nstatus: pending\ntitle: earmarked and left alone\nassigned_to: \"cloud-alpha\"\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	if found := findingsMentioning(report, verifyCategoryAssignedElsewhereClaimedHere); len(found) != 0 {
		t.Fatalf("an assigned REQ sitting in the queue is the normal case, got %d findings:\n%s",
			len(found), renderVerifyReport(report))
	}
}

// A terminally-resolved REQ stranded in do-work/working/ while still carrying
// assigned_to used to trip both probes, and the assigned-elsewhere remedy then told
// the user to clear or release a claim on work that was already done. The
// stranded-finished probe owns that state (it is fixable, and cleanup Pass 0 sweeps
// it), so the assigned-elsewhere probe must stay silent on it. Both halves of the
// terminal set are covered here — `completed` and `cancelled` — because the
// carve-out keys on isTerminalResolvedStatus, not on completion alone. The converse
// case (a non-terminal assigned REQ in working/ must still be flagged) is
// TestVerifyFlagsAssignedRequestClaimedHere, above.
func TestVerifyStrandedFinishedAssignedRequestFiresOnlyTheStrandedProbe(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/working/REQ-072-finished-but-earmarked.md",
			"---\nid: REQ-072\nstatus: completed\ntitle: stranded completed claim\nassigned_to: \"cloud-alpha\"\ncompleted_at: 2026-08-01T10:00:00Z\ncommit: abc1234\n---\n"},
		{"do-work/working/REQ-073-cancelled-but-earmarked.md",
			"---\nid: REQ-073\nstatus: cancelled\ntitle: stranded cancelled claim\nassigned_to: \"cloud-beta\"\ncompleted_at: 2026-08-01T11:00:00Z\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	strandedFindings := findingsMentioning(report, verifyCategoryStrandedFinishedRequest)
	if len(strandedFindings) != 2 {
		t.Fatalf("got %d stranded-finished findings, want 2 (one completed, one cancelled):\n%s",
			len(strandedFindings), renderVerifyReport(report))
	}
	if found := findingsMentioning(report, verifyCategoryAssignedElsewhereClaimedHere); len(found) != 0 {
		t.Fatalf("a terminally-resolved REQ is not being built here — the assigned-elsewhere probe must not double-fire, got %d findings:\n%s",
			len(found), renderVerifyReport(report))
	}
}

func TestVerifyAllowsReviewGeneratedMemberUnderClosedUserRequest(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/UR-095/input.md", "---\nid: UR-095\ntitle: closed UR with review follow-ups\nrequests: [REQ-090]\n---\n"},
		{"do-work/queue/REQ-090-review-generated.md",
			"---\nid: REQ-090\nstatus: pending\ntitle: legitimate queued review follow-up\nuser_request: UR-095\nreview_generated: true\n---\n"},
		{"do-work/working/REQ-091-review-generated.md",
			"---\nid: REQ-091\nstatus: claimed\ntitle: legitimate working review follow-up\nuser_request: UR-095\nreview_generated: true\nclaimed_at: " + time.Now().UTC().Format(time.RFC3339) + "\n---\n"},
		{"do-work/queue/REQ-092-ordinary.md",
			"---\nid: REQ-092\nstatus: pending\ntitle: ordinary live member\nuser_request: UR-095\nreview_generated: false\n---\n"},
		{"do-work/working/REQ-093-noncanonical-marker.md",
			"---\nid: REQ-093\nstatus: claimed\ntitle: noncanonical marker is ordinary\nuser_request: UR-095\nreview_generated: \"TRUE\"\nclaimed_at: " + time.Now().UTC().Format(time.RFC3339) + "\n---\n"},
		{"do-work/queue/REQ-094-finished-review-generated.md",
			"---\nid: REQ-094\nstatus: completed\ntitle: stranded finished review follow-up\nuser_request: UR-095\nreview_generated: true\ncompleted_at: 2026-08-01T10:00:00Z\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	liveMemberFindings := findingsMentioning(report, verifyCategoryArchivedUserRequestLiveMember)
	if len(liveMemberFindings) != 1 {
		t.Fatalf("got %d archived-UR-live-member findings, want 1 for the ordinary siblings:\n%s",
			len(liveMemberFindings), renderVerifyReport(report))
	}
	detail := liveMemberFindings[0].Detail
	for _, ordinaryRequestId := range []string{"REQ-092", "REQ-093"} {
		if !strings.Contains(detail, ordinaryRequestId) {
			t.Errorf("ordinary member %s must remain reported, got %q", ordinaryRequestId, detail)
		}
	}
	for _, exemptRequestId := range []string{"REQ-090", "REQ-091", "REQ-094"} {
		if strings.Contains(detail, exemptRequestId) {
			t.Errorf("review-generated or terminal member %s must not be reported as an ordinary live anomaly, got %q", exemptRequestId, detail)
		}
	}
	if liveMemberFindings[0].Fixable {
		t.Error("resolving an ordinary live member under a closed UR is a human decision — not mechanical")
	}
	if strings.Contains(liveMemberFindings[0].Remedy, "user-requests/") ||
		strings.Contains(strings.ToLower(liveMemberFindings[0].Remedy), "bring the ur folder back") {
		t.Errorf("remedy must keep the archived UR closed, got %q", liveMemberFindings[0].Remedy)
	}
	if !strings.Contains(strings.ToLower(liveMemberFindings[0].Remedy), "archived") {
		t.Errorf("remedy must explicitly preserve the UR's archived state, got %q", liveMemberFindings[0].Remedy)
	}

	strandedFindings := findingsMentioning(report, verifyCategoryStrandedFinishedRequest)
	if len(strandedFindings) != 1 || !strings.Contains(strandedFindings[0].Detail, "REQ-094") {
		t.Fatalf("terminal review-generated member must remain owned by stranded-finished, got %d findings:\n%s",
			len(strandedFindings), renderVerifyReport(report))
	}
}

// A UR moves to do-work/archive/ only once every member REQ is terminally
// resolved (actions/work.md Step 8). A member still in queue/ or working/ after
// the UR was archived means the closure check passed on stale information, or a
// cleanup pass moved the folder by hand — and the live REQ is now orphaned from
// the input.md that explains it.
func TestVerifyFlagsArchivedUserRequestWithALiveMember(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/UR-090/input.md", "---\nid: UR-090\ntitle: archived too early\nrequests: [REQ-080]\n---\n"},
		{"do-work/archive/UR-090/REQ-079-done.md",
			"---\nid: REQ-079\nstatus: completed\ntitle: done\nuser_request: UR-090\ncompleted_at: 2026-08-01T10:00:00Z\ncommit: abc1234\n---\n"},
		{"do-work/queue/REQ-080-still-live.md",
			"---\nid: REQ-080\nstatus: pending\ntitle: still live\nuser_request: UR-090\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	liveMemberFindings := findingsMentioning(report, verifyCategoryArchivedUserRequestLiveMember)
	if len(liveMemberFindings) != 1 {
		t.Fatalf("got %d archived-UR-live-member findings, want 1:\n%s",
			len(liveMemberFindings), renderVerifyReport(report))
	}
	detail := liveMemberFindings[0].Detail
	if !strings.Contains(detail, "UR-090") || !strings.Contains(detail, "REQ-080") {
		t.Errorf("finding must name both the archived UR and the live member, got %q", detail)
	}
	if strings.Contains(detail, "REQ-079") {
		t.Errorf("the terminally-resolved member is not the problem and must not be listed, got %q", detail)
	}
	if liveMemberFindings[0].Fixable {
		t.Error("un-archiving a UR or force-resolving a live REQ are both human decisions — not mechanical")
	}
}

func TestVerifyIgnoresAnArchivedUserRequestWhoseMembersAreAllResolved(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/UR-091/input.md", "---\nid: UR-091\ntitle: properly closed\nrequests: [REQ-081]\n---\n"},
		{"do-work/archive/UR-091/REQ-081-done.md",
			"---\nid: REQ-081\nstatus: completed\ntitle: done\nuser_request: UR-091\ncompleted_at: 2026-08-01T10:00:00Z\ncommit: abc1234\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	if found := findingsMentioning(report, verifyCategoryArchivedUserRequestLiveMember); len(found) != 0 {
		t.Fatalf("a fully-closed archived UR is the normal case, got %d findings:\n%s",
			len(found), renderVerifyReport(report))
	}
}

// The UR's own `requests:` array is capture-time-only and can be wrong in both
// directions; the closure predicate actions/work.md Step 8 evaluates is a scan of
// `user_request:` frontmatter. This fixture makes the two disagree: the array omits
// the live member entirely, so a probe reading the array sees a closed UR.
func TestVerifyLiveMemberProbeScansUserRequestFrontmatterNotTheUrArray(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/UR-092/input.md", "---\nid: UR-092\ntitle: array disagrees with the tree\nrequests: [REQ-082]\n---\n"},
		{"do-work/archive/UR-092/REQ-082-done.md",
			"---\nid: REQ-082\nstatus: completed\ntitle: done\nuser_request: UR-092\ncompleted_at: 2026-08-01T10:00:00Z\ncommit: abc1234\n---\n"},
		// Points at UR-092 by frontmatter, absent from its requests: array.
		{"do-work/working/REQ-083-unlisted-live-member.md",
			"---\nid: REQ-083\nstatus: claimed\ntitle: unlisted but live\nuser_request: UR-092\nclaimed_at: " +
				time.Now().UTC().Format(time.RFC3339) + "\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	liveMemberFindings := findingsMentioning(report, verifyCategoryArchivedUserRequestLiveMember)
	if len(liveMemberFindings) != 1 {
		t.Fatalf("got %d findings, want 1 — the probe must scan user_request: frontmatter, not the UR's requests: array:\n%s",
			len(liveMemberFindings), renderVerifyReport(report))
	}
	if !strings.Contains(liveMemberFindings[0].Detail, "REQ-083") {
		t.Errorf("finding must name the unlisted live member, got %q", liveMemberFindings[0].Detail)
	}
}

// resolveRepoRootOrDefault returns an explicit --repo-root override verbatim, so
// `verify --repo-root .` produces ticket paths with no leading separator. The
// archived-UR probe used to require a leading "/do-work/archive/" and therefore
// recognized nothing in that mode — silently, which is the worst way for a probe to
// fail. Every other test here builds fixtures under t.TempDir() (absolute), so none
// of them could catch it. Found by review on PR #128.
func TestVerifyFlagsArchivedUserRequestLiveMemberUnderARelativeRepoRoot(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/UR-093/input.md", "---\nid: UR-093\ntitle: archived early\nrequests: [REQ-084]\n---\n"},
		{"do-work/archive/UR-093/REQ-084-done.md",
			"---\nid: REQ-084\nstatus: completed\ntitle: done\nuser_request: UR-093\ncompleted_at: 2026-08-01T10:00:00Z\ncommit: abc1234\n---\n"},
		{"do-work/queue/REQ-085-still-live.md",
			"---\nid: REQ-085\nstatus: pending\ntitle: still live\nuser_request: UR-093\n---\n"},
	})

	originalWorkingDirectory, getwdError := os.Getwd()
	if getwdError != nil {
		t.Fatalf("Getwd: %v", getwdError)
	}
	if chdirError := os.Chdir(repoRoot); chdirError != nil {
		t.Fatalf("Chdir: %v", chdirError)
	}
	t.Cleanup(func() { _ = os.Chdir(originalWorkingDirectory) })

	// "." is the shape actions/forensics.md Check 14 can pass, and the shape the
	// leading-slash match could never satisfy.
	report, verifyError := runVerifyProbes(".", time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	liveMemberFindings := findingsMentioning(report, verifyCategoryArchivedUserRequestLiveMember)
	if len(liveMemberFindings) != 1 {
		t.Fatalf("got %d findings under a relative repo root, want 1 — the probe must not depend on a leading separator:\n%s",
			len(liveMemberFindings), renderVerifyReport(report))
	}
	if !strings.Contains(liveMemberFindings[0].Detail, "REQ-085") {
		t.Errorf("finding must name the live member, got %q", liveMemberFindings[0].Detail)
	}
}

// The same probe must still refuse a directory merely named "archive" that is not
// this project's do-work archive — the reason the match is separator-anchored.
func TestIsArchivedUserRequestPathRejectsLookalikeDirectories(t *testing.T) {
	for _, archivedPath := range []string{
		"do-work/archive/UR-001/input.md",
		"/abs/repo/do-work/archive/UR-001/input.md",
		"./do-work/archive/UR-001/input.md",
	} {
		if !isArchivedUserRequestPath(archivedPath) {
			t.Errorf("isArchivedUserRequestPath(%q) = false, want true", archivedPath)
		}
	}
	for _, livePath := range []string{
		"do-work/user-requests/UR-001/input.md",
		"/abs/repo/do-work/user-requests/UR-001/input.md",
		"my-do-work/archive/UR-001/input.md", // not this project's do-work/
		"archive/UR-001/input.md",            // bare archive/, no do-work/ segment
		"docs/archive/UR-001/input.md",
	} {
		if isArchivedUserRequestPath(livePath) {
			t.Errorf("isArchivedUserRequestPath(%q) = true, want false", livePath)
		}
	}
}

// verify was blind to completion anomalies until REQ-214: it reported
// "OK: no findings" on a tree whose summary showed a flagged anomaly strip.
// The probe forwards buildBoard's structured evidence, one finding per ticket.
func TestVerifyLiftsCompletionAnomaliesIntoFindings(t *testing.T) {
	board := &Board{}
	board.Columns.CompletionAnomalies = []*RequestTicket{{
		RequestId:               "REQ-9330",
		Status:                  "completed",
		CompletionAnomaly:       true,
		CompletionAnomalyReason: reversedSpanAnomalyReason(t),
	}}
	report := VerifyReport{}
	appendCompletionAnomalyFindings(&report, ".", board)
	if len(report.Findings) != 1 {
		t.Fatalf("findings = %d, want 1", len(report.Findings))
	}
	finding := report.Findings[0]
	if finding.Category != verifyCategoryCompletionAnomaly {
		t.Fatalf("category = %q, want %q", finding.Category, verifyCategoryCompletionAnomaly)
	}
	if !strings.Contains(finding.Detail, "REQ-9330") || !strings.Contains(finding.Detail, "is earlier than claimed_at") {
		t.Fatalf("detail = %q, want the ticket id and its reason forwarded", finding.Detail)
	}

	cleanReport := VerifyReport{}
	appendCompletionAnomalyFindings(&cleanReport, ".", &Board{})
	if len(cleanReport.Findings) != 0 {
		t.Fatalf("clean board produced %d findings, want 0", len(cleanReport.Findings))
	}
}

// When git (or the repository) is unavailable, a hash-only anomaly is
// indistinguishable from a healthy record — the same dating probe fails for
// every valid hash — so it routes to SkippedProbes per the ExitCode contract.
// Classes that are on-disk defects regardless of git (a reversed span here)
// still fail the check in the same environment.
func TestHashOnlyAnomalySkippedWhenGitUnavailable(t *testing.T) {
	nonRepoRoot := t.TempDir()
	board := &Board{}
	board.Columns.CompletionAnomalies = []*RequestTicket{
		{
			RequestId:               "REQ-9331",
			Status:                  "completed",
			CommitHash:              "deadbeef",
			CompletionTimeSource:    CompletionUnresolved,
			CompletionAnomalyReason: `commit "deadbeef" could not be dated — the hash is unknown to git, or git/the repository is unavailable`,
		},
		{
			RequestId:               "REQ-9332",
			Status:                  "completed",
			ClaimedAt:               "2026-01-02T10:00:00Z",
			CompletedAt:             "2026-01-01T10:00:00Z",
			CompletionAnomaly:       true,
			CompletionAnomalyReason: reversedSpanAnomalyReason(t),
		},
	}
	report := VerifyReport{}
	appendCompletionAnomalyFindings(&report, nonRepoRoot, board)
	if len(report.Findings) != 1 || !strings.Contains(report.Findings[0].Detail, "REQ-9332") {
		t.Fatalf("findings = %v, want exactly the reversed-span ticket REQ-9332", report.Findings)
	}
	sawSkip := false
	for _, skippedProbe := range report.SkippedProbes {
		if strings.Contains(skippedProbe, "REQ-9331") {
			sawSkip = true
		}
	}
	if !sawSkip {
		t.Fatalf("SkippedProbes = %v, want the hash-only ticket REQ-9331 routed there", report.SkippedProbes)
	}
}

// The captured RED from REQ-280: an archived REQ claimed before it was created,
// with a forward claimed_at → completed_at span so detectCompletionAnomaly stays
// silent. Before the ordering probe this fixture produced `OK: no findings` and
// exit 0 — a fabricated claimed_at was invisible to every check the suite had
// unless it happened to invert the one pair model.go already compared.
func TestVerifyFlagsCreatedAfterClaimed(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/REQ-800-ordering-wedge.md",
			"---\nid: REQ-800\ntitle: fixture\nstatus: completed\n" +
				"created_at: 2026-08-19T12:00:00Z\nclaimed_at: 2026-08-18T09:00:00Z\n" +
				"completed_at: 2026-08-19T14:00:00Z\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	orderingFindings := findingsMentioning(report, verifyCategoryTimestampOrdering)
	if len(orderingFindings) != 1 {
		t.Fatalf("got %d ordering findings, want 1:\n%s", len(orderingFindings), renderVerifyReport(report))
	}
	// Both field names and both raw values, so the reader can repair without reopening the file.
	for _, required := range []string{"REQ-800", "created_at", "claimed_at",
		"2026-08-19T12:00:00Z", "2026-08-18T09:00:00Z"} {
		if !strings.Contains(orderingFindings[0].Detail, required) {
			t.Errorf("ordering finding detail %q omits %q", orderingFindings[0].Detail, required)
		}
	}
	// An archived file's remedy must name the archive repair, not hand git archaeology.
	if !strings.Contains(orderingFindings[0].Remedy, "audit-archive-timestamps.sh") {
		t.Errorf("archive remedy %q does not name audit-archive-timestamps.sh", orderingFindings[0].Remedy)
	}
	// `do-work cleanup` cannot rewrite stamps, so the fixable count must not claim it can.
	if orderingFindings[0].Fixable {
		t.Error("ordering finding is marked Fixable, but cleanup does not repair stamps")
	}
	if report.ExitCode() != 1 {
		t.Errorf("ordering violation exit code = %d, want 1", report.ExitCode())
	}
}

// A REQ can be wrong in both pairs. Each is its own finding, because collapsing
// them would show the reader half the repair.
func TestVerifyReportsEachViolatedTimestampPairSeparately(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/REQ-801-both-pairs.md",
			"---\nid: REQ-801\ntitle: fixture\nstatus: completed\n" +
				"created_at: 2026-08-19T12:00:00Z\nclaimed_at: 2026-08-18T09:00:00Z\n" +
				"completed_at: 2026-08-17T08:00:00Z\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	orderingFindings := findingsMentioning(report, verifyCategoryTimestampOrdering)
	if len(orderingFindings) != 2 {
		t.Fatalf("got %d ordering findings, want 2 (one per violated pair):\n%s",
			len(orderingFindings), renderVerifyReport(report))
	}
}

// Queue and working files are repaired by the SessionStart hook, so their remedy
// must say so rather than sending the reader to the archive-only tool.
func TestVerifyRoutesQueueOrderingRemedyToTheSessionStartRepairer(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/queue/REQ-802-queue-ordering.md",
			"---\nid: REQ-802\ntitle: fixture\nstatus: pending\n" +
				"created_at: 2026-08-19T12:00:00Z\nclaimed_at: 2026-08-18T09:00:00Z\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	orderingFindings := findingsMentioning(report, verifyCategoryTimestampOrdering)
	if len(orderingFindings) != 1 {
		t.Fatalf("got %d ordering findings, want 1:\n%s", len(orderingFindings), renderVerifyReport(report))
	}
	if !strings.Contains(orderingFindings[0].Remedy, "repair-req-timestamps.sh") {
		t.Errorf("queue remedy %q does not name the SessionStart repairer", orderingFindings[0].Remedy)
	}
	if strings.Contains(orderingFindings[0].Remedy, "audit-archive-timestamps.sh") {
		t.Errorf("queue remedy %q sends the reader to the archive-only tool", orderingFindings[0].Remedy)
	}
}

// Equal stamps are legal and must not be reported: Step 2's claim and Step 3.6's
// estimate legitimately read the same instant, and an absent or unparseable stamp
// is other checks' territory — the same carve-out detectCompletionAnomaly holds.
// Without this the probe would fire on a large share of a real queue.
func TestVerifyAcceptsEqualAbsentAndUnparseableTimestamps(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/REQ-803-equal.md",
			"---\nid: REQ-803\ntitle: fixture\nstatus: completed\n" +
				"created_at: 2026-08-19T12:00:00Z\nclaimed_at: 2026-08-19T12:00:00Z\n" +
				"completed_at: 2026-08-19T12:00:00Z\n---\n"},
		{"do-work/archive/REQ-804-absent.md",
			"---\nid: REQ-804\ntitle: fixture\nstatus: completed\n" +
				"created_at: 2026-08-19T12:00:00Z\ncompleted_at: 2026-08-19T14:00:00Z\n---\n"},
		{"do-work/archive/REQ-805-unparseable.md",
			"---\nid: REQ-805\ntitle: fixture\nstatus: completed\n" +
				"created_at: 2026-08-19T12:00:00Z\nclaimed_at: not-a-timestamp\n" +
				"completed_at: 2026-08-19T14:00:00Z\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	if orderingFindings := findingsMentioning(report, verifyCategoryTimestampOrdering); len(orderingFindings) != 0 {
		t.Errorf("got %d ordering findings, want 0:\n%s", len(orderingFindings), renderVerifyReport(report))
	}
}

// The outer pair: with claimed_at absent, nothing spans created_at and completed_at,
// so an impossible ordering would pass every other comparison. Checked only in that
// case, so a REQ with a parseable claimed_at never gets a third finding for a defect
// the two inner pairs already report.
func TestVerifyFlagsCreatedAfterCompletedWhenClaimedIsAbsent(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/REQ-806-outer-pair.md",
			"---\nid: REQ-806\ntitle: fixture\nstatus: completed\n" +
				"created_at: 2026-08-19T12:00:00Z\ncompleted_at: 2026-08-18T09:00:00Z\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	orderingFindings := findingsMentioning(report, verifyCategoryTimestampOrdering)
	if len(orderingFindings) != 1 {
		t.Fatalf("got %d ordering findings, want 1:\n%s", len(orderingFindings), renderVerifyReport(report))
	}
	if !strings.Contains(orderingFindings[0].Detail, "created_at") ||
		!strings.Contains(orderingFindings[0].Detail, "completed_at") {
		t.Errorf("outer-pair detail %q does not name both fields", orderingFindings[0].Detail)
	}
}

// The guard: a fully-reversed REQ with a parseable claimed_at reports exactly the
// two inner pairs, never a third for the outer one.
func TestVerifyDoesNotDoubleReportTheOuterPair(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/REQ-807-fully-reversed.md",
			"---\nid: REQ-807\ntitle: fixture\nstatus: completed\n" +
				"created_at: 2026-08-19T12:00:00Z\nclaimed_at: 2026-08-18T09:00:00Z\n" +
				"completed_at: 2026-08-17T08:00:00Z\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	if orderingFindings := findingsMentioning(report, verifyCategoryTimestampOrdering); len(orderingFindings) != 2 {
		t.Errorf("got %d ordering findings, want exactly 2:\n%s",
			len(orderingFindings), renderVerifyReport(report))
	}
}

// The captured RED from REQ-281, the REQ-233 shape observed live in this repo: a
// 10-minute span logged as 70. Before this probe the fixture exited 0 — nothing in
// the suite read calibration-log.tsv except the estimator at recalibration time,
// which re-fits the scoring table from the rows without ever checking them.
// The second row pins the tolerance: it agrees within a minute and must stay silent.
func TestVerifyFlagsCalibrationRowDisagreeingWithFrontmatter(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/REQ-233-logged-wrong.md",
			"---\nid: REQ-233\ntitle: fixture\nstatus: completed\n" +
				"claimed_at: 2026-08-18T11:00:00Z\ncompleted_at: 2026-08-18T11:10:30Z\n---\n"},
		{"do-work/archive/REQ-234-logged-right.md",
			"---\nid: REQ-234\ntitle: fixture\nstatus: completed\n" +
				"claimed_at: 2026-08-18T11:00:00Z\ncompleted_at: 2026-08-18T11:30:40Z\n---\n"},
		{"do-work/calibration-log.tsv",
			"req_id\troute\testimated_p50_minutes\twall_minutes\tcompleted_at\n" +
				"REQ-233\tB\t25\t70\t2026-08-18T11:10:30Z\n" +
				"REQ-234\tB\t25\t30\t2026-08-18T11:30:40Z\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	mismatches := findingsMentioning(report, verifyCategoryCalibrationLogMismatch)
	if len(mismatches) != 1 {
		t.Fatalf("got %d calibration mismatches, want 1 (the agreeing row must stay silent):\n%s",
			len(mismatches), renderVerifyReport(report))
	}
	for _, required := range []string{"REQ-233", "70", "10"} {
		if !strings.Contains(mismatches[0].Detail, required) {
			t.Errorf("calibration mismatch detail %q omits %q", mismatches[0].Detail, required)
		}
	}
	// It must not pick a winner: the frontmatter can legitimately have been rewritten.
	if !strings.Contains(mismatches[0].Remedy, "either record may be the correct one") {
		t.Errorf("calibration remedy %q does not say either record may be correct", mismatches[0].Remedy)
	}
	if mismatches[0].Fixable {
		t.Error("calibration mismatch is marked Fixable, but no cleanup pass resolves it")
	}
	if report.ExitCode() != 1 {
		t.Errorf("calibration mismatch exit code = %d, want 1", report.ExitCode())
	}
}

// A row naming a REQ that exists nowhere, and a row whose REQ cannot yield a span,
// are their own findings — never disagreements. Reporting them as a mismatch would
// print a number next to a value that was never computed.
func TestVerifyReportsUnreconcilableCalibrationRowsSeparately(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/REQ-236-no-stamps.md",
			"---\nid: REQ-236\ntitle: fixture\nstatus: completed\n---\n"},
		{"do-work/calibration-log.tsv",
			"req_id\troute\testimated_p50_minutes\twall_minutes\tcompleted_at\n" +
				"REQ-235\tB\t25\t30\t2026-08-18T11:30:00Z\n" + // names no REQ in the tree
				"REQ-236\tB\t25\t30\t2026-08-18T11:30:00Z\n" + // present, but no parseable stamps
				"REQ-237\tB\t25\tnot-a-number\t2026-08-18T11:30:00Z\n" + // malformed wall_minutes
				"REQ-238\tB\n"}, // too few columns
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	if unreconcilable := findingsMentioning(report, verifyCategoryCalibrationRowUnreconcilable); len(unreconcilable) != 4 {
		t.Fatalf("got %d unreconcilable-row findings, want 4:\n%s",
			len(unreconcilable), renderVerifyReport(report))
	}
	if mismatches := findingsMentioning(report, verifyCategoryCalibrationLogMismatch); len(mismatches) != 0 {
		t.Errorf("got %d mismatches, want 0 — an unreconcilable row is never a disagreement:\n%s",
			len(mismatches), renderVerifyReport(report))
	}
}

// No log is normal in a repo that has archived nothing yet. It is not a finding —
// but it is also not a verified invariant, so it must be reported as a skipped probe
// rather than passing silently.
func TestVerifySkipsCalibrationProbeWhenTheLogIsAbsent(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	if findings := findingsMentioning(report, verifyCategoryCalibrationLogMismatch); len(findings) != 0 {
		t.Errorf("absent log produced %d mismatches, want 0", len(findings))
	}
	skippedMentionsCalibration := false
	for _, skipped := range report.SkippedProbes {
		if strings.Contains(skipped, "calibration-log probe") {
			skippedMentionsCalibration = true
		}
	}
	if !skippedMentionsCalibration {
		t.Errorf("absent log was not reported as a skipped probe:\n%s", renderVerifyReport(report))
	}
}

// The captured RED from REQ-284's first half. collectVerifyFindings takes a board
// the caller already built, and `now` stays a parameter — so the SAME board, with no
// file mtime changing between calls, must start reporting a stale claim once enough
// wall-clock time has passed. That is precisely the case an mtime-keyed cache cannot
// see, and the reason serve calls this outside its cache.
func TestCollectVerifyFindingsSeesAClaimGoStaleWithoutAnyFileChanging(t *testing.T) {
	claimedAt := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/working/REQ-810-in-flight.md",
			"---\nid: REQ-810\ntitle: fixture\nstatus: claimed\n" +
				"claimed_at: 2026-08-19T09:00:00Z\n---\n"},
	})

	// One board, built once, reused for both calls — no re-parse, no mtime change.
	freshMoment := claimedAt.Add(1 * time.Minute)
	board, buildError := buildBoard(repoRoot, freshMoment, defaultRecentWindow, lookupGitCommitDate)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}

	fresh := collectVerifyFindings(repoRoot, board, freshMoment)
	if claimFindings := findingsMentioning(fresh, verifyCategoryClaimNeedsAttention); len(claimFindings) != 0 {
		t.Fatalf("a one-minute-old claim already needs attention: %s", renderVerifyReport(fresh))
	}

	stale := collectVerifyFindings(repoRoot, board, claimedAt.Add(4*time.Hour))
	if claimFindings := findingsMentioning(stale, verifyCategoryClaimNeedsAttention); len(claimFindings) != 1 {
		t.Fatalf("got %d claim findings after advancing now by 4h, want 1 — the same board must age:\n%s",
			len(claimFindings), renderVerifyReport(stale))
	}
}

// runVerifyProbes keeps its signature and its behavior: it is now the wrapper that
// builds the board and delegates. The CLI contract (non-zero exit, same report)
// depends on this staying true.
func TestRunVerifyProbesStillReportsWhatCollectDoes(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/archive/REQ-811-reversed.md",
			"---\nid: REQ-811\ntitle: fixture\nstatus: completed\n" +
				"created_at: 2026-08-19T12:00:00Z\nclaimed_at: 2026-08-18T09:00:00Z\n" +
				"completed_at: 2026-08-19T14:00:00Z\n---\n"},
	})
	moment := time.Date(2026, 8, 19, 15, 0, 0, 0, time.UTC)

	wrapped, verifyError := runVerifyProbes(repoRoot, moment)
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	board, buildError := buildBoard(repoRoot, moment, defaultRecentWindow, lookupGitCommitDate)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	direct := collectVerifyFindings(repoRoot, board, moment)

	if len(wrapped.Findings) != len(direct.Findings) {
		t.Fatalf("wrapper reported %d findings, direct call %d — they must not diverge",
			len(wrapped.Findings), len(direct.Findings))
	}
	if wrapped.ExitCode() != 1 {
		t.Errorf("runVerifyProbes exit code = %d, want 1 — the CLI contract", wrapped.ExitCode())
	}
	if wrapped.RepoRoot != repoRoot {
		t.Errorf("wrapper RepoRoot = %q, want %q", wrapped.RepoRoot, repoRoot)
	}
}

// structuralDamageFixture is the user's reported shape, verbatim in structure:
// several REQ files in one queue carrying delimiter/field damage — including one
// whose opening frontmatter fence is broken so its status, title and user_request
// all parse empty — plus the three files that must stay silent (a healthy REQ, a
// stakeholder REQ that omits user_request by design, and a legacy archived REQ
// that predates the field).
//
// Before the structural probes, this exact tree printed `OK: no findings` and
// exited 0.
func structuralDamageFixture(t *testing.T) string {
	t.Helper()
	return writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		// Opening fence broken: `--` instead of `---`, so nothing below it is read.
		{"do-work/queue/REQ-901-broken-opening-fence.md",
			"--\nid: REQ-901\ntitle: \"Broken opening fence\"\nstatus: pending\nuser_request: UR-900\n---\n\n# Broken\n"},
		{"do-work/queue/REQ-902-empty-status.md",
			"---\nid: REQ-902\ntitle: \"Empty status\"\nstatus:\nuser_request: UR-900\n---\n\n# Empty status\n"},
		{"do-work/queue/REQ-903-unrecognized-status.md",
			"---\nid: REQ-903\ntitle: \"Typo status\"\nstatus: pnding\nuser_request: UR-900\n---\n\n# Typo\n"},
		{"do-work/queue/REQ-904-empty-id.md",
			"---\nid:\ntitle: \"Empty id\"\nstatus: pending\nuser_request: UR-900\n---\n\n# Empty id\n"},
		{"do-work/queue/REQ-905-missing-user-request.md",
			"---\nid: REQ-905\ntitle: \"No upward pointer\"\nstatus: pending\n---\n\n# No pointer\n"},
		{"do-work/queue/REQ-906-healthy.md",
			"---\nid: REQ-906\ntitle: \"Healthy\"\nstatus: pending\nuser_request: UR-900\n---\n\n# Healthy\n"},
		// Legitimate absence 1 — actions/work-reference.md → Stakeholder REQ Template
		// omits user_request deliberately, so UR membership cannot hold the source UR open.
		{"do-work/queue/REQ-907-stakeholder-questions-priya.md",
			"---\nid: REQ-907\ntitle: \"Stakeholder questions: Priya (design)\"\nstatus: blocked\n" +
				"stakeholder: \"Priya (design)\"\nblocked_by: \"answers from Priya (design)\"\n---\n\n# Stakeholder Questions\n"},
		// Legitimate absence 2 — source: code-review is a documented UR-less shape.
		{"do-work/archive/legacy/REQ-001-legacy.md",
			"---\nid: REQ-001\ntitle: \"Legacy archived work\"\nstatus: completed\nsource: code-review\ncompleted_at: 2026-05-01T10:00:00Z\n---\n\n# Legacy\n"},
		{"do-work/user-requests/UR-900/input.md", "---\nid: UR-900\ntitle: \"Fixture UR\"\n---\n\n# Fixture UR\n"},
	})
}

// findingsNaming returns every finding whose detail names a REQ id, across all
// categories — the question a carve-out assertion actually asks ("did anything at
// all fire on this file?"), which a per-category filter cannot answer.
func findingsNaming(report VerifyReport, requestId string) []VerifyFinding {
	var matched []VerifyFinding
	for _, finding := range report.Findings {
		if strings.Contains(finding.Detail, requestId) {
			matched = append(matched, finding)
		}
	}
	return matched
}

// The captured RED: verify exited 0 on a queue where structural damage had eaten
// the fields the pipeline routes on. Each damage shape must now produce a finding
// that names the broken field and carries a remedy, so the operator can act
// without opening this file.
func TestVerifyFlagsEachStructuralDamageShape(t *testing.T) {
	report, verifyError := runVerifyProbes(structuralDamageFixture(t), time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	if report.ExitCode() == 0 {
		t.Fatalf("structurally damaged queue still exits 0:\n%s", renderVerifyReport(report))
	}

	for _, damageCase := range []struct {
		requestId      string
		category       string
		detailFragment string
		remedyFragment string
	}{
		{"REQ-901", verifyCategoryStructurallyDamagedRequest, "no leading frontmatter fence", "opening `---`"},
		{"REQ-902", verifyCategoryUnrecognizedRequestStatus, "empty or absent status:", "Schema Read Contract"},
		{"REQ-903", verifyCategoryUnrecognizedRequestStatus, `unrecognized status: value "pnding"`, "Schema Read Contract"},
		{"REQ-904", verifyCategoryStructurallyDamagedRequest, "caution: its id was recovered from the filename", "id: REQ-904"},
		{"REQ-905", verifyCategoryStructurallyDamagedRequest, "carries no user_request:", "user_request: UR-NNN"},
	} {
		matched := findingsNaming(report, damageCase.requestId)
		if len(matched) != 1 {
			t.Errorf("%s produced %d findings, want exactly 1:\n%s",
				damageCase.requestId, len(matched), renderVerifyReport(report))
			continue
		}
		finding := matched[0]
		if finding.Category != damageCase.category {
			t.Errorf("%s category = %q, want %q", damageCase.requestId, finding.Category, damageCase.category)
		}
		if !strings.Contains(finding.Detail, damageCase.detailFragment) {
			t.Errorf("%s detail must name the broken field (%q), got: %s",
				damageCase.requestId, damageCase.detailFragment, finding.Detail)
		}
		if !strings.Contains(finding.Remedy, damageCase.remedyFragment) {
			t.Errorf("%s remedy must tell the operator what to write (%q), got: %s",
				damageCase.requestId, damageCase.remedyFragment, finding.Remedy)
		}
		if finding.Fixable {
			t.Errorf("%s must not advertise `do-work cleanup` as a mechanical fix: %+v",
				damageCase.requestId, finding)
		}
	}
}

// The broken-fence file's id, status and user_request are all empty BECAUSE the
// fence is gone, and the fence remedy repairs all three. Reporting each field
// separately would be four findings for one defect — the double-report the
// timestamp probe's outer-pair carve-out avoids. Pinned separately from the shape
// table above because "exactly one finding" is the property, not an incidental count.
func TestVerifyReportsABrokenFenceOnceNotOncePerEmptiedField(t *testing.T) {
	report, verifyError := runVerifyProbes(structuralDamageFixture(t), time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	fenceFindings := findingsNaming(report, "REQ-901")
	if len(fenceFindings) != 1 {
		t.Fatalf("broken fence produced %d findings, want 1 — its emptied fields must not each report:\n%s",
			len(fenceFindings), renderVerifyReport(report))
	}
	for _, unwantedFragment := range []string{"empty or absent id", "carries no user_request", "unrecognized status"} {
		if strings.Contains(fenceFindings[0].Detail, unwantedFragment) {
			t.Errorf("the fence finding restates a consequence (%q) instead of the cause: %s",
				unwantedFragment, fenceFindings[0].Detail)
		}
	}
}

// The exemptions keep the probe trustworthy. Each is an affirmative documented
// schema shape, rather than an archive location, so an ordinary missing field
// stays visible regardless of where its file lives.
func TestVerifyStaysSilentOnLegitimateAbsenceOfUserRequest(t *testing.T) {
	report, verifyError := runVerifyProbes(structuralDamageFixture(t), time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	for _, mustStaySilent := range []struct {
		requestId string
		why       string
	}{
		{"REQ-906", "a healthy REQ carries every field"},
		{"REQ-907", "a stakeholder REQ omits user_request by design (Stakeholder REQ Template)"},
		{"REQ-001", "source: code-review is a documented UR-less shape"},
	} {
		if matched := findingsNaming(report, mustStaySilent.requestId); len(matched) != 0 {
			t.Errorf("%s was flagged but must not be — %s; got: %+v",
				mustStaySilent.requestId, mustStaySilent.why, matched)
		}
	}
}

func TestVerifyStaysSilentOnEveryDocumentedUserRequestlessSchema(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/queue/REQ-908-source-code-review.md",
			"---\nid: REQ-908\ntitle: \"Code review\"\nstatus: pending\nsource: code-review\n---\n"},
		{"do-work/queue/REQ-909-scoped-generated-review.md",
			"---\nid: REQ-909\ntitle: \"Generated review\"\nstatus: pending\nreview_generated: true\nscope: queue-kanban\n---\n"},
		{"do-work/archive/REQ-910-context-reference.md",
			"---\nid: REQ-910\ntitle: \"Context reference\"\nstatus: completed\ncontext_ref: do-work/runs/review.md\ncompleted_at: 2026-05-01T10:00:00Z\n---\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	for _, requestId := range []string{"REQ-908", "REQ-909", "REQ-910"} {
		if matched := findingsNaming(report, requestId); len(matched) != 0 {
			t.Errorf("%s was flagged despite its documented UR-less schema: %+v", requestId, matched)
		}
	}
}

// The exemptions must key on their declared schema discriminators, not merely
// happen to pass because of status or location. Removing each discriminator must
// restore the user_request finding; an ordinary legacy-path file has no exemption.
func TestVerifyUserRequestExemptionsRequireTheirSchemaDiscriminator(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		// Same file as the exempt stakeholder REQ, minus the stakeholder: marker.
		{"do-work/queue/REQ-907-stakeholder-questions-priya.md",
			"---\nid: REQ-907\ntitle: \"Stakeholder questions: Priya (design)\"\nstatus: blocked\n" +
				"blocked_by: \"answers from Priya (design)\"\n---\n\n# Stakeholder Questions\n"},
		// Source discriminator removed.
		{"do-work/archive/REQ-001-code-review.md",
			"---\nid: REQ-001\ntitle: \"Code review\"\nstatus: completed\ncompleted_at: 2026-05-01T10:00:00Z\n---\n\n# Review\n"},
		// review_generated needs its paired nonempty scope discriminator.
		{"do-work/archive/REQ-002-generated-review.md",
			"---\nid: REQ-002\ntitle: \"Generated review\"\nstatus: completed\nreview_generated: true\ncompleted_at: 2026-05-01T10:00:00Z\n---\n\n# Review\n"},
		// context_ref is an affirmative legacy shape even at archive root; this
		// mutation removes it.
		{"do-work/archive/REQ-003-context-reference.md",
			"---\nid: REQ-003\ntitle: \"Context reference\"\nstatus: completed\ncompleted_at: 2026-05-01T10:00:00Z\n---\n\n# Context\n"},
		// Location alone must never exempt a missing pointer.
		{"do-work/archive/legacy/REQ-004-path-only.md",
			"---\nid: REQ-004\ntitle: \"Path-only legacy\"\nstatus: completed\ncompleted_at: 2026-05-01T10:00:00Z\n---\n\n# Legacy\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	for _, mustBeFlagged := range []string{"REQ-907", "REQ-001", "REQ-002", "REQ-003", "REQ-004"} {
		matched := findingsMentioning(report, verifyCategoryStructurallyDamagedRequest)
		found := false
		for _, finding := range matched {
			if strings.Contains(finding.Detail, mustBeFlagged) && strings.Contains(finding.Detail, "user_request") {
				found = true
			}
		}
		if !found {
			t.Errorf("%s lost its exemption's discriminator and was still not flagged — the carve-out is not keyed on it:\n%s",
				mustBeFlagged, renderVerifyReport(report))
		}
	}
}

// An intact opening fence with no closing fence reaches the parser as raw body
// text. Verify must retain that leniency, inspect only those retained bytes, and
// tell the operator which delimiter is actually missing.
func TestVerifyDistinguishesMissingClosingFrontmatterFence(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/queue/REQ-908-missing-closing-fence.md",
			"\ufeff---\r\nid: REQ-908\r\nstatus: pending\r\nuser_request: UR-900\r\n# still parsed as body\r\n"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	matched := findingsNaming(report, "REQ-908")
	if len(matched) != 1 {
		t.Fatalf("missing closing fence produced %d findings, want 1:\n%s", len(matched), renderVerifyReport(report))
	}
	if !strings.Contains(matched[0].Detail, "opening frontmatter fence but no closing fence") {
		t.Errorf("missing-closing detail = %q, want the missing closing fence named", matched[0].Detail)
	}
	if strings.Contains(matched[0].Detail, "no leading frontmatter fence") {
		t.Errorf("missing-closing detail describes the wrong delimiter: %q", matched[0].Detail)
	}
}

// A bare EOF delimiter is not an opening frontmatter fence under
// splitFrontmatter's contract: an opening fence must be its own newline-terminated
// first line. Keep verifier wording aligned with that parser boundary.
func TestVerifyDoesNotTreatABareEOFDelimiterAsAnOpeningFence(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/queue/REQ-909-bare-delimiter.md", "---"},
	})

	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	matched := findingsNaming(report, "REQ-909")
	if len(matched) != 1 {
		t.Fatalf("bare delimiter produced %d findings, want 1:\n%s", len(matched), renderVerifyReport(report))
	}
	if !strings.Contains(matched[0].Detail, "no leading frontmatter fence") {
		t.Errorf("bare delimiter detail = %q, want no-leading-fence classification", matched[0].Detail)
	}
}

// Detection must not cost the parser its leniency: the point is to REPORT the
// damage, not to start rejecting files. Every one of the eight fixture REQs — the
// broken-fence file included — must still parse and still reach the board.
func TestStructuralDamageStillParsesAndStillReachesTheBoard(t *testing.T) {
	repoRoot := structuralDamageFixture(t)
	board, buildError := buildBoard(repoRoot, time.Now(), defaultRecentWindow, lookupGitCommitDate)
	if buildError != nil {
		t.Fatalf("buildBoard on a damaged tree must not error: %v", buildError)
	}
	if len(board.AllRequests) != 8 {
		t.Fatalf("damaged tree yielded %d tickets, want all 8 — a REQ with one bad line must still appear on the board", len(board.AllRequests))
	}
	for _, requiredId := range []string{"REQ-901", "REQ-902", "REQ-903", "REQ-904", "REQ-905", "REQ-906", "REQ-907", "REQ-001"} {
		if _, present := board.RequestsById[requiredId]; !present {
			t.Errorf("%s did not reach the board — the parser stopped being lenient", requiredId)
		}
	}
}

// Both probes must forward the board's structured evidence rather than re-walking
// the tree or matching warning prose — the rule appendCompletionAnomalyFindings
// states for itself, pinned the way the stray-file probe pins it: an in-memory
// board whose repo root does not exist on disk, so a filesystem re-walk would have
// nothing to find, and an empty Warnings slice, so prose matching would have
// nothing to match.
func TestStructuralProbesUseStructuredEvidenceNotWarningProse(t *testing.T) {
	absentRepoRoot := filepath.Join(t.TempDir(), "missing-repo")
	if _, statError := os.Stat(absentRepoRoot); !os.IsNotExist(statError) {
		t.Fatalf("fixture repo root must remain absent so a re-walk has no evidence; stat error = %v", statError)
	}
	fencelessTicket := &RequestTicket{
		RequestId:           "REQ-911",
		FrontmatterMarkdown: "",
		FilePath:            filepath.Join(absentRepoRoot, "do-work", "queue", "REQ-911-fenceless.md"),
		TreeSection:         "queue",
	}
	typoStatusTicket := &RequestTicket{
		RequestId:           "REQ-912",
		FrontmatterMarkdown: "---\nid: REQ-912\nstatus: pnding\nuser_request: UR-900\n---\n",
		OriginalStatus:      "pnding",
		Status:              "pnding",
		StatusUnrecognized:  true,
		UserRequestId:       "UR-900",
		FilePath:            filepath.Join(absentRepoRoot, "do-work", "queue", "REQ-912-typo.md"),
		TreeSection:         "queue",
	}
	board := &Board{
		RepoRoot:    absentRepoRoot,
		AllRequests: []*RequestTicket{fencelessTicket, typoStatusTicket},
		Warnings:    nil,
	}

	var report VerifyReport
	appendStructuralDamageFindings(&report, board)
	appendUnrecognizedStatusFindings(&report, board)

	if len(findingsNaming(report, "REQ-911")) != 1 {
		t.Errorf("the fenceless ticket produced %d findings from structured evidence alone, want 1: %+v",
			len(findingsNaming(report, "REQ-911")), report.Findings)
	}
	statusFindings := findingsNaming(report, "REQ-912")
	if len(statusFindings) != 1 {
		t.Fatalf("the unrecognized-status ticket produced %d findings from structured evidence alone, want 1: %+v",
			len(statusFindings), report.Findings)
	}
	if statusFindings[0].Category != verifyCategoryUnrecognizedRequestStatus {
		t.Errorf("category = %q, want %q", statusFindings[0].Category, verifyCategoryUnrecognizedRequestStatus)
	}
}

// ungatedOverlapFixtureBoard builds a queue of pending REQs that all write one file,
// wired with whatever depends_on edges the case supplies.
func ungatedOverlapFixtureBoard(t *testing.T, dependsOnByRequestId map[string]string) *Board {
	t.Helper()
	fixtureFiles := []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
	}
	for _, requestId := range []string{"REQ-901", "REQ-902", "REQ-903"} {
		fixtureFiles = append(fixtureFiles, verifyFixtureFile{
			"do-work/queue/" + requestId + "-fixture.md",
			"---\nid: " + requestId + "\ntitle: fixture\nstatus: pending\n" +
				"depends_on: " + dependsOnByRequestId[requestId] + "\n" +
				"write_set:\n  - shared/contended.go\n---\n",
		})
	}
	repoRoot := writeVerifyFixture(t, fixtureFiles)
	board, buildError := buildBoard(repoRoot, time.Date(2026, 8, 27, 9, 0, 0, 0, time.UTC),
		defaultRecentWindow, lookupGitCommitDate)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	return board
}

// Two REQs that declare the same file and can be dispatched together are
// reported; the same two, once ordered, are not.
//
// The third case is the one that matters and the reason this probe walks the
// graph instead of checking adjacency: REQ-901 → REQ-902 → REQ-903 leaves 901
// and 903 perfectly serialized with no edge between them. A direct-edge test
// reports that correctly-ordered chain on every single run, and a probe that
// cries wolf on healthy state is a probe people learn to skip.
func TestVerifyReportsOnlyWriteSetOverlapsThatCanRunTogether(t *testing.T) {
	repoRoot := t.TempDir()
	moment := time.Date(2026, 8, 27, 9, 0, 0, 0, time.UTC)

	for _, testCase := range []struct {
		name           string
		dependsOn      map[string]string
		wantPairCount  int
		wantMentioning []string
	}{
		{
			name:      "no edges at all",
			dependsOn: map[string]string{"REQ-901": "[]", "REQ-902": "[]", "REQ-903": "[]"},
			// Every pair of the three collides: 901/902, 901/903, 902/903.
			wantPairCount:  3,
			wantMentioning: []string{"REQ-901", "REQ-902", "REQ-903", "shared/contended.go"},
		},
		{
			name:          "one direct edge orders one pair",
			dependsOn:     map[string]string{"REQ-901": "[]", "REQ-902": "[REQ-901]", "REQ-903": "[]"},
			wantPairCount: 2, // 901/903 and 902/903 remain
		},
		{
			name:          "a chain orders every pair, including the transitive one",
			dependsOn:     map[string]string{"REQ-901": "[]", "REQ-902": "[REQ-901]", "REQ-903": "[REQ-902]"},
			wantPairCount: 0,
		},
		{
			name:          "an edge pointing the other way orders the pair just as well",
			dependsOn:     map[string]string{"REQ-901": "[REQ-903]", "REQ-902": "[REQ-901]", "REQ-903": "[]"},
			wantPairCount: 0,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			board := ungatedOverlapFixtureBoard(t, testCase.dependsOn)
			report := collectVerifyFindings(repoRoot, board, moment)
			overlapFindings := findingsMentioning(report, verifyCategoryUngatedWriteSetOverlap)
			if len(overlapFindings) != testCase.wantPairCount {
				t.Errorf("reported %d ungated pairs, want %d:\n%s",
					len(overlapFindings), testCase.wantPairCount, renderVerifyReport(report))
			}
			for _, wantText := range testCase.wantMentioning {
				if !strings.Contains(renderVerifyReport(report), wantText) {
					t.Errorf("the report never names %q — a reader cannot act on it:\n%s",
						wantText, renderVerifyReport(report))
				}
			}
		})
	}
}

// A depends_on cycle must not hang verify. The graph is authored by hand, so a
// cycle is a thing a person can write, and a reachability walk with no visited
// set revisits the same ids forever.
//
// The FOURTH REQ is what makes this test work. With only the cycle, every query
// is between two members and the target is always found before the walk can
// loop — so dropping the visited set changes nothing and the test passes on a
// hanging implementation. REQ-904 sits outside the cycle, so proving it
// unreachable means walking the cycle to exhaustion, which is the only shape
// that actually hangs.
func TestVerifyWriteSetOverlapProbeTerminatesOnADependencyCycle(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/queue/REQ-901-cycle.md", cycleFixtureRequest("REQ-901", "[REQ-903]")},
		{"do-work/queue/REQ-902-cycle.md", cycleFixtureRequest("REQ-902", "[REQ-901]")},
		{"do-work/queue/REQ-903-cycle.md", cycleFixtureRequest("REQ-903", "[REQ-902]")},
		{"do-work/queue/REQ-904-outside.md", cycleFixtureRequest("REQ-904", "[]")},
	})
	moment := time.Date(2026, 8, 27, 9, 0, 0, 0, time.UTC)
	board, buildError := buildBoard(repoRoot, moment, defaultRecentWindow, lookupGitCommitDate)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}

	// Returning at all is half the assertion; a hang fails by test timeout.
	report := collectVerifyFindings(repoRoot, board, moment)
	overlapFindings := findingsMentioning(report, verifyCategoryUngatedWriteSetOverlap)
	// REQ-904 is ordered against nothing, so it collides with all three cycle
	// members. The cycle members reach each other, so they are not reported.
	if len(overlapFindings) != 3 {
		t.Errorf("reported %d ungated pairs, want the three against REQ-904:\n%s",
			len(overlapFindings), renderVerifyReport(report))
	}
}

func cycleFixtureRequest(requestId string, dependsOn string) string {
	return "---\nid: " + requestId + "\ntitle: fixture\nstatus: pending\n" +
		"depends_on: " + dependsOn + "\nwrite_set:\n  - shared/contended.go\n---\n"
}

// A write_set entry is a PATTERN. Two REQs collide when their patterns
// intersect, which does not require them to be equal — so the finding has to
// name both sides, or it reports a collision and points at nothing.
func TestVerifyWriteSetOverlapNamesBothSidesOfAGlobCollision(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/queue/REQ-911-glob.md",
			"---\nid: REQ-911\ntitle: fixture\nstatus: pending\ndepends_on: []\n" +
				"write_set:\n  - web/*.js\n---\n"},
		{"do-work/queue/REQ-912-literal.md",
			"---\nid: REQ-912\ntitle: fixture\nstatus: pending\ndepends_on: []\n" +
				"write_set:\n  - web/board-detail.js\n---\n"},
	})
	moment := time.Date(2026, 8, 27, 9, 0, 0, 0, time.UTC)
	board, buildError := buildBoard(repoRoot, moment, defaultRecentWindow, lookupGitCommitDate)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}

	report := collectVerifyFindings(repoRoot, board, moment)
	overlapFindings := findingsMentioning(report, verifyCategoryUngatedWriteSetOverlap)
	if len(overlapFindings) != 1 {
		t.Fatalf("reported %d ungated pairs, want the one glob collision:\n%s",
			len(overlapFindings), renderVerifyReport(report))
	}
	for _, wantText := range []string{"web/*.js", "web/board-detail.js"} {
		if !strings.Contains(overlapFindings[0].Detail, wantText) {
			t.Errorf("the finding does not name %q, so it points at no file a reader can change: %s",
				wantText, overlapFindings[0].Detail)
		}
	}
}

// The finding must name the path it actually covers.
//
// An auto-wave computes its ready set from depends_on, so an edge keeps two
// REQs out of one wave. A TARGETED run does not — an explicitly-named REQ
// enters the wave regardless of depends_on (actions/work-reference.md →
// Auto-wave, condition 2) — so `do-work run --fan-out REQ-901 REQ-902` can
// dispatch two REQs this probe calls ordered. An unqualified "--fan-out may
// dispatch them concurrently" overclaims in one direction and an unqualified
// silence overclaims in the other; naming the wave is what keeps the finding
// true. Pinned because scope is the claim here, not phrasing.
func TestVerifyWriteSetOverlapFindingNamesTheAutoWave(t *testing.T) {
	board := ungatedOverlapFixtureBoard(t, map[string]string{
		"REQ-901": "[]", "REQ-902": "[]", "REQ-903": "[REQ-902]",
	})
	report := collectVerifyFindings(t.TempDir(), board, time.Date(2026, 8, 27, 9, 0, 0, 0, time.UTC))
	overlapFindings := findingsMentioning(report, verifyCategoryUngatedWriteSetOverlap)
	if len(overlapFindings) == 0 {
		t.Fatalf("no ungated pair reported, so the wording below is unchecked:\n%s", renderVerifyReport(report))
	}
	for _, finding := range overlapFindings {
		if !strings.Contains(finding.Detail, "auto-wave") {
			t.Errorf("the finding claims concurrency without naming the auto-wave, which is the only path depends_on gates: %s",
				finding.Detail)
		}
	}
}
