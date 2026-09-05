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

// TestBuildActivityRowsEmitsOneRowPerLifecycleStamp is REQ-572's captured
// RED/GREEN case: REQ-570 was captured at 22:52 and claimed eight minutes
// later, and the Activity view showed only the claim. Every parseable stamp is
// now its own row, so the whole path a REQ took is readable on one surface.
func TestBuildActivityRowsEmitsOneRowPerLifecycleStamp(t *testing.T) {
	rows := buildActivityRows([]*RequestTicket{{
		RequestId: "REQ-570",
		Status:    "claimed",
		CreatedAt: "2026-09-04T22:52:00Z",
		ClaimedAt: "2026-09-04T23:00:17Z",
	}})

	want := []ActivityRow{
		{RequestId: "REQ-570", StampField: "claimed_at", StampTime: time.Date(2026, 9, 4, 23, 0, 17, 0, time.UTC), Transition: "claimed"},
		{RequestId: "REQ-570", StampField: "created_at", StampTime: time.Date(2026, 9, 4, 22, 52, 0, 0, time.UTC), Transition: "captured"},
	}
	if len(rows) != len(want) {
		t.Fatalf("rows = %d, want %d (one per parseable stamp): %+v", len(rows), len(want), rows)
	}
	for index, wantRow := range want {
		got := rows[index]
		if got.StampField != wantRow.StampField || got.Transition != wantRow.Transition || !got.StampTime.Equal(wantRow.StampTime) {
			t.Fatalf("row %d = %+v, want %+v", index, got, wantRow)
		}
	}
}

// TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition is REQ-568's
// captured GREEN condition in Go, widened by REQ-572: a REQ whose status flip
// lands inside the window is on the surface, carries the transition its stamp
// records, and sorts by that stamp rather than by status or by queue order —
// and its earlier claim now sits on the same surface too.
func TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition(t *testing.T) {
	// The instants are the ones REQ-568's Why section cites for the 16:56 board
	// that prompted it, not invented offsets: REQ-504's status flipped at 16:38,
	// REQ-505 was the only claimed card and was 20 minutes old, REQ-485 was the
	// newest Recently-done card at two hours old, and REQ-567 and REQ-503 carry
	// the flip stamps their own frontmatter records. (The status those three
	// flipped TO was retired from the schema by REQ-570; `blocked` stands in for
	// it, because what this test measures is the stamp, not the value.)
	now := time.Date(2026, 9, 4, 16, 56, 0, 0, time.UTC)
	at := func(hour, minute int) string {
		return time.Date(2026, 9, 4, hour, minute, 0, 0, time.UTC).Format(time.RFC3339)
	}
	tickets := []*RequestTicket{
		// The claimed REQ the board already showed.
		{RequestId: "REQ-505", Status: "claimed", ClaimedAt: at(16, 36)},
		// The three REQs the board showed nowhere. Each sorts by its flip stamp,
		// not by its older claim.
		{RequestId: "REQ-504", Status: "blocked", ClaimedAt: at(15, 40), StatusChangedAt: at(16, 38)},
		{RequestId: "REQ-503", Status: "blocked", ClaimedAt: at(14, 57), StatusChangedAt: at(15, 4)},
		{RequestId: "REQ-567", Status: "blocked", ClaimedAt: at(15, 8), StatusChangedAt: at(15, 19)},
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

	// Newest first by stamp, and every stamp is its own row. The capture listed
	// the REQs as 505, 504, 503, 567, 485, but the stamps it also cites put
	// REQ-504's 16:38 flip above REQ-505's 16:36 claim and REQ-567's 15:19 flip
	// above REQ-503's 15:04. The stamps decide the order, not the capture's
	// approximate listing — and the claims that preceded each flip interleave
	// with the other REQs by time rather than hiding behind the newer stamp.
	wantOrder := []ActivityRow{
		{RequestId: "REQ-504", StampField: "status_changed_at", Transition: "status changed to blocked"},
		{RequestId: "REQ-505", StampField: "claimed_at", Transition: "claimed"},
		{RequestId: "REQ-504", StampField: "claimed_at", Transition: "claimed"},
		{RequestId: "REQ-567", StampField: "status_changed_at", Transition: "status changed to blocked"},
		{RequestId: "REQ-567", StampField: "claimed_at", Transition: "claimed"},
		{RequestId: "REQ-503", StampField: "status_changed_at", Transition: "status changed to blocked"},
		{RequestId: "REQ-503", StampField: "claimed_at", Transition: "claimed"},
		{RequestId: "REQ-485", StampField: "completed_at", Transition: "completed"},
	}
	if len(inWindow) != len(wantOrder) {
		t.Fatalf("in-window rows = %d, want %d: %+v", len(inWindow), len(wantOrder), inWindow)
	}
	// Status must not filter the surface, and each row must name the transition
	// its stamp records — not merely carry a time.
	for index, wantRow := range wantOrder {
		got := inWindow[index]
		if got.RequestId != wantRow.RequestId || got.StampField != wantRow.StampField || got.Transition != wantRow.Transition {
			t.Fatalf("row %d = %s/%s/%q, want %s/%s/%q (full order %+v)",
				index, got.RequestId, got.StampField, got.Transition,
				wantRow.RequestId, wantRow.StampField, wantRow.Transition, inWindow)
		}
	}
}

// TestBuildActivityRowsOrdersOneTicketsStampsNewestFirst pins that one REQ's
// own rows are ordered by their instants, not by the order the fields are
// declared in. created_at is first in the lifecycle list and every ticket has
// one, so a reader that emitted rows in declaration order would put the oldest
// transition at the top of the surface and still pass a count assertion.
func TestBuildActivityRowsOrdersOneTicketsStampsNewestFirst(t *testing.T) {
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
	wantSequence := []struct {
		stampField string
		hoursAgo   int
	}{
		{"review_at", 1},
		{"planning_at", 8},
		{"claimed_at", 9},
		{"created_at", 72},
	}
	if len(rows) != len(wantSequence) {
		t.Fatalf("rows = %d, want %d (one per parseable stamp): %+v", len(rows), len(wantSequence), rows)
	}
	for index, wantRow := range wantSequence {
		wantInstant := now.Add(-time.Duration(wantRow.hoursAgo) * time.Hour)
		if rows[index].StampField != wantRow.stampField || !rows[index].StampTime.Equal(wantInstant) {
			t.Fatalf("row %d = %s at %s, want %s at %s", index, rows[index].StampField, rows[index].StampTime, wantRow.stampField, wantInstant)
		}
	}
}

// TestBuildActivityRowsStraddlesTheWindowBoundary derives both samples from
// activityWindow instead of picking two comfortably-separated times, so a
// second, wrong cutoff cannot classify the pair the same way and pass.
func TestBuildActivityRowsStraddlesTheWindowBoundary(t *testing.T) {
	now := time.Date(2026, 9, 4, 18, 0, 0, 0, time.UTC)
	inside := &RequestTicket{
		RequestId:       "REQ-902",
		Status:          "blocked",
		StatusChangedAt: now.Add(-activityWindow + time.Minute).Format(time.RFC3339),
	}
	outside := &RequestTicket{
		RequestId:       "REQ-903",
		Status:          "blocked",
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

	// One REQ carrying two stamps at one instant is the common case, not an
	// exotic one: the work loop writes completed_at and status_changed_at
	// together. The REQ id cannot separate those two rows, so the stamp field is
	// the third key. The blocked_at/claimed_at pair is the one that bites: the
	// field order decides it against the declaration order in
	// lifecycleTimestampFields, so an implementation that leans on the emit
	// order or on the sort's stability puts claimed_at first and fails here.
	sameTicketRows := buildActivityRows([]*RequestTicket{{
		RequestId:       "REQ-909",
		Status:          "completed",
		BlockedAt:       "2026-09-04T09:00:00Z",
		ClaimedAt:       "2026-09-04T09:00:00Z",
		CompletedAt:     "2026-09-04T12:00:00Z",
		StatusChangedAt: "2026-09-04T12:00:00Z",
	}})
	wantStampOrder := []string{"completed_at", "status_changed_at", "blocked_at", "claimed_at"}
	if len(sameTicketRows) != len(wantStampOrder) {
		t.Fatalf("rows = %d, want %d: %+v", len(sameTicketRows), len(wantStampOrder), sameTicketRows)
	}
	for index, wantField := range wantStampOrder {
		if sameTicketRows[index].StampField != wantField {
			t.Fatalf("same-instant row %d = %q, want %q (stamp field must break a same-REQ tie): %+v",
				index, sameTicketRows[index].StampField, wantField, sameTicketRows)
		}
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

	// And the aggregation reads the same list: every declared stamp becomes its
	// own row, and no row cites a field the list does not declare. A count alone
	// would pass on a reader that emitted one field fourteen times, so the two
	// sets are compared rather than their sizes.
	rows := buildActivityRows([]*RequestTicket{everyStamp})
	if len(rows) != len(declared) {
		t.Fatalf("rows = %d, want %d — one per declared lifecycle stamp: %+v", len(rows), len(declared), rows)
	}
	rowFields := map[string]int{}
	for _, row := range rows {
		rowFields[row.StampField]++
	}
	for _, field := range declared {
		if rowFields[field.FieldName] != 1 {
			t.Fatalf("declared lifecycle stamp %q produced %d activity rows, want exactly 1: %+v",
				field.FieldName, rowFields[field.FieldName], rows)
		}
	}
}
