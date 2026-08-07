//go:build aix || darwin || dragonfly || freebsd || illumos || ios || linux || netbsd || openbsd || solaris

package main

import "os"

// replaceFileAtomically uses the same-directory atomic rename guaranteed by
// Unix filesystems for this existing-file replacement.
func replaceFileAtomically(temporaryPath string, destinationPath string) error {
	return os.Rename(temporaryPath, destinationPath)
}
