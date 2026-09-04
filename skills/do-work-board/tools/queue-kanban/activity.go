package main

import (
	"sort"
	"time"
)

// Activity aggregation for the board's Activity view.
//
// One row per REQ: its newest lifecycle stamp, the field that stamp lives in,
// and the transition it records. Newest first, and STATUS DOES NOT FILTER IT —
// a REQ that was claimed, held for heavy testing, blocked, completed, cancelled
// or failed inside the window all belong on the same surface.
//
// The question this answers is "what changed on the queue in the last N hours,
// and why", which no other surface answers. Recently done is terminal states
// only. The Timeline draws wait and work spans but not the transitions between
// them. The Calendar dates a REQ from its claim or resolve day. When three REQs
// were claimed, built, merged and held between 14:57 and 16:39 UTC, every one of
// those surfaces showed a two-hour hole and the answer needed `git log`.
//
// Everything here comes from stamps the board already parsed onto each ticket:
// no new frontmatter field, no second walk of the tree, the same posture
// timeline.go and durations.go take.
//
// THE WINDOW IS NOT APPLIED HERE. Every row ships and the client filters
// against the wall clock at render time, the same way recentlyDoneIds does in
// web/board-cards.js — a tab left open past the snapshot must not keep counting
// "last 24 hours" from page-generation time. Callers window with
// isWithinRecentWindow (model.go) so there is one definition of "inside the
// window" rather than a second one here.

// ActivityRow is one REQ's most recent lifecycle transition.
type ActivityRow struct {
	RequestId string

	// StampField is the frontmatter field StampTime was read from. It ships to
	// the client so a reader can go straight to the line that produced the row,
	// which is what makes a missing stamp diagnosable rather than merely absent.
	StampField string
	StampTime  time.Time

	// Transition is what happened, in words — "held for heavy testing", not
	// "status_changed_at". Decided in Go by lifecycleTimestampFields so the
	// client never becomes a second definition of what a stamp means.
	Transition string
}

// buildActivityRows returns one row per ticket carrying at least one parseable
// lifecycle stamp, newest first.
//
// A ticket with no parseable stamp is SKIPPED rather than dated from the zero
// time: a fabricated instant would sort it last forever and read as real
// evidence that nothing happened to it, when the truth is that the board cannot
// tell. Absence stays absence.
func buildActivityRows(tickets []*RequestTicket) []ActivityRow {
	rows := make([]ActivityRow, 0, len(tickets))
	for _, ticket := range tickets {
		if ticket == nil {
			continue
		}
		newest, ok := newestLifecycleStamp(ticket)
		if !ok {
			continue
		}
		rows = append(rows, newest)
	}
	sort.Slice(rows, func(first, second int) bool {
		if !rows[first].StampTime.Equal(rows[second].StampTime) {
			return rows[first].StampTime.After(rows[second].StampTime)
		}
		// Ties are broken in the comparator, never left to the sort's stability:
		// otherwise the input order decides, and the order silently changes the
		// next time the tree walk returns tickets in a different sequence.
		return rows[first].RequestId > rows[second].RequestId
	})
	return rows
}

// newestLifecycleStamp picks the latest parseable stamp on one ticket. It reads
// every field lifecycleTimestampFields names and compares them all, rather than
// stopping at the first present one — created_at is first in that list and
// almost every ticket has one, so a first-match reader would report "captured"
// for a REQ that was claimed, built and merged this afternoon.
func newestLifecycleStamp(ticket *RequestTicket) (ActivityRow, bool) {
	var newest ActivityRow
	found := false
	for _, field := range lifecycleTimestampFields(ticket) {
		parsedInstant, parsedOk := parseTimestamp(field.RawValue)
		if !parsedOk {
			continue
		}
		if found && !parsedInstant.After(newest.StampTime) {
			continue
		}
		newest = ActivityRow{
			RequestId:  ticket.RequestId,
			StampField: field.FieldName,
			StampTime:  parsedInstant,
			Transition: field.Transition,
		}
		found = true
	}
	return newest, found
}
