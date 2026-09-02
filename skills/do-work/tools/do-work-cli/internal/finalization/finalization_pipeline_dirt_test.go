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

// Serial primary_commit finalization is the first commit after the claim, so the
// working REQ it archives is always dirty with the run's own plan, review, and
// lessons sections. REQ-513 hit exactly this: finalize refused its own preimage
// with FINALIZATION-LIFECYCLE-APPLY, and both run and run-with-recovery replayed
// the same refusal at startup.
func TestFinalizeAcceptsWorkingRequestDirtWrittenByThePipeline(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	requestPath := "do-work/working/REQ-513.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	writeFinalizationFile(t, repositoryRoot, requestPath, "---\nid: REQ-513\ntitle: Fixture\nstatus: claimed\nclaimed_at: 2026-09-02T20:45:12Z\n---\nBody\n")
	writeFinalizationFile(t, repositoryRoot, checkpointPath, "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-513: Fixture — claimed now — writer: host:/repo\n")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "claim")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "after\n")
	writeFinalizationFile(t, repositoryRoot, requestPath, readFinalizationFile(t, repositoryRoot, requestPath)+"\n## Review\n\nApproved after the claim commit.\n")

	manifest := Manifest{
		RequestID: "REQ-513", RequestPath: requestPath, WriterLabel: "host:/repo", Transition: "complete",
		TerminalStatus: "completed", CompletedAt: "2026-09-02T21:33:21Z",
		ExpectedRequestSHA256: digestFile(t, repositoryRoot, requestPath), ExpectedCheckpointSHA256: digestFile(t, repositoryRoot, checkpointPath),
		CommitPaths:    []string{requestPath, "do-work/archive/REQ-513.md", checkpointPath, "implementation.txt"},
		CommitMessage:  "[REQ-513] finalize fixture",
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
		t.Fatalf("finalize refused the pipeline's own working-REQ dirt: %#v", result)
	}
	archived := readFinalizationFile(t, repositoryRoot, "do-work/archive/REQ-513.md")
	if !strings.Contains(archived, "status: completed") || !strings.Contains(archived, "## Review") {
		t.Fatalf("archived request lost the uncommitted review section or terminal status:\n%s", archived)
	}
	if status := runFinalizationGit(t, repositoryRoot, "status", "--porcelain=v1"); status != "" {
		t.Fatalf("finalization left dirt: %q", status)
	}
}
