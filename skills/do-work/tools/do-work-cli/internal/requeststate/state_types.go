package requeststate

import (
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type Transition string

const (
	TransitionClaim    Transition = "claim"
	TransitionRecover  Transition = "recover-claim"
	TransitionUnblock  Transition = "unblock"
	TransitionComplete Transition = "complete"
	TransitionFail     Transition = "fail"
	TransitionCancel   Transition = "cancel"
)

type SelectionProvenance string

const (
	ProvenanceDefault    SelectionProvenance = "default"
	ProvenanceExplicit   SelectionProvenance = "explicit-req"
	ProvenanceURExpanded SelectionProvenance = "ur-expanded"
)

type UnblockSource string

const (
	UnblockProbe   UnblockSource = "successful-probe"
	UnblockClarify UnblockSource = "user-via-clarify"
)

type StateOptions struct {
	Transition            Transition
	RequestID             string
	RequestPath           string
	Provenance            SelectionProvenance
	OriginalStatus        string
	ProbeStatus           resultmodel.SelectionProbeStatus
	UnblockRequired       bool
	UnblockSource         UnblockSource
	TerminalStatus        string
	ImplementationHash    string
	FailureError          string
	FailureType           string
	CancellationReason    string
	CancellationSummary   string
	CancellationConfirmed bool
	DependentDisposition  string
	WriterLabel           string
	CheckpointWriter      string
	CheckpointUnlabeled   bool
	CheckpointAbsent      bool
	CheckpointAllEntries  bool
	AssumeSoleWriter      bool
	Now                   time.Time
	DryRun                bool
	Commit                bool
	RecordCommitHashOnly  bool
	// AcceptedPreimageDigests maps an exact target path to the hex SHA-256 of
	// dirty tracked bytes a journal recorded as its own preimage. A dirty target
	// whose current bytes hash to its entry is transaction input rather than a
	// refusal; every other dirty target keeps the default refusal.
	AcceptedPreimageDigests map[string]string
}

type StateRefusal struct {
	Code   string
	Reason string
	Paths  []string
}

type FileMove struct {
	SourcePath      string
	DestinationPath string
	ExpectedBytes   []byte
}

type StatePlan struct {
	RepositoryRoot               string
	Transition                   Transition
	Options                      StateOptions
	Target                       *repositorymodel.RequestFile
	SourcePath                   string
	DestinationPath              string
	ExpectedTargetBytes          []byte
	TargetPaths                  []string
	ExistingUntrackedTargetPaths []string
	ExistingDirtyTargetPaths     []string
	CreatedDirectoryPaths        []string
	AdditionalMoves              []FileMove
	CheckpointPath               string
	CheckpointBytes              []byte
	CheckpointExisted            bool
	CalibrationPath              string
	CalibrationBytes             []byte
	CalibrationExisted           bool
	Changes                      []resultmodel.RecordedChange
	SkippedWork                  []resultmodel.SkippedWork
	Refusal                      *StateRefusal
}

func (plan StatePlan) Runnable() bool { return plan.Refusal == nil && plan.Target != nil }

// PlannedFileImage is an exact filesystem postcondition emitted by the
// lifecycle planner for a higher-level resumable transaction. A missing image
// is represented by Exists=false; Bytes are otherwise the complete file.
type PlannedFileImage struct {
	Path   string
	Exists bool
	Bytes  []byte
	Mode   uint32
}
