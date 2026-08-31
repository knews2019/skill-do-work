package requeststate

import (
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/repositorymodel"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

type Transition string

const (
	TransitionClaim    Transition = "claim"
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
	Now                   time.Time
	DryRun                bool
	Commit                bool
	RecordCommitHashOnly  bool
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
