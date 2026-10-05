#!/usr/bin/env bash
# Focused tests for REQ-633: Panel B exclusion rule, verdict readers, and the calibration-log row builders.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work-board/tools/queue-kanban"
QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_BROWSER_PROBES=off \
  go test -count=1 -run 'Duration|Timeline|ImplementationSpan|UserRequestProgress|DoneCard|Calibration' .
go test -C "$root/skills/do-work/tools/do-work-cli" -count=1 ./internal/requeststate/
