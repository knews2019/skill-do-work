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

// SetAsideReasonCode marks a finalization record this run excluded from
// selection instead of stopping for. It rides in the record's reason codes, so
// every selector that already reads them sees the exclusion without a new
// field. Claim recovery reads it too: a REQ this run excluded keeps its claim
// (internal/lifecycleadvance).
const SetAsideReasonCode = "FINALIZATION-SET-ASIDE"

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
		if !consumeRecoveryRecord(&aggregate, journal.Manifest.RequestID, result) {
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
			if !consumeRecoveryRecord(&aggregate, journal.Manifest.RequestID, result) {
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

// consumeRecoveryRecord folds one journal's outcome into the aggregate and
// reports whether recovery may keep draining the remaining records. A refusal
// that belongs to a single REQ becomes that REQ's exclusion; anything else is
// evidence the next REQ would trip over, so it stops the run.
func consumeRecoveryRecord(aggregate *resultmodel.CommandResult, requestID string, result resultmodel.CommandResult) bool {
	aggregate.Changes = append(aggregate.Changes, result.Changes...)
	if result.Outcome != resultmodel.OutcomeSuccess && requestScopedRefusal(requestID, result) {
		record, finding := setAsideProjection(requestID, result)
		aggregate.Findings = append(aggregate.Findings, finding)
		appendFinalizationRecord(aggregate, record)
		return true
	}
	aggregate.Findings = append(aggregate.Findings, result.Findings...)
	appendFinalizationResult(aggregate, result)
	if result.Outcome == resultmodel.OutcomeSuccess {
		return true
	}
	aggregate.Outcome = result.Outcome
	return false
}

// requestScopedRefusal answers whether a refused record belongs to exactly one
// REQ. Ownership is REQ-514's test — every finding names this REQ and no other
// — plus proof that the attempt left nothing behind: an incomplete rollback is
// residue on paths the next claim would write through, and a command-level
// failure is not a per-record verdict at all. A refusal whose cause is shared
// state is produced without ownership (sharedStateRefusal, discoveryRefusal),
// so this returns false for it and the whole run stops, which is the one stop
// the maintainer's principle allows.
func requestScopedRefusal(requestID string, result resultmodel.CommandResult) bool {
	if requestID == "" || result.Finalization == nil || len(result.Findings) == 0 {
		return false
	}
	if result.Outcome == resultmodel.OutcomeFailure || result.Rollback.Status == resultmodel.RollbackIncomplete {
		return false
	}
	for _, finding := range result.Findings {
		if !namesOnlyRequestID(finding.AffectedIDs, requestID) {
			return false
		}
	}
	return true
}

// setAsideProjection turns one REQ-scoped refusal into the exclusion a selector
// reads. The record keeps its own reason codes and gains the set-aside code; the
// finding names no next verb, because the verb that resolves a refused
// finalization tail is a judgment the exit summary makes, not one this command
// can pick (REQ-514: a refusal never names itself as the fix).
func setAsideProjection(requestID string, result resultmodel.CommandResult) (resultmodel.FinalizationResult, resultmodel.CommandFinding) {
	record := *result.Finalization
	record.ReasonCodes = append(append([]string(nil), record.ReasonCodes...), SetAsideReasonCode)
	record.NextArgv = nil
	evidence := []string{}
	for _, finding := range result.Findings {
		evidence = append(evidence, finding.Evidence...)
	}
	evidence = append(evidence, "rollback: "+rollbackLabel(result.Rollback.Status))
	return record, resultmodel.CommandFinding{
		Code: SetAsideReasonCode, Severity: resultmodel.SeverityWarning,
		AffectedIDs: []string{requestID}, AffectedPaths: append([]string(nil), record.BlockedPaths...),
		Evidence: evidence, Fixability: resultmodel.FixabilityManual,
		AutomationStopReason: requestID + " is set aside for this run; its finalization tail refused and the remaining REQs continue",
		VerificationArgv:     append([]string(nil), record.VerificationArgv...),
	}
}

func rollbackLabel(status resultmodel.RollbackStatus) string {
	if status == "" {
		return string(resultmodel.RollbackNotNeeded)
	}
	return string(status)
}

// namesOnlyRequestID is exclusive, not a membership test: a finding that names
// this REQ alongside another one is not this REQ's private exclusion, and
// setting it aside would hide the other REQ's evidence.
func namesOnlyRequestID(affectedIDs []string, requestID string) bool {
	return len(affectedIDs) == 1 && affectedIDs[0] == requestID
}

func appendFinalizationResult(aggregate *resultmodel.CommandResult, result resultmodel.CommandResult) {
	if result.Finalization == nil {
		return
	}
	appendFinalizationRecord(aggregate, *result.Finalization)
}

func appendFinalizationRecord(aggregate *resultmodel.CommandResult, record resultmodel.FinalizationResult) {
	aggregate.Finalizations = append(aggregate.Finalizations, record)
	if len(aggregate.Finalizations) == 1 {
		first := aggregate.Finalizations[0]
		aggregate.Finalization = &first
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
