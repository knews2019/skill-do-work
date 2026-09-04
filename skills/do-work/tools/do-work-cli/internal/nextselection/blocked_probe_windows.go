//go:build windows

package nextselection

import (
	"fmt"
	"io"
	"time"
)

func runOwnedProbe(_ string, _ []byte, _ time.Duration, _ io.Writer) (int, error) {
	return BlockedProbeLaunchStatus, fmt.Errorf("standard-library process-tree ownership is unavailable on windows")
}
