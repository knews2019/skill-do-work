// Package atomicfile provides complete same-directory replacement, best-effort
// pre-publish change detection, and exclusive creation for durable files.
package atomicfile

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	// ErrUnsafeTarget identifies a symlink, directory, or special target.
	ErrUnsafeTarget = errors.New("unsafe atomic-file target")
	// ErrTargetChanged identifies a target change observed before publication.
	ErrTargetChanged = errors.New("atomic-file target changed")
)

// beforeAtomicReplace is a test seam for deterministic pre-publish changes.
var beforeAtomicReplace = func(string) {}

// ReplaceExisting publishes either the complete old file or the complete new
// file using the platform's same-directory atomic replacement primitive. It
// refuses changes observed while validating the target, but portable standard
// library APIs cannot compare-and-swap against a non-cooperating writer after
// the final validation and before publication.
func ReplaceExisting(filePath string, fileContents []byte) error {
	originalInfo, statError := os.Lstat(filePath)
	if statError != nil {
		return fmt.Errorf("checking atomic-file target %s: %w", filePath, statError)
	}
	if !originalInfo.Mode().IsRegular() {
		return fmt.Errorf("%w: %s is not a regular file", ErrUnsafeTarget, filePath)
	}
	originalDigest, digestError := regularFileDigest(filePath, originalInfo)
	if digestError != nil {
		return digestError
	}

	parentDirectory := filepath.Dir(filePath)
	temporaryFile, createError := os.CreateTemp(parentDirectory, "."+filepath.Base(filePath)+".tmp-*")
	if createError != nil {
		return fmt.Errorf("creating temporary file beside %s: %w", filePath, createError)
	}
	temporaryPath := temporaryFile.Name()
	defer os.Remove(temporaryPath)

	if _, writeError := temporaryFile.Write(fileContents); writeError != nil {
		temporaryFile.Close()
		return fmt.Errorf("writing temporary file for %s: %w", filePath, writeError)
	}
	if chmodError := temporaryFile.Chmod(originalInfo.Mode().Perm()); chmodError != nil {
		temporaryFile.Close()
		return fmt.Errorf("preserving permissions for %s: %w", filePath, chmodError)
	}
	if syncError := temporaryFile.Sync(); syncError != nil {
		temporaryFile.Close()
		return fmt.Errorf("syncing temporary file for %s: %w", filePath, syncError)
	}
	if closeError := temporaryFile.Close(); closeError != nil {
		return fmt.Errorf("closing temporary file for %s: %w", filePath, closeError)
	}

	beforeAtomicReplace(filePath)
	currentInfo, currentStatError := os.Lstat(filePath)
	if currentStatError != nil {
		return fmt.Errorf("%w: rechecking %s: %v", ErrTargetChanged, filePath, currentStatError)
	}
	if !currentInfo.Mode().IsRegular() || !os.SameFile(originalInfo, currentInfo) {
		return fmt.Errorf("%w: %s", ErrTargetChanged, filePath)
	}
	currentDigest, currentDigestError := regularFileDigest(filePath, originalInfo)
	if currentDigestError != nil {
		return currentDigestError
	}
	if currentDigest != originalDigest {
		return fmt.Errorf("%w: contents of %s changed", ErrTargetChanged, filePath)
	}
	if replaceError := replaceAtomicFile(temporaryPath, filePath); replaceError != nil {
		return fmt.Errorf("replacing %s: %w", filePath, replaceError)
	}
	return nil
}

// CreateExclusive creates a new regular file and fails when the path exists.
func CreateExclusive(filePath string, fileContents []byte, fileMode fs.FileMode) error {
	parentDirectory := filepath.Dir(filePath)
	directoryRoot, rootError := os.OpenRoot(parentDirectory)
	if rootError != nil {
		return fmt.Errorf("opening exclusive-file directory %s: %w", parentDirectory, rootError)
	}
	defer directoryRoot.Close()
	return CreateExclusiveAt(directoryRoot, filepath.Base(filePath), fileContents, fileMode)
}

// CreateExclusiveAt creates a file relative to an already-open rooted directory.
func CreateExclusiveAt(directoryRoot *os.Root, fileName string, fileContents []byte, fileMode fs.FileMode) error {
	if directoryRoot == nil {
		return fmt.Errorf("rooted directory is required")
	}
	createdFile, createError := directoryRoot.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, fileMode.Perm())
	if createError != nil {
		return fmt.Errorf("creating exclusive file %s: %w", fileName, createError)
	}
	keepFile := false
	defer func() {
		if !keepFile {
			directoryRoot.Remove(fileName)
		}
	}()
	if _, writeError := createdFile.Write(fileContents); writeError != nil {
		createdFile.Close()
		return fmt.Errorf("writing exclusive file %s: %w", fileName, writeError)
	}
	if syncError := createdFile.Sync(); syncError != nil {
		createdFile.Close()
		return fmt.Errorf("syncing exclusive file %s: %w", fileName, syncError)
	}
	if closeError := createdFile.Close(); closeError != nil {
		return fmt.Errorf("closing exclusive file %s: %w", fileName, closeError)
	}
	keepFile = true
	return nil
}

func regularFileDigest(filePath string, expectedInfo fs.FileInfo) ([sha256.Size]byte, error) {
	var emptyDigest [sha256.Size]byte
	targetFile, openError := os.Open(filePath)
	if openError != nil {
		return emptyDigest, fmt.Errorf("%w: opening %s for change detection: %v", ErrTargetChanged, filePath, openError)
	}
	defer targetFile.Close()
	beforeInfo, beforeStatError := targetFile.Stat()
	if beforeStatError != nil || !beforeInfo.Mode().IsRegular() || !os.SameFile(expectedInfo, beforeInfo) {
		return emptyDigest, fmt.Errorf("%w: %s changed before content verification", ErrTargetChanged, filePath)
	}
	digestHash := sha256.New()
	if _, copyError := io.Copy(digestHash, targetFile); copyError != nil {
		return emptyDigest, fmt.Errorf("hashing atomic-file target %s: %w", filePath, copyError)
	}
	afterInfo, afterStatError := targetFile.Stat()
	if afterStatError != nil || !os.SameFile(beforeInfo, afterInfo) || beforeInfo.Size() != afterInfo.Size() || !beforeInfo.ModTime().Equal(afterInfo.ModTime()) {
		return emptyDigest, fmt.Errorf("%w: %s changed during content verification", ErrTargetChanged, filePath)
	}
	var contentDigest [sha256.Size]byte
	copy(contentDigest[:], digestHash.Sum(nil))
	return contentDigest, nil
}
