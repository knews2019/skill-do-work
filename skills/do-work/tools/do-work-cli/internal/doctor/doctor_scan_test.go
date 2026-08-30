package doctor

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
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
