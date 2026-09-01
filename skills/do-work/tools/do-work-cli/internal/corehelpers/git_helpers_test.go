package corehelpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

import "github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"

func TestStageExactDeletionStagesOnlyNamedPath(t *testing.T) {
	repository := newGitFixture(t)
	first := filepath.Join(repository, "first.txt")
	second := filepath.Join(repository, "second.txt")
	_ = os.Remove(first)
	_ = os.Remove(second)
	result := handleStageDeletion(testContext(repository), []string{"--path", "first.txt"})
	if result.Outcome != "success" {
		t.Fatalf("result=%+v", result)
	}
	output, err := exec.Command("git", "-C", repository, "diff", "--cached", "--name-only").Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "first.txt\n" {
		t.Fatalf("staged=%q", output)
	}
}

func TestMutatingGitHelpersDryRunLeavesIndexAndExcludeUnchanged(t *testing.T) {
	repository := newGitFixture(t)
	if err := os.Remove(filepath.Join(repository, "first.txt")); err != nil {
		t.Fatal(err)
	}
	stageResult := handleStageDeletion(testContext(repository), []string{"--path", "first.txt", "--dry-run"})
	if stageResult.Outcome != "success" || !hasFinding(stageResult, "DELETION-DRY-RUN") {
		t.Fatalf("stage dry-run=%#v", stageResult)
	}
	if cached := runGitOutput(t, repository, "diff", "--cached", "--name-only"); cached != "" {
		t.Fatalf("dry-run changed index: %q", cached)
	}
	excludeOutput := runGitOutput(t, repository, "rev-parse", "--git-path", "info/exclude")
	excludePath := strings.TrimSpace(excludeOutput)
	if !filepath.IsAbs(excludePath) {
		excludePath = filepath.Join(repository, excludePath)
	}
	before, _ := os.ReadFile(excludePath)
	excludeResult := handleAddExclude(testContext(repository), []string{"--probe-path", "private/cache", "--dry-run"})
	if excludeResult.Outcome != "success" || !hasFinding(excludeResult, "GIT-EXCLUDE-DRY-RUN") {
		t.Fatalf("exclude dry-run=%#v", excludeResult)
	}
	after, _ := os.ReadFile(excludePath)
	if string(after) != string(before) {
		t.Fatalf("dry-run changed exclude: before=%q after=%q", before, after)
	}
}

func runGitOutput(t *testing.T, repository string, args ...string) string {
	t.Helper()
	output, err := exec.Command("git", append([]string{"-C", repository}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
	return string(output)
}
func newGitFixture(t *testing.T) string {
	t.Helper()
	repository := t.TempDir()
	run := func(args ...string) {
		command := exec.Command("git", append([]string{"-C", repository}, args...)...)
		command.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.com", "GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.com")
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, output)
		}
	}
	run("init", "-q")
	_ = os.WriteFile(filepath.Join(repository, "first.txt"), []byte("1"), 0o600)
	_ = os.WriteFile(filepath.Join(repository, "second.txt"), []byte("2"), 0o600)
	run("add", ".")
	run("commit", "-qm", "fixture")
	return repository
}
func testContext(root string) commandruntime.ExecutionContext {
	return commandruntime.ExecutionContext{RepositoryRoot: root}
}
