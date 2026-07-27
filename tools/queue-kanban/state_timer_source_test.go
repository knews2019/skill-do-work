package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// stateTimerSyntheticBoard builds a board over a tree seeded with the timestamp
// sources the pending-card state timer resolves from:
//
//	REQ-9501  pending, status_changed_at stamped (clarify-style flip)  → payload carries it
//	REQ-9502  pending, no stamp, mtime forced hours past created_at    → payload carries fileModifiedAt
//
// The frontend's stateTimerSpecFor prefers statusChangedAt, then the later of
// createdAt/fileModifiedAt — this test pins the Go side of that contract: the
// fields parse and survive into the JSON payload.
func stateTimerSyntheticBoard(t *testing.T) *Board {
	t.Helper()
	repoRoot := t.TempDir()

	writeFixture := func(relativePath string, content string) string {
		absolutePath := filepath.Join(repoRoot, relativePath)
		if mkdirError := os.MkdirAll(filepath.Dir(absolutePath), 0o755); mkdirError != nil {
			t.Fatalf("mkdir %s: %v", relativePath, mkdirError)
		}
		if writeError := os.WriteFile(absolutePath, []byte(content), 0o644); writeError != nil {
			t.Fatalf("write %s: %v", relativePath, writeError)
		}
		return absolutePath
	}
	requestContent := func(requestId string, extraFrontmatter string) string {
		return "---\nid: " + requestId + "\ntitle: Fixture " + requestId +
			"\nstatus: pending\ncreated_at: 2026-06-29T10:00:00Z\n" + extraFrontmatter +
			"---\n\nBody for " + requestId + ".\n"
	}

	writeFixture(filepath.Join("do-work", "queue", "REQ-9501-stamped-flip.md"),
		requestContent("REQ-9501", "status_changed_at: 2026-06-30T09:30:00Z\n"))
	editedPath := writeFixture(filepath.Join("do-work", "queue", "REQ-9502-edited-later.md"),
		requestContent("REQ-9502", ""))
	forcedModTime := time.Date(2026, 6, 30, 8, 0, 0, 0, time.UTC)
	if chtimesError := os.Chtimes(editedPath, forcedModTime, forcedModTime); chtimesError != nil {
		t.Fatalf("chtimes %s: %v", editedPath, chtimesError)
	}

	neverCalledGitLookup := func(string, string) (time.Time, bool) { return time.Time{}, false }
	fixedNow := time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC)
	board, buildError := buildBoard(repoRoot, fixedNow, 7*24*time.Hour, neverCalledGitLookup)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	return board
}

func TestStatusChangedAtAndFileMtimeParsedIntoModel(t *testing.T) {
	board := stateTimerSyntheticBoard(t)

	stampedTicket := board.RequestsById["REQ-9501"]
	if stampedTicket == nil {
		t.Fatal("REQ-9501 not parsed")
	}
	if stampedTicket.StatusChangedAt != "2026-06-30T09:30:00Z" {
		t.Fatalf("REQ-9501 StatusChangedAt = %q, want the stamped instant", stampedTicket.StatusChangedAt)
	}

	editedTicket := board.RequestsById["REQ-9502"]
	if editedTicket == nil {
		t.Fatal("REQ-9502 not parsed")
	}
	if editedTicket.StatusChangedAt != "" {
		t.Fatalf("REQ-9502 StatusChangedAt = %q, want empty", editedTicket.StatusChangedAt)
	}
	expectedModTime := time.Date(2026, 6, 30, 8, 0, 0, 0, time.UTC)
	if !editedTicket.FileModifiedAt.Equal(expectedModTime) {
		t.Fatalf("REQ-9502 FileModifiedAt = %v, want %v", editedTicket.FileModifiedAt, expectedModTime)
	}
}

func TestStatusChangedAtAndFileMtimeSurviveIntoPayload(t *testing.T) {
	board := stateTimerSyntheticBoard(t)

	generatedData, buildError := buildGeneratedBoardData(board)
	if buildError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", buildError)
	}

	stampedPayload := generatedData.Requests["REQ-9501"]
	if stampedPayload.StatusChangedAt != "2026-06-30T09:30:00Z" {
		t.Fatalf("REQ-9501 payload statusChangedAt = %q, want the stamped instant", stampedPayload.StatusChangedAt)
	}

	editedPayload := generatedData.Requests["REQ-9502"]
	if editedPayload.FileModifiedAt != "2026-06-30T08:00:00Z" {
		t.Fatalf("REQ-9502 payload fileModifiedAt = %q, want the forced mtime as RFC3339 UTC", editedPayload.FileModifiedAt)
	}
	if stampedPayload.FileModifiedAt == "" {
		t.Fatal("REQ-9501 payload fileModifiedAt empty — mtime should be captured for every parsed REQ")
	}
}
