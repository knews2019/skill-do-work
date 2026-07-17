package main

import (
	"fmt"
	"os"
	"strings"
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
func appendTestingProfile(testersFilePath string, rawProfileName string) ([]string, error) {
	profileName, validateError := validateTesterProfileName(rawProfileName)
	if validateError != nil {
		return nil, validateError
	}

	existingProfiles := loadTestingProfiles(testersFilePath)
	for _, existingProfile := range existingProfiles {
		if strings.EqualFold(existingProfile, profileName) {
			return existingProfiles, nil
		}
	}

	fileHandle, openError := os.OpenFile(testersFilePath, os.O_CREATE|os.O_WRONLY, 0o644)
	if openError != nil {
		return nil, fmt.Errorf("opening testers file: %w", openError)
	}
	defer fileHandle.Close()

	fileInfo, statError := fileHandle.Stat()
	if statError != nil {
		return nil, fmt.Errorf("stat testers file: %w", statError)
	}
	if _, seekError := fileHandle.Seek(0, 2); seekError != nil {
		return nil, fmt.Errorf("seeking testers file: %w", seekError)
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

	if writeError := os.WriteFile(filePath, []byte(strings.Join(fileLines, "\n")), 0o644); writeError != nil {
		return fmt.Errorf("writing %s: %w", filePath, writeError)
	}
	return nil
}

// applyFrontmatterFieldUpdate performs one field's upsert/removal against the
// frontmatter lines (indexes 1..closingFenceIndex-1) and returns the updated
// line slice plus the (possibly shifted) closing-fence index.
func applyFrontmatterFieldUpdate(fileLines []string, closingFenceIndex int, fieldUpdate frontmatterFieldUpdate) ([]string, int) {
	for lineIndex := 1; lineIndex < closingFenceIndex; lineIndex++ {
		keyName, isKeyLine := topLevelKeyName(strings.TrimRight(fileLines[lineIndex], "\r"))
		if !isKeyLine || keyName != fieldUpdate.FieldKey {
			continue
		}
		continuationEndIndex := frontmatterContinuationEnd(fileLines, lineIndex+1, closingFenceIndex)
		removedLineCount := continuationEndIndex - lineIndex
		if fieldUpdate.RemoveField {
			fileLines = append(fileLines[:lineIndex], fileLines[continuationEndIndex:]...)
			return fileLines, closingFenceIndex - removedLineCount
		}
		replacementLine := fieldUpdate.FieldKey + ": " + fieldUpdate.FieldValue
		fileLines = append(fileLines[:lineIndex+1], fileLines[continuationEndIndex:]...)
		fileLines[lineIndex] = replacementLine
		return fileLines, closingFenceIndex - (removedLineCount - 1)
	}

	if fieldUpdate.RemoveField {
		return fileLines, closingFenceIndex // nothing to remove
	}
	insertedLine := fieldUpdate.FieldKey + ": " + fieldUpdate.FieldValue
	fileLines = append(fileLines[:closingFenceIndex],
		append([]string{insertedLine}, fileLines[closingFenceIndex:]...)...)
	return fileLines, closingFenceIndex + 1
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
