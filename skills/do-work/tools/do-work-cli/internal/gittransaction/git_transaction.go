package gittransaction

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type FailureKind string

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
	// ExistingUntrackedTargetPaths is a narrow opt-in for exact durable state
	// that is intentionally untracked in some consumer repositories. Each path
	// must also be a target. Its bytes and mode are snapshotted for rollback;
	// every existing untracked target not named here retains the default refusal.
	ExistingUntrackedTargetPaths []string
	CreatedDirectoryPaths        []string
	DryRun                       bool
	Commit                       bool
	CommitMessage                string
	PostCommitVerify             func(context.Context, string) error
}

type MutationRecorder struct {
	allowedPaths              map[string]struct{}
	creatablePaths            map[string]struct{}
	touchedPaths              map[string]struct{}
	createdPaths              map[string]struct{}
	allowedCreatedDirectories map[string]struct{}
	createdDirectories        map[string]struct{}
}

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
		result.Failure = &TransactionFailure{Kind: FailureNotGit, Reason: err.Error()}
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
	return nil
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
	return nil
}

type targetState struct {
	path                     string
	tracked                  bool
	existed                  bool
	existingUntrackedAllowed bool
	originalBytes            []byte
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
		return failTransaction(result, resultmodel.OutcomeFailure, FailureNotGit, err.Error())
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
	for _, path := range allowedExistingUntracked {
		if _, targeted := targetSet[path]; !targeted {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions,
				fmt.Sprintf("existing untracked opt-in path %q is not a declared target", path), path)
		}
	}
	for stateIndex := range states {
		state := &states[stateIndex]
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
			state.originalMode = fileInfo.Mode().Perm()
		}
	}
	for _, state := range states {
		dirty, statusErr := targetIsDirty(ctx, repositoryRoot, state.path)
		if statusErr != nil {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, statusErr.Error(), state.path)
		}
		if dirty && !state.existingUntrackedAllowed {
			return failTransaction(result, resultmodel.OutcomeRefused, FailureDirtyTarget, fmt.Sprintf("target path %q is already dirty", state.path), state.path)
		}
		if state.existed && !state.tracked && !state.existingUntrackedAllowed {
			return failTransaction(result, resultmodel.OutcomeRefused, FailureDirtyTarget, fmt.Sprintf("target path %q exists but cannot be restored from Git", state.path), state.path)
		}
	}
	if options.Commit {
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
	recorder := newMutationRecorder(states, createdDirectories)
	if mutationErr := mutate(recorder); mutationErr != nil {
		return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureMutation, mutationErr)
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
	return result
}

func committableChangedPaths(repositoryRoot string, states []targetState, changedPaths []string) ([]string, error) {
	stateByPath := make(map[string]targetState, len(states))
	for _, state := range states {
		stateByPath[state.path] = state
	}
	paths := make([]string, 0, len(changedPaths))
	for _, path := range changedPaths {
		state := stateByPath[path]
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
		paths = append(paths, path)
	}
	return paths, nil
}

func newMutationRecorder(states []targetState, createdDirectories []string) *MutationRecorder {
	allowed := make(map[string]struct{}, len(states))
	creatable := make(map[string]struct{}, len(states))
	for _, state := range states {
		allowed[state.path] = struct{}{}
		if !state.existed {
			creatable[state.path] = struct{}{}
		}
	}
	return &MutationRecorder{
		allowedPaths:              allowed,
		creatablePaths:            creatable,
		touchedPaths:              map[string]struct{}{},
		createdPaths:              map[string]struct{}{},
		allowedCreatedDirectories: stringSet(createdDirectories),
		createdDirectories:        map[string]struct{}{},
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
	return filepath.Clean(strings.TrimSpace(output)), nil
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
		if existed && info.IsDir() {
			return nil, fmt.Errorf("target %q is a directory; declare exact files instead", path)
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
	for _, state := range states {
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
	rollback := resultmodel.RollbackResult{Status: resultmodel.RollbackSucceeded, Actions: []string{}, Errors: []string{}}
	for _, state := range states {
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
			if writeError := os.WriteFile(absolutePath, state.originalBytes, state.originalMode); writeError != nil {
				rollback.Errors = append(rollback.Errors, fmt.Sprintf("restore existing untracked target %s: %v", state.path, writeError))
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
		absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(path))
		if _, err := runGit(ctx, repositoryRoot, "rm", "--cached", "--ignore-unmatch", "--", path); err != nil {
			rollback.Errors = append(rollback.Errors, fmt.Sprintf("unstage created target %s: %v", path, err))
		}
		if err := os.Remove(absolutePath); err != nil && !os.IsNotExist(err) {
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
		absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(path))
		if err := os.Remove(absolutePath); err != nil && !os.IsNotExist(err) {
			rollback.Errors = append(rollback.Errors, fmt.Sprintf("remove created directory %s: %v", path, err))
		} else {
			rollback.Actions = append(rollback.Actions, "removed created directory "+path)
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
	var standardOutput bytes.Buffer
	var standardError bytes.Buffer
	command.Stdout = &standardOutput
	command.Stderr = &standardError
	if err := command.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(standardError.String()))
	}
	return standardOutput.String(), nil
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
