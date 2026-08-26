package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
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

// calibratedLiveArchiveFindings returns one line per pinned figure the aggregate
// disagrees with, and nothing when every pin holds. The assertions live here
// rather than inline in the test so a second test can prove they still BITE —
// a pinned check that silently stopped biting reads exactly like a passing one.
//
// The pins themselves are the point: past days are immutable — no future REQ can
// complete on 2026-07-31 — so they catch a regression in the duration rule
// without going stale as new work lands.
func calibratedLiveArchiveFindings(aggregate DurationAggregate) []string {
	findings := []string{}

	if len(aggregate.Samples) < 195 {
		findings = append(findings,
			fmt.Sprintf("the archive carried 195 measurable samples at capture and only grows; got %d", len(aggregate.Samples)))
	}

	pausedDay, pausedDayFound := lookupDurationDay(aggregate, "2026-07-31")
	switch {
	case !pausedDayFound:
		findings = append(findings, "no aggregated day for 2026-07-31")
	default:
		if math.Abs(pausedDay.MedianMinutes-2.5) > 0.05 {
			findings = append(findings,
				fmt.Sprintf("2026-07-31 must report the ruled median 2.5 min (REQ-064's 655 min excluded), got %.4f", pausedDay.MedianMinutes))
		}
		if pausedDay.CompletedCount != 2 || pausedDay.KeptCount != 1 {
			findings = append(findings,
				fmt.Sprintf("2026-07-31 should be 2 completed with 1 inside the rule, got %d and %d", pausedDay.CompletedCount, pausedDay.KeptCount))
		}
	}

	busiestDay, busiestDayFound := lookupDurationDay(aggregate, "2026-08-15")
	switch {
	case !busiestDayFound:
		findings = append(findings, "no aggregated day for 2026-08-15")
	default:
		if busiestDay.CompletedCount != 25 {
			findings = append(findings,
				fmt.Sprintf("2026-08-15 should carry 25 completed REQs, got %d", busiestDay.CompletedCount))
		}
		if math.Abs(busiestDay.MedianMinutes-19.6) > 0.1 {
			findings = append(findings,
				fmt.Sprintf("2026-08-15 should report a median of 19.6 min, got %.4f", busiestDay.MedianMinutes))
		}
	}

	return findings
}

// lookupDurationDay is findDurationDay without a *testing.T — a missing day is a
// finding for the caller to report, not a fatal here.
func lookupDurationDay(aggregate DurationAggregate, dayKey string) (DurationDay, bool) {
	for _, day := range aggregate.Days {
		if day.DayKey == dayKey {
			return day, true
		}
	}
	return DurationDay{}, false
}

// The live archive is the corpus the view actually renders — which is why this
// runs against it rather than a fixture, and why it applies only where that
// corpus is ours (suiteCheckoutSkipReason, board_live_test.go).
func TestLiveArchiveDurationsMatchTheCalibratedFigures(t *testing.T) {
	repoRoot := liveRepoRoot(t)
	if skipReason := suiteCheckoutSkipReason(repoRoot); skipReason != "" {
		t.Skip(skipReason)
	}
	board := liveBoardAt(t, repoRoot)

	for _, finding := range calibratedLiveArchiveFindings(buildDurationAggregate(board.AllRequests)) {
		t.Error(finding)
	}
}

// TestLiveArchiveCalibrationRunsInASuiteCheckoutAndSkipsElsewhere asserts BOTH
// halves in one test, the way TestReleaseProbesRunInASuiteCheckoutAndAreNot​ApplicableElsewhere
// does for the release probes: a skip path that also fired in a suite checkout
// would silence the pinned figures permanently and read as clean. Splitting these
// would let either half be satisfied by breaking the other.
//
// It is itself repo-independent — it asserts about the gate and the findings
// function, never about this archive's contents — so it runs everywhere.
func TestLiveArchiveCalibrationRunsInASuiteCheckoutAndSkipsElsewhere(t *testing.T) {
	// --- Half one: a suite checkout runs the assertions, and they still bite. ---
	suiteRoot := t.TempDir()
	if mkdirError := os.MkdirAll(filepath.Join(suiteRoot, "skills", "do-work", "actions"), 0o755); mkdirError != nil {
		t.Fatalf("mkdir suite actions: %v", mkdirError)
	}
	if writeError := os.WriteFile(filepath.Join(suiteRoot, "skills", "do-work", "actions", "version.md"),
		[]byte("**Current version**: 1.4.0\n"), 0o644); writeError != nil {
		t.Fatalf("write suite version.md: %v", writeError)
	}
	if skipReason := suiteCheckoutSkipReason(suiteRoot); skipReason != "" {
		t.Fatalf("a suite checkout is exactly where these figures apply, got skip %q", skipReason)
	}

	// Biting is the half a silenced check cannot fake. An aggregate carrying the
	// pinned days with the WRONG figures must produce a finding for each.
	wrongFigures := DurationAggregate{
		Samples: make([]DurationSample, 195),
		Days: []DurationDay{
			{DayKey: "2026-07-31", MedianMinutes: 9.0833, HasMedian: true, KeptCount: 12, CompletedCount: 12},
			{DayKey: "2026-08-15", MedianMinutes: 41.2, HasMedian: true, KeptCount: 3, CompletedCount: 3},
		},
	}
	wrongFindings := calibratedLiveArchiveFindings(wrongFigures)
	if len(wrongFindings) != 4 {
		t.Fatalf("both pinned days carry a wrong median AND a wrong count — want 4 findings, got %d: %v",
			len(wrongFindings), wrongFindings)
	}

	// And an aggregate missing the pinned days entirely is a finding, never silence.
	if len(calibratedLiveArchiveFindings(DurationAggregate{})) == 0 {
		t.Fatal("an aggregate with none of the pinned days must report findings, not pass")
	}

	// --- Half two: a consumer root skips, with the condition named. ---
	// do-work/ at the root, the suite vendored under .claude/skills/ — the layout
	// where the resolved archive is the consumer's own.
	consumerRoot := t.TempDir()
	if mkdirError := os.MkdirAll(filepath.Join(consumerRoot, "do-work", "queue"), 0o755); mkdirError != nil {
		t.Fatalf("mkdir consumer queue: %v", mkdirError)
	}
	if mkdirError := os.MkdirAll(filepath.Join(consumerRoot, ".claude", "skills", "do-work", "actions"), 0o755); mkdirError != nil {
		t.Fatalf("mkdir vendored suite: %v", mkdirError)
	}
	if writeError := os.WriteFile(filepath.Join(consumerRoot, ".claude", "skills", "do-work", "actions", "version.md"),
		[]byte("**Current version**: 1.4.0\n"), 0o644); writeError != nil {
		t.Fatalf("write vendored version.md: %v", writeError)
	}
	consumerSkipReason := suiteCheckoutSkipReason(consumerRoot)
	if consumerSkipReason == "" {
		t.Fatal("a consumer root must skip the pinned figures — its archive is not the one they were calibrated against")
	}
	if !strings.Contains(consumerSkipReason, "suite checkout") {
		t.Errorf("the skip reason must name the suite-checkout condition, got %q", consumerSkipReason)
	}
	// No path in the message: a path reads as "look here and fix it", and nothing
	// is missing — the consumer's archive is simply theirs.
	if strings.Contains(consumerSkipReason, consumerRoot) || strings.Contains(consumerSkipReason, "version.md") {
		t.Errorf("the skip reason must not name a path — nothing is missing: %q", consumerSkipReason)
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

// Descent below the baseline for Panel B's 11px slowest-day annotation. It
// exists only in that test's clearance question, so it lives here.
const durationsLabelTextDescentUnits = 2.0

// The x-axis domain is anchored to whole UTC days, the first completion
// floored to its UTC midnight and the midnight AFTER the last, because the
// renderer's day buckets sit at their days' midnights and a domain that began at
// the first completion INSTANT put every bucket left of its samples and pushed
// Panel B off canvas at one or two active days.
//
// It lives in the test file so the parity assertion keeps a second side to be
// held against without leaving an unused definition in the shipped code. If the
// renderer's domain moves, TestJavaScriptBehaviorDurationsDayBucketsStayInsideThePlot
// fails; if this moves, the same test fails from the other direction.
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
	rangeStart = rangeStart.UTC().Truncate(24 * time.Hour)
	rangeEnd = rangeEnd.UTC().Truncate(24 * time.Hour).Add(24 * time.Hour)
	return rangeStart, rangeEnd, true
}

// durationLabelPlotX places one completion instant in plot-local user units, 0 at
// the plot's left edge and durationsPlotWidthUnits at its right. A zero-width
// domain collapses to a single column, exactly as the renderer's `timeSpan || 1`
// guard does.
func durationLabelPlotX(completionTime time.Time, rangeStart time.Time, rangeEnd time.Time) float64 {
	domainSeconds := rangeEnd.Sub(rangeStart).Seconds()
	if domainSeconds <= 0 {
		return 0
	}
	return (completionTime.Sub(rangeStart).Seconds() / domainSeconds) * durationsPlotWidthUnits
}

// ---- the done card's implementation span -----------------------------------

// The Recently-Done card's duration reading and the Durations view are two
// READERS of one read-time rule, never two definitions of it. The boundary is
// pinned against analysisOutlierCeiling itself: a restated "4 hours" here would
// keep passing with the ceiling moved and would prove nothing about the rule.
func TestImplementationSpanVerdictBoundaryReadsTheOutlierCeiling(t *testing.T) {
	claimInstant := time.Date(2026, 7, 4, 9, 0, 0, 0, time.UTC)
	spanEndingAfter := func(offset time.Duration) ImplementationSpan {
		return measureImplementationSpan(durationTicket("REQ-401", "B",
			claimInstant.Format(time.RFC3339),
			claimInstant.Add(offset).Format(time.RFC3339)))
	}

	atCeiling := spanEndingAfter(analysisOutlierCeiling)
	if !atCeiling.StampsParsed {
		t.Fatalf("both stamps parse, yet the span measured nothing: %#v", atCeiling)
	}
	if atCeiling.WallMinutes != analysisOutlierCeiling.Minutes() {
		t.Errorf("at-ceiling span = %v min, want %v", atCeiling.WallMinutes, analysisOutlierCeiling.Minutes())
	}
	if atCeiling.ExclusionReason != "" {
		t.Errorf("a span exactly at the ceiling read %q, want the plain verdict — the rule excludes spans OVER the ceiling, not at it", atCeiling.ExclusionReason)
	}

	overCeiling := spanEndingAfter(analysisOutlierCeiling + time.Minute)
	if overCeiling.ExclusionReason != "paused" {
		t.Errorf("a span one minute over the ceiling read %q, want \"paused\"", overCeiling.ExclusionReason)
	}

	underCeiling := spanEndingAfter(analysisOutlierCeiling - time.Minute)
	if underCeiling.ExclusionReason != "" {
		t.Errorf("a span one minute under the ceiling read %q, want the plain verdict", underCeiling.ExclusionReason)
	}
}

// A reversed span is data, not an error to swallow, and an unparseable stamp is
// "unmeasured" rather than a span of zero — a zero would print "0.0 min" on the
// card and read as instant work.
func TestImplementationSpanMarksReversedStampsAndRefusesUnparseableOnes(t *testing.T) {
	reversed := measureImplementationSpan(durationTicket("REQ-403", "B",
		"2026-07-04T12:45:00Z", "2026-07-04T10:05:00Z"))
	if !reversed.StampsParsed {
		t.Fatalf("a reversed pair of parseable stamps measured nothing: %#v", reversed)
	}
	if reversed.ExclusionReason != "reversed" {
		t.Errorf("reversed span read %q, want \"reversed\"", reversed.ExclusionReason)
	}
	if reversed.WallMinutes != -160 {
		t.Errorf("reversed span = %v min, want the raw signed -160 (clamping to zero hides the broken stamp)", reversed.WallMinutes)
	}

	unmeasurableTickets := []struct {
		caseName string
		ticket   *RequestTicket
	}{
		{"no claim stamp", durationTicket("REQ-404", "B", "", "2026-07-04T10:05:00Z")},
		{"unparseable claim stamp", durationTicket("REQ-405", "B", "yesterday", "2026-07-04T10:05:00Z")},
		{"no completion stamp", durationTicket("REQ-406", "B", "2026-07-04T10:05:00Z", "")},
		{"unparseable completion stamp", durationTicket("REQ-407", "B", "2026-07-04T10:05:00Z", "not-a-time")},
		{"no ticket at all", nil},
	}
	for _, unmeasurable := range unmeasurableTickets {
		span := measureImplementationSpan(unmeasurable.ticket)
		if span.StampsParsed {
			t.Errorf("%s measured a span (%v min); want no span at all", unmeasurable.caseName, span.WallMinutes)
		}
		if span.WallMinutes != 0 || span.ExclusionReason != "" {
			t.Errorf("%s carried %v min / %q alongside an unmeasured span; want the zero value",
				unmeasurable.caseName, span.WallMinutes, span.ExclusionReason)
		}
	}
}

// One definition, two readers. A second ceiling — or a second subtraction order —
// introduced on either side breaks this and nothing else in the suite would.
//
// The ceiling half of that claim needs the fixture to STRADDLE the ceiling, not
// merely to span it widely: three samples at 40 min, 18 h and −3 h agree under any
// second ceiling anywhere in (40 min, 18 h), which is most of the plausible ones.
// The straddling pair below is derived FROM analysisOutlierCeiling, so a second
// definition disagrees with the real one at the only place a threshold can be
// caught — its own boundary — and moving the real ceiling moves the pair with it.
func TestImplementationSpanAgreesWithTheDurationsAggregate(t *testing.T) {
	ceilingClaim := time.Date(2026, 7, 5, 6, 0, 0, 0, time.UTC)
	tickets := []*RequestTicket{
		durationTicket("REQ-410", "A", "2026-07-05T09:00:00Z", "2026-07-05T09:40:00Z"),
		durationTicket("REQ-411", "B", "2026-07-05T09:00:00Z", "2026-07-06T03:00:00Z"),
		durationTicket("REQ-412", "C", "2026-07-05T12:00:00Z", "2026-07-05T09:00:00Z"),
		durationTicket("REQ-413", "B", ceilingClaim.Format(time.RFC3339),
			ceilingClaim.Add(analysisOutlierCeiling).Format(time.RFC3339)),
		durationTicket("REQ-414", "B", ceilingClaim.Format(time.RFC3339),
			ceilingClaim.Add(analysisOutlierCeiling+time.Minute).Format(time.RFC3339)),
	}
	aggregate := buildDurationAggregate(tickets)
	if len(aggregate.Samples) != len(tickets) {
		t.Fatalf("aggregate holds %d samples for %d tickets; the agreement has nothing to compare", len(aggregate.Samples), len(tickets))
	}

	verdictsWitnessed := map[string]bool{}
	for _, ticket := range tickets {
		sample := findDurationSample(t, aggregate, ticket.RequestId)
		span := measureImplementationSpan(ticket)
		if span.WallMinutes != sample.WallMinutes {
			t.Errorf("%s: card span = %v min, Durations sample = %v min", ticket.RequestId, span.WallMinutes, sample.WallMinutes)
		}
		if span.ExclusionReason != sample.DayMedianExclusion {
			t.Errorf("%s: card verdict = %q, Durations verdict = %q", ticket.RequestId, span.ExclusionReason, sample.DayMedianExclusion)
		}
		verdictsWitnessed[span.ExclusionReason] = true
	}
	// Vacuity guard: agreement over one verdict is not agreement about the rule.
	for _, requiredVerdict := range []string{"", "paused", "reversed"} {
		if !verdictsWitnessed[requiredVerdict] {
			t.Fatalf("the fixture never produced the %q verdict, so this test cannot witness a disagreement about it", requiredVerdict)
		}
	}
	// Second vacuity guard, for the ceiling half specifically: the straddling pair
	// must actually land on opposite sides of the real ceiling, or a second ceiling
	// could sit between them and agree with the first everywhere this test looks.
	atCeiling := findDurationSample(t, aggregate, "REQ-413")
	pastCeiling := findDurationSample(t, aggregate, "REQ-414")
	if atCeiling.DayMedianExclusion != "" || pastCeiling.DayMedianExclusion != "paused" {
		t.Fatalf("the straddling pair read %q / %q, want \"\" / \"paused\" — it no longer brackets the ceiling, so a second ceiling would pass unnoticed",
			atCeiling.DayMedianExclusion, pastCeiling.DayMedianExclusion)
	}
}
