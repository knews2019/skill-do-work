package main

import "sort"

// Distribution math and band flagging for the audit's metric summaries. The
// audit report quotes these numbers verbatim, so they are deliberately simple:
// nearest-rank percentiles (never interpolated — interpolation would print a
// line count no file actually has) and strict-greater band thresholds.

// nearestRankPercentile returns the nearest-rank percentile of sortedValues
// (ascending): the value at 1-indexed rank ceil(percent/100 * n). An empty
// slice yields 0.
func nearestRankPercentile(sortedValues []int, percent int) int {
	if len(sortedValues) == 0 {
		return 0
	}
	rank := (percent*len(sortedValues) + 99) / 100
	if rank < 1 {
		rank = 1
	}
	if rank > len(sortedValues) {
		rank = len(sortedValues)
	}
	return sortedValues[rank-1]
}

// distributionSummary is the row the audit's distribution tables print for one
// metric.
type distributionSummary struct {
	Median int
	P90    int
	P95    int
	Max    int
}

// summarizeDistribution computes median/p90/p95/max over values (any order;
// the input slice is not mutated).
func summarizeDistribution(values []int) distributionSummary {
	sortedValues := append([]int(nil), values...)
	sort.Ints(sortedValues)
	return distributionSummary{
		Median: nearestRankPercentile(sortedValues, 50),
		P90:    nearestRankPercentile(sortedValues, 90),
		P95:    nearestRankPercentile(sortedValues, 95),
		Max:    nearestRankPercentile(sortedValues, 100),
	}
}

// bandThresholdUnset marks a band threshold the caller did not pass. Bands come
// ONLY from flags — no flag, no band output — so every threshold flag defaults
// to this sentinel (thresholds are counts, never negative).
const bandThresholdUnset = -1

// bandLabelForValue applies WATCH/FLAG thresholds to one value. Strictly
// greater is flagged; equal is NOT — a file sitting exactly at the threshold is
// at the line, not over it. FLAG wins over WATCH when both are exceeded.
func bandLabelForValue(value int, watchThreshold int, flagThreshold int) string {
	if flagThreshold != bandThresholdUnset && value > flagThreshold {
		return "FLAG"
	}
	if watchThreshold != bandThresholdUnset && value > watchThreshold {
		return "WATCH"
	}
	return ""
}
