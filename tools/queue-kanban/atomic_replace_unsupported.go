//go:build !aix && !darwin && !dragonfly && !freebsd && !illumos && !ios && !linux && !netbsd && !openbsd && !solaris && !windows

package main

import (
	"fmt"
	"runtime"
)

// replaceFileAtomically fails closed where the project has no platform-specific
// primitive whose atomic replacement contract it can substantiate.
func replaceFileAtomically(temporaryPath string, destinationPath string) error {
	return fmt.Errorf("atomic file replacement is unsupported on %s", runtime.GOOS)
}
