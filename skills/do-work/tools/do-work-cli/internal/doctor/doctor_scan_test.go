package doctor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func writeDoctorFixture(t *testing.T, repositoryRoot, relativePath, contents string) {
	t.Helper()
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolutePath, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func doctorRequest(requestID, status, extra, body string) string {
	return "---\nid: " + requestID + "\nstatus: " + status + "\n" + extra + "---\n" + body + "\n"
}

func initDoctorGit(t *testing.T, repositoryRoot string) {
	t.Helper()
	for _, arguments := range [][]string{{"init", "-q"}, {"config", "user.email", "doctor@example.com"}, {"config", "user.name", "Doctor Test"}} {
		command := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...)
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", arguments, err, output)
		}
	}
}

func commitDoctorFixture(t *testing.T, repositoryRoot, message string) string {
	t.Helper()
	for _, arguments := range [][]string{{"add", "-A"}, {"commit", "-q", "-m", message}} {
		command := exec.Command("git", append([]string{"-C", repositoryRoot}, arguments...)...)
		command.Env = append(os.Environ(), "GIT_AUTHOR_DATE=2026-08-30T19:00:00Z", "GIT_COMMITTER_DATE=2026-08-30T19:00:00Z")
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", arguments, err, output)
		}
	}
	command := exec.Command("git", "-C", repositoryRoot, "rev-parse", "HEAD")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(output))
}

func TestScanRepositoryClassifiesCombinedCleanBlankTimestampAndCollisionFixture(t *testing.T) {
	repositoryRoot := t.TempDir()
	initDoctorGit(t, repositoryRoot)
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-001-clean.md", doctorRequest("REQ-001", "pending", "created_at: 2026-08-30T19:00:00Z\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-002-first.md", doctorRequest("REQ-003", "pending", "created_at: 2099-01-01T00:00:00Z\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-003-second.md", doctorRequest("REQ-004", "completed", "created_at: 2026-08-30T19:00:00Z\n", "## Implementation Summary\n**Files changed:**\n- `src/x.go`\n\n## Qualification"))
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-005-blank.md", doctorRequest("REQ-005", "completed", "created_at: 2026-08-30T19:00:00Z\n", "Recovered"))
	commitDoctorFixture(t, repositoryRoot, "fixture")
	if err := os.WriteFile(filepath.Join(repositoryRoot, "do-work/archive/REQ-005-blank.md"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	commitDoctorFixture(t, repositoryRoot, "[REQ-005] record commit hash abcdef1")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	result := ScanRepository(context.Background(), snapshot, ScanOptions{Now: time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC)})
	codes := map[string]bool{}
	for _, finding := range result.Findings {
		codes[finding.Code] = true
	}
	for _, code := range []string{"REQUEST-ID-COLLISION", "BLANKED-RECORD", "TIMESTAMP-FUTURE"} {
		if !codes[code] {
			t.Fatalf("missing %s in %#v", code, result.Findings)
		}
	}
}

func TestScanRepositoryTreatsIncompleteInspectionAsFinding(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-010-malformed.md", "---\nid: REQ-010\n")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	result := ScanRepository(context.Background(), snapshot, ScanOptions{Now: time.Now().UTC()})
	if result.Outcome == "success" || len(result.Findings) == 0 {
		t.Fatalf("incomplete inspection reported clean: %#v", result)
	}
}

func TestImplementationEvidenceRejectsAbsoluteAndEscapingPaths(t *testing.T) {
	body := "## Implementation Summary\n- `/etc/passwd` (new)\n- `../escape.md` (new)\n- `src/kept.go` (modified)\n"
	paths := implementationPaths(body)
	if len(paths) != 1 || paths[0] != "src/kept.go" {
		t.Fatalf("implementation paths = %v", paths)
	}
}

func TestImplementationPathsUseOnlyPathLedFileBulletsAndDeduplicate(t *testing.T) {
	body := "## Implementation Summary\n\n" +
		"**Files changed:**\n" +
		"- `src/doctor.go` (modified)\n" +
		"- `src/doctor.go` (modified; repeated in a generated handback)\n" +
		"- `src/doctor_test.go` (new)\n\n" +
		"Processed `(2,557) +` historical records.\n" +
		"**What was done:** `See` the typed result; `(modified) —` is prose, not a path.\n" +
		"- Explanation only: `not/a/file-list-bullet.go` must not become ownership evidence.\n\n" +
		"## Qualification\n"
	paths := implementationPaths(body)
	want := []string{"src/doctor.go", "src/doctor_test.go"}
	if len(paths) != len(want) || paths[0] != want[0] || paths[1] != want[1] {
		t.Fatalf("implementation paths = %v, want %v", paths, want)
	}
}

func TestCanonicalPredicatesRetainStuckHollowAndStaleEvidence(t *testing.T) {
	repositoryRoot := t.TempDir()
	workingExtra := "title: Resume migration\nroute: b\ncreated_at: 2026-08-20T00:00:00Z\n"
	writeDoctorFixture(t, repositoryRoot, "do-work/working/REQ-100-missing-claim.md", doctorRequest("REQ-100", "in-progress", workingExtra, "## Plan\nReady\n\n## Testing\nInterrupted"))
	writeDoctorFixture(t, repositoryRoot, "do-work/working/REQ-101-bad-claim.md", doctorRequest("REQ-101", "in-progress", workingExtra+"claimed_at: yesterday\n", "## Exploration\nObserved"))
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-102-heading-only.md", doctorRequest("REQ-102", "completed", "created_at: 2026-08-20T00:00:00Z\n", "## Implementation Summary"))
	writeDoctorFixture(t, repositoryRoot, "do-work/archive/REQ-103-do-work-only.md", doctorRequest("REQ-103", "completed", "created_at: 2026-08-20T00:00:00Z\n", "## Implementation Summary\n- `do-work/archive/REQ-103-do-work-only.md` (modified)"))
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-104-pending-info.md", doctorRequest("REQ-104", "pending-answers", "created_at: 2026-08-26T00:00:00Z\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-105-pending-warning.md", doctorRequest("REQ-105", "pending-answers", "created_at: 2026-08-22T00:00:00Z\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-106-blocked-info.md", doctorRequest("REQ-106", "blocked", "created_at: 2026-08-01T00:00:00Z\nblocked_at: 2026-08-20T00:00:00Z\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-107-blocked-warning.md", doctorRequest("REQ-107", "blocked", "created_at: 2026-08-01T00:00:00Z\nblocked_at: 2026-08-15T00:00:00Z\n", "Body"))
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	result := ScanRepository(context.Background(), snapshot, ScanOptions{Now: time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)})
	byIDAndCode := map[string]repositoryFinding{}
	for _, finding := range result.Findings {
		for _, id := range finding.AffectedIDs {
			byIDAndCode[id+"/"+finding.Code] = repositoryFinding{severity: string(finding.Severity), evidence: strings.Join(finding.Evidence, " ")}
		}
	}
	for _, id := range []string{"REQ-100", "REQ-101"} {
		finding, found := byIDAndCode[id+"/STUCK-WORK"]
		if !found || !strings.Contains(finding.evidence, "title=Resume migration") || !strings.Contains(finding.evidence, "route=B") || !strings.Contains(finding.evidence, "last_phase=") {
			t.Fatalf("%s stuck evidence = %#v", id, finding)
		}
	}
	for _, id := range []string{"REQ-102", "REQ-103"} {
		if _, found := byIDAndCode[id+"/HOLLOW-COMPLETION"]; !found {
			t.Fatalf("missing hollow completion for %s: %#v", id, result.Findings)
		}
	}
	wantSeverity := map[string]string{"REQ-104": "info", "REQ-105": "warning", "REQ-106": "info", "REQ-107": "warning"}
	for id, severity := range wantSeverity {
		code := "STALE-PENDING-ANSWERS"
		if strings.HasPrefix(id, "REQ-106") || strings.HasPrefix(id, "REQ-107") {
			code = "STALE-BLOCKED"
		}
		if got := byIDAndCode[id+"/"+code].severity; got != severity {
			t.Fatalf("%s severity = %q, want %q", id, got, severity)
		}
	}
}

func TestStuckWorkSkipsTerminalAndRecentlyModifiedClaims(t *testing.T) {
	repositoryRoot := t.TempDir()
	now := time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC)
	terminalPath := "do-work/working/REQ-130-terminal.md"
	recentPath := "do-work/working/REQ-131-recent.md"
	inactivePath := "do-work/working/REQ-132-inactive.md"
	writeDoctorFixture(t, repositoryRoot, terminalPath, doctorRequest("REQ-130", "completed", "created_at: 2026-08-20T00:00:00Z\nclaimed_at: 2026-08-20T01:00:00Z\n", "## Qualification\nDone"))
	writeDoctorFixture(t, repositoryRoot, recentPath, doctorRequest("REQ-131", "claimed", "created_at: 2026-08-20T00:00:00Z\nclaimed_at: 2026-08-20T01:00:00Z\n", "## Testing\nStill running"))
	writeDoctorFixture(t, repositoryRoot, inactivePath, doctorRequest("REQ-132", "claimed", "created_at: 2026-08-20T00:00:00Z\nclaimed_at: 2026-08-20T01:00:00Z\n", "## Testing\nInterrupted"))
	for relativePath, modificationTime := range map[string]time.Time{
		terminalPath: now.Add(-26 * time.Hour),
		recentPath:   now.Add(-30 * time.Minute),
		inactivePath: now.Add(-26 * time.Hour),
	} {
		if err := os.Chtimes(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)), modificationTime, modificationTime); err != nil {
			t.Fatal(err)
		}
	}

	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	result := ScanRepository(context.Background(), snapshot, ScanOptions{Now: now})
	byIDAndCode := map[string]repositoryFinding{}
	for _, finding := range result.Findings {
		for _, id := range finding.AffectedIDs {
			byIDAndCode[id+"/"+finding.Code] = repositoryFinding{severity: string(finding.Severity), evidence: strings.Join(finding.Evidence, " ")}
		}
	}
	for _, id := range []string{"REQ-130", "REQ-131"} {
		if _, found := byIDAndCode[id+"/STUCK-WORK"]; found {
			t.Fatalf("%s unexpectedly reported as stuck: %#v", id, result.Findings)
		}
	}
	if _, found := byIDAndCode["REQ-130/STRANDED-TERMINAL-REQUEST"]; !found {
		t.Fatalf("terminal working request lost its location finding: %#v", result.Findings)
	}
	if finding, found := byIDAndCode["REQ-132/STUCK-WORK"]; !found || finding.severity != "error" {
		t.Fatalf("inactive request stuck finding = %#v", finding)
	}
}

type repositoryFinding struct {
	severity string
	evidence string
}

func TestManualFindingsKeepPathsAvoidDuplicateStatusAndAdvanceInspection(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-120-invalid.md", doctorRequest("REQ-120", "pnding", "created_at: 2026-08-30T00:00:00Z\n", "Body"))
	writeDoctorFixture(t, repositoryRoot, "do-work/queue/REQ-121-malformed.md", "---\nid: REQ-121\n")
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	result := ScanRepository(context.Background(), snapshot, ScanOptions{Now: time.Date(2026, 8, 30, 1, 0, 0, 0, time.UTC)})
	invalidStatusCount := 0
	for _, finding := range result.Findings {
		if finding.Code == "INVALID-STATUS" {
			invalidStatusCount++
		}
		if finding.Code == "DOCTOR-INSPECTION-WARNING" && strings.Contains(strings.Join(finding.Evidence, " "), "status") {
			t.Fatalf("status warning duplicated generic finding: %#v", finding)
		}
		if string(finding.Fixability) == "manual" {
			if len(finding.AffectedPaths) == 0 {
				t.Fatalf("manual finding lost exact path: %#v", finding)
			}
			if strings.Join(finding.NextArgv, "\x00") == strings.Join(doctorArgv(), "\x00") {
				t.Fatalf("manual finding reruns identical doctor: %#v", finding)
			}
		}
	}
	if invalidStatusCount != 1 {
		t.Fatalf("invalid status findings = %d: %#v", invalidStatusCount, result.Findings)
	}
}

func TestScanRepositoryReportsFinalizationTailsWithoutCommittedBlankCommitFalsePositive(t *testing.T) {
	repositoryRoot := t.TempDir()
	initDoctorGit(t, repositoryRoot)
	requestBody := "## Implementation Summary\n**Files changed:**\n- `src/x.go` (modified)\n\n## Qualification\nVerified"
	cleanPath := "do-work/archive/REQ-700-clean.md"
	modifiedPath := "do-work/archive/REQ-701-modified.md"
	writeDoctorFixture(t, repositoryRoot, cleanPath, doctorRequest("REQ-700", "completed", "created_at: 2026-08-30T19:00:00Z\ncommit:\n", requestBody))
	writeDoctorFixture(t, repositoryRoot, modifiedPath, doctorRequest("REQ-701", "completed-with-issues", "created_at: 2026-08-30T19:00:00Z\ncommit:\n", requestBody))
	commitDoctorFixture(t, repositoryRoot, "archive blank-commit controls")
	modifiedAbsolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(modifiedPath))
	modifiedBytes, err := os.ReadFile(modifiedAbsolutePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(modifiedAbsolutePath, append(modifiedBytes, []byte("\npost-archive tail\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	untrackedPath := "do-work/archive/REQ-702-untracked.md"
	writeDoctorFixture(t, repositoryRoot, untrackedPath, doctorRequest("REQ-702", "completed", "created_at: 2026-08-30T19:00:00Z\ncommit:\n", requestBody))

	gitPath := exec.Command("git", "-C", repositoryRoot, "rev-parse", "--git-path", "do-work-finalization")
	journalRootBytes, err := gitPath.Output()
	if err != nil {
		t.Fatal(err)
	}
	journalRoot := strings.TrimSpace(string(journalRootBytes))
	if !filepath.IsAbs(journalRoot) {
		journalRoot = filepath.Join(repositoryRoot, journalRoot)
	}
	if err := os.MkdirAll(journalRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	journalPath := filepath.Join(journalRoot, "REQ-703.json")
	journalBytes := []byte("{\"phase\":\"release_applied\",\"manifest\":{\"request_id\":\"REQ-703\"}}\n")
	if err := os.WriteFile(journalPath, journalBytes, 0o600); err != nil {
		t.Fatal(err)
	}

	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	result := ScanRepository(context.Background(), snapshot, ScanOptions{Now: time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC)})
	findings := map[string]resultmodel.CommandFinding{}
	for _, finding := range result.Findings {
		for _, requestID := range finding.AffectedIDs {
			findings[requestID+"/"+finding.Code] = finding
		}
	}
	if _, found := findings["REQ-700/ARCHIVED-WITHOUT-COMMIT"]; found {
		t.Fatalf("committed blank commit produced false positive: %#v", result.Findings)
	}
	for _, requestID := range []string{"REQ-701", "REQ-702"} {
		finding, found := findings[requestID+"/ARCHIVED-WITHOUT-COMMIT"]
		if !found || finding.Fixability != resultmodel.FixabilityRefused {
			t.Fatalf("%s archived tail finding = %#v", requestID, finding)
		}
		if !strings.Contains(strings.Join(finding.Evidence, " "), requestID) || !strings.Contains(strings.Join(finding.Evidence, " "), "do-work run-with-recovery") {
			t.Fatalf("%s archived tail lost recovery evidence: %#v", requestID, finding)
		}
	}
	journalFinding, found := findings["REQ-703/UNFINISHED-FINALIZATION"]
	if !found || journalFinding.Fixability != resultmodel.FixabilityRefused {
		t.Fatalf("journal tail finding = %#v", journalFinding)
	}
	if evidence := strings.Join(journalFinding.Evidence, " "); !strings.Contains(evidence, "REQ-703") || !strings.Contains(evidence, "release_applied") {
		t.Fatalf("journal tail lost request or phase evidence: %#v", journalFinding)
	}
	wantRecoveryArgv := []string{"do-work-cli", "--format", "json", "recover-finalization", "--discover"}
	if !reflect.DeepEqual(journalFinding.NextArgv, wantRecoveryArgv) || !reflect.DeepEqual(findings["REQ-701/ARCHIVED-WITHOUT-COMMIT"].NextArgv, wantRecoveryArgv) {
		t.Fatalf("tail recovery argv drifted: journal=%v archive=%v", journalFinding.NextArgv, findings["REQ-701/ARCHIVED-WITHOUT-COMMIT"].NextArgv)
	}
	journalAfter, err := os.ReadFile(journalPath)
	if err != nil || !reflect.DeepEqual(journalAfter, journalBytes) {
		t.Fatalf("doctor changed journal bytes: err=%v after=%q", err, journalAfter)
	}
}

func TestFinalizationTailFindingsPreserveMalformedAndUnreadableJournalEvidence(t *testing.T) {
	repositoryRoot := t.TempDir()
	initDoctorGit(t, repositoryRoot)
	if err := os.MkdirAll(filepath.Join(repositoryRoot, "do-work", "queue"), 0o755); err != nil {
		t.Fatal(err)
	}
	gitPath := exec.Command("git", "-C", repositoryRoot, "rev-parse", "--git-path", "do-work-finalization")
	journalRootBytes, err := gitPath.Output()
	if err != nil {
		t.Fatal(err)
	}
	journalRoot := strings.TrimSpace(string(journalRootBytes))
	if !filepath.IsAbs(journalRoot) {
		journalRoot = filepath.Join(repositoryRoot, journalRoot)
	}
	if err := os.MkdirAll(journalRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	malformedPath := filepath.Join(journalRoot, "REQ-720.json")
	malformedBytes := []byte("{malformed\n")
	if err := os.WriteFile(malformedPath, malformedBytes, 0o644); err != nil {
		t.Fatal(err)
	}
	unreadablePath := filepath.Join(journalRoot, "REQ-721.json")
	if err := os.Symlink("missing-journal", unreadablePath); err != nil {
		t.Fatal(err)
	}

	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	findings := FinalizationTailFindings(context.Background(), snapshot)
	byRequest := map[string]resultmodel.CommandFinding{}
	for _, finding := range findings {
		for _, requestID := range finding.AffectedIDs {
			byRequest[requestID] = finding
		}
	}
	for _, requestID := range []string{"REQ-720", "REQ-721"} {
		finding, found := byRequest[requestID]
		if !found || finding.Code != "UNFINISHED-FINALIZATION" || finding.Fixability != resultmodel.FixabilityRefused {
			t.Fatalf("%s missing typed unfinished-journal evidence: %#v", requestID, finding)
		}
		if evidence := strings.Join(finding.Evidence, " "); !strings.Contains(evidence, requestID) || !strings.Contains(evidence, "inspection failed") {
			t.Fatalf("%s missing inspection failure detail: %#v", requestID, finding)
		}
		if !reflect.DeepEqual(finding.NextArgv, finalizationRecoveryArgv()) {
			t.Fatalf("%s recovery argv=%v", requestID, finding.NextArgv)
		}
	}
	malformedAfter, err := os.ReadFile(malformedPath)
	if err != nil || !reflect.DeepEqual(malformedAfter, malformedBytes) {
		t.Fatalf("malformed journal changed: err=%v after=%q", err, malformedAfter)
	}
	linkAfter, err := os.Readlink(unreadablePath)
	if err != nil || linkAfter != "missing-journal" {
		t.Fatalf("unreadable journal object changed: err=%v target=%q", err, linkAfter)
	}
}

func TestFinalizationTailFindingsUseOneArchiveGitInventory(t *testing.T) {
	repositoryRoot := t.TempDir()
	initDoctorGit(t, repositoryRoot)
	requestBody := "## Implementation Summary\n- `src/x.go` (modified)\n\n## Qualification\nVerified"
	for requestNumber := 800; requestNumber < 1000; requestNumber++ {
		requestID := fmt.Sprintf("REQ-%d", requestNumber)
		writeDoctorFixture(t, repositoryRoot, "do-work/archive/"+requestID+"-clean.md",
			doctorRequest(requestID, "completed", "created_at: 2026-08-30T19:00:00Z\ncommit:\n", requestBody))
	}
	commitDoctorFixture(t, repositoryRoot, "committed blank-provenance scale fixture")

	realGit, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	probeRoot := t.TempDir()
	probeLog := filepath.Join(probeRoot, "git-calls.log")
	wrapper := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$DO_WORK_GIT_PROBE_LOG\"\nexec \"$DO_WORK_REAL_GIT\" \"$@\"\n"
	if err := os.WriteFile(filepath.Join(probeRoot, "git"), []byte(wrapper), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DO_WORK_REAL_GIT", realGit)
	t.Setenv("DO_WORK_GIT_PROBE_LOG", probeLog)
	t.Setenv("PATH", probeRoot+string(os.PathListSeparator)+os.Getenv("PATH"))

	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if findings := FinalizationTailFindings(context.Background(), snapshot); len(findings) != 0 {
		t.Fatalf("committed scale fixture produced findings: %#v", findings)
	}
	probeBytes, err := os.ReadFile(probeLog)
	if err != nil {
		t.Fatal(err)
	}
	archiveStatusProbes := 0
	for _, line := range strings.Split(string(probeBytes), "\n") {
		if strings.Contains(line, " status ") && strings.Contains(line, "do-work/archive") {
			archiveStatusProbes++
		}
	}
	if archiveStatusProbes != 1 {
		t.Fatalf("archive Git inventory probes=%d, want 1; calls:\n%s", archiveStatusProbes, probeBytes)
	}
}
