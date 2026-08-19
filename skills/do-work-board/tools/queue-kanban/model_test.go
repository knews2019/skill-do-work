package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNormalizeStatus(t *testing.T) {
	testCases := []struct {
		raw  string
		want string
	}{
		{"completed", "completed"},
		{"complete", "completed"},
		{"done", "completed"},
		{"finished", "completed"},
		{"closed", "completed"},
		{"  Complete ", "completed"},
		{"completed-with-issues", "completed-with-issues"},
		{"cancelled", "cancelled"},
		{"canceled", "cancelled"},
		{"abandoned", "cancelled"},
		{"wont-do", "cancelled"},
		{"wontfix", "cancelled"},
		{"  Cancelled ", "cancelled"},
		{"pending", "pending"},
		{"pending-answers", "pending-answers"},
		{"blocked", "blocked"},
		{"claimed", "claimed"},
		{"custom-status", "custom-status"},
	}
	for _, testCase := range testCases {
		if got := normalizeStatus(testCase.raw); got != testCase.want {
			t.Fatalf("normalizeStatus(%q) = %q, want %q", testCase.raw, got, testCase.want)
		}
	}
}

func TestStatusClassifiers(t *testing.T) {
	if !isCompletedStatus("completed") || !isCompletedStatus("completed-with-issues") {
		t.Fatalf("completed* statuses should classify as completed")
	}
	if isCompletedStatus("pending") {
		t.Fatalf("pending must not classify as completed")
	}
	for _, typoStatus := range []string{"completed-wth-issues", "completedish", "completed-with-issues-again"} {
		if isCompletedStatus(typoStatus) {
			t.Fatalf("%q is outside the terminal-success enum (completed | completed-with-issues) — a completed-prefixed typo must route through the unrecognized-status warning path, not classify as done", typoStatus)
		}
	}
	for _, blocked := range []string{"pending-answers", "blocked", "blocked-archive-collision", "blocked-dependency-cycle", "failed"} {
		if !isNeedsInputOrBlockedStatus(blocked) {
			t.Fatalf("%q should be a needs-input/blocked status", blocked)
		}
	}
	if isNeedsInputOrBlockedStatus("pending") || isNeedsInputOrBlockedStatus("claimed") {
		t.Fatalf("pending/claimed are their own columns, not needs-input/blocked")
	}
	if isNeedsInputOrBlockedStatus("deferred") {
		t.Fatalf("deferred is not in the Schema Read Contract enum (actions/work-reference.md) — it must route through the unrecognized-status warning path, not the recognized list")
	}
	// Bare `blocked` is a recognized needs-input status, but a `blocked-<reason>`
	// variant that is NOT one of the two canonical holding states must still fall
	// through to the unrecognized-status warning path.
	if isNeedsInputOrBlockedStatus("blocked-on-lm-studio") {
		t.Fatalf("blocked-on-lm-studio is not in the Schema Read Contract enum — only bare `blocked` plus the two blocked-* holding states are recognized")
	}
	if isCompletedStatus("blocked") || isTerminalResolvedStatus("blocked") {
		t.Fatalf("blocked is a non-terminal holding status — it must never classify as completed or terminally resolved")
	}
	if !isCancelledStatus("cancelled") {
		t.Fatalf("cancelled should classify as the cancelled terminal status")
	}
	for _, notCancelled := range []string{"cancel", "cancelledish", "cancelled-maybe"} {
		if isCancelledStatus(notCancelled) {
			t.Fatalf("%q is outside the canonical enum — a cancelled-prefixed typo must route through the unrecognized-status warning path", notCancelled)
		}
	}
	if isCompletedStatus("cancelled") {
		t.Fatalf("cancelled must NOT classify as terminal success — success-readers exclude it (Terminal-success status set, actions/work-reference.md)")
	}
	if isNeedsInputOrBlockedStatus("cancelled") {
		t.Fatalf("cancelled is terminally resolved — it belongs with done work, not the Needs-input/Blocked column")
	}
	for _, resolved := range []string{"completed", "completed-with-issues", "cancelled"} {
		if !isTerminalResolvedStatus(resolved) {
			t.Fatalf("%q should classify as terminally resolved (Terminal-resolved status set, actions/work-reference.md)", resolved)
		}
	}
	for _, notResolved := range []string{"failed", "pending", "claimed", "pending-answers"} {
		if isTerminalResolvedStatus(notResolved) {
			t.Fatalf("%q must not classify as terminally resolved", notResolved)
		}
	}
}

func TestResolveCommitHashVariants(t *testing.T) {
	testCases := []struct {
		name      string
		fields    map[string]any
		want      string
		wantField string
	}{
		{"canonical commit", map[string]any{"commit": "abc123"}, "abc123", "commit"},
		{"commit_hash variant", map[string]any{"commit_hash": "def456"}, "def456", "commit_hash"},
		{"green_commit variant", map[string]any{"green_commit": "aaa111"}, "aaa111", "green_commit"},
		{"commit_green variant", map[string]any{"commit_green": "bbb222"}, "bbb222", "commit_green"},
		{"impl_commit variant", map[string]any{"impl_commit": "ccc333"}, "ccc333", "impl_commit"},
		{"canonical wins over variant", map[string]any{"commit": "primary", "commit_hash": "secondary"}, "primary", "commit"},
		{"none present", map[string]any{}, "", ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, gotField := resolveCommitHash(testCase.fields)
			if got != testCase.want || gotField != testCase.wantField {
				t.Fatalf("resolveCommitHash = (%q, %q), want (%q, %q)", got, gotField, testCase.want, testCase.wantField)
			}
		})
	}
}

func TestResolveDependsOnPrefersCanonical(t *testing.T) {
	canonical := resolveDependsOn(map[string]any{
		"depends_on":   []any{"REQ-10"},
		"dependencies": []any{"REQ-99"},
	})
	if !reflect.DeepEqual(canonical, []string{"REQ-10"}) {
		t.Fatalf("depends_on should win, got %v", canonical)
	}
	legacy := resolveDependsOn(map[string]any{"dependencies": []any{"REQ-99"}})
	if !reflect.DeepEqual(legacy, []string{"REQ-99"}) {
		t.Fatalf("legacy dependencies should be used when depends_on absent, got %v", legacy)
	}
}

func TestDeriveRequestIdFromFilename(t *testing.T) {
	testCases := map[string]string{
		"/x/do-work/queue/REQ-1207-queue-kanban-parser.md": "REQ-1207",
		"/x/archive/UR-446/REQ-1203-modal-shell.md":        "REQ-1203",
	}
	for path, want := range testCases {
		if got := deriveRequestIdFromFilename(path); got != want {
			t.Fatalf("deriveRequestIdFromFilename(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestIdentifierLessNumericOrder(t *testing.T) {
	ids := []string{"REQ-100", "REQ-9", "REQ-21"}
	sortRequestIdList(ids)
	want := []string{"REQ-9", "REQ-21", "REQ-100"}
	if !reflect.DeepEqual(ids, want) {
		t.Fatalf("numeric id order = %v, want %v", ids, want)
	}
}

// TestResolveCompletionTimeFallbackChain exercises every step of the fallback
// chain (frontmatter → git → unresolved) deterministically, with the git
// lookup injected so no subprocess is spawned. File mtime is deliberately NOT
// in the chain — a clone/checkout/extraction resets it, fabricating dates.
func TestResolveCompletionTimeFallbackChain(t *testing.T) {
	temporaryDirectory := t.TempDir()
	knownModificationTime := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)

	existingFile := filepath.Join(temporaryDirectory, "REQ-1-existing.md")
	if writeError := os.WriteFile(existingFile, []byte("body"), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}
	if chtimesError := os.Chtimes(existingFile, knownModificationTime, knownModificationTime); chtimesError != nil {
		t.Fatalf("chtimes fixture: %v", chtimesError)
	}

	gitTime := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)
	stubGitLookup := func(repoRoot string, commitHash string) (time.Time, bool) {
		if commitHash == "deadbeef" {
			return gitTime, true
		}
		return time.Time{}, false
	}

	t.Run("frontmatter completed_at wins", func(t *testing.T) {
		ticket := &RequestTicket{CompletedAt: "2026-06-10T14:00:00Z", CommitHash: "deadbeef", FilePath: existingFile}
		got, source := resolveCompletionTime(ticket, temporaryDirectory, stubGitLookup)
		if source != CompletionFromFrontmatter {
			t.Fatalf("source = %q, want frontmatter", source)
		}
		want, _ := parseTimestamp("2026-06-10T14:00:00Z")
		if !got.Equal(want) {
			t.Fatalf("time = %v, want %v", got, want)
		}
	})

	t.Run("git committer date is the second step", func(t *testing.T) {
		ticket := &RequestTicket{CommitHash: "deadbeef", FilePath: existingFile}
		got, source := resolveCompletionTime(ticket, temporaryDirectory, stubGitLookup)
		if source != CompletionFromGitLog {
			t.Fatalf("source = %q, want git", source)
		}
		if !got.Equal(gitTime) {
			t.Fatalf("time = %v, want %v", got, gitTime)
		}
	})

	t.Run("file mtime is NOT a fallback", func(t *testing.T) {
		// The file exists with a known old mtime, but no frontmatter timestamp and
		// no resolvable commit — the completion must stay unresolved instead of
		// adopting the mtime (which a clone/checkout/extraction would have reset).
		ticket := &RequestTicket{FilePath: existingFile}
		got, source := resolveCompletionTime(ticket, temporaryDirectory, stubGitLookup)
		if source != CompletionUnresolved {
			t.Fatalf("source = %q, want unresolved (mtime must not be used)", source)
		}
		if !got.IsZero() {
			t.Fatalf("time = %v, want zero (mtime must not be used)", got)
		}
	})

	t.Run("unresolved when nothing is available", func(t *testing.T) {
		ticket := &RequestTicket{FilePath: filepath.Join(temporaryDirectory, "does-not-exist.md")}
		got, source := resolveCompletionTime(ticket, temporaryDirectory, stubGitLookup)
		if source != CompletionUnresolved {
			t.Fatalf("source = %q, want unresolved", source)
		}
		if !got.IsZero() {
			t.Fatalf("time = %v, want zero", got)
		}
	})
}

// TestDedupeTicketsByRequestId covers the queue+archive id-collision state the
// skill explicitly models (blocked-archive-collision): exactly one copy per id
// may reach the views (the id-keyed JSON map can only carry one), the active
// copy wins, and the duplicate is surfaced as a warning — never dropped silently.
func TestDedupeTicketsByRequestId(t *testing.T) {
	archiveCopy := &RequestTicket{RequestId: "REQ-42", Status: "completed", TreeSection: "archive", FilePath: "/a/REQ-42.md"}
	queueCopy := &RequestTicket{RequestId: "REQ-42", Status: "pending", TreeSection: "queue", FilePath: "/q/REQ-42.md"}
	unrelated := &RequestTicket{RequestId: "REQ-7", Status: "pending", TreeSection: "queue", FilePath: "/q/REQ-7.md"}

	// Archive walks first in the real tree order — the later queue copy must still win.
	winners, warnings := dedupeTicketsByRequestId([]*RequestTicket{archiveCopy, queueCopy, unrelated})
	if len(winners) != 2 {
		t.Fatalf("winners = %d tickets, want 2", len(winners))
	}
	if winners[0] != queueCopy {
		t.Fatalf("winner for REQ-42 = %s copy, want the queue copy", winners[0].TreeSection)
	}
	if len(warnings) != 1 ||
		!strings.Contains(warnings[0], "REQ-42") ||
		!strings.Contains(warnings[0], "/q/REQ-42.md") ||
		!strings.Contains(warnings[0], "/a/REQ-42.md") {
		t.Fatalf("expected one duplicate warning naming both copies, got %v", warnings)
	}
}

func TestIsPlausibleCommitHashRejectsOptionShapedValues(t *testing.T) {
	for _, valid := range []string{"deadbeef", "096dacba", "0123456789abcdefABCDEF00"} {
		if !isPlausibleCommitHash(valid) {
			t.Fatalf("isPlausibleCommitHash(%q) = false, want true", valid)
		}
	}
	for _, invalid := range []string{"", "abc", "--all", "--output=/tmp/pwned", "HEAD", "main", "dead beef", strings.Repeat("a", 65)} {
		if isPlausibleCommitHash(invalid) {
			t.Fatalf("isPlausibleCommitHash(%q) = true, want false", invalid)
		}
	}
}

func TestParseRequestTicketNormalizesAndResolves(t *testing.T) {
	temporaryDirectory := t.TempDir()
	fixturePath := filepath.Join(temporaryDirectory, "REQ-555-legacy-complete.md")
	fixtureContent := `---
id: REQ-555
title: Legacy complete with commit variant
status: complete
commit_hash: feedface
user_request: UR-77
domain: frontend
depends_on: [REQ-500]
dependencies: [REQ-499]
related: [REQ-501]
write_set: [src/a.ts, src/b.ts]
---

# Body heading

Some body text.
`
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}

	ticket, parseError := parseRequestTicket(fixturePath, "archive")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if ticket.RequestId != "REQ-555" {
		t.Fatalf("RequestId = %q", ticket.RequestId)
	}
	if ticket.OriginalStatus != "complete" || ticket.Status != "completed" {
		t.Fatalf("status normalization wrong: original=%q normalized=%q", ticket.OriginalStatus, ticket.Status)
	}
	if ticket.CommitHash != "feedface" {
		t.Fatalf("CommitHash = %q, want feedface", ticket.CommitHash)
	}
	if ticket.UserRequestId != "UR-77" {
		t.Fatalf("UserRequestId = %q", ticket.UserRequestId)
	}
	if !reflect.DeepEqual(ticket.DependsOn, []string{"REQ-500"}) {
		t.Fatalf("DependsOn = %v, want [REQ-500] (depends_on wins)", ticket.DependsOn)
	}
	if !reflect.DeepEqual(ticket.WriteSet, []string{"src/a.ts", "src/b.ts"}) {
		t.Fatalf("WriteSet = %v, want [src/a.ts src/b.ts] (read verbatim, no normalization)", ticket.WriteSet)
	}
	if ticket.TreeSection != "archive" {
		t.Fatalf("TreeSection = %q", ticket.TreeSection)
	}
	if !strings.Contains(ticket.BodyMarkdown, "# Body heading") {
		t.Fatalf("body not preserved: %q", ticket.BodyMarkdown)
	}
}

// TestParseRequestTicketCrlfFileRoundTripsByteExactly pins the Copy contract
// on a Windows-authored (CRLF) REQ file: FrontmatterMarkdown + BodyMarkdown
// must equal the file bytes exactly. The old fence-by-subtraction arithmetic
// measured the fence against a CRLF-normalized body, so the fence stole bytes
// from the body's start and the body itself came back with Unix endings —
// the drawer's Copy payload no longer matched the file on disk.
func TestParseRequestTicketCrlfFileRoundTripsByteExactly(t *testing.T) {
	temporaryDirectory := t.TempDir()
	fixturePath := filepath.Join(temporaryDirectory, "REQ-556-crlf.md")
	fixtureContent := "---\r\nid: REQ-556\r\ntitle: CRLF fixture\r\nstatus: pending\r\n---\r\n\r\n# CRLF body\r\n\r\nAuthored on Windows.\r\n"
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}

	ticket, parseError := parseRequestTicket(fixturePath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if ticket.RequestId != "REQ-556" || ticket.Status != "pending" {
		t.Fatalf("frontmatter fields lost on CRLF file: id=%q status=%q", ticket.RequestId, ticket.Status)
	}
	if got := ticket.FrontmatterMarkdown + ticket.BodyMarkdown; got != fixtureContent {
		t.Fatalf("Copy payload must round-trip the CRLF file byte-for-byte:\ngot  %q\nwant %q", got, fixtureContent)
	}
	if !strings.HasSuffix(ticket.FrontmatterMarkdown, "---\r\n") {
		t.Fatalf("fence must end at the closing --- line, got %q", ticket.FrontmatterMarkdown)
	}
}

func TestBucketColumns(t *testing.T) {
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	window := 48 * time.Hour
	recentDone := &RequestTicket{RequestId: "REQ-1", Status: "completed", CompletionTime: now.Add(-1 * time.Hour)}
	oldDone := &RequestTicket{RequestId: "REQ-2", Status: "completed", CompletionTime: now.Add(-200 * time.Hour)}
	recentCancelled := &RequestTicket{RequestId: "REQ-9", Status: "cancelled", CompletionTime: now.Add(-2 * time.Hour)}
	oldCancelled := &RequestTicket{RequestId: "REQ-10", Status: "cancelled", CompletionTime: now.Add(-300 * time.Hour)}
	tickets := []*RequestTicket{
		{RequestId: "REQ-3", Status: "pending"},
		{RequestId: "REQ-4", Status: "claimed"},
		{RequestId: "REQ-11", Status: "reserved", OriginalStatus: "reserved"}, // legacy status from the removed reserve action — unrecognized now, must land visible with a warning
		{RequestId: "REQ-5", Status: "pending-answers"},
		{RequestId: "REQ-13", Status: "blocked", BlockedBy: []string{"LM Studio running locally"}, BlockedAt: "2026-06-27T12:00:00Z"},         // recognized holding status — Needs-input/Blocked, NO warning
		{RequestId: "REQ-6", Status: "deferred", OriginalStatus: "deferred"},                                                                  // hand-edited status outside the Schema Read Contract enum — must still land in Needs-input/Blocked, now via the unrecognized-status warning path
		{RequestId: "REQ-7", Status: "pnding", OriginalStatus: "pnding"},                                                                      // typo'd status — must never be silently dropped
		{RequestId: "REQ-8", Status: "completed-wth-issues", OriginalStatus: "completed-wth-issues", CompletionTime: now.Add(-1 * time.Hour)}, // completed-prefixed typo — must warn, never pass as terminal success
		recentDone,
		oldDone,
		recentCancelled,
		oldCancelled,
	}
	columns, statusWarnings := bucketColumns(tickets, now, window)
	if len(columns.Pending) != 1 || columns.Pending[0].RequestId != "REQ-3" {
		t.Fatalf("Pending = %+v", columns.Pending)
	}
	if len(columns.Claimed) != 1 || columns.Claimed[0].RequestId != "REQ-4" {
		t.Fatalf("Claimed should hold only the claimed REQ, got %+v", columns.Claimed)
	}
	if len(columns.NeedsInputOrBlocked) != 6 {
		t.Fatalf("NeedsInputOrBlocked should hold pending-answers + blocked + the reserved, deferred, pnding, and completed-wth-issues unrecognized statuses, got %d", len(columns.NeedsInputOrBlocked))
	}
	if len(columns.RecentlyDone) != 2 ||
		columns.RecentlyDone[0].RequestId != "REQ-1" || columns.RecentlyDone[1].RequestId != "REQ-9" {
		t.Fatalf("RecentlyDone should hold the in-window completion then the in-window cancellation (most-recent first), got %+v", columns.RecentlyDone)
	}
	if len(statusWarnings) != 4 {
		t.Fatalf("expected four warnings (reserved + deferred + pnding + completed-wth-issues unrecognized statuses), got %d: %v", len(statusWarnings), statusWarnings)
	}
	foundDeferredWarning := false
	foundTypoWarning := false
	foundCompletedTypoWarning := false
	foundLegacyReservedWarning := false
	for _, warning := range statusWarnings {
		if strings.Contains(warning, "REQ-6") && strings.Contains(warning, "deferred") {
			foundDeferredWarning = true
		}
		if strings.Contains(warning, "REQ-7") && strings.Contains(warning, "pnding") {
			foundTypoWarning = true
		}
		if strings.Contains(warning, "REQ-8") && strings.Contains(warning, "completed-wth-issues") {
			foundCompletedTypoWarning = true
		}
		if strings.Contains(warning, "REQ-11") && strings.Contains(warning, "reserved") {
			foundLegacyReservedWarning = true
		}
	}
	if !foundDeferredWarning {
		t.Fatalf("expected an unrecognized-status warning naming REQ-6/deferred, got %v", statusWarnings)
	}
	if !foundTypoWarning {
		t.Fatalf("expected an unrecognized-status warning naming REQ-7/pnding, got %v", statusWarnings)
	}
	if !foundCompletedTypoWarning {
		t.Fatalf("expected an unrecognized-status warning naming REQ-8/completed-wth-issues, got %v", statusWarnings)
	}
	if !foundLegacyReservedWarning {
		t.Fatalf("expected an unrecognized-status warning naming REQ-11/reserved (legacy status from the removed reserve action), got %v", statusWarnings)
	}
	for _, warning := range statusWarnings {
		if strings.Contains(warning, "REQ-13") {
			t.Fatalf("blocked REQ-13 is a recognized status — it must NOT produce an unrecognized-status warning, got %v", statusWarnings)
		}
	}
}

func TestParseRequestTicketReadsBlockedFields(t *testing.T) {
	temporaryDirectory := t.TempDir()
	fixturePath := filepath.Join(temporaryDirectory, "REQ-778-blocked-on-lm-studio.md")
	fixtureContent := `---
id: REQ-778
title: Waiting on LM Studio
status: blocked
blocked_by: "LM Studio running locally"
blocked_at: 2026-07-18T10:00:00Z
blocked_check: "curl -sf http://localhost:1234/v1/models"
---

Body.
`
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}
	ticket, parseError := parseRequestTicket(fixturePath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if ticket.Status != "blocked" {
		t.Fatalf("Status = %q, want blocked (a recognized Schema Read Contract status, never unrecognized)", ticket.Status)
	}
	// Free-text blocked_by parses as a one-element list through coerceToStringList.
	if len(ticket.BlockedBy) != 1 || ticket.BlockedBy[0] != "LM Studio running locally" {
		t.Fatalf("BlockedBy = %+v, want [\"LM Studio running locally\"]", ticket.BlockedBy)
	}
	if ticket.BlockedAt == "" {
		t.Fatalf("BlockedAt not parsed from frontmatter")
	}
	if ticket.BlockedCheck != "curl -sf http://localhost:1234/v1/models" {
		t.Fatalf("BlockedCheck = %q, want the probe command verbatim", ticket.BlockedCheck)
	}
}

// writeSetOverlapFixture is one REQ to seed for the overlap annotation: an id, a
// normalized status, and its declared write_set.
type writeSetOverlapFixture struct {
	RequestId   string
	Status      string
	WriteSetted []string
}

// annotateWriteSetOverlapFixtures runs the annotation over the fixtures and
// returns id → annotated overlaps, so assertions read as "who does REQ-N contend
// with?" rather than as index arithmetic.
func annotateWriteSetOverlapFixtures(fixtures []writeSetOverlapFixture) map[string][]string {
	tickets := make([]*RequestTicket, 0, len(fixtures))
	for _, fixture := range fixtures {
		tickets = append(tickets, &RequestTicket{
			RequestId: fixture.RequestId,
			Status:    fixture.Status,
			WriteSet:  fixture.WriteSetted,
		})
	}
	annotateWriteSetOverlap(tickets)

	overlapsById := map[string][]string{}
	for _, ticket := range tickets {
		overlapsById[ticket.RequestId] = ticket.WriteSetOverlaps
	}
	return overlapsById
}

func TestAnnotateWriteSetOverlapPairsContendingRequests(t *testing.T) {
	overlapsById := annotateWriteSetOverlapFixtures([]writeSetOverlapFixture{
		{RequestId: "REQ-1", Status: "pending", WriteSetted: []string{"web/board.css", "web/board.js"}},
		{RequestId: "REQ-2", Status: "claimed", WriteSetted: []string{"web/board.css"}},         // literal overlap with REQ-1
		{RequestId: "REQ-3", Status: "pending", WriteSetted: []string{"docs/board-guide.md"}},   // disjoint from everything
		{RequestId: "REQ-4", Status: "pending", WriteSetted: []string{"web/*.js"}},              // glob catching REQ-1's literal
		{RequestId: "REQ-5", Status: "pending", WriteSetted: []string{}},                        // declared nothing — unknown, never a badge
		{RequestId: "REQ-6", Status: "completed", WriteSetted: []string{"web/board.css"}},       // terminal — not a dispatch candidate
		{RequestId: "REQ-7", Status: "reserved", WriteSetted: []string{"web/board.css"}},        // legacy status outside the vocabulary — not a candidate tier, excluded
		{RequestId: "REQ-8", Status: "pending-answers", WriteSetted: []string{"web/board.css"}}, // needs-input tier — excluded
	})

	expectedOverlapsById := map[string][]string{
		"REQ-1": {"REQ-2", "REQ-4"},
		"REQ-2": {"REQ-1"},
		"REQ-3": nil,
		"REQ-4": {"REQ-1"},
		"REQ-5": nil,
		"REQ-6": nil,
		"REQ-7": nil,
		"REQ-8": nil,
	}
	for requestId, expectedOverlaps := range expectedOverlapsById {
		if !reflect.DeepEqual(overlapsById[requestId], expectedOverlaps) {
			t.Errorf("%s overlaps = %v, want %v", requestId, overlapsById[requestId], expectedOverlaps)
		}
	}
	for requestId, actualOverlaps := range overlapsById {
		for _, overlappingId := range actualOverlaps {
			if overlappingId == requestId {
				t.Errorf("%s lists itself as an overlap: %v", requestId, actualOverlaps)
			}
		}
	}
}

// A glob must catch a literal no matter which side declared it — the pair is
// compared in both directions, so badge presence can't depend on REQ ordering.
func TestAnnotateWriteSetOverlapMatchesGlobsInBothDirections(t *testing.T) {
	globFirst := annotateWriteSetOverlapFixtures([]writeSetOverlapFixture{
		{RequestId: "REQ-1", Status: "pending", WriteSetted: []string{"src/auth/*.ts"}},
		{RequestId: "REQ-2", Status: "pending", WriteSetted: []string{"src/auth/session.ts"}},
	})
	literalFirst := annotateWriteSetOverlapFixtures([]writeSetOverlapFixture{
		{RequestId: "REQ-1", Status: "pending", WriteSetted: []string{"src/auth/session.ts"}},
		{RequestId: "REQ-2", Status: "pending", WriteSetted: []string{"src/auth/*.ts"}},
	})
	for _, overlapsById := range []map[string][]string{globFirst, literalFirst} {
		if !reflect.DeepEqual(overlapsById["REQ-1"], []string{"REQ-2"}) ||
			!reflect.DeepEqual(overlapsById["REQ-2"], []string{"REQ-1"}) {
			t.Fatalf("glob-vs-literal must pair both ways, got %v", overlapsById)
		}
	}

	// A glob that does NOT cover the literal stays disjoint (the match is real
	// pattern matching, not a substring or prefix guess).
	nonMatching := annotateWriteSetOverlapFixtures([]writeSetOverlapFixture{
		{RequestId: "REQ-1", Status: "pending", WriteSetted: []string{"src/auth/*.ts"}},
		{RequestId: "REQ-2", Status: "pending", WriteSetted: []string{"src/billing/invoice.ts"}},
	})
	if nonMatching["REQ-1"] != nil || nonMatching["REQ-2"] != nil {
		t.Fatalf("a glob matching nothing in the other set must not annotate, got %v", nonMatching)
	}
}

// path.Match's `*` matches within a single path segment and never crosses `/`,
// so a segment glob must not badge a file nested a directory deeper. The
// positive control (`web/*` DOES catch `web/app.css` in the same segment) proves
// the boundary is the slash, not the glob failing outright — this pins the
// OS-independent slash semantics that filepath.Match would break on Windows.
func TestAnnotateWriteSetOverlapGlobStarNeverCrossesSlash(t *testing.T) {
	nested := annotateWriteSetOverlapFixtures([]writeSetOverlapFixture{
		{RequestId: "REQ-1", Status: "pending", WriteSetted: []string{"web/*"}},
		{RequestId: "REQ-2", Status: "pending", WriteSetted: []string{"web/a/b.css"}},
	})
	if nested["REQ-1"] != nil || nested["REQ-2"] != nil {
		t.Fatalf("`web/*` must not cross `/` to match `web/a/b.css`, got %v", nested)
	}

	sameSegment := annotateWriteSetOverlapFixtures([]writeSetOverlapFixture{
		{RequestId: "REQ-1", Status: "pending", WriteSetted: []string{"web/*"}},
		{RequestId: "REQ-2", Status: "pending", WriteSetted: []string{"web/app.css"}},
	})
	if !reflect.DeepEqual(sameSegment["REQ-1"], []string{"REQ-2"}) ||
		!reflect.DeepEqual(sameSegment["REQ-2"], []string{"REQ-1"}) {
		t.Fatalf("`web/*` must match `web/app.css` in the same segment, got %v", sameSegment)
	}
}

// A malformed glob (unterminated character class) makes path.Match return
// ErrBadPattern; writeSetPatternsIntersect treats that as no-match for that
// direction rather than propagating the error, so the pair simply does not badge.
func TestAnnotateWriteSetOverlapMalformedPatternMatchesNothing(t *testing.T) {
	overlapsById := annotateWriteSetOverlapFixtures([]writeSetOverlapFixture{
		{RequestId: "REQ-1", Status: "pending", WriteSetted: []string{"web/["}},
		{RequestId: "REQ-2", Status: "pending", WriteSetted: []string{"web/app.css"}},
	})
	if overlapsById["REQ-1"] != nil || overlapsById["REQ-2"] != nil {
		t.Fatalf("a malformed pattern must match nothing (ErrBadPattern ⇒ no-match), got %v", overlapsById)
	}
}

// The annotation is display-only and must stay out of column placement: two
// pending REQs contending on the same file are both still Ready to work.
func TestWriteSetOverlapNeverAffectsColumnPlacement(t *testing.T) {
	repoRoot := t.TempDir()
	queueDirectory := filepath.Join(repoRoot, "do-work", "queue")
	if mkdirError := os.MkdirAll(queueDirectory, 0o755); mkdirError != nil {
		t.Fatalf("mkdir: %v", mkdirError)
	}
	for _, fixture := range []struct{ requestId, writeSet string }{
		{"REQ-1", "[web/board.css, web/board.js]"},
		{"REQ-2", "[web/board.css]"},
	} {
		fixtureContent := "---\nid: " + fixture.requestId + "\ntitle: Fixture " + fixture.requestId +
			"\nstatus: pending\nwrite_set: " + fixture.writeSet + "\n---\n\nBody.\n"
		fixturePath := filepath.Join(queueDirectory, fixture.requestId+"-fixture.md")
		if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
			t.Fatalf("write %s: %v", fixturePath, writeError)
		}
	}

	board, buildError := buildBoard(repoRoot, time.Date(2026, 7, 29, 0, 0, 0, 0, time.UTC), defaultRecentWindow, nil)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	if len(board.Columns.PendingReady) != 2 {
		t.Fatalf("both contending REQs must stay Ready (display-only annotation), got %+v", board.Columns.PendingReady)
	}
	if !reflect.DeepEqual(board.RequestsById["REQ-1"].WriteSetOverlaps, []string{"REQ-2"}) ||
		!reflect.DeepEqual(board.RequestsById["REQ-2"].WriteSetOverlaps, []string{"REQ-1"}) {
		t.Fatalf("buildBoard must annotate both cards with each other, got %v / %v",
			board.RequestsById["REQ-1"].WriteSetOverlaps, board.RequestsById["REQ-2"].WriteSetOverlaps)
	}
}

// assigned_to is the advisory cooperative claim marker: a pending REQ earmarked for a
// named session, which another session's default scan skips and reports. The board reads
// it VERBATIM and DISPLAY-ONLY — same class as write_set. These tests pin both halves:
// the value survives parsing unaltered, and it never moves a card between columns.
func TestParseRequestTicketReadsAssignedToVerbatim(t *testing.T) {
	temporaryDirectory := t.TempDir()
	fixturePath := filepath.Join(temporaryDirectory, "REQ-560-earmarked.md")
	// Mixed case with surrounding whitespace inside the quotes: a normalizing reader would
	// fold or trim it, and the verbatim-read contract says nothing may.
	fixtureContent := `---
id: REQ-560
title: Earmarked for another checkout
status: pending
assigned_to: "Cloud-Alpha_2"
---

Body.
`
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}

	ticket, parseError := parseRequestTicket(fixturePath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if ticket.AssignedTo != "Cloud-Alpha_2" {
		t.Fatalf("AssignedTo = %q, want %q (read verbatim, never normalized)", ticket.AssignedTo, "Cloud-Alpha_2")
	}
}

func TestParseRequestTicketLeavesAssignedToEmptyWhenAbsent(t *testing.T) {
	temporaryDirectory := t.TempDir()
	fixturePath := filepath.Join(temporaryDirectory, "REQ-561-unassigned.md")
	fixtureContent := "---\nid: REQ-561\ntitle: Nobody's yet\nstatus: pending\n---\n\nBody.\n"
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}

	ticket, parseError := parseRequestTicket(fixturePath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if ticket.AssignedTo != "" {
		t.Fatalf("AssignedTo = %q, want empty — absence must read as unassigned, not as a value", ticket.AssignedTo)
	}
}

func TestAssignedToNeverAffectsColumnPlacement(t *testing.T) {
	temporaryDirectory := t.TempDir()
	unassignedPath := filepath.Join(temporaryDirectory, "REQ-562-plain.md")
	assignedPath := filepath.Join(temporaryDirectory, "REQ-563-earmarked.md")
	if writeError := os.WriteFile(unassignedPath,
		[]byte("---\nid: REQ-562\ntitle: Plain\nstatus: pending\n---\n\nBody.\n"), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}
	if writeError := os.WriteFile(assignedPath,
		[]byte("---\nid: REQ-563\ntitle: Earmarked\nstatus: pending\nassigned_to: \"cloud-alpha\"\n---\n\nBody.\n"),
		0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}

	unassignedTicket, parseError := parseRequestTicket(unassignedPath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	assignedTicket, parseError := parseRequestTicket(assignedPath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if assignedTicket.Status != unassignedTicket.Status {
		t.Fatalf("assigned ticket Status = %q, unassigned = %q — assigned_to must not touch status, which is what buckets the card",
			assignedTicket.Status, unassignedTicket.Status)
	}
}

// The verbatim-read contract's load-bearing half is that NOTHING is normalized
// against a vocabulary: case is preserved and no alias map applies. Surrounding
// whitespace is the one documented exception, shared by every field in the class
// (write_set and prime_files trim per element through the same coercion), because
// padding survives only explicit YAML quoting and means nothing in a name. Raised by
// review on PR #128 — the original contract text over-claimed "no trimming".
func TestParseRequestTicketPreservesAssignedToCaseAndTrimsOnlyPadding(t *testing.T) {
	temporaryDirectory := t.TempDir()
	fixturePath := filepath.Join(temporaryDirectory, "REQ-570-padded.md")
	fixtureContent := "---\nid: REQ-570\ntitle: Padded and mixed case\nstatus: pending\n" +
		"assigned_to: \"  Cloud-ALPHA_2  \"\n---\n\nBody.\n"
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}

	ticket, parseError := parseRequestTicket(fixturePath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if ticket.AssignedTo != "Cloud-ALPHA_2" {
		t.Fatalf("AssignedTo = %q, want %q — padding trimmed, case and internal characters untouched",
			ticket.AssignedTo, "Cloud-ALPHA_2")
	}
}

// TestNormalizeSchemaFieldCoversContractAliases exercises the Schema Read
// Contract fields the shared table normalizes (the seven that had no normalizer
// before REQ-111, plus later schema additions such as effort_estimate) — the
// contract table lives in actions/work-reference.md's Schema Read Contract and
// is the source of truth for every row below.
func TestNormalizeSchemaFieldCoversContractAliases(t *testing.T) {
	testCases := []struct {
		fieldName string
		raw       string
		want      string
	}{
		// domain — the field the board was reading verbatim.
		{"domain", "backend", "backend"},
		{"domain", "back-end", "backend"},
		{"domain", "back_end", "backend"},
		{"domain", "front-end", "frontend"},
		{"domain", "front_end", "frontend"},
		{"domain", "ui_design", "ui-design"},
		{"domain", "sec", "security"},
		{"domain", "test", "testing"},
		{"domain", "cms", "cms"},
		{"domain", "CMS", "cms"},
		{"domain", "content-management", "cms"},
		{"domain", "content_management", "cms"},
		{"domain", "  Back-End  ", "backend"},
		// route — lowercase letters uppercase.
		{"route", "a", "A"},
		{"route", "B", "B"},
		{"route", " c ", "C"},
		// caveman — truthy strings plus the intensity alias.
		{"caveman", "yes", "true"},
		{"caveman", "on", "true"},
		{"caveman", "light", "lite"},
		{"caveman", "ultra", "ultra"},
		{"caveman", "false", "false"},
		// maintenance — YAML boolean plus truthy/falsey strings.
		{"maintenance", "yes", "true"},
		{"maintenance", "t", "true"},
		{"maintenance", "off", "false"},
		{"maintenance", "f", "false"},
		// tdd — same, plus test_first.
		{"tdd", "test_first", "true"},
		{"tdd", "yes", "true"},
		{"tdd", "no", "false"},
		{"tdd", "TRUE", "true"},
		// error_type — no aliases identified by the contract; canonical passthrough.
		{"error_type", "environment", "environment"},
		{"error_type", " Code ", "code"},
		// kb_status — two aliases.
		{"kb_status", "skip", "skipped"},
		{"kb_status", "rejected", "declined"},
		{"kb_status", "promoted", "promoted"},
		// effort_estimate — closed two-value enum, case-folded, plus the two
		// read-only legacy aliases that keep every REQ written before the
		// impact/effort rename resolving unchanged.
		{"effort_estimate", "effort-mechanical", "effort-mechanical"},
		{"effort_estimate", " EFFORT-MECHANICAL ", "effort-mechanical"},
		{"effort_estimate", "Effort-Substantive", "effort-substantive"},
		{"effort_estimate", "trivial", "effort-mechanical"},
		{"effort_estimate", " TRIVIAL ", "effort-mechanical"},
		{"effort_estimate", "Normal", "effort-substantive"},
		// impact — no aliases; new prefix-unique vocabulary, case-folded.
		{"impact", "impact-critical", "impact-critical"},
		{"impact", " Impact-Negligible ", "impact-negligible"},
	}
	for _, testCase := range testCases {
		if got := normalizeSchemaField(testCase.fieldName, testCase.raw); got != testCase.want {
			t.Fatalf("normalizeSchemaField(%q, %q) = %q, want %q",
				testCase.fieldName, testCase.raw, got, testCase.want)
		}
	}
}

// TestResolveSchemaFieldFallsBackWithoutSilentRemap covers the never-silently-drop
// leg of the contract: an unrecognized value resolves to the documented default
// AND reports itself unrecognized, so the caller can emit the warning.
func TestResolveSchemaFieldFallsBackWithoutSilentRemap(t *testing.T) {
	testCases := []struct {
		fieldName      string
		raw            string
		wantValue      string
		wantRecognized bool
	}{
		{"domain", "back-end", "backend", true},
		{"domain", "quantum", "general", false},
		{"caveman", "medium", "false", false},
		{"maintenance", "maybe", "false", false},
		{"tdd", "sometimes", "false", false},
		{"error_type", "cosmic-ray", "code", false},
		{"kb_status", "half-promoted", "pending", false},
		{"effort_estimate", "huge", "effort-substantive", false},
		{"effort_estimate", "", "effort-substantive", true},
		{"impact", "sort-of-important", "impact-user-visible", false},
		{"impact", "", "impact-user-visible", true},
		{"route", "Z", "", false},
		// An absent field is not an unrecognized value — it resolves to the
		// default and must NOT warn, or every REQ omitting an optional field
		// would emit noise.
		{"domain", "", "general", true},
		{"tdd", "", "false", true},
	}
	for _, testCase := range testCases {
		gotValue, gotRecognized := resolveSchemaField(testCase.fieldName, testCase.raw)
		if gotValue != testCase.wantValue || gotRecognized != testCase.wantRecognized {
			t.Fatalf("resolveSchemaField(%q, %q) = (%q, %v), want (%q, %v)",
				testCase.fieldName, testCase.raw,
				gotValue, gotRecognized, testCase.wantValue, testCase.wantRecognized)
		}
	}
}

// TestNormalizeSchemaFieldDispatchesStatusToItsOwnNormalizer keeps the two
// pre-existing normalizers authoritative for their fields rather than forking
// their alias maps into the new table.
func TestNormalizeSchemaFieldDispatchesStatusToItsOwnNormalizer(t *testing.T) {
	if got := normalizeSchemaField("status", "done"); got != "completed" {
		t.Fatalf(`normalizeSchemaField("status", "done") = %q, want "completed"`, got)
	}
	if got := normalizeSchemaField("testing_status", "in_testing"); got != "in-testing" {
		t.Fatalf(`normalizeSchemaField("testing_status", "in_testing") = %q, want "in-testing"`, got)
	}
}

// TestParseRequestTicketNormalizesDomain is REQ-111's stated RED case: the board
// read domain verbatim via coerceScalarToString, so an aliased value reached the
// card unchanged and silently mis-selected the crew file at work.md Step 6.
func TestParseRequestTicketNormalizesDomain(t *testing.T) {
	temporaryDirectory := t.TempDir()
	fixturePath := filepath.Join(temporaryDirectory, "REQ-901-aliased-domain.md")
	fixtureContent := `---
id: REQ-901
title: Aliased domain value
status: pending
domain: back-end
---

Body.
`
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}
	ticket, parseError := parseRequestTicket(fixturePath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if ticket.Domain != "backend" {
		t.Fatalf("Domain = %q, want %q — domain must normalize per the Schema Read Contract",
			ticket.Domain, "backend")
	}
}

// TestParseRequestTicketPreservesAbsentDomain guards the regression the REQ-111
// wiring nearly shipped: resolveSchemaField maps an ABSENT field to the
// contract's default, which is right for a reader that must pick a crew file
// (work.md Step 6) and wrong for the renderer. board-cards.js and
// board-filters.js gate on `if (request.domain)`, so defaulting absence
// to "general" would give every domain-less card a badge and a filter entry it
// never had.
func TestParseRequestTicketPreservesAbsentDomain(t *testing.T) {
	temporaryDirectory := t.TempDir()
	fixturePath := filepath.Join(temporaryDirectory, "REQ-902-no-domain.md")
	fixtureContent := `---
id: REQ-902
title: No domain field at all
status: pending
---

Body.
`
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}
	ticket, parseError := parseRequestTicket(fixturePath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if ticket.Domain != "" {
		t.Fatalf("Domain = %q, want %q — an absent domain must stay absent for the renderer",
			ticket.Domain, "")
	}
}

// TestUnrecognizedDomainFlagsAndWarns is REQ-117's stated RED case, and it is the
// exact shape TestUnrecognizedTestingStatusFlagsAndWarns already holds for the
// sibling field. Before REQ-111 the board read domain verbatim, so a typo was at
// least visible on the card; wiring it through resolveSchemaField silently
// substituted the contract's default and discarded the recognized flag, which made
// the typo *less* visible than it had been. The value may default — the contract
// says `general` — but the footprint is not optional (Schema Read Contract item 3,
// "Never silently drop").
func TestUnrecognizedDomainFlagsAndWarns(t *testing.T) {
	repoRoot := t.TempDir()
	queueDir := filepath.Join(repoRoot, "do-work", "queue")
	if mkdirError := os.MkdirAll(queueDir, 0o755); mkdirError != nil {
		t.Fatalf("mkdir: %v", mkdirError)
	}
	reqFileContent := "---\nid: REQ-0103\ntitle: Fixture\nstatus: pending\ndomain: quantum\n---\nbody\n"
	if writeError := os.WriteFile(filepath.Join(queueDir, "REQ-0103-bad-domain.md"), []byte(reqFileContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}

	board, buildError := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, nil)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	ticket := board.RequestsById["REQ-0103"]
	if ticket == nil {
		t.Fatalf("REQ-0103 not parsed")
	}
	if !ticket.DomainUnrecognized {
		t.Errorf("DomainUnrecognized = false, want true")
	}
	if ticket.Domain != "general" {
		t.Errorf("Domain = %q, want %q — the contract's documented default still applies", ticket.Domain, "general")
	}
	if ticket.OriginalDomain != "quantum" {
		t.Errorf("OriginalDomain = %q, want %q — the warning has to name what was actually written",
			ticket.OriginalDomain, "quantum")
	}
	warningFound := false
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "domain") && strings.Contains(warningText, "quantum") {
			warningFound = true
		}
	}
	if !warningFound {
		t.Errorf("no domain warning naming the written value; warnings=%v", board.Warnings)
	}
}

// TestUnrecognizedRouteFlagsAndWarns is REQ-119's stated RED case, and it closes
// the asymmetry REQ-116 and REQ-117 left between two fields of the same class:
// REQ-116 taught the board to normalize `route` before REQ-117 established the
// warning channel, so `domain: quantum` warned while `route: z` passed in silence.
//
// The resolved value deliberately stays `Z` rather than being defaulted — route's
// documented default is the empty string, so substituting it would make an
// unrecognized letter indistinguishable from an absent field and destroy the
// evidence a re-triage reads. The footprint is what was missing, not the value.
func TestUnrecognizedRouteFlagsAndWarns(t *testing.T) {
	repoRoot := t.TempDir()
	queueDir := filepath.Join(repoRoot, "do-work", "queue")
	if mkdirError := os.MkdirAll(queueDir, 0o755); mkdirError != nil {
		t.Fatalf("mkdir: %v", mkdirError)
	}
	reqFileContent := "---\nid: REQ-0106\ntitle: Fixture\nstatus: pending\nroute: z\n---\nbody\n"
	if writeError := os.WriteFile(filepath.Join(queueDir, "REQ-0106-bad-route.md"), []byte(reqFileContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}

	board, buildError := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, nil)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	ticket := board.RequestsById["REQ-0106"]
	if ticket == nil {
		t.Fatalf("REQ-0106 not parsed")
	}
	if !ticket.RouteUnrecognized {
		t.Errorf("RouteUnrecognized = false, want true")
	}
	if ticket.Route != "Z" {
		t.Errorf("Route = %q, want %q — an unrecognized route is reported case-folded, never blanked or defaulted",
			ticket.Route, "Z")
	}
	warningFound := false
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "route") && strings.Contains(warningText, "z") {
			warningFound = true
		}
	}
	if !warningFound {
		t.Errorf("no route warning naming the written value; warnings=%v", board.Warnings)
	}
}

// TestRecognizedRouteRaisesNoWarning is route's silence half — the alias case
// (`a` → `A`) and the absent case must both stay quiet, for the same reason the
// domain equivalent does.
func TestRecognizedRouteRaisesNoWarning(t *testing.T) {
	repoRoot := t.TempDir()
	queueDir := filepath.Join(repoRoot, "do-work", "queue")
	if mkdirError := os.MkdirAll(queueDir, 0o755); mkdirError != nil {
		t.Fatalf("mkdir: %v", mkdirError)
	}
	for _, fixture := range []struct {
		fileName        string
		requestId       string
		frontmatterLine string
	}{
		{"REQ-0107-lowercase-route.md", "REQ-0107", "route: a\n"},
		{"REQ-0108-no-route.md", "REQ-0108", ""},
	} {
		reqFileContent := "---\nid: " + fixture.requestId +
			"\ntitle: Fixture\nstatus: pending\n" + fixture.frontmatterLine + "---\nbody\n"
		if writeError := os.WriteFile(filepath.Join(queueDir, fixture.fileName), []byte(reqFileContent), 0o644); writeError != nil {
			t.Fatalf("write fixture %s: %v", fixture.fileName, writeError)
		}
	}

	board, buildError := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, nil)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "route") {
			t.Errorf("unexpected route warning for a recognized alias / absent field: %q", warningText)
		}
	}
	aliasTicket := board.RequestsById["REQ-0107"]
	if aliasTicket == nil {
		t.Fatalf("REQ-0107 not parsed")
	}
	if aliasTicket.Route != "A" || aliasTicket.RouteUnrecognized {
		t.Errorf("alias ticket: Route = %q, RouteUnrecognized = %v — want %q, false",
			aliasTicket.Route, aliasTicket.RouteUnrecognized, "A")
	}
	absentTicket := board.RequestsById["REQ-0108"]
	if absentTicket == nil {
		t.Fatalf("REQ-0108 not parsed")
	}
	if absentTicket.Route != "" || absentTicket.RouteUnrecognized {
		t.Errorf("absent-route ticket: Route = %q, RouteUnrecognized = %v — want empty, false",
			absentTicket.Route, absentTicket.RouteUnrecognized)
	}
}

// TestUnrecognizedEffortEstimateFlagsAndWarns mirrors the domain shape for the
// effort_estimate field (REQ-122): an off-vocabulary PRESENT value resolves to
// the contract default (normal), keeps the verbatim original, and leaves a
// footprint in board.Warnings — never a silent drop.
func TestUnrecognizedEffortEstimateFlagsAndWarns(t *testing.T) {
	repoRoot := t.TempDir()
	queueDir := filepath.Join(repoRoot, "do-work", "queue")
	if mkdirError := os.MkdirAll(queueDir, 0o755); mkdirError != nil {
		t.Fatalf("mkdir: %v", mkdirError)
	}
	reqFileContent := "---\nid: REQ-0110\ntitle: Fixture\nstatus: pending\neffort_estimate: huge\n---\nbody\n"
	if writeError := os.WriteFile(filepath.Join(queueDir, "REQ-0110-bad-effort.md"), []byte(reqFileContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}

	board, buildError := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, nil)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	ticket := board.RequestsById["REQ-0110"]
	if ticket == nil {
		t.Fatalf("REQ-0110 not parsed")
	}
	if !ticket.EffortEstimateUnrecognized {
		t.Errorf("EffortEstimateUnrecognized = false, want true")
	}
	if ticket.EffortEstimate != "effort-substantive" {
		t.Errorf("EffortEstimate = %q, want %q — the contract's documented default still applies", ticket.EffortEstimate, "effort-substantive")
	}
	if ticket.OriginalEffortEstimate != "huge" {
		t.Errorf("OriginalEffortEstimate = %q, want %q — the warning has to name what was actually written",
			ticket.OriginalEffortEstimate, "huge")
	}
	warningFound := false
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "effort_estimate") && strings.Contains(warningText, "huge") {
			warningFound = true
		}
	}
	if !warningFound {
		t.Errorf("no effort_estimate warning naming the written value; warnings=%v", board.Warnings)
	}
}

// TestImpactFieldFollowsPresentValueOnlyContract covers the field's three states
// at the buildBoard seam. The absent case is the one that matters most: absence
// must leave "" here so the card chip never fires on a REQ that never carried the
// field, even though every reader asking what an absent impact MEANS gets
// impact-user-visible — and never impact-negligible, which would let a filter
// built on this field skip the whole pre-rename queue.
func TestImpactFieldFollowsPresentValueOnlyContract(t *testing.T) {
	repoRoot := t.TempDir()
	queueDir := filepath.Join(repoRoot, "do-work", "queue")
	if mkdirError := os.MkdirAll(queueDir, 0o755); mkdirError != nil {
		t.Fatalf("mkdir: %v", mkdirError)
	}
	for _, fixture := range []struct {
		fileName        string
		requestId       string
		frontmatterLine string
	}{
		{"REQ-0120-bad-impact.md", "REQ-0120", "impact: quite-important\n"},
		{"REQ-0121-cased-impact.md", "REQ-0121", "impact: Impact-Negligible\n"},
		{"REQ-0122-no-impact.md", "REQ-0122", ""},
	} {
		reqFileContent := "---\nid: " + fixture.requestId +
			"\ntitle: Fixture\nstatus: pending\n" + fixture.frontmatterLine + "---\nbody\n"
		if writeError := os.WriteFile(filepath.Join(queueDir, fixture.fileName), []byte(reqFileContent), 0o644); writeError != nil {
			t.Fatalf("write fixture %s: %v", fixture.fileName, writeError)
		}
	}

	board, buildError := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, nil)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}

	unrecognizedTicket := board.RequestsById["REQ-0120"]
	if unrecognizedTicket == nil {
		t.Fatalf("REQ-0120 not parsed")
	}
	if !unrecognizedTicket.ImpactUnrecognized {
		t.Errorf("ImpactUnrecognized = false, want true — a typo'd impact may not vanish silently")
	}
	if unrecognizedTicket.Impact != "impact-user-visible" {
		t.Errorf("Impact = %q, want %q — the contract's documented default still applies",
			unrecognizedTicket.Impact, "impact-user-visible")
	}
	if unrecognizedTicket.OriginalImpact != "quite-important" {
		t.Errorf("OriginalImpact = %q, want %q — the warning has to name what was actually written",
			unrecognizedTicket.OriginalImpact, "quite-important")
	}
	impactWarningFound := false
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "impact") && strings.Contains(warningText, "quite-important") {
			impactWarningFound = true
		}
	}
	if !impactWarningFound {
		t.Errorf("no impact warning naming the written value; got %v", board.Warnings)
	}

	casedTicket := board.RequestsById["REQ-0121"]
	if casedTicket == nil {
		t.Fatalf("REQ-0121 not parsed")
	}
	if casedTicket.Impact != "impact-negligible" || casedTicket.ImpactUnrecognized {
		t.Errorf("cased ticket: Impact = %q, ImpactUnrecognized = %v — want %q, false",
			casedTicket.Impact, casedTicket.ImpactUnrecognized, "impact-negligible")
	}

	absentTicket := board.RequestsById["REQ-0122"]
	if absentTicket == nil {
		t.Fatalf("REQ-0122 not parsed")
	}
	if absentTicket.Impact != "" || absentTicket.ImpactUnrecognized {
		t.Errorf("absent-impact ticket: Impact = %q, ImpactUnrecognized = %v — want empty, false; "+
			"defaulting at the parse layer would chip every legacy card",
			absentTicket.Impact, absentTicket.ImpactUnrecognized)
	}
}

// TestRecognizedEffortEstimateRaisesNoWarning is the silence half: a case-folded
// canonical value and an absent field must both stay quiet, and the absent field
// must stay empty at the parse layer (board-cards.js chips only on the literal
// "effort-mechanical", so defaulting absence would be invisible today — the guard
// is what keeps that true if the frontend gate ever loosens). The cased fixture
// spells the PRE-RENAME value on purpose: `Trivial` is now a read-only legacy
// alias, and this is where the forty-odd archived REQs carrying it are proved to
// still resolve — through buildBoard, the seam a real board read goes through.
func TestRecognizedEffortEstimateRaisesNoWarning(t *testing.T) {
	repoRoot := t.TempDir()
	queueDir := filepath.Join(repoRoot, "do-work", "queue")
	if mkdirError := os.MkdirAll(queueDir, 0o755); mkdirError != nil {
		t.Fatalf("mkdir: %v", mkdirError)
	}
	for _, fixture := range []struct {
		fileName        string
		requestId       string
		frontmatterLine string
	}{
		{"REQ-0111-cased-effort.md", "REQ-0111", "effort_estimate: Trivial\n"},
		{"REQ-0112-no-effort.md", "REQ-0112", ""},
	} {
		reqFileContent := "---\nid: " + fixture.requestId +
			"\ntitle: Fixture\nstatus: pending\n" + fixture.frontmatterLine + "---\nbody\n"
		if writeError := os.WriteFile(filepath.Join(queueDir, fixture.fileName), []byte(reqFileContent), 0o644); writeError != nil {
			t.Fatalf("write fixture %s: %v", fixture.fileName, writeError)
		}
	}

	board, buildError := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, nil)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "effort_estimate") {
			t.Errorf("unexpected effort_estimate warning for a recognized / absent value: %q", warningText)
		}
	}
	casedTicket := board.RequestsById["REQ-0111"]
	if casedTicket == nil {
		t.Fatalf("REQ-0111 not parsed")
	}
	if casedTicket.EffortEstimate != "effort-mechanical" || casedTicket.EffortEstimateUnrecognized {
		t.Errorf("cased ticket: EffortEstimate = %q, EffortEstimateUnrecognized = %v — want %q, false — "+
			"the legacy `trivial` alias must still resolve, or every archived REQ predating the rename breaks",
			casedTicket.EffortEstimate, casedTicket.EffortEstimateUnrecognized, "effort-mechanical")
	}
	absentTicket := board.RequestsById["REQ-0112"]
	if absentTicket == nil {
		t.Fatalf("REQ-0112 not parsed")
	}
	if absentTicket.EffortEstimate != "" || absentTicket.EffortEstimateUnrecognized {
		t.Errorf("absent-effort ticket: EffortEstimate = %q, EffortEstimateUnrecognized = %v — want empty, false",
			absentTicket.EffortEstimate, absentTicket.EffortEstimateUnrecognized)
	}
}

// TestEffortEstimateNeverAffectsColumnPlacement pins the display-only guarantee:
// the chip is information for a human's pick, and a trivial REQ buckets exactly
// like its unmarked twin.
func TestEffortEstimateNeverAffectsColumnPlacement(t *testing.T) {
	temporaryDirectory := t.TempDir()
	plainPath := filepath.Join(temporaryDirectory, "REQ-572-plain.md")
	trivialPath := filepath.Join(temporaryDirectory, "REQ-573-trivial.md")
	if writeError := os.WriteFile(plainPath,
		[]byte("---\nid: REQ-572\ntitle: Plain\nstatus: pending\n---\n\nBody.\n"), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}
	if writeError := os.WriteFile(trivialPath,
		[]byte("---\nid: REQ-573\ntitle: Trivial twin\nstatus: pending\neffort_estimate: trivial\n---\n\nBody.\n"),
		0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}

	plainTicket, parseError := parseRequestTicket(plainPath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	trivialTicket, parseError := parseRequestTicket(trivialPath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if trivialTicket.Status != plainTicket.Status {
		t.Fatalf("trivial ticket Status = %q, plain = %q — effort_estimate must not touch status, which is what buckets the card",
			trivialTicket.Status, plainTicket.Status)
	}
}

// TestRecognizedDomainRaisesNoWarning is the other half, and the one that keeps
// the channel worth reading: a documented alias must resolve in silence. If every
// `domain: back-end` also warned, the warnings list would fill with noise on a
// real queue and readers would learn to ignore it — which is the failure mode the
// contract's absent-field carve-out exists to prevent.
func TestRecognizedDomainRaisesNoWarning(t *testing.T) {
	repoRoot := t.TempDir()
	queueDir := filepath.Join(repoRoot, "do-work", "queue")
	if mkdirError := os.MkdirAll(queueDir, 0o755); mkdirError != nil {
		t.Fatalf("mkdir: %v", mkdirError)
	}
	for _, fixture := range []struct {
		fileName        string
		requestId       string
		frontmatterLine string
	}{
		{"REQ-0104-alias-domain.md", "REQ-0104", "domain: back-end\n"},
		{"REQ-0105-no-domain.md", "REQ-0105", ""},
	} {
		reqFileContent := "---\nid: " + fixture.requestId +
			"\ntitle: Fixture\nstatus: pending\n" + fixture.frontmatterLine + "---\nbody\n"
		if writeError := os.WriteFile(filepath.Join(queueDir, fixture.fileName), []byte(reqFileContent), 0o644); writeError != nil {
			t.Fatalf("write fixture %s: %v", fixture.fileName, writeError)
		}
	}

	board, buildError := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, nil)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "domain") {
			t.Errorf("unexpected domain warning for a recognized alias / absent field: %q", warningText)
		}
	}
	aliasTicket := board.RequestsById["REQ-0104"]
	if aliasTicket == nil {
		t.Fatalf("REQ-0104 not parsed")
	}
	if aliasTicket.Domain != "backend" || aliasTicket.DomainUnrecognized {
		t.Errorf("alias ticket: Domain = %q, DomainUnrecognized = %v — want %q, false",
			aliasTicket.Domain, aliasTicket.DomainUnrecognized, "backend")
	}
	absentTicket := board.RequestsById["REQ-0105"]
	if absentTicket == nil {
		t.Fatalf("REQ-0105 not parsed")
	}
	if absentTicket.Domain != "" || absentTicket.DomainUnrecognized {
		t.Errorf("absent-domain ticket: Domain = %q, DomainUnrecognized = %v — want empty, false",
			absentTicket.Domain, absentTicket.DomainUnrecognized)
	}
}

// TestParseRequestTicketNormalizesRoute is REQ-116's stated RED case. REQ-111
// added the route contract row and its normalizer but wired only `domain` at the
// read site, so the board kept reading route through coerceScalarToString — which
// trims and nothing else. A lowercase route letter therefore reached the card
// badge and the drawer row as written, in a field the contract says is uppercase.
//
// Parse-level on purpose: the normalizer's own table test (route a→A, B→B, " c "→C)
// passed the whole time the board was wrong, so only a test that goes through
// parseRequestTicket can hold this line.
func TestParseRequestTicketNormalizesRoute(t *testing.T) {
	for _, testCase := range []struct {
		name          string
		routeLine     string
		expectedRoute string
	}{
		{"lowercase letter uppercases", "route: a", "A"},
		{"canonical letter is unchanged", "route: B", "B"},
		{"padded lowercase letter uppercases", "route: \" c \"", "C"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			temporaryDirectory := t.TempDir()
			fixturePath := filepath.Join(temporaryDirectory, "REQ-903-route-case.md")
			fixtureContent := `---
id: REQ-903
title: Route letter case
status: pending
` + testCase.routeLine + `
---

Body.
`
			if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
				t.Fatalf("write fixture: %v", writeError)
			}
			ticket, parseError := parseRequestTicket(fixturePath, "queue")
			if parseError != nil {
				t.Fatalf("parseRequestTicket: %v", parseError)
			}
			if ticket.Route != testCase.expectedRoute {
				t.Fatalf("Route = %q, want %q — route must normalize per the Schema Read Contract",
					ticket.Route, testCase.expectedRoute)
			}
		})
	}
}

// TestParseRequestTicketPreservesAbsentRoute is the absent-field half of the same
// wiring, and the reason the read site must not call resolveSchemaField: that
// helper substitutes the field's documented default, and route's default is the
// empty string, so absence and an unrecognized letter would become the same
// value. board-cards.js and board-detail.js gate their route surfaces on `if (request.route)`.
func TestParseRequestTicketPreservesAbsentRoute(t *testing.T) {
	temporaryDirectory := t.TempDir()
	fixturePath := filepath.Join(temporaryDirectory, "REQ-904-no-route.md")
	fixtureContent := `---
id: REQ-904
title: No route field at all
status: pending
---

Body.
`
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}
	ticket, parseError := parseRequestTicket(fixturePath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if ticket.Route != "" {
		t.Fatalf("Route = %q, want %q — an absent route must stay absent for the renderer",
			ticket.Route, "")
	}
}

// TestParseRequestTicketReportsUnrecognizedRouteUnchanged pins the no-default half
// of route's contract row: an unrecognized letter means the REQ needs re-triage,
// which is not a value this parser can invent. Blanking it — what
// resolveSchemaField would do, since route's documented default is "" — would
// destroy the evidence that re-triage reads.
func TestParseRequestTicketReportsUnrecognizedRouteUnchanged(t *testing.T) {
	temporaryDirectory := t.TempDir()
	fixturePath := filepath.Join(temporaryDirectory, "REQ-905-unknown-route.md")
	fixtureContent := `---
id: REQ-905
title: Route letter outside the enum
status: pending
route: z
---

Body.
`
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}
	ticket, parseError := parseRequestTicket(fixturePath, "queue")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if ticket.Route != "Z" {
		t.Fatalf("Route = %q, want %q — an unrecognized route is reported case-folded, never blanked",
			ticket.Route, "Z")
	}
}

func TestParseRequestTicketCarriesFailureDetails(t *testing.T) {
	temporaryDirectory := t.TempDir()
	fixturePath := filepath.Join(temporaryDirectory, "REQ-906-failed.md")
	fixtureContent := `---
id: REQ-906
title: Failed request
status: failed
error: "  compiler exploded  "
error_type: environment
---

Body.
`
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}
	ticket, parseError := parseRequestTicket(fixturePath, "archive")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if ticket.Error != "compiler exploded" {
		t.Fatalf("Error = %q, want %q — failure text must use the row-less scalar read without schema normalization",
			ticket.Error, "compiler exploded")
	}
	if ticket.ErrorType != "environment" || ticket.OriginalErrorType != "environment" || ticket.ErrorTypeUnrecognized {
		t.Fatalf("error_type provenance = (%q, %q, %v), want (%q, %q, false)",
			ticket.ErrorType, ticket.OriginalErrorType, ticket.ErrorTypeUnrecognized,
			"environment", "environment")
	}
}

func TestParseRequestTicketPreservesAbsentErrorType(t *testing.T) {
	temporaryDirectory := t.TempDir()
	fixturePath := filepath.Join(temporaryDirectory, "REQ-907-unclassified-failure.md")
	fixtureContent := `---
id: REQ-907
title: Unclassified failure
status: failed
error: tests never reached green
---

Body.
`
	if writeError := os.WriteFile(fixturePath, []byte(fixtureContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}
	ticket, parseError := parseRequestTicket(fixturePath, "archive")
	if parseError != nil {
		t.Fatalf("parseRequestTicket: %v", parseError)
	}
	if ticket.Error != "tests never reached green" {
		t.Fatalf("Error = %q, want the recorded failure reason", ticket.Error)
	}
	if ticket.ErrorType != "" || ticket.OriginalErrorType != "" || ticket.ErrorTypeUnrecognized {
		t.Fatalf("absent error_type became (%q, %q, %v) — absence must never fabricate the contract default 'code'",
			ticket.ErrorType, ticket.OriginalErrorType, ticket.ErrorTypeUnrecognized)
	}
	for _, warningText := range collectSchemaFieldWarnings([]*RequestTicket{ticket}) {
		if strings.Contains(warningText, "error_type") {
			t.Fatalf("absent error_type emitted a warning: %q", warningText)
		}
	}
}

func TestUnrecognizedErrorTypeFlagsAndWarns(t *testing.T) {
	repoRoot := t.TempDir()
	archiveDirectory := filepath.Join(repoRoot, "do-work", "archive")
	if mkdirError := os.MkdirAll(archiveDirectory, 0o755); mkdirError != nil {
		t.Fatalf("mkdir: %v", mkdirError)
	}
	requestFileContent := "---\nid: REQ-0908\ntitle: Bad failure class\nstatus: failed\nerror: fixture failure\nerror_type: cosmic-ray\n---\nbody\n"
	if writeError := os.WriteFile(filepath.Join(archiveDirectory, "REQ-0908-bad-error-type.md"), []byte(requestFileContent), 0o644); writeError != nil {
		t.Fatalf("write fixture: %v", writeError)
	}

	board, buildError := buildBoard(repoRoot, time.Now(), 7*24*time.Hour, nil)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	ticket := board.RequestsById["REQ-0908"]
	if ticket == nil {
		t.Fatalf("REQ-0908 not parsed")
	}
	if ticket.ErrorType != "code" || ticket.OriginalErrorType != "cosmic-ray" || !ticket.ErrorTypeUnrecognized {
		t.Fatalf("unrecognized error_type provenance = (%q, %q, %v), want (%q, %q, true)",
			ticket.ErrorType, ticket.OriginalErrorType, ticket.ErrorTypeUnrecognized,
			"code", "cosmic-ray")
	}
	warningFound := false
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "error_type") && strings.Contains(warningText, "cosmic-ray") {
			warningFound = true
		}
	}
	if !warningFound {
		t.Fatalf("no error_type warning naming the written value; warnings=%v", board.Warnings)
	}
}
