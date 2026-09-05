#!/usr/bin/env bash
# Focused probe recorded as REQ-589 run evidence: the exact argv this run executed.
# Runs in the detached worktree at the merge revision 34422032 (the shared main tree carries a
# sibling session's uncommitted do-work/ edits). The Node-lane selection holds the three slim-band
# cases, REQ-578's hide-on-Activity case and the assembly structure test the template change touches;
# the generate tests pin the verify payload the strip renders.
set -euo pipefail
cd .git/work-run-20260905-1820/gate-34422032/skills/do-work-board/tools/queue-kanban && QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 go test -count=1 -run 'TestJavaScriptBehaviorVerifyFindings|TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip|TestBoardJavaScriptAssemblyStructure|TestGeneratedBoardDataCarriesVerifyFindings|TestGeneratedVerifyPayloadCarriesNoAbsolutePaths' .
