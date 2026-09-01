package hookcommands

import (
	"io"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	CommandSessionStart       = "session-start"
	CommandMemorySessionStart = "memory-session-start"
	CommandMemoryStopCapture  = "memory-stop-capture"
)

func Handlers(input io.Reader) map[string]commandruntime.CommandHandler {
	if input == nil {
		input = io.LimitReader(zeroReader{}, 0)
	}
	return map[string]commandruntime.CommandHandler{
		CommandSessionStart:       handleSessionStart,
		CommandMemorySessionStart: handleMemorySessionStart,
		CommandMemoryStopCapture: func(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
			return handleMemoryStopCapture(executionContext, arguments, input)
		},
	}
}

type zeroReader struct{}

func (zeroReader) Read([]byte) (int, error) { return 0, io.EOF }

func protocolResult(output string, operations ...resultmodel.CommandResult) resultmodel.CommandResult {
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, ProtocolOutput: &output}
	for _, operation := range operations {
		result.Findings = append(result.Findings, operation.Findings...)
		result.Changes = append(result.Changes, operation.Changes...)
		result.SkippedWork = append(result.SkippedWork, operation.SkippedWork...)
	}
	return result
}

func hookFinding(code, path, evidence, stopReason string) resultmodel.CommandFinding {
	return resultmodel.CommandFinding{
		Code: code, Severity: resultmodel.SeverityWarning, AffectedPaths: []string{path},
		Evidence: []string{evidence}, Fixability: resultmodel.FixabilityManual,
		AutomationStopReason: stopReason,
	}
}

func usageResult(commandName, evidence string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{{
		Code: "HOOK-USAGE", Severity: resultmodel.SeverityError, Evidence: []string{evidence},
		Fixability: resultmodel.FixabilityManual, AutomationStopReason: "the hook command line is invalid",
		NextArgv: []string{"do-work-cli", commandName}, VerificationArgv: []string{"do-work-cli", "--format", "json", commandName},
	}}}
}
