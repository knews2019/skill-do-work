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
		{"pending", "pending"},
		{"pending-answers", "pending-answers"},
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
	for _, blocked := range []string{"pending-answers", "blocked-archive-collision", "blocked-dependency-cycle", "failed"} {
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
}

func TestResolveCommitHashVariants(t *testing.T) {
	testCases := []struct {
		name   string
		fields map[string]any
		want   string
	}{
		{"canonical commit", map[string]any{"commit": "abc123"}, "abc123"},
		{"commit_hash variant", map[string]any{"commit_hash": "def456"}, "def456"},
		{"green_commit variant", map[string]any{"green_commit": "aaa111"}, "aaa111"},
		{"commit_green variant", map[string]any{"commit_green": "bbb222"}, "bbb222"},
		{"impl_commit variant", map[string]any{"impl_commit": "ccc333"}, "ccc333"},
		{"canonical wins over variant", map[string]any{"commit": "primary", "commit_hash": "secondary"}, "primary"},
		{"none present", map[string]any{}, ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := resolveCommitHash(testCase.fields); got != testCase.want {
				t.Fatalf("resolveCommitHash = %q, want %q", got, testCase.want)
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
	if ticket.TreeSection != "archive" {
		t.Fatalf("TreeSection = %q", ticket.TreeSection)
	}
	if !strings.Contains(ticket.BodyMarkdown, "# Body heading") {
		t.Fatalf("body not preserved: %q", ticket.BodyMarkdown)
	}
}

func TestBucketColumns(t *testing.T) {
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	window := 48 * time.Hour
	recentDone := &RequestTicket{RequestId: "REQ-1", Status: "completed", CompletionTime: now.Add(-1 * time.Hour)}
	oldDone := &RequestTicket{RequestId: "REQ-2", Status: "completed", CompletionTime: now.Add(-200 * time.Hour)}
	tickets := []*RequestTicket{
		{RequestId: "REQ-3", Status: "pending"},
		{RequestId: "REQ-4", Status: "claimed"},
		{RequestId: "REQ-5", Status: "pending-answers"},
		{RequestId: "REQ-6", Status: "deferred", OriginalStatus: "deferred"}, // hand-edited status outside the Schema Read Contract enum — must still land in Needs-input/Blocked, now via the unrecognized-status warning path
		{RequestId: "REQ-7", Status: "pnding", OriginalStatus: "pnding"},     // typo'd status — must never be silently dropped
		recentDone,
		oldDone,
	}
	columns, statusWarnings := bucketColumns(tickets, now, window)
	if len(columns.Pending) != 1 || columns.Pending[0].RequestId != "REQ-3" {
		t.Fatalf("Pending = %+v", columns.Pending)
	}
	if len(columns.Claimed) != 1 || columns.Claimed[0].RequestId != "REQ-4" {
		t.Fatalf("Claimed = %+v", columns.Claimed)
	}
	if len(columns.NeedsInputOrBlocked) != 3 {
		t.Fatalf("NeedsInputOrBlocked should hold pending-answers + the deferred and pnding unrecognized statuses, got %d", len(columns.NeedsInputOrBlocked))
	}
	if len(columns.RecentlyDone) != 1 || columns.RecentlyDone[0].RequestId != "REQ-1" {
		t.Fatalf("RecentlyDone should hold only the in-window completion, got %+v", columns.RecentlyDone)
	}
	if len(statusWarnings) != 2 {
		t.Fatalf("expected two unrecognized-status warnings (deferred + pnding), got %d: %v", len(statusWarnings), statusWarnings)
	}
	foundDeferredWarning := false
	foundTypoWarning := false
	for _, warning := range statusWarnings {
		if strings.Contains(warning, "REQ-6") && strings.Contains(warning, "deferred") {
			foundDeferredWarning = true
		}
		if strings.Contains(warning, "REQ-7") && strings.Contains(warning, "pnding") {
			foundTypoWarning = true
		}
	}
	if !foundDeferredWarning {
		t.Fatalf("expected an unrecognized-status warning naming REQ-6/deferred, got %v", statusWarnings)
	}
	if !foundTypoWarning {
		t.Fatalf("expected an unrecognized-status warning naming REQ-7/pnding, got %v", statusWarnings)
	}
}
