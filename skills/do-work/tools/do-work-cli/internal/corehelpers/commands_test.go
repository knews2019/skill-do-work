package corehelpers

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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
