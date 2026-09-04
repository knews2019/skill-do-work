// Package finalization owns resumable REQ archive, release, commit, and
// provenance tails. Lifecycle and release policy remain with their existing
// planners; this package journals and composes their exact results.
package finalization

import (
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/publication"
)

const journalVersion = 1

const (
	ProvenancePrimaryCommit  = "primary_commit"
	ProvenanceSuppliedCommit = "supplied_commit"
)

type Phase string

const (
	PhasePrepared          Phase = "prepared"
	PhaseLifecycleApplied  Phase = "lifecycle_applied"
	PhaseReleaseApplied    Phase = "release_applied"
	PhasePrimaryCommitted  Phase = "primary_committed"
	PhaseMetadataCommitted Phase = "metadata_committed"
	PhaseVerified          Phase = "verified"
	PhaseCleanupComplete   Phase = "cleanup_complete"
	PhaseDiscoveryRefused  Phase = "discovery_refused"
)

type Manifest struct {
	RequestID                string   `json:"request_id"`
	RequestPath              string   `json:"request_path"`
	WriterLabel              string   `json:"writer_label"`
	Transition               string   `json:"transition"`
	TerminalStatus           string   `json:"terminal_status,omitempty"`
	FailureError             string   `json:"failure_error,omitempty"`
	FailureType              string   `json:"failure_type,omitempty"`
	CompletedAt              string   `json:"completed_at"`
	ExpectedRequestSHA256    string   `json:"expected_request_sha256"`
	ExpectedCheckpointSHA256 string   `json:"expected_checkpoint_sha256"`
	CommitPaths              []string `json:"commit_paths"`
	CommitMessage            string   `json:"commit_message"`
	ProvenanceMode           string   `json:"provenance_mode"`
	ImplementationHash       string   `json:"implementation_hash,omitempty"`
	ReleaseManifestPath      string   `json:"release_manifest_path,omitempty"`
	ReleaseAt                string   `json:"release_at,omitempty"`
}

type FileImage struct {
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
	Bytes  []byte `json:"bytes,omitempty"`
	Mode   uint32 `json:"mode,omitempty"`
}

type Journal struct {
	Version                int                   `json:"version"`
	CreatedAt              time.Time             `json:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at"`
	Phase                  Phase                 `json:"phase"`
	Manifest               Manifest              `json:"manifest"`
	ManifestSHA256         string                `json:"manifest_sha256"`
	ImageSetSHA256         string                `json:"image_set_sha256"`
	JournalPath            string                `json:"journal_path"`
	PayloadDirectory       string                `json:"payload_directory,omitempty"`
	ArchivedPath           string                `json:"archived_path"`
	LifecyclePreimages     []FileImage           `json:"lifecycle_preimages"`
	LifecyclePostimages    []FileImage           `json:"lifecycle_postimages"`
	ReleaseManifest        *publication.Manifest `json:"release_manifest,omitempty"`
	ReleasePreimages       []FileImage           `json:"release_preimages"`
	ReleasePostimages      []FileImage           `json:"release_postimages"`
	EffectiveCommitPaths   []string              `json:"effective_commit_paths"`
	PreparedHead           string                `json:"prepared_head"`
	PreparedDiffSHA256     string                `json:"prepared_diff_sha256"`
	PrimaryCommit          string                `json:"primary_commit,omitempty"`
	MetadataCommit         string                `json:"metadata_commit,omitempty"`
	Discovered             bool                  `json:"discovered,omitempty"`
	SoleReleaserAttributed []string              `json:"sole_releaser_attributed,omitempty"`
	CreatedPrimaryCommit   string                `json:"-"`
	CreatedMetadataCommit  string                `json:"-"`
}

type commandOptions struct {
	ManifestPath            string
	Discover                bool
	AssumeSoleReleaser      bool
	DiscardJournalRequestID string
}
