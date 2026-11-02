package finalization

import (
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// The inbox note the 2026-11-02 UR-161 run's recover refused on.
const unownedInboxNotePath = "do-work/inbox/2026-10-28_do-work-upstream-suggestion-month-folder-report-layer.md"

// Pins the 2026-11-02 UR-161 run refusal: recover refused FINALIZATION-DISCOVERY-AMBIGUOUS on the untracked note do-work/inbox/2026-10-28_...md.
func TestRecoverFinalizationIgnoresNotesNoRequestRecordOwns(t *testing.T) {
	t.Run("no terminal request", func(t *testing.T) {
		repositoryRoot := newFinalizationRepository(t)
		writeFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Session Checkpoint\n")
		runFinalizationGit(t, repositoryRoot, "add", ".")
		runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
		writeFinalizationFile(t, repositoryRoot, unownedInboxNotePath, "# Note\n")
		// Shares a prefix with the working folder without being inside it.
		writeFinalizationFile(t, repositoryRoot, "do-work/working-notes.md", "# Note\n")

		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
		if result.Outcome != resultmodel.OutcomeSuccess || len(result.Finalizations) != 0 {
			t.Fatalf("an unowned note stopped discovery: %#v", result)
		}
	})
	t.Run("legacy tail recovers around an untracked inbox note", func(t *testing.T) {
		repositoryRoot := newFinalizationRepository(t)
		seedSemanticLegacyTail(t, repositoryRoot)
		writeFinalizationFile(t, repositoryRoot, unownedInboxNotePath, "# Note\n")

		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
		if result.Outcome != resultmodel.OutcomeSuccess || len(result.Finalizations) != 1 || result.Finalizations[0].Phase != string(PhaseCleanupComplete) {
			t.Fatalf("an unowned note blocked legacy recovery: %#v", result)
		}
		assertNoteLeftForCommit(t, repositoryRoot, result.Finalizations[0], unownedInboxNotePath, "??")
	})
	t.Run("legacy tail recovers around a modified tracked handoff note its summary names", func(t *testing.T) {
		repositoryRoot := newFinalizationRepository(t)
		handoffPath := "do-work/HANDDOWN-UR-1.md"
		writeFinalizationFile(t, repositoryRoot, handoffPath, "# Handoff\n")
		// Naming the note in the summary must not make it the request's: only
		// record locations are finalization evidence.
		writeFinalizationFile(t, repositoryRoot, "do-work/archive/REQ-750-fixture.md", "---\nid: REQ-750\ntitle: Fixture\nstatus: completed\ncompleted_at: 2026-09-02T09:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `req-750.txt` (modified)\n- `"+handoffPath+"` (modified)\n")
		writeFinalizationFile(t, repositoryRoot, "req-750.txt", "before\n")
		runFinalizationGit(t, repositoryRoot, "add", ".")
		runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
		writeFinalizationFile(t, repositoryRoot, "req-750.txt", "after\n")
		writeFinalizationFile(t, repositoryRoot, handoffPath, "# Handoff\n\nResume here.\n")

		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
		if result.Outcome != resultmodel.OutcomeSuccess || len(result.Finalizations) != 1 || result.Finalizations[0].Phase != string(PhaseCleanupComplete) {
			t.Fatalf("a tracked handoff note blocked legacy recovery: %#v", result)
		}
		assertNoteLeftForCommit(t, repositoryRoot, result.Finalizations[0], handoffPath, " M")
	})
}

// Pins the other side of the UR-161 fix: real unfinished finalization state next
// to a note still refuses exactly as before, and the refusal names that state,
// not the note.
func TestRecoverFinalizationStillRefusesRecordStateNextToANote(t *testing.T) {
	t.Run("foreign claim line in the checkpoint", func(t *testing.T) {
		repositoryRoot := newFinalizationRepository(t)
		seedSimpleDiscoveredTail(t, repositoryRoot, "REQ-740", "req-740.txt")
		writeFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md", "# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-999: foreign — writer: other:/repo\n")
		writeFinalizationFile(t, repositoryRoot, unownedInboxNotePath, "# Note\n")

		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
		assertDiscoveryReason(t, result, "FINALIZATION-DISCOVERY-AMBIGUOUS", "do-work/CHECKPOINT.md")
	})
	// Every place a request record lives stays evidence; a missing row here is
	// a record location the predicate forgot.
	for _, recordPath := range []string{
		"do-work/queue/scratch.md",
		"do-work/working/scratch.md",
		"do-work/archive/scratch.md",
		"do-work/user-requests/UR-900/scratch.md",
		"do-work/.req-reservations/900",
		"do-work/runs/work-2026-11-02-000000/notes.md",
		"do-work/CHECKPOINT.md",
		"do-work/calibration-log.tsv",
	} {
		t.Run("untracked "+recordPath, func(t *testing.T) {
			repositoryRoot := newFinalizationRepository(t)
			writeFinalizationFile(t, repositoryRoot, "seed.txt", "seed\n")
			runFinalizationGit(t, repositoryRoot, "add", ".")
			runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
			writeFinalizationFile(t, repositoryRoot, recordPath, "unfinished\n")
			writeFinalizationFile(t, repositoryRoot, unownedInboxNotePath, "# Note\n")

			result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discover"})
			assertDiscoveryReason(t, result, "FINALIZATION-DISCOVERY-AMBIGUOUS", recordPath)
		})
	}
}

// assertNoteLeftForCommit checks that recovery neither committed nor changed a
// note, so the porcelain status it had before recovery is still its status.
func assertNoteLeftForCommit(t *testing.T, repositoryRoot string, record resultmodel.FinalizationResult, notePath, wantStatus string) {
	t.Helper()
	for _, path := range record.CommitPaths {
		if path == notePath {
			t.Fatalf("recovery associated the note with %s: %v", record.RequestID, record.CommitPaths)
		}
	}
	status := runFinalizationGit(t, repositoryRoot, "status", "--porcelain", "--untracked-files=all", "--", notePath)
	if strings.TrimRight(status, "\n") != wantStatus+" "+notePath {
		t.Fatalf("recovery touched or committed the note: %q", status)
	}
}
