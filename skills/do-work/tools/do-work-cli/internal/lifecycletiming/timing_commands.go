package lifecycletiming

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	CommandRecordTimingEvent = "record-timing-event"
	CommandRunTimedCommand   = "run-timed-command"
	CommandFoldTimingSummary = "fold-timing-summary"
)

// Handlers registers the three timing verbs. childOutput receives a wrapped
// command's own stdout and stderr, which keeps stdout free for the rendered
// CommandResult; main passes os.Stderr.
func Handlers(childOutput io.Writer) map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{
		CommandRecordTimingEvent: handleRecordTimingEvent,
		CommandRunTimedCommand: func(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
			return handleRunTimedCommand(executionContext, arguments, childOutput)
		},
		CommandFoldTimingSummary: handleFoldTimingSummary,
	}
}

// timingOptions is the parsed argv shared by all three verbs; each verb reads
// only the options it accepts and refuses the rest.
type timingOptions struct {
	requestID     string
	runIdentifier string
	category      string
	operation     string
	startedAt     time.Time
	hasStartedAt  bool
	outcome       string
	revision      string
	agent         string
	exitStatus    *int
	commandText   string
	requestPath   string
	separatorSeen bool
	commandArgv   []string
}

func handleRecordTimingEvent(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, parseError := parseTimingOptions(arguments)
	if parseError != nil {
		return timingUsageFailure(CommandRecordTimingEvent, parseError.Error())
	}
	if options.separatorSeen {
		return timingUsageFailure(CommandRecordTimingEvent, "record-timing-event takes no argv after --; use run-timed-command to time a command")
	}
	if options.category == "" || options.operation == "" {
		return timingUsageFailure(CommandRecordTimingEvent, "--category and --operation are required")
	}
	eventRequest := EventRequest{
		RunIdentifier: options.runIdentifier, RequestID: options.requestID,
		Category: options.category, Operation: options.operation,
		StartedAt: options.startedAt, HasExplicitStart: options.hasStartedAt,
		Outcome: options.outcome, Revision: options.revision, ResponsibleAgent: options.agent,
		ExitStatus: options.exitStatus, CommandArgv: strings.Fields(options.commandText),
	}
	timing, err := AppendTimingEvent(executionContext.RepositoryRoot, eventRequest)
	if err != nil {
		return timingFailure(CommandRecordTimingEvent, "TIMING-EVENT-REFUSED", err.Error())
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, LifecycleTiming: &timing}
}

func handleRunTimedCommand(executionContext commandruntime.ExecutionContext, arguments []string, childOutput io.Writer) resultmodel.CommandResult {
	options, parseError := parseTimingOptions(arguments)
	if parseError != nil {
		return timingUsageFailure(CommandRunTimedCommand, parseError.Error())
	}
	if options.category == "" {
		return timingUsageFailure(CommandRunTimedCommand, "--category is required")
	}
	if !options.separatorSeen || len(options.commandArgv) == 0 {
		return timingUsageFailure(CommandRunTimedCommand, "the command to time is required after --")
	}
	if options.exitStatus != nil || options.commandText != "" {
		return timingUsageFailure(CommandRunTimedCommand, "run-timed-command observes the exit status and command itself; drop --exit-status and --command")
	}
	eventRequest := EventRequest{
		RunIdentifier: options.runIdentifier, RequestID: options.requestID,
		Category: options.category, Operation: options.operation,
		Revision: options.revision, ResponsibleAgent: options.agent,
	}
	timing, exitStatus, err := RunTimedCommand(executionContext.RepositoryRoot, eventRequest, options.commandArgv, childOutput)
	if err != nil {
		if errors.Is(err, errTimedCommandInput) {
			return timingUsageFailure(CommandRunTimedCommand, err.Error())
		}
		if errors.Is(err, errTimingRecording) {
			outcome := resultmodel.OutcomeSuccess
			if exitStatus != 0 {
				outcome = resultmodel.OutcomeFindings
			}
			return resultmodel.CommandResult{
				Outcome: outcome, ExitCodeOverride: exitStatus,
				Findings: []resultmodel.CommandFinding{{
					Code: "TIMED-COMMAND-RECORDING-FAILED", Severity: resultmodel.SeverityWarning,
					AffectedIDs: []string{options.requestID}, Evidence: []string{err.Error()},
					Fixability: resultmodel.FixabilityManual,
				}},
			}
		}
		// A command that never launched exits 127 like a shell's "command not
		// found", so a caller reading the process status still learns that the
		// gate did not run rather than that it failed.
		launchFailure := timingFailure(CommandRunTimedCommand, "TIMED-COMMAND-LAUNCH-FAILED", err.Error())
		launchFailure.ExitCodeOverride = exitStatus
		return launchFailure
	}
	if exitStatus != 0 {
		return resultmodel.CommandResult{
			Outcome: resultmodel.OutcomeFindings, ExitCodeOverride: exitStatus, LifecycleTiming: &timing,
			Findings: []resultmodel.CommandFinding{{
				Code: "TIMED-COMMAND-NONZERO-EXIT", Severity: resultmodel.SeverityWarning,
				AffectedIDs: []string{timing.RequestID},
				Evidence:    []string{fmt.Sprintf("%s exited %d after %ds", timing.RecordedEvent.CommandIdentity, exitStatus, timing.RecordedEvent.ElapsedSeconds)},
				Fixability:  resultmodel.FixabilityManual,
			}},
		}
	}
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, LifecycleTiming: &timing}
}

func handleFoldTimingSummary(executionContext commandruntime.ExecutionContext, arguments []string) resultmodel.CommandResult {
	options, parseError := parseTimingOptions(arguments)
	if parseError != nil {
		return timingUsageFailure(CommandFoldTimingSummary, parseError.Error())
	}
	if options.separatorSeen || options.category != "" || options.operation != "" || options.exitStatus != nil {
		return timingUsageFailure(CommandFoldTimingSummary, "fold-timing-summary accepts only --request, --run and --request-path")
	}
	if options.requestPath == "" {
		return timingUsageFailure(CommandFoldTimingSummary, "--request-path is required")
	}
	timing, skipped, err := FoldTimingSummary(executionContext.RepositoryRoot, FoldRequest{
		RunIdentifier: options.runIdentifier, RequestID: options.requestID, RequestPath: options.requestPath,
	})
	if err != nil {
		return timingFailure(CommandFoldTimingSummary, "TIMING-FOLD-FAILED", err.Error())
	}
	result := resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, LifecycleTiming: &timing, SkippedWork: skipped}
	if timing.SectionWritten {
		result.Changes = []resultmodel.RecordedChange{{
			Path: timing.RequestPath, Kind: "modified",
			Detail: fmt.Sprintf("folded a %s Timing summary over %d events and removed the raw stream", formatElapsedSeconds(timing.TotalObservedSeconds), timing.EventCount),
		}}
	}
	return result
}

// parseTimingOptions reads the union of the three verbs' options. Everything
// after a bare -- is the command to time; nothing before it may be positional.
func parseTimingOptions(arguments []string) (timingOptions, error) {
	options := timingOptions{}
	for index := 0; index < len(arguments); index++ {
		argument := arguments[index]
		if argument == "--" {
			options.separatorSeen = true
			options.commandArgv = append([]string(nil), arguments[index+1:]...)
			break
		}
		name, value, valueError := timingOptionValue(arguments, &index)
		if valueError != nil {
			return options, valueError
		}
		switch name {
		case "--request":
			options.requestID = value
		case "--run":
			options.runIdentifier = value
		case "--category":
			options.category = value
		case "--operation":
			options.operation = value
		case "--started-at":
			parsed, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return options, fmt.Errorf("--started-at requires an RFC3339 instant: %w", err)
			}
			options.startedAt, options.hasStartedAt = parsed.UTC(), true
		case "--outcome":
			options.outcome = value
		case "--revision":
			options.revision = value
		case "--agent":
			options.agent = value
		case "--exit-status":
			parsed, err := strconv.Atoi(value)
			if err != nil {
				return options, fmt.Errorf("--exit-status requires an integer: %w", err)
			}
			options.exitStatus = &parsed
		case "--command":
			options.commandText = value
		case "--request-path":
			options.requestPath = value
		default:
			return options, fmt.Errorf("unknown option %q", name)
		}
	}
	if options.requestID == "" {
		return options, fmt.Errorf("--request REQ-NNN is required")
	}
	return options, nil
}

// timingOptionValue accepts both --name value and --name=value spellings.
func timingOptionValue(arguments []string, index *int) (string, string, error) {
	argument := arguments[*index]
	if !strings.HasPrefix(argument, "--") {
		return "", "", fmt.Errorf("unexpected positional argument %q", argument)
	}
	if name, value, found := strings.Cut(argument, "="); found {
		if value == "" {
			return name, "", fmt.Errorf("%s requires a value", name)
		}
		return name, value, nil
	}
	*index++
	if *index >= len(arguments) {
		return argument, "", fmt.Errorf("%s requires a value", argument)
	}
	return argument, arguments[*index], nil
}

func timingUsageFailure(command, evidence string) resultmodel.CommandResult {
	return timingFailure(command, "TIMING-USAGE", evidence)
}

func timingFailure(command, code, evidence string) resultmodel.CommandResult {
	return resultmodel.CommandResult{
		Outcome: resultmodel.OutcomeFailure,
		Findings: []resultmodel.CommandFinding{{
			Code: code, Severity: resultmodel.SeverityError, Evidence: []string{evidence},
			Fixability:           resultmodel.FixabilityManual,
			AutomationStopReason: "the timing invocation is not valid",
			NextArgv:             []string{"do-work-cli", command, "--request", "REQ-NNN"},
			VerificationArgv:     []string{"do-work-cli", "--format", "json", command, "--request", "REQ-NNN"},
		}},
	}
}
