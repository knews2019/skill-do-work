# Hand-back — REQ-626
**Branch:** worktree-agent-REQ-626-board-page-urls  **Base:** 7917efbe  **Head:** 3100dea9

## File manifest
- skills/do-work-board/tools/queue-kanban/web/board-controls.js (modified) — fragment parse/format, guarded boot read, click-only write, hashchange listener
- skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go (modified) — 4 TestJavaScriptBehavior probes driving the real wireControls
- skills/do-work-board/tools/queue-kanban/page_fragment_browser_probe_test.go (new) — real-browser probe: probe.html#timeline opens Timeline, Activity click writes #activity
- skills/do-work-board/tools/queue-kanban/browser_probe_test.go (modified) — helper `probePageAddressWithoutFragment`, `startTrustedInputBrowserSessionAtFragment`, 3 URL checks use the helper
- skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go (modified) — 2 URL checks use the helper
- skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go (modified, granted widening) — 9 URL checks use the helper
- skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go (modified, granted widening) — 1 URL check uses the helper
- skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go (modified, granted widening) — 1 URL check uses the helper
- skills/do-work-board/tools/queue-kanban/user_request_progress_browser_probe_test.go (modified, granted widening) — 1 URL check uses the helper
- skills/do-work-board/tools/queue-kanban/prime-do-kanban.md (modified) — one Traps bullet

## Integration seams
None. Release paths (VERSION, changelogs, version action) untouched as instructed; the orchestrator owns the release.

## P-A-U
- [PLAN] Five small functions in board-controls.js: `boardStateFromFragment` (table lookup, null for unknown), `boardFragmentFromState`, `currentAddressHash` (guarded), `applyBoardStateSelection` (sets view/lens/fold, resets `renderedOnce.userRequestLens` when a lens is named, re-syncs both button groups, renders nothing), `writeBoardFragment` (guarded replaceState with location.replace fallback, skips when unchanged). `wireControls` calls write after page/lens clicks, reads the hash once at boot, and adds a hashchange listener that applies known fragments and calls `applyView()`.
- [APPLY] Only the declared Scope files plus the granted one-line URL-check substitution in timeline, clipboard, user_request_clipboard and user_request_progress probe tests. board.js and template.html untouched. `applyView` unchanged.
- [UNIFY] `git diff --stat HEAD~1`: 10 files changed, 473 insertions(+), 19 deletions(-). `gofmt -l .` empty, `go vet .` exit 0. Checked board-controls.js (no stub left, no braces in string literals in new functions), every substituted URL check (diff reviewed line by line), the prime bullet.

## Red-green evidence
RED captured against a compiling stub (functions present, doing nothing; stub never committed).
- TestJavaScriptBehaviorBoardFragmentOpensItsPageAtBoot: RED `opening #timeline: got view="board" lens="flat" folded=false pressed view="board" lens="flat"; want {View:timeline ...}` (same for the other 6 non-default fragments) → GREEN at 3100dea9
- TestJavaScriptBehaviorBoardClickWritesItsFragment: RED `timelineClick wrote [], want ["history.replaceState #timeline"]` ... `refusedReplaceState wrote [], want ["location.replace #calendar"]` → GREEN at 3100dea9
- TestJavaScriptBehaviorBoardHashChangeSwitchesKnownPagesOnly: RED `hashchange to #durations did not switch and render the page: {View:board ...AppliedViews:[]}` → GREEN at 3100dea9
- TestJavaScriptBehaviorUnknownBoardFragmentKeepsTheBoard: passes against the stub by design (it pins the fallback, which the stub already has); guards against a reader that accepts `#board-main`, `#constructor` or `#BOARD`.
- TestBrowserBehaviorPageFragmentOpensAndFollowsThePage: RED `probe.html#timeline did not open the Timeline page: {Hash:#timeline TimelineShown:false BoardShown:true TimelinePressed:false BoardPressed:true TimelineSummary:}` → GREEN at 3100dea9 (2.43s, Google Chrome at /Applications)

## Tests run
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 -run TestJavaScriptBehavior .` → exit 0, 77 tests passed, 7.0s wall (new tests 0.05s each)
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 ./...` → exit 0, `ok ... 57.215s`, 57.5s wall
- `QUEUE_KANBAN_BROWSER_PROBES=on QUEUE_KANBAN_BROWSER=".../Google Chrome" QUEUE_KANBAN_JAVASCRIPT_PROBES=off go test -count=1 -run 'Browser|Probe|Citation|Clipboard|Contrast|Priority|Progress|Scroll' .` → exit 0, 70 top-level passes, 0 failures, 183.8s wall. Skips are only JavaScript-probe-gated tests (lane has JS probes off) and the legacy strict lane. Every existing timeline/durations/activity probe that clicks a view passed with the fragment now in the URL.
- `gofmt -l .` → empty; `go vet .` → exit 0
- Per-file: new page_fragment_browser_probe_test.go 2.4s; new tests in javascript_behavior_d_test.go 0.2s combined. Pre-existing slow test noted, not mine: TestLegacyStrictEntryPointsExecuteProbes 81.1s (strict_behavior_test.go).

## Decisions
- D-01 (DECIDE & STATE): Back behaviour. replaceState adds no history entry, so Back leaves the board as before. The browser probe asserts `history.length` does not grow on a click.
- D-02 (DECIDE & STATE): The boot read uses a new `applyBoardStateSelection` instead of `applyLensSelection`, because `applyLensSelection` calls `applyLens()` (a UR-lens render) before board.js has run `renderColumns()`. The new helper resets `renderedOnce.userRequestLens` the same way and renders nothing; lens clicks still go through `applyLensSelection`.
- D-03 (DECIDE & STATE): A non-Board fragment carries `lens: null`, so `#activity` keeps whatever lens was selected; `#board` explicitly selects the flat lens.
- D-04 (DECIDE & STATE): Added `startTrustedInputBrowserSessionAtFragment` in browser_probe_test.go (the old function delegates with ""), because the new probe must really open `probe.html#timeline`. The attach-target URL match also goes through the helper so a fragment URL is found. Small addition beyond the "one shared helper"; inside a Scope file.
- D-05 (DECIDE & STATE): Fragment lookup uses `hasOwnProperty` and is case-sensitive, so `#constructor` and `#BOARD` read as unknown.

## Discovered Tasks
None

## Lessons read
- `_dev/primes/lessons-kanban-board.md`: lines 41, 43, 71, 76 plus the surrounding index (lines 35-80 skimmed).
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`: not read beyond the exploration's note that it has nothing about views.

## Blockers
None
