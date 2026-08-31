package cleanup

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type ApplyOptions struct {
	DryRun           bool
	Commit           bool
	CommitMessage    string
	PostCommitVerify func(context.Context, string) error
}

func ApplyPlan(ctx context.Context, plan CleanupPlan, options ApplyOptions) resultmodel.CommandResult {
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, RepositoryRoot: plan.RepositoryRoot,
		Findings: append([]resultmodel.CommandFinding(nil), plan.Findings...), Changes: []resultmodel.RecordedChange{}, SkippedWork: []resultmodel.SkippedWork{},
		Rollback: resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded}}
	if plan.RepositoryRoot == "" {
		result.Outcome = resultmodel.OutcomeFailure
		return result
	}
	groupCodeCounts := map[string]int{}
	groupIndexByCode := map[string]int{}
	for groupIndex, group := range plan.Groups {
		if group.Code == "" {
			continue
		}
		groupCodeCounts[group.Code]++
		groupIndexByCode[group.Code] = groupIndex
	}
	directlyEligible := make([]bool, len(plan.Groups))
	scratchCandidates := make([]bool, len(plan.Groups))
	for groupIndex, group := range plan.Groups {
		if group.Code != "" && groupCodeCounts[group.Code] > 1 {
			result.Findings = append(result.Findings, refusedGroupFinding(group, "", "group code is duplicated in the cleanup plan"))
			continue
		}
		targetPaths := groupTargetPaths(group)
		if collisionPath := existingDestination(plan.RepositoryRoot, group); collisionPath != "" {
			result.Findings = append(result.Findings, refusedGroupFinding(group, collisionPath, "destination already exists; cleanup never overwrites"))
			continue
		}
		preflight := gittransaction.PreflightTargets(ctx, plan.RepositoryRoot, targetPaths, options.Commit)
		if preflight.Failure != nil {
			if preflight.Failure.Kind == gittransaction.FailureDirtyTarget && isUntrackedConsumedScratch(ctx, plan.RepositoryRoot, group) {
				directlyEligible[groupIndex] = true
				scratchCandidates[groupIndex] = true
				continue
			}
			path := ""
			if len(preflight.Failure.Paths) > 0 {
				path = preflight.Failure.Paths[0]
			}
			finding := refusedGroupFinding(group, path, preflight.Failure.Reason)
			if preflight.Failure.Kind == gittransaction.FailureDirtyIndex {
				finding.NextArgv = []string{"git", "diff", "--cached", "--name-only"}
				finding.VerificationArgv = []string{"git", "diff", "--cached", "--quiet", "--exit-code"}
			}
			result.Findings = append(result.Findings, finding)
			continue
		}
		directlyEligible[groupIndex] = true
	}

	dependencyStates := make([]groupDependencyState, len(plan.Groups))
	dependencyBlockers := make([]groupDependencyBlocker, len(plan.Groups))
	for groupIndex := range plan.Groups {
		if !directlyEligible[groupIndex] {
			dependencyStates[groupIndex] = groupDependencyBlocked
			dependencyBlockers[groupIndex] = groupDependencyBlocker{Code: plan.Groups[groupIndex].Code, Reason: "required group did not pass direct preflight"}
		}
	}
	eligibleGroups := make([]OperationGroup, 0, len(plan.Groups))
	scratchGroups := []OperationGroup{}
	proposedChanges := []resultmodel.RecordedChange{}
	for groupIndex, group := range plan.Groups {
		if !directlyEligible[groupIndex] {
			continue
		}
		eligible, blocker := resolveGroupDependencies(groupIndex, plan.Groups, directlyEligible, groupCodeCounts, groupIndexByCode, dependencyStates, dependencyBlockers)
		if !eligible {
			result.Findings = append(result.Findings, dependencyRefusedGroupFinding(group, blocker))
			continue
		}
		if scratchCandidates[groupIndex] {
			scratchGroups = append(scratchGroups, group)
		} else {
			eligibleGroups = append(eligibleGroups, group)
		}
		for _, operation := range group.Operations {
			proposedChanges = append(proposedChanges, changeForOperation(operation, true))
		}
	}
	if len(eligibleGroups) == 0 && len(scratchGroups) == 0 {
		if len(result.Findings) > 0 {
			result.Outcome = resultmodel.OutcomeFindings
		}
		return result
	}
	transactionResult := gittransaction.TransactionResult{Outcome: resultmodel.OutcomeSuccess, Rollback: resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded}}
	if len(eligibleGroups) > 0 {
		targetPaths := groupsTargetPaths(eligibleGroups)
		createdDirectories := plannedCreatedDirectories(plan.RepositoryRoot, eligibleGroups)
		transactionResult = gittransaction.ExecuteTransaction(ctx, gittransaction.TransactionOptions{
			RepositoryRoot: plan.RepositoryRoot, TargetPaths: targetPaths, CreatedDirectoryPaths: createdDirectories,
			DryRun: options.DryRun, Commit: options.Commit, CommitMessage: options.CommitMessage, PostCommitVerify: options.PostCommitVerify,
		}, func(recorder *gittransaction.MutationRecorder) error {
			repositoryHandle, rootError := os.OpenRoot(plan.RepositoryRoot)
			if rootError != nil {
				return fmt.Errorf("open rooted repository for cleanup: %w", rootError)
			}
			defer repositoryHandle.Close()
			for _, directory := range createdDirectories {
				if err := repositoryHandle.Mkdir(filepath.FromSlash(directory), 0o755); err != nil {
					return err
				}
				if err := recorder.RecordCreatedDirectory(directory); err != nil {
					return err
				}
			}
			for _, group := range eligibleGroups {
				for _, operation := range group.Operations {
					if err := applyOperation(plan.RepositoryRoot, recorder, operation); err != nil {
						return fmt.Errorf("%s: %w", group.Code, err)
					}
				}
			}
			return nil
		})
	}
	result.Rollback = transactionResult.Rollback
	if transactionResult.Failure != nil {
		result.Outcome = transactionResult.Outcome
		changeState := "not applied"
		if transactionResult.Outcome == resultmodel.OutcomeRolledBack {
			changeState = "rolled back"
		} else if transactionResult.Outcome == resultmodel.OutcomeRisk && transactionResult.CommitSHA != "" {
			changeState = "committed state at risk"
		} else if transactionResult.Outcome == resultmodel.OutcomeRisk {
			changeState = "rollback incomplete"
		}
		result.Changes = relabelChanges(proposedChanges, changeState)
		markGroupChangeState(result.Changes, scratchGroups, "not applied")
		nextArgv := []string{"do-work-cli", "cleanup", "--dry-run"}
		if len(transactionResult.RevertArgv) > 0 {
			nextArgv = append([]string(nil), transactionResult.RevertArgv...)
		}
		result.Findings = append(result.Findings, resultmodel.CommandFinding{Code: "CLEANUP-TRANSACTION-" + strings.ToUpper(string(transactionResult.Failure.Kind)),
			Severity: resultmodel.SeverityError, AffectedPaths: transactionResult.Failure.Paths, Evidence: []string{transactionResult.Failure.Reason, "commit=" + transactionResult.CommitSHA},
			Fixability: resultmodel.FixabilityManual, AutomationStopReason: "the guarded cleanup transaction did not complete",
			NextArgv: nextArgv, VerificationArgv: []string{"git", "status", "--short"}})
		return result
	}
	if options.DryRun {
		result.Changes = proposedChanges
	} else {
		result.Changes = relabelChanges(proposedChanges, "applied")
	}
	if len(scratchGroups) > 0 {
		markConsumedScratchChanges(result.Changes, scratchGroups, options.DryRun)
		if !options.DryRun {
			if scratchError := applyConsumedScratch(plan.RepositoryRoot, scratchGroups); scratchError != nil {
				result.Outcome = resultmodel.OutcomeRisk
				markGroupChangeState(result.Changes, scratchGroups, "non-rollback outcome requires verification")
				nextArgv := []string{"do-work-cli", "cleanup", "--dry-run"}
				if len(transactionResult.RevertArgv) > 0 {
					nextArgv = transactionResult.RevertArgv
				}
				result.Findings = append(result.Findings, resultmodel.CommandFinding{Code: "CONSUMED-SCRATCH-DELETE-RISK", Severity: resultmodel.SeverityError,
					Evidence: []string{scratchError.Error()}, Fixability: resultmodel.FixabilityManual, AutomationStopReason: "spent scratch deletion stopped after a non-rollback mutation",
					NextArgv: nextArgv, VerificationArgv: []string{"do-work-cli", "cleanup", "--dry-run"}})
				return result
			}
		}
	}
	if !options.DryRun {
		removeEmptySourceDirectories(plan.RepositoryRoot, eligibleGroups)
	}
	if len(result.Findings) > 0 {
		result.Outcome = resultmodel.OutcomeFindings
	}
	return result
}

type groupDependencyState uint8

const (
	groupDependencyUnresolved groupDependencyState = iota
	groupDependencyResolving
	groupDependencyEligible
	groupDependencyBlocked
)

type groupDependencyBlocker struct {
	Code   string
	Reason string
}

func resolveGroupDependencies(groupIndex int, groups []OperationGroup, directlyEligible []bool, groupCodeCounts map[string]int, groupIndexByCode map[string]int, states []groupDependencyState, blockers []groupDependencyBlocker) (bool, groupDependencyBlocker) {
	switch states[groupIndex] {
	case groupDependencyEligible:
		return true, groupDependencyBlocker{}
	case groupDependencyBlocked:
		return false, blockers[groupIndex]
	case groupDependencyResolving:
		return false, groupDependencyBlocker{Code: groups[groupIndex].Code, Reason: "prerequisite cycle includes " + groups[groupIndex].Code}
	}
	if !directlyEligible[groupIndex] {
		return false, blockers[groupIndex]
	}

	states[groupIndex] = groupDependencyResolving
	seenPrerequisites := map[string]bool{}
	for _, requiredCode := range groups[groupIndex].RequiredGroupCodes {
		if seenPrerequisites[requiredCode] {
			states[groupIndex] = groupDependencyBlocked
			blockers[groupIndex] = groupDependencyBlocker{Code: requiredCode, Reason: "required group code is duplicated"}
			return false, blockers[groupIndex]
		}
		seenPrerequisites[requiredCode] = true
		if requiredCode == "" || groupCodeCounts[requiredCode] == 0 {
			states[groupIndex] = groupDependencyBlocked
			blockers[groupIndex] = groupDependencyBlocker{Code: requiredCode, Reason: "required group is missing from the cleanup plan"}
			return false, blockers[groupIndex]
		}
		if groupCodeCounts[requiredCode] > 1 {
			states[groupIndex] = groupDependencyBlocked
			blockers[groupIndex] = groupDependencyBlocker{Code: requiredCode, Reason: "required group code is duplicated in the cleanup plan"}
			return false, blockers[groupIndex]
		}
		requiredIndex := groupIndexByCode[requiredCode]
		if !directlyEligible[requiredIndex] {
			states[groupIndex] = groupDependencyBlocked
			blockers[groupIndex] = groupDependencyBlocker{Code: requiredCode, Reason: "required group did not pass direct preflight"}
			return false, blockers[groupIndex]
		}
		eligible, blocker := resolveGroupDependencies(requiredIndex, groups, directlyEligible, groupCodeCounts, groupIndexByCode, states, blockers)
		if !eligible {
			states[groupIndex] = groupDependencyBlocked
			blockers[groupIndex] = blocker
			return false, blocker
		}
	}
	states[groupIndex] = groupDependencyEligible
	return true, groupDependencyBlocker{}
}

func markConsumedScratchChanges(changes []resultmodel.RecordedChange, groups []OperationGroup, dryRun bool) {
	state := "applied non-rollback spent-scratch deletion"
	if dryRun {
		state = "planned non-rollback spent-scratch deletion"
	}
	markGroupChangeState(changes, groups, state)
}

func markGroupChangeState(changes []resultmodel.RecordedChange, groups []OperationGroup, state string) {
	paths := map[string]bool{}
	for _, group := range groups {
		for _, operation := range group.Operations {
			paths[operationPath(operation)] = true
		}
	}
	for index := range changes {
		if paths[changes[index].Path] {
			changes[index].Detail = state
		}
	}
}

func applyOperation(repositoryRoot string, recorder *gittransaction.MutationRecorder, operation CleanupOperation) error {
	sourcePath := filepath.Join(repositoryRoot, filepath.FromSlash(operation.SourcePath))
	switch operation.Kind {
	case OperationMove:
		if len(operation.Contents) > 0 {
			if err := recorder.RecordTouched(operation.SourcePath); err != nil {
				return err
			}
			if err := atomicfile.ReplaceExisting(sourcePath, operation.Contents); err != nil {
				return err
			}
		}
		if err := recorder.RecordTouched(operation.SourcePath); err != nil {
			return err
		}
		if err := recorder.RecordCreated(operation.DestinationPath); err != nil {
			return err
		}
		return moveWithoutOverwrite(repositoryRoot, operation.SourcePath, operation.DestinationPath)
	case OperationReplace:
		if err := recorder.RecordTouched(operation.SourcePath); err != nil {
			return err
		}
		return atomicfile.ReplaceExisting(sourcePath, operation.Contents)
	case OperationRewriteLinks:
		contents, err := os.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("read documentation for link rewrite: %w", err)
		}
		updated, _ := rewriteMarkdownTargets(contents, operation.SourcePath, operation.LinkMoveTargets)
		if err := recorder.RecordTouched(operation.SourcePath); err != nil {
			return err
		}
		return atomicfile.ReplaceExisting(sourcePath, updated)
	case OperationDelete:
		if err := recorder.RecordTouched(operation.SourcePath); err != nil {
			return err
		}
		return os.Remove(sourcePath)
	default:
		return fmt.Errorf("unknown cleanup operation %q", operation.Kind)
	}
}

func moveWithoutOverwrite(repositoryRoot, sourceRelativePath, destinationRelativePath string) error {
	sourcePath := filepath.Join(repositoryRoot, filepath.FromSlash(sourceRelativePath))
	destinationPath := filepath.Join(repositoryRoot, filepath.FromSlash(destinationRelativePath))
	sourceInfo, statError := os.Lstat(sourcePath)
	if statError != nil || !sourceInfo.Mode().IsRegular() {
		return fmt.Errorf("move source %s is not a regular file: %v", sourcePath, statError)
	}
	sourceContents, readError := os.ReadFile(sourcePath)
	if readError != nil {
		return fmt.Errorf("read move source %s: %w", sourcePath, readError)
	}
	afterReadInfo, afterReadError := os.Lstat(sourcePath)
	if afterReadError != nil || !afterReadInfo.Mode().IsRegular() || !os.SameFile(sourceInfo, afterReadInfo) {
		return fmt.Errorf("move source %s changed during read", sourcePath)
	}
	repositoryHandle, rootError := os.OpenRoot(repositoryRoot)
	if rootError != nil {
		return fmt.Errorf("open rooted repository for move: %w", rootError)
	}
	defer repositoryHandle.Close()
	destinationDirectory := filepath.ToSlash(filepath.Dir(filepath.FromSlash(destinationRelativePath)))
	destinationDirectoryInfo, directoryStatError := repositoryHandle.Lstat(filepath.FromSlash(destinationDirectory))
	if directoryStatError != nil || !destinationDirectoryInfo.IsDir() || destinationDirectoryInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("move destination directory %s is not a contained real directory", destinationDirectory)
	}
	destinationHandle, directoryRootError := repositoryHandle.OpenRoot(filepath.FromSlash(destinationDirectory))
	if directoryRootError != nil {
		return fmt.Errorf("open rooted move destination %s: %w", destinationDirectory, directoryRootError)
	}
	defer destinationHandle.Close()
	openedDirectoryInfo, openedStatError := destinationHandle.Stat(".")
	if openedStatError != nil || !os.SameFile(destinationDirectoryInfo, openedDirectoryInfo) {
		return fmt.Errorf("move destination directory %s changed while opening", destinationDirectory)
	}
	destinationName := filepath.Base(destinationRelativePath)
	if createError := atomicfile.CreateExclusiveAt(destinationHandle, destinationName, sourceContents, sourceInfo.Mode()); createError != nil {
		return fmt.Errorf("publish move destination %s without overwrite: %w", destinationPath, createError)
	}
	currentDirectoryInfo, currentDirectoryError := repositoryHandle.Lstat(filepath.FromSlash(destinationDirectory))
	if currentDirectoryError != nil || !currentDirectoryInfo.IsDir() || !os.SameFile(destinationDirectoryInfo, currentDirectoryInfo) {
		_ = destinationHandle.Remove(destinationName)
		return fmt.Errorf("move destination directory %s changed before publication", destinationDirectory)
	}
	currentContents, currentReadError := os.ReadFile(sourcePath)
	currentInfo, currentStatError := os.Lstat(sourcePath)
	if currentReadError != nil || currentStatError != nil || !currentInfo.Mode().IsRegular() || !os.SameFile(sourceInfo, currentInfo) || sha256.Sum256(currentContents) != sha256.Sum256(sourceContents) {
		_ = destinationHandle.Remove(destinationName)
		return fmt.Errorf("move source %s changed before deletion", sourcePath)
	}
	if removeError := os.Remove(sourcePath); removeError != nil {
		_ = destinationHandle.Remove(destinationName)
		return fmt.Errorf("remove move source %s: %w", sourcePath, removeError)
	}
	return nil
}

func groupTargetPaths(group OperationGroup) []string {
	return groupsTargetPaths([]OperationGroup{group})
}

func groupsTargetPaths(groups []OperationGroup) []string {
	paths := map[string]bool{}
	for _, group := range groups {
		for _, operation := range group.Operations {
			if operation.SourcePath != "" {
				paths[operation.SourcePath] = true
			}
			if operation.DestinationPath != "" {
				paths[operation.DestinationPath] = true
			}
		}
	}
	result := make([]string, 0, len(paths))
	for path := range paths {
		result = append(result, path)
	}
	sort.Strings(result)
	return result
}

func existingDestination(repositoryRoot string, group OperationGroup) string {
	for _, operation := range group.Operations {
		if operation.DestinationPath == "" {
			continue
		}
		if _, err := os.Lstat(filepath.Join(repositoryRoot, filepath.FromSlash(operation.DestinationPath))); err == nil {
			return operation.DestinationPath
		}
	}
	return ""
}

func plannedCreatedDirectories(repositoryRoot string, groups []OperationGroup) []string {
	directories := map[string]bool{}
	for _, group := range groups {
		for _, operation := range group.Operations {
			if operation.DestinationPath == "" {
				continue
			}
			directory := filepath.Dir(filepath.FromSlash(operation.DestinationPath))
			for directory != "." {
				if _, err := os.Lstat(filepath.Join(repositoryRoot, directory)); err == nil {
					break
				}
				directories[filepath.ToSlash(directory)] = true
				directory = filepath.Dir(directory)
			}
		}
	}
	result := make([]string, 0, len(directories))
	for directory := range directories {
		result = append(result, directory)
	}
	sort.Slice(result, func(left, right int) bool {
		leftDepth, rightDepth := strings.Count(result[left], "/"), strings.Count(result[right], "/")
		if leftDepth != rightDepth {
			return leftDepth < rightDepth
		}
		return result[left] < result[right]
	})
	return result
}

func removeEmptySourceDirectories(repositoryRoot string, groups []OperationGroup) {
	directories := map[string]bool{}
	preserved := map[string]bool{}
	for _, relativeRoot := range []string{"do-work", "do-work/queue", "do-work/working", "do-work/user-requests", "do-work/archive", "do-work/runs"} {
		preserved[filepath.Join(repositoryRoot, filepath.FromSlash(relativeRoot))] = true
	}
	for _, group := range groups {
		for _, operation := range group.Operations {
			if operation.Kind != OperationMove && operation.Kind != OperationDelete {
				continue
			}
			directory := filepath.Dir(filepath.Join(repositoryRoot, filepath.FromSlash(operation.SourcePath)))
			for directory != repositoryRoot && strings.HasPrefix(directory, repositoryRoot+string(filepath.Separator)) {
				directories[directory] = true
				directory = filepath.Dir(directory)
			}
		}
	}
	ordered := make([]string, 0, len(directories))
	for directory := range directories {
		ordered = append(ordered, directory)
	}
	sort.Slice(ordered, func(left, right int) bool {
		return strings.Count(ordered[left], string(filepath.Separator)) > strings.Count(ordered[right], string(filepath.Separator))
	})
	for _, directory := range ordered {
		if preserved[directory] {
			continue
		}
		_ = os.Remove(directory)
	}
}

func relabelChanges(changes []resultmodel.RecordedChange, state string) []resultmodel.RecordedChange {
	result := make([]resultmodel.RecordedChange, len(changes))
	for index, change := range changes {
		change.Detail = state + strings.TrimPrefix(change.Detail, "planned")
		result[index] = change
	}
	return result
}

func isUntrackedConsumedScratch(ctx context.Context, repositoryRoot string, group OperationGroup) bool {
	if group.PassNumber != 4 || len(group.Operations) == 0 {
		return false
	}
	if _, err := cleanupGit(ctx, repositoryRoot, "rev-parse", "--is-inside-work-tree"); err != nil {
		return false
	}
	for _, operation := range group.Operations {
		if operation.Kind != OperationDelete {
			return false
		}
		if _, err := cleanupGit(ctx, repositoryRoot, "ls-files", "--error-unmatch", "--", operation.SourcePath); err == nil {
			return false
		}
	}
	return true
}

func applyConsumedScratch(repositoryRoot string, groups []OperationGroup) error {
	repositoryHandle, rootError := os.OpenRoot(repositoryRoot)
	if rootError != nil {
		return rootError
	}
	defer repositoryHandle.Close()
	for _, group := range groups {
		if len(group.Operations) == 0 {
			continue
		}
		runDirectory := filepath.Dir(group.Operations[0].SourcePath)
		manifestRelativePath := filepath.Join(runDirectory, "manifest.md")
		manifestBytes, readError := repositoryHandle.ReadFile(filepath.FromSlash(manifestRelativePath))
		if readError != nil || !manifestIsConsumed(manifestBytes) {
			return fmt.Errorf("%s no longer has an exact consumed manifest", runDirectory)
		}
		planned := map[string]bool{}
		for _, operation := range group.Operations {
			planned[operation.SourcePath] = true
		}
		actual := map[string]bool{}
		walkError := fs.WalkDir(repositoryHandle.FS(), filepath.ToSlash(runDirectory), func(path string, entry fs.DirEntry, entryError error) error {
			if entryError != nil {
				return entryError
			}
			if entry.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("symlink %s entered consumed scratch", path)
			}
			if entry.IsDir() {
				return nil
			}
			if !entry.Type().IsRegular() {
				return fmt.Errorf("special file %s entered consumed scratch", path)
			}
			actual[filepath.ToSlash(path)] = true
			return nil
		})
		if walkError != nil || len(actual) != len(planned) {
			return fmt.Errorf("%s changed after planning: %v", runDirectory, walkError)
		}
		for path := range actual {
			if !planned[path] {
				return fmt.Errorf("%s gained unplanned path %s", runDirectory, path)
			}
		}
		for _, operation := range group.Operations {
			if removeError := repositoryHandle.Remove(filepath.FromSlash(operation.SourcePath)); removeError != nil {
				return removeError
			}
		}
	}
	removeEmptySourceDirectories(repositoryRoot, groups)
	return nil
}

func manifestIsConsumed(contents []byte) bool {
	for _, line := range strings.Split(string(contents), "\n") {
		name, value, found := strings.Cut(strings.TrimSpace(line), ":")
		if found && strings.EqualFold(strings.TrimSpace(name), "status") {
			return strings.EqualFold(strings.TrimSpace(value), "consumed")
		}
	}
	return false
}

func changeForOperation(operation CleanupOperation, dryRun bool) resultmodel.RecordedChange {
	detailPrefix := "applied"
	if dryRun {
		detailPrefix = "planned"
	}
	switch operation.Kind {
	case OperationMove:
		return resultmodel.RecordedChange{Path: operation.DestinationPath, Kind: string(operation.Kind), Detail: detailPrefix + " move from " + operation.SourcePath}
	case OperationReplace:
		return resultmodel.RecordedChange{Path: operation.SourcePath, Kind: string(operation.Kind), Detail: detailPrefix + " atomic replacement"}
	case OperationRewriteLinks:
		return resultmodel.RecordedChange{Path: operation.SourcePath, Kind: string(operation.Kind), Detail: detailPrefix + " documentation link rewrite"}
	default:
		return resultmodel.RecordedChange{Path: operation.SourcePath, Kind: string(operation.Kind), Detail: detailPrefix + " consumed scratch deletion"}
	}
}

func refusedGroupFinding(group OperationGroup, path, reason string) resultmodel.CommandFinding {
	paths := []string{}
	if path != "" {
		paths = append(paths, path)
	}
	ids := []string{}
	if group.AffectedID != "" {
		ids = append(ids, group.AffectedID)
	}
	return resultmodel.CommandFinding{Code: "CLEANUP-GROUP-REFUSED", Severity: resultmodel.SeverityWarning, AffectedIDs: ids, AffectedPaths: paths,
		Evidence: []string{group.Code + ": " + reason}, Fixability: resultmodel.FixabilityRefused, AutomationStopReason: "this operation group did not pass exact-target guards",
		NextArgv: []string{"git", "status", "--short", "--", path}, VerificationArgv: []string{"do-work-cli", "cleanup", "--dry-run"}}
}

func dependencyRefusedGroupFinding(group OperationGroup, blocker groupDependencyBlocker) resultmodel.CommandFinding {
	ids := []string{}
	if group.AffectedID != "" {
		ids = append(ids, group.AffectedID)
	}
	blockingCode := blocker.Code
	if blockingCode == "" {
		blockingCode = "<empty>"
	}
	return resultmodel.CommandFinding{Code: "CLEANUP-GROUP-REFUSED", Severity: resultmodel.SeverityWarning, AffectedIDs: ids,
		Evidence: []string{group.Code + ": prerequisite " + blockingCode + " blocked cleanup: " + blocker.Reason}, Fixability: resultmodel.FixabilityRefused,
		AutomationStopReason: "this operation group has a prerequisite that is not eligible",
		NextArgv:             []string{"do-work-cli", "cleanup", "--dry-run"}, VerificationArgv: []string{"do-work-cli", "cleanup", "--dry-run"}}
}
