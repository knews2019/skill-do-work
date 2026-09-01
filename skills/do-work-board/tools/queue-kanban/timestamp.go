package main

import (
	"fmt"
	"io"
	"time"
)

// The `now` subcommand emits the current UTC instant in the exact shape the
// skill's Timestamp rule prescribes (actions/work-reference.md → Request File
// Schema — Full Frontmatter): YYYY-MM-DDTHH:MM:SSZ, whole seconds, always UTC.
//
// Why this lives in the tool at all, when `date -u +%Y-%m-%dT%H:%M:%SZ` already
// works: that command is a POSIX-ism with no Windows `cmd` equivalent, so the
// prescribed stamp is silently unobtainable there. Putting the writer beside the
// readers (parseTimestamp in model.go, the future-stamp guard) also means writer
// and readers agree by construction — timestamp_test.go asserts the round trip.
//
// This board-local spelling is retained read-only compatibility. Shipped guidance uses
// the build-on-demand core `do-work-cli now` command as the canonical owner.

// canonicalTimestampLayout is the Timestamp rule's shape. Deliberately not
// time.RFC3339: that layout emits a numeric offset for a non-UTC instant
// (+09:00) and would carry sub-second digits, neither of which the schema's
// *_at fields accept.
const canonicalTimestampLayout = "2006-01-02T15:04:05Z"

// formatCanonicalTimestamp renders an instant as the schema's UTC stamp.
//
// The instant is converted to UTC, never relabelled: writing local wall-clock
// time with a Z suffix is one of the two corruptions the board diagnoses
// (futureStampCauseClause names both; fabricating a value is the other, and no
// writer can prevent that one) — in any zone east of UTC it yields a *future*
// instant, which freezes the board's state timer and flags the card with a
// clock-skew warning. Sub-second precision is truncated rather than rounded,
// so a stamp can never round forward past the board's own `now`.
func formatCanonicalTimestamp(instant time.Time) string {
	return instant.UTC().Truncate(time.Second).Format(canonicalTimestampLayout)
}

// writeCanonicalTimestamp prints the stamp and a single newline — nothing else,
// so compatibility callers can capture `queue-kanban now` directly, the way
// `REQ-$(queue-kanban next-req)` already is.
func writeCanonicalTimestamp(outputWriter io.Writer, instant time.Time) {
	fmt.Fprintf(outputWriter, "%s\n", formatCanonicalTimestamp(instant))
}
