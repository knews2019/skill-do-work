#!/usr/bin/env bash
# GREEN checks for REQ-680 (integrator test gate): the absent-target case exists inside the total-failure test and passes (the case name is part of the brief). Fails before the build.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
case_name='TestTotalFailurePreservesTheTargetAndLeavesNoScratch/absent_target'
output="$(go test -C "$root/skills/do-work/tools/do-work-cli" -count=1 -v -run "^${case_name%%/*}\$/^${case_name#*/}\$" ./internal/archivefetch/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
grep -q -- "--- PASS: $case_name " <<<"$output" || { printf '%s\n' "$output"; echo "the absent-target case did not run and pass"; exit 1; }
