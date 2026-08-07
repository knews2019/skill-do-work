//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var replaceFileProcedure = syscall.NewLazyDLL("kernel32.dll").NewProc("ReplaceFileW")

// replaceFileAtomically uses Windows' single-file document replacement API.
// Both paths are in the same directory, satisfying ReplaceFileW's same-volume
// requirement without falling back to Go's non-atomic Windows Rename behavior.
func replaceFileAtomically(temporaryPath string, destinationPath string) error {
	destinationPointer, destinationError := syscall.UTF16PtrFromString(destinationPath)
	if destinationError != nil {
		return fmt.Errorf("encoding destination path: %w", destinationError)
	}
	temporaryPointer, temporaryError := syscall.UTF16PtrFromString(temporaryPath)
	if temporaryError != nil {
		return fmt.Errorf("encoding replacement path: %w", temporaryError)
	}

	result, _, callError := replaceFileProcedure.Call(
		uintptr(unsafe.Pointer(destinationPointer)),
		uintptr(unsafe.Pointer(temporaryPointer)),
		0,
		0,
		0,
		0,
	)
	if result != 0 {
		return nil
	}
	if callError == syscall.Errno(0) {
		return fmt.Errorf("ReplaceFileW failed without an error code")
	}
	return callError
}
