#!/usr/bin/env bash
# Focused tests for REQ-637: Panel B's stamp-gap verdict and the drawer gap-row label.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
go test -C "$root/skills/do-work-board/tools/queue-kanban" -count=1 -run 'ImplementationSpan|DayMedian|Durations' .
QUEUE_KANBAN_JAVASCRIPT_PROBES=on go test -C "$root/skills/do-work-board/tools/queue-kanban" -count=1 -run 'TestJavaScriptBehaviorDetailStatesTheLargestIdleGap' .
