package toolboxcommands

import (
	"context"
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
	transaction := gittransaction.ExecuteTransaction(context.Background(), gittransaction.TransactionOptions{
		RepositoryRoot: repositoryRoot, TargetPaths: targets, CreatedDirectoryPaths: directories,
		DryRun: dryRun, Commit: commit, CommitMessage: message,
	}, mutate)
	return transactionResult(command, transaction, "")
}
