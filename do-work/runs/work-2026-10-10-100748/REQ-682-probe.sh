#!/usr/bin/env bash
# GREEN checks for REQ-682 (integrator test gate): renderVisibleRows opens with a hidden-panel guard, the Timeline Node-lane tests still pass on the stub DOM, the scroll-surface browser probe passes (plus any kept hidden-view probe, named TestBrowserBehaviorTimelineHiddenView*) with the browser lane on, and a skip is not a pass. Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
board_dir="$root/skills/do-work-board/tools/queue-kanban"
render_head="$(sed -n '/function renderVisibleRows() {/,/rowsSvg.textContent = ""/p' "$board_dir/web/board-timeline.js")"
grep -q 'getElementById("view-timeline")' <<<"$render_head" || { echo "renderVisibleRows has no view-timeline hidden guard before it rebuilds rows"; exit 1; }
go test -C "$board_dir" -count=1 -run 'TestJavaScriptBehaviorTimeline' .
export QUEUE_KANBAN_BROWSER_PROBES=on
export QUEUE_KANBAN_BROWSER="${QUEUE_KANBAN_BROWSER:-/Applications/Google Chrome.app/Contents/MacOS/Google Chrome}"
probe_tests=(TestBrowserBehaviorTimelineViewHasOneScrollSurface)
scroll_probe_source="$(<"$board_dir/timeline_scroll_browser_probe_test.go")"
while IFS= read -r hidden_probe; do
  probe_tests+=("${hidden_probe#func }")
done < <(grep -o 'func TestBrowserBehaviorTimelineHiddenView[A-Za-z]*' <<<"$scroll_probe_source" || true)
output="$(go test -C "$board_dir" -count=1 -v -run "^($(IFS='|'; echo "${probe_tests[*]}"))\$" . 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${probe_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done
