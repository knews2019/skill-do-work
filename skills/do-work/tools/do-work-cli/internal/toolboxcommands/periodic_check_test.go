package toolboxcommands

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func periodicFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, args := range [][]string{{"init", "-q"}, {"-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "--allow-empty", "-qm", "initial"}} {
		if output, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git: %v %s", err, output)
		}
	}
	checklist := periodicChecklist{CompatibilityVersion: "security-v1", CheckRecords: []periodicCheckDefinition{
		{StableID: "dependencies", CompatibilityVersion: "v1", CadenceSeconds: 604800, QuickProfile: true, CheckCommand: "scanner --read-only", ResultInterpretation: "Findings are not execution errors"},
		{StableID: "images", CompatibilityVersion: "v1", CadenceSeconds: 604800, CheckCommand: "image-scanner --read-only", ResultInterpretation: "Missing credentials means incomplete"},
	}}
	data, _ := json.Marshal(checklist)
	if err := os.WriteFile(filepath.Join(root, "prime-periodic-check.md"), []byte("# Checks\n\n```periodic-check\n"+string(data)+"\n```\n"), 0644); err != nil {
		t.Fatal(err)
	}
	return root
}
func periodicInvoke(t *testing.T, root string, args ...string) string {
	t.Helper()
	output, err := periodicCommand(root, args)
	if err != nil {
		t.Fatal(err)
	}
	return output
}
func periodicStart(t *testing.T, root string, args ...string) string {
	t.Helper()
	output := periodicInvoke(t, root, append([]string{"start"}, args...)...)
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "Run: ") {
			return strings.TrimPrefix(line, "Run: ")
		}
	}
	t.Fatal("no run path")
	return ""
}
func periodicInput(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "input.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	return path
}
func periodicRecord(t *testing.T, root, run, id string, complete bool, findings ...periodicFinding) {
	t.Helper()
	periodicInvoke(t, root, "record", "--run", run, "--input", periodicInput(t, periodicCheckResult{CheckID: id, CheckCompleted: complete, EvidenceText: "actual scanner evidence", CheckFindings: findings}))
}
func periodicRead(t *testing.T, root, run string) periodicRun {
	t.Helper()
	record, err := periodicReadRun(root, run)
	if err != nil {
		t.Fatal(err)
	}
	return record
}
func periodicIssue() periodicFinding {
	return periodicFinding{StableID: "advisory/package", FindingSummary: "Package vulnerable", EvidenceText: "scanner output"}
}
func periodicFinish(t *testing.T, root, run string) periodicRun {
	t.Helper()
	periodicInvoke(t, root, "finish", "--run", run)
	return periodicRead(t, root, run)
}

func TestPeriodicHistorySeparatesProfilesAndPreservesCleanAfterInterruption(t *testing.T) {
	root := periodicFixture(t)
	before, _ := auditGit(root, "status", "--porcelain=v1", "--untracked-files=all")
	output := periodicInvoke(t, root, "status")
	if !strings.Contains(output, "no recorded run") {
		t.Fatal(output)
	}
	after, _ := auditGit(root, "status", "--porcelain=v1", "--untracked-files=all")
	if before != after {
		t.Fatal("status wrote files")
	}
	full := periodicStart(t, root)
	periodicRecord(t, root, full, "dependencies", true)
	periodicRecord(t, root, full, "images", true)
	if run := periodicFinish(t, root, full); run.RunOutcome != "clean" {
		t.Fatal(run)
	}
	quick := periodicStart(t, root, "--profile", "quick")
	periodicRecord(t, root, quick, "dependencies", true, periodicIssue())
	if run := periodicFinish(t, root, quick); run.RunOutcome != "completed-with-findings" {
		t.Fatal(run)
	}
	interrupted := periodicStart(t, root)
	periodicRecord(t, root, interrupted, "dependencies", false)
	output = periodicInvoke(t, root, "status")
	if !strings.Contains(output, "last completed:") || !strings.Contains(output, full) || !strings.Contains(output, interrupted) {
		t.Fatal(output)
	}
	run := periodicFinish(t, root, interrupted)
	if run.RunOutcome != "incomplete" || len(run.MissingCoverage) != 2 {
		t.Fatal(run)
	}
	if _, err := periodicCommand(root, []string{"finish", "--run", full}); err == nil {
		t.Fatal("finished run was editable")
	}
}

func TestPeriodicFindingsResolveOnlyWithCoverageAndLinkRepairs(t *testing.T) {
	root := periodicFixture(t)
	requestDir := filepath.Join(root, "do-work", "queue")
	if err := os.MkdirAll(requestDir, 0755); err != nil {
		t.Fatal(err)
	}
	requestPath := filepath.Join(requestDir, "REQ-001-package-repair.md")
	if err := os.WriteFile(requestPath, []byte("---\nid: REQ-001\nstatus: pending\ntitle: Repair package\n---\n"), 0644); err != nil {
		t.Fatal(err)
	}
	first := periodicStart(t, root, "--profile", "quick", "--capture")
	periodicRecord(t, root, first, "dependencies", true, periodicIssue())
	periodicInvoke(t, root, "link", "--run", first, "--input", periodicInput(t, map[string]string{"check_id": "dependencies", "finding_id": "advisory/package", "request_id": "REQ-001", "disposition": "created", "note": "Captured with canonical transaction"}))
	periodicFinish(t, root, first)
	for _, status := range []string{"pending", "in-progress", "completed"} {
		destination := requestPath
		if status == "in-progress" {
			destination = filepath.Join(root, "do-work", "working", "REQ-001-package-repair.md")
		}
		if status == "completed" {
			destination = filepath.Join(root, "do-work", "archive", "UR-001", "REQ-001-package-repair.md")
		}
		if destination != requestPath {
			if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.Rename(requestPath, destination); err != nil {
				t.Fatal(err)
			}
			requestPath = destination
		}
		if err := os.WriteFile(requestPath, []byte("---\nid: REQ-001\nstatus: "+status+"\ntitle: Repair package\n---\n"), 0644); err != nil {
			t.Fatal(err)
		}
		next := periodicStart(t, root, "--profile", "quick", "--capture")
		periodicRecord(t, root, next, "dependencies", true, periodicIssue())
		finding := periodicRead(t, root, next).CheckRecords[0].CheckFindings[0]
		want := "link"
		if status == "completed" {
			want = "verify-resolution"
		}
		if finding.FindingDelta != "unchanged" || finding.CaptureAction != want || finding.RepairRequest != "REQ-001" {
			t.Fatal(finding)
		}
		periodicFinish(t, root, next)
	}
	skipped := periodicStart(t, root, "--profile", "quick")
	periodicRecord(t, root, skipped, "dependencies", false)
	periodicFinish(t, root, skipped)
	if !strings.Contains(periodicInvoke(t, root, "status"), "unresolved dependencies/advisory/package") {
		t.Fatal("skipped check resolved issue")
	}
	clean := periodicStart(t, root, "--profile", "quick")
	periodicRecord(t, root, clean, "dependencies", true)
	if finding := periodicRead(t, root, clean).CheckRecords[0].CheckFindings[0]; finding.FindingDelta != "resolved" {
		t.Fatal(finding)
	}
	periodicFinish(t, root, clean)
	recurrence := periodicStart(t, root, "--profile", "quick", "--capture")
	periodicRecord(t, root, recurrence, "dependencies", true, periodicIssue())
	finding := periodicRead(t, root, recurrence).CheckRecords[0].CheckFindings[0]
	if finding.CaptureAction != "recurrence" || finding.PreviousRequest != "REQ-001" {
		t.Fatal(finding)
	}
}

func TestPeriodicConcurrentAttemptsAndSharedCheckout(t *testing.T) {
	root := periodicFixture(t)
	var wait sync.WaitGroup
	errors := make(chan error, 4)
	for index := 0; index < 4; index++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			_, err := periodicCommand(root, []string{"start", "--profile", "quick"})
			errors <- err
		}()
	}
	wait.Wait()
	close(errors)
	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}
	history, err := periodicReadHistory(root)
	if err != nil || len(history) != 4 {
		t.Fatalf("history %d: %v", len(history), err)
	}
	// Share actual records through Git, not a machine-specific memory location.
	for _, args := range [][]string{{"add", "prime-periodic-check.md", "do-work/deliverables"}, {"-c", "user.name=Test", "-c", "user.email=test@example.invalid", "commit", "-qm", "share history"}} {
		if output, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("%v %s", err, output)
		}
	}
	second := filepath.Join(t.TempDir(), "second")
	if output, err := exec.Command("git", "clone", "-q", root, second).CombinedOutput(); err != nil {
		t.Fatalf("%v %s", err, output)
	}
	output := periodicInvoke(t, second, "status")
	if !strings.Contains(output, history[0].RecordPath) {
		t.Fatal(output)
	}
}

func TestPeriodicMissingConfigurationAndUnattendedEquivalence(t *testing.T) {
	root := periodicFixture(t)
	manual := periodicStart(t, root, "--profile", "quick")
	hosted := periodicStart(t, root, "--profile", "quick", "--unattended")
	for _, run := range []string{manual, hosted} {
		periodicRecord(t, root, run, "dependencies", false)
		if periodicFinish(t, root, run).RunOutcome != "incomplete" {
			t.Fatal("missing prerequisite succeeded")
		}
	}
	if err := os.Remove(filepath.Join(root, "prime-periodic-check.md")); err != nil {
		t.Fatal(err)
	}
	run := periodicRead(t, root, periodicStart(t, root, "--unattended"))
	if run.RunOutcome != "incomplete" || run.ConfigurationError == "" || run.FinishedAt == nil {
		t.Fatal(run)
	}
}

func TestPeriodicChecklistChangeInvalidatesCoverageAndCaptureOverride(t *testing.T) {
	root := periodicFixture(t)
	path := filepath.Join(root, "prime-periodic-check.md")
	data, _ := os.ReadFile(path)
	data = []byte(strings.Replace(string(data), `"capture":false`, `"capture":true`, 1))
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	run := periodicStart(t, root, "--profile", "quick", "--capture", "--report-only")
	if periodicRead(t, root, run).CaptureEnabled {
		t.Fatal("report-only did not override")
	}
	periodicRecord(t, root, run, "dependencies", true)
	periodicFinish(t, root, run)
	data = []byte(strings.Replace(string(data), "security-v1", "security-v2", 1))
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	if output := periodicInvoke(t, root, "status"); !strings.Contains(output, "historical checklist") || !strings.Contains(output, "due: no equivalent completed check") {
		t.Fatal(output)
	}
}

func TestPeriodicConcurrentWriterAndLinkedHistoryRefuse(t *testing.T) {
	root := periodicFixture(t)
	run := periodicStart(t, root)
	if err := os.Mkdir(filepath.Join(root, run, "writer-lock"), 0700); err != nil {
		t.Fatal(err)
	}
	if _, err := periodicCommand(root, []string{"finish", "--run", run}); err == nil {
		t.Fatal("competing writer accepted")
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, periodicHistoryDirectory, "linked-history")); err != nil {
		t.Fatal(err)
	}
	if _, err := periodicCommand(root, []string{"status"}); err == nil {
		t.Fatal("linked history accepted")
	}
}

func TestPeriodicAgeUsesAbsoluteUTC(t *testing.T) {
	stamp := time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)
	if got := periodicTimestamp(stamp, stamp.Add(time.Hour)); got != "2026-10-29T12:00:00Z (1h0m0s ago)" {
		t.Fatal(got)
	}
}

func TestPeriodicRecurringRepairBecomesTheNextDuplicateTarget(t *testing.T) {
	root := periodicFixture(t)
	for _, id := range []string{"REQ-001", "REQ-002"} {
		path := filepath.Join(root, "do-work", "queue", id+"-repair.md")
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("---\nid: "+id+"\nstatus: pending\ntitle: Repair\n---\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	for _, id := range []string{"REQ-001", "REQ-002"} {
		run := periodicStart(t, root, "--profile", "quick", "--capture")
		periodicRecord(t, root, run, "dependencies", true, periodicIssue())
		periodicInvoke(t, root, "link", "--run", run, "--input", periodicInput(t, map[string]string{"check_id": "dependencies", "finding_id": "advisory/package", "request_id": id, "disposition": "created", "note": "canonical capture confirmed"}))
		periodicFinish(t, root, run)
		if id == "REQ-001" {
			clean := periodicStart(t, root, "--profile", "quick")
			periodicRecord(t, root, clean, "dependencies", true)
			periodicFinish(t, root, clean)
		}
	}
	repeated := periodicStart(t, root, "--profile", "quick", "--capture")
	periodicRecord(t, root, repeated, "dependencies", true, periodicIssue())
	finding := periodicRead(t, root, repeated).CheckRecords[0].CheckFindings[0]
	if finding.RepairRequest != "REQ-002" || finding.CaptureAction != "link" || finding.FindingDelta != "unchanged" {
		t.Fatal(finding)
	}
}

func TestPeriodicStatusDoesNotRewriteExistingRecordsAndUndefinedQuickIsIncomplete(t *testing.T) {
	root := periodicFixture(t)
	run := periodicStart(t, root, "--profile", "quick")
	periodicRecord(t, root, run, "dependencies", true)
	periodicFinish(t, root, run)
	recordPath := filepath.Join(root, run, "record.json")
	before, _ := os.ReadFile(recordPath)
	periodicInvoke(t, root, "status")
	after, _ := os.ReadFile(recordPath)
	if string(before) != string(after) {
		t.Fatal("status changed history")
	}
	primePath := filepath.Join(root, "prime-periodic-check.md")
	data, _ := os.ReadFile(primePath)
	if err := os.WriteFile(primePath, []byte(strings.ReplaceAll(string(data), `"quick":true`, `"quick":false`)), 0644); err != nil {
		t.Fatal(err)
	}
	undefined := periodicRead(t, root, periodicStart(t, root, "--profile", "quick"))
	if undefined.RunOutcome != "incomplete" || undefined.ConfigurationError != "quick subset is undefined" {
		t.Fatal(undefined)
	}
}

func TestPeriodicCompletedOrderUsesFinishTimeAndChangedFindingIsNotNew(t *testing.T) {
	root := periodicFixture(t)
	first := periodicStart(t, root, "--profile", "quick")
	second := periodicStart(t, root, "--profile", "quick")
	periodicRecord(t, root, second, "dependencies", true, periodicIssue())
	periodicFinish(t, root, second)
	changed := periodicIssue()
	changed.FindingSummary = "Package vulnerable at critical severity"
	periodicRecord(t, root, first, "dependencies", true, changed)
	if finding := periodicRead(t, root, first).CheckRecords[0].CheckFindings[0]; finding.FindingDelta != "changed" {
		t.Fatal(finding)
	}
	periodicFinish(t, root, first)
	for _, line := range strings.Split(periodicInvoke(t, root, "status"), "\n") {
		if strings.Contains(line, "last completed:") && strings.Contains(line, periodicHistoryDirectory) && !strings.Contains(line, first) {
			t.Fatal(line)
		}
	}
}

func TestPeriodicRecordCannotRedirectReportOrMintFromIncompleteCoverage(t *testing.T) {
	root := periodicFixture(t)
	run := periodicStart(t, root, "--profile", "quick", "--capture")
	periodicRecord(t, root, run, "dependencies", false, periodicIssue())
	payload := periodicInput(t, map[string]string{"check_id": "dependencies", "finding_id": "advisory/package", "request_id": "REQ-999", "disposition": "created", "note": "not actually captured"})
	if _, err := periodicCommand(root, []string{"link", "--run", run, "--input", payload}); err == nil {
		t.Fatal("nonexistent capture accepted")
	}
	record := periodicRead(t, root, run)
	record.ReportPath = "unexpected.md"
	data, _ := json.Marshal(record)
	if err := os.WriteFile(filepath.Join(root, run, "record.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := periodicCommand(root, []string{"finish", "--run", run}); err == nil {
		t.Fatal("report redirection accepted")
	}
}
