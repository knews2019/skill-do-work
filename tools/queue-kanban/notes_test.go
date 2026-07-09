package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// writeNotesTree seeds a repo root containing do-work/ with one pending REQ (so
// the board builds at all) and, when notesContent is non-empty, a top-level
// do-work/notes.md carrying it.
func writeNotesTree(t *testing.T, notesContent string) string {
	t.Helper()
	repoRoot := t.TempDir()

	writeFixture := func(relativePath string, content string) {
		absolutePath := filepath.Join(repoRoot, relativePath)
		if mkdirError := os.MkdirAll(filepath.Dir(absolutePath), 0o755); mkdirError != nil {
			t.Fatalf("mkdir %s: %v", relativePath, mkdirError)
		}
		if writeError := os.WriteFile(absolutePath, []byte(content), 0o644); writeError != nil {
			t.Fatalf("write %s: %v", relativePath, writeError)
		}
	}

	writeFixture(filepath.Join("do-work", "queue", "REQ-9001-pending.md"),
		"---\nid: REQ-9001\ntitle: Fixture\nstatus: pending\n---\n\nBody.\n")
	if notesContent != "" {
		writeFixture(filepath.Join("do-work", "notes.md"), notesContent)
	}
	return repoRoot
}

// TestParseQueueNoteLineShapes covers the canonical `- [date] text` shape plus
// the hand-edited drift the file tolerates: a missing bullet, a missing date, a
// bracketed prefix that is not a date, and a task-list marker.
func TestParseQueueNoteLineShapes(t *testing.T) {
	testCases := []struct {
		inputLine    string
		expectedDate string
		expectedText string
	}{
		{"- [2026-07-09] check the retry budget", "2026-07-09", "check the retry budget"},
		{"* [2026-01-02] star bullet still parses", "2026-01-02", "star bullet still parses"},
		{"[2026-07-09] no bullet at all", "2026-07-09", "no bullet at all"},
		{"- undated hint with no bracket", "", "undated hint with no bracket"},
		{"- [not-a-date] keeps its bracket text", "", "[not-a-date] keeps its bracket text"},
		{"- [ ] a task marker is not a date", "", "[ ] a task marker is not a date"},
		{"- [2026-07-09] revisit REQ-1 [after] the merge", "2026-07-09", "revisit REQ-1 [after] the merge"},
	}
	for _, testCase := range testCases {
		parsedNote := parseQueueNoteLine(testCase.inputLine)
		if parsedNote.NoteDate != testCase.expectedDate || parsedNote.NoteText != testCase.expectedText {
			t.Errorf("parseQueueNoteLine(%q) = {%q, %q}, want {%q, %q}",
				testCase.inputLine, parsedNote.NoteDate, parsedNote.NoteText,
				testCase.expectedDate, testCase.expectedText)
		}
	}
}

// TestLoadQueueNotesPreservesAppendOrder asserts blank lines are skipped and the
// surviving lines keep file order — `do-work note` only appends, so file order
// is chronological order and must not be sorted or deduped.
func TestLoadQueueNotesPreservesAppendOrder(t *testing.T) {
	repoRoot := writeNotesTree(t, "- [2026-07-01] first\n\n- [2026-06-01] second, older, still second\n\n")
	notes := loadQueueNotes(filepath.Join(repoRoot, "do-work", "notes.md"))

	if len(notes) != 2 {
		t.Fatalf("loadQueueNotes returned %d notes, want 2", len(notes))
	}
	if notes[0].NoteText != "first" || notes[1].NoteText != "second, older, still second" {
		t.Fatalf("notes out of append order: %+v", notes)
	}
}

// TestLoadQueueNotesMissingFileIsNotAnError asserts an absent notes.md yields no
// notes rather than failing — every do-work repo predates the note command.
func TestLoadQueueNotesMissingFileIsNotAnError(t *testing.T) {
	if notes := loadQueueNotes(""); notes != nil {
		t.Fatalf("loadQueueNotes(\"\") = %+v, want nil", notes)
	}
	if notes := loadQueueNotes(filepath.Join(t.TempDir(), "absent.md")); notes != nil {
		t.Fatalf("loadQueueNotes(absent path) = %+v, want nil", notes)
	}
}

// TestEnumerateFindsOnlyTopLevelNotes asserts the walker binds notes.md at the
// top level of do-work/ and ignores a same-named scratch file nested under a UR.
func TestEnumerateFindsOnlyTopLevelNotes(t *testing.T) {
	repoRoot := writeNotesTree(t, "- [2026-07-09] top level\n")
	nestedNotesPath := filepath.Join(repoRoot, "do-work", "archive", "UR-100", "notes.md")
	if mkdirError := os.MkdirAll(filepath.Dir(nestedNotesPath), 0o755); mkdirError != nil {
		t.Fatalf("mkdir nested: %v", mkdirError)
	}
	if writeError := os.WriteFile(nestedNotesPath, []byte("- [2026-07-09] nested scratch\n"), 0o644); writeError != nil {
		t.Fatalf("write nested: %v", writeError)
	}

	discovered, enumerateError := enumerateDoWorkTree(repoRoot)
	if enumerateError != nil {
		t.Fatalf("enumerateDoWorkTree: %v", enumerateError)
	}
	expectedPath := filepath.Join(repoRoot, "do-work", "notes.md")
	if discovered.NotesFilePath != expectedPath {
		t.Fatalf("NotesFilePath = %q, want %q", discovered.NotesFilePath, expectedPath)
	}
}

// TestBoardCarriesNotesIntoGeneratedData walks the whole path a note travels —
// tree walk → board model → JSON data island — and asserts the note text lands
// in board-data.js, where board.js reads it.
func TestBoardCarriesNotesIntoGeneratedData(t *testing.T) {
	repoRoot := writeNotesTree(t, "- [2026-07-09] check the retry budget\n")
	board, buildError := buildBoard(repoRoot, time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC), defaultRecentWindow, nil)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	if len(board.Notes) != 1 || board.Notes[0].NoteDate != "2026-07-09" {
		t.Fatalf("board.Notes = %+v, want one dated note", board.Notes)
	}

	boardData, projectError := buildGeneratedBoardData(board)
	if projectError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", projectError)
	}
	encoded, encodeError := encodeBoardDataForJsAssignment(boardData)
	if encodeError != nil {
		t.Fatalf("encodeBoardDataForJsAssignment: %v", encodeError)
	}
	if !strings.Contains(encoded, `"text":"check the retry budget"`) {
		t.Fatalf("board-data.js is missing the note text:\n%s", encoded)
	}
}

// TestNotesOmittedFromDataWhenAbsent asserts a repo with no notes.md emits no
// `notes` key at all, so board.js leaves the strip hidden.
func TestNotesOmittedFromDataWhenAbsent(t *testing.T) {
	repoRoot := writeNotesTree(t, "")
	board, buildError := buildBoard(repoRoot, time.Now(), defaultRecentWindow, nil)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	boardData, projectError := buildGeneratedBoardData(board)
	if projectError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", projectError)
	}
	encoded, encodeError := encodeBoardDataForJsAssignment(boardData)
	if encodeError != nil {
		t.Fatalf("encodeBoardDataForJsAssignment: %v", encodeError)
	}
	if strings.Contains(encoded, `"notes"`) {
		t.Fatalf("board-data.js carries a notes key with no notes.md present:\n%s", encoded)
	}
}

// TestServeFingerprintTracksNotesFile asserts appending a note invalidates the
// live server's mtime cache. Without notes.md in the fingerprint, `serve` would
// keep replaying the stale strip until an unrelated REQ file changed.
func TestServeFingerprintTracksNotesFile(t *testing.T) {
	repoRoot := writeNotesTree(t, "- [2026-07-09] first\n")
	notesPath := filepath.Join(repoRoot, "do-work", "notes.md")

	discovered, enumerateError := enumerateDoWorkTree(repoRoot)
	if enumerateError != nil {
		t.Fatalf("enumerateDoWorkTree: %v", enumerateError)
	}
	beforeFingerprint := buildTreeMtimeFingerprint(discovered)
	if _, tracked := beforeFingerprint[notesPath]; !tracked {
		t.Fatalf("fingerprint does not track %s", notesPath)
	}

	// Stamp a distinctly later mtime rather than relying on filesystem clock
	// granularity — HFS+/ext3 truncate to the second, so a same-second rewrite
	// would leave the mtime unchanged and make this test flaky.
	if writeError := os.WriteFile(notesPath, []byte("- [2026-07-09] first\n- [2026-07-10] second\n"), 0o644); writeError != nil {
		t.Fatalf("append note: %v", writeError)
	}
	laterInstant := time.Now().Add(2 * time.Second)
	if chtimesError := os.Chtimes(notesPath, laterInstant, laterInstant); chtimesError != nil {
		t.Fatalf("chtimes: %v", chtimesError)
	}

	afterFingerprint := buildTreeMtimeFingerprint(discovered)
	if treeMtimeFingerprintsEqual(beforeFingerprint, afterFingerprint) {
		t.Fatal("fingerprint unchanged after appending a note — serve would cache a stale board")
	}
}
