package corehelpers

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestEveryRemainingUtilityHasOneHandler(t *testing.T) {
	expected := []string{CommandPreflight, CommandQualify, CommandScopeDrift, CommandInventory, CommandAssociate, CommandProtectedInventory, CommandRecordCommit, CommandCaptureScreenshot, CommandAtomicDownload, CommandAddExclude, CommandBlockedCheck, CommandShowCommitDiff, CommandStageDeletion, CommandCleanupReservations, CommandRepairTimestamps, CommandAuditTimestamps, CommandHandoffSurvey}
	handlers := Handlers()
	if len(handlers) != len(expected) {
		t.Fatalf("handlers=%d want %d", len(handlers), len(expected))
	}
	for _, name := range expected {
		if handlers[name] == nil {
			t.Errorf("missing handler %s", name)
		}
	}
}

func TestNonInformationalFindingsReceiveCommandSpecificActions(t *testing.T) {
	tests := []struct {
		code       string
		paths      []string
		wantNext   string
		wantVerify string
	}{
		{"SCOPE-UNDECLARED-TOUCH", []string{"internal/new.go"}, "git status --short -- internal/new.go", "git diff --name-status -- internal/new.go"},
		{"QUALIFY-DEBUG-TOKEN-INTRODUCED", []string{"internal/new.go"}, "git diff -- internal/new.go", "git diff --check"},
		{"RESERVATION-MALFORMED", []string{"do-work/.req-reservations/bad"}, "do-work-cli cleanup-req-reservations --dry-run", "do-work-cli --format json cleanup-req-reservations --dry-run"},
		{"PREFLIGHT-NODE-MODULES-MISSING", []string{"package.json"}, "do-work-cli preflight --dry-run", "do-work-cli --format json preflight --dry-run"},
	}
	for _, test := range tests {
		finding := helperFinding(test.code, resultmodel.SeverityWarning, test.paths, "evidence", resultmodel.FixabilityManual, "blocked", nil, nil)
		if strings.Join(finding.NextArgv, " ") != test.wantNext || strings.Join(finding.VerificationArgv, " ") != test.wantVerify {
			t.Errorf("%s actions next=%q verify=%q", test.code, finding.NextArgv, finding.VerificationArgv)
		}
		if strings.Contains(strings.Join(append(append([]string{}, finding.NextArgv...), finding.VerificationArgv...), " "), "doctor") {
			t.Errorf("%s retained generic doctor placeholder", test.code)
		}
	}
}

func TestRealCommandRendersActionableTextAndJSON(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	moduleRoot := filepath.Clean(filepath.Join(workingDirectory, "..", ".."))
	repositoryRoot := t.TempDir()
	for _, format := range []string{"text", "json"} {
		command := exec.Command("go", "run", "./cmd/do-work-cli", "--repo-root", repositoryRoot, "--format", format, CommandPreflight)
		command.Dir = moduleRoot
		output, runError := command.CombinedOutput()
		if runError != nil {
			t.Fatalf("%s command: %v\n%s", format, runError, output)
		}
		if format == "text" {
			if !strings.Contains(string(output), "preflight: success") || !strings.Contains(string(output), "PREFLIGHT-GIT-UNAVAILABLE") {
				t.Fatalf("text=%s", output)
			}
			continue
		}
		var decoded map[string]any
		if err := json.Unmarshal(output, &decoded); err != nil {
			t.Fatalf("json: %v\n%s", err, output)
		}
		if decoded["command"] != CommandPreflight || decoded["findings"] == nil || decoded["changes"] == nil {
			t.Fatalf("json=%v", decoded)
		}
	}
}

func TestDryRunSurfacesDoNotMutateBaselineDownloadOrTimestamps(t *testing.T) {
	repository := newGitFixture(t)
	marker := filepath.Join(repository, "baseline-command-ran")
	preflight := handlePreflight(testContext(repository), []string{"--dry-run", "touch " + marker})
	if preflight.Outcome != "success" || !hasFinding(preflight, "PREFLIGHT-BASELINE-DRY-RUN") {
		t.Fatalf("preflight=%#v", preflight)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("preflight dry-run executed command: %v", err)
	}
	target := filepath.Join(repository, "download.bin")
	download := handleAtomicDownload(testContext(repository), []string{"--source-url", "https://example.invalid/archive", "--target-path", target, "--dry-run"})
	if download.Outcome != "success" || !hasFinding(download, "DOWNLOAD-DRY-RUN") {
		t.Fatalf("download=%#v", download)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("download dry-run published target: %v", err)
	}
	relativeRequest := "do-work/queue/REQ-888-future.md"
	absoluteRequest := filepath.Join(repository, filepath.FromSlash(relativeRequest))
	if err := os.MkdirAll(filepath.Dir(absoluteRequest), 0o755); err != nil {
		t.Fatal(err)
	}
	requestBytes := []byte("---\nid: REQ-888\nstatus: pending\ncreated_at: 2099-01-01T00:00:00Z\n---\nBody\n")
	if err := os.WriteFile(absoluteRequest, requestBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	runFixtureGitCommand(t, repository, "add", relativeRequest)
	runFixtureGitCommand(t, repository, "commit", "-qm", "future request")
	timestamp := handleRepairTimestamps(testContext(repository), []string{"--dry-run"})
	if timestamp.Outcome != "findings" || len(timestamp.Changes) == 0 {
		t.Fatalf("timestamp=%#v", timestamp)
	}
	if after, err := os.ReadFile(absoluteRequest); err != nil || string(after) != string(requestBytes) {
		t.Fatalf("timestamp dry-run changed bytes: %q err=%v", after, err)
	}
}

func TestAllSeventeenPublicCommandsRunInTextAndJSONWithStableStatusAndNoDryRunEffects(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	moduleRoot := filepath.Clean(filepath.Join(workingDirectory, "..", ".."))
	binary := filepath.Join(t.TempDir(), "do-work-cli")
	build := exec.Command("go", "build", "-o", binary, "./cmd/do-work-cli")
	build.Dir = moduleRoot
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build real CLI: %v: %s", err, output)
	}
	repository := newGitFixture(t)
	branch := strings.TrimSpace(runGitOutput(t, repository, "branch", "--show-current"))
	writeMatrixFile(t, repository, "do-work/working/REQ-901-matrix.md", "---\nid: REQ-901\nstatus: claimed\n---\n\n## AI Execution State (P-A-U Loop)\n- [x] **[PLAN]:** done\n- [x] **[APPLY]:** done\n- [x] **[UNIFY]:** done\n\n## Scope\n- `matrix-new.txt`\n\n## Implementation Summary\n- `matrix-new.txt` (new)\n")
	writeMatrixFile(t, repository, "do-work/archive/REQ-900-owner.md", "---\nid: REQ-900\nstatus: completed\ncompleted_at: 2026-08-31T21:00:00Z\ncommit: oldhash0\n---\n\n## Implementation Summary\n- `first.txt` (modified)\n")
	writeMatrixFile(t, repository, "do-work/queue/REQ-902-future.md", "---\nid: REQ-902\nstatus: pending\ncreated_at: 2099-01-01T00:00:00Z\n---\nBody\n")
	writeMatrixFile(t, repository, "do-work/archive/REQ-903-future.md", "---\nid: REQ-903\nstatus: completed\ncreated_at: 2099-01-01T00:00:00Z\ncompleted_at: 2099-01-02T00:00:00Z\n---\nBody\n")
	runFixtureGitCommand(t, repository, "add", "do-work")
	runFixtureGitCommand(t, repository, "commit", "-qm", "matrix authorities")
	implementationHash := strings.TrimSpace(runGitOutput(t, repository, "rev-parse", "HEAD"))
	writeMatrixFile(t, repository, "matrix-new.txt", "new behavior\n")
	writeMatrixFile(t, repository, "paths.txt", "first.txt\n")
	writeMatrixFile(t, repository, "probe.sh", "exit 0\n")
	writeMatrixFile(t, repository, "staging/source.png", "image\n")
	writeMatrixFile(t, repository, ".env.matrix", "TOKEN=x\n")
	if err := os.Remove(filepath.Join(repository, "first.txt")); err != nil {
		t.Fatal(err)
	}
	reservationRoot := filepath.Join(repository, "do-work", ".req-reservations")
	if err := os.MkdirAll(reservationRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(reservationRoot, "REQ-000999")
	if err := os.WriteFile(marker, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	past := time.Now().Add(-49 * time.Hour)
	_ = os.Chtimes(marker, past, past)
	commands := []struct {
		name string
		args []string
	}{
		{CommandPreflight, []string{"--dry-run"}},
		{CommandQualify, []string{"--request-path", "do-work/working/REQ-901-matrix.md"}},
		{CommandScopeDrift, []string{"--request-path", "do-work/working/REQ-901-matrix.md"}},
		{CommandInventory, nil},
		{CommandAssociate, []string{"--paths-file", "paths.txt"}},
		{CommandProtectedInventory, []string{"start", "--dry-run"}},
		{CommandRecordCommit, []string{"--request-path", "do-work/archive/REQ-900-owner.md", "--implementation-hash", implementationHash, "--dry-run"}},
		{CommandCaptureScreenshot, []string{"--staged", "--source", "staging/source.png", "--destination", "assets/result.png", "--dry-run"}},
		{CommandAtomicDownload, []string{"--source-url", "https://example.invalid/archive", "--target-path", filepath.Join(repository, "download.bin"), "--dry-run"}},
		{CommandAddExclude, []string{"--probe-path", "private/cache", "--dry-run"}},
		{CommandBlockedCheck, []string{"--probe-file", filepath.Join(repository, "probe.sh")}},
		{CommandShowCommitDiff, []string{"--commit", "HEAD"}},
		{CommandStageDeletion, []string{"--path", "first.txt", "--dry-run"}},
		{CommandCleanupReservations, []string{"--dry-run"}},
		{CommandRepairTimestamps, []string{"--dry-run"}},
		{CommandAuditTimestamps, []string{"--fix", "--dry-run"}},
		{CommandHandoffSurvey, []string{"--integration-branch", branch}},
	}
	beforeStatus := runGitOutput(t, repository, "status", "--porcelain=v1", "--untracked-files=all")
	for _, commandCase := range commands {
		t.Run(commandCase.name, func(t *testing.T) {
			statuses := map[string]int{}
			for _, format := range []string{"text", "json"} {
				arguments := []string{"--repo-root", repository, "--format", format, commandCase.name}
				arguments = append(arguments, commandCase.args...)
				command := exec.Command(binary, arguments...)
				output, runError := command.CombinedOutput()
				status := commandExitStatus(runError)
				statuses[format] = status
				if status > 1 {
					t.Fatalf("%s status=%d err=%v output=%s", format, status, runError, output)
				}
				if format == "json" {
					var decoded map[string]any
					if err := json.Unmarshal(output, &decoded); err != nil || decoded["command"] != commandCase.name || decoded["findings"] == nil || decoded["changes"] == nil {
						t.Fatalf("JSON err=%v decoded=%v output=%s", err, decoded, output)
					}
				} else if !strings.Contains(string(output), commandCase.name+":") {
					t.Fatalf("text output lacks command identity: %s", output)
				}
			}
			if statuses["text"] != statuses["json"] {
				t.Fatalf("format status mismatch: %v", statuses)
			}
		})
	}
	afterStatus := runGitOutput(t, repository, "status", "--porcelain=v1", "--untracked-files=all")
	if afterStatus != beforeStatus {
		t.Fatalf("command matrix changed repository state:\nbefore=%s\nafter=%s", beforeStatus, afterStatus)
	}
	if _, err := os.Stat(filepath.Join(repository, "assets", "result.png")); !os.IsNotExist(err) {
		t.Fatalf("screenshot dry-run published: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repository, "download.bin")); !os.IsNotExist(err) {
		t.Fatalf("download dry-run published: %v", err)
	}
}

func writeMatrixFile(t *testing.T, repository, relativePath, contents string) {
	t.Helper()
	absolutePath := filepath.Join(repository, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolutePath, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}

func commandExitStatus(err error) int {
	if err == nil {
		return 0
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return 127
}
