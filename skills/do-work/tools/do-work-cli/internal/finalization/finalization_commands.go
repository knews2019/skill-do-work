package finalization

import (
	"context"
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

func handleRecoverFinalization(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	if len(arguments) != 0 {
		return commandFailure(executionContext.RepositoryRoot, CommandRecoverFinalization, "FINALIZATION-USAGE", "recover-finalization accepts no options in the journal-replay slice")
	}
	paths, err := listJournals(executionContext.RepositoryRoot)
	if err != nil {
		return commandFailure(executionContext.RepositoryRoot, CommandRecoverFinalization, "FINALIZATION-JOURNAL-LIST", err.Error())
	}
	if len(paths) == 0 {
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Findings: []resultmodel.CommandFinding{{
			Code: "FINALIZATION-NONE", Severity: resultmodel.SeverityInfo, Evidence: []string{"no unfinished finalization journals"},
			Fixability: resultmodel.FixabilityAutomatic, VerificationArgv: []string{"do-work-cli", "--format", "json", CommandRecoverFinalization},
		}}}
	}
	aggregate := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess}
	for _, path := range paths {
		journal, readError := readJournal(path)
		if readError != nil {
			return commandFailure(executionContext.RepositoryRoot, CommandRecoverFinalization, "FINALIZATION-JOURNAL-INVALID", readError.Error())
		}
		result := advanceJournal(context.Background(), executionContext.RepositoryRoot, journal, true)
		aggregate.Findings = append(aggregate.Findings, result.Findings...)
		aggregate.Changes = append(aggregate.Changes, result.Changes...)
		aggregate.Finalization = result.Finalization
		if result.Outcome != resultmodel.OutcomeSuccess {
			aggregate.Outcome = result.Outcome
			return aggregate
		}
	}
	return aggregate
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
