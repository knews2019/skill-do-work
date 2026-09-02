package finalization

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestRecoverFinalizationResumesJournalAfterLifecycleInterruption(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	requestPath := "do-work/working/REQ-700.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	writeFinalizationFile(t, repositoryRoot, requestPath, "---\nid: REQ-700\ntitle: Fixture\nstatus: claimed\nclaimed_at: 2026-09-02T08:00:00Z\n---\nBody\n")
	writeFinalizationFile(t, repositoryRoot, checkpointPath, "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-700: Fixture — claimed now — writer: host:/repo\n")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "after\n")

	manifest := Manifest{
		RequestID: "REQ-700", RequestPath: requestPath, WriterLabel: "host:/repo", Transition: "complete",
		TerminalStatus: "completed", CompletedAt: "2026-09-02T09:00:00Z",
		ExpectedRequestSHA256: digestFile(t, repositoryRoot, requestPath), ExpectedCheckpointSHA256: digestFile(t, repositoryRoot, checkpointPath),
		CommitPaths:    []string{requestPath, "do-work/archive/REQ-700.md", checkpointPath, "implementation.txt"},
		CommitMessage:  "[REQ-700] finalize fixture",
		ProvenanceMode: ProvenancePrimaryCommit,
	}
	manifestPath := filepath.Join(t.TempDir(), "manifest.json")
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, manifestBytes, 0o600); err != nil {
		t.Fatal(err)
	}

	previousHook := afterFinalizationPhase
	afterFinalizationPhase = func(phase Phase) error {
		if phase == PhaseLifecycleApplied {
			return context.Canceled
		}
		return nil
	}
	first := handleFinalize(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--manifest", manifestPath})
	afterFinalizationPhase = previousHook
	t.Cleanup(func() { afterFinalizationPhase = previousHook })
	if first.Outcome != resultmodel.OutcomeRolledBack || first.Finalization == nil || first.Rollback.Status != resultmodel.RollbackSucceeded {
		t.Fatalf("interrupted result = %#v", first)
	}

	recovered := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
	if recovered.Outcome != resultmodel.OutcomeSuccess || recovered.Finalization == nil || !recovered.Finalization.Resumed {
		t.Fatalf("recovered result = %#v", recovered)
	}
	primary := recovered.Finalization.PrimaryCommit
	metadata := recovered.Finalization.MetadataCommit
	if primary == "" || metadata == "" || primary == metadata {
		t.Fatalf("commit hashes primary=%q metadata=%q", primary, metadata)
	}
	archived := readFinalizationFile(t, repositoryRoot, "do-work/archive/REQ-700.md")
	if !strings.Contains(archived, "status: completed") || !strings.Contains(archived, "commit: "+primary) {
		t.Fatalf("archived request lacks terminal provenance:\n%s", archived)
	}
	if got := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "show", primary+":implementation.txt")); got != "after" {
		t.Fatalf("primary implementation bytes = %q", got)
	}
	journalPath, _, err := journalLocations(repositoryRoot, "REQ-700")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(journalPath); !os.IsNotExist(err) {
		t.Fatalf("journal remains after recovery: %v", err)
	}
	if status := runFinalizationGit(t, repositoryRoot, "status", "--porcelain=v1"); status != "" {
		t.Fatalf("recovery left dirt: %q", status)
	}
}

func TestReadJournalRejectsNoncanonicalPayloadDirectory(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	journalPath, _, err := journalLocations(repositoryRoot, "REQ-701")
	if err != nil {
		t.Fatal(err)
	}
	journal := Journal{
		Version: journalVersion, Phase: PhasePrepared, JournalPath: journalPath, PayloadDirectory: repositoryRoot,
		Manifest: Manifest{
			RequestID: "REQ-701", RequestPath: "do-work/working/REQ-701.md", WriterLabel: "writer", Transition: "complete",
			TerminalStatus: "completed", CompletedAt: "2026-09-02T09:00:00Z", ExpectedRequestSHA256: strings.Repeat("a", 64),
			ExpectedCheckpointSHA256: strings.Repeat("b", 64), CommitPaths: []string{"one"}, CommitMessage: "commit",
			ProvenanceMode: ProvenancePrimaryCommit,
		},
	}
	contents, _ := json.Marshal(journal)
	if err := os.WriteFile(journalPath, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readJournal(repositoryRoot, journalPath); err == nil || !strings.Contains(err.Error(), "payload directory") {
		t.Fatalf("unsafe journal error = %v", err)
	}
}

func newFinalizationRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runFinalizationGit(t, root, "init", "-q")
	runFinalizationGit(t, root, "config", "user.name", "Finalization Fixture")
	runFinalizationGit(t, root, "config", "user.email", "finalization@example.invalid")
	return root
}

func runFinalizationGit(t *testing.T, root string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, arguments...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", arguments, err, output)
	}
	return string(output)
}

func writeFinalizationFile(t *testing.T, root, path, contents string) {
	t.Helper()
	absolute := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolute, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFinalizationFile(t *testing.T, root, path string) string {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		t.Fatal(err)
	}
	return string(contents)
}

func digestFile(t *testing.T, root, path string) string {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		t.Fatal(err)
	}
	return digestBytes(contents)
}
