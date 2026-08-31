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

type SelectionProbeStatus string

const (
	ProbeNotApplicable SelectionProbeStatus = "not_applicable"
	ProbeMissing       SelectionProbeStatus = "missing"
	ProbeSucceeded     SelectionProbeStatus = "succeeded"
	ProbeFailed        SelectionProbeStatus = "failed"
	ProbeTimedOut      SelectionProbeStatus = "timed_out"
	ProbeLaunchFailed  SelectionProbeStatus = "launch_failed"
)

// SelectionRecord is one request the caller may process in this invocation.
// Commands are carried as argv rather than display strings so every renderer
// preserves the same pasteable next action.
type SelectionRecord struct {
	RequestID        string               `json:"request_id"`
	RequestPath      string               `json:"request_path"`
	Title            string               `json:"title"`
	Provenance       string               `json:"provenance"`
	OriginalStatus   string               `json:"original_status"`
	ProbeStatus      SelectionProbeStatus `json:"probe_status"`
	ProbeAttempted   bool                 `json:"probe_attempted"`
	ProbeExitCode    int                  `json:"probe_exit_code"`
	UnblockRequired  bool                 `json:"unblock_required"`
	DependencyDepth  int                  `json:"dependency_depth"`
	Dependencies     []string             `json:"dependencies"`
	EstimateMinutes  int                  `json:"estimate_minutes"`
	EstimateKnown    bool                 `json:"estimate_known"`
	NextArgv         []string             `json:"next_argv"`
	NextJustRecipe   string               `json:"next_just_recipe"`
	VerificationArgv []string             `json:"verification_argv"`
}

// SelectionExclusion is one considered request that cannot be selected. Code
// is stable for machines; Reason is actionable for people.
type SelectionExclusion struct {
	RequestID        string               `json:"request_id"`
	RequestPath      string               `json:"request_path"`
	Title            string               `json:"title"`
	Provenance       string               `json:"provenance"`
	OriginalStatus   string               `json:"original_status"`
	ProbeStatus      SelectionProbeStatus `json:"probe_status"`
	ProbeAttempted   bool                 `json:"probe_attempted"`
	ProbeExitCode    int                  `json:"probe_exit_code"`
	UnblockRequired  bool                 `json:"unblock_required"`
	Code             string               `json:"code"`
	Reason           string               `json:"reason"`
	NextArgv         []string             `json:"next_argv"`
	NextJustRecipe   string               `json:"next_just_recipe"`
	VerificationArgv []string             `json:"verification_argv"`
}

// SelectionSummary is the queue projection rendered beside selected and
// excluded records. It is computed from the same snapshot as the records.
type SelectionSummary struct {
	Pending                 int `json:"pending"`
	FinishedAwaitingArchive int `json:"finished_awaiting_archive"`
	PendingAnswers          int `json:"pending_answers"`
	Blocked                 int `json:"blocked"`
	BlockedArchiveCollision int `json:"blocked_archive_collision"`
	BlockedDependencyCycle  int `json:"blocked_dependency_cycle"`
	Probed                  int `json:"probed"`
	ProbeSucceeded          int `json:"probe_succeeded"`
	SkippedImpactNegligible int `json:"skipped_impact_negligible"`
	TotalEstimatedMinutes   int `json:"total_estimated_minutes"`
	UnknownEstimateCount    int `json:"unknown_estimate_count"`
}

type RollbackStatus string

const (
	RollbackNotNeeded  RollbackStatus = "not_needed"
	RollbackSucceeded  RollbackStatus = "succeeded"
	RollbackIncomplete RollbackStatus = "incomplete"
)

type RollbackResult struct {
	Status  RollbackStatus `json:"status"`
	Actions []string       `json:"actions"`
	Errors  []string       `json:"errors"`
}

type CommandResult struct {
	SchemaVersion    int                  `json:"schema_version"`
	Command          string               `json:"command"`
	Outcome          CommandOutcome       `json:"outcome"`
	RepositoryRoot   string               `json:"repository_root"`
	Findings         []CommandFinding     `json:"findings"`
	Changes          []RecordedChange     `json:"changes"`
	SkippedWork      []SkippedWork        `json:"skipped_work"`
	Selected         []SelectionRecord    `json:"selected"`
	Excluded         []SelectionExclusion `json:"excluded"`
	SelectionSummary SelectionSummary     `json:"selection_summary"`
	Rollback         RollbackResult       `json:"rollback"`
}

// ExitCode is the single authority for the 0-4 process status contract. Nothing else in
// this module may map an outcome to a number; a second table is how the two drift apart.
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
	if result.Selected == nil {
		result.Selected = []SelectionRecord{}
	}
	if result.Excluded == nil {
		result.Excluded = []SelectionExclusion{}
	}
	for index := range result.Selected {
		selection := &result.Selected[index]
		if selection.ProbeStatus == "" {
			selection.ProbeStatus = ProbeNotApplicable
		}
		if !selection.ProbeAttempted {
			selection.ProbeExitCode = -1
		}
		if selection.Dependencies == nil {
			selection.Dependencies = []string{}
		}
		if selection.NextArgv == nil {
			selection.NextArgv = []string{}
		}
		if selection.VerificationArgv == nil {
			selection.VerificationArgv = []string{}
		}
	}
	for index := range result.Excluded {
		exclusion := &result.Excluded[index]
		if exclusion.ProbeStatus == "" {
			exclusion.ProbeStatus = ProbeNotApplicable
		}
		if !exclusion.ProbeAttempted {
			exclusion.ProbeExitCode = -1
		}
		if exclusion.NextArgv == nil {
			exclusion.NextArgv = []string{}
		}
		if exclusion.VerificationArgv == nil {
			exclusion.VerificationArgv = []string{}
		}
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
	if result.Command == "next" || len(result.Selected) > 0 || len(result.Excluded) > 0 || result.SelectionSummary.Pending > 0 || result.SelectionSummary.Blocked > 0 {
		fmt.Fprintf(&output, "queue: %d pending | %d finished (awaiting archive) | %d pending-answers | %d blocked | %d blocked-archive-collision | %d blocked-dependency-cycle\n",
			result.SelectionSummary.Pending, result.SelectionSummary.FinishedAwaitingArchive,
			result.SelectionSummary.PendingAnswers, result.SelectionSummary.Blocked,
			result.SelectionSummary.BlockedArchiveCollision, result.SelectionSummary.BlockedDependencyCycle)
	}
	selectedIDs := make([]string, 0, len(result.Selected))
	for _, selection := range result.Selected {
		estimate := "not yet estimated"
		if selection.EstimateKnown {
			estimate = fmt.Sprintf("%d min", selection.EstimateMinutes)
		}
		fmt.Fprintf(&output, "selected %s [%s, depth %d, %s]: %s\n", selection.RequestID, selection.Provenance, selection.DependencyDepth, estimate, selection.Title)
		fmt.Fprintf(&output, "  request: %s (original status: %s)\n", selection.RequestPath, selection.OriginalStatus)
		fmt.Fprintf(&output, "  probe %s (attempted: %t, exit: %d)", selection.ProbeStatus, selection.ProbeAttempted, selection.ProbeExitCode)
		if selection.UnblockRequired {
			fmt.Fprint(&output, "; unblock required")
		}
		fmt.Fprintln(&output)
		fmt.Fprintf(&output, "  next: %s\n", joinArgv(selection.NextArgv))
		if selection.NextJustRecipe != "" {
			fmt.Fprintf(&output, "  just: just %s\n", selection.NextJustRecipe)
		}
		fmt.Fprintf(&output, "  verify: %s\n", joinArgv(selection.VerificationArgv))
		selectedIDs = append(selectedIDs, selection.RequestID)
	}
	for _, exclusion := range result.Excluded {
		fmt.Fprintf(&output, "excluded %s [%s] %s: %s\n", exclusion.RequestID, exclusion.Provenance, exclusion.Code, exclusion.Reason)
		fmt.Fprintf(&output, "  request: %s (original status: %s)\n", exclusion.RequestPath, exclusion.OriginalStatus)
		fmt.Fprintf(&output, "  probe %s (attempted: %t, exit: %d)", exclusion.ProbeStatus, exclusion.ProbeAttempted, exclusion.ProbeExitCode)
		if exclusion.UnblockRequired {
			fmt.Fprint(&output, "; unblock required")
		}
		fmt.Fprintln(&output)
		if len(exclusion.NextArgv) > 0 {
			fmt.Fprintf(&output, "  next: %s\n", joinArgv(exclusion.NextArgv))
		}
		if exclusion.NextJustRecipe != "" {
			fmt.Fprintf(&output, "  just: just %s\n", exclusion.NextJustRecipe)
		}
		if len(exclusion.VerificationArgv) > 0 {
			fmt.Fprintf(&output, "  verify: %s\n", joinArgv(exclusion.VerificationArgv))
		}
	}
	if result.Command == "next" || len(result.Selected) > 0 || len(result.Excluded) > 0 {
		fmt.Fprintf(&output, "run_set: %s\n", strings.Join(selectedIDs, " "))
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
