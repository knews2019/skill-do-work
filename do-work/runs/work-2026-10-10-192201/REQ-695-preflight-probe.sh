#!/usr/bin/env bash
# Pre-flight baseline for REQ-695 (run before the build): the run-status tests this REQ extends pass today, including the two pins the build must keep green (C3 with a landed hand-back, and no run folder means absent run fields).
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module_dir="$root/skills/do-work/tools/do-work-cli"
output="$(go test -C "$module_dir" -count=1 -v -run '^TestRunStatus' ./internal/runstatus/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for test_name in TestRunStatusLandedHandbackIsClassC3 TestRunStatusReportsQueueRowsWithoutARunDirectory TestRunStatusStaleClaimNamesRecoveryOnlyInTheStopReason; do
  grep -qF -- "--- PASS: $test_name (" <<<"$output" || { printf 'no PASS line for %s:\n%s\n' "$test_name" "$output"; exit 1; }
done
printf 'REQ-695 pre-flight: run-status tests pass\n'
