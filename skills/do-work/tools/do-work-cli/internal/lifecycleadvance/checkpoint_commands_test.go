package lifecycleadvance

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAdvanceCheckpointChangesOnlyCheckpointAndPreservesLiveEntries(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeAdvanceRequest(t, repositoryRoot, "queue", "REQ-714", "pending", "", "")
	writeAdvanceFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-800: foreign — writer: other:/checkout\n  keep foreign detail\n- REQ-801: unknown owner\n  keep unknown detail\n\n## Session Notes\n\nOld note.\n")
	writeAdvanceFile(t, repositoryRoot, "project.txt", "base\n")
	runAdvanceGit(t, repositoryRoot, "init", "-q")
	runAdvanceGit(t, repositoryRoot, "config", "user.name", "Checkpoint Test")
	runAdvanceGit(t, repositoryRoot, "config", "user.email", "checkpoint@example.invalid")
	runAdvanceGit(t, repositoryRoot, "add", ".")
	runAdvanceGit(t, repositoryRoot, "commit", "-qm", "fixture")
	if err := os.WriteFile(filepath.Join(repositoryRoot, "project.txt"), []byte("base\nforeign dirt\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	command := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "advance", "--checkpoint")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("advance checkpoint: %v\n%s", err, output)
	}
	var result struct {
		Checkpoint struct {
			CheckpointPath  string `json:"checkpoint_path"`
			PreservedClaims int    `json:"preserved_claims"`
		} `json:"checkpoint"`
		Changes []struct {
			Path string `json:"path"`
		} `json:"changes"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("decode checkpoint result: %v\n%s", err, output)
	}
	if result.Checkpoint.CheckpointPath != "do-work/CHECKPOINT.md" || result.Checkpoint.PreservedClaims != 2 || len(result.Changes) != 1 || result.Changes[0].Path != "do-work/CHECKPOINT.md" {
		t.Fatalf("checkpoint result did not expose the exact mutation: %#v", result)
	}
	checkpoint, _ := os.ReadFile(filepath.Join(repositoryRoot, "do-work", "CHECKPOINT.md"))
	for _, exact := range []string{"- REQ-800: foreign — writer: other:/checkout\n  keep foreign detail", "- REQ-801: unknown owner\n  keep unknown detail"} {
		if !strings.Contains(string(checkpoint), exact) {
			t.Fatalf("checkpoint lost live record %q:\n%s", exact, checkpoint)
		}
	}
	project, _ := os.ReadFile(filepath.Join(repositoryRoot, "project.txt"))
	if string(project) != "base\nforeign dirt\n" {
		t.Fatalf("project dirt changed: %q", project)
	}
	status := string(runAdvanceGit(t, repositoryRoot, "status", "--porcelain=v1", "--untracked-files=all"))
	if status != " M do-work/CHECKPOINT.md\n M project.txt\n" {
		t.Fatalf("unexpected changed paths:\n%s", status)
	}
}

func TestWorkingAdvanceRemainsReadOnlyAfterCheckpointMode(t *testing.T) {
	repositoryRoot := t.TempDir()
	writeAdvanceRequest(t, repositoryRoot, "working", "REQ-715", "claimed", "", "")
	writeAdvanceFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Session Checkpoint\n\n## In Progress (interrupted)\n")
	before := advanceTreeDigest(t, repositoryRoot)
	command := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "advance", "REQ-715")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("ordinary advance: %v\n%s", err, output)
	}
	if after := advanceTreeDigest(t, repositoryRoot); before != after {
		t.Fatalf("ordinary advance changed bytes")
	}
}
