package gittransaction

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestMutationRequiresGit(t *testing.T) {
	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: t.TempDir(),
		TargetPaths:    []string{"target.txt"},
	}, func(*MutationRecorder) error { return nil })
	if result.ExitCode != 2 || result.Failure == nil || result.Failure.Kind != FailureNotGit {
		t.Fatalf("result = %#v", result)
	}
}

func TestDirtyTargetIsRefusedButUnrelatedDirtIsAllowed(t *testing.T) {
	repositoryRoot := newRepository(t)
	writeFile(t, repositoryRoot, "target.txt", "initial\n")
	writeFile(t, repositoryRoot, "unrelated.txt", "initial\n")
	commitAll(t, repositoryRoot, "initial")
	writeFile(t, repositoryRoot, "target.txt", "user change\n")
	called := false
	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"target.txt"},
	}, func(*MutationRecorder) error { called = true; return nil })
	if called || result.ExitCode != 1 || result.Failure == nil || result.Failure.Kind != FailureDirtyTarget {
		t.Fatalf("dirty target result = %#v, called = %v", result, called)
	}

	runFixtureGit(t, repositoryRoot, "restore", "--", "target.txt")
	writeFile(t, repositoryRoot, "unrelated.txt", "user change\n")
	result = ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"target.txt"},
	}, func(recorder *MutationRecorder) error {
		called = true
		if err := recorder.RecordTouched("target.txt"); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(repositoryRoot, "target.txt"), []byte("tool change\n"), 0o644)
	})
	if result.ExitCode != 0 {
		t.Fatalf("unrelated dirt result = %#v", result)
	}
	if got := readFile(t, repositoryRoot, "unrelated.txt"); got != "user change\n" {
		t.Fatalf("unrelated path changed: %q", got)
	}
}

func TestDryRunDoesNotMutateAndCannotCommit(t *testing.T) {
	repositoryRoot := newRepository(t)
	called := false
	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"new.txt"},
		DryRun:         true,
	}, func(*MutationRecorder) error { called = true; return nil })
	if result.ExitCode != 0 || called {
		t.Fatalf("dry-run result = %#v, called = %v", result, called)
	}

	result = ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"new.txt"},
		DryRun:         true,
		Commit:         true,
	}, func(*MutationRecorder) error { return nil })
	if result.ExitCode != 2 || result.Failure == nil || result.Failure.Kind != FailureInvalidOptions {
		t.Fatalf("dry-run commit result = %#v", result)
	}
}

func TestCommitRequiresInitiallyEmptyIndex(t *testing.T) {
	repositoryRoot := newRepository(t)
	writeFile(t, repositoryRoot, "staged.txt", "staged\n")
	runFixtureGit(t, repositoryRoot, "add", "staged.txt")
	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"target.txt"},
		Commit:         true,
		CommitMessage:  "target change",
	}, func(*MutationRecorder) error { return nil })
	if result.ExitCode != 1 || result.Failure == nil || result.Failure.Kind != FailureDirtyIndex {
		t.Fatalf("result = %#v", result)
	}
}

func TestPreCommitFailureRestoresTrackedAndRemovesOnlyCreatedTargets(t *testing.T) {
	repositoryRoot := newRepository(t)
	writeFile(t, repositoryRoot, "tracked.txt", "initial\n")
	commitAll(t, repositoryRoot, "initial")
	unrelatedPath := filepath.Join(repositoryRoot, "unrelated.txt")
	writeFile(t, repositoryRoot, "unrelated.txt", "keep\n")
	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"tracked.txt", "created.txt"},
	}, func(recorder *MutationRecorder) error {
		if err := recorder.RecordTouched("tracked.txt"); err != nil {
			return err
		}
		if err := recorder.RecordCreated("created.txt"); err != nil {
			return err
		}
		writeFile(t, repositoryRoot, "tracked.txt", "mutated\n")
		writeFile(t, repositoryRoot, "created.txt", "created\n")
		return errors.New("forced mutation failure")
	})
	if result.ExitCode != 3 || result.Rollback.Status != RollbackSucceeded {
		t.Fatalf("result = %#v", result)
	}
	if got := readFile(t, repositoryRoot, "tracked.txt"); got != "initial\n" {
		t.Fatalf("tracked file after rollback = %q", got)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "created.txt")); !os.IsNotExist(err) {
		t.Fatalf("created target remains: %v", err)
	}
	if got := readPath(t, unrelatedPath); got != "keep\n" {
		t.Fatalf("unrelated path changed: %q", got)
	}
}

func TestIncompleteRollbackReportsRiskWithoutRecursiveDeletion(t *testing.T) {
	repositoryRoot := newRepository(t)
	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"created"},
	}, func(recorder *MutationRecorder) error {
		if err := recorder.RecordCreated("created"); err != nil {
			return err
		}
		if err := os.Mkdir(filepath.Join(repositoryRoot, "created"), 0o755); err != nil {
			return err
		}
		writeFile(t, repositoryRoot, "created/child.txt", "must remain\n")
		return errors.New("forced mutation failure")
	})
	if result.ExitCode != 4 || result.Rollback.Status != RollbackIncomplete {
		t.Fatalf("result = %#v", result)
	}
	if got := readFile(t, repositoryRoot, "created/child.txt"); got != "must remain\n" {
		t.Fatalf("recursive rollback touched child: %q", got)
	}
}

func TestCommitContainsOnlyChangedTargetsAndPostCommitFailureReportsRevert(t *testing.T) {
	repositoryRoot := newRepository(t)
	writeFile(t, repositoryRoot, "target.txt", "initial\n")
	writeFile(t, repositoryRoot, "unrelated.txt", "initial\n")
	commitAll(t, repositoryRoot, "initial")
	writeFile(t, repositoryRoot, "unrelated.txt", "user change\n")
	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"target.txt", "created.txt"},
		Commit:         true,
		CommitMessage:  "exact target change",
		PostCommitVerify: func(context.Context, string) error {
			return errors.New("forced post-commit verification failure")
		},
	}, func(recorder *MutationRecorder) error {
		if err := recorder.RecordTouched("target.txt"); err != nil {
			return err
		}
		if err := recorder.RecordCreated("created.txt"); err != nil {
			return err
		}
		writeFile(t, repositoryRoot, "target.txt", "committed\n")
		writeFile(t, repositoryRoot, "created.txt", "committed\n")
		return nil
	})
	if result.ExitCode != 4 || result.Failure == nil || result.Failure.Kind != FailureCommittedRisk {
		t.Fatalf("result = %#v", result)
	}
	if result.CommitSHA == "" || !reflect.DeepEqual(result.RevertArgv, []string{"git", "revert", result.CommitSHA}) {
		t.Fatalf("commit/revert result = %#v", result)
	}
	if got := commitPaths(t, repositoryRoot, result.CommitSHA); !reflect.DeepEqual(got, []string{"created.txt", "target.txt"}) {
		t.Fatalf("committed paths = %#v", got)
	}
	if got := readFile(t, repositoryRoot, "target.txt"); got != "committed\n" {
		t.Fatalf("post-commit failure rewrote history/worktree: %q", got)
	}
	if got := readFile(t, repositoryRoot, "unrelated.txt"); got != "user change\n" {
		t.Fatalf("unrelated dirt changed: %q", got)
	}
}

func newRepository(t *testing.T) string {
	t.Helper()
	repositoryRoot := t.TempDir()
	runFixtureGit(t, repositoryRoot, "init", "-q")
	runFixtureGit(t, repositoryRoot, "config", "user.name", "Do Work Test")
	runFixtureGit(t, repositoryRoot, "config", "user.email", "do-work@example.invalid")
	return repositoryRoot
}

func commitAll(t *testing.T, repositoryRoot, message string) {
	t.Helper()
	runFixtureGit(t, repositoryRoot, "add", "-A")
	runFixtureGit(t, repositoryRoot, "commit", "-q", "-m", message)
}

func runFixtureGit(t *testing.T, repositoryRoot string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repositoryRoot}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}

func writeFile(t *testing.T, repositoryRoot, relativePath, content string) {
	t.Helper()
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolutePath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, repositoryRoot, relativePath string) string {
	t.Helper()
	return readPath(t, filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)))
}

func readPath(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func commitPaths(t *testing.T, repositoryRoot, commitSHA string) []string {
	t.Helper()
	output := runFixtureGit(t, repositoryRoot, "diff-tree", "--no-commit-id", "--name-only", "-r", commitSHA)
	paths := strings.Fields(output)
	sort.Strings(paths)
	return paths
}
