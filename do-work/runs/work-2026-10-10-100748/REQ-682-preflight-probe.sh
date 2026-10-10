#!/usr/bin/env bash
# Pre-flight baseline for REQ-682 (run before the build): the Timeline Node-lane tests and the existing scroll-surface browser probe pass today with the browser lane on and are not skipped.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
board_dir="$root/skills/do-work-board/tools/queue-kanban"
go test -C "$board_dir" -count=1 -run 'TestJavaScriptBehaviorTimeline' .
export QUEUE_KANBAN_BROWSER_PROBES=on
export QUEUE_KANBAN_BROWSER="${QUEUE_KANBAN_BROWSER:-/Applications/Google Chrome.app/Contents/MacOS/Google Chrome}"
output="$(go test -C "$board_dir" -count=1 -v -run '^TestBrowserBehaviorTimelineViewHasOneScrollSurface$' . 2>&1)" || { printf '%s\n' "$output"; exit 1; }
grep -q -- '^--- PASS: TestBrowserBehaviorTimelineViewHasOneScrollSurface ' <<<"$output" || { printf '%s\n' "$output"; echo "scroll-surface probe did not pass (skipped or missing)"; exit 1; }
