package cleanup

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type commandOptions struct {
	dryRun           bool
	commit           bool
	restoreBlanked   []string
	discardWorktrees []string
}

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{"cleanup": handleCleanup}
}

func handleCleanup(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, parseError := parseCommandOptions(arguments)
	if parseError != nil {
		return commandFailure(executionContext.RepositoryRoot, "CLEANUP-USAGE", parseError.Error())
	}
	snapshot, discoveryError := repositorymodel.DiscoverRepository(executionContext.RepositoryRoot)
	if discoveryError != nil {
		return commandFailure(executionContext.RepositoryRoot, "CLEANUP-DISCOVERY-FAILED", discoveryError.Error())
	}
	plan := BuildPlan(snapshot)
	for _, warning := range snapshot.WarningMessages {
		plan.Findings = append(plan.Findings, manualFinding("CLEANUP-DISCOVERY-WARNING", nil, nil, warning))
	}
	AddBlankedRepairs(context.Background(), &plan, options.restoreBlanked)
	EnrichDocumentationLinks(context.Background(), &plan)
	result := ApplyPlan(context.Background(), plan, ApplyOptions{DryRun: options.dryRun, Commit: options.commit,
		CommitMessage: "do-work: cleanup safe repository state"})
	if result.Outcome != resultmodel.OutcomeFailure && result.Outcome != resultmodel.OutcomeRolledBack && result.Outcome != resultmodel.OutcomeRisk {
		worktreeChanges, worktreeFindings := ApplyWorktreeRepairs(context.Background(), snapshot.RepositoryRoot, WorktreeRepairOptions{DryRun: options.dryRun, DiscardNames: options.discardWorktrees})
		result.Changes = append(result.Changes, worktreeChanges...)
		result.Findings = append(result.Findings, worktreeFindings...)
		if len(result.Findings) > 0 {
			result.Outcome = resultmodel.OutcomeFindings
		}
	}
	return result
}

func parseCommandOptions(arguments []string) (commandOptions, error) {
	options := commandOptions{}
	for index := 0; index < len(arguments); index++ {
		switch arguments[index] {
		case "--dry-run":
			options.dryRun = true
		case "--commit":
			options.commit = true
		case "--restore-blanked":
			index++
			if index >= len(arguments) {
				return options, fmt.Errorf("--restore-blanked requires an exact repository-relative path")
			}
			path, err := exactTarget(arguments[index])
			if err != nil {
				return options, err
			}
			options.restoreBlanked = append(options.restoreBlanked, path)
		case "--discard-worktree":
			index++
			if index >= len(arguments) || !strings.HasPrefix(arguments[index], "worktree-agent-") || strings.ContainsAny(arguments[index], "/\\") {
				return options, fmt.Errorf("--discard-worktree requires an exact worktree-agent-* branch name")
			}
			options.discardWorktrees = append(options.discardWorktrees, arguments[index])
		default:
			return options, fmt.Errorf("unknown cleanup option %q", arguments[index])
		}
	}
	if options.dryRun && options.commit {
		return options, fmt.Errorf("--dry-run and --commit cannot be combined")
	}
	if options.commit && len(options.discardWorktrees) > 0 {
		return options, fmt.Errorf("--commit and --discard-worktree cannot be combined; discard worktrees in a separate explicit run")
	}
	return options, nil
}

func exactTarget(value string) (string, error) {
	if value == "" || filepath.IsAbs(value) {
		return "", fmt.Errorf("destructive repair targets must be exact repository-relative paths")
	}
	cleaned := filepath.Clean(value)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("destructive repair target %q escapes the repository", value)
	}
	return filepath.ToSlash(cleaned), nil
}

func commandFailure(repositoryRoot, code, evidence string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, RepositoryRoot: repositoryRoot,
		Findings: []resultmodel.CommandFinding{{Code: code, Severity: resultmodel.SeverityError, Evidence: []string{evidence},
			Fixability: resultmodel.FixabilityManual, AutomationStopReason: "cleanup could not start safely",
			NextArgv: []string{"do-work-cli", "cleanup", "--dry-run"}, VerificationArgv: []string{"do-work-cli", "--format", "json", "cleanup", "--dry-run"}}}}
}
