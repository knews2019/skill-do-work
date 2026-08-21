package main

import (
	"math"
	"sort"
	"strconv"
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

	// One verdict per direct-label band. The placed labels carry their row on
	// the sample; these carry what could not be placed.
}

// buildDurationAggregate derives the Durations view's data from tickets the
// board already parsed. It is a pass over memory, not over the archive.
func buildDurationAggregate(tickets []*RequestTicket) DurationAggregate {
	aggregate := DurationAggregate{}

	for _, ticket := range tickets {
		if ticket == nil || !isCompletedStatus(ticket.Status) {
			continue
		}
		claimedInstant, claimedParsed := parseTimestamp(ticket.ClaimedAt)
		completedInstant, completedParsed := parseTimestamp(ticket.CompletedAt)
		if !claimedParsed || !completedParsed {
			continue
		}
		wallSpan := completedInstant.Sub(claimedInstant)
		aggregate.Samples = append(aggregate.Samples, DurationSample{
			RequestId:          ticket.RequestId,
			Route:              ticket.Route,
			CompletionTime:     completedInstant.UTC(),
			DayKey:             completedInstant.UTC().Format("2006-01-02"),
			WallMinutes:        wallSpan.Minutes(),
			DayMedianExclusion: dayMedianExclusionReason(wallSpan),
			EffortEstimate:     ticket.EffortEstimate,
		})
	}

	sort.SliceStable(aggregate.Samples, func(earlier, later int) bool {
		return aggregate.Samples[earlier].CompletionTime.Before(aggregate.Samples[later].CompletionTime)
	})
	aggregate.Days = buildDurationDays(aggregate.Samples)
	return aggregate
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

// ---- direct-label placement -------------------------------------------------
//
// Panel A draws a direct label beside a mark only where the mark carries a value
// its y cannot: the overflow lane, where every mark sits at one y so the text is
// the only carrier of magnitude, and the reversed band. WHICH of those marks get
// a label is decided here rather than in the renderer, so the rule has one
// reader — the lesson REQ-219 recorded about this very view.
//
// The geometry below is the SVG's own user-unit space, the space its viewBox
// defines. The chart is width:100% over a fixed viewBox, so a label's share of
// the plot is identical at every browser zoom and window width; user units are
// the only frame in which the question has a stable answer.

const (
	// durationsPlotWidthUnits mirrors the renderer's DURATIONS_VIEW_WIDTH minus
	// its two margins, and durationsOverflowCeilingMinutes its
	// DURATIONS_CEILING_MINUTES. TestDurationLabelGeometryMatchesTheRenderer
	// pins both against web/board-durations.js so they cannot drift apart.
	durationsPlotWidthUnits         = 1128.0
	durationsOverflowCeilingMinutes = 60.0
)

// / durationLabelText is the label the renderer draws for one sample. It exists
// here because the placement rule has to size the text it is placing.
func durationLabelText(sample DurationSample) string {
	return sample.RequestId + " " + formatDurationLabelMinutes(sample.WallMinutes)
}

// formatDurationLabelMinutes mirrors the renderer's formatDurationMinutes. Only
// the character count matters to placement, so this stays a width model rather
// than becoming a second definition of the view's copy.
//
// One known divergence: the sign below is an ASCII hyphen where the renderer
// draws U+2212 MINUS SIGN, and the two glyphs are far from the same width — on
// Chromium 141.0.7390.37 (Playwright 1.56.1, REQ-252) the minus measures
// 9.2015 units in the 11px face against the hyphen's 3.9642, and even that
// delta is per-browser (an earlier build measured it at 1.73). Width-neutral
// today, because both are one character and the model above counts characters —
// but a per-glyph width model (attempted and abandoned by REQ-241) would
// under-state every reversed label unless it models the minus the renderer
// actually draws, not the hyphen this string carries.
func formatDurationLabelMinutes(minutes float64) string {
	sign := ""
	if minutes < 0 {
		sign = "-"
	}
	// Mirrors the renderer's rounding order, not just its branches: round to the
	// displayed precision first, then split. Splitting first let the remainder
	// round to 60 and print "1h 60m" for 119.5 — one character wider than the
	// "2h 0m" the renderer should draw, so the width model was wrong too.
	displayedMinutes := math.Round(math.Abs(minutes)*10) / 10
	if displayedMinutes < 60 {
		return sign + strconv.FormatFloat(displayedMinutes, 'f', 1, 64) + " min"
	}
	wholeMinutes := int(math.Round(displayedMinutes))
	return sign + strconv.Itoa(wholeMinutes/60) + "h " + strconv.Itoa(wholeMinutes%60) + "m"
}
