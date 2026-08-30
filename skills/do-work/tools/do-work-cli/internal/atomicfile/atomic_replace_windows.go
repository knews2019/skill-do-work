//go:build windows

package atomicfile

import (
	"fmt"
	"syscall"
	"unsafe"
)

var replaceFileProcedure = syscall.NewLazyDLL("kernel32.dll").NewProc("ReplaceFileW")

func replaceAtomicFile(temporaryPath string, destinationPath string) error {
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
		0, 0, 0, 0,
	)
	if result != 0 {
		return nil
	}
	if callError == syscall.Errno(0) {
		return fmt.Errorf("ReplaceFileW failed without an error code")
	}
	return callError
}
