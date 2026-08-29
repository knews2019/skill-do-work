package gittransaction

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
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

type RollbackStatus string

const (
	RollbackNotNeeded  RollbackStatus = "not_needed"
	RollbackSucceeded  RollbackStatus = "succeeded"
	RollbackIncomplete RollbackStatus = "incomplete"
)

type RollbackResult struct {
	Status  RollbackStatus
	Actions []string
	Errors  []string
}

type TransactionFailure struct {
	Kind   FailureKind
	Reason string
}

type TransactionResult struct {
	ExitCode       int
	RepositoryRoot string
	ChangedPaths   []string
	CommitSHA      string
	RevertArgv     []string
	Rollback       RollbackResult
	Failure        *TransactionFailure
}

type TransactionOptions struct {
	RepositoryRoot   string
	TargetPaths      []string
	DryRun           bool
	Commit           bool
	CommitMessage    string
	PostCommitVerify func(context.Context, string) error
}

type MutationRecorder struct {
	allowedPaths   map[string]struct{}
	creatablePaths map[string]struct{}
	touchedPaths   map[string]struct{}
	createdPaths   map[string]struct{}
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
	path    string
	tracked bool
	existed bool
}

func ExecuteTransaction(ctx context.Context, options TransactionOptions, mutate func(*MutationRecorder) error) TransactionResult {
	result := TransactionResult{Rollback: RollbackResult{Status: RollbackNotNeeded, Actions: []string{}, Errors: []string{}}}
	if options.DryRun && options.Commit {
		return fail(result, 2, FailureInvalidOptions, "--dry-run and --commit cannot be combined")
	}
	if len(options.TargetPaths) == 0 {
		return fail(result, 2, FailureInvalidOptions, "at least one exact target path is required")
	}
	repositoryRoot, err := resolveRepositoryRoot(ctx, options.RepositoryRoot)
	if err != nil {
		return fail(result, 2, FailureNotGit, err.Error())
	}
	result.RepositoryRoot = repositoryRoot
	targetPaths, err := normalizeTargetPaths(options.TargetPaths)
	if err != nil {
		return fail(result, 2, FailureInvalidOptions, err.Error())
	}
	states, err := inspectTargets(ctx, repositoryRoot, targetPaths)
	if err != nil {
		return fail(result, 2, FailureInvalidOptions, err.Error())
	}
	for _, state := range states {
		dirty, statusErr := targetIsDirty(ctx, repositoryRoot, state.path)
		if statusErr != nil {
			return fail(result, 2, FailureInvalidOptions, statusErr.Error())
		}
		if dirty {
			return fail(result, 1, FailureDirtyTarget, fmt.Sprintf("target path %q is already dirty", state.path))
		}
		if state.existed && !state.tracked {
			return fail(result, 1, FailureDirtyTarget, fmt.Sprintf("target path %q exists but cannot be restored from Git", state.path))
		}
	}
	if options.Commit {
		empty, indexErr := indexIsEmpty(ctx, repositoryRoot)
		if indexErr != nil {
			return fail(result, 2, FailureInvalidOptions, indexErr.Error())
		}
		if !empty {
			return fail(result, 1, FailureDirtyIndex, "--commit requires an empty existing index")
		}
		if strings.TrimSpace(options.CommitMessage) == "" {
			return fail(result, 2, FailureInvalidOptions, "--commit requires a non-empty commit message")
		}
	}
	if options.DryRun {
		return result
	}
	if mutate == nil {
		return fail(result, 2, FailureInvalidOptions, "mutation callback is required")
	}
	recorder := newMutationRecorder(states)
	if mutationErr := mutate(recorder); mutationErr != nil {
		return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureMutation, mutationErr)
	}
	changedPaths, err := changedTargets(ctx, repositoryRoot, targetPaths)
	if err != nil {
		return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureMutation, err)
	}
	result.ChangedPaths = changedPaths
	for _, path := range changedPaths {
		if _, recorded := recorder.touchedPaths[path]; !recorded {
			return rollbackFailure(ctx, result, repositoryRoot, states, recorder, FailureMutation,
				fmt.Errorf("changed path %q was not recorded by the mutation", path))
		}
	}
	if !options.Commit || len(changedPaths) == 0 {
		return result
	}
	if _, err := runGit(ctx, repositoryRoot, append([]string{"add", "-A", "--"}, changedPaths...)...); err != nil {
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
	if !equalStrings(committedPaths, changedPaths) {
		return committedRisk(result,
			fmt.Sprintf("committed paths %q do not match exact touched paths %q", committedPaths, changedPaths), commitSHA)
	}
	if options.PostCommitVerify != nil {
		if verifyErr := options.PostCommitVerify(ctx, commitSHA); verifyErr != nil {
			return committedRisk(result, verifyErr.Error(), commitSHA)
		}
	}
	return result
}

func newMutationRecorder(states []targetState) *MutationRecorder {
	allowed := make(map[string]struct{}, len(states))
	creatable := make(map[string]struct{}, len(states))
	for _, state := range states {
		allowed[state.path] = struct{}{}
		if !state.existed {
			creatable[state.path] = struct{}{}
		}
	}
	return &MutationRecorder{
		allowedPaths:   allowed,
		creatablePaths: creatable,
		touchedPaths:   map[string]struct{}{},
		createdPaths:   map[string]struct{}{},
	}
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

func indexIsEmpty(ctx context.Context, repositoryRoot string) (bool, error) {
	command := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "diff", "--cached", "--quiet", "--exit-code")
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

func changedTargets(ctx context.Context, repositoryRoot string, paths []string) ([]string, error) {
	changed := make([]string, 0, len(paths))
	for _, path := range paths {
		dirty, err := targetIsDirty(ctx, repositoryRoot, path)
		if err != nil {
			return nil, err
		}
		if dirty {
			changed = append(changed, path)
		}
	}
	return changed, nil
}

func rollbackFailure(ctx context.Context, result TransactionResult, repositoryRoot string, states []targetState, recorder *MutationRecorder, failureKind FailureKind, operationError error) TransactionResult {
	rollback := RollbackResult{Status: RollbackSucceeded, Actions: []string{}, Errors: []string{}}
	for _, state := range states {
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
	if empty, err := indexIsEmpty(ctx, repositoryRoot); err != nil {
		rollback.Errors = append(rollback.Errors, err.Error())
	} else if !empty {
		rollback.Errors = append(rollback.Errors, "Git index is not empty after rollback")
	}
	result.Rollback = rollback
	if len(rollback.Errors) > 0 {
		result.Rollback.Status = RollbackIncomplete
		return fail(result, 4, FailureRollback, operationError.Error())
	}
	return fail(result, 3, failureKind, operationError.Error())
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
	return fail(result, 4, FailureCommittedRisk, reason)
}

func fail(result TransactionResult, exitCode int, kind FailureKind, reason string) TransactionResult {
	result.ExitCode = exitCode
	result.Failure = &TransactionFailure{Kind: kind, Reason: reason}
	return result
}

func runGit(ctx context.Context, repositoryRoot string, args ...string) (string, error) {
	commandArgs := append([]string{"-C", repositoryRoot, "--literal-pathspecs"}, args...)
	command := exec.CommandContext(ctx, "git", commandArgs...)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return string(output), nil
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
