package nextselection

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	BlockedProbeTimeoutStatus = 124
	BlockedProbeLaunchStatus  = 125
)

// BlockedProbeInterruption distinguishes an invocation signal from a probe's
// own non-zero status so queue selection can stop instead of excluding one REQ.
type BlockedProbeInterruption struct {
	ExitStatus int
}

func (interruption BlockedProbeInterruption) Error() string {
	return fmt.Sprintf("blocked probe interrupted with exit %d", interruption.ExitStatus)
}

func (interruption BlockedProbeInterruption) InterruptionExitStatus() int {
	return interruption.ExitStatus
}

// RunBlockedProbe executes one materialized probe while the platform-specific
// implementation owns its complete descendant tree. Raw legacy status is returned
// as evidence; public command rendering remains in the standard 0-4 envelope.
func RunBlockedProbe(probeBytes []byte, timeoutSeconds int) (int, error) {
	repositoryRoot, err := os.Getwd()
	if err != nil {
		return BlockedProbeLaunchStatus, err
	}
	return RunBlockedProbeAtRoot(repositoryRoot, probeBytes, timeoutSeconds)
}

// RunBlockedProbeAtRoot executes one materialized probe relative to the selected repository.
func RunBlockedProbeAtRoot(repositoryRoot string, probeBytes []byte, timeoutSeconds int) (int, error) {
	if repositoryRoot == "" {
		return BlockedProbeLaunchStatus, fmt.Errorf("repository root is empty")
	}
	if timeoutSeconds <= 0 {
		return BlockedProbeLaunchStatus, fmt.Errorf("timeout must be positive")
	}
	if len(probeBytes) == 0 {
		return BlockedProbeLaunchStatus, fmt.Errorf("probe is empty")
	}
	return runOwnedProbe(repositoryRoot, probeBytes, time.Duration(timeoutSeconds)*time.Second)
}

func blockedProbeInterruptionStatus(err error) (int, bool) {
	var interruption interface {
		error
		InterruptionExitStatus() int
	}
	if !errors.As(err, &interruption) {
		return 0, false
	}
	status := interruption.InterruptionExitStatus()
	return status, status == 129 || status == 130 || status == 143
}
