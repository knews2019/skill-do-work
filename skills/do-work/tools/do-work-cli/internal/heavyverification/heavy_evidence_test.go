//go:build unix

// Every test in this file executes a real lane whenever reuse is refused, and
// the runner fails closed on a platform that cannot prove descendant process
// ownership, so the whole file is scoped to the platforms where a lane can run.

package heavyverification

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/knews2019/skill-do-work/do-work-cli/internal/resultmodel"
)

// laneEvidenceTestManifest declares two independently fingerprinted lanes and
// one lane that declares no fingerprint inputs at all. Each fingerprinted lane
// probes its own toolchain file, which no lane covers, so a toolchain change
// and a covered-fixture change move the fingerprint by different routes.
const laneEvidenceTestManifest = `{
  "schema_version": 1,
  "lanes": [
    {
      "id": "alpha-lane",
      "argv": ["sh", "lanes/alpha.sh"],
      "coverage": [
        {"kind": "subtree", "path": "alpha"},
        {"kind": "exact", "path": "lanes/alpha.sh"}
      ],
      "fingerprint": {
        "toolchain_probes": [["cat", "toolchain/alpha-version.txt"]],
        "environment_variables": ["HEAVY_LANE_ALPHA_BROWSER"]
      }
    },
    {
      "id": "beta-lane",
      "argv": ["sh", "lanes/beta.sh"],
      "coverage": [
        {"kind": "subtree", "path": "beta"},
        {"kind": "exact", "path": "lanes/beta.sh"}
      ],
      "fingerprint": {
        "toolchain_probes": [["cat", "toolchain/beta-version.txt"]]
      }
    },
    {
      "id": "unfingerprinted-lane",
      "argv": ["sh", "lanes/unfingerprinted.sh"],
      "coverage": [{"kind": "exact", "path": "lanes/unfingerprinted.sh"}]
    },
    {
      "id": "red-lane",
      "argv": ["sh", "lanes/red.sh"],
      "coverage": [{"kind": "exact", "path": "lanes/red.sh"}],
      "fingerprint": {
        "toolchain_probes": [["cat", "toolchain/alpha-version.txt"]]
      }
    },
    {
      "id": "skip-lane",
      "argv": ["sh", "lanes/skip.sh"],
      "coverage": [{"kind": "exact", "path": "lanes/skip.sh"}],
      "fingerprint": {
        "toolchain_probes": [["cat", "toolchain/alpha-version.txt"]]
      }
    }
  ],
  "non_heavy_coverage": [{"kind": "subtree", "path": "docs"}]
}`

// newLaneEvidenceRepository commits lanes that each drop an ignored marker file
// when they execute, so a test can prove from the filesystem whether a lane ran
// rather than trusting the record the runner reports.
func newLaneEvidenceRepository(t *testing.T) string {
	t.Helper()
	repositoryRoot, _ := newHeavyTestRepository(t, laneEvidenceTestManifest)
	writeHeavyTestFile(t, repositoryRoot, ".gitignore", "*.ran\n")
	writeHeavyTestFile(t, repositoryRoot, "lanes/alpha.sh", "echo ran > alpha-lane.ran\nexit 0\n")
	writeHeavyTestFile(t, repositoryRoot, "lanes/beta.sh", "echo ran > beta-lane.ran\nexit 0\n")
	writeHeavyTestFile(t, repositoryRoot, "lanes/unfingerprinted.sh", "echo ran > unfingerprinted-lane.ran\nexit 0\n")
	writeHeavyTestFile(t, repositoryRoot, "lanes/red.sh", "echo ran > red-lane.ran\nexit 3\n")
	writeHeavyTestFile(t, repositoryRoot, "lanes/skip.sh", "echo ran > skip-lane.ran\nprintf 'SKIP: no engine here\\n'\nexit 0\n")
	writeHeavyTestFile(t, repositoryRoot, "alpha/fixture.txt", "alpha fixture\n")
	writeHeavyTestFile(t, repositoryRoot, "beta/fixture.txt", "beta fixture\n")
	writeHeavyTestFile(t, repositoryRoot, "toolchain/alpha-version.txt", "alpha toolchain 1\n")
	writeHeavyTestFile(t, repositoryRoot, "toolchain/beta-version.txt", "beta toolchain 1\n")
	commitHeavyTestChanges(t, repositoryRoot, "lane scripts, fixtures, and toolchain versions")
	return repositoryRoot
}

// verifyLanes runs the named lanes with reuse enabled at the given instant.
func verifyLanes(t *testing.T, repositoryRoot string, evaluatedAt time.Time, laneIDs ...string) resultmodel.HeavyVerificationRun {
	t.Helper()
	return verifyLanesWithReuse(t, repositoryRoot, evaluatedAt, true, laneIDs...)
}

func verifyLanesWithReuse(t *testing.T, repositoryRoot string, evaluatedAt time.Time, evidenceReuse bool, laneIDs ...string) resultmodel.HeavyVerificationRun {
	t.Helper()
	run, _, err := RunLanes(LaneRunRequest{
		RepositoryRoot: repositoryRoot, ManifestPath: "heavy-lanes.json", LaneIDs: laneIDs,
		LaneTimeout: 30 * time.Second, EvidenceReuse: evidenceReuse, EvaluatedAt: evaluatedAt,
	})
	if err != nil {
		t.Fatalf("run lanes %v: %v", laneIDs, err)
	}
	return run
}

// laneExecutionRecord finds one lane in a run so a test names the lane it means.
func laneExecutionRecord(t *testing.T, run resultmodel.HeavyVerificationRun, laneID string) resultmodel.HeavyLaneExecution {
	t.Helper()
	for _, lane := range run.Lanes {
		if lane.LaneID == laneID {
			return lane
		}
	}
	t.Fatalf("run has no lane %s: %#v", laneID, run.Lanes)
	return resultmodel.HeavyLaneExecution{}
}

// takeLaneRanMarker reports whether the lane's process actually ran since the
// last call, and clears the marker so the next assertion starts from nothing.
func takeLaneRanMarker(t *testing.T, repositoryRoot, laneID string) bool {
	t.Helper()
	markerPath := filepath.Join(repositoryRoot, laneID+".ran")
	if _, err := os.Stat(markerPath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		t.Fatalf("inspect lane marker %s: %v", markerPath, err)
	}
	if err := os.Remove(markerPath); err != nil {
		t.Fatalf("clear lane marker %s: %v", markerPath, err)
	}
	return true
}

func assertLaneDisposition(t *testing.T, lane resultmodel.HeavyLaneExecution, disposition, reason string) {
	t.Helper()
	if lane.Disposition != disposition || lane.DispositionReason != reason {
		t.Fatalf("lane %s disposition = %s/%s, want %s/%s", lane.LaneID, lane.Disposition, lane.DispositionReason, disposition, reason)
	}
}

// TestRunLanesReusesMatchingEvidenceWhileChangedLanesRerun is the request's
// central proof: within four hours an identical lane is reused and recorded as
// reused, while a lane whose covered fixture changed executes in the same pass.
func TestRunLanesReusesMatchingEvidenceWhileChangedLanesRerun(t *testing.T) {
	repositoryRoot := newLaneEvidenceRepository(t)
	recordedAt := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)

	firstRun := verifyLanes(t, repositoryRoot, recordedAt, "alpha-lane", "beta-lane")
	assertLaneDisposition(t, laneExecutionRecord(t, firstRun, "alpha-lane"), LaneDispositionExecuted, laneReasonNoPriorEvidence)
	assertLaneDisposition(t, laneExecutionRecord(t, firstRun, "beta-lane"), LaneDispositionExecuted, laneReasonNoPriorEvidence)
	if !takeLaneRanMarker(t, repositoryRoot, "alpha-lane") || !takeLaneRanMarker(t, repositoryRoot, "beta-lane") {
		t.Fatal("the first run did not execute both lanes")
	}

	reusedRun := verifyLanes(t, repositoryRoot, recordedAt.Add(time.Hour), "alpha-lane", "beta-lane")
	reusedAlpha := laneExecutionRecord(t, reusedRun, "alpha-lane")
	assertLaneDisposition(t, reusedAlpha, LaneDispositionReused, laneReasonFingerprintMatch)
	if reusedAlpha.ExitStatus != 0 || reusedAlpha.EvidenceRecordedAt != "2026-09-04T12:00:00Z" || reusedAlpha.EvidenceRevision == "" {
		t.Fatalf("reused alpha lane did not name the run it inherited: %#v", reusedAlpha)
	}
	if takeLaneRanMarker(t, repositoryRoot, "alpha-lane") || takeLaneRanMarker(t, repositoryRoot, "beta-lane") {
		t.Fatal("a lane executed even though its evidence was reported as reused")
	}

	// Only alpha's covered fixture changes, so only alpha's fingerprint moves.
	writeHeavyTestFile(t, repositoryRoot, "alpha/fixture.txt", "alpha fixture changed\n")
	commitHeavyTestChanges(t, repositoryRoot, "change one lane's fixture")

	partialRun := verifyLanes(t, repositoryRoot, recordedAt.Add(2*time.Hour), "alpha-lane", "beta-lane")
	assertLaneDisposition(t, laneExecutionRecord(t, partialRun, "alpha-lane"), LaneDispositionExecuted, laneReasonFingerprintMismatch)
	assertLaneDisposition(t, laneExecutionRecord(t, partialRun, "beta-lane"), LaneDispositionReused, laneReasonFingerprintMatch)
	if !takeLaneRanMarker(t, repositoryRoot, "alpha-lane") {
		t.Fatal("the lane whose fixture changed did not rerun")
	}
	if takeLaneRanMarker(t, repositoryRoot, "beta-lane") {
		t.Fatal("the unaffected lane reran instead of reusing its own evidence")
	}
}

// TestRunLanesRerunsExpiredEvidenceAndReuseNeverExtendsTheWindow pins the
// four-hour ceiling as a maximum age measured from the run that actually
// executed the lane: reusing at three hours must not restamp the record.
func TestRunLanesRerunsExpiredEvidenceAndReuseNeverExtendsTheWindow(t *testing.T) {
	repositoryRoot := newLaneEvidenceRepository(t)
	recordedAt := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)

	verifyLanes(t, repositoryRoot, recordedAt, "alpha-lane")
	takeLaneRanMarker(t, repositoryRoot, "alpha-lane")

	stillFreshRun := verifyLanes(t, repositoryRoot, recordedAt.Add(3*time.Hour), "alpha-lane")
	assertLaneDisposition(t, laneExecutionRecord(t, stillFreshRun, "alpha-lane"), LaneDispositionReused, laneReasonFingerprintMatch)

	// Nothing about the lane changed, so only age can force this rerun — and it
	// is measured against the original execution, not against the reuse above.
	expiredRun := verifyLanes(t, repositoryRoot, recordedAt.Add(4*time.Hour+time.Minute), "alpha-lane")
	expiredAlpha := laneExecutionRecord(t, expiredRun, "alpha-lane")
	assertLaneDisposition(t, expiredAlpha, LaneDispositionExecuted, laneReasonEvidenceExpired)
	if !takeLaneRanMarker(t, repositoryRoot, "alpha-lane") {
		t.Fatal("expired evidence was reported as executed but the lane never ran")
	}
	if expiredAlpha.FingerprintSHA256 == "" {
		t.Fatal("an expired lane still has a determinable fingerprint")
	}
}

// TestRunLanesRerunsWhenToolchainOrEnvironmentChanges proves the fingerprint
// covers more than the repository's own files: a lane whose covered bytes are
// identical still reruns when its toolchain or required environment moves.
func TestRunLanesRerunsWhenToolchainOrEnvironmentChanges(t *testing.T) {
	for _, testCase := range []struct {
		name           string
		changeTheInput func(t *testing.T, repositoryRoot string)
	}{
		{
			name: "toolchain probe output changes",
			changeTheInput: func(t *testing.T, repositoryRoot string) {
				writeHeavyTestFile(t, repositoryRoot, "toolchain/alpha-version.txt", "alpha toolchain 2\n")
				commitHeavyTestChanges(t, repositoryRoot, "upgrade the lane toolchain")
			},
		},
		{
			name: "required environment variable is set",
			changeTheInput: func(t *testing.T, repositoryRoot string) {
				t.Setenv("HEAVY_LANE_ALPHA_BROWSER", "/usr/bin/chromium")
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			repositoryRoot := newLaneEvidenceRepository(t)
			recordedAt := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
			firstRun := verifyLanes(t, repositoryRoot, recordedAt, "alpha-lane")
			takeLaneRanMarker(t, repositoryRoot, "alpha-lane")

			testCase.changeTheInput(t, repositoryRoot)

			secondRun := verifyLanes(t, repositoryRoot, recordedAt.Add(time.Minute), "alpha-lane")
			secondAlpha := laneExecutionRecord(t, secondRun, "alpha-lane")
			assertLaneDisposition(t, secondAlpha, LaneDispositionExecuted, laneReasonFingerprintMismatch)
			if secondAlpha.FingerprintSHA256 == laneExecutionRecord(t, firstRun, "alpha-lane").FingerprintSHA256 {
				t.Fatal("the fingerprint did not move when the lane's declared input changed")
			}
			if !takeLaneRanMarker(t, repositoryRoot, "alpha-lane") {
				t.Fatal("the lane was reported as executed but never ran")
			}
		})
	}
}

// TestRunLanesNeverReusesWithoutADeterminableFingerprint covers the fail-closed
// half of the contract: a lane that declares no fingerprint inputs has an
// undetermined toolchain, so a recent green never authorizes reuse.
func TestRunLanesNeverReusesWithoutADeterminableFingerprint(t *testing.T) {
	repositoryRoot := newLaneEvidenceRepository(t)
	recordedAt := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)

	verifyLanes(t, repositoryRoot, recordedAt, "unfingerprinted-lane")
	if !takeLaneRanMarker(t, repositoryRoot, "unfingerprinted-lane") {
		t.Fatal("the first run did not execute the lane")
	}

	secondRun := verifyLanes(t, repositoryRoot, recordedAt.Add(time.Minute), "unfingerprinted-lane")
	secondLane := laneExecutionRecord(t, secondRun, "unfingerprinted-lane")
	assertLaneDisposition(t, secondLane, LaneDispositionExecuted, laneReasonFingerprintUncertain)
	if secondLane.FingerprintSHA256 != "" {
		t.Fatalf("an undetermined fingerprint was reported anyway: %#v", secondLane)
	}
	if !takeLaneRanMarker(t, repositoryRoot, "unfingerprinted-lane") {
		t.Fatal("the lane with no fingerprint inputs was skipped by age alone")
	}
}

// TestRunLanesStoresNoEvidenceForRedOrSkippedLanes keeps a failure and a lane
// that announced it did not run out of the cache: neither is a result a later
// run may inherit.
func TestRunLanesStoresNoEvidenceForRedOrSkippedLanes(t *testing.T) {
	repositoryRoot := newLaneEvidenceRepository(t)
	recordedAt := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)

	verifyLanes(t, repositoryRoot, recordedAt, "red-lane", "skip-lane")
	takeLaneRanMarker(t, repositoryRoot, "red-lane")
	takeLaneRanMarker(t, repositoryRoot, "skip-lane")

	secondRun := verifyLanes(t, repositoryRoot, recordedAt.Add(time.Minute), "red-lane", "skip-lane")
	assertLaneDisposition(t, laneExecutionRecord(t, secondRun, "red-lane"), LaneDispositionExecuted, laneReasonNoPriorEvidence)
	assertLaneDisposition(t, laneExecutionRecord(t, secondRun, "skip-lane"), LaneDispositionExecuted, laneReasonNoPriorEvidence)
	if !takeLaneRanMarker(t, repositoryRoot, "red-lane") || !takeLaneRanMarker(t, repositoryRoot, "skip-lane") {
		t.Fatal("a red or skipped lane was reused instead of rerun")
	}
}

// TestRunLanesWithoutEvidenceReuseExecutesAndStillRefreshesTheRecord pins the
// override: --no-evidence-reuse forces execution, and the record it writes is
// available to the next run that does allow reuse.
func TestRunLanesWithoutEvidenceReuseExecutesAndStillRefreshesTheRecord(t *testing.T) {
	repositoryRoot := newLaneEvidenceRepository(t)
	recordedAt := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)

	forcedRun := verifyLanesWithReuse(t, repositoryRoot, recordedAt, false, "alpha-lane")
	assertLaneDisposition(t, laneExecutionRecord(t, forcedRun, "alpha-lane"), LaneDispositionExecuted, laneReasonReuseDisabled)
	if !takeLaneRanMarker(t, repositoryRoot, "alpha-lane") {
		t.Fatal("the forced run did not execute the lane")
	}

	reusingRun := verifyLanes(t, repositoryRoot, recordedAt.Add(time.Hour), "alpha-lane")
	assertLaneDisposition(t, laneExecutionRecord(t, reusingRun, "alpha-lane"), LaneDispositionReused, laneReasonFingerprintMatch)
	if takeLaneRanMarker(t, repositoryRoot, "alpha-lane") {
		t.Fatal("the record written by the forced run was not reusable")
	}
}

// TestRunHeavyVerificationCommandReusesEvidenceThroughItsOwnSeam exercises the
// command a drain actually invokes, not just RunLanes: a second invocation with
// no extra flags reuses, and --no-evidence-reuse forces the lane to execute.
func TestRunHeavyVerificationCommandReusesEvidenceThroughItsOwnSeam(t *testing.T) {
	repositoryRoot := newLaneEvidenceRepository(t)

	firstResult := runHeavyLanes(t, repositoryRoot, "--lane", "alpha-lane")
	if firstResult.Outcome != resultmodel.OutcomeSuccess || firstResult.HeavyVerificationRun == nil {
		t.Fatalf("first command result = %#v", firstResult)
	}
	assertLaneDisposition(t, laneExecutionRecord(t, *firstResult.HeavyVerificationRun, "alpha-lane"), LaneDispositionExecuted, laneReasonNoPriorEvidence)
	if !takeLaneRanMarker(t, repositoryRoot, "alpha-lane") {
		t.Fatal("the first command invocation did not execute the lane")
	}

	secondResult := runHeavyLanes(t, repositoryRoot, "--lane", "alpha-lane")
	assertLaneDisposition(t, laneExecutionRecord(t, *secondResult.HeavyVerificationRun, "alpha-lane"), LaneDispositionReused, laneReasonFingerprintMatch)
	if takeLaneRanMarker(t, repositoryRoot, "alpha-lane") {
		t.Fatal("the command executed a lane it reported as reused")
	}

	forcedResult := runHeavyLanes(t, repositoryRoot, "--lane", "alpha-lane", "--no-evidence-reuse")
	assertLaneDisposition(t, laneExecutionRecord(t, *forcedResult.HeavyVerificationRun, "alpha-lane"), LaneDispositionExecuted, laneReasonReuseDisabled)
	if !takeLaneRanMarker(t, repositoryRoot, "alpha-lane") {
		t.Fatal("--no-evidence-reuse did not force the lane to execute")
	}
}

// TestRunLanesRefusesTamperedEvidenceInsteadOfTrustingIt keeps a stored record
// that no longer describes a successful run out of the reuse path.
func TestRunLanesRefusesTamperedEvidenceInsteadOfTrustingIt(t *testing.T) {
	repositoryRoot := newLaneEvidenceRepository(t)
	recordedAt := time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC)
	verifyLanes(t, repositoryRoot, recordedAt, "alpha-lane")
	takeLaneRanMarker(t, repositoryRoot, "alpha-lane")

	store, err := openLaneEvidenceStore(repositoryRoot)
	if err != nil {
		t.Fatal(err)
	}
	record, exists, err := store.ReadLaneEvidence("alpha-lane")
	if err != nil || !exists {
		t.Fatalf("stored alpha evidence: exists=%t err=%v", exists, err)
	}
	record.ExitStatus = 1
	if err := store.WriteLaneEvidence(record); err != nil {
		t.Fatal(err)
	}

	secondRun := verifyLanes(t, repositoryRoot, recordedAt.Add(time.Minute), "alpha-lane")
	assertLaneDisposition(t, laneExecutionRecord(t, secondRun, "alpha-lane"), LaneDispositionExecuted, laneReasonEvidenceUnusable)
	if !takeLaneRanMarker(t, repositoryRoot, "alpha-lane") {
		t.Fatal("a record that no longer describes a green run was reused")
	}
}
