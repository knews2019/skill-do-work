package main

import (
	"sort"
	"time"
)

// Timeline aggregation for the board's Timeline view.
//
// One row per REQ, two spans per row: created_at→claimed_at is the WAIT, and
// claimed_at→completed_at is the WORK. Everything here comes from stamps the
// board already parsed onto each ticket — no new frontmatter field, no second
// walk of the archive, the same posture durations.go takes.
//
// Two things separate this from the Durations view rather than duplicating it.
// The wait is measured nowhere else on the board, and on a queue with real
// backlog it is usually the larger share of calendar time. And a REQ that is
// claimed and still running contributes to no other panel of any view, so an
// open span is a first-class shape here rather than a missing value.
//
// Spans are SIGNED and never clamped. A negative span is the reversed-stamp
// anomaly the board already surfaces, and this file consumes that verdict from
// the ticket rather than restating the rule: detectCompletionAnomaly
// (model.go) decides what is broken, this decides how it is measured.

// TimelineRow is one REQ's bar: the instants it is drawn between, both spans,
// and the flags saying which spans are still running.
type TimelineRow struct {
	RequestId string

	CreatedTime   time.Time // required — a ticket without a parseable created_at has no row
	ClaimedTime   time.Time // zero when never claimed
	CompletedTime time.Time // zero when not completed

	// WaitMinutes is claimed_at − created_at, or now − created_at when the REQ
	// was never claimed, in which case WaitOpen is true.
	WaitMinutes float64
	WaitOpen    bool

	// HasWork is false for a REQ that was never claimed. It has no work segment
	// and none may be invented for it — a projected segment is a separate
	// concern with its own REQ.
	HasWork bool

	// WorkMinutes is completed_at − claimed_at, or now − claimed_at while the
	// REQ is in flight, in which case WorkOpen is true.
	WorkMinutes float64
	WorkOpen    bool

	// Copied from the ticket, never recomputed. The board decides what counts
	// as broken bookkeeping in exactly one place.
	Anomaly       bool
	AnomalyReason string
}

// TimelineAggregate is the whole view's data: rows in creation order, the time
// range they span, and the single instant every open span was measured against.
type TimelineAggregate struct {
	Rows []TimelineRow

	// RangeStart is the earliest created_at and RangeEnd the latest instant any
	// row reaches, Now included — an open bar must fit inside the fitted view.
	RangeStart time.Time
	RangeEnd   time.Time

	// Now is the board's generation instant, and it is the ONLY now in this
	// view: the same value measures every open span here and positions the
	// client's now-line. A live board regenerates per request, so its now is the
	// request instant; a snapshot's is frozen at generation, which the header's
	// existing relative "Generated … ago" makes visible.
	Now time.Time
}

// buildTimelineAggregate derives the Timeline view's rows from tickets the board
// already parsed. Like buildDurationAggregate it is a pass over memory.
func buildTimelineAggregate(tickets []*RequestTicket, now time.Time) TimelineAggregate {
	aggregate := TimelineAggregate{Now: now.UTC()}

	for _, ticket := range tickets {
		if ticket == nil {
			continue
		}
		createdInstant, createdParsed := parseTimestamp(ticket.CreatedAt)
		if !createdParsed {
			continue
		}
		row := TimelineRow{
			RequestId:     ticket.RequestId,
			CreatedTime:   createdInstant.UTC(),
			Anomaly:       ticket.CompletionAnomaly,
			AnomalyReason: ticket.CompletionAnomalyReason,
		}

		claimedInstant, claimedParsed := parseTimestamp(ticket.ClaimedAt)
		if !claimedParsed {
			// Never claimed: the wait is still running against the board's now,
			// and there is no work segment to draw.
			row.WaitMinutes = aggregate.Now.Sub(row.CreatedTime).Minutes()
			row.WaitOpen = true
			aggregate.Rows = append(aggregate.Rows, row)
			continue
		}

		row.ClaimedTime = claimedInstant.UTC()
		row.WaitMinutes = row.ClaimedTime.Sub(row.CreatedTime).Minutes()
		row.HasWork = true

		completedInstant, completedParsed := parseTimestamp(ticket.CompletedAt)
		if !completedParsed {
			row.WorkMinutes = aggregate.Now.Sub(row.ClaimedTime).Minutes()
			row.WorkOpen = true
			aggregate.Rows = append(aggregate.Rows, row)
			continue
		}
		row.CompletedTime = completedInstant.UTC()
		row.WorkMinutes = row.CompletedTime.Sub(row.ClaimedTime).Minutes()
		aggregate.Rows = append(aggregate.Rows, row)
	}

	// Chronological by created_at, oldest first, with the id as the tiebreak so
	// two REQs captured in the same second cannot swap places between builds.
	// The view states this ordering rather than leaving the reader to infer it.
	sort.SliceStable(aggregate.Rows, func(earlier, later int) bool {
		if aggregate.Rows[earlier].CreatedTime.Equal(aggregate.Rows[later].CreatedTime) {
			return aggregate.Rows[earlier].RequestId < aggregate.Rows[later].RequestId
		}
		return aggregate.Rows[earlier].CreatedTime.Before(aggregate.Rows[later].CreatedTime)
	})
	aggregate.RangeStart, aggregate.RangeEnd = timelineRange(aggregate.Rows, aggregate.Now)
	return aggregate
}

// timelineRange is the span a fitted view has to cover: the earliest instant any
// row starts to the latest any row reaches. Now is included whenever a row is
// still open, so an open bar cannot run off the end of a fitted view.
func timelineRange(rows []TimelineRow, now time.Time) (time.Time, time.Time) {
	if len(rows) == 0 {
		return now, now
	}
	rangeStart := rows[0].CreatedTime
	rangeEnd := rows[0].CreatedTime
	for _, row := range rows {
		if row.CreatedTime.Before(rangeStart) {
			rangeStart = row.CreatedTime
		}
		for _, instant := range []time.Time{row.CreatedTime, row.ClaimedTime, row.CompletedTime} {
			if !instant.IsZero() && instant.After(rangeEnd) {
				rangeEnd = instant
			}
		}
		if (row.WaitOpen || row.WorkOpen) && now.After(rangeEnd) {
			rangeEnd = now
		}
	}
	return rangeStart, rangeEnd
}
