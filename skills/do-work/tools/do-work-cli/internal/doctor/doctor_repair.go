package doctor

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

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
			result.Findings = append(result.Findings, doctorFinding(gittransaction.FindingCode(preflight.Failure.Kind), resultmodel.SeverityWarning,
				nil, preflight.Failure.Paths, preflight.Failure.Reason, resultmodel.FixabilityRefused,
				"the exact timestamp target is dirty, untracked, or cannot satisfy the commit guard",
				[]string{"git", "status", "--short", "--", plan.RelativePath}, []string{"git", "diff", "--quiet", "--exit-code", "--", plan.RelativePath}))
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
