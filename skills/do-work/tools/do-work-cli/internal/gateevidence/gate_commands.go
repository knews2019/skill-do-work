package gateevidence

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	CommandCheckGreenGate  = "check-green-gate"
	CommandRecordGreenGate = "record-green-gate"
)

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{
		CommandCheckGreenGate:  handleCheckGreenGate,
		CommandRecordGreenGate: handleRecordGreenGate,
	}
}

func handleCheckGreenGate(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	gateCommand, err := parseCheckArguments(arguments)
	if err != nil {
		return gateEvidenceFailure(CommandCheckGreenGate, resultmodel.GateEvidenceResult{}, "GATE-EVIDENCE-USAGE", err)
	}
	evidence, err := CheckGreenGate(executionContext.RepositoryRoot, gateCommand)
	if err != nil {
		return gateEvidenceFailure(CommandCheckGreenGate, evidence, "GATE-EVIDENCE-CHECK-FAILED", err)
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, GateEvidence: &evidence}
}

func handleRecordGreenGate(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	gateExitStatus, gateCommand, err := parseRecordArguments(arguments)
	if err != nil {
		return gateEvidenceFailure(CommandRecordGreenGate, resultmodel.GateEvidenceResult{}, "GATE-EVIDENCE-USAGE", err)
	}
	if gateExitStatus != 0 {
		evidence := resultmodel.GateEvidenceResult{
			GateCommand: append([]string(nil), gateCommand...), GateCommandSHA256: GateCommandSHA256(gateCommand),
			GateExitStatus: gateExitStatus, State: resultmodel.GateEvidenceNotGreen, MatchBasis: "none",
		}
		verificationArgv := append([]string{"do-work-cli", "--format", "json", CommandCheckGreenGate, "--"}, gateCommand...)
		finding := resultmodel.CommandFinding{
			Code: "GATE-EVIDENCE-NOT-GREEN", Severity: resultmodel.SeverityError,
			Evidence:   []string{fmt.Sprintf("direct gate exit status was %d; green evidence was not written", gateExitStatus)},
			Fixability: resultmodel.FixabilityRefused, AutomationStopReason: "only a direct zero gate exit authorizes a green record",
			NextArgv: append([]string(nil), gateCommand...), VerificationArgv: verificationArgv,
		}
		return resultmodel.CommandResult{Outcome: resultmodel.OutcomeRefused, Findings: []resultmodel.CommandFinding{finding}, GateEvidence: &evidence}
	}
	evidence, err := RecordGreenGate(executionContext.RepositoryRoot, gateCommand)
	if err != nil {
		return gateEvidenceFailure(CommandRecordGreenGate, evidence, "GATE-EVIDENCE-RECORD-FAILED", err)
	}
	change := resultmodel.RecordedChange{Path: evidence.RecordPath, Kind: "git-private", Detail: "recorded the current revision for the exact green gate argv"}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Changes: []resultmodel.RecordedChange{change}, GateEvidence: &evidence}
}

func parseCheckArguments(arguments []string) ([]string, error) {
	if len(arguments) < 2 || arguments[0] != "--" {
		return nil, fmt.Errorf("usage: %s -- <gate argv...>", CommandCheckGreenGate)
	}
	return append([]string(nil), arguments[1:]...), nil
}

func parseRecordArguments(arguments []string) (int, []string, error) {
	gateExitStatus := 0
	statusSeen := false
	separatorIndex := -1
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		if argument == "--" {
			separatorIndex = index
			break
		}
		var value string
		switch {
		case argument == "--gate-exit-status":
			index++
			if index >= len(arguments) {
				return 0, nil, fmt.Errorf("--gate-exit-status requires an integer")
			}
			value = arguments[index]
		case strings.HasPrefix(argument, "--gate-exit-status="):
			value = strings.TrimPrefix(argument, "--gate-exit-status=")
		default:
			return 0, nil, fmt.Errorf("unknown record-green-gate option %q", argument)
		}
		if statusSeen {
			return 0, nil, fmt.Errorf("--gate-exit-status may be supplied only once")
		}
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, nil, fmt.Errorf("--gate-exit-status requires an integer")
		}
		gateExitStatus = parsed
		statusSeen = true
	}
	if !statusSeen || separatorIndex < 0 || separatorIndex == len(arguments)-1 {
		return 0, nil, fmt.Errorf("usage: %s --gate-exit-status <status> -- <gate argv...>", CommandRecordGreenGate)
	}
	return gateExitStatus, append([]string(nil), arguments[separatorIndex+1:]...), nil
}

func gateEvidenceFailure(commandName string, evidence resultmodel.GateEvidenceResult, code string, err error) resultmodel.CommandResult {
	if evidence.State == "" {
		evidence.State = resultmodel.GateEvidenceInvalidRecord
	}
	if evidence.MatchBasis == "" {
		evidence.MatchBasis = "none"
	}
	finding := resultmodel.CommandFinding{
		Code: code, Severity: resultmodel.SeverityError, Evidence: []string{err.Error()},
		Fixability: resultmodel.FixabilityManual, AutomationStopReason: "green-gate evidence is unverifiable",
		NextArgv: []string{"git", "status", "--short"}, VerificationArgv: []string{"git", "rev-parse", "--verify", "HEAD"},
	}
	if code == "GATE-EVIDENCE-USAGE" {
		finding.AutomationStopReason = "the command line is invalid"
		finding.NextArgv = []string{"do-work-cli", commandName, "--", "<gate argv...>"}
		finding.VerificationArgv = []string{"do-work-cli", "--format", "json", commandName, "--", "<gate argv...>"}
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{finding}, GateEvidence: &evidence}
}
