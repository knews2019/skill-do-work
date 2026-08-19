package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// currentVersionLinePrefix is the anchor for the version this tool reads and
// writes. It is a PREFIX, not a path: do-work keeps its version inside an action
// file (actions/version.md) rather than a dedicated data file, and anchoring on
// the marker rather than the location is what lets --version-file point anywhere
// the same convention is used.
const currentVersionLinePrefix = "**Current version**: "

// defaultVersionFileRelativePath is where a skill-development checkout keeps its
// version line. In a CONSUMER install this path does not exist at the repo root
// (the board skill lives under .claude/skills/do-work-board/ while do-work/ sits at the
// consumer root), and that is deliberate: next-version then reports and writes
// nothing, and the calling action falls back to the Changelog Entry Procedure's
// own version-source resolution (actions/work-reference.md). This tool is an
// accelerator for a repo that matches do-work's convention, never a dependency.
const defaultVersionFileRelativePath = "actions/version.md"

// suiteVersionFileRelativePath is where the SAME version line lives in the suite's own
// development checkout since the four-skill split: the core package became a subdirectory,
// so `actions/version.md` stopped resolving from the repo root while the release CHANGELOG.md
// stayed at the root. One additional known location, not a search — and deliberately NOT
// consulted by resolveVersionFilePath, because that function is next-version's writer path
// and teaching it this location would make next-version start rewriting a file it correctly
// finds nothing at today (REQ-282).
const suiteVersionFileRelativePath = "skills/do-work/actions/version.md"

// semanticVersionPattern matches a bare X.Y.Z. Pre-release and build metadata are
// deliberately unsupported: do-work versions are plain triples, and quietly
// accepting 1.2.3-rc1 would mean quietly deciding how to bump it.
var semanticVersionPattern = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)

// changelogEntryHeadingPattern matches the house entry key:
//
//	## X.Y.Z — Short Descriptive Title (YYYY-MM-DD)
//
// The em dash is the house separator; a hyphen is accepted too so an entry typed
// with the wrong dash is still read rather than silently ignored (ignoring it
// would make the newest-entry probe compare against the wrong entry).
var changelogEntryHeadingPattern = regexp.MustCompile(`^##\s+(\d+\.\d+\.\d+)\s*[—-]\s*(.+?)\s*\((\d{4}-\d{2}-\d{2})\)\s*$`)

// ChangelogEntry is one parsed house-format changelog heading, newest first in
// the slice readChangelogEntries returns.
type ChangelogEntry struct {
	Version string
	Title   string
	Date    string
}

// resolveVersionFilePath returns the version file to read or write: the override
// when non-empty, else <repoRoot>/actions/version.md. Read-and-WRITE, so its resolution is
// deliberately narrow — see resolveReleaseProbeVersionFilePath for the read-only variant.
func resolveVersionFilePath(repoRoot string, versionFileOverride string) string {
	if strings.TrimSpace(versionFileOverride) != "" {
		return versionFileOverride
	}
	return filepath.Join(repoRoot, defaultVersionFileRelativePath)
}

// resolveReleaseProbeVersionFilePath finds the version file the READ-ONLY release probes
// should compare against, and reports whether this repo root is a suite checkout at all.
//
// Two known locations, tried in order — the pre-split root path first so nothing that
// already worked changes, then the modular suite path. Existence is the whole test: the
// REQ's constraint is that detecting "not a suite checkout" must not become its own
// inference engine, and it does not need to be, because these are the only two layouts the
// suite has ever had.
//
// Neither present means this root is not a suite checkout — a consumer install, where the
// root CHANGELOG.md belongs to the consumer and the release ritual being probed is not
// theirs. That is reported as NOT APPLICABLE rather than skipped: a skipped probe is an
// unverified invariant someone should act on, and there is nothing here to act on.
// Deliberately NOT a walk up to a vendored suite root — see the REQ-282 constraint and
// decisions/records/adr-019-four-skill-suite-contract.md, which already assigns install
// integrity to the updater.
func resolveReleaseProbeVersionFilePath(repoRoot string) (versionFilePath string, isSuiteCheckout bool) {
	for _, relativePath := range []string{defaultVersionFileRelativePath, suiteVersionFileRelativePath} {
		candidatePath := filepath.Join(repoRoot, relativePath)
		if fileInfo, statError := os.Stat(candidatePath); statError == nil && !fileInfo.IsDir() {
			return candidatePath, true
		}
	}
	return "", false
}

// readCurrentVersion returns the X.Y.Z on the file's `**Current version**: `
// line. A missing file, or a file with no such line, is an error — never a
// guessed-at default, because a guess here would be written back as fact.
func readCurrentVersion(versionFilePath string) (string, error) {
	fileBytes, readError := os.ReadFile(versionFilePath)
	if readError != nil {
		return "", fmt.Errorf("no version file readable at %s: %w", versionFilePath, readError)
	}
	for _, line := range strings.Split(string(fileBytes), "\n") {
		trimmedLine := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmedLine, currentVersionLinePrefix) {
			continue
		}
		versionText := strings.TrimSpace(strings.TrimPrefix(trimmedLine, currentVersionLinePrefix))
		if !semanticVersionPattern.MatchString(versionText) {
			return "", fmt.Errorf("%s carries a %q line but %q is not an X.Y.Z version",
				versionFilePath, strings.TrimSuffix(currentVersionLinePrefix, " "), versionText)
		}
		return versionText, nil
	}
	return "", fmt.Errorf("%s has no %q line — this tool only allocates versions for a repo that uses that marker",
		versionFilePath, strings.TrimSuffix(currentVersionLinePrefix, " "))
}

// bumpSemanticVersion applies an explicitly named bump size. The size is an
// ARGUMENT and never an inference: patch vs minor vs major is a judgment about
// what the change did to consumers (actions/work-reference.md → Changelog Entry
// Procedure), which a tool reading a version string cannot make. An empty or
// unrecognized size is an error rather than a default to patch — a wrong-but-
// plausible bump is worse than a refusal, because it ships.
func bumpSemanticVersion(currentVersion string, bumpSize string) (string, error) {
	match := semanticVersionPattern.FindStringSubmatch(strings.TrimSpace(currentVersion))
	if match == nil {
		return "", fmt.Errorf("%q is not an X.Y.Z version", currentVersion)
	}
	majorNumber, _ := strconv.Atoi(match[1])
	minorNumber, _ := strconv.Atoi(match[2])
	patchNumber, _ := strconv.Atoi(match[3])

	switch strings.TrimSpace(bumpSize) {
	case "patch":
		patchNumber++
	case "minor":
		minorNumber++
		patchNumber = 0
	case "major":
		majorNumber++
		minorNumber = 0
		patchNumber = 0
	default:
		return "", fmt.Errorf("bump size must be one of patch | minor | major (got %q) — the tool must not infer it", bumpSize)
	}
	return fmt.Sprintf("%d.%d.%d", majorNumber, minorNumber, patchNumber), nil
}

// allocateNextVersion bumps the version file's `**Current version**: ` line by
// the named size, writes it back, re-reads the file to confirm the value landed,
// and returns the new version.
//
// The read-back is the point (REQ-072 requirement 6): reporting a number the
// caller then puts in a changelog heading, when the write silently failed or was
// clobbered, would produce exactly the version/changelog mismatch `verify` exists
// to catch.
//
// This is the tool's ONLY write outside the board's testing fields, and its whole
// write surface is this one line in this one file: nothing under do-work/, no REQ
// frontmatter, no status, and never CHANGELOG.md.
func allocateNextVersion(versionFilePath string, bumpSize string) (string, error) {
	currentVersion, readError := readCurrentVersion(versionFilePath)
	if readError != nil {
		return "", readError
	}
	bumpedVersion, bumpError := bumpSemanticVersion(currentVersion, bumpSize)
	if bumpError != nil {
		return "", bumpError
	}

	originalBytes, readError := os.ReadFile(versionFilePath)
	if readError != nil {
		return "", readError
	}

	rewrittenLines := strings.Split(string(originalBytes), "\n")
	replacedLineCount := 0
	for lineIndex, line := range rewrittenLines {
		if !strings.HasPrefix(strings.TrimSpace(line), currentVersionLinePrefix) {
			continue
		}
		// Preserve whatever leading whitespace the line had; only the version
		// text changes. The surrounding file is a prompt an agent reads, so any
		// edit beyond this line would be a silent instruction change.
		leadingWhitespace := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
		rewrittenLines[lineIndex] = leadingWhitespace + currentVersionLinePrefix + bumpedVersion
		replacedLineCount++
	}
	if replacedLineCount != 1 {
		return "", fmt.Errorf("%s has %d %q lines — expected exactly 1; refusing to write",
			versionFilePath, replacedLineCount, strings.TrimSuffix(currentVersionLinePrefix, " "))
	}

	if writeError := writeFileAtomically(versionFilePath, []byte(strings.Join(rewrittenLines, "\n"))); writeError != nil {
		return "", writeError
	}

	confirmedVersion, confirmError := readCurrentVersion(versionFilePath)
	if confirmError != nil {
		return "", fmt.Errorf("wrote %s but could not read it back: %w", bumpedVersion, confirmError)
	}
	if confirmedVersion != bumpedVersion {
		return "", fmt.Errorf("wrote %s to %s but it reads back as %s — the write did not land",
			bumpedVersion, versionFilePath, confirmedVersion)
	}
	return bumpedVersion, nil
}

// readChangelogEntries parses every house-format entry heading in a CHANGELOG.md,
// in file order (newest first, per the house convention). A repo whose changelog
// follows a different convention yields no entries and no error — the caller
// decides what to do with that, which is how the release probes skip instead of
// inventing findings (actions/work-reference.md → Changelog Entry Procedure's
// precedence check: match the existing format, never impose the house one).
func readChangelogEntries(changelogPath string) ([]ChangelogEntry, error) {
	fileBytes, readError := os.ReadFile(changelogPath)
	if readError != nil {
		return nil, readError
	}
	var entries []ChangelogEntry
	for _, line := range strings.Split(string(fileBytes), "\n") {
		match := changelogEntryHeadingPattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		entries = append(entries, ChangelogEntry{
			Version: match[1],
			Title:   strings.TrimSpace(match[2]),
			Date:    match[3],
		})
	}
	return entries, nil
}

// compareSemanticVersions returns 1, 0, or -1 for left >, ==, or < right,
// comparing component by component as numbers. String comparison would order
// 0.100.0 below 0.99.0, which is the mistake this exists to avoid.
func compareSemanticVersions(leftVersion string, rightVersion string) (int, error) {
	leftMatch := semanticVersionPattern.FindStringSubmatch(strings.TrimSpace(leftVersion))
	rightMatch := semanticVersionPattern.FindStringSubmatch(strings.TrimSpace(rightVersion))
	if leftMatch == nil {
		return 0, fmt.Errorf("%q is not an X.Y.Z version", leftVersion)
	}
	if rightMatch == nil {
		return 0, fmt.Errorf("%q is not an X.Y.Z version", rightVersion)
	}
	for componentIndex := 1; componentIndex <= 3; componentIndex++ {
		leftComponent, _ := strconv.Atoi(leftMatch[componentIndex])
		rightComponent, _ := strconv.Atoi(rightMatch[componentIndex])
		if leftComponent > rightComponent {
			return 1, nil
		}
		if leftComponent < rightComponent {
			return -1, nil
		}
	}
	return 0, nil
}
