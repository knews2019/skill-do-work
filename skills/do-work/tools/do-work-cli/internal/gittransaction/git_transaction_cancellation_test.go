//go:build unix

package gittransaction

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// A pre-commit hook that ignores SIGTERM used to outlive the transaction that
// launched it. Cancellation signalled the whole group at once, so git died
// first, the hook was orphaned to init, and the later SIGKILL left an unreaped
// zombie — which kill(pid, 0) still reports as alive. The escalation also ran
// in a detached goroutine, so runGit returned before any of it happened.
//
// The invariant is pinned here, at the seam every Git process in this module is
// launched through, and not only through the report-image caller. The hook and
// the process the hook starts are both TERM-deaf, so nothing but the escalation
// can end them, and both must be gone rather than zombies when the call
// returns. The hook waits on its child, so the child's killed status is the
// hook's status: git aborts the commit and the existing commit-failure path
// rolls the target back to its preimage.
func TestCancelledCommitKillsTermDeafHookBeforeReturning(t *testing.T) {
	repositoryRoot := newRepository(t)
	writeFile(t, repositoryRoot, "tracked.txt", "original\n")
	commitAll(t, repositoryRoot, "seed")
	hooksDirectory := strings.TrimSpace(runFixtureGit(t, repositoryRoot, "rev-parse", "--git-path", "hooks"))
	if !filepath.IsAbs(hooksDirectory) {
		hooksDirectory = filepath.Join(repositoryRoot, hooksDirectory)
	}
	if err := os.MkdirAll(hooksDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	// The markers live outside the worktree so the hook cannot dirty it.
	markerDirectory := t.TempDir()
	hookMarker := filepath.Join(markerDirectory, "hook.pid")
	descendantMarker := filepath.Join(markerDirectory, "descendant.pid")
	hook := "#!/bin/sh\ntrap '' TERM\n" +
		"(trap '' TERM; sleep 30) &\n" +
		"descendant=$!\n" +
		"echo $descendant > " + descendantMarker + "\n" +
		"echo $$ > " + hookMarker + "\n" +
		"wait \"$descendant\"\n"
	if err := os.WriteFile(filepath.Join(hooksDirectory, "pre-commit"), []byte(hook), 0o755); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan TransactionResult, 1)
	go func() {
		done <- ExecuteTransaction(ctx, TransactionOptions{
			RepositoryRoot: repositoryRoot,
			TargetPaths:    []string{"tracked.txt"},
			Commit:         true,
			CommitMessage:  "cancelled while a TERM-deaf hook runs",
		}, func(recorder *MutationRecorder) error {
			writeFile(t, repositoryRoot, "tracked.txt", "updated\n")
			return recorder.RecordTouched("tracked.txt")
		})
	}()

	hookPID := awaitRecordedPID(t, hookMarker)
	descendantPID := awaitRecordedPID(t, descendantMarker)

	cancel()
	var result TransactionResult
	select {
	case result = <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("cancelled commit did not return")
	}

	for name, processID := range map[string]int{"hook": hookPID, "hook descendant": descendantPID} {
		if syscall.Kill(processID, 0) == nil {
			_ = syscall.Kill(processID, syscall.SIGKILL)
			t.Errorf("%s %d survived cancellation", name, processID)
		}
	}
	if result.Outcome != resultmodel.OutcomeRolledBack {
		t.Errorf("outcome = %q, want %q (failure %#v)", result.Outcome, resultmodel.OutcomeRolledBack, result.Failure)
	}
	if result.Failure == nil || result.Failure.Kind != FailureCommit {
		t.Errorf("failure = %#v, want kind %q", result.Failure, FailureCommit)
	}
	if result.CommitSHA != "" {
		t.Errorf("commit SHA = %q, want none for a cancelled commit", result.CommitSHA)
	}
	if result.Rollback.Status != resultmodel.RollbackSucceeded {
		t.Errorf("rollback status = %q, errors %v", result.Rollback.Status, result.Rollback.Errors)
	}
	if contents := readFile(t, repositoryRoot, "tracked.txt"); contents != "original\n" {
		t.Errorf("cancelled commit did not roll back to the preimage: %q", contents)
	}
	if status := runFixtureGit(t, repositoryRoot, "status", "--porcelain"); status != "" {
		t.Errorf("worktree is not clean after the cancelled commit rolled back: %q", status)
	}
}

// A commit that lands and then reports failure must not be unwound as though
// nothing happened: rolling the worktree back from the advanced HEAD would
// restore the committed bytes and call them the preimage. This hook exits zero
// once its descendants are killed, which is what makes the ref update land
// while the caller is already cancelling.
func TestCancelledCommitThatLandsReportsCommittedRisk(t *testing.T) {
	repositoryRoot := newRepository(t)
	writeFile(t, repositoryRoot, "tracked.txt", "original\n")
	commitAll(t, repositoryRoot, "seed")
	seedHead := strings.TrimSpace(runFixtureGit(t, repositoryRoot, "rev-parse", "HEAD"))
	hooksDirectory := strings.TrimSpace(runFixtureGit(t, repositoryRoot, "rev-parse", "--git-path", "hooks"))
	if !filepath.IsAbs(hooksDirectory) {
		hooksDirectory = filepath.Join(repositoryRoot, hooksDirectory)
	}
	if err := os.MkdirAll(hooksDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	markerDirectory := t.TempDir()
	hookMarker := filepath.Join(markerDirectory, "hook.pid")
	// `wait` with no operand reports success however its children ended, so the
	// hook passes after the escalation kills the process it was waiting on.
	hook := "#!/bin/sh\ntrap '' TERM\n" +
		"(trap '' TERM; sleep 30) &\n" +
		"echo $$ > " + hookMarker + "\n" +
		"wait\n"
	if err := os.WriteFile(filepath.Join(hooksDirectory, "pre-commit"), []byte(hook), 0o755); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan TransactionResult, 1)
	go func() {
		done <- ExecuteTransaction(ctx, TransactionOptions{
			RepositoryRoot: repositoryRoot,
			TargetPaths:    []string{"tracked.txt"},
			Commit:         true,
			CommitMessage:  "cancelled after the commit landed",
		}, func(recorder *MutationRecorder) error {
			writeFile(t, repositoryRoot, "tracked.txt", "updated\n")
			return recorder.RecordTouched("tracked.txt")
		})
	}()

	awaitRecordedPID(t, hookMarker)
	cancel()
	var result TransactionResult
	select {
	case result = <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("cancelled commit did not return")
	}

	// Whether git reaches the ref update before the escalation is a race, so
	// both branches are asserted rather than skipping one. The defect this pins
	// is the advanced-HEAD branch reporting a rollback and leaving the
	// committed bytes in the worktree as though they were the preimage.
	head := strings.TrimSpace(runFixtureGit(t, repositoryRoot, "rev-parse", "HEAD"))
	if head == seedHead {
		if result.Outcome != resultmodel.OutcomeRolledBack {
			t.Fatalf("HEAD did not advance, so outcome = %q, want %q (failure %#v)",
				result.Outcome, resultmodel.OutcomeRolledBack, result.Failure)
		}
		if contents := readFile(t, repositoryRoot, "tracked.txt"); contents != "original\n" {
			t.Errorf("aborted commit did not roll back to the preimage: %q", contents)
		}
		return
	}
	if result.Outcome != resultmodel.OutcomeRisk {
		t.Fatalf("HEAD advanced to %s, so outcome = %q, want %q (failure %#v)",
			head, result.Outcome, resultmodel.OutcomeRisk, result.Failure)
	}
	if result.Failure == nil || result.Failure.Kind != FailureCommittedRisk {
		t.Fatalf("failure = %#v, want kind %q", result.Failure, FailureCommittedRisk)
	}
	if result.CommitSHA != head {
		t.Errorf("commit SHA = %q, want the advanced HEAD %q", result.CommitSHA, head)
	}
	if len(result.RevertArgv) != 3 || result.RevertArgv[2] != head {
		t.Errorf("revert argv = %v, want git revert %s", result.RevertArgv, head)
	}
	if contents := readFile(t, repositoryRoot, "tracked.txt"); contents != "updated\n" {
		t.Errorf("committed bytes were unwound as a preimage restore: %q", contents)
	}
}

func awaitRecordedPID(t *testing.T, markerPath string) int {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		contents, err := os.ReadFile(markerPath)
		if err == nil {
			if processID, convertError := strconv.Atoi(strings.TrimSpace(string(contents))); convertError == nil && processID > 0 {
				return processID
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("no process id was recorded at %s", markerPath)
	return 0
}
