//go:build unix

package heavyverification

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFailedForcedRunRevokesPriorSuccess(t *testing.T) {
	for _, outcome := range []string{"exit 3", "echo 'SKIP: unavailable'; exit 0"} {
		t.Run(outcome, func(t *testing.T) {
			root := newLaneEvidenceRepository(t)
			writeHeavyTestFile(t, root, "lanes/alpha.sh", "echo ran > alpha-lane.ran\nif [ -f failure.ran ]; then "+outcome+"; fi\n")
			commitHeavyTestChanges(t, root, "conditional lane failure")
			now := time.Now()
			verifyLanes(t, root, now, "alpha-lane")
			writeHeavyTestFile(t, root, "failure.ran", "fail")
			verifyLanesWithReuse(t, root, now.Add(time.Minute), false, "alpha-lane")
			os.Remove(filepath.Join(root, "failure.ran"))
			takeLaneRanMarker(t, root, "alpha-lane")
			result := verifyLanes(t, root, now.Add(2*time.Minute), "alpha-lane")
			if laneExecutionRecord(t, result, "alpha-lane").Disposition != LaneDispositionExecuted || !takeLaneRanMarker(t, root, "alpha-lane") {
				t.Fatal("failed/skipped forced run left prior successful evidence reusable")
			}
		})
	}
}

func TestEachLaneChecksExpiryAtItsOwnDecision(t *testing.T) {
	root := newLaneEvidenceRepository(t)
	writeHeavyTestFile(t, root, "lanes/alpha.sh", "sleep 2\necho ran > alpha-lane.ran\n")
	commitHeavyTestChanges(t, root, "slow preceding lane")
	verifyLanes(t, root, time.Now(), "beta-lane")
	takeLaneRanMarker(t, root, "beta-lane")
	store, err := openLaneEvidenceStore(root)
	if err != nil {
		t.Fatal(err)
	}
	record, _, err := store.ReadLaneEvidence("beta-lane")
	if err != nil {
		t.Fatal(err)
	}
	record.RecordedAt = time.Now().Truncate(time.Second).Add(-4*time.Hour + 2*time.Second).UTC().Format(time.RFC3339)
	if err := store.WriteLaneEvidence(record); err != nil {
		t.Fatal(err)
	}
	result := verifyLanes(t, root, time.Time{}, "alpha-lane", "beta-lane")
	if laneExecutionRecord(t, result, "beta-lane").Disposition != LaneDispositionExecuted || !takeLaneRanMarker(t, root, "beta-lane") {
		t.Fatal("later lane reused evidence older than four hours at its decision")
	}
}

func TestUnknownChangedPathInvalidatesAllLaneEvidence(t *testing.T) {
	root := newLaneEvidenceRepository(t)
	now := time.Now()
	verifyLanes(t, root, now, "alpha-lane", "beta-lane")
	writeHeavyTestFile(t, root, "unknown-helper.sh", "changed runtime input\n")
	commitHeavyTestChanges(t, root, "unclassified input")
	result := verifyLanes(t, root, now.Add(time.Minute), "alpha-lane", "beta-lane")
	for _, lane := range result.Lanes {
		if lane.Disposition != LaneDispositionExecuted {
			t.Fatalf("unknown path did not invalidate %s", lane.LaneID)
		}
	}
}

func TestUndeclaredEnvironmentChangesInvalidateEvidence(t *testing.T) {
	root := newLaneEvidenceRepository(t)
	now := time.Now()
	verifyLanes(t, root, now, "alpha-lane")
	t.Setenv("DO_WORK_GO_TEST_EXCLUDE_PREFIXES", "TestSomething")
	result := verifyLanes(t, root, now.Add(time.Minute), "alpha-lane")
	if laneExecutionRecord(t, result, "alpha-lane").Disposition != LaneDispositionExecuted {
		t.Fatal("an undeclared inherited environment input did not invalidate evidence")
	}
}

func TestInvalidationFailureRefusesBeforeExecuting(t *testing.T) {
	root := newLaneEvidenceRepository(t)
	now := time.Now()
	verifyLanes(t, root, now, "alpha-lane")
	takeLaneRanMarker(t, root, "alpha-lane")
	store, err := openLaneEvidenceStore(root)
	if err != nil {
		t.Fatal(err)
	}
	recordPath := store.recordPath("alpha-lane")
	if err := os.Remove(recordPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(recordPath, 0700); err != nil {
		t.Fatal(err)
	}
	_, _, err = RunLanes(LaneRunRequest{RepositoryRoot: root, ManifestPath: "heavy-lanes.json", LaneIDs: []string{"alpha-lane"}, LaneTimeout: time.Second})
	if err == nil || LaneRunRefusalCode(err) != "HEAVY-RUN-EVIDENCE-INVALIDATION" {
		t.Fatalf("invalidation refusal = %v", err)
	}
	if takeLaneRanMarker(t, root, "alpha-lane") {
		t.Fatal("lane ran despite failed evidence invalidation")
	}
}

func TestFingerprintProbeTimeoutOwnsDescendants(t *testing.T) {
	root := t.TempDir()
	startedAt := time.Now()
	_, err := runFingerprintProbe(root, []string{"sh", "-c", "sleep 10 & wait"}, 50*time.Millisecond)
	if err == nil {
		t.Fatal("hung probe was accepted")
	}
	if time.Since(startedAt) > 4*time.Second {
		t.Fatal("probe exceeded its bound plus termination grace")
	}
}

func TestUntrackedLaneInputsCannotReuseCommittedEvidence(t *testing.T) {
	for _, ignored := range []bool{false, true} {
		t.Run(fmt.Sprintf("ignored=%t", ignored), func(t *testing.T) {
			root := newLaneEvidenceRepository(t)
			if ignored {
				writeHeavyTestFile(t, root, ".gitignore", "*.ran\nalpha/untracked_test.go\n")
				commitHeavyTestChanges(t, root, "ignored source fixture")
			}
			now := time.Now()
			verifyLanes(t, root, now, "alpha-lane", "beta-lane")
			writeHeavyTestFile(t, root, "alpha/untracked_test.go", "package alpha\n")
			result := verifyLanes(t, root, now.Add(time.Minute), "alpha-lane", "beta-lane")
			assertLaneDisposition(t, laneExecutionRecord(t, result, "alpha-lane"), LaneDispositionExecuted, laneReasonFingerprintMismatch)
			// The unrelated lane may still reuse its complete committed inputs.
			assertLaneDisposition(t, laneExecutionRecord(t, result, "beta-lane"), LaneDispositionReused, laneReasonFingerprintMatch)
		})
	}
}

func TestTrackedSymlinkInputsRemainUncertain(t *testing.T) {
	root := newLaneEvidenceRepository(t)
	externalPath := filepath.Join(t.TempDir(), "external.go")
	if err := os.WriteFile(externalPath, []byte("package alpha\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(externalPath, filepath.Join(root, "alpha/linked.go")); err != nil {
		t.Fatal(err)
	}
	commitHeavyTestChanges(t, root, "tracked source symlink")
	result := verifyLanes(t, root, time.Now(), "alpha-lane")
	assertLaneDisposition(t, laneExecutionRecord(t, result, "alpha-lane"), LaneDispositionExecuted, laneReasonFingerprintUncertain)
}
