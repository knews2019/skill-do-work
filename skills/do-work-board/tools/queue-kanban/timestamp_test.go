package main

import (
	"strings"
	"testing"
	"time"
)

// The `now` subcommand's whole job is to emit the one timestamp shape the skill's
// Timestamp rule prescribes (actions/work-reference.md → Request File Schema), so
// the tests assert that shape and nothing about wall-clock time: the clock is
// injected, exactly as LoadBoard and runVerifyProbes take theirs.

func TestFormatCanonicalTimestampMatchesSchemaShape(t *testing.T) {
	fixedInstant := time.Date(2026, 8, 3, 16, 9, 6, 0, time.UTC)

	formatted := formatCanonicalTimestamp(fixedInstant)

	if formatted != "2026-08-03T16:09:06Z" {
		t.Fatalf("formatCanonicalTimestamp = %q, want 2026-08-03T16:09:06Z", formatted)
	}
}

// The board's own reader must accept what the writer emits. This is the whole
// point of putting the writer next to the readers — writer and readers agree by
// construction rather than by convention (REQ-076).
func TestFormatCanonicalTimestampRoundTripsThroughBoardParser(t *testing.T) {
	fixedInstant := time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)

	formatted := formatCanonicalTimestamp(fixedInstant)
	parsed, parseOk := parseTimestamp(formatted)

	if !parseOk {
		t.Fatalf("parseTimestamp rejected %q — the board could not read a stamp the tool wrote", formatted)
	}
	if !parsed.Equal(fixedInstant) {
		t.Fatalf("round trip changed the instant: wrote %s, read back %s", fixedInstant, parsed)
	}
}

// A local-zone instant must be converted, never labelled. Writing local
// wall-clock time with a Z suffix is one of the two corruptions the board
// diagnoses (futureStampCauseClause names both; fabricating a value is the
// other, and no writer can prevent that one): east of UTC it produces a future
// instant, which freezes the board's stopwatch and raises a clock-skew warning.
// This is the corruption a correct writer rules out, which is why the test is
// here.
func TestFormatCanonicalTimestampConvertsNonUTCZones(t *testing.T) {
	easternOfUTC := time.FixedZone("UTC+9", 9*60*60)
	localInstant := time.Date(2026, 8, 3, 9, 0, 0, 0, easternOfUTC)

	formatted := formatCanonicalTimestamp(localInstant)

	if formatted != "2026-08-03T00:00:00Z" {
		t.Fatalf("formatCanonicalTimestamp = %q, want 2026-08-03T00:00:00Z (converted to UTC, not relabelled)", formatted)
	}
}

// Sub-second precision is truncated, not rounded: the schema's shape carries
// whole seconds, and rounding could push a stamp past the board's `now`.
func TestFormatCanonicalTimestampTruncatesSubSecondPrecision(t *testing.T) {
	instantWithNanos := time.Date(2026, 8, 3, 16, 9, 6, 999_999_999, time.UTC)

	formatted := formatCanonicalTimestamp(instantWithNanos)

	if formatted != "2026-08-03T16:09:06Z" {
		t.Fatalf("formatCanonicalTimestamp = %q, want the second truncated to 16:09:06Z", formatted)
	}
}

// The command prints exactly the stamp plus one newline, so `REQ_STAMP=$(queue-kanban now)`
// is directly usable the way `REQ-$(queue-kanban next-req)` already is.
func TestWriteCanonicalTimestampEmitsOnlyTheStampAndOneNewline(t *testing.T) {
	fixedInstant := time.Date(2026, 8, 3, 16, 9, 6, 0, time.UTC)
	var capturedOutput strings.Builder

	writeCanonicalTimestamp(&capturedOutput, fixedInstant)

	if capturedOutput.String() != "2026-08-03T16:09:06Z\n" {
		t.Fatalf("writeCanonicalTimestamp wrote %q, want the stamp and a single newline", capturedOutput.String())
	}
}

// --- Future-stamp diagnosis wording -----------------------------------------
//
// The board diagnoses a stamp that cannot be real in three places — the
// generate-time data warning, the reversed-span completion-anomaly reason, and
// verify's future-`claimed_at` finding. Each must name BOTH observed causes: a
// fabricated value, guessed or extrapolated instead of read from the clock, and
// local wall-clock time written with a `Z` suffix. Naming only the timezone
// cause tells the fabricating reader to go hunting for a clock bug that is not
// there (REQ-245).
//
// The three are asserted in one place because the contract is that they agree;
// split across the three renderers' own test files, a fourth cause added to one
// of them would pass. They live in this file rather than a renderer's because
// what they share is a claim about timestamps, which is this file's subject.

// futureStampBoardWarning returns the generate-time data warning for the
// synthetic tree's future-claimed REQ-9401.
func futureStampBoardWarning(t *testing.T) string {
	t.Helper()
	board := futureStampSyntheticBoard(t)
	for _, warningText := range board.Warnings {
		if strings.Contains(warningText, "REQ-9401 has future-dated timestamp(s)") {
			return warningText
		}
	}
	t.Fatalf("no future-dated-timestamp warning for REQ-9401; warnings: %v", board.Warnings)
	return ""
}

// reversedSpanAnomalyReason returns the completion-anomaly reason for a
// completed_at that precedes its claimed_at.
func reversedSpanAnomalyReason(t *testing.T) string {
	t.Helper()
	flagged, reason := detectCompletionAnomaly(&RequestTicket{
		RequestId:            "REQ-9340",
		Status:               "completed",
		ClaimedAt:            "2026-01-02T10:00:00Z",
		CompletedAt:          "2026-01-01T10:00:00Z",
		CompletionTime:       time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
		CompletionTimeSource: CompletionFromFrontmatter,
	})
	if !flagged {
		t.Fatalf("reversed claimed→completed span was not flagged, so there is no reason text to check")
	}
	return reason
}

// futureClaimVerifyFinding returns verify's finding for a claimed REQ whose
// claimed_at parses past the skew horizon.
func futureClaimVerifyFinding(t *testing.T) VerifyFinding {
	t.Helper()
	fixedNow := time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"actions/version.md", cleanVersionFile},
		{"CHANGELOG.md", cleanChangelog},
		{"do-work/working/REQ-9341-future-claim.md",
			"---\nid: REQ-9341\nstatus: claimed\ntitle: future\nclaimed_at: " +
				fixedNow.Add(3*time.Hour).Format(time.RFC3339) + "\n---\n"},
	})
	report, verifyError := runVerifyProbes(repoRoot, fixedNow)
	if verifyError != nil {
		t.Fatalf("runVerifyProbes: %v", verifyError)
	}
	for _, finding := range findingsMentioning(report, verifyCategoryClaimNeedsAttention) {
		if strings.Contains(finding.Detail, "future-dated claimed_at") {
			return finding
		}
	}
	t.Fatalf("verify produced no future-dated claimed_at finding: %+v", report.Findings)
	return VerifyFinding{}
}

func TestFutureStampDiagnosesNameBothCauses(t *testing.T) {
	verifyFinding := futureClaimVerifyFinding(t)

	for _, diagnosis := range []struct {
		renderer string
		message  string
	}{
		{"board future-stamp data warning (model.go)", futureStampBoardWarning(t)},
		{"reversed-span completion-anomaly reason (model.go)", reversedSpanAnomalyReason(t)},
		{"verify future-claimed_at finding (verify.go)", verifyFinding.Detail},
	} {
		for _, requiredCause := range []string{"fabricated", "wall-clock", "Z suffix"} {
			if !strings.Contains(diagnosis.message, requiredCause) {
				t.Errorf("%s does not name %q as a possible cause — a reader hitting the other cause is sent to the wrong fix.\ngot: %s",
					diagnosis.renderer, requiredCause, diagnosis.message)
			}
		}
	}
}

// The fix — rewrite the stamp with the current UTC instant — is correct for a
// fabricated value and for a mislabelled local one alike, so naming the second
// cause must not disturb any of the three fix instructions.
func TestFutureStampDiagnosesKeepTheirFixInstruction(t *testing.T) {
	boardWarning := futureStampBoardWarning(t)
	for _, expectedFragment := range []string{
		"rewrite with the current UTC instant",
		"YYYY-MM-DDTHH:MM:SSZ",
		"Timestamp rule in actions/work-reference.md",
	} {
		if !strings.Contains(boardWarning, expectedFragment) {
			t.Errorf("board future-stamp warning lost %q; got: %s", expectedFragment, boardWarning)
		}
	}

	anomalyReason := reversedSpanAnomalyReason(t)
	if !strings.Contains(anomalyReason, "rewrite the wrong stamp with the true UTC instant") {
		t.Errorf("reversed-span reason lost its fix instruction; got: %s", anomalyReason)
	}

	verifyRemedy := futureClaimVerifyFinding(t).Remedy
	if !strings.Contains(verifyRemedy, "queue-kanban now") {
		t.Errorf("verify future-claim remedy lost its `queue-kanban now` instruction; got: %s", verifyRemedy)
	}
}
