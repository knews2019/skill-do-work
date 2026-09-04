package finalization

import (
	"context"
	"errors"
	"fmt"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	CommandFinalize            = "finalize"
	CommandRecoverFinalization = "recover-finalization"
)

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{
		CommandFinalize:            handleFinalize,
		CommandRecoverFinalization: handleRecoverFinalization,
	}
}

func handleFinalize(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	manifestPath, err := parseFinalizeArguments(arguments)
	if err != nil {
		return commandFailure(executionContext.RepositoryRoot, CommandFinalize, "FINALIZATION-USAGE", err.Error())
	}
	journal, resumed, err := prepareJournal(context.Background(), executionContext.RepositoryRoot, manifestPath)
	if err != nil {
		return commandFailure(executionContext.RepositoryRoot, CommandFinalize, "FINALIZATION-PREPARE", err.Error())
	}
	return advanceJournal(context.Background(), executionContext.RepositoryRoot, journal, resumed)
}

// FinalizeBound runs the canonical finalization transaction only when the
// action-authored manifest identifies the request selected by the caller.
// prepareBoundJournal performs that comparison from the manifest's single
// decode before journal, index, or repository mutation.
func FinalizeBound(executionContext commandruntime.ExecutionContext, manifestPath, expectedRequestID, expectedRequestPath string) resultmodel.CommandResult {
	journal, resumed, err := prepareBoundJournal(context.Background(), executionContext.RepositoryRoot, manifestPath, expectedRequestID, expectedRequestPath)
	if err != nil {
		var bindingError requestBindingError
		if !errors.As(err, &bindingError) {
			return commandFailure(executionContext.RepositoryRoot, CommandFinalize, "FINALIZATION-PREPARE", err.Error())
		}
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, RepositoryRoot: executionContext.RepositoryRoot, Findings: []resultmodel.CommandFinding{{
			Code: "FINALIZATION-REQUEST-MISMATCH", Severity: resultmodel.SeverityError,
			AffectedIDs: []string{expectedRequestID}, AffectedPaths: []string{expectedRequestPath, manifestPath}, Evidence: []string{err.Error()},
			Fixability: resultmodel.FixabilityRefused, AutomationStopReason: "the finalization manifest is not bound to the selected request",
			VerificationArgv: []string{"do-work-cli", "--format", "json", CommandFinalize, "--manifest", manifestPath},
		}}}
	}
	return advanceJournal(context.Background(), executionContext.RepositoryRoot, journal, resumed)
}

func handleRecoverFinalization(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, err := parseRecoverArguments(arguments)
	if err != nil {
		return commandFailure(executionContext.RepositoryRoot, CommandRecoverFinalization, "FINALIZATION-USAGE", err.Error())
	}
	paths, err := listJournals(executionContext.RepositoryRoot)
	if err != nil {
		return commandFailure(executionContext.RepositoryRoot, CommandRecoverFinalization, "FINALIZATION-JOURNAL-LIST", err.Error())
	}
	if len(paths) == 0 && !options.Discover {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: []resultmodel.CommandFinding{{
			Code: "FINALIZATION-NONE", Severity: resultmodel.SeverityInfo, Evidence: []string{"no unfinished finalization journals"},
			Fixability: resultmodel.FixabilityAutomatic, VerificationArgv: []string{"do-work-cli", "--format", "json", CommandRecoverFinalization},
		}}}
	}
	aggregate := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Finalizations: []resultmodel.FinalizationResult{}}
	for _, path := range paths {
		journal, readError := readJournal(executionContext.RepositoryRoot, path)
		if readError != nil {
			return commandFailure(executionContext.RepositoryRoot, CommandRecoverFinalization, "FINALIZATION-JOURNAL-INVALID", readError.Error())
		}
		result := advanceJournal(context.Background(), executionContext.RepositoryRoot, journal, true)
		aggregate.Findings = append(aggregate.Findings, result.Findings...)
		aggregate.Changes = append(aggregate.Changes, result.Changes...)
		appendFinalizationResult(&aggregate, result)
		if result.Outcome != resultmodel.OutcomeSuccess {
			aggregate.Outcome = result.Outcome
			return aggregate
		}
	}
	if options.Discover {
		journals, discoveryResult := discoverFinalizationJournals(executionContext.RepositoryRoot, options.AssumeSoleReleaser)
		if discoveryResult != nil {
			for _, record := range aggregate.Finalizations {
				discoveryResult.Finalizations = append([]resultmodel.FinalizationResult{record}, discoveryResult.Finalizations...)
			}
			if len(discoveryResult.Finalizations) == 1 {
				discoveryResult.Finalization = &discoveryResult.Finalizations[0]
			}
			return *discoveryResult
		}
		for _, journal := range journals {
			result := advanceJournal(context.Background(), executionContext.RepositoryRoot, journal, true)
			aggregate.Findings = append(aggregate.Findings, result.Findings...)
			aggregate.Changes = append(aggregate.Changes, result.Changes...)
			appendFinalizationResult(&aggregate, result)
			if result.Outcome != resultmodel.OutcomeSuccess {
				aggregate.Outcome = result.Outcome
				return aggregate
			}
		}
	}
	if len(aggregate.Finalizations) == 0 {
		aggregate.Findings = append(aggregate.Findings, resultmodel.CommandFinding{Code: "FINALIZATION-NONE", Severity: resultmodel.SeverityInfo,
			Evidence: []string{"no unfinished or safely discoverable finalizations"}, Fixability: resultmodel.FixabilityAutomatic,
			VerificationArgv: []string{"do-work-cli", "--format", "json", CommandRecoverFinalization, "--discover"}})
	}
	return aggregate
}

func parseRecoverArguments(arguments []string) (commandOptions, error) {
	options := commandOptions{}
	for _, argument := range arguments {
		switch argument {
		case "--discover":
			if options.Discover {
				return commandOptions{}, fmt.Errorf("--discover may be supplied only once")
			}
			options.Discover = true
		case "--assume-sole-releaser":
			if options.AssumeSoleReleaser {
				return commandOptions{}, fmt.Errorf("--assume-sole-releaser may be supplied only once")
			}
			options.AssumeSoleReleaser = true
		default:
			return commandOptions{}, fmt.Errorf("unknown recover-finalization option %q", argument)
		}
	}
	if options.AssumeSoleReleaser && !options.Discover {
		return commandOptions{}, fmt.Errorf("--assume-sole-releaser requires --discover")
	}
	return options, nil
}

func appendFinalizationResult(aggregate *resultmodel.CommandResult, result resultmodel.CommandResult) {
	if result.Finalization == nil {
		return
	}
	aggregate.Finalizations = append(aggregate.Finalizations, *result.Finalization)
	if len(aggregate.Finalizations) == 1 {
		record := aggregate.Finalizations[0]
		aggregate.Finalization = &record
	} else {
		aggregate.Finalization = nil
	}
}

func parseFinalizeArguments(arguments []string) (string, error) {
	manifestPath := ""
	for index := 0; index < len(arguments); index++ {
		if arguments[index] != "--manifest" {
			return "", fmt.Errorf("unknown finalize option %q", arguments[index])
		}
		index++
		if index >= len(arguments) || manifestPath != "" {
			return "", fmt.Errorf("--manifest requires one JSON file")
		}
		manifestPath = arguments[index]
	}
	if manifestPath == "" {
		return "", fmt.Errorf("--manifest requires one JSON file")
	}
	return manifestPath, nil
}

func commandFailure(repositoryRoot, command, code, reason string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, RepositoryRoot: repositoryRoot, Findings: []resultmodel.CommandFinding{{
		Code: code, Severity: resultmodel.SeverityError, Evidence: []string{reason}, Fixability: resultmodel.FixabilityManual,
		AutomationStopReason: "resumable finalization could not start safely", NextArgv: []string{"do-work-cli", command},
		VerificationArgv: []string{"do-work-cli", "--format", "json", command},
	}}}
}
