package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// writeFileAtomically writes fileContents to a synced temporary file beside an
// existing regular target, then hands the final replacement to the platform
// primitive. A reader or process crash therefore observes the complete old file
// or the complete new file, never a truncate-then-write intermediate state.
func writeFileAtomically(filePath string, fileContents []byte) error {
	originalInfo, lstatError := os.Lstat(filePath)
	if lstatError != nil {
		return fmt.Errorf("checking atomic-write target %s: %w", filePath, lstatError)
	}
	if !originalInfo.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file — refusing to replace a symlink or special file", filePath)
	}

	parentDirectory := filepath.Dir(filePath)
	temporaryFile, createError := os.CreateTemp(parentDirectory, "."+filepath.Base(filePath)+".tmp-*")
	if createError != nil {
		return fmt.Errorf("creating temp file in %s: %w", parentDirectory, createError)
	}
	temporaryPath := temporaryFile.Name()
	defer os.Remove(temporaryPath) // no-op once the replacement has landed

	if _, writeError := temporaryFile.Write(fileContents); writeError != nil {
		temporaryFile.Close()
		return fmt.Errorf("writing %s: %w", temporaryPath, writeError)
	}
	// CreateTemp opens 0600; retain the permission behavior of writing the
	// existing file directly. Platform replacement may preserve more metadata.
	if chmodError := temporaryFile.Chmod(originalInfo.Mode().Perm()); chmodError != nil {
		temporaryFile.Close()
		return fmt.Errorf("setting mode on %s: %w", temporaryPath, chmodError)
	}
	if syncError := temporaryFile.Sync(); syncError != nil {
		temporaryFile.Close()
		return fmt.Errorf("syncing %s: %w", temporaryPath, syncError)
	}
	if closeError := temporaryFile.Close(); closeError != nil {
		return fmt.Errorf("closing %s: %w", temporaryPath, closeError)
	}

	currentInfo, currentLstatError := os.Lstat(filePath)
	if currentLstatError != nil {
		return fmt.Errorf("rechecking atomic-write target %s: %w", filePath, currentLstatError)
	}
	if !currentInfo.Mode().IsRegular() || !os.SameFile(originalInfo, currentInfo) {
		return fmt.Errorf("%s changed before atomic replacement — refusing to overwrite it", filePath)
	}
	if replaceError := replaceFileAtomically(temporaryPath, filePath); replaceError != nil {
		return fmt.Errorf("replacing %s: %w", filePath, replaceError)
	}
	return nil
}
