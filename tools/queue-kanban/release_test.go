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

// TestParseNextVersionArgumentsAcceptsFlagsOnEitherSideOfTheBumpSize pins the
// defect that shipped in 0.166.0: flag.FlagSet.Parse halts at the first non-flag
// argument, so every flag placed AFTER the positional bump size was silently
// discarded — and the command then bumped whatever repo it was launched from and
// exited 0. The skill's own prescribed invocation put --repo-root last, so the
// documented form was the broken one. The first case below is that exact shape.
func TestParseNextVersionArgumentsAcceptsFlagsOnEitherSideOfTheBumpSize(t *testing.T) {
	testCases := []struct {
		caseName            string
		arguments           []string
		expectedBumpSize    string
		expectedRepoRoot    string
		expectedVersionFile string
	}{
		{
			caseName:            "flags after the positional (the invocation actions/work.md prescribes)",
			arguments:           []string{"patch", "--repo-root", "/tmp/x", "--version-file", "/tmp/x/actions/version.md"},
			expectedBumpSize:    "patch",
			expectedRepoRoot:    "/tmp/x",
			expectedVersionFile: "/tmp/x/actions/version.md",
		},
		{
			caseName:            "flags before the positional",
			arguments:           []string{"--repo-root", "/tmp/x", "--version-file", "/tmp/x/actions/version.md", "minor"},
			expectedBumpSize:    "minor",
			expectedRepoRoot:    "/tmp/x",
			expectedVersionFile: "/tmp/x/actions/version.md",
		},
		{
			caseName:            "interleaved around the positional",
			arguments:           []string{"--repo-root", "/tmp/x", "major", "--version-file", "/tmp/x/VERSION"},
			expectedBumpSize:    "major",
			expectedRepoRoot:    "/tmp/x",
			expectedVersionFile: "/tmp/x/VERSION",
		},
		{
			caseName:            "equals form on either side",
			arguments:           []string{"--repo-root=/tmp/x", "patch", "--version-file=/tmp/x/VERSION"},
			expectedBumpSize:    "patch",
			expectedRepoRoot:    "/tmp/x",
			expectedVersionFile: "/tmp/x/VERSION",
		},
		{
			caseName:         "bare bump size, no flags",
			arguments:        []string{"patch"},
			expectedBumpSize: "patch",
		},
	}

	for _, testCase := range testCases {
		parsed, parseError := parseNextVersionArguments(testCase.arguments)
		if parseError != nil {
			t.Errorf("%s: parseNextVersionArguments(%q) returned error %v, want success",
				testCase.caseName, testCase.arguments, parseError)
			continue
		}
		if parsed.BumpSize != testCase.expectedBumpSize {
			t.Errorf("%s: BumpSize = %q, want %q", testCase.caseName, parsed.BumpSize, testCase.expectedBumpSize)
		}
		if parsed.RepoRootOverride != testCase.expectedRepoRoot {
			t.Errorf("%s: RepoRootOverride = %q, want %q", testCase.caseName, parsed.RepoRootOverride, testCase.expectedRepoRoot)
		}
		if parsed.VersionFileOverride != testCase.expectedVersionFile {
			t.Errorf("%s: VersionFileOverride = %q, want %q", testCase.caseName, parsed.VersionFileOverride, testCase.expectedVersionFile)
		}
	}
}

// TestParseNextVersionArgumentsRejectsRatherThanIgnores covers the other half of
// the same defect: unconsumed tokens were left sitting in Arg(1)/Arg(2), neither
// rejected nor warned about, which is how a discarded --repo-root produced a
// successful-looking bump of the wrong tree.
func TestParseNextVersionArgumentsRejectsRatherThanIgnores(t *testing.T) {
	testCases := []struct {
		caseName            string
		arguments           []string
		expectedErrorSubstr string
	}{
		{
			caseName:            "missing bump size",
			arguments:           []string{},
			expectedErrorSubstr: "name the bump size",
		},
		{
			caseName:            "only flags, no bump size",
			arguments:           []string{"--repo-root", "/tmp/x"},
			expectedErrorSubstr: "name the bump size",
		},
		{
			caseName:            "stray second positional",
			arguments:           []string{"patch", "minor"},
			expectedErrorSubstr: "unrecognized argument(s)",
		},
		{
			caseName:            "stray positional after a flag",
			arguments:           []string{"patch", "--repo-root", "/tmp/x", "stray"},
			expectedErrorSubstr: "unrecognized argument(s)",
		},
		{
			caseName:            "unknown flag before the positional",
			arguments:           []string{"--nope", "patch"},
			expectedErrorSubstr: "not defined",
		},
		{
			caseName:            "unknown flag after the positional",
			arguments:           []string{"patch", "--nope"},
			expectedErrorSubstr: "not defined",
		},
	}

	for _, testCase := range testCases {
		_, parseError := parseNextVersionArguments(testCase.arguments)
		if parseError == nil {
			t.Errorf("%s: parseNextVersionArguments(%q) succeeded, want an error", testCase.caseName, testCase.arguments)
			continue
		}
		if !strings.Contains(parseError.Error(), testCase.expectedErrorSubstr) {
			t.Errorf("%s: error = %q, want it to contain %q",
				testCase.caseName, parseError.Error(), testCase.expectedErrorSubstr)
		}
	}
}

// TestRejectLeftoverArgumentsIsTheSharedRule pins the condition rather than
// today's subcommand list: any subcommand finishing its parse with tokens left
// over must fail. Flags-only subcommands are reachable by the same shape — a
// stray token placed first halts Parse and every flag after it is discarded.
func TestRejectLeftoverArgumentsIsTheSharedRule(t *testing.T) {
	if rejectLeftoverArguments("verify", nil) != nil {
		t.Error("no leftover arguments must not be an error")
	}
	if rejectLeftoverArguments("verify", []string{}) != nil {
		t.Error("an empty leftover slice must not be an error")
	}
	leftoverError := rejectLeftoverArguments("verify", []string{"stray"})
	if leftoverError == nil {
		t.Fatal("a leftover argument must be an error")
	}
	if !strings.Contains(leftoverError.Error(), "verify") || !strings.Contains(leftoverError.Error(), "stray") {
		t.Errorf("error must name the subcommand and the offending token, got %q", leftoverError.Error())
	}
}
