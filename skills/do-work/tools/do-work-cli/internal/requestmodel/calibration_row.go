package requestmodel

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"time"
)

// CalibrationLogHeader is the header a NEW do-work/calibration-log.tsv gets.
// The header is the log's version marker: a log whose header ends at
// completed_at (every log written before max_stamp_gap_minutes existed) keeps
// that header and keeps receiving five-column rows. Nothing rewrites an
// existing header, because finalization proves the log append as "the old
// bytes plus exactly one row".
const CalibrationLogHeader = "req_id\troute\testimated_p50_minutes\twall_minutes\tcompleted_at\tmax_stamp_gap_minutes\n"

// CalibrationGapStampFields names the lifecycle stamps, besides completed_at,
// whose consecutive gaps max_stamp_gap_minutes measures: claimed_at and the
// optional phase stamps. release_at is out because it records shipping, after
// the work. The board's Panel B rule reads the same set
// (queue-kanban durations.go phaseMilestonesOf minus completed_at and
// release_at); _dev/tests/contracts/queue-kanban.sh pins the two lists equal.
var CalibrationGapStampFields = []string{"claimed_at", "planning_at", "dispatch_at", "builder_handback_at", "integration_at", "review_at", "remediation_at", "re_review_at"}

// CalibrationLogCarriesGapColumn reports whether the log's header row ends
// with max_stamp_gap_minutes, which decides the shape of every row appended
// to it. Pass the log as it will read after the append (header included).
func CalibrationLogCarriesGapColumn(logBytes []byte) bool {
	headerLine, _, _ := bytes.Cut(logBytes, []byte("\n"))
	return bytes.HasSuffix(bytes.TrimRight(headerLine, "\r"), []byte("\tmax_stamp_gap_minutes"))
}

// FormatCalibrationRow is the one calibration-log row builder. The lifecycle
// appender, the post-archive verifier, and finalization's append proof all
// call it, so the three can never disagree about a row. completedAt is passed
// in because the appender runs before completed_at is written to the record.
func FormatCalibrationRow(record RequestRecord, completedAt time.Time, logBytes []byte) (string, error) {
	estimateMinutes, estimateError := strconv.Atoi(record.FieldEvidenceByName["estimate"].NestedValues["p50_active_minutes"])
	claimedAt, claimedError := ParseTimestamp(record.ClaimedAt)
	if estimateError != nil || claimedError != nil {
		return "", fmt.Errorf("estimate or claimed_at is not parseable")
	}
	route := record.RouteValue
	if route == "" {
		route = "-"
	}
	row := fmt.Sprintf("%s\t%s\t%d\t%d\t%s", record.RequestID, route, estimateMinutes, int(completedAt.Sub(claimedAt).Minutes()), CanonicalTimestamp(completedAt))
	if CalibrationLogCarriesGapColumn(logBytes) {
		row += fmt.Sprintf("\t%d", int(largestCalibrationStampGap(record, claimedAt, completedAt).Minutes()))
	}
	return row + "\n", nil
}

// largestCalibrationStampGap sorts the parseable stamps by instant and returns
// the longest stretch between two consecutive ones. With no phase stamps the
// only gap is claimed_at → completed_at. Unparseable stamps are skipped.
func largestCalibrationStampGap(record RequestRecord, claimedAt, completedAt time.Time) time.Duration {
	instants := []time.Time{claimedAt, completedAt}
	for _, fieldName := range CalibrationGapStampFields {
		if fieldName == "claimed_at" {
			continue
		}
		if instant, parseError := ParseTimestamp(record.FieldEvidenceByName[fieldName].ScalarValue); parseError == nil {
			instants = append(instants, instant)
		}
	}
	sort.Slice(instants, func(left, right int) bool { return instants[left].Before(instants[right]) })
	var largestGap time.Duration
	for index := 1; index < len(instants); index++ {
		if gap := instants[index].Sub(instants[index-1]); gap > largestGap {
			largestGap = gap
		}
	}
	return largestGap
}
