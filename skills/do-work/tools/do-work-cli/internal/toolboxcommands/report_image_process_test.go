//go:build !windows

package toolboxcommands

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestOwnedProcessCancellationTerminatesAndReaps(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	oldGrace := reportImageGracePeriod
	reportImageGracePeriod = 100 * time.Millisecond
	t.Cleanup(func() { reportImageGracePeriod = oldGrace })
	done := make(chan ownedProcessResult, 1)
	go func() { done <- runOwnedProcess(ctx, "", "sh", "-c", "trap '' TERM; sleep 30 & wait") }()
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case result := <-done:
		if !result.Interrupted {
			t.Fatalf("result=%+v", result)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("owned process tree was not terminated and reaped")
	}
}

func TestRemediationCancellationReachesMediaGitCommitAndRollback(t *testing.T) {
	repository := toolboxTestRepository(t)
	target := filepath.Join(repository, "image.png")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	toolboxTestGit(t, repository, "add", "image.png")
	toolboxTestGit(t, repository, "commit", "-m", "baseline")
	hooks := strings.TrimSpace(toolboxTestGit(t, repository, "rev-parse", "--git-path", "hooks"))
	if !filepath.IsAbs(hooks) {
		hooks = filepath.Join(repository, hooks)
	}
	marker := filepath.Join(repository, "hook.pid")
	hook := "#!/bin/sh\necho $$ > " + marker + "\ntrap '' TERM\nsleep 30\n"
	if err := os.WriteFile(filepath.Join(hooks, "pre-commit"), []byte(hook), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan resultmodel.CommandResult, 1)
	go func() {
		done <- runTransactionContext(ctx, CommandReportImage, repository, []string{"image.png"}, nil, false, true, "media", func(recorder *gittransaction.MutationRecorder) error {
			if err := rootedPublishFile(repository, "image.png", []byte("new"), 0o644, true); err != nil {
				return err
			}
			return recorder.RecordTouched("image.png")
		})
	}()
	deadline := time.Now().Add(3 * time.Second)
	for {
		if _, err := os.Stat(marker); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("commit hook did not start")
		}
		time.Sleep(20 * time.Millisecond)
	}
	cancel()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("cancelled media Git transaction did not return")
	}
	pidBytes, _ := os.ReadFile(marker)
	pid, _ := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	if pid > 0 && syscall.Kill(pid, 0) == nil {
		_ = syscall.Kill(pid, syscall.SIGKILL)
		t.Fatal("media commit hook survived cancellation")
	}
	if contents, _ := os.ReadFile(target); string(contents) != "old" {
		t.Fatalf("cancelled media transaction did not roll back: %q", contents)
	}
}

func TestRemediationOwnedProcessDoesNotLaunchAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := runOwnedProcess(ctx, "", "sh", "-c", "exit 99")
	if !result.Interrupted {
		t.Fatalf("pre-cancelled launch result=%+v", result)
	}
}

func TestRemediationLeaderExitStillKillsTermDeafDescendant(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	oldGrace := reportImageGracePeriod
	reportImageGracePeriod = 100 * time.Millisecond
	t.Cleanup(func() { reportImageGracePeriod = oldGrace })
	done := make(chan ownedProcessResult, 1)
	go func() {
		done <- runOwnedProcess(ctx, "", "sh", "-c", "trap 'exit 0' TERM; (trap '' TERM; sleep 30) & wait")
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case result := <-done:
		if !result.Interrupted {
			t.Fatalf("result=%+v", result)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("TERM-deaf descendant survived leader exit")
	}
}
