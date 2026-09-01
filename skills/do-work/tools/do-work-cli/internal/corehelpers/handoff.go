package corehelpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func handleHandoffSurvey(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	integrationBranch := ""
	for index := 0; index < len(arguments); index++ {
		if arguments[index] == "--integration-branch" || strings.HasPrefix(arguments[index], "--integration-branch=") {
			value, err := optionValue(arguments, &index, "--integration-branch")
			if err != nil {
				return usageResult(CommandHandoffSurvey, err.Error())
			}
			integrationBranch = value
		} else {
			return usageResult(CommandHandoffSurvey, "unknown option "+arguments[index])
		}
	}
	if integrationBranch == "" {
		for _, candidate := range []string{"main", "master", "trunk"} {
			if _, err := gitOutput(executionContext.RepositoryRoot, "rev-parse", "--verify", "--quiet", "refs/heads/"+candidate); err == nil {
				integrationBranch = candidate
				break
			}
		}
	}
	if integrationBranch == "" {
		return usageResult(CommandHandoffSurvey, "no integration branch found")
	}
	if _, err := gitOutput(executionContext.RepositoryRoot, "rev-parse", "--verify", "--quiet", "refs/heads/"+integrationBranch); err != nil {
		return usageResult(CommandHandoffSurvey, "integration branch does not resolve")
	}
	findings := []resultmodel.CommandFinding{helperFinding("HANDOFF-INTEGRATION-BRANCH", resultmodel.SeverityInfo, nil, integrationBranch, resultmodel.FixabilityAutomatic, "", nil, []string{"git", "rev-parse", "--verify", "refs/heads/" + integrationBranch})}
	for _, query := range []struct {
		code string
		args []string
	}{{"HANDOFF-RECENT-HISTORY", []string{"log", "--oneline", "-15"}}, {"HANDOFF-WORKTREES", []string{"worktree", "list"}}, {"HANDOFF-MERGED-BUILDERS", []string{"branch", "--merged", integrationBranch, "--list", "worktree-agent-*"}}, {"HANDOFF-UNMERGED-BUILDERS", []string{"branch", "--no-merged", integrationBranch, "--list", "worktree-agent-*"}}} {
		output, err := gitOutput(executionContext.RepositoryRoot, query.args...)
		if err != nil {
			return usageResult(CommandHandoffSurvey, err.Error())
		}
		findings = append(findings, helperFinding(query.code, resultmodel.SeverityInfo, nil, string(output), resultmodel.FixabilityAutomatic, "", nil, query.args))
	}
	porcelain, err := gitOutput(executionContext.RepositoryRoot, "worktree", "list", "--porcelain")
	if err != nil {
		return usageResult(CommandHandoffSurvey, err.Error())
	}
	for _, line := range strings.Split(string(porcelain), "\n") {
		if !strings.HasPrefix(line, "worktree ") {
			continue
		}
		worktreePath := strings.TrimPrefix(line, "worktree ")
		if info, err := os.Stat(worktreePath); err != nil || !info.IsDir() {
			findings = append(findings, helperFinding("HANDOFF-WORKTREE-MISSING", resultmodel.SeverityWarning, []string{worktreePath}, "worktree path is missing or prunable", resultmodel.FixabilityManual, "prune or restore the worktree", []string{"git", "worktree", "prune", "--dry-run"}, nil))
			continue
		}
		status, err := gitOutput(worktreePath, "status", "--short", "--untracked-files=all")
		if err != nil {
			return usageResult(CommandHandoffSurvey, fmt.Sprintf("status %s: %v", worktreePath, err))
		}
		code, evidence := "HANDOFF-WORKTREE-CLEAN", "clean"
		severity := resultmodel.SeverityInfo
		if len(status) > 0 {
			code, evidence, severity = "HANDOFF-WORKTREE-DIRTY", string(status), resultmodel.SeverityWarning
		}
		relative := worktreePath
		if value, err := filepath.Rel(executionContext.RepositoryRoot, worktreePath); err == nil {
			relative = value
		}
		findings = append(findings, helperFinding(code, severity, []string{relative}, evidence, resultmodel.FixabilityManual, map[bool]string{true: "dirty files require disposition before handoff", false: ""}[len(status) > 0], nil, []string{"git", "-C", worktreePath, "status", "--short", "--untracked-files=all"}))
	}
	return successResult(nil, findings)
}
