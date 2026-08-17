#!/usr/bin/env bash
# p50-estimator-determinism.sh — lock-in suite for the shipped P50 estimator
# (skills/do-work/tools/estimate-p50.sh). Pins the contracts REQ-208 promised:
# deterministic output, nearest-5 rounding, the 5-minute floor, representative
# route pins (small Route A, focused Route B, integrated Route C, browser-heavy
# QA), dependency-graph critical-path math (chain + diamond), and the
# backwards-compatibility guarantee that the estimator writes nothing (the
# estimate: block stays strictly additive, so legacy REQs remain valid).
#
# Exit 0: every probe passed. Exit 1: at least one FAIL line above.
set -uo pipefail

script_path="${BASH_SOURCE[0]}"
script_directory="${script_path%/*}"
if [ "$script_directory" = "$script_path" ]; then
  script_directory='.'
fi
repo_root="$(cd "$script_directory/../.." && pwd)"
estimator="$repo_root/skills/do-work/tools/estimate-p50.sh"

fail_count=0

report_failure() {
  printf 'FAIL: %s\n' "$1" >&2
  fail_count=$((fail_count + 1))
}

if [ ! -x "$estimator" ]; then
  report_failure "estimator missing or not executable at skills/do-work/tools/estimate-p50.sh"
  printf 'p50 estimator suite: %s probes failed.\n' "$fail_count" >&2
  exit 1
fi

# --- Determinism: identical flags twice must be byte-identical. -------------
first_run="$("$estimator" --route C --write-set 12 --browser --persistence --full-suite)"
second_run="$("$estimator" --route C --write-set 12 --browser --persistence --full-suite)"
if [ "$first_run" != "$second_run" ]; then
  report_failure "identical flags produced different output (determinism broken)"
fi

# --- Calibrated worked example: Route C + 12-file write set + browser +
# --- persistence + full-suite pins at 50 medium (archive-calibrated table;
# --- the spec's uncalibrated example printed 125). ------------------------
expect_line() {
  local probe_name="$1" output="$2" wanted_line="$3"
  if ! printf '%s\n' "$output" | grep -qxF -- "$wanted_line"; then
    report_failure "$probe_name: expected line '$wanted_line' missing from output: $(printf '%s' "$output" | tr '\n' '|')"
  fi
}

expect_line "spec example" "$first_run" "p50_active_minutes: 50"
expect_line "spec example" "$first_run" "confidence: medium"
expect_line "spec example" "$first_run" "- Route C"
expect_line "spec example" "$first_run" "- 12-file write set"
expect_line "spec example" "$first_run" "- browser evidence"

# --- Small Route A: the measured route median (= floor) + high confidence. --
route_a_output="$("$estimator" --route A)"
expect_line "small Route A" "$route_a_output" "p50_active_minutes: 5"
expect_line "small Route A" "$route_a_output" "confidence: high"

# --- Focused Route B: 10 + 3*1 + 2*1 = 15, medium. --------------------------
route_b_output="$("$estimator" --route B --write-set 3 --acceptance 2)"
expect_line "focused Route B" "$route_b_output" "p50_active_minutes: 15"
expect_line "focused Route B" "$route_b_output" "confidence: medium"

# --- Browser-heavy QA: 10 + 5*1 + 8 + 6 = 29 rounds to 30, medium. ----------
browser_heavy_output="$("$estimator" --route B --write-set 5 --browser --async-behavior)"
expect_line "browser-heavy QA" "$browser_heavy_output" "p50_active_minutes: 30"
expect_line "browser-heavy QA" "$browser_heavy_output" "- async lifecycle behavior"

# --- Rounding: 5 + 4*1 = 9 rounds to nearest five = 10. ---------------------
rounding_output="$("$estimator" --route A --acceptance 4)"
expect_line "nearest-5 rounding" "$rounding_output" "p50_active_minutes: 10"

# --- Trivial short-circuit: exactly the floor, high confidence. -------------
trivial_output="$("$estimator" --trivial)"
expect_line "trivial short-circuit" "$trivial_output" "p50_active_minutes: 5"
expect_line "trivial short-circuit" "$trivial_output" "confidence: high"
expect_line "trivial short-circuit" "$trivial_output" "- trivial short-circuit"

# --- Wide Route C drops confidence to low (write-set >= 15). ----------------
wide_route_c_output="$("$estimator" --route C --write-set 16)"
expect_line "wide Route C" "$wide_route_c_output" "p50_active_minutes: 35"
expect_line "wide Route C" "$wide_route_c_output" "confidence: low"

# --- Every emitted estimate is a multiple of five. --------------------------
for probe_output in "$first_run" "$route_a_output" "$route_b_output" \
  "$browser_heavy_output" "$rounding_output" "$trivial_output" "$wide_route_c_output"; do
  minutes_value="$(printf '%s\n' "$probe_output" | sed -n 's/^p50_active_minutes: //p')"
  if [ -z "$minutes_value" ] || [ $((minutes_value % 5)) -ne 0 ] || [ "$minutes_value" -lt 5 ]; then
    report_failure "estimate '$minutes_value' is not a >=5 multiple of five"
  fi
done

# --- Critical path, chain: sums along the only path. ------------------------
chain_output="$("$estimator" critical-path REQ-1:10 REQ-2:20:REQ-1 REQ-3:30:REQ-2)"
expect_line "chain graph" "$chain_output" "total_estimated_effort_minutes: 60"
expect_line "chain graph" "$chain_output" "critical_path_minutes: 60"

# --- Critical path, diamond: longest branch, never the parallel sum. --------
diamond_output="$("$estimator" critical-path REQ-1:10 REQ-2:50:REQ-1 REQ-3:20:REQ-1 REQ-4:10:REQ-2,REQ-3)"
expect_line "diamond graph" "$diamond_output" "total_estimated_effort_minutes: 90"
expect_line "diamond graph" "$diamond_output" "critical_path_minutes: 70"

# --- Unknown dependency ids contribute zero (already-archived REQs). --------
unknown_dep_output="$("$estimator" critical-path REQ-7:20:REQ-999)"
expect_line "unknown dependency" "$unknown_dep_output" "total_estimated_effort_minutes: 20"
expect_line "unknown dependency" "$unknown_dep_output" "critical_path_minutes: 20"

# --- A dependency cycle must fail loudly, never hang or fabricate. ----------
if "$estimator" critical-path REQ-1:10:REQ-2 REQ-2:10:REQ-1 >/dev/null 2>&1; then
  report_failure "dependency cycle was accepted instead of rejected"
fi

# --- Unrecognized flags are rejected, mirroring the work action's parser. ---
if "$estimator" --route B --no-such-flag >/dev/null 2>&1; then
  report_failure "unrecognized flag was accepted instead of rejected"
fi

# --- Backwards compatibility: the estimator only prints — it must never -----
# --- write, so a legacy REQ without an estimate: block stays untouched. -----
backcompat_workdir="$(mktemp -d)"
printf -- '---\nid: REQ-001\ntitle: Legacy request\nstatus: pending\n---\n\n# Legacy\n' \
  > "$backcompat_workdir/REQ-001-legacy.md"
before_listing="$(cd "$backcompat_workdir" && find . -type f | sort && cat REQ-001-legacy.md)"
(cd "$backcompat_workdir" && "$estimator" --route A >/dev/null)
after_listing="$(cd "$backcompat_workdir" && find . -type f | sort && cat REQ-001-legacy.md)"
if [ "$before_listing" != "$after_listing" ]; then
  report_failure "estimator invocation mutated the working directory (must be print-only)"
fi
rm -rf "$backcompat_workdir"

if [ "$fail_count" -gt 0 ]; then
  printf 'p50 estimator suite: %s probes failed.\n' "$fail_count" >&2
  exit 1
fi
printf 'p50 estimator suite: all probes passed.\n'
