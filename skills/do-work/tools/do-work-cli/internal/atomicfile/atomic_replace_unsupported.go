//go:build !aix && !darwin && !dragonfly && !freebsd && !illumos && !ios && !linux && !netbsd && !openbsd && !solaris && !windows

package atomicfile

import (
	"fmt"
	"runtime"
)

func replaceAtomicFile(temporaryPath string, destinationPath string) error {
	return fmt.Errorf("atomic file replacement is unsupported on %s", runtime.GOOS)
}
