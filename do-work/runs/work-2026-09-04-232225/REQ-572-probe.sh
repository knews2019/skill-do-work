#!/usr/bin/env bash
# Focused probe recorded as REQ-572 run evidence: the exact argv this run executed.
set -euo pipefail
cd skills/do-work-board/tools/queue-kanban && QUEUE_KANBAN_BROWSER_PROBES=off QUEUE_KANBAN_JAVASCRIPT_PROBES=on go test -count=1 -run 'TestBuildActivityRows|TestLifecycleTimestampFieldsIsTheOneListBothReadersUse|TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests' .
