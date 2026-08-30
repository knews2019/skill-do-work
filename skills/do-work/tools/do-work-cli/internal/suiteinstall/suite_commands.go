package suiteinstall

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/archivefetch"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/managedsection"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/suitemanifest"
)

// The five command names this package registers. They are the argv tokens, so they stay
// exactly as a user types them.
const (
	CommandInstallSuite     = "install-suite"
	CommandUpdateSuite      = "update-suite"
	CommandReplaceSection   = "replace-section"
	CommandValidateManifest = "validate-manifest"
	CommandFetchArchive     = "fetch-archive"
)

// Handlers returns the command table main.go registers. Narration and confirmation are
// injected rather than reached for, so a test drives the same handlers the binary does.
//
// stdout carries only the rendered CommandResult; every progress line, review diff and
// confirmation prompt goes to the narration writer. That split is the contract the later
// commands in this family inherit.
func Handlers(narration io.Writer, confirmationInput io.Reader) map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{
		CommandInstallSuite: func(context commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
			return handleInstallSuite(context, arguments, narration, confirmationInput)
		},
		CommandUpdateSuite: func(context commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
			return handleUpdateSuite(context, arguments, narration, confirmationInput)
		},
		CommandReplaceSection: func(context commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
			return handleReplaceSection(context, arguments)
		},
		CommandValidateManifest: func(context commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
			return handleValidateManifest(context, arguments)
		},
		CommandFetchArchive: func(context commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
			return handleFetchArchive(context, arguments, narration)
		},
	}
}

// optionValues parses long-form options only. Command flags carry no single-letter aliases
// and take no positional arguments: a legacy positional argv stays in the compatibility
// launcher, where it belongs.
type optionValues struct {
	values map[string]string
	flags  map[string]bool
}

func parseCommandOptions(arguments []string, valueOptions, booleanOptions []string) (optionValues, error) {
	parsed := optionValues{values: map[string]string{}, flags: map[string]bool{}}
	valueSet := map[string]struct{}{}
	for _, option := range valueOptions {
		valueSet[option] = struct{}{}
	}
	booleanSet := map[string]struct{}{}
	for _, option := range booleanOptions {
		booleanSet[option] = struct{}{}
	}
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		name, inlineValue, hasInlineValue := strings.Cut(argument, "=")
		if _, isValueOption := valueSet[name]; isValueOption {
			if hasInlineValue {
				parsed.values[name] = inlineValue
				continue
			}
			index++
			if index >= len(arguments) {
				return parsed, fmt.Errorf("%s requires a value", name)
			}
			parsed.values[name] = arguments[index]
			continue
		}
		if _, isBooleanOption := booleanSet[argument]; isBooleanOption {
			parsed.flags[argument] = true
			continue
		}
		return parsed, fmt.Errorf("option %q is not available", argument)
	}
	return parsed, nil
}

func handleInstallSuite(context commandruntime.ExecutionContext, arguments []string, narration io.Writer, confirmationInput io.Reader) resultmodel.CommandResult {
	parsed, err := parseCommandOptions(arguments, []string{"--archive", "--upstream-url"}, nil)
	if err != nil {
		return usageResult(CommandInstallSuite, context.RepositoryRoot, err.Error())
	}
	result := RunInstall(commandContext(), InstallOptions{
		ProjectRoot:         context.RepositoryRoot,
		SuppliedArchivePath: parsed.values["--archive"],
		UpstreamURL:         parsed.values["--upstream-url"],
		ToolDirectory:       shippedToolsDirectory(),
		Narration:           narration,
		ConfirmationInput:   confirmationInput,
	})
	return installResultToCommandResult(CommandInstallSuite, context.RepositoryRoot, result)
}

func installResultToCommandResult(commandName, repositoryRoot string, result InstallResult) resultmodel.CommandResult {
	commandResult := resultmodel.CommandResult{
		Command:        commandName,
		Outcome:        result.Outcome,
		RepositoryRoot: repositoryRoot,
		Changes:        result.Changes,
		SkippedWork:    result.SkippedWork,
		Rollback:       result.Rollback,
	}
	if result.FailureReason != "" {
		commandResult.Findings = []resultmodel.CommandFinding{{
			Code:                 findingCodeForInstallOutcome(result.Outcome),
			Severity:             resultmodel.SeverityError,
			AffectedPaths:        result.FailurePaths,
			Evidence:             []string{result.FailureReason},
			Fixability:           resultmodel.FixabilityManual,
			AutomationStopReason: automationStopReasonFor(result.Outcome),
			NextArgv:             []string{"do-work-cli", "--format", "json", commandName},
			VerificationArgv:     []string{"git", "status", "--short"},
		}}
	}
	return commandResult
}

func findingCodeForInstallOutcome(outcome resultmodel.CommandOutcome) string {
	switch outcome {
	case resultmodel.OutcomeRolledBack:
		return "SUITE-INSTALL-ROLLED-BACK"
	case resultmodel.OutcomeRisk:
		return "SUITE-INSTALL-RECOVERY-INCOMPLETE"
	default:
		return "SUITE-INSTALL-FAILED"
	}
}

func automationStopReasonFor(outcome resultmodel.CommandOutcome) string {
	switch outcome {
	case resultmodel.OutcomeRolledBack:
		return "the install failed and every managed path plus the Git index was restored"
	case resultmodel.OutcomeRisk:
		return "automatic recovery was incomplete, so the project needs a person before any retry"
	default:
		return "the install refused before any client file was written"
	}
}

func handleUpdateSuite(context commandruntime.ExecutionContext, arguments []string, narration io.Writer, confirmationInput io.Reader) resultmodel.CommandResult {
	parsed, err := parseCommandOptions(arguments, []string{"--skill-root", "--upstream-url"}, nil)
	if err != nil {
		return usageResult(CommandUpdateSuite, context.RepositoryRoot, err.Error())
	}
	skillRoot := parsed.values["--skill-root"]
	if skillRoot == "" {
		skillRoot = filepath.Join(context.RepositoryRoot, ".claude", "skills", "do-work")
	}
	result := RunUpdate(commandContext(), UpdateOptions{
		ProjectRoot:        context.RepositoryRoot,
		InstalledSkillRoot: skillRoot,
		UpstreamURL:        parsed.values["--upstream-url"],
		ToolDirectory:      shippedToolsDirectory(),
		Narration:          narration,
		ConfirmationInput:  confirmationInput,
	})
	commandResult := resultmodel.CommandResult{
		Command:        CommandUpdateSuite,
		Outcome:        result.Outcome,
		RepositoryRoot: context.RepositoryRoot,
		Changes:        result.Changes,
		SkippedWork:    result.SkippedWork,
		Rollback:       result.Rollback,
	}
	if result.FailureReason != "" {
		commandResult.Findings = []resultmodel.CommandFinding{{
			Code:                 "SUITE-UPDATE-FAILED",
			Severity:             resultmodel.SeverityError,
			AffectedPaths:        result.FailurePaths,
			Evidence:             []string{result.FailureReason},
			Fixability:           resultmodel.FixabilityManual,
			AutomationStopReason: automationStopReasonFor(result.Outcome),
			NextArgv:             []string{"do-work-cli", "--format", "json", CommandUpdateSuite},
			VerificationArgv:     []string{"git", "status", "--short"},
		}}
	}
	return commandResult
}

// handleReplaceSection is repository-independent: it operates on the files it is given, so
// the global --repo-root only labels the result.
func handleReplaceSection(context commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	parsed, err := parseCommandOptions(arguments,
		[]string{"--target", "--section-file", "--template-file", "--begin-marker", "--end-marker"},
		[]string{"--reject-recipe-collisions"})
	if err != nil {
		return usageResult(CommandReplaceSection, context.RepositoryRoot, err.Error())
	}
	outcome, replaceErr := managedsection.ReplaceSection(managedsection.ReplaceRequest{
		TargetPath:             parsed.values["--target"],
		SectionFilePath:        parsed.values["--section-file"],
		TemplateFilePath:       parsed.values["--template-file"],
		RejectRecipeCollisions: parsed.flags["--reject-recipe-collisions"],
		BeginMarker:            parsed.values["--begin-marker"],
		EndMarker:              parsed.values["--end-marker"],
	})
	if replaceErr != nil {
		failure, isSectionFailure := replaceErr.(*managedsection.SectionFailure)
		code := managedsection.FailureInvalidInput
		commandOutcome := resultmodel.OutcomeFailure
		fixability := resultmodel.FixabilityManual
		stopReason := "the managed section could not be replaced, so the target was left untouched"
		if isSectionFailure {
			code = failure.Code
			if failure.Code == managedsection.FailureReservedRecipeCollision {
				commandOutcome = resultmodel.OutcomeRefused
				fixability = resultmodel.FixabilityRefused
				stopReason = "the target already owns a name the managed section reserves"
			}
		}
		return resultmodel.CommandResult{
			Command:        CommandReplaceSection,
			Outcome:        commandOutcome,
			RepositoryRoot: context.RepositoryRoot,
			Rollback:       resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
			Findings: []resultmodel.CommandFinding{{
				Code:                 code,
				Severity:             resultmodel.SeverityError,
				AffectedPaths:        []string{parsed.values["--target"]},
				Evidence:             []string{replaceErr.Error()},
				Fixability:           fixability,
				AutomationStopReason: stopReason,
				NextArgv: []string{"do-work-cli", "--format", "text", CommandReplaceSection,
					"--target", parsed.values["--target"], "--section-file", parsed.values["--section-file"]},
				VerificationArgv: []string{"do-work-cli", "--format", "json", CommandReplaceSection,
					"--target", parsed.values["--target"], "--section-file", parsed.values["--section-file"]},
			}},
		}
	}
	result := resultmodel.CommandResult{
		Command:        CommandReplaceSection,
		Outcome:        resultmodel.OutcomeSuccess,
		RepositoryRoot: context.RepositoryRoot,
		Rollback:       resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
	}
	if outcome.Changed {
		result.Changes = []resultmodel.RecordedChange{{
			Path:   parsed.values["--target"],
			Kind:   outcome.Kind,
			Detail: "managed section written from " + parsed.values["--section-file"],
		}}
	}
	return result
}

func handleValidateManifest(context commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	parsed, err := parseCommandOptions(arguments, []string{"--root"}, nil)
	if err != nil {
		return usageResult(CommandValidateManifest, context.RepositoryRoot, err.Error())
	}
	archiveRoot := parsed.values["--root"]
	if archiveRoot == "" {
		return usageResult(CommandValidateManifest, context.RepositoryRoot,
			"usage: validate-suite-manifest.sh --root <archive-root>")
	}
	validation, validateErr := suitemanifest.ValidateSuite(archiveRoot)
	if validateErr != nil {
		return resultmodel.CommandResult{
			Command:        CommandValidateManifest,
			Outcome:        resultmodel.OutcomeFailure,
			RepositoryRoot: context.RepositoryRoot,
			Rollback:       resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
			Findings: []resultmodel.CommandFinding{{
				Code:                 "SUITE-MANIFEST-INVALID",
				Severity:             resultmodel.SeverityError,
				AffectedPaths:        []string{archiveRoot},
				Evidence:             []string{validateErr.Error()},
				Fixability:           resultmodel.FixabilityManual,
				AutomationStopReason: "the suite archive does not describe a complete four-module do-work suite",
				NextArgv:             []string{"do-work-cli", "--format", "text", CommandValidateManifest, "--root", archiveRoot},
				VerificationArgv:     []string{"do-work-cli", "--format", "json", CommandValidateManifest, "--root", archiveRoot},
			}},
		}
	}
	return resultmodel.CommandResult{
		Command:        CommandValidateManifest,
		Outcome:        resultmodel.OutcomeSuccess,
		RepositoryRoot: context.RepositoryRoot,
		Rollback:       resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
		Changes: []resultmodel.RecordedChange{{
			Path:   archiveRoot,
			Kind:   "validated",
			Detail: validation.SummaryLine(),
		}},
	}
}

// handleFetchArchive reports the winning route as a recorded change on the archive path, so
// the route name stays machine-readable on stdout without inventing an info-severity finding
// on a success outcome.
func handleFetchArchive(context commandruntime.ExecutionContext, arguments []string, narration io.Writer) resultmodel.CommandResult {
	parsed, err := parseCommandOptions(arguments, []string{"--target", "--url", "--repo-url"}, nil)
	if err != nil {
		return usageResult(CommandFetchArchive, context.RepositoryRoot, err.Error())
	}
	targetPath := parsed.values["--target"]
	tarballURL := parsed.values["--url"]
	if targetPath == "" || tarballURL == "" {
		return usageResult(CommandFetchArchive, context.RepositoryRoot,
			"usage: fetch-upstream-archive.sh <archive-target-path> <upstream-tarball-url> [upstream-repo-url]")
	}
	result, fetchErr := archivefetch.FetchArchive(commandContext(), archivefetch.Request{
		ArchiveTargetPath:     targetPath,
		UpstreamTarballURL:    tarballURL,
		UpstreamRepositoryURL: parsed.values["--repo-url"],
		AtomicDownloadScript:  archivefetch.LocateAtomicDownloadScript(shippedToolsDirectory()),
	})
	if fetchErr != nil {
		if narration != nil {
			fmt.Fprintf(narration, "%s\n", fetchErr.Error())
		}
		return resultmodel.CommandResult{
			Command:        CommandFetchArchive,
			Outcome:        resultmodel.OutcomeFailure,
			RepositoryRoot: context.RepositoryRoot,
			Rollback:       resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
			Findings: []resultmodel.CommandFinding{{
				Code:                 "UPSTREAM-ARCHIVE-UNREACHABLE",
				Severity:             resultmodel.SeverityError,
				AffectedPaths:        []string{targetPath},
				Evidence:             []string{fetchErr.Error()},
				Fixability:           resultmodel.FixabilityManual,
				AutomationStopReason: "no fetch route produced a readable suite archive, so the target was left untouched",
				NextArgv: []string{"do-work-cli", "--format", "text", CommandFetchArchive,
					"--target", targetPath, "--url", tarballURL},
				VerificationArgv: []string{"tar", "tzf", targetPath},
			}},
		}
	}
	return resultmodel.CommandResult{
		Command:        CommandFetchArchive,
		Outcome:        resultmodel.OutcomeSuccess,
		RepositoryRoot: context.RepositoryRoot,
		Rollback:       resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
		Changes: []resultmodel.RecordedChange{{
			Path:   targetPath,
			Kind:   "fetched",
			Detail: result.RouteDescription,
		}},
	}
}

func usageResult(commandName, repositoryRoot, evidence string) resultmodel.CommandResult {
	return resultmodel.CommandResult{
		Command:        commandName,
		Outcome:        resultmodel.OutcomeFailure,
		RepositoryRoot: repositoryRoot,
		Rollback:       resultmodel.RollbackResult{Status: resultmodel.RollbackNotNeeded},
		Findings: []resultmodel.CommandFinding{{
			Code:                 "INVALID-COMMAND-OPTIONS",
			Severity:             resultmodel.SeverityError,
			Evidence:             []string{evidence},
			Fixability:           resultmodel.FixabilityManual,
			AutomationStopReason: "the command line is not valid",
			NextArgv:             []string{"do-work-cli", "--format", "text", commandName},
			VerificationArgv:     []string{"do-work-cli", "--format", "json", commandName},
		}},
	}
}

// shippedToolsDirectory answers where the launcher scripts sit, which is what the two
// atomic-download probes are relative to. The binary lives one level below them, inside its
// own module directory, in both the repository and the installed layout. It is resolved at
// handler entry because the install transaction later removes the directory the running
// binary sits in.
func shippedToolsDirectory() string {
	executablePath, err := os.Executable()
	if err != nil {
		return ""
	}
	if resolved, err := filepath.EvalSymlinks(executablePath); err == nil {
		executablePath = resolved
	}
	return filepath.Dir(filepath.Dir(executablePath))
}

func commandContext() context.Context { return context.Background() }
