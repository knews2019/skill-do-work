//go:build !aix && !darwin && !dragonfly && !freebsd && !ios && !linux && !windows

package main

import (
	"fmt"
	"runtime"
)

// measureDiskSpace has no free-space call to make here, so the probe reports
// itself skipped rather than reading as checked and clean.
func measureDiskSpace(directory string) (diskSpaceMeasurement, error) {
	return diskSpaceMeasurement{}, fmt.Errorf("%w (%s)", errDiskSpaceUnsupported, runtime.GOOS)
}
