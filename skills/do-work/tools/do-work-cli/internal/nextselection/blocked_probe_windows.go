//go:build windows

package nextselection

import (
	"fmt"
	"time"
)

func runOwnedProbe(_ []byte, _ time.Duration) (int, error) {
	return BlockedProbeLaunchStatus, fmt.Errorf("standard-library process-tree ownership is unavailable on windows")
}
