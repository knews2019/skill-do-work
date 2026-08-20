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
	// The calendar holds every REQ, not just the finished ones. Asserting the
	// exact sequence pins the render-order contract board-calendar.js depends on:
	// it groups days by walking contiguous DayKeys, so a band emitted out of
	// order renders as two separate sections of the same band.
	if got := len(board.Calendar); got != len(board.AllRequests) {
		t.Fatalf("calendar entries = %d, want %d (one per REQ — queued and claimed work belong on the calendar too)",
			got, len(board.AllRequests))
	}
	wantCalendar := []struct {
		requestId string
		dayKey    string
	}{
		{"REQ-9001", queuedCalendarDayKey},  // pending — never started
		{"REQ-9006", queuedCalendarDayKey},  // blocked — also never started
		{"REQ-9005", "2026-06-28"},          // cancelled, dated by completed_at
		{"REQ-9004", "2026-06-10"},          // completed, dated by completed_at
		{"REQ-9003", "2026-03-04"},          // completed, dated by the deadbeef commit
		{"REQ-9002", undatedCalendarDayKey}, // claimed with no claimed_at to place it on
	}
	for index, want := range wantCalendar {
		got := board.Calendar[index]
		if got.RequestId != want.requestId || got.DayKey != want.dayKey {
			t.Fatalf("calendar[%d] = %s on %q, want %s on %q (order is queued band, then days newest-first, then undated)",
				index, got.RequestId, got.DayKey, want.requestId, want.dayKey)
		}
	}
}

// TestCalendarDatesClaimsAndFailuresFromTheirOwnStamps pins the two placements
// that do NOT come from the shared completion-resolution chain.
//
// A claimed REQ is dated by claimed_at — the day work STARTED — because it has
// no completion instant to be dated by; getting this wrong would either hide
// in-flight work or date it by a stale leftover field. A failed REQ is dated by
// completed_at read here rather than by ticket.CompletionTime, because the
// resolution step deliberately skips `failed` so the detail drawer never shows a
// "Completed" row for work that did not complete; before this it appeared on no
// day at all. Either stamp missing means the entry falls to the undated bucket,
// never to a fabricated date and never out of the view.
func TestCalendarDatesClaimsAndFailuresFromTheirOwnStamps(t *testing.T) {
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
	request := func(requestId string, status string, extraFrontmatter string) string {
		return "---\nid: " + requestId + "\ntitle: Fixture " + requestId +
			"\nstatus: " + status + "\n" + extraFrontmatter + "---\n\nBody.\n"
	}

	// The claimed REQ carries a completed_at as well, left over from an earlier
	// run: the claim must still be dated by claimed_at, not by that stamp.
	writeFixture(filepath.Join("do-work", "working", "REQ-9101-claimed.md"),
		request("REQ-9101", "claimed", "claimed_at: 2026-06-29T09:00:00Z\ncompleted_at: 2026-06-20T09:00:00Z\n"))
	writeFixture(filepath.Join("do-work", "archive", "REQ-9102-failed.md"),
		request("REQ-9102", "failed", "completed_at: 2026-06-27T10:00:00Z\nerror: exploded\n"))
	writeFixture(filepath.Join("do-work", "archive", "REQ-9103-failed-unstamped.md"),
		request("REQ-9103", "failed", "error: exploded with no stamp\n"))

	noGitLookup := func(_ string, _ string) (time.Time, bool) { return time.Time{}, false }
	board, buildError := buildBoard(repoRoot, time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC), 7*24*time.Hour, noGitLookup)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}

	byRequestId := map[string]CalendarEntry{}
	for _, entry := range board.Calendar {
		byRequestId[entry.RequestId] = entry
	}

	claimed := byRequestId["REQ-9101"]
	if claimed.Kind != CalendarClaimEntry || claimed.DayKey != "2026-06-29" {
		t.Fatalf("claimed REQ-9101 = %s on %q, want a claim entry on 2026-06-29 (its claimed_at day, not its stale completed_at)",
			claimed.Kind, claimed.DayKey)
	}

	failed := byRequestId["REQ-9102"]
	if failed.DayKey != "2026-06-27" || failed.TimeSource != CompletionFromFrontmatter {
		t.Fatalf("failed REQ-9102 = %q via %q, want 2026-06-27 via frontmatter — a failure happened on a real day",
			failed.DayKey, failed.TimeSource)
	}

	unstamped := byRequestId["REQ-9103"]
	if unstamped.DayKey != undatedCalendarDayKey {
		t.Fatalf("failed-without-completed_at REQ-9103 = %q, want %q (kept and visible, never given a fabricated date)",
			unstamped.DayKey, undatedCalendarDayKey)
	}

	// The drawer must not gain a "Completed" row for failed work.
	for _, ticket := range board.AllRequests {
		if ticket.Status == "failed" && !ticket.CompletionTime.IsZero() {
			t.Fatalf("%s is failed but carries CompletionTime %s — that field feeds the drawer's Completed row",
				ticket.RequestId, ticket.CompletionTime)
		}
	}
}

// TestCalendarDayKeySentinelsMatchTheShippedClient keeps the two sentinel day
// keys in lock-step across the language boundary. board-calendar.js matches the
// queued key as a literal to pick the "In the queue" label, and reaches the
// undated label by the same key failing to parse as a date. Renaming either
// constant on the Go side alone compiles, passes every Go test, and silently
// relabels a whole band in the browser — the queued band would render under
// "Undated", which is the opposite of what it means.
func TestCalendarDayKeySentinelsMatchTheShippedClient(t *testing.T) {
	calendarSource, readError := embeddedWebAssets.ReadFile("web/board-calendar.js")
	if readError != nil {
		t.Fatalf("read web/board-calendar.js: %v", readError)
	}
	if !strings.Contains(string(calendarSource), `dayKey === "`+queuedCalendarDayKey+`"`) {
		t.Errorf("web/board-calendar.js does not test for the queued day key %q, so the queued band will render with the undated label",
			queuedCalendarDayKey)
	}
	// The undated branch is the client's fallback, so what must hold is that the
	// key cannot parse as a date — otherwise it would be grouped as a real day.
	if _, parseError := time.Parse("2006-01-02", undatedCalendarDayKey); parseError == nil {
		t.Errorf("undated day key %q parses as a date, so board-calendar.js would render it as a day rather than the Undated band",
			undatedCalendarDayKey)
	}
}
