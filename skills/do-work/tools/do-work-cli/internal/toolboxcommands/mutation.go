package toolboxcommands

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func absentTransactionDirectories(repositoryRoot string, candidates ...string) []string {
	seen := map[string]struct{}{}
	var absent []string
	for _, candidate := range candidates {
		for candidate = filepath.ToSlash(filepath.Clean(candidate)); candidate != "."; candidate = filepath.ToSlash(filepath.Dir(candidate)) {
			if _, duplicate := seen[candidate]; duplicate {
				continue
			}
			seen[candidate] = struct{}{}
			if _, err := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(candidate))); os.IsNotExist(err) {
				absent = append(absent, candidate)
			}
		}
	}
	sort.Slice(absent, func(i, j int) bool {
		leftDepth := strings.Count(absent[i], "/")
		rightDepth := strings.Count(absent[j], "/")
		if leftDepth != rightDepth {
			return leftDepth < rightDepth
		}
		return absent[i] < absent[j]
	})
	return absent
}

func transactionResult(command string, transaction gittransaction.TransactionResult, exactText string) resultmodel.CommandResult {
	result := gittransaction.BuildCommandResult(command, transaction)
	if result.Outcome == resultmodel.OutcomeSuccess && exactText != "" {
		result.ExactTextOutput = &exactText
	}
	return result
}

func runTransaction(command, repositoryRoot string, targets, directories []string, dryRun, commit bool, message string, mutate func(*gittransaction.MutationRecorder) error) resultmodel.CommandResult {
	return runTransactionContext(context.Background(), command, repositoryRoot, targets, directories, dryRun, commit, message, mutate)
}

func runTransactionContext(ctx context.Context, command, repositoryRoot string, targets, directories []string, dryRun, commit bool, message string, mutate func(*gittransaction.MutationRecorder) error) resultmodel.CommandResult {
	return runTransactionContextWithExisting(ctx, command, repositoryRoot, targets, nil, directories, dryRun, commit, message, mutate)
}

func runTransactionContextWithExisting(ctx context.Context, command, repositoryRoot string, targets, existingUntracked, directories []string, dryRun, commit bool, message string, mutate func(*gittransaction.MutationRecorder) error) resultmodel.CommandResult {
	transaction := gittransaction.ExecuteTransaction(ctx, gittransaction.TransactionOptions{
		RepositoryRoot: repositoryRoot, TargetPaths: targets, ExistingUntrackedTargetPaths: existingUntracked, CreatedDirectoryPaths: directories,
		DryRun: dryRun, Commit: commit, CommitMessage: message,
	}, mutate)
	return transactionResult(command, transaction, "")
}

// validateNoLinkedAncestors rejects linked path components while os.Root
// supplies the race-safe confinement boundary used by the actual operation.
func validateNoLinkedAncestors(repositoryRoot, relative string, includeLeaf bool) error {
	root, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return err
	}
	defer root.Close()
	clean := filepath.ToSlash(filepath.Clean(relative))
	parts := strings.Split(clean, "/")
	limit := len(parts)
	if !includeLeaf {
		limit--
	}
	for index := 1; index <= limit; index++ {
		component := filepath.FromSlash(strings.Join(parts[:index], "/"))
		info, statErr := root.Lstat(component)
		if os.IsNotExist(statErr) {
			continue
		}
		if statErr != nil {
			return fmt.Errorf("inspect confined path %s: %w", component, statErr)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("confined path contains a symbolic link: %s", component)
		}
		if index < len(parts) && !info.IsDir() {
			return fmt.Errorf("confined path ancestor is not a directory: %s", component)
		}
	}
	return nil
}

func rootedMkdirAll(repositoryRoot, relative string, mode os.FileMode) error {
	if err := validateNoLinkedAncestors(repositoryRoot, relative, false); err != nil {
		return err
	}
	root, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return err
	}
	defer root.Close()
	return root.MkdirAll(filepath.FromSlash(relative), mode)
}

func rootedMkdirExclusive(repositoryRoot, relative string, mode os.FileMode) (os.FileInfo, error) {
	if err := validateNoLinkedAncestors(repositoryRoot, relative, false); err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	if err := root.Mkdir(filepath.FromSlash(relative), mode); err != nil {
		return nil, err
	}
	return root.Lstat(filepath.FromSlash(relative))
}

// rootedPublishFile publishes through a parent handle pinned beneath the
// repository. Creation uses an exclusive hard-link claim; replacement checks
// the observed inode and bytes immediately before the final rename.
func rootedPublishFile(repositoryRoot, relative string, data []byte, mode os.FileMode, replace bool) error {
	if err := validateNoLinkedAncestors(repositoryRoot, relative, false); err != nil {
		return err
	}
	root, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return err
	}
	defer root.Close()
	parentName, leaf := filepath.Dir(filepath.FromSlash(relative)), filepath.Base(filepath.FromSlash(relative))
	parent, err := root.OpenRoot(parentName)
	if err != nil {
		return fmt.Errorf("open confined parent for %s: %w", relative, err)
	}
	defer parent.Close()
	var originalInfo os.FileInfo
	var originalBytes []byte
	if replace {
		originalInfo, err = parent.Lstat(leaf)
		if err != nil || !originalInfo.Mode().IsRegular() {
			return fmt.Errorf("replacement target is not a regular file: %s", relative)
		}
		originalBytes, err = parent.ReadFile(leaf)
		if err != nil {
			return err
		}
	}
	random := make([]byte, 12)
	if _, err := rand.Read(random); err != nil {
		return err
	}
	temporary := "." + leaf + ".publishing-" + hex.EncodeToString(random)
	handle, err := parent.OpenFile(temporary, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	keep := false
	defer func() {
		if !keep {
			_ = parent.Remove(temporary)
		}
	}()
	if _, err := handle.Write(data); err != nil {
		handle.Close()
		return err
	}
	if err := handle.Sync(); err != nil {
		handle.Close()
		return err
	}
	if err := handle.Close(); err != nil {
		return err
	}
	if replace {
		currentInfo, statErr := parent.Lstat(leaf)
		currentBytes, readErr := parent.ReadFile(leaf)
		if statErr != nil || readErr != nil || !os.SameFile(originalInfo, currentInfo) || !bytes.Equal(originalBytes, currentBytes) {
			return fmt.Errorf("target changed before confined replacement: %s", relative)
		}
		if err := parent.Rename(temporary, leaf); err != nil {
			return err
		}
		keep = true
		return nil
	}
	if err := parent.Link(temporary, leaf); err != nil {
		return err
	}
	if err := parent.Remove(temporary); err != nil {
		_ = parent.Remove(leaf)
		return err
	}
	keep = true
	return nil
}

func rootedReadFile(repositoryRoot, relative string) ([]byte, error) {
	if err := validateNoLinkedAncestors(repositoryRoot, relative, true); err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	handle, err := root.Open(filepath.FromSlash(relative))
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	info, err := handle.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("confined source is not a regular file: %s", relative)
	}
	return io.ReadAll(handle)
}

func rootedPublishInOwnedDirectory(repositoryRoot, directoryRelative string, ownedDirectory os.FileInfo, leaf string, data []byte, mode os.FileMode) error {
	root, err := os.OpenRoot(repositoryRoot)
	if err != nil {
		return err
	}
	defer root.Close()
	current, err := root.Lstat(filepath.FromSlash(directoryRelative))
	if err != nil || !os.SameFile(ownedDirectory, current) {
		return fmt.Errorf("owned publication directory changed: %s", directoryRelative)
	}
	directory, err := root.OpenRoot(filepath.FromSlash(directoryRelative))
	if err != nil {
		return err
	}
	defer directory.Close()
	opened, err := directory.Stat(".")
	if err != nil || !os.SameFile(ownedDirectory, opened) {
		return fmt.Errorf("owned publication directory changed while opening: %s", directoryRelative)
	}
	random := make([]byte, 12)
	if _, err := rand.Read(random); err != nil {
		return err
	}
	temporary := "." + leaf + ".publishing-" + hex.EncodeToString(random)
	handle, err := directory.OpenFile(temporary, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	if _, err := handle.Write(data); err != nil {
		handle.Close()
		_ = directory.Remove(temporary)
		return err
	}
	if err := handle.Sync(); err != nil {
		handle.Close()
		_ = directory.Remove(temporary)
		return err
	}
	if err := handle.Close(); err != nil {
		_ = directory.Remove(temporary)
		return err
	}
	if err := directory.Link(temporary, leaf); err != nil {
		_ = directory.Remove(temporary)
		return err
	}
	if err := directory.Remove(temporary); err != nil {
		_ = directory.Remove(leaf)
		return err
	}
	return nil
}
