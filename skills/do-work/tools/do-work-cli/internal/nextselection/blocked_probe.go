package nextselection

import (
	"fmt"
	"time"
)

const (
	BlockedProbeTimeoutStatus = 124
	BlockedProbeLaunchStatus  = 125
)

// RunBlockedProbe executes one materialized probe while the platform-specific
// implementation owns its complete descendant tree. Raw legacy status is returned
// as evidence; public command rendering remains in the standard 0-4 envelope.
func RunBlockedProbe(probeBytes []byte, timeoutSeconds int) (int, error) {
	if timeoutSeconds <= 0 {
		return BlockedProbeLaunchStatus, fmt.Errorf("timeout must be positive")
	}
	if len(probeBytes) == 0 {
		return BlockedProbeLaunchStatus, fmt.Errorf("probe is empty")
	}
	return runOwnedProbe(probeBytes, time.Duration(timeoutSeconds)*time.Second)
}
