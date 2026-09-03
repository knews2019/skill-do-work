package requeststate

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
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
		ExistingDirtyTargetPaths:     plan.ExistingDirtyTargetPaths,
		CommitExistingDirtyTargets:   plan.Transition == TransitionRecover && len(plan.ExistingDirtyTargetPaths) > 0,
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
		result.Findings[findingIndex].NextJustRecipe = stateJustRecipe(plan.Transition, plan.Target.TypedRecord.RequestID)
	}
	if transactionResult.Failure == nil && plan.Options.Commit && plan.Transition == TransitionComplete && plan.Options.ImplementationHash == "" {
		return applyCommitMetadata(ctx, plan, transactionResult, result)
	}
	return result
}

func stateJustRecipe(transition Transition, requestID string) string {
	if transition == TransitionRecover {
		return ""
	}
	return "do-work-" + string(transition) + " " + requestID
}

// PlannedPostimages projects the exact lifecycle result before mutation. It is
// consumed by finalization journaling so a restart can distinguish a phase
// that never started from one that completed before its journal update.
func PlannedPostimages(plan StatePlan) ([]PlannedFileImage, error) {
	if !plan.Runnable() {
		return nil, fmt.Errorf("lifecycle plan is not runnable")
	}
	requestBytes, err := lifecycleRequestBytes(plan)
	if err != nil {
		return nil, err
	}
	images := map[string]PlannedFileImage{}
	put := func(path string, exists bool, contents []byte, mode uint32) {
		images[path] = PlannedFileImage{Path: path, Exists: exists, Bytes: append([]byte(nil), contents...), Mode: mode}
	}
	if plan.SourcePath == plan.DestinationPath {
		put(plan.SourcePath, true, requestBytes, 0o644)
	} else {
		put(plan.SourcePath, false, nil, 0)
		put(plan.DestinationPath, true, requestBytes, 0o644)
	}
	for _, move := range plan.AdditionalMoves {
		put(move.SourcePath, false, nil, 0)
		put(move.DestinationPath, true, move.ExpectedBytes, 0o644)
	}
	if plan.CheckpointPath != "" {
		put(plan.CheckpointPath, plan.CheckpointExisted || len(plan.CheckpointBytes) > 0, plan.CheckpointBytes, 0o644)
	}
	if plan.CalibrationPath != "" {
		put(plan.CalibrationPath, plan.CalibrationExisted || len(plan.CalibrationBytes) > 0, plan.CalibrationBytes, 0o644)
	}
	paths := make([]string, 0, len(images))
	for path := range images {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	result := make([]PlannedFileImage, 0, len(paths))
	for _, path := range paths {
		result = append(result, images[path])
	}
	return result, nil
}

var provenanceHashPattern = regexp.MustCompile(`^[0-9a-f]{7,40}$`)

// RecordCommitProvenance is the standalone guarded provenance authority shared by the
// compatibility-shaped command and request-state lifecycle code. It rewrites only the
// top-level commit scalar and refuses lossy or ambiguous documents.
func RecordCommitProvenance(ctx context.Context, repositoryRoot, requestPath, implementationHash string, verifyOnly, dryRun bool) resultmodel.CommandResult {
	failure := func(code, evidence string) resultmodel.CommandResult {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, RepositoryRoot: repositoryRoot, Findings: []resultmodel.CommandFinding{{
			Code: code, Severity: resultmodel.SeverityError, AffectedPaths: nonEmptyPaths(requestPath), Evidence: []string{evidence}, Fixability: resultmodel.FixabilityManual,
			AutomationStopReason: "commit provenance could not be recorded safely",
			NextArgv:             []string{"do-work-cli", "record-commit-hash", "--request-path", requestPath, "--implementation-hash", implementationHash},
			VerificationArgv:     []string{"do-work-cli", "record-commit-hash", "--verify", "--request-path", requestPath, "--implementation-hash", implementationHash},
		}}}
	}
	if requestPath == "" || !provenanceHashPattern.MatchString(implementationHash) {
		return failure("PROVENANCE-USAGE", "request path and a 7-40 character lowercase hex hash are required")
	}
	absolutePath := requestPath
	if !filepath.IsAbs(absolutePath) {
		absolutePath = filepath.Join(repositoryRoot, filepath.FromSlash(requestPath))
	}
	info, err := os.Lstat(absolutePath)
	if err != nil || !info.Mode().IsRegular() {
		return failure("PROVENANCE-PATH-UNSAFE", "request path must be a regular non-symlink file")
	}
	contents, err := os.ReadFile(absolutePath)
	if err != nil || len(contents) == 0 {
		return failure("PROVENANCE-READ-FAILED", "request bytes are empty or unreadable")
	}
	if bytes.Contains(contents, []byte("\r\n")) {
		return failure("PROVENANCE-CRLF-REFUSED", "normalize CRLF before guarded provenance recording")
	}
	document, err := requestmodel.ParseDocument(contents)
	if err != nil {
		return failure("PROVENANCE-FRONTMATTER-INVALID", err.Error())
	}
	record := document.TypedRecord()
	if !strings.HasPrefix(record.RequestID, "REQ-") {
		return failure("PROVENANCE-NOT-REQUEST", "frontmatter id is not a REQ")
	}
	commitEvidence := record.FieldEvidenceByName["commit"]
	if commitEvidence.DuplicateCount > 1 {
		return failure("PROVENANCE-COMMIT-AMBIGUOUS", "frontmatter has more than one commit field")
	}
	if err := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "rev-parse", "--verify", "--quiet", implementationHash+"^{commit}").Run(); err != nil {
		return failure("PROVENANCE-HASH-UNRESOLVED", "implementation hash does not resolve to a commit")
	}
	relativePath := requestPath
	if value, relativeError := filepath.Rel(repositoryRoot, absolutePath); relativeError == nil {
		relativePath = filepath.ToSlash(value)
	}
	if relativePath == ".." || strings.HasPrefix(relativePath, "../") || filepath.IsAbs(relativePath) {
		return failure("PROVENANCE-PATH-OUTSIDE-REPOSITORY", "request path must remain inside the repository")
	}
	if verifyOnly {
		committed, showError := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "cat-file", "blob", "HEAD:"+relativePath).Output()
		if showError != nil {
			return failure("PROVENANCE-VERIFY-UNTRACKED", "HEAD does not contain the request path")
		}
		if !bytes.Equal(committed, contents) {
			return failure("PROVENANCE-VERIFY-BYTES", "committed bytes differ from the worktree")
		}
		if commitEvidence.ScalarValue != implementationHash {
			return failure("PROVENANCE-VERIFY-HASH", "committed frontmatter does not contain the expected hash")
		}
		if verifyError := verifyProvenanceCommitPatch(ctx, repositoryRoot, relativePath, implementationHash); verifyError != nil {
			return failure("PROVENANCE-VERIFY-PATCH", verifyError.Error())
		}
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: repositoryRoot, Findings: []resultmodel.CommandFinding{{Code: "PROVENANCE-VERIFIED", Severity: resultmodel.SeverityInfo, AffectedIDs: []string{record.RequestID}, AffectedPaths: []string{relativePath}, Evidence: []string{"HEAD bytes match the worktree and HEAD changed only the exact top-level commit field"}, Fixability: resultmodel.FixabilityAutomatic, VerificationArgv: []string{"git", "diff", "--unified=0", "HEAD^", "HEAD", "--", relativePath}}}}
	}
	if guardError := guardProvenancePreimage(ctx, repositoryRoot, relativePath, contents); guardError != nil {
		return failure("PROVENANCE-PREIMAGE-GUARD", guardError.Error())
	}
	if commitEvidence.ScalarValue == implementationHash {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: repositoryRoot, Findings: []resultmodel.CommandFinding{{Code: "PROVENANCE-NOOP", Severity: resultmodel.SeverityInfo, AffectedIDs: []string{record.RequestID}, AffectedPaths: []string{relativePath}, Evidence: []string{"expected commit hash is already recorded"}, Fixability: resultmodel.FixabilityAutomatic}}}
	}
	if err := document.SetScalar("commit", implementationHash); err != nil {
		return failure("PROVENANCE-EDIT-FAILED", err.Error())
	}
	updated := document.DocumentBytes()
	updatedDocument, err := requestmodel.ParseDocument(updated)
	if err != nil || updatedDocument.TypedRecord().RequestID != record.RequestID || updatedDocument.TypedRecord().RequestStatus != record.RequestStatus || updatedDocument.TypedRecord().FieldEvidenceByName["commit"].ScalarValue != implementationHash {
		return failure("PROVENANCE-POSTIMAGE-INVALID", "the one-field rewrite did not preserve request identity and status")
	}
	originalWithoutCommit, stripOriginalError := requestmodel.ParseDocument(contents)
	updatedWithoutCommit, stripUpdatedError := requestmodel.ParseDocument(updated)
	if stripOriginalError != nil || stripUpdatedError != nil || originalWithoutCommit.DeleteField("commit") != nil || updatedWithoutCommit.DeleteField("commit") != nil || !bytes.Equal(originalWithoutCommit.DocumentBytes(), updatedWithoutCommit.DocumentBytes()) {
		return failure("PROVENANCE-DELTA-GUARD", "rewrite changed bytes outside the top-level commit field")
	}
	plannedChange := resultmodel.RecordedChange{Path: relativePath, Kind: "modified", Detail: "recorded implementation commit hash"}
	if dryRun {
		plannedChange.Detail = "would record implementation commit hash after all guards"
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: repositoryRoot, Changes: []resultmodel.RecordedChange{plannedChange}, Findings: []resultmodel.CommandFinding{{Code: "PROVENANCE-DRY-RUN", Severity: resultmodel.SeverityInfo, AffectedIDs: []string{record.RequestID}, AffectedPaths: []string{relativePath}, Evidence: []string{"all preimage and one-field delta guards passed; no bytes were written"}, Fixability: resultmodel.FixabilityAutomatic}}}
	}
	if err := atomicfile.ReplaceExisting(absolutePath, updated); err != nil {
		return failure("PROVENANCE-PUBLISH-FAILED", err.Error())
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: repositoryRoot, Changes: []resultmodel.RecordedChange{{Path: relativePath, Kind: "modified", Detail: "recorded implementation commit hash"}}, Findings: []resultmodel.CommandFinding{{Code: "PROVENANCE-RECORDED", Severity: resultmodel.SeverityInfo, AffectedIDs: []string{record.RequestID}, AffectedPaths: []string{relativePath}, Evidence: []string{"only the top-level commit scalar was changed"}, Fixability: resultmodel.FixabilityAutomatic, VerificationArgv: []string{"do-work-cli", "record-commit-hash", "--verify", "--request-path", relativePath, "--implementation-hash", implementationHash}}}}
}

func guardProvenancePreimage(ctx context.Context, repositoryRoot, relativePath string, contents []byte) error {
	if err := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "rev-parse", "--verify", "--quiet", "HEAD^{commit}").Run(); err != nil {
		return fmt.Errorf("HEAD does not resolve; Git object guards cannot run")
	}
	tracked := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "ls-files", "--error-unmatch", "--", relativePath).Run() == nil
	if !tracked {
		return nil
	}
	blob := "HEAD:" + relativePath
	if err := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "cat-file", "-e", blob).Run(); err != nil {
		return nil
	}
	sizeOutput, err := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "cat-file", "-s", blob).Output()
	if err != nil {
		return fmt.Errorf("HEAD blob exists but its size could not be read")
	}
	var size int64
	if _, err := fmt.Sscan(strings.TrimSpace(string(sizeOutput)), &size); err != nil || size < 0 {
		return fmt.Errorf("HEAD blob size is not a valid integer")
	}
	if size > 0 && int64(len(contents))*2 < size {
		return fmt.Errorf("worktree request is less than half the committed size")
	}
	numstat, err := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "diff", "--numstat", "--no-renames", "HEAD", "--", relativePath).Output()
	if err != nil {
		return fmt.Errorf("Git numstat preimage guard could not run")
	}
	if line := strings.TrimSpace(string(numstat)); line != "" {
		fields := strings.Split(line, "\t")
		if len(fields) < 3 || !decimalString(fields[0]) || !decimalString(fields[1]) {
			return fmt.Errorf("Git reports a binary or malformed pending delta")
		}
	}
	return nil
}

func verifyProvenanceCommitPatch(ctx context.Context, repositoryRoot, relativePath, implementationHash string) error {
	for _, revision := range []string{"HEAD^{commit}", "HEAD^"} {
		if err := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "rev-parse", "--verify", "--quiet", revision).Run(); err != nil {
			return fmt.Errorf("%s does not resolve; metadata patch cannot be isolated", revision)
		}
	}
	if exec.CommandContext(ctx, "git", "-C", repositoryRoot, "rev-parse", "--verify", "--quiet", "HEAD^2").Run() == nil {
		return fmt.Errorf("HEAD is a merge commit, not an isolatable metadata commit")
	}
	parentBlob := "HEAD^:" + relativePath
	parent, err := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "cat-file", "blob", parentBlob).Output()
	if err != nil {
		return fmt.Errorf("request does not have a readable parent blob")
	}
	if _, err := requestmodel.ParseDocument(parent); err != nil {
		return fmt.Errorf("parent request frontmatter is not safely parseable: %w", err)
	}
	parentCommitLine, err := topLevelFieldLine(parent, "commit")
	if err != nil {
		return err
	}
	patch, err := exec.CommandContext(ctx, "git", "-C", repositoryRoot, "--no-pager", "diff", "--unified=0", "--no-color", "--no-ext-diff", "--no-textconv", "--no-renames", "HEAD^", "HEAD", "--", relativePath).Output()
	if err != nil || len(patch) == 0 {
		return fmt.Errorf("HEAD does not expose a readable metadata patch for the request")
	}
	added, removed := netPatchLines(patch)
	if len(added) != 1 || added[0] != "commit: "+implementationHash {
		return fmt.Errorf("metadata commit net-added %d line(s), not the exact expected commit field", len(added))
	}
	if parentCommitLine == "" {
		if len(removed) != 0 {
			return fmt.Errorf("metadata insert removed %d unexpected line(s)", len(removed))
		}
	} else if len(removed) != 1 || removed[0] != parentCommitLine {
		return fmt.Errorf("metadata replacement did not remove exactly HEAD^'s frontmatter commit field")
	}
	return nil
}

func topLevelFieldLine(contents []byte, field string) (string, error) {
	lines := strings.Split(strings.TrimSuffix(string(contents), "\n"), "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return "", fmt.Errorf("request has no closed frontmatter")
	}
	for index := 1; index < len(lines); index++ {
		if lines[index] == "---" {
			return "", nil
		}
		if strings.HasPrefix(lines[index], field+":") {
			return lines[index], nil
		}
	}
	return "", fmt.Errorf("request frontmatter is not closed")
}

func netPatchLines(patch []byte) ([]string, []string) {
	added, removed := []string{}, []string{}
	inHunk := false
	for _, line := range strings.Split(string(patch), "\n") {
		if strings.HasPrefix(line, "@@") {
			inHunk = true
			continue
		}
		if !inHunk || line == "" {
			continue
		}
		if line[0] == '+' {
			added = append(added, line[1:])
		} else if line[0] == '-' {
			removed = append(removed, line[1:])
		}
	}
	return cancelMatchingLines(added, removed)
}

func cancelMatchingLines(added, removed []string) ([]string, []string) {
	removedCounts := map[string]int{}
	for _, line := range removed {
		removedCounts[line]++
	}
	netAdded := []string{}
	for _, line := range added {
		if removedCounts[line] > 0 {
			removedCounts[line]--
		} else {
			netAdded = append(netAdded, line)
		}
	}
	addedCounts := map[string]int{}
	for _, line := range added {
		addedCounts[line]++
	}
	netRemoved := []string{}
	for _, line := range removed {
		if addedCounts[line] > 0 {
			addedCounts[line]--
		} else {
			netRemoved = append(netRemoved, line)
		}
	}
	return netAdded, netRemoved
}

func decimalString(value string) bool {
	if value == "" {
		return false
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

func nonEmptyPaths(path string) []string {
	if path == "" {
		return nil
	}
	return []string{path}
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
	case TransitionRecover:
		blocked := plan.Target.TypedRecord.RequestStatus == "blocked"
		hasScope := markdownSectionExists(document.BodyBytes(), "Scope")
		if !blocked {
			status := "pending"
			if markdownSectionContains(document.BodyBytes(), "Open Questions", regexp.MustCompile(`(?m)^- \[ \]`)) {
				status = "pending-answers"
			}
			if err := document.SetScalar("status", status); err != nil {
				return nil, err
			}
			if err := document.SetScalar("status_changed_at", timestamp); err != nil {
				return nil, err
			}
		}
		for _, field := range []string{"claimed_at", "route", "planning_at", "dispatch_at", "builder_handback_at", "integration_at", "review_at", "remediation_at", "re_review_at", "release_at"} {
			if err := document.DeleteField(field); err != nil {
				return nil, err
			}
		}
		if hasScope {
			if err := document.DeleteField("write_set"); err != nil {
				return nil, err
			}
		}
		return stripGeneratedRecoverySections(document.DocumentBytes())
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
	wantStatus := map[Transition]string{TransitionClaim: "claimed", TransitionRecover: recoveredStatus(plan), TransitionUnblock: "pending", TransitionComplete: plan.Options.TerminalStatus, TransitionFail: "failed", TransitionCancel: "cancelled"}[plan.Transition]
	if wantStatus == "" && plan.Transition == TransitionComplete {
		wantStatus = "completed"
	}
	if status != wantStatus {
		return fmt.Errorf("post-transaction status %s, want %s", status, wantStatus)
	}
	return nil
}

func recoveredStatus(plan StatePlan) string {
	if plan.Target.TypedRecord.RequestStatus == "blocked" {
		return "blocked"
	}
	document, parseError := requestmodel.ParseDocument(plan.Target.ContentBytes)
	if parseError == nil && markdownSectionContains(document.BodyBytes(), "Open Questions", regexp.MustCompile(`(?m)^- \[ \]`)) {
		return "pending-answers"
	}
	return "pending"
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
	nextRecipe := stateJustRecipe(plan.Transition, plan.Options.RequestID)
	return resultmodel.CommandResult{Command: string(plan.Transition), Outcome: resultmodel.OutcomeRefused, RepositoryRoot: plan.RepositoryRoot, Findings: []resultmodel.CommandFinding{{
		Code: plan.Refusal.Code, Severity: resultmodel.SeverityWarning, AffectedIDs: []string{plan.Options.RequestID}, AffectedPaths: plan.Refusal.Paths,
		Evidence: []string{plan.Refusal.Reason}, Fixability: resultmodel.FixabilityRefused, AutomationStopReason: "the lifecycle precondition did not hold",
		NextArgv: []string{"do-work-cli", string(plan.Transition), plan.Options.RequestID}, NextJustRecipe: nextRecipe,
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
		NextArgv: []string{"do-work-cli", string(plan.Transition), plan.Options.RequestID}, NextJustRecipe: stateJustRecipe(plan.Transition, plan.Options.RequestID),
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
	updated, _ := RemoveOwnedCheckpointClaim(existing, requestID, writer)
	return updated
}

// RemoveOwnedCheckpointClaim removes one writer-labelled entry from the real
// In Progress section, including that entry's indented continuation lines.
func RemoveOwnedCheckpointClaim(existing []byte, requestID, writer string) ([]byte, bool) {
	return checkpointWithoutAuthorizedClaim(existing, requestID, writer, false)
}

func checkpointWithoutAuthorizedClaim(existing []byte, requestID, writer string, unlabeled bool) ([]byte, bool) {
	if len(existing) == 0 {
		return existing, false
	}
	lines := strings.Split(string(existing), "\n")
	headingLine, sectionEnd, found := sectionLineBounds(lines, "In Progress (interrupted)")
	if !found {
		return existing, false
	}
	filtered := append([]string(nil), lines[:headingLine+1]...)
	removed := false
	for lineIndex := headingLine + 1; lineIndex < sectionEnd; lineIndex++ {
		line := lines[lineIndex]
		writerMatches := strings.HasSuffix(strings.TrimRight(line, "\r"), "writer: "+writer)
		if unlabeled {
			writerMatches = !strings.Contains(line, " — writer: ")
		}
		if strings.HasPrefix(line, "- "+requestID+":") && writerMatches {
			removed = true
			for lineIndex+1 < sectionEnd && strings.TrimSpace(lines[lineIndex+1]) != "" && (strings.HasPrefix(lines[lineIndex+1], " ") || strings.HasPrefix(lines[lineIndex+1], "\t")) {
				lineIndex++
			}
			continue
		}
		filtered = append(filtered, line)
	}
	filtered = append(filtered, lines[sectionEnd:]...)
	return []byte(strings.Join(filtered, "\n")), removed
}

func checkpointHasRequestEntry(existing []byte, requestID string) bool {
	lines := strings.Split(string(existing), "\n")
	headingLine, sectionEnd, found := sectionLineBounds(lines, "In Progress (interrupted)")
	if !found {
		return false
	}
	for _, line := range lines[headingLine+1 : sectionEnd] {
		if strings.HasPrefix(line, "- "+requestID+":") {
			return true
		}
	}
	return false
}

var generatedRecoveryHeading = regexp.MustCompile(`(?m)^## (Triage|Exploration|Plan|Scope|Pre-Flight|Implementation Summary|Qualification|Testing|Review|Lessons Learned|Orientation|Decisions|Discovered Tasks)[ \t]*\r?$`)
var markdownHeading = regexp.MustCompile(`(?m)^## [^\r\n]+[ \t]*\r?$`)

func stripGeneratedRecoverySections(contents []byte) ([]byte, error) {
	document, parseError := requestmodel.ParseDocument(contents)
	if parseError != nil {
		return nil, parseError
	}
	body := document.BodyBytes()
	matches := generatedRecoveryHeading.FindAllIndex(body, -1)
	for matchIndex := len(matches) - 1; matchIndex >= 0; matchIndex-- {
		start := matches[matchIndex][0]
		end := len(body)
		if next := markdownHeading.FindIndex(body[matches[matchIndex][1]:]); next != nil {
			end = matches[matchIndex][1] + next[0]
		}
		if replaceError := document.ReplaceBodySpan(start, end, nil); replaceError != nil {
			return nil, replaceError
		}
		body = document.BodyBytes()
		matches = generatedRecoveryHeading.FindAllIndex(body, -1)
	}
	return document.DocumentBytes(), nil
}

func markdownSectionExists(body []byte, section string) bool {
	return markdownSectionBytes(body, section) != nil
}

func markdownSectionContains(body []byte, section string, pattern *regexp.Regexp) bool {
	sectionBytes := markdownSectionBytes(body, section)
	return sectionBytes != nil && pattern.Match(sectionBytes)
}

func markdownSectionBytes(body []byte, section string) []byte {
	heading := regexp.MustCompile(`(?m)^## ` + regexp.QuoteMeta(section) + `[ \t]*\r?$`)
	match := heading.FindIndex(body)
	if match == nil {
		return nil
	}
	end := len(body)
	if next := markdownHeading.FindIndex(body[match[1]:]); next != nil {
		end = match[1] + next[0]
	}
	return body[match[0]:end]
}

func appendSectionEntry(contents []byte, section, entry string) []byte {
	text := strings.TrimRight(string(contents), "\n")
	lines := strings.Split(text, "\n")
	_, sectionEnd, found := sectionLineBounds(lines, section)
	if !found {
		heading := "## " + section
		return []byte(text + "\n\n" + heading + "\n\n" + entry + "\n")
	}
	if sectionEnd == len(lines) {
		return []byte(text + "\n\n" + entry + "\n")
	}
	sectionText := strings.TrimRight(strings.Join(lines[:sectionEnd], "\n"), "\n")
	followingText := strings.Join(lines[sectionEnd:], "\n")
	return []byte(sectionText + "\n\n" + entry + "\n\n" + followingText + "\n")
}

func sectionLineBounds(lines []string, section string) (int, int, bool) {
	heading := "## " + section
	for lineIndex, line := range lines {
		if line != heading {
			continue
		}
		for sectionEnd := lineIndex + 1; sectionEnd < len(lines); sectionEnd++ {
			if strings.HasPrefix(lines[sectionEnd], "## ") {
				return lineIndex, sectionEnd, true
			}
		}
		return lineIndex, len(lines), true
	}
	return 0, 0, false
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
