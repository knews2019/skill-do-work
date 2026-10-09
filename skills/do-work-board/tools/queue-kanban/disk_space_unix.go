//go:build aix || darwin || dragonfly || freebsd || ios || linux

package main

import (
	"fmt"
	"syscall"
)

// measureDiskSpace reads the filesystem holding directory with one Statfs call.
// The tag list is narrower than atomic_replace_unix.go's on purpose: illumos,
// solaris and netbsd have no syscall.Statfs, and openbsd's Statfs_t names its
// fields differently, so those fall to disk_space_unsupported.go.
//
// Free space is Bavail (blocks this user may write), not Bfree (which counts the
// root-reserved blocks a builder cannot use). Field types differ per OS, so
// every conversion happens here.
func measureDiskSpace(directory string) (diskSpaceMeasurement, error) {
	var filesystemStatus syscall.Statfs_t
	if statfsError := syscall.Statfs(directory, &filesystemStatus); statfsError != nil {
		return diskSpaceMeasurement{}, fmt.Errorf("statfs: %w", statfsError)
	}

	// Bavail is signed on some platforms and may go negative when the
	// root-reserved blocks are in use; that is zero free for this user.
	var freeBytes uint64
	if availableBlocks := int64(filesystemStatus.Bavail); availableBlocks > 0 {
		freeBytes = uint64(availableBlocks) * uint64(filesystemStatus.Bsize)
	}
	return diskSpaceMeasurement{
		freeBytes:  freeBytes,
		totalBytes: uint64(filesystemStatus.Blocks) * uint64(filesystemStatus.Bsize),
	}, nil
}
