package corehelpers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func handleAddExclude(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	probePath, pattern := "", ""
	for index := 0; index < len(arguments); index++ {
		switch {
		case arguments[index] == "--probe-path" || strings.HasPrefix(arguments[index], "--probe-path="):
			value, err := optionValue(arguments, &index, "--probe-path")
			if err != nil {
				return usageResult(CommandAddExclude, err.Error())
			}
			probePath = value
		case arguments[index] == "--pattern" || strings.HasPrefix(arguments[index], "--pattern="):
			value, err := optionValue(arguments, &index, "--pattern")
			if err != nil {
				return usageResult(CommandAddExclude, err.Error())
			}
			pattern = value
		default:
			return usageResult(CommandAddExclude, "unknown option "+arguments[index])
		}
	}
	if probePath == "" {
		return usageResult(CommandAddExclude, "--probe-path is required")
	}
	if strings.ContainsAny(probePath+pattern, "\r\n") {
		return usageResult(CommandAddExclude, "paths and patterns must not contain newlines")
	}
	if pattern == "" {
		pattern = "**/" + strings.TrimPrefix(filepath.ToSlash(probePath), "./")
	}
	gitDirectory, err := gitOutput(executionContext.RepositoryRoot, "rev-parse", "--git-dir")
	if err != nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: []resultmodel.CommandFinding{helperFinding("GIT-EXCLUDE-NOT-A-REPOSITORY", resultmodel.SeverityWarning, nil, "outside a Git repository; no exclude was needed", resultmodel.FixabilityAutomatic, "", nil, nil)}}
	}
	_ = gitDirectory
	check := exec.Command("git", "-C", executionContext.RepositoryRoot, "check-ignore", "-q", "--", probePath)
	if check.Run() == nil {
		return successResult(nil, []resultmodel.CommandFinding{helperFinding("GIT-EXCLUDE-ALREADY-EFFECTIVE", resultmodel.SeverityInfo, []string{probePath}, "existing ignore rules already cover the probe", resultmodel.FixabilityAutomatic, "", nil, []string{"git", "check-ignore", "--", probePath})})
	}
	excludeOutput, err := gitOutput(executionContext.RepositoryRoot, "rev-parse", "--git-path", "info/exclude")
	if err != nil {
		return usageResult(CommandAddExclude, err.Error())
	}
	excludePath := strings.TrimSpace(string(excludeOutput))
	if !filepath.IsAbs(excludePath) {
		excludePath = filepath.Join(executionContext.RepositoryRoot, excludePath)
	}
	contents, _ := os.ReadFile(excludePath)
	for _, line := range strings.Split(string(contents), "\n") {
		if line == pattern {
			return successResult(nil, nil)
		}
	}
	handle, err := os.OpenFile(excludePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if err != nil {
		return usageResult(CommandAddExclude, err.Error())
	}
	prefix := ""
	if len(contents) > 0 && !bytes.HasSuffix(contents, []byte("\n")) {
		prefix = "\n"
	}
	_, writeErr := handle.WriteString(prefix + pattern + "\n")
	closeErr := handle.Close()
	if writeErr == nil {
		writeErr = closeErr
	}
	if writeErr != nil {
		return usageResult(CommandAddExclude, writeErr.Error())
	}
	return successResult([]resultmodel.RecordedChange{{Path: excludePath, Kind: "git-private", Detail: "appended exact local exclude pattern"}}, nil)
}

func handleShowCommitDiff(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	commit, parseResult := requiredPathOption(arguments, "--commit", CommandShowCommitDiff)
	if parseResult != nil {
		return *parseResult
	}
	resolved, err := gitOutput(executionContext.RepositoryRoot, "rev-parse", "--verify", commit+"^{commit}")
	if err != nil {
		return usageResult(CommandShowCommitDiff, "commit does not resolve")
	}
	revision := strings.TrimSpace(string(resolved))
	args := []string{"show", revision}
	if exec.Command("git", "-C", executionContext.RepositoryRoot, "rev-parse", "--verify", "-q", revision+"^2").Run() == nil {
		args = []string{"show", "--first-parent", "-m", revision}
	}
	output, err := gitOutput(executionContext.RepositoryRoot, args...)
	if err != nil {
		return usageResult(CommandShowCommitDiff, err.Error())
	}
	return successResult(nil, []resultmodel.CommandFinding{helperFinding("COMMIT-DIFF", resultmodel.SeverityInfo, nil, string(output), resultmodel.FixabilityAutomatic, "", nil, []string{"git", "show", "--stat", revision})})
}

func handleStageDeletion(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	path, parseResult := requiredPathOption(arguments, "--path", CommandStageDeletion)
	if parseResult != nil {
		return *parseResult
	}
	if exactCachedDeletion(executionContext.RepositoryRoot, path) {
		return successResult(nil, []resultmodel.CommandFinding{helperFinding("DELETION-ALREADY-STAGED", resultmodel.SeverityInfo, []string{path}, "exact cached deletion already exists", resultmodel.FixabilityAutomatic, "", nil, []string{"git", "diff", "--cached", "--name-status", "--", path})})
	}
	command := exec.Command("git", "-C", executionContext.RepositoryRoot, "--literal-pathspecs", "add", "-u", "--", path)
	if output, err := command.CombinedOutput(); err != nil {
		return usageResult(CommandStageDeletion, fmt.Sprintf("git add -u: %s: %v", output, err))
	}
	if !exactCachedDeletion(executionContext.RepositoryRoot, path) {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{helperFinding("DELETION-STAGE-VERIFY-FAILED", resultmodel.SeverityError, []string{path}, "index does not contain one exact deletion", resultmodel.FixabilityManual, "unrelated index state was not accepted as success", nil, []string{"git", "diff", "--cached", "--name-status", "--", path})}}
	}
	return successResult([]resultmodel.RecordedChange{{Path: path, Kind: "staged-deletion", Detail: "staged one exact tracked deletion"}}, nil)
}

func exactCachedDeletion(repositoryRoot, path string) bool {
	output, err := gitOutput(repositoryRoot, "diff", "--cached", "--name-status", "-z", "--", path)
	if err != nil {
		return false
	}
	records := bytes.Split(output, []byte{0})
	return len(records) >= 2 && string(records[0]) == "D" && string(records[1]) == path && len(records) == 3
}
