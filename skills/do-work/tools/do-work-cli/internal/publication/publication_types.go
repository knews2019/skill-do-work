// Package publication owns deterministic capture, answer, and release mutations.
package publication

import (
	"io/fs"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type OperationName string

const (
	OperationCaptureFiles OperationName = "capture-files"
	OperationAnswer       OperationName = "answer"
	OperationRelease      OperationName = "release"
	OperationDeferGate    OperationName = "defer-gate"
)

type Manifest struct {
	Operation     OperationName      `json:"operation"`
	CommitMessage string             `json:"commit_message,omitempty"`
	Capture       *CaptureManifest   `json:"capture,omitempty"`
	Answer        *AnswerManifest    `json:"answer,omitempty"`
	Release       *ReleaseManifest   `json:"release,omitempty"`
	DeferGate     *DeferGateManifest `json:"defer_gate,omitempty"`
}

type PayloadFile struct {
	SourcePath string `json:"source_path"`
	SHA256     string `json:"sha256,omitempty"`
}

type PublishedFile struct {
	Path    string      `json:"path"`
	Payload PayloadFile `json:"payload"`
	Mode    uint32      `json:"mode,omitempty"`
}

type ReplacementFile struct {
	Path            string      `json:"path"`
	ExpectedPayload PayloadFile `json:"expected_payload"`
	NewPayload      PayloadFile `json:"new_payload"`
	AllowUntracked  bool        `json:"allow_untracked,omitempty"`
}

type CaptureRequest struct {
	ID              string        `json:"id"`
	UserRequestID   string        `json:"user_request_id"`
	File            PublishedFile `json:"file"`
	ReservationPath string        `json:"reservation_path"`
}

type CaptureManifest struct {
	UserRequestID string            `json:"user_request_id"`
	UserRequest   PublishedFile     `json:"user_request"`
	RawInput      *PayloadFile      `json:"raw_input,omitempty"`
	Requests      []CaptureRequest  `json:"requests"`
	Assets        []PublishedFile   `json:"assets,omitempty"`
	Folds         []ReplacementFile `json:"folds,omitempty"`
}

type QuestionAnswer struct {
	QuestionID     string       `json:"question_id,omitempty"`
	ExpectedLine   string       `json:"expected_line,omitempty"`
	InsertQuestion bool         `json:"insert_question,omitempty"`
	Outcome        string       `json:"outcome"`
	Summary        string       `json:"summary"`
	RawAnswer      *PayloadFile `json:"raw_answer,omitempty"`
}

type StakeholderReportEvidence struct {
	BlockedBy      string      `json:"blocked_by"`
	ReportsHistory PayloadFile `json:"reports_history"`
}

type StakeholderTerminalEvidence struct {
	BlockedHistory PayloadFile `json:"blocked_history"`
	Implementation PayloadFile `json:"implementation"`
}

type AnswerManifest struct {
	RequestPath         string                       `json:"request_path"`
	ExpectedStatus      string                       `json:"expected_status"`
	Mode                string                       `json:"mode"`
	Answers             []QuestionAnswer             `json:"answers"`
	Report              *PublishedFile               `json:"report,omitempty"`
	StakeholderReport   *StakeholderReportEvidence   `json:"stakeholder_report,omitempty"`
	StakeholderTerminal *StakeholderTerminalEvidence `json:"stakeholder_terminal,omitempty"`
	OverrideCapture     *CaptureManifest             `json:"override_capture,omitempty"`
	// OverrideCreates and OverrideFolds are decoded only to produce a stable
	// refusal for obsolete unstructured manifests. New callers use OverrideCapture.
	OverrideCreates  []PublishedFile   `json:"override_creates,omitempty"`
	OverrideFolds    []ReplacementFile `json:"override_folds,omitempty"`
	CloseUserRequest bool              `json:"close_user_request,omitempty"`
	ArchivePath      string            `json:"archive_path,omitempty"`
	UserRequestPath  string            `json:"user_request_path,omitempty"`
	ArchiveDirectory string            `json:"archive_directory,omitempty"`
}

type ReleaseTarget struct {
	Path            string      `json:"path"`
	ExpectedPayload PayloadFile `json:"expected_payload"`
	NewPayload      PayloadFile `json:"new_payload"`
	OldVersion      string      `json:"old_version"`
	NewVersion      string      `json:"new_version"`
}

type ChangelogTarget struct {
	Path            string      `json:"path"`
	Create          bool        `json:"create,omitempty"`
	ExpectedPayload PayloadFile `json:"expected_payload"`
	NewPayload      PayloadFile `json:"new_payload"`
	InsertionAnchor string      `json:"insertion_anchor"`
	EntryKey        string      `json:"entry_key"`
	EntryTitle      string      `json:"entry_title"`
}

type ReleaseManifest struct {
	MaintainerRelease   bool              `json:"maintainer_release,omitempty"`
	OldVersion          string            `json:"old_version,omitempty"`
	NewVersion          string            `json:"new_version,omitempty"`
	ProjectOwnedTargets []string          `json:"project_owned_targets,omitempty"`
	RequiredMirrors     []string          `json:"required_mirrors,omitempty"`
	Targets             []ReleaseTarget   `json:"targets"`
	Changelogs          []ChangelogTarget `json:"changelogs"`
}

type DeferGateManifest struct {
	ParentID                    string       `json:"parent_id"`
	ParentPath                  string       `json:"parent_path"`
	ExpectedParent              PayloadFile  `json:"expected_parent"`
	ExpectedStatus              string       `json:"expected_status"`
	CheckpointPath              string       `json:"checkpoint_path"`
	ExpectedCheckpoint          PayloadFile  `json:"expected_checkpoint"`
	WriterLabel                 string       `json:"writer_label"`
	GateCommand                 []string     `json:"gate_command"`
	GateExitStatus              int          `json:"gate_exit_status"`
	DiagnosticFingerprint       string       `json:"diagnostic_fingerprint"`
	DiagnosticEvidence          []string     `json:"diagnostic_evidence"`
	SweepKey                    string       `json:"sweep_key"`
	RepairID                    string       `json:"repair_id"`
	RepairPath                  string       `json:"repair_path"`
	RepairTitle                 string       `json:"repair_title"`
	RepairCreatedAt             string       `json:"repair_created_at"`
	ReservationPath             string       `json:"reservation_path"`
	ExpectedRepair              *PayloadFile `json:"expected_repair,omitempty"`
	DeferredImplementationBase  string       `json:"deferred_implementation_base,omitempty"`
	DeferredImplementationMerge string       `json:"deferred_implementation_merge,omitempty"`
}

type MutationKind string

const (
	MutationCreate  MutationKind = "create"
	MutationReplace MutationKind = "replace"
	MutationMove    MutationKind = "move"
)

type PlannedMutation struct {
	Kind            MutationKind
	Path            string
	DestinationPath string
	ExpectedBytes   []byte
	Contents        []byte
	Mode            fs.FileMode
	AllowUntracked  bool
}

type Refusal struct {
	Code   string
	Reason string
	IDs    []string
	Paths  []string
}

type PublicationPlan struct {
	Operation                    OperationName
	RepositoryRoot               string
	ManifestPath                 string
	AnswerAt                     string
	CommitMessage                string
	Mutations                    []PlannedMutation
	TargetPaths                  []string
	ExistingUntrackedTargetPaths []string
	ExistingDirtyTargetPaths     []string
	CreatedDirectoryPaths        []string
	Changes                      []resultmodel.RecordedChange
	Refusal                      *Refusal
	GateDeferral                 *resultmodel.GateDeferralResult
}

func (plan PublicationPlan) Runnable() bool {
	return plan.Refusal == nil && plan.RepositoryRoot != "" && len(plan.Mutations) > 0
}
