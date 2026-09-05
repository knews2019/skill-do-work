//go:build windows

package nextselection

import (
	"fmt"
	"io"
	"time"
)

// runOwnedProbe never launches anything here, so it states the launch status and
// an explicitly false pair of execution facts.
func runOwnedProbe(_ string, _ []byte, _ time.Duration, _ io.Writer) (BlockedProbeEvidence, error) {
	return BlockedProbeEvidence{ExitStatus: BlockedProbeLaunchStatus, Launched: false, TimedOut: false},
		fmt.Errorf("standard-library process-tree ownership is unavailable on windows")
}
