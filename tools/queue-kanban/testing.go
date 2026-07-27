package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Testing track — who tested a finished REQ, and where it stands.
//
// The queue's Markdown files ARE the database: the board's testing view writes
// a small set of placeholder frontmatter fields into the REQ file itself
// (mirrored in actions/work-reference.md's Request File Schema), and tester
// profiles live in do-work/testers.md as plain bullet lines. There is
// deliberately no locking or concurrency control — every write lands in the
// working tree where git provides the history, review, and rollback story.
//
//	testing_status: in-testing | tested | returned   # absent → not tested yet
//	tested_by: "Alice"                               # tester profile (raw user text, always quoted)
//	testing_updated_at: 2026-07-17T10:00:00Z         # stamped by the server on every transition
//	testing_feedback: "…"                            # present only while status is returned
//
// The testing track is orthogonal to the work pipeline's `status` field: the
// board never touches `status`, and the work pipeline never touches the
// testing placeholders.

// Canonical testing_status vocabulary. An absent field means "not tested yet";
// there is no explicit "untested" value on disk.
const (
	testingStatusInTesting = "in-testing"
	testingStatusTested    = "tested"
	testingStatusReturned  = "returned"
)

// testersFileRelativePath is where tester profiles live, relative to the
// do-work/ tree root — a sibling of notes.md, discovered by the same walk.
const testersFileRelativePath = "testers.md"

// testingWriteMutex serializes the board's two write surfaces (REQ frontmatter
// upserts and testers.md appends). The HTTP server runs each request in its own
// goroutine, so without this a double-submitted UI action could interleave two
// read-modify-write cycles on the same file.
var testingWriteMutex sync.Mutex

// testersFileHeader is written once when the profile store is first created, so
// a hand-opened do-work/testers.md explains itself.
const testersFileHeader = "# Testers\n\nTester profiles for the do-work board's testing view (one `- Name` bullet per\nprofile). The board appends new profiles; edit or remove lines by hand.\n\n"

// normalizeTestingStatus collapses natural spelling variants of the testing
// vocabulary onto the canonical enum, lower-casing/trimming the rest. It
// mirrors the Schema Read Contract's normalize-and-warn pattern
// (actions/work-reference.md): a value that still doesn't match after
// normalization is NOT silently remapped — the caller flags it (see
// collectTestingWarnings) and the ticket renders as not-yet-tested with an
// invalid marker.
func normalizeTestingStatus(rawTestingStatus string) string {
	normalized := strings.ToLower(strings.TrimSpace(rawTestingStatus))
	switch normalized {
	case "in_testing", "in testing", "testing", "selected-for-testing", "selected for testing":
		return testingStatusInTesting
	case "returned-with-feedback", "returned_with_feedback", "returned with feedback":
		return testingStatusReturned
	default:
		return normalized
	}
}

// isKnownTestingStatus reports whether a normalized testing status is in the
// canonical vocabulary. The empty string is "not tested yet" — valid, but not a
// member of the enum.
func isKnownTestingStatus(normalizedTestingStatus string) bool {
	switch normalizedTestingStatus {
	case testingStatusInTesting, testingStatusTested, testingStatusReturned:
		return true
	default:
		return false
	}
}

// collectTestingWarnings emits one warning per ticket whose testing_status is
// present but outside the canonical vocabulary — the never-silently-drop leg of
// the normalize-and-warn contract. The ticket keeps rendering (as not-yet-tested
// with an invalid flag); the warning is the feedback channel.
func collectTestingWarnings(tickets []*RequestTicket) []string {
	var testingWarnings []string
	for _, ticket := range tickets {
		if ticket.TestingStatusUnrecognized {
			testingWarnings = append(testingWarnings, fmt.Sprintf(
				"%s has unrecognized testing_status %q — expected one of [in-testing, tested, returned]; shown as not tested",
				ticket.RequestId, ticket.OriginalTestingStatus))
		}
	}
	return testingWarnings
}

// loadTestingProfiles reads do-work/testers.md into an ordered profile list —
// one profile per bullet line, file order preserved (the file is append-only
// from the board's side, so file order is creation order). Best-effort like
// loadQueueNotes: an absent or unreadable file yields no profiles.
func loadTestingProfiles(testersFilePath string) []string {
	if testersFilePath == "" {
		return nil
	}
	contentBytes, readError := os.ReadFile(testersFilePath)
	if readError != nil {
		return nil
	}

	var profiles []string
	seenProfiles := map[string]bool{}
	for _, rawLine := range strings.Split(string(contentBytes), "\n") {
		trimmedLine := strings.TrimSpace(rawLine)
		profileName := ""
		for _, bulletPrefix := range []string{"- ", "* ", "+ "} {
			if strings.HasPrefix(trimmedLine, bulletPrefix) {
				profileName = strings.TrimSpace(strings.TrimPrefix(trimmedLine, bulletPrefix))
				break
			}
		}
		if profileName == "" {
			continue
		}
		dedupeKey := strings.ToLower(profileName)
		if seenProfiles[dedupeKey] {
			continue
		}
		seenProfiles[dedupeKey] = true
		profiles = append(profiles, profileName)
	}
	return profiles
}

// validateTesterProfileName trims and validates a tester profile name coming
// from the browser: non-empty, at most 80 characters, and free of control
// characters (which would break the one-bullet-per-line store and the one-line
// frontmatter placeholder).
func validateTesterProfileName(rawProfileName string) (string, error) {
	profileName := strings.TrimSpace(rawProfileName)
	if profileName == "" {
		return "", fmt.Errorf("tester name must not be empty")
	}
	if len(profileName) > 80 {
		return "", fmt.Errorf("tester name must be at most 80 characters")
	}
	for _, character := range profileName {
		if character < 0x20 || character == 0x7f {
			return "", fmt.Errorf("tester name must not contain control characters")
		}
	}
	return profileName, nil
}

// appendTestingProfile adds a profile to do-work/testers.md, creating the file
// (with its explanatory header) on first use. Adding an already-known name
// (case-insensitive) is a no-op, not an error — the caller re-reads the file
// either way. Returns the updated profile list.
func appendTestingProfile(repoRoot string, rawProfileName string) ([]string, error) {
	profileName, validateError := validateTesterProfileName(rawProfileName)
	if validateError != nil {
		return nil, validateError
	}
	testersFilePath := filepath.Join(repoRoot, "do-work", testersFileRelativePath)

	// Serialize with every other testing write: two concurrent appends (a
	// double-clicked "add tester") would otherwise both see an empty file and
	// both write the header, or interleave their read-check-append cycles.
	testingWriteMutex.Lock()
	defer testingWriteMutex.Unlock()
	if parentError := validateTestingWriteParent(repoRoot, testersFilePath); parentError != nil {
		return nil, parentError
	}

	existingProfiles := loadTestingProfiles(testersFilePath)
	for _, existingProfile := range existingProfiles {
		if strings.EqualFold(existingProfile, profileName) {
			return existingProfiles, nil
		}
	}

	// Guard the append the same way the REQ write path is guarded: a
	// testers.md that is a symlink (or any non-regular file) planted in the
	// tree must not redirect the write elsewhere.
	if lstatInfo, lstatError := os.Lstat(testersFilePath); lstatError == nil && !lstatInfo.Mode().IsRegular() {
		return nil, fmt.Errorf("%s is not a regular file — refusing to write through it", testersFilePath)
	}

	fileHandle, openError := os.OpenFile(testersFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if openError != nil {
		return nil, fmt.Errorf("opening testers file: %w", openError)
	}
	defer fileHandle.Close()

	fileInfo, statError := fileHandle.Stat()
	if statError != nil {
		return nil, fmt.Errorf("stat testers file: %w", statError)
	}
	appendText := "- " + profileName + "\n"
	if fileInfo.Size() == 0 {
		appendText = testersFileHeader + appendText
	}
	if _, writeError := fileHandle.WriteString(appendText); writeError != nil {
		return nil, fmt.Errorf("appending tester profile: %w", writeError)
	}
	return append(existingProfiles, profileName), nil
}

// validateTestingWriteTarget rejects a REQ write target that would escape the
// repository's do-work tree. The path itself always comes from the parsed tree
// (never from client input), but a checked-out repo can contain a symlink named
// REQ-*.md — or a symlinked parent directory — whose target is outside the
// repository. The file must be regular (no symlink), and its parent must pass
// the same containment check used by the testers.md writer.
func validateTestingWriteTarget(repoRoot string, requestFilePath string) error {
	lstatInfo, lstatError := os.Lstat(requestFilePath)
	if lstatError != nil {
		return fmt.Errorf("stat %s: %w", requestFilePath, lstatError)
	}
	if !lstatInfo.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file — refusing to write testing placeholders through a symlink", requestFilePath)
	}
	return validateTestingWriteParent(repoRoot, requestFilePath)
}

// validateTestingWriteParent permits a testing write only when the resolved
// do-work root remains inside the resolved repository and the target's resolved
// parent remains inside that do-work root. Checking both boundaries matters:
// an otherwise-valid testers.md path can escape when do-work/ itself is a
// symlink outside the repository.
func validateTestingWriteParent(repoRoot string, testingFilePath string) error {
	absoluteRepoRoot, absoluteRootError := filepath.Abs(repoRoot)
	if absoluteRootError != nil {
		return fmt.Errorf("resolving repository root: %w", absoluteRootError)
	}
	resolvedRepoRoot, repoEvalError := filepath.EvalSymlinks(absoluteRepoRoot)
	if repoEvalError != nil {
		return fmt.Errorf("resolving repository root: %w", repoEvalError)
	}
	resolvedDoWorkRoot, rootEvalError := filepath.EvalSymlinks(filepath.Join(absoluteRepoRoot, "do-work"))
	if rootEvalError != nil {
		return fmt.Errorf("resolving do-work root: %w", rootEvalError)
	}
	if !resolvedPathIsWithinDirectory(resolvedRepoRoot, resolvedDoWorkRoot) {
		return fmt.Errorf("do-work root resolves outside the repository — refusing to write %s", testingFilePath)
	}

	absoluteParentDirectory, absoluteParentError := filepath.Abs(filepath.Dir(testingFilePath))
	if absoluteParentError != nil {
		return fmt.Errorf("resolving %s: %w", filepath.Dir(testingFilePath), absoluteParentError)
	}
	resolvedParentDirectory, parentEvalError := filepath.EvalSymlinks(absoluteParentDirectory)
	if parentEvalError != nil {
		return fmt.Errorf("resolving %s: %w", filepath.Dir(testingFilePath), parentEvalError)
	}
	if !resolvedPathIsWithinDirectory(resolvedDoWorkRoot, resolvedParentDirectory) {
		return fmt.Errorf("%s resolves outside the do-work tree — refusing to write testing data", testingFilePath)
	}
	return nil
}

func resolvedPathIsWithinDirectory(resolvedDirectory string, resolvedCandidatePath string) bool {
	relativePath, relativeError := filepath.Rel(resolvedDirectory, resolvedCandidatePath)
	return relativeError == nil && relativePath != ".." && !strings.HasPrefix(relativePath, ".."+string(filepath.Separator))
}

// frontmatterFieldUpdate is one placeholder mutation for
// upsertFrontmatterFields: set FieldKey to FieldValueLine's value, or remove
// the field entirely when RemoveField is true.
type frontmatterFieldUpdate struct {
	FieldKey    string
	FieldValue  string // the already-encoded YAML value text (ignored when RemoveField)
	RemoveField bool
}

// upsertFrontmatterFields rewrites exactly the named top-level frontmatter
// fields of a REQ file in place, preserving every other line verbatim — the
// file is the database, so the edit must be surgical. An existing key's line is
// replaced where it sits (its indented/list continuation lines, if any, are
// dropped with it — the new value is always a one-line scalar); a missing key
// is inserted just above the closing fence; RemoveField deletes the key and its
// continuation lines. A file without a well-formed frontmatter block is an
// error, never a guessed edit.
func upsertFrontmatterFields(filePath string, fieldUpdates []frontmatterFieldUpdate) error {
	// Serialize with every other testing write — two concurrent upserts on the
	// same REQ would race their read-modify-write cycles and one edit would
	// silently vanish.
	testingWriteMutex.Lock()
	defer testingWriteMutex.Unlock()

	contentBytes, readError := os.ReadFile(filePath)
	if readError != nil {
		return fmt.Errorf("reading %s: %w", filePath, readError)
	}

	fileLines := strings.Split(string(contentBytes), "\n")
	if len(fileLines) == 0 || strings.TrimRight(fileLines[0], "\r") != "---" {
		return fmt.Errorf("%s has no frontmatter block to update", filePath)
	}
	closingFenceIndex := -1
	for lineIndex := 1; lineIndex < len(fileLines); lineIndex++ {
		if strings.TrimRight(fileLines[lineIndex], "\r") == "---" {
			closingFenceIndex = lineIndex
			break
		}
	}
	if closingFenceIndex < 0 {
		return fmt.Errorf("%s has no closing frontmatter fence", filePath)
	}

	for _, fieldUpdate := range fieldUpdates {
		fileLines, closingFenceIndex = applyFrontmatterFieldUpdate(fileLines, closingFenceIndex, fieldUpdate)
	}

	// Atomic replace, not a truncating WriteFile: the file IS the database, and
	// an interrupted truncate-then-write (crash, disk full) would leave a
	// zero-byte REQ that git cannot restore if it was never committed.
	if writeError := writeFileAtomically(filePath, []byte(strings.Join(fileLines, "\n"))); writeError != nil {
		return fmt.Errorf("writing %s: %w", filePath, writeError)
	}
	return nil
}

// writeFileAtomically writes fileContents to a temporary file in the target's
// directory, then renames it over filePath — so a reader (or a crash) sees
// either the complete old file or the complete new one, never a truncated
// in-between. The dot-prefixed temp name keeps a crash leftover out of the
// board's REQ-*.md walk.
func writeFileAtomically(filePath string, fileContents []byte) error {
	parentDirectory := filepath.Dir(filePath)
	temporaryFile, createError := os.CreateTemp(parentDirectory, "."+filepath.Base(filePath)+".tmp-*")
	if createError != nil {
		return fmt.Errorf("creating temp file in %s: %w", parentDirectory, createError)
	}
	temporaryPath := temporaryFile.Name()
	defer os.Remove(temporaryPath) // no-op once the rename has landed

	if _, writeError := temporaryFile.Write(fileContents); writeError != nil {
		temporaryFile.Close()
		return fmt.Errorf("writing %s: %w", temporaryPath, writeError)
	}
	// CreateTemp opens 0600; match the 0644 the direct write used to produce.
	if chmodError := temporaryFile.Chmod(0o644); chmodError != nil {
		temporaryFile.Close()
		return fmt.Errorf("setting mode on %s: %w", temporaryPath, chmodError)
	}
	if syncError := temporaryFile.Sync(); syncError != nil {
		temporaryFile.Close()
		return fmt.Errorf("syncing %s: %w", temporaryPath, syncError)
	}
	if closeError := temporaryFile.Close(); closeError != nil {
		return fmt.Errorf("closing %s: %w", temporaryPath, closeError)
	}
	if renameError := os.Rename(temporaryPath, filePath); renameError != nil {
		return fmt.Errorf("replacing %s: %w", filePath, renameError)
	}
	return nil
}

// applyFrontmatterFieldUpdate performs one field's upsert/removal against the
// frontmatter lines (indexes 1..closingFenceIndex-1) and returns the updated
// line slice plus the (possibly shifted) closing-fence index.
//
// EVERY occurrence of the key is consumed, not just the first: real REQ files
// carry duplicate top-level keys, and the YAML reader's duplicate-key recovery
// (parseFrontmatterFields) keeps the LAST value — so an edit that touched only
// the first occurrence would look successful yet read back as the untouched
// last value. The replacement value lands once, at the first occurrence's
// position; the duplicates are dropped with their continuation lines.
func applyFrontmatterFieldUpdate(fileLines []string, closingFenceIndex int, fieldUpdate frontmatterFieldUpdate) ([]string, int) {
	type frontmatterLineSpan struct{ startIndex, endIndex int } // [start, end)
	var occurrenceSpans []frontmatterLineSpan
	for lineIndex := 1; lineIndex < closingFenceIndex; {
		keyName, isKeyLine := topLevelKeyName(strings.TrimRight(fileLines[lineIndex], "\r"))
		if !isKeyLine || keyName != fieldUpdate.FieldKey {
			lineIndex++
			continue
		}
		continuationEndIndex := frontmatterContinuationEnd(fileLines, lineIndex+1, closingFenceIndex)
		occurrenceSpans = append(occurrenceSpans, frontmatterLineSpan{lineIndex, continuationEndIndex})
		lineIndex = continuationEndIndex
	}

	if len(occurrenceSpans) == 0 {
		if fieldUpdate.RemoveField {
			return fileLines, closingFenceIndex // nothing to remove
		}
		insertedLine := fieldUpdate.FieldKey + ": " + fieldUpdate.FieldValue
		fileLines = append(fileLines[:closingFenceIndex],
			append([]string{insertedLine}, fileLines[closingFenceIndex:]...)...)
		return fileLines, closingFenceIndex + 1
	}

	rebuiltLines := make([]string, 0, len(fileLines))
	replacementInserted := false
	spanCursor := 0
	for lineIndex := 0; lineIndex < len(fileLines); lineIndex++ {
		if spanCursor < len(occurrenceSpans) && lineIndex == occurrenceSpans[spanCursor].startIndex {
			if !fieldUpdate.RemoveField && !replacementInserted {
				rebuiltLines = append(rebuiltLines, fieldUpdate.FieldKey+": "+fieldUpdate.FieldValue)
				replacementInserted = true
			}
			lineIndex = occurrenceSpans[spanCursor].endIndex - 1 // loop increment lands on endIndex
			spanCursor++
			continue
		}
		rebuiltLines = append(rebuiltLines, fileLines[lineIndex])
	}

	removedLineCount := 0
	for _, span := range occurrenceSpans {
		removedLineCount += span.endIndex - span.startIndex
	}
	newClosingFenceIndex := closingFenceIndex - removedLineCount
	if replacementInserted {
		newClosingFenceIndex++
	}
	return rebuiltLines, newClosingFenceIndex
}

// frontmatterContinuationEnd returns the index of the first line at or after
// startIndex that is NOT a continuation of the preceding key's value — i.e. the
// first blank line, comment, top-level key, or the closing fence. Indented
// lines and list items belong to the key above them (block scalars and block
// lists), so a replace/remove takes them along instead of orphaning them into
// the neighboring field.
func frontmatterContinuationEnd(fileLines []string, startIndex int, closingFenceIndex int) int {
	for lineIndex := startIndex; lineIndex < closingFenceIndex; lineIndex++ {
		line := strings.TrimRight(fileLines[lineIndex], "\r")
		if line == "" || strings.HasPrefix(line, "#") {
			return lineIndex
		}
		firstByte := line[0]
		if firstByte != ' ' && firstByte != '\t' && firstByte != '-' {
			return lineIndex
		}
	}
	return closingFenceIndex
}

// encodeYamlDoubleQuotedScalar renders raw user text as a single-line YAML
// double-quoted scalar, so a tester name or a multi-paragraph feedback note
// always fits the one-line-per-placeholder model that upsertFrontmatterFields
// maintains. Newlines survive as \n escapes (yaml.v3 restores them on read);
// carriage returns are dropped and other control characters become spaces.
func encodeYamlDoubleQuotedScalar(rawText string) string {
	var encodedBuilder strings.Builder
	encodedBuilder.WriteByte('"')
	for _, character := range rawText {
		switch character {
		case '\\':
			encodedBuilder.WriteString(`\\`)
		case '"':
			encodedBuilder.WriteString(`\"`)
		case '\n':
			encodedBuilder.WriteString(`\n`)
		case '\t':
			encodedBuilder.WriteString(`\t`)
		case '\r':
			// dropped — CRLF feedback reads back as plain LF
		default:
			if character < 0x20 || character == 0x7f {
				encodedBuilder.WriteByte(' ')
			} else {
				encodedBuilder.WriteRune(character)
			}
		}
	}
	encodedBuilder.WriteByte('"')
	return encodedBuilder.String()
}

// testingClearState is the API's explicit "undo the testing track" action —
// it removes every testing placeholder rather than setting a status.
const testingClearState = "clear"

// buildTestingFieldUpdates maps a validated testing transition onto the
// placeholder mutations to apply. `returned` carries the feedback; the other
// set-states drop any stale feedback line (its history lives in git);
// `clear` removes the whole track.
func buildTestingFieldUpdates(testingState string, testedBy string, feedbackText string, updateInstant time.Time) []frontmatterFieldUpdate {
	if testingState == testingClearState {
		return []frontmatterFieldUpdate{
			{FieldKey: "testing_status", RemoveField: true},
			{FieldKey: "tested_by", RemoveField: true},
			{FieldKey: "testing_updated_at", RemoveField: true},
			{FieldKey: "testing_feedback", RemoveField: true},
		}
	}
	fieldUpdates := []frontmatterFieldUpdate{
		{FieldKey: "testing_status", FieldValue: testingState},
		{FieldKey: "tested_by", FieldValue: encodeYamlDoubleQuotedScalar(testedBy)},
		{FieldKey: "testing_updated_at", FieldValue: updateInstant.UTC().Format(time.RFC3339)},
	}
	if testingState == testingStatusReturned {
		fieldUpdates = append(fieldUpdates, frontmatterFieldUpdate{
			FieldKey: "testing_feedback", FieldValue: encodeYamlDoubleQuotedScalar(feedbackText),
		})
	} else {
		fieldUpdates = append(fieldUpdates, frontmatterFieldUpdate{
			FieldKey: "testing_feedback", RemoveField: true,
		})
	}
	return fieldUpdates
}
