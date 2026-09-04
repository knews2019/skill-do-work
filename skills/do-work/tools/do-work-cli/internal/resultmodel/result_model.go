package resultmodel

import (
	"encoding/json"
	"fmt"
	"path/filepath"
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

type HeavyLaneSelection struct {
	LaneID      string   `json:"lane_id"`
	CommandArgv []string `json:"command_argv"`
	Reasons     []string `json:"reasons"`
}

type HeavySourceRange struct {
	BaseRevision   string `json:"base_revision"`
	TargetRevision string `json:"target_revision"`
}

type HeavyVerificationPlan struct {
	Mode              string               `json:"mode,omitempty"`
	SourceRanges      []HeavySourceRange   `json:"source_ranges,omitempty"`
	ManifestRevision  string               `json:"manifest_revision,omitempty"`
	ExecutionRevision string               `json:"execution_revision,omitempty"`
	ManifestPath      string               `json:"manifest_path"`
	BaseRevision      string               `json:"base_revision"`
	TargetRevision    string               `json:"target_revision"`
	ForcedAll         bool                 `json:"forced_all"`
	Uncertain         bool                 `json:"uncertain"`
	ChangedPaths      []string             `json:"changed_paths"`
	UncoveredPaths    []string             `json:"uncovered_paths"`
	SelectedLanes     []HeavyLaneSelection `json:"selected_lanes"`
}

// HeavyLaneExecution is one lane of a run: either executed now, or reported
// from matching stored evidence. Its first five fields are the durable per-lane
// evidence the work action's drain copies onto the claimed record; the reuse
// fields below are run-local and are not part of that evidence.
type HeavyLaneExecution struct {
	LaneID      string   `json:"lane_id"`
	CommandArgv []string `json:"command_argv"`
	ExitStatus  int      `json:"exit_status"`
	Skipped     bool     `json:"skipped,omitempty"`
	WallSeconds int      `json:"wall_seconds"`
	// Disposition is "executed" or "reused", so a reader can tell which greens
	// this run measured and which it inherited from stored evidence.
	Disposition string `json:"disposition"`
	// DispositionReason names the exact condition that decided the
	// disposition: a fingerprint match, or the one check that failed.
	DispositionReason string `json:"disposition_reason"`
	// FingerprintSHA256 is the deterministic digest of the lane's command, its
	// covered committed files, its toolchain probes, and its required
	// environment. Empty when that fingerprint could not be determined.
	FingerprintSHA256 string `json:"fingerprint_sha256,omitempty"`
	// EvidenceRevision and EvidenceRecordedAt name the run a reused lane
	// inherits its result from. Both stay empty for an executed lane.
	EvidenceRevision   string `json:"evidence_revision,omitempty"`
	EvidenceRecordedAt string `json:"evidence_recorded_at,omitempty"`
}

// HeavyVerificationRun is the typed record of one heavy-lane verification pass.
// Planning stays with HeavyVerificationPlan; this type carries only what ran.
type HeavyVerificationRun struct {
	ManifestPath      string               `json:"manifest_path"`
	ExecutionRevision string               `json:"execution_revision"`
	Lanes             []HeavyLaneExecution `json:"lanes"`
}

// AuditMetricsResult is the machine projection of the maintainability-audit
// measurements. Compatibility Markdown and JSON are rendered from this same
// ordered data; consumers never need to parse a table back into numbers.
type AuditMetricsResult struct {
	Kind             string                   `json:"kind"`
	SinceWindow      string                   `json:"since_window"`
	ShallowClone     bool                     `json:"shallow_clone"`
	CommitCount      int                      `json:"commit_count"`
	BinaryCount      int                      `json:"binary_count"`
	UnreadableCount  int                      `json:"unreadable_count"`
	Aggregates       []AuditAggregate         `json:"aggregates"`
	Distributions    []AuditDistribution      `json:"distributions"`
	Files            []AuditFileMeasurement   `json:"files"`
	Folders          []AuditFolderMeasurement `json:"folders"`
	Bands            []AuditBand              `json:"bands"`
	Churn            []AuditChurnMeasurement  `json:"churn"`
	Hotspots         []AuditHotspot           `json:"hotspots"`
	UnavailablePaths []string                 `json:"unavailable_paths"`
}

type AuditAggregate struct {
	Extension string `json:"extension"`
	Files     int    `json:"files"`
	Lines     int    `json:"lines"`
	Words     int    `json:"words"`
}

type AuditDistribution struct {
	Metric string `json:"metric"`
	Median int    `json:"median"`
	P90    int    `json:"p90"`
	P95    int    `json:"p95"`
	Max    int    `json:"max"`
}

type AuditFileMeasurement struct {
	Path   string `json:"path"`
	Lines  int    `json:"lines"`
	Words  int    `json:"words"`
	Binary bool   `json:"binary"`
}

type AuditFolderMeasurement struct {
	Folder string `json:"folder"`
	Files  int    `json:"files"`
}

type AuditBand struct {
	Path   string `json:"path"`
	Metric string `json:"metric"`
	Value  int    `json:"value"`
	Band   string `json:"band"`
}

type AuditChurnMeasurement struct {
	Path    string `json:"path"`
	Commits int    `json:"commits"`
}

type AuditHotspot struct {
	Path    string `json:"path"`
	Commits int    `json:"commits"`
	Lines   int    `json:"lines"`
	Score   int    `json:"score"`
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
	RequestID         string               `json:"request_id"`
	RequestPath       string               `json:"request_path"`
	Title             string               `json:"title"`
	Provenance        string               `json:"provenance"`
	SelectionPriority string               `json:"selection_priority"`
	RequestPriority   string               `json:"priority"`
	OriginalStatus    string               `json:"original_status"`
	ProbeStatus       SelectionProbeStatus `json:"probe_status"`
	ProbeAttempted    bool                 `json:"probe_attempted"`
	ProbeExitCode     int                  `json:"probe_exit_code"`
	UnblockRequired   bool                 `json:"unblock_required"`
	DependencyDepth   int                  `json:"dependency_depth"`
	Dependencies      []string             `json:"dependencies"`
	EstimateMinutes   int                  `json:"estimate_minutes"`
	EstimateKnown     bool                 `json:"estimate_known"`
	NextArgv          []string             `json:"next_argv"`
	NextJustRecipe    string               `json:"next_just_recipe"`
	VerificationArgv  []string             `json:"verification_argv"`
}

// SelectionClaimEvidence is one exact ownership fact that vetoed selection.
type SelectionClaimEvidence struct {
	Source     string `json:"source"`
	ClaimedAt  string `json:"claimed_at"`
	Writer     string `json:"writer"`
	Path       string `json:"path"`
	SourceLine int    `json:"source_line"`
	HeaderText string `json:"header_text"`
}

// SelectionExclusion is one considered request that cannot be selected. Code
// is stable for machines; Reason is actionable for people.
type SelectionExclusion struct {
	RequestID         string                   `json:"request_id"`
	RequestPath       string                   `json:"request_path"`
	Title             string                   `json:"title"`
	Provenance        string                   `json:"provenance"`
	SelectionPriority string                   `json:"selection_priority"`
	RequestPriority   string                   `json:"priority"`
	OriginalStatus    string                   `json:"original_status"`
	ProbeStatus       SelectionProbeStatus     `json:"probe_status"`
	ProbeAttempted    bool                     `json:"probe_attempted"`
	ProbeExitCode     int                      `json:"probe_exit_code"`
	UnblockRequired   bool                     `json:"unblock_required"`
	Code              string                   `json:"code"`
	Reason            string                   `json:"reason"`
	ClaimEvidence     []SelectionClaimEvidence `json:"claim_evidence"`
	NextArgv          []string                 `json:"next_argv"`
	NextJustRecipe    string                   `json:"next_just_recipe"`
	VerificationArgv  []string                 `json:"verification_argv"`
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

type GateDeferralResult struct {
	ParentID                    string   `json:"parent_id"`
	ParentPath                  string   `json:"parent_path"`
	RepairID                    string   `json:"repair_id"`
	RepairPath                  string   `json:"repair_path"`
	CheckpointPath              string   `json:"checkpoint_path"`
	RepairOutcome               string   `json:"repair_outcome"`
	RepairDependency            string   `json:"repair_dependency"`
	DiagnosticFingerprint       string   `json:"diagnostic_fingerprint"`
	SweepKey                    string   `json:"sweep_key"`
	GateCommand                 []string `json:"gate_command"`
	GateExitStatus              int      `json:"gate_exit_status"`
	DeferredImplementationBase  string   `json:"deferred_implementation_base,omitempty"`
	DeferredImplementationMerge string   `json:"deferred_implementation_merge,omitempty"`
}

type GateEvidenceState string

const (
	GateEvidenceRecorded                    GateEvidenceState = "recorded"
	GateEvidenceMissing                     GateEvidenceState = "missing"
	GateEvidenceExactRevisionMatch          GateEvidenceState = "exact_revision_match"
	GateEvidenceLogDescendantMatch          GateEvidenceState = "gate_log_descendant_match"
	GateEvidenceDifferentRepository         GateEvidenceState = "different_repository"
	GateEvidenceDifferentArgv               GateEvidenceState = "different_argv"
	GateEvidenceRecordedRevisionMissing     GateEvidenceState = "recorded_revision_missing"
	GateEvidenceRecordedRevisionNotAncestor GateEvidenceState = "recorded_revision_not_ancestor"
	GateEvidenceInvalidated                 GateEvidenceState = "invalidated_by_non_gate_log_commit"
	GateEvidenceInvalidRecord               GateEvidenceState = "invalid_record"
	GateEvidenceNotGreen                    GateEvidenceState = "gate_not_green"
)

type GateEvidenceResult struct {
	RepositoryIdentity string            `json:"repository_identity"`
	GateCommand        []string          `json:"gate_command"`
	GateCommandSHA256  string            `json:"gate_command_sha256"`
	RecordPath         string            `json:"record_path"`
	RecordProvenance   string            `json:"record_provenance"`
	GateExitStatus     int               `json:"gate_exit_status"`
	RecordedRevision   string            `json:"recorded_revision"`
	HeadRevision       string            `json:"head_revision"`
	TargetRevision     string            `json:"target_revision"`
	State              GateEvidenceState `json:"state"`
	Matches            bool              `json:"matches"`
	MatchBasis         string            `json:"match_basis"`
	BaselineRevision   string            `json:"baseline_revision"`
}

type FocusedTestBaselineState string

const (
	FocusedBaselineNotCompared FocusedTestBaselineState = "not_compared"
	FocusedBaselineMissing     FocusedTestBaselineState = "baseline_missing"
	FocusedBaselineUnusable    FocusedTestBaselineState = "baseline_unusable"
	FocusedBaselineMatchingRed FocusedTestBaselineState = "matching_red"
	FocusedBaselineNewRed      FocusedTestBaselineState = "new_red"
	FocusedBaselineGreen       FocusedTestBaselineState = "green"
)

// FocusedTestResult is the bounded, machine-comparable result of one owned
// focused-test probe. The diagnostic hash is computed from the same bounded
// normalized bytes carried in Diagnostic.
type FocusedTestResult struct {
	ProbeFile        string                   `json:"probe_file"`
	ExitStatus       int                      `json:"exit_status"`
	Launched         bool                     `json:"launched"`
	TimedOut         bool                     `json:"timed_out"`
	Diagnostic       string                   `json:"diagnostic"`
	DiagnosticSHA256 string                   `json:"diagnostic_sha256"`
	BaselineState    FocusedTestBaselineState `json:"baseline_state"`
	BaselineStatus   int                      `json:"baseline_exit_status"`
	BaselineLaunched bool                     `json:"baseline_launched"`
	CommandText      string                   `json:"command_text"`
}

// AlreadyGreenRepairValidation is the single machine projection consumed by
// work's TDD exception and review's no-diff exception. Both decisions are
// derived from the same observation set; review additionally requires exact
// canonical-completion staging.
type AlreadyGreenRepairValidation struct {
	RequestID                string   `json:"request_id"`
	RequestPath              string   `json:"request_path"`
	TDDAllowed               bool     `json:"tdd_allowed"`
	ReviewAllowed            bool     `json:"review_allowed"`
	IntakeFingerprint        string   `json:"intake_fingerprint"`
	ExpectedFingerprint      string   `json:"expected_fingerprint"`
	GateCommand              []string `json:"gate_command"`
	RecordedRevision         string   `json:"recorded_revision"`
	CanonicalCompletionPaths []string `json:"canonical_completion_paths"`
	StagedPaths              []string `json:"staged_paths"`
	ProjectChangedPaths      []string `json:"project_changed_paths"`
	ReasonCodes              []string `json:"reason_codes"`
	OffendingPaths           []string `json:"offending_paths"`
	Writer                   string   `json:"writer"`
	PlannedAt                string   `json:"planned_at"`
}

// FinalizationResult is the stable machine projection for one resumable REQ
// release tail. Phase names are durable journal states; callers never need to
// infer recovery progress from Git status or prose findings.
type FinalizationResult struct {
	RequestID             string   `json:"request_id"`
	RequestPath           string   `json:"request_path"`
	ArchivePath           string   `json:"archive_path"`
	JournalPath           string   `json:"journal_path"`
	Phase                 string   `json:"phase"`
	TerminalStatus        string   `json:"terminal_status"`
	Resumed               bool     `json:"resumed"`
	Discovered            bool     `json:"discovered"`
	CommitPaths           []string `json:"commit_paths"`
	PrimaryCommit         string   `json:"primary_commit"`
	MetadataCommit        string   `json:"metadata_commit"`
	CreatedPrimaryCommit  string   `json:"created_primary_commit"`
	CreatedMetadataCommit string   `json:"created_metadata_commit"`
	BlockedPaths          []string `json:"blocked_paths"`
	ReasonCodes           []string `json:"reason_codes"`
	NextArgv              []string `json:"next_argv"`
	VerificationArgv      []string `json:"verification_argv"`
	CollectionArgv        []string `json:"collection_argv"`
}

type AdvancePhaseKind string

const (
	AdvancePhaseMechanical    AdvancePhaseKind = "mechanical"
	AdvancePhaseAgentJudgment AdvancePhaseKind = "agent_judgment"
	AdvancePhaseComplete      AdvancePhaseKind = "complete"
)

type AdvanceGateState string

const (
	AdvanceGateNeedsInput AdvanceGateState = "needs_input"
	AdvanceGateSatisfied  AdvanceGateState = "satisfied"
	AdvanceGateFindings   AdvanceGateState = "findings"
	AdvanceGateFailed     AdvanceGateState = "failed"
)

type AdvanceGateProvenance string

const (
	AdvanceGateExistingEvidence AdvanceGateProvenance = "existing_evidence"
	AdvanceGateExecuted         AdvanceGateProvenance = "advance_executed"
	AdvanceGateBaselineRecord   AdvanceGateProvenance = "baseline_record"
	AdvanceGateMergedRange      AdvanceGateProvenance = "merged_range"
)

// AdvanceGateRecord binds one subordinate evidence command to the request the
// caller selected. Its collections stay typed so an action never parses a
// compatibility paragraph back into lifecycle state.
type AdvanceGateRecord struct {
	RequestID        string                `json:"request_id"`
	RequestPath      string                `json:"request_path"`
	GateID           string                `json:"gate_id"`
	Provenance       AdvanceGateProvenance `json:"provenance"`
	State            AdvanceGateState      `json:"state"`
	Outcome          CommandOutcome        `json:"outcome"`
	Findings         []CommandFinding      `json:"findings"`
	Changes          []RecordedChange      `json:"changes"`
	OutputLines      []string              `json:"output_lines"`
	NextArgv         []string              `json:"next_argv"`
	VerificationArgv []string              `json:"verification_argv"`
	FocusedTest      *FocusedTestResult    `json:"focused_test,omitempty"`
	GreenGate        *GateEvidenceResult   `json:"green_gate,omitempty"`
}

// AdvanceMissingEvidence identifies one durable file, field, or Markdown
// section that moves a request into its next lifecycle phase.
type AdvanceMissingEvidence struct {
	Kind     string `json:"kind"`
	Path     string `json:"path"`
	Field    string `json:"field,omitempty"`
	Section  string `json:"section,omitempty"`
	Expected string `json:"expected"`
}

// AdvanceLifecycleResult is the read-only per-REQ lifecycle projection. Exact
// argv stays tokenized so text and JSON consumers never parse a display string.
type AdvanceLifecycleResult struct {
	RequestID        string                   `json:"request_id"`
	RequestPath      string                   `json:"request_path"`
	TreeSection      string                   `json:"tree_section"`
	Status           string                   `json:"status"`
	Route            string                   `json:"route"`
	Phase            string                   `json:"phase"`
	PhaseKind        AdvancePhaseKind         `json:"phase_kind"`
	MissingEvidence  []AdvanceMissingEvidence `json:"missing_evidence"`
	NextArgv         []string                 `json:"next_argv"`
	VerificationArgv []string                 `json:"verification_argv"`
	GateRecords      []AdvanceGateRecord      `json:"gate_records"`
}

// QueueAdvanceMember is one member of a queue invocation's frozen target set.
// Consumed means this invocation committed the member's claim or terminal hold;
// a continuation never expands this set from a changing UR or queue.
type QueueAdvanceMember struct {
	RequestID   string `json:"request_id"`
	RequestPath string `json:"request_path"`
	Provenance  string `json:"provenance"`
	Consumed    bool   `json:"consumed"`
}

// QueueAdvancePhase records one committed or refused queue mutation. Each
// request gets its own transaction, so a later refusal cannot imply rollback
// of an earlier claim.
type QueueAdvancePhase struct {
	RequestID string           `json:"request_id"`
	Phase     string           `json:"phase"`
	Outcome   CommandOutcome   `json:"outcome"`
	Changes   []RecordedChange `json:"changes"`
	Findings  []CommandFinding `json:"findings"`
}

// QueueAdvanceResult is the typed selection-and-claim session owned by
// advance. Original non-fan-out options are retained separately from the
// dispatch bound so replay can observe, project frozen membership, then bound.
type QueueAdvanceResult struct {
	TargetTokens         []string             `json:"target_tokens"`
	WaveDepth            *int                 `json:"wave_depth,omitempty"`
	SkipImpactNegligible bool                 `json:"skip_impact_negligible"`
	SimpleOnly           bool                 `json:"simple_only"`
	DispatchBound        int                  `json:"dispatch_bound"`
	FrozenMembers        []QueueAdvanceMember `json:"frozen_members"`
	Claimed              []QueueAdvanceMember `json:"claimed"`
	Phases               []QueueAdvancePhase  `json:"phases"`
	Partial              bool                 `json:"partial"`
	ContinuationArgv     []string             `json:"continuation_argv"`
	VerificationArgv     []string             `json:"verification_argv"`
}

// RecoveryClaimResult is one structurally observed working claim and the
// authority decision made for it by the canonical recovery command.
type RecoveryClaimResult struct {
	RequestID          string                   `json:"request_id"`
	RequestPath        string                   `json:"request_path"`
	CheckpointEvidence []SelectionClaimEvidence `json:"checkpoint_evidence"`
	Decision           string                   `json:"decision"`
	Recovered          bool                     `json:"recovered"`
	// HeldForHeavyLanes marks a claim recovery deliberately preserved: the
	// request is waiting for this session's heavy-lane drain, not interrupted.
	HeldForHeavyLanes bool `json:"held_for_heavy_lanes"`
}

// RecoveryResult is the ordered public recovery composition. Finalization is
// always attempted before these claim records are classified or changed.
type RecoveryResult struct {
	AuthorityMode      string                `json:"authority_mode"`
	TakeOverRequestID  string                `json:"take_over_request_id,omitempty"`
	FinalizationPassed bool                  `json:"finalization_passed"`
	Claims             []RecoveryClaimResult `json:"claims"`
	NextArgv           []string              `json:"next_argv"`
	VerificationArgv   []string              `json:"verification_argv"`
}

// CheckpointResult describes the sole mutating advance mode.
type CheckpointResult struct {
	CheckpointPath  string `json:"checkpoint_path"`
	PreservedClaims int    `json:"preserved_claims"`
	WrittenAt       string `json:"written_at"`
}

type CommandResult struct {
	SchemaVersion        int                           `json:"schema_version"`
	Command              string                        `json:"command"`
	Outcome              CommandOutcome                `json:"outcome"`
	RepositoryRoot       string                        `json:"repository_root"`
	Findings             []CommandFinding              `json:"findings"`
	Changes              []RecordedChange              `json:"changes"`
	SkippedWork          []SkippedWork                 `json:"skipped_work"`
	AvailableCommands    []string                      `json:"available_commands,omitempty"`
	Selected             []SelectionRecord             `json:"selected"`
	Excluded             []SelectionExclusion          `json:"excluded"`
	SelectionSummary     SelectionSummary              `json:"selection_summary"`
	Rollback             RollbackResult                `json:"rollback"`
	ProtocolOutput       *string                       `json:"protocol_output,omitempty"`
	AuditMetrics         *AuditMetricsResult           `json:"audit_metrics,omitempty"`
	GateDeferral         *GateDeferralResult           `json:"gate_deferral,omitempty"`
	GateEvidence         *GateEvidenceResult           `json:"gate_evidence,omitempty"`
	FocusedTest          *FocusedTestResult            `json:"focused_test,omitempty"`
	HeavyVerification    *HeavyVerificationPlan        `json:"heavy_verification,omitempty"`
	HeavyVerificationRun *HeavyVerificationRun         `json:"heavy_verification_run,omitempty"`
	AlreadyGreenRepair   *AlreadyGreenRepairValidation `json:"already_green_repair,omitempty"`
	Finalization         *FinalizationResult           `json:"finalization,omitempty"`
	Finalizations        []FinalizationResult          `json:"finalizations"`
	Advance              *AdvanceLifecycleResult       `json:"advance,omitempty"`
	QueueAdvance         *QueueAdvanceResult           `json:"queue_advance,omitempty"`
	Recovery             *RecoveryResult               `json:"recovery,omitempty"`
	Checkpoint           *CheckpointResult             `json:"checkpoint,omitempty"`
	// ExactTextOutput preserves compatibility-shaped stdout without polluting
	// JSON with an opaque duplicate. It must be derived from the same typed
	// observation carried by the result.
	ExactTextOutput *string `json:"-"`
	// ExitCodeOverride preserves compatibility commands whose public contract is
	// an underlying process status (including cleaned-up interruptions).
	ExitCodeOverride int `json:"-"`
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
	if result.Finalizations == nil {
		result.Finalizations = []FinalizationResult{}
	}
	if result.HeavyVerification != nil {
		if result.HeavyVerification.SourceRanges == nil && result.HeavyVerification.Mode != "" {
			result.HeavyVerification.SourceRanges = []HeavySourceRange{}
		}
		if result.HeavyVerification.ChangedPaths == nil {
			result.HeavyVerification.ChangedPaths = []string{}
		}
		if result.HeavyVerification.UncoveredPaths == nil {
			result.HeavyVerification.UncoveredPaths = []string{}
		}
		if result.HeavyVerification.SelectedLanes == nil {
			result.HeavyVerification.SelectedLanes = []HeavyLaneSelection{}
		}
		for laneIndex := range result.HeavyVerification.SelectedLanes {
			lane := &result.HeavyVerification.SelectedLanes[laneIndex]
			if lane.CommandArgv == nil {
				lane.CommandArgv = []string{}
			}
			if lane.Reasons == nil {
				lane.Reasons = []string{}
			}
		}
	}
	if result.HeavyVerificationRun != nil {
		if result.HeavyVerificationRun.Lanes == nil {
			result.HeavyVerificationRun.Lanes = []HeavyLaneExecution{}
		}
		for laneIndex := range result.HeavyVerificationRun.Lanes {
			lane := &result.HeavyVerificationRun.Lanes[laneIndex]
			if lane.CommandArgv == nil {
				lane.CommandArgv = []string{}
			}
			if lane.Disposition == "" {
				// A lane with no stated disposition was executed: reuse is the
				// case that has to be declared, never the one inferred.
				lane.Disposition = "executed"
			}
		}
	}
	if len(result.Finalizations) == 0 && result.Finalization != nil {
		result.Finalizations = append(result.Finalizations, *result.Finalization)
	}
	if len(result.Finalizations) == 1 {
		record := result.Finalizations[0]
		result.Finalization = &record
	} else if len(result.Finalizations) > 1 {
		result.Finalization = nil
	}
	for index := range result.Finalizations {
		normalizeFinalization(&result.Finalizations[index])
	}
	if result.Finalization != nil {
		normalizeFinalization(result.Finalization)
	}
	for index := range result.Selected {
		selection := &result.Selected[index]
		if selection.SelectionPriority == "" {
			selection.SelectionPriority = "ordinary"
		}
		if selection.RequestPriority == "" {
			selection.RequestPriority = "next"
		}
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
		if exclusion.SelectionPriority == "" {
			exclusion.SelectionPriority = "ordinary"
		}
		if exclusion.RequestPriority == "" {
			exclusion.RequestPriority = "next"
		}
		if exclusion.ProbeStatus == "" {
			exclusion.ProbeStatus = ProbeNotApplicable
		}
		if !exclusion.ProbeAttempted {
			exclusion.ProbeExitCode = -1
		}
		if exclusion.NextArgv == nil {
			exclusion.NextArgv = []string{}
		}
		if exclusion.ClaimEvidence == nil {
			exclusion.ClaimEvidence = []SelectionClaimEvidence{}
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
	if result.AuditMetrics != nil {
		metrics := result.AuditMetrics
		if metrics.Aggregates == nil {
			metrics.Aggregates = []AuditAggregate{}
		}
		if metrics.Distributions == nil {
			metrics.Distributions = []AuditDistribution{}
		}
		if metrics.Files == nil {
			metrics.Files = []AuditFileMeasurement{}
		}
		if metrics.Folders == nil {
			metrics.Folders = []AuditFolderMeasurement{}
		}
		if metrics.Bands == nil {
			metrics.Bands = []AuditBand{}
		}
		if metrics.Churn == nil {
			metrics.Churn = []AuditChurnMeasurement{}
		}
		if metrics.Hotspots == nil {
			metrics.Hotspots = []AuditHotspot{}
		}
		if metrics.UnavailablePaths == nil {
			metrics.UnavailablePaths = []string{}
		}
	}
	if result.GateDeferral != nil && result.GateDeferral.GateCommand == nil {
		result.GateDeferral.GateCommand = []string{}
	}
	if result.GateEvidence != nil && result.GateEvidence.GateCommand == nil {
		result.GateEvidence.GateCommand = []string{}
	}
	if result.AlreadyGreenRepair != nil {
		validation := result.AlreadyGreenRepair
		if validation.GateCommand == nil {
			validation.GateCommand = []string{}
		}
		if validation.CanonicalCompletionPaths == nil {
			validation.CanonicalCompletionPaths = []string{}
		}
		if validation.StagedPaths == nil {
			validation.StagedPaths = []string{}
		}
		if validation.ProjectChangedPaths == nil {
			validation.ProjectChangedPaths = []string{}
		}
		if validation.ReasonCodes == nil {
			validation.ReasonCodes = []string{}
		}
		if validation.OffendingPaths == nil {
			validation.OffendingPaths = []string{}
		}
	}
	if result.Advance != nil {
		if result.Advance.MissingEvidence == nil {
			result.Advance.MissingEvidence = []AdvanceMissingEvidence{}
		}
		if result.Advance.NextArgv == nil {
			result.Advance.NextArgv = []string{}
		}
		if result.Advance.VerificationArgv == nil {
			result.Advance.VerificationArgv = []string{}
		}
		if result.Advance.GateRecords == nil {
			result.Advance.GateRecords = []AdvanceGateRecord{}
		}
		for gateIndex := range result.Advance.GateRecords {
			gate := &result.Advance.GateRecords[gateIndex]
			if gate.Findings == nil {
				gate.Findings = []CommandFinding{}
			}
			if gate.Changes == nil {
				gate.Changes = []RecordedChange{}
			}
			if gate.OutputLines == nil {
				gate.OutputLines = []string{}
			}
			if gate.NextArgv == nil {
				gate.NextArgv = []string{}
			}
			if gate.VerificationArgv == nil {
				gate.VerificationArgv = []string{}
			}
		}
	}
	if result.QueueAdvance != nil {
		queueAdvance := result.QueueAdvance
		if queueAdvance.TargetTokens == nil {
			queueAdvance.TargetTokens = []string{}
		}
		if queueAdvance.FrozenMembers == nil {
			queueAdvance.FrozenMembers = []QueueAdvanceMember{}
		}
		if queueAdvance.Claimed == nil {
			queueAdvance.Claimed = []QueueAdvanceMember{}
		}
		if queueAdvance.Phases == nil {
			queueAdvance.Phases = []QueueAdvancePhase{}
		}
		for phaseIndex := range queueAdvance.Phases {
			if queueAdvance.Phases[phaseIndex].Changes == nil {
				queueAdvance.Phases[phaseIndex].Changes = []RecordedChange{}
			}
			if queueAdvance.Phases[phaseIndex].Findings == nil {
				queueAdvance.Phases[phaseIndex].Findings = []CommandFinding{}
			}
		}
		if queueAdvance.ContinuationArgv == nil {
			queueAdvance.ContinuationArgv = []string{}
		}
		if queueAdvance.VerificationArgv == nil {
			queueAdvance.VerificationArgv = []string{}
		}
	}
	if result.Recovery != nil {
		if result.Recovery.Claims == nil {
			result.Recovery.Claims = []RecoveryClaimResult{}
		}
		for claimIndex := range result.Recovery.Claims {
			if result.Recovery.Claims[claimIndex].CheckpointEvidence == nil {
				result.Recovery.Claims[claimIndex].CheckpointEvidence = []SelectionClaimEvidence{}
			}
		}
		if result.Recovery.NextArgv == nil {
			result.Recovery.NextArgv = []string{}
		}
		if result.Recovery.VerificationArgv == nil {
			result.Recovery.VerificationArgv = []string{}
		}
	}
	setAsideSelfRefusal := false
	allRefusalBlockersOwned := true
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
		if finding.Fixability == FixabilityRefused && len(finding.AffectedIDs) == 0 {
			allRefusalBlockersOwned = false
		}
		isRefusalFinding := result.Outcome == OutcomeRefused || finding.Fixability == FixabilityRefused
		if isRefusalFinding && nextCommandVerb(finding.NextArgv) == result.Command {
			finding.NextArgv = []string{}
			finding.NextJustRecipe = ""
			if finding.AutomationStopReason == "" {
				finding.AutomationStopReason = "no distinct resolving command is available"
			}
			if len(finding.AffectedIDs) > 0 {
				setAsideSelfRefusal = true
			}
		}
	}
	if setAsideSelfRefusal && allRefusalBlockersOwned {
		result.Outcome = OutcomeFindings
	}
	return result
}

func nextCommandVerb(arguments []string) string {
	if len(arguments) == 0 {
		return ""
	}
	executable := strings.TrimSuffix(filepath.Base(arguments[0]), ".sh")
	if executable != "do-work-cli" {
		return executable
	}
	for index := 1; index < len(arguments); index++ {
		argument := arguments[index]
		if argument == "--format" || argument == "--repo-root" {
			index++
			continue
		}
		if strings.HasPrefix(argument, "--format=") || strings.HasPrefix(argument, "--repo-root=") || strings.HasPrefix(argument, "-") {
			continue
		}
		return argument
	}
	return ""
}

func normalizeFinalization(record *FinalizationResult) {
	if record.CommitPaths == nil {
		record.CommitPaths = []string{}
	}
	if record.BlockedPaths == nil {
		record.BlockedPaths = []string{}
	}
	if record.ReasonCodes == nil {
		record.ReasonCodes = []string{}
	}
	if record.NextArgv == nil {
		record.NextArgv = []string{}
	}
	if record.VerificationArgv == nil {
		record.VerificationArgv = []string{}
	}
	if record.CollectionArgv == nil {
		record.CollectionArgv = []string{}
	}
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
	if result.ProtocolOutput != nil {
		return []byte(*result.ProtocolOutput)
	}
	if result.ExactTextOutput != nil {
		return []byte(*result.ExactTextOutput)
	}
	var output strings.Builder
	fmt.Fprintf(&output, "%s: %s\n", result.Command, result.Outcome)
	fmt.Fprintf(&output, "repository: %s\n", result.RepositoryRoot)
	if len(result.AvailableCommands) > 0 {
		fmt.Fprintf(&output, "available commands: %s\n", strings.Join(result.AvailableCommands, ", "))
	}
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
	if result.Advance != nil {
		advance := result.Advance
		fmt.Fprintf(&output, "advance %s [%s, %s, route %s]: %s\n", advance.RequestID, advance.TreeSection, advance.Status, advance.Route, advance.Phase)
		fmt.Fprintf(&output, "  request: %s\n", advance.RequestPath)
		fmt.Fprintf(&output, "  phase kind: %s\n", advance.PhaseKind)
		for _, evidence := range advance.MissingEvidence {
			coordinate := evidence.Path
			if evidence.Field != "" {
				coordinate += " field=" + evidence.Field
			}
			if evidence.Section != "" {
				coordinate += " section=" + strconv.Quote(evidence.Section)
			}
			fmt.Fprintf(&output, "  missing %s: %s expected=%s\n", evidence.Kind, coordinate, evidence.Expected)
		}
		if len(advance.NextArgv) > 0 {
			fmt.Fprintf(&output, "  next: %s\n", joinArgv(advance.NextArgv))
		}
		if len(advance.VerificationArgv) > 0 {
			fmt.Fprintf(&output, "  verify: %s\n", joinArgv(advance.VerificationArgv))
		}
		for _, gate := range advance.GateRecords {
			fmt.Fprintf(&output, "  gate %s [%s, %s]: %s\n", gate.GateID, gate.Provenance, gate.State, gate.Outcome)
			for _, line := range gate.OutputLines {
				fmt.Fprintf(&output, "    evidence: %s\n", line)
			}
			if gate.FocusedTest != nil {
				fmt.Fprintf(&output, "    focused test: status=%d baseline=%s diagnostic=%s\n", gate.FocusedTest.ExitStatus, gate.FocusedTest.BaselineState, gate.FocusedTest.DiagnosticSHA256)
			}
			if gate.GreenGate != nil {
				fmt.Fprintf(&output, "    green gate: state=%s matches=%t\n", gate.GreenGate.State, gate.GreenGate.Matches)
			}
			if len(gate.NextArgv) > 0 {
				fmt.Fprintf(&output, "    next: %s\n", joinArgv(gate.NextArgv))
			}
		}
	}
	for _, record := range result.Finalizations {
		fmt.Fprintf(&output, "finalization %s [%s, %s]\n", record.RequestID, record.Phase, record.TerminalStatus)
		fmt.Fprintf(&output, "  request: %s\n", record.RequestPath)
		fmt.Fprintf(&output, "  archive: %s\n", record.ArchivePath)
		fmt.Fprintf(&output, "  journal: %s\n", record.JournalPath)
		fmt.Fprintf(&output, "  resumed: %t\n", record.Resumed)
		fmt.Fprintf(&output, "  discovered: %t\n", record.Discovered)
		fmt.Fprintf(&output, "  commit paths: %s\n", strings.Join(record.CommitPaths, ", "))
		fmt.Fprintf(&output, "  primary commit: %s\n", record.PrimaryCommit)
		fmt.Fprintf(&output, "  metadata commit: %s\n", record.MetadataCommit)
		fmt.Fprintf(&output, "  created primary commit: %s\n", record.CreatedPrimaryCommit)
		fmt.Fprintf(&output, "  created metadata commit: %s\n", record.CreatedMetadataCommit)
		fmt.Fprintf(&output, "  blocked paths: %s\n", strings.Join(record.BlockedPaths, ", "))
		fmt.Fprintf(&output, "  reason codes: %s\n", strings.Join(record.ReasonCodes, ", "))
		if len(record.NextArgv) > 0 {
			fmt.Fprintf(&output, "  next: %s\n", joinArgv(record.NextArgv))
		}
		if len(record.VerificationArgv) > 0 {
			fmt.Fprintf(&output, "  verify: %s\n", joinArgv(record.VerificationArgv))
		}
		if len(record.CollectionArgv) > 0 {
			fmt.Fprintf(&output, "  collect: %s\n", joinArgv(record.CollectionArgv))
		}
	}
	if result.QueueAdvance != nil {
		queueAdvance := result.QueueAdvance
		fmt.Fprintf(&output, "queue advance: members=%d claimed=%d bound=%d partial=%t\n", len(queueAdvance.FrozenMembers), len(queueAdvance.Claimed), queueAdvance.DispatchBound, queueAdvance.Partial)
		for _, phase := range queueAdvance.Phases {
			fmt.Fprintf(&output, "  %s %s: %s\n", phase.RequestID, phase.Phase, phase.Outcome)
		}
		if len(queueAdvance.ContinuationArgv) > 0 {
			fmt.Fprintf(&output, "  continue: %s\n", joinArgv(queueAdvance.ContinuationArgv))
		}
		if len(queueAdvance.VerificationArgv) > 0 {
			fmt.Fprintf(&output, "  verify: %s\n", joinArgv(queueAdvance.VerificationArgv))
		}
	}
	if result.Recovery != nil {
		recovery := result.Recovery
		fmt.Fprintf(&output, "recovery [%s]: finalization_passed=%t claims=%d\n", recovery.AuthorityMode, recovery.FinalizationPassed, len(recovery.Claims))
		for _, claim := range recovery.Claims {
			fmt.Fprintf(&output, "  claim %s [%s]: %s\n", claim.RequestID, claim.RequestPath, claim.Decision)
		}
		if len(recovery.NextArgv) > 0 {
			fmt.Fprintf(&output, "  next: %s\n", joinArgv(recovery.NextArgv))
		}
		if len(recovery.VerificationArgv) > 0 {
			fmt.Fprintf(&output, "  verify: %s\n", joinArgv(recovery.VerificationArgv))
		}
	}
	if result.Checkpoint != nil {
		fmt.Fprintf(&output, "checkpoint %s: preserved_claims=%d written_at=%s\n", result.Checkpoint.CheckpointPath, result.Checkpoint.PreservedClaims, result.Checkpoint.WrittenAt)
	}
	for _, change := range result.Changes {
		fmt.Fprintf(&output, "change %s [%s]: %s\n", change.Path, change.Kind, change.Detail)
	}
	if result.GateDeferral != nil {
		gate := result.GateDeferral
		fmt.Fprintf(&output, "gate deferral: parent=%s repair=%s outcome=%s fingerprint=%s\n", gate.ParentID, gate.RepairID, gate.RepairOutcome, gate.DiagnosticFingerprint)
		fmt.Fprintf(&output, "  parent path: %s\n", gate.ParentPath)
		fmt.Fprintf(&output, "  repair path: %s\n", gate.RepairPath)
		fmt.Fprintf(&output, "  checkpoint path: %s\n", gate.CheckpointPath)
		fmt.Fprintf(&output, "  repair dependency: %s\n", gate.RepairDependency)
		fmt.Fprintf(&output, "  sweep key: %s\n", gate.SweepKey)
		fmt.Fprintf(&output, "  gate command: %s (exit: %d)\n", joinArgv(gate.GateCommand), gate.GateExitStatus)
		if gate.DeferredImplementationBase != "" {
			fmt.Fprintf(&output, "  implementation range: %s..%s\n", gate.DeferredImplementationBase, gate.DeferredImplementationMerge)
		}
	}
	if result.GateEvidence != nil {
		gate := result.GateEvidence
		fmt.Fprintf(&output, "gate evidence: state=%s matches=%t basis=%s\n", gate.State, gate.Matches, gate.MatchBasis)
		fmt.Fprintf(&output, "  repository identity: %s\n", gate.RepositoryIdentity)
		fmt.Fprintf(&output, "  gate command: %s (sha256: %s, exit: %d)\n", joinArgv(gate.GateCommand), gate.GateCommandSHA256, gate.GateExitStatus)
		fmt.Fprintf(&output, "  record: %s (provenance: %s)\n", gate.RecordPath, gate.RecordProvenance)
		fmt.Fprintf(&output, "  revisions: recorded=%s head=%s baseline=%s target=%s\n", gate.RecordedRevision, gate.HeadRevision, gate.BaselineRevision, gate.TargetRevision)
	}
	if result.HeavyVerification != nil {
		plan := result.HeavyVerification
		fmt.Fprintf(&output, "heavy verification: selected=%d uncertain=%t force_all=%t\n", len(plan.SelectedLanes), plan.Uncertain, plan.ForcedAll)
		fmt.Fprintf(&output, "  manifest: %s\n", plan.ManifestPath)
		if plan.Mode == "" {
			fmt.Fprintf(&output, "  revisions: %s..%s\n", plan.BaseRevision, plan.TargetRevision)
		} else {
			fmt.Fprintf(&output, "  mode: %s\n", plan.Mode)
			fmt.Fprintf(&output, "  manifest revision: %s\n", plan.ManifestRevision)
			fmt.Fprintf(&output, "  execution revision: %s\n", plan.ExecutionRevision)
			for _, sourceRange := range plan.SourceRanges {
				fmt.Fprintf(&output, "  source range: %s..%s\n", sourceRange.BaseRevision, sourceRange.TargetRevision)
			}
		}
		fmt.Fprintf(&output, "  changed paths: %s\n", strings.Join(plan.ChangedPaths, ", "))
		if len(plan.UncoveredPaths) > 0 {
			fmt.Fprintf(&output, "  uncovered paths: %s\n", strings.Join(plan.UncoveredPaths, ", "))
		}
		for _, lane := range plan.SelectedLanes {
			fmt.Fprintf(&output, "  lane %s: %s\n", lane.LaneID, joinArgv(lane.CommandArgv))
			for _, reason := range lane.Reasons {
				fmt.Fprintf(&output, "    reason: %s\n", reason)
			}
		}
	}
	if result.HeavyVerificationRun != nil {
		run := result.HeavyVerificationRun
		fmt.Fprintf(&output, "heavy verification run: lanes=%d\n", len(run.Lanes))
		fmt.Fprintf(&output, "  manifest: %s\n", run.ManifestPath)
		fmt.Fprintf(&output, "  execution revision: %s\n", run.ExecutionRevision)
		for _, lane := range run.Lanes {
			laneOutcome := fmt.Sprintf("exit %d", lane.ExitStatus)
			if lane.Skipped {
				laneOutcome = "skipped"
			}
			laneDisposition := lane.Disposition
			if lane.DispositionReason != "" {
				laneDisposition = laneDisposition + ": " + lane.DispositionReason
			}
			fmt.Fprintf(&output, "  lane %s: %s in %ds [%s] — %s\n", lane.LaneID, laneOutcome, lane.WallSeconds, laneDisposition, joinArgv(lane.CommandArgv))
			if lane.EvidenceRevision != "" || lane.EvidenceRecordedAt != "" {
				fmt.Fprintf(&output, "    inherited evidence: revision=%s recorded_at=%s\n", lane.EvidenceRevision, lane.EvidenceRecordedAt)
			}
		}
	}
	if result.AlreadyGreenRepair != nil {
		validation := result.AlreadyGreenRepair
		fmt.Fprintf(&output, "already-green repair: request=%s tdd_allowed=%t review_allowed=%t\n", validation.RequestID, validation.TDDAllowed, validation.ReviewAllowed)
		fmt.Fprintf(&output, "  request path: %s\n", validation.RequestPath)
		fmt.Fprintf(&output, "  fingerprints: intake=%s expected=%s\n", validation.IntakeFingerprint, validation.ExpectedFingerprint)
		fmt.Fprintf(&output, "  gate command: %s\n", joinArgv(validation.GateCommand))
		fmt.Fprintf(&output, "  recorded revision: %s\n", validation.RecordedRevision)
		fmt.Fprintf(&output, "  completion paths: %s\n", strings.Join(validation.CanonicalCompletionPaths, ", "))
		fmt.Fprintf(&output, "  staged paths: %s\n", strings.Join(validation.StagedPaths, ", "))
		fmt.Fprintf(&output, "  project changed paths: %s\n", strings.Join(validation.ProjectChangedPaths, ", "))
		fmt.Fprintf(&output, "  reason codes: %s\n", strings.Join(validation.ReasonCodes, ", "))
		fmt.Fprintf(&output, "  offending paths: %s\n", strings.Join(validation.OffendingPaths, ", "))
		fmt.Fprintf(&output, "  completion authority: writer=%s at=%s\n", validation.Writer, validation.PlannedAt)
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
		fmt.Fprintf(&output, "selected %s [%s, %s, depth %d, %s]: %s\n", selection.RequestID, selection.Provenance, selection.SelectionPriority, selection.DependencyDepth, estimate, selection.Title)
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
		fmt.Fprintf(&output, "excluded %s [%s, %s] %s: %s\n", exclusion.RequestID, exclusion.Provenance, exclusion.SelectionPriority, exclusion.Code, exclusion.Reason)
		fmt.Fprintf(&output, "  request: %s (original status: %s)\n", exclusion.RequestPath, exclusion.OriginalStatus)
		fmt.Fprintf(&output, "  probe %s (attempted: %t, exit: %d)", exclusion.ProbeStatus, exclusion.ProbeAttempted, exclusion.ProbeExitCode)
		if exclusion.UnblockRequired {
			fmt.Fprint(&output, "; unblock required")
		}
		fmt.Fprintln(&output)
		for _, claim := range exclusion.ClaimEvidence {
			fmt.Fprintf(&output, "  claim %s: claimed_at=%s writer=%s path=%s line=%d\n", claim.Source, claim.ClaimedAt, claim.Writer, claim.Path, claim.SourceLine)
			if claim.HeaderText != "" {
				fmt.Fprintf(&output, "    header: %s\n", claim.HeaderText)
			}
		}
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
