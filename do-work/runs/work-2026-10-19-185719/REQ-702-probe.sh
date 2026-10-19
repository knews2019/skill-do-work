#!/usr/bin/env bash
# GREEN checks for REQ-702 (test gate): the foreign-backend launcher mirrors backend output and
# writes heartbeats to its stderr sink while the backend runs, and the existing second-opinion
# tests stay green. Names the REQ's new test plus the package's SecondOpinion tests; a skip is
# not a pass. Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module="$root/skills/do-work/tools/do-work-cli"
new_test=TestSecondOpinionCaptureStreamsProgressWhileBackendRuns
output="$(go test -C "$module" -count=1 -v -run 'SecondOpinion' ./internal/toolboxcommands 2>&1)" || { printf '%s\n' "$output"; exit 1; }
grep -q -- "^--- PASS: $new_test " <<<"$output" || { printf '%s\n' "$output"; echo "$new_test did not pass (skipped or missing)"; exit 1; }
echo "REQ-702 probe: progress test and every SecondOpinion test passed"
