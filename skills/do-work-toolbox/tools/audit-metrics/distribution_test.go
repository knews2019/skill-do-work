package main

import "testing"

// TestNearestRankPercentile pins the distribution math the audit report quotes:
// median/p90/p95/max are nearest-rank (ceil(P/100*n), 1-indexed) over the
// ascending values, never interpolated — interpolation would print line counts
// no file actually has.
func TestNearestRankPercentile(t *testing.T) {
	testCases := []struct {
		name         string
		sortedValues []int
		percent      int
		want         int
	}{
		{"median of even-length slice", []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}, 50, 50},
		{"p90 of ten values", []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}, 90, 90},
		{"p95 of ten values rounds up to max", []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}, 95, 100},
		{"p100 is the max", []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}, 100, 100},
		{"median of odd-length slice", []int{1, 2, 3, 4, 5}, 50, 3},
		{"p90 of five values", []int{1, 2, 3, 4, 5}, 90, 5},
		{"single value is every percentile", []int{7}, 50, 7},
		{"empty slice yields zero", []int{}, 95, 0},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := nearestRankPercentile(testCase.sortedValues, testCase.percent)
			if got != testCase.want {
				t.Fatalf("nearestRankPercentile(%v, %d) = %d, want %d",
					testCase.sortedValues, testCase.percent, got, testCase.want)
			}
		})
	}
}
