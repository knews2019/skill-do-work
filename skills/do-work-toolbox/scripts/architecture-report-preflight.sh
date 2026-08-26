#!/usr/bin/env bash
# Pre-flight and publication for dated, immutable architecture reports.
#
# Two verbs, one for each end of `actions/architecture-report.md`:
#
#   --scan <reports-directory>
#       Emit the run's ground facts as `key=value` lines: the current commit, the report
#       slug, the unsuffixed candidate directory, and the newest prior report with the
#       watermark hash the drift scope is computed against.
#
#   --publish <draft-path> <candidate-directory>
#       Create the first free `-<n>` sibling of the candidate directory, put the finished
#       draft inside it as the report file, and print the published path. Never touches an
#       existing report.
#
# A report is one directory named `<yyyy-mm-dd>_<hhmm>_architecture-report` holding
# `architecture-report.md`, matching the bundle shape `ai-report` publishes beside it.
# Times are UTC (`date -u`), matching the suite's stamp convention, so the slug a run
# picks does not depend on the runner's timezone.
#
# Suffix escalation lives here and only here: `--scan` deliberately emits the unsuffixed
# slug so there is exactly one implementation of "first free path" to get right. The
# no-clobber rule it serves is
# `actions/completed-work-presentation-reference.md` -> Collision-Safe Publication.
#
# `mkdir` is the exclusive primitive: it fails when the path already exists, so the
# directory a run creates is one no other run had, without a check-then-write window.
set -u

usage() {
  printf 'Usage: %s --scan <reports-directory>\n' "$0" >&2
  printf '       %s --publish <draft-path> <candidate-directory>\n' "$0" >&2
  exit 2
}

report_name_suffix='_architecture-report'
report_file_name='architecture-report.md'

# Prints "<sortable-slug> <sequence>" for a report directory name, or nothing when the
# name is not a report this action publishes. The sequence is what makes the ordering
# numeric: `-10` sorts after `-2` here, where a lexical sort puts it before. The
# `yyyy-mm-dd_hhmm` stem is fixed-width, so it sorts lexically as itself.
parse_report_basename() {
  local basename_text="$1"
  local stem_text
  local sequence_part

  case "$basename_text" in
    *"$report_name_suffix") stem_text="${basename_text%"$report_name_suffix"}"; sequence_part=1 ;;
    *"$report_name_suffix"-*)
      sequence_part="${basename_text##*"$report_name_suffix"-}"
      stem_text="${basename_text%"$report_name_suffix"-"$sequence_part"}"
      case "$sequence_part" in *[!0-9]* | '') return 1 ;; esac
      ;;
    *) return 1 ;;
  esac

  # `yyyy-mm-dd_hhmm` exactly: reject anything else so a hand-named directory beside the
  # reports cannot become the baseline a run verifies against.
  case "$stem_text" in
    [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]_[0-9][0-9][0-9][0-9]) ;;
    *) return 1 ;;
  esac
  printf '%s %s\n' "$stem_text" "$sequence_part"
}

scan_reports() {
  local reports_directory="$1"
  local head_hash
  local report_slug
  local candidate_path
  local prior_report=''
  local prior_hash=''
  local prior_hash_resolves='n/a'
  local newest_key=''
  local report_path
  local report_key
  local watermark_line
  local prior_report_file

  # Ask the question that has a real answer first: a repository without a resolvable HEAD
  # cannot be watermarked, and an empty hash must never reach the report as one.
  if ! head_hash="$(git rev-parse --short HEAD 2>/dev/null)" || [ -z "$head_hash" ]; then
    printf 'architecture-report-preflight: no resolvable HEAD commit to watermark against.\n' >&2
    return 2
  fi
  if ! report_slug="$(date -u +%Y-%m-%d_%H%M)" || [ -z "$report_slug" ]; then
    printf 'architecture-report-preflight: the date command produced no UTC timestamp.\n' >&2
    return 2
  fi
  candidate_path="$reports_directory/$report_slug$report_name_suffix"

  if [ -d "$reports_directory" ]; then
    for report_path in "$reports_directory"/*"$report_name_suffix" "$reports_directory"/*"$report_name_suffix"-*; do
      [ -d "$report_path" ] || continue
      report_key="$(parse_report_basename "${report_path##*/}")" || continue
      # Pad the sequence so `002` orders below `010`; the stem is already fixed-width.
      report_key="$(printf '%s %012d' "${report_key%% *}" "${report_key##* }")"
      if [ -z "$newest_key" ] || [ "$report_key" \> "$newest_key" ]; then
        newest_key="$report_key"
        prior_report="$report_path/$report_file_name"
      fi
    done
  fi

  if [ -n "$prior_report" ]; then
    # `unreadable`, never empty: an unparseable or missing watermark means every prior
    # claim must be re-verified, which is the opposite decision from "no prior report".
    prior_hash='unreadable'
    prior_report_file="$prior_report"
    if [ -f "$prior_report_file" ]; then
      watermark_line="$(grep -m1 -E '^verified-at:[[:space:]]+[0-9a-f]{7,40}([[:space:]]|$)' "$prior_report_file")"
      if [ -n "$watermark_line" ]; then
        watermark_line="${watermark_line#verified-at:}"
        watermark_line="${watermark_line#"${watermark_line%%[![:space:]]*}"}"
        prior_hash="${watermark_line%%[![:alnum:]]*}"
      fi
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
  printf 'report_slug=%s\n' "$report_slug"
  printf 'report_candidate=%s\n' "$candidate_path"
  printf 'prior_report=%s\n' "$prior_report"
  printf 'prior_hash=%s\n' "$prior_hash"
  printf 'prior_hash_resolves=%s\n' "$prior_hash_resolves"
}

publish_report() {
  local draft_path="$1"
  local candidate_path="$2"
  local sequence_number=1
  local published_directory
  local published_path

  if [ ! -f "$draft_path" ]; then
    printf 'architecture-report-preflight: draft is not a regular file: %s\n' "$draft_path" >&2
    return 2
  fi
  if ! mkdir -p "$(dirname "$candidate_path")"; then
    printf 'architecture-report-preflight: reports directory could not be created: %s\n' \
      "$(dirname "$candidate_path")" >&2
    return 2
  fi

  while :; do
    if [ "$sequence_number" -eq 1 ]; then
      published_directory="$candidate_path"
    else
      published_directory="${candidate_path}-${sequence_number}"
    fi

    # Plain `mkdir`, never `mkdir -p`: -p succeeds on an existing directory, which would
    # publish this run's report into the previous run's bundle.
    if mkdir "$published_directory" 2>/dev/null; then
      published_path="$published_directory/$report_file_name"
      if ! cp "$draft_path" "$published_path"; then
        printf 'architecture-report-preflight: report could not be written: %s\n' \
          "$published_path" >&2
        rmdir "$published_directory" 2>/dev/null
        return 1
      fi
      if ! cmp -s "$draft_path" "$published_path"; then
        printf 'architecture-report-preflight: published report does not match the draft: %s\n' \
          "$published_path" >&2
        return 1
      fi
      printf '%s\n' "$published_path"
      return 0
    fi
    if [ -e "$published_directory" ] || [ -L "$published_directory" ]; then
      sequence_number=$((sequence_number + 1))
      continue
    fi
    printf 'architecture-report-preflight: report directory could not be created: %s\n' \
      "$published_directory" >&2
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
