package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// futureStampSyntheticBoard builds a board over a tree seeded with the
// future-dated-timestamp shapes plus sane controls:
//
//	REQ-9401  claimed, claimed_at 2h past `now`            → flagged (the local-time-plus-Z signature)
//	REQ-9402  claimed, claimed_at 1min past `now`          → NOT flagged (inside the 2min skew allowance)
//	REQ-9403  pending, created_at + blocked-shape fields sane, testing_updated_at 3h past → flagged on that one field only
//	REQ-9404  claimed, claimed_at unparseable              → NOT flagged (unparseable is not future)
//	REQ-9405  pending, status_changed_at 2h past `now`     → flagged (flip stamps are guarded too)
//
// `now` is fixed so the assertions are deterministic.
func futureStampSyntheticBoard(t *testing.T) *Board {
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

	// fixedNow is 2026-06-30T12:00:00Z in every stamp below.
	writeFixture(filepath.Join("do-work", "working", "REQ-9401-future-claim.md"),
		requestContent("REQ-9401", "claimed", "created_at: 2026-06-30T09:00:00Z\nclaimed_at: 2026-06-30T14:00:00Z\n"))
	writeFixture(filepath.Join("do-work", "working", "REQ-9402-inside-skew.md"),
		requestContent("REQ-9402", "claimed", "created_at: 2026-06-30T09:00:00Z\nclaimed_at: 2026-06-30T12:01:00Z\n"))
	writeFixture(filepath.Join("do-work", "queue", "REQ-9403-future-testing-stamp.md"),
		requestContent("REQ-9403", "pending", "created_at: 2026-06-30T09:00:00Z\ntesting_updated_at: 2026-06-30T15:00:00Z\n"))
	writeFixture(filepath.Join("do-work", "working", "REQ-9404-unparseable-claim.md"),
		requestContent("REQ-9404", "claimed", "created_at: 2026-06-30T09:00:00Z\nclaimed_at: not-a-real-instant\n"))
	writeFixture(filepath.Join("do-work", "queue", "REQ-9405-future-flip-stamp.md"),
		requestContent("REQ-9405", "pending", "created_at: 2026-06-30T09:00:00Z\nstatus_changed_at: 2026-06-30T14:00:00Z\n"))

	neverCalledGitLookup := func(string, string) (time.Time, bool) { return time.Time{}, false }
	fixedNow := time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC)
	board, buildError := buildBoard(repoRoot, fixedNow, 7*24*time.Hour, neverCalledGitLookup)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	return board
}

func TestFutureTimestampFieldsFlaggedInBoardModel(t *testing.T) {
	board := futureStampSyntheticBoard(t)

	expectedFutureFields := map[string][]string{
		"REQ-9401": {"claimed_at 2026-06-30T14:00:00Z"},
		"REQ-9402": nil,
		"REQ-9403": {"testing_updated_at 2026-06-30T15:00:00Z"},
		"REQ-9404": nil,
		"REQ-9405": {"status_changed_at 2026-06-30T14:00:00Z"},
	}
	for requestId, expectedEntries := range expectedFutureFields {
		ticket := board.RequestsById[requestId]
		if ticket == nil {
			t.Fatalf("%s not parsed", requestId)
		}
		if len(ticket.FutureTimestampFields) != len(expectedEntries) {
			t.Fatalf("%s FutureTimestampFields = %v, want %v",
				requestId, ticket.FutureTimestampFields, expectedEntries)
		}
		for entryIndex, expectedEntry := range expectedEntries {
			if ticket.FutureTimestampFields[entryIndex] != expectedEntry {
				t.Fatalf("%s FutureTimestampFields = %v, want %v",
					requestId, ticket.FutureTimestampFields, expectedEntries)
			}
		}
	}
}

func TestFutureTimestampWarningNamesFieldAndFix(t *testing.T) {
	board := futureStampSyntheticBoard(t)

	var matchedWarnings []string
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "future-dated timestamp") {
			matchedWarnings = append(matchedWarnings, warningText)
		}
	}
	if len(matchedWarnings) != 3 {
		t.Fatalf("want exactly 3 future-timestamp warnings (REQ-9401, REQ-9403, REQ-9405), got %d: %v",
			len(matchedWarnings), matchedWarnings)
	}
	for _, expectedFragment := range []string{
		"REQ-9401 has future-dated timestamp(s): claimed_at 2026-06-30T14:00:00Z",
		"date -u +%Y-%m-%dT%H:%M:%SZ",
	} {
		found := false
		for _, warningText := range matchedWarnings {
			if strings.Contains(warningText, expectedFragment) {
				found = true
			}
		}
		if !found {
			t.Fatalf("no future-timestamp warning contains %q; warnings: %v", expectedFragment, matchedWarnings)
		}
	}
}
