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
	waitForDescendantToDisappear(t, pid)
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
	waitForDescendantToDisappear(t, pid)
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
	waitForDescendantToDisappear(t, pid)
}

// descendantReapBudget bounds the wait for a killed descendant to disappear. Its
// own parent is already gone by then, so init does the reaping and a zombie still
// satisfies kill(pid, 0) until it does. On a loaded machine that has been measured
// close to two seconds, so the budget is deliberately generous: it proves the
// descendant does not survive, and it is not a latency assertion.
const descendantReapBudget = 10 * time.Second

func waitForDescendantToDisappear(t *testing.T, pid int) {
	t.Helper()
	deadline := time.Now().Add(descendantReapBudget)
	for time.Now().Before(deadline) {
		if syscall.Kill(pid, 0) != nil {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("descendant %d survived %s", pid, descendantReapBudget)
}

func runBlockedProbeFixture(repositoryRoot string, probeBytes []byte, timeoutSeconds int) (int, error) {
	return RunBlockedProbeAtRoot(repositoryRoot, probeBytes, timeoutSeconds)
}

// TestBlockedProbeEvidencePreservesOrdinaryReservedExitValues pins the REQ-506
// defect from the other side: 124 and 125 are ordinary exit values a probe may
// choose for itself, so evidence must report them as a launched, non-timed-out
// completion instead of rebuilding "timed out" or "never launched" from the
// integer the child happened to pick.
func TestBlockedProbeEvidencePreservesOrdinaryReservedExitValues(t *testing.T) {
	for _, reservedStatus := range []int{BlockedProbeTimeoutStatus, BlockedProbeLaunchStatus} {
		t.Run(fmt.Sprintf("exit %d", reservedStatus), func(t *testing.T) {
			repositoryRoot := t.TempDir()
			evidence, err := RunBlockedProbeEvidenceAtRoot(repositoryRoot, []byte(fmt.Sprintf("exit %d", reservedStatus)), 5)
			if err != nil || evidence.ExitStatus != reservedStatus || !evidence.Launched || evidence.TimedOut {
				t.Fatalf("evidence=%#v err=%v", evidence, err)
			}
		})
	}
}

// TestBlockedProbeEvidenceRefusesUnrunnableInputsAsUnlaunched pins the input
// guards that return before any process exists. They report the launch status
// today only because false is the zero value of both booleans; once launch and
// timeout are set from observation these refusals must keep stating the same
// facts explicitly.
func TestBlockedProbeEvidenceRefusesUnrunnableInputsAsUnlaunched(t *testing.T) {
	for _, test := range []struct {
		name           string
		repositoryRoot string
		probeBytes     []byte
		timeoutSeconds int
	}{
		{name: "empty repository root", repositoryRoot: "", probeBytes: []byte("exit 0"), timeoutSeconds: 2},
		{name: "nonpositive timeout", repositoryRoot: t.TempDir(), probeBytes: []byte("exit 0"), timeoutSeconds: 0},
		{name: "empty probe", repositoryRoot: t.TempDir(), probeBytes: nil, timeoutSeconds: 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			evidence, err := RunBlockedProbeEvidenceAtRoot(test.repositoryRoot, test.probeBytes, test.timeoutSeconds)
			if err == nil || evidence.ExitStatus != BlockedProbeLaunchStatus || evidence.Launched || evidence.TimedOut {
				t.Fatalf("evidence=%#v err=%v", evidence, err)
			}
		})
	}
}
