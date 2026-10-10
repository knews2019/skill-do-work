# REQ-682 hand-back (hidden Timeline stops redrawing on scroll)

- Branch: worktree-agent-REQ-682-hidden-timeline-scroll-guard
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-682-hidden-timeline-scroll-guard
- Base: b629e5cd. Commits: 1e4c10a7 (guard, re-entry nudge, two probes) and b0588380 (stub `dispatchEvent` line in javascript_behavior_c_test.go, its own commit, approved as D-02). Both subjects start `[REQ-682]`.

## File manifest (all under skills/do-work-board/tools/queue-kanban/)
- web/board-timeline.js (modified): null-safe early return at the top of `renderVisibleRows` when `#view-timeline` is hidden.
- web/board-controls.js (modified): after the Timeline arrival scroll reset, `boardMain.dispatchEvent(new Event("scroll"))` so the still-live listener redraws (re-entry fix, see Red-green).
- timeline_scroll_browser_probe_test.go (modified): two new browser probes; the fixture/site setup of `TestBrowserBehaviorTimelineViewHasOneScrollSurface` moved verbatim into `buildTimelineScrollProbeSite` so the new probes share it (that is most of the diff-stat churn).
- javascript_behavior_c_test.go (modified, commit b0588380, D-02 approved): one line `dispatchEvent: function () {},` in the stub node. Superseded note: formerly an uncommitted patch: one stub line in javascript_behavior_c_test.go, patch at /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/2033d20b-9fe7-4b47-aa5d-67ee6d4581f2/scratchpad/REQ-682-stub-dispatchevent.patch (`git apply` checked; with it the whole JS Node lane passes).

## P-A-U
[PLAN] One guard in renderVisibleRows (covers scroll and capturing toggle listeners). Two browser probes first (RED), then guard, then check re-entry; fix only what re-entry needs.
[APPLY] As planned. Re-entry failed with the guard, so one dispatchEvent line in board-controls.js. No listener release, no state, no observer.
[UNIFY] `git diff b629e5cd --stat`: timeline_scroll_browser_probe_test.go 292 (+/-), board-controls.js 7, board-timeline.js 8; 3 files, 241 insertions, 66 deletions. Checks: probe script REQ-682-probe.sh exit 0 (7 s); gofmt -l empty (exit 0); go vet ./... exit 0; git diff --check exit 0; full package browser-lane-off `go test -count=1 .` ok 88 s (but see Node-lane note); `-v -run Timeline ./...` with browser vars: PASS 54 s, all 24 TestBrowserBehaviorTimeline* pass, no browser SKIP. Node lane with QUEUE_KANBAN_JAVASCRIPT_PROBES=on: `-run TestJavaScriptBehaviorTimeline` 23 PASS (3 s); full package with it on: one failure, TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip (stub lacks dispatchEvent, fixed by the uncommitted patch above; with the patch `-run TestJavaScriptBehavior` passes, 7 s). TestLegacyStrictEntryPointsRejectZeroProbes failed once in a loaded full run and passed alone (21.8 s). Files checked: the three above plus javascript_behavior_c_test.go.

## Red-green record (Chrome via QUEUE_KANBAN_BROWSER, 1600x900; each result carries location.href in the log)
- `TestBrowserBehaviorTimelineHiddenViewIgnoresBoardScroll` (Timeline visited, switch to Testing, MutationObserver on #timeline-scroll subtree, 10 scroll steps): base 2330 childList records, 2310 nodes added, 527 removed (FAIL). The REQ's 20/700 came from a shorter Testing view; this probe site's findings strip makes Testing 32140 px scrollable. With guard: 0/0/0 (PASS).
- `TestBrowserBehaviorTimelineHiddenViewRedrawsOnReturnAfterScrollClamp` (Timeline scrolled 1500 px into the rows, switch to Activity which clamps scrollTop to 0, return): base PASS (the hidden redraw happened to land on the top rows), guard only FAIL (rows after return are exactly the scrolled-down rows, no scroll event at scrollTop 0), guard + dispatchEvent PASS. So the guard alone is a regression and the one-line nudge is needed. The oracle is "top row back and no row of the scrolled-down window", not row-for-row equality: the first draw had 39 rows and a redraw 44.

## Decisions
- D-01 DECIDE & STATE: re-entry fix is `dispatchEvent(new Event("scroll"))` in board-controls.js, not a direct `renderVisibleRows()` call, because that function is local to `renderTimelineView` and exposing it would add state. Only the Timeline listens to scroll on #board-main (grep), so nothing else reacts.
- D-02 ESCALATE, approved by team-lead and committed as b0588380 (no `if (boardMain.dispatchEvent)` guard in production): add `dispatchEvent: function () {},` to the stub in javascript_behavior_c_test.go (outside my boundary). Value: Node lane green (one test breaks without it). Risk: tiny merge seam if another builder edits that file; reversible. Without it the committed branch breaks that Node test when QUEUE_KANBAN_JAVASCRIPT_PROBES=on.
- D-03 DECIDE & STATE: moved the existing test's setup into `buildTimelineScrollProbeSite` rather than copy 60 lines twice.
- D-04 DECIDE & STATE: no probe-name change; both new tests start with `TestBrowserBehaviorTimelineHiddenView`, so REQ-682-probe.sh runs both.

## Discovered Tasks
- impact-user-visible: the brief's Node-lane commands and REQ-682-probe.sh line `go test -run TestJavaScriptBehaviorTimeline` skip silently unless QUEUE_KANBAN_JAVASCRIPT_PROBES=on, so the probe script's Node check is vacuous and would not have caught the stub break above. → report only
- impact-user-visible: at base a Timeline scroll while hidden redraws using hidden-panel geometry (rect of a display:none panel), so the rows drawn there are only right by luck; this guard removes that path. → report only

## Lessons read
`_dev/primes/lessons-releases.md` (whole); primes prime-kanban-board.md and prime-releases.md; crew-members general, coding-guardrails, shared-principles, communication-style, frontend, testing. Board satellites not read (dropped for budget; touched code is not named in their traps).

## Anti-bloat check
Added that the REQ did not name: `buildTimelineScrollProbeSite` (test helper, replaces duplicated setup), `timelineHiddenScrollMeasurement` and `timelineReentryMeasurement` (test result structs), and the board-controls.js dispatch (REQ allowed it for re-entry). No new production function, constant, option, flag or file.

## Proposed CHANGELOG entry
**Hidden Timeline no longer redraws its rows while you scroll other views.** After opening the Timeline once, every scroll in Board, Calendar or Testing rebuilt the Timeline's invisible rows (about one frame of work per scroll in Testing).
- `renderVisibleRows` returns at once when the Timeline panel is hidden.
- Arriving at the Timeline asks for one redraw, so a Timeline left scrolled down and returned to after a short view clamped the scroll position shows the top rows.
- Two browser probes pin both behaviors.

## Proposed lesson bullet
- [family: hidden-panel-listener-guard] [REQ-682: a guard that drops events while a panel is hidden also drops the event a later re-entry relied on; prove the return path with a probe, and give a stubbed node every method a production change calls](../../do-work/archive/UR-152/REQ-682-hidden-timeline-scroll-guard.md) -> lessons-kanban-board.md

## Integration seams and wall times
No shared file with REQ-679 (it owns timeline_browser_probe_test.go). Only seam: javascript_behavior_c_test.go if D-02 is approved. Wall times: Timeline browser lane 54 s, probe script 7 s, Node Timeline lane 3 s, full package 88 s.

## Final verification after b0588380 (head b0588380)
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=on go test -count=1 -v -run TestJavaScriptBehavior .`: 80 PASS, 0 SKIP, 0 FAIL, ok 7.7 s.
- Browser vars on, `-v -run TestBrowserBehaviorTimeline .`: 18 top-level PASS, 0 SKIP, 0 FAIL, ok 71 s (includes both new HiddenView probes).
- Integration seam: javascript_behavior_c_test.go (one added stub line) is now touched by this branch; no other builder in this batch is known to edit it. REQ-682-probe.sh needs QUEUE_KANBAN_JAVASCRIPT_PROBES=on (integrator's fix, not mine).
