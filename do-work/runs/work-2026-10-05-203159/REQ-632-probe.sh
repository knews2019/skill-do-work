#!/usr/bin/env bash
# Focused tests for REQ-632: implementation span, durations, generated payload, activity correlation, card and drawer behaviour.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)/skills/do-work-board/tools/queue-kanban"
QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_BROWSER_PROBES=off \
  go test -count=1 -run 'ImplementationSpan|PhaseBreakdown|Activity|DoneCard|Drawer|GeneratedRequest' .
