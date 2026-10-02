#!/usr/bin/env bash
set -euo pipefail
cd /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban
QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 -run 'TestJavaScriptBehavior(Board|UnknownBoard)' .
