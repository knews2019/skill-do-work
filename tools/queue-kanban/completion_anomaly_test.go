package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// anomalySyntheticBoard builds a board over a tree seeded with every
// completion-anomaly shape plus a healthy dated-but-old control:
//
//	REQ-9301  completed, no completed_at, no commit hash          → anomaly (nothing to resolve)
//	REQ-9302  completed, completed_at unparseable, no commit hash → anomaly (bad completed_at)
//	REQ-9303  cancelled, no completed_at, commit git can't find   → anomaly (bad commit hash)
//	REQ-9304  completed, completed_at parses but is months old    → NOT an anomaly, NOT recently done
//
// The git lookup is stubbed to fail for every hash, and `now` is fixed, so the
// assertions are deterministic and never shell out.
func anomalySyntheticBoard(t *testing.T) *Board {
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
			"\nstatus: " + status + "\n" + extraFrontmatter + "---\n\nBody for " + requestId + ".\n"
	}

	writeFixture(filepath.Join("do-work", "archive", "REQ-9301-no-instant.md"),
		requestContent("REQ-9301", "completed", ""))
	writeFixture(filepath.Join("do-work", "archive", "REQ-9302-bad-timestamp.md"),
		requestContent("REQ-9302", "completed", "completed_at: not-a-real-instant\n"))
	writeFixture(filepath.Join("do-work", "archive", "REQ-9303-bad-commit.md"),
		requestContent("REQ-9303", "cancelled", "commit: badc0ffee\n"))
	writeFixture(filepath.Join("do-work", "archive", "REQ-9304-dated-old.md"),
		requestContent("REQ-9304", "completed", "completed_at: 2026-01-05T10:00:00Z\n"))

	alwaysFailingGitLookup := func(string, string) (time.Time, bool) { return time.Time{}, false }
	fixedNow := time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC)
	board, buildError := buildBoard(repoRoot, fixedNow, 7*24*time.Hour, alwaysFailingGitLookup)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	return board
}

func TestCompletionAnomaliesFlaggedInBoardModel(t *testing.T) {
	board := anomalySyntheticBoard(t)

	expectedReasonFragments := map[string]string{
		"REQ-9301": "terminal status but no completed_at and no resolvable commit hash",
		"REQ-9302": `completed_at "not-a-real-instant" does not parse`,
		"REQ-9303": `commit "badc0ffee" cannot be resolved by git`,
	}
	for requestId, reasonFragment := range expectedReasonFragments {
		ticket := board.RequestsById[requestId]
		if ticket == nil {
			t.Fatalf("%s not parsed", requestId)
		}
		if !ticket.CompletionAnomaly {
			t.Fatalf("%s should be flagged CompletionAnomaly", requestId)
		}
		if !strings.Contains(ticket.CompletionAnomalyReason, reasonFragment) {
			t.Fatalf("%s anomaly reason = %q, want it to contain %q",
				requestId, ticket.CompletionAnomalyReason, reasonFragment)
		}
		if !columnContainsRequestId(board.Columns.CompletionAnomalies, requestId) {
			t.Fatalf("%s missing from the CompletionAnomalies column", requestId)
		}
		// Never counted as completed "now": no fabricated instant, no
		// Recently-done membership.
		if !ticket.CompletionTime.IsZero() {
			t.Fatalf("%s must carry a zero CompletionTime, got %v", requestId, ticket.CompletionTime)
		}
		if columnContainsRequestId(board.Columns.RecentlyDone, requestId) {
			t.Fatalf("%s (anomalous) must not enter RecentlyDone as if completed now", requestId)
		}
	}
	if got := len(board.Columns.CompletionAnomalies); got != 3 {
		t.Fatalf("CompletionAnomalies count = %d, want 3", got)
	}

	// Every anomaly leaves a warning footprint carrying reason + fix.
	for requestId := range expectedReasonFragments {
		sawWarning := false
		for _, warningText := range board.Warnings {
			if strings.Contains(warningText, requestId) &&
				strings.Contains(warningText, "completion anomaly") &&
				strings.Contains(warningText, "completed_at") {
				sawWarning = true
			}
		}
		if !sawWarning {
			t.Fatalf("expected a completion-anomaly warning naming %s with the fix; got %v", requestId, board.Warnings)
		}
	}
}

// Regression: a terminal ticket whose completed_at parses fine but is older
// than the recent window is neither an anomaly nor recently done — it belongs
// to the calendar only, exactly as before anomaly detection existed.
func TestDatedButOldTicketIsNeitherAnomalyNorRecentlyDone(t *testing.T) {
	board := anomalySyntheticBoard(t)
	datedOldTicket := board.RequestsById["REQ-9304"]
	if datedOldTicket == nil {
		t.Fatalf("REQ-9304 not parsed")
	}
	if datedOldTicket.CompletionAnomaly {
		t.Fatalf("REQ-9304 (parseable completed_at) must not be flagged as an anomaly")
	}
	if columnContainsRequestId(board.Columns.CompletionAnomalies, "REQ-9304") {
		t.Fatalf("REQ-9304 must stay out of the CompletionAnomalies column")
	}
	if columnContainsRequestId(board.Columns.RecentlyDone, "REQ-9304") {
		t.Fatalf("REQ-9304 (completed 2026-01-05, window 7d) must stay out of RecentlyDone")
	}
	if datedOldTicket.CompletionTimeSource != CompletionFromFrontmatter {
		t.Fatalf("REQ-9304 CompletionTimeSource = %q, want frontmatter", datedOldTicket.CompletionTimeSource)
	}
}

// A broken completed_at with a commit hash git CAN resolve keeps the git date
// (so window behavior is unchanged) but is still flagged — the field on disk
// is wrong and the user should fix it.
func TestUnparseableCompletedAtStillFlaggedWhenGitRescuesTheDate(t *testing.T) {
	ticket := &RequestTicket{
		RequestId:            "REQ-9310",
		Status:               "completed",
		CompletedAt:          "yesterday-ish",
		CommitHash:           "deadbeef",
		CommitHashField:      "commit_hash",
		CompletionTime:       time.Date(2026, 6, 29, 8, 0, 0, 0, time.UTC),
		CompletionTimeSource: CompletionFromGitLog,
	}
	flagged, reason := detectCompletionAnomaly(ticket)
	if !flagged {
		t.Fatalf("git-dated ticket with unparseable completed_at should still be flagged")
	}
	if !strings.Contains(reason, `completed_at "yesterday-ish" does not parse`) {
		t.Fatalf("reason = %q, want it to name the bad completed_at", reason)
	}
	if strings.Contains(reason, "deadbeef") {
		t.Fatalf("reason = %q must not blame the commit hash that resolved fine", reason)
	}
}

func TestCompletionAnomaliesInGeneratedPayload(t *testing.T) {
	board := anomalySyntheticBoard(t)
	boardData, buildError := buildGeneratedBoardData(board)
	if buildError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", buildError)
	}

	wantAnomalyIds := []string{"REQ-9301", "REQ-9302", "REQ-9303"}
	if len(boardData.Columns.CompletionAnomalies) != len(wantAnomalyIds) {
		t.Fatalf("payload completionAnomalies = %v, want %v", boardData.Columns.CompletionAnomalies, wantAnomalyIds)
	}
	for _, requestId := range wantAnomalyIds {
		if !stringSliceContains(boardData.Columns.CompletionAnomalies, requestId) {
			t.Fatalf("payload completionAnomalies %v missing %s", boardData.Columns.CompletionAnomalies, requestId)
		}
		generated := boardData.Requests[requestId]
		if !generated.CompletionAnomaly {
			t.Fatalf("payload request %s should carry completionAnomaly", requestId)
		}
		if generated.CompletionAnomalyReason == "" {
			t.Fatalf("payload request %s should carry a completionAnomalyReason", requestId)
		}
		if generated.CompletionTime != "" {
			t.Fatalf("payload request %s completionTime = %q, want empty (no fabricated instant)",
				requestId, generated.CompletionTime)
		}
	}
	if boardData.Requests["REQ-9304"].CompletionAnomaly {
		t.Fatalf("payload request REQ-9304 (dated, old) must not carry completionAnomaly")
	}
	if stringSliceContains(boardData.Columns.RecentlyDone, "REQ-9301") {
		t.Fatalf("payload recentlyDone must not contain the anomalous REQ-9301")
	}
}

func TestCompletionAnomaliesInSummaryOutput(t *testing.T) {
	board := anomalySyntheticBoard(t)
	var summaryBuffer bytes.Buffer
	writeBoardSummary(&summaryBuffer, board)
	summaryText := summaryBuffer.String()

	if !strings.Contains(summaryText, "completion anomalies : 3") {
		t.Fatalf("summary should print the anomaly count; got:\n%s", summaryText)
	}
	for _, requestId := range []string{"REQ-9301", "REQ-9302", "REQ-9303"} {
		if !strings.Contains(summaryText, requestId) {
			t.Fatalf("summary should list anomalous %s; got:\n%s", requestId, summaryText)
		}
	}
	if strings.Contains(summaryText, "REQ-9304 —") {
		t.Fatalf("summary must not list the healthy REQ-9304 as an anomaly; got:\n%s", summaryText)
	}
}

// The clean shared synthetic tree must stay anomaly-free: its terminal tickets
// all resolve a completion instant (frontmatter or the stubbed git lookup).
func TestCleanSyntheticTreeHasNoCompletionAnomalies(t *testing.T) {
	board := syntheticBoard(t)
	if got := len(board.Columns.CompletionAnomalies); got != 0 {
		t.Fatalf("clean tree CompletionAnomalies = %d, want 0", got)
	}
}
