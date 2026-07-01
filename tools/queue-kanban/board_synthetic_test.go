package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// createSyntheticDoWorkTree builds a deterministic, repo-independent do-work tree
// in a temp dir and returns the repo root (the directory holding do-work/). It
// seeds exactly the shapes the old live tests asserted against the source
// monorepo (REQ-1207 / UR-448 / banded archive / >=900 tickets) so those
// exact-parse behaviours stay covered WITHOUT depending on this repo's actual
// queue contents:
//
//	do-work/queue/REQ-9001-pending.md                      status pending
//	do-work/working/REQ-9002-claimed.md                    status claimed
//	do-work/archive/UR-100/input.md                        flat UR
//	do-work/archive/UR-100/REQ-9003-flat.md                legacy "complete" → completed, grouped under UR-100
//	do-work/archive/UR-200-209/UR-205/input.md             banded UR
//	do-work/archive/UR-200-209/UR-205/REQ-9004-banded.md   completed (completed_at), grouped under UR-205
func createSyntheticDoWorkTree(t *testing.T) string {
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

	requestContent := func(requestId string, status string, extraFrontmatter string) string {
		return "---\nid: " + requestId + "\ntitle: Fixture " + requestId +
			"\nstatus: " + status + "\n" + extraFrontmatter + "---\n\n## What\n\nBody for " + requestId + ".\n"
	}
	userRequestContent := func(userRequestId string) string {
		return "---\nid: " + userRequestId + "\ntitle: Fixture user request " + userRequestId + "\n---\n\nOriginal request text.\n"
	}

	writeFixture(filepath.Join("do-work", "queue", "REQ-9001-pending.md"), requestContent("REQ-9001", "pending", ""))
	writeFixture(filepath.Join("do-work", "working", "REQ-9002-claimed.md"), requestContent("REQ-9002", "claimed", ""))

	writeFixture(filepath.Join("do-work", "archive", "UR-100", "input.md"), userRequestContent("UR-100"))
	writeFixture(filepath.Join("do-work", "archive", "UR-100", "REQ-9003-flat.md"),
		requestContent("REQ-9003", "complete", "user_request: UR-100\ncommit_hash: deadbeef\n"))

	writeFixture(filepath.Join("do-work", "archive", "UR-200-209", "UR-205", "input.md"), userRequestContent("UR-205"))
	writeFixture(filepath.Join("do-work", "archive", "UR-200-209", "UR-205", "REQ-9004-banded.md"),
		requestContent("REQ-9004", "completed", "user_request: UR-205\ncompleted_at: 2026-06-10T14:00:00Z\n"))

	return repoRoot
}

// syntheticBoard builds the board from the synthetic tree with a stubbed git
// lookup (resolving the seeded "deadbeef" commit) and a fixed `now`, so every
// assertion is deterministic and never shells out to git.
func syntheticBoard(t *testing.T) *Board {
	t.Helper()
	repoRoot := createSyntheticDoWorkTree(t)
	gitCommitTime := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)
	stubGitLookup := func(_ string, commitHash string) (time.Time, bool) {
		if commitHash == "deadbeef" {
			return gitCommitTime, true
		}
		return time.Time{}, false
	}
	fixedNow := time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC)
	board, buildError := buildBoard(repoRoot, fixedNow, 7*24*time.Hour, stubGitLookup)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	return board
}

// columnContainsRequestId reports whether a board column holds the given REQ id.
func columnContainsRequestId(column []*RequestTicket, requestId string) bool {
	for _, ticket := range column {
		if ticket.RequestId == requestId {
			return true
		}
	}
	return false
}

func TestSyntheticParsesBothArchiveShapes(t *testing.T) {
	board := syntheticBoard(t)
	sawBanded := false
	sawFlat := false
	for _, ticket := range board.AllRequests {
		if ticket.TreeSection != "archive" {
			continue
		}
		if pathHasBandedArchiveSegment(ticket.FilePath) {
			sawBanded = true
		}
		if pathIsFlatArchiveRequest(ticket.FilePath) {
			sawFlat = true
		}
	}
	if !sawBanded {
		t.Fatalf("banded archive/UR-NNN-MMM/ REQ was not parsed from the synthetic tree")
	}
	if !sawFlat {
		t.Fatalf("flat archive/UR-NNN/ REQ was not parsed from the synthetic tree")
	}
}

func TestSyntheticUserRequestLinkage(t *testing.T) {
	board := syntheticBoard(t)
	groupings := []struct {
		userRequestId string
		requestId     string
	}{
		{"UR-100", "REQ-9003"},
		{"UR-205", "REQ-9004"},
	}
	for _, grouping := range groupings {
		userRequest := board.UserRequestsById[grouping.userRequestId]
		if userRequest == nil {
			t.Fatalf("%s not present in the board", grouping.userRequestId)
		}
		if !stringSliceContains(userRequest.RequestIds, grouping.requestId) {
			t.Fatalf("%s should group %s, got %v", grouping.userRequestId, grouping.requestId, userRequest.RequestIds)
		}
	}
}

func TestSyntheticColumnBucketing(t *testing.T) {
	board := syntheticBoard(t)
	if !columnContainsRequestId(board.Columns.Pending, "REQ-9001") {
		t.Fatalf("REQ-9001 (pending) missing from the Pending column")
	}
	if !columnContainsRequestId(board.Columns.Claimed, "REQ-9002") {
		t.Fatalf("REQ-9002 (claimed) missing from the Claimed column")
	}
}

func TestSyntheticLegacyCompleteNormalized(t *testing.T) {
	board := syntheticBoard(t)
	ticket := board.RequestsById["REQ-9003"]
	if ticket == nil {
		t.Fatalf("REQ-9003 not parsed")
	}
	if ticket.OriginalStatus != "complete" {
		t.Fatalf("REQ-9003 OriginalStatus = %q, want complete", ticket.OriginalStatus)
	}
	if ticket.Status != "completed" {
		t.Fatalf("legacy 'complete' on REQ-9003 normalized to %q, want completed", ticket.Status)
	}
}

func TestSyntheticCountsAndCalendar(t *testing.T) {
	board := syntheticBoard(t)
	if got := len(board.AllRequests); got != 4 {
		t.Fatalf("AllRequests = %d, want 4", got)
	}
	archivedCompleted := 0
	for _, ticket := range board.AllRequests {
		if ticket.TreeSection == "archive" && isCompletedStatus(ticket.Status) {
			archivedCompleted++
		}
	}
	if archivedCompleted != 2 {
		t.Fatalf("archived completed = %d, want 2", archivedCompleted)
	}
	if got := len(board.Calendar); got != 2 {
		t.Fatalf("calendar entries = %d, want 2 (both completed REQs resolve a completion time)", got)
	}
}
