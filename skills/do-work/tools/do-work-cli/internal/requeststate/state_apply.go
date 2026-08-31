package requeststate

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func ApplyPlan(ctx context.Context, plan StatePlan) resultmodel.CommandResult {
	if plan.Refusal != nil {
		return refusalResult(plan)
	}
	if !plan.Runnable() {
		return commandFailure(plan.RepositoryRoot, plan.Transition, "STATE-PLAN-INVALID", "lifecycle plan has no target")
	}
	if plan.Options.DryRun {
		return resultmodel.CommandResult{Command: string(plan.Transition), Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: plan.RepositoryRoot,
			Changes: append([]resultmodel.RecordedChange(nil), plan.Changes...), SkippedWork: append([]resultmodel.SkippedWork(nil), plan.SkippedWork...),
			Findings: []resultmodel.CommandFinding{stateSuccessFinding(plan, true)}, Rollback: resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded}}
	}
	transactionResult := gittransaction.ExecuteTransaction(ctx, gittransaction.TransactionOptions{
		RepositoryRoot: plan.RepositoryRoot, TargetPaths: plan.TargetPaths,
		ExistingUntrackedTargetPaths: plan.ExistingUntrackedTargetPaths,
		CreatedDirectoryPaths:        plan.CreatedDirectoryPaths, Commit: plan.Options.Commit,
		CommitMessage:    fmt.Sprintf("[%s] %s request lifecycle", plan.Target.TypedRecord.RequestID, plan.Transition),
		PostCommitVerify: func(context.Context, string) error { return verifyAppliedState(plan) },
	}, func(recorder *gittransaction.MutationRecorder) error {
		for _, directory := range plan.CreatedDirectoryPaths {
			if makeError := os.Mkdir(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(directory)), 0o755); makeError != nil && !os.IsExist(makeError) {
				return makeError
			}
			if recordError := recorder.RecordCreatedDirectory(directory); recordError != nil {
				return recordError
			}
		}
		currentBytes, readError := os.ReadFile(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(plan.SourcePath)))
		if readError != nil || !bytes.Equal(currentBytes, plan.ExpectedTargetBytes) {
			return fmt.Errorf("request snapshot changed before %s", plan.Transition)
		}
		updatedBytes, updateError := lifecycleRequestBytes(plan)
		if updateError != nil {
			return updateError
		}
		if replaceError := atomicfile.ReplaceExisting(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(plan.SourcePath)), updatedBytes); replaceError != nil {
			return replaceError
		}
		if recordError := recorder.RecordTouched(plan.SourcePath); recordError != nil {
			return recordError
		}
		if plan.DestinationPath != plan.SourcePath {
			if renameError := os.Rename(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(plan.SourcePath)), filepath.Join(plan.RepositoryRoot, filepath.FromSlash(plan.DestinationPath))); renameError != nil {
				return renameError
			}
			if recordError := recorder.RecordCreated(plan.DestinationPath); recordError != nil {
				return recordError
			}
		}
		for _, move := range plan.AdditionalMoves {
			currentMoveBytes, moveReadError := os.ReadFile(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(move.SourcePath)))
			if moveReadError != nil || !bytes.Equal(currentMoveBytes, move.ExpectedBytes) {
				return fmt.Errorf("UR closure source changed: %s", move.SourcePath)
			}
			if recordError := recorder.RecordTouched(move.SourcePath); recordError != nil {
				return recordError
			}
			if renameError := os.Rename(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(move.SourcePath)), filepath.Join(plan.RepositoryRoot, filepath.FromSlash(move.DestinationPath))); renameError != nil {
				return renameError
			}
			if recordError := recorder.RecordCreated(move.DestinationPath); recordError != nil {
				return recordError
			}
		}
		removedURDirectories := map[string]bool{}
		for _, move := range plan.AdditionalMoves {
			if strings.HasPrefix(move.SourcePath, "do-work/user-requests/") {
				removedURDirectories[filepath.Dir(move.SourcePath)] = true
			}
		}
		for directory := range removedURDirectories {
			if removeError := os.Remove(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(directory))); removeError != nil {
				return fmt.Errorf("active UR directory was not empty after planned closure: %s: %w", directory, removeError)
			}
		}
		if verifyError := verifyArchivedCalibrationEvidence(plan); verifyError != nil {
			return verifyError
		}
		if writeError := writeCoupledFile(plan.RepositoryRoot, recorder, plan.CheckpointPath, plan.CheckpointBytes, plan.CheckpointExisted); writeError != nil {
			return writeError
		}
		if writeError := writeCoupledFile(plan.RepositoryRoot, recorder, plan.CalibrationPath, plan.CalibrationBytes, plan.CalibrationExisted); writeError != nil {
			return writeError
		}
		return nil
	})
	result := gittransaction.BuildCommandResult(string(plan.Transition), transactionResult)
	result.SkippedWork = append(result.SkippedWork, plan.SkippedWork...)
	if transactionResult.Failure == nil {
		result.Changes = appliedChanges(plan.Changes, transactionResult.CommitSHA)
		result.Findings = append(result.Findings, stateSuccessFinding(plan, false))
	}
	for findingIndex := range result.Findings {
		result.Findings[findingIndex].AffectedIDs = []string{plan.Target.TypedRecord.RequestID}
		result.Findings[findingIndex].NextJustRecipe = "do-work-" + string(plan.Transition) + " " + plan.Target.TypedRecord.RequestID
	}
	if transactionResult.Failure == nil && plan.Options.Commit && plan.Transition == TransitionComplete && plan.Options.ImplementationHash == "" {
		return applyCommitMetadata(ctx, plan, transactionResult, result)
	}
	return result
}

func verifyArchivedCalibrationEvidence(plan StatePlan) error {
	if plan.CalibrationPath == "" {
		return nil
	}
	archivedBytes, readError := os.ReadFile(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(plan.DestinationPath)))
	if readError != nil {
		return fmt.Errorf("re-read archived calibration source: %w", readError)
	}
	document, parseError := requestmodel.ParseDocument(archivedBytes)
	if parseError != nil {
		return parseError
	}
	record := document.TypedRecord()
	claimedAt, claimedError := requestmodel.ParseTimestamp(record.ClaimedAt)
	completedAt, completedError := requestmodel.ParseTimestamp(record.CompletedAt)
	if claimedError != nil || completedError != nil {
		return fmt.Errorf("archived calibration timestamps are not parseable")
	}
	estimate := record.FieldEvidenceByName["estimate"].NestedValues["p50_active_minutes"]
	route := record.RouteValue
	if route == "" {
		route = "-"
	}
	wantRow := fmt.Sprintf("%s\t%s\t%s\t%d\t%s\n", record.RequestID, route, estimate, int(completedAt.Sub(claimedAt).Minutes()), requestmodel.CanonicalTimestamp(completedAt))
	if !bytes.HasSuffix(plan.CalibrationBytes, []byte(wantRow)) {
		return fmt.Errorf("planned calibration row does not match archived lifecycle stamps")
	}
	return nil
}

func lifecycleRequestBytes(plan StatePlan) ([]byte, error) {
	document, parseError := requestmodel.ParseDocument(plan.ExpectedTargetBytes)
	if parseError != nil {
		return nil, parseError
	}
	timestamp := requestmodel.CanonicalTimestamp(plan.Options.Now)
	switch plan.Transition {
	case TransitionClaim:
		if err := document.SetScalar("status", "claimed"); err != nil {
			return nil, err
		}
		if err := document.SetScalar("claimed_at", timestamp); err != nil {
			return nil, err
		}
		if plan.Options.Provenance == ProvenanceExplicit {
			if err := document.DeleteField("assigned_to"); err != nil {
				return nil, err
			}
		}
	case TransitionUnblock:
		condition := plan.Target.TypedRecord.FieldEvidenceByName["blocked_by"].ScalarValue
		if err := document.SetScalar("status", "pending"); err != nil {
			return nil, err
		}
		if err := document.SetScalar("status_changed_at", timestamp); err != nil {
			return nil, err
		}
		if err := document.DeleteField("blocked_by"); err != nil {
			return nil, err
		}
		if err := document.DeleteField("blocked_at"); err != nil {
			return nil, err
		}
		source := "probe"
		if plan.Options.UnblockSource == UnblockClarify {
			source = "user via clarify"
		}
		return appendSectionEntry(document.DocumentBytes(), "Blocked", fmt.Sprintf("- [%s] blocked on %q — cleared by %s", timestamp[:10], condition, source)), nil
	case TransitionComplete:
		if plan.Options.RecordCommitHashOnly {
			if err := document.SetScalar("commit", plan.Options.ImplementationHash); err != nil {
				return nil, err
			}
			return document.DocumentBytes(), nil
		}
		status := plan.Options.TerminalStatus
		if status == "" {
			status = "completed"
		}
		if err := document.SetScalar("status", status); err != nil {
			return nil, err
		}
		if err := document.SetScalar("completed_at", timestamp); err != nil {
			return nil, err
		}
		if plan.Options.ImplementationHash != "" {
			if err := document.SetScalar("commit", plan.Options.ImplementationHash); err != nil {
				return nil, err
			}
		}
	case TransitionFail:
		if err := document.SetScalar("status", "failed"); err != nil {
			return nil, err
		}
		if err := document.SetScalar("completed_at", timestamp); err != nil {
			return nil, err
		}
		if err := document.SetScalar("error", plan.Options.FailureError); err != nil {
			return nil, err
		}
		if err := document.SetScalar("error_type", plan.Options.FailureType); err != nil {
			return nil, err
		}
	case TransitionCancel:
		priorStatus := plan.Target.TypedRecord.RequestStatus
		priorCompleted := plan.Target.TypedRecord.CompletedAt
		if err := document.SetScalar("status", "cancelled"); err != nil {
			return nil, err
		}
		if err := document.SetScalar("completed_at", timestamp); err != nil {
			return nil, err
		}
		entry := fmt.Sprintf("- **When:** %s\n%s\n- **Decided by:** user, via `do-work abandon`", timestamp, cancellationReasonBlock(plan.Options.CancellationReason, plan.Options.CancellationSummary))
		if priorStatus == "failed" {
			prior := "failure instant unrecorded"
			if priorCompleted != "" {
				prior = "failed at " + priorCompleted
			}
			classification := ""
			if errorType, found := plan.Target.TypedRecord.FieldEvidenceByName["error_type"]; found {
				classification = fmt.Sprintf(" (`error_type: %s`)", errorType.ScalarValue)
			}
			entry += "\n- **Previously:** failed" + classification + " — " + prior + " — resolved by decision not to retry"
		}
		return appendSectionEntry(document.DocumentBytes(), "Cancelled", entry), nil
	}
	return document.DocumentBytes(), nil
}

func writeCoupledFile(repositoryRoot string, recorder *gittransaction.MutationRecorder, relativePath string, contents []byte, existed bool) error {
	if relativePath == "" {
		return nil
	}
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(relativePath))
	if existed {
		if replaceError := atomicfile.ReplaceExisting(absolutePath, contents); replaceError != nil {
			return replaceError
		}
		return recorder.RecordTouched(relativePath)
	}
	if createError := atomicfile.CreateExclusive(absolutePath, contents, 0o644); createError != nil {
		return createError
	}
	return recorder.RecordCreated(relativePath)
}

func verifyAppliedState(plan StatePlan) error {
	destinationBytes, readError := os.ReadFile(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(plan.DestinationPath)))
	if readError != nil {
		return readError
	}
	document, parseError := requestmodel.ParseDocument(destinationBytes)
	if parseError != nil {
		return parseError
	}
	expectedBytes, expectedError := lifecycleRequestBytes(plan)
	if expectedError != nil || !bytes.Equal(destinationBytes, expectedBytes) {
		return fmt.Errorf("post-transaction request bytes do not match the lifecycle plan")
	}
	if plan.SourcePath != plan.DestinationPath {
		if _, statError := os.Lstat(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(plan.SourcePath))); !os.IsNotExist(statError) {
			return fmt.Errorf("post-transaction source still exists: %s", plan.SourcePath)
		}
	}
	for _, move := range plan.AdditionalMoves {
		movedBytes, moveError := os.ReadFile(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(move.DestinationPath)))
		if moveError != nil || !bytes.Equal(movedBytes, move.ExpectedBytes) {
			return fmt.Errorf("post-transaction closure destination does not match: %s", move.DestinationPath)
		}
		if _, statError := os.Lstat(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(move.SourcePath))); !os.IsNotExist(statError) {
			return fmt.Errorf("post-transaction closure source still exists: %s", move.SourcePath)
		}
	}
	for path, expected := range map[string][]byte{plan.CheckpointPath: plan.CheckpointBytes, plan.CalibrationPath: plan.CalibrationBytes} {
		if path == "" {
			continue
		}
		actual, coupledError := os.ReadFile(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(path)))
		if coupledError != nil || !bytes.Equal(actual, expected) {
			return fmt.Errorf("post-transaction coupled file does not match: %s", path)
		}
	}
	status := document.TypedRecord().RequestStatus
	if plan.Options.RecordCommitHashOnly {
		commitEvidence, found := document.FieldValue("commit")
		if !found || commitEvidence.ScalarValue != plan.Options.ImplementationHash {
			return fmt.Errorf("post-transaction commit hash does not match %s", plan.Options.ImplementationHash)
		}
		return nil
	}
	wantStatus := map[Transition]string{TransitionClaim: "claimed", TransitionUnblock: "pending", TransitionComplete: plan.Options.TerminalStatus, TransitionFail: "failed", TransitionCancel: "cancelled"}[plan.Transition]
	if wantStatus == "" && plan.Transition == TransitionComplete {
		wantStatus = "completed"
	}
	if status != wantStatus {
		return fmt.Errorf("post-transaction status %s, want %s", status, wantStatus)
	}
	return nil
}

func applyCommitMetadata(ctx context.Context, plan StatePlan, lifecycle gittransaction.TransactionResult, result resultmodel.CommandResult) resultmodel.CommandResult {
	metadataPath := plan.DestinationPath
	metadataPreimage, metadataReadError := os.ReadFile(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(metadataPath)))
	if metadataReadError != nil {
		result.Outcome = resultmodel.OutcomeRisk
		result.Findings = append(result.Findings, resultmodel.CommandFinding{Code: "LIFECYCLE-COMMIT-METADATA-RISK", Severity: resultmodel.SeverityError,
			AffectedIDs: []string{plan.Target.TypedRecord.RequestID}, AffectedPaths: []string{metadataPath}, Evidence: []string{metadataReadError.Error()}, Fixability: resultmodel.FixabilityManual,
			AutomationStopReason: "the lifecycle commit exists but its provenance metadata could not start", NextArgv: []string{"git", "revert", lifecycle.CommitSHA}, VerificationArgv: []string{"git", "show", "--stat", lifecycle.CommitSHA}})
		return result
	}
	metadata := gittransaction.ExecuteTransaction(ctx, gittransaction.TransactionOptions{
		RepositoryRoot: plan.RepositoryRoot, TargetPaths: []string{metadataPath}, Commit: true,
		CommitMessage: fmt.Sprintf("[%s] record commit hash %s", plan.Target.TypedRecord.RequestID, lifecycle.CommitSHA),
		PostCommitVerify: func(context.Context, string) error {
			verificationPlan := plan
			verificationPlan.ExpectedTargetBytes = metadataPreimage
			verificationPlan.Options.RecordCommitHashOnly = true
			verificationPlan.Options.ImplementationHash = lifecycle.CommitSHA
			verificationPlan.DestinationPath = metadataPath
			return verifyAppliedState(verificationPlan)
		},
	}, func(recorder *gittransaction.MutationRecorder) error {
		absolutePath := filepath.Join(plan.RepositoryRoot, filepath.FromSlash(metadataPath))
		contents, readError := os.ReadFile(absolutePath)
		if readError != nil {
			return readError
		}
		document, parseError := requestmodel.ParseDocument(contents)
		if parseError != nil {
			return parseError
		}
		if setError := document.SetScalar("commit", lifecycle.CommitSHA); setError != nil {
			return setError
		}
		if replaceError := atomicfile.ReplaceExisting(absolutePath, document.DocumentBytes()); replaceError != nil {
			return replaceError
		}
		return recorder.RecordTouched(metadataPath)
	})
	if metadata.Failure == nil {
		for findingIndex := range result.Findings {
			result.Findings[findingIndex].Evidence = append(result.Findings[findingIndex].Evidence,
				"lifecycle commit "+lifecycle.CommitSHA, "metadata commit "+metadata.CommitSHA)
		}
		return result
	}
	result.Outcome = resultmodel.OutcomeRisk
	result.Rollback = metadata.Rollback
	revertArgv := []string{"git", "revert"}
	if metadata.CommitSHA != "" {
		revertArgv = append(revertArgv, metadata.CommitSHA)
	}
	revertArgv = append(revertArgv, lifecycle.CommitSHA)
	result.Findings = append(result.Findings, resultmodel.CommandFinding{Code: "LIFECYCLE-COMMIT-METADATA-RISK", Severity: resultmodel.SeverityError,
		AffectedIDs: []string{plan.Target.TypedRecord.RequestID}, AffectedPaths: []string{metadataPath}, Evidence: []string{metadata.Failure.Reason}, Fixability: resultmodel.FixabilityManual,
		AutomationStopReason: "the lifecycle commit exists but its provenance metadata did not complete", NextArgv: revertArgv, VerificationArgv: []string{"git", "show", "--stat", lifecycle.CommitSHA}})
	return result
}

func refusalResult(plan StatePlan) resultmodel.CommandResult {
	return resultmodel.CommandResult{Command: string(plan.Transition), Outcome: resultmodel.OutcomeRefused, RepositoryRoot: plan.RepositoryRoot, Findings: []resultmodel.CommandFinding{{
		Code: plan.Refusal.Code, Severity: resultmodel.SeverityWarning, AffectedIDs: []string{plan.Options.RequestID}, AffectedPaths: plan.Refusal.Paths,
		Evidence: []string{plan.Refusal.Reason}, Fixability: resultmodel.FixabilityRefused, AutomationStopReason: "the lifecycle precondition did not hold",
		NextArgv: []string{"do-work-cli", string(plan.Transition), plan.Options.RequestID}, NextJustRecipe: "do-work-" + string(plan.Transition) + " " + plan.Options.RequestID,
		VerificationArgv: []string{"do-work-cli", "--format", "json", string(plan.Transition), plan.Options.RequestID},
	}}}
}

func commandFailure(repositoryRoot string, transition Transition, code, reason string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Command: string(transition), Outcome: resultmodel.OutcomeFailure, RepositoryRoot: repositoryRoot, Findings: []resultmodel.CommandFinding{{Code: code, Severity: resultmodel.SeverityError,
		Evidence: []string{reason}, Fixability: resultmodel.FixabilityManual, AutomationStopReason: "the lifecycle command could not start safely", NextArgv: []string{"do-work-cli", string(transition)}, VerificationArgv: []string{"do-work-cli", "--format", "json", string(transition)}}}}
}

func stateSuccessFinding(plan StatePlan, dryRun bool) resultmodel.CommandFinding {
	code, evidence := "STATE-APPLIED", "lifecycle transaction applied"
	if dryRun {
		code, evidence = "STATE-DRY-RUN", "lifecycle transaction planned without changing bytes"
	}
	verification := []string{"git", "status", "--short", "--"}
	verification = append(verification, plan.TargetPaths...)
	return resultmodel.CommandFinding{
		Code: code, Severity: resultmodel.SeverityInfo, AffectedIDs: []string{plan.Options.RequestID}, AffectedPaths: append([]string(nil), plan.TargetPaths...),
		Evidence: []string{evidence}, Fixability: resultmodel.FixabilityAutomatic,
		NextArgv: []string{"do-work-cli", string(plan.Transition), plan.Options.RequestID}, NextJustRecipe: "do-work-" + string(plan.Transition) + " " + plan.Options.RequestID,
		VerificationArgv: verification,
	}
}

func checkpointWithClaim(existing []byte, requestID, title, claimedAt, writer string) []byte {
	if len(existing) == 0 {
		existing = []byte("# Session Checkpoint\n\n## In Progress (interrupted)\n")
	}
	without := checkpointWithoutClaim(existing, requestID, writer)
	entry := fmt.Sprintf("- %s: %s — claimed %s — writer: %s", requestID, title, claimedAt, writer)
	return appendSectionEntry(without, "In Progress (interrupted)", entry)
}

func checkpointWithoutClaim(existing []byte, requestID, writer string) []byte {
	if len(existing) == 0 {
		return existing
	}
	lines := strings.Split(string(existing), "\n")
	filtered := lines[:0]
	for _, line := range lines {
		if strings.HasPrefix(line, "- "+requestID+":") && strings.Contains(line, "writer: "+writer) {
			continue
		}
		filtered = append(filtered, line)
	}
	return []byte(strings.Join(filtered, "\n"))
}

func appendSectionEntry(contents []byte, section, entry string) []byte {
	text := strings.TrimRight(string(contents), "\n")
	heading := "## " + section
	position := strings.Index(text, heading)
	if position < 0 {
		return []byte(text + "\n\n" + heading + "\n\n" + entry + "\n")
	}
	sectionEnd := strings.Index(text[position+len(heading):], "\n## ")
	if sectionEnd < 0 {
		return []byte(text + "\n\n" + entry + "\n")
	}
	insertAt := position + len(heading) + sectionEnd
	return []byte(strings.TrimRight(text[:insertAt], "\n") + "\n\n" + entry + "\n" + text[insertAt:] + "\n")
}

func cancellationReasonBlock(reason, summary string) string {
	if strings.TrimSpace(reason) == "" {
		return "- **Why:** no reason given"
	}
	if !strings.Contains(reason, "\n") {
		return "- **Why:** " + reason
	}
	return "- **Why:** " + summary + "\n" + containedOutsideText(reason)
}

func containedOutsideText(value string) string {
	longestRun, currentRun := 0, 0
	for _, character := range value {
		if character == '`' {
			currentRun++
			if currentRun > longestRun {
				longestRun = currentRun
			}
		} else {
			currentRun = 0
		}
	}
	fenceLength := longestRun + 1
	if fenceLength < 3 {
		fenceLength = 3
	}
	fence := strings.Repeat("`", fenceLength)
	var contained strings.Builder
	contained.WriteString("> " + fence + "\n")
	for _, line := range strings.Split(value, "\n") {
		contained.WriteString("> " + line + "\n")
	}
	contained.WriteString("> " + fence)
	return contained.String()
}

func appliedChanges(changes []resultmodel.RecordedChange, commitSHA string) []resultmodel.RecordedChange {
	applied := append([]resultmodel.RecordedChange(nil), changes...)
	for index := range applied {
		applied[index].Detail = "applied"
		if commitSHA != "" {
			applied[index].Detail = "committed in " + commitSHA
		}
	}
	return applied
}
