package finalization

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/corehelpers"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/dependencygraph"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/publication"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requestmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

var afterFinalizationPhase = func(Phase) error { return nil }

func advanceJournal(ctx context.Context, repositoryRoot string, journal *Journal, resumed bool) (result resultmodel.CommandResult) {
	defer func() {
		if journal.Discovered || result.Outcome == resultmodel.OutcomeSuccess || journal.Phase == PhasePrimaryCommitted || journal.Phase == PhaseMetadataCommitted || journal.Phase == PhaseVerified || journal.Phase == PhaseCleanupComplete {
			return
		}
		actions, rollbackErrors := rollbackBeforePrimary(repositoryRoot, journal)
		result.Rollback.Actions = actions
		result.Rollback.Errors = rollbackErrors
		if len(rollbackErrors) == 0 {
			result.Outcome = resultmodel.OutcomeRolledBack
			result.Rollback.Status = resultmodel.RollbackSucceeded
		} else {
			result.Rollback.Status = resultmodel.RollbackIncomplete
			if result.Finalization != nil {
				result.Finalization.ReasonCodes = append(result.Finalization.ReasonCodes, "FINALIZATION-ROLLBACK-INCOMPLETE")
			}
		}
		if result.Finalization != nil {
			result.Finalization.Phase = string(journal.Phase)
		}
		if len(result.Finalizations) == 1 {
			result.Finalizations[0] = *result.Finalization
		}
	}()
	if journal.Phase == PhasePrepared {
		state, err := imageSetState(repositoryRoot, journal.LifecyclePreimages, journal.LifecyclePostimages)
		if err != nil {
			return finalizationFailure(journal, resumed, "FINALIZATION-LIFECYCLE-CONFLICT", err.Error(), imagePaths(journal.LifecyclePostimages))
		}
		switch state {
		case "pre":
			plan, planError := lifecyclePlan(repositoryRoot, journal.Manifest)
			if planError != nil {
				return finalizationFailure(journal, resumed, "FINALIZATION-LIFECYCLE-PLAN", planError.Error(), nil)
			}
			applied := requeststate.ApplyPlan(ctx, plan)
			if applied.Outcome != resultmodel.OutcomeSuccess {
				return resultFromNestedFailure(journal, resumed, "FINALIZATION-LIFECYCLE-APPLY", applied)
			}
		case "mixed":
			if err := convergeImages(repositoryRoot, journal.LifecyclePreimages, journal.LifecyclePostimages); err != nil {
				return finalizationFailure(journal, resumed, "FINALIZATION-LIFECYCLE-RECOVERY", err.Error(), imagePaths(journal.LifecyclePostimages))
			}
		}
		journal.Phase = PhaseLifecycleApplied
		if err := persistPhase(journal); err != nil {
			return finalizationFailure(journal, resumed, "FINALIZATION-JOURNAL-WRITE", err.Error(), nil)
		}
	}

	if journal.Phase == PhaseLifecycleApplied {
		if journal.ReleaseManifest != nil {
			state, err := imageSetState(repositoryRoot, journal.ReleasePreimages, journal.ReleasePostimages)
			if err != nil {
				return finalizationFailure(journal, resumed, "FINALIZATION-RELEASE-CONFLICT", err.Error(), imagePaths(journal.ReleasePostimages))
			}
			switch state {
			case "pre":
				plan := publication.BuildReleasePlan(repositoryRoot, *journal.ReleaseManifest)
				if !plan.Runnable() {
					reason := "release plan is not runnable"
					if plan.Refusal != nil {
						reason = plan.Refusal.Code + ": " + plan.Refusal.Reason
					}
					return finalizationFailure(journal, resumed, "FINALIZATION-RELEASE-PLAN", reason, imagePaths(journal.ReleasePostimages))
				}
				applied := publication.ApplyPlan(ctx, plan, false, false)
				if applied.Outcome != resultmodel.OutcomeSuccess {
					return resultFromNestedFailure(journal, resumed, "FINALIZATION-RELEASE-APPLY", applied)
				}
			case "mixed":
				if err := convergeImages(repositoryRoot, journal.ReleasePreimages, journal.ReleasePostimages); err != nil {
					return finalizationFailure(journal, resumed, "FINALIZATION-RELEASE-RECOVERY", err.Error(), imagePaths(journal.ReleasePostimages))
				}
			}
			if err := stampReleaseAt(repositoryRoot, journal.ArchivedPath, journal.Manifest.ReleaseAt); err != nil {
				return finalizationFailure(journal, resumed, "FINALIZATION-RELEASE-STAMP", err.Error(), []string{journal.ArchivedPath})
			}
		}
		journal.Phase = PhaseReleaseApplied
		if err := persistPhase(journal); err != nil {
			return finalizationFailure(journal, resumed, "FINALIZATION-JOURNAL-WRITE", err.Error(), nil)
		}
	}

	if journal.Phase == PhaseReleaseApplied {
		if journal.PreparedHead == "" || journal.PreparedDiffSHA256 == "" {
			preparedHead, preparedDiff, err := preparedCommitIdentity(repositoryRoot, journal.EffectiveCommitPaths)
			if err != nil {
				return finalizationFailure(journal, resumed, "FINALIZATION-PREPARED-IDENTITY", err.Error(), journal.EffectiveCommitPaths)
			}
			journal.PreparedHead, journal.PreparedDiffSHA256 = preparedHead, preparedDiff
			if err := writeJournal(journal); err != nil {
				return finalizationFailure(journal, resumed, "FINALIZATION-JOURNAL-WRITE", err.Error(), nil)
			}
		}
		if recoveredSHA, ok := matchingHeadCommit(repositoryRoot, journal); ok {
			journal.PrimaryCommit = recoveredSHA
		} else {
			if blockedCode, blockedReason, blockedPaths := commitSafety(repositoryRoot, journal); blockedCode != "" {
				return finalizationFailure(journal, resumed, blockedCode, blockedReason, blockedPaths)
			}
			transaction := gittransaction.CommitExactPaths(ctx, repositoryRoot, journal.EffectiveCommitPaths, journal.Manifest.CommitMessage, nil)
			if transaction.Failure != nil {
				return finalizationFailure(journal, resumed, "FINALIZATION-PRIMARY-COMMIT", transaction.Failure.Reason, transaction.Failure.Paths)
			}
			journal.PrimaryCommit = transaction.CommitSHA
			journal.CreatedPrimaryCommit = transaction.CommitSHA
		}
		journal.Phase = PhasePrimaryCommitted
		if err := persistPhase(journal); err != nil {
			return finalizationFailure(journal, resumed, "FINALIZATION-JOURNAL-WRITE", err.Error(), nil)
		}
	}

	if journal.Phase == PhasePrimaryCommitted {
		implementationHash := journal.Manifest.ImplementationHash
		if implementationHash == "" {
			implementationHash = journal.PrimaryCommit
		}
		if journal.Manifest.Transition == "complete" && journal.Manifest.ImplementationHash == "" {
			if verified := requeststate.RecordCommitProvenance(ctx, repositoryRoot, journal.ArchivedPath, implementationHash, true, false); verified.Outcome == resultmodel.OutcomeSuccess {
				journal.MetadataCommit = currentHead(repositoryRoot)
			} else {
				recorded := requeststate.RecordCommitProvenance(ctx, repositoryRoot, journal.ArchivedPath, implementationHash, false, false)
				if recorded.Outcome != resultmodel.OutcomeSuccess {
					return resultFromNestedFailure(journal, resumed, "FINALIZATION-PROVENANCE-RECORD", recorded)
				}
				message := fmt.Sprintf("[%s] record commit hash %s", journal.Manifest.RequestID, implementationHash)
				transaction := gittransaction.CommitExactPaths(ctx, repositoryRoot, []string{journal.ArchivedPath}, message, func(context.Context, string) error {
					verified := requeststate.RecordCommitProvenance(ctx, repositoryRoot, journal.ArchivedPath, implementationHash, true, false)
					if verified.Outcome != resultmodel.OutcomeSuccess {
						return fmt.Errorf("metadata commit verification failed")
					}
					return nil
				})
				if transaction.Failure != nil {
					return finalizationFailure(journal, resumed, "FINALIZATION-METADATA-COMMIT", transaction.Failure.Reason, transaction.Failure.Paths)
				}
				journal.MetadataCommit = transaction.CommitSHA
				journal.CreatedMetadataCommit = transaction.CommitSHA
			}
		} else if journal.Manifest.Transition == "complete" {
			if err := verifyRecordedHash(repositoryRoot, journal.ArchivedPath, implementationHash); err != nil {
				return finalizationFailure(journal, resumed, "FINALIZATION-PROVENANCE-VERIFY", err.Error(), []string{journal.ArchivedPath})
			}
		}
		journal.Phase = PhaseMetadataCommitted
		if err := persistPhase(journal); err != nil {
			return finalizationFailure(journal, resumed, "FINALIZATION-JOURNAL-WRITE", err.Error(), nil)
		}
	}

	if journal.Phase == PhaseMetadataCommitted {
		if err := verifyFinalState(repositoryRoot, journal); err != nil {
			return finalizationFailure(journal, resumed, "FINALIZATION-VERIFY", err.Error(), []string{journal.ArchivedPath})
		}
		journal.Phase = PhaseVerified
		if err := persistPhase(journal); err != nil {
			return finalizationFailure(journal, resumed, "FINALIZATION-JOURNAL-WRITE", err.Error(), nil)
		}
	}
	if journal.Phase == PhaseVerified {
		if err := verifyFinalState(repositoryRoot, journal); err != nil {
			return finalizationFailure(journal, resumed, "FINALIZATION-VERIFY", err.Error(), []string{journal.ArchivedPath})
		}
		journal.Phase = PhaseCleanupComplete
		if err := persistPhase(journal); err != nil {
			return finalizationFailure(journal, resumed, "FINALIZATION-JOURNAL-WRITE", err.Error(), nil)
		}
	}
	if journal.Phase != PhaseCleanupComplete {
		return finalizationFailure(journal, resumed, "FINALIZATION-PHASE-INCOMPLETE", "finalization did not reach cleanup", []string{journal.JournalPath})
	}
	if err := verifyFinalState(repositoryRoot, journal); err != nil {
		return finalizationFailure(journal, resumed, "FINALIZATION-VERIFY", err.Error(), []string{journal.ArchivedPath})
	}
	if err := removeJournal(journal); err != nil {
		return finalizationFailure(journal, resumed, "FINALIZATION-JOURNAL-CLEANUP", err.Error(), []string{journal.JournalPath})
	}
	return finalizationSuccess(journal, resumed)
}

func rollbackBeforePrimary(repositoryRoot string, journal *Journal) ([]string, []string) {
	actions := []string{}
	rollbackErrors := []string{}
	if len(journal.ReleasePostimages) > 0 {
		if err := convergeImages(repositoryRoot, journal.ReleasePostimages, journal.ReleasePreimages); err != nil {
			rollbackErrors = append(rollbackErrors, "release rollback: "+err.Error())
		} else {
			actions = append(actions, "restored exact release images")
		}
	}
	if len(rollbackErrors) == 0 && len(journal.LifecyclePostimages) > 0 {
		if err := convergeImages(repositoryRoot, journal.LifecyclePostimages, journal.LifecyclePreimages); err != nil {
			rollbackErrors = append(rollbackErrors, "lifecycle rollback: "+err.Error())
		} else {
			actions = append(actions, "restored exact lifecycle images")
		}
	}
	if len(rollbackErrors) > 0 {
		return actions, rollbackErrors
	}
	ownedPaths := uniqueSorted(append(imagePaths(journal.LifecyclePostimages), imagePaths(journal.ReleasePostimages)...))
	if len(ownedPaths) > 0 {
		arguments := append([]string{"-C", repositoryRoot, "reset", "-q", "HEAD", "--"}, ownedPaths...)
		if output, err := exec.Command("git", arguments...).CombinedOutput(); err != nil {
			return actions, append(rollbackErrors, "index rollback: "+strings.TrimSpace(string(output)))
		}
		actions = append(actions, "cleared finalizer-owned index entries")
	}
	journal.Phase = PhasePrepared
	journal.PreparedHead = ""
	journal.PreparedDiffSHA256 = ""
	if err := writeJournal(journal); err != nil {
		return actions, append(rollbackErrors, "journal rollback: "+err.Error())
	}
	return actions, rollbackErrors
}

func lifecyclePlan(repositoryRoot string, manifest Manifest) (requeststate.StatePlan, error) {
	snapshot, err := repositorymodel.DiscoverRepository(repositoryRoot)
	if err != nil {
		return requeststate.StatePlan{}, err
	}
	completedAt, _ := time.Parse(time.RFC3339, manifest.CompletedAt)
	options := requeststate.StateOptions{RequestID: manifest.RequestID, RequestPath: manifest.RequestPath, WriterLabel: manifest.WriterLabel, Now: completedAt, ImplementationHash: manifest.ImplementationHash}
	if manifest.Transition == "complete" {
		options.Transition = requeststate.TransitionComplete
		options.TerminalStatus = manifest.TerminalStatus
	} else {
		options.Transition = requeststate.TransitionFail
		options.FailureError = manifest.FailureError
		options.FailureType = manifest.FailureType
	}
	plan := requeststate.BuildPlan(snapshot, dependencygraph.BuildGraph(snapshot), options)
	if !plan.Runnable() {
		if plan.Refusal != nil {
			return plan, fmt.Errorf("%s: %s", plan.Refusal.Code, plan.Refusal.Reason)
		}
		return plan, fmt.Errorf("lifecycle plan is not runnable")
	}
	return plan, nil
}

func persistPhase(journal *Journal) error {
	if err := writeJournal(journal); err != nil {
		return err
	}
	return afterFinalizationPhase(journal.Phase)
}

func imageSetState(repositoryRoot string, preimages, postimages []FileImage) (string, error) {
	pre := imagesByPath(preimages)
	post := imagesByPath(postimages)
	if len(pre) != len(post) {
		return "", fmt.Errorf("journal image sets have different path counts")
	}
	preCount, postCount := 0, 0
	for path, before := range pre {
		after, exists := post[path]
		if !exists {
			return "", fmt.Errorf("journal postimage is missing %s", path)
		}
		current, err := currentImage(repositoryRoot, path)
		if err != nil {
			return "", err
		}
		matchesPre := equalImage(current, before)
		matchesPost := equalImage(current, after)
		if !matchesPre && !matchesPost {
			return "", fmt.Errorf("%s matches neither the journal preimage nor postimage", path)
		}
		if matchesPre {
			preCount++
		}
		if matchesPost {
			postCount++
		}
	}
	if postCount == len(post) {
		return "post", nil
	}
	if preCount == len(pre) {
		return "pre", nil
	}
	return "mixed", nil
}

func convergeImages(repositoryRoot string, preimages, postimages []FileImage) error {
	pre := imagesByPath(preimages)
	post := imagesByPath(postimages)
	paths := imagePaths(postimages)
	for _, path := range paths {
		current, err := currentImage(repositoryRoot, path)
		if err != nil {
			return err
		}
		if !equalImage(current, pre[path]) && !equalImage(current, post[path]) {
			return fmt.Errorf("%s changed outside the journal", path)
		}
	}
	// Publish every desired file before removing move sources, so another
	// interruption preserves at least one complete copy.
	for _, path := range paths {
		target := post[path]
		if !target.Exists {
			continue
		}
		current, _ := currentImage(repositoryRoot, path)
		if equalImage(current, target) {
			continue
		}
		absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
			return err
		}
		if current.Exists {
			if err := atomicfile.ReplaceExisting(absolutePath, target.Bytes); err != nil {
				return err
			}
			if err := os.Chmod(absolutePath, os.FileMode(target.Mode)); err != nil {
				return err
			}
		} else {
			mode := os.FileMode(target.Mode)
			if mode == 0 {
				mode = 0o644
			}
			if err := atomicfile.CreateExclusive(absolutePath, target.Bytes, mode); err != nil {
				return err
			}
		}
	}
	for _, path := range paths {
		if post[path].Exists {
			continue
		}
		current, _ := currentImage(repositoryRoot, path)
		if !current.Exists {
			continue
		}
		if !equalImage(current, pre[path]) {
			return fmt.Errorf("refusing to remove changed recovery source %s", path)
		}
		absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(path))
		if err := os.Remove(absolutePath); err != nil {
			return err
		}
	}
	return nil
}

func currentImage(repositoryRoot, path string) (FileImage, error) {
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(path))
	info, err := os.Lstat(absolutePath)
	if os.IsNotExist(err) {
		return FileImage{Path: path}, nil
	}
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return FileImage{}, fmt.Errorf("recovery target is not a regular file: %s", path)
	}
	contents, err := os.ReadFile(absolutePath)
	if err != nil {
		return FileImage{}, err
	}
	return FileImage{Path: path, Exists: true, Bytes: contents, Mode: uint32(info.Mode())}, nil
}

func imagesByPath(images []FileImage) map[string]FileImage {
	result := map[string]FileImage{}
	for _, image := range images {
		result[image.Path] = image
	}
	return result
}

func equalImage(left, right FileImage) bool {
	return left.Exists == right.Exists && (!left.Exists || bytes.Equal(left.Bytes, right.Bytes) && left.Mode == right.Mode)
}

func imagePaths(images []FileImage) []string {
	paths := make([]string, 0, len(images))
	for _, image := range images {
		paths = append(paths, image.Path)
	}
	sort.Strings(paths)
	return paths
}

func stampReleaseAt(repositoryRoot, archivedPath, releaseAt string) error {
	if releaseAt == "" {
		return nil
	}
	absolutePath := filepath.Join(repositoryRoot, filepath.FromSlash(archivedPath))
	contents, err := os.ReadFile(absolutePath)
	if err != nil {
		return err
	}
	document, err := requestmodel.ParseDocument(contents)
	if err != nil {
		return err
	}
	existing := document.TypedRecord().FieldEvidenceByName["release_at"].ScalarValue
	if existing != "" && existing != releaseAt {
		return fmt.Errorf("archived request already carries a different release_at")
	}
	if existing == releaseAt {
		return nil
	}
	if err := document.SetScalar("release_at", releaseAt); err != nil {
		return err
	}
	return atomicfile.ReplaceExisting(absolutePath, document.DocumentBytes())
}

func commitSafety(repositoryRoot string, journal *Journal) (string, string, []string) {
	if err := exec.Command("git", "-C", repositoryRoot, "diff", "--cached", "--quiet", "--exit-code").Run(); err != nil {
		return "FINALIZATION-DIRTY-INDEX", "finalization requires an empty existing index", nil
	}
	rows, err := corehelpers.ReadProtectedInventory(repositoryRoot)
	if err != nil {
		return "FINALIZATION-INVENTORY-FAILED", err.Error(), nil
	}
	allowed := map[string]bool{}
	for _, path := range journal.EffectiveCommitPaths {
		allowed[path] = true
	}
	if journal.Discovered {
		if journalPaths, listError := listJournals(repositoryRoot); listError == nil {
			for _, journalPath := range journalPaths {
				other, readError := readJournal(repositoryRoot, journalPath)
				if readError == nil && other.Discovered {
					for _, path := range other.EffectiveCommitPaths {
						allowed[path] = true
					}
				}
			}
		}
	}
	releasePaths := map[string]bool{}
	for _, image := range journal.ReleasePostimages {
		releasePaths[image.Path] = true
	}
	blocked := []string{}
	for _, row := range rows {
		if row.Classification == "X" && allowed[row.Path] {
			blocked = append(blocked, row.Path)
			continue
		}
		if !allowed[row.Path] && (row.Path == "do-work" || strings.HasPrefix(row.Path, "do-work/") || releasePaths[row.Path]) {
			blocked = append(blocked, row.Path)
		}
	}
	blocked = uniqueSorted(blocked)
	if len(blocked) > 0 {
		return "FINALIZATION-AMBIGUOUS-SHARED-STATE", "shared lifecycle, release, or protected paths remain outside the exact recovery group", blocked
	}
	return "", "", nil
}

func matchingHeadCommit(repositoryRoot string, journal *Journal) (string, bool) {
	commits, err := exec.Command("git", "-C", repositoryRoot, "rev-list", "--reverse", "--ancestry-path", journal.PreparedHead+"..HEAD").Output()
	if err != nil || len(strings.Fields(string(commits))) == 0 {
		return "", false
	}
	allowed := map[string]bool{}
	for _, path := range journal.EffectiveCommitPaths {
		allowed[path] = true
	}
	for _, candidate := range strings.Fields(string(commits)) {
		arguments := append([]string{"-C", repositoryRoot, "diff", "--binary", journal.PreparedHead, candidate, "--"}, journal.EffectiveCommitPaths...)
		diff, diffError := exec.Command("git", arguments...).Output()
		if diffError != nil || digestBytes(diff) != journal.PreparedDiffSHA256 {
			continue
		}
		changed, changedError := exec.Command("git", "-C", repositoryRoot, "diff-tree", "--no-commit-id", "--name-only", "-r", candidate).Output()
		if changedError != nil {
			continue
		}
		exact := true
		for _, path := range strings.Fields(string(changed)) {
			if !allowed[path] {
				exact = false
				break
			}
		}
		if exact {
			return candidate, true
		}
	}
	return "", false
}

func currentHead(repositoryRoot string) string {
	output, _ := exec.Command("git", "-C", repositoryRoot, "rev-parse", "HEAD").Output()
	return strings.TrimSpace(string(output))
}

func verifyRecordedHash(repositoryRoot, archivedPath, expectedHash string) error {
	contents, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(archivedPath)))
	if err != nil {
		return err
	}
	document, err := requestmodel.ParseDocument(contents)
	if err != nil {
		return err
	}
	if document.TypedRecord().FieldEvidenceByName["commit"].ScalarValue != expectedHash {
		return fmt.Errorf("archived request does not record implementation hash %s", expectedHash)
	}
	return nil
}

func verifyFinalState(repositoryRoot string, journal *Journal) error {
	contents, err := os.ReadFile(filepath.Join(repositoryRoot, filepath.FromSlash(journal.ArchivedPath)))
	if err != nil {
		return err
	}
	document, err := requestmodel.ParseDocument(contents)
	if err != nil {
		return err
	}
	record := document.TypedRecord()
	if record.RequestID != journal.Manifest.RequestID {
		return fmt.Errorf("archived request identity changed")
	}
	if journal.Manifest.Transition == "complete" {
		if record.RequestStatus != journal.Manifest.TerminalStatus {
			return fmt.Errorf("archived request status is %s", record.RequestStatus)
		}
		expectedHash := journal.Manifest.ImplementationHash
		if expectedHash == "" {
			expectedHash = journal.PrimaryCommit
		}
		if record.FieldEvidenceByName["commit"].ScalarValue != expectedHash {
			return fmt.Errorf("archived request provenance is incomplete")
		}
		if journal.Manifest.ReleaseAt != "" && record.FieldEvidenceByName["release_at"].ScalarValue != journal.Manifest.ReleaseAt {
			return fmt.Errorf("archived request release_at is incomplete")
		}
	} else if record.RequestStatus != "failed" {
		return fmt.Errorf("failed finalization archived status is %s", record.RequestStatus)
	}
	if len(journal.ReleasePostimages) > 0 {
		releaseImages := make([]FileImage, 0, len(journal.ReleasePostimages))
		for _, image := range journal.ReleasePostimages {
			if image.Path != journal.ArchivedPath {
				releaseImages = append(releaseImages, image)
			}
		}
		state, err := imageSetState(repositoryRoot, releaseImages, releaseImages)
		if err != nil || state != "post" {
			return fmt.Errorf("release postimages are not current")
		}
	}
	return nil
}

func finalizationSuccess(journal *Journal, resumed bool) resultmodel.CommandResult {
	verification := recoveryArgv(journal)
	record := finalizationRecord(journal, resumed, nil, nil)
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: []resultmodel.CommandFinding{{
		Code: "FINALIZATION-COMPLETE", Severity: resultmodel.SeverityInfo, AffectedIDs: []string{journal.Manifest.RequestID},
		AffectedPaths: []string{journal.ArchivedPath}, Evidence: []string{"resumable finalization reached verified cleanup-complete state"},
		Fixability: resultmodel.FixabilityAutomatic, VerificationArgv: verification,
	}}, Finalization: &record, Finalizations: []resultmodel.FinalizationResult{record}}
}

func finalizationFailure(journal *Journal, resumed bool, code, reason string, paths []string) resultmodel.CommandResult {
	verification := recoveryArgv(journal)
	record := finalizationRecord(journal, resumed, paths, []string{code})
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{{
		Code: code, Severity: resultmodel.SeverityError, AffectedIDs: []string{journal.Manifest.RequestID}, AffectedPaths: paths,
		Evidence: []string{reason}, Fixability: resultmodel.FixabilityRefused, AutomationStopReason: "finalization evidence is incomplete or ambiguous",
		NextArgv: verification, VerificationArgv: verification,
	}}, Finalization: &record, Finalizations: []resultmodel.FinalizationResult{record}}
}

func finalizationRecord(journal *Journal, resumed bool, paths, reasonCodes []string) resultmodel.FinalizationResult {
	verification := recoveryArgv(journal)
	terminalStatus := journal.Manifest.TerminalStatus
	if terminalStatus == "" && journal.Manifest.Transition == "fail" {
		terminalStatus = "failed"
	}
	return resultmodel.FinalizationResult{
		RequestID: journal.Manifest.RequestID, RequestPath: journal.Manifest.RequestPath, ArchivePath: journal.ArchivedPath,
		JournalPath: journal.JournalPath, Phase: string(journal.Phase), TerminalStatus: terminalStatus,
		Resumed: resumed, Discovered: journal.Discovered, CommitPaths: append([]string(nil), journal.EffectiveCommitPaths...),
		PrimaryCommit: journal.PrimaryCommit, MetadataCommit: journal.MetadataCommit,
		CreatedPrimaryCommit: journal.CreatedPrimaryCommit, CreatedMetadataCommit: journal.CreatedMetadataCommit,
		BlockedPaths: append([]string(nil), paths...), ReasonCodes: append([]string(nil), reasonCodes...),
		NextArgv: verification, VerificationArgv: verification,
		CollectionArgv: []string{"do-work-cli", "--format", "json", "uncommitted-inventory"},
	}
}

func recoveryArgv(journal *Journal) []string {
	arguments := []string{"do-work-cli", "--format", "json", "recover-finalization"}
	if journal.Discovered {
		arguments = append(arguments, "--discover")
	}
	return arguments
}

func resultFromNestedFailure(journal *Journal, resumed bool, code string, nested resultmodel.CommandResult) resultmodel.CommandResult {
	reason := "nested transaction did not succeed"
	paths := []string{}
	if len(nested.Findings) > 0 {
		if len(nested.Findings[0].Evidence) > 0 {
			reason = nested.Findings[0].Evidence[0]
		}
		paths = nested.Findings[0].AffectedPaths
	}
	return finalizationFailure(journal, resumed, code, reason, paths)
}
