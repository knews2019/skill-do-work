#!/usr/bin/env bash
# Fixture execution proofs for repair-req-timestamps.
# shellcheck source=_dev/tests/prescribed-shell-harness.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/prescribed-shell-harness.sh"

# repair-req-timestamps: a future created_at in the queue is rewritten to the
# file's own mtime — the actual write instant — and the correction is logged.
repair_mtime_project="$fixture_root/repair-mtime-project"
mkdir -p "$repair_mtime_project/do-work/queue"
printf -- '---\nid: REQ-801\nstatus: pending\ncreated_at: 2093-01-01T00:00:00Z\n---\n\nbody\n' \
  > "$repair_mtime_project/do-work/queue/REQ-801-future.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_mtime_project/do-work/queue/REQ-801-future.md"
repair_mtime_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_mtime_project")" \
  || fail_case 'repair-req-timestamps future-stamp case returned nonzero'
grep -q '^created_at: 2026-08-10T12:00:00Z$' "$repair_mtime_project/do-work/queue/REQ-801-future.md" \
  || fail_case 'repair-req-timestamps future-stamp case did not rewrite the stamp to the file mtime'
printf '%s' "$repair_mtime_output" \
  | grep -q 'REQ-801-future.md created_at: 2093-01-01T00:00:00Z -> 2026-08-10T12:00:00Z (file mtime)' \
  || fail_case 'repair-req-timestamps future-stamp case did not log the correction'

# repair-req-timestamps: impossible orderings in working/ are repaired and
# clamped so created_at <= claimed_at <= completed_at <= now — here the derived
# mtime precedes created_at, so both later fields land exactly on the clamp floor.
repair_order_project="$fixture_root/repair-order-project"
mkdir -p "$repair_order_project/do-work/working"
printf -- '---\nid: REQ-802\nstatus: completed\ncreated_at: 2026-08-10T12:00:00Z\nclaimed_at: 2026-08-01T09:00:00Z\ncompleted_at: 2026-08-03T10:00:00Z\n---\nbody\n' \
  > "$repair_order_project/do-work/working/REQ-802-order.md"
TZ=UTC touch -m -t 202608050800.00 "$repair_order_project/do-work/working/REQ-802-order.md"
repair_order_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_order_project")" \
  || fail_case 'repair-req-timestamps ordering case returned nonzero'
grep -q '^claimed_at: 2026-08-10T12:00:00Z$' "$repair_order_project/do-work/working/REQ-802-order.md" \
  || fail_case 'repair-req-timestamps ordering case did not clamp claimed_at up to created_at'
grep -q '^completed_at: 2026-08-10T12:00:00Z$' "$repair_order_project/do-work/working/REQ-802-order.md" \
  || fail_case 'repair-req-timestamps ordering case did not clamp completed_at up to the repaired claimed_at'
printf '%s' "$repair_order_output" | grep -q 'clamped to 2026-08-10T12:00:00Z' \
  || fail_case 'repair-req-timestamps ordering case did not log the clamp'

# repair-req-timestamps: a committed file that matches HEAD is repaired to the
# author time of the commit that introduced the stamp line, not to a fresh clone's
# meaningless mtime.
repair_blame_project="$fixture_root/repair-blame-project"
fixture_repo_init "$repair_blame_project"
mkdir -p "$repair_blame_project/do-work/queue"
printf -- '---\nid: REQ-803\nstatus: pending\ncreated_at: 2026-08-14T09:00:00Z\nclaimed_at: "2093-02-02T02:02:02Z"\n---\nbody\n' \
  > "$repair_blame_project/do-work/queue/REQ-803-committed.md"
git -C "$repair_blame_project" add -A
GIT_AUTHOR_DATE='2026-08-15T14:00:00Z' GIT_COMMITTER_DATE='2026-08-15T14:05:00Z' \
  git -C "$repair_blame_project" commit -qm fixture
repair_blame_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_blame_project")" \
  || fail_case 'repair-req-timestamps committed-file case returned nonzero'
grep -q '^claimed_at: 2026-08-15T14:00:00Z$' "$repair_blame_project/do-work/queue/REQ-803-committed.md" \
  || fail_case 'repair-req-timestamps committed-file case did not use the introducing commit author time (quoted stamps must be repairable too)'
printf '%s' "$repair_blame_output" | grep -q 'author time' \
  || fail_case 'repair-req-timestamps committed-file case did not name the commit author time as the replacement source'

# repair-req-timestamps: a clean fixture passes through byte-identical — including
# the shapes the repairer must not touch: a nested (indented) calculated_at, a
# numeric-offset value it refuses permanently because the timezone arithmetic is
# the risk and not the obstacle (REQ-257), and an archive-scope directory it must
# never scan.
repair_clean_project="$fixture_root/repair-clean-project"
mkdir -p "$repair_clean_project/do-work/queue" "$repair_clean_project/do-work/archive"
printf -- '---\nid: REQ-804\nstatus: pending\ncreated_at: 2026-08-10T12:00:00Z   # trailing comment\nblocked_at: 2026-08-11T09:00:00+09:00\nestimate:\n  calculated_at: 2093-06-06T06:06:06Z\n---\nbody\n' \
  > "$repair_clean_project/do-work/queue/REQ-804-clean.md"
printf -- '---\nid: REQ-805\nstatus: completed\ncreated_at: 2093-03-03T03:03:03Z\n---\nbody\n' \
  > "$repair_clean_project/do-work/archive/REQ-805-archived.md"
cp "$repair_clean_project/do-work/queue/REQ-804-clean.md" "$fixture_root/repair-clean-before.md"
cp "$repair_clean_project/do-work/archive/REQ-805-archived.md" "$fixture_root/repair-archive-before.md"
repair_clean_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_clean_project")" \
  || fail_case 'repair-req-timestamps clean-fixture case returned nonzero'
[ -z "$repair_clean_output" ] \
  || fail_case 'repair-req-timestamps clean-fixture case printed output for a no-op run'
cmp -s "$fixture_root/repair-clean-before.md" "$repair_clean_project/do-work/queue/REQ-804-clean.md" \
  || fail_case 'repair-req-timestamps clean-fixture case changed a file with nothing provably wrong'
cmp -s "$fixture_root/repair-archive-before.md" "$repair_clean_project/do-work/archive/REQ-805-archived.md" \
  || fail_case 'repair-req-timestamps clean-fixture case wrote into archive scope'

# repair-req-timestamps: a tripped guard leaves the file byte-identical and exits
# nonzero — here the truncation floor: a file at less than half its committed size
# lost content before the run, and repairing a stamp in the remains would help
# commit the loss.
repair_guard_project="$fixture_root/repair-guard-project"
fixture_repo_init "$repair_guard_project"
mkdir -p "$repair_guard_project/do-work/queue"
{
  printf -- '---\nid: REQ-806\nstatus: pending\ncreated_at: 2093-04-04T04:04:04Z\n---\n'
  awk 'BEGIN { for (line_index = 1; line_index <= 200; line_index++) print "ballast decision-trail line " line_index }'
} > "$repair_guard_project/do-work/queue/REQ-806-truncated.md"
fixture_repo_commit_all "$repair_guard_project" fixture
head -n 5 "$repair_guard_project/do-work/queue/REQ-806-truncated.md" > "$fixture_root/repair-truncated.tmp"
mv "$fixture_root/repair-truncated.tmp" "$repair_guard_project/do-work/queue/REQ-806-truncated.md"
cp "$repair_guard_project/do-work/queue/REQ-806-truncated.md" "$fixture_root/repair-guard-before.md"
repair_guard_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_guard_project")" \
  && fail_case 'repair-req-timestamps tripped-guard case exited zero on a truncated file'
cmp -s "$fixture_root/repair-guard-before.md" "$repair_guard_project/do-work/queue/REQ-806-truncated.md" \
  || fail_case 'repair-req-timestamps tripped-guard case modified the truncated file'
printf '%s' "$repair_guard_output" | grep -q 'content was lost before this run' \
  || fail_case 'repair-req-timestamps tripped-guard case did not name the truncation as the reason'

# repair-req-timestamps: an unquoted space-separated future instant is repaired
# whole — never half-rewritten into a date plus a phantom time-of-day — and the
# audit line reports the full old value (REQ-255 I1, the corrupting shape).
repair_space_project="$fixture_root/repair-space-project"
mkdir -p "$repair_space_project/do-work/queue"
printf -- '---\nid: REQ-807\nstatus: pending\ncreated_at: 2093-01-01 00:00:00\n---\nbody\n' \
  > "$repair_space_project/do-work/queue/REQ-807-space.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_space_project/do-work/queue/REQ-807-space.md"
repair_space_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_space_project")" \
  || fail_case 'repair-req-timestamps space-separated case returned nonzero'
grep -q '^created_at: 2026-08-10T12:00:00Z$' "$repair_space_project/do-work/queue/REQ-807-space.md" \
  || fail_case 'repair-req-timestamps space-separated case did not rewrite the whole value to the canonical form'
grep -q '00:00:00$' "$repair_space_project/do-work/queue/REQ-807-space.md" \
  && fail_case 'repair-req-timestamps space-separated case left a phantom time-of-day suffix behind'
printf '%s' "$repair_space_output" \
  | grep -q 'REQ-807-space.md created_at: 2093-01-01 00:00:00 -> 2026-08-10T12:00:00Z (file mtime)' \
  || fail_case 'repair-req-timestamps space-separated case did not report the full old value in the audit line'

# repair-req-timestamps: a quoted space-separated future instant is repaired to
# the canonical unquoted form instead of silently truncating at the unmatched
# quote and passing through (REQ-255 I1's quoted sibling).
repair_quoted_space_project="$fixture_root/repair-quoted-space-project"
mkdir -p "$repair_quoted_space_project/do-work/queue"
printf -- '---\nid: REQ-808\nstatus: pending\ncreated_at: "2093-01-01 00:00:00"\n---\nbody\n' \
  > "$repair_quoted_space_project/do-work/queue/REQ-808-quoted-space.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_quoted_space_project/do-work/queue/REQ-808-quoted-space.md"
repair_quoted_space_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_quoted_space_project")" \
  || fail_case 'repair-req-timestamps quoted-space case returned nonzero'
grep -q '^created_at: 2026-08-10T12:00:00Z$' "$repair_quoted_space_project/do-work/queue/REQ-808-quoted-space.md" \
  || fail_case 'repair-req-timestamps quoted-space case did not repair the quoted instant to the canonical unquoted form'
printf '%s' "$repair_quoted_space_output" | grep -q 'REQ-808-quoted-space.md created_at' \
  || fail_case 'repair-req-timestamps quoted-space case did not log the correction'

# repair-req-timestamps: a CRLF-fenced file is scanned like the board scans it,
# and a repair preserves every line's CRLF ending — Windows agents are the
# likeliest source of both CRLF files and wrong local-time stamps (REQ-255 I2).
repair_crlf_project="$fixture_root/repair-crlf-project"
mkdir -p "$repair_crlf_project/do-work/queue"
printf -- '---\r\nid: REQ-809\r\nstatus: pending\r\ncreated_at: 2093-03-03T03:03:03Z\r\n---\r\nbody\r\n' \
  > "$repair_crlf_project/do-work/queue/REQ-809-crlf.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_crlf_project/do-work/queue/REQ-809-crlf.md"
repair_crlf_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_crlf_project")" \
  || fail_case 'repair-req-timestamps CRLF case returned nonzero'
grep -q $'^created_at: 2026-08-10T12:00:00Z\r$' "$repair_crlf_project/do-work/queue/REQ-809-crlf.md" \
  || fail_case 'repair-req-timestamps CRLF case did not repair the stamp behind the CRLF fence (or dropped the CR)'
[ "$(grep -c $'\r$' "$repair_crlf_project/do-work/queue/REQ-809-crlf.md")" -eq 6 ] \
  || fail_case 'repair-req-timestamps CRLF case did not preserve every CRLF line ending'
printf '%s' "$repair_crlf_output" | grep -q 'REQ-809-crlf.md created_at: 2093-03-03T03:03:03Z -> 2026-08-10T12:00:00Z' \
  || fail_case 'repair-req-timestamps CRLF case did not log the correction'

# repair-req-timestamps: a BOM-prefixed file is scanned like the board scans it
# (the board strips the BOM before the fence match), and a repair keeps the BOM
# bytes in place (REQ-255 I2).
repair_bom_project="$fixture_root/repair-bom-project"
mkdir -p "$repair_bom_project/do-work/queue"
printf -- '\xef\xbb\xbf---\nid: REQ-810\nstatus: pending\ncreated_at: 2093-04-04T04:04:04Z\n---\nbody\n' \
  > "$repair_bom_project/do-work/queue/REQ-810-bom.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_bom_project/do-work/queue/REQ-810-bom.md"
repair_bom_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_bom_project")" \
  || fail_case 'repair-req-timestamps BOM case returned nonzero'
grep -q '^created_at: 2026-08-10T12:00:00Z$' "$repair_bom_project/do-work/queue/REQ-810-bom.md" \
  || fail_case 'repair-req-timestamps BOM case did not repair the stamp behind the BOM-prefixed fence'
[ "$(head -c 3 "$repair_bom_project/do-work/queue/REQ-810-bom.md")" = "$(printf '\xef\xbb\xbf')" ] \
  || fail_case 'repair-req-timestamps BOM case did not keep the BOM bytes in place'
printf '%s' "$repair_bom_output" | grep -q 'REQ-810-bom.md created_at: 2093-04-04T04:04:04Z -> 2026-08-10T12:00:00Z' \
  || fail_case 'repair-req-timestamps BOM case did not log the correction'

# repair-req-timestamps: a shape-valid but calendar-impossible stamp is left
# byte-identical for diagnosis — the board's parser rejects it, so erasing it
# to a derived instant would destroy the malformed evidence while claiming
# parity (REQ-255, PR #145 external review). The range check must match the
# read side's real calendar: month, day-in-that-month, leap years, and time
# components — a real leap-day future stamp must still be repaired.
repair_calendar_project="$fixture_root/repair-calendar-project"
mkdir -p "$repair_calendar_project/do-work/queue"
printf -- '---\nid: REQ-811\nstatus: pending\ncreated_at: 9999-99-99T99:99:99Z\n---\nbody\n' \
  > "$repair_calendar_project/do-work/queue/REQ-811-impossible.md"
printf -- '---\nid: REQ-812\nstatus: pending\ncreated_at: 2093-04-31T10:00:00Z\n---\nbody\n' \
  > "$repair_calendar_project/do-work/queue/REQ-812-april-31.md"
printf -- '---\nid: REQ-813\nstatus: pending\ncreated_at: 2093-02-29T10:00:00Z\n---\nbody\n' \
  > "$repair_calendar_project/do-work/queue/REQ-813-not-a-leap-year.md"
printf -- '---\nid: REQ-814\nstatus: pending\ncreated_at: 2092-02-29T10:00:00Z\n---\nbody\n' \
  > "$repair_calendar_project/do-work/queue/REQ-814-real-leap-day.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_calendar_project/do-work/queue/"REQ-81*.md
for impossible_fixture in REQ-811-impossible REQ-812-april-31 REQ-813-not-a-leap-year; do
  cp "$repair_calendar_project/do-work/queue/$impossible_fixture.md" "$fixture_root/$impossible_fixture-before.md"
done
repair_calendar_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_calendar_project")" \
  || fail_case 'repair-req-timestamps calendar case returned nonzero'
for impossible_fixture in REQ-811-impossible REQ-812-april-31 REQ-813-not-a-leap-year; do
  cmp -s "$fixture_root/$impossible_fixture-before.md" "$repair_calendar_project/do-work/queue/$impossible_fixture.md" \
    || fail_case "repair-req-timestamps calendar case erased the impossible stamp in $impossible_fixture instead of leaving it for diagnosis"
done
grep -q '^created_at: 2026-08-10T12:00:00Z$' "$repair_calendar_project/do-work/queue/REQ-814-real-leap-day.md" \
  || fail_case 'repair-req-timestamps calendar case refused a real leap-day future stamp the board parses'
printf '%s' "$repair_calendar_output" | grep -q 'REQ-811\|REQ-812\|REQ-813' \
  && fail_case 'repair-req-timestamps calendar case logged a correction for a value it must not touch'

# repair-req-timestamps: a numeric UTC offset or fractional seconds is refused
# permanently — the decided answer to the residual, not a to-do (REQ-257). Every
# such shape the read side parses and future-badges is left byte-identical, is
# never logged, and does not fail the run. This case is what fails if someone
# quietly teaches comparison_key_for offset arithmetic: a refusal that lived only
# in a header comment pinned nothing, and the risk being refused is real —
# repairing here would rewrite an instant the read side already reads correctly.
repair_refused_shape_project="$fixture_root/repair-refused-shape-project"
mkdir -p "$repair_refused_shape_project/do-work/queue"
printf -- '---\nid: REQ-817\nstatus: pending\ncreated_at: 2093-01-01T00:00:00+02:00\n---\nbody\n' \
  > "$repair_refused_shape_project/do-work/queue/REQ-817-offset-ahead.md"
printf -- '---\nid: REQ-818\nstatus: pending\ncreated_at: 2093-01-01T00:00:00-05:00\n---\nbody\n' \
  > "$repair_refused_shape_project/do-work/queue/REQ-818-offset-behind.md"
printf -- '---\nid: REQ-819\nstatus: pending\ncreated_at: "2093-01-01T00:00:00+02:00"\n---\nbody\n' \
  > "$repair_refused_shape_project/do-work/queue/REQ-819-quoted-offset.md"
printf -- '---\nid: REQ-820\nstatus: pending\ncreated_at: 2093-01-01T00:00:00.500Z\n---\nbody\n' \
  > "$repair_refused_shape_project/do-work/queue/REQ-820-fractional-zulu.md"
printf -- '---\nid: REQ-821\nstatus: pending\ncreated_at: 2093-01-01T00:00:00.5\n---\nbody\n' \
  > "$repair_refused_shape_project/do-work/queue/REQ-821-fractional-bare.md"
printf -- '---\nid: REQ-822\nstatus: pending\ncreated_at: 2093-01-01 00:00:00.5\n---\nbody\n' \
  > "$repair_refused_shape_project/do-work/queue/REQ-822-fractional-space.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_refused_shape_project/do-work/queue/"REQ-8[12]*.md
for refused_fixture in "$repair_refused_shape_project/do-work/queue/"REQ-8[12]*.md; do
  cp "$refused_fixture" "$fixture_root/refused-$(basename "$refused_fixture")"
done
repair_refused_shape_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_refused_shape_project")" \
  || fail_case 'repair-req-timestamps refused-shape case returned nonzero for shapes it deliberately does not repair'
[ -z "$repair_refused_shape_output" ] \
  || fail_case 'repair-req-timestamps refused-shape case logged a correction for an offset or fractional stamp'
for refused_fixture in "$repair_refused_shape_project/do-work/queue/"REQ-8[12]*.md; do
  cmp -s "$fixture_root/refused-$(basename "$refused_fixture")" "$refused_fixture" \
    || fail_case "repair-req-timestamps refused-shape case rewrote $(basename "$refused_fixture") — the offset/fractional refusal is a decided permanent answer; changing it means re-deciding it, not widening comparison_key_for"
done

# repair-req-timestamps: the read-side layout list the refusal is scoped against
# stays in lock-step — the residual is the GAP between the board's parser and
# this script's, so pinning only the refusing side would let a widened board grow
# the gap silently. A new layout here means the offset/fractional decision (and
# the header paragraph stating it) gets re-made, not inherited.
# The extraction keys on the CONDITION — every element line inside parseTimestamp's layout
# slice, whatever its spelling — never on an enumeration of the spellings that happen to be
# there today (REQ-271; _dev/primes/prime-shell-commands.md § Closed Enumerations Go Stale).
# The predecessor matched only `time.RFC3339` or a `"2006…"`-prefixed literal, so REQ-257's
# own review added `time.RFC3339Nano` and `time.DateTime` to the slice and the suite stayed
# green — blind in its headline scenario, since RFC3339Nano is exactly what someone reaching
# for fractional-second support would add. A guard that cannot fail is read as coverage.
# Structural, so a new spelling cannot hide: enter the function, enter the `[]string{`
# composite literal, take every non-empty line until the closing brace, and strip only
# indentation, a trailing line comment, and the trailing comma. POSIX awk throughout —
# `[ \t]` rather than a character class, and no GNU-only regex construct anywhere, because
# the predecessor's `\|` BRE alternation matched nothing under BSD/macOS sed and turned the
# whole maintainer gate red there (REQ-216's macOS lesson, one file over).
board_timestamp_layout_lines="$(awk '
  /^func parseTimestamp\(/ { inside_function = 1 }
  inside_function && /\[\]string/ { inside_layout_slice = 1; next }
  inside_layout_slice && /^[ \t]*\}/ { exit }
  inside_layout_slice {
    layout_element = $0
    sub(/[ \t]*\/\/.*$/, "", layout_element)
    sub(/^[ \t]+/, "", layout_element)
    sub(/,[ \t]*$/, "", layout_element)
    if (layout_element != "") print layout_element
  }
' "$repo_root/skills/do-work-board/tools/queue-kanban/model.go")"
board_timestamp_layout_count="$(printf '%s\n' "$board_timestamp_layout_lines" | grep -c . || true)"
# Anti-vacuity: an extraction that finds nothing would compare equal to nothing and could
# never fail. Assert it found layouts before trusting what it says about them.
[ "${board_timestamp_layout_count:-0}" -gt 0 ] \
  || fail_case "repair-req-timestamps read-side-layout case extracted ZERO layouts from parseTimestamp — the slice moved or was restructured, so this guard is blind rather than passing; fix the extraction before trusting a green run"
board_timestamp_layouts="$(printf '%s\n' "$board_timestamp_layout_lines" | tr '\n' ' ')"
[ "$board_timestamp_layouts" = 'time.RFC3339 "2006-01-02T15:04:05" "2006-01-02 15:04:05" "2006-01-02" ' ] \
  || fail_case "repair-req-timestamps read-side-layout case: the board's parseTimestamp layouts are now [$board_timestamp_layouts] — re-decide what repair-req-timestamps.sh refuses before widening the read side"

# repair-req-timestamps: a duplicated anchor key follows the last occurrence,
# exactly like the read side (the board's YAML dedup keeps the LAST value of a
# repeated top-level key) — a later-then-earlier claimed_at pair is a real
# ordering defect on the board and must not be reported clean; and a future
# FIRST occurrence shadowed by a clean last one is invisible to every YAML
# reader and must stay untouched (REQ-255, PR #145 external review).
repair_duplicate_project="$fixture_root/repair-duplicate-project"
mkdir -p "$repair_duplicate_project/do-work/working"
printf -- '---\nid: REQ-815\nstatus: claimed\ncreated_at: 2026-08-10T12:00:00Z\nclaimed_at: 2026-08-11T12:00:00Z\nclaimed_at: 2026-08-01T09:00:00Z\n---\nbody\n' \
  > "$repair_duplicate_project/do-work/working/REQ-815-duplicate-anchor.md"
printf -- '---\nid: REQ-816\nstatus: pending\ncreated_at: 2026-08-01T09:00:00Z\nblocked_at: 2093-01-01T00:00:00Z\nblocked_at: 2026-08-02T09:00:00Z\n---\nbody\n' \
  > "$repair_duplicate_project/do-work/working/REQ-816-shadowed-first.md"
TZ=UTC touch -m -t 202608121200.00 "$repair_duplicate_project/do-work/working/"REQ-81*.md
cp "$repair_duplicate_project/do-work/working/REQ-816-shadowed-first.md" "$fixture_root/repair-shadowed-before.md"
repair_duplicate_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_duplicate_project")" \
  || fail_case 'repair-req-timestamps duplicate-anchor case returned nonzero'
grep -q '^claimed_at: 2026-08-12T12:00:00Z$' "$repair_duplicate_project/do-work/working/REQ-815-duplicate-anchor.md" \
  || fail_case 'repair-req-timestamps duplicate-anchor case reported clean instead of repairing the effective (last) occurrence'
grep -q '^claimed_at: 2026-08-11T12:00:00Z$' "$repair_duplicate_project/do-work/working/REQ-815-duplicate-anchor.md" \
  || fail_case 'repair-req-timestamps duplicate-anchor case rewrote the shadowed first occurrence'
cmp -s "$fixture_root/repair-shadowed-before.md" "$repair_duplicate_project/do-work/working/REQ-816-shadowed-first.md" \
  || fail_case 'repair-req-timestamps duplicate-anchor case touched a future first occurrence no YAML reader can see'
printf '%s' "$repair_duplicate_output" | grep -q 'REQ-815-duplicate-anchor.md claimed_at: 2026-08-01T09:00:00Z -> 2026-08-12T12:00:00Z' \
  || fail_case 'repair-req-timestamps duplicate-anchor case did not log the effective-occurrence correction'

# repair-req-timestamps: a file whose opening fence is never closed is refused
# whole, because the board's splitFrontmatter reports NO frontmatter for that
# shape and renders every line as body (REQ-267 I1). The second fixture is the
# shape that could wedge the repair permanently: the file ends on the defective
# stamp with no trailing newline, so a last-line rewrite can never produce the
# final-newline diff pair the changed-line guard expects — the repair was
# rejected and the script exited 1 on EVERY run, with nothing able to heal it.
# The run must therefore be clean and silent, not merely non-destructive.
repair_unterminated_project="$fixture_root/repair-unterminated-project"
mkdir -p "$repair_unterminated_project/do-work/queue"
printf -- '---\nid: REQ-823\nstatus: pending\n\n# Body prose, and the fence above was never closed\n\ncreated_at: 2093-01-01T00:00:00Z\n' \
  > "$repair_unterminated_project/do-work/queue/REQ-823-unterminated-body.md"
printf -- '---\nid: REQ-824\nstatus: pending\ncreated_at: 2093-01-01T00:00:00Z' \
  > "$repair_unterminated_project/do-work/queue/REQ-824-unterminated-eof.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_unterminated_project/do-work/queue/"REQ-82*.md
for unterminated_fixture in REQ-823-unterminated-body REQ-824-unterminated-eof; do
  cp "$repair_unterminated_project/do-work/queue/$unterminated_fixture.md" "$fixture_root/$unterminated_fixture-before.md"
done
repair_unterminated_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_unterminated_project")" \
  || fail_case 'repair-req-timestamps unterminated-fence case exited nonzero — a shape that fails every session with no self-heal is exactly the defect'
[ -z "$repair_unterminated_output" ] \
  || fail_case 'repair-req-timestamps unterminated-fence case printed output for a file the read side sees no frontmatter in'
for unterminated_fixture in REQ-823-unterminated-body REQ-824-unterminated-eof; do
  cmp -s "$fixture_root/$unterminated_fixture-before.md" "$repair_unterminated_project/do-work/queue/$unterminated_fixture.md" \
    || fail_case "repair-req-timestamps unterminated-fence case rewrote $unterminated_fixture — the board reads that whole file as body text"
done

# repair-req-timestamps: a stamp padded INSIDE its quotes is repaired, because
# the read side unquotes and then trims and so parses and future-badges it
# (REQ-267 I2). The non-ASCII-padded sibling pins the stated boundary of that
# trim — this script matches bytes under LC_ALL=C, so a U+00A0-padded value the
# read side still parses stays refused byte-identical, and the header says so.
repair_padded_quote_project="$fixture_root/repair-padded-quote-project"
mkdir -p "$repair_padded_quote_project/do-work/queue"
printf -- '---\nid: REQ-825\nstatus: pending\ncreated_at: "2093-01-01 00:00:00 "\n---\nbody\n' \
  > "$repair_padded_quote_project/do-work/queue/REQ-825-padded-quote.md"
printf -- '---\nid: REQ-826\nstatus: pending\ncreated_at: "2093-01-01T00:00:00Z\xc2\xa0"\n---\nbody\n' \
  > "$repair_padded_quote_project/do-work/queue/REQ-826-non-ascii-pad.md"
TZ=UTC touch -m -t 202608101200.00 "$repair_padded_quote_project/do-work/queue/"REQ-82*.md
cp "$repair_padded_quote_project/do-work/queue/REQ-826-non-ascii-pad.md" "$fixture_root/repair-non-ascii-pad-before.md"
repair_padded_quote_output="$("$core_scripts/repair-req-timestamps.sh" "$repair_padded_quote_project")" \
  || fail_case 'repair-req-timestamps padded-quote case returned nonzero'
grep -q '^created_at: 2026-08-10T12:00:00Z$' "$repair_padded_quote_project/do-work/queue/REQ-825-padded-quote.md" \
  || fail_case 'repair-req-timestamps padded-quote case refused a padded quoted instant the board parses and future-badges'
cmp -s "$fixture_root/repair-non-ascii-pad-before.md" "$repair_padded_quote_project/do-work/queue/REQ-826-non-ascii-pad.md" \
  || fail_case 'repair-req-timestamps padded-quote case repaired a non-ASCII-padded value — the header states that residual is refused'
printf '%s' "$repair_padded_quote_output" | grep -q 'REQ-825-padded-quote.md created_at: "2093-01-01 00:00:00 " -> 2026-08-10T12:00:00Z' \
  || fail_case 'repair-req-timestamps padded-quote case did not report the full padded old value in the audit line'
printf '%s' "$repair_padded_quote_output" | grep -q 'REQ-826' \
  && fail_case 'repair-req-timestamps padded-quote case logged a correction for the refused non-ASCII-padded value'

# The retired implementation's two awk-failure probes now belong to the Go
# requestmodel/doctor unit tests. Keeping PATH-injected awk assertions here would
# test a dependency the compatibility launcher deliberately no longer has.

# The retired post-rename `git diff --numstat` seam is replaced by the Go plan
# delta guard, exercised directly in doctor_repair_test.go. Pre-existing dirty
# lines are therefore excluded from the repair's own change budget.

# repair-req-timestamps: an unreadable HEAD blob size is a failed inspection, not an
# absent baseline. Both used to collapse to 0 behind `|| echo 0`, which skipped the
# truncation floor — the guard that refuses to repair a stamp inside a file that already
# lost content (REQ-268 review, finding I1). A tracked path with no blob in HEAD is the
# real absence, and it must still repair.
repair_floor_project="$fixture_root/repair-floor-project"
fixture_repo_init "$repair_floor_project"
mkdir -p "$repair_floor_project/do-work/queue"
printf -- '---\nid: REQ-924\nstatus: pending\ncreated_at: 2093-01-01T00:00:00Z\n---\nbody\n' \
  > "$repair_floor_project/do-work/queue/REQ-924-future.md"
fixture_repo_commit_all "$repair_floor_project" fixture
cp "$repair_floor_project/do-work/queue/REQ-924-future.md" "$fixture_root/repair-floor-before.md"
repair_floor_bin="$fixture_root/repair-floor-bin"
mkdir -p "$repair_floor_bin"
{
  printf '%s\n' '#!/usr/bin/env bash'
  printf 'real_git=%q\n' "$(command -v git)"
  # `cat-file -e` still answers (the blob exists); only the size read fails.
  printf '%s\n' \
    'if [ "${1:-}" = "-C" ] && [ "${3:-}" = "cat-file" ] && [ "${4:-}" = "-s" ]; then' \
    '  exit 9' \
    'fi' \
    'exec "$real_git" "$@"'
} > "$repair_floor_bin/git"
chmod +x "$repair_floor_bin/git"
repair_floor_output="$(PATH="$repair_floor_bin:$PATH" "$core_scripts/repair-req-timestamps.sh" "$repair_floor_project" 2>&1)" \
  && fail_case 'repair-req-timestamps unreadable-floor case exited zero with the truncation floor unchecked'
printf '%s' "$repair_floor_output" | grep -q 'truncation floor' \
  || fail_case 'repair-req-timestamps unreadable-floor case did not say which guard could not run'
cmp -s "$fixture_root/repair-floor-before.md" "$repair_floor_project/do-work/queue/REQ-924-future.md" \
  || fail_case 'repair-req-timestamps unreadable-floor case repaired a file whose truncation floor it never read'
# The real absence — tracked, staged, never committed — still has no floor and still repairs.
repair_uncommitted_project="$fixture_root/repair-uncommitted-project"
fixture_repo_init "$repair_uncommitted_project"
mkdir -p "$repair_uncommitted_project/do-work/queue"
printf -- '---\nid: REQ-925\nstatus: pending\ncreated_at: 2093-01-01T00:00:00Z\n---\nbody\n' \
  > "$repair_uncommitted_project/do-work/queue/REQ-925-future.md"
git -C "$repair_uncommitted_project" add -A
"$core_scripts/repair-req-timestamps.sh" "$repair_uncommitted_project" >/dev/null 2>&1 \
  || fail_case 'repair-req-timestamps uncommitted-baseline case failed a staged-but-uncommitted file that simply has no HEAD blob'
grep -q '^created_at: 2093-01-01T00:00:00Z$' "$repair_uncommitted_project/do-work/queue/REQ-925-future.md" \
  && fail_case 'repair-req-timestamps uncommitted-baseline case left a future stamp unrepaired because HEAD held no blob'

# repair-req-timestamps: the 2-minute future-skew constant stays in lock-step
# with the board's futureTimestampSkewAllowance — a fourth hand-kept copy of
# the same allowance, pinned the way the repo pins cause-clause pairs
# (REQ-255 rider; REQ-246 review nit).
grep -q '^const timestampFutureSkewAllowance = 2 \* time\.Minute$' \
  "$repo_root/skills/do-work/tools/do-work-cli/internal/doctor/doctor_timestamps.go" \
  || fail_case 'repair-req-timestamps skew-constant case: the canonical repairer allowance moved or changed'
grep -q '^const futureTimestampSkewAllowance = 2 \* time\.Minute$' \
  "$repo_root/skills/do-work-board/tools/queue-kanban/model.go" \
  || fail_case 'repair-req-timestamps skew-constant case: the board constant moved or changed — keep the two allowances in lock-step'

prescribed_shell_finish
