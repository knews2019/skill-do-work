package finalization

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// A journaled finalize declares its whole write set in the manifest, so tree
// dirt at any other path cannot be a torn tail of this transaction. REQ-560:
// commit safety used to refuse every uncommitted do-work/ path outside the
// allowlist, which stopped a run whenever a concurrent session left an
// untracked draft or a modified shared log behind — and pushed the pipeline
// into committing that foreign work under its own name to get past the check.
func TestFinalizeIgnoresForeignTreeDirtOutsideTheManifest(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	requestPath := "do-work/working/REQ-560.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	foreignTracked := "do-work/calibration-log.tsv"
	foreignUntracked := "do-work/audits/maintainability-draft.md"
	writeFinalizationFile(t, repositoryRoot, requestPath, "---\nid: REQ-560\ntitle: Fixture\nstatus: claimed\nclaimed_at: 2026-09-04T18:15:54Z\n---\nBody\n")
	writeFinalizationFile(t, repositoryRoot, checkpointPath, "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-560: Fixture — claimed now — writer: host:/repo\n")
	writeFinalizationFile(t, repositoryRoot, foreignTracked, "req\tminutes\nREQ-1\t10\n")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "claim")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "after\n")

	// Dirt this REQ does not own: another session's edit to a shared tracked
	// log, and another session's untracked draft. Neither is in the manifest.
	writeFinalizationFile(t, repositoryRoot, foreignTracked, "req\tminutes\nREQ-1\t10\nREQ-2\t25\n")
	writeFinalizationFile(t, repositoryRoot, foreignUntracked, "Concurrent audit draft, half written.\n")

	manifest := Manifest{
		RequestID: "REQ-560", RequestPath: requestPath, WriterLabel: "host:/repo", Transition: "complete",
		TerminalStatus: "completed", CompletedAt: "2026-09-04T19:02:10Z",
		ExpectedRequestSHA256: digestFile(t, repositoryRoot, requestPath), ExpectedCheckpointSHA256: digestFile(t, repositoryRoot, checkpointPath),
		CommitPaths:    []string{requestPath, "do-work/archive/REQ-560.md", checkpointPath, "implementation.txt"},
		CommitMessage:  "[REQ-560] finalize fixture",
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

	result := handleFinalize(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--manifest", manifestPath})
	if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || result.Finalization.PrimaryCommit == "" {
		t.Fatalf("finalize refused on foreign dirt outside its manifest: %#v", result)
	}

	// The foreign paths are left exactly as found: still dirty, never committed.
	status := runFinalizationGit(t, repositoryRoot, "status", "--porcelain=v1", "--untracked-files=all")
	// Only " M " is correct here. "M  " would mean the path had been staged,
	// which is the opposite of leaving it alone.
	if !strings.Contains(status, " M "+foreignTracked) {
		t.Fatalf("foreign tracked dirt was not left alone as an unstaged modification:\n%s", status)
	}
	if !strings.Contains(status, "?? "+foreignUntracked) {
		t.Fatalf("foreign untracked draft was not left alone:\n%s", status)
	}
	committed := runFinalizationGit(t, repositoryRoot, "diff-tree", "--no-commit-id", "--name-only", "-r", result.Finalization.PrimaryCommit)
	for _, path := range []string{foreignTracked, foreignUntracked} {
		if strings.Contains(committed, path) {
			t.Fatalf("the primary commit swept in a foreign path %s:\n%s", path, committed)
		}
	}
}

// REQ-560 narrowed the shared-remainder refusal to discovered recovery groups
// and deliberately kept two refusals that nothing else pins: a discovered group
// still refuses on shared dirt outside it, because there the group is inferred
// from the tree and an unattributed shared path really can be a torn tail of
// the same interrupted transaction. Without this, a later edit can widen
// recovery silently — the one thing this REQ's Constraints forbid.
func TestCommitSafetyStillRefusesSharedDirtForADiscoveredGroup(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	declared := "do-work/working/REQ-560.md"
	foreign := "do-work/calibration-log.tsv"
	writeFinalizationFile(t, repositoryRoot, declared, "declared\n")
	writeFinalizationFile(t, repositoryRoot, foreign, "req\tminutes\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "base")
	writeFinalizationFile(t, repositoryRoot, declared, "declared, edited\n")
	writeFinalizationFile(t, repositoryRoot, foreign, "req\tminutes\nREQ-2\t25\n")

	journal := &Journal{EffectiveCommitPaths: []string{declared}, Discovered: true}
	code, _, blocked := commitSafety(repositoryRoot, journal)
	if code != "FINALIZATION-AMBIGUOUS-SHARED-STATE" {
		t.Fatalf("a discovered group must still refuse on shared dirt outside it; got code %q blocked %v", code, blocked)
	}
	if len(blocked) != 1 || blocked[0] != foreign {
		t.Fatalf("the refusal must name the shared path outside the group; got %v", blocked)
	}

	// The same tree under a journaled (manifest-declared) group is the narrowing
	// this REQ shipped: the foreign path is not this transaction's to judge.
	journal.Discovered = false
	code, _, blocked = commitSafety(repositoryRoot, journal)
	if code != "" || len(blocked) != 0 {
		t.Fatalf("a journaled group must ignore foreign dirt outside its manifest; got code %q blocked %v", code, blocked)
	}
}
