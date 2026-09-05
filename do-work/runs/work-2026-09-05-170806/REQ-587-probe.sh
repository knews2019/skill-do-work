#!/usr/bin/env bash
# Focused probe recorded as REQ-587 run evidence: the exact argv this run executed.
# Runs in a detached worktree at the merge revision 8fad73b2, because the shared main tree
# carries sibling sessions' uncommitted edits in this very package. The browser engine is
# named explicitly: Chrome is installed but not on PATH under any name the probe looks for,
# and without the variable every browser probe reports skipped, which is not a pass.
set -euo pipefail
scratch_root="/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/f4a29cd7-255d-49fa-bcd0-8b59fe0da814/scratchpad"
cd "$scratch_root/verify-587/skills/do-work-board/tools/queue-kanban"
QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" \
QUEUE_KANBAN_BROWSER_PROBES=on \
  go test -count=1 -run 'TestBrowserBehaviorTimelineViewHasOneScrollSurface|TestBrowserBehaviorActivityViewHasOneScrollSurface|TestBrowserBehaviorTimelineBarsSurviveTheDetailDrawerOpening|TestBrowserBehaviorTimelineRowListIsOneTabStop|TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving|TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage|TestBrowserBehaviorTimelineListsRowsBeneathUserRequestHeaders' -v . 2>&1 | grep -E '^(=== RUN|--- (PASS|FAIL|SKIP)|ok|FAIL|PASS)' 
