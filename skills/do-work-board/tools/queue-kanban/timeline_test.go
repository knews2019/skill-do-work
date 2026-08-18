package main

import (
	"fmt"
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

// projectionTicket builds one ticket for the projection fixtures.
func projectionTicket(requestId string, status string, effortEstimate string, dependsOn []string) *RequestTicket {
	return &RequestTicket{
		RequestId:      requestId,
		Title:          "fixture " + requestId,
		Status:         status,
		EffortEstimate: effortEstimate,
		DependsOn:      dependsOn,
		CreatedAt:      "2026-06-01T09:00:00Z",
	}
}

// completedSpanTicket is one finished REQ contributing a measured span to the
// medians the projection is taken over.
func completedSpanTicket(requestId string, effortEstimate string, completedAt time.Time, spanMinutes float64) *RequestTicket {
	claimedAt := completedAt.Add(-time.Duration(spanMinutes * float64(time.Minute)))
	ticket := timelineTicket(requestId, "completed",
		"2026-05-01T09:00:00Z", claimedAt.Format(time.RFC3339), completedAt.Format(time.RFC3339))
	ticket.EffortEstimate = effortEstimate
	return ticket
}

// The REQ's own Red-Green Proof, verbatim: a known history with one paused and
// one reversed span that must not reach either median, and four pending REQs of
// which one depends on another and one is blocked. Before this there was no
// projection of any kind — buildDurationAggregate skips every non-terminal
// ticket, so a pending REQ contributed to no computation anywhere.
func TestTimelineProjectionChainsTheQueueSerially(t *testing.T) {
	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	tickets := []*RequestTicket{}

	// Six normal spans with median 40, and six trivial spans with median 10.
	for spanIndex, spanMinutes := range []float64{20, 30, 35, 45, 50, 60} {
		tickets = append(tickets, completedSpanTicket(
			fmt.Sprintf("REQ-1%02d", spanIndex), "normal",
			now.Add(-time.Duration(spanIndex+1)*time.Hour), spanMinutes))
	}
	for spanIndex, spanMinutes := range []float64{5, 8, 9, 11, 12, 15} {
		tickets = append(tickets, completedSpanTicket(
			fmt.Sprintf("REQ-2%02d", spanIndex), "trivial",
			now.Add(-time.Duration(spanIndex+1)*time.Hour), spanMinutes))
	}
	// One paused (over the four-hour ceiling) and one reversed span. Both are
	// normal-bucket and both are enormous, so if either reached the median the
	// normal figure could not stay at 40.
	tickets = append(tickets, completedSpanTicket("REQ-300", "normal", now.Add(-30*time.Minute), 900))
	tickets = append(tickets, completedSpanTicket("REQ-301", "normal", now.Add(-20*time.Minute), -120))

	// Four pending REQs. REQ-402 depends on REQ-403, so it must be scheduled
	// after it even though its id sorts first.
	tickets = append(tickets,
		projectionTicket("REQ-401", "pending", "trivial", nil),
		projectionTicket("REQ-402", "pending", "normal", []string{"REQ-403"}),
		projectionTicket("REQ-403", "pending", "normal", nil),
		projectionTicket("REQ-404", "blocked", "normal", nil),
	)
	resolveUnmetDependenciesForTest(tickets)

	projection := buildTimelineProjection(tickets, buildDurationAggregate(tickets), now)

	if !projection.Confident {
		t.Fatalf("twelve in-rule samples must clear the confidence floor; declined because %q", projection.DeclinedReason)
	}
	if projection.NormalMedianMinutes != 40 {
		t.Fatalf("normal median = %.1f, want 40 — the paused and reversed spans must not reach it",
			projection.NormalMedianMinutes)
	}
	if projection.TrivialMedianMinutes != 10 {
		t.Fatalf("trivial median = %.1f, want 10", projection.TrivialMedianMinutes)
	}

	wantChain := []struct {
		requestId   string
		spanMinutes float64
	}{
		{"REQ-401", 10}, // trivial, no dependencies, lowest id
		{"REQ-403", 40}, // normal, no dependencies
		{"REQ-402", 40}, // normal, but only after REQ-403
	}
	if len(projection.Rows) != len(wantChain) {
		t.Fatalf("chain = %+v, want %d rows", projection.Rows, len(wantChain))
	}
	chainCursor := now
	for chainIndex, want := range wantChain {
		row := projection.Rows[chainIndex]
		if row.RequestId != want.requestId {
			t.Fatalf("chain[%d] = %s, want %s (dependency order, then id order)",
				chainIndex, row.RequestId, want.requestId)
		}
		if !row.StartTime.Equal(chainCursor) {
			t.Fatalf("%s starts %s, want %s — the chain is serial with no gaps or overlap",
				row.RequestId, row.StartTime, chainCursor)
		}
		wantEnd := chainCursor.Add(time.Duration(want.spanMinutes * float64(time.Minute)))
		if !row.EndTime.Equal(wantEnd) {
			t.Fatalf("%s ends %s, want %s (its own bucket's median)", row.RequestId, row.EndTime, wantEnd)
		}
		chainCursor = wantEnd
	}

	if len(projection.Excluded) != 1 || projection.Excluded[0].RequestId != "REQ-404" {
		t.Fatalf("excluded = %+v, want only the blocked REQ-404", projection.Excluded)
	}
	if !projection.QueueEnd.Equal(chainCursor) {
		t.Fatalf("queue end = %s, want the last projected finish %s", projection.QueueEnd, chainCursor)
	}
}

// resolveUnmetDependenciesForTest fills the field buildBoard normally computes,
// so the fixtures exercise the same input the projection consumes in production
// rather than a hand-built shortcut.
func resolveUnmetDependenciesForTest(tickets []*RequestTicket) {
	statusById := map[string]string{}
	for _, ticket := range tickets {
		statusById[ticket.RequestId] = ticket.Status
	}
	for _, ticket := range tickets {
		ticket.UnmetDependencies = nil
		for _, dependencyId := range ticket.DependsOn {
			if !isCompletedStatus(statusById[dependencyId]) {
				ticket.UnmetDependencies = append(ticket.UnmetDependencies, dependencyId)
			}
		}
	}
}

// Two data points are not a forecast. Below the floor the projection declines
// with a reason rather than drawing a confident end date off them.
func TestTimelineProjectionDeclinesOnThinHistory(t *testing.T) {
	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	tickets := []*RequestTicket{
		completedSpanTicket("REQ-100", "normal", now.Add(-2*time.Hour), 30),
		completedSpanTicket("REQ-101", "normal", now.Add(-1*time.Hour), 50),
		projectionTicket("REQ-401", "pending", "normal", nil),
	}
	resolveUnmetDependenciesForTest(tickets)

	projection := buildTimelineProjection(tickets, buildDurationAggregate(tickets), now)
	if projection.Confident {
		t.Fatal("two samples must not produce a confident forecast")
	}
	if len(projection.Rows) != 0 {
		t.Fatalf("declined projection drew %d bars; it must draw none", len(projection.Rows))
	}
	if projection.DeclinedReason == "" {
		t.Fatal("declining silently is the failure; the reason must be stated")
	}
	if !projection.QueueEnd.IsZero() {
		t.Fatalf("declined projection reported a queue end of %s", projection.QueueEnd)
	}
}

// The subtle way this view could lie: a pending REQ whose prerequisite is itself
// unschedulable. It looks ready by status, so folding it into the chain would
// move the queue-end figure earlier than anything can actually deliver.
func TestTimelineProjectionExcludesREQsBlockedBehindAnUnschedulableOne(t *testing.T) {
	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	tickets := []*RequestTicket{}
	for spanIndex, spanMinutes := range []float64{20, 30, 40, 50, 60, 70} {
		tickets = append(tickets, completedSpanTicket(
			fmt.Sprintf("REQ-1%02d", spanIndex), "normal",
			now.Add(-time.Duration(spanIndex+1)*time.Hour), spanMinutes))
	}
	tickets = append(tickets,
		projectionTicket("REQ-401", "pending", "normal", nil),
		projectionTicket("REQ-402", "pending", "normal", []string{"REQ-403"}), // behind a blocked REQ
		projectionTicket("REQ-403", "blocked", "normal", nil),
		projectionTicket("REQ-404", "pending", "normal", []string{"REQ-405"}), // dangling dependency
		projectionTicket("REQ-406", "pending-answers", "normal", nil),
	)
	resolveUnmetDependenciesForTest(tickets)

	projection := buildTimelineProjection(tickets, buildDurationAggregate(tickets), now)

	if len(projection.Rows) != 1 || projection.Rows[0].RequestId != "REQ-401" {
		t.Fatalf("chain = %+v, want only REQ-401 — nothing else can be given an honest start time",
			projection.Rows)
	}
	excludedById := map[string]string{}
	for _, exclusion := range projection.Excluded {
		excludedById[exclusion.RequestId] = exclusion.Reason
	}
	for _, wantExcluded := range []string{"REQ-402", "REQ-403", "REQ-404", "REQ-406"} {
		if excludedById[wantExcluded] == "" {
			t.Fatalf("%s is not listed as excluded; every unschedulable REQ must be named with a reason (got %+v)",
				wantExcluded, projection.Excluded)
		}
	}
	if !projection.QueueEnd.Equal(projection.Rows[0].EndTime) {
		t.Fatalf("queue end = %s, want the only scheduled REQ's finish", projection.QueueEnd)
	}
}

// The chain starts after work already in flight, not at now, or the forecast
// claims the queue starts draining while a builder is still mid-REQ.
func TestTimelineProjectionStartsAfterWorkAlreadyInFlight(t *testing.T) {
	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	tickets := []*RequestTicket{}
	for spanIndex, spanMinutes := range []float64{20, 30, 40, 50, 60, 70} {
		tickets = append(tickets, completedSpanTicket(
			fmt.Sprintf("REQ-1%02d", spanIndex), "normal",
			now.Add(-time.Duration(spanIndex+1)*time.Hour), spanMinutes))
	}
	// Claimed 10 minutes ago; the normal median is 45, so it has 35 left.
	inFlight := timelineTicket("REQ-300", "claimed",
		"2026-06-01T09:00:00Z", now.Add(-10*time.Minute).Format(time.RFC3339), "")
	inFlight.EffortEstimate = "normal"
	tickets = append(tickets, inFlight, projectionTicket("REQ-401", "pending", "normal", nil))
	resolveUnmetDependenciesForTest(tickets)

	projection := buildTimelineProjection(tickets, buildDurationAggregate(tickets), now)
	wantStart := now.Add(35 * time.Minute)
	if !projection.ChainStart.Equal(wantStart) {
		t.Fatalf("chain starts %s, want %s — the in-flight REQ's remaining median",
			projection.ChainStart, wantStart)
	}
	if len(projection.Rows) != 1 || !projection.Rows[0].StartTime.Equal(wantStart) {
		t.Fatalf("first projected REQ starts %+v, want %s", projection.Rows, wantStart)
	}
}

// A pending REQ whose prerequisite is CLAIMED is schedulable, not stuck: the
// chain already starts after in-flight work finishes, so by the time the chain
// runs that dependency has resolved. Treating it as unschedulable dropped the
// REQ from the queue-end figure and reported a reason — "unschedulable, missing,
// or circular" — that was actively wrong about a REQ someone was working on.
func TestTimelineProjectionSchedulesDependentsOfClaimedWork(t *testing.T) {
	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	tickets := []*RequestTicket{}
	for spanIndex, spanMinutes := range []float64{20, 30, 40, 50, 60, 70} {
		tickets = append(tickets, completedSpanTicket(
			fmt.Sprintf("REQ-1%02d", spanIndex), "normal",
			now.Add(-time.Duration(spanIndex+1)*time.Hour), spanMinutes))
	}
	// Claimed 10 minutes ago against a normal median of 45, so it has 35 left.
	inFlight := timelineTicket("REQ-300", "claimed",
		"2026-06-01T09:00:00Z", now.Add(-10*time.Minute).Format(time.RFC3339), "")
	inFlight.EffortEstimate = "normal"
	tickets = append(tickets,
		inFlight,
		projectionTicket("REQ-401", "pending", "normal", []string{"REQ-300"}),
		projectionTicket("REQ-402", "pending", "normal", []string{"REQ-401"}),
	)
	resolveUnmetDependenciesForTest(tickets)

	projection := buildTimelineProjection(tickets, buildDurationAggregate(tickets), now)

	if len(projection.Rows) != 2 {
		t.Fatalf("chain = %+v, want both REQs scheduled behind the in-flight one", projection.Rows)
	}
	if projection.Rows[0].RequestId != "REQ-401" || projection.Rows[1].RequestId != "REQ-402" {
		t.Fatalf("chain order = %s, %s; want REQ-401 then REQ-402",
			projection.Rows[0].RequestId, projection.Rows[1].RequestId)
	}
	wantStart := now.Add(35 * time.Minute)
	if !projection.Rows[0].StartTime.Equal(wantStart) {
		t.Fatalf("REQ-401 starts %s, want %s — after the in-flight REQ's remaining median",
			projection.Rows[0].StartTime, wantStart)
	}
	if len(projection.Excluded) != 0 {
		t.Fatalf("excluded = %+v, want none — a claimed prerequisite is in flight, not unschedulable",
			projection.Excluded)
	}
	if projection.QueueEnd.IsZero() {
		t.Fatal("queue end went absent; the estimate must cover both scheduled REQs")
	}

	// A pending REQ behind a genuinely unschedulable one still drops out, so the
	// fix widens what counts as resolvable without weakening the exclusion.
	blockedTickets := append([]*RequestTicket{}, tickets...)
	blockedTickets = append(blockedTickets,
		projectionTicket("REQ-403", "blocked", "normal", nil),
		projectionTicket("REQ-404", "pending", "normal", []string{"REQ-403"}),
	)
	resolveUnmetDependenciesForTest(blockedTickets)
	blockedProjection := buildTimelineProjection(blockedTickets, buildDurationAggregate(blockedTickets), now)
	excludedIds := map[string]bool{}
	for _, exclusion := range blockedProjection.Excluded {
		excludedIds[exclusion.RequestId] = true
	}
	if !excludedIds["REQ-404"] || !excludedIds["REQ-403"] {
		t.Fatalf("excluded = %+v, want REQ-403 and REQ-404 still held out",
			blockedProjection.Excluded)
	}
}

// The chain must run in the order the work action would actually claim, and that
// order is numeric. Lexically REQ-1000 sorts before REQ-999, so past four digits
// a lexical chain predicts the wrong next REQ — the one thing REQ-228's own
// lesson says the forecast must never get wrong.
func TestTimelineProjectionChainOrderIsNumericNotLexical(t *testing.T) {
	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	tickets := []*RequestTicket{}
	for spanIndex, spanMinutes := range []float64{20, 30, 40, 50, 60, 70} {
		tickets = append(tickets, completedSpanTicket(
			fmt.Sprintf("REQ-1%02d", spanIndex), "normal",
			now.Add(-time.Duration(spanIndex+1)*time.Hour), spanMinutes))
	}
	tickets = append(tickets,
		projectionTicket("REQ-1000", "pending", "normal", nil),
		projectionTicket("REQ-999", "pending", "normal", nil),
	)
	resolveUnmetDependenciesForTest(tickets)

	projection := buildTimelineProjection(tickets, buildDurationAggregate(tickets), now)
	if len(projection.Rows) != 2 {
		t.Fatalf("chain = %+v, want both REQs", projection.Rows)
	}
	if projection.Rows[0].RequestId != "REQ-999" {
		t.Fatalf("chain starts with %s; work.md takes the lowest NUMERIC id, so REQ-999 precedes REQ-1000",
			projection.Rows[0].RequestId)
	}
}
