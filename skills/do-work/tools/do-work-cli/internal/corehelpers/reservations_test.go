package corehelpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReservationCleanupRemovesOldAndPreservesFresh(t *testing.T) {
	repository := t.TempDir()
	root := filepath.Join(repository, "do-work", ".req-reservations")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	old := filepath.Join(root, "REQ-000101")
	fresh := filepath.Join(root, "REQ-000102")
	_ = os.WriteFile(old, nil, 0o600)
	_ = os.WriteFile(fresh, nil, 0o600)
	past := time.Now().Add(-49 * time.Hour)
	_ = os.Chtimes(old, past, past)
	result := handleCleanupReservations(testContext(repository), nil)
	if result.Outcome != "success" {
		t.Fatalf("result=%+v", result)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Fatalf("old marker remains: %v", err)
	}
	if _, err := os.Stat(fresh); err != nil {
		t.Fatalf("fresh marker removed: %v", err)
	}
}

func TestReservationCleanupRequiresCommittedRequestInGitRepository(t *testing.T) {
	repository := newGitFixture(t)
	root := filepath.Join(repository, "do-work", ".req-reservations")
	if err := os.MkdirAll(filepath.Join(repository, "do-work", "working"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(root, "REQ-000203")
	request := filepath.Join(repository, "do-work", "working", "REQ-203-capture.md")
	_ = os.WriteFile(marker, nil, 0o600)
	_ = os.WriteFile(request, []byte("---\nid: REQ-203\n---\n"), 0o600)
	result := handleCleanupReservations(testContext(repository), nil)
	if result.Outcome != "success" {
		t.Fatalf("uncommitted result=%+v", result)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("uncommitted capture marker removed: %v", err)
	}
	runFixtureGitCommand(t, repository, "add", "do-work/working/REQ-203-capture.md")
	runFixtureGitCommand(t, repository, "commit", "-qm", "land request")
	result = handleCleanupReservations(testContext(repository), nil)
	if result.Outcome != "success" {
		t.Fatalf("committed result=%+v", result)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("landed marker remains: %v", err)
	}
}

func TestReservationCleanupDryRunPreservesEligibleMarker(t *testing.T) {
	repository := t.TempDir()
	root := filepath.Join(repository, "do-work", ".req-reservations")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(root, "REQ-000777")
	if err := os.WriteFile(marker, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	past := time.Now().Add(-49 * time.Hour)
	_ = os.Chtimes(marker, past, past)
	result := handleCleanupReservations(testContext(repository), []string{"--dry-run"})
	if result.Outcome != "success" || len(result.Changes) != 1 || !strings.Contains(result.Changes[0].Detail, "would remove") {
		t.Fatalf("result=%#v", result)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("dry-run removed marker: %v", err)
	}
}
