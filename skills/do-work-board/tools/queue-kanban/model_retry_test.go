package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBuildBoardRecoversRequestMovedDuringRead(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"do-work/queue/REQ-701-moved.md", "---\nid: REQ-701\nstatus: claimed\ntitle: moved request\n---\n"},
		{"do-work/queue/REQ-703-dependent.md", "---\nid: REQ-703\nstatus: pending\ntitle: dependent request\ndepends_on: [REQ-701]\n---\n"},
		{"do-work/CHECKPOINT.md", "# In Progress\n- REQ-701: moved request\n"},
	})
	queuePath := filepath.Join(repoRoot, "do-work", "queue", "REQ-701-moved.md")
	workingPath := filepath.Join(repoRoot, "do-work", "working", "REQ-701-moved.md")
	if err := os.MkdirAll(filepath.Dir(workingPath), 0o755); err != nil {
		t.Fatal(err)
	}
	moved := false
	installTicketReader(t, func(filePath string) ([]byte, error) {
		if filePath == queuePath && !moved {
			moved = true
			if err := os.Rename(queuePath, workingPath); err != nil {
				t.Fatal(err)
			}
		}
		return os.ReadFile(filePath)
	})

	board, err := LoadBoard(repoRoot, time.Now(), 24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	ticket := board.RequestsById["REQ-701"]
	if ticket == nil {
		t.Error("REQ-701 vanished from RequestsById after moving from queue to working")
	} else if ticket.TreeSection != "working" || ticket.FilePath != workingPath {
		t.Errorf("moved ticket location = %s, %s; want working, %s", ticket.TreeSection, ticket.FilePath, workingPath)
	}
	for _, warning := range board.Warnings {
		if strings.Contains(warning, "REQ-701") && strings.Contains(warning, "not in the do-work tree") {
			t.Errorf("moved request caused false missing dependency warning: %s", warning)
		}
	}
	var report VerifyReport
	appendCheckpointGhostFindings(&report, repoRoot, board)
	if len(report.Findings) != 0 {
		t.Errorf("moved request caused false checkpoint ghost: %+v", report.Findings)
	}
}

func TestBuildBoardRecoversUserRequestArchivedDuringRead(t *testing.T) {
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"do-work/queue/REQ-701-linked.md", "---\nid: REQ-701\nstatus: pending\ntitle: before archive\nuser_request: UR-701\n---\n"},
		{"do-work/user-requests/UR-701/input.md", "---\nid: UR-701\ntitle: archived input\n---\nOriginal input.\n"},
	})
	activeDirectory := filepath.Join(repoRoot, "do-work", "user-requests", "UR-701")
	archiveDirectory := filepath.Join(repoRoot, "do-work", "archive", "UR-701")
	activePath := filepath.Join(activeDirectory, "input.md")
	if err := os.MkdirAll(filepath.Dir(archiveDirectory), 0o755); err != nil {
		t.Fatal(err)
	}
	moved := false
	installTicketReader(t, func(filePath string) ([]byte, error) {
		if filePath == activePath && !moved {
			moved = true
			if err := os.Rename(activeDirectory, archiveDirectory); err != nil {
				t.Fatal(err)
			}
			requestPath := filepath.Join(repoRoot, "do-work", "queue", "REQ-701-linked.md")
			if err := os.WriteFile(requestPath, []byte("---\nid: REQ-701\nstatus: pending\ntitle: after archive\nuser_request: UR-701\n---\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return os.ReadFile(filePath)
	})

	board, err := LoadBoard(repoRoot, time.Now(), 24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	userRequest := board.UserRequestsById["UR-701"]
	if userRequest == nil || !userRequest.InputFilePresent || userRequest.FilePath != filepath.Join(archiveDirectory, "input.md") {
		t.Errorf("archived UR input was dropped: %+v", userRequest)
	}
	if ticket := board.RequestsById["REQ-701"]; ticket == nil || ticket.Title != "after archive" {
		t.Errorf("UR retry did not rebuild the REQ pass from scratch: %+v", ticket)
	}
	if len(board.AllRequests) != 1 || len(board.UserRequests) != 1 {
		t.Errorf("retry retained first-pass tickets: requests=%d, user requests=%d", len(board.AllRequests), len(board.UserRequests))
	}
}

func TestBuildBoardSkipsNonMissingReadErrorsWithoutRetry(t *testing.T) {
	for _, ticketType := range []string{"request", "user request"} {
		t.Run(ticketType, func(t *testing.T) {
			repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
				{"do-work/queue/REQ-701-unreadable.md", "---\nid: REQ-701\nstatus: pending\n---\n"},
				{"do-work/user-requests/UR-701/input.md", "---\nid: UR-701\ntitle: input\n---\n"},
			})
			failedPath := filepath.Join(repoRoot, "do-work", "queue", "REQ-701-unreadable.md")
			if ticketType == "user request" {
				failedPath = filepath.Join(repoRoot, "do-work", "user-requests", "UR-701", "input.md")
			}
			readCount := 0
			installTicketReader(t, func(filePath string) ([]byte, error) {
				if filePath == failedPath {
					readCount++
					return nil, &os.PathError{Op: "read", Path: filePath, Err: fs.ErrPermission}
				}
				return os.ReadFile(filePath)
			})
			board, err := LoadBoard(repoRoot, time.Now(), 24*time.Hour)
			if err != nil {
				t.Fatal(err)
			}
			if readCount != 1 {
				t.Errorf("unreadable %s was read %d times, want 1", ticketType, readCount)
			}
			if ticketType == "request" && board.RequestsById["REQ-701"] != nil {
				t.Error("unreadable REQ was not skipped")
			}
			if ticketType == "user request" && board.UserRequestsById["UR-701"] != nil {
				t.Error("unreadable UR was not skipped")
			}
		})
	}
}

func TestBuildBoardStopsAfterOneMissingFileRetry(t *testing.T) {
	for _, ticketType := range []string{"request", "user request"} {
		t.Run(ticketType, func(t *testing.T) {
			repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
				{"do-work/queue/REQ-701-missing.md", "---\nid: REQ-701\nstatus: pending\n---\n"},
				{"do-work/queue/REQ-703-kept.md", "---\nid: REQ-703\nstatus: pending\n---\n"},
				{"do-work/user-requests/UR-701/input.md", "---\nid: UR-701\n---\n"},
			})
			failedPath := filepath.Join(repoRoot, "do-work", "queue", "REQ-701-missing.md")
			if ticketType == "user request" {
				failedPath = filepath.Join(repoRoot, "do-work", "user-requests", "UR-701", "input.md")
			}
			readCount := 0
			installTicketReader(t, func(filePath string) ([]byte, error) {
				if filePath == failedPath {
					readCount++
					return nil, fmt.Errorf("ticket disappeared: %w", fs.ErrNotExist)
				}
				return os.ReadFile(filePath)
			})
			board, err := LoadBoard(repoRoot, time.Now(), 24*time.Hour)
			if err != nil {
				t.Fatal(err)
			}
			if readCount != 2 {
				t.Errorf("missing %s was read %d times, want exactly 2", ticketType, readCount)
			}
			if board.RequestsById["REQ-703"] == nil {
				t.Error("second pass dropped the readable request")
			}
			if ticketType == "request" && board.RequestsById["REQ-701"] != nil {
				t.Error("second-pass missing REQ was not skipped")
			}
			if ticketType == "user request" && board.UserRequestsById["UR-701"] != nil {
				t.Error("second-pass missing UR was not skipped")
			}
		})
	}
}

func installTicketReader(t *testing.T, reader func(string) ([]byte, error)) {
	t.Helper()
	originalReader := readTicketFile
	readTicketFile = reader
	t.Cleanup(func() { readTicketFile = originalReader })
}
