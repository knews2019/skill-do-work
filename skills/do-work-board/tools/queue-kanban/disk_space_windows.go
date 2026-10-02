//go:build windows

package main

import (
	"fmt"
	"hash/fnv"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var getDiskFreeSpaceProcedure = syscall.NewLazyDLL("kernel32.dll").NewProc("GetDiskFreeSpaceExW")

// measureDiskSpace asks Windows for the free and total space of the volume
// holding directory. Free space is the caller's available bytes, which honor
// disk quotas, not the volume's raw free count.
//
// The device identity hashes the upper-cased volume name (`C:`, or a UNC share),
// so a volume mounted in a folder below a drive letter is not told apart from
// that drive and may be reported once for both.
func measureDiskSpace(directory string) (diskSpaceMeasurement, error) {
	absoluteDirectory, absoluteError := filepath.Abs(directory)
	if absoluteError != nil {
		return diskSpaceMeasurement{}, absoluteError
	}
	directoryPointer, encodeError := syscall.UTF16PtrFromString(absoluteDirectory)
	if encodeError != nil {
		return diskSpaceMeasurement{}, fmt.Errorf("encoding directory path: %w", encodeError)
	}

	var availableBytes, totalBytes, volumeFreeBytes uint64
	result, _, callError := getDiskFreeSpaceProcedure.Call(
		uintptr(unsafe.Pointer(directoryPointer)),
		uintptr(unsafe.Pointer(&availableBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&volumeFreeBytes)),
	)
	if result == 0 {
		if callError == syscall.Errno(0) {
			return diskSpaceMeasurement{}, fmt.Errorf("GetDiskFreeSpaceExW failed without an error code")
		}
		return diskSpaceMeasurement{}, callError
	}

	volumeHash := fnv.New64a()
	volumeHash.Write([]byte(strings.ToUpper(filepath.VolumeName(absoluteDirectory))))
	return diskSpaceMeasurement{
		freeBytes:      availableBytes,
		totalBytes:     totalBytes,
		deviceIdentity: volumeHash.Sum64(),
	}, nil
}
