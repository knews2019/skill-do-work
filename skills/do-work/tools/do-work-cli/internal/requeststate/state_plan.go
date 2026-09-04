package requeststate

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/schemanormalization"
)

func BuildPlan(snapshot *repositorymodel.RepositorySnapshot, graph *dependencygraph.DependencyGraph, options StateOptions) StatePlan {
	plan := StatePlan{Transition: options.Transition, Options: options}
	if snapshot == nil {
		plan.Refusal = refuse("STATE-DISCOVERY-FAILED", "repository snapshot is required")
		return plan
	}
	plan.RepositoryRoot = snapshot.RepositoryRoot
	target, targetRefusal := resolvedPlanTarget(snapshot, options)
	if targetRefusal != nil {
		plan.Refusal = targetRefusal
		return plan
	}
	plan.Target = target
	plan.SourcePath = requestPathFor(target)
	plan.ExpectedTargetBytes = append([]byte(nil), target.ContentBytes...)
	if options.Now.IsZero() {
		options.Now = time.Now().UTC()
	}
	options.Now = options.Now.UTC().Truncate(time.Second)
	plan.Options = options
	status := target.TypedRecord.RequestStatus
	if options.OriginalStatus != "" && options.OriginalStatus != status {
		plan.Refusal = refuse("REQUEST-SNAPSHOT-STALE", fmt.Sprintf("expected status %s, found %s", options.OriginalStatus, status), plan.SourcePath)
		return plan
	}
	if transitionRefusal := validateTransition(target, graph, options); transitionRefusal != nil {
		plan.Refusal = transitionRefusal
		return plan
	}

	plan.DestinationPath = plan.SourcePath
	if options.RecordCommitHashOnly {
		planTargets(&plan)
		return plan
	}
	switch options.Transition {
	case TransitionClaim:
		plan.DestinationPath = filepath.ToSlash(filepath.Join("do-work", "working", filepath.Base(target.RelativePath)))
	case TransitionRecover:
		plan.DestinationPath = filepath.ToSlash(filepath.Join("do-work", "queue", filepath.Base(target.RelativePath)))
	case TransitionComplete, TransitionCancel:
		if target.TreeSection == "archive" {
			plan.DestinationPath = plan.SourcePath
		} else if target.TypedRecord.UserRequestID != "" && projectedURClosed(snapshot, target, options) {
			plan.DestinationPath = filepath.ToSlash(filepath.Join("do-work", "archive", target.TypedRecord.UserRequestID, filepath.Base(target.RelativePath)))
			plan.AdditionalMoves = closureMoves(snapshot, target, target.TypedRecord.UserRequestID)
		} else {
			plan.DestinationPath = filepath.ToSlash(filepath.Join("do-work", "archive", filepath.Base(target.RelativePath)))
		}
	case TransitionFail:
		plan.DestinationPath = filepath.ToSlash(filepath.Join("do-work", "archive", filepath.Base(target.RelativePath)))
	}
	if plan.DestinationPath != plan.SourcePath {
		if _, statError := os.Lstat(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(plan.DestinationPath))); statError == nil {
			plan.Refusal = refuse("ARCHIVE-COLLISION", "destination already exists; lifecycle transactions never overwrite", plan.DestinationPath)
			return plan
		} else if !os.IsNotExist(statError) {
			plan.Refusal = refuse("TARGET-INSPECTION-FAILED", statError.Error(), plan.DestinationPath)
			return plan
		}
	}
	for _, move := range plan.AdditionalMoves {
		if _, statError := os.Lstat(filepath.Join(plan.RepositoryRoot, filepath.FromSlash(move.DestinationPath))); statError == nil {
			plan.Refusal = refuse("ARCHIVE-COLLISION", "UR closure destination already exists", move.DestinationPath)
			return plan
		} else if !os.IsNotExist(statError) {
			plan.Refusal = refuse("TARGET-INSPECTION-FAILED", statError.Error(), move.DestinationPath)
			return plan
		}
	}
	planCheckpoint(snapshot, &plan)
	planCalibration(&plan)
	planTargets(&plan)
	return plan
}

func resolvedPlanTarget(snapshot *repositorymodel.RepositorySnapshot, options StateOptions) (*repositorymodel.RequestFile, *StateRefusal) {
	if options.ResolvedTarget == nil {
		return ResolveTarget(snapshot, options.RequestID, options.RequestPath)
	}
	target := options.ResolvedTarget
	found := false
	for _, requestFile := range snapshot.RequestFiles {
		if requestFile == target {
			found = true
			break
		}
	}
	if !found || target.TypedRecord.RequestID != options.RequestID || requestPathFor(target) != options.RequestPath || target.ParseFailure != "" {
		return nil, refuse("REQUEST-SNAPSHOT-STALE", "resolved request no longer belongs to this repository snapshot", options.RequestPath)
	}
	return target, nil
}

func validateTransition(target *repositorymodel.RequestFile, graph *dependencygraph.DependencyGraph, options StateOptions) *StateRefusal {
	status := target.TypedRecord.RequestStatus
	sourcePath := requestPathFor(target)
	switch options.Transition {
	case TransitionClaim:
		if target.TreeSection != "queue" || status != "pending" {
			return refuse("CLAIM-STATUS", "claim requires one pending queue REQ", sourcePath)
		}
		if options.Provenance == "" {
			options.Provenance = ProvenanceDefault
		}
		if options.Provenance != ProvenanceDefault && options.Provenance != ProvenanceExplicit && options.Provenance != ProvenanceURExpanded {
			return refuse("CLAIM-PROVENANCE-INVALID", "claim provenance must be default, explicit-req, or ur-expanded", sourcePath)
		}
		if options.Provenance != ProvenanceExplicit {
			node := graph.NodesByID[target.TypedRecord.RequestID]
			if node == nil || !node.DependenciesSatisfied {
				reason := "dependencies are not satisfied"
				if node != nil {
					reason = fmt.Sprintf("dependencies are not satisfied: unmet=%v missing=%v ambiguous=%v cyclic=%t", node.UnmetDependencies, node.MissingTargets, node.AmbiguousTargets, node.IsCyclic)
				}
				return refuse("CLAIM-DEPENDENCY-FAILED", reason, sourcePath)
			}
		}
	case TransitionRecover:
		if target.TreeSection != "working" || (status != "claimed" && (status != "blocked" || strings.TrimSpace(target.TypedRecord.FieldEvidenceByName["blocked_by"].ScalarValue) == "")) {
			return refuse("RECOVER-CLAIM-STATUS", "recover-claim requires one claimed working REQ, or a blocked working REQ with blocked_by", sourcePath)
		}
		if !options.AssumeSoleWriter {
			return refuse("RECOVER-CLAIM-AUTHORITY", "recover-claim requires the invocation-scoped --assume-sole-writer assertion", sourcePath)
		}
		if !options.DryRun && !options.Commit {
			return refuse("RECOVER-CLAIM-COMMIT-REQUIRED", "recover-claim must commit its exact ownership-transfer paths before selection resumes", sourcePath)
		}
		evidenceModes := 0
		if strings.TrimSpace(options.CheckpointWriter) != "" {
			evidenceModes++
		}
		if options.CheckpointUnlabeled {
			evidenceModes++
		}
		if options.CheckpointAbsent {
			evidenceModes++
		}
		if options.CheckpointAllEntries {
			evidenceModes++
		}
		if evidenceModes != 1 {
			return refuse("RECOVER-CLAIM-CHECKPOINT-EVIDENCE", "recover-claim requires exactly one structural checkpoint evidence mode", sourcePath)
		}
	case TransitionUnblock:
		if target.TreeSection != "queue" || status != "blocked" {
			return refuse("UNBLOCK-STATUS", "unblock requires one blocked queue REQ", sourcePath)
		}
		switch options.UnblockSource {
		case UnblockProbe:
			if options.OriginalStatus != "blocked" || options.ProbeStatus != resultmodel.ProbeSucceeded || !options.UnblockRequired {
				return refuse("UNBLOCK-EVIDENCE", "successful-probe unblock requires the exact REQ-411 evidence", sourcePath)
			}
		case UnblockClarify:
			if !options.CancellationConfirmed {
				return refuse("UNBLOCK-EVIDENCE", "user-via-clarify unblock requires explicit confirmation", sourcePath)
			}
		default:
			return refuse("UNBLOCK-EVIDENCE", "unblock source must be successful-probe or user-via-clarify", sourcePath)
		}
	case TransitionComplete:
		if options.RecordCommitHashOnly {
			if target.TreeSection != "archive" || !schemanormalization.IsTerminalSuccess(status) || !validCommitHash(options.ImplementationHash) {
				return refuse("COMMIT-HASH-PRECONDITION", "commit-hash recording requires an archived terminal-success REQ and an exact implementation hash", sourcePath)
			}
			return nil
		}
		if target.TreeSection != "working" || (status != "claimed" && status != "completed-with-issues") {
			return refuse("COMPLETE-STATUS", "complete requires a claimed working REQ", sourcePath)
		}
		if options.TerminalStatus == "" {
			options.TerminalStatus = "completed"
		}
		if options.TerminalStatus != "completed" && options.TerminalStatus != "completed-with-issues" {
			return refuse("COMPLETE-STATUS", "terminal status must be completed or completed-with-issues", sourcePath)
		}
		if options.ImplementationHash != "" && !validCommitHash(options.ImplementationHash) {
			return refuse("COMMIT-HASH-INVALID", "implementation hash must be 7-40 lowercase hexadecimal characters", sourcePath)
		}
	case TransitionFail:
		if target.TreeSection != "working" || status != "claimed" {
			return refuse("FAIL-STATUS", "fail requires a claimed working REQ", sourcePath)
		}
		if strings.TrimSpace(options.FailureError) == "" {
			return refuse("FAIL-CLASSIFICATION-MISSING", "fail requires an action-classified error", sourcePath)
		}
		if !canonicalFailureType(options.FailureType) {
			return refuse("FAIL-CLASSIFICATION-INVALID", "error_type must be exactly intent, spec, code, or environment", sourcePath)
		}
	case TransitionCancel:
		if !options.CancellationConfirmed || strings.TrimSpace(options.DependentDisposition) == "" {
			return refuse("CANCEL-CONFIRMATION-MISSING", "cancel requires explicit confirmation and dependent disposition", sourcePath)
		}
		if schemanormalization.IsTerminalSuccess(status) || status == "cancelled" {
			return refuse("CANCEL-STATUS", "completed or already-cancelled work is not cancellable", sourcePath)
		}
		if target.TreeSection == "archive" && status != "failed" {
			return refuse("CANCEL-STATUS", "only an archived failed REQ can be cancelled in place", sourcePath)
		}
		if options.DependentDisposition != "leave" && options.DependentDisposition != "repoint" && options.DependentDisposition != "cascade" {
			return refuse("CANCEL-DISPOSITION-INVALID", "dependent disposition must be leave, repoint, or cascade", sourcePath)
		}
		if reasonError := validateOutsideText(options.CancellationReason); reasonError != nil {
			return refuse("CANCEL-REASON-UNSAFE", reasonError.Error(), sourcePath)
		}
		if strings.Contains(options.CancellationReason, "\n") {
			if strings.TrimSpace(options.CancellationSummary) == "" || strings.ContainsAny(options.CancellationSummary, "\r\n") {
				return refuse("CANCEL-REASON-SUMMARY-MISSING", "multiline cancellation reason requires one safe summary line", sourcePath)
			}
			if summaryError := validateOutsideText(options.CancellationSummary); summaryError != nil {
				return refuse("CANCEL-REASON-UNSAFE", summaryError.Error(), sourcePath)
			}
		}
	case TransitionHoldArchiveCollision, TransitionHoldDependencyCycle:
		if target.TreeSection != "queue" || status != "pending" {
			return refuse("HOLD-STATUS", "queue hold requires one pending queue REQ", sourcePath)
		}
	default:
		return refuse("STATE-USAGE", "unknown lifecycle transition", sourcePath)
	}
	return nil
}

func canonicalFailureType(value string) bool {
	switch value {
	case "intent", "spec", "code", "environment":
		return true
	default:
		return false
	}
}

func validateOutsideText(value string) error {
	for _, character := range []byte(value) {
		if (character < 0x20 && character != '\n' && character != '\t') || character == 0x7f {
			return fmt.Errorf("outside text contains unsupported control byte 0x%02x", character)
		}
	}
	return nil
}

func projectedURClosed(snapshot *repositorymodel.RepositorySnapshot, target *repositorymodel.RequestFile, options StateOptions) bool {
	userRequestID := target.TypedRecord.UserRequestID
	if userRequestID == "" {
		return false
	}
	for _, requestFile := range snapshot.RequestFiles {
		if requestFile.TypedRecord.UserRequestID != userRequestID {
			continue
		}
		if requestFile == target {
			continue
		}
		if !schemanormalization.IsTerminalResolved(requestFile.TypedRecord.RequestStatus) {
			return false
		}
	}
	return options.Transition == TransitionComplete || options.Transition == TransitionCancel
}

func closureMoves(snapshot *repositorymodel.RepositorySnapshot, target *repositorymodel.RequestFile, userRequestID string) []FileMove {
	var moves []FileMove
	archivePrefix := filepath.ToSlash(filepath.Join("do-work", "archive", userRequestID))
	for _, requestFile := range snapshot.RequestFiles {
		if requestFile == target || requestFile.TypedRecord.UserRequestID != userRequestID {
			continue
		}
		source := requestPathFor(requestFile)
		destination := filepath.ToSlash(filepath.Join(archivePrefix, filepath.Base(requestFile.RelativePath)))
		if source != destination {
			moves = append(moves, FileMove{SourcePath: source, DestinationPath: destination, ExpectedBytes: append([]byte(nil), requestFile.ContentBytes...)})
		}
	}
	activeURDir := filepath.Join(snapshot.RepositoryRoot, "do-work", "user-requests", userRequestID)
	if info, err := os.Stat(activeURDir); err == nil && info.IsDir() {
		_ = filepath.WalkDir(activeURDir, func(path string, entry os.DirEntry, err error) error {
			if err != nil || entry.IsDir() {
				return nil
			}
			rel, relErr := filepath.Rel(activeURDir, path)
			if relErr != nil {
				return nil
			}
			source := filepath.ToSlash(filepath.Join("do-work", "user-requests", userRequestID, rel))
			destination := filepath.ToSlash(filepath.Join(archivePrefix, rel))
			fileBytes, readErr := os.ReadFile(path)
			if readErr != nil {
				return nil
			}
			moves = append(moves, FileMove{SourcePath: source, DestinationPath: destination, ExpectedBytes: fileBytes})
			return nil
		})
	}
	sort.Slice(moves, func(i, j int) bool { return moves[i].SourcePath < moves[j].SourcePath })
	return moves
}

func planCheckpoint(snapshot *repositorymodel.RepositorySnapshot, plan *StatePlan) {
	needsCheckpoint := plan.Transition == TransitionClaim || (plan.Target.TreeSection == "working" && (plan.Transition == TransitionRecover || plan.Transition == TransitionComplete || plan.Transition == TransitionFail || plan.Transition == TransitionCancel))
	if !needsCheckpoint {
		return
	}
	plan.CheckpointPath = "do-work/CHECKPOINT.md"
	absolutePath := filepath.Join(snapshot.RepositoryRoot, filepath.FromSlash(plan.CheckpointPath))
	existingBytes, readError := os.ReadFile(absolutePath)
	if readError == nil {
		plan.CheckpointExisted = true
	} else if !os.IsNotExist(readError) {
		plan.Refusal = refuse("CHECKPOINT-READ-FAILED", readError.Error(), plan.CheckpointPath)
		return
	}
	if plan.Transition == TransitionRecover {
		planRecoverCheckpoint(existingBytes, plan)
		return
	}
	if plan.Transition != TransitionClaim && !plan.CheckpointExisted {
		plan.CheckpointPath = ""
		plan.SkippedWork = append(plan.SkippedWork, resultmodel.SkippedWork{Code: "CHECKPOINT-NOT-PRESENT", Reason: "there is no checkpoint entry file to remove"})
		return
	}
	writerLabel := plan.Options.WriterLabel
	if writerLabel == "" {
		hostname, _ := os.Hostname()
		writerLabel = hostname + ":" + snapshot.RepositoryRoot
		plan.Options.WriterLabel = writerLabel
	}
	if plan.Transition == TransitionClaim {
		plan.CheckpointBytes = checkpointWithClaim(existingBytes, plan.Target.TypedRecord.RequestID, plan.Target.TypedRecord.RequestTitle, requestmodel.CanonicalTimestamp(plan.Options.Now), writerLabel)
		return
	}
	plan.CheckpointBytes = checkpointWithoutClaim(existingBytes, plan.Target.TypedRecord.RequestID, writerLabel)
	if bytes.Equal(plan.CheckpointBytes, existingBytes) {
		// Nothing this writer may remove: the entry is absent or carries another
		// writer's label, which only a human clears. Targets are declared only for
		// changed bytes, so writing the unchanged file would touch an undeclared
		// path and roll the whole transition back on its own no-op.
		plan.CheckpointPath = ""
		plan.SkippedWork = append(plan.SkippedWork, resultmodel.SkippedWork{Code: "CHECKPOINT-ENTRY-NOT-PRESENT", Reason: "the checkpoint holds no entry this writer may remove for the REQ"})
	}
}

func planRecoverCheckpoint(existingBytes []byte, plan *StatePlan) {
	requestID := plan.Target.TypedRecord.RequestID
	if !plan.CheckpointExisted {
		if !plan.Options.CheckpointAbsent {
			plan.Refusal = refuse("RECOVER-CLAIM-CHECKPOINT-EVIDENCE", "the asserted checkpoint entry does not exist", plan.SourcePath, plan.CheckpointPath)
			return
		}
		plan.CheckpointPath = ""
		plan.SkippedWork = append(plan.SkippedWork, resultmodel.SkippedWork{Code: "CHECKPOINT-ENTRY-NOT-PRESENT", Reason: "the caller explicitly proved no checkpoint entry was present"})
		return
	}
	if plan.Options.CheckpointAbsent {
		if checkpointHasRequestEntry(existingBytes, requestID) {
			plan.Refusal = refuse("RECOVER-CLAIM-CHECKPOINT-EVIDENCE", "checkpoint-absent was asserted but an entry for the REQ exists", plan.SourcePath, plan.CheckpointPath)
			return
		}
		plan.CheckpointPath = ""
		plan.SkippedWork = append(plan.SkippedWork, resultmodel.SkippedWork{Code: "CHECKPOINT-ENTRY-NOT-PRESENT", Reason: "the checkpoint contains no entry for the recovered REQ"})
		return
	}
	updated, removed := checkpointWithoutAuthorizedClaim(existingBytes, requestID, plan.Options.CheckpointWriter, plan.Options.CheckpointUnlabeled)
	if plan.Options.CheckpointAllEntries {
		updated, removed = RemoveAllCheckpointClaims(existingBytes, requestID)
	}
	if !removed {
		plan.Refusal = refuse("RECOVER-CLAIM-CHECKPOINT-EVIDENCE", "the exact asserted checkpoint entry does not exist", plan.SourcePath, plan.CheckpointPath)
		return
	}
	plan.CheckpointBytes = updated
}

func planCalibration(plan *StatePlan) {
	if plan.Transition != TransitionComplete {
		return
	}
	estimateEvidence, found := plan.Target.TypedRecord.FieldEvidenceByName["estimate"]
	if !found || estimateEvidence.NestedValues == nil {
		plan.SkippedWork = append(plan.SkippedWork, resultmodel.SkippedWork{Code: "CALIBRATION-NOT-APPLICABLE", Reason: "request has no estimate block"})
		return
	}
	estimateMinutes, estimateError := strconv.Atoi(estimateEvidence.NestedValues["p50_active_minutes"])
	claimedAt, claimedError := requestmodel.ParseTimestamp(plan.Target.TypedRecord.ClaimedAt)
	if estimateError != nil || claimedError != nil {
		plan.SkippedWork = append(plan.SkippedWork, resultmodel.SkippedWork{Code: "CALIBRATION-INVALID-EVIDENCE", Reason: "estimate or claimed_at is not parseable"})
		return
	}
	completedAt := plan.Options.Now
	wallMinutes := int(completedAt.Sub(claimedAt).Minutes())
	plan.CalibrationPath = "do-work/calibration-log.tsv"
	absolutePath := filepath.Join(plan.RepositoryRoot, filepath.FromSlash(plan.CalibrationPath))
	existingBytes, readError := os.ReadFile(absolutePath)
	if readError == nil {
		plan.CalibrationExisted = true
	} else if !os.IsNotExist(readError) {
		plan.Refusal = refuse("CALIBRATION-READ-FAILED", readError.Error(), plan.CalibrationPath)
		return
	}
	if len(existingBytes) == 0 {
		existingBytes = []byte("req_id\troute\testimated_p50_minutes\twall_minutes\tcompleted_at\n")
	}
	route := plan.Target.TypedRecord.RouteValue
	if route == "" {
		route = "-"
	}
	row := fmt.Sprintf("%s\t%s\t%d\t%d\t%s\n", plan.Target.TypedRecord.RequestID, route, estimateMinutes, wallMinutes, requestmodel.CanonicalTimestamp(completedAt))
	plan.CalibrationBytes = append(existingBytes, []byte(row)...)
}

func planTargets(plan *StatePlan) {
	pathSet := map[string]bool{plan.SourcePath: true}
	if plan.DestinationPath != plan.SourcePath {
		pathSet[plan.DestinationPath] = true
		plan.Changes = append(plan.Changes,
			resultmodel.RecordedChange{Path: plan.SourcePath, Kind: "moved", Detail: "planned source"},
			resultmodel.RecordedChange{Path: plan.DestinationPath, Kind: "created", Detail: "planned destination"})
	} else {
		plan.Changes = append(plan.Changes, resultmodel.RecordedChange{Path: plan.SourcePath, Kind: "modified", Detail: "planned lifecycle update"})
	}
	for _, move := range plan.AdditionalMoves {
		pathSet[move.SourcePath], pathSet[move.DestinationPath] = true, true
		plan.Changes = append(plan.Changes, resultmodel.RecordedChange{Path: move.SourcePath, Kind: "moved", Detail: "planned UR closure source"}, resultmodel.RecordedChange{Path: move.DestinationPath, Kind: "created", Detail: "planned UR closure destination"})
	}
	if plan.CheckpointPath != "" && (!plan.CheckpointExisted || !bytes.Equal(plan.CheckpointBytes, mustReadFile(plan.RepositoryRoot, plan.CheckpointPath))) {
		pathSet[plan.CheckpointPath] = true
		plan.Changes = append(plan.Changes, resultmodel.RecordedChange{Path: plan.CheckpointPath, Kind: map[bool]string{true: "modified", false: "created"}[plan.CheckpointExisted], Detail: "planned checkpoint synchronization"})
	}
	if plan.CalibrationPath != "" {
		pathSet[plan.CalibrationPath] = true
		plan.Changes = append(plan.Changes, resultmodel.RecordedChange{Path: plan.CalibrationPath, Kind: map[bool]string{true: "modified", false: "created"}[plan.CalibrationExisted], Detail: "planned calibration append"})
	}
	for path := range pathSet {
		plan.TargetPaths = append(plan.TargetPaths, path)
	}
	sort.Strings(plan.TargetPaths)
	plan.ExistingUntrackedTargetPaths = existingUntrackedPaths(plan.RepositoryRoot, plan.TargetPaths)
	if plan.Transition == TransitionRecover {
		plan.ExistingDirtyTargetPaths = existingDirtyTrackedPaths(plan.RepositoryRoot, plan.TargetPaths)
	} else if len(plan.Options.AcceptedPreimageDigests) > 0 {
		plan.ExistingDirtyTargetPaths = preimageProvenDirtyPaths(plan.RepositoryRoot, existingDirtyTrackedPaths(plan.RepositoryRoot, plan.TargetPaths), plan.Options.AcceptedPreimageDigests)
	}
	plan.CreatedDirectoryPaths = missingParentDirectories(plan.RepositoryRoot, plan.TargetPaths)
}

func existingDirtyTrackedPaths(repositoryRoot string, paths []string) []string {
	var dirty []string
	for _, path := range paths {
		if command := exec.Command("git", "-C", repositoryRoot, "ls-files", "--error-unmatch", "--", path); command.Run() != nil {
			continue
		}
		output, statusError := exec.Command("git", "-C", repositoryRoot, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--", path).Output()
		if statusError == nil && len(output) > 0 {
			dirty = append(dirty, path)
		}
	}
	return dirty
}

// preimageProvenDirtyPaths keeps only the dirty tracked paths whose current bytes
// hash to the digest a journal recorded for them. Acceptance is by recorded hash,
// never by path class, so a dirty deletion or foreign edit stays refused.
func preimageProvenDirtyPaths(repositoryRoot string, dirtyPaths []string, acceptedDigests map[string]string) []string {
	var proven []string
	for _, path := range dirtyPaths {
		expected, recorded := acceptedDigests[path]
		if !recorded {
			continue
		}
		current, readError := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(path)))
		if readError != nil {
			continue
		}
		digest := sha256.Sum256(current)
		if hex.EncodeToString(digest[:]) == expected {
			proven = append(proven, path)
		}
	}
	return proven
}

func existingUntrackedPaths(repositoryRoot string, paths []string) []string {
	var untracked []string
	for _, path := range paths {
		if _, statError := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(path))); statError != nil {
			continue
		}
		command := exec.Command("git", "-C", repositoryRoot, "--literal-pathspecs", "ls-files", "--error-unmatch", "--", path)
		if command.Run() != nil {
			untracked = append(untracked, path)
		}
	}
	return untracked
}

func missingParentDirectories(repositoryRoot string, paths []string) []string {
	directorySet := map[string]bool{}
	for _, path := range paths {
		for directory := filepath.Dir(filepath.FromSlash(path)); directory != "."; directory = filepath.Dir(directory) {
			slashDirectory := filepath.ToSlash(directory)
			if _, statError := os.Stat(filepath.Join(repositoryRoot, directory)); os.IsNotExist(statError) {
				directorySet[slashDirectory] = true
			}
			if directory == filepath.Dir(directory) {
				break
			}
		}
	}
	var directories []string
	for directory := range directorySet {
		directories = append(directories, directory)
	}
	sort.Slice(directories, func(i, j int) bool { return strings.Count(directories[i], "/") < strings.Count(directories[j], "/") })
	return directories
}

func validCommitHash(value string) bool {
	if len(value) < 7 || len(value) > 40 {
		return false
	}
	for _, character := range value {
		if !(character >= '0' && character <= '9') && !(character >= 'a' && character <= 'f') {
			return false
		}
	}
	return true
}

func mustReadFile(repositoryRoot, relativePath string) []byte {
	contents, _ := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(relativePath)))
	return contents
}
