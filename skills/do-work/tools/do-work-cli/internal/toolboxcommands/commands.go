package toolboxcommands

import (
	"fmt"
	"strings"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

const (
	CommandNote             = "do-work-note"
	CommandArchitecture     = "architecture-report-preflight"
	CommandReportImage      = "generate-report-image"
	CommandReportImageBatch = "generate-report-image-batch"
	CommandPortfolio        = "publish-portfolio-summary"
	CommandLast30Days       = "install-last30days"
	CommandAuditMetrics     = "audit-metrics"
)

func Handlers() map[string]commandruntime.CommandHandler {
	return map[string]commandruntime.CommandHandler{
		CommandNote:             handleNote,
		CommandArchitecture:     handleArchitecture,
		CommandReportImage:      handleReportImage,
		CommandReportImageBatch: handleReportImageBatch,
		CommandPortfolio:        handlePortfolio,
		CommandLast30Days:       handleLast30Days,
		CommandAuditMetrics:     handleAuditMetrics,
	}
}

func usageResult(command, evidence string) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeFailure, Findings: []resultmodel.CommandFinding{
		toolboxFinding(command, "TOOLBOX-USAGE", resultmodel.SeverityError, nil, evidence,
			resultmodel.FixabilityManual, "the command line is invalid"),
	}}
}

func toolboxFinding(command, code string, severity resultmodel.FindingSeverity, paths []string, evidence string, fixability resultmodel.FindingFixability, stop string) resultmodel.CommandFinding {
	return resultmodel.CommandFinding{
		Code: code, Severity: severity, AffectedPaths: append([]string(nil), paths...),
		Evidence: []string{evidence}, Fixability: fixability, AutomationStopReason: stop,
		NextArgv:         []string{"do-work-cli", command},
		VerificationArgv: []string{"do-work-cli", "--format", "json", command},
	}
}

func parseMutationFlags(arguments []string) (rest []string, dryRun, commit bool, err error) {
	for _, argument := range arguments {
		switch argument {
		case "--dry-run":
			dryRun = true
		case "--commit":
			commit = true
		default:
			rest = append(rest, argument)
		}
	}
	if dryRun && commit {
		err = fmt.Errorf("--dry-run and --commit cannot be combined")
	}
	return
}

func exactOutputResult(output string, changes []resultmodel.RecordedChange) resultmodel.CommandResult {
	return resultmodel.CommandResult{Outcome: resultmodel.OutcomeSuccess, Changes: changes, ExactTextOutput: &output}
}

func optionValue(arguments []string, index *int, name string) (string, error) {
	argument := arguments[*index]
	if strings.HasPrefix(argument, name+"=") {
		value := strings.TrimPrefix(argument, name+"=")
		if value == "" {
			return "", fmt.Errorf("%s requires a value", name)
		}
		return value, nil
	}
	*index++
	if *index >= len(arguments) {
		return "", fmt.Errorf("%s requires a value", name)
	}
	return arguments[*index], nil
}
