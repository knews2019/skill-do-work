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

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestMutationRequiresGit(t *testing.T) {
	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: t.TempDir(),
		TargetPaths:    []string{"target.txt"},
	}, func(*MutationRecorder) error { return nil })
	if resultmodel.ExitCode(result.Outcome) != 2 || result.Failure == nil || result.Failure.Kind != FailureNotGit {
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
	if called || resultmodel.ExitCode(result.Outcome) != 1 || result.Failure == nil || result.Failure.Kind != FailureDirtyTarget {
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
	if resultmodel.ExitCode(result.Outcome) != 0 {
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
	if resultmodel.ExitCode(result.Outcome) != 0 || called {
		t.Fatalf("dry-run result = %#v, called = %v", result, called)
	}

	result = ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"new.txt"},
		DryRun:         true,
		Commit:         true,
	}, func(*MutationRecorder) error { return nil })
	if resultmodel.ExitCode(result.Outcome) != 2 || result.Failure == nil || result.Failure.Kind != FailureInvalidOptions {
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
	if resultmodel.ExitCode(result.Outcome) != 1 || result.Failure == nil || result.Failure.Kind != FailureDirtyIndex {
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
	if resultmodel.ExitCode(result.Outcome) != 3 || result.Rollback.Status != resultmodel.RollbackSucceeded {
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
	if resultmodel.ExitCode(result.Outcome) != 4 || result.Rollback.Status != resultmodel.RollbackIncomplete {
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
	if resultmodel.ExitCode(result.Outcome) != 4 || result.Failure == nil || result.Failure.Kind != FailureCommittedRisk {
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

// A malformed .gitattributes line makes every later git invocation print a warning on
// stderr while stdout stays empty. Folding the two streams together made that warning
// read as porcelain content, so a clean target was refused as dirty.
func TestGitStderrWarningIsNotReadAsTargetDirt(t *testing.T) {
	repositoryRoot := newRepository(t)
	writeFile(t, repositoryRoot, "target.txt", "initial\n")
	writeFile(t, repositoryRoot, ".gitattributes", "*.txt [attr]bogus\n")
	commitAll(t, repositoryRoot, "initial")
	called := false
	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"target.txt"},
	}, func(recorder *MutationRecorder) error {
		called = true
		if err := recorder.RecordTouched("target.txt"); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(repositoryRoot, "target.txt"), []byte("tool change\n"), 0o644)
	})
	if !called || resultmodel.ExitCode(result.Outcome) != 0 || result.Failure != nil {
		t.Fatalf("git stderr warning was treated as target dirt: result = %#v, called = %v", result, called)
	}
	if got := readFile(t, repositoryRoot, "target.txt"); got != "tool change\n" {
		t.Fatalf("target after mutation = %q", got)
	}
}

// The user's own staged work is unrelated dirt, which the transaction contract allows.
// Checking the whole index after a rollback turned a complete rollback into a
// committed-state risk purely because someone else had a file staged.
func TestUnrelatedStagedWorkDoesNotBreakRollback(t *testing.T) {
	repositoryRoot := newRepository(t)
	writeFile(t, repositoryRoot, "tracked.txt", "initial\n")
	commitAll(t, repositoryRoot, "initial")
	writeFile(t, repositoryRoot, "unrelated.txt", "staged by the user\n")
	runFixtureGit(t, repositoryRoot, "add", "unrelated.txt")
	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"tracked.txt"},
	}, func(recorder *MutationRecorder) error {
		if err := recorder.RecordTouched("tracked.txt"); err != nil {
			return err
		}
		writeFile(t, repositoryRoot, "tracked.txt", "mutated\n")
		return errors.New("forced mutation failure")
	})
	if resultmodel.ExitCode(result.Outcome) != 3 || result.Rollback.Status != resultmodel.RollbackSucceeded {
		t.Fatalf("unrelated staged work broke a complete rollback: result = %#v", result)
	}
	if got := readFile(t, repositoryRoot, "tracked.txt"); got != "initial\n" {
		t.Fatalf("tracked file after rollback = %q", got)
	}
	if got := readFile(t, repositoryRoot, "unrelated.txt"); got != "staged by the user\n" {
		t.Fatalf("unrelated staged file changed: %q", got)
	}
}

func newRepository(t *testing.T) string {
	t.Helper()
	// Fixtures measure exactly what git writes to each stream, so an ambient global or
	// system config must not be able to add or suppress a warning.
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	t.Setenv("GIT_TERMINAL_PROMPT", "0")
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

// REQ-406 built the committing path but registered no command, so nothing observed a
// SUCCESSFUL --commit end to end. This drives a real repository: the commit must land, hold
// exactly the declared targets, leave an empty index, and report a revert command.
func TestSuccessfulCommitLandsExactlyTheDeclaredTargets(t *testing.T) {
	repositoryRoot := newRepository(t)
	writeFile(t, repositoryRoot, "tracked.txt", "original\n")
	writeFile(t, repositoryRoot, "untouched.txt", "leave me\n")
	commitAll(t, repositoryRoot, "seed")
	headBeforeCommit := runFixtureGit(t, repositoryRoot, "rev-parse", "HEAD")

	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"tracked.txt", "created.txt"},
		Commit:         true,
		CommitMessage:  "exact-path commit fixture",
	}, func(recorder *MutationRecorder) error {
		writeFile(t, repositoryRoot, "tracked.txt", "updated\n")
		if err := recorder.RecordTouched("tracked.txt"); err != nil {
			return err
		}
		writeFile(t, repositoryRoot, "created.txt", "brand new\n")
		return recorder.RecordCreated("created.txt")
	})

	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("outcome = %q, failure = %#v", result.Outcome, result.Failure)
	}
	if result.CommitSHA == "" || result.CommitSHA == headBeforeCommit {
		t.Fatalf("commit SHA = %q, want a new commit", result.CommitSHA)
	}
	if want := []string{"created.txt", "tracked.txt"}; !equalStrings(commitPaths(t, repositoryRoot, result.CommitSHA), want) {
		t.Errorf("committed paths = %v, want %v", commitPaths(t, repositoryRoot, result.CommitSHA), want)
	}
	if len(result.RevertArgv) != 3 || result.RevertArgv[1] != "revert" || result.RevertArgv[2] != result.CommitSHA {
		t.Errorf("revert argv = %v, want git revert %s", result.RevertArgv, result.CommitSHA)
	}
	if result.Rollback.Status != resultmodel.RollbackNotNeeded {
		t.Errorf("rollback status = %q, want %q", result.Rollback.Status, resultmodel.RollbackNotNeeded)
	}
	if status := runFixtureGit(t, repositoryRoot, "status", "--porcelain"); status != "" {
		t.Errorf("worktree is not clean after a successful commit: %q", status)
	}
	if readFile(t, repositoryRoot, "untouched.txt") != "leave me\n" {
		t.Errorf("an undeclared path was changed by the commit")
	}
	commandResult := BuildCommandResult("apply", result)
	if len(commandResult.Changes) != 2 {
		t.Fatalf("changes = %#v, want both declared targets", commandResult.Changes)
	}
	for _, change := range commandResult.Changes {
		if change.Detail != "committed in "+result.CommitSHA {
			t.Errorf("change %q detail = %q, want it to name the commit", change.Path, change.Detail)
		}
	}
}

// commit_failed had a remediation template but no behavioural test. A commit hook that
// refuses makes git commit fail for real, so the transaction must roll back every declared
// target and report the kind rather than a generic mutation failure.
func TestRefusedCommitRollsBackAndReportsCommitFailed(t *testing.T) {
	repositoryRoot := newRepository(t)
	writeFile(t, repositoryRoot, "tracked.txt", "original\n")
	commitAll(t, repositoryRoot, "seed")
	hooksDirectory := runFixtureGit(t, repositoryRoot, "rev-parse", "--git-path", "hooks")
	if !filepath.IsAbs(hooksDirectory) {
		hooksDirectory = filepath.Join(repositoryRoot, hooksDirectory)
	}
	if err := os.MkdirAll(hooksDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hooksDirectory, "pre-commit"),
		[]byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"tracked.txt", "created.txt"},
		Commit:         true,
		CommitMessage:  "this commit is refused by a hook",
	}, func(recorder *MutationRecorder) error {
		writeFile(t, repositoryRoot, "tracked.txt", "updated\n")
		if err := recorder.RecordTouched("tracked.txt"); err != nil {
			return err
		}
		writeFile(t, repositoryRoot, "created.txt", "brand new\n")
		return recorder.RecordCreated("created.txt")
	})

	if result.Outcome != resultmodel.OutcomeRolledBack {
		t.Fatalf("outcome = %q, want %q (failure %#v)", result.Outcome, resultmodel.OutcomeRolledBack, result.Failure)
	}
	if result.Failure == nil || result.Failure.Kind != FailureCommit {
		t.Fatalf("failure = %#v, want kind %q", result.Failure, FailureCommit)
	}
	if result.Rollback.Status != resultmodel.RollbackSucceeded {
		t.Fatalf("rollback status = %q, errors %v", result.Rollback.Status, result.Rollback.Errors)
	}
	if readFile(t, repositoryRoot, "tracked.txt") != "original\n" {
		t.Errorf("a refused commit left the tracked target modified")
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "created.txt")); !os.IsNotExist(err) {
		t.Errorf("a refused commit left the created target behind")
	}
	if status := runFixtureGit(t, repositoryRoot, "status", "--porcelain"); status != "" {
		t.Errorf("worktree is not clean after rollback: %q", status)
	}
	finding := BuildCommandResult("apply", result).Findings[0]
	if finding.Code != FindingCode(FailureCommit) {
		t.Errorf("finding code = %q, want %q", finding.Code, FindingCode(FailureCommit))
	}
}

// A path the mutation created but recorded only as touched used to report success, so the
// result described a creation the transaction never saw. The success path now consults
// state.existed and treats it as an unrecorded mutation.
func TestAnUnrecordedCreationIsRolledBackRatherThanReportedAsSuccess(t *testing.T) {
	repositoryRoot := newRepository(t)
	writeFile(t, repositoryRoot, "tracked.txt", "original\n")
	commitAll(t, repositoryRoot, "seed")

	result := ExecuteTransaction(context.Background(), TransactionOptions{
		RepositoryRoot: repositoryRoot,
		TargetPaths:    []string{"created.txt"},
	}, func(recorder *MutationRecorder) error {
		writeFile(t, repositoryRoot, "created.txt", "brand new\n")
		return recorder.RecordTouched("created.txt")
	})

	if result.Outcome == resultmodel.OutcomeSuccess {
		t.Fatalf("an unrecorded creation reported success: %#v", result)
	}
	if result.Failure == nil || result.Failure.Kind != FailureMutation {
		t.Fatalf("failure = %#v, want kind %q", result.Failure, FailureMutation)
	}
	if _, err := os.Stat(filepath.Join(repositoryRoot, "created.txt")); err != nil {
		t.Logf("the unrecorded creation was removed by rollback, which is acceptable")
	}
}
