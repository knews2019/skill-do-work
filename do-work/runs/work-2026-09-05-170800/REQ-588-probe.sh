#!/usr/bin/env bash
# Focused probe recorded as REQ-588 run evidence: the exact argv this run executed.
# Runs in the detached worktree at the merge revision ab251f24 (the shared main tree carries a
# sibling session's uncommitted do-work/ edits). The Node-lane selection holds the REQ's new
# RED/GREEN case beside REQ-579's one-row-list case and REQ-578's hide-on-Activity case; the Go
# selection holds the producer's new detail rule, the verify suite whose assertions moved to
# Subject, the timeline claim assertion D4 forced, and the generate payload tests.
set -euo pipefail
cd .git/work-run-20260905-1708/gate-ab251f24/skills/do-work-board/tools/queue-kanban && QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 go test -count=1 -run 'TestJavaScriptBehaviorVerifyFindingRemedyIsItsOwnLineAfterTheDetail|TestJavaScriptBehaviorVerifyFindingsRenderAsOneRowList|TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip|TestVerify|TestTimelineProjectionReservesTimeForUntimedClaimedWork|TestGeneratedBoardDataCarriesVerifyFindings|TestGeneratedVerifyPayloadCarriesNoAbsolutePaths' .
