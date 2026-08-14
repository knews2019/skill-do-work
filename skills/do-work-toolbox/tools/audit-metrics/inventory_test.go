package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Fixture helpers: git-behavior tests need a REAL git repo — plain temp dirs
// make git probes silently skip or fail (queue-kanban lesson REQ-083). Every
// git identity detail is pinned in the fixture so the tests pass on a bare CI
// container with no global config.

// runGitInFixture runs one git command while building a fixture. Setup only —
// the code under test does its own shelling out.
func runGitInFixture(t *testing.T, workingDirectory string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", arguments...)
	command.Dir = workingDirectory
	if output, runError := command.CombinedOutput(); runError != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(arguments, " "), runError, output)
	}
}

// newFixtureRepo initializes an empty real git repo with a pinned identity and
// default branch. Skips when git is unavailable, matching how the tool itself
// degrades (it cannot run at all without git).
func newFixtureRepo(t *testing.T) string {
	t.Helper()
	if _, lookupError := exec.LookPath("git"); lookupError != nil {
		t.Skip("git is not on PATH — audit-metrics cannot run without it")
	}
	repoRoot := t.TempDir()
	runGitInFixture(t, repoRoot, "init", "--quiet", "--initial-branch=main")
	runGitInFixture(t, repoRoot, "config", "user.email", "fixture@example.test")
	runGitInFixture(t, repoRoot, "config", "user.name", "Audit Fixture")
	return repoRoot
}

// writeFixtureFile writes one file (creating parent directories) inside a
// fixture repo.
func writeFixtureFile(t *testing.T, repoRoot string, relativePath string, content string) {
	t.Helper()
	fullPath := filepath.Join(repoRoot, filepath.FromSlash(relativePath))
	if mkdirError := os.MkdirAll(filepath.Dir(fullPath), 0o755); mkdirError != nil {
		t.Fatalf("mkdir for %s: %v", relativePath, mkdirError)
	}
	if writeError := os.WriteFile(fullPath, []byte(content), 0o644); writeError != nil {
		t.Fatalf("write %s: %v", relativePath, writeError)
	}
}

// commitFixtureAll stages and commits everything in the fixture repo.
func commitFixtureAll(t *testing.T, repoRoot string, message string) {
	t.Helper()
	runGitInFixture(t, repoRoot, "add", "-A")
	runGitInFixture(t, repoRoot, "commit", "--quiet", "-m", message)
}

// TestInventoryExcludePathHonored pins the exclude contract: --exclude-path is
// a repo-relative prefix, excluded files vanish from every number, and the
// default (no excludes) counts everything tracked. Also pins the line/word
// math on a known file.
func TestInventoryExcludePathHonored(t *testing.T) {
	repoRoot := newFixtureRepo(t)
	writeFixtureFile(t, repoRoot, "keep.md", "one two\nthree\n")
	writeFixtureFile(t, repoRoot, "sub/drop.md", "dropped words here\n")
	commitFixtureAll(t, repoRoot, "seed")

	fullReport, fullError := computeInventoryReport(repoRoot, nil)
	if fullError != nil {
		t.Fatalf("computeInventoryReport (no excludes): %v", fullError)
	}
	if len(fullReport.Files) != 2 {
		t.Fatalf("no-exclude file count = %d, want 2", len(fullReport.Files))
	}

	excludedReport, excludedError := computeInventoryReport(repoRoot, []string{"sub/"})
	if excludedError != nil {
		t.Fatalf("computeInventoryReport (excludes): %v", excludedError)
	}
	if len(excludedReport.Files) != 1 {
		t.Fatalf("excluded file count = %d, want 1", len(excludedReport.Files))
	}
	keptFile := excludedReport.Files[0]
	if keptFile.Path != "keep.md" {
		t.Fatalf("kept path = %q, want keep.md", keptFile.Path)
	}
	if keptFile.Lines != 2 || keptFile.Words != 3 {
		t.Fatalf("keep.md measured lines=%d words=%d, want lines=2 words=3", keptFile.Lines, keptFile.Words)
	}
}

// TestInventoryBinarySniff pins the binary handling: a NUL byte marks the file
// binary, it still counts as a file, and its lines/words stay zero so word
// counts on binaries never pollute the distributions.
func TestInventoryBinarySniff(t *testing.T) {
	repoRoot := newFixtureRepo(t)
	writeFixtureFile(t, repoRoot, "blob.bin", "head\x00tail")
	commitFixtureAll(t, repoRoot, "seed")

	report, computeError := computeInventoryReport(repoRoot, nil)
	if computeError != nil {
		t.Fatalf("computeInventoryReport: %v", computeError)
	}
	if len(report.Files) != 1 || report.BinaryCount != 1 {
		t.Fatalf("files=%d binaryCount=%d, want 1 and 1", len(report.Files), report.BinaryCount)
	}
	binaryFile := report.Files[0]
	if !binaryFile.Binary || binaryFile.Lines != 0 || binaryFile.Words != 0 {
		t.Fatalf("binary measurement = %+v, want Binary=true with zero lines/words", binaryFile)
	}
}

// TestInventoryBandSectionOnlyWithFlags pins the "bands come only from flags"
// contract at the renderer seam: no band flag → no band section at all; a
// passed threshold → the section with the strictly-greater rows.
func TestInventoryBandSectionOnlyWithFlags(t *testing.T) {
	report := inventoryReport{Files: []fileMeasurement{
		{Path: "small.md", Lines: 10, Words: 40},
		{Path: "large.md", Lines: 300, Words: 900},
	}}

	var withoutBands strings.Builder
	writeInventoryReport(&withoutBands, report, fileBandThresholds{
		WatchLines: bandThresholdUnset, FlagLines: bandThresholdUnset,
		WatchWords: bandThresholdUnset, FlagWords: bandThresholdUnset,
	}, 10)
	if strings.Contains(withoutBands.String(), "Band Flags") {
		t.Fatalf("no band flags passed, but output contains a Band Flags section:\n%s", withoutBands.String())
	}

	var withBands strings.Builder
	writeInventoryReport(&withBands, report, fileBandThresholds{
		WatchLines: 100, FlagLines: bandThresholdUnset,
		WatchWords: bandThresholdUnset, FlagWords: bandThresholdUnset,
	}, 10)
	bandOutput := withBands.String()
	if !strings.Contains(bandOutput, "Band Flags") {
		t.Fatalf("band flag passed, but no Band Flags section:\n%s", bandOutput)
	}
	if !strings.Contains(bandOutput, "| large.md | lines | 300 | WATCH |") {
		t.Fatalf("large.md (300 > 100) missing from band rows:\n%s", bandOutput)
	}
	if strings.Contains(bandOutput, "| small.md | lines") {
		t.Fatalf("small.md (10 <= 100) must not be flagged:\n%s", bandOutput)
	}
}
