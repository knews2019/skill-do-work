package main

import (
	"sort"
	"time"
)

// Activity aggregation for the board's Activity view.
//
// One row per lifecycle STAMP: the field it lives in, the instant it holds, and
// the transition it records. A REQ captured, claimed, dispatched, merged,
// reviewed, released and completed inside the window carries seven rows, so the
// whole path it took is readable here instead of only in its frontmatter or in
// `git log`. Newest first, and STATUS DOES NOT FILTER IT — a REQ that was
// claimed, blocked, completed, cancelled or failed inside the window all belong
// on the same surface.
//
// The question this answers is "what changed on the queue in the last N hours,
// and why", which no other surface answers. Recently done is terminal states
// only. The Timeline draws wait and work spans but not the transitions between
// them. The Calendar dates a REQ from its claim or resolve day. The detail
// drawer prints Created, Claimed and Completed and nothing between them. When
// three REQs were claimed, built, merged and held between 14:57 and 16:39 UTC,
// every one of those surfaces showed a two-hour hole and the answer needed
// `git log`.
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

// ActivityRow is one lifecycle transition of one REQ. A REQ contributes as many
// rows as it carries parseable stamps, so RequestId is not unique across rows.
type ActivityRow struct {
	RequestId string

	// StampField is the frontmatter field StampTime was read from. It ships to
	// the client so a reader can go straight to the line that produced the row,
	// which is what makes a missing stamp diagnosable rather than merely absent.
	StampField string
	StampTime  time.Time

	// Transition is what happened, in words — "work merged", not
	// "integration_at". Decided in Go by lifecycleTimestampFields so the client
	// never becomes a second definition of what a stamp means.
	Transition string
}

// buildActivityRows returns one row per parseable lifecycle stamp across every
// ticket, newest first. lifecycleTimestampFields (model.go) is the one
// enumeration of which stamps exist and what each records, so a stamp added to
// the schema reaches this surface without an edit here.
//
// A stamp that does not parse is SKIPPED rather than dated from the zero time:
// a fabricated instant would sort it last forever and read as real evidence
// that nothing happened, when the truth is that the board cannot tell. Absence
// stays absence, which is also why a ticket carrying no parseable stamp at all
// contributes no rows.
func buildActivityRows(tickets []*RequestTicket) []ActivityRow {
	rows := make([]ActivityRow, 0, len(tickets))
	for _, ticket := range tickets {
		if ticket == nil {
			continue
		}
		for _, field := range lifecycleTimestampFields(ticket) {
			parsedInstant, parsedOk := parseTimestamp(field.RawValue)
			if !parsedOk {
				continue
			}
			rows = append(rows, ActivityRow{
				RequestId:  ticket.RequestId,
				StampField: field.FieldName,
				StampTime:  parsedInstant,
				Transition: field.Transition,
			})
		}
	}
	sort.Slice(rows, func(first, second int) bool {
		if !rows[first].StampTime.Equal(rows[second].StampTime) {
			return rows[first].StampTime.After(rows[second].StampTime)
		}
		// Ties are broken in the comparator, never left to the sort's stability:
		// otherwise the input order decides, and the order silently changes the
		// next time the tree walk returns tickets in a different sequence.
		if rows[first].RequestId != rows[second].RequestId {
			return rows[first].RequestId > rows[second].RequestId
		}
		// The REQ id cannot separate two stamps of ONE ticket written at one
		// instant, which is the common case rather than the exotic one: the work
		// loop writes completed_at and status_changed_at together. The stamp
		// field is the last key. It sorts ascending because a field name carries
		// no notion of newer or older — the direction is arbitrary and only the
		// determinism matters, unlike the id above, where the larger number is
		// the newer REQ.
		return rows[first].StampField < rows[second].StampField
	})
	return rows
}
