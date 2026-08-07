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
// wall-clock time with a Z suffix is the specific corruption the Timestamp rule
// names: east of UTC it produces a future instant, which freezes the board's
// stopwatch and raises a clock-skew warning.
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
