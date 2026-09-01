package corehelpers

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	history, err := gitOutput(executionContext.RepositoryRoot, "log", "-15", "-z", "--format=%H%x00%s")
	if err != nil {
		return usageResult(CommandHandoffSurvey, err.Error())
	}
	historyFields := nonemptyNULFields(history)
	for index := 0; index+1 < len(historyFields); index += 2 {
		findings = append(findings, resultmodel.CommandFinding{Code: "HANDOFF-RECENT-COMMIT", Severity: resultmodel.SeverityInfo, AffectedIDs: []string{historyFields[index]}, Evidence: []string{"subject=" + historyFields[index+1]}, Fixability: resultmodel.FixabilityAutomatic, VerificationArgv: []string{"git", "show", "--stat", historyFields[index]}})
	}
	for _, branchSet := range []struct {
		code string
		flag string
	}{{"HANDOFF-MERGED-BUILDER", "--merged"}, {"HANDOFF-UNMERGED-BUILDER", "--no-merged"}} {
		output, branchError := gitOutput(executionContext.RepositoryRoot, "branch", "--format=%(refname:short)", branchSet.flag, integrationBranch, "--list", "worktree-agent-*")
		if branchError != nil {
			return usageResult(CommandHandoffSurvey, branchError.Error())
		}
		for _, branch := range nonblankLines(output) {
			findings = append(findings, resultmodel.CommandFinding{Code: branchSet.code, Severity: resultmodel.SeverityInfo, AffectedIDs: []string{branch}, Evidence: []string{"branch=" + branch}, Fixability: resultmodel.FixabilityAutomatic, VerificationArgv: []string{"git", "branch", branchSet.flag, integrationBranch, "--list", branch}})
		}
	}
	porcelain, err := gitOutput(executionContext.RepositoryRoot, "worktree", "list", "--porcelain")
	if err != nil {
		return usageResult(CommandHandoffSurvey, err.Error())
	}
	for _, record := range strings.Split(strings.TrimSpace(string(porcelain)), "\n\n") {
		worktreePath, head, branch := "", "", ""
		for _, line := range strings.Split(record, "\n") {
			switch {
			case strings.HasPrefix(line, "worktree "):
				worktreePath = strings.TrimPrefix(line, "worktree ")
			case strings.HasPrefix(line, "HEAD "):
				head = strings.TrimPrefix(line, "HEAD ")
			case strings.HasPrefix(line, "branch "):
				branch = strings.TrimPrefix(line, "branch refs/heads/")
			}
		}
		if worktreePath == "" {
			continue
		}
		findings = append(findings, resultmodel.CommandFinding{Code: "HANDOFF-WORKTREE", Severity: resultmodel.SeverityInfo, AffectedIDs: nonemptyStrings(branch, head), AffectedPaths: []string{worktreePath}, Evidence: []string{"path=" + worktreePath, "head=" + head, "branch=" + branch}, Fixability: resultmodel.FixabilityAutomatic, VerificationArgv: []string{"git", "-C", worktreePath, "rev-parse", "HEAD"}})
		if info, err := os.Stat(worktreePath); err != nil || !info.IsDir() {
			findings = append(findings, helperFinding("HANDOFF-WORKTREE-MISSING", resultmodel.SeverityWarning, []string{worktreePath}, "worktree path is missing or prunable", resultmodel.FixabilityManual, "prune or restore the worktree", []string{"git", "worktree", "prune", "--dry-run"}, nil))
			continue
		}
		status, err := gitOutput(worktreePath, "status", "--porcelain=v1", "-z", "--untracked-files=all")
		if err != nil {
			return usageResult(CommandHandoffSurvey, fmt.Sprintf("status %s: %v", worktreePath, err))
		}
		relative := worktreePath
		if value, err := filepath.Rel(executionContext.RepositoryRoot, worktreePath); err == nil {
			relative = value
		}
		dirtyRows, parseError := parseHandoffStatus(status)
		if parseError != nil {
			return usageResult(CommandHandoffSurvey, fmt.Sprintf("status %s: %v", worktreePath, parseError))
		}
		if len(dirtyRows) == 0 {
			findings = append(findings, helperFinding("HANDOFF-WORKTREE-CLEAN", resultmodel.SeverityInfo, []string{relative}, "clean", resultmodel.FixabilityManual, "", nil, []string{"git", "-C", worktreePath, "status", "--short", "--untracked-files=all"}))
			continue
		}
		for _, row := range dirtyRows {
			evidence := []string{"worktree=" + strconv.Quote(worktreePath), "xy=" + row.status, "path=" + strconv.Quote(row.path)}
			affected := []string{row.path}
			if row.origin != "" {
				evidence = append(evidence, "origin="+strconv.Quote(row.origin))
				affected = append(affected, row.origin)
			}
			findings = append(findings, resultmodel.CommandFinding{Code: "HANDOFF-WORKTREE-DIRTY", Severity: resultmodel.SeverityWarning, AffectedPaths: affected, Evidence: evidence, Fixability: resultmodel.FixabilityManual, AutomationStopReason: "this exact dirty path requires disposition before handoff", NextArgv: []string{"git", "-C", worktreePath, "diff", "--", row.path}, VerificationArgv: []string{"git", "-C", worktreePath, "status", "--short", "--", row.path}})
		}
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: findings}
}

type handoffStatusRow struct {
	status string
	path   string
	origin string
}

func parseHandoffStatus(contents []byte) ([]handoffStatusRow, error) {
	records := bytes.Split(contents, []byte{0})
	rows := []handoffStatusRow{}
	for index := 0; index < len(records); index++ {
		record := records[index]
		if len(record) == 0 {
			continue
		}
		if len(record) < 4 || record[2] != ' ' {
			return nil, fmt.Errorf("short porcelain record")
		}
		row := handoffStatusRow{status: string(record[:2]), path: string(record[3:])}
		if strings.ContainsAny(row.status, "RC") {
			index++
			if index >= len(records) || len(records[index]) == 0 {
				return nil, fmt.Errorf("rename/copy origin missing")
			}
			row.origin = string(records[index])
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func nonemptyNULFields(contents []byte) []string {
	fields := []string{}
	for _, field := range strings.Split(string(contents), "\x00") {
		if field = strings.TrimSpace(field); field != "" {
			fields = append(fields, field)
		}
	}
	return fields
}

func nonemptyStrings(values ...string) []string {
	result := []string{}
	for _, value := range values {
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}
