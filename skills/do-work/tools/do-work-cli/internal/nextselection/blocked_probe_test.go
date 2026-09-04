//go:build unix

package nextselection

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestBlockedProbePreservesRawStatus(t *testing.T) {
	status, err := RunBlockedProbe([]byte("exit 37"), 2)
	if err != nil || status != 37 {
		t.Fatalf("status=%d err=%v", status, err)
	}
}

func TestBlockedProbeEvidenceBoundsAndNormalizesDiagnostics(t *testing.T) {
	repositoryRoot := t.TempDir()
	probe := fmt.Sprintf("printf '%s\\r\\n'; yes x | head -c 70000; exit 29", repositoryRoot)
	evidence, err := RunBlockedProbeEvidenceAtRoot(repositoryRoot, []byte(probe), 2)
	if err != nil {
		t.Fatal(err)
	}
	if evidence.ExitStatus != 29 || !evidence.Launched || evidence.TimedOut {
		t.Fatalf("evidence=%#v", evidence)
	}
	if strings.Contains(evidence.Diagnostic, repositoryRoot) || !strings.Contains(evidence.Diagnostic, "<repo-root>") || !strings.Contains(evidence.Diagnostic, "[diagnostic truncated]") {
		t.Fatalf("diagnostic was not normalized and bounded: %q", evidence.Diagnostic)
	}
	if evidence.DiagnosticSHA256 == "" || len(evidence.Diagnostic) > blockedProbeDiagnosticLimit+64 {
		t.Fatalf("diagnostic evidence=%#v", evidence)
	}
}
func TestBlockedProbeTimeoutKillsDescendantGroup(t *testing.T) {
	directory := t.TempDir()
	pidPath := filepath.Join(directory, "child.pid")
	probe := fmt.Sprintf("(trap '' TERM; sleep 30) & echo $! > %q; wait", pidPath)
	status, err := RunBlockedProbe([]byte(probe), 1)
	if err != nil || status != BlockedProbeTimeoutStatus {
		t.Fatalf("status=%d err=%v", status, err)
	}
	pidBytes, readError := os.ReadFile(pidPath)
	if readError != nil {
		t.Fatal(readError)
	}
	pid, _ := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if syscall.Kill(pid, 0) != nil {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("descendant %d survived timeout", pid)
}

func TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits(t *testing.T) {
	directory := t.TempDir()
	pidPath := filepath.Join(directory, "child.pid")
	probe := fmt.Sprintf("sleep 30 & echo $! > %q; exit 0", pidPath)
	status, err := RunBlockedProbe([]byte(probe), 2)
	if err != nil || status != 0 {
		t.Fatalf("status=%d err=%v", status, err)
	}
	pidBytes, readError := os.ReadFile(pidPath)
	if readError != nil {
		t.Fatal(readError)
	}
	pid, _ := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if syscall.Kill(pid, 0) != nil {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("background descendant %d survived successful probe", pid)
}

func TestBlockedProbeInterruptionIsTypedAndReapsDescendants(t *testing.T) {
	directory := t.TempDir()
	pidPath := filepath.Join(directory, "child.pid")
	probe := fmt.Sprintf("sleep 30 & echo $! > %q; wait", pidPath)
	guardSignals := make(chan os.Signal, 1)
	signal.Notify(guardSignals, syscall.SIGINT)
	defer signal.Stop(guardSignals)
	type probeOutcome struct {
		status int
		err    error
	}
	outcomes := make(chan probeOutcome, 1)
	go func() {
		status, err := runBlockedProbeFixture(directory, []byte(probe), 3)
		outcomes <- probeOutcome{status: status, err: err}
	}()

	var pid int
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		pidBytes, err := os.ReadFile(pidPath)
		if err == nil {
			pid, _ = strconv.Atoi(strings.TrimSpace(string(pidBytes)))
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if pid == 0 {
		t.Fatal("probe descendant did not start")
	}
	// The guard above makes an early delivery safe on the RED implementation;
	// the short delay lets that implementation reach its late signal.Notify.
	time.Sleep(100 * time.Millisecond)
	if err := syscall.Kill(os.Getpid(), syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	var outcome probeOutcome
	select {
	case outcome = <-outcomes:
	case <-time.After(5 * time.Second):
		t.Fatal("interrupted probe did not return")
	}
	var interruption interface{ InterruptionExitStatus() int }
	if outcome.status != 130 || !errors.As(outcome.err, &interruption) || interruption.InterruptionExitStatus() != 130 {
		t.Fatalf("status=%d err=%T %v, want typed interruption 130", outcome.status, outcome.err, outcome.err)
	}
	deadline = time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if syscall.Kill(pid, 0) != nil {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("descendant %d survived interruption", pid)
}

func runBlockedProbeFixture(repositoryRoot string, probeBytes []byte, timeoutSeconds int) (int, error) {
	return RunBlockedProbeAtRoot(repositoryRoot, probeBytes, timeoutSeconds)
}
