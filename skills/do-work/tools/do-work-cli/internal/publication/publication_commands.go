package publication

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/atomicfile"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/gittransaction"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type commandOptions struct {
	manifestPath string
	dryRun       bool
	commit       bool
	answerTime   time.Time
}

var beforePublicationMutation = func(int, PlannedMutation) error { return nil }
var afterPublicationMutation = func(int, PlannedMutation) error { return nil }
var afterPublicationCommit = func(PublicationPlan) error { return nil }

func Handlers() map[string]commandruntime.CommandHandler {
	handlers := map[string]commandruntime.CommandHandler{}
	for _, operation := range []OperationName{OperationCaptureFiles, OperationAnswer, OperationRelease, OperationDeferGate} {
		selectedOperation := operation
		handlers[string(operation)] = func(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
			return handlePublicationCommand(executionContext, selectedOperation, arguments)
		}
	}
	return handlers
}

func handlePublicationCommand(executionContext commandruntime.ExecutionContext, operation OperationName, arguments []string) resultmodel.CommandResult {
	options, parseError := parseCommandOptions(arguments)
	if parseError != nil {
		return commandFailureWithManifest(executionContext.RepositoryRoot, operation, options.manifestPath, options.answerTime, "PUBLICATION-USAGE", parseError.Error())
	}
	manifestFile, openError := openManifest(executionContext.RepositoryRoot, options.manifestPath)
	if openError != nil {
		return commandFailureWithManifest(executionContext.RepositoryRoot, operation, options.manifestPath, options.answerTime, "PUBLICATION-MANIFEST-OPEN", openError.Error())
	}
	defer manifestFile.Close()
	manifest, decodeError := DecodeManifest(manifestFile, operation)
	if decodeError != nil {
		return commandFailureWithManifest(executionContext.RepositoryRoot, operation, options.manifestPath, options.answerTime, "PUBLICATION-MANIFEST-INVALID", decodeError.Error())
	}
	if options.commit && strings.TrimSpace(manifest.CommitMessage) == "" {
		return commandFailureWithManifest(executionContext.RepositoryRoot, operation, options.manifestPath, options.answerTime, "PUBLICATION-COMMIT-MESSAGE-MISSING", "--commit requires manifest commit_message")
	}
	var plan PublicationPlan
	switch operation {
	case OperationCaptureFiles:
		plan = BuildCapturePlan(executionContext.RepositoryRoot, manifest)
	case OperationAnswer:
		plan = BuildAnswerPlan(executionContext.RepositoryRoot, manifest, options.answerTime)
	case OperationRelease:
		plan = BuildReleasePlan(executionContext.RepositoryRoot, manifest)
	case OperationDeferGate:
		plan = BuildDeferGatePlan(executionContext.RepositoryRoot, manifest)
	}
	plan.ManifestPath = options.manifestPath
	if !options.answerTime.IsZero() {
		plan.AnswerAt = options.answerTime.Format(time.RFC3339)
	}
	return ApplyPlan(context.Background(), plan, options.dryRun, options.commit)
}

func parseCommandOptions(arguments []string) (commandOptions, error) {
	options := commandOptions{}
	for index := 0; index < len(arguments); index++ {
		switch arguments[index] {
		case "--dry-run":
			options.dryRun = true
		case "--commit":
			options.commit = true
		case "--manifest", "--at":
			option := arguments[index]
			index++
			if index >= len(arguments) {
				return options, fmt.Errorf("%s requires a value", option)
			}
			if option == "--manifest" {
				if options.manifestPath != "" {
					return options, fmt.Errorf("--manifest may be supplied once")
				}
				options.manifestPath = arguments[index]
			} else {
				parsed, err := time.Parse(time.RFC3339, arguments[index])
				if err != nil {
					return options, fmt.Errorf("--at requires RFC3339: %w", err)
				}
				options.answerTime = parsed
			}
		default:
			return options, fmt.Errorf("unknown publication option %q", arguments[index])
		}
	}
	if options.manifestPath == "" {
		return options, fmt.Errorf("--manifest requires a JSON file")
	}
	if options.dryRun && options.commit {
		return options, fmt.Errorf("--dry-run and --commit cannot be combined")
	}
	return options, nil
}

func openManifest(repositoryRoot, manifestPath string) (*os.File, error) {
	path := manifestPath
	if !filepath.IsAbs(path) {
		path = filepath.Join(repositoryRoot, filepath.FromSlash(path))
	}
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("manifest must be a regular non-symlink file")
	}
	return os.Open(path)
}

func ApplyPlan(ctx context.Context, plan PublicationPlan, dryRun, commit bool) resultmodel.CommandResult {
	if plan.Refusal != nil {
		return refusalResult(plan)
	}
	if !plan.Runnable() {
		return commandFailureWithManifest(plan.RepositoryRoot, plan.Operation, plan.ManifestPath, parseAnswerAt(plan.AnswerAt), "PUBLICATION-PLAN-INVALID", "publication plan has no mutations")
	}
	transaction := gittransaction.ExecuteTransaction(ctx, gittransaction.TransactionOptions{RepositoryRoot: plan.RepositoryRoot, TargetPaths: plan.TargetPaths, ExistingUntrackedTargetPaths: plan.ExistingUntrackedTargetPaths, ExistingDirtyTargetPaths: plan.ExistingDirtyTargetPaths, CreatedDirectoryPaths: plan.CreatedDirectoryPaths, DryRun: dryRun, Commit: commit, CommitMessage: plan.CommitMessage, PostCommitVerify: func(context.Context, string) error {
		if err := afterPublicationCommit(plan); err != nil {
			return err
		}
		return verifyPublicationPlan(plan)
	}}, func(recorder *gittransaction.MutationRecorder) error {
		repositoryHandle, openError := os.OpenRoot(plan.RepositoryRoot)
		if openError != nil {
			return openError
		}
		defer repositoryHandle.Close()
		for _, directory := range plan.CreatedDirectoryPaths {
			if recordError := recorder.RecordCreatedDirectory(directory); recordError != nil {
				return recordError
			}
			if makeError := repositoryHandle.Mkdir(filepath.FromSlash(directory), 0o755); makeError != nil {
				return makeError
			}
			if recordError := recorder.RecordCreatedDirectory(directory); recordError != nil {
				return recordError
			}
		}
		movedSourceDirectories := map[string]bool{}
		for mutationIndex, mutation := range plan.Mutations {
			if injectedError := beforePublicationMutation(mutationIndex, mutation); injectedError != nil {
				return injectedError
			}
			switch mutation.Kind {
			case MutationCreate:
				if createError := createRootedFile(plan.RepositoryRoot, mutation.Path, mutation.Contents, mutation.Mode); createError != nil {
					return createError
				}
				if recordError := recorder.RecordCreated(mutation.Path); recordError != nil {
					return recordError
				}
			case MutationReplace:
				if replaceError := replaceRootedFile(plan.RepositoryRoot, mutation.Path, mutation.ExpectedBytes, mutation.Contents); replaceError != nil {
					return replaceError
				}
				if recordError := recorder.RecordTouched(mutation.Path); recordError != nil {
					return recordError
				}
			case MutationMove:
				if len(mutation.Contents) > 0 && !bytes.Equal(mutation.Contents, mutation.ExpectedBytes) {
					if replaceError := replaceRootedFile(plan.RepositoryRoot, mutation.Path, mutation.ExpectedBytes, mutation.Contents); replaceError != nil {
						return replaceError
					}
					if recordError := recorder.RecordTouched(mutation.Path); recordError != nil {
						return recordError
					}
				}
				expectedMoveBytes := mutation.ExpectedBytes
				if len(mutation.Contents) > 0 {
					expectedMoveBytes = mutation.Contents
				}
				if moveError := moveRootedFile(plan.RepositoryRoot, mutation.Path, mutation.DestinationPath, expectedMoveBytes); moveError != nil {
					return moveError
				}
				if recordError := recorder.RecordTouched(mutation.Path); recordError != nil {
					return recordError
				}
				if recordError := recorder.RecordCreated(mutation.DestinationPath); recordError != nil {
					return recordError
				}
				movedSourceDirectories[filepath.ToSlash(filepath.Dir(filepath.FromSlash(mutation.Path)))] = true
			default:
				return fmt.Errorf("unknown publication mutation %q", mutation.Kind)
			}
			if injectedError := afterPublicationMutation(mutationIndex, mutation); injectedError != nil {
				return injectedError
			}
		}
		if verifyError := verifyPublicationPlan(plan); verifyError != nil {
			return verifyError
		}
		removeMovedDirectories(plan.RepositoryRoot, movedSourceDirectories)
		return nil
	})
	result := gittransaction.BuildCommandResult(string(plan.Operation), transaction)
	if transaction.Failure == nil {
		result.Changes = append([]resultmodel.RecordedChange(nil), plan.Changes...)
		detail := "applied"
		code, evidence := "PUBLICATION-APPLIED", "publication transaction applied"
		if dryRun {
			detail, code, evidence = "planned", "PUBLICATION-DRY-RUN", "publication transaction planned without changing bytes"
		}
		if transaction.CommitSHA != "" {
			detail = "committed in " + transaction.CommitSHA
		}
		for index := range result.Changes {
			result.Changes[index].Detail = detail
		}
		result.Findings = append(result.Findings, successFinding(plan, code, evidence))
		if plan.GateDeferral != nil {
			gateDeferral := *plan.GateDeferral
			gateDeferral.GateCommand = append([]string(nil), gateDeferral.GateCommand...)
			result.GateDeferral = &gateDeferral
		}
	}
	return result
}

func createRootedFile(repositoryRoot, relativePath string, contents []byte, mode os.FileMode) error {
	repositoryHandle, rootError := os.OpenRoot(repositoryRoot)
	if rootError != nil {
		return rootError
	}
	defer repositoryHandle.Close()
	parentPath := filepath.Dir(filepath.FromSlash(relativePath))
	parentRoot, beforeInfo, openError := openRootedParent(repositoryHandle, parentPath)
	if openError != nil {
		return openError
	}
	defer parentRoot.Close()
	if createError := atomicfile.CreateExclusiveAt(parentRoot, filepath.Base(relativePath), contents, mode); createError != nil {
		return createError
	}
	afterInfo, afterError := repositoryHandle.Lstat(parentPath)
	if afterError != nil || !afterInfo.IsDir() || !os.SameFile(beforeInfo, afterInfo) {
		_ = parentRoot.Remove(filepath.Base(relativePath))
		return fmt.Errorf("destination parent identity changed: %s", parentPath)
	}
	return nil
}

func moveRootedFile(repositoryRoot, sourcePath, destinationPath string, expectedBytes []byte) error {
	repositoryHandle, rootError := os.OpenRoot(repositoryRoot)
	if rootError != nil {
		return rootError
	}
	defer repositoryHandle.Close()
	sourceParentPath := filepath.Dir(filepath.FromSlash(sourcePath))
	sourceParent, sourceParentInfo, sourceParentError := openRootedParent(repositoryHandle, sourceParentPath)
	if sourceParentError != nil {
		return sourceParentError
	}
	defer sourceParent.Close()
	destinationParentPath := filepath.Dir(filepath.FromSlash(destinationPath))
	destinationParent, destinationParentInfo, destinationParentError := openRootedParent(repositoryHandle, destinationParentPath)
	if destinationParentError != nil {
		return destinationParentError
	}
	defer destinationParent.Close()
	sourceName := filepath.Base(sourcePath)
	destinationName := filepath.Base(destinationPath)
	sourceFile, openError := sourceParent.Open(sourceName)
	if openError != nil {
		return openError
	}
	defer sourceFile.Close()
	sourceInfo, statError := sourceFile.Stat()
	if statError != nil || !sourceInfo.Mode().IsRegular() {
		return fmt.Errorf("move source is not a regular file: %s", sourcePath)
	}
	contents, readError := sourceParent.ReadFile(sourceName)
	if readError != nil || !bytes.Equal(contents, expectedBytes) {
		return fmt.Errorf("move source changed: %s", sourcePath)
	}
	if !rootedParentIdentity(repositoryHandle, sourceParentPath, sourceParentInfo) || !rootedParentIdentity(repositoryHandle, destinationParentPath, destinationParentInfo) {
		return fmt.Errorf("move parent identity changed")
	}
	if createError := atomicfile.CreateExclusiveAt(destinationParent, destinationName, contents, sourceInfo.Mode()); createError != nil {
		return createError
	}
	keepDestination := false
	defer func() {
		if !keepDestination {
			_ = destinationParent.Remove(destinationName)
		}
	}()
	if !rootedParentIdentity(repositoryHandle, destinationParentPath, destinationParentInfo) {
		return fmt.Errorf("destination parent identity changed: %s", destinationParentPath)
	}
	currentInfo, currentError := sourceParent.Lstat(sourceName)
	if currentError != nil || !os.SameFile(sourceInfo, currentInfo) {
		return fmt.Errorf("move source identity changed: %s", sourcePath)
	}
	if !rootedParentIdentity(repositoryHandle, sourceParentPath, sourceParentInfo) {
		return fmt.Errorf("source parent identity changed: %s", sourceParentPath)
	}
	if removeError := sourceParent.Remove(sourceName); removeError != nil {
		return removeError
	}
	keepDestination = true
	return nil
}

func replaceRootedFile(repositoryRoot, path string, expected, contents []byte) error {
	repositoryHandle, rootError := os.OpenRoot(repositoryRoot)
	if rootError != nil {
		return rootError
	}
	defer repositoryHandle.Close()
	parentPath := filepath.Dir(filepath.FromSlash(path))
	parentRoot, parentInfo, parentError := openRootedParent(repositoryHandle, parentPath)
	if parentError != nil {
		return parentError
	}
	defer parentRoot.Close()
	name := filepath.Base(path)
	targetFile, openError := parentRoot.Open(name)
	if openError != nil {
		return openError
	}
	targetInfo, statError := targetFile.Stat()
	targetFile.Close()
	if statError != nil || !targetInfo.Mode().IsRegular() {
		return fmt.Errorf("target is not a regular non-symlink file: %s", path)
	}
	current, readError := parentRoot.ReadFile(name)
	if readError != nil || !bytes.Equal(current, expected) {
		return fmt.Errorf("target preimage changed: %s", path)
	}
	temporaryName, nameError := rootedTemporaryName(name)
	if nameError != nil {
		return nameError
	}
	if createError := atomicfile.CreateExclusiveAt(parentRoot, temporaryName, contents, targetInfo.Mode()); createError != nil {
		return createError
	}
	defer parentRoot.Remove(temporaryName)
	if !rootedParentIdentity(repositoryHandle, parentPath, parentInfo) {
		return fmt.Errorf("target parent identity changed: %s", parentPath)
	}
	currentInfo, currentInfoError := parentRoot.Lstat(name)
	current, readError = parentRoot.ReadFile(name)
	if currentInfoError != nil || !os.SameFile(targetInfo, currentInfo) || readError != nil || !bytes.Equal(current, expected) {
		return fmt.Errorf("target changed before publication: %s", path)
	}
	if renameError := parentRoot.Rename(temporaryName, name); renameError != nil {
		return renameError
	}
	if !rootedParentIdentity(repositoryHandle, parentPath, parentInfo) {
		return fmt.Errorf("target parent identity changed after publication: %s", parentPath)
	}
	return nil
}

func openRootedParent(repositoryHandle *os.Root, parentPath string) (*os.Root, os.FileInfo, error) {
	info, statError := repositoryHandle.Lstat(parentPath)
	if statError != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, nil, fmt.Errorf("destination parent is not a real repository directory: %s", parentPath)
	}
	parentRoot, openError := repositoryHandle.OpenRoot(parentPath)
	if openError != nil {
		return nil, nil, fmt.Errorf("opening rooted destination parent %s: %w", parentPath, openError)
	}
	return parentRoot, info, nil
}

func rootedParentIdentity(repositoryHandle *os.Root, parentPath string, expected os.FileInfo) bool {
	current, err := repositoryHandle.Lstat(parentPath)
	return err == nil && current.IsDir() && current.Mode()&os.ModeSymlink == 0 && os.SameFile(expected, current)
}

func rootedTemporaryName(base string) (string, error) {
	var random [12]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", err
	}
	return "." + base + ".publication-" + hex.EncodeToString(random[:]), nil
}
func absolutePath(root, path string) string { return filepath.Join(root, filepath.FromSlash(path)) }

func verifyPublicationPlan(plan PublicationPlan) error {
	for _, mutation := range plan.Mutations {
		path, expected := mutation.Path, mutation.Contents
		if mutation.Kind == MutationMove {
			path = mutation.DestinationPath
			if len(expected) == 0 {
				expected = mutation.ExpectedBytes
			}
			if _, err := os.Lstat(absolutePath(plan.RepositoryRoot, mutation.Path)); !os.IsNotExist(err) {
				return fmt.Errorf("move source still exists: %s", mutation.Path)
			}
		}
		current, readError := os.ReadFile(absolutePath(plan.RepositoryRoot, path))
		if readError != nil || !bytes.Equal(current, expected) {
			return fmt.Errorf("publication verification failed: %s", path)
		}
	}
	return nil
}

func removeMovedDirectories(repositoryRoot string, directories map[string]bool) {
	paths := make([]string, 0, len(directories))
	for path := range directories {
		paths = append(paths, path)
	}
	sort.Slice(paths, func(i, j int) bool { return strings.Count(paths[i], "/") > strings.Count(paths[j], "/") })
	for _, path := range paths {
		for current := path; strings.HasPrefix(current, "do-work/user-requests/"); current = filepath.ToSlash(filepath.Dir(filepath.FromSlash(current))) {
			if os.Remove(absolutePath(repositoryRoot, current)) != nil {
				break
			}
		}
	}
}

func refusalResult(plan PublicationPlan) resultmodel.CommandResult {
	refusal := plan.Refusal
	nextArgv, verificationArgv, recipe := publicationProtocol(plan.Operation, plan.ManifestPath, plan.AnswerAt)
	return resultmodel.CommandResult{Command: string(plan.Operation), Outcome: resultmodel.OutcomeRefused, RepositoryRoot: plan.RepositoryRoot,
		Findings: []resultmodel.CommandFinding{{Code: refusal.Code, Severity: resultmodel.SeverityWarning, AffectedIDs: refusal.IDs, AffectedPaths: refusal.Paths,
			Evidence: []string{refusal.Reason}, Fixability: resultmodel.FixabilityRefused, AutomationStopReason: "the publication precondition did not hold",
			NextArgv: nextArgv, NextJustRecipe: recipe, VerificationArgv: verificationArgv}}}
}

func commandFailure(repositoryRoot string, operation OperationName, code, reason string) resultmodel.CommandResult {
	return commandFailureWithManifest(repositoryRoot, operation, "publication-manifest.json", time.Time{}, code, reason)
}

func commandFailureWithManifest(repositoryRoot string, operation OperationName, manifestPath string, answerTime time.Time, code, reason string) resultmodel.CommandResult {
	answerAt := ""
	if !answerTime.IsZero() {
		answerAt = answerTime.Format(time.RFC3339)
	}
	nextArgv, verificationArgv, recipe := publicationProtocol(operation, manifestPath, answerAt)
	return resultmodel.CommandResult{Command: string(operation), Outcome: resultmodel.OutcomeFailure, RepositoryRoot: repositoryRoot,
		Findings: []resultmodel.CommandFinding{{Code: code, Severity: resultmodel.SeverityError, Evidence: []string{reason}, Fixability: resultmodel.FixabilityManual,
			AutomationStopReason: "the publication command could not start safely", NextArgv: nextArgv, NextJustRecipe: recipe, VerificationArgv: verificationArgv}}}
}

func successFinding(plan PublicationPlan, code, evidence string) resultmodel.CommandFinding {
	nextArgv, verification, recipe := publicationProtocol(plan.Operation, plan.ManifestPath, plan.AnswerAt)
	affectedIDs := []string(nil)
	if plan.GateDeferral != nil {
		affectedIDs = []string{plan.GateDeferral.ParentID, plan.GateDeferral.RepairID}
		evidence = fmt.Sprintf("%s; repair %s; fingerprint %s", evidence, plan.GateDeferral.RepairOutcome, plan.GateDeferral.DiagnosticFingerprint)
	}
	return resultmodel.CommandFinding{Code: code, Severity: resultmodel.SeverityInfo, AffectedIDs: affectedIDs, AffectedPaths: plan.TargetPaths, Evidence: []string{evidence},
		Fixability: resultmodel.FixabilityAutomatic, NextArgv: nextArgv, NextJustRecipe: recipe, VerificationArgv: verification}
}

func publicationProtocol(operation OperationName, manifestPath, answerAt string) ([]string, []string, string) {
	if strings.TrimSpace(manifestPath) == "" {
		manifestPath = "publication-manifest.json"
	}
	next := []string{"do-work-cli", string(operation), "--manifest", manifestPath}
	if operation == OperationAnswer && answerAt != "" {
		next = append(next, "--at", answerAt)
	}
	verification := append([]string{"do-work-cli", "--format", "json"}, next[1:]...)
	verification = append(verification, "--dry-run")
	recipe := "do-work-" + string(operation) + " --manifest " + quotePublicationRecipeArgument(manifestPath)
	if operation == OperationAnswer && answerAt != "" {
		recipe += " --at " + quotePublicationRecipeArgument(answerAt)
	}
	return next, verification, recipe
}

func quotePublicationRecipeArgument(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func parseAnswerAt(value string) time.Time {
	parsed, _ := time.Parse(time.RFC3339, value)
	return parsed
}
