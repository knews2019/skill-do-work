package doctor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// ApplyUncommittedTimestampPlans is the guarded fail-soft authority for the retained
// repair/audit helpers. Unlike doctor repair it intentionally accepts dirty and untracked
// request files, never stages, and verifies each captured preimage immediately before the
// atomic replacement.
func ApplyUncommittedTimestampPlans(snapshot *repositorymodel.RepositorySnapshot, plans []TimestampRepairPlan) resultmodel.CommandResult {
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess}
	if snapshot == nil {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{
			doctorFinding("DOCTOR-SNAPSHOT-MISSING", resultmodel.SeverityError, nil, nil, "repository snapshot is required", resultmodel.FixabilityManual, "timestamp repair could not start", doctorArgv(), doctorJSONArgv()),
		}}
	}
	result.RepositoryRoot = snapshot.RepositoryRoot
	for _, plan := range plans {
		absolutePath := filepath.Join(snapshot.RepositoryRoot, filepath.FromSlash(plan.RelativePath))
		current, err := os.ReadFile(absolutePath)
		if err != nil || !bytes.Equal(current, plan.ExpectedBytes) {
			evidence := "request changed after inspection"
			if err != nil {
				evidence = err.Error()
			}
			result.Findings = append(result.Findings, doctorFinding("TIMESTAMP-PREIMAGE-CHANGED", resultmodel.SeverityError, nil, []string{plan.RelativePath}, evidence, resultmodel.FixabilityManual, "this file was left byte-identical by the repairer", doctorArgv(), doctorJSONArgv()))
			continue
		}
		if err := atomicfile.ReplaceExisting(absolutePath, plan.UpdatedBytes); err != nil {
			result.Findings = append(result.Findings, doctorFinding("TIMESTAMP-REPLACE-FAILED", resultmodel.SeverityError, nil, []string{plan.RelativePath}, err.Error(), resultmodel.FixabilityManual, "this file was left byte-identical by the repairer", doctorArgv(), doctorJSONArgv()))
			continue
		}
		for _, change := range plan.Changes {
			result.Changes = append(result.Changes, resultmodel.RecordedChange{Path: plan.RelativePath, Kind: "modified", Detail: fmt.Sprintf("repaired %s line %d: %s -> %s (%s)", change.FieldName, change.LineNumber, change.OldValue, change.NewValue, change.Source)})
		}
	}
	if len(result.Findings) > 0 {
		result.Outcome = resultmodel.OutcomeFindings
	}
	return result
}

type RepairOptions struct {
	DryRun bool
	Commit bool
}

func ApplyTimestampPlan(ctx context.Context, snapshot *repositorymodel.RepositorySnapshot, plans []TimestampRepairPlan, options RepairOptions) resultmodel.CommandResult {
	return applyTimestampPlan(ctx, snapshot, plans, options).result
}

type repairExecution struct {
	result     resultmodel.CommandResult
	commitSHA  string
	revertArgv []string
}

func applyTimestampPlan(ctx context.Context, snapshot *repositorymodel.RepositorySnapshot, plans []TimestampRepairPlan, options RepairOptions) repairExecution {
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Rollback: resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded}}
	if snapshot == nil {
		result.Outcome = resultmodel.OutcomeFailure
		result.Findings = []resultmodel.CommandFinding{doctorFinding("DOCTOR-SNAPSHOT-MISSING", resultmodel.SeverityError, nil, nil,
			"repository snapshot is required", resultmodel.FixabilityManual, "timestamp repair could not start", doctorArgv(), doctorJSONArgv())}
		return repairExecution{result: result}
	}
	result.RepositoryRoot = snapshot.RepositoryRoot
	eligiblePlans := []TimestampRepairPlan{}
	for _, plan := range plans {
		preflight := gittransaction.PreflightTargets(ctx, snapshot.RepositoryRoot, []string{plan.RelativePath}, options.Commit)
		if preflight.Failure != nil {
			preflightResult := gittransaction.BuildCommandResult("doctor", gittransaction.TransactionResult{
				Outcome:        resultmodel.OutcomeRefused,
				RepositoryRoot: preflight.RepositoryRoot,
				Failure:        preflight.Failure,
			})
			result.Findings = append(result.Findings, preflightResult.Findings...)
			continue
		}
		eligiblePlans = append(eligiblePlans, plan)
	}
	if len(eligiblePlans) == 0 {
		if len(result.Findings) > 0 {
			result.Outcome = resultmodel.OutcomeRefused
		}
		return repairExecution{result: result}
	}
	targetPaths := make([]string, len(eligiblePlans))
	for index, plan := range eligiblePlans {
		targetPaths[index] = plan.RelativePath
	}
	transactionResult := gittransaction.ExecuteTransaction(ctx, gittransaction.TransactionOptions{
		RepositoryRoot: snapshot.RepositoryRoot, TargetPaths: targetPaths, DryRun: options.DryRun, Commit: options.Commit,
		CommitMessage: "do-work: repair request timestamps",
	}, func(recorder *gittransaction.MutationRecorder) error {
		for _, plan := range eligiblePlans {
			if recordError := recorder.RecordTouched(plan.RelativePath); recordError != nil {
				return recordError
			}
			absolutePath := filepath.Join(snapshot.RepositoryRoot, filepath.FromSlash(plan.RelativePath))
			if replaceError := atomicfile.ReplaceExisting(absolutePath, plan.UpdatedBytes); replaceError != nil {
				return fmt.Errorf("repair %s: %w", plan.RelativePath, replaceError)
			}
		}
		return nil
	})
	if transactionResult.Failure != nil {
		transactionCommandResult := gittransaction.BuildCommandResult("doctor", transactionResult)
		transactionCommandResult.Findings = append(result.Findings, transactionCommandResult.Findings...)
		return repairExecution{result: transactionCommandResult, commitSHA: transactionResult.CommitSHA, revertArgv: transactionResult.RevertArgv}
	}
	result.Rollback = transactionResult.Rollback
	for _, plan := range eligiblePlans {
		for _, change := range plan.Changes {
			detail := fmt.Sprintf("repaired %s line %d: %s -> %s (%s)", change.FieldName, change.LineNumber, change.OldValue, change.NewValue, change.Source)
			if options.DryRun {
				detail = "planned " + detail
			}
			if transactionResult.CommitSHA != "" {
				detail += " committed in " + transactionResult.CommitSHA
			}
			result.Changes = append(result.Changes, resultmodel.RecordedChange{Path: plan.RelativePath, Kind: "modified", Detail: detail})
		}
	}
	if len(result.Findings) > 0 {
		result.Outcome = resultmodel.OutcomeFindings
	}
	return repairExecution{result: result, commitSHA: transactionResult.CommitSHA, revertArgv: append([]string(nil), transactionResult.RevertArgv...)}
}
