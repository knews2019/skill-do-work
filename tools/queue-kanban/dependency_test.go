package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// dependencyFixture is one REQ to seed: an id, a status, and its depends_on list.
type dependencyFixture struct {
	RequestId    string
	Status       string
	DependsOnIds []string
}

// buildDependencyBoard seeds a queue of REQs with the given statuses and
// depends_on edges, then builds the board with a fixed `now` and no git lookup.
func buildDependencyBoard(t *testing.T, fixtures []dependencyFixture) *Board {
	t.Helper()
	repoRoot := t.TempDir()

	for _, fixture := range fixtures {
		frontmatter := "---\nid: " + fixture.RequestId + "\ntitle: Fixture " + fixture.RequestId +
			"\nstatus: " + fixture.Status + "\n"
		if len(fixture.DependsOnIds) > 0 {
			frontmatter += "depends_on: [" + strings.Join(fixture.DependsOnIds, ", ") + "]\n"
		}
		frontmatter += "---\n\nBody.\n"

		filePath := filepath.Join(repoRoot, "do-work", "queue", fixture.RequestId+"-fixture.md")
		if mkdirError := os.MkdirAll(filepath.Dir(filePath), 0o755); mkdirError != nil {
			t.Fatalf("mkdir: %v", mkdirError)
		}
		if writeError := os.WriteFile(filePath, []byte(frontmatter), 0o644); writeError != nil {
			t.Fatalf("write %s: %v", filePath, writeError)
		}
	}

	board, buildError := buildBoard(repoRoot, time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC), defaultRecentWindow, nil)
	if buildError != nil {
		t.Fatalf("buildBoard: %v", buildError)
	}
	return board
}

func requestIdSet(tickets []*RequestTicket) map[string]bool {
	ids := map[string]bool{}
	for _, ticket := range tickets {
		ids[ticket.RequestId] = true
	}
	return ids
}

// TestPendingSplitsOnDependencyReadiness is the core contract: a pending REQ is
// Ready only when every depends_on target reached terminal success.
func TestPendingSplitsOnDependencyReadiness(t *testing.T) {
	board := buildDependencyBoard(t, []dependencyFixture{
		{RequestId: "REQ-1", Status: "completed"},
		{RequestId: "REQ-2", Status: "pending", DependsOnIds: []string{"REQ-1"}},  // dep done → ready
		{RequestId: "REQ-3", Status: "pending", DependsOnIds: []string{"REQ-4"}},  // dep pending → waiting
		{RequestId: "REQ-4", Status: "pending"},                                   // no deps → ready
		{RequestId: "REQ-5", Status: "pending", DependsOnIds: []string{"REQ-99"}}, // dangling → waiting
	})

	readyIds := requestIdSet(board.Columns.PendingReady)
	waitingIds := requestIdSet(board.Columns.PendingWaiting)

	for _, expectedReadyId := range []string{"REQ-2", "REQ-4"} {
		if !readyIds[expectedReadyId] {
			t.Errorf("%s should be ready; ready=%v waiting=%v", expectedReadyId, readyIds, waitingIds)
		}
	}
	for _, expectedWaitingId := range []string{"REQ-3", "REQ-5"} {
		if !waitingIds[expectedWaitingId] {
			t.Errorf("%s should be waiting; ready=%v waiting=%v", expectedWaitingId, readyIds, waitingIds)
		}
	}
	if len(board.Columns.Pending) != len(board.Columns.PendingReady)+len(board.Columns.PendingWaiting) {
		t.Fatalf("ready+waiting must partition pending: %d + %d != %d",
			len(board.Columns.PendingReady), len(board.Columns.PendingWaiting), len(board.Columns.Pending))
	}
}

// TestCancelledDependencyNeverSatisfiesGating pins the rule from
// actions/work-reference.md: `cancelled` is terminal but not success, so a REQ
// depending on it stays waiting rather than quietly reading as ready.
func TestCancelledDependencyNeverSatisfiesGating(t *testing.T) {
	board := buildDependencyBoard(t, []dependencyFixture{
		{RequestId: "REQ-1", Status: "cancelled"},
		{RequestId: "REQ-2", Status: "pending", DependsOnIds: []string{"REQ-1"}},
	})

	if !requestIdSet(board.Columns.PendingWaiting)["REQ-2"] {
		t.Fatal("REQ-2 depends on a cancelled REQ and must stay waiting, not ready")
	}
}

// TestCompletedWithIssuesSatisfiesGating covers the other half of the
// terminal-success pair — it must unblock dependents exactly like `completed`.
func TestCompletedWithIssuesSatisfiesGating(t *testing.T) {
	board := buildDependencyBoard(t, []dependencyFixture{
		{RequestId: "REQ-1", Status: "completed-with-issues"},
		{RequestId: "REQ-2", Status: "pending", DependsOnIds: []string{"REQ-1"}},
	})

	if !requestIdSet(board.Columns.PendingReady)["REQ-2"] {
		t.Fatal("REQ-2 depends on a completed-with-issues REQ and must be ready")
	}
}

// TestDanglingDependencyWarnsAndBlocks asserts a depends_on pointing at nothing
// fails closed (waiting, never ready) AND surfaces a warning — the pointer can
// never self-resolve, so silence would strand the REQ invisibly.
func TestDanglingDependencyWarnsAndBlocks(t *testing.T) {
	board := buildDependencyBoard(t, []dependencyFixture{
		{RequestId: "REQ-2", Status: "pending", DependsOnIds: []string{"REQ-9999"}},
	})

	if !requestIdSet(board.Columns.PendingWaiting)["REQ-2"] {
		t.Fatal("a REQ with a dangling dependency must be waiting, not ready")
	}
	foundWarning := false
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "REQ-9999") && strings.Contains(warningText, "not in the do-work tree") {
			foundWarning = true
		}
	}
	if !foundWarning {
		t.Fatalf("no dangling-dependency warning; warnings=%v", board.Warnings)
	}
}

// TestDependentsRecordTheReverseEdge asserts the "blocked on me" view: the ids of
// every REQ whose depends_on names this one, in id order and de-duplicated.
func TestDependentsRecordTheReverseEdge(t *testing.T) {
	board := buildDependencyBoard(t, []dependencyFixture{
		{RequestId: "REQ-1", Status: "pending"},
		{RequestId: "REQ-2", Status: "pending", DependsOnIds: []string{"REQ-1", "REQ-1"}}, // repeat must not double-count
		{RequestId: "REQ-3", Status: "pending", DependsOnIds: []string{"REQ-1"}},
	})

	dependents := board.RequestsById["REQ-1"].Dependents
	if len(dependents) != 2 || dependents[0] != "REQ-2" || dependents[1] != "REQ-3" {
		t.Fatalf("REQ-1.Dependents = %v, want [REQ-2 REQ-3]", dependents)
	}
	if len(board.RequestsById["REQ-2"].UnmetDependencies) != 1 {
		t.Fatalf("a repeated depends_on entry must collapse to one unmet dependency: %v",
			board.RequestsById["REQ-2"].UnmetDependencies)
	}
}

// TestDependencyStateReachesGeneratedData asserts the annotations survive the
// projection into board-data.js, where board.js styles chips and badges from them.
func TestDependencyStateReachesGeneratedData(t *testing.T) {
	board := buildDependencyBoard(t, []dependencyFixture{
		{RequestId: "REQ-1", Status: "pending"},
		{RequestId: "REQ-2", Status: "pending", DependsOnIds: []string{"REQ-1"}},
	})
	boardData, projectError := buildGeneratedBoardData(board)
	if projectError != nil {
		t.Fatalf("buildGeneratedBoardData: %v", projectError)
	}

	if got := boardData.Requests["REQ-2"].UnmetDependencies; len(got) != 1 || got[0] != "REQ-1" {
		t.Fatalf("REQ-2 unmetDependencies = %v, want [REQ-1]", got)
	}
	if got := boardData.Requests["REQ-1"].Dependents; len(got) != 1 || got[0] != "REQ-2" {
		t.Fatalf("REQ-1 dependents = %v, want [REQ-2]", got)
	}
	if len(boardData.Columns.PendingReady) != 1 || boardData.Columns.PendingReady[0] != "REQ-1" {
		t.Fatalf("pendingReady = %v, want [REQ-1]", boardData.Columns.PendingReady)
	}
	if len(boardData.Columns.PendingWaiting) != 1 || boardData.Columns.PendingWaiting[0] != "REQ-2" {
		t.Fatalf("pendingWaiting = %v, want [REQ-2]", boardData.Columns.PendingWaiting)
	}
}
