#!/usr/bin/env bash
# Pre-flight and publication for dated, immutable architecture reports.
#
# Two verbs, one for each end of `actions/architecture-report.md`:
#
#   --scan <reports-directory>
#       Emit the run's ground facts as `key=value` lines: the current commit, the report
#       date, the unsuffixed candidate path, and the newest prior report with the
#       watermark hash the drift scope is computed against.
#
#   --publish <draft-path> <candidate-path>
#       Link a finished draft to the first free `_<n>` sibling of the candidate and print
#       the path it landed on. Never touches an existing report.
#
# Dates are UTC (`date -u`), matching the suite's stamp convention, so the filename a run
# picks does not depend on the runner's timezone.
#
# Suffix escalation lives here and only here: `--scan` deliberately emits the unsuffixed
# name so there is exactly one implementation of "first free path" to get right. The
# separator is `_`, continuing the filename's existing date separator rather than
# introducing a second one; the no-clobber rule it serves is
# `actions/completed-work-presentation-reference.md` -> Collision-Safe Publication.
#
# `ln` treats an existing directory operand as a container rather than a collision, so a
# successful link is verified after the fact instead of being taken on its exit status.
set -u

usage() {
  printf 'Usage: %s --scan <reports-directory>\n' "$0" >&2
  printf '       %s --publish <draft-path> <candidate-path>\n' "$0" >&2
  exit 2
}

report_name_prefix='architecture-report_'
report_name_suffix='.md'

# Prints "<date-digits> <sequence>" for a report basename, or nothing when the name is not
# a report this action publishes. The sequence is what makes the ordering numeric:
# `_10` sorts after `_2` here, where a lexical sort of the filenames puts it before.
parse_report_basename() {
  local basename_text="$1"
  local stem_text

  case "$basename_text" in
    "$report_name_prefix"*"$report_name_suffix") ;;
    *) return 1 ;;
  esac
  stem_text="${basename_text#"$report_name_prefix"}"
  stem_text="${stem_text%"$report_name_suffix"}"

  case "$stem_text" in
    *_*)
      local date_part="${stem_text%%_*}"
      local sequence_part="${stem_text#*_}"
      case "$date_part" in *[!0-9]* | '') return 1 ;; esac
      case "$sequence_part" in *[!0-9]* | '') return 1 ;; esac
      printf '%s %s\n' "$date_part" "$sequence_part"
      ;;
    *)
      case "$stem_text" in *[!0-9]* | '') return 1 ;; esac
      printf '%s 1\n' "$stem_text"
      ;;
  esac
}

scan_reports() {
  local reports_directory="$1"
  local head_hash
  local report_date
  local candidate_path
  local prior_report=''
  local prior_hash=''
  local prior_hash_resolves='n/a'
  local newest_key=''
  local report_path
  local report_key
  local watermark_line

  # Ask the question that has a real answer first: a repository without a resolvable HEAD
  # cannot be watermarked, and an empty hash must never reach the report as one.
  if ! head_hash="$(git rev-parse --short HEAD 2>/dev/null)" || [ -z "$head_hash" ]; then
    printf 'architecture-report-preflight: no resolvable HEAD commit to watermark against.\n' >&2
    return 2
  fi
  if ! report_date="$(date -u +%Y%m%d)" || [ -z "$report_date" ]; then
    printf 'architecture-report-preflight: the date command produced no UTC date.\n' >&2
    return 2
  fi
  candidate_path="$reports_directory/$report_name_prefix$report_date$report_name_suffix"

  if [ -d "$reports_directory" ]; then
    for report_path in "$reports_directory/$report_name_prefix"*"$report_name_suffix"; do
      [ -f "$report_path" ] || continue
      report_key="$(parse_report_basename "${report_path##*/}")" || continue
      # Zero-pad both fields to one comparable sort key: dates are already fixed-width,
      # and the sequence is padded so `002` orders below `010`.
      report_key="$(printf '%s %012d' "${report_key%% *}" "${report_key##* }")"
      if [ -z "$newest_key" ] || [ "$report_key" \> "$newest_key" ]; then
        newest_key="$report_key"
        prior_report="$report_path"
      fi
    done
  fi

  if [ -n "$prior_report" ]; then
    # `unreadable`, never empty: an unparseable watermark means every prior claim must be
    # re-verified, which is the opposite decision from "there is no prior report".
    prior_hash='unreadable'
    watermark_line="$(grep -m1 -E '^verified-at:[[:space:]]+[0-9a-f]{7,40}([[:space:]]|$)' "$prior_report")"
    if [ -n "$watermark_line" ]; then
      watermark_line="${watermark_line#verified-at:}"
      watermark_line="${watermark_line#"${watermark_line%%[![:space:]]*}"}"
      prior_hash="${watermark_line%%[![:alnum:]]*}"
    fi
    if [ "$prior_hash" = 'unreadable' ]; then
      prior_hash_resolves='no'
    elif git rev-parse --verify -q "$prior_hash^{commit}" >/dev/null 2>&1; then
      prior_hash_resolves='yes'
    else
      prior_hash_resolves='no'
    fi
  fi

  printf 'head_hash=%s\n' "$head_hash"
  printf 'report_date=%s\n' "$report_date"
  printf 'report_candidate=%s\n' "$candidate_path"
  printf 'prior_report=%s\n' "$prior_report"
  printf 'prior_hash=%s\n' "$prior_hash"
  printf 'prior_hash_resolves=%s\n' "$prior_hash_resolves"
}

publish_report() {
  local draft_path="$1"
  local candidate_path="$2"
  local candidate_directory
  local candidate_stem
  local candidate_extension
  local sequence_number=1
  local published_path
  local nested_path

  if [ ! -f "$draft_path" ]; then
    printf 'architecture-report-preflight: draft is not a regular file: %s\n' "$draft_path" >&2
    return 2
  fi
  candidate_directory="$(dirname "$candidate_path")"
  case "${candidate_path##*/}" in
    *.*)
      candidate_stem="${candidate_path%.*}"
      candidate_extension=".${candidate_path##*.}"
      ;;
    *)
      candidate_stem="$candidate_path"
      candidate_extension=''
      ;;
  esac
  if ! mkdir -p "$candidate_directory"; then
    printf 'architecture-report-preflight: reports directory could not be created: %s\n' \
      "$candidate_directory" >&2
    return 2
  fi

  while :; do
    if [ "$sequence_number" -eq 1 ]; then
      published_path="$candidate_path"
    else
      published_path="${candidate_stem}_${sequence_number}${candidate_extension}"
    fi

    if ln "$draft_path" "$published_path" 2>/dev/null; then
      nested_path="$published_path/${draft_path##*/}"
      if [ -e "$nested_path" ]; then
        # The candidate is a directory, so `ln` linked into it rather than colliding.
        rm -f "$nested_path"
        sequence_number=$((sequence_number + 1))
        continue
      fi
      if [ ! -f "$published_path" ]; then
        printf 'architecture-report-preflight: published report is not a regular file: %s\n' \
          "$published_path" >&2
        return 1
      fi
      printf '%s\n' "$published_path"
      return 0
    fi
    if [ -e "$published_path" ] || [ -L "$published_path" ]; then
      sequence_number=$((sequence_number + 1))
      continue
    fi
    printf 'architecture-report-preflight: report could not be published exclusively: %s\n' \
      "$published_path" >&2
    return 1
  done
}

case "${1:-}" in
  --scan)
    [ "$#" -eq 2 ] || usage
    scan_reports "$2"
    ;;
  --publish)
    [ "$#" -eq 3 ] || usage
    publish_report "$2" "$3"
    ;;
  *) usage ;;
esac
