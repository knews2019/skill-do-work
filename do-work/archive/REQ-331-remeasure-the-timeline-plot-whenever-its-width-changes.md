---
id: REQ-331
title: "[impact-critical] Remeasure the timeline plot whenever its width changes"
status: completed
created_at: 2026-08-23T12:20:00Z
claimed_at: 2026-08-23T13:05:00Z
completed_at: 2026-08-23T13:40:00Z
commit: 406c64e
route: B
user_request: UR-066
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# Remeasure the Timeline Plot Whenever Its Width Changes

## What

The plot width is measured once per render and memoised, and the memo is dropped only by `renderAll`. The
comment at `web/board-timeline.js:1246-1248` claims "the resize listener IS `renderAll`, so the two moments
the width can change are both covered". There are more than two moments. Two of them destroy the chart
outright: **opening the detail drawer blanks it**, and **a browser resize while another view is on screen
leaves it permanently crushed into a 120-pixel strip**.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`. Chose a ResizeObserver on the scroll host over enumerating callers, because the condition ("this box changed") is what the rule is about and an enumeration would need the next layout change added to it. Chose a render-level guard over a defensive branch in `plotWidth`, because the branch would be unreachable.
- [x] **[APPLY]:** `web/board-timeline.js`, plus `generate_test.go` and `timeline_browser_probe_test.go` for the two probes. **`web/board-controls.js` was NOT touched** — the observer fires when the host gains a box on view activation, so the `renderedOnce.timeline` gate needed no change, and a file shared by five views stays untouched. `requestPanRender` renamed to `requestFrameRender` now that it serves two callers.
- [x] **[UNIFY]:** Verified:
  - `web/board-timeline.js` — `node --check` clean; no stale `requestPanRender` / `timelinePanRenderFrame` reference; the `ResizeObserver` construction is feature-guarded for hosts that lack it; the observer is disconnected through the existing teardown registry.
  - `generate_test.go` / `timeline_browser_probe_test.go` — `gofmt` and `go vet` clean; `plotEdgeText` added after a real failure message printed a pointer address instead of a measurement.
  - `bash _dev/tests/maintainer-verify.sh` exit 0.
  - Mutations: removing the observer, and neutering `invalidatePlotWidth`, each fail the browser probe with "the chart is blank"; removing the render guard fails the Node probe with the exact eight-row truncation from the report. The zero-width branch in `plotWidth` survived every mutation, so it was deleted rather than shipped untestable.
  - Live render before/after: 55 segments / 0 visible after a row click, now 52 segments / 52 visible re-laid out against the 866px host; the hidden-resize case 8 segments at x≈301, now 49 at x 1033–1066.

## Why

Judged `impact-critical`. Clicking a row is the affordance `.timeline-hint` explicitly advertises — "Click a
row for its full detail" — and doing it makes every bar on the chart disappear. There is no error, no empty
state, and no hint that a window move would repair it. A reader who follows the instruction printed under the
chart loses the chart.

## Context

Measured in Chromium 1194 at a 1600×900 viewport against a board generated from this repo (317 REQs):

| Moment | `#timeline-scroll` clientWidth | Segments drawn | Segments actually inside the host |
|---|---|---|---|
| Timeline opened | 1496 | 50 | 50 |
| after clicking one row | **866** | 50 | **0** |
| after resizing to 1200×800 while the Board view was on screen, then returning | 1096 | **8** | 8 (all at x≈301) |

**The drawer.** `.detail-drawer` occupies its own grid column at `--detail-panel-width: 620px`
(`board.css:75`, `:1526`). Opening it re-lays out the grid and narrows the scroll host by 630px. Nothing
calls `invalidatePlotWidth` (`:1256`) — its only caller is `renderAll` (`:1713`) — so `xOfEpoch` keeps
converting instants against the old 1300px plot. The bars keep x values of 1402.7…1453.5, which is now past
the right edge, and `drawSegment`'s clamp (`:1289`) drops every one of them.

**The hidden resize.** `addTimelineListener(window, "resize", renderAll)` (`:2004`) fires while
`#view-timeline` is `hidden`, so `plotWidth` (`:1252`) measures `clientWidth` 0 and memoises
`Math.max(120, 0 - 184 - 12)` = **120**, and `timelineVisibleRowRange` gets `clientHeight` 0 and returns 8
rows. Returning to the view does not repair it: `board-controls.js:49-51` re-renders the Timeline only when
`renderedOnce.timeline` is false, and it is already true. Only a window move (a chip, a zoom, a drag) calls
`renderAll` again and fixes it.

Two smaller members of the same family, in scope because the fix decides them:

- The aria-live row readout (`#timeline-readout`) is written only by `mousemove` / `mouseleave`
  (`:1810-1826`). A window move that removes the described row leaves it announcing a REQ that is no longer
  drawn anywhere. `renderAll` refreshes everything else and never touches it.
- A mouse sweep across the rows rewrites that same `role="status" aria-live="polite"` node once per
  `mousemove`, queueing dozens of announcements for one gesture.

## Detailed Requirements

1. The plot width is remeasured whenever the plot's box can have changed — keyed on the condition, not on a
   list of callers. A `ResizeObserver` on the scroll host is the mechanism that makes the condition the rule;
   if the implementation instead enumerates callers, the comment at `:1246` must be replaced by one that
   states why that enumeration is complete, and the drawer must be in it.
2. A measurement taken while the view is hidden is never memoised, and never becomes the row range. A
   zero-size box is not a width; treat it as "not measurable yet" and re-measure on activation.
3. Re-entering the Timeline after any layout change repairs the chart without the reader having to move the
   window. `renderedOnce.timeline` may still gate the *first* render, but it must not gate a repair.
4. `#timeline-readout` is cleared by `renderAll` like every other piece of prose the view owns, so it cannot
   describe a row the current window does not draw.
5. The hover readout does not queue an announcement per `mousemove`. Write it only when the row under the
   pointer actually changes.

## Constraints

- Keep the measurement to one place. The memo exists because `clientWidth` forces a synchronous layout and
  `xOfEpoch` is called several times per row per frame (`:1243-1249`); a fix that measures per call would
  cost a drag its frame rate.
- Do not make the row SVG's x scale differ from the axis SVG's. Both read `plotWidth()` on the scroll host,
  and that is what keeps the axis aligned with the bars.
- `board-controls.js` is shared by five views. Any change there must not re-render the other four.

## Red-Green Proof

**RED prompt/case:** Open the board, open Timeline, count `rect.timeline-segment` nodes whose bounding box
intersects `#timeline-scroll`'s box. Click any row. Count again. Separately: from Timeline switch to Board,
resize the window, switch back, count the segments.

**Why RED now:** 50 intersecting segments before the click, **0** after. After the hidden resize, 8 segments
exist at all, every one at x≈301, and they stay that way until the window is moved.

**GREEN when:**
- The intersecting-segment count after opening the drawer is greater than zero and the bars are laid out
  against the narrowed plot — the same instants, at the new scale.
- Closing the drawer restores the wide layout without a window move.
- After resizing while another view is on screen, returning to Timeline draws the same row count the window
  admits (317 at Fit all on this board), not 8.
- `#timeline-readout` is empty after any window move that drops the hovered row.
- A browser behaviour probe asserts the drawer case by measuring `getBoundingClientRect()` intersections
  before and after the drawer opens, and returns `location.href` alongside the numbers.

**Validation:** Inferred during capture; every number in the table above was measured in a browser this
session, not estimated.

## Full Context

See `do-work/user-requests/UR-066/input.md`.
