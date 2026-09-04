package finalization

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/commandruntime"
	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// REQ-547: the checkpoint's contents stopped deciding whether a REQ could be
// finalized at all. A REQ whose entry is missing planned one fewer path than it
// applied, the journal's two image sets disagreed on the count, and finalization
// refused with FINALIZATION-LIFECYCLE-CONFLICT on every replay. The fixture
// carries an estimate so the calibration append joins the target set, which is
// the exact four-path shape the refused release journal had on 2026-09-03.
func TestFinalizeCompletesWhenTheCheckpointHoldsNoEntryForTheRequest(t *testing.T) {
	repositoryRoot, manifestPath := seedCheckpointShapedFinalization(t,
		"# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-999: Another REQ — claimed now — writer: host:/repo\n", "host:/repo")

	journal, _, err := prepareJournal(context.Background(), repositoryRoot, manifestPath)
	if err != nil {
		t.Fatalf("prepare refused a REQ with no checkpoint entry: %v", err)
	}
	if preimagePaths, postimagePaths := imagePaths(journal.LifecyclePreimages), imagePaths(journal.LifecyclePostimages); !reflect.DeepEqual(preimagePaths, postimagePaths) {
		t.Fatalf("journal image sets name different paths:\npre  = %v\npost = %v", preimagePaths, postimagePaths)
	}

	result := handleFinalize(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--manifest", manifestPath})
	if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || result.Finalization.Phase != string(PhaseCleanupComplete) {
		t.Fatalf("finalize with no checkpoint entry = %#v", result)
	}
	if _, statError := os.Stat(filepath.Join(repositoryRoot, "do-work", "archive", "REQ-740.md")); statError != nil {
		t.Fatalf("the REQ was not archived: %v", statError)
	}
	if checkpoint := readFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md"); !strings.Contains(checkpoint, "REQ-999") {
		t.Fatalf("another REQ's checkpoint entry was lost:\n%s", checkpoint)
	}
}

// REQ-547: finalize used to remove only an entry carrying the manifest's own
// writer label. A manifest whose label did not match the one advance wrote at
// claim time returned typed success and archived the REQ while the checkpoint
// went on listing it as in flight, with nothing reported. Hit three times on
// 2026-09-04.
func TestFinalizeClearsACheckpointEntryLabelledByAnotherWriter(t *testing.T) {
	repositoryRoot, manifestPath := seedCheckpointShapedFinalization(t,
		"# Session Checkpoint\n\n## In Progress (interrupted)\n\n"+
			"- REQ-740: Checkpoint fixture — claimed now — writer: vm:/other/checkout\n"+
			"  - **Last known state:** claimed by a session with another label\n"+
			"- REQ-999: Another REQ — claimed now — writer: host:/repo\n", "host:/repo")

	result := handleFinalize(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--manifest", manifestPath})
	if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || result.Finalization.Phase != string(PhaseCleanupComplete) {
		t.Fatalf("finalize with a foreign writer label = %#v", result)
	}
	checkpoint := readFinalizationFile(t, repositoryRoot, "do-work/CHECKPOINT.md")
	if strings.Contains(checkpoint, "REQ-740") || strings.Contains(checkpoint, "another label") {
		t.Fatalf("the archived REQ is still listed as in flight:\n%s", checkpoint)
	}
	if !strings.Contains(checkpoint, "REQ-999") {
		t.Fatalf("another REQ's checkpoint entry was removed with it:\n%s", checkpoint)
	}
}

// REQ-547: a journal written in a refusing state had no exit but moving the
// Git-private file by hand. --discard-journal is that exit, and it is bounded:
// it removes a journal only while the transaction is provably pre-mutation.
func TestRecoverFinalizationDiscardsOnlyAPreMutationJournal(t *testing.T) {
	t.Run("discards a prepared journal and lets finalize run again", func(t *testing.T) {
		repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenancePrimaryCommit)
		journal, _, err := prepareJournal(context.Background(), repositoryRoot, manifestPath)
		if err != nil {
			t.Fatal(err)
		}
		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discard-journal", "REQ-720"})
		if result.Outcome != resultmodel.OutcomeSuccess || result.Finalization == nil || !containsFinalizationReason(result.Finalization.ReasonCodes, "FINALIZATION-JOURNAL-DISCARDED") {
			t.Fatalf("discard result = %#v", result)
		}
		if _, statError := os.Stat(journal.JournalPath); !os.IsNotExist(statError) {
			t.Fatalf("the journal survived its discard: %v", statError)
		}
		again := handleFinalize(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--manifest", manifestPath})
		if again.Outcome != resultmodel.OutcomeSuccess || again.Finalization == nil || again.Finalization.Phase != string(PhaseCleanupComplete) {
			t.Fatalf("finalize after discard = %#v", again)
		}
	})

	t.Run("refuses a journal that already applied a phase", func(t *testing.T) {
		repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenancePrimaryCommit)
		journal, _, err := prepareJournal(context.Background(), repositoryRoot, manifestPath)
		if err != nil {
			t.Fatal(err)
		}
		journal.Phase = PhaseLifecycleApplied
		if writeError := writeJournal(journal); writeError != nil {
			t.Fatal(writeError)
		}
		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discard-journal", "REQ-720"})
		if result.Outcome != resultmodel.OutcomeRefused || result.Finalization == nil || !containsFinalizationReason(result.Finalization.ReasonCodes, "FINALIZATION-DISCARD-REFUSED") {
			t.Fatalf("discard of an applied journal = %#v", result)
		}
		if _, statError := os.Stat(journal.JournalPath); statError != nil {
			t.Fatalf("a refused discard removed the journal anyway: %v", statError)
		}
	})

	t.Run("refuses once a lifecycle path left its preimage", func(t *testing.T) {
		repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenancePrimaryCommit)
		journal, _, err := prepareJournal(context.Background(), repositoryRoot, manifestPath)
		if err != nil {
			t.Fatal(err)
		}
		writeFinalizationFile(t, repositoryRoot, "do-work/working/REQ-720.md", "mutation started after the journal was written\n")
		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discard-journal", "REQ-720"})
		if result.Outcome != resultmodel.OutcomeRefused || result.Finalization == nil || !containsFinalizationReason(result.Finalization.ReasonCodes, "FINALIZATION-DISCARD-REFUSED") {
			t.Fatalf("discard after a lifecycle mutation = %#v", result)
		}
		if _, statError := os.Stat(journal.JournalPath); statError != nil {
			t.Fatalf("a refused discard removed the journal anyway: %v", statError)
		}
	})

	t.Run("reports a REQ with no journal instead of discarding something else", func(t *testing.T) {
		repositoryRoot, _ := seedPlannedFinalization(t, ProvenancePrimaryCommit)
		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, []string{"--discard-journal", "REQ-721"})
		if result.Outcome != resultmodel.OutcomeFailure || len(result.Findings) != 1 || result.Findings[0].Code != "FINALIZATION-JOURNAL-MISSING" {
			t.Fatalf("discard without a journal = %#v", result)
		}
	})
}

// REQ-547: the count check in imageSetState is what catches a journal whose two
// sets genuinely disagree about a path, so widening the checkpoint case must not
// have blunted it. Both halves of the journal are covered because each has its
// own image-set comparison.
func TestFinalizationStillRefusesJournalImageSetsThatDisagree(t *testing.T) {
	t.Run("lifecycle", func(t *testing.T) {
		repositoryRoot, manifestPath := seedPlannedFinalization(t, ProvenancePrimaryCommit)
		journal, _, err := prepareJournal(context.Background(), repositoryRoot, manifestPath)
		if err != nil {
			t.Fatal(err)
		}
		journal.LifecyclePostimages = append(journal.LifecyclePostimages, FileImage{Path: "do-work/never-planned.md", Exists: true, Bytes: []byte("unplanned\n"), Mode: 0o644})
		rewriteJournalImages(t, journal)
		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
		assertLifecycleConflictWithoutArchive(t, repositoryRoot, result, "do-work/archive/REQ-720.md", "FINALIZATION-LIFECYCLE-CONFLICT")
	})

	t.Run("release", func(t *testing.T) {
		repositoryRoot, manifestPath := seedPlannedReleaseFinalization(t)
		journal, _, err := prepareJournal(context.Background(), repositoryRoot, manifestPath)
		if err != nil {
			t.Fatal(err)
		}
		journal.ReleasePostimages = append(journal.ReleasePostimages, FileImage{Path: "NOTICE.md", Exists: true, Bytes: []byte("unplanned\n"), Mode: 0o644})
		rewriteJournalImages(t, journal)
		result := handleRecoverFinalization(commandruntime.ExecutionContext{RepositoryRoot: repositoryRoot}, nil)
		assertLifecycleConflictWithoutArchive(t, repositoryRoot, result, "NOTICE.md", "FINALIZATION-RELEASE-CONFLICT")
	})
}

func assertLifecycleConflictWithoutArchive(t *testing.T, repositoryRoot string, result resultmodel.CommandResult, unwrittenPath, expectedCode string) {
	t.Helper()
	if len(result.Finalizations) != 1 || !containsFinalizationReason(result.Finalizations[0].ReasonCodes, expectedCode) {
		t.Fatalf("disagreeing image sets did not refuse with %s: %#v", expectedCode, result)
	}
	if _, statError := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(unwrittenPath))); !os.IsNotExist(statError) {
		t.Fatalf("the refused journal still wrote %s: %v", unwrittenPath, statError)
	}
}

// rewriteJournalImages persists hand-edited image sets under a recomputed
// digest, so the journal fails its own path comparison rather than the stored
// image integrity check that guards a tampered file.
func rewriteJournalImages(t *testing.T, journal *Journal) {
	t.Helper()
	journal.ImageSetSHA256 = ""
	if err := writeJournal(journal); err != nil {
		t.Fatal(err)
	}
}

// seedCheckpointShapedFinalization plans a complete for a REQ carrying an
// estimate, so the transaction spans the request source, its archive
// destination, the calibration append, and the checkpoint - the shape whose
// path counts disagreed. The checkpoint body is the case under test.
func seedCheckpointShapedFinalization(t *testing.T, checkpointBody, writerLabel string) (string, string) {
	t.Helper()
	repositoryRoot := newFinalizationRepository(t)
	requestPath := "do-work/working/REQ-740.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	writeFinalizationFile(t, repositoryRoot, requestPath, "---\nid: REQ-740\ntitle: Checkpoint fixture\nstatus: claimed\nclaimed_at: 2026-09-02T08:00:00Z\nestimate:\n  p50_active_minutes: 30\ncommit:\n---\n\n## Implementation Summary\n- `implementation.txt` (modified)\n")
	writeFinalizationFile(t, repositoryRoot, checkpointPath, checkpointBody)
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "after\n")
	manifest := Manifest{
		RequestID: "REQ-740", RequestPath: requestPath, WriterLabel: writerLabel, Transition: "complete", TerminalStatus: "completed",
		CompletedAt: "2026-09-02T09:00:00Z", ExpectedRequestSHA256: digestFile(t, repositoryRoot, requestPath), ExpectedCheckpointSHA256: digestFile(t, repositoryRoot, checkpointPath),
		CommitPaths:   []string{requestPath, "do-work/archive/REQ-740.md", checkpointPath, "do-work/calibration-log.tsv", "implementation.txt"},
		CommitMessage: "[REQ-740] finalize the checkpoint fixture", ProvenanceMode: ProvenancePrimaryCommit,
	}
	manifestPath := filepath.Join(t.TempDir(), "manifest.json")
	contents, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, contents, 0o600); err != nil {
		t.Fatal(err)
	}
	return repositoryRoot, manifestPath
}
