package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

// durationTicket builds one archived-REQ ticket for the aggregation tests.
func durationTicket(requestId string, route string, claimedAt string, completedAt string) *RequestTicket {
	return &RequestTicket{
		RequestId:   requestId,
		Status:      "completed",
		Route:       route,
		ClaimedAt:   claimedAt,
		CompletedAt: completedAt,
	}
}

func findDurationDay(t *testing.T, aggregate DurationAggregate, dayKey string) DurationDay {
	t.Helper()
	for _, day := range aggregate.Days {
		if day.DayKey == dayKey {
			return day
		}
	}
	t.Fatalf("no aggregated day for %s", dayKey)
	return DurationDay{}
}

func findDurationSample(t *testing.T, aggregate DurationAggregate, requestId string) DurationSample {
	t.Helper()
	for _, sample := range aggregate.Samples {
		if sample.RequestId == requestId {
			return sample
		}
	}
	t.Fatalf("no sample for %s", requestId)
	return DurationSample{}
}

// A REQ contributes a sample only when it reached terminal success AND both
// stamps parse. Anything else is absent from the view rather than plotted at
// zero, because a zero would read as "instant" instead of "unmeasured".
func TestDurationAggregateAdmitsOnlyMeasurableCompletedRequests(t *testing.T) {
	tickets := []*RequestTicket{
		durationTicket("REQ-001", "A", "2026-07-01T10:00:00Z", "2026-07-01T10:05:00Z"),
		durationTicket("REQ-002", "B", "", "2026-07-01T11:00:00Z"),           // no claim stamp
		durationTicket("REQ-003", "B", "2026-07-01T12:00:00Z", ""),           // no completion stamp
		durationTicket("REQ-004", "C", "2026-07-01T13:00:00Z", "not-a-time"), // unparseable
		{RequestId: "REQ-005", Status: "pending", ClaimedAt: "2026-07-01T14:00:00Z", CompletedAt: "2026-07-01T14:30:00Z"},
		{RequestId: "REQ-006", Status: "cancelled", ClaimedAt: "2026-07-01T15:00:00Z", CompletedAt: "2026-07-01T15:30:00Z"},
		{RequestId: "REQ-007", Status: "completed-with-issues", Route: "C", ClaimedAt: "2026-07-01T16:00:00Z", CompletedAt: "2026-07-01T16:20:00Z"},
	}

	aggregate := buildDurationAggregate(tickets)

	if len(aggregate.Samples) != 2 {
		t.Fatalf("expected 2 measurable samples, got %d", len(aggregate.Samples))
	}
	if aggregate.Samples[0].RequestId != "REQ-001" || aggregate.Samples[1].RequestId != "REQ-007" {
		t.Fatalf("samples are not the two measurable terminal-success REQs in completion order: %+v", aggregate.Samples)
	}
	if aggregate.Samples[1].Route != "C" {
		t.Fatalf("completed-with-issues is a terminal success and must keep its route, got %q", aggregate.Samples[1].Route)
	}
}

// The read-time rule is what stops one paused session inventing a five-hour day.
// This pins the rule as APPLIED, not merely present: the same day's naive median
// is asserted alongside the ruled one so a regression that drops the rule fails
// on a value, not on a missing symbol.
func TestDurationDayMedianAppliesTheReadTimeOutlierRule(t *testing.T) {
	tickets := []*RequestTicket{
		// 2.5 minutes — an ordinary span.
		durationTicket("REQ-101", "A", "2026-07-31T08:00:00Z", "2026-07-31T08:02:31Z"),
		// 655.2 minutes — over the four-hour ceiling, so an assumed pause.
		durationTicket("REQ-064", "C", "2026-07-31T09:00:00Z", "2026-07-31T19:55:12Z"),
	}

	aggregate := buildDurationAggregate(tickets)
	day := findDurationDay(t, aggregate, "2026-07-31")

	if !day.HasMedian {
		t.Fatal("the day has one in-rule sample, so it must carry a median")
	}
	if math.Abs(day.MedianMinutes-2.5167) > 0.01 {
		t.Fatalf("day median should be the in-rule sample's 2.5 min, got %.4f", day.MedianMinutes)
	}
	if day.KeptCount != 1 {
		t.Fatalf("expected 1 sample inside the rule, got %d", day.KeptCount)
	}
	// Panel C counts everything, rule-excluded or not — it is a count, not a duration.
	if day.CompletedCount != 2 {
		t.Fatalf("expected both REQs counted for the day, got %d", day.CompletedCount)
	}

	// Without the rule the same day would report 328.9 — the number the rule exists
	// to prevent. Asserting the gap is what proves the rule ran.
	naiveMedian, _ := medianOf([]float64{2.5167, 655.2})
	if math.Abs(naiveMedian-328.9) > 0.5 {
		t.Fatalf("test fixture drifted: naive median should be ~328.9, got %.4f", naiveMedian)
	}
	if math.Abs(day.MedianMinutes-naiveMedian) < 1 {
		t.Fatal("the day median matches the unruled median — the read-time rule is not being applied")
	}

	excluded := findDurationSample(t, aggregate, "REQ-064")
	if excluded.DayMedianExclusion != "paused" {
		t.Fatalf("an over-ceiling span must be excluded as a pause, got %q", excluded.DayMedianExclusion)
	}
	if math.Abs(excluded.WallMinutes-655.2) > 0.01 {
		t.Fatalf("panel A must still carry the raw span, got %.4f", excluded.WallMinutes)
	}
}

// A reversed span (completed before claimed) is the anomaly class the board
// exists to surface. It must survive aggregation raw and negative — never
// clamped to zero — and must stay out of the day medians. The archive carries
// none today, which is exactly why this is pinned by a fixture.
func TestDurationAggregateKeepsReversedSpansRawAndOutOfMedians(t *testing.T) {
	tickets := []*RequestTicket{
		durationTicket("REQ-091", "B", "2026-08-02T12:00:00Z", "2026-08-02T11:30:00Z"), // −30 min
		durationTicket("REQ-092", "A", "2026-08-02T13:00:00Z", "2026-08-02T13:10:00Z"), // +10 min
	}

	aggregate := buildDurationAggregate(tickets)
	reversed := findDurationSample(t, aggregate, "REQ-091")

	if reversed.WallMinutes >= 0 {
		t.Fatalf("a reversed span must stay negative, got %.4f", reversed.WallMinutes)
	}
	if math.Abs(reversed.WallMinutes+30) > 0.01 {
		t.Fatalf("expected the raw −30 min span, got %.4f", reversed.WallMinutes)
	}
	if reversed.DayMedianExclusion != "reversed" {
		t.Fatalf("a reversed span must be excluded as a broken stamp, got %q", reversed.DayMedianExclusion)
	}

	day := findDurationDay(t, aggregate, "2026-08-02")
	if day.KeptCount != 1 || math.Abs(day.MedianMinutes-10) > 0.01 {
		t.Fatalf("the day median must be the one sound sample's 10 min over 1 sample, got %.4f over %d", day.MedianMinutes, day.KeptCount)
	}
	if day.CompletedCount != 2 {
		t.Fatalf("the day count must include the reversed REQ, got %d", day.CompletedCount)
	}
}

// A day whose only samples are rule-excluded has no median. That is a distinct
// state from a median of zero, and the panel must be able to tell them apart.
func TestDurationDayWithOnlyExcludedSamplesHasNoMedian(t *testing.T) {
	aggregate := buildDurationAggregate([]*RequestTicket{
		durationTicket("REQ-201", "C", "2026-08-03T01:00:00Z", "2026-08-03T09:00:00Z"), // 8h — paused
	})

	day := findDurationDay(t, aggregate, "2026-08-03")
	if day.HasMedian {
		t.Fatal("a day with no in-rule samples must report no median rather than zero")
	}
	if day.CompletedCount != 1 || day.KeptCount != 0 {
		t.Fatalf("expected 1 completed and 0 kept, got %d and %d", day.CompletedCount, day.KeptCount)
	}
}

// Days are chronological and gaps are real: the aggregate reports only active
// days, and the view plots them on a linear calendar axis so idle stretches stay
// visible. Compressing to one column per active day would destroy the cadence
// answer the view exists to give.
func TestDurationDaysAreChronologicalAndOnlyActiveDaysAppear(t *testing.T) {
	aggregate := buildDurationAggregate([]*RequestTicket{
		durationTicket("REQ-301", "A", "2026-08-10T10:00:00Z", "2026-08-10T10:05:00Z"),
		durationTicket("REQ-302", "A", "2026-06-01T10:00:00Z", "2026-06-01T10:05:00Z"),
		durationTicket("REQ-303", "A", "2026-07-04T10:00:00Z", "2026-07-04T10:05:00Z"),
	})

	if len(aggregate.Days) != 3 {
		t.Fatalf("expected exactly the 3 active days, got %d", len(aggregate.Days))
	}
	wantOrder := []string{"2026-06-01", "2026-07-04", "2026-08-10"}
	for index, dayKey := range wantOrder {
		if aggregate.Days[index].DayKey != dayKey {
			t.Fatalf("day %d should be %s, got %s", index, dayKey, aggregate.Days[index].DayKey)
		}
		if !aggregate.Days[index].DayTime.Equal(mustParseDay(t, dayKey)) {
			t.Fatalf("%s must carry its midnight-UTC instant for the shared calendar axis", dayKey)
		}
	}
}

func mustParseDay(t *testing.T, dayKey string) time.Time {
	t.Helper()
	parsed, parseError := time.Parse("2006-01-02", dayKey)
	if parseError != nil {
		t.Fatalf("parse %s: %v", dayKey, parseError)
	}
	return parsed.UTC()
}

// The live archive is the corpus the view actually renders. Past days are
// immutable — no future REQ can complete on 2026-07-31 — so pinning two of them
// catches a regression in the rule without going stale as new work lands.
func TestLiveArchiveDurationsMatchTheCalibratedFigures(t *testing.T) {
	board := liveBoard(t)
	aggregate := buildDurationAggregate(board.AllRequests)

	if len(aggregate.Samples) < 195 {
		t.Fatalf("the archive carried 195 measurable samples at capture and only grows; got %d", len(aggregate.Samples))
	}

	pausedDay := findDurationDay(t, aggregate, "2026-07-31")
	if math.Abs(pausedDay.MedianMinutes-2.5) > 0.05 {
		t.Fatalf("2026-07-31 must report the ruled median 2.5 min (REQ-064's 655 min excluded), got %.4f", pausedDay.MedianMinutes)
	}
	if pausedDay.CompletedCount != 2 || pausedDay.KeptCount != 1 {
		t.Fatalf("2026-07-31 should be 2 completed with 1 inside the rule, got %d and %d", pausedDay.CompletedCount, pausedDay.KeptCount)
	}

	busiestDay := findDurationDay(t, aggregate, "2026-08-15")
	if busiestDay.CompletedCount != 25 {
		t.Fatalf("2026-08-15 should carry 25 completed REQs, got %d", busiestDay.CompletedCount)
	}
	if math.Abs(busiestDay.MedianMinutes-19.6) > 0.1 {
		t.Fatalf("2026-08-15 should report a median of 19.6 min, got %.4f", busiestDay.MedianMinutes)
	}
}

// denseOverflowTickets builds a board whose overflow band is saturated: every
// REQ ran well over the ceiling and they all completed inside a two-day window,
// so consecutive marks sit far closer together than their label text is wide.
// The live archive cannot express this — it carries three overflow samples out
// of ~205 — which is why this fixture is synthetic rather than pinned.
func denseOverflowTickets(sampleCount int) []*RequestTicket {
	windowStart := time.Date(2026, 5, 4, 6, 0, 0, 0, time.UTC)
	windowLength := 48 * time.Hour
	tickets := make([]*RequestTicket, 0, sampleCount)
	for sampleIndex := 0; sampleIndex < sampleCount; sampleIndex++ {
		completedAt := windowStart.Add(time.Duration(sampleIndex) * windowLength / time.Duration(sampleCount-1))
		claimedAt := completedAt.Add(-95 * time.Minute)
		tickets = append(tickets, durationTicket(
			fmt.Sprintf("REQ-%03d", 400+sampleIndex),
			"C",
			claimedAt.Format(time.RFC3339),
			completedAt.Format(time.RFC3339),
		))
	}
	return tickets
}

// placedDurationLabel is one drawn label's occupied x-interval together with the
// sample it belongs to: the geometry an overlap check needs, plus the magnitude
// a "was this span passed over for a shorter one" check needs.
type placedDurationLabel struct {
	requestId string
	magnitude float64
	spanLeft  float64
	spanRight float64
}

// placedDurationLabelsByRow returns one band's drawn labels grouped by text row.
// It is the tests' single definition of where a drawn label actually sits, so no
// assertion below re-derives that geometry for itself.
func placedDurationLabelsByRow(aggregate DurationAggregate, bandName string) map[int][]placedDurationLabel {
	rangeStart, rangeEnd, _ := durationLabelTimeRange(aggregate.Samples)
	placedByRow := map[int][]placedDurationLabel{}
	for _, sample := range aggregate.Samples {
		if durationLabelBandOf(sample) != bandName || sample.LabelRow == durationsLabelRowUnplaced {
			continue
		}
		markX := durationLabelPlotX(sample.CompletionTime, rangeStart, rangeEnd)
		textWidth := durationLabelWidthUnits(sample)
		spanLeft, spanRight, _ := durationLabelSpan(markX, textWidth, sample.LabelAnchor, durationsPlotWidthUnits)
		placedByRow[sample.LabelRow] = append(placedByRow[sample.LabelRow], placedDurationLabel{
			requestId: sample.RequestId,
			magnitude: math.Abs(sample.WallMinutes),
			spanLeft:  spanLeft,
			spanRight: spanRight,
		})
	}
	return placedByRow
}

// labelledDurationSpans returns each placed label's occupied x-interval, per row,
// for one band — the geometry an overlap check needs.
func labelledDurationSpans(aggregate DurationAggregate, bandName string) map[int][][2]float64 {
	spansByRow := map[int][][2]float64{}
	for rowIndex, placedLabels := range placedDurationLabelsByRow(aggregate, bandName) {
		for _, placed := range placedLabels {
			spansByRow[rowIndex] = append(spansByRow[rowIndex], [2]float64{placed.spanLeft, placed.spanRight})
		}
	}
	return spansByRow
}

// The defect this pins: the renderer labelled every overflow sample and chose
// its slot from the sample's array index, so N overflow marks produced N labels
// in four recycled slots — an unreadable blob at any real density. Placement is
// now a geometric decision made once, on the Go side, and every label it cannot
// fit is counted rather than dropped in silence.
func TestDenseOverflowLabelsStayBoundedAndNeverOverlap(t *testing.T) {
	const overflowSampleCount = 40
	aggregate := buildDurationAggregate(denseOverflowTickets(overflowSampleCount))

	if len(aggregate.Samples) != overflowSampleCount {
		t.Fatalf("fixture produced %d samples, want %d", len(aggregate.Samples), overflowSampleCount)
	}
	labelledCount := 0
	for _, sample := range aggregate.Samples {
		if durationLabelBandOf(sample) != "overflow" {
			t.Fatalf("%s is not in the overflow band; fixture is wrong", sample.RequestId)
		}
		if sample.LabelRow != durationsLabelRowUnplaced {
			labelledCount++
		}
	}

	// (a) Bounded by what the lane can physically hold. The narrowest label in
	// this fixture, plus the gutter it needs, gives the tightest legal pitch;
	// two rows of that is the ceiling no correct placement can exceed.
	narrowestLabelWidth := math.Inf(1)
	for _, sample := range aggregate.Samples {
		narrowestLabelWidth = math.Min(narrowestLabelWidth, durationLabelWidthUnits(sample))
	}
	physicalCapacity := durationsLabelRowCount *
		int(durationsPlotWidthUnits/(narrowestLabelWidth+durationsLabelSeparationUnits)+1)
	if labelledCount > physicalCapacity {
		t.Fatalf("placed %d labels, but two rows of %.0f units hold at most %d",
			labelledCount, durationsPlotWidthUnits, physicalCapacity)
	}

	// (b) No two labels on a row may touch.
	for rowIndex, spans := range labelledDurationSpans(aggregate, "overflow") {
		sort.Slice(spans, func(earlier, later int) bool { return spans[earlier][0] < spans[later][0] })
		for spanIndex := 1; spanIndex < len(spans); spanIndex++ {
			previousRight := spans[spanIndex-1][1]
			currentLeft := spans[spanIndex][0]
			if currentLeft < previousRight+durationsLabelSeparationUnits {
				t.Fatalf("row %d: label starting at %.1f overlaps the one ending at %.1f",
					rowIndex, currentLeft, previousRight)
			}
		}
	}

	// Every sample the rule dropped is accounted for, so the chart can say so.
	if labelledCount+aggregate.OverflowLabels.HiddenCount != overflowSampleCount {
		t.Fatalf("%d labelled + %d hidden = %d, want %d — dropped labels are unaccounted for",
			labelledCount, aggregate.OverflowLabels.HiddenCount,
			labelledCount+aggregate.OverflowLabels.HiddenCount, overflowSampleCount)
	}
	if aggregate.OverflowLabels.HiddenCount == 0 {
		t.Fatal("a 40-sample two-day overflow burst must not fit entirely; the fixture is not dense enough")
	}
}

// durationsWindowStart is the fixture window both label fixtures share, so the
// two aggregates below plot against an identical x-domain and their placements
// are comparable rather than merely similar.
var durationsWindowStart = time.Date(2026, 5, 4, 6, 0, 0, 0, time.UTC)

const durationsWindowLength = 48 * time.Hour

// reversedStampTickets pins three broken-stamp REQs to the window's start,
// middle and end, so they exist at the same x positions with or without an
// overflow burst beside them.
func reversedStampTickets() []*RequestTicket {
	tickets := []*RequestTicket{}
	for offsetIndex, offset := range []time.Duration{0, durationsWindowLength / 2, durationsWindowLength} {
		completedAt := durationsWindowStart.Add(offset)
		claimedAt := completedAt.Add(20 * time.Minute) // completed before claimed: a reversed stamp
		tickets = append(tickets, durationTicket(
			fmt.Sprintf("REQ-9%02d", offsetIndex),
			"A",
			claimedAt.Format(time.RFC3339),
			completedAt.Format(time.RFC3339),
		))
	}
	return tickets
}

func durationLabelPlacementByRequestId(aggregate DurationAggregate) map[string][2]string {
	placement := map[string][2]string{}
	for _, sample := range aggregate.Samples {
		placement[sample.RequestId] = [2]string{strconv.Itoa(sample.LabelRow), sample.LabelAnchor}
	}
	return placement
}

// The old rule ran ONE counter across both bands, so a reversed label's slot was
// decided by how many overflow samples happened to precede it — two visually
// unrelated groups sharing a phase. Placement is now per band: saturating the
// overflow lane must leave the reversed band's answer byte-identical.
func TestReversedLabelPlacementIsIndependentOfOverflowDensity(t *testing.T) {
	reversedOnly := buildDurationAggregate(reversedStampTickets())
	withOverflowBurst := buildDurationAggregate(append(denseOverflowTickets(40), reversedStampTickets()...))

	alonePlacement := durationLabelPlacementByRequestId(reversedOnly)
	besidePlacement := durationLabelPlacementByRequestId(withOverflowBurst)
	for _, ticket := range reversedStampTickets() {
		if alonePlacement[ticket.RequestId] != besidePlacement[ticket.RequestId] {
			t.Fatalf("%s placed %v alone but %v beside a dense overflow lane; the bands share state",
				ticket.RequestId, alonePlacement[ticket.RequestId], besidePlacement[ticket.RequestId])
		}
	}
	if withOverflowBurst.OverflowLabels.HiddenCount == 0 {
		t.Fatal("the overflow band must be saturated for the comparison above to mean anything")
	}
	if reversedOnly.ReversedLabels.HiddenCount != 0 {
		t.Fatalf("three well-separated reversed stamps must all fit; %d were hidden",
			reversedOnly.ReversedLabels.HiddenCount)
	}
	if withOverflowBurst.ReversedLabels.HiddenCount != 0 {
		t.Fatalf("a saturated overflow lane must not consume the reversed band's rows; %d were hidden",
			withOverflowBurst.ReversedLabels.HiddenCount)
	}
}

// rendererNumericConstant reads one `var NAME = NUMBER;` out of an embedded view
// fragment. Parity tests and JavaScript behavior probes read the renderer's own
// values through it, so neither can drift into asserting against a hand-copied
// number. It is why every constant a test needs is written as a plain literal
// rather than as an expression.
func rendererNumericConstant(t *testing.T, assetPath string, constantName string) float64 {
	t.Helper()
	rendererText, readError := embeddedWebAssets.ReadFile(assetPath)
	if readError != nil {
		t.Fatalf("read %s: %v", assetPath, readError)
	}
	pattern := regexp.MustCompile(`(?m)^\s*var ` + regexp.QuoteMeta(constantName) + ` = (-?[0-9.]+);`)
	match := pattern.FindSubmatch(rendererText)
	if match == nil {
		t.Fatalf("%s declares no numeric constant %s", assetPath, constantName)
	}
	value, parseError := strconv.ParseFloat(string(match[1]), 64)
	if parseError != nil {
		t.Fatalf("parse %s: %v", constantName, parseError)
	}
	return value
}

func durationsRendererConstant(t *testing.T, constantName string) float64 {
	t.Helper()
	return rendererNumericConstant(t, "web/board-durations.js", constantName)
}

// Descent below the baseline for the renderer's 11px label face — the number
// REQ-226 used to describe a label's text box. The ascent is the renderer's own
// DURATIONS_LABEL_TEXT_ASCENT constant, read below; the descent exists only in
// this test's geometric question, so it lives here.
const durationsLabelTextDescentUnits = 2.0

// ---- the rendered face, measured ------------------------------------------
//
// Every number below comes from the browser, because the face is the browser's
// answer and no amount of arithmetic over the constants under test can produce
// it — computing a guarantee from the constant it is meant to judge is circular.
// They are recorded here so a Go test can hold the model to the face without a
// browser in the loop, and each is rounded AWAY from the model (up), so a
// passing assertion can never be an artefact of the rounding.
//
// Procedure, reproducible from any board directory `queue-kanban generate` wrote:
// load index.html, activate the Durations view, append an SVG <text> carrying
// class "durations-mark-label" to the chart's own <svg>, and read
// getComputedTextLength() and getBBox() off it, at the class's declared 11px
// over the board's --font-sans stack; the SVG is a fixed viewBox at width:100%,
// so user units are zoom- and window-independent.
//
// A measured face is PER-BROWSER — the same probe returns different numbers on
// different Chromium builds, because the face resolves through the build and
// the machine's font stack. So each constant's own doc comment names the build
// its number was taken on (TestDurationsMeasuredConstantsNameTheirChromiumBuild
// enforces this for every durationsMeasured constant in the package), a
// re-measurement on another build may only RAISE a constant, never lower it,
// and any new browser measurement takes the durationsMeasured prefix so the
// same test covers it.

// The SUPREMUM of units per character over the whole label space — not the worst
// case of a sweep. The distinction is the entire point, and getting it wrong is
// what this constant's own history records: 6.2 was fitted to a plausible string,
// and 6.71 came from a sweep that varied REQ-id digits but held the hour count at
// two, which bounded nothing.
//
// The space is infinite in two independent directions (id digit count, hour
// count) and per-character width rises with length in both, because the narrow
// fixed characters are a fixed cost diluted by every added digit. So the sweep
// has to be closed by an argument, not by more sampling. The argument:
//
//  1. A label is "REQ-" + id + " " + duration, and only DIGITS can repeat without
//     limit — every wide fixed character (R, E, Q, m, h, −) appears at most twice.
//  2. So as a label grows, the fixed characters amortize away and per-character
//     width converges to that of a pure digit run, from below.
//  3. A pure run of the widest digit ("4") measures 7.1441 units per character.
//
// Both ends were then measured rather than assumed. Exhaustively over the bounded
// region — every digit, id lengths 1-40, hour counts 1-40, both duration forms,
// both signs, mixed-digit ids, 280 800 labels — the maximum is 7.0643. Pushing
// both unbounded parameters to the limit gives 7.1384 at 2 010 characters, 7.1411
// at 10 010 and 7.1417 at 40 010: rising, and still under 7.1441. A randomized
// search over 4 000 mixed-digit labels confirms a uniform "4" run is the worst
// case, so the exhaustive region's shape is right.
//
// Build: measured for REQ-241 on the Chromium bundled by its Playwright 1.59,
// recorded only as browser build chromium-1228 — no version number survives.
// On Chromium 141.0.7390.37 (Playwright 1.56.1, REQ-252) the same "4" run
// measures 6.9865 and the worst minus-bearing short label 6.8952 per character:
// smaller, so the larger recorded supremum stands.
const durationsMeasuredLabelWidthSupremumUnits = 7.1441

// How far above the supremum the model may sit. A width model must over-estimate,
// but unbounded over-estimation is not free — it silently drops labels the rows
// could have carried — so the pin is two-sided and this is the slack it allows.
const durationsLabelWidthModelSlackUnits = 0.25

// The same face's rendered line box — the height getBBox() reports for a
// .durations-mark-label <text>, constant whether or not the string carries
// descenders, because it is the line box and not the ink. It is exactly that
// face's ascent plus its descent, and generate_test.go already bounds both
// parts for the whole package, as durationsMeasuredMarkLabelAscentUnits and
// durationsMeasuredMarkLabelDescentUnits. The box gets a constant of its own
// anyway because the SUM of those two bounds is 13.3: each was rounded up from
// the build that maximised it, the maxima fall on different builds, and no
// build has ever drawn both together — so the sum over-bounds every face ever
// measured, while the pitch floor below needs a bound a real face can be held
// to.
//
// Sampled builds, ascent + descent = box:
//
//	chromium-1228 (Playwright 1.59, REQ-241)             10.4278 + 2.4064 = 12.8343
//	Chromium 146 (Playwright 1.59, REQ-242)              10.1853 + 2.7778 = 12.9631
//	Chromium 141.0.7390.37 (Playwright 1.56.1, REQ-252)  10.1853 + 2.7778 = 12.9631
//
// REQ-265 re-ran the procedure above on Chromium 141.0.7390.37 and read the
// same three numbers to four decimals. 12.97 is the largest sampled box,
// rounded up.
//
// What this constant is NOT is a supremum over the face space, and that gap
// matters more than the 0.13 units the REQ-265 raise moved it. The face is
// whatever --font-sans resolves to, and board.css ends that stack in the open
// generic sans-serif; every sample above is one Linux container's answer to it.
// The design's own cap is the row pitch: 13 units over an 11-unit font size is
// an (ascent+descent)/em ratio of 1.1818, and the largest ratio ever sampled is
// 1.1785 — 0.03 units of slack under the floor below. So this number is wrong,
// with the pitch wrong behind it, the first time the board renders somewhere
// --font-sans resolves to a face whose ratio clears 1.1818. Padding it toward
// that cap would buy nothing and hide exactly that, so it is not padded: the
// raise-only rule keeps it honest, and a measurement that exceeds it is meant
// to cost somebody an edit.
const durationsMeasuredLabelBoxHeightUnits = 12.97

// The 12px axis-title face's ascent is declared once for the whole package, in
// generate_test.go's measured-face block, because REQ-241 and REQ-242 each
// measured it independently and landed on different numbers. The clearance test
// below reads it from there — see that block for which measurement won and why.

// The defect this pins (REQ-241): durationsLabelCharacterWidthUnits was 6.2
// against a face whose labels reach 7.14 units per character, and its comment
// called that "deliberately generous" — a claimed safety margin pointing the
// wrong way. A per-character width model is only honest as an UPPER bound:
// over-estimating drops a label, which the remainder sentence counts, while
// under-estimating overprints one, which is the whole defect placement exists to
// prevent.
//
// The pin is two-sided on purpose. A one-sided "at least the measured worst case"
// assertion passed while the constant sat at 6.75 and the true supremum was
// 7.1441, and it would have kept passing if the constant were dropped back to the
// stale reference value — the assertion cannot distinguish "correct" from "equal
// to whatever number the last sweep happened to produce". The lower bound is the
// correctness one; the upper bound is what stops the constant being inflated to
// buy safety with labels nobody chose to spend.
func TestDurationsLabelWidthEstimateCoversTheRenderedFace(t *testing.T) {
	if durationsLabelCharacterWidthUnits < durationsMeasuredLabelWidthSupremumUnits {
		t.Fatalf("width model assumes %.4f units per character, but label text reaches %.4f in the rendered face — the estimate under-states the text it is placing",
			durationsLabelCharacterWidthUnits, durationsMeasuredLabelWidthSupremumUnits)
	}
	ceiling := durationsMeasuredLabelWidthSupremumUnits + durationsLabelWidthModelSlackUnits
	if durationsLabelCharacterWidthUnits > ceiling {
		t.Fatalf("width model assumes %.4f units per character against a supremum of %.4f — over-estimating by more than %.2f drops labels the rows could carry",
			durationsLabelCharacterWidthUnits, durationsMeasuredLabelWidthSupremumUnits, durationsLabelWidthModelSlackUnits)
	}
}

// The defect this pins (REQ-241): DURATIONS_LABEL_ROW_HEIGHT was 12 while the
// same file declared an 11-unit ascent and this file a 2-unit descent — a 13-unit
// box on a 12-unit pitch, so two label rows' boxes overlapped by 0.83 units in a
// live render. Nothing visibly collided, because a line box is padding rather
// than ink; what it cost was the ability to assert row-against-row separation at
// all, the way TestDurationsLabelRowsClearTheMarkBands asserts row-against-mark.
// The pitch must clear BOTH boxes: the one the code declares and the one the
// browser draws.
func TestDurationsLabelRowPitchClearsTheLabelTextBox(t *testing.T) {
	rowHeight := durationsRendererConstant(t, "DURATIONS_LABEL_ROW_HEIGHT")
	declaredBoxHeight := durationsRendererConstant(t, "DURATIONS_LABEL_TEXT_ASCENT") + durationsLabelTextDescentUnits
	if rowHeight < declaredBoxHeight {
		t.Fatalf("row pitch %.2f is smaller than the %.2f-unit text box the renderer declares (ascent + %.2f descent) — consecutive rows share vertical space",
			rowHeight, declaredBoxHeight, durationsLabelTextDescentUnits)
	}
	if rowHeight < durationsMeasuredLabelBoxHeightUnits {
		t.Fatalf("row pitch %.2f is smaller than the %.2f-unit line box the browser draws — consecutive rows share vertical space",
			rowHeight, durationsMeasuredLabelBoxHeightUnits)
	}
}

// The row pitch has a ceiling as well as a floor, and this is it. The reversed
// band's rows grow DOWNWARD into the gap above Panel B's title, so every unit
// added to the pitch is taken from that gap — in this assertion's own model the
// shipped pitch of 13 leaves 0.10 units between the last row's box bottom and
// the title's box top, and a pitch of 14 would put it through.
//
// The descent below is the package's one bound for the 11px
// .durations-mark-label face — durationsMeasuredMarkLabelDescentUnits, in
// generate_test.go, which the slowest-day annotation's clearance reads for the
// same face and the same reason. It is deliberately the DRAWN descent and not
// durationsLabelTextDescentUnits, the smaller box the renderer declares; a
// clearance question has to use what the browser puts on the page. Until
// REQ-265 this test read a second constant declared here at 2.41 while that one
// already stood at 2.8 — two numbers for one quantity of one face, in two files
// of one package, which is the same shape that made REQ-241's and REQ-242's
// title-face measurements collide. REQ-252 made every such constant name its
// build; the duplicate is gone so there is one bound per face and quantity for
// a build name to be attached to.
//
// That budget is PER-BROWSER, and the live headroom has measured very
// differently on different builds: REQ-241 read 1.364 units on its Chromium
// (recorded only as browser build chromium-1228, Playwright 1.59), REQ-242's
// review read 0.185 on Chromium 146 over byte-identical SVG, and REQ-252 read
// 1.111 on Chromium 141.0.7390.37 (Playwright 1.56.1) after REQ-248's
// day-anchoring. Anyone about to spend this budget must re-measure it on their
// own build first — the roughly 7x spread between recorded values is real, not
// a mistake.
//
// REQ-241 established the budget by measurement and then left it in prose,
// where nothing enforced it: the row-pitch floor above and this ceiling are the
// two halves of one constraint, and a future change that raises the pitch to
// clear a bigger face would satisfy the floor and silently eat the ceiling. The
// band's last row is the one that matters because the marks sit level with the
// first.
func TestDurationsLastLabelRowClearsPanelBTitle(t *testing.T) {
	reversedRowY := durationsRendererConstant(t, "DURATIONS_REVERSED_LABEL_ROW_Y")
	rowCount := durationsRendererConstant(t, "DURATIONS_LABEL_ROW_COUNT")
	rowHeight := durationsRendererConstant(t, "DURATIONS_LABEL_ROW_HEIGHT")
	panelBTitleY := durationsRendererConstant(t, "DURATIONS_MEDIAN_TITLE_Y")

	lastRowBoxBottom := reversedRowY + (rowCount-1)*rowHeight + durationsMeasuredMarkLabelDescentUnits
	panelBTitleBoxTop := panelBTitleY - durationsMeasuredAxisTitleAscentUnits
	if lastRowBoxBottom >= panelBTitleBoxTop {
		t.Fatalf("the reversed band's last label row ends at %.2f but Panel B's title starts at %.2f — the label rows have grown into the title",
			lastRowBoxBottom, panelBTitleBoxTop)
	}
}

// The defect this pins (REQ-231): REQ-226 stopped labels overprinting each
// other, but the mark band and the first label row still shared vertical space,
// so on a dense board the DOTS overprinted the text instead. Both bands' label
// rows must sit wholly clear of their band's marks, at any density — the rule is
// geometric, so it is asserted against the renderer's own constants the same way
// TestDurationLabelGeometryMatchesTheRenderer reads them.
func TestDurationsLabelRowsClearTheMarkBands(t *testing.T) {
	markRadius := durationsRendererConstant(t, "DURATIONS_BAND_MARK_RADIUS")
	rowCount := int(durationsRendererConstant(t, "DURATIONS_LABEL_ROW_COUNT"))
	rowHeight := durationsRendererConstant(t, "DURATIONS_LABEL_ROW_HEIGHT")
	textAscent := durationsRendererConstant(t, "DURATIONS_LABEL_TEXT_ASCENT")
	for _, band := range []struct {
		bandName         string
		markConstantName string
		rowConstantName  string
	}{
		{"overflow", "DURATIONS_LANE_MARK_Y", "DURATIONS_LANE_LABEL_ROW_Y"},
		{"reversed", "DURATIONS_BELOW_ZERO_Y", "DURATIONS_REVERSED_LABEL_ROW_Y"},
	} {
		markY := durationsRendererConstant(t, band.markConstantName)
		firstRowY := durationsRendererConstant(t, band.rowConstantName)
		for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
			baseline := firstRowY + float64(rowIndex)*rowHeight
			textTop := baseline - textAscent
			textBottom := baseline + durationsLabelTextDescentUnits
			if textBottom >= markY-markRadius && textTop <= markY+markRadius {
				t.Fatalf("%s band: row %d's text box [%.0f, %.0f] intersects the mark band [%.0f, %.0f] — a neighbouring dot can overprint the label",
					band.bandName, rowIndex, textTop, textBottom, markY-markRadius, markY+markRadius)
			}
		}
	}
}

// The defect this pins (REQ-252): a browser measurement is only arguable while
// it names the browser it came from. REQ-241 and REQ-242 measured the same 12px
// axis-title face on different Chromium builds, got 11.2300 and 12.0372, and
// declared the same constant in different files of one package — a collision
// git could not see, because the two edits never touched adjacent lines. The
// package's convention closes that class: a constant holding a browser's answer
// takes the durationsMeasured prefix, and its own doc comment names the
// Chromium build the number was taken on. This test enforces the second half
// over every Go file in the package. The prefix half stays convention — a
// browser number smuggled into a differently named constant is caught only by
// review, because no comment-reading test can know where a number came from.
func TestDurationsMeasuredConstantsNameTheirChromiumBuild(t *testing.T) {
	// Matches "Chromium 146", "Chromium 141.0.7390.37" and the Playwright
	// browser-build form "chromium-1228" — a browser name with no digits after
	// it (e.g. "a different Chromium") is not a build.
	chromiumBuildPattern := regexp.MustCompile(`(?i)chromium[ -]\d`)
	fileSet := token.NewFileSet()
	packageEntries, readDirError := os.ReadDir(".")
	if readDirError != nil {
		t.Fatalf("read package directory: %v", readDirError)
	}
	measuredConstantCount := 0
	for _, packageEntry := range packageEntries {
		if packageEntry.IsDir() || !strings.HasSuffix(packageEntry.Name(), ".go") {
			continue
		}
		parsedFile, parseError := parser.ParseFile(fileSet, packageEntry.Name(), nil, parser.ParseComments)
		if parseError != nil {
			t.Fatalf("parse %s: %v", packageEntry.Name(), parseError)
		}
		for _, topLevelDecl := range parsedFile.Decls {
			constDecl, isGenDecl := topLevelDecl.(*ast.GenDecl)
			if !isGenDecl || constDecl.Tok != token.CONST {
				continue
			}
			for _, constSpec := range constDecl.Specs {
				valueSpec, isValueSpec := constSpec.(*ast.ValueSpec)
				if !isValueSpec {
					continue
				}
				for _, constName := range valueSpec.Names {
					if !strings.HasPrefix(constName.Name, "durationsMeasured") {
						continue
					}
					measuredConstantCount++
					provenanceComment := ""
					if valueSpec.Doc != nil {
						provenanceComment = valueSpec.Doc.Text()
					} else if constDecl.Doc != nil {
						provenanceComment = constDecl.Doc.Text()
					}
					if !chromiumBuildPattern.MatchString(provenanceComment) {
						t.Errorf("%s: %s is a browser measurement whose doc comment names no Chromium build — the number cannot be re-derived or argued with on another machine",
							packageEntry.Name(), constName.Name)
					}
				}
			}
		}
	}
	if measuredConstantCount == 0 {
		t.Fatal("found no durationsMeasured constants in the package — the walk is broken or the naming convention changed; fix this test to scan the real declarations")
	}
}

// variedOverflowTickets is the dense fixture with a magnitude gradient: spans
// grow left to right, so the longest samples sit at the crowded right edge and
// a first-fit walk would spend both rows on the SHORT leftmost spans.
func variedOverflowTickets(sampleCount int) []*RequestTicket {
	tickets := make([]*RequestTicket, 0, sampleCount)
	for sampleIndex := 0; sampleIndex < sampleCount; sampleIndex++ {
		completedAt := durationsWindowStart.Add(time.Duration(sampleIndex) * durationsWindowLength / time.Duration(sampleCount-1))
		claimedAt := completedAt.Add(-time.Duration(65+7*sampleIndex) * time.Minute)
		tickets = append(tickets, durationTicket(
			fmt.Sprintf("REQ-%03d", 500+sampleIndex),
			"B",
			claimedAt.Format(time.RFC3339),
			completedAt.Format(time.RFC3339),
		))
	}
	return tickets
}

// durationLabelRowRightEdge mirrors placement's own right-edge rule: the band's
// last row gives up durationsLabelRemainderReserveUnits to the remainder
// sentence, and only when there IS a remainder to print. A check that ignored
// the reservation would call the reserved strip "free space".
func durationLabelRowRightEdge(rowIndex int, band DurationLabelBand) float64 {
	if band.HiddenCount > 0 && rowIndex == durationsLabelRowCount-1 {
		return durationsPlotWidthUnits - durationsLabelRemainderReserveUnits
	}
	return durationsPlotWidthUnits
}

// durationLabelSpansConflict reports whether two x-intervals on one row come
// closer than the separation the rule demands.
func durationLabelSpansConflict(firstLeft float64, firstRight float64, secondLeft float64, secondRight float64) bool {
	return firstLeft < secondRight+durationsLabelSeparationUnits &&
		secondLeft < firstRight+durationsLabelSeparationUnits
}

// assertDurationLabelPriority is the band-level statement of "labels go to the
// longest spans", written so it holds however many labels the rows end up
// carrying. For every sample that got none, every row/anchor slot its text could
// legally occupy must be blocked by a drawn label AT LEAST AS LONG as it. Two
// distinct defects fail it: a free slot means the rows stopped short of what they
// hold, and a slot blocked only by shorter labels means a short span took a long
// span's place, which is the left-edge first-fit sampling REQ-231 removed.
func assertDurationLabelPriority(t *testing.T, aggregate DurationAggregate, bandName string, band DurationLabelBand) {
	t.Helper()
	rangeStart, rangeEnd, _ := durationLabelTimeRange(aggregate.Samples)
	placedByRow := placedDurationLabelsByRow(aggregate, bandName)
	for _, sample := range aggregate.Samples {
		if durationLabelBandOf(sample) != bandName || sample.LabelRow != durationsLabelRowUnplaced {
			continue
		}
		markX := durationLabelPlotX(sample.CompletionTime, rangeStart, rangeEnd)
		textWidth := durationLabelWidthUnits(sample)
		magnitude := math.Abs(sample.WallMinutes)
		for rowIndex := 0; rowIndex < durationsLabelRowCount; rowIndex++ {
			rowRightEdge := durationLabelRowRightEdge(rowIndex, band)
			for _, anchor := range []string{"end", "start"} {
				spanLeft, spanRight, anchorFits := durationLabelSpan(markX, textWidth, anchor, rowRightEdge)
				if !anchorFits {
					continue
				}
				longestBlocker := math.Inf(-1)
				blockerId := ""
				for _, placed := range placedByRow[rowIndex] {
					if !durationLabelSpansConflict(spanLeft, spanRight, placed.spanLeft, placed.spanRight) {
						continue
					}
					if placed.magnitude > longestBlocker {
						longestBlocker = placed.magnitude
						blockerId = placed.requestId
					}
				}
				if blockerId == "" {
					t.Fatalf("%s band: %s (%.0f min) carries no label, but row %d anchor %q is free at [%.0f, %.0f] — the rows stopped short of what they hold",
						bandName, sample.RequestId, sample.WallMinutes, rowIndex, anchor, spanLeft, spanRight)
				}
				if longestBlocker < magnitude {
					t.Fatalf("%s band: %s (%.0f min) was passed over on row %d anchor %q, blocked only by shorter labels (longest blocker %s at %.0f min) — labels are not going to the longest spans",
						bandName, sample.RequestId, sample.WallMinutes, rowIndex, anchor, blockerId, longestBlocker)
				}
			}
		}
	}
}

func labelledDurationSampleCount(aggregate DurationAggregate, bandName string) int {
	labelledCount := 0
	for _, sample := range aggregate.Samples {
		if durationLabelBandOf(sample) == bandName && sample.LabelRow != durationsLabelRowUnplaced {
			labelledCount++
		}
	}
	return labelledCount
}

// The lane's text answers "where are the outliers", so its labels must go to the
// longest spans, not to whichever spans a left-to-right walk packed first — the
// rule REQ-231 introduced. REQ-231 spelt it as a fixed top-6 candidate set, which
// REQ-237 replaced with a descending-magnitude walk that keeps going until the
// rows are full; the invariant that survives both is the priority one, so this
// asserts that rather than a candidate count.
func TestOverflowLabelsGoToTheLongestSpans(t *testing.T) {
	const overflowSampleCount = 40
	aggregate := buildDurationAggregate(variedOverflowTickets(overflowSampleCount))

	assertDurationLabelPriority(t, aggregate, "overflow", aggregate.OverflowLabels)

	labelledCount := labelledDurationSampleCount(aggregate, "overflow")
	if labelledCount == 0 {
		t.Fatal("the rule must still label something on a dense board")
	}
	if labelledCount+aggregate.OverflowLabels.HiddenCount != overflowSampleCount {
		t.Fatalf("%d labelled + %d hidden ≠ %d — unlabelled samples must join the drawn remainder",
			labelledCount, aggregate.OverflowLabels.HiddenCount, overflowSampleCount)
	}
}

// The defect this pins (REQ-237): selection took the six longest spans and
// placement then dropped whichever of them collided, with nothing offering the
// freed space to the seventh. On a board where magnitude correlates with
// completion time every candidate lands in the same crowded corner, so the two
// rows drew 2 labels where they physically hold ten times that and the remainder
// count carried all 38. The rows are a fixed budget; a dropped span's space
// belongs to the next-longest span, not to nobody.
func TestClusteredOverflowLabelsFillBothLabelRows(t *testing.T) {
	const overflowSampleCount = 40
	aggregate := buildDurationAggregate(variedOverflowTickets(overflowSampleCount))

	// The fixture spaces its samples evenly across the plot, so one row's greedy
	// capacity is arithmetic rather than a guess: a label needs its own width plus
	// the separation clear, which is every Nth sample at that pitch. Two rows must
	// carry at least what a single row holds — a deliberately slack floor, since
	// the second row interleaves with the first.
	widestLabelWidth := 0.0
	for _, sample := range aggregate.Samples {
		if durationLabelBandOf(sample) == "overflow" {
			widestLabelWidth = math.Max(widestLabelWidth, durationLabelWidthUnits(sample))
		}
	}
	sampleStrideUnits := durationsPlotWidthUnits / float64(overflowSampleCount-1)
	samplesPerLabel := math.Ceil((widestLabelWidth + durationsLabelSeparationUnits) / sampleStrideUnits)
	singleRowCapacity := 1 + int(float64(overflowSampleCount-1)/samplesPerLabel)

	labelledCount := labelledDurationSampleCount(aggregate, "overflow")
	if labelledCount < singleRowCapacity {
		t.Fatalf("%d labels drawn across %d rows, but one row alone holds %d at this fixture's pitch — collided spans are not being backfilled",
			labelledCount, durationsLabelRowCount, singleRowCapacity)
	}
}

// Placement decides in the renderer's user-unit space, so the two files agree on
// that space or every label is sized against a plot it does not land on. The
// numbers are duplicated by necessity — one side computes, the other draws — and
// this is what stops the duplicate becoming a divergence.
func TestDurationLabelGeometryMatchesTheRenderer(t *testing.T) {
	rendererPlotWidth := durationsRendererConstant(t, "DURATIONS_VIEW_WIDTH") -
		durationsRendererConstant(t, "DURATIONS_MARGIN_LEFT") -
		durationsRendererConstant(t, "DURATIONS_MARGIN_RIGHT")
	if rendererPlotWidth != durationsPlotWidthUnits {
		t.Fatalf("renderer plot width = %.1f, placement assumes %.1f", rendererPlotWidth, durationsPlotWidthUnits)
	}
	for _, pair := range []struct {
		rendererName string
		goValue      float64
	}{
		{"DURATIONS_CEILING_MINUTES", durationsOverflowCeilingMinutes},
		{"DURATIONS_LABEL_ROW_COUNT", durationsLabelRowCount},
		{"DURATIONS_LABEL_GAP", durationsLabelGapUnits},
	} {
		if rendererValue := durationsRendererConstant(t, pair.rendererName); rendererValue != pair.goValue {
			t.Fatalf("renderer %s = %v, placement assumes %v", pair.rendererName, rendererValue, pair.goValue)
		}
	}
}
