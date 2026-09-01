package corehelpers

import (
	"os"
	"os/exec"
	"path/filepath"
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
