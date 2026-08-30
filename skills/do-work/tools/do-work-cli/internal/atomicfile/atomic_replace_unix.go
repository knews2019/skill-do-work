//go:build aix || darwin || dragonfly || freebsd || illumos || ios || linux || netbsd || openbsd || solaris

package atomicfile

import "os"

func replaceAtomicFile(temporaryPath string, destinationPath string) error {
	return os.Rename(temporaryPath, destinationPath)
}
