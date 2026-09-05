#!/usr/bin/env bash
# Focused probe recorded as REQ-585 run evidence: the exact argv this run executed.
# Runs in the detached worktree at the merge revision c08ac2b4, because the shared main tree carries a
# sibling session's uncommitted edits to the same package; the browser probe is the REQ's GREEN condition,
# and the generate and assembly tests pin the stylesheet shape the change edits.
set -euo pipefail
cd .git/work-run-20260905-1248/gate-c08ac2b4/skills/do-work-board/tools/queue-kanban && QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_BROWSER_PROBES=on QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" go test -count=1 -run 'TestBrowserBehaviorActivityViewHasOneScrollSurface|TestGenerateInlines|TestBoardJavaScriptAssemblyStructure|TestJavaScriptBehaviorActivity' .
