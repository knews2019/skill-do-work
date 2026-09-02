package gittransaction

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// CommitExactPaths commits the current changes at an explicit path allowlist.
// It is intentionally separate from ExecuteTransaction: the caller already
// owns the worktree mutations and needs Git to commit them without treating
// those bytes as a fresh transaction target. Unrelated unstaged work remains
// untouched, while any pre-existing staged entry refuses the operation.
func CommitExactPaths(ctx context.Context, repositoryRoot string, paths []string, message string, postCommitVerify func(context.Context, string) error) TransactionResult {
	result := TransactionResult{
		Outcome:  resultmodel.OutcomeSuccess,
		Rollback: resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded, Actions: []string{}, Errors: []string{}},
	}
	resolvedRoot, err := resolveRepositoryRoot(ctx, repositoryRoot)
	if err != nil {
		return failTransaction(result, resultmodel.OutcomeFailure, repositoryRootFailureKind(err), err.Error())
	}
	result.RepositoryRoot = resolvedRoot
	normalizedPaths, err := normalizeTargetPaths(paths)
	if err != nil || len(normalizedPaths) == 0 {
		if err == nil {
			err = errors.New("at least one exact commit path is required")
		}
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, err.Error())
	}
	if strings.TrimSpace(message) == "" {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, "an exact-path commit requires a non-empty message")
	}
	empty, err := indexIsEmpty(ctx, resolvedRoot)
	if err != nil {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, err.Error())
	}
	if !empty {
		return failTransaction(result, resultmodel.OutcomeRefused, FailureDirtyIndex, "an exact-path commit requires an empty existing index")
	}

	dirtyPaths := make([]string, 0, len(normalizedPaths))
	for _, path := range normalizedPaths {
		dirty, statusError := targetIsDirty(ctx, resolvedRoot, path)
		if statusError != nil {
			return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, statusError.Error(), path)
		}
		if dirty {
			dirtyPaths = append(dirtyPaths, path)
		}
	}
	if len(dirtyPaths) == 0 {
		return failTransaction(result, resultmodel.OutcomeFailure, FailureInvalidOptions, "none of the exact commit paths has a pending change", normalizedPaths...)
	}
	if _, err := runGit(ctx, resolvedRoot, append([]string{"add", "-A", "--"}, dirtyPaths...)...); err != nil {
		unstageExactPaths(ctx, resolvedRoot, dirtyPaths)
		return rolledBackExactCommit(result, FailureCommit, err, dirtyPaths)
	}
	beforeCommit, err := runGit(ctx, resolvedRoot, "rev-parse", "HEAD")
	if err != nil {
		unstageExactPaths(ctx, resolvedRoot, dirtyPaths)
		return rolledBackExactCommit(result, FailureCommit, err, dirtyPaths)
	}
	beforeCommit = strings.TrimSpace(beforeCommit)
	if _, err := runGit(ctx, resolvedRoot, "commit", "-m", message); err != nil {
		afterCommit, afterError := runGit(ctx, resolvedRoot, "rev-parse", "HEAD")
		if afterError == nil && strings.TrimSpace(afterCommit) != beforeCommit {
			return committedRisk(result, "Git reported a commit failure after HEAD advanced", strings.TrimSpace(afterCommit))
		}
		unstageExactPaths(ctx, resolvedRoot, dirtyPaths)
		return rolledBackExactCommit(result, FailureCommit, err, dirtyPaths)
	}
	commitSHA, err := runGit(ctx, resolvedRoot, "rev-parse", "HEAD")
	if err != nil {
		return committedRisk(result, "the exact-path commit succeeded but its ID could not be read", "HEAD")
	}
	commitSHA = strings.TrimSpace(commitSHA)
	result.CommitSHA = commitSHA
	result.RevertArgv = []string{"git", "revert", commitSHA}
	if empty, indexError := indexIsEmpty(ctx, resolvedRoot); indexError != nil {
		return committedRisk(result, indexError.Error(), commitSHA)
	} else if !empty {
		return committedRisk(result, "the Git index is not empty after the exact-path commit", commitSHA)
	}
	committed, err := committedPaths(ctx, resolvedRoot, commitSHA)
	if err != nil {
		return committedRisk(result, err.Error(), commitSHA)
	}
	if !equalStrings(committed, dirtyPaths) {
		return committedRisk(result, fmt.Sprintf("committed paths %q do not match exact dirty paths %q", committed, dirtyPaths), commitSHA)
	}
	if postCommitVerify != nil {
		if err := postCommitVerify(ctx, commitSHA); err != nil {
			return committedRisk(result, err.Error(), commitSHA)
		}
	}
	result.ChangedPaths = append([]string(nil), dirtyPaths...)
	return result
}

func unstageExactPaths(ctx context.Context, repositoryRoot string, paths []string) {
	_, _ = runGit(ctx, repositoryRoot, append([]string{"restore", "--staged", "--"}, paths...)...)
}

func rolledBackExactCommit(result TransactionResult, kind FailureKind, err error, paths []string) TransactionResult {
	result.Outcome = resultmodel.OutcomeRolledBack
	result.Failure = &TransactionFailure{Kind: kind, Reason: err.Error(), Paths: append([]string(nil), paths...)}
	result.Rollback = resultmodel.RollbackResult{Status: resultmodel.RollbackSucceeded, Actions: []string{"restored the exact paths to the pre-commit index state"}, Errors: []string{}}
	return result
}
