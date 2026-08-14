package main

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestChurnRenameAttribution pins requirement: pre-rename touches count toward
// the file's CURRENT path, not the dead one. History built: add alpha → edit
// alpha → rename alpha→beta → edit beta. All four commits must land on
// beta.md, and alpha.md must not appear at all.
func TestChurnRenameAttribution(t *testing.T) {
	repoRoot := newFixtureRepo(t)
	writeFixtureFile(t, repoRoot, "alpha.md", "first draft of the alpha file\n")
	commitFixtureAll(t, repoRoot, "add alpha")
	writeFixtureFile(t, repoRoot, "alpha.md", "first draft of the alpha file\nplus an edit\n")
	commitFixtureAll(t, repoRoot, "edit alpha")
	runGitInFixture(t, repoRoot, "mv", "alpha.md", "beta.md")
	commitFixtureAll(t, repoRoot, "rename alpha to beta")
	writeFixtureFile(t, repoRoot, "beta.md", "first draft of the alpha file\nplus an edit\nplus another\n")
	commitFixtureAll(t, repoRoot, "edit beta")

	report, computeError := computeChurnReport(repoRoot, "10 years", nil)
	if computeError != nil {
		t.Fatalf("computeChurnReport: %v", computeError)
	}
	if report.ShallowClone {
		t.Fatalf("full local fixture reported as shallow")
	}
	if got := report.TouchesByPath["beta.md"]; got != 4 {
		t.Fatalf("beta.md touches = %d, want 4 (2 alpha-era + rename + edit); full map: %v", got, report.TouchesByPath)
	}
	if _, deadPathPresent := report.TouchesByPath["alpha.md"]; deadPathPresent {
		t.Fatalf("dead path alpha.md still present: %v", report.TouchesByPath)
	}
}

// TestChurnStagedCopyMigrationAttribution pins the staged-migration rule the
// skills/ restructure needed: when a file is COPIED to its new home and the
// original is deleted in a LATER commit (so -M never sees a rename), the dead
// original's whole history — pre-copy and between copy and delete — still
// lands on the surviving copy.
func TestChurnStagedCopyMigrationAttribution(t *testing.T) {
	repoRoot := newFixtureRepo(t)
	writeFixtureFile(t, repoRoot, "old/home.md", "the file's original content, long enough to match\n")
	commitFixtureAll(t, repoRoot, "add original")
	writeFixtureFile(t, repoRoot, "new/home.md", "the file's original content, long enough to match\n")
	commitFixtureAll(t, repoRoot, "stage copy")
	writeFixtureFile(t, repoRoot, "old/home.md", "the file's original content, long enough to match\nedited during staging\n")
	commitFixtureAll(t, repoRoot, "edit original during staging")
	runGitInFixture(t, repoRoot, "rm", "--quiet", "old/home.md")
	commitFixtureAll(t, repoRoot, "retire original")
	writeFixtureFile(t, repoRoot, "new/home.md", "the file's original content, long enough to match\nedited after cutover\n")
	commitFixtureAll(t, repoRoot, "edit survivor")

	report, computeError := computeChurnReport(repoRoot, "10 years", nil)
	if computeError != nil {
		t.Fatalf("computeChurnReport: %v", computeError)
	}
	// add original + stage copy + edit original + edit survivor = 4; the
	// deletion is not a touch.
	if got := report.TouchesByPath["new/home.md"]; got != 4 {
		t.Fatalf("new/home.md touches = %d, want 4; full map: %v", got, report.TouchesByPath)
	}
	if _, deadPathPresent := report.TouchesByPath["old/home.md"]; deadPathPresent {
		t.Fatalf("dead path old/home.md still present: %v", report.TouchesByPath)
	}
}

// TestChurnDeletedPathDropped pins the other half of the rename rule: only
// paths deleted OUTRIGHT are dropped — their history attaches to nothing.
func TestChurnDeletedPathDropped(t *testing.T) {
	repoRoot := newFixtureRepo(t)
	writeFixtureFile(t, repoRoot, "gone.md", "short-lived\n")
	writeFixtureFile(t, repoRoot, "stay.md", "long-lived\n")
	commitFixtureAll(t, repoRoot, "add both")
	runGitInFixture(t, repoRoot, "rm", "--quiet", "gone.md")
	commitFixtureAll(t, repoRoot, "delete gone")

	report, computeError := computeChurnReport(repoRoot, "10 years", nil)
	if computeError != nil {
		t.Fatalf("computeChurnReport: %v", computeError)
	}
	if _, deadPathPresent := report.TouchesByPath["gone.md"]; deadPathPresent {
		t.Fatalf("outright-deleted gone.md still present: %v", report.TouchesByPath)
	}
	if got := report.TouchesByPath["stay.md"]; got != 1 {
		t.Fatalf("stay.md touches = %d, want 1; full map: %v", got, report.TouchesByPath)
	}
}

// TestChurnExcludePathHonored pins the ceremony-exclude contract: an excluded
// prefix drops the path from churn entirely (the release ritual's synchronized
// files would otherwise top every churn table).
func TestChurnExcludePathHonored(t *testing.T) {
	repoRoot := newFixtureRepo(t)
	writeFixtureFile(t, repoRoot, "CHANGELOG.md", "## 1.0.0\n")
	writeFixtureFile(t, repoRoot, "code.md", "real work\n")
	commitFixtureAll(t, repoRoot, "seed")

	report, computeError := computeChurnReport(repoRoot, "10 years", []string{"CHANGELOG.md"})
	if computeError != nil {
		t.Fatalf("computeChurnReport: %v", computeError)
	}
	if _, excludedPresent := report.TouchesByPath["CHANGELOG.md"]; excludedPresent {
		t.Fatalf("excluded CHANGELOG.md still present: %v", report.TouchesByPath)
	}
	if got := report.TouchesByPath["code.md"]; got != 1 {
		t.Fatalf("code.md touches = %d, want 1; full map: %v", got, report.TouchesByPath)
	}
}

// TestChurnShallowCloneReported pins requirement: a shallow clone is DETECTED
// and visibly reported, never silently truncated. The fixture is a real
// --depth 1 clone over the file:// transport (a plain local-path clone ignores
// --depth and stays full).
func TestChurnShallowCloneReported(t *testing.T) {
	sourceRoot := newFixtureRepo(t)
	writeFixtureFile(t, sourceRoot, "history.md", "v1\n")
	commitFixtureAll(t, sourceRoot, "first")
	writeFixtureFile(t, sourceRoot, "history.md", "v1\nv2\n")
	commitFixtureAll(t, sourceRoot, "second")

	cloneParent := t.TempDir()
	shallowRoot := filepath.Join(cloneParent, "shallow-clone")
	runGitInFixture(t, cloneParent, "clone", "--quiet", "--depth", "1",
		"file://"+filepath.ToSlash(sourceRoot), shallowRoot)

	report, computeError := computeChurnReport(shallowRoot, "10 years", nil)
	if computeError != nil {
		t.Fatalf("computeChurnReport on shallow clone: %v", computeError)
	}
	if !report.ShallowClone {
		t.Fatalf("shallow clone not detected")
	}

	var renderedOutput strings.Builder
	writeChurnReport(&renderedOutput, report, 10)
	if !strings.Contains(renderedOutput.String(), "WARNING: shallow clone detected") {
		t.Fatalf("churn output missing the shallow warning:\n%s", renderedOutput.String())
	}
}

// TestBandLabelForValueEdges pins the band-edge rule at its source: value ==
// threshold is NOT flagged, strictly greater is, and FLAG outranks WATCH.
// (Lives here rather than a renderer test so the rule holds for every caller.)
func TestBandLabelForValueEdges(t *testing.T) {
	testCases := []struct {
		name           string
		value          int
		watchThreshold int
		flagThreshold  int
		want           string
	}{
		{"equal to watch is not flagged", 10, 10, bandThresholdUnset, ""},
		{"one over watch is WATCH", 11, 10, bandThresholdUnset, "WATCH"},
		{"equal to flag stays WATCH", 20, 10, 20, "WATCH"},
		{"one over flag is FLAG", 21, 10, 20, "FLAG"},
		{"no thresholds means no band", 1000, bandThresholdUnset, bandThresholdUnset, ""},
		{"flag alone works without watch", 21, bandThresholdUnset, 20, "FLAG"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := bandLabelForValue(testCase.value, testCase.watchThreshold, testCase.flagThreshold)
			if got != testCase.want {
				t.Fatalf("bandLabelForValue(%d, %d, %d) = %q, want %q",
					testCase.value, testCase.watchThreshold, testCase.flagThreshold, got, testCase.want)
			}
		})
	}
}
