//go:build unix

// Every fingerprint here runs a real bounded toolchain probe, and the probe
// runner fails closed on a platform that cannot prove descendant process
// ownership, so the whole file is scoped to the platforms where it can run.

package heavyverification

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The fixture declares two independently covered stages plus one that declares
// no fingerprint inputs at all. Each stage's toolchain probe reads a file under
// do-work/, which is the one tree no seal includes, so a toolchain-change case
// moves the probe output and nothing else. Everything the manifest classifies
// nowhere — gate/verify.sh and shared/notes.md — is sealed into BOTH stages, so
// a change there has to force both of them.
const fastStageTestManifest = `{
  "schema_version": 1,
  "stages": [
    {
      "id": "alpha-stage",
      "argv": ["sh", "module-alpha/run.sh"],
      "coverage": [{"kind": "subtree", "path": "module-alpha"}],
      "fingerprint": {
        "toolchain_probes": [["cat", "do-work/alpha-toolchain.txt"]],
        "environment_variables": ["FAST_STAGE_TEST_SELECTOR"]
      }
    },
    {
      "id": "beta-stage",
      "argv": ["sh", "module-beta/run.sh"],
      "coverage": [{"kind": "subtree", "path": "module-beta"}],
      "fingerprint": {
        "toolchain_probes": [["cat", "do-work/beta-toolchain.txt"]]
      }
    },
    {
      "id": "unfingerprinted-stage",
      "argv": ["sh", "module-alpha/run.sh"],
      "coverage": [{"kind": "subtree", "path": "module-alpha"}]
    }
  ],
  "non_stage_coverage": [{"kind": "subtree", "path": "do-work"}]
}`

const fastStageTestManifestPath = "fast-stages.json"

var fastStageTestArgv = map[string][]string{
	"alpha-stage":           {"sh", "module-alpha/run.sh"},
	"beta-stage":            {"sh", "module-beta/run.sh"},
	"unfingerprinted-stage": {"sh", "module-alpha/run.sh"},
}

// buildFastStageTemplateRepository initializes one repository that every case
// below copies. One `git init` per test function rather than one per case: the
// cost of these fixtures is subprocess spawning, and the copy is a tenth of it.
func buildFastStageTemplateRepository(t *testing.T) string {
	t.Helper()
	templateRoot := filepath.Join(t.TempDir(), "template")
	if err := os.Mkdir(templateRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	runHeavyTestGit(t, templateRoot, "init", "-q")
	runHeavyTestGit(t, templateRoot, "config", "user.name", "Fast Stage Test")
	runHeavyTestGit(t, templateRoot, "config", "user.email", "fast-stage@example.test")
	writeHeavyTestFile(t, templateRoot, ".gitignore", "*.generated\n")
	writeHeavyTestFile(t, templateRoot, fastStageTestManifestPath, fastStageTestManifest)
	writeHeavyTestFile(t, templateRoot, "module-alpha/run.sh", "exit 0\n")
	writeHeavyTestFile(t, templateRoot, "module-alpha/source.txt", "alpha source 1\n")
	writeHeavyTestFile(t, templateRoot, "module-alpha/testdata/fixture.txt", "alpha fixture 1\n")
	writeHeavyTestFile(t, templateRoot, "module-beta/run.sh", "exit 0\n")
	writeHeavyTestFile(t, templateRoot, "module-beta/source.txt", "beta source 1\n")
	writeHeavyTestFile(t, templateRoot, "gate/verify.sh", "# gate script\n")
	writeHeavyTestFile(t, templateRoot, "shared/notes.md", "shared notes 1\n")
	writeHeavyTestFile(t, templateRoot, "do-work/queue/REQ-001-fixture.md", "queue state 1\n")
	writeHeavyTestFile(t, templateRoot, "do-work/alpha-toolchain.txt", "alpha toolchain 1\n")
	writeHeavyTestFile(t, templateRoot, "do-work/beta-toolchain.txt", "beta toolchain 1\n")
	commitHeavyTestChanges(t, templateRoot, "fast-stage fixture")
	return templateRoot
}

func copyFastStageRepository(t *testing.T, templateRoot string) string {
	t.Helper()
	repositoryRoot := filepath.Join(t.TempDir(), "repository")
	if err := os.CopyFS(repositoryRoot, os.DirFS(templateRoot)); err != nil {
		t.Fatal(err)
	}
	return repositoryRoot
}

// decideFastStageFields runs the real decision command path and splits the one
// line the gate parses, so every case asserts the shipped output shape too.
func decideFastStageFields(t *testing.T, repositoryRoot, stageID string, evaluatedAt time.Time) (disposition, reason, fingerprint, recordedAt string) {
	t.Helper()
	line, err := DecideFastStage(FastStageDecisionRequest{
		RepositoryRoot: repositoryRoot, ManifestPath: fastStageTestManifestPath,
		StageID: stageID, SuppliedArgv: fastStageTestArgv[stageID], EvaluatedAt: evaluatedAt,
	})
	if err != nil {
		t.Fatalf("decide %s: %v", stageID, err)
	}
	if !strings.HasSuffix(line, "\n") {
		t.Fatalf("decision line for %s is not newline-terminated: %q", stageID, line)
	}
	fields := strings.Fields(line)
	if len(fields) != 4 {
		t.Fatalf("decision line for %s has %d fields, want 4: %q", stageID, len(fields), line)
	}
	return fields[0], fields[1], fields[2], fields[3]
}

// establishFastStageGreen records one successful result for a stage, stamped at
// the instant the caller names, exactly as the gate would after a zero exit.
func establishFastStageGreen(t *testing.T, repositoryRoot, stageID string, recordedAt time.Time) string {
	t.Helper()
	_, _, fingerprint, _ := decideFastStageFields(t, repositoryRoot, stageID, recordedAt)
	if fingerprint == "-" {
		t.Fatalf("stage %s has no determinable fingerprint to record", stageID)
	}
	if err := RecordFastStage(FastStageRecordRequest{
		RepositoryRoot: repositoryRoot, ManifestPath: fastStageTestManifestPath,
		StageID: stageID, SuppliedArgv: fastStageTestArgv[stageID],
		SuppliedFingerprint: fingerprint, StageExitStatus: 0, RecordedAt: recordedAt,
	}); err != nil {
		t.Fatalf("record %s: %v", stageID, err)
	}
	return fingerprint
}

func writeFastStageRecord(t *testing.T, repositoryRoot string, record storedFastStageEvidence) {
	t.Helper()
	store, err := openFastStageEvidenceStore(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.WriteFastStageEvidence(record); err != nil {
		t.Fatal(err)
	}
}

func readFastStageRecord(t *testing.T, repositoryRoot, stageID string) storedFastStageEvidence {
	t.Helper()
	store, err := openFastStageEvidenceStore(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	record, exists, err := store.ReadFastStageEvidence(stageID)
	if err != nil || !exists {
		t.Fatalf("read %s evidence: exists=%v err=%v", stageID, exists, err)
	}
	return record
}

// TestFastStageReuseDecisionTable is the acceptance table for requirement 2.
// Positive cases prove reuse happens; every other case proves it does not, and
// each asserts the exact reason code, because a reuse mechanism whose tests only
// prove that reuse happens cannot show it fails closed.
func TestFastStageReuseDecisionTable(t *testing.T) {
	templateRoot := buildFastStageTemplateRepository(t)
	recordedAt := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)
	evaluatedAt := recordedAt.Add(time.Minute)

	decisionCases := []struct {
		name              string
		stageID           string
		mutate            func(t *testing.T, repositoryRoot string)
		skipPriorEvidence bool
		expectDisposition string
		expectReason      string
	}{
		{
			name: "nothing changed", stageID: "alpha-stage",
			mutate:            func(*testing.T, string) {},
			expectDisposition: LaneDispositionReused, expectReason: laneReasonFingerprintMatch,
		},
		{
			// The one declared exemption. Queue bookkeeping is not stage input,
			// and the orchestrator rewrites it while the gate runs.
			name: "queue state changed", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				writeHeavyTestFile(t, root, "do-work/queue/REQ-002-fixture.md", "queue state 2\n")
			},
			expectDisposition: LaneDispositionReused, expectReason: laneReasonFingerprintMatch,
		},
		{
			// Per-stage separation: a change confined to another stage's covered
			// subtree leaves this stage's inputs provably unchanged.
			name: "only the other stage's module changed", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				writeHeavyTestFile(t, root, "module-beta/source.txt", "beta source 2\n")
			},
			expectDisposition: LaneDispositionReused, expectReason: laneReasonFingerprintMatch,
		},
		{
			name: "covered source changed and committed", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				writeHeavyTestFile(t, root, "module-alpha/source.txt", "alpha source 2\n")
				commitHeavyTestChanges(t, root, "change alpha source")
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintMismatch,
		},
		{
			// Requirement 2 names uncommitted changes explicitly. This is the
			// case a committed-object seal would report as a false green.
			name: "covered source changed and left uncommitted", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				writeHeavyTestFile(t, root, "module-alpha/source.txt", "alpha source 2\n")
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintMismatch,
		},
		{
			name: "covered fixture changed", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				writeHeavyTestFile(t, root, "module-alpha/testdata/fixture.txt", "alpha fixture 2\n")
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintMismatch,
		},
		{
			name: "gate script changed", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				writeHeavyTestFile(t, root, "gate/verify.sh", "# gate script, edited\n")
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintMismatch,
		},
		{
			name: "unclassified path changed", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				writeHeavyTestFile(t, root, "shared/notes.md", "shared notes 2\n")
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintMismatch,
		},
		{
			name: "untracked file added under coverage", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				writeHeavyTestFile(t, root, "module-alpha/scratch.txt", "new input\n")
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintMismatch,
		},
		{
			name: "untracked unclassified file added", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				writeHeavyTestFile(t, root, "shared/scratch.txt", "new input\n")
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintMismatch,
		},
		{
			name: "environment variable added", stageID: "alpha-stage",
			mutate: func(t *testing.T, _ string) {
				t.Setenv("FAST_STAGE_TEST_ADDED", "1")
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintMismatch,
		},
		{
			name: "declared selector variable changed", stageID: "alpha-stage",
			mutate: func(t *testing.T, _ string) {
				t.Setenv("FAST_STAGE_TEST_SELECTOR", "off")
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintMismatch,
		},
		{
			name: "toolchain probe output changed", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				writeHeavyTestFile(t, root, "do-work/alpha-toolchain.txt", "alpha toolchain 2\n")
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintMismatch,
		},
		{
			name: "toolchain probe cannot run", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				if err := os.Remove(filepath.Join(root, "do-work", "alpha-toolchain.txt")); err != nil {
					t.Fatal(err)
				}
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintUncertain,
		},
		{
			name: "tracked covered path deleted from the worktree", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				if err := os.Remove(filepath.Join(root, "module-alpha", "source.txt")); err != nil {
					t.Fatal(err)
				}
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintUncertain,
		},
		{
			name: "covered path replaced by a symlink", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				sourcePath := filepath.Join(root, "module-alpha", "source.txt")
				if err := os.Remove(sourcePath); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(filepath.Join(root, "shared", "notes.md"), sourcePath); err != nil {
					t.Fatal(err)
				}
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintUncertain,
		},
		{
			// A stage with no declared toolchain inputs can never be reused: an
			// undeclared toolchain is an input this command cannot determine.
			name: "stage declares no fingerprint inputs", stageID: "unfingerprinted-stage",
			mutate:            func(*testing.T, string) {},
			skipPriorEvidence: true,
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonFingerprintUncertain,
		},
		{
			name: "no prior evidence", stageID: "alpha-stage",
			mutate:            func(*testing.T, string) {},
			skipPriorEvidence: true,
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonNoPriorEvidence,
		},
		{
			name: "record is stamped with a non-zero exit status", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				record := readFastStageRecord(t, root, "alpha-stage")
				record.ExitStatus = 1
				writeFastStageRecord(t, root, record)
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonEvidenceUnusable,
		},
		{
			name: "record carries no fingerprint", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				record := readFastStageRecord(t, root, "alpha-stage")
				record.FingerprintSHA256 = ""
				writeFastStageRecord(t, root, record)
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonEvidenceUnusable,
		},
		{
			name: "record belongs to another working tree", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				record := readFastStageRecord(t, root, "alpha-stage")
				record.WorkingTreeRoot = filepath.Join(record.WorkingTreeRoot, "..", "other-worktree")
				writeFastStageRecord(t, root, record)
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonEvidenceUnusable,
		},
		{
			name: "record carries a different command", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				record := readFastStageRecord(t, root, "alpha-stage")
				record.CommandArgv = []string{"sh", "module-alpha/other.sh"}
				writeFastStageRecord(t, root, record)
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonEvidenceUnusable,
		},
		{
			// A heavy-lane record written at a fast stage's key must not read as
			// fast-stage evidence: the two guarantees are not interchangeable.
			name: "record carries a foreign schema version", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				record := readFastStageRecord(t, root, "alpha-stage")
				record.SchemaVersion = laneEvidenceSchemaVersion + 99
				writeFastStageRecord(t, root, record)
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonEvidenceUnusable,
		},
		{
			name: "record is stamped in the future", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				record := readFastStageRecord(t, root, "alpha-stage")
				record.RecordedAt = evaluatedAt.Add(time.Hour).UTC().Format(time.RFC3339)
				writeFastStageRecord(t, root, record)
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonEvidenceUnusable,
		},
		{
			name: "record is readable beyond its owner", stageID: "alpha-stage",
			mutate: func(t *testing.T, root string) {
				store, err := openFastStageEvidenceStore(root)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.Chmod(store.recordPath("alpha-stage"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			expectDisposition: LaneDispositionExecuted, expectReason: laneReasonEvidenceUnusable,
		},
	}

	for _, decisionCase := range decisionCases {
		t.Run(decisionCase.name, func(t *testing.T) {
			repositoryRoot := copyFastStageRepository(t, templateRoot)
			if !decisionCase.skipPriorEvidence {
				establishFastStageGreen(t, repositoryRoot, decisionCase.stageID, recordedAt)
			}
			decisionCase.mutate(t, repositoryRoot)
			disposition, reason, _, _ := decideFastStageFields(t, repositoryRoot, decisionCase.stageID, evaluatedAt)
			if disposition != decisionCase.expectDisposition || reason != decisionCase.expectReason {
				t.Fatalf("decision = %s/%s, want %s/%s",
					disposition, reason, decisionCase.expectDisposition, decisionCase.expectReason)
			}
		})
	}
}

// TestFastStageEvidenceExpiresIndependentlyOfTheFingerprint pins the ceiling as
// a second, independent condition: a matching fingerprint does not extend it,
// and a reuse never refreshes the record that authorized it.
func TestFastStageEvidenceExpiresIndependentlyOfTheFingerprint(t *testing.T) {
	repositoryRoot := copyFastStageRepository(t, buildFastStageTemplateRepository(t))
	recordedAt := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)
	establishFastStageGreen(t, repositoryRoot, "alpha-stage", recordedAt)

	disposition, reason, _, reportedRecordedAt := decideFastStageFields(t, repositoryRoot, "alpha-stage", recordedAt.Add(laneEvidenceMaximumAge))
	if disposition != LaneDispositionReused || reason != laneReasonFingerprintMatch {
		t.Fatalf("at the ceiling the decision = %s/%s, want reused/fingerprint_match", disposition, reason)
	}
	if reportedRecordedAt != recordedAt.Format(time.RFC3339) {
		t.Fatalf("reuse reported recorded_at %s, want %s", reportedRecordedAt, recordedAt.Format(time.RFC3339))
	}
	disposition, reason, _, _ = decideFastStageFields(t, repositoryRoot, "alpha-stage", recordedAt.Add(laneEvidenceMaximumAge+time.Second))
	if disposition != LaneDispositionExecuted || reason != laneReasonEvidenceExpired {
		t.Fatalf("past the ceiling the decision = %s/%s, want executed/evidence_expired", disposition, reason)
	}
	if readFastStageRecord(t, repositoryRoot, "alpha-stage").RecordedAt != recordedAt.Format(time.RFC3339) {
		t.Fatal("a reuse refreshed the record that authorized it")
	}
}

// TestFastStageRecordingRefusesEverythingButAMeasuredGreen covers requirement 3:
// a skipped, failed, or interrupted stage supplies no evidence, and a stage that
// modified its own inputs while running records nothing either.
func TestFastStageRecordingRefusesEverythingButAMeasuredGreen(t *testing.T) {
	templateRoot := buildFastStageTemplateRepository(t)
	recordedAt := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)

	t.Run("a non-zero stage status records nothing", func(t *testing.T) {
		repositoryRoot := copyFastStageRepository(t, templateRoot)
		_, _, fingerprint, _ := decideFastStageFields(t, repositoryRoot, "alpha-stage", recordedAt)
		if err := RecordFastStage(FastStageRecordRequest{
			RepositoryRoot: repositoryRoot, ManifestPath: fastStageTestManifestPath,
			StageID: "alpha-stage", SuppliedArgv: fastStageTestArgv["alpha-stage"],
			SuppliedFingerprint: fingerprint, StageExitStatus: 3, RecordedAt: recordedAt,
		}); err == nil {
			t.Fatal("recording accepted a stage that exited 3")
		}
		if _, reason, _, _ := decideFastStageFields(t, repositoryRoot, "alpha-stage", recordedAt); reason != laneReasonNoPriorEvidence {
			t.Fatalf("after a refused recording the reason = %s, want no_prior_evidence", reason)
		}
	})

	t.Run("a stage that changed its own inputs records nothing", func(t *testing.T) {
		repositoryRoot := copyFastStageRepository(t, templateRoot)
		_, _, fingerprint, _ := decideFastStageFields(t, repositoryRoot, "alpha-stage", recordedAt)
		writeHeavyTestFile(t, repositoryRoot, "module-alpha/source.txt", "rewritten by the stage itself\n")
		if err := RecordFastStage(FastStageRecordRequest{
			RepositoryRoot: repositoryRoot, ManifestPath: fastStageTestManifestPath,
			StageID: "alpha-stage", SuppliedArgv: fastStageTestArgv["alpha-stage"],
			SuppliedFingerprint: fingerprint, StageExitStatus: 0, RecordedAt: recordedAt,
		}); err == nil {
			t.Fatal("recording accepted a fingerprint the tree no longer produces")
		}
		if _, reason, _, _ := decideFastStageFields(t, repositoryRoot, "alpha-stage", recordedAt); reason != laneReasonNoPriorEvidence {
			t.Fatalf("after a refused recording the reason = %s, want no_prior_evidence", reason)
		}
	})

	t.Run("invalidation revokes a prior success", func(t *testing.T) {
		repositoryRoot := copyFastStageRepository(t, templateRoot)
		establishFastStageGreen(t, repositoryRoot, "alpha-stage", recordedAt)
		if err := InvalidateFastStage(repositoryRoot, "alpha-stage"); err != nil {
			t.Fatal(err)
		}
		if _, reason, _, _ := decideFastStageFields(t, repositoryRoot, "alpha-stage", recordedAt); reason != laneReasonNoPriorEvidence {
			t.Fatalf("after invalidation the reason = %s, want no_prior_evidence", reason)
		}
		if err := InvalidateFastStage(repositoryRoot, "alpha-stage"); err != nil {
			t.Fatalf("invalidating an already-absent record must succeed: %v", err)
		}
	})
}

// TestFastStageRefusesWhenTheCallerAndManifestDisagree pins the drift guard: the
// gate supplies the command it is about to run, so a call site that no longer
// matches the manifest is refused rather than answered from either version.
func TestFastStageRefusesWhenTheCallerAndManifestDisagree(t *testing.T) {
	repositoryRoot := copyFastStageRepository(t, buildFastStageTemplateRepository(t))
	if _, err := DecideFastStage(FastStageDecisionRequest{
		RepositoryRoot: repositoryRoot, ManifestPath: fastStageTestManifestPath,
		StageID: "alpha-stage", SuppliedArgv: []string{"sh", "module-alpha/run.sh", "--extra"},
	}); err == nil {
		t.Fatal("a decision was answered for a command the manifest does not declare")
	}
	if _, err := DecideFastStage(FastStageDecisionRequest{
		RepositoryRoot: repositoryRoot, ManifestPath: fastStageTestManifestPath,
		StageID: "absent-stage", SuppliedArgv: []string{"sh", "module-alpha/run.sh"},
	}); err == nil {
		t.Fatal("a decision was answered for a stage the manifest does not declare")
	}
	if _, err := DecideFastStage(FastStageDecisionRequest{
		RepositoryRoot: repositoryRoot, ManifestPath: "absent-manifest.json",
		StageID: "alpha-stage", SuppliedArgv: []string{"sh", "module-alpha/run.sh"},
	}); err == nil {
		t.Fatal("a decision was answered from a manifest that could not be read")
	}
}

// TestFastStageManifestDecodingRefusesAmbiguity keeps the manifest as strict as
// the heavy-lane one: an unknown field, a duplicate id, or a stage that covers
// nothing is a refusal rather than a silently narrower selection.
func TestFastStageManifestDecodingRefusesAmbiguity(t *testing.T) {
	for _, manifestCase := range []struct {
		name     string
		contents string
	}{
		{"unknown field", `{"schema_version": 1, "stages": [], "lanes": []}`},
		{"wrong schema version", `{"schema_version": 2, "stages": [{"id": "a", "argv": ["sh"], "coverage": [{"kind": "subtree", "path": "a"}]}]}`},
		{"no stages", `{"schema_version": 1, "stages": []}`},
		{"duplicate stage id", `{"schema_version": 1, "stages": [
			{"id": "a", "argv": ["sh"], "coverage": [{"kind": "subtree", "path": "a"}]},
			{"id": "a", "argv": ["sh"], "coverage": [{"kind": "subtree", "path": "b"}]}]}`},
		{"empty coverage", `{"schema_version": 1, "stages": [{"id": "a", "argv": ["sh"], "coverage": []}]}`},
		{"empty argv", `{"schema_version": 1, "stages": [{"id": "a", "argv": [], "coverage": [{"kind": "subtree", "path": "a"}]}]}`},
		{"unsupported coverage kind", `{"schema_version": 1, "stages": [{"id": "a", "argv": ["sh"], "coverage": [{"kind": "glob", "path": "a"}]}]}`},
		{"fingerprint without a probe", `{"schema_version": 1, "stages": [{"id": "a", "argv": ["sh"], "coverage": [{"kind": "subtree", "path": "a"}], "fingerprint": {"toolchain_probes": []}}]}`},
		{"trailing value", `{"schema_version": 1, "stages": [{"id": "a", "argv": ["sh"], "coverage": [{"kind": "subtree", "path": "a"}]}]} {}`},
	} {
		t.Run(manifestCase.name, func(t *testing.T) {
			if _, err := decodeFastStageManifest([]byte(manifestCase.contents)); err == nil {
				t.Fatalf("decoding accepted %s", manifestCase.name)
			}
		})
	}
}

// TestFastStageCoversNothingIsUncertain pins the last fail-closed edge: a stage
// whose declared coverage matches no path on disk has no inputs to compare, so
// no record can be trusted against it.
func TestFastStageCoversNothingIsUncertain(t *testing.T) {
	repositoryRoot := copyFastStageRepository(t, buildFastStageTemplateRepository(t))
	if err := os.RemoveAll(filepath.Join(repositoryRoot, "module-alpha")); err != nil {
		t.Fatal(err)
	}
	runHeavyTestGit(t, repositoryRoot, "add", "-A")
	if _, reason, _, _ := decideFastStageFields(t, repositoryRoot, "alpha-stage", time.Now()); reason != laneReasonFingerprintUncertain {
		t.Fatalf("a stage covering nothing reported %s, want fingerprint_uncertain", reason)
	}
}
