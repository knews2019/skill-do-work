package main

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

// The disk-space probe is tested through fake measurements only: real free space
// changes under the test, so a real-disk assertion could never be deterministic.
// TestMain (generate_test.go) installs a plenty-free fake for every other test.

const fixtureDiskTotalBytes = 500 << 30

const expectedDiskSpaceRemedy = "free space: clear regenerable QA output, finished builder worktrees (do-work cleanup), browser caches; " +
	"`du -sh * | sort -h` at the repo root shows the largest directories"

// fakeDiskMeasurement is one fake answer for one directory: the measurement, or
// the error the measurer returns instead.
type fakeDiskMeasurement struct {
	measurement      diskSpaceMeasurement
	measurementError error
}

// fakeDiskSpaceMeasurer answers from a fixed table and fails the test on any
// directory the table does not name, so a probe that measures the wrong path
// cannot pass by accident.
func fakeDiskSpaceMeasurer(t *testing.T, answersByDirectory map[string]fakeDiskMeasurement) func(string) (diskSpaceMeasurement, error) {
	t.Helper()
	return func(directory string) (diskSpaceMeasurement, error) {
		answer, isKnown := answersByDirectory[directory]
		if !isKnown {
			t.Errorf("the probe measured an unexpected directory %q", directory)
			return diskSpaceMeasurement{}, fmt.Errorf("unexpected directory %q", directory)
		}
		return answer.measurement, answer.measurementError
	}
}

func freeOnDevice(freeBytes uint64, deviceIdentity uint64) fakeDiskMeasurement {
	return fakeDiskMeasurement{measurement: diskSpaceMeasurement{
		freeBytes: freeBytes, totalBytes: fixtureDiskTotalBytes, deviceIdentity: deviceIdentity,
	}}
}

// Pins "probe silently clean when the disk is low": the repo root below each
// threshold must be reported at that threshold's level, and above both must not.
func TestDiskSpaceProbeReportsLowFreeSpaceAtEachThreshold(t *testing.T) {
	repoRoot := "/fixture/repo"
	cases := []struct {
		name          string
		freeBytes     uint64
		wantFinding   bool
		wantThreshold string
	}{
		{"2 GiB free is critical", 2 << 30, true, "below the critical threshold (3.0 GiB)"},
		{"9 GiB free is a warning", 9 << 30, true, "below the warning threshold (10.0 GiB)"},
		{"20 GiB free is clean", 20 << 30, false, ""},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			report := VerifyReport{}
			measurer := fakeDiskSpaceMeasurer(t, map[string]fakeDiskMeasurement{repoRoot: freeOnDevice(testCase.freeBytes, 1)})
			appendDiskSpaceFindings(&report, repoRoot, nil, measurer)

			findings := findingsMentioning(report, verifyCategoryLowDiskSpace)
			if !testCase.wantFinding {
				if len(findings) != 0 {
					t.Fatalf("expected no low-disk-space finding, got %+v", findings)
				}
				return
			}
			if len(findings) != 1 {
				t.Fatalf("expected exactly one low-disk-space finding, got %+v", report.Findings)
			}
			finding := findings[0]
			if finding.Subject != repoRoot {
				t.Errorf("Subject = %q, want the repo root %q", finding.Subject, repoRoot)
			}
			if finding.Fixable {
				t.Error("a low-disk-space finding was advertised as fixable; do-work cleanup cannot free disk space")
			}
			if !strings.Contains(finding.Detail, testCase.wantThreshold) {
				t.Errorf("Detail %q does not name %q", finding.Detail, testCase.wantThreshold)
			}
			if strings.HasPrefix(finding.Detail, finding.Subject) {
				t.Errorf("Detail %q repeats its Subject", finding.Detail)
			}
			if finding.Remedy != expectedDiskSpaceRemedy {
				t.Errorf("Remedy = %q, want %q", finding.Remedy, expectedDiskSpaceRemedy)
			}
		})
	}
}

// Pins an off-by-one at the constants: "below" is strict, so exactly 3 GiB is
// still only a warning and exactly 10 GiB is clean.
func TestDiskSpaceProbeThresholdBoundariesAreStrict(t *testing.T) {
	repoRoot := "/fixture/repo"

	atCritical := VerifyReport{}
	appendDiskSpaceFindings(&atCritical, repoRoot, nil,
		fakeDiskSpaceMeasurer(t, map[string]fakeDiskMeasurement{repoRoot: freeOnDevice(lowDiskSpaceCriticalBytes, 1)}))
	criticalFindings := findingsMentioning(atCritical, verifyCategoryLowDiskSpace)
	if len(criticalFindings) != 1 || !strings.Contains(criticalFindings[0].Detail, "warning threshold") {
		t.Errorf("exactly 3 GiB free should be one warning-level finding, got %+v", atCritical.Findings)
	}

	atWarning := VerifyReport{}
	appendDiskSpaceFindings(&atWarning, repoRoot, nil,
		fakeDiskSpaceMeasurer(t, map[string]fakeDiskMeasurement{repoRoot: freeOnDevice(lowDiskSpaceWarningBytes, 1)}))
	if len(atWarning.Findings) != 0 {
		t.Errorf("exactly 10 GiB free should be clean, got %+v", atWarning.Findings)
	}
	if lowDiskSpaceCriticalBytes != 3<<30 || lowDiskSpaceWarningBytes != 10<<30 {
		t.Errorf("thresholds moved: critical %d, warning %d", lowDiskSpaceCriticalBytes, lowDiskSpaceWarningBytes)
	}
}

// Pins duplicate findings per device: two worktrees on one filesystem are one
// disk, so they are one finding, named after the first worktree in sorted order.
func TestDiskSpaceProbeReportsOneFindingPerDevice(t *testing.T) {
	repoRoot := "/fixture/repo"
	worktreePathsByName := map[string]string{
		"worktree-agent-REQ-902": "/fixture/worktrees/worktree-agent-REQ-902",
		"worktree-agent-REQ-901": "/fixture/worktrees/worktree-agent-REQ-901",
	}
	report := VerifyReport{}
	appendDiskSpaceFindings(&report, repoRoot, worktreePathsByName, fakeDiskSpaceMeasurer(t, map[string]fakeDiskMeasurement{
		repoRoot: freeOnDevice(200<<30, 1),
		"/fixture/worktrees/worktree-agent-REQ-901": freeOnDevice(2<<30, 2),
		"/fixture/worktrees/worktree-agent-REQ-902": freeOnDevice(2<<30, 2),
	}))

	findings := findingsMentioning(report, verifyCategoryLowDiskSpace)
	if len(findings) != 1 {
		t.Fatalf("two worktrees on one device should yield one finding, got %+v", report.Findings)
	}
	if findings[0].Subject != "worktree-agent-REQ-901" {
		t.Errorf("Subject = %q, want the first worktree name, never its path", findings[0].Subject)
	}
}

// Pins "unknown reads as clean": a platform the probe cannot measure on is a
// skipped probe, never a quiet pass.
func TestDiskSpaceProbeUnsupportedPlatformIsSkippedNotClean(t *testing.T) {
	repoRoot := "/fixture/repo"
	report := VerifyReport{}
	unsupportedError := fmt.Errorf("measuring: %w", errDiskSpaceUnsupported)
	appendDiskSpaceFindings(&report, repoRoot, map[string]string{"worktree-agent-REQ-901": "/fixture/worktrees/worktree-agent-REQ-901"},
		func(string) (diskSpaceMeasurement, error) { return diskSpaceMeasurement{}, unsupportedError })

	if len(report.Findings) != 0 {
		t.Errorf("an unsupported platform produced findings: %+v", report.Findings)
	}
	wantSkip := "disk-space probe: unsupported on " + runtime.GOOS
	if len(report.SkippedProbes) != 1 || report.SkippedProbes[0] != wantSkip {
		t.Errorf("SkippedProbes = %q, want exactly [%q]", report.SkippedProbes, wantSkip)
	}
}

// Pins a crash or a lost probe on one unreadable directory: that directory is a
// skip, and every other directory is still measured and reported.
func TestDiskSpaceProbeMeasurementErrorSkipsOnlyThatDirectory(t *testing.T) {
	repoRoot := "/fixture/repo"
	brokenWorktreePath := "/fixture/worktrees/worktree-agent-REQ-901"
	report := VerifyReport{}
	appendDiskSpaceFindings(&report, repoRoot, map[string]string{"worktree-agent-REQ-901": brokenWorktreePath},
		fakeDiskSpaceMeasurer(t, map[string]fakeDiskMeasurement{
			repoRoot:           freeOnDevice(2<<30, 1),
			brokenWorktreePath: {measurementError: errors.New("permission denied")},
		}))

	findings := findingsMentioning(report, verifyCategoryLowDiskSpace)
	if len(findings) != 1 || findings[0].Subject != repoRoot {
		t.Errorf("the repo root should still be reported, got %+v", report.Findings)
	}
	wantSkip := "disk-space probe for " + brokenWorktreePath + ": permission denied"
	if len(report.SkippedProbes) != 1 || report.SkippedProbes[0] != wantSkip {
		t.Errorf("SkippedProbes = %q, want exactly [%q]", report.SkippedProbes, wantSkip)
	}
}

// Pins the probe being unwired: the REQ's red case runs the real verify entry
// point over a clean fixture with the package measurer faked low.
func TestDiskSpaceProbeReachesTheVerifyReport(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
	})
	previousMeasurer := diskSpaceMeasurer
	t.Cleanup(func() { diskSpaceMeasurer = previousMeasurer })

	diskSpaceMeasurer = func(string) (diskSpaceMeasurement, error) {
		return diskSpaceMeasurement{freeBytes: 2 << 30, totalBytes: fixtureDiskTotalBytes, deviceIdentity: 1}, nil
	}
	report, verifyError := runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	findings := findingsMentioning(report, verifyCategoryLowDiskSpace)
	if len(findings) != 1 || findings[0].Subject != repoRoot || !strings.Contains(findings[0].Detail, "critical threshold") {
		t.Errorf("expected one critical low-disk-space finding for the repo root, got %+v", report.Findings)
	}

	diskSpaceMeasurer = func(string) (diskSpaceMeasurement, error) { return diskSpaceMeasurement{}, errDiskSpaceUnsupported }
	report, verifyError = runVerifyProbes(repoRoot, time.Now())
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	unsupportedSkips := 0
	for _, skipped := range report.SkippedProbes {
		if skipped == "disk-space probe: unsupported on "+runtime.GOOS {
			unsupportedSkips++
		}
	}
	if unsupportedSkips != 1 {
		t.Errorf("expected one unsupported disk-space skip, got %q", report.SkippedProbes)
	}
}
