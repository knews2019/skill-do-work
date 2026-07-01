package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// liveBoard builds the board against the REAL do-work tree, resolving the repo
// root by walking up from the test's working directory. The git lookup is
// stubbed to false so the whole-tree build is deterministic and fast: completed
// REQs missing a parseable completed_at resolve as undated rather than spawning
// one git process per file. (lookupGitCommitDate itself is exercised separately.)
func liveBoard(t *testing.T) *Board {
	t.Helper()
	workingDirectory, getwdError := os.Getwd()
	if getwdError != nil {
		t.Fatalf("getwd: %v", getwdError)
	}
	repoRoot, resolveError := resolveRepoRoot(workingDirectory)
	if resolveError != nil {
		t.Fatalf("resolveRepoRoot: %v", resolveError)
	}
	stubGitLookup := func(string, string) (time.Time, bool) { return time.Time{}, false }
	board, buildError := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, stubGitLookup)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	return board
}

func TestLiveTreeExcludesMirrorAndDeliverables(t *testing.T) {
	board := liveBoard(t)
	for _, ticket := range board.AllRequests {
		if strings.Contains(ticket.FilePath, filepath.Join("kb", "wiki", "sources")) {
			t.Fatalf("kb mirror leaked into the board: %s", ticket.FilePath)
		}
		if strings.Contains(ticket.FilePath, filepath.Join("do-work", "deliverables")) {
			t.Fatalf("deliverables leaked into the board: %s", ticket.FilePath)
		}
	}
}

// TestLiveTreeArchiveShapeClassifierInvariant asserts a repo-independent
// invariant: the banded and flat archive-shape classifiers are mutually
// exclusive on every archive REQ path (no path is classified as both). It does
// NOT require either shape to be present — a freshly-extracted repo may have no
// banded archive bands yet. Exact "both shapes parse" coverage lives in
// TestSyntheticParsesBothArchiveShapes, which seeds a deterministic tree.
func TestLiveTreeArchiveShapeClassifierInvariant(t *testing.T) {
	board := liveBoard(t)
	for _, ticket := range board.AllRequests {
		if ticket.TreeSection != "archive" {
			continue
		}
		if pathHasBandedArchiveSegment(ticket.FilePath) && pathIsFlatArchiveRequest(ticket.FilePath) {
			t.Fatalf("archive REQ %s classified as BOTH banded and flat: %s", ticket.RequestId, ticket.FilePath)
		}
	}
}

func TestLiveTreeStatusNormalization(t *testing.T) {
	board := liveBoard(t)
	for _, ticket := range board.AllRequests {
		if isCompletedStatus(ticket.Status) && ticket.Status != "completed" && ticket.Status != "completed-with-issues" {
			t.Fatalf("unexpected completed* status %q on %s", ticket.Status, ticket.RequestId)
		}
		// Forward invariant: any legacy 'complete' present in this repo must
		// normalize to 'completed'. We assert the RULE without requiring such a REQ
		// to exist (TestParseRequestTicketNormalizesAndResolves and
		// TestSyntheticLegacyCompleteNormalized cover the mapping on seeded fixtures).
		if strings.ToLower(strings.TrimSpace(ticket.OriginalStatus)) == "complete" && ticket.Status != "completed" {
			t.Fatalf("legacy 'complete' on %s normalized to %q, want completed", ticket.RequestId, ticket.Status)
		}
	}
}

func TestLiveTreeColumnBucketingMatchesStatus(t *testing.T) {
	board := liveBoard(t)

	// Invariant, not a snapshot: the status-driven columns must contain exactly
	// the tickets whose normalized status maps to them — regardless of which REQs
	// happen to be in the queue today. (This test previously hard-coded
	// REQ-1207..1210 as Pending; that broke the moment those REQs were archived to
	// `completed`. Assert the bucketing RULE — status maps to column — never
	// today's queue contents.)
	idSet := func(column []*RequestTicket) map[string]bool {
		ids := map[string]bool{}
		for _, ticket := range column {
			ids[ticket.RequestId] = true
		}
		return ids
	}
	pendingColumn := idSet(board.Columns.Pending)
	claimedColumn := idSet(board.Columns.Claimed)
	needsInputColumn := idSet(board.Columns.NeedsInputOrBlocked)

	// Forward direction: every parsed ticket lands in the column its status dictates.
	for _, ticket := range board.AllRequests {
		switch {
		case ticket.Status == "pending":
			if !pendingColumn[ticket.RequestId] {
				t.Fatalf("pending %s is missing from the Pending column", ticket.RequestId)
			}
		case ticket.Status == "claimed":
			if !claimedColumn[ticket.RequestId] {
				t.Fatalf("claimed %s is missing from the Claimed column", ticket.RequestId)
			}
		case isNeedsInputOrBlockedStatus(ticket.Status):
			if !needsInputColumn[ticket.RequestId] {
				t.Fatalf("needs-input/blocked %s is missing from the Needs-input column", ticket.RequestId)
			}
		}
	}

	// Reverse direction: nothing foreign sneaks into a status-driven column.
	for _, ticket := range board.Columns.Pending {
		if ticket.Status != "pending" {
			t.Fatalf("non-pending status %q in the Pending column (%s)", ticket.Status, ticket.RequestId)
		}
	}
	for _, ticket := range board.Columns.Claimed {
		if ticket.Status != "claimed" {
			t.Fatalf("non-claimed status %q in the Claimed column (%s)", ticket.Status, ticket.RequestId)
		}
	}
	for _, ticket := range board.Columns.NeedsInputOrBlocked {
		if !isNeedsInputOrBlockedStatus(ticket.Status) {
			t.Fatalf("status %q does not belong in the Needs-input column (%s)", ticket.Status, ticket.RequestId)
		}
	}
	// Recently-done holds only completed* tickets (the within-window subset; the
	// rest of the completed history lives in the calendar).
	for _, ticket := range board.Columns.RecentlyDone {
		if !isCompletedStatus(ticket.Status) {
			t.Fatalf("non-completed status %q in the Recently-done column (%s)", ticket.Status, ticket.RequestId)
		}
	}
}

func TestLiveTreeUserRequestLinkage(t *testing.T) {
	board := liveBoard(t)

	// Every REQ that points at a UR must be grouped under exactly that UR. This is
	// a repo-independent invariant; the seeded UR↔REQ grouping assertions (which
	// previously hard-coded REQ-1207/UR-448 from the source monorepo) now live in
	// TestSyntheticUserRequestLinkage.
	linkageChecks := 0
	for _, ticket := range board.AllRequests {
		if ticket.UserRequestId == "" {
			continue
		}
		parent := board.UserRequestsById[ticket.UserRequestId]
		if parent == nil {
			t.Fatalf("%s points at missing UR %s", ticket.RequestId, ticket.UserRequestId)
		}
		if !stringSliceContains(parent.RequestIds, ticket.RequestId) {
			t.Fatalf("%s not grouped under its UR %s", ticket.RequestId, ticket.UserRequestId)
		}
		linkageChecks++
		if linkageChecks >= 50 {
			break // a representative sample is enough; the loop is O(n) otherwise
		}
	}
}

// TestLiveTreeCompletionTimeConsistent asserts the completion-resolution
// contract: a resolved (non-zero) instant always carries a real source
// (frontmatter or git), and an unresolved completion is exactly the zero time
// with the unresolved source. Mtime is never a source — a clone or tarball
// extraction resets mtimes, so it fabricates completion dates.
func TestLiveTreeCompletionTimeConsistent(t *testing.T) {
	board := liveBoard(t)
	for _, ticket := range board.AllRequests {
		if !isCompletedStatus(ticket.Status) {
			continue
		}
		if ticket.CompletionTime.IsZero() != (ticket.CompletionTimeSource == CompletionUnresolved) {
			t.Fatalf("completed REQ %s: zero-time=%v disagrees with source=%q",
				ticket.RequestId, ticket.CompletionTime.IsZero(), ticket.CompletionTimeSource)
		}
		if !ticket.CompletionTime.IsZero() &&
			ticket.CompletionTimeSource != CompletionFromFrontmatter && ticket.CompletionTimeSource != CompletionFromGitLog {
			t.Fatalf("completed REQ %s resolved via unexpected source %q", ticket.RequestId, ticket.CompletionTimeSource)
		}
	}
}

// TestLiveTreeCalendarCoversCompletions asserts a repo-independent invariant: the
// completion calendar holds exactly one entry per completed* REQ (each of which
// resolves a completion time — see TestLiveTreeCompletionTimeResolved). The old
// absolute ">= 900 tickets" ballpark was a source-monorepo snapshot that breaks
// in this 33-REQ extraction; exact counts now live in TestSyntheticCountsAndCalendar.
func TestLiveTreeCalendarCoversCompletions(t *testing.T) {
	board := liveBoard(t)
	if len(board.AllRequests) == 0 {
		t.Fatalf("live tree parsed zero REQ tickets — the do-work walk found nothing")
	}

	completedCount := 0
	for _, ticket := range board.AllRequests {
		if isCompletedStatus(ticket.Status) {
			completedCount++
		}
	}
	if len(board.Calendar) != completedCount {
		t.Fatalf("calendar entries = %d, want %d (one per completed REQ)", len(board.Calendar), completedCount)
	}
}

// TestLiveGitCommitDateLookupBestEffort exercises the real git fallback step
// against the repo's own HEAD commit. It is skipped (not failed) when git is
// unavailable, honoring the best-effort contract.
func TestLiveGitCommitDateLookupBestEffort(t *testing.T) {
	workingDirectory, getwdError := os.Getwd()
	if getwdError != nil {
		t.Fatalf("getwd: %v", getwdError)
	}
	repoRoot, resolveError := resolveRepoRoot(workingDirectory)
	if resolveError != nil {
		t.Fatalf("resolveRepoRoot: %v", resolveError)
	}

	headBytes, revParseError := exec.Command("git", "-C", repoRoot, "rev-parse", "HEAD").Output()
	if revParseError != nil {
		t.Skip("git unavailable or not a repo; skipping live git lookup")
	}
	headHash := strings.TrimSpace(string(headBytes))

	committedAt, ok := lookupGitCommitDate(repoRoot, headHash)
	if !ok {
		t.Fatalf("lookupGitCommitDate(HEAD=%s) returned not-ok", headHash)
	}
	if committedAt.IsZero() {
		t.Fatalf("lookupGitCommitDate(HEAD) returned the zero time")
	}

	if _, missingOk := lookupGitCommitDate(repoRoot, "0000000000000000000000000000000000000000"); missingOk {
		t.Fatalf("lookupGitCommitDate of an unknown hash should be not-ok")
	}
}

func stringSliceContains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

// pathHasBandedArchiveSegment reports whether any path segment is a banded UR
// folder of the form "UR-NNN-MMM" (the archive band shape).
func pathHasBandedArchiveSegment(path string) bool {
	for _, segment := range strings.Split(path, string(filepath.Separator)) {
		if isBandedUrFolderName(segment) {
			return true
		}
	}
	return false
}

// pathIsFlatArchiveRequest reports whether a REQ path lives directly inside a
// flat "archive/UR-NNN/" folder (parent is a single-number UR folder whose own
// parent is "archive") — the non-banded archive shape.
func pathIsFlatArchiveRequest(path string) bool {
	parentDirectory := filepath.Dir(path)
	parentName := filepath.Base(parentDirectory)
	grandparentName := filepath.Base(filepath.Dir(parentDirectory))
	return grandparentName == "archive" && isPlainUrFolderName(parentName)
}

func isBandedUrFolderName(segment string) bool {
	if !strings.HasPrefix(segment, "UR-") {
		return false
	}
	parts := strings.Split(strings.TrimPrefix(segment, "UR-"), "-")
	if len(parts) != 2 {
		return false
	}
	return isAllDigits(parts[0]) && isAllDigits(parts[1])
}

func isPlainUrFolderName(segment string) bool {
	if !strings.HasPrefix(segment, "UR-") {
		return false
	}
	return isAllDigits(strings.TrimPrefix(segment, "UR-"))
}

func isAllDigits(text string) bool {
	if text == "" {
		return false
	}
	for _, character := range text {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}
