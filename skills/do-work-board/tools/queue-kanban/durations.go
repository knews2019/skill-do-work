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

	// The REQ's effort_estimate bucket, normalized ("trivial" or "normal";
	// absent reads as normal, per the closed two-value enum in model.go). Panel
	// A and the medians ignore it; the Timeline's forward projection splits its
	// medians by it, and carrying it here is what lets that projection read the
	// bucket off a sample the read-time rule has ALREADY classified rather than
	// re-walking the tickets and re-deciding what counts.
	EffortEstimate string

	// Direct-label verdict, decided by planDurationLabels. LabelRow is the text
	// row inside this sample's band, or durationsLabelRowUnplaced when the
	// collision rule could not fit a label beside this mark. LabelAnchor is the
	// SVG text-anchor the placement assumed, "" when unplaced.
	LabelRow    int
	LabelAnchor string
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
	OverflowLabels DurationLabelBand
	ReversedLabels DurationLabelBand
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
	aggregate.OverflowLabels, aggregate.ReversedLabels = planDurationLabels(aggregate.Samples)
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

	// A label's width is estimated per character rather than measured. The exact
	// answer needs getComputedTextLength, which exists only after render and
	// would move this decision back into the client. Over-estimating drops a
	// label, which is visible and counted; under-estimating overprints one, which
	// is the defect this rule exists to fix — so the constant is an UPPER BOUND on
	// units per character over the whole label space, not the worst case of some
	// sample of it.
	//
	// The label space is INFINITE: a label is "REQ-" + id + " " + duration, and
	// neither the id's digit count nor the hour count is bounded by anything —
	// broken stamps produce arbitrarily large magnitudes. Per-character width
	// rises with length, because the narrow fixed characters ("-", the spaces,
	// "i", ".") are a fixed cost diluted by every digit added. So no sweep over
	// "realistic" labels can establish this bound, and two earlier attempts that
	// tried are why this comment is long.
	//
	// What makes it boundable anyway: only DIGITS can repeat without limit. Every
	// wide fixed character (R, E, Q, m, h, −) appears at most a couple of times
	// per label, so its contribution is amortized away as the label grows, and
	// per-character width is dragged toward — and cannot pass — the per-character
	// width of a pure digit run. That is a measurable quantity: 7.1441 user units,
	// the widest digit ("4") in the 11px sans face. Measured labels approach it
	// from below and never reach it (7.1417 at 40 010 characters), so 7.15 is
	// above the supremum of the whole space.
	// TestDurationsLabelWidthEstimateCoversTheRenderedFace pins this value on both
	// sides — under the supremum it under-states the text, and far over it drops
	// labels for nothing.
	//
	// It stood at 6.2 until REQ-241 and at 6.75 briefly within it; both were below
	// the face, and 6.2's comment claimed a generosity that ran the other way.
	durationsLabelCharacterWidthUnits = 7.15

	// Mark centre to the near edge of its text, matching the renderer's offset.
	durationsLabelGapUnits = 9.0

	// Minimum clear space between two labels sharing a row. With the width model
	// calibrated above the face (REQ-241) this is real whitespace rather than the
	// slack that used to absorb an under-estimate: at 6.2 units per character the
	// tightest gap in a saturated lane measured 3.08 units against the 6 the rule
	// claims to hold.
	durationsLabelSeparationUnits = 6.0

	// Text rows available in each band. Each band packs its own rows; they have
	// unrelated local densities and never share a counter.
	durationsLabelRowCount = 2

	// Width held back at the right edge of each band's LAST row for the remainder
	// sentence, sized from the widest sentence the renderer can compose rather
	// than from the composed string: the count is not known until placement has
	// finished. The last row rather than the first because the marks themselves
	// sit level with the first, so a sentence there is legible only while the
	// band is sparse — which is exactly when there is no remainder to print.
	durationsLabelRemainderReserveUnits = 24 * durationsLabelCharacterWidthUnits

	// LabelRow for a sample the collision rule could not place.
	durationsLabelRowUnplaced = -1
)

// DurationLabelBand is one direct-label band's verdict: how many of its samples
// carry no label because the collision rule could not fit one. A nonzero count is
// drawn on the chart as a remainder, so a reader can never mistake the visible
// labels for all of them.
type DurationLabelBand struct {
	HiddenCount int
}

// durationLabelWidthUnits estimates the space a sample's label text occupies.
func durationLabelWidthUnits(sample DurationSample) float64 {
	return float64(len(durationLabelText(sample))) * durationsLabelCharacterWidthUnits
}

// durationLabelText is the label the renderer draws for one sample. It exists
// here because the placement rule has to size the text it is placing.
func durationLabelText(sample DurationSample) string {
	return sample.RequestId + " " + formatDurationLabelMinutes(sample.WallMinutes)
}

// formatDurationLabelMinutes mirrors the renderer's formatDurationMinutes. Only
// the character count matters to placement, so this stays a width model rather
// than becoming a second definition of the view's copy.
func formatDurationLabelMinutes(minutes float64) string {
	magnitude := math.Abs(minutes)
	sign := ""
	if minutes < 0 {
		sign = "-"
	}
	if magnitude < 60 {
		return sign + strconv.FormatFloat(magnitude, 'f', 1, 64) + " min"
	}
	hours := int(magnitude) / 60
	remainder := int(math.Round(magnitude - float64(hours*60)))
	return sign + strconv.Itoa(hours) + "h " + strconv.Itoa(remainder) + "m"
}

// durationLabelTimeRange is the x-axis domain the renderer uses: first to last
// completion instant across every sample, overflow and ordinary alike. Placement
// has to share that domain or a label would be sized against a different plot
// than the one it lands on.
func durationLabelTimeRange(samples []DurationSample) (time.Time, time.Time, bool) {
	if len(samples) == 0 {
		return time.Time{}, time.Time{}, false
	}
	rangeStart := samples[0].CompletionTime
	rangeEnd := samples[0].CompletionTime
	for _, sample := range samples {
		if sample.CompletionTime.Before(rangeStart) {
			rangeStart = sample.CompletionTime
		}
		if sample.CompletionTime.After(rangeEnd) {
			rangeEnd = sample.CompletionTime
		}
	}
	return rangeStart, rangeEnd, true
}

// durationLabelPlotX places one completion instant in plot-local user units,
// 0 at the plot's left edge and durationsPlotWidthUnits at its right. A
// zero-width domain collapses to a single column, exactly as the renderer's
// `timeSpan || 1` guard does.
func durationLabelPlotX(completionTime time.Time, rangeStart time.Time, rangeEnd time.Time) float64 {
	domainSeconds := rangeEnd.Sub(rangeStart).Seconds()
	if domainSeconds <= 0 {
		return 0
	}
	return (completionTime.Sub(rangeStart).Seconds() / domainSeconds) * durationsPlotWidthUnits
}

// durationLabelBandOf reports which direct-label band a sample belongs to, or ""
// when it carries no direct label at all. The two bands sit at different heights
// with unrelated local densities, so they are packed independently.
func durationLabelBandOf(sample DurationSample) string {
	switch {
	case sample.WallMinutes < 0:
		return "reversed"
	case sample.WallMinutes > durationsOverflowCeilingMinutes:
		return "overflow"
	default:
		return ""
	}
}

// planDurationLabels decides which samples get a direct label. Each band packs
// its own rows in DESCENDING MAGNITUDE order: the longest span is offered a row
// first, and each span after it takes the first row where its text touches
// nothing already placed. The lane's text answers "where are the outliers", so
// the order is what sends it to superlatives — and because the walk simply keeps
// going down that order until nothing more fits, a span the geometry rejects
// frees its space to the next-longest span rather than to nobody (REQ-237,
// amending REQ-231's fixed six-candidate selection). A sample that fits nowhere
// is counted rather than drawn, and the renderer prints that count, so a reader
// can never mistake the visible labels for all of them.
func planDurationLabels(samples []DurationSample) (DurationLabelBand, DurationLabelBand) {
	for index := range samples {
		samples[index].LabelRow = durationsLabelRowUnplaced
		samples[index].LabelAnchor = ""
	}
	rangeStart, rangeEnd, hasRange := durationLabelTimeRange(samples)
	if !hasRange {
		return DurationLabelBand{}, DurationLabelBand{}
	}
	return packDurationLabelBand(samples, "overflow", rangeStart, rangeEnd),
		packDurationLabelBand(samples, "reversed", rangeStart, rangeEnd)
}

// packDurationLabelBand runs the greedy pass, at most twice. The remainder
// sentence needs room at the LAST row's right edge (durationsLabelRemainderReserveUnits
// says why there rather than the first), but whether there IS a remainder is
// only known once placement has finished — so the first pass uses the full
// width, and only a pass that actually dropped something is redone with the
// reservation held back. A board with no remainder pays nothing for one.
func packDurationLabelBand(samples []DurationSample, bandName string, rangeStart time.Time, rangeEnd time.Time) DurationLabelBand {
	labelOrder := durationLabelMagnitudeOrder(samples, bandName)
	band := placeDurationLabelBand(samples, labelOrder, rangeStart, rangeEnd, false)
	if band.HiddenCount == 0 {
		return band
	}
	return placeDurationLabelBand(samples, labelOrder, rangeStart, rangeEnd, true)
}

// durationLabelMagnitudeOrder lists one band's sample indexes longest span
// first — the order placement offers rows in, and so the whole of the "labels go
// to the outliers" rule. The sort is stable over completion order, so equal spans
// keep their left-to-right precedence and the choice is deterministic on an
// archive full of identical durations.
func durationLabelMagnitudeOrder(samples []DurationSample, bandName string) []int {
	labelOrder := []int{}
	for index := range samples {
		if durationLabelBandOf(samples[index]) == bandName {
			labelOrder = append(labelOrder, index)
		}
	}
	sort.SliceStable(labelOrder, func(first, second int) bool {
		return math.Abs(samples[labelOrder[first]].WallMinutes) >
			math.Abs(samples[labelOrder[second]].WallMinutes)
	})
	return labelOrder
}

// durationLabelInterval is one placed label's occupied x-interval on a row.
type durationLabelInterval struct {
	spanLeft  float64
	spanRight float64
}

// durationLabelSpanIsBlocked reports whether a candidate interval comes closer
// than durationsLabelSeparationUnits to a label already on that row. The walk is
// ordered by magnitude rather than by x, so a row's occupancy is no longer a
// single "reaches this far right" number and every interval has to be consulted.
func durationLabelSpanIsBlocked(occupied []durationLabelInterval, spanLeft float64, spanRight float64) bool {
	for _, placed := range occupied {
		if spanLeft < placed.spanRight+durationsLabelSeparationUnits &&
			placed.spanLeft < spanRight+durationsLabelSeparationUnits {
			return true
		}
	}
	return false
}

// placeDurationLabelBand is one greedy pass over a single band, taking its
// samples in the order given. Every sample the order names is offered a row, so
// the pass stops only when the rows are full: a span the geometry rejects costs
// the band nothing but its own label.
func placeDurationLabelBand(
	samples []DurationSample,
	labelOrder []int,
	rangeStart time.Time,
	rangeEnd time.Time,
	reserveRemainder bool,
) DurationLabelBand {
	occupied := make([][]durationLabelInterval, durationsLabelRowCount)

	band := DurationLabelBand{}
	for _, index := range labelOrder {
		samples[index].LabelRow = durationsLabelRowUnplaced
		samples[index].LabelAnchor = ""

		markX := durationLabelPlotX(samples[index].CompletionTime, rangeStart, rangeEnd)
		textWidth := durationLabelWidthUnits(samples[index])
		for rowIndex := 0; rowIndex < durationsLabelRowCount; rowIndex++ {
			rowRightEdge := durationsPlotWidthUnits
			if reserveRemainder && rowIndex == durationsLabelRowCount-1 {
				rowRightEdge -= durationsLabelRemainderReserveUnits
			}
			// Anchoring BEFORE the mark is tried first, with after-the-mark as the
			// fallback that keeps the leftmost sample labellable, since its own
			// label would otherwise start off-plot. REQ-231 chose that order as a
			// packing decision — the walk ran left to right, so a label drawn to
			// the mark's left reused space the walk had already passed. This walk
			// runs by magnitude, so that argument no longer applies and the order
			// is kept as a consistency one: a label sits on the same side of its
			// mark unless the plot edge forbids it. Neither anchor is an
			// alternation: both are geometric decisions about one specific label.
			spanLeft, spanRight, anchorFits := durationLabelSpan(markX, textWidth, "end", rowRightEdge)
			anchor := "end"
			if !anchorFits || durationLabelSpanIsBlocked(occupied[rowIndex], spanLeft, spanRight) {
				spanLeft, spanRight, anchorFits = durationLabelSpan(markX, textWidth, "start", rowRightEdge)
				anchor = "start"
			}
			if !anchorFits || durationLabelSpanIsBlocked(occupied[rowIndex], spanLeft, spanRight) {
				continue
			}
			samples[index].LabelRow = rowIndex
			samples[index].LabelAnchor = anchor
			occupied[rowIndex] = append(occupied[rowIndex], durationLabelInterval{spanLeft: spanLeft, spanRight: spanRight})
			break
		}
		if samples[index].LabelRow == durationsLabelRowUnplaced {
			band.HiddenCount++
		}
	}
	return band
}

// durationLabelSpan is the x-interval one label would occupy at one anchor, and
// whether that interval stays inside the plot.
func durationLabelSpan(markX float64, textWidth float64, anchor string, rowRightEdge float64) (float64, float64, bool) {
	spanLeft := markX + durationsLabelGapUnits
	spanRight := spanLeft + textWidth
	if anchor == "end" {
		spanRight = markX - durationsLabelGapUnits
		spanLeft = spanRight - textWidth
	}
	return spanLeft, spanRight, spanLeft >= 0 && spanRight <= rowRightEdge
}
