package doctor

import (
	"context"
	"fmt"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type commandOptions struct {
	repairTimestamps bool
	dryRun           bool
	commit           bool
}

var discoverRepository = repositorymodel.DiscoverRepository

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{"doctor": handleDoctor, "blanked-req-scan": handleDoctor}
}

func handleDoctor(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, parseError := parseCommandOptions(arguments)
	if parseError != nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{
			doctorFinding("DOCTOR-USAGE", resultmodel.SeverityError, nil, nil, parseError.Error(), resultmodel.FixabilityManual,
				"the command line does not express a valid diagnosis or repair", doctorArgv(), doctorJSONArgv()),
		}}
	}
	snapshot, discoveryError := discoverRepository(executionContext.RepositoryRoot)
	if discoveryError != nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{
			doctorFinding("DOCTOR-DISCOVERY-FAILED", resultmodel.SeverityError, nil, nil, discoveryError.Error(), resultmodel.FixabilityManual,
				"the repository snapshot could not be built", doctorArgv(), doctorJSONArgv()),
		}}
	}
	now := time.Now().UTC().Truncate(time.Second)
	if !options.repairTimestamps {
		return ScanRepository(context.Background(), snapshot, ScanOptions{Now: now})
	}
	plans, planFindings := BuildTimestampPlan(context.Background(), snapshot, now)
	repairExecution := applyTimestampPlan(context.Background(), snapshot, plans, RepairOptions{DryRun: options.dryRun, Commit: options.commit})
	repairResult := repairExecution.result
	if options.dryRun || repairResult.Outcome == resultmodel.OutcomeFailure || repairResult.Outcome == resultmodel.OutcomeRolledBack || repairResult.Outcome == resultmodel.OutcomeRisk || repairResult.Outcome == resultmodel.OutcomeRefused {
		repairResult.Findings = append(planFindings, repairResult.Findings...)
		if len(repairResult.Findings) > 0 && repairResult.Outcome == resultmodel.OutcomeSuccess {
			repairResult.Outcome = resultmodel.OutcomeFindings
		}
		sortFindings(repairResult.Findings)
		return repairResult
	}
	updatedSnapshot, rediscoveryError := discoverRepository(executionContext.RepositoryRoot)
	if rediscoveryError != nil {
		nextArgv := doctorArgv()
		affectedPaths := changedPaths(repairResult.Changes)
		stopReason := "repair landed but verification could not rescan"
		if repairExecution.commitSHA != "" {
			repairResult.Outcome = resultmodel.OutcomeRisk
			nextArgv = append([]string(nil), repairExecution.revertArgv...)
			stopReason = "the committed repair could not be verified; revert the exact repair commit before retrying"
		} else {
			repairResult.Outcome = resultmodel.OutcomeFailure
		}
		repairResult.Findings = append(repairResult.Findings, doctorFinding("DOCTOR-POST-REPAIR-DISCOVERY-FAILED", resultmodel.SeverityError,
			nil, affectedPaths, rediscoveryError.Error(), resultmodel.FixabilityManual, stopReason, nextArgv, doctorJSONArgv()))
		return repairResult
	}
	scanResult := ScanRepository(context.Background(), updatedSnapshot, ScanOptions{Now: now})
	scanResult.Changes = repairResult.Changes
	scanResult.Rollback = repairResult.Rollback
	scanResult.Findings = append(repairResult.Findings, scanResult.Findings...)
	if len(scanResult.Findings) > 0 {
		scanResult.Outcome = resultmodel.OutcomeFindings
	}
	return scanResult
}

func changedPaths(changes []resultmodel.RecordedChange) []string {
	paths := []string{}
	seen := map[string]bool{}
	for _, change := range changes {
		if change.Path != "" && !seen[change.Path] {
			seen[change.Path] = true
			paths = append(paths, change.Path)
		}
	}
	return paths
}

func parseCommandOptions(arguments []string) (commandOptions, error) {
	options := commandOptions{}
	for _, argument := range arguments {
		switch argument {
		case "--repair-timestamps":
			options.repairTimestamps = true
		case "--dry-run":
			options.dryRun = true
		case "--commit":
			options.commit = true
		default:
			return options, fmt.Errorf("unknown doctor option %q", argument)
		}
	}
	if (options.dryRun || options.commit) && !options.repairTimestamps {
		return options, fmt.Errorf("--dry-run and --commit are valid only with --repair-timestamps")
	}
	if options.dryRun && options.commit {
		return options, fmt.Errorf("--dry-run and --commit cannot be combined")
	}
	return options, nil
}
