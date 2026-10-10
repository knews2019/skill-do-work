#!/usr/bin/env bash
# GREEN checks for REQ-679 (integrator test gate): the two live-queue timeline browser probes build from a fixed fixture and pass in a real browser, and a skip is not a pass. The Previous and Next assertions live inside the NowAndFitAll probe. Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
export QUEUE_KANBAN_BROWSER_PROBES=on
export QUEUE_KANBAN_BROWSER="${QUEUE_KANBAN_BROWSER:-/Applications/Google Chrome.app/Contents/MacOS/Google Chrome}"
probe_tests=(TestBrowserBehaviorTimelineProseDescribesOnlyTheWindowOnScreen TestBrowserBehaviorTimelineNowAndFitAllLandSomewhereReadable)
output="$(go test -C "$root/skills/do-work-board/tools/queue-kanban" -count=1 -v -run "^($(IFS='|'; echo "${probe_tests[*]}"))\$" . 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${probe_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done
