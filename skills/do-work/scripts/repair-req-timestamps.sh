#!/usr/bin/env bash
# Mechanically repair detectably wrong `*_at` stamps in the REQ files under
# do-work/queue/ and do-work/working/. No agent judgment anywhere in this path.
#
# WHY THIS EXISTS. Detection of a wrong timestamp was already everywhere on the
# READ side — the board's future-stamp badge and data warning, `do-work forensics`
# Check 11, the work action's takeover prompt — but repair depended on the same
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
#      futureTimestampSkewAllowance and `do-work forensics` Check 11 use. The rule
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
# SHAPES LEFT ALONE, deliberately. A value is only comparable when it reads
# `YYYY-MM-DD`, `YYYY-MM-DDTHH:MM:SS`, or `YYYY-MM-DDTHH:MM:SSZ` (a space may
# stand in for the `T`), optionally wrapped in one matching pair of quotes —
# the schema's YAML readers unquote, so a quoted future stamp is flagged
# read-side and must be repairable here too (its replacement is written in the
# canonical unquoted form the Timestamp rule prescribes). A numeric UTC offset,
# fractional seconds, or anything unparseable is NOT provably wrong without
# timezone arithmetic, so it is never touched — the conservative direction.
# Indented keys are skipped too: `estimate.calculated_at` is a nested field,
# and every other reader in this skill anchors frontmatter keys at column zero.
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

# Matches the board's futureTimestampSkewAllowance and forensics Check 11.
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
# ---------------------------------------------------------------------------
timestamp_repair_apply_mode="${timestamp_repair_apply_mode:-1}"
timestamp_repair_git_only="${timestamp_repair_git_only:-0}"
replacement_source_names='git blame or the file mtime'
[ "$timestamp_repair_git_only" -eq 1 ] && replacement_source_names='git blame'

project_root="${1:-${CLAUDE_PROJECT_DIR:-.}}"
work_root="$project_root/do-work"
failure_count=0
repair_count=0
pending_repair_count=0

report_failure() {
  printf 'do-work: FAILED to repair %s — %s The file is unchanged.\n' "$1" "$2"
  failure_count=$((failure_count + 1))
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

# A comparison key is a canonical `YYYY-MM-DDTHH:MM:SSZ` string, which orders
# correctly under a plain string comparison — that is what lets this script
# compare instants without any date-parsing dependency. Empty output means the
# value is not comparable and must be left alone.
comparison_key_for() {
  local raw_value="$1" separator_normalized
  # One matching pair of wrapping quotes is stripped first: YAML readers hand
  # the board the unquoted value, so a quoted stamp is comparable here too.
  case "$raw_value" in
    \"*\") raw_value="${raw_value#\"}"; raw_value="${raw_value%\"}" ;;
    \'*\') raw_value="${raw_value#\'}"; raw_value="${raw_value%\'}" ;;
  esac
  # A space separator is folded to `T` so the patterns below never have to
  # carry a literal space inside a bracket expression, where its quoting is
  # shell-dependent.
  separator_normalized="${raw_value/ /T}"
  case "$separator_normalized" in
    [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9])
      printf '%sT00:00:00Z' "$separator_normalized" ;;
    [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9])
      printf '%sZ' "$separator_normalized" ;;
    [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z)
      printf '%s' "$separator_normalized" ;;
  esac
  return 0
}

earlier_of() {
  if [[ "$1" < "$2" ]]; then printf '%s' "$1"; else printf '%s' "$2"; fi
}
later_of() {
  if [[ "$1" > "$2" ]]; then printf '%s' "$1"; else printf '%s' "$2"; fi
}

# Frontmatter scope is line 1 `---` up to the next `---`, and keys are anchored at
# column zero. A `status:` in body prose or an indented `  calculated_at:` under
# `estimate:` is therefore unreachable — which is what keeps this script off the
# nested fields no other reader in this skill treats as schema.
frontmatter_value_for() {
  awk -v field_name="$2" '
    NR == 1 && $0 == "---" { inside_frontmatter = 1; next }
    inside_frontmatter && $0 == "---" { exit }
    inside_frontmatter && index($0, field_name ":") == 1 {
      field_value = substr($0, length(field_name) + 2)
      sub(/^[ \t]+/, "", field_value)
      sub(/[ \t]+$/, "", field_value)
      print field_value
      exit
    }
  ' "$1"
}

# Emits one `<line-number>\t<field-name>\t<value-token>` row per top-level
# frontmatter key whose name ends in `_at`. The value token is the first
# whitespace-delimited word after the colon, so a trailing YAML comment is never
# part of it — and never part of what gets rewritten either.
extract_timestamp_fields() {
  awk '
    NR == 1 && $0 == "---" { inside_frontmatter = 1; next }
    inside_frontmatter && $0 == "---" { exit }
    inside_frontmatter {
      colon_index = index($0, ":")
      if (colon_index < 2) next
      field_name = substr($0, 1, colon_index - 1)
      if (field_name !~ /^[A-Za-z_][A-Za-z0-9_]*_at$/) next
      field_rest = substr($0, colon_index + 1)
      sub(/^[ \t]+/, "", field_rest)
      value_token = field_rest
      sub(/[ \t].*$/, "", value_token)
      if (value_token == "") next
      print NR "\t" field_name "\t" value_token
    }
  ' "$1"
}

git_available=0
git -C "$project_root" rev-parse --git-dir >/dev/null 2>&1 && git_available=1

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
  local ordered_field_names ordered_name planned_count awk_status
  local -a field_line_numbers=() field_names=() field_tokens=()
  local -a field_keys=() field_new_values=() field_sources=()

  # A symlinked REQ file is not a shape this skill writes, and rewriting one by
  # atomic rename would replace the link with a regular file.
  [ -L "$request_file" ] && return 0
  [ -f "$request_file" ] || return 0
  [ -s "$request_file" ] || return 0

  field_rows="$(extract_timestamp_fields "$request_file")"
  [ -n "$field_rows" ] || return 0

  field_count=0
  while IFS=$'\t' read -r line_number field_name value_token; do
    [ -n "${line_number:-}" ] || continue
    field_line_numbers[field_count]="$line_number"
    field_names[field_count]="$field_name"
    field_tokens[field_count]="$value_token"
    field_keys[field_count]="$(comparison_key_for "$value_token")"
    field_new_values[field_count]=''
    field_sources[field_count]=''
    field_count=$((field_count + 1))
  done <<< "$field_rows"
  [ "$field_count" -gt 0 ] || return 0

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
  if [ "$path_tracked" -eq 1 ]; then
    head_blob_bytes="$(git -C "$project_root" cat-file -s \
      "HEAD:$(git -C "$project_root" ls-files --full-name -- "$request_file")" 2>/dev/null || echo 0)"
    # The truncation floor: a timestamp repair changes single lines, so a file at
    # less than half its committed size lost content BEFORE this run — repairing
    # a stamp in the remains would help commit the loss.
    if [ "${head_blob_bytes:-0}" -gt 0 ] && [ "$((pre_edit_bytes * 2))" -lt "$head_blob_bytes" ]; then
      report_failure "$display_path" \
        "the file is $pre_edit_bytes bytes on disk but $head_blob_bytes bytes in HEAD — content was lost before this run; recover it first (git checkout HEAD -- <file>)."
      return 0
    fi
    read -r pending_insertions pending_deletions _ <<< "$(git -C "$project_root" \
      diff --numstat --no-renames HEAD -- "$request_file" 2>/dev/null)" || true
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
      plan_spec="$plan_spec${field_line_numbers[$index]}=${field_new_values[$index]};"
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
  # spacing, the new value, then everything after the old token (a trailing YAML
  # comment survives verbatim). Every other line streams through untouched.
  awk -v plan_spec="$plan_spec" '
    BEGIN {
      plan_entry_count = split(plan_spec, plan_entries, ";")
      for (plan_index = 1; plan_index <= plan_entry_count; plan_index++) {
        if (plan_entries[plan_index] == "") continue
        separator_index = index(plan_entries[plan_index], "=")
        planned_line = substr(plan_entries[plan_index], 1, separator_index - 1) + 0
        planned_value[planned_line] = substr(plan_entries[plan_index], separator_index + 1)
      }
    }
    NR in planned_value {
      colon_index = index($0, ":")
      line_prefix = substr($0, 1, colon_index)
      line_rest = substr($0, colon_index + 1)
      match(line_rest, /^[ \t]*/)
      value_spacing = substr(line_rest, 1, RLENGTH)
      token_and_suffix = substr(line_rest, RLENGTH + 1)
      match(token_and_suffix, /^[^ \t]+/)
      line_suffix = substr(token_and_suffix, RLENGTH + 1)
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
    while IFS=$'\t' read -r line_number field_name value_token; do
      [ -n "${line_number:-}" ] || continue
      index=0
      while [ "$index" -lt "$field_count" ]; do
        if [ "${field_line_numbers[$index]}" = "$line_number" ] && [ -n "${field_new_values[$index]}" ] && \
          [ "$value_token" != "${field_new_values[$index]}" ]; then
          guard_verdict="line $line_number reads '$value_token' after the rewrite; expected '${field_new_values[$index]}'"
        fi
        index=$((index + 1))
      done
    done <<< "$(extract_timestamp_fields "$temp_file")"
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
  if [ "$path_tracked" -eq 1 ]; then
    read -r post_insertions post_deletions _ <<< "$(git -C "$project_root" \
      diff --numstat --no-renames HEAD -- "$request_file" 2>/dev/null)" || true
    case "${post_insertions:-0}${post_deletions:-0}" in
      *[!0-9]*) post_insertions=0; post_deletions=0 ;;
    esac
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
