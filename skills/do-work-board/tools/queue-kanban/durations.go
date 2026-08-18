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
	// collision rule could not fit a label. LabelAnchor is the SVG text-anchor
	// the placement assumed, "" when unplaced.
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
	// would move this decision back into the client. The estimate is deliberately
	// generous for the 11px sans face: over-estimating drops a label, which is
	// visible and counted, while under-estimating overprints one, which is the
	// defect this rule exists to fix.
	durationsLabelCharacterWidthUnits = 6.2

	// Mark centre to the near edge of its text, matching the renderer's offset.
	durationsLabelGapUnits = 9.0

	// Minimum clear space between two labels sharing a row.
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

	// At most this many samples per band are label candidates, chosen by
	// magnitude (REQ-231, Alternative 2). The lane's text answers "where are
	// the outliers", so it goes to the longest spans rather than to whichever
	// spans a left-to-right walk happened to pack first; everything else is
	// carried by the remainder count, the hover readout, and the table.
	durationsLabelTopCount = 6

	// LabelRow for a sample that carries no direct label — either selection
	// did not pick it, or the collision rule could not place it.
	durationsLabelRowUnplaced = -1
)

// DurationLabelBand is one direct-label band's verdict: how many of its samples
// carry no label, whether selection did not pick them or the collision rule
// could not place them. A nonzero count is drawn on the chart as a remainder, so
// a reader can never mistake the visible labels for all of them.
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

// planDurationLabels decides which samples get a direct label. Each band first
// narrows to its durationsLabelTopCount longest spans — the lane's text answers
// "where are the outliers", so it goes to superlatives rather than to whichever
// spans a walk packed first — then packs its own rows: walk the candidates left
// to right and give each the first row where its text touches nothing already
// placed there. A sample that is not selected, or that fits nowhere, is counted
// rather than drawn, and the renderer prints that count, so a reader can never
// mistake the visible labels for all of them.
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
// sentence needs room at row 0's right edge, but whether there IS a remainder is
// only known once placement has finished — so the first pass uses the full
// width, and only a pass that actually dropped something is redone with the
// reservation held back. A board with no remainder pays nothing for one.
func packDurationLabelBand(samples []DurationSample, bandName string, rangeStart time.Time, rangeEnd time.Time) DurationLabelBand {
	selected := selectDurationLabelCandidates(samples, bandName)
	band := placeDurationLabelBand(samples, bandName, selected, rangeStart, rangeEnd, false)
	if band.HiddenCount == 0 {
		return band
	}
	return placeDurationLabelBand(samples, bandName, selected, rangeStart, rangeEnd, true)
}

// selectDurationLabelCandidates picks one band's label candidates: at most
// durationsLabelTopCount sample indexes, by magnitude. The sort is stable over
// completion order, so equal spans keep their left-to-right precedence and the
// choice is deterministic on an archive full of identical durations.
func selectDurationLabelCandidates(samples []DurationSample, bandName string) map[int]bool {
	candidateIndexes := []int{}
	for index := range samples {
		if durationLabelBandOf(samples[index]) == bandName {
			candidateIndexes = append(candidateIndexes, index)
		}
	}
	sort.SliceStable(candidateIndexes, func(first, second int) bool {
		return math.Abs(samples[candidateIndexes[first]].WallMinutes) >
			math.Abs(samples[candidateIndexes[second]].WallMinutes)
	})
	if len(candidateIndexes) > durationsLabelTopCount {
		candidateIndexes = candidateIndexes[:durationsLabelTopCount]
	}
	selected := make(map[int]bool, len(candidateIndexes))
	for _, candidateIndex := range candidateIndexes {
		selected[candidateIndex] = true
	}
	return selected
}

// placeDurationLabelBand is one greedy left-to-right pass over a single band.
// Samples arrive in completion order, so their x positions are non-decreasing
// and a row's occupancy is fully described by how far right it already reaches.
func placeDurationLabelBand(
	samples []DurationSample,
	bandName string,
	selected map[int]bool,
	rangeStart time.Time,
	rangeEnd time.Time,
	reserveRemainder bool,
) DurationLabelBand {
	occupiedTo := make([]float64, durationsLabelRowCount)
	for rowIndex := range occupiedTo {
		occupiedTo[rowIndex] = math.Inf(-1)
	}

	band := DurationLabelBand{}
	for index := range samples {
		if durationLabelBandOf(samples[index]) != bandName {
			continue
		}
		samples[index].LabelRow = durationsLabelRowUnplaced
		samples[index].LabelAnchor = ""
		if !selected[index] {
			band.HiddenCount++
			continue
		}

		markX := durationLabelPlotX(samples[index].CompletionTime, rangeStart, rangeEnd)
		textWidth := durationLabelWidthUnits(samples[index])
		for rowIndex := 0; rowIndex < durationsLabelRowCount; rowIndex++ {
			rowRightEdge := durationsPlotWidthUnits
			if reserveRemainder && rowIndex == durationsLabelRowCount-1 {
				rowRightEdge -= durationsLabelRemainderReserveUnits
			}
			// Anchoring BEFORE the mark is tried first, and the preference is a
			// packing decision rather than a stylistic one: this walk moves left
			// to right, so a label drawn to the mark's left reuses space the walk
			// has already passed, while one drawn to its right consumes space the
			// next mark still needs. Anchoring after the mark is the fallback that
			// keeps the leftmost sample labellable, since its own label would
			// otherwise start off-plot. Neither is an alternation: both are
			// geometric decisions about one specific label.
			spanLeft, spanRight, anchorFits := durationLabelSpan(markX, textWidth, "end", rowRightEdge)
			anchor := "end"
			if !anchorFits || spanLeft < occupiedTo[rowIndex]+durationsLabelSeparationUnits {
				spanLeft, spanRight, anchorFits = durationLabelSpan(markX, textWidth, "start", rowRightEdge)
				anchor = "start"
			}
			if !anchorFits || spanLeft < occupiedTo[rowIndex]+durationsLabelSeparationUnits {
				continue
			}
			samples[index].LabelRow = rowIndex
			samples[index].LabelAnchor = anchor
			occupiedTo[rowIndex] = spanRight
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
