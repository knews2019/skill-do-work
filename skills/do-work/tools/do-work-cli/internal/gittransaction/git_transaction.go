package gittransaction

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type FailureKind string

var errRepositoryRootMismatch = errors.New("supplied repository root does not identify the Git worktree root")

const (
	FailureInvalidOptions FailureKind = "invalid_options"
	FailureNotGit         FailureKind = "not_git_repository"
	FailureDirtyTarget    FailureKind = "dirty_target"
	FailureDirtyIndex     FailureKind = "dirty_index"
	FailureMutation       FailureKind = "mutation_failed"
	FailureRollback       FailureKind = "rollback_incomplete"
	FailureCommit         FailureKind = "commit_failed"
	FailureCommittedRisk  FailureKind = "committed_state_risk"
)

type TransactionFailure struct {
	Kind   FailureKind
	Reason string
	Paths  []string
}

// TransactionResult reports an outcome, not a number. resultmodel.ExitCode is the only
// place an outcome becomes a process status.
type TransactionResult struct {
	Outcome            resultmodel.CommandOutcome
	RepositoryRoot     string
	ChangedPaths       []string
	CreatedPaths       []string
	CreatedDirectories []string
	CommitSHA          string
	RevertArgv         []string
	Rollback           resultmodel.RollbackResult
	Failure            *TransactionFailure
}

type TransactionOptions struct {
	RepositoryRoot string
	TargetPaths    []string
	// ExistingDirtyTargetPaths is a narrow opt-in for exact tracked files whose
	// unstaged bytes are intentional transaction input. Staged targets remain
	// refused. A tracked deletion is also an eligible exact preimage. The exact
	// presence, bytes, mode, and identity are restored on rollback.
	ExistingDirtyTargetPaths []string
	// CommitExistingDirtyTargets permits a semantic transaction owner to commit
	// only the postimages it produced from ExistingDirtyTargetPaths. The default
	// remains refusal so ordinary callers cannot launder arbitrary dirty bytes.
	CommitExistingDirtyTargets bool
	// ExistingUntrackedTargetPaths is a narrow opt-in for exact durable state
	// that is intentionally untracked in some consumer repositories. Each path
	// must also be a target. Its bytes and mode are snapshotted for rollback;
	// every existing untracked target not named here retains the default refusal.
	ExistingUntrackedTargetPaths []string
	// PrivateUntrackedTargetPaths declares exact intentionally-untracked durable
	// files which participate in change detection and rollback but must never be
	// staged. Unlike ExistingUntrackedTargetPaths, these targets may be absent and
	// ignored before the transaction.
	PrivateUntrackedTargetPaths []string
	CreatedDirectoryPaths       []string
	DryRun                      bool
	Commit                      bool
	CommitMessage               string
	PostCommitVerify            func(context.Context, string) error
}

type MutationRecorder struct {
	allowedPaths              map[string]struct{}
	creatablePaths            map[string]struct{}
	touchedPaths              map[string]struct{}
	createdPaths              map[string]struct{}
	allowedCreatedDirectories map[string]struct{}
	createdDirectories        map[string]struct{}
	repositoryRoot            string
	publishedPrivate          map[string]publishedPrivateState
	publishedTracked          map[string]publishedTrackedState
	publishedDirectories      map[string]os.FileInfo
	createdObjects            map[string]createdObjectIdentity
	dirtyTrackedPaths         map[string]struct{}
}

// createdObjectIdentity binds a created path to the exact filesystem object this
// invocation published there. Inode identity alone is not enough: removing a file
// and recreating it under the same name commonly reuses the inode, so the content
// digest is part of the binding.
type createdObjectIdentity struct {
	info   os.FileInfo
	digest [sha256.Size]byte
}

type publishedPrivateState struct {
	info   os.FileInfo
	digest [sha256.Size]byte
}

type publishedTrackedState struct {
	existed bool
	info    os.FileInfo
	digest  [sha256.Size]byte
}

// privateTransactionTestHook is nil outside deterministic adversarial tests.
// Stage names describe identity boundaries, never production behavior.
var privateTransactionTestHook func(stage, path string)

// RecordCreatedDirectory records an exact directory created by this invocation.
// Rollback removes it only if it is empty, deepest first.
func (recorder *MutationRecorder) RecordCreatedDirectory(path string) error {
	normalized, err := normalizeTargetPath(path)
	if err != nil {
		return err
	}
	if _, allowed := recorder.allowedCreatedDirectories[normalized]; !allowed {
		return fmt.Errorf("directory %q is outside the declared transaction directories", path)
	}
	recorder.createdDirectories[normalized] = struct{}{}
	root, rootError := os.OpenRoot(recorder.repositoryRoot)
	if rootError != nil {
		return fmt.Errorf("open transaction root: %w", rootError)
	}
	defer root.Close()
	info, statError := root.Lstat(filepath.FromSlash(normalized))
	if os.IsNotExist(statError) {
		// Legacy callers declare the directory immediately before rooted creation.
		// Its identity is captured by the first subsequently recorded child.
		return nil
	}
	if statError != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("created directory %q is not an owned real directory", normalized)
	}
	recorder.publishedDirectories[normalized] = info
	return nil
}

// TargetPreflight is a read-only exact-path eligibility check used to keep one
// dirty operation group from blocking unrelated cleanup repairs.
type TargetPreflight struct {
	RepositoryRoot string
	TargetPaths    []string
	Failure        *TransactionFailure
}

// PreflightTargets applies the transaction's Git and dirty-target guards without mutating.
func PreflightTargets(ctx context.Context, repositoryRoot string, targetPaths []string, requireEmptyIndex bool) TargetPreflight {
	result := TargetPreflight{}
	resolvedRoot, err := resolveRepositoryRoot(ctx, repositoryRoot)
	if err != nil {
		result.Failure = &TransactionFailure{Kind: repositoryRootFailureKind(err), Reason: err.Error()}
		return result
	}
	result.RepositoryRoot = resolvedRoot
	normalizedPaths, err := normalizeTargetPaths(targetPaths)
	if err != nil || len(normalizedPaths) == 0 {
		if err == nil {
			err = errors.New("at least one exact target path is required")
		}
		result.Failure = &TransactionFailure{Kind: FailureInvalidOptions, Reason: err.Error()}
		return result
	}
	result.TargetPaths = normalizedPaths
	if requireEmptyIndex {
		empty, indexError := indexIsEmpty(ctx, resolvedRoot)
		if indexError != nil || !empty {
			reason := "--commit requires an empty existing index"
			if indexError != nil {
				reason = indexError.Error()
			}
			result.Failure = &TransactionFailure{Kind: FailureDirtyIndex, Reason: reason}
			return result
		}
	}
	states, err := inspectTargets(ctx, resolvedRoot, normalizedPaths)
	if err != nil {
		result.Failure = &TransactionFailure{Kind: FailureInvalidOptions, Reason: err.Error(), Paths: normalizedPaths}
		return result
	}
	for _, state := range states {
		dirty, statusError := targetIsDirty(ctx, resolvedRoot, state.path)
		if statusError != nil {
			result.Failure = &TransactionFailure{Kind: FailureInvalidOptions, Reason: statusError.Error(), Paths: []string{state.path}}
			return result
		}
		if dirty || (state.existed && !state.tracked) {
			result.Failure = &TransactionFailure{Kind: FailureDirtyTarget, Reason: fmt.Sprintf("target path %q is already dirty or not restorable from Git", state.path), Paths: []string{state.path}}
			return result
		}
	}
	return result
}

func (recorder *MutationRecorder) RecordTouched(path string) error {
	normalized, err := normalizeTargetPath(path)
	if err != nil {
		return err
	}
	if _, allowed := recorder.allowedPaths[normalized]; !allowed {
		return fmt.Errorf("path %q is outside the declared transaction targets", path)
	}
	recorder.touchedPaths[normalized] = struct{}{}
	root, rootError := os.OpenRoot(recorder.repositoryRoot)
	if rootError != nil {
		return rootError
	}
	defer root.Close()
	if _, dirtyTracked := recorder.dirtyTrackedPaths[normalized]; dirtyTracked {
		info, digest, snapshotError := rootedRegularSnapshot(root, normalized)
		if isMissingPathError(snapshotError) {
			recorder.publishedTracked[normalized] = publishedTrackedState{}
		} else if snapshotError != nil {
			return snapshotError
		} else {
			recorder.publishedTracked[normalized] = publishedTrackedState{existed: true, info: info, digest: digest}
		}
	}
	// This transaction is told about each of its own later mutations, so re-capturing the
	// recorded path here keeps ownership across our own republication (which renames, and
	// therefore changes the inode). A foreign writer's swap never routes through the
	// recorder, which is what makes the remaining mismatch below a real swap.
	if _, created := recorder.createdPaths[normalized]; created {
		if captureError := recorder.captureCreatedObject(root, normalized); captureError != nil {
			return captureError
		}
	}
	return recorder.revalidateCreatedObjects(root, normalized)
}

// captureCreatedObject binds the object now standing at a created path to this
// invocation. An absent path records no binding, and rollback then preserves
// whatever it finds there instead of removing it by pathname.
func (recorder *MutationRecorder) captureCreatedObject(root *os.Root, path string) error {
	info, digest, err := rootedCreatedTargetSnapshot(root, path)
	if isMissingPathError(err) {
		delete(recorder.createdObjects, path)
		return nil
	}
	if err != nil {
		return fmt.Errorf("identity-record created object: %w", err)
	}
	recorder.createdObjects[path] = createdObjectIdentity{info: info, digest: digest}
	return nil
}

// revalidateCreatedObjects proves every other already-created path still holds the object
// this invocation published there, so a swap surfaces at the next mutation point instead of
// only at rollback. Paths are checked in sorted order so two simultaneous swaps always
// report the same path.
func (recorder *MutationRecorder) revalidateCreatedObjects(root *os.Root, excludedPath string) error {
	recordedPaths := make([]string, 0, len(recorder.createdObjects))
	for path := range recorder.createdObjects {
		recordedPaths = append(recordedPaths, path)
	}
	sort.Strings(recordedPaths)
	for _, path := range recordedPaths {
		if path == excludedPath {
			continue
		}
		if !createdObjectStillOwned(root, path, recorder.createdObjects[path]) {
			return fmt.Errorf("created target %q changed after publication", path)
		}
	}
	return nil
}

// createdObjectOwnership says what stands at a created path relative to the object this
// invocation published there. Absence is its own answer: rollback wants the path gone, so
// nothing standing there is nothing to preserve and nothing to remove.
type createdObjectOwnership int

const (
	createdObjectAbsent createdObjectOwnership = iota
	createdObjectOwned
	createdObjectReplaced
)

func inspectCreatedObject(root *os.Root, path string, identity createdObjectIdentity) createdObjectOwnership {
	if root == nil {
		return createdObjectReplaced
	}
	info, digest, err := rootedCreatedTargetSnapshot(root, path)
	switch {
	case isMissingPathError(err):
		return createdObjectAbsent
	case err == nil && os.SameFile(identity.info, info) && identity.digest == digest:
		return createdObjectOwned
	default:
		return createdObjectReplaced
	}
}

func createdObjectStillOwned(root *os.Root, path string, identity createdObjectIdentity) bool {
	return inspectCreatedObject(root, path, identity) == createdObjectOwned
}

func (recorder *MutationRecorder) RecordCreated(path string) error {
	if err := recorder.RecordTouched(path); err != nil {
		return err
	}
	normalized, _ := normalizeTargetPath(path)
	if _, creatable := recorder.creatablePaths[normalized]; !creatable {
		return fmt.Errorf("path %q existed before the transaction and cannot be recorded as created", path)
	}
	recorder.createdPaths[normalized] = struct{}{}
	root, rootError := os.OpenRoot(recorder.repositoryRoot)
	if rootError != nil {
		return fmt.Errorf("open transaction root: %w", rootError)
	}
	defer root.Close()
	// The binding is captured last so a failing directory capture cannot leave a recorded
	// identity behind for a creation this call never finished registering.
	if err := recorder.captureCreatedDirectoryIdentities(); err != nil {
		return err
	}
	return recorder.captureCreatedObject(root, normalized)
}

func (recorder *MutationRecorder) captureCreatedDirectoryIdentities() error {
	root, err := os.OpenRoot(recorder.repositoryRoot)
	if err != nil {
		return err
	}
	defer root.Close()
	for path := range recorder.createdDirectories {
		if _, captured := recorder.publishedDirectories[path]; captured {
			continue
		}
		info, statError := root.Lstat(filepath.FromSlash(path))
		if os.IsNotExist(statError) {
			continue
		}
		if statError != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("created directory %q is not an owned real directory", path)
		}
		recorder.publishedDirectories[path] = info
	}
	return nil
}

// RecordPublished binds the current private target identity and bytes to this
// transaction. Rollback will remove or replace only this exact object; a later
// writer's replacement is preserved and reported as incomplete rollback.
func (recorder *MutationRecorder) RecordPublished(path string) error {
	normalized, err := normalizeTargetPath(path)
	if err != nil {
		return err
	}
	if _, allowed := recorder.allowedPaths[normalized]; !allowed {
		return fmt.Errorf("path %q is outside the declared transaction targets", path)
	}
	root, err := os.OpenRoot(recorder.repositoryRoot)
	if err != nil {
		return fmt.Errorf("open transaction root: %w", err)
	}
	defer root.Close()
	info, digest, err := rootedRegularSnapshot(root, normalized)
	if err != nil {
		return err
	}
	recorder.publishedPrivate[normalized] = publishedPrivateState{info: info, digest: digest}
	return nil
}

type targetState struct {
	path                     string
	tracked                  bool
	existed                  bool
	existingUntrackedAllowed bool
	privateUntracked         bool
	existingDirtyAllowed     bool
	originalBytes            []byte
	originalDigest           [sha256.Size]byte
	originalMode             os.FileMode
}

func ExecuteTransaction(ctx context.Context, options TransactionOptions, mutate func(*MutationRecorder) error) TransactionResult {
	// CommandOutcome's zero value is the empty string, which resultmodel.ExitCode reads as a
	// failure. Every success path returns this construction unchanged, so it names success here.
	result := TransactionResult{
		Outcome:  resultmodel.OutcomeSuccess,
		Rollback: resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded, Actions: []string{}, Errors: []string{}},
	}
	if options.DryRun && options.Commit {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, "--dry-run and --commit cannot be combined")
	}
	if len(options.TargetPaths) == 0 {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, "at least one exact target path is required")
	}
	repositoryRoot, err := resolveRepositoryRoot(ctx, options.RepositoryRoot)
	if err != nil {
		return failTransaction(result, resultmodel.OutcomeFailure, repositoryRootFailureKind(err), err.Error())
	}
	result.RepositoryRoot = repositoryRoot
	targetPaths, err := normalizeTargetPaths(options.TargetPaths)
	if err != nil {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, err.Error())
	}
	states, err := inspectTargets(ctx, repositoryRoot, targetPaths)
	if err != nil {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, err.Error(), targetPaths...)
	}
	allowedExistingUntracked, err := normalizeTargetPaths(options.ExistingUntrackedTargetPaths)
	if err != nil {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, err.Error())
	}
	targetSet := stringSet(targetPaths)
	allowedSet := stringSet(allowedExistingUntracked)
	allowedDirtyTargets, err := normalizeTargetPaths(options.ExistingDirtyTargetPaths)
	if err != nil {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, err.Error())
	}
	dirtyAllowedSet := stringSet(allowedDirtyTargets)
	privatePaths, err := normalizeTargetPaths(options.PrivateUntrackedTargetPaths)
	if err != nil {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, err.Error())
	}
	privateSet := stringSet(privatePaths)
	for _, path := range allowedExistingUntracked {
		if _, targeted := targetSet[path]; !targeted {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions,
				fmt.Sprintf("existing untracked opt-in path %q is not a declared target", path), path)
		}
	}
	for _, path := range privatePaths {
		if _, targeted := targetSet[path]; !targeted {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions,
				fmt.Sprintf("private untracked path %q is not a declared target", path), path)
		}
		if _, legacy := allowedSet[path]; legacy {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions,
				fmt.Sprintf("private untracked path %q cannot use both untracked options", path), path)
		}
	}
	for _, path := range allowedDirtyTargets {
		if _, targeted := targetSet[path]; !targeted {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions,
				fmt.Sprintf("existing dirty opt-in path %q is not a declared target", path), path)
		}
		if _, legacy := allowedSet[path]; legacy {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions,
				fmt.Sprintf("existing dirty path %q cannot also be an untracked target", path), path)
		}
		if _, private := privateSet[path]; private {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions,
				fmt.Sprintf("existing dirty path %q cannot also be a private target", path), path)
		}
	}
	root, rootErr := os.OpenRoot(repositoryRoot)
	if rootErr != nil {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, rootErr.Error())
	}
	defer root.Close()
	for stateIndex := range states {
		state := &states[stateIndex]
		if _, allowed := dirtyAllowedSet[state.path]; allowed {
			if !state.tracked {
				return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions,
					fmt.Sprintf("existing dirty opt-in path %q must be tracked", state.path), state.path)
			}
			indexEmpty, indexError := indexIsEmpty(ctx, repositoryRoot, state.path)
			if indexError != nil || !indexEmpty {
				reason := fmt.Sprintf("existing dirty opt-in path %q must not be staged", state.path)
				if indexError != nil {
					reason = indexError.Error()
				}
				return failTransaction(result, resultmodel.OutcomeRefused, FailureDirtyIndex, reason, state.path)
			}
			dirty, statusError := targetIsDirty(ctx, repositoryRoot, state.path)
			if statusError != nil || !dirty {
				reason := fmt.Sprintf("existing dirty opt-in path %q is not dirty", state.path)
				if statusError != nil {
					reason = statusError.Error()
				}
				return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, reason, state.path)
			}
			if state.existed {
				fileInfo, digest, originalBytes, snapshotError := rootedRegularPreimage(root, state.path)
				if snapshotError != nil {
					return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, snapshotError.Error(), state.path)
				}
				state.originalBytes = originalBytes
				state.originalMode = completeRegularFileMode(fileInfo.Mode())
				state.originalDigest = digest
			}
			state.existingDirtyAllowed = true
			continue
		}
		if _, private := privateSet[state.path]; private {
			if state.tracked {
				return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions,
					fmt.Sprintf("private target %q is tracked by Git", state.path), state.path)
			}
			state.privateUntracked = true
			if state.existed {
				fileInfo, digest, originalBytes, snapshotError := rootedRegularPreimage(root, state.path)
				if snapshotError != nil {
					return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, snapshotError.Error(), state.path)
				}
				state.originalBytes = originalBytes
				state.originalMode = completeRegularFileMode(fileInfo.Mode())
				state.originalDigest = digest
			}
			continue
		}
		if _, allowed := allowedSet[state.path]; allowed {
			if !state.existed || state.tracked {
				return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions,
					fmt.Sprintf("existing untracked opt-in path %q must exist and be untracked", state.path), state.path)
			}
			absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(state.path))
			fileInfo, statError := os.Lstat(absolutePath)
			if statError != nil || !fileInfo.Mode().IsRegular() {
				return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions,
					fmt.Sprintf("existing untracked opt-in path %q is not a regular file", state.path), state.path)
			}
			originalBytes, readError := os.ReadFile(absolutePath)
			if readError != nil {
				return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, readError.Error(), state.path)
			}
			state.existingUntrackedAllowed = true
			state.originalBytes = originalBytes
			state.originalMode = completeRegularFileMode(fileInfo.Mode())
		}
	}
	for _, state := range states {
		dirty, statusErr := targetIsDirty(ctx, repositoryRoot, state.path)
		if statusErr != nil {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, statusErr.Error(), state.path)
		}
		if dirty && !state.existingUntrackedAllowed && !state.privateUntracked && !state.existingDirtyAllowed {
			return failTransaction(result, resultmodel.OutcomeRefused, FailureDirtyTarget, fmt.Sprintf("target path %q is already dirty", state.path), state.path)
		}
		if state.existed && !state.tracked && !state.existingUntrackedAllowed && !state.privateUntracked {
			return failTransaction(result, resultmodel.OutcomeRefused, FailureDirtyTarget, fmt.Sprintf("target path %q exists but cannot be restored from Git", state.path), state.path)
		}
	}
	if options.Commit {
		if len(allowedDirtyTargets) > 0 && !options.CommitExistingDirtyTargets {
			return failTransaction(result, resultmodel.OutcomeRefused, FailureDirtyTarget, "--commit cannot include pre-existing dirty target bytes", allowedDirtyTargets...)
		}
		if options.CommitExistingDirtyTargets && len(allowedDirtyTargets) == 0 {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, "dirty-target commit authority requires at least one exact dirty target")
		}
		empty, indexErr := indexIsEmpty(ctx, repositoryRoot)
		if indexErr != nil {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, indexErr.Error())
		}
		if !empty {
			return failTransaction(result, resultmodel.OutcomeRefused, FailureDirtyIndex, "--commit requires an empty existing index")
		}
		if strings.TrimSpace(options.CommitMessage) == "" {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, "--commit requires a non-empty commit message")
		}
	}
	if options.DryRun {
		return result
	}
	if mutate == nil {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, "mutation callback is required")
	}
	createdDirectories, err := normalizeCreatedDirectories(options.CreatedDirectoryPaths, repositoryRoot)
	if err != nil {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, err.Error())
	}
	recorder := newMutationRecorder(repositoryRoot, states, createdDirectories)
	if mutationErr := mutate(recorder); mutationErr != nil {
		return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureMutation, mutationErr)
	}
	if captureError := recorder.captureTrackedPublications(root); captureError != nil {
		return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureMutation, captureError)
	}
	changedPaths, err := changedTargets(ctx, repositoryRoot, states)
	if err != nil {
		return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureMutation, err)
	}
	result.ChangedPaths = changedPaths
	result.CreatedPaths = sortedKeys(recorder.createdPaths)
	result.CreatedDirectories = sortedKeys(recorder.createdDirectories)
	// A changed path must have been recorded, and a path that did NOT exist beforehand must
	// have been recorded as CREATED. Without the second half, a file the mutation brought
	// into existence through RecordTouched alone reports success, and the transaction claims
	// to describe a creation it never saw.
	existedBefore := make(map[string]bool, len(states))
	for _, state := range states {
		existedBefore[state.path] = state.existed
	}
	for _, path := range changedPaths {
		if _, recorded := recorder.touchedPaths[path]; !recorded {
			return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureMutation,
				fmt.Errorf("changed path %q was not recorded by the mutation", path))
		}
		if _, created := recorder.createdPaths[path]; !created && !existedBefore[path] {
			return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureMutation,
				fmt.Errorf("changed path %q was created but was not recorded as created", path))
		}
	}
	if verifyError := verifyPublishedPrivateTargets(root, states, recorder, changedPaths); verifyError != nil {
		return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureMutation, verifyError)
	}
	if !options.Commit || len(changedPaths) == 0 {
		return result
	}
	commitPaths, err := committableChangedPaths(repositoryRoot, states, changedPaths)
	if err != nil {
		return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureCommit, err)
	}
	if len(commitPaths) == 0 {
		return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureCommit, errors.New("the transaction changed no paths Git can commit"))
	}
	if _, err := runGit(ctx, repositoryRoot, append([]string{"add", "-A", "--"}, commitPaths...)...); err != nil {
		return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureCommit, err)
	}
	if _, err := runGit(ctx, repositoryRoot, "commit", "-m", options.CommitMessage); err != nil {
		return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureCommit, err)
	}
	commitSHA, err := runGit(ctx, repositoryRoot, "rev-parse", "HEAD")
	if err != nil {
		return committedRisk(result, "the commit succeeded but its ID could not be read", "HEAD")
	}
	commitSHA = strings.TrimSpace(commitSHA)
	result.CommitSHA = commitSHA
	result.RevertArgv = []string{"git", "revert", commitSHA}
	if empty, indexErr := indexIsEmpty(ctx, repositoryRoot); indexErr != nil {
		return committedRisk(result, indexErr.Error(), commitSHA)
	} else if !empty {
		return committedRisk(result, "the Git index is not empty after the exact-path commit", commitSHA)
	}
	committedPaths, verifyErr := committedPaths(ctx, repositoryRoot, commitSHA)
	if verifyErr != nil {
		return committedRisk(result, verifyErr.Error(), commitSHA)
	}
	if !equalStrings(committedPaths, commitPaths) {
		return committedRisk(result,
			fmt.Sprintf("committed paths %q do not match exact committable paths %q", committedPaths, commitPaths), commitSHA)
	}
	if options.PostCommitVerify != nil {
		if verifyErr := options.PostCommitVerify(ctx, commitSHA); verifyErr != nil {
			return committedRisk(result, verifyErr.Error(), commitSHA)
		}
	}
	if verifyError := verifyPublishedPrivateTargets(root, states, recorder, changedPaths); verifyError != nil {
		return committedRisk(result, verifyError.Error(), commitSHA)
	}
	return result
}

func (recorder *MutationRecorder) captureTrackedPublications(root *os.Root) error {
	for path := range recorder.touchedPaths {
		if _, captured := recorder.publishedTracked[path]; captured {
			continue
		}
		info, digest, err := rootedRegularSnapshot(root, path)
		if isMissingPathError(err) {
			recorder.publishedTracked[path] = publishedTrackedState{}
			continue
		}
		if err != nil {
			return fmt.Errorf("identity-record changed target %q: %w", path, err)
		}
		recorder.publishedTracked[path] = publishedTrackedState{existed: true, info: info, digest: digest}
	}
	return nil
}

func verifyPublishedPrivateTargets(root *os.Root, states []targetState, recorder *MutationRecorder, changedPaths []string) error {
	changed := stringSet(changedPaths)
	for _, state := range states {
		if !state.privateUntracked {
			continue
		}
		if _, isChanged := changed[state.path]; !isChanged {
			continue
		}
		published, recorded := recorder.publishedPrivate[state.path]
		if !recorded {
			return fmt.Errorf("changed private target %q was not identity-recorded", state.path)
		}
		if privateTransactionTestHook != nil {
			privateTransactionTestHook("before-private-final-verify", state.path)
		}
		info, digest, err := rootedRegularSnapshot(root, state.path)
		if err != nil || !os.SameFile(published.info, info) || published.digest != digest {
			return fmt.Errorf("private target %q changed after publication", state.path)
		}
	}
	return nil
}

func committableChangedPaths(repositoryRoot string, states []targetState, changedPaths []string) ([]string, error) {
	stateByPath := make(map[string]targetState, len(states))
	for _, state := range states {
		stateByPath[state.path] = state
	}
	paths := make([]string, 0, len(changedPaths))
	for _, path := range changedPaths {
		state := stateByPath[path]
		if state.privateUntracked {
			continue
		}
		if state.existingUntrackedAllowed {
			_, statError := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(path)))
			if os.IsNotExist(statError) {
				// This path never existed in the index, so its disappearance has no Git
				// deletion to stage. A moved destination remains a separate exact target.
				continue
			}
			if statError != nil {
				return nil, statError
			}
		}
		dirty, statusError := targetIsDirty(context.Background(), repositoryRoot, path)
		if statusError != nil {
			return nil, statusError
		}
		if !dirty {
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func newMutationRecorder(repositoryRoot string, states []targetState, createdDirectories []string) *MutationRecorder {
	allowed := make(map[string]struct{}, len(states))
	creatable := make(map[string]struct{}, len(states))
	dirtyTracked := make(map[string]struct{})
	for _, state := range states {
		allowed[state.path] = struct{}{}
		if !state.existed {
			creatable[state.path] = struct{}{}
		}
		if state.existingDirtyAllowed {
			dirtyTracked[state.path] = struct{}{}
		}
	}
	return &MutationRecorder{
		allowedPaths:              allowed,
		creatablePaths:            creatable,
		touchedPaths:              map[string]struct{}{},
		createdPaths:              map[string]struct{}{},
		allowedCreatedDirectories: stringSet(createdDirectories),
		createdDirectories:        map[string]struct{}{},
		repositoryRoot:            repositoryRoot,
		publishedPrivate:          map[string]publishedPrivateState{},
		publishedTracked:          map[string]publishedTrackedState{},
		publishedDirectories:      map[string]os.FileInfo{},
		createdObjects:            map[string]createdObjectIdentity{},
		dirtyTrackedPaths:         dirtyTracked,
	}
}

func normalizeCreatedDirectories(paths []string, repositoryRoot string) ([]string, error) {
	normalized, err := normalizeTargetPaths(paths)
	if err != nil {
		return nil, err
	}
	for _, path := range normalized {
		info, statError := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(path)))
		if statError == nil {
			if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				return nil, fmt.Errorf("created directory target %q is not a real directory", path)
			}
			return nil, fmt.Errorf("created directory target %q already exists", path)
		}
		if !os.IsNotExist(statError) {
			return nil, fmt.Errorf("inspect created directory target %q: %w", path, statError)
		}
	}
	return normalized, nil
}

func resolveRepositoryRoot(ctx context.Context, suppliedRoot string) (string, error) {
	if suppliedRoot == "" {
		var err error
		suppliedRoot, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("read working directory: %w", err)
		}
	}
	absoluteRoot, err := filepath.Abs(suppliedRoot)
	if err != nil {
		return "", fmt.Errorf("resolve repository root: %w", err)
	}
	output, err := runGit(ctx, absoluteRoot, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", errors.New("mutating commands require a Git repository")
	}
	gitRoot := filepath.Clean(strings.TrimSpace(output))
	suppliedInfo, suppliedError := os.Stat(absoluteRoot)
	gitInfo, gitError := os.Stat(gitRoot)
	if suppliedError != nil || gitError != nil || !suppliedInfo.IsDir() || !gitInfo.IsDir() {
		return "", fmt.Errorf("inspect supplied repository root %q and Git worktree root %q", absoluteRoot, gitRoot)
	}
	if !os.SameFile(suppliedInfo, gitInfo) {
		return "", fmt.Errorf("%w: supplied=%q git=%q", errRepositoryRootMismatch, absoluteRoot, gitRoot)
	}
	return gitRoot, nil
}

func repositoryRootFailureKind(err error) FailureKind {
	if errors.Is(err, errRepositoryRootMismatch) {
		return FailureInvalidOptions
	}
	return FailureNotGit
}

func normalizeTargetPaths(paths []string) ([]string, error) {
	seen := map[string]struct{}{}
	normalized := make([]string, 0, len(paths))
	for _, path := range paths {
		cleaned, err := normalizeTargetPath(path)
		if err != nil {
			return nil, err
		}
		if _, duplicate := seen[cleaned]; duplicate {
			continue
		}
		seen[cleaned] = struct{}{}
		normalized = append(normalized, cleaned)
	}
	sort.Strings(normalized)
	return normalized, nil
}

func normalizeTargetPath(path string) (string, error) {
	if path == "" || filepath.IsAbs(path) {
		return "", fmt.Errorf("target path %q must be a non-empty repository-relative path", path)
	}
	cleaned := filepath.Clean(path)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("target path %q escapes or names the repository root", path)
	}
	return filepath.ToSlash(cleaned), nil
}

func inspectTargets(ctx context.Context, repositoryRoot string, paths []string) ([]targetState, error) {
	states := make([]targetState, 0, len(paths))
	for _, path := range paths {
		absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(path))
		info, statErr := os.Lstat(absolutePath)
		existed := statErr == nil
		if statErr != nil && !os.IsNotExist(statErr) {
			return nil, fmt.Errorf("inspect target %q: %w", path, statErr)
		}
		if existed && !info.Mode().IsRegular() {
			return nil, fmt.Errorf("target %q is not a regular file; declare exact regular files instead", path)
		}
		_, trackedErr := runGit(ctx, repositoryRoot, "ls-files", "--error-unmatch", "--", path)
		states = append(states, targetState{path: path, tracked: trackedErr == nil, existed: existed})
	}
	return states, nil
}

func targetIsDirty(ctx context.Context, repositoryRoot, path string) (bool, error) {
	output, err := runGit(ctx, repositoryRoot, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--", path)
	if err != nil {
		return false, err
	}
	return len(output) > 0, nil
}

// indexIsEmpty answers whether anything is staged. With no pathspecs it asks about the
// whole index; with pathspecs it asks only about those paths, which is what a post-rollback
// check needs — the user's unrelated staged work is dirt this transaction must leave alone,
// not evidence that the rollback failed.
func indexIsEmpty(ctx context.Context, repositoryRoot string, pathspecs ...string) (bool, error) {
	commandArgs := append([]string{
		"-C", repositoryRoot, "--literal-pathspecs",
		"diff", "--cached", "--quiet", "--exit-code", "--",
	}, pathspecs...)
	command := exec.CommandContext(ctx, "git", commandArgs...)
	err := command.Run()
	if err == nil {
		return true, nil
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
		return false, nil
	}
	return false, fmt.Errorf("inspect Git index: %w", err)
}

func changedTargets(ctx context.Context, repositoryRoot string, states []targetState) ([]string, error) {
	changed := make([]string, 0, len(states))
	root, rootError := os.OpenRoot(repositoryRoot)
	if rootError != nil {
		return nil, rootError
	}
	defer root.Close()
	for _, state := range states {
		if state.existingDirtyAllowed {
			if !state.existed {
				_, statError := root.Lstat(filepath.FromSlash(state.path))
				if statError == nil {
					changed = append(changed, state.path)
					continue
				}
				if !os.IsNotExist(statError) {
					return nil, statError
				}
				continue
			}
			currentInfo, currentDigest, snapshotError := rootedRegularSnapshot(root, state.path)
			if isMissingPathError(snapshotError) || snapshotError == nil && (currentDigest != state.originalDigest || completeRegularFileMode(currentInfo.Mode()) != state.originalMode) {
				changed = append(changed, state.path)
				continue
			}
			if snapshotError != nil {
				return nil, snapshotError
			}
			continue
		}
		if state.privateUntracked {
			if state.existed {
				currentInfo, currentDigest, snapshotError := rootedRegularSnapshot(root, state.path)
				if isMissingPathError(snapshotError) || snapshotError == nil && (currentDigest != state.originalDigest || completeRegularFileMode(currentInfo.Mode()) != state.originalMode) {
					changed = append(changed, state.path)
					continue
				}
				if snapshotError != nil {
					return nil, snapshotError
				}
				continue
			}
			_, statError := root.Lstat(filepath.FromSlash(state.path))
			if statError == nil {
				changed = append(changed, state.path)
				continue
			}
			if !os.IsNotExist(statError) {
				return nil, statError
			}
			continue
		}
		if state.existingUntrackedAllowed {
			currentBytes, readError := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(state.path)))
			if os.IsNotExist(readError) || readError == nil && !bytes.Equal(currentBytes, state.originalBytes) {
				changed = append(changed, state.path)
				continue
			}
			if readError != nil {
				return nil, readError
			}
			continue
		}
		dirty, err := targetIsDirty(ctx, repositoryRoot, state.path)
		if err != nil {
			return nil, err
		}
		if dirty {
			changed = append(changed, state.path)
		}
	}
	return changed, nil
}

func rollbackFailure(ctx context.Context, result TransactionResult, repositoryRoot string, states []targetState, recorder *MutationRecorder, failureKind FailureKind, operationError error) TransactionResult {
	// Cancellation stops the requested operation, never the cleanup needed to
	// restore or safely preserve exact targets.
	ctx = context.WithoutCancel(ctx)
	rollback := resultmodel.RollbackResult{Status: resultmodel.RollbackSucceeded, Actions: []string{}, Errors: []string{}}
	root, rootError := os.OpenRoot(repositoryRoot)
	if rootError != nil {
		rollback.Errors = append(rollback.Errors, "open rollback root: "+rootError.Error())
	}
	if root != nil {
		defer root.Close()
	}
	for _, state := range states {
		if state.existingDirtyAllowed {
			if _, unstageError := runGit(ctx, repositoryRoot, "restore", "--staged", "--", state.path); unstageError != nil {
				rollback.Errors = append(rollback.Errors, fmt.Sprintf("unstage dirty tracked target %s: %v", state.path, unstageError))
				continue
			}
			published, recorded := recorder.publishedTracked[state.path]
			if !recorded {
				if root != nil && privateStateStillOriginal(root, state) {
					continue
				}
				rollback.Errors = append(rollback.Errors, "dirty tracked target was not identity-recorded: "+state.path)
				continue
			}
			action, restoreError := rollbackDirtyTracked(root, state, published)
			if restoreError != nil {
				rollback.Errors = append(rollback.Errors, restoreError.Error())
				continue
			}
			rollback.Actions = append(rollback.Actions, action)
			continue
		}
		if state.privateUntracked {
			published, recorded := recorder.publishedPrivate[state.path]
			if !recorded {
				if root != nil && privateStateStillOriginal(root, state) {
					continue
				}
				rollback.Errors = append(rollback.Errors, "private target was not identity-recorded: "+state.path)
				continue
			}
			action, privateRollbackError := quarantineAndRollbackPrivate(root, state, published)
			if privateRollbackError != nil {
				rollback.Errors = append(rollback.Errors, privateRollbackError.Error())
				continue
			}
			rollback.Actions = append(rollback.Actions, action)
			continue
		}
		if state.existingUntrackedAllowed {
			if _, err := runGit(ctx, repositoryRoot, "rm", "--cached", "--ignore-unmatch", "--", state.path); err != nil {
				rollback.Errors = append(rollback.Errors, fmt.Sprintf("unstage existing untracked target %s: %v", state.path, err))
			}
			absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(state.path))
			if removeError := os.Remove(absolutePath); removeError != nil && !os.IsNotExist(removeError) {
				rollback.Errors = append(rollback.Errors, fmt.Sprintf("remove changed existing untracked target %s: %v", state.path, removeError))
				continue
			}
			if makeError := os.MkdirAll(filepath.Dir(absolutePath), 0o755); makeError != nil {
				rollback.Errors = append(rollback.Errors, fmt.Sprintf("recreate parent for existing untracked target %s: %v", state.path, makeError))
				continue
			}
			if writeError := os.WriteFile(absolutePath, state.originalBytes, state.originalMode.Perm()); writeError != nil {
				rollback.Errors = append(rollback.Errors, fmt.Sprintf("restore existing untracked target %s: %v", state.path, writeError))
			} else if chmodError := os.Chmod(absolutePath, state.originalMode); chmodError != nil {
				rollback.Errors = append(rollback.Errors, fmt.Sprintf("restore mode for existing untracked target %s: %v", state.path, chmodError))
			} else {
				rollback.Actions = append(rollback.Actions, "restored existing untracked target "+state.path)
			}
			continue
		}
		if !state.tracked {
			continue
		}
		dirty, err := targetIsDirty(ctx, repositoryRoot, state.path)
		if err != nil {
			rollback.Errors = append(rollback.Errors, err.Error())
			continue
		}
		if !dirty {
			continue
		}
		published, recorded := recorder.publishedTracked[state.path]
		if recorded && !trackedPublicationStillOwned(root, state.path, published) {
			if _, unstageError := runGit(ctx, repositoryRoot, "restore", "--staged", "--", state.path); unstageError != nil {
				rollback.Errors = append(rollback.Errors, unstageError.Error())
			}
			rollback.Errors = append(rollback.Errors, "tracked target changed after publication; preserved replacement: "+state.path)
			continue
		}
		if _, err := runGit(ctx, repositoryRoot, "restore", "--source=HEAD", "--staged", "--worktree", "--", state.path); err != nil {
			rollback.Errors = append(rollback.Errors, err.Error())
		} else {
			rollback.Actions = append(rollback.Actions, "restored "+state.path+" from HEAD")
		}
	}
	createdPaths := mapKeys(recorder.createdPaths)
	sort.Slice(createdPaths, func(first, second int) bool {
		return strings.Count(createdPaths[first], "/") > strings.Count(createdPaths[second], "/")
	})
	for _, path := range createdPaths {
		alreadyRestored := false
		for _, state := range states {
			if state.path == path && (state.privateUntracked || state.existingDirtyAllowed || state.existingUntrackedAllowed) {
				alreadyRestored = true
				break
			}
		}
		if alreadyRestored {
			continue
		}
		if _, err := runGit(ctx, repositoryRoot, "rm", "--cached", "--ignore-unmatch", "--", path); err != nil {
			rollback.Errors = append(rollback.Errors, fmt.Sprintf("unstage created target %s: %v", path, err))
		}
		published, recorded := recorder.publishedTracked[path]
		if recorded && !trackedPublicationStillOwned(root, path, published) {
			rollback.Errors = append(rollback.Errors, "created target changed after tracked publication; preserved replacement: "+path)
			continue
		}
		identity, identityRecorded := recorder.createdObjects[path]
		if !identityRecorded {
			// An absent path holds nothing this invocation could wrongly remove; an object
			// standing there is not provably ours, so it is preserved and reported.
			if root != nil {
				if _, statError := root.Lstat(filepath.FromSlash(path)); os.IsNotExist(statError) {
					continue
				}
			}
			rollback.Errors = append(rollback.Errors, "created target was not identity-recorded; preserved object: "+path)
			continue
		}
		switch inspectCreatedObject(root, path, identity) {
		case createdObjectAbsent:
			// Already gone, which is the state this loop is trying to reach. The transaction
			// removing its own creation is a completed removal, never a foreign replacement.
			continue
		case createdObjectReplaced:
			rollback.Errors = append(rollback.Errors, "created target changed after created-object capture; preserved replacement: "+path)
			continue
		}
		if err := root.Remove(filepath.FromSlash(path)); err != nil && !os.IsNotExist(err) {
			rollback.Errors = append(rollback.Errors, fmt.Sprintf("remove created target %s: %v", path, err))
		} else {
			rollback.Actions = append(rollback.Actions, "removed created target "+path)
		}
	}
	createdDirectories := mapKeys(recorder.createdDirectories)
	sort.Slice(createdDirectories, func(first, second int) bool {
		return strings.Count(createdDirectories[first], "/") > strings.Count(createdDirectories[second], "/")
	})
	for _, path := range createdDirectories {
		ownedInfo, recorded := recorder.publishedDirectories[path]
		if !recorded || root == nil {
			rollback.Errors = append(rollback.Errors, "created directory was not identity-recorded: "+path)
			continue
		}
		currentInfo, statError := root.Lstat(filepath.FromSlash(path))
		if statError != nil || !os.SameFile(ownedInfo, currentInfo) {
			rollback.Errors = append(rollback.Errors, "created directory changed after publication; preserved replacement: "+path)
			continue
		}
		if err := root.Remove(filepath.FromSlash(path)); err != nil && !os.IsNotExist(err) {
			rollback.Errors = append(rollback.Errors, fmt.Sprintf("remove owned created directory %s: %v", path, err))
		} else {
			rollback.Actions = append(rollback.Actions, "removed owned created directory "+path)
		}
	}
	rolledBackPaths := make([]string, 0, len(states))
	for _, state := range states {
		rolledBackPaths = append(rolledBackPaths, state.path)
	}
	if empty, err := indexIsEmpty(ctx, repositoryRoot, rolledBackPaths...); err != nil {
		rollback.Errors = append(rollback.Errors, err.Error())
	} else if !empty {
		rollback.Errors = append(rollback.Errors, "the declared target paths are still staged after rollback")
	}
	result.Rollback = rollback
	if len(rollback.Errors) > 0 {
		result.Rollback.Status = resultmodel.RollbackIncomplete
		return failTransaction(result, resultmodel.OutcomeRisk, FailureRollback, operationError.Error(), rolledBackPaths...)
	}
	return failTransaction(result, resultmodel.OutcomeRolledBack, failureKind, operationError.Error(), rolledBackPaths...)
}

func rollbackDirtyTracked(root *os.Root, state targetState, published publishedTrackedState) (string, error) {
	if root == nil {
		return "", errors.New("rollback root is unavailable")
	}
	if !published.existed {
		if _, statError := root.Lstat(filepath.FromSlash(state.path)); !os.IsNotExist(statError) {
			return "", fmt.Errorf("dirty tracked target changed after publication; preserved replacement: %s", state.path)
		}
		if err := rootedCreateRegular(root, state.path, state.originalBytes, state.originalMode); err != nil {
			return "", fmt.Errorf("restore dirty tracked target %s: %w", state.path, err)
		}
		return "restored dirty tracked target " + state.path, nil
	}
	privatePublished := publishedPrivateState{info: published.info, digest: published.digest}
	return quarantineAndRollbackPrivate(root, state, privatePublished)
}

func trackedPublicationStillOwned(root *os.Root, path string, published publishedTrackedState) bool {
	if root == nil {
		return false
	}
	if !published.existed {
		_, err := root.Lstat(filepath.FromSlash(path))
		return os.IsNotExist(err)
	}
	info, digest, err := rootedRegularSnapshot(root, path)
	return err == nil && os.SameFile(published.info, info) && published.digest == digest
}

func quarantineAndRollbackPrivate(root *os.Root, state targetState, published publishedPrivateState) (string, error) {
	if privateTransactionTestHook != nil {
		privateTransactionTestHook("before-private-rollback-quarantine", state.path)
	}
	random := make([]byte, 12)
	if _, err := rand.Read(random); err != nil {
		return "", fmt.Errorf("allocate private rollback quarantine: %w", err)
	}
	directory := ".do-work-private-rollback-" + hex.EncodeToString(random)
	if err := root.Mkdir(directory, 0o700); err != nil {
		return "", fmt.Errorf("create private rollback quarantine: %w", err)
	}
	quarantined := filepath.Join(directory, "object")
	cleanupDirectory := func() { _ = root.Remove(directory) }
	if err := root.Rename(filepath.FromSlash(state.path), quarantined); err != nil {
		cleanupDirectory()
		return "", fmt.Errorf("quarantine private target %s: %w", state.path, err)
	}
	currentInfo, currentDigest, snapshotError := rootedRegularSnapshot(root, quarantined)
	if snapshotError != nil || !os.SameFile(published.info, currentInfo) || published.digest != currentDigest {
		if _, targetError := root.Lstat(filepath.FromSlash(state.path)); os.IsNotExist(targetError) {
			if restoreError := root.Rename(quarantined, filepath.FromSlash(state.path)); restoreError != nil {
				return "", fmt.Errorf("private target changed after publication; replacement retained at %s: %v", quarantined, restoreError)
			}
			cleanupDirectory()
			return "", fmt.Errorf("private target changed after publication; preserved replacement: %s", state.path)
		}
		return "", fmt.Errorf("private target changed after publication; replacements preserved at %s and %s", state.path, quarantined)
	}
	if state.existed {
		if err := rootedCreateRegular(root, state.path, state.originalBytes, state.originalMode); err != nil {
			return "", fmt.Errorf("restore private target %s (published bytes retained at %s): %w", state.path, quarantined, err)
		}
	}
	if err := root.Remove(quarantined); err != nil {
		return "", fmt.Errorf("remove quarantined private publication %s: %w", quarantined, err)
	}
	cleanupDirectory()
	if state.existed {
		return "restored private target " + state.path, nil
	}
	return "removed owned private target " + state.path, nil
}

func privateStateStillOriginal(root *os.Root, state targetState) bool {
	if !state.existed {
		_, err := root.Lstat(filepath.FromSlash(state.path))
		return os.IsNotExist(err)
	}
	info, digest, err := rootedRegularSnapshot(root, state.path)
	return err == nil && digest == state.originalDigest && completeRegularFileMode(info.Mode()) == state.originalMode
}

func rootedRegularSnapshot(root *os.Root, path string) (os.FileInfo, [sha256.Size]byte, error) {
	info, digest, _, err := rootedOpenSnapshot(root, path, "private target", "")
	return info, digest, err
}

// rootedCreatedTargetSnapshot is rootedRegularSnapshot for a path this invocation created.
// It differs only in how its failures name the target: a created cleanup or publication
// target is not a private target.
func rootedCreatedTargetSnapshot(root *os.Root, path string) (os.FileInfo, [sha256.Size]byte, error) {
	info, digest, _, err := rootedOpenSnapshot(root, path, "created target", "")
	return info, digest, err
}

func rootedRegularPreimage(root *os.Root, path string) (os.FileInfo, [sha256.Size]byte, []byte, error) {
	return rootedOpenSnapshot(root, path, "private target", "after-private-preimage-lstat")
}

func rootedOpenSnapshot(root *os.Root, path, targetDescription, hookStage string) (os.FileInfo, [sha256.Size]byte, []byte, error) {
	var empty [sha256.Size]byte
	if root == nil {
		return nil, empty, nil, errors.New("rooted filesystem handle is unavailable")
	}
	rootPath := filepath.FromSlash(path)
	lstatInfo, err := root.Lstat(rootPath)
	if err != nil {
		return nil, empty, nil, fmt.Errorf("inspect %s %q: %w", targetDescription, path, err)
	}
	if !lstatInfo.Mode().IsRegular() {
		return nil, empty, nil, fmt.Errorf("%s %q is not a regular file", targetDescription, path)
	}
	if hookStage != "" && privateTransactionTestHook != nil {
		privateTransactionTestHook(hookStage, path)
	}
	file, err := root.Open(rootPath)
	if err != nil {
		return nil, empty, nil, fmt.Errorf("open %s %q: %w", targetDescription, path, err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil || !openedInfo.Mode().IsRegular() || !os.SameFile(lstatInfo, openedInfo) {
		return nil, empty, nil, fmt.Errorf("%s %q changed while it was opened", targetDescription, path)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, empty, nil, fmt.Errorf("read %s %q: %w", targetDescription, path, err)
	}
	finalInfo, err := file.Stat()
	if err != nil || !os.SameFile(openedInfo, finalInfo) || openedInfo.Size() != finalInfo.Size() || !openedInfo.ModTime().Equal(finalInfo.ModTime()) || int64(len(data)) != finalInfo.Size() {
		return nil, empty, nil, fmt.Errorf("%s %q changed while it was read", targetDescription, path)
	}
	return finalInfo, sha256.Sum256(data), data, nil
}

func rootedCreateRegular(root *os.Root, path string, contents []byte, mode os.FileMode) error {
	rootPath := filepath.FromSlash(path)
	parent := filepath.Dir(rootPath)
	if parent != "." {
		if err := root.MkdirAll(parent, 0o755); err != nil {
			return err
		}
	}
	file, err := root.OpenFile(rootPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode.Perm())
	if err != nil {
		return err
	}
	keep := false
	defer func() {
		if !keep {
			_ = root.Remove(rootPath)
		}
	}()
	if _, err := file.Write(contents); err != nil {
		file.Close()
		return err
	}
	if err := file.Chmod(mode); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	keep = true
	return nil
}

func completeRegularFileMode(fileMode os.FileMode) os.FileMode {
	return fileMode.Perm() | (fileMode & (os.ModeSetuid | os.ModeSetgid | os.ModeSticky))
}

func isMissingPathError(err error) bool {
	return os.IsNotExist(err) || errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT)
}

func committedPaths(ctx context.Context, repositoryRoot, commitSHA string) ([]string, error) {
	output, err := runGit(ctx, repositoryRoot, "diff-tree", "--root", "--no-commit-id", "--name-only", "-r", "-z", commitSHA)
	if err != nil {
		return nil, err
	}
	paths := splitNUL(output)
	sort.Strings(paths)
	return paths, nil
}

func committedRisk(result TransactionResult, reason, commitSHA string) TransactionResult {
	result.CommitSHA = commitSHA
	result.RevertArgv = []string{"git", "revert", commitSHA}
	return failTransaction(result, resultmodel.OutcomeRisk, FailureCommittedRisk, reason)
}

func failTransaction(result TransactionResult, outcome resultmodel.CommandOutcome, kind FailureKind, reason string, paths ...string) TransactionResult {
	result.Outcome = outcome
	result.Failure = &TransactionFailure{Kind: kind, Reason: reason, Paths: paths}
	return result
}

// runGit returns only stdout as the command's answer. git writes warnings — a malformed
// .gitattributes line, for one — to stderr while the porcelain answer stays on stdout, so
// folding the two streams together makes a warning read as content and a clean target
// read as dirty. stderr is kept for the error text and nothing else.
func runGit(ctx context.Context, repositoryRoot string, args ...string) (string, error) {
	commandArgs := append([]string{"-C", repositoryRoot, "--literal-pathspecs"}, args...)
	command := exec.CommandContext(ctx, "git", commandArgs...)
	configureCancellableProcessGroup(command)
	var standardOutput bytes.Buffer
	var standardError bytes.Buffer
	command.Stdout = &standardOutput
	command.Stderr = &standardError
	if err := command.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(standardError.String()))
	}
	return standardOutput.String(), nil
}

// configureCancellableProcessGroup uses reflection so the shared file still
// cross-compiles on platforms whose syscall.SysProcAttr has no Setpgid field.
// On Unix, cancellation reaches Git hooks and their descendants as one owned
// group, escalates after grace, and WaitDelay prevents inherited pipes from
// holding the caller forever.
func configureCancellableProcessGroup(command *exec.Cmd) {
	attributes := &syscall.SysProcAttr{}
	setProcessGroup := reflect.ValueOf(attributes).Elem().FieldByName("Setpgid")
	if !setProcessGroup.IsValid() || !setProcessGroup.CanSet() || setProcessGroup.Kind() != reflect.Bool {
		return
	}
	setProcessGroup.SetBool(true)
	command.SysProcAttr = attributes
	command.WaitDelay = 2 * time.Second
	command.Cancel = func() error {
		if command.Process == nil {
			return os.ErrProcessDone
		}
		processGroup, findError := os.FindProcess(-command.Process.Pid)
		if findError != nil {
			return findError
		}
		termError := processGroup.Signal(syscall.Signal(15))
		go func(processID int) {
			timer := time.NewTimer(time.Second)
			defer timer.Stop()
			<-timer.C
			if group, err := os.FindProcess(-processID); err == nil {
				_ = group.Signal(os.Kill)
			}
		}(command.Process.Pid)
		return termError
	}
}

func splitNUL(value string) []string {
	parts := strings.Split(value, "\x00")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			result = append(result, filepath.ToSlash(part))
		}
	}
	return result
}

func mapKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func sortedKeys(values map[string]struct{}) []string {
	keys := mapKeys(values)
	sort.Strings(keys)
	return keys
}

func equalStrings(first, second []string) bool {
	if len(first) != len(second) {
		return false
	}
	for index := range first {
		if first[index] != second[index] {
			return false
		}
	}
	return true
}

func stringSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}
