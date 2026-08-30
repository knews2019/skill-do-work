package cleanup

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestCleanupHandlerIsRegisteredAndJSONUsesSharedResult(t *testing.T) {
	repositoryRoot := cleanupRepository(t)
	writeCleanupFile(t, repositoryRoot, "do-work/queue/REQ-108-done.md", cleanupRequest("REQ-108", "completed", ""))
	commitCleanupFixture(t, repositoryRoot)
	var output bytes.Buffer
	runtime := commandruntime.NewRuntime(&output, Handlers())
	exitCode := runtime.Run([]string{"--repo-root", repositoryRoot, "--format", "json", "cleanup", "--dry-run"})
	if exitCode != 0 {
		t.Fatalf("exit = %d output=%s", exitCode, output.String())
	}
	var result resultmodel.CommandResult
	if err := json.Unmarshal(output.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Command != "cleanup" || len(result.Changes) != 1 {
		t.Fatalf("result = %#v", result)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "do-work/queue/REQ-108-done.md")); err != nil {
		t.Fatalf("dry run mutated: %v", err)
	}
}

func TestCleanupDestructiveFlagsRequireExactTargets(t *testing.T) {
	if _, err := parseCommandOptions([]string{"--restore-blanked", "../outside"}); err == nil {
		t.Fatal("traversal restore target accepted")
	}
	if _, err := parseCommandOptions([]string{"--discard-worktree", "main"}); err == nil {
		t.Fatal("non-builder branch accepted")
	}
	if _, err := parseCommandOptions([]string{"--dry-run", "--commit"}); err == nil {
		t.Fatal("dry-run commit combination accepted")
	}
	if _, err := parseCommandOptions([]string{"--commit", "--discard-worktree", "worktree-agent-REQ-210-test"}); err == nil {
		t.Fatal("commit plus destructive worktree discard accepted")
	}
}
