# REQ-679 hand-back

- Branch: worktree-agent-REQ-679-timeline-probes-fixed-data
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-679-timeline-probes-fixed-data
- Base: b629e5cd. Commit: bab949ea `[REQ-679] timeline probes use fixed fixture data and assert Previous and Next` (one commit).

## File manifest
- skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go (modified): the shared range-end fixture tree gains two rows; the prose probe and the Now/Fit all probe build from it; the Now/Fit all probe presses Previous then Next and the forward refusal is asserted outright.

## P-A-U
**[PLAN]** Reuse `generateLiveSiteInDirAtRangeEnd` and add two rows to its one tree instead of building a second tree. Trailing-windows probe only asserts "range ends at now" and "range starts more than 30 days before now", so older rows cannot break it (re-run, passes). Rows: `REQ-164` (completed, absolute dates 2026-07-28 to 07-29, so the fixed typed week draws it and the `REQ-164` filter fits to a few days) and `REQ-0003` (pending, created 2026-07-28, so it stays open through the typed week, is the long span that makes the filtered fit under half the unfiltered one, and fills the seven days before now so Previous lands on drawn data). Nine other probes read: see Proof record.
**[APPLY]** As planned: fixture rows plus comment, two `generateLiveSiteInDir` calls switched, two script presses added, two struct fields and two entries in the `states` loop, clause (3) rewritten without the two-branch conditional, new clause (3b). Stale comments ("Late July is the busiest stretch of this repo's own archive", "depends on the live queue") updated.
**[UNIFY]** `git diff b629e5cd --stat`: 1 file changed, 86 insertions(+), 45 deletions(-). Checks: probe.sh exit 0 (3.5 s); `gofmt -l` empty; `go vet` exit 0; `git diff --check` exit 0; `git status` clean after commit; no `web/` path in the diff (mutation reverted with `git checkout -- web/`). File checked: timeline_browser_probe_test.go (read the full diff for debug artifacts, consumer commit IDs and REQ numbers in new comments: none; the new comments name REQ-164 and REQ-0003 only as fixture row ids).

## Proof record
- RED, prose probe: scratch detached worktree at b629e5cd, deleted the 54 files in do-work/archive|queue|working whose created_at or completed_at falls 2026-07-27..2026-08-02. BASE probe: FAIL at `timeline_browser_probe_test.go:2482` "Nothing was drawn between 2026-07-27 00:00 UTC and 2026-08-03 00:00 UTC ... 616 REQs are outside it". My file copied into the same scratch tree: both TestBrowserBehaviorTimelineProseDescribesOnlyTheWindowOnScreen and ...NowAndFitAllLandSomewhereReadable PASS. Scratch worktree removed and pruned.
- RED, step assertions: throwaway edit in web/board-timeline.js `applyTrailingWindowStep` moving by `stepCount * 86400000` (one day) instead of the stepped screenful. NowAndFitAll FAILS at :2279 "stepping back from 2026-10-03 10:29 UTC → 2026-10-10 10:29 UTC landed on 2026-10-02 10:29 UTC → 2026-10-09 10:29 UTC; both endpoints should move back exactly seven days (604800000 ms)" and at :2291 (Next disabled one day back). Reverted; web/ clean. A second mutation (step count divided by 7 in `steppedWindowFor`) also failed :2279 and :2295.
- GREEN: `go test -count=1 -v -run Timeline ./...` with browser variables: exit 0, 49 s, 38 PASS, 0 FAIL, 22 `TestBrowserBehavior*` PASS. 23 SKIP lines are all `TestJavaScriptBehaviorTimeline*` ("JavaScript behavior probes are heavy-only", another lane, not the browser lane); no TestBrowserBehavior test skipped. Each of the three edited-fixture probes also passed on their own first run (1.6 to 1.9 s). No load-related reruns were needed.
- Nine other `generateLiveSiteInDir` probes (verdict: none moved; none has a hardcoded date, only `2099-12-31` which is out-of-range on purpose):
  - :629 PressBecomesAPanOnlyAfterMoving: needs one row under the pointer; no counts or dates. Leave.
  - :901 DragRendersOncePerFrame: counts renders and width reads. Leave.
  - :1328 BarsSurviveTheDetailDrawerOpening: needs some segments and a row. Leave.
  - :1599 RangeFieldsShowTheWindowTheChartIsDrawnAt: relative to the board's own range and Last 30 days. Leave.
  - :2545 PointerAndKeyboardPathsStayAlive: needs one row and a pannable window. Leave.
  - :3176 PointerCaptureWaitsForThePanEngage: needs a row to press. Leave.
  - :3382 RowListIsOneTabStop: needs at least 2 rows. A weak count dependence that only fails on a near-empty queue. Leave.
  - :3564 ListsRowsBeneathUserRequestHeaders: needs 2+ UR groups, a table taller than the viewport (virtualization) and Fit all vs Last day listing different REQs. This is the only count-dependent one; its own comment says it is deliberately rendered by the shipped board, not a fixture. It cannot go stale from one archive prune; moving it needs a large fixture (bloat). Leave (see D-02).
  - :3832 GroupHeadersReadInBothThemes: colour and geometry of a group header; needs at least one header. Leave.

## Decisions
- D-01 DECIDE & STATE: two added rows, not three or a second tree. Both named probes and the trailing-windows probe share the one tree. Reason: the trailing probe only needs range end at now and start more than 30 days back, and it passes unchanged. The REQ's Exploration also listed "a long-span row elsewhere"; REQ-0003 serves that and the open-in-week need and the drawn-previous-week need in one row.
- D-02 DECIDE & STATE: zero of the nine moved (REQ expected two or three at most). Only :3564 depends on counts and the queue (not dates) and it is a documented design choice. If the maintainer wants it moved, that is a separate fixture of dozens of rows.
- D-03 DECIDE & STATE: kept the helper name `generateLiveSiteInDirAtRangeEnd` (YAGNI, no rename); its doc comment now says the tree also serves the prose and Now/Fit all probes.
- D-04 DECIDE & STATE: no `timelineProbeInstant` helper (the patch adds one). Failure messages quote the readouts, which already show instants.
- D-05 DECIDE & STATE: dropped the patch's extra "nothing drawn after Next" assertion; the REQ names Previous drawn, Next returns exact start instants.

## Discovered Tasks
- The fixture rows use absolute 2026 dates for the typed week; the probe will stay valid as long as the real clock is after 2026-08-02. No action. → report only
- `REQ-679-probe.sh` finishes in 3.5 s, so it cannot be running the browser lane; the GREEN above is from the direct go test run. → report only

## Lessons read
`_dev/primes/lessons-releases.md` (whole), prime-releases.md, prime-kanban-board.md (Conventions, Traps); crew members general, coding-guardrails, shared-principles, communication-style, testing. Go satellites were not read (no Go production code touched).

## Anti-bloat check
`git diff b629e5cd --stat`: 1 file changed, 86 insertions(+), 45 deletions(-). Added that the REQ did not name: none (no new function, constant, option, file or test). New local names inside the test only: `sevenDaysMs`, `trailing`, `afterPrev`, `afterNext`; two struct fields and two script captures for Previous/Next, which the REQ names. Two fixture rows added to the existing tree (REQ names "add those rows to a fixture tree in the same style").

## Proposed CHANGELOG entry
Timeline browser probes no longer depend on which REQs the checkout holds, and step navigation is covered

The prose and Now/Fit all browser probes read the live queue, so a checkout with nothing dated 2026-07-27 to 2026-08-02 failed the strict lane for no product reason, and no probe pressed Previous or Next.
- Both probes build their page from the fixed range-end fixture tree (two rows added: a short completed REQ-164 and a long open REQ-0003).
- Now/Fit all presses Previous then Next on the trailing 7 days: exact 7-day move on both endpoints in epoch milliseconds, same span, data drawn, Next returns to the start.
- The forward arrow refusal on the trailing 7 days is asserted outright instead of in a two-branch conditional.

## Proposed lesson bullet
- [family: probe-reads-live-queue] [REQ-679: a browser probe that types fixed dates against the live queue passes only while the archive happens to have rows there; build it from the fixture tree, and prove step buttons by pressing them, since permanently disabled arrows satisfy a refusal-only assertion](../../do-work/archive/UR-152/REQ-679-timeline-probes-fixed-data.md) (satellite: `_dev/primes/lessons-kanban-board.md`)

## Integration seams and wall times
No overlap with REQ-682 (it edits web/board-timeline.js and timeline_scroll_browser_probe_test.go). Wall times: GREEN Timeline run 49 s; the three edited probes together 5.3 s; RED scratch prose run 4.3 s; mutation run about 2 s each.
