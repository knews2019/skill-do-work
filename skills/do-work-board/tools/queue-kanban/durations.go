package main

import (
	"fmt"
	"sort"
	"time"
)

// Duration aggregation for the board's Durations view.
//
// Everything here is derived from tickets the board has already parsed — no new
// frontmatter field, no second walk of the archive. A REQ contributes one sample
// when it reached terminal success and BOTH its claim and completion stamps
// parse.
//
// The wall span is recorded raw and signed. A negative span is data, not an
// error to swallow: it is the reversed-stamp anomaly class the board learned to
// surface, and rounding it up to zero would hide exactly the bookkeeping bug the
// reader is meant to see.
//
// Panels A and B deliberately disagree about which samples count, and that split
// is the point. Panel A plots every sample raw, so an outlier stays visible.
// Panel B applies the calibration's documented read-time rule — a span over four
// hours is an assumed paused session, a negative span is a broken stamp — so one
// paused session cannot invent a five-hour day. The rule is stated once, in
// skills/do-work/actions/estimate-reference.md → Calibration; this is its second
// reader, not a second definition.

// analysisOutlierCeiling is the read-time rule's upper bound: a wall span longer
// than this is assumed to include a pause rather than four solid hours of work.
const analysisOutlierCeiling = 4 * time.Hour

// implementationSpanPausedBadgeText turns the read-time ceiling into the
// human-facing marker carried by done cards. The client receives the completed
// label rather than a numeric ceiling, so Go remains the only place that decides
// which spans cross the rule.
func implementationSpanPausedBadgeText(ceiling time.Duration) string {
	wholeHours := int(ceiling / time.Hour)
	remainingMinutes := int((ceiling % time.Hour) / time.Minute)
	thresholdText := ""
	switch {
	case wholeHours > 0 && remainingMinutes > 0:
		thresholdText = fmt.Sprintf("%dh%dm", wholeHours, remainingMinutes)
	case wholeHours > 0:
		thresholdText = fmt.Sprintf("%dh", wholeHours)
	default:
		thresholdText = fmt.Sprintf("%dm", remainingMinutes)
	}
	return "over " + thresholdText + " · assumed pause"
}

// DurationSample is one archived REQ's measured wall span.
type DurationSample struct {
	RequestId      string
	Route          string    // normalized route ("A"/"B"/"C"), "" when the REQ predates routing
	CompletionTime time.Time // the parsed completed_at instant
	DayKey         string    // "2006-01-02" UTC bucket, shared with the calendar view
	WallMinutes    float64   // completed_at − claimed_at, raw and signed

	// Why the read-time rule excluded this sample from the day medians:
	// "paused" (over the ceiling), "reversed" (negative), or "" when it counts.
	DayMedianExclusion string

	// The REQ's effort_estimate bucket, normalized (effortMechanical or
	// effortSubstantive; absent reads as substantive, per the closed two-value
	// enum in model.go, whose read-only aliases carry the pre-rename
	// trivial/normal spellings on archived REQs). Panel
	// A and the medians ignore it; the Timeline's forward projection splits its
	// medians by it, and carrying it here is what lets that projection read the
	// bucket off a sample the read-time rule has ALREADY classified rather than
	// re-walking the tickets and re-deciding what counts.
	EffortEstimate string
}

// ExcludedFromDayMedian reports whether the read-time rule holds this sample out
// of the per-day medians. Panel A still plots it.
func (sample DurationSample) ExcludedFromDayMedian() bool {
	return sample.DayMedianExclusion != ""
}

// DurationDay is one active calendar day: the median of the samples the
// read-time rule keeps, and the count of every sample that landed that day.
type DurationDay struct {
	DayKey  string
	DayTime time.Time // midnight UTC of DayKey, so the day plots on the shared calendar axis

	// MedianMinutes over the samples the rule kept. HasMedian is false when the
	// day's only samples were excluded — a real state, not a zero.
	MedianMinutes float64
	HasMedian     bool

	KeptCount      int // samples inside the rule, i.e. the median's sample size
	CompletedCount int // every sample that day, rule-excluded ones included
}

// DurationAggregate is the whole Durations view's data: one sample per
// qualifying REQ in completion order, and one entry per active day.
type DurationAggregate struct {
	Samples []DurationSample
	Days    []DurationDay
}

// buildDurationAggregate derives the Durations view's data from tickets the
// board already parsed. It is a pass over memory, not over the archive.
func buildDurationAggregate(tickets []*RequestTicket) DurationAggregate {
	aggregate := DurationAggregate{}

	for _, ticket := range tickets {
		if ticket == nil || !isCompletedStatus(ticket.Status) {
			continue
		}
		measuredSpan := measureImplementationSpan(ticket)
		if !measuredSpan.StampsParsed {
			continue
		}
		aggregate.Samples = append(aggregate.Samples, DurationSample{
			RequestId:          ticket.RequestId,
			Route:              ticket.Route,
			CompletionTime:     measuredSpan.CompletionInstant,
			DayKey:             measuredSpan.CompletionInstant.Format("2006-01-02"),
			WallMinutes:        measuredSpan.WallMinutes,
			DayMedianExclusion: measuredSpan.ExclusionReason,
			EffortEstimate:     ticket.EffortEstimate,
		})
	}

	sort.SliceStable(aggregate.Samples, func(earlier, later int) bool {
		return aggregate.Samples[earlier].CompletionTime.Before(aggregate.Samples[later].CompletionTime)
	})
	aggregate.Days = buildDurationDays(aggregate.Samples)
	return aggregate
}

// ImplementationSpan is one REQ's measured claim-to-completion span with the
// read-time rule's verdict already applied. It is the single place that span and
// that verdict are decided: the Durations view's samples and the Recently-Done
// card's duration reading are both readers of it, so the ceiling above keeps
// exactly one definition and a card can never disagree with the chart.
//
// StampsParsed false is a real state, not a span of zero — a zero would print as
// instant work on the card instead of as unmeasured.
type ImplementationSpan struct {
	StampsParsed      bool
	CompletionInstant time.Time // parsed completed_at, in UTC
	WallMinutes       float64   // completed_at − claimed_at, raw and signed
	ExclusionReason   string    // "paused" (over the ceiling), "reversed" (negative), "" when it reads plainly
}

// measureImplementationSpan reads both instants off the ticket's FRONTMATTER and
// nowhere else. Deliberately not RequestTicket.CompletionTime: that field falls
// back to the commit's git committer date (see model.go's resolveCompletionTime),
// which measures when a commit landed rather than how long the work took, and
// would make a card state a duration for exactly the REQs the Durations view
// excludes. A REQ whose completion instant came from git therefore has no span.
//
// The caller decides WHICH tickets to ask about — this says nothing about status.
func measureImplementationSpan(ticket *RequestTicket) ImplementationSpan {
	if ticket == nil {
		return ImplementationSpan{}
	}
	claimedInstant, claimedParsed := parseTimestamp(ticket.ClaimedAt)
	completedInstant, completedParsed := parseTimestamp(ticket.CompletedAt)
	if !claimedParsed || !completedParsed {
		return ImplementationSpan{}
	}
	wallSpan := completedInstant.Sub(claimedInstant)
	return ImplementationSpan{
		StampsParsed:      true,
		CompletionInstant: completedInstant.UTC(),
		WallMinutes:       wallSpan.Minutes(),
		ExclusionReason:   dayMedianExclusionReason(wallSpan),
	}
}

// dayMedianExclusionReason applies the calibration's read-time rule to one span.
func dayMedianExclusionReason(wallSpan time.Duration) string {
	switch {
	case wallSpan < 0:
		return "reversed"
	case wallSpan > analysisOutlierCeiling:
		return "paused"
	default:
		return ""
	}
}

// buildDurationDays groups samples into active days, in chronological order.
func buildDurationDays(samples []DurationSample) []DurationDay {
	if len(samples) == 0 {
		return nil
	}

	keptMinutesByDay := map[string][]float64{}
	completedCountByDay := map[string]int{}
	dayOrder := []string{}
	for _, sample := range samples {
		if _, seen := completedCountByDay[sample.DayKey]; !seen {
			dayOrder = append(dayOrder, sample.DayKey)
		}
		completedCountByDay[sample.DayKey]++
		if !sample.ExcludedFromDayMedian() {
			keptMinutesByDay[sample.DayKey] = append(keptMinutesByDay[sample.DayKey], sample.WallMinutes)
		}
	}
	sort.Strings(dayOrder)

	days := make([]DurationDay, 0, len(dayOrder))
	for _, dayKey := range dayOrder {
		dayTime, parseError := time.Parse("2006-01-02", dayKey)
		if parseError != nil {
			continue
		}
		keptMinutes := keptMinutesByDay[dayKey]
		medianMinutes, hasMedian := medianOf(keptMinutes)
		days = append(days, DurationDay{
			DayKey:         dayKey,
			DayTime:        dayTime.UTC(),
			MedianMinutes:  medianMinutes,
			HasMedian:      hasMedian,
			KeptCount:      len(keptMinutes),
			CompletedCount: completedCountByDay[dayKey],
		})
	}
	return days
}

// medianOf returns the median of the supplied values. The false second result is
// "no samples", which is not the same as a median of zero.
func medianOf(values []float64) (float64, bool) {
	if len(values) == 0 {
		return 0, false
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	middle := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return sorted[middle], true
	}
	return (sorted[middle-1] + sorted[middle]) / 2, true
}

// durationsPlotWidthUnits mirrors the renderer's fixed viewBox plot width. It
// remains as the independent side of the day-bucket projection test; no label
// sizing or placement reads it.
const durationsPlotWidthUnits = 1128.0
