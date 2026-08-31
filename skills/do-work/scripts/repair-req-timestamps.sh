#!/usr/bin/env bash
# Mechanically repair detectably wrong `*_at` stamps in the REQ files under
# do-work/queue/ and do-work/working/. No agent judgment anywhere in this path.
#
# WHY THIS EXISTS. Detection of a wrong timestamp was already everywhere on the
# READ side — the board's future-stamp badge and data warning, doctor's
# `TIMESTAMP-FUTURE` finding, the work action's takeover prompt — but repair depended on the same
# agent that wrote the bad stamp. A fabricated or timezone-shifted stamp therefore
# survived for as long as that agent stayed unaware of it, and every elapsed-time
# reading built on it (queue wait, claim stopwatch, implementation span) stayed
# wrong. This script closes that loop without asking anyone: it runs from the
# SessionStart hook, before an agent or a board render ever sees the file.
#
#   bash scripts/repair-req-timestamps.sh [project-root]
#
# WHAT COUNTS AS WRONG. Only provable wrongness, never a stamp that is merely
# surprising:
#
#   1. FUTURE — any top-level frontmatter field whose name ends in `_at` whose
#      value parses to later than now + 2 minutes. The 2-minute allowance absorbs
#      ordinary clock skew and is the same one the board's
#      futureTimestampSkewAllowance and doctor's `TIMESTAMP-FUTURE` finding use. The rule
#      is the `_at` SUFFIX, not a list of field names: a schema that grows a new
#      stamp field is covered the day it is added, with nothing to remember.
#   2. IMPOSSIBLE ORDERING — `claimed_at` earlier than `created_at`, or
#      `completed_at` earlier than `claimed_at`. A request cannot be claimed
#      before it exists, or completed before it was claimed.
#
# WHERE THE REPLACEMENT VALUE COMES FROM, decided by file state and not by
# preference:
#
#   - The file differs from HEAD (or is untracked, or there is no git) → the
#     file's own mtime. For a REQ still being edited that IS the write instant.
#   - The file matches HEAD → the author time of the commit that introduced that
#     stamp's line, read with `git blame`. That is the closest mechanical record
#     of when the line was really written.
#
# The derived value is then clamped so the repaired set satisfies
# `created_at <= claimed_at <= completed_at <= now`. An ordering repair rewrites
# the LATER field of the offending pair — the earlier one is the anchor.
#
# SHAPES LEFT ALONE, deliberately — and the parity rule that scopes the list:
# a value is only comparable when it reads `YYYY-MM-DD`, `YYYY-MM-DDTHH:MM:SS`,
# or `YYYY-MM-DDTHH:MM:SSZ` (a space may stand in for the `T`; the whole value
# after the colon is read, comment-aware, so a space-separated instant is
# repaired whole — never split at the space and half-rewritten), optionally
# wrapped in one matching pair of quotes with any ASCII whitespace padding
# inside them trimmed — the schema's YAML readers unquote AND THEN TRIM, so
# `created_at: "2093-01-01 00:00:00 "` is a value the board parses and
# future-badges, and it must be repairable here too (its replacement is written
# in the canonical unquoted form the Timestamp rule prescribes). Everything
# else is REFUSED byte-identical, never half-rewritten:
#   - A numeric UTC offset (`2093-01-01T00:00:00+02:00`) or fractional seconds
#     (`2093-01-01T00:00:00.500Z`): refused PERMANENTLY — a settled answer, not
#     a gap waiting on someone. The read side parses both, so a future one
#     keeps its board badge while this script leaves it alone. That residual
#     is real and it is accepted, because the arithmetic that would close it is
#     the risk: comparison here is dependency-free string ordering, and
#     normalizing an offset means carrying civil-date arithmetic in this file.
#     `2026-08-19T00:29:11+05:00` denotes 2026-08-18T19:29:11Z; a repairer that
#     reads the wall clock and ignores the offset sees a value five hours later
#     than the instant, which is how a CORRECT stamp gets erased as
#     future-dated. Refusing can only fail to fix; repairing can destroy,
#     unattended, from a SessionStart hook. The population is close to empty
#     besides: what this script exists to repair is a local wall clock stamped
#     `Z`, or a fabricated value, and neither carries an offset or a fraction —
#     a value that carries one came from a formatter that got the instant
#     right. Widening the read side re-opens this decision rather than
#     inheriting it, and both halves are pinned:
#     _dev/tests/prescribed-shell-scripts-behavior.sh fails if this refusal is
#     quietly dropped, or if the board's parseTimestamp layouts change under it.
#   - A shape-valid but calendar-impossible instant (a 99th month, April 31, a
#     non-leap-year February 29): the read-side parser rejects it, so erasing
#     it to a derived instant would destroy the malformed evidence the board
#     leaves visible for diagnosis.
#   - A value padded with a NON-ASCII space inside its quotes (U+00A0 and the
#     rest of Unicode's whitespace): Go's strings.TrimSpace removes those read
#     side, but this file matches bytes under LC_ALL=C and only trims ASCII
#     whitespace, so such a value stays refused. Refusing can only fail to fix.
#   - EVERY field of a file whose opening `---` fence is never closed. The
#     board's splitFrontmatter reports NO frontmatter for exactly that shape
#     and renders every line as body, so scanning to EOF here would rewrite
#     prose the read side never treats as schema. The refusal also closes the
#     one way this script could fail forever: when such a file ends on the
#     defective stamp with no trailing newline, the changed-line guard expects
#     the final-newline diff pair that a last-line rewrite can never produce,
#     so the repair was rejected and the script exited 1 — printing a FAILED
#     line into every session's start banner, with nothing able to heal it
#     (REQ-267). A closed fence always puts at least the fence line after a
#     repaired stamp, so no planned line can be the last line of the file.
#   - Anything else unparseable.
# Indented keys are skipped too: `estimate.calculated_at` is a nested field,
# and every other reader in this skill anchors frontmatter keys at column zero.
# A repeated top-level key is read by its LAST occurrence — the value every
# YAML reader effectively sees — and the shadowed earlier lines are invisible
# to those readers, so they are never examined and never rewritten. CRLF line
# endings and a leading UTF-8 BOM are tolerated and preserved, exactly as the
# board's splitFrontmatter tolerates them.
#
# SILENCE IS NEVER A CLEAN ANSWER. Every command and process substitution in this
# file is judged by its EXIT STATUS first and its content second: empty output from
# a tool that failed means "nothing was inspected", never "nothing was wrong". That
# is the condition, not a list of the sites that once broke it — a new substitution
# added below inherits it.
#
# REFUSALS ARE COUNTED ALWAYS AND PRINTED ON REQUEST. A refused value (the SHAPES
# LEFT ALONE list below) is never a repair and never a failure, so `refusal_count`
# rises while the exit status does not. Printing is opt-in through
# `timestamp_repair_voice_refusals` because the two callers want opposite things:
# this file runs unattended from the SessionStart hook, where a permanent refusal
# would print the same unhealable line into every session's start banner forever —
# REQ-267's wedge, and REQ-274's live complaint — while
# scripts/audit-archive-timestamps.sh is a deliberate human inspection where that
# same line IS the product, and swallowing it is what let an archive holding
# nothing but a refused stamp report as clean (REQ-268 instance 1).
#
# GUARD STYLE is deliberately `tools/checks/record-commit-hash.sh`'s, and that is
# a hard requirement rather than a stylistic nod: free-form frontmatter edits once
# truncated six archived REQ files to 0 bytes in a consumer repo. Every guard runs
# BEFORE the file is replaced, the replacement is an atomic rename, a tripped
# guard leaves the file byte-identical, and any failure exits nonzero.
#
# This script never stages or commits. Deletions and edits are left for the next
# housekeeping commit, exactly like scripts/cleanup-req-reservations.sh.
#
# Exit 0 — every scanned file is either clean or was repaired.
# Exit 1 — at least one file needed repair and did not get it. The named files are
#          byte-identical to how they were found.
set -uo pipefail

# ${#var} must count BYTES for the size arithmetic below, and byte-oriented
# matching keeps awk from erroring on an invalid multibyte sequence inside REQ
# prose. Only ASCII patterns are ever matched.
export LC_ALL=C

# Matches the board's futureTimestampSkewAllowance and doctor's
# `TIMESTAMP-FUTURE` finding.
future_stamp_skew_seconds=120

# ---------------------------------------------------------------------------
# Shared-machinery switches. scripts/audit-archive-timestamps.sh sources this
# file for the detection predicate, the derivation, and the guarded atomic
# rewrite, then runs its own archive scan — sourcing stops at the library
# return guard above the queue/working scan at the bottom. Sharing by sourcing
# is deliberate: the recognized value shapes live in comparison_key_for below,
# so widening that one function widens every tool built on this file at once
# instead of leaving a sibling with a narrower hand-rolled recognizer.
# The hook-run repairer keeps these defaults; a sourcing tool sets its own
# values BEFORE the `source`.
#
#   timestamp_repair_apply_mode  0 = plan and report only, write nothing
#   timestamp_repair_git_only    1 = replacements come from git blame alone;
#                                    the file-mtime fallback is disabled
#   timestamp_repair_voice_refusals
#                                1 = print one line per refused value; the count
#                                    is kept either way for the caller's summary
# ---------------------------------------------------------------------------
timestamp_repair_apply_mode="${timestamp_repair_apply_mode:-1}"
timestamp_repair_git_only="${timestamp_repair_git_only:-0}"
timestamp_repair_voice_refusals="${timestamp_repair_voice_refusals:-0}"
replacement_source_names='git blame or the file mtime'
[ "$timestamp_repair_git_only" -eq 1 ] && replacement_source_names='git blame'

project_root="${1:-${CLAUDE_PROJECT_DIR:-.}}"
work_root="$project_root/do-work"
failure_count=0
repair_count=0
pending_repair_count=0

refusal_count=0

report_failure() {
  printf 'do-work: FAILED to repair %s — %s The file is unchanged.\n' "$1" "$2"
  failure_count=$((failure_count + 1))
}

# A value this file recognizes as a stamp field but deliberately will not touch (a
# numeric offset, fractional seconds, a calendar-impossible instant, non-ASCII
# padding — the SHAPES LEFT ALONE list in the header). Always counted, printed only
# when the caller asked to hear it; never touches the exit status either way.
report_refusal() {
  refusal_count=$((refusal_count + 1))
  [ "$timestamp_repair_voice_refusals" -eq 1 ] || return 0
  printf 'do-work: refused %s %s (%s) — not a shape this repairer will rewrite; left byte-identical.\n' \
    "$1" "$2" "$3"
}

# Never write through a link: a symlinked do-work/ or scan directory could put an
# atomic rename anywhere on the host. The reservation cleaner refuses the same.
[ -L "$work_root" ] && exit 0
[ -d "$work_root" ] || exit 0

# ---------------------------------------------------------------------------
# Clock and epoch conversion. Both conversion spellings are PROBED against a
# known answer rather than guessed from the platform: GNU date wants -d @EPOCH
# and reads -r as a reference FILE, while BSD date wants -r EPOCH and reads -d as
# a daylight-saving flag — so a wrong guess on BSD would silently return the
# CURRENT time and the repair would write a value it never derived.
# ---------------------------------------------------------------------------
epoch_conversion_mode=''
if [ "$(date -u -d '@0' +%Y-%m-%dT%H:%M:%SZ 2>/dev/null)" = '1970-01-01T00:00:00Z' ]; then
  epoch_conversion_mode='gnu'
elif [ "$(date -u -r 0 +%Y-%m-%dT%H:%M:%SZ 2>/dev/null)" = '1970-01-01T00:00:00Z' ]; then
  epoch_conversion_mode='bsd'
fi

epoch_to_utc_stamp() {
  case "$epoch_conversion_mode" in
    gnu) date -u -d "@$1" +%Y-%m-%dT%H:%M:%SZ 2>/dev/null ;;
    bsd) date -u -r "$1" +%Y-%m-%dT%H:%M:%SZ 2>/dev/null ;;
    *) return 1 ;;
  esac
}

# The Timestamp rule's own command (actions/work-reference.md): the current UTC
# instant, read from the clock at the moment of use and never carried forward.
now_stamp="$(date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null)"
now_epoch="$(date -u +%s 2>/dev/null)"
case "${now_stamp:-}" in
  [0-9][0-9][0-9][0-9]-*) ;;
  *) printf 'do-work: cannot read the current UTC instant; no timestamp was inspected.\n' >&2; exit 1 ;;
esac
future_horizon_stamp="$now_stamp"
case "${now_epoch:-}" in
  '' | *[!0-9]*) ;;
  *) future_horizon_stamp="$(epoch_to_utc_stamp "$((now_epoch + future_stamp_skew_seconds))")" \
       || future_horizon_stamp="$now_stamp" ;;
esac

file_modified_epoch() {
  local modified_epoch
  modified_epoch="$(stat -c %Y -- "$1" 2>/dev/null)" || modified_epoch=''
  case "$modified_epoch" in
    '' | *[!0-9]*) modified_epoch="$(stat -f %m -- "$1" 2>/dev/null)" || modified_epoch='' ;;
  esac
  case "$modified_epoch" in
    '' | *[!0-9]*) return 1 ;;
  esac
  printf '%s' "$modified_epoch"
}

# The read-side parser (the board's parseTimestamp, backed by Go's time.Parse)
# rejects a shape-valid but calendar-impossible instant — 9999-99-99, April 31,
# February 29 outside a leap year — so treating one as comparable here would
# "repair" (erase) exactly the malformed evidence the board leaves visible for
# diagnosis. Components are validated against the real calendar before any
# string comparison; an impossible value is not comparable and stays untouched.
calendar_components_valid() {
  local canonical_stamp="$1" month_day_ceiling
  local year_number=$((10#${canonical_stamp:0:4})) month_number=$((10#${canonical_stamp:5:2}))
  local day_number=$((10#${canonical_stamp:8:2})) hour_number=$((10#${canonical_stamp:11:2}))
  local minute_number=$((10#${canonical_stamp:14:2})) second_number=$((10#${canonical_stamp:17:2}))
  [ "$month_number" -ge 1 ] && [ "$month_number" -le 12 ] || return 1
  [ "$hour_number" -le 23 ] || return 1
  [ "$minute_number" -le 59 ] || return 1
  [ "$second_number" -le 59 ] || return 1
  case "$month_number" in
    4 | 6 | 9 | 11) month_day_ceiling=30 ;;
    2)
      month_day_ceiling=28
      if [ $((year_number % 4)) -eq 0 ] && \
        { [ $((year_number % 100)) -ne 0 ] || [ $((year_number % 400)) -eq 0 ]; }; then
        month_day_ceiling=29
      fi
      ;;
    *) month_day_ceiling=31 ;;
  esac
  [ "$day_number" -ge 1 ] && [ "$day_number" -le "$month_day_ceiling" ]
}

# A comparison key is a canonical `YYYY-MM-DDTHH:MM:SSZ` string, which orders
# correctly under a plain string comparison — that is what lets this script
# compare instants without any date-parsing dependency. Empty output means the
# value is not comparable and must be left alone.
comparison_key_for() {
  local raw_value="$1" separator_normalized candidate_key
  # One matching pair of wrapping quotes is stripped first: YAML readers hand
  # the board the unquoted value, so a quoted stamp is comparable here too.
  case "$raw_value" in
    \"*\") raw_value="${raw_value#\"}"; raw_value="${raw_value%\"}" ;;
    \'*\') raw_value="${raw_value#\'}"; raw_value="${raw_value%\'}" ;;
  esac
  # The read side unquotes and THEN trims (coerceScalarToString and
  # parseTimestamp both run strings.TrimSpace), so padding inside the quotes —
  # `"2093-01-01 00:00:00 "` — is a value the board parses and future-badges.
  # Padding outside them was already trimmed by extract_timestamp_fields. Under
  # LC_ALL=C this trims the ASCII whitespace only; the header records the
  # non-ASCII residual.
  raw_value="${raw_value#"${raw_value%%[![:space:]]*}"}"
  raw_value="${raw_value%"${raw_value##*[![:space:]]}"}"
  # A space separator is folded to `T` so the patterns below never have to
  # carry a literal space inside a bracket expression, where its quoting is
  # shell-dependent. With whole-value extraction this is what makes the
  # board-parseable `YYYY-MM-DD HH:MM:SS` layout comparable — and repairable.
  separator_normalized="${raw_value/ /T}"
  candidate_key=''
  case "$separator_normalized" in
    [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9])
      candidate_key="${separator_normalized}T00:00:00Z" ;;
    [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9])
      candidate_key="${separator_normalized}Z" ;;
    [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z)
      candidate_key="$separator_normalized" ;;
    *) return 0 ;;
  esac
  calendar_components_valid "$candidate_key" || return 0
  printf '%s' "$candidate_key"
  return 0
}

earlier_of() {
  if [[ "$1" < "$2" ]]; then printf '%s' "$1"; else printf '%s' "$2"; fi
}
later_of() {
  if [[ "$1" > "$2" ]]; then printf '%s' "$1"; else printf '%s' "$2"; fi
}

# Frontmatter scope is line 1 `---` up to the next `---`, and keys are anchored
# at column zero. A `status:` in body prose or an indented `  calculated_at:`
# under `estimate:` is therefore unreachable — which is what keeps this script
# off the nested fields no other reader in this skill treats as schema.
#
# Emits one `<line-number>\t<field-name>\t<value>` row per top-level frontmatter
# key whose name ends in `_at`. The value is everything after the colon up to a
# trailing YAML comment (a `#` preceded by whitespace — the read-side YAML
# parsers use the same boundary), trimmed of surrounding whitespace, so a
# space-separated instant survives whole instead of truncating at the first
# space — the truncation is what once half-rewrote `2093-01-01 00:00:00` into
# an unparseable date-plus-phantom-suffix. The comment itself is never part of
# the value, and never part of what gets rewritten either. Fence matching
# tolerates what the board's splitFrontmatter tolerates: a UTF-8 BOM before the
# opening fence and a CRLF ending on any line — Windows agents are the
# likeliest source of both AND of the wrong local-time stamps this script
# exists to repair. The BOM and every CR live outside the value span, so a
# rewrite leaves them byte-for-byte in place.
#
# An opening fence that is never closed emits NOTHING: splitFrontmatter reports
# no frontmatter at all for that shape, so every line of such a file is body to
# the read side and none of it is this script's to rewrite. Rows are therefore
# buffered and only printed once the closing fence has been seen.
extract_timestamp_fields() {
  awk '
    BEGIN { utf8_bom = sprintf("%c%c%c", 239, 187, 191) }
    NR == 1 && index($0, utf8_bom) == 1 { $0 = substr($0, length(utf8_bom) + 1) }
    { line_body = $0; sub(/\r$/, "", line_body) }
    NR == 1 && line_body == "---" { inside_frontmatter = 1; next }
    inside_frontmatter && line_body == "---" { frontmatter_closed = 1; exit }
    inside_frontmatter {
      colon_index = index(line_body, ":")
      if (colon_index < 2) next
      field_name = substr(line_body, 1, colon_index - 1)
      if (field_name !~ /^[A-Za-z_][A-Za-z0-9_]*_at$/) next
      field_rest = substr(line_body, colon_index + 1)
      comment_start = match(field_rest, /[ \t]#/)
      if (comment_start > 0) field_rest = substr(field_rest, 1, comment_start - 1)
      sub(/^[ \t]+/, "", field_rest)
      sub(/[ \t]+$/, "", field_rest)
      if (field_rest == "") next
      buffered_row_count++
      buffered_rows[buffered_row_count] = NR "\t" field_name "\t" field_rest
    }
    END {
      if (!frontmatter_closed) exit
      for (row_index = 1; row_index <= buffered_row_count; row_index++) print buffered_rows[row_index]
    }
  ' "$1"
}

git_available=0
git -C "$project_root" rev-parse --git-dir >/dev/null 2>&1 && git_available=1
# Every guard below that compares against HEAD needs one to exist. A repo with no
# commit yet has no baseline at all — a real absence, not a git that failed to
# answer — and the two must not be folded together: reading the absence as a
# failure refuses to repair a perfectly ordinary staged-but-uncommitted REQ, and
# reading a failure as an absence is the fail-open this REQ exists to close.
head_commit_exists=0
[ "$git_available" -eq 1 ] && git -C "$project_root" rev-parse --verify -q HEAD >/dev/null 2>&1 \
  && head_commit_exists=1

# ---------------------------------------------------------------------------
# One REQ file. Detects, derives, guards, and rewrites; prints one audit line per
# correction. Returns nonzero only through report_failure's counter, never by
# leaving a half-written file behind.
# ---------------------------------------------------------------------------
repair_request_file() {
  local request_file="$1" display_path="$2"
  local field_rows field_count index
  local line_number field_name value_token comparison_key
  local derived_stamp derived_source clamped_stamp predecessor_key effective_key
  local plan_spec guard_verdict temp_file backup_file
  local expected_bytes expected_lines pre_edit_bytes post_edit_bytes pre_edit_lines post_edit_lines
  local changed_line_count expected_changed_lines byte_delta trailing_newline_added
  local file_matches_head path_tracked head_blob_bytes
  local pending_insertions pending_deletions post_insertions post_deletions
  local ordered_field_names ordered_name planned_count awk_status field_slot
  local extraction_status reparse_rows reparse_status verified_count
  local head_relative_path pending_numstat post_numstat numstat_status
  local -a field_line_numbers=() field_names=() field_tokens=()
  local -a field_keys=() field_new_values=() field_sources=()

  # A symlinked REQ file is not a shape this skill writes, and rewriting one by
  # atomic rename would replace the link with a regular file.
  [ -L "$request_file" ] && return 0
  [ -f "$request_file" ] || return 0
  [ -s "$request_file" ] || return 0

  field_rows="$(extract_timestamp_fields "$request_file")"
  extraction_status=$?
  if [ "$extraction_status" -ne 0 ]; then
    report_failure "$display_path" \
      "the frontmatter scan exited $extraction_status, so no field in this file was inspected."
    return 0
  fi
  [ -n "$field_rows" ] || return 0

  # A repeated top-level key keeps only its LAST occurrence, because that is
  # what every YAML reader on the read side effectively sees (the board's
  # duplicate-key recovery keeps the last value). Later rows overwrite the
  # earlier slot, so a shadowed first occurrence is never examined and never
  # rewritten — matching the readers is what makes a duplicated anchor's real
  # ordering defect detectable instead of reported clean.
  field_count=0
  while IFS=$'\t' read -r line_number field_name value_token; do
    [ -n "${line_number:-}" ] || continue
    field_slot="$field_count"
    index=0
    while [ "$index" -lt "$field_count" ]; do
      [ "${field_names[$index]}" = "$field_name" ] && field_slot="$index"
      index=$((index + 1))
    done
    [ "$field_slot" -eq "$field_count" ] && field_count=$((field_count + 1))
    field_line_numbers[field_slot]="$line_number"
    field_names[field_slot]="$field_name"
    field_tokens[field_slot]="$value_token"
    field_keys[field_slot]="$(comparison_key_for "$value_token")"
    field_new_values[field_slot]=''
    field_sources[field_slot]=''
  done <<< "$field_rows"
  [ "$field_count" -gt 0 ] || return 0

  # Every field whose value the recognizer would not accept. Voiced here, before any
  # defect pass, because the passes below skip an empty comparison key — which is
  # exactly how a file holding nothing but a refused stamp used to read as clean.
  index=0
  while [ "$index" -lt "$field_count" ]; do
    [ -n "${field_keys[$index]}" ] || \
      report_refusal "$display_path" "${field_names[$index]}" "${field_tokens[$index]}"
    index=$((index + 1))
  done

  path_tracked=0
  file_matches_head=0
  if [ "$git_available" -eq 1 ] && \
    git -C "$project_root" ls-files --error-unmatch -- "$request_file" >/dev/null 2>&1; then
    path_tracked=1
    if git -C "$project_root" diff --quiet HEAD -- "$request_file" 2>/dev/null; then
      file_matches_head=1
    fi
  fi

  # Replacement source, decided by file state. `git blame --line-porcelain` prints
  # the commit on its first line and `author-time <epoch>` in the header block; a
  # blame that cannot answer (shallow clone, an uncommitted line) falls back to the
  # mtime rather than inventing a value.
  derive_replacement_stamp() {
    local blame_line_number="$1" blame_output author_epoch blame_commit modified_epoch
    derived_stamp=''
    derived_source=''
    if [ "$file_matches_head" -eq 1 ]; then
      blame_output="$(git -C "$project_root" blame --line-porcelain \
        -L "$blame_line_number,$blame_line_number" -- "$request_file" 2>/dev/null)" || blame_output=''
      if [ -n "$blame_output" ]; then
        author_epoch="$(awk '/^author-time /{ epoch_value = $2 } END { print epoch_value }' <<< "$blame_output")"
        blame_commit="$(awk 'NR == 1 { print substr($1, 1, 7) }' <<< "$blame_output")"
        case "${author_epoch:-}" in
          '' | *[!0-9]*) ;;
          *)
            derived_stamp="$(epoch_to_utc_stamp "$author_epoch")" || derived_stamp=''
            [ -n "$derived_stamp" ] && derived_source="commit $blame_commit author time"
            ;;
        esac
      fi
    fi
    # Git-only mode (the archive audit): a checkout resets committed files'
    # mtimes, so for archived content the fallback below would read noise as
    # signal — an unanswerable blame must fail the derivation instead.
    if [ -z "$derived_stamp" ] && [ "$timestamp_repair_git_only" -eq 1 ]; then
      return 1
    fi
    if [ -z "$derived_stamp" ]; then
      modified_epoch="$(file_modified_epoch "$request_file")" || modified_epoch=''
      case "${modified_epoch:-}" in
        '' | *[!0-9]*) ;;
        *)
          derived_stamp="$(epoch_to_utc_stamp "$modified_epoch")" || derived_stamp=''
          [ -n "$derived_stamp" ] && derived_source='file mtime'
          ;;
      esac
    fi
    [ -n "$derived_stamp" ]
  }

  # --- Defect pass 1: any `*_at` later than now + the skew allowance. ---
  index=0
  while [ "$index" -lt "$field_count" ]; do
    comparison_key="${field_keys[$index]}"
    if [ -n "$comparison_key" ] && [[ "$comparison_key" > "$future_horizon_stamp" ]]; then
      if derive_replacement_stamp "${field_line_numbers[$index]}"; then
        field_new_values[$index]="$(earlier_of "$derived_stamp" "$now_stamp")"
        field_sources[$index]="$derived_source"
      else
        report_failure "$display_path" \
          "${field_names[$index]} is future-dated (${field_tokens[$index]}) but no replacement instant could be derived from $replacement_source_names."
        return 0
      fi
    fi
    index=$((index + 1))
  done

  # --- Defect pass 2: impossible ordering across the schema's three anchors. ---
  # Walked left to right on POST-repair values, so the LATER field of an offending
  # pair is the one rewritten: a request cannot be claimed before it exists, or
  # completed before it was claimed, and the earlier field is the anchor.
  predecessor_key=''
  ordered_field_names='created_at claimed_at completed_at'
  for ordered_name in $ordered_field_names; do
    index=0
    while [ "$index" -lt "$field_count" ]; do
      if [ "${field_names[$index]}" = "$ordered_name" ] && [ -n "${field_keys[$index]}" ]; then
        effective_key="${field_new_values[$index]:-${field_keys[$index]}}"
        if [ -n "$predecessor_key" ] && [[ "$effective_key" < "$predecessor_key" ]]; then
          if derive_replacement_stamp "${field_line_numbers[$index]}"; then
            clamped_stamp="$(earlier_of "$derived_stamp" "$now_stamp")"
            if [[ "$clamped_stamp" < "$predecessor_key" ]]; then
              # The derived instant still precedes its anchor; the clamp floor is
              # what makes the repaired set satisfy created <= claimed <= completed.
              clamped_stamp="$predecessor_key"
              field_sources[$index]="$derived_source, clamped to $predecessor_key"
            else
              field_sources[$index]="$derived_source"
            fi
            field_new_values[$index]="$clamped_stamp"
            effective_key="$clamped_stamp"
          else
            report_failure "$display_path" \
              "$ordered_name (${field_tokens[$index]}) precedes the field before it, but no replacement instant could be derived from $replacement_source_names."
            return 0
          fi
        fi
        predecessor_key="$effective_key"
        break
      fi
      index=$((index + 1))
    done
  done

  planned_count=0
  index=0
  while [ "$index" -lt "$field_count" ]; do
    [ -n "${field_new_values[$index]}" ] && planned_count=$((planned_count + 1))
    index=$((index + 1))
  done
  [ "$planned_count" -gt 0 ] || return 0

  # A planned value equal to what is already on the line would make the byte
  # arithmetic below expect a change diff cannot show. It cannot happen — a
  # future stamp is strictly later than its replacement, an ordering repair is
  # strictly later than the offending value — so treat one as the guard trip it
  # would be, not as a silent no-op.
  index=0
  while [ "$index" -lt "$field_count" ]; do
    if [ -n "${field_new_values[$index]}" ] && \
      [ "${field_new_values[$index]}" = "${field_tokens[$index]}" ]; then
      report_failure "$display_path" \
        "the derived replacement for ${field_names[$index]} equals the value already judged wrong (${field_tokens[$index]})."
      return 0
    fi
    index=$((index + 1))
  done

  # Report-only mode: the plan is printed and nothing is written. The caller's
  # exit contract still holds — a detectably wrong stamp exists unrepaired.
  if [ "$timestamp_repair_apply_mode" -eq 0 ]; then
    index=0
    while [ "$index" -lt "$field_count" ]; do
      if [ -n "${field_new_values[$index]}" ]; then
        printf 'do-work: would repair %s %s: %s -> %s (%s)\n' \
          "$display_path" "${field_names[$index]}" "${field_tokens[$index]}" \
          "${field_new_values[$index]}" "${field_sources[$index]}"
        pending_repair_count=$((pending_repair_count + 1))
      fi
      index=$((index + 1))
    done
    return 0
  fi

  # --- Pre-edit measurements the guards below compare against. ---
  pre_edit_bytes="$(wc -c < "$request_file" | tr -d '[:space:]')"
  pre_edit_lines="$(wc -l < "$request_file" | tr -d '[:space:]')"
  trailing_newline_added=0
  if [ -n "$(tail -c 1 "$request_file")" ]; then
    # awk's print restores the missing final newline; the size and diff guards
    # carry that byte explicitly, the same accounting record-commit-hash.sh uses.
    trailing_newline_added=1
  fi

  pending_insertions=0
  pending_deletions=0
  if [ "$path_tracked" -eq 1 ] && [ "$head_commit_exists" -eq 1 ]; then
    # The truncation floor: a timestamp repair changes single lines, so a file at
    # less than half its committed size lost content BEFORE this run — repairing
    # a stamp in the remains would help commit the loss.
    #
    # `-e` asks the question `-s` cannot answer alone: a tracked path with no blob
    # in HEAD (staged-but-never-committed) has no floor to compare against, which is
    # a real absence, while a blob that EXISTS and whose size will not read is a
    # failed inspection. Folded together behind `|| echo 0` those two were one
    # value, and the guard below silently skipped for both.
    head_relative_path="$(git -C "$project_root" ls-files --full-name -- "$request_file" 2>/dev/null)"
    if [ -n "$head_relative_path" ] && \
      git -C "$project_root" cat-file -e "HEAD:$head_relative_path" 2>/dev/null; then
      head_blob_bytes="$(git -C "$project_root" cat-file -s "HEAD:$head_relative_path" 2>/dev/null)"
      case "${head_blob_bytes:-}" in
        '' | *[!0-9]*)
          report_failure "$display_path" \
            "its HEAD blob exists but its size could not be read, so the truncation floor was never checked."
          return 0
          ;;
      esac
      if [ "$head_blob_bytes" -gt 0 ] && [ "$((pre_edit_bytes * 2))" -lt "$head_blob_bytes" ]; then
        report_failure "$display_path" \
          "the file is $pre_edit_bytes bytes on disk but $head_blob_bytes bytes in HEAD — content was lost before this run; recover it first (git checkout HEAD -- <file>)."
        return 0
      fi
    fi
    # The baseline the post-rename guard measures against. Run as its own command so
    # its exit status survives: consumed straight into `read`, only `read`'s status
    # was ever visible (and it is nonzero on the ordinary no-pending-changes case),
    # so a git that could not answer became a silent 0/0 baseline.
    pending_numstat="$(git -C "$project_root" diff --numstat --no-renames HEAD -- "$request_file" 2>/dev/null)"
    numstat_status=$?
    if [ "$numstat_status" -ne 0 ]; then
      report_failure "$display_path" \
        "its pending-change baseline could not be read (git diff --numstat exited $numstat_status); the rewrite was not attempted."
      return 0
    fi
    read -r pending_insertions pending_deletions _ <<< "$pending_numstat" || true
    case "${pending_insertions:-0}${pending_deletions:-0}" in
      *[!0-9]*) pending_insertions=0; pending_deletions=0 ;;
    esac
    pending_insertions="${pending_insertions:-0}"
    pending_deletions="${pending_deletions:-0}"
  fi

  byte_delta=0
  plan_spec=''
  index=0
  while [ "$index" -lt "$field_count" ]; do
    if [ -n "${field_new_values[$index]}" ]; then
      byte_delta=$((byte_delta + ${#field_new_values[$index]} - ${#field_tokens[$index]}))
      plan_spec="$plan_spec${field_line_numbers[$index]}=${#field_tokens[$index]}:${field_new_values[$index]};"
    fi
    index=$((index + 1))
  done
  expected_bytes=$((pre_edit_bytes + byte_delta + trailing_newline_added))
  expected_lines=$((pre_edit_lines + trailing_newline_added))
  expected_changed_lines=$((2 * planned_count + 2 * trailing_newline_added))

  # --- The edit. Temps live in the REQ's own directory so `mv` is a
  # same-filesystem atomic rename, and `cp -p` seeds the temp with the original's
  # mode so mktemp's 0600 cannot ride the rename onto the REQ file. ---
  temp_file="$(mktemp "$(dirname "$request_file")/.repair-req-timestamps.XXXXXX")" || {
    report_failure "$display_path" 'cannot create a temp file beside it — is the directory writable?'
    return 0
  }
  backup_file="$(mktemp "$(dirname "$request_file")/.repair-req-timestamps.orig.XXXXXX")" || {
    rm -f -- "$temp_file"
    report_failure "$display_path" 'cannot create a backup file beside it.'
    return 0
  }
  if ! cp -p "$request_file" "$backup_file" || ! cp -p "$request_file" "$temp_file"; then
    rm -f -- "$temp_file" "$backup_file"
    report_failure "$display_path" 'could not seed the temp files from the original.'
    return 0
  fi

  # Only the planned lines are rebuilt: prefix through the colon, the original
  # spacing, the new value, then everything after the old value's byte span (a
  # trailing YAML comment — and a CRLF ending's carriage return — survives
  # verbatim). The span length rides in the plan because the old value may
  # contain spaces; re-guessing a token boundary here is what once split a
  # space-separated instant. Every other line streams through untouched.
  awk -v plan_spec="$plan_spec" '
    BEGIN {
      plan_entry_count = split(plan_spec, plan_entries, ";")
      for (plan_index = 1; plan_index <= plan_entry_count; plan_index++) {
        if (plan_entries[plan_index] == "") continue
        separator_index = index(plan_entries[plan_index], "=")
        planned_line = substr(plan_entries[plan_index], 1, separator_index - 1) + 0
        plan_entry_rest = substr(plan_entries[plan_index], separator_index + 1)
        length_boundary = index(plan_entry_rest, ":")
        planned_old_length[planned_line] = substr(plan_entry_rest, 1, length_boundary - 1) + 0
        planned_value[planned_line] = substr(plan_entry_rest, length_boundary + 1)
      }
    }
    NR in planned_value {
      colon_index = index($0, ":")
      line_prefix = substr($0, 1, colon_index)
      line_rest = substr($0, colon_index + 1)
      match(line_rest, /^[ \t]*/)
      value_spacing = substr(line_rest, 1, RLENGTH)
      token_and_suffix = substr(line_rest, RLENGTH + 1)
      line_suffix = substr(token_and_suffix, planned_old_length[NR] + 1)
      print line_prefix value_spacing planned_value[NR] line_suffix
      next
    }
    { print }
  ' "$request_file" > "$temp_file"
  awk_status=$?
  if [ "$awk_status" -ne 0 ]; then
    rm -f -- "$temp_file" "$backup_file"
    report_failure "$display_path" "the line rewrite failed (awk exit $awk_status)."
    return 0
  fi

  # --- Exact post-edit arithmetic. Any other changed byte fails this, and the
  # rejected temp is discarded before it can touch the file. ---
  post_edit_bytes="$(wc -c < "$temp_file" | tr -d '[:space:]')"
  post_edit_lines="$(wc -l < "$temp_file" | tr -d '[:space:]')"
  changed_line_count="$(diff "$request_file" "$temp_file" | grep -c '^[<>]' || true)"
  guard_verdict='ok'
  [ -s "$temp_file" ] || guard_verdict='the rewritten file is empty'
  [ "$post_edit_bytes" -eq "$expected_bytes" ] || guard_verdict="the rewrite changed $post_edit_bytes-byte/$expected_bytes-byte size arithmetic"
  [ "$post_edit_lines" -eq "$expected_lines" ] || guard_verdict='the rewrite changed the line count'
  [ "$changed_line_count" -eq "$expected_changed_lines" ] || guard_verdict="the rewrite touched $changed_line_count diff lines; expected $expected_changed_lines"
  # Re-parse the rewrite with the same extractor that planned it: every planned
  # field must now read exactly its planned value.
  if [ "$guard_verdict" = 'ok' ]; then
    # Materialized before the loop, and its status read: consumed through a
    # herestring the extractor could fail, produce zero rows, and leave this guard
    # verifying nothing while reporting ok. The verified count is what proves the
    # loop actually saw each planned line, rather than a silent no-iteration pass.
    reparse_rows="$(extract_timestamp_fields "$temp_file")"
    reparse_status=$?
    verified_count=0
    if [ "$reparse_status" -ne 0 ]; then
      guard_verdict="the post-rewrite re-parse exited $reparse_status, so the rewrite was never verified"
    else
      while IFS=$'\t' read -r line_number field_name value_token; do
        [ -n "${line_number:-}" ] || continue
        index=0
        while [ "$index" -lt "$field_count" ]; do
          if [ "${field_line_numbers[$index]}" = "$line_number" ] && [ -n "${field_new_values[$index]}" ]; then
            if [ "$value_token" != "${field_new_values[$index]}" ]; then
              guard_verdict="line $line_number reads '$value_token' after the rewrite; expected '${field_new_values[$index]}'"
            else
              verified_count=$((verified_count + 1))
            fi
          fi
          index=$((index + 1))
        done
      done <<< "$reparse_rows"
      [ "$guard_verdict" != 'ok' ] || [ "$verified_count" -eq "$planned_count" ] || \
        guard_verdict="the post-rewrite re-parse confirmed $verified_count of $planned_count planned field(s)"
    fi
  fi
  if [ "$guard_verdict" != 'ok' ]; then
    rm -f -- "$temp_file" "$backup_file"
    report_failure "$display_path" "$guard_verdict — the edit was rejected and discarded."
    return 0
  fi

  if ! mv "$temp_file" "$request_file"; then
    rm -f -- "$temp_file" "$backup_file"
    report_failure "$display_path" 'the atomic rename failed.'
    return 0
  fi

  # Post-rename numstat guard: against HEAD the pending delta may only have grown
  # by the planned lines. A failure here restores the pre-edit bytes from the
  # backup, so even this last guard leaves the file byte-identical.
  if [ "$path_tracked" -eq 1 ] && [ "$head_commit_exists" -eq 1 ]; then
    # Same separation as the pre-edit baseline, and it matters more here: this runs
    # AFTER the rename, so a git that could not answer used to fall back to 0/0, and
    # `[ 0 -gt threshold ]` is false for every threshold — the last guard standing
    # between a bad rewrite and a `repaired` line silently passed. A guard that
    # could not run is a tripped guard: restore the file and say so.
    post_numstat="$(git -C "$project_root" diff --numstat --no-renames HEAD -- "$request_file" 2>/dev/null)"
    numstat_status=$?
    read -r post_insertions post_deletions _ <<< "$post_numstat" || true
    case "${post_insertions:-0}${post_deletions:-0}" in
      *[!0-9]*) post_insertions=0; post_deletions=0 ;;
    esac
    if [ "$numstat_status" -ne 0 ]; then
      if cp -p "$backup_file" "$request_file"; then
        report_failure "$display_path" \
          "the post-rewrite diff check could not run (git diff --numstat exited $numstat_status); the file was RESTORED to its pre-edit content."
      else
        printf 'do-work: FAILED to repair %s — the post-rewrite diff check could not run AND the restore failed; recover with git checkout HEAD -- <file>.\n' "$display_path"
        failure_count=$((failure_count + 1))
      fi
      rm -f -- "$backup_file"
      return 0
    fi
    if [ "${post_insertions:-0}" -gt "$((pending_insertions + planned_count + trailing_newline_added))" ] || \
      [ "${post_deletions:-0}" -gt "$((pending_deletions + planned_count + trailing_newline_added))" ]; then
      if cp -p "$backup_file" "$request_file"; then
        report_failure "$display_path" \
          "git diff --numstat reported +${post_insertions}/-${post_deletions} against a +$pending_insertions/-$pending_deletions baseline; the file was RESTORED to its pre-edit content."
      else
        printf 'do-work: FAILED to repair %s — the numstat guard tripped AND the restore failed; recover with git checkout HEAD -- <file>.\n' "$display_path"
        failure_count=$((failure_count + 1))
      fi
      rm -f -- "$backup_file"
      return 0
    fi
  fi
  rm -f -- "$backup_file"

  # --- Audit trail: one line per corrected field. ---
  index=0
  while [ "$index" -lt "$field_count" ]; do
    if [ -n "${field_new_values[$index]}" ]; then
      printf 'do-work: repaired %s %s: %s -> %s (%s)\n' \
        "$display_path" "${field_names[$index]}" "${field_tokens[$index]}" \
        "${field_new_values[$index]}" "${field_sources[$index]}"
      repair_count=$((repair_count + 1))
    fi
    index=$((index + 1))
  done
  return 0
}

# ---------------------------------------------------------------------------
# Library mode ends here: when this file is sourced (scripts/
# audit-archive-timestamps.sh), the shared machinery above is defined and
# initialized and the sourcing tool owns its own scan, summary, and exit code.
# Only direct execution continues into the queue/working scan below.
# ---------------------------------------------------------------------------
if [ "${BASH_SOURCE[0]}" != "$0" ]; then
  return 0
fi

# ---------------------------------------------------------------------------
# The scan: queue + working only. Archived files are deliberately out of scope.
# ---------------------------------------------------------------------------
for scan_directory in "$work_root/queue" "$work_root/working"; do
  [ -L "$scan_directory" ] && continue
  [ -d "$scan_directory" ] || continue
  for request_file in "$scan_directory"/REQ-*.md; do
    [ -e "$request_file" ] || continue
    repair_request_file "$request_file" "${request_file#"$project_root"/}"
  done
done

if [ "$failure_count" -gt 0 ]; then
  exit 1
fi
if [ "$repair_count" -gt 0 ]; then
  printf 'do-work: repaired %s detectably wrong timestamp(s) — review and commit the correction(s) with the next housekeeping commit.\n' "$repair_count"
fi
exit 0
