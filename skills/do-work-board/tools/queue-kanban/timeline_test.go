package main

import (
	"math"
	"testing"
	"time"
)

func timelineTicket(requestId string, status string, createdAt string, claimedAt string, completedAt string) *RequestTicket {
	return &RequestTicket{
		RequestId:   requestId,
		Title:       "fixture " + requestId,
		Status:      status,
		CreatedAt:   createdAt,
		ClaimedAt:   claimedAt,
		CompletedAt: completedAt,
	}
}

func findTimelineRow(t *testing.T, aggregate TimelineAggregate, requestId string) TimelineRow {
	t.Helper()
	for _, row := range aggregate.Rows {
		if row.RequestId == requestId {
			return row
		}
	}
	t.Fatalf("no timeline row for %s", requestId)
	return TimelineRow{}
}

// The four bar shapes the view has to draw. Before this aggregation the board's
// only duration reader was buildDurationAggregate, which skips every ticket that
// is not terminal-success and measures only completed_at − claimed_at — so the
// wait span existed nowhere, and an in-flight REQ contributed to no view at all.
func TestTimelineAggregateProducesTheFourBarShapes(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	finished := timelineTicket("REQ-401", "completed",
		"2026-06-01T09:00:00Z", "2026-06-03T09:00:00Z", "2026-06-03T10:30:00Z")
	inFlight := timelineTicket("REQ-402", "claimed",
		"2026-06-05T09:00:00Z", "2026-06-09T09:00:00Z", "")
	waitingOnly := timelineTicket("REQ-403", "pending",
		"2026-06-08T09:00:00Z", "", "")
	reversed := timelineTicket("REQ-404", "completed",
		"2026-06-02T09:00:00Z", "2026-06-04T09:00:00Z", "2026-06-04T08:00:00Z")
	reversed.CompletionAnomaly = true
	reversed.CompletionAnomalyReason = "completed_at is earlier than claimed_at"

	aggregate := buildTimelineAggregate(
		[]*RequestTicket{finished, inFlight, waitingOnly, reversed}, now)

	if len(aggregate.Rows) != 4 {
		t.Fatalf("aggregate carried %d rows, want 4", len(aggregate.Rows))
	}

	finishedRow := findTimelineRow(t, aggregate, "REQ-401")
	assertTimelineMinutes(t, "REQ-401 wait", finishedRow.WaitMinutes, 2*24*60)
	if finishedRow.WaitOpen {
		t.Fatal("REQ-401 was claimed, so its wait is closed")
	}
	if !finishedRow.HasWork {
		t.Fatal("REQ-401 was claimed, so it has a work segment")
	}
	assertTimelineMinutes(t, "REQ-401 work", finishedRow.WorkMinutes, 90)
	if finishedRow.WorkOpen {
		t.Fatal("REQ-401 completed, so its work segment is closed")
	}

	inFlightRow := findTimelineRow(t, aggregate, "REQ-402")
	assertTimelineMinutes(t, "REQ-402 wait", inFlightRow.WaitMinutes, 4*24*60)
	if !inFlightRow.HasWork || !inFlightRow.WorkOpen {
		t.Fatal("REQ-402 is claimed and unfinished, so its work segment is open")
	}
	// now − claimed_at, so the open span is measured against the board's one now.
	assertTimelineMinutes(t, "REQ-402 open work", inFlightRow.WorkMinutes, 27*60)

	waitingRow := findTimelineRow(t, aggregate, "REQ-403")
	if !waitingRow.WaitOpen {
		t.Fatal("REQ-403 was never claimed, so its wait is still running")
	}
	assertTimelineMinutes(t, "REQ-403 open wait", waitingRow.WaitMinutes, 2*24*60+3*60)
	if waitingRow.HasWork {
		t.Fatal("REQ-403 was never claimed; inventing a work segment for it is REQ-228's job, not this one's")
	}

	reversedRow := findTimelineRow(t, aggregate, "REQ-404")
	if !reversedRow.Anomaly {
		t.Fatal("REQ-404's reversed stamp must be flagged, not swallowed")
	}
	if reversedRow.AnomalyReason != "completed_at is earlier than claimed_at" {
		t.Fatalf("REQ-404 reason = %q, want the ticket's own verdict verbatim", reversedRow.AnomalyReason)
	}
	if reversedRow.WorkMinutes >= 0 {
		t.Fatalf("REQ-404 work = %.1f min; a reversed span must stay negative rather than be clamped to zero",
			reversedRow.WorkMinutes)
	}
	assertTimelineMinutes(t, "REQ-404 work", reversedRow.WorkMinutes, -60)
}

func assertTimelineMinutes(t *testing.T, label string, got float64, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.001 {
		t.Fatalf("%s = %.4f min, want %.4f", label, got, want)
	}
}

// Row order is a stated property of the view, so it is pinned rather than left
// to sort stability: equal created_at instants break by id, which is what stops
// two REQs captured in the same second swapping places between builds.
func TestTimelineRowsAreChronologicalWithAStableTiebreak(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	aggregate := buildTimelineAggregate([]*RequestTicket{
		timelineTicket("REQ-503", "pending", "2026-06-04T09:00:00Z", "", ""),
		timelineTicket("REQ-501", "pending", "2026-06-02T09:00:00Z", "", ""),
		timelineTicket("REQ-500", "pending", "2026-06-02T09:00:00Z", "", ""),
		timelineTicket("REQ-502", "pending", "2026-06-03T09:00:00Z", "", ""),
	}, now)

	gotOrder := make([]string, 0, len(aggregate.Rows))
	for _, row := range aggregate.Rows {
		gotOrder = append(gotOrder, row.RequestId)
	}
	wantOrder := []string{"REQ-500", "REQ-501", "REQ-502", "REQ-503"}
	for orderIndex := range wantOrder {
		if gotOrder[orderIndex] != wantOrder[orderIndex] {
			t.Fatalf("row order = %v, want %v", gotOrder, wantOrder)
		}
	}
}

// A fitted view has to contain every bar it draws. An open bar runs to now, so
// now belongs to the range whenever any row is open — otherwise the one bar the
// reader most wants to see is the one hanging off the right edge.
func TestTimelineRangeReachesNowWhileAnyBarIsOpen(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	openAggregate := buildTimelineAggregate([]*RequestTicket{
		timelineTicket("REQ-601", "claimed", "2026-06-01T09:00:00Z", "2026-06-02T09:00:00Z", ""),
	}, now)
	if !openAggregate.RangeEnd.Equal(now) {
		t.Fatalf("range end = %s with an open bar, want now (%s)", openAggregate.RangeEnd, now)
	}

	// With every bar closed the range stops at the last real instant, so a board
	// nobody has touched for a month does not draw a month of empty axis.
	closedAggregate := buildTimelineAggregate([]*RequestTicket{
		timelineTicket("REQ-602", "completed", "2026-06-01T09:00:00Z", "2026-06-02T09:00:00Z", "2026-06-02T11:00:00Z"),
	}, now)
	wantEnd := time.Date(2026, 6, 2, 11, 0, 0, 0, time.UTC)
	if !closedAggregate.RangeEnd.Equal(wantEnd) {
		t.Fatalf("range end = %s with every bar closed, want the last completion (%s)",
			closedAggregate.RangeEnd, wantEnd)
	}
	if !closedAggregate.RangeStart.Equal(time.Date(2026, 6, 1, 9, 0, 0, 0, time.UTC)) {
		t.Fatalf("range start = %s, want the earliest created_at", closedAggregate.RangeStart)
	}
}

// created_at is the row's left edge, so a ticket without a parseable one has no
// bar to draw. Dropping it is the honest outcome; placing it at zero or at now
// would invent a wait that was never measured.
func TestTimelineSkipsTicketsWithoutAParseableCreatedAt(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	aggregate := buildTimelineAggregate([]*RequestTicket{
		timelineTicket("REQ-701", "pending", "", "", ""),
		timelineTicket("REQ-702", "pending", "not a timestamp", "", ""),
		timelineTicket("REQ-703", "pending", "2026-06-02T09:00:00Z", "", ""),
		nil,
	}, now)
	if len(aggregate.Rows) != 1 || aggregate.Rows[0].RequestId != "REQ-703" {
		t.Fatalf("rows = %+v, want only REQ-703", aggregate.Rows)
	}
}
