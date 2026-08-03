package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeReleaseFixture seeds a repo root with an actions/version.md carrying the
// given version line and an optional CHANGELOG.md, and returns the root.
func writeReleaseFixture(t *testing.T, versionLine string, changelogBody string) string {
	t.Helper()
	repoRoot := t.TempDir()
	actionsDirectory := filepath.Join(repoRoot, "actions")
	if mkdirError := os.MkdirAll(actionsDirectory, 0o755); mkdirError != nil {
		t.Fatalf("mkdir actions: %v", mkdirError)
	}
	if versionLine != "" {
		versionFileBody := "# Version Action\n\n> Reports the installed version.\n\n" + versionLine + "\n\nMore prose below.\n"
		if writeError := os.WriteFile(filepath.Join(actionsDirectory, "version.md"), []byte(versionFileBody), 0o644); writeError != nil {
			t.Fatalf("write version.md: %v", writeError)
		}
	}
	if changelogBody != "" {
		if writeError := os.WriteFile(filepath.Join(repoRoot, "CHANGELOG.md"), []byte(changelogBody), 0o644); writeError != nil {
			t.Fatalf("write CHANGELOG.md: %v", writeError)
		}
	}
	return repoRoot
}

func TestBumpSemanticVersion(t *testing.T) {
	testCases := []struct {
		caseName        string
		currentVersion  string
		bumpSize        string
		expectedVersion string
	}{
		{"patch on 0.163.3", "0.163.3", "patch", "0.163.4"},
		{"minor on 0.163.3 zeroes the patch", "0.163.3", "minor", "0.164.0"},
		{"major on 0.163.3 zeroes minor and patch", "0.163.3", "major", "1.0.0"},
		{"patch rolls into a two-digit patch", "1.2.9", "patch", "1.2.10"},
		{"minor rolls into a two-digit minor", "1.9.4", "minor", "1.10.0"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			bumpedVersion, bumpError := bumpSemanticVersion(testCase.currentVersion, testCase.bumpSize)
			if bumpError != nil {
				t.Fatalf("bumpSemanticVersion: unexpected error: %v", bumpError)
			}
			if bumpedVersion != testCase.expectedVersion {
				t.Errorf("bumpSemanticVersion(%q, %q) = %q, want %q",
					testCase.currentVersion, testCase.bumpSize, bumpedVersion, testCase.expectedVersion)
			}
		})
	}
}

// The bump size is an argument, never an inference (requirement 2): an absent or
// unrecognized size is an error, never a silent default to patch.
func TestBumpSemanticVersionRejectsAnUnnamedBumpSize(t *testing.T) {
	for _, badBumpSize := range []string{"", "PATCH-ish", "bugfix", "1", "auto"} {
		if _, bumpError := bumpSemanticVersion("0.163.3", badBumpSize); bumpError == nil {
			t.Errorf("bumpSemanticVersion(_, %q) succeeded; the bump size must be an explicit patch|minor|major", badBumpSize)
		}
	}
}

func TestReadCurrentVersionAnchorsOnTheCurrentVersionPrefix(t *testing.T) {
	repoRoot := writeReleaseFixture(t, "**Current version**: 0.163.3", "")
	versionFilePath := filepath.Join(repoRoot, "actions", "version.md")

	currentVersion, readError := readCurrentVersion(versionFilePath)
	if readError != nil {
		t.Fatalf("readCurrentVersion: %v", readError)
	}
	if currentVersion != "0.163.3" {
		t.Errorf("readCurrentVersion = %q, want %q", currentVersion, "0.163.3")
	}
}

func TestReadCurrentVersionErrorsWhenNoVersionLineExists(t *testing.T) {
	repoRoot := writeReleaseFixture(t, "This file carries no version line.", "")
	versionFilePath := filepath.Join(repoRoot, "actions", "version.md")

	if _, readError := readCurrentVersion(versionFilePath); readError == nil {
		t.Error("readCurrentVersion succeeded on a file with no **Current version**: line; it must report instead")
	}
}

// next-version writes the bumped value back and reads it back to confirm the
// write landed before reporting success (requirement 6).
func TestAllocateNextVersionWritesBackAndConfirms(t *testing.T) {
	repoRoot := writeReleaseFixture(t, "**Current version**: 0.163.3", "")
	versionFilePath := filepath.Join(repoRoot, "actions", "version.md")

	allocatedVersion, allocateError := allocateNextVersion(versionFilePath, "patch")
	if allocateError != nil {
		t.Fatalf("allocateNextVersion: %v", allocateError)
	}
	if allocatedVersion != "0.163.4" {
		t.Fatalf("allocateNextVersion = %q, want %q", allocatedVersion, "0.163.4")
	}

	rewrittenBytes, readError := os.ReadFile(versionFilePath)
	if readError != nil {
		t.Fatalf("read back: %v", readError)
	}
	rewrittenText := string(rewrittenBytes)
	if !strings.Contains(rewrittenText, "**Current version**: 0.163.4") {
		t.Error("the bumped version was not written into the file")
	}
	if strings.Contains(rewrittenText, "0.163.3") {
		t.Error("the old version survived the rewrite")
	}
	// Only the version line may change: the surrounding prose is a prompt an
	// agent reads, and rewriting it would be a silent instruction edit.
	if !strings.Contains(rewrittenText, "> Reports the installed version.") ||
		!strings.Contains(rewrittenText, "More prose below.") {
		t.Error("allocateNextVersion disturbed prose outside the version line")
	}
}

// next-version writes exactly one file (the Constraints section): nothing under
// do-work/, and no CHANGELOG.md — not even to create it.
func TestAllocateNextVersionWritesNothingElse(t *testing.T) {
	repoRoot := writeReleaseFixture(t, "**Current version**: 1.0.0", "# Changelog\n\n## 1.0.0 — Seed (2026-01-01)\n")
	if mkdirError := os.MkdirAll(filepath.Join(repoRoot, "do-work", "queue"), 0o755); mkdirError != nil {
		t.Fatalf("mkdir: %v", mkdirError)
	}
	queueFile := filepath.Join(repoRoot, "do-work", "queue", "REQ-001-untouched.md")
	if writeError := os.WriteFile(queueFile, []byte("---\nid: REQ-001\nstatus: pending\n---\n"), 0o644); writeError != nil {
		t.Fatalf("write queue file: %v", writeError)
	}

	changelogBefore, _ := os.ReadFile(filepath.Join(repoRoot, "CHANGELOG.md"))
	queueBefore, _ := os.ReadFile(queueFile)

	if _, allocateError := allocateNextVersion(filepath.Join(repoRoot, "actions", "version.md"), "minor"); allocateError != nil {
		t.Fatalf("allocateNextVersion: %v", allocateError)
	}

	changelogAfter, _ := os.ReadFile(filepath.Join(repoRoot, "CHANGELOG.md"))
	if string(changelogAfter) != string(changelogBefore) {
		t.Error("allocateNextVersion touched CHANGELOG.md — the changelog is an owner-only, human-authored write")
	}
	queueAfter, _ := os.ReadFile(queueFile)
	if string(queueAfter) != string(queueBefore) {
		t.Error("allocateNextVersion touched a file under do-work/ — its write surface is exactly one version file")
	}
}

func TestReadChangelogEntries(t *testing.T) {
	changelogBody := `# Changelog

What's new.

---

## 0.163.3 — Board Copy Includes the REQ Id (2026-08-02)

Lead sentence.

## 0.163.2 — Audit Sweep (2026-08-01)

Another.

## 0.100.0 — Board Copy Includes the REQ Id (2026-07-01)

A reused title, deliberately.
`
	repoRoot := writeReleaseFixture(t, "**Current version**: 0.163.4", changelogBody)

	entries, readError := readChangelogEntries(filepath.Join(repoRoot, "CHANGELOG.md"))
	if readError != nil {
		t.Fatalf("readChangelogEntries: %v", readError)
	}
	if len(entries) != 3 {
		t.Fatalf("parsed %d entries, want 3", len(entries))
	}
	if entries[0].Version != "0.163.3" {
		t.Errorf("newest entry version = %q, want 0.163.3 (newest is first)", entries[0].Version)
	}
	if entries[0].Title != "Board Copy Includes the REQ Id" {
		t.Errorf("newest entry title = %q, want %q", entries[0].Title, "Board Copy Includes the REQ Id")
	}
	if entries[2].Version != "0.100.0" {
		t.Errorf("oldest entry version = %q, want 0.100.0", entries[2].Version)
	}
}

func TestCompareSemanticVersions(t *testing.T) {
	testCases := []struct {
		leftVersion    string
		rightVersion   string
		expectedResult int
	}{
		{"0.163.4", "0.163.3", 1},
		{"0.163.3", "0.163.3", 0},
		{"0.163.3", "0.163.4", -1},
		{"0.164.0", "0.163.99", 1},
		{"1.0.0", "0.999.999", 1},
		{"0.100.0", "0.99.0", 1}, // numeric compare, not lexicographic
	}

	for _, testCase := range testCases {
		comparisonResult, compareError := compareSemanticVersions(testCase.leftVersion, testCase.rightVersion)
		if compareError != nil {
			t.Fatalf("compareSemanticVersions(%q, %q): %v", testCase.leftVersion, testCase.rightVersion, compareError)
		}
		if comparisonResult != testCase.expectedResult {
			t.Errorf("compareSemanticVersions(%q, %q) = %d, want %d",
				testCase.leftVersion, testCase.rightVersion, comparisonResult, testCase.expectedResult)
		}
	}
}
