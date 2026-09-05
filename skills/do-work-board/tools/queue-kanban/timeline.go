package main

import (
	"fmt"
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

	// WaitMinutes is claimed_at − created_at, or now − created_at when a REQ that
	// is still in the queue has not been claimed, in which case WaitOpen is true.
	//
	// OPEN MEANS STILL RUNNING, and that is a question about the REQ's STATE, not
	// about whether a stamp string happened to parse. Inferring it from a missing
	// stamp made this board describe 25 of 26 rows as "still open" that were
	// `completed` or `cancelled`, each drawn as a dashed bar running to the
	// now-line — REQ-059, completed, drew 24.8 days of work in flight.
	WaitMinutes float64
	WaitOpen    bool

	// HasWork is false for a REQ that was never claimed. It has no work segment
	// and none may be invented for it — a projected segment is a separate
	// concern with its own REQ.
	HasWork bool

	// WorkMinutes is completed_at − claimed_at, or now − claimed_at while the
	// REQ is genuinely in flight, in which case WorkOpen is true. Same rule as
	// WaitOpen above.
	WorkMinutes float64
	WorkOpen    bool

	// Copied from the ticket, never recomputed. The board decides what counts
	// as broken bookkeeping in exactly one place.
	Anomaly       bool
	AnomalyReason string

	// A REQ that has STOPPED and whose end instant the board could not resolve
	// carries NO flag of its own, deliberately. It is already fully described:
	// CompletedTime is zero and neither open flag is set, and no other row shape
	// can say that — a running one sets an open flag, a resolved one carries the
	// instant. The client reads that pair and draws the break marker the summary
	// and the legend already promise. A third field would be a second way to say
	// the same thing, and the first thing to drift.
}

// TimelineAggregate is the whole view's data: rows newest-created first, the time
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

		// Has this REQ stopped? The board decides that in one place, and this
		// consumes the verdict rather than listing statuses of its own. It is the
		// question every "open" flag below turns on: a span is open because the
		// work is still running, never because a stamp failed to parse.
		//
		// isStoppedStatus, NOT isTerminalResolvedStatus. The narrower one excludes
		// `failed`, and reading it here made every failed REQ draw an open bar to the
		// now-line with its completed_at sitting right there — a regression this file
		// introduced and Codex caught on the pull request.
		hasStopped := isStoppedStatus(ticket.Status)
		// WHERE it stopped.
		//
		// CompletionTime is the instant buildBoard already worked out — from frontmatter
		// or from the commit hash's committer date (resolveCompletionTime, model.go) —
		// and reading it rather than re-parsing completed_at is what stopped REQ-311
		// and REQ-302 being drawn as still waiting beside calendar entries dating them
		// days earlier.
		//
		// completed_at is the fallback for the one status buildBoard declines to
		// resolve: a `failed` REQ, whose own stamp is where buildCalendar dates it from
		// too, so this mirrors that rather than inventing a second rule.
		//
		// The ORDER of the two is not observable, and saying so is worth more than a
		// claim that it is: resolveCompletionTime derives CompletionTime FROM
		// completed_at whenever that stamp parses, so a git-derived instant only exists
		// where the stamp did not parse. They cannot disagree. CompletionTime is read
		// first because it is the board's own verdict, not because the outcome depends
		// on it — swapping the branches passes every test in this file, deliberately.
		endInstant := time.Time{}
		endResolved := false
		if hasStopped {
			if !ticket.CompletionTime.IsZero() {
				endInstant = ticket.CompletionTime.UTC()
				endResolved = true
			} else if stoppedAt, stoppedParsed := parseTimestamp(ticket.CompletedAt); stoppedParsed {
				endInstant = stoppedAt.UTC()
				endResolved = true
			}
		}
		claimedInstant, claimedParsed := parseTimestamp(ticket.ClaimedAt)
		if !claimedParsed {
			// Never claimed. Still in the queue, so the wait is genuinely running
			// against the board's now; stopped, so it ended when the REQ did.
			// Either way there is no work segment to draw.
			switch {
			case !hasStopped:
				row.WaitMinutes = aggregate.Now.Sub(row.CreatedTime).Minutes()
				row.WaitOpen = true
			case endResolved:
				row.WaitMinutes = endInstant.Sub(row.CreatedTime).Minutes()
				row.CompletedTime = endInstant
			}
			aggregate.Rows = append(aggregate.Rows, row)
			continue
		}

		row.ClaimedTime = claimedInstant.UTC()
		row.WaitMinutes = row.ClaimedTime.Sub(row.CreatedTime).Minutes()
		row.HasWork = true

		switch {
		case !hasStopped:
			// In flight: measured to the board's now, and the only case that may
			// say so.
			row.WorkMinutes = aggregate.Now.Sub(row.ClaimedTime).Minutes()
			row.WorkOpen = true
		case endResolved:
			row.CompletedTime = endInstant
			row.WorkMinutes = endInstant.Sub(row.ClaimedTime).Minutes()
		default:
			// Stopped, no resolvable end: the work segment has no width to draw, so
			// WorkMinutes stays 0 and CompletedTime stays zero. HasWork remains true
			// — work WAS started — and a zero CompletedTime with no open flag is what
			// tells the client to draw a break rather than a bar.
		}
		aggregate.Rows = append(aggregate.Rows, row)
	}

	// By created_at, NEWEST first, with the id as the tiebreak so two REQs
	// captured in the same second cannot swap places between builds. The view
	// states this ordering rather than leaving the reader to infer it.
	//
	// Newest-first because the reader is almost always asking about current
	// work: on a queue carrying hundreds of archived REQs, oldest-first put
	// today's row thousands of pixels below the fold and opened every visit on
	// REQ-001. The tiebreak reverses WITH the instants rather than staying
	// ascending — one instant's rows read the same direction as the list around
	// them, and within a second the higher id was captured later. It stays
	// NUMERIC through requestIdLess: past four digits a lexical compare puts
	// REQ-999 above REQ-1000 and claims the older id came last.
	sort.SliceStable(aggregate.Rows, func(earlier, later int) bool {
		if aggregate.Rows[earlier].CreatedTime.Equal(aggregate.Rows[later].CreatedTime) {
			return requestIdLess(aggregate.Rows[later].RequestId, aggregate.Rows[earlier].RequestId)
		}
		return aggregate.Rows[earlier].CreatedTime.After(aggregate.Rows[later].CreatedTime)
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

// ---- forward projection -----------------------------------------------------
//
// The measured half of this view answers "what happened". This answers "what is
// left, in what order, and roughly when does it end" — deliberately with the
// crudest model that can be honest: one REQ at a time, each taking the median of
// recent work in its own effort bucket. No parallelism knob, no per-REQ estimate
// field, no cleverness.
//
// The crudeness is the point. A forecast is the kind of artifact people
// screenshot and quote, and this one ignores parallel builders, review loops,
// blocked-then-unblocked churn, and the fact that the queue grows while it
// drains. The mitigation is that the view states what it assumed and declines
// when the history is too thin — not that the number is made cleverer.
//
// Three rules this file CONSUMES and does not restate:
//   - which spans count: DurationSample carries the read-time rule's verdict,
//     decided once in dayMedianExclusionReason;
//   - which REQs are ready: RequestTicket.UnmetDependencies, computed against
//     depends_on by buildBoard;
//   - what order they run in: actions/work.md Step 1 takes dependency-ready
//     REQs in numeric id order, which is exactly the chain below. If the two can
//     disagree, the ordering source here is wrong.

const (
	// The rolling window: the most recent in-rule completions the medians are
	// taken over. Long enough to survive a slow week, short enough that a
	// six-month-old pace does not outvote this month's.
	timelineProjectionWindowSize = 60

	// Below this many samples a median is a coincidence. Applied per bucket
	// (fall back to the window's overall median) and again to the window as a
	// whole (decline entirely).
	timelineProjectionMinimumSamples = 5
)

// TimelineProjectedRow is one unstarted REQ's forecast slot.
type TimelineProjectedRow struct {
	RequestId string
	StartTime time.Time
	EndTime   time.Time
	Bucket    string // effortMechanical or effortSubstantive — the effort bucket its span came from
	Position  int    // 1-based place in the chain
}

// TimelineExclusion is one REQ the projection refuses to schedule, and why.
// Every one is listed: silently folding these into the chain is the single
// easiest way to make this view lie.
type TimelineExclusion struct {
	RequestId string
	Reason    string
}

// TimelineProjection is the forward half of the view.
type TimelineProjection struct {
	Rows       []TimelineProjectedRow
	Excluded   []TimelineExclusion
	QueueEnd   time.Time // zero when nothing was scheduled
	ChainStart time.Time

	// What the forecast rests on, so the view can state it rather than imply it.
	WindowSize           int
	WindowSampleCount    int
	TrivialSampleCount   int
	NormalSampleCount    int
	TrivialMedianMinutes float64
	NormalMedianMinutes  float64
	MinimumSamples       int

	// Confident is false when the window holds too little history to forecast
	// from. Rows and QueueEnd are then empty and DeclinedReason says why.
	Confident      bool
	DeclinedReason string
}

// buildTimelineProjection forecasts the unstarted queue.
func buildTimelineProjection(tickets []*RequestTicket, aggregate DurationAggregate, now time.Time) TimelineProjection {
	projection := TimelineProjection{
		WindowSize:     timelineProjectionWindowSize,
		MinimumSamples: timelineProjectionMinimumSamples,
		ChainStart:     now.UTC(),
	}

	trivialMinutes, normalMinutes, windowMinutes := timelineProjectionWindow(aggregate)
	projection.TrivialSampleCount = len(trivialMinutes)
	projection.NormalSampleCount = len(normalMinutes)
	projection.WindowSampleCount = len(windowMinutes)

	windowMedian, hasWindowMedian := medianOf(windowMinutes)
	if !hasWindowMedian || len(windowMinutes) < timelineProjectionMinimumSamples {
		projection.DeclinedReason = fmt.Sprintf(
			"only %d completed REQ%s inside the read-time rule; %d are needed before a median means anything",
			len(windowMinutes), pluralSuffix(len(windowMinutes)), timelineProjectionMinimumSamples)
		return projection
	}
	projection.Confident = true
	// A bucket too thin to speak for itself borrows the window's overall median
	// rather than inventing one from two samples.
	projection.TrivialMedianMinutes = timelineBucketMedian(trivialMinutes, windowMedian)
	projection.NormalMedianMinutes = timelineBucketMedian(normalMinutes, windowMedian)

	projection.ChainStart = timelineChainStart(tickets, projection, now)
	projection.Rows, projection.Excluded = timelineChain(tickets, projection)
	if len(projection.Rows) > 0 {
		projection.QueueEnd = projection.Rows[len(projection.Rows)-1].EndTime
	}
	return projection
}

// timelineProjectionWindow takes the most recent in-rule samples and splits them
// by effort bucket. The samples arrive already classified by the read-time rule,
// so nothing here decides what a paused or reversed span is.
func timelineProjectionWindow(aggregate DurationAggregate) ([]float64, []float64, []float64) {
	var trivialMinutes, normalMinutes, windowMinutes []float64
	for sampleIndex := len(aggregate.Samples) - 1; sampleIndex >= 0; sampleIndex-- {
		sample := aggregate.Samples[sampleIndex]
		if sample.ExcludedFromDayMedian() {
			continue
		}
		if len(windowMinutes) >= timelineProjectionWindowSize {
			break
		}
		windowMinutes = append(windowMinutes, sample.WallMinutes)
		if sample.EffortEstimate == effortMechanical {
			trivialMinutes = append(trivialMinutes, sample.WallMinutes)
		} else {
			normalMinutes = append(normalMinutes, sample.WallMinutes)
		}
	}
	return trivialMinutes, normalMinutes, windowMinutes
}

// timelineBucketMedian is a bucket's own median once it has enough samples to
// mean anything, and the window's overall median until then.
func timelineBucketMedian(bucketMinutes []float64, windowMedian float64) float64 {
	if len(bucketMinutes) < timelineProjectionMinimumSamples {
		return windowMedian
	}
	bucketMedian, hasBucketMedian := medianOf(bucketMinutes)
	if !hasBucketMedian {
		return windowMedian
	}
	return bucketMedian
}

// timelineProjectedSpan is how long one REQ is forecast to take.
func timelineProjectedSpan(ticket *RequestTicket, projection TimelineProjection) (time.Duration, string) {
	// effort_estimate is a closed two-value triage bit whose documented default
	// is effort-substantive; absent reads as substantive rather than as a third
	// case. Compared against the named constants on purpose: a bare literal here
	// is what makes a rename of the vocabulary compile and silently drop every
	// REQ into the substantive median.
	if ticket.EffortEstimate == effortMechanical {
		return time.Duration(projection.TrivialMedianMinutes * float64(time.Minute)), effortMechanical
	}
	return time.Duration(projection.NormalMedianMinutes * float64(time.Minute)), effortSubstantive
}

// timelineChainStart is when the first unstarted REQ can begin: after whatever is
// already running. The in-flight REQ's own `estimate:` block would be the better
// offset for exactly this bar, but the board parses no nested frontmatter blocks,
// and adding that surface for one bar is the sophistication this REQ trades for a
// stated assumption.
func timelineChainStart(tickets []*RequestTicket, projection TimelineProjection, now time.Time) time.Time {
	chainStart := now.UTC()
	for _, ticket := range tickets {
		if ticket == nil || ticket.Status != "claimed" {
			continue
		}
		projectedSpan, _ := timelineProjectedSpan(ticket, projection)
		claimedInstant, claimedParsed := parseTimestamp(ticket.ClaimedAt)
		projectedStart := now.UTC()
		if claimedParsed {
			projectedStart = claimedInstant.UTC()
		} else {
			// The timestamp defect is reported by verify; the forecast neither repairs
			// nor trusts it. Reserving one whole bucket median from now is the
			// conservative alternative to pretending the active work takes no time.
		}
		projectedFinish := projectedStart.Add(projectedSpan)
		if projectedFinish.After(chainStart) {
			chainStart = projectedFinish
		}
	}
	return chainStart
}

// timelineChain places every schedulable pending REQ, and reports the rest.
//
// The placement rule is work.md's: among REQs whose dependencies have all
// resolved or been placed already, take the lowest id — NUMERICALLY, through
// requestIdLess, because past four digits a lexical order puts REQ-1000 ahead of
// REQ-999 and the forecast would name the wrong next REQ. Repeating that until
// nothing more can be placed is what makes the chain agree with what `do-work
// run` would actually claim next.
func timelineChain(tickets []*RequestTicket, projection TimelineProjection) ([]TimelineProjectedRow, []TimelineExclusion) {
	var pendingTickets []*RequestTicket
	var exclusions []TimelineExclusion
	for _, ticket := range tickets {
		if ticket == nil {
			continue
		}
		switch {
		case ticket.Status == "pending":
			pendingTickets = append(pendingTickets, ticket)
		case isNeedsInputOrBlockedStatus(ticket.Status) && ticket.Status != "failed",
			!isKnownSchemaFieldValue("status", ticket.Status):
			// Two kinds of REQ that get no start time but must still be NAMED. One
			// waits on a person or an external condition. The other carries a status
			// outside the Schema Read Contract vocabulary — a typo, or a value
			// retired from the schema and still sitting in an old file. This is the
			// only arm that names anything, so whatever falls past it leaves the
			// chain, the queue-end figure AND the exclusion list at once, which reads
			// exactly like having no such REQ at all. bucketColumns (model.go)
			// already warns about the second kind; the forecast owes that reader a
			// line rather than silence.
			exclusions = append(exclusions, TimelineExclusion{
				RequestId: ticket.RequestId,
				Reason:    timelineExclusionReason(ticket),
			})
		}
	}
	sort.SliceStable(pendingTickets, func(earlier, later int) bool {
		return requestIdLess(pendingTickets[earlier].RequestId, pendingTickets[later].RequestId)
	})

	// Seeded with the REQs already in flight. UnmetDependencies holds anything
	// that has not reached terminal success, which includes a CLAIMED
	// prerequisite — but the chain starts after in-flight work finishes
	// (timelineChainStart), so by the time any of this runs that dependency has
	// resolved. Without the seed a REQ waiting on claimed work could never
	// become ready: it dropped out of the chain, out of the queue-end figure,
	// and into the exclusion list under a reason that was wrong about it.
	placedIds := map[string]bool{}
	for _, ticket := range tickets {
		if ticket != nil && ticket.Status == "claimed" {
			placedIds[ticket.RequestId] = true
		}
	}
	var rows []TimelineProjectedRow
	chainCursor := projection.ChainStart
	for {
		nextIndex := -1
		for ticketIndex, ticket := range pendingTickets {
			if ticket == nil {
				continue
			}
			ready := true
			for _, dependencyId := range ticket.UnmetDependencies {
				if !placedIds[dependencyId] {
					ready = false
					break
				}
			}
			if ready {
				nextIndex = ticketIndex
				break
			}
		}
		if nextIndex == -1 {
			break
		}
		ticket := pendingTickets[nextIndex]
		pendingTickets[nextIndex] = nil
		projectedSpan, bucket := timelineProjectedSpan(ticket, projection)
		rows = append(rows, TimelineProjectedRow{
			RequestId: ticket.RequestId,
			StartTime: chainCursor,
			EndTime:   chainCursor.Add(projectedSpan),
			Bucket:    bucket,
			Position:  len(rows) + 1,
		})
		placedIds[ticket.RequestId] = true
		chainCursor = chainCursor.Add(projectedSpan)
	}

	// Whatever the loop could not place can never be placed: its dependency
	// chain reaches something excluded, dangling, or circular. Reporting it is
	// what keeps the queue-end figure from quietly describing a subset.
	for _, ticket := range pendingTickets {
		if ticket == nil {
			continue
		}
		exclusions = append(exclusions, TimelineExclusion{
			RequestId: ticket.RequestId,
			Reason:    "waiting on a dependency that is itself unschedulable, missing, or circular",
		})
	}
	sort.SliceStable(exclusions, func(earlier, later int) bool {
		return requestIdLess(exclusions[earlier].RequestId, exclusions[later].RequestId)
	})
	return rows, exclusions
}

// timelineExclusionReason says in plain words why a REQ cannot be given a start
// time, rather than echoing its status back at the reader.
func timelineExclusionReason(ticket *RequestTicket) string {
	switch ticket.Status {
	case "pending-answers":
		return "waiting on an answer from you"
	case "blocked":
		return "waiting on an external condition"
	case "blocked-archive-collision":
		return "held: its id is already archived"
	case "blocked-dependency-cycle":
		return "held: its dependency chain is circular"
	default:
		// Every recognized status that reaches here has a case above, so this is
		// the off-vocabulary arm: say the status is the problem instead of
		// implying the REQ is merely resting.
		return "its status is not in the schema vocabulary, so it cannot be scheduled"
	}
}

func pluralSuffix(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
