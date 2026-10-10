package finalization

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// autoManifestFixture is the REQ-760 planned-release fixture plus the judged
// inputs an action supplies: a commit message file and the release manifest the
// hand-built fixture manifest already names.
func autoManifestFixture(t *testing.T) (repositoryRoot, messagePath, releasePath, emitPath string) {
	t.Helper()
	repositoryRoot, handBuiltPath := seedPlannedReleaseFinalization(t)
	handBuiltBytes, err := os.ReadFile(handBuiltPath)
	if err != nil {
		t.Fatal(err)
	}
	var handBuilt Manifest
	if err := json.Unmarshal(handBuiltBytes, &handBuilt); err != nil {
		t.Fatal(err)
	}
	scratch := t.TempDir()
	messagePath = filepath.Join(scratch, "message.txt")
	if err := os.WriteFile(messagePath, []byte("[REQ-760] finalize planned release\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return repositoryRoot, messagePath, handBuilt.ReleaseManifestPath, filepath.Join(scratch, "manifest.json")
}

func runAutoManifest(repositoryRoot string, arguments ...string) resultmodel.CommandResult {
	return Handlers()[CommandFinalize](commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, arguments)
}

// assertNothingWritten pins the preflight contract: a refusal leaves no journal,
// no Git-private payload directory and no emitted manifest behind.
func assertNothingWritten(t *testing.T, repositoryRoot, requestID, emitPath string) {
	t.Helper()
	journalPath, payloadDirectory, err := journalLocations(repositoryRoot, requestID)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{journalPath, payloadDirectory, emitPath} {
		if _, statError := os.Lstat(path); !os.IsNotExist(statError) {
			t.Fatalf("refusal left %s behind (stat error %v)", path, statError)
		}
	}
}

func findingText(result resultmodel.CommandResult) string {
	parts := []string{}
	for _, finding := range result.Findings {
		parts = append(parts, finding.Code)
		parts = append(parts, finding.AffectedPaths...)
		parts = append(parts, finding.Evidence...)
	}
	return strings.Join(parts, "\n")
}

// Pins the reported failure: sessions hand-built this manifest and learned
// commit_paths from refusals. The emitted file must carry the planner's exact
// required set and be accepted unchanged by the existing finalize --manifest.
func TestFinalizeAutoManifestEmitsAManifestFinalizeAccepts(t *testing.T) {
	repositoryRoot, messagePath, releasePath, emitPath := autoManifestFixture(t)
	result := runAutoManifest(repositoryRoot, "--auto-manifest", "REQ-760", "--transition", "complete", "--terminal-status", "completed",
		"--message-file", messagePath, "--provenance", "primary_commit", "--release-manifest", releasePath, "--emit", emitPath)
	if result.Outcome != resultmodel.OutcomeSuccess {
		t.Fatalf("auto-manifest emit result = %#v", result)
	}
	journalPath, _, err := journalLocations(repositoryRoot, "REQ-760")
	if err != nil {
		t.Fatal(err)
	}
	if _, statError := os.Lstat(journalPath); !os.IsNotExist(statError) {
		t.Fatalf("--emit wrote a journal at %s", journalPath)
	}
	emittedBytes, err := os.ReadFile(emitPath)
	if err != nil {
		t.Fatal(err)
	}
	var emitted Manifest
	if err := json.Unmarshal(emittedBytes, &emitted); err != nil {
		t.Fatal(err)
	}
	wantPaths := []string{"CHANGELOG.md", "VERSION", "do-work/CHECKPOINT.md", "do-work/archive/REQ-760.md", "do-work/working/REQ-760.md"}
	if !slices.Equal(emitted.CommitPaths, wantPaths) {
		t.Fatalf("emitted commit_paths = %v, want %v", emitted.CommitPaths, wantPaths)
	}
	if emitted.ReleaseAt == "" || emitted.ReleaseAt != emitted.CompletedAt {
		t.Fatalf("release_at = %q, completed_at = %q", emitted.ReleaseAt, emitted.CompletedAt)
	}

	finalized := runAutoManifest(repositoryRoot, "--manifest", emitPath)
	if finalized.Outcome != resultmodel.OutcomeSuccess || len(finalized.Finalizations) != 1 || finalized.Finalizations[0].Phase != string(PhaseCleanupComplete) {
		t.Fatalf("finalize --manifest on the emitted file = %#v", finalized)
	}
}

// Pins the preflight order: a staged unrelated path must refuse before any
// journal, payload directory or emitted file exists, and name what is staged.
func TestFinalizeAutoManifestRefusesAStagedPathBeforeAnyJournal(t *testing.T) {
	repositoryRoot, messagePath, releasePath, emitPath := autoManifestFixture(t)
	writeFinalizationFile(t, repositoryRoot, "unrelated.txt", "staged by someone else\n")
	runFinalizationGit(t, repositoryRoot, "add", "unrelated.txt")
	result := runAutoManifest(repositoryRoot, "--auto-manifest", "REQ-760", "--transition", "complete", "--terminal-status", "completed",
		"--message-file", messagePath, "--provenance", "primary_commit", "--release-manifest", releasePath, "--emit", emitPath)
	if result.Outcome != resultmodel.OutcomeRefused || resultmodel.ExitCode(result.Outcome) != 1 {
		t.Fatalf("staged-path result = %#v", result)
	}
	if !strings.Contains(findingText(result), "unrelated.txt") {
		t.Fatalf("refusal does not name the staged path: %s", findingText(result))
	}
	assertNothingWritten(t, repositoryRoot, "REQ-760", emitPath)
}

// Pins "the judged fields are never invented": without the commit message the
// command refuses and names the missing flag.
func TestFinalizeAutoManifestRefusesWithoutAMessageFile(t *testing.T) {
	repositoryRoot, _, releasePath, emitPath := autoManifestFixture(t)
	result := runAutoManifest(repositoryRoot, "--auto-manifest", "REQ-760", "--transition", "complete", "--terminal-status", "completed",
		"--provenance", "primary_commit", "--release-manifest", releasePath, "--emit", emitPath)
	if result.Outcome != resultmodel.OutcomeRefused || resultmodel.ExitCode(result.Outcome) != 1 {
		t.Fatalf("missing message-file result = %#v", result)
	}
	if !strings.Contains(findingText(result), "--message-file") {
		t.Fatalf("refusal does not name --message-file: %s", findingText(result))
	}
	if _, statError := os.Lstat(emitPath); !os.IsNotExist(statError) {
		t.Fatalf("refusal emitted %s", emitPath)
	}
}

// Pins the version preflight: when the project version no longer matches the
// release manifest's old version, the release planner's refusal surfaces before
// any payload is adopted into the Git-private directory.
func TestFinalizeAutoManifestRefusesAStaleReleaseVersionBeforeWriting(t *testing.T) {
	repositoryRoot, messagePath, releasePath, emitPath := autoManifestFixture(t)
	writeFinalizationFile(t, repositoryRoot, "VERSION", "1.0.1\n")
	runFinalizationGit(t, repositoryRoot, "commit", "-qam", "someone else released 1.0.1")
	result := runAutoManifest(repositoryRoot, "--auto-manifest", "REQ-760", "--transition", "complete", "--terminal-status", "completed",
		"--message-file", messagePath, "--provenance", "primary_commit", "--release-manifest", releasePath, "--emit", emitPath)
	if result.Outcome != resultmodel.OutcomeRefused || resultmodel.ExitCode(result.Outcome) != 1 {
		t.Fatalf("stale release result = %#v", result)
	}
	if !strings.Contains(findingText(result), "RELEASE-PREIMAGE-STALE") {
		t.Fatalf("refusal does not carry RELEASE-PREIMAGE-STALE: %s", findingText(result))
	}
	assertNothingWritten(t, repositoryRoot, "REQ-760", emitPath)
}
