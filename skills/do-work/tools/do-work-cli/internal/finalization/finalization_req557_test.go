package finalization

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// REQ-557 / D-02: the helper that computed the missing commit paths used to be
// `uniqueSorted`, whose whole body was `result, _ := normalizeRepositoryPaths(paths);
// return result`. Throwing the validator's error away made the guard on the next line
// unreachable: one required path that is empty, absolute, or escaping turned the WHOLE
// required set into nil, the subtraction then found nothing missing, and prepare
// continued as if commit_paths covered every planned target.
//
// The three tests below pin three different things, and none of them substitutes for
// another:
//
//   - This first test pins the propagation INSIDE missingCommitPaths. It fails if the
//     helper's body discards the validator's error again. It calls the helper directly,
//     so it says nothing about what prepareBoundJournal does with the two results.
//   - The second pins the helper's answer on a usable set: which required targets
//     commit_paths does not cover.
//   - The third drives prepareBoundJournal and pins the CALL SITE: that prepare still
//     refuses a manifest whose commit_paths omits a planned target. It is the only one
//     of the three that fails if the `if len(missing) > 0` refusal is deleted.
//
// A call-site test cannot also pin the propagation. The required set prepare passes in
// is tool-derived, not manifest-derived: statePlan.TargetPaths only holds paths
// requeststate resolved out of the repository snapshot (a path outside it refuses with
// REQUEST-SNAPSHOT-STALE first), and every release target passed publication's
// containedPath check (RELEASE-PATH-UNSAFE) before prepare sees it. No manifest the
// validators accept can put an empty, absolute, or escaping path into that set, so a
// discard written at the call site instead of inside the helper is invisible from here.
// The first test is what covers it.
func TestMissingCommitPathsRefusesAnUnusableRequiredPathInsteadOfEmptyingTheSet(t *testing.T) {
	for _, unusableRequiredPath := range []string{"", "/etc/passwd", "../outside.md"} {
		requiredCommitPaths := []string{"do-work/archive/REQ-557.md", unusableRequiredPath}

		missing, missingError := missingCommitPaths(requiredCommitPaths, nil)

		if missingError == nil {
			t.Fatalf("missingCommitPaths accepted the unusable required path %q and returned missing=%#v; "+
				"the commit_paths guard is disabled again", unusableRequiredPath, missing)
		}
		if missing != nil {
			t.Fatalf("missingCommitPaths(%q) returned both an error and %#v", unusableRequiredPath, missing)
		}
	}
}

// The same helper on a usable set still answers the question the guard asks: which
// planned lifecycle or release targets commit_paths does not cover.
func TestMissingCommitPathsNamesTheRequiredTargetsCommitPathsOmits(t *testing.T) {
	missing, missingError := missingCommitPaths(
		[]string{"do-work/CHECKPOINT.md", "do-work/archive/REQ-557.md", "do-work/CHECKPOINT.md"},
		[]string{"do-work/CHECKPOINT.md"},
	)
	if missingError != nil {
		t.Fatalf("missingCommitPaths refused a usable set: %v", missingError)
	}
	if want := []string{"do-work/archive/REQ-557.md"}; !reflect.DeepEqual(missing, want) {
		t.Fatalf("missingCommitPaths = %#v, want %#v", missing, want)
	}
}

// The guard the helper feeds, driven end to end. A manifest that plans the archive
// move but leaves the archive path out of commit_paths must be refused before any
// journal is written: finalization would otherwise commit a lifecycle transition whose
// destination file is not in the allowlist it commits.
func TestPrepareBoundJournalRefusesCommitPathsThatOmitThePlannedArchiveTarget(t *testing.T) {
	repositoryRoot := newFinalizationRepository(t)
	requestPath := "do-work/working/REQ-721.md"
	checkpointPath := "do-work/CHECKPOINT.md"
	archivePath := "do-work/archive/REQ-721.md"
	writeFinalizationFile(t, repositoryRoot, requestPath,
		"---\nid: REQ-721\ntitle: Omitted archive target\nstatus: claimed\nclaimed_at: 2026-09-02T08:00:00Z\ncommit:\n---\n\n## Implementation Summary\n- `implementation.txt` (modified)\n")
	writeFinalizationFile(t, repositoryRoot, checkpointPath,
		"# Session Checkpoint\n\n## In Progress (interrupted)\n\n- REQ-721: Omitted archive target — claimed now — writer: host:/repo\n")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "before\n")
	runFinalizationGit(t, repositoryRoot, "add", ".")
	runFinalizationGit(t, repositoryRoot, "commit", "-qm", "seed")
	writeFinalizationFile(t, repositoryRoot, "implementation.txt", "after\n")

	manifest := Manifest{
		RequestID: "REQ-721", RequestPath: requestPath, WriterLabel: "host:/repo", Transition: "complete",
		TerminalStatus: "completed", CompletedAt: "2026-09-02T09:00:00Z",
		ExpectedRequestSHA256:    digestFile(t, repositoryRoot, requestPath),
		ExpectedCheckpointSHA256: digestFile(t, repositoryRoot, checkpointPath),
		// The planned archive destination is deliberately absent.
		CommitPaths:   []string{requestPath, checkpointPath, "implementation.txt"},
		CommitMessage: "[REQ-721] finalize without the archive path", ProvenanceMode: ProvenancePrimaryCommit,
	}
	manifestPath := filepath.Join(t.TempDir(), "manifest.json")
	manifestBytes, marshalError := json.Marshal(manifest)
	if marshalError != nil {
		t.Fatal(marshalError)
	}
	if writeError := os.WriteFile(manifestPath, manifestBytes, 0o600); writeError != nil {
		t.Fatal(writeError)
	}

	journal, resumed, prepareError := prepareBoundJournal(context.Background(), repositoryRoot, manifestPath, "REQ-721", requestPath, "complete")

	if prepareError == nil {
		t.Fatalf("prepare accepted commit_paths that omit %s and wrote a journal for %s with effective paths %v (resumed=%t); "+
			"the omitted-target guard no longer fires", archivePath, journal.ArchivedPath, journal.EffectiveCommitPaths, resumed)
	}
	if !strings.Contains(prepareError.Error(), "commit_paths omits planned lifecycle or release targets") {
		t.Fatalf("prepare error = %q, want the omitted-target refusal", prepareError)
	}
	if !strings.Contains(prepareError.Error(), archivePath) {
		t.Fatalf("prepare error = %q, want it to name the omitted target %s", prepareError, archivePath)
	}
	if journal != nil {
		t.Fatalf("a refused prepare returned a journal: %#v", journal)
	}
	journalPath, _, locationError := journalLocations(repositoryRoot, "REQ-721")
	if locationError != nil {
		t.Fatal(locationError)
	}
	if _, statError := os.Stat(journalPath); !os.IsNotExist(statError) {
		t.Fatalf("a refused prepare left a journal behind: %v", statError)
	}
}
