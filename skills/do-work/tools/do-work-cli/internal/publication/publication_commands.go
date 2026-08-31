package publication

import (
	"bytes"
	"context"
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

func Handlers() map[string]commandruntime.CommandHandler {
	handlers := map[string]commandruntime.CommandHandler{}
	for _, operation := range []OperationName{OperationCaptureFiles, OperationAnswer, OperationRelease} {
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
		return commandFailure(executionContext.RepositoryRoot, operation, "PUBLICATION-USAGE", parseError.Error())
	}
	manifestFile, openError := openManifest(executionContext.RepositoryRoot, options.manifestPath)
	if openError != nil {
		return commandFailure(executionContext.RepositoryRoot, operation, "PUBLICATION-MANIFEST-OPEN", openError.Error())
	}
	defer manifestFile.Close()
	manifest, decodeError := DecodeManifest(manifestFile, operation)
	if decodeError != nil {
		return commandFailure(executionContext.RepositoryRoot, operation, "PUBLICATION-MANIFEST-INVALID", decodeError.Error())
	}
	if options.commit && strings.TrimSpace(manifest.CommitMessage) == "" {
		return commandFailure(executionContext.RepositoryRoot, operation, "PUBLICATION-COMMIT-MESSAGE-MISSING", "--commit requires manifest commit_message")
	}
	var plan PublicationPlan
	switch operation {
	case OperationCaptureFiles:
		plan = BuildCapturePlan(executionContext.RepositoryRoot, manifest)
	case OperationAnswer:
		plan = BuildAnswerPlan(executionContext.RepositoryRoot, manifest, options.answerTime)
	case OperationRelease:
		plan = BuildReleasePlan(executionContext.RepositoryRoot, manifest)
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
		return commandFailure(plan.RepositoryRoot, plan.Operation, "PUBLICATION-PLAN-INVALID", "publication plan has no mutations")
	}
	transaction := gittransaction.ExecuteTransaction(ctx, gittransaction.TransactionOptions{RepositoryRoot: plan.RepositoryRoot, TargetPaths: plan.TargetPaths, ExistingUntrackedTargetPaths: plan.ExistingUntrackedTargetPaths, CreatedDirectoryPaths: plan.CreatedDirectoryPaths, DryRun: dryRun, Commit: commit, CommitMessage: plan.CommitMessage, PostCommitVerify: func(context.Context, string) error { return verifyPublicationPlan(plan) }}, func(recorder *gittransaction.MutationRecorder) error {
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
		}
		movedSourceDirectories := map[string]bool{}
		for mutationIndex, mutation := range plan.Mutations {
			if injectedError := beforePublicationMutation(mutationIndex, mutation); injectedError != nil {
				return injectedError
			}
			switch mutation.Kind {
			case MutationCreate:
				if recordError := recorder.RecordCreated(mutation.Path); recordError != nil {
					return recordError
				}
				if createError := createRootedFile(plan.RepositoryRoot, mutation.Path, mutation.Contents, mutation.Mode); createError != nil {
					return createError
				}
			case MutationReplace:
				if revalidateError := revalidateFile(plan.RepositoryRoot, mutation.Path, mutation.ExpectedBytes); revalidateError != nil {
					return revalidateError
				}
				if recordError := recorder.RecordTouched(mutation.Path); recordError != nil {
					return recordError
				}
				if replaceError := atomicfile.ReplaceExisting(absolutePath(plan.RepositoryRoot, mutation.Path), mutation.Contents); replaceError != nil {
					return replaceError
				}
			case MutationMove:
				if revalidateError := revalidateFile(plan.RepositoryRoot, mutation.Path, mutation.ExpectedBytes); revalidateError != nil {
					return revalidateError
				}
				if recordError := recorder.RecordTouched(mutation.Path); recordError != nil {
					return recordError
				}
				if len(mutation.Contents) > 0 && !bytes.Equal(mutation.Contents, mutation.ExpectedBytes) {
					if replaceError := atomicfile.ReplaceExisting(absolutePath(plan.RepositoryRoot, mutation.Path), mutation.Contents); replaceError != nil {
						return replaceError
					}
				}
				if recordError := recorder.RecordCreated(mutation.DestinationPath); recordError != nil {
					return recordError
				}
				expectedMoveBytes := mutation.ExpectedBytes
				if len(mutation.Contents) > 0 {
					expectedMoveBytes = mutation.Contents
				}
				if moveError := moveRootedFile(plan.RepositoryRoot, mutation.Path, mutation.DestinationPath, expectedMoveBytes); moveError != nil {
					return moveError
				}
				movedSourceDirectories[filepath.ToSlash(filepath.Dir(filepath.FromSlash(mutation.Path)))] = true
			default:
				return fmt.Errorf("unknown publication mutation %q", mutation.Kind)
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
	}
	return result
}

func createRootedFile(repositoryRoot, relativePath string, contents []byte, mode os.FileMode) error {
	parentPath := filepath.Dir(filepath.FromSlash(relativePath))
	absoluteParent := filepath.Join(repositoryRoot, parentPath)
	beforeInfo, statError := os.Stat(absoluteParent)
	if statError != nil || !beforeInfo.IsDir() {
		return fmt.Errorf("destination parent is not a real directory: %s", parentPath)
	}
	parentRoot, openError := os.OpenRoot(absoluteParent)
	if openError != nil {
		return openError
	}
	defer parentRoot.Close()
	if createError := atomicfile.CreateExclusiveAt(parentRoot, filepath.Base(relativePath), contents, mode); createError != nil {
		return createError
	}
	afterInfo, afterError := os.Stat(absoluteParent)
	if afterError != nil || !os.SameFile(beforeInfo, afterInfo) {
		return fmt.Errorf("destination parent identity changed: %s", parentPath)
	}
	return nil
}

func moveRootedFile(repositoryRoot, sourcePath, destinationPath string, expectedBytes []byte) error {
	sourceAbsolute := absolutePath(repositoryRoot, sourcePath)
	sourceFile, openError := os.Open(sourceAbsolute)
	if openError != nil {
		return openError
	}
	defer sourceFile.Close()
	sourceInfo, statError := sourceFile.Stat()
	if statError != nil || !sourceInfo.Mode().IsRegular() {
		return fmt.Errorf("move source is not a regular file: %s", sourcePath)
	}
	contents, readError := os.ReadFile(sourceAbsolute)
	if readError != nil || !bytes.Equal(contents, expectedBytes) {
		return fmt.Errorf("move source changed: %s", sourcePath)
	}
	if createError := createRootedFile(repositoryRoot, destinationPath, contents, sourceInfo.Mode()); createError != nil {
		return createError
	}
	currentInfo, currentError := os.Lstat(sourceAbsolute)
	if currentError != nil || !os.SameFile(sourceInfo, currentInfo) {
		return fmt.Errorf("move source identity changed: %s", sourcePath)
	}
	return os.Remove(sourceAbsolute)
}

func revalidateFile(repositoryRoot, path string, expected []byte) error {
	current, err := os.ReadFile(absolutePath(repositoryRoot, path))
	if err != nil || !bytes.Equal(current, expected) {
		return fmt.Errorf("target preimage changed: %s", path)
	}
	info, err := os.Lstat(absolutePath(repositoryRoot, path))
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("target is not a regular non-symlink file: %s", path)
	}
	return nil
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
	return resultmodel.CommandResult{Command: string(plan.Operation), Outcome: resultmodel.OutcomeRefused, RepositoryRoot: plan.RepositoryRoot,
		Findings: []resultmodel.CommandFinding{{Code: refusal.Code, Severity: resultmodel.SeverityWarning, AffectedIDs: refusal.IDs, AffectedPaths: refusal.Paths,
			Evidence: []string{refusal.Reason}, Fixability: resultmodel.FixabilityRefused, AutomationStopReason: "the publication precondition did not hold",
			NextArgv: []string{"do-work-cli", string(plan.Operation), "--manifest", "<manifest.json>"}, VerificationArgv: []string{"do-work-cli", "--format", "json", string(plan.Operation), "--manifest", "<manifest.json>", "--dry-run"}}}}
}

func commandFailure(repositoryRoot string, operation OperationName, code, reason string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Command: string(operation), Outcome: resultmodel.OutcomeFailure, RepositoryRoot: repositoryRoot,
		Findings: []resultmodel.CommandFinding{{Code: code, Severity: resultmodel.SeverityError, Evidence: []string{reason}, Fixability: resultmodel.FixabilityManual,
			AutomationStopReason: "the publication command could not start safely", NextArgv: []string{"do-work-cli", string(operation), "--manifest", "<manifest.json>"},
			VerificationArgv: []string{"do-work-cli", "--format", "json", string(operation), "--manifest", "<manifest.json>", "--dry-run"}}}}
}

func successFinding(plan PublicationPlan, code, evidence string) resultmodel.CommandFinding {
	verification := append([]string{"git", "status", "--short", "--"}, plan.TargetPaths...)
	return resultmodel.CommandFinding{Code: code, Severity: resultmodel.SeverityInfo, AffectedPaths: plan.TargetPaths, Evidence: []string{evidence},
		Fixability: resultmodel.FixabilityAutomatic, NextArgv: []string{"do-work-cli", string(plan.Operation), "--manifest", "<manifest.json>"}, VerificationArgv: verification}
}
