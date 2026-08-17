package main

import (
	"math"
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
