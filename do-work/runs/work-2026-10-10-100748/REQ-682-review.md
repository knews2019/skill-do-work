## Review: REQ-682

**Approve** — the hidden Timeline no longer redraws on another view's scroll, re-entry still draws the top window, and the diff adds nothing the fix does not use.
Route B | range a7b78f62..32552b5b

### What's built
- `renderVisibleRows` (`web/board-timeline.js:2139-2146`) returns at once when `#view-timeline` exists and is hidden. One early return, null-safe, no listener release, no new state, no observer.
- Arrival at the Timeline (`web/board-controls.js:44-49`) sends one synthetic `scroll` event to `#board-main` after the scroll reset, so a Timeline whose position was clamped to 0 while hidden redraws for scrollTop 0.
- Two browser probes in `timeline_scroll_browser_probe_test.go` pin both behaviors.

### Decisions / risks for you
- None that need a decision. The one undeclared file (`javascript_behavior_c_test.go`, one stub line) was escalated and approved as D-02 (builder decision 2: give the Node stub node a `dispatchEvent` method instead of guarding production code).

### Findings

**Important:**
- None.

**Minor:**
- F1. Anti-bloat count. Beyond the REQ's named items the diff adds exactly three things, all used: `buildTimelineScrollProbeSite` (the existing test's 60-line setup moved verbatim, now called by three tests; the alternative was two 60-line copies), and the two result structs `timelineHiddenScrollMeasurement` and `timelineReentryMeasurement` (the file's own pattern for typed `decodeResult` decoding, as `timelineAnchorMeasurement` already does). No new production function, constant, flag, option, config or file. No decorative test: one probe per named failure. Judged needed, not bloat. — impact-negligible → report only
- F2. Two comments in the moved helper are now slightly wrong. `timeline_scroll_browser_probe_test.go:309-310` says "Both Timeline scroll probes start from it", but three tests call it (lines 167, 767, 841). Line 335 says "the measurement below reports whichever element actually came first", but that measurement is in `TestBrowserBehaviorTimelineViewHasOneScrollSurface`, which is now above the helper. — impact-negligible → report only
- F3. The synthetic scroll event is synchronous, so on an arrival where a shared-filter change reset `renderedOnce.timeline` (`web/board-filters.js:176-187`), it runs the previous render's `renderVisibleRows` closure once before `renderTimelineView` rebuilds everything (`board-controls.js:49` runs before `:84`). On an arrival from a non-zero scrollTop it adds a second render next to the native async scroll event. Both are wasted work of about one row render (under 1 ms per the file's own measurement), and the following render overwrites the rows, so nothing wrong is left on screen. — impact-negligible → report only
- F4. The probe script's Node-lane step runs only `-run TestJavaScriptBehaviorTimeline`. The test the stub line exists for, `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip`, is outside that filter. It is covered by the builder's 80-PASS run and by the planned `queue-kanban-javascript` heavy lane, not by the integrator probe. — impact-negligible → report only

### Requirements Checklist

- [x] Early return at the top of `renderVisibleRows` when the Timeline panel is hidden — delivered (`board-timeline.js:2143-2146`). Both listeners that call it (scroll at `:2878`, capturing toggle at `:2881`) are covered by the one guard. Null-safe: the Node lane's `javascript_behavior_a_test.go` stub returns `null` for unknown ids, and `timelinePanel && timelinePanel.hidden` treats that as not hidden.
- [x] Stale-window re-entry checked — delivered. The guard alone regressed it (hand-back: rows after return were the scrolled-down rows); the smallest fix is one `dispatchEvent` line that reuses the still-live listener.
- [x] No listener release, new state or visibility observer — delivered. The diff adds one local variable in each production file and nothing else stateful.
- [x] Ten hidden scroll steps cause zero childList mutations — delivered (`TestBrowserBehaviorTimelineHiddenViewIgnoresBoardScroll`, 2330 records at base, 0 after).
- [x] Visible Timeline still redraws on scroll — delivered (`TestBrowserBehaviorTimelineViewHasOneScrollSurface` passes; builder ran 18 Timeline browser probes, 0 skipped).
- [x] Node lane still passes — delivered through the D-02 stub line.

Review focus checks:
- Dispatch side effects: `grep` across `web/` finds exactly one `scroll` listener on `#board-main` (`board-timeline.js:2878`) and no `onscroll`, single-quoted `'scroll'` or capture-phase scroll listener on `document`/`window`. `new Event("scroll")` does not bubble. Nothing else can react.
- Ordering: `board-controls.js:28-31` un-hides every panel before `:44-49` runs, so the guard reads `hidden === false` and does not drop the synthetic event. On a first visit no listener exists yet, so the dispatch is a no-op and `renderTimelineView` draws as before. The width-change path (`renderIfPlotWidthChanged`, `:2806`) does not cover this case, because the host returns to the same width it rendered at; that confirms the dispatch is needed.
- Probes: both return `location.href` in the same evaluate call that measures (`_dev/primes/prime-kanban-board.md:25`). Both have guards against a vacuous pass: the hidden probe fails fast unless the panel is hidden, rows were drawn and the board actually scrolled (scrollTop end ≥ 10); the re-entry probe fails fast unless the scrolled-down window shares no row with the top window and scrollTop was clamped to 0 both away and on return. The hidden probe's observer watches `#timeline-scroll` with `subtree: true`, so it sees the rows svg. RED was shown for both (hidden: at base; re-entry: with the guard alone).

### Acceptance Testing

**Result: Pass** (implementation and integration stages)
- `bash do-work/runs/work-2026-10-10-100748/REQ-682-probe.sh`: ran once, about 11 s, no failure output (every failure path in the script prints and exits 1). This covers the guard text check, the Timeline Node-lane tests with `QUEUE_KANBAN_JAVASCRIPT_PROBES=on` and no skips, and the scroll-surface probe plus both new hidden-view probes with the browser lane on, each required to print `--- PASS`.
- `gofmt -l` on both changed Go files: empty.
- Repository gate at the merge (recorded in `## Testing`, not rerun here): exit 0.

### Suggested Additional Testing

- Integration: run the planned `queue-kanban-javascript` heavy lane; it is the run that covers `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip` (F4).
- Manual verification: open the Timeline, scroll down, switch to Activity, change a shared filter, return. The Timeline should show the top rows for the new filter with no flicker (F3 path).
- Deployment: unassessed. A consumer install picks this up only after release and `do-work-update`.
- Live acceptance: unassessed. In a consumer board, open the Timeline once, then scroll Testing with the browser performance panel open; scroll handlers should not show Timeline row rebuilds.

### Scores (on the record — not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All three requirements and five acceptance criteria met |
| Code Quality | 93% | Minimal guard and one dispatch line; two stale comments in the moved helper (F2) |
| Test Adequacy | 95% | One probe per named failure, RED shown for each, anti-vacuity guards |
| Scope | 92% | One approved undeclared stub line; helper move justified |
| Risk | Low | Wasted stale-closure render on filter-reset arrival (F3) |
| Acceptance | Pass | Probe script and gofmt clean; implementation and integration stages |

### Follow-ups created
- None (4 findings report only)

## Review

**Overall: 95%** | <timestamp>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 93% |
| Test Adequacy | 95% |
| Scope | 92% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- F1. Anti-bloat: beyond the REQ's named items the diff adds only `buildTimelineScrollProbeSite` (verbatim move, three callers) and the two result structs `timelineHiddenScrollMeasurement` / `timelineReentryMeasurement` (typed decode, the file's pattern); all used, no new production helper, flag, option or file, no decorative test — impact-negligible → report only
- F2. `timeline_scroll_browser_probe_test.go:309-310` says "Both" probes use the helper (three do) and `:335` says "the measurement below" for a test that is now above it — impact-negligible → report only
- F3. The synchronous `dispatchEvent` at `board-controls.js:49` runs the previous render's `renderVisibleRows` once before `renderTimelineView` on a filter-reset arrival, and adds a second render beside the native scroll event on a non-zero-scrollTop arrival; wasted work under 1 ms, overwritten at once — impact-negligible → report only
- F4. The integrator probe's Node-lane filter `TestJavaScriptBehaviorTimeline` does not run `TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip`, the test the D-02 stub line exists for; covered by the builder's run and the planned heavy lane — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages: REQ-682-probe.sh (guard check, Timeline Node lane unskipped, scroll-surface and both hidden-view browser probes) passed, gofmt clean, merge gate green
**Restatement sweep:** redefined what the Timeline does while hidden (the scroll and toggle listeners now no-op) and on arrival (a synthetic scroll redraw); searched `skills/do-work-board/` (docs, `web/` comments, Go tests), `_dev/primes/prime-kanban-board.md`, `_dev/primes/lessons-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` and `CHANGELOG.md` for scroll listener, listener lifetime/outlive/release, view switch, renderVisibleRows, hidden Timeline, scroll reset, re-entry/arrival and `#board-main` scroll; no stale statement. `board-timeline.js:2876-2877` ("binding here leaks nothing") and `:1998-1999` ("repairs re-entry without a second gate in board-controls.js", about the zero-width case only) still hold; the three `dispatchEvent(new Event("scroll"))` settle calls in `timeline_browser_probe_test.go` run on a visible Timeline
**Suggested testing:** 4 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*
