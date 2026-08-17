#!/usr/bin/env bash
# estimate-p50.sh — deterministic P50 active-duration estimator for REQs.
# The agent extracts judgment signals from a REQ (route, write-set size,
# evidence requirements, ...); this script does the arithmetic, so the same
# normalized signals always produce the same estimate. Signal-extraction
# guidance and the frontmatter block template live in
# actions/estimate-reference.md — this script is the deterministic half.
#
# Estimate mode:
#   tools/estimate-p50.sh --route A|B|C [--write-set N] [--new-files N]
#     [--subsystems N] [--acceptance N] [--deps-depth N] [--browser]
#     [--persistence] [--async-behavior] [--performance] [--regression]
#     [--full-suite]
#   tools/estimate-p50.sh --trivial          # short-circuit: floor estimate
#
#   Prints p50_active_minutes (a multiple of five, never below the floor),
#   confidence (low|medium|high), and the basis lines for the estimate: block.
#   Independent review and ordinary remediation cost are folded into the route
#   base — they are not separate flags. No P80 or other percentiles, ever.
#
# Critical-path mode:
#   tools/estimate-p50.sh critical-path ID:MINUTES[:DEP[,DEP...]] ...
#
#   Prints total_estimated_effort_minutes (plain sum) and
#   critical_path_minutes (longest path through the depends_on graph — never
#   the sum of parallel branches). A dependency id not present among the
#   arguments contributes zero (treat archived/completed REQs that way). A
#   cycle is an error (exit 1) — the work action has its own cycle handling.
#
# Exit 0: estimate printed. Exit 1: invalid graph. Exit 2: usage error.
set -euo pipefail

# Scoring scale calibrated 2026-08-17 to the archive's measured actuals: route
# bases equal the per-route medians of 188 claimed_at→completed_at spans (>4h
# and negative spans excluded as assumed pauses/broken stamps); signal weights
# stretch heavy REQs toward the per-route p80. Provenance and re-fit method:
# actions/estimate-reference.md → Calibration.
floor_minutes=5

usage_error() {
  printf 'usage: %s --route A|B|C [signal flags] | --trivial | critical-path ID:MIN[:DEP,...] ...\n' "$0" >&2
  printf 'error: %s\n' "$1" >&2
  exit 2
}

round_to_nearest_five() {
  local raw_minutes="$1"
  printf '%s' $(( ((raw_minutes + 2) / 5) * 5 ))
}

# ---------------------------------------------------------------------------
# Critical-path mode
# ---------------------------------------------------------------------------
if [ "${1:-}" = "critical-path" ]; then
  shift
  [ "$#" -ge 1 ] || usage_error "critical-path needs at least one ID:MINUTES[:DEPS] argument"
  for graph_entry in "$@"; do
    case "$graph_entry" in
      *:*) ;;
      *) usage_error "malformed graph entry '$graph_entry' (want ID:MINUTES[:DEP,...])" ;;
    esac
  done
  printf '%s\n' "$@" | awk -F ':' '
    {
      node_id = $1
      node_minutes = $2
      if (node_minutes !~ /^[0-9]+$/) {
        printf "error: minutes for %s must be a non-negative integer\n", node_id > "/dev/stderr"
        exit_code = 2; exit 2
      }
      minutes_by_id[node_id] = node_minutes
      dependency_list[node_id] = (NF >= 3) ? $3 : ""
      total_minutes += node_minutes
    }
    function longest_path_ending_at(node_id,    best_prefix, dependency_count, dependency_ids, i, candidate) {
      if (node_id in memoized_path) return memoized_path[node_id]
      if (!(node_id in minutes_by_id)) return 0   # unknown dep: archived/complete, costs nothing here
      if (node_id in in_progress) {
        printf "error: dependency cycle involving %s\n", node_id > "/dev/stderr"
        exit_code = 1; exit 1
      }
      in_progress[node_id] = 1
      best_prefix = 0
      dependency_count = split(dependency_list[node_id], dependency_ids, ",")
      for (i = 1; i <= dependency_count; i++) {
        if (dependency_ids[i] == "") continue
        candidate = longest_path_ending_at(dependency_ids[i])
        if (candidate > best_prefix) best_prefix = candidate
      }
      delete in_progress[node_id]
      memoized_path[node_id] = best_prefix + minutes_by_id[node_id]
      return memoized_path[node_id]
    }
    END {
      if (exit_code) exit exit_code
      critical_path = 0
      for (node_id in minutes_by_id) {
        candidate = longest_path_ending_at(node_id)
        if (candidate > critical_path) critical_path = candidate
      }
      printf "total_estimated_effort_minutes: %d\n", total_minutes
      printf "critical_path_minutes: %d\n", critical_path
    }
  '
  exit "$?"
fi

# ---------------------------------------------------------------------------
# Estimate mode
# ---------------------------------------------------------------------------
route_class=""
write_set_count=0
new_file_count=0
subsystem_count=1
acceptance_count=0
dependency_depth=0
browser_evidence=0
persistence_changes=0
async_behavior=0
performance_work=0
regression_gates=0
full_suite_verification=0
trivial_short_circuit=0

require_integer() {
  case "${2:-}" in
    ''|*[!0-9]*) usage_error "$1 needs a non-negative integer argument" ;;
  esac
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --route)
      case "${2:-}" in
        A|a) route_class="A" ;;
        B|b) route_class="B" ;;
        C|c) route_class="C" ;;
        *) usage_error "--route must be A, B, or C" ;;
      esac
      shift 2 ;;
    --write-set)    require_integer --write-set "${2:-}";  write_set_count="$2";  shift 2 ;;
    --new-files)    require_integer --new-files "${2:-}";  new_file_count="$2";   shift 2 ;;
    --subsystems)   require_integer --subsystems "${2:-}"; subsystem_count="$2";  shift 2 ;;
    --acceptance)   require_integer --acceptance "${2:-}"; acceptance_count="$2"; shift 2 ;;
    --deps-depth)   require_integer --deps-depth "${2:-}"; dependency_depth="$2"; shift 2 ;;
    --browser)        browser_evidence=1;        shift ;;
    --persistence)    persistence_changes=1;     shift ;;
    --async-behavior) async_behavior=1;          shift ;;
    --performance)    performance_work=1;        shift ;;
    --regression)     regression_gates=1;        shift ;;
    --full-suite)     full_suite_verification=1; shift ;;
    --trivial)        trivial_short_circuit=1;   shift ;;
    *) usage_error "unrecognized argument '$1'" ;;
  esac
done

if [ "$trivial_short_circuit" -eq 1 ]; then
  printf 'p50_active_minutes: %d\n' "$floor_minutes"
  printf 'confidence: high\n'
  printf 'basis:\n- trivial short-circuit\n'
  exit 0
fi

[ -n "$route_class" ] || usage_error "--route is required (or use --trivial)"

case "$route_class" in
  A) raw_minutes=5 ;;
  B) raw_minutes=10 ;;
  C) raw_minutes=20 ;;
esac
raw_minutes=$(( raw_minutes + write_set_count * 1 ))
raw_minutes=$(( raw_minutes + new_file_count * 2 ))
if [ "$subsystem_count" -gt 1 ]; then
  raw_minutes=$(( raw_minutes + (subsystem_count - 1) * 3 ))
fi
raw_minutes=$(( raw_minutes + acceptance_count * 1 ))
raw_minutes=$(( raw_minutes + dependency_depth * 2 ))
[ "$browser_evidence" -eq 0 ]        || raw_minutes=$(( raw_minutes + 8 ))
[ "$persistence_changes" -eq 0 ]     || raw_minutes=$(( raw_minutes + 6 ))
[ "$async_behavior" -eq 0 ]          || raw_minutes=$(( raw_minutes + 6 ))
[ "$performance_work" -eq 0 ]        || raw_minutes=$(( raw_minutes + 4 ))
[ "$regression_gates" -eq 0 ]        || raw_minutes=$(( raw_minutes + 4 ))
[ "$full_suite_verification" -eq 0 ] || raw_minutes=$(( raw_minutes + 4 ))

rounded_minutes="$(round_to_nearest_five "$raw_minutes")"
if [ "$rounded_minutes" -lt "$floor_minutes" ]; then
  rounded_minutes="$floor_minutes"
fi

# Deterministic confidence rubric — same inputs, same answer, always.
confidence_level="medium"
if [ "$route_class" = "A" ] && [ "$raw_minutes" -le 10 ]; then
  confidence_level="high"
fi
if [ "$route_class" = "C" ]; then
  if [ "$write_set_count" -ge 15 ] || [ "$subsystem_count" -ge 3 ] || [ "$raw_minutes" -ge 75 ]; then
    confidence_level="low"
  fi
fi

printf 'p50_active_minutes: %d\n' "$rounded_minutes"
printf 'confidence: %s\n' "$confidence_level"
printf 'basis:\n'
printf -- '- Route %s\n' "$route_class"
[ "$write_set_count" -eq 0 ]  || printf -- '- %d-file write set\n' "$write_set_count"
[ "$new_file_count" -eq 0 ]   || printf -- '- %d new files\n' "$new_file_count"
[ "$subsystem_count" -le 1 ]  || printf -- '- %d subsystems involved\n' "$subsystem_count"
[ "$acceptance_count" -eq 0 ] || printf -- '- %d acceptance criteria\n' "$acceptance_count"
[ "$dependency_depth" -eq 0 ] || printf -- '- dependency depth %d\n' "$dependency_depth"
[ "$browser_evidence" -eq 0 ]        || printf -- '- browser evidence\n'
[ "$persistence_changes" -eq 0 ]     || printf -- '- persistence changes\n'
[ "$async_behavior" -eq 0 ]          || printf -- '- async lifecycle behavior\n'
[ "$performance_work" -eq 0 ]        || printf -- '- performance instrumentation\n'
[ "$regression_gates" -eq 0 ]        || printf -- '- cross-route regression gates\n'
[ "$full_suite_verification" -eq 0 ] || printf -- '- full-suite verification\n'
