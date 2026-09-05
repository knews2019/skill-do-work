#!/usr/bin/env bash
# Focused probe recorded as REQ-586 run evidence: the exact argv this run executed.
# Runs in the detached worktree at the merge revision 2ea0b150 (the shared main tree carries a sibling
# session's uncommitted edits to the same package). The two Node-lane tests are the REQ's RED/GREEN pair;
# REQ-585's browser probe must stay green with the chips beside the summary; the generate and assembly
# tests pin the template placeholder and fragment manifest the change touches.
set -euo pipefail
cd .git/work-run-20260905-1248/gate-2ea0b150/skills/do-work-board/tools/queue-kanban && QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 QUEUE_KANBAN_BROWSER_PROBES=on QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" go test -count=1 -run 'TestJavaScriptBehaviorActivityWindowChipsRenderInsideTheActivityView|TestJavaScriptBehaviorTopBarIdentityIsOneLineWithAFullStampTooltip|TestJavaScriptBehaviorActivity|TestBrowserBehaviorActivityViewHasOneScrollSurface|TestGenerate|TestBoardJavaScriptAssemblyStructure' .
