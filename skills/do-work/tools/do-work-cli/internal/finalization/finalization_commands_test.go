package finalization

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

func TestRecoverFinalizationResumesJournalAfterLifecycleInterruption(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	queuePath := "do-work/queue/REQ-700.md"
	requestPath := "do-work/working/REQ-700.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	writeFinalizationFile(t, repositoryRoot, queuePath, "---\nid: REQ-700\ntitle: Fixture\nstatus: pending\n---\nBody\n")
	writeFinalizationFile(t, repositoryRoot, checkpointPath, "# Session Checkpoint\n\n## In Progress (interrupted)\n")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	claim := requeststate.Handlers()["claim"](commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{
		"REQ-700", "--request-path", queuePath, "--provenance", "default", "--writer", "host:/repo", "--commit", "--at", "2026-09-02T08:00:00Z",
	})
	if claim.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("claim result = %#v", claim)
	}
	if subject := strings.TrimSpace(runFinalizationGit(t, repositoryRoot, "log", "-1", "--format=%s")); subject != "[REQ-700] claim request lifecycle" {
		t.Fatalf("claim commit subject = %q", subject)
	}
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "after\n")
	writeFinalizationFile(t, repositoryRoot, requestPath, readFinalizationFile(t, repositoryRoot, requestPath)+"\n## Review\n\nApproved after the claim commit.\n")

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
	journal.ImageSetSHA256 = journalImageDigest(&journal)
	contents, _ := json.Marshal(journal)
	if err := os.WriteFile(journalPath, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readJournal(repositoryRoot, journalPath); err == nil || !strings.Contains(err.Error(), "payload directory") {
		t.Fatalf("unsafe journal error = %v", err)
	}
}

// finalizationRepositoryTemplate is an initialized, configured, empty repository built
// once for the whole test binary and copied per fixture.
//
// This package calls newFinalizationRepository 55 times. Building each fixture with
// `git init` and two `git config` runs cost three subprocess spawns every time, and
// subprocess spawning — not the assertions — is what put this package's largest test
// files near the gate's 30s per-file duration budget. One directory copy over the ten
// small files a template-free `git init` produces replaces all three (REQ-574).
var finalizationRepositoryTemplate string

func TestMain(m *testing.M) {
	templateRoot, err := os.MkdirTemp("", "finalization-git-template-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "finalization fixture template: %v\n", err)
		os.Exit(1)
	}
	finalizationRepositoryTemplate = templateRoot
	// `--template=` empty skips the sample hooks git would otherwise copy in, which is
	// most of the files a fresh .git holds and none of what these tests read.
	for _, arguments := range [][]string{
		{"init", "-q", "--template="},
		{"config", "user.name", "Finalization Fixture"},
		{"config", "user.email", "finalization@example.invalid"},
	} {
		command := exec.Command("git", append([]string{"-C", templateRoot}, arguments...)...)
		if output, err := command.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "finalization fixture template: git %v: %v: %s\n", arguments, err, output)
			os.RemoveAll(templateRoot)
			os.Exit(1)
		}
	}
	// `--template=` also skips .git/hooks, which an ordinary `git init` always creates
	// and which TestRecoverFinalizationResumesAfterRealPreCommitHookFailure writes a
	// real hook into. Recreate the empty directory so a copied fixture is shaped like
	// a repository git made.
	if err := os.MkdirAll(filepath.Join(templateRoot, ".git", "hooks"), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "finalization fixture template hooks: %v\n", err)
		os.RemoveAll(templateRoot)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(templateRoot)
	os.Exit(code)
}

func newFinalizationRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS(finalizationRepositoryTemplate)); err != nil {
		t.Fatalf("copy finalization fixture template: %v", err)
	}
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
