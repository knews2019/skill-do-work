package finalization

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/nextselection"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestRecoverFinalizationDiscoversLegacyNoJournalTail(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	workingPath := "do-work/working/REQ-700-legacy.md"
	archivePath := "do-work/archive/REQ-700-legacy.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	requestBefore := "---\nid: REQ-700\ntitle: Legacy fixture\nstatus: claimed\nclaimed_at: 2026-09-02T08:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `implementation.txt` (modified)\n"
	requestAfter := "---\nid: REQ-700\ntitle: Legacy fixture\nstatus: completed\ncompleted_at: 2026-09-02T09:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `implementation.txt` (modified)\n"
	writeFinalizationFile(t, repositoryRoot, workingPath, requestBefore)
	writeFinalizationFile(t, repositoryRoot, checkpointPath, "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-700: Legacy fixture — claimed now — writer: host:/repo\n")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "before\n")
	writeFinalizationFile(t, repositoryRoot, "notes.txt", "before\n")
	writeFinalizationFile(t, repositoryRoot, "do-work/queue/REQ-701-next.md", "---\nid: REQ-701\ntitle: Next fixture\nstatus: pending\n---\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")

	if err := os.MkdirAll(filepath.Join(repositoryRoot, "do-work", "archive"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(repositoryRoot, filepath.FromSlash(workingPath)), filepath.Join(repositoryRoot, filepath.FromSlash(archivePath))); err != nil {
		t.Fatal(err)
	}
	writeFinalizationFile(t, repositoryRoot, archivePath, requestAfter)
	writeFinalizationFile(t, repositoryRoot, checkpointPath, "# Session Checkpoint\n\n## In Progress (interrupted)\n")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "after\n")
	writeFinalizationFile(t, repositoryRoot, "notes.txt", "unrelated\n")

	recovered := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if recovered.Outcome != resultmodel.OutcomeSuccess || recovered.Finalization == nil || len(recovered.Finalizations) != 1 {
		t.Fatalf("discovery result = %#v", recovered)
	}
	record := recovered.Finalizations[0]
	if record.RequestID != "REQ-700" || !record.Discovered || !record.Resumed {
		t.Fatalf("discovered record = %#v", record)
	}
	if record.PrimaryCommit == "" || record.MetadataCommit == "" || record.PrimaryCommit == record.MetadataCommit {
		t.Fatalf("commit identities = %#v", record)
	}
	archived := readFinalizationFile(t, repositoryRoot, archivePath)
	if !strings.Contains(archived, "commit: "+record.PrimaryCommit) {
		t.Fatalf("canonical provenance missing:\n%s", archived)
	}
	if got := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "show", record.PrimaryCommit+":implementation.txt")); got != "after" {
		t.Fatalf("primary implementation bytes = %q", got)
	}
	if got := readFinalizationFile(t, repositoryRoot, "notes.txt"); got != "unrelated\n" {
		t.Fatalf("unrelated unstaged file changed: %q", got)
	}
	journalPath, _, err := journalLocations(repositoryRoot, "REQ-700")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(journalPath); !os.IsNotExist(err) {
		t.Fatalf("journal remains after discovery: %v", err)
	}

	second := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if second.Outcome != resultmodel.OutcomeSuccess || second.Finalization != nil || len(second.Finalizations) != 0 {
		t.Fatalf("second discovery is not idempotent: %#v", second)
	}
	selected := nextselection.Handlers()[nextselection.CommandNext](commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
	if selected.Outcome != resultmodel.OutcomeSuccess || len(selected.Selected) != 1 || selected.Selected[0].RequestID != "REQ-701" {
		t.Fatalf("selection after recovery = %#v", selected)
	}
	claimed := requeststate.Handlers()["claim"](commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{
		"REQ-701", "--request-path", selected.Selected[0].RequestPath, "--provenance", "default", "--writer", "host:/repo", "--at", "2026-09-02T10:00:00Z",
	})
	if claimed.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("claim after recovery = %#v", claimed)
	}
}

func TestValidateManifestRequiresExplicitProvenanceMode(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	writeFinalizationFile(t, repositoryRoot, "seed.txt", "seed\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	manifest := Manifest{
		RequestID: "REQ-701", RequestPath: "do-work/working/REQ-701.md", WriterLabel: "writer", Transition: "complete",
		TerminalStatus: "completed", CompletedAt: "2026-09-02T09:00:00Z", ExpectedRequestSHA256: strings.Repeat("a", 64),
		ExpectedCheckpointSHA256: strings.Repeat("b", 64), CommitPaths: []string{"seed.txt"}, CommitMessage: "fixture",
	}
	if err := validateManifest(repositoryRoot, manifest); err == nil || !strings.Contains(err.Error(), "provenance_mode") {
		t.Fatalf("missing provenance mode error = %v", err)
	}
	manifest.ProvenanceMode = ProvenancePrimaryCommit
	manifest.ImplementationHash = "abcdef0"
	if err := validateManifest(repositoryRoot, manifest); err == nil || !strings.Contains(err.Error(), "forbids") {
		t.Fatalf("primary hash error = %v", err)
	}
	manifest.ProvenanceMode = ProvenanceSuppliedCommit
	manifest.ImplementationHash = strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "rev-parse", "--short", "HEAD"))
	if err := validateManifest(repositoryRoot, manifest); err != nil {
		t.Fatalf("valid supplied commit refused: %v", err)
	}
	manifest.ImplementationHash = "ABCDEF0"
	if err := validateManifest(repositoryRoot, manifest); err == nil || !strings.Contains(err.Error(), "lowercase") {
		t.Fatalf("uppercase supplied hash error = %v", err)
	}
}

func TestRecoverFinalizationDiscoveryRefusesExistingIndex(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	writeFinalizationFile(t, repositoryRoot, "ordinary.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "ordinary.txt", "staged\n")
	runFinalizationGit(t, repositoryRoot, "add", "ordinary.txt")

	result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
	if result.Outcome != resultmodel.OutcomeRefused || len(result.Findings) != 1 || result.Findings[0].Code != "FINALIZATION-DISCOVERY-STAGED" {
		t.Fatalf("staged discovery result = %#v", result)
	}
	if got := readFinalizationFile(t, repositoryRoot, "ordinary.txt"); got != "staged\n" {
		t.Fatalf("staged path changed: %q", got)
	}
}

func TestMatchingHeadCommitUsesPreparedDiffIdentityAcrossLaterCommit(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "before\n")
	writeFinalizationFile(t, repositoryRoot, "later.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "owned.txt", "after\n")
	preparedHead, preparedDiff, err := preparedCommitIdentity(repositoryRoot, []string{"owned.txt"})
	if err != nil {
		t.Fatal(err)
	}
	runFinalizationGit(t, repositoryRoot, "add", "owned.txt")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "primary")
	primaryCommit := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "rev-parse", "HEAD"))
	writeFinalizationFile(t, repositoryRoot, "later.txt", "after\n")
	runFinalizationGit(t, repositoryRoot, "add", "later.txt")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "later")

	journal := &Journal{PreparedHead: preparedHead, PreparedDiffSHA256: preparedDiff, EffectiveCommitPaths: []string{"owned.txt"}}
	if matched, ok := matchingHeadCommit(repositoryRoot, journal); !ok || matched != primaryCommit {
		t.Fatalf("matched commit = %q, %t; want %q", matched, ok, primaryCommit)
	}
}
