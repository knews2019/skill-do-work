package resultmodel

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const SchemaVersion = 1

type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
)

type CommandOutcome string

const (
	OutcomeSuccess    CommandOutcome = "success"
	OutcomeFindings   CommandOutcome = "findings"
	OutcomeRefused    CommandOutcome = "refused"
	OutcomeFailure    CommandOutcome = "failure"
	OutcomeRolledBack CommandOutcome = "rolled_back"
	OutcomeRisk       CommandOutcome = "committed_state_risk"
)

type FindingSeverity string

const (
	SeverityInfo    FindingSeverity = "info"
	SeverityWarning FindingSeverity = "warning"
	SeverityError   FindingSeverity = "error"
)

type FindingFixability string

const (
	FixabilityAutomatic FindingFixability = "automatic"
	FixabilityManual    FindingFixability = "manual"
	FixabilityRefused   FindingFixability = "safely_refused"
)

type CommandFinding struct {
	Code                 string            `json:"code"`
	Severity             FindingSeverity   `json:"severity"`
	AffectedIDs          []string          `json:"affected_ids"`
	AffectedPaths        []string          `json:"affected_paths"`
	Evidence             []string          `json:"observed_evidence"`
	Fixability           FindingFixability `json:"fixability"`
	AutomationStopReason string            `json:"automation_stop_reason"`
	NextArgv             []string          `json:"next_argv"`
	NextJustRecipe       string            `json:"next_just_recipe"`
	VerificationArgv     []string          `json:"verification_argv"`
}

type RecordedChange struct {
	Path   string `json:"path"`
	Kind   string `json:"kind"`
	Detail string `json:"detail"`
}

type SkippedWork struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
}

type RollbackResult struct {
	Status  string   `json:"status"`
	Actions []string `json:"actions"`
	Errors  []string `json:"errors"`
}

type CommandResult struct {
	SchemaVersion  int              `json:"schema_version"`
	Command        string           `json:"command"`
	Outcome        CommandOutcome   `json:"outcome"`
	RepositoryRoot string           `json:"repository_root"`
	Findings       []CommandFinding `json:"findings"`
	Changes        []RecordedChange `json:"changes"`
	SkippedWork    []SkippedWork    `json:"skipped_work"`
	Rollback       RollbackResult   `json:"rollback"`
}

func ExitCode(outcome CommandOutcome) int {
	switch outcome {
	case OutcomeSuccess:
		return 0
	case OutcomeFindings, OutcomeRefused:
		return 1
	case OutcomeRolledBack:
		return 3
	case OutcomeRisk:
		return 4
	case OutcomeFailure:
		return 2
	default:
		return 2
	}
}

func NormalizeResult(result CommandResult) CommandResult {
	result.SchemaVersion = SchemaVersion
	if result.Findings == nil {
		result.Findings = []CommandFinding{}
	}
	if result.Changes == nil {
		result.Changes = []RecordedChange{}
	}
	if result.SkippedWork == nil {
		result.SkippedWork = []SkippedWork{}
	}
	if result.Rollback.Actions == nil {
		result.Rollback.Actions = []string{}
	}
	if result.Rollback.Errors == nil {
		result.Rollback.Errors = []string{}
	}
	for index := range result.Findings {
		finding := &result.Findings[index]
		if finding.AffectedIDs == nil {
			finding.AffectedIDs = []string{}
		}
		if finding.AffectedPaths == nil {
			finding.AffectedPaths = []string{}
		}
		if finding.Evidence == nil {
			finding.Evidence = []string{}
		}
		if finding.NextArgv == nil {
			finding.NextArgv = []string{}
		}
		if finding.VerificationArgv == nil {
			finding.VerificationArgv = []string{}
		}
	}
	return result
}

func RenderResult(result CommandResult, outputFormat OutputFormat) ([]byte, error) {
	normalized := NormalizeResult(result)
	switch outputFormat {
	case FormatJSON:
		output, err := json.MarshalIndent(normalized, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("render JSON result: %w", err)
		}
		return append(output, '\n'), nil
	case FormatText:
		return renderText(normalized), nil
	default:
		return nil, fmt.Errorf("unsupported output format %q", outputFormat)
	}
}

func renderText(result CommandResult) []byte {
	var output strings.Builder
	fmt.Fprintf(&output, "%s: %s\n", result.Command, result.Outcome)
	fmt.Fprintf(&output, "repository: %s\n", result.RepositoryRoot)
	for _, finding := range result.Findings {
		fmt.Fprintf(&output, "finding %s [%s]: %s\n", finding.Code, finding.Severity, strings.Join(finding.Evidence, "; "))
		if len(finding.AffectedIDs) > 0 {
			fmt.Fprintf(&output, "  ids: %s\n", strings.Join(finding.AffectedIDs, ", "))
		}
		if len(finding.AffectedPaths) > 0 {
			fmt.Fprintf(&output, "  paths: %s\n", strings.Join(finding.AffectedPaths, ", "))
		}
		fmt.Fprintf(&output, "  fixability: %s\n", finding.Fixability)
		if finding.AutomationStopReason != "" {
			fmt.Fprintf(&output, "  stopped: %s\n", finding.AutomationStopReason)
		}
		if len(finding.NextArgv) > 0 {
			fmt.Fprintf(&output, "  next: %s\n", joinArgv(finding.NextArgv))
		}
		if finding.NextJustRecipe != "" {
			fmt.Fprintf(&output, "  just: just %s\n", finding.NextJustRecipe)
		}
		if len(finding.VerificationArgv) > 0 {
			fmt.Fprintf(&output, "  verify: %s\n", joinArgv(finding.VerificationArgv))
		}
	}
	for _, change := range result.Changes {
		fmt.Fprintf(&output, "change %s [%s]: %s\n", change.Path, change.Kind, change.Detail)
	}
	for _, skipped := range result.SkippedWork {
		fmt.Fprintf(&output, "skipped %s: %s\n", skipped.Code, skipped.Reason)
	}
	if result.Rollback.Status != "" {
		fmt.Fprintf(&output, "rollback: %s\n", result.Rollback.Status)
	}
	for _, action := range result.Rollback.Actions {
		fmt.Fprintf(&output, "  rollback action: %s\n", action)
	}
	for _, rollbackError := range result.Rollback.Errors {
		fmt.Fprintf(&output, "  rollback error: %s\n", rollbackError)
	}
	return []byte(output.String())
}

func joinArgv(argv []string) string {
	quoted := make([]string, len(argv))
	for index, argument := range argv {
		if argument != "" && strings.IndexFunc(argument, func(character rune) bool {
			return !(character == '-' || character == '_' || character == '/' || character == '.' || character == ':' ||
				character >= 'a' && character <= 'z' || character >= 'A' && character <= 'Z' || character >= '0' && character <= '9')
		}) == -1 {
			quoted[index] = argument
		} else {
			quoted[index] = strconv.Quote(argument)
		}
	}
	return strings.Join(quoted, " ")
}
