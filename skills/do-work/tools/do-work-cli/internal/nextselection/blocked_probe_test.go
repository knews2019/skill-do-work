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

// leakedDescendantBody is the background descendant every descendant-cleanup test
// below plants. It releases its inherited stdout and stderr before sleeping,
// because a descendant that keeps them holds the runner's own diagnostic pipe
// open: the runner would then be unable to return until that descendant had
// already exited, and no assertion could ever observe a leak. It sleeps far past
// descendantReapBudget so a descendant the probe failed to kill is still running
// when the poll looks for it.
const leakedDescendantBody = "exec 1>&- 2>&-; sleep 30"

// probeLeaderHoldSeconds keeps the probe's shell leader alive past the runner's
// own timeout and past the interrupt below, which is what makes those branches
// fire at all. When process-group teardown is broken nothing else ends the
// leader, so this also bounds how long a broken runner can block before the
// descendant assertion gets to run.
const probeLeaderHoldSeconds = 4

func TestBlockedProbeTimeoutKillsDescendantGroup(t *testing.T) {
	directory := t.TempDir()
	pidPath := filepath.Join(directory, "child.pid")
	probe := fmt.Sprintf("(trap '' TERM; %s) & echo $! > %q; sleep %d", leakedDescendantBody, pidPath, probeLeaderHoldSeconds)
	status, err := RunBlockedProbe([]byte(probe), 1)
	if err != nil || status != BlockedProbeTimeoutStatus {
		t.Fatalf("status=%d err=%v", status, err)
	}
	waitForDescendantToDisappear(t, waitForDescendantPid(t, pidPath))
}

func TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits(t *testing.T) {
	directory := t.TempDir()
	pidPath := filepath.Join(directory, "child.pid")
	probe := fmt.Sprintf("(%s) & echo $! > %q; exit 0", leakedDescendantBody, pidPath)
	status, err := RunBlockedProbe([]byte(probe), 2)
	if err != nil || status != 0 {
		t.Fatalf("status=%d err=%v", status, err)
	}
	waitForDescendantToDisappear(t, waitForDescendantPid(t, pidPath))
}

func TestBlockedProbeInterruptionIsTypedAndReapsDescendants(t *testing.T) {
	directory := t.TempDir()
	pidPath := filepath.Join(directory, "child.pid")
	probe := fmt.Sprintf("(%s) & echo $! > %q; sleep %d", leakedDescendantBody, pidPath, probeLeaderHoldSeconds)
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

	pid := waitForDescendantPid(t, pidPath)
	// The guard above makes an early delivery safe on the RED implementation;
	// the short delay lets that implementation reach its late signal.Notify.
	time.Sleep(100 * time.Millisecond)
	if err := syscall.Kill(os.Getpid(), syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	var outcome probeOutcome
	select {
	case outcome = <-outcomes:
	case <-time.After(interruptedProbeReturnBudget):
		t.Fatal("interrupted probe did not return")
	}
	var interruption interface{ InterruptionExitStatus() int }
	if outcome.status != 130 || !errors.As(outcome.err, &interruption) || interruption.InterruptionExitStatus() != 130 {
		t.Fatalf("status=%d err=%T %v, want typed interruption 130", outcome.status, outcome.err, outcome.err)
	}
	waitForDescendantToDisappear(t, pid)
}

// interruptedProbeReturnBudget bounds the interrupted runner's return. It must
// outlast probeLeaderHoldSeconds so that a runner which fails to tear the process
// group down still returns and reaches the descendant assertion, instead of
// failing here on a bound that names no surviving process.
const interruptedProbeReturnBudget = 15 * time.Second

// descendantReapBudget bounds how long init may take to reap a descendant the
// probe already killed: the descendant's own parent is gone by then, and a zombie
// still satisfies kill(pid, 0) until init collects it. That reap has been measured
// close to two seconds on a loaded machine, so the budget is generous. It is not
// what proves the kill — leakedDescendantBody sleeps far longer than this budget,
// so a pid still answering kill(pid, 0) at the deadline names a process that is
// genuinely running rather than one waiting to be reaped.
const descendantReapBudget = 10 * time.Second

// waitForDescendantPid reads the pid the probe recorded, waiting for the file when
// the probe is still running. It kills whatever the pid names once the test ends,
// so an assertion that fails on a real leak does not leave the leaked process
// running on the machine.
func waitForDescendantPid(t *testing.T, pidPath string) int {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for {
		pidBytes, readError := os.ReadFile(pidPath)
		if readError == nil {
			if pid, parseError := strconv.Atoi(strings.TrimSpace(string(pidBytes))); parseError == nil && pid > 0 {
				t.Cleanup(func() { _ = syscall.Kill(pid, syscall.SIGKILL) })
				return pid
			}
		}
		if !time.Now().Before(deadline) {
			t.Fatalf("probe descendant pid never appeared at %q", pidPath)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

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
