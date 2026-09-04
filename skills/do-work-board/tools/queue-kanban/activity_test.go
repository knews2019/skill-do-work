package main

import (
	"strings"
	"testing"
	"time"
)

// activityWindow is the window the captured GREEN condition names ("last 24h").
// Both boundary cases below are derived from it rather than restated beside it:
// a fixture that spans a threshold widely does not test the threshold, because
// any cutoff in the gap classifies every sample the same way.
const activityWindow = 24 * time.Hour

// TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition is REQ-568's
// captured GREEN condition in Go: a pending-heavy-testing REQ whose hold lands
// inside the window is on the surface, carries the transition its stamp
// records, and sorts by that stamp rather than by status or by queue order.
func TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition(t *testing.T) {
	// The instants are the ones REQ-568's Why section cites for the 16:56 board
	// that prompted it, not invented offsets: REQ-504 was held at 16:38, REQ-505
	// was the only claimed card and was 20 minutes old, REQ-485 was the newest
	// Recently-done card at two hours old, and REQ-567 and REQ-503 carry the
	// hold stamps their own frontmatter records.
	now := time.Date(2026, 9, 4, 16, 56, 0, 0, time.UTC)
	at := func(hour, minute int) string {
		return time.Date(2026, 9, 4, hour, minute, 0, 0, time.UTC).Format(time.RFC3339)
	}
	tickets := []*RequestTicket{
		// The claimed REQ the board already showed.
		{RequestId: "REQ-505", Status: "claimed", ClaimedAt: at(16, 36)},
		// The three held REQs the board showed nowhere. Each sorts by its hold
		// stamp, not by its older claim.
		{RequestId: "REQ-504", Status: "pending-heavy-testing", ClaimedAt: at(15, 40), StatusChangedAt: at(16, 38)},
		{RequestId: "REQ-503", Status: "pending-heavy-testing", ClaimedAt: at(14, 57), StatusChangedAt: at(15, 4)},
		{RequestId: "REQ-567", Status: "pending-heavy-testing", ClaimedAt: at(15, 8), StatusChangedAt: at(15, 19)},
		// The terminal REQ that was the newest thing "Recently done" could show.
		{RequestId: "REQ-485", Status: "completed", CompletedAt: at(14, 56), CompletionTime: time.Date(2026, 9, 4, 14, 56, 0, 0, time.UTC)},
	}

	rows := buildActivityRows(tickets)

	var inWindow []ActivityRow
	for _, row := range rows {
		if isWithinRecentWindow(row.StampTime, now, activityWindow) {
			inWindow = append(inWindow, row)
		}
	}

	// Newest first by stamp. This is the captured GREEN's five REQs; the capture
	// listed them as 505, 504, 503, 567, 485, but the stamps it also cites put
	// REQ-504's 16:38 hold above REQ-505's 16:36 claim and REQ-567's 15:19 hold
	// above REQ-503's 15:04. The stamps decide the order, not the capture's
	// approximate listing.
	wantOrder := []string{"REQ-504", "REQ-505", "REQ-567", "REQ-503", "REQ-485"}
	if len(inWindow) != len(wantOrder) {
		t.Fatalf("in-window rows = %d, want %d: %+v", len(inWindow), len(wantOrder), inWindow)
	}
	for index, wantId := range wantOrder {
		if inWindow[index].RequestId != wantId {
			t.Fatalf("row %d = %s, want %s (full order %+v)", index, inWindow[index].RequestId, wantId, inWindow)
		}
	}

	// Status must not filter the surface, and each row must name the transition
	// its stamp records — not merely carry a time.
	byId := map[string]ActivityRow{}
	for _, row := range inWindow {
		byId[row.RequestId] = row
	}
	for _, heldId := range []string{"REQ-504", "REQ-503", "REQ-567"} {
		held := byId[heldId]
		if held.StampField != "status_changed_at" {
			t.Fatalf("%s newest stamp field = %q, want status_changed_at", heldId, held.StampField)
		}
		if held.Transition != "held for heavy testing" {
			t.Fatalf("%s transition = %q, want %q", heldId, held.Transition, "held for heavy testing")
		}
	}
	if byId["REQ-505"].Transition != "claimed" {
		t.Fatalf("REQ-505 transition = %q, want %q", byId["REQ-505"].Transition, "claimed")
	}
	if byId["REQ-485"].Transition != "completed" {
		t.Fatalf("REQ-485 transition = %q, want %q", byId["REQ-485"].Transition, "completed")
	}
}

// TestBuildActivityRowsPicksTheNewestStampNotTheFirstDeclaredField pins that
// the rule is "newest stamp", not "first field in the table". created_at is
// first in the lifecycle list and every ticket has one, so a reader that
// stopped at the first present field would pass every other assertion here.
func TestBuildActivityRowsPicksTheNewestStampNotTheFirstDeclaredField(t *testing.T) {
	now := time.Date(2026, 9, 4, 18, 0, 0, 0, time.UTC)
	ticket := &RequestTicket{
		RequestId:   "REQ-901",
		Status:      "claimed",
		CreatedAt:   now.Add(-72 * time.Hour).Format(time.RFC3339),
		ClaimedAt:   now.Add(-9 * time.Hour).Format(time.RFC3339),
		ReviewAt:    now.Add(-1 * time.Hour).Format(time.RFC3339),
		PlanningAt:  now.Add(-8 * time.Hour).Format(time.RFC3339),
		BlockedAt:   "",
		CompletedAt: "",
	}
	rows := buildActivityRows([]*RequestTicket{ticket})
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	if rows[0].StampField != "review_at" {
		t.Fatalf("newest stamp field = %q, want review_at", rows[0].StampField)
	}
	if !rows[0].StampTime.Equal(now.Add(-1 * time.Hour)) {
		t.Fatalf("newest stamp time = %s, want %s", rows[0].StampTime, now.Add(-1*time.Hour))
	}
}

// TestBuildActivityRowsStraddlesTheWindowBoundary derives both samples from
// activityWindow instead of picking two comfortably-separated times, so a
// second, wrong cutoff cannot classify the pair the same way and pass.
func TestBuildActivityRowsStraddlesTheWindowBoundary(t *testing.T) {
	now := time.Date(2026, 9, 4, 18, 0, 0, 0, time.UTC)
	inside := &RequestTicket{
		RequestId:       "REQ-902",
		Status:          "pending-heavy-testing",
		StatusChangedAt: now.Add(-activityWindow + time.Minute).Format(time.RFC3339),
	}
	outside := &RequestTicket{
		RequestId:       "REQ-903",
		Status:          "pending-heavy-testing",
		StatusChangedAt: now.Add(-activityWindow - time.Minute).Format(time.RFC3339),
	}
	rows := buildActivityRows([]*RequestTicket{inside, outside})
	if len(rows) != 2 {
		t.Fatalf("aggregation must return both rows and leave windowing to its caller, got %d", len(rows))
	}
	if !isWithinRecentWindow(rows[0].StampTime, now, activityWindow) {
		t.Fatalf("REQ-902 one minute inside the window was excluded")
	}
	if isWithinRecentWindow(rows[1].StampTime, now, activityWindow) {
		t.Fatalf("REQ-903 one minute outside the window was included")
	}
}

// TestBuildActivityRowsSkipsTicketsWithNoParseableStamp keeps an unparseable or
// stampless ticket off the surface rather than dating it from zero time, which
// would sort it last forever and read as real evidence of no activity.
func TestBuildActivityRowsSkipsTicketsWithNoParseableStamp(t *testing.T) {
	rows := buildActivityRows([]*RequestTicket{
		{RequestId: "REQ-904", Status: "pending"},
		{RequestId: "REQ-905", Status: "pending", CreatedAt: "not-a-timestamp"},
		{RequestId: "REQ-906", Status: "pending", CreatedAt: "2026-09-04T10:00:00Z"},
	})
	if len(rows) != 1 || rows[0].RequestId != "REQ-906" {
		t.Fatalf("rows = %+v, want only REQ-906", rows)
	}
	if rows[0].Transition != "captured" {
		t.Fatalf("created_at transition = %q, want %q", rows[0].Transition, "captured")
	}
}

// TestBuildActivityRowsBreaksStampTiesDeterministically pins the tiebreak in
// the comparator rather than leaving it to the sort's stability, which would
// let the input order silently decide and would change under a reordered walk.
func TestBuildActivityRowsBreaksStampTiesDeterministically(t *testing.T) {
	sameInstant := "2026-09-04T12:00:00Z"
	forward := buildActivityRows([]*RequestTicket{
		{RequestId: "REQ-907", Status: "claimed", ClaimedAt: sameInstant},
		{RequestId: "REQ-908", Status: "claimed", ClaimedAt: sameInstant},
	})
	reversed := buildActivityRows([]*RequestTicket{
		{RequestId: "REQ-908", Status: "claimed", ClaimedAt: sameInstant},
		{RequestId: "REQ-907", Status: "claimed", ClaimedAt: sameInstant},
	})
	if len(forward) != 2 || len(reversed) != 2 {
		t.Fatalf("expected two rows from each ordering")
	}
	if forward[0].RequestId != reversed[0].RequestId || forward[1].RequestId != reversed[1].RequestId {
		t.Fatalf("tie order depends on input order: %s,%s vs %s,%s",
			forward[0].RequestId, forward[1].RequestId, reversed[0].RequestId, reversed[1].RequestId)
	}
}

// TestLifecycleTimestampFieldsIsTheOneListBothReadersUse is the anti-drift pin.
// detectFutureTimestampFields and the activity aggregation must read the same
// enumeration, so adding a stamp to the schema cannot light up one surface and
// leave the other silently stale.
func TestLifecycleTimestampFieldsIsTheOneListBothReadersUse(t *testing.T) {
	future := time.Now().UTC().Add(72 * time.Hour).Format(time.RFC3339)
	for _, field := range lifecycleTimestampFields(&RequestTicket{}) {
		if field.Transition == "" {
			t.Fatalf("%s carries no transition name; every lifecycle stamp must say what it records", field.FieldName)
		}
	}
	// This literal is the tripwire, and it is deliberately a second spelling of
	// the list: adding a stamp to lifecycleTimestampFields without adding it
	// here fails the count below, which is the prompt to decide what the new
	// stamp's transition should say. Every field is set to the same future
	// instant, so detectFutureTimestampFields must flag all of them — proving it
	// reads the whole shared list rather than a subset of it.
	everyStamp := &RequestTicket{
		RequestId:         "REQ-909",
		CreatedAt:         future,
		ClaimedAt:         future,
		CompletedAt:       future,
		PlanningAt:        future,
		DispatchAt:        future,
		BuilderHandbackAt: future,
		IntegrationAt:     future,
		ReviewAt:          future,
		RemediationAt:     future,
		ReReviewAt:        future,
		ReleaseAt:         future,
		StatusChangedAt:   future,
		BlockedAt:         future,
		TestingUpdatedAt:  future,
	}
	declared := lifecycleTimestampFields(everyStamp)
	flagged := detectFutureTimestampFields(everyStamp, time.Now().UTC())
	if len(flagged) != len(declared) {
		t.Fatalf("future-stamp reader flagged %d of %d declared lifecycle stamps — the readers have drifted, or this test's ticket literal is missing a newly declared stamp: %v",
			len(flagged), len(declared), flagged)
	}
	for _, field := range declared {
		found := false
		for _, entry := range flagged {
			if strings.HasPrefix(entry, field.FieldName+" ") {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("declared lifecycle stamp %q never reached the future-stamp reader: %v", field.FieldName, flagged)
		}
	}

	// And the aggregation reads the same list: with every stamp at one instant,
	// the newest-stamp pick must land on a field the list declares.
	rows := buildActivityRows([]*RequestTicket{everyStamp})
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1", len(rows))
	}
	declaredField := false
	for _, field := range declared {
		if field.FieldName == rows[0].StampField {
			declaredField = true
			break
		}
	}
	if !declaredField {
		t.Fatalf("activity row cited %q, which lifecycleTimestampFields does not declare", rows[0].StampField)
	}
}
