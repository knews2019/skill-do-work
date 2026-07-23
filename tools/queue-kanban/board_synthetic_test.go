package main

import (
	"os"
	"path/filepath"
	"strings"
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
//	do-work/archive/UR-100/REQ-9005-cancelled.md           cancelled (completed_at inside the window), grouped under UR-100
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
	writeFixture(filepath.Join("do-work", "queue", "REQ-9006-blocked.md"),
		requestContent("REQ-9006", "blocked", "blocked_by: \"LM Studio running locally\"\nblocked_at: 2026-06-28T10:00:00Z\nblocked_check: \"curl -sf http://localhost:1234/v1/models\"\n"))

	writeFixture(filepath.Join("do-work", "archive", "UR-100", "input.md"), userRequestContent("UR-100"))
	writeFixture(filepath.Join("do-work", "archive", "UR-100", "REQ-9003-flat.md"),
		requestContent("REQ-9003", "complete", "user_request: UR-100\ncommit_hash: deadbeef\n"))
	writeFixture(filepath.Join("do-work", "archive", "UR-100", "REQ-9005-cancelled.md"),
		requestContent("REQ-9005", "cancelled", "user_request: UR-100\ncompleted_at: 2026-06-28T10:00:00Z\n"))

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
	if !columnContainsRequestId(board.Columns.NeedsInputOrBlocked, "REQ-9006") {
		t.Fatalf("REQ-9006 (blocked on external condition) missing from the Needs-input/Blocked column")
	}
	blockedTicket := board.RequestsById["REQ-9006"]
	if blockedTicket == nil || blockedTicket.StatusUnrecognized {
		t.Fatalf("REQ-9006 (blocked) must be a recognized status, never flagged unrecognized — got %+v", blockedTicket)
	}
	if blockedTicket != nil && (len(blockedTicket.BlockedBy) != 1 || blockedTicket.BlockedCheck == "") {
		t.Fatalf("REQ-9006 blocked fields not parsed: BlockedBy=%v BlockedCheck=%q", blockedTicket.BlockedBy, blockedTicket.BlockedCheck)
	}
	if !columnContainsRequestId(board.Columns.RecentlyDone, "REQ-9005") {
		t.Fatalf("REQ-9005 (cancelled, completed_at inside the window) missing from the Recently-done column")
	}
	if columnContainsRequestId(board.Columns.NeedsInputOrBlocked, "REQ-9005") {
		t.Fatalf("REQ-9005 (cancelled) must not land in Needs-input/Blocked — cancelled is a recognized terminal status")
	}
	if columnContainsRequestId(board.Columns.RecentlyDone, "REQ-9004") {
		t.Fatalf("REQ-9004 (completed 2026-06-10) is outside the recent window and belongs to the calendar only")
	}
	if len(board.Warnings) != 0 {
		t.Fatalf("synthetic tree should produce no data warnings, got %v", board.Warnings)
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

// TestSyntheticUnrecognizedStatusFlagged pins the off-vocabulary contract
// (found via REQ-950's review feedback): a status outside the Schema Read
// Contract set is parked in Needs input / Blocked, flagged StatusUnrecognized
// for the frontend's invalid-status highlight, and produces a warning that
// carries the fix prompt — while recognized statuses stay unflagged. It seeds
// its own tree because the shared synthetic tree asserts zero warnings.
func TestSyntheticUnrecognizedStatusFlagged(t *testing.T) {
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
	writeFixture(filepath.Join("do-work", "queue", "REQ-9101-pending.md"),
		"---\nid: REQ-9101\ntitle: Fixture REQ-9101\nstatus: pending\n---\n\nBody.\n")
	writeFixture(filepath.Join("do-work", "working", "REQ-9102-off-vocab.md"),
		"---\nid: REQ-9102\ntitle: Fixture REQ-9102\nstatus: in-progress\n---\n\nBody.\n")

	stubGitLookup := func(string, string) (time.Time, bool) { return time.Time{}, false }
	fixedNow := time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC)
	board, buildError := buildBoard(repoRoot, fixedNow, 7*24*time.Hour, stubGitLookup)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}

	offVocabTicket := board.RequestsById["REQ-9102"]
	if offVocabTicket == nil {
		t.Fatalf("REQ-9102 not parsed")
	}
	if !offVocabTicket.StatusUnrecognized {
		t.Fatalf("REQ-9102 (status in-progress) should be flagged StatusUnrecognized")
	}
	if !columnContainsRequestId(board.Columns.NeedsInputOrBlocked, "REQ-9102") {
		t.Fatalf("REQ-9102 (unrecognized status) must be parked in Needs-input/Blocked, never dropped")
	}

	pendingTicket := board.RequestsById["REQ-9101"]
	if pendingTicket == nil || pendingTicket.StatusUnrecognized {
		t.Fatalf("REQ-9101 (recognized status pending) must not be flagged StatusUnrecognized")
	}

	sawFixPromptWarning := false
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "REQ-9102") &&
			strings.Contains(warningText, `"in-progress"`) &&
			strings.Contains(warningText, "do-work forensics") {
			sawFixPromptWarning = true
		}
	}
	if !sawFixPromptWarning {
		t.Fatalf("expected a warning naming REQ-9102, its status, and the fix prompt; got %v", board.Warnings)
	}
}

// TestSyntheticStrayRequestFlagged reproduces the invisible-REQ bug: a work
// agent archived a completed REQ to do-work/user-requests/UR-NNN/ instead of
// do-work/archive/, so the board (which buckets only queue/working/archive
// files) rendered no card. The REQ must never be silently dropped — a data
// warning naming the id and its location must fire, and it must NOT sneak into
// AllRequests as if it were a real card.
func TestSyntheticStrayRequestFlagged(t *testing.T) {
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
	writeFixture(filepath.Join("do-work", "queue", "REQ-9201-pending.md"),
		"---\nid: REQ-9201\ntitle: Fixture REQ-9201\nstatus: pending\n---\n\nBody.\n")
	// The misplaced REQ — under user-requests/, which the walk visits (for
	// input.md) but which bucketing never scans for REQ cards.
	writeFixture(filepath.Join("do-work", "user-requests", "UR-301", "REQ-1213.md"),
		"---\nid: REQ-1213\ntitle: Fixture REQ-1213\nstatus: completed\n---\n\nBody.\n")

	stubGitLookup := func(string, string) (time.Time, bool) { return time.Time{}, false }
	fixedNow := time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC)
	board, buildError := buildBoard(repoRoot, fixedNow, 7*24*time.Hour, stubGitLookup)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}

	if _, snuckIn := board.RequestsById["REQ-1213"]; snuckIn {
		t.Fatalf("stray REQ-1213 must not be parsed into the board as a card")
	}

	sawStrayWarning := false
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "REQ-1213") &&
			strings.Contains(warningText, "user-requests/UR-301/REQ-1213.md") &&
			strings.Contains(warningText, "invisible") {
			sawStrayWarning = true
		}
	}
	if !sawStrayWarning {
		t.Fatalf("expected a warning naming REQ-1213, its location, and that it is invisible; got %v", board.Warnings)
	}
}

func TestSyntheticCountsAndCalendar(t *testing.T) {
	board := syntheticBoard(t)
	if got := len(board.AllRequests); got != 6 {
		t.Fatalf("AllRequests = %d, want 6", got)
	}
	archivedCompleted := 0
	archivedResolved := 0
	for _, ticket := range board.AllRequests {
		if ticket.TreeSection != "archive" {
			continue
		}
		if isCompletedStatus(ticket.Status) {
			archivedCompleted++
		}
		if isTerminalResolvedStatus(ticket.Status) {
			archivedResolved++
		}
	}
	if archivedCompleted != 2 {
		t.Fatalf("archived completed = %d, want 2 (cancelled must NOT count as terminal success)", archivedCompleted)
	}
	if archivedResolved != 3 {
		t.Fatalf("archived terminally resolved = %d, want 3 (completed pair + cancelled)", archivedResolved)
	}
	if got := len(board.Calendar); got != 3 {
		t.Fatalf("calendar entries = %d, want 3 (both completed REQs and the cancelled REQ resolve a completion time)", got)
	}
}
