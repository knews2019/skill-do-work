# REQ-587 implementation plan — give the Timeline view one scroll surface

Shape is fixed by **D-01**: re-point the scroll host to `.board-main`. This plan does not re-open it.

Repo root: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
Board root (all `web/` paths below hang off it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban/`

Every line number in the exploration report was re-checked against the source. All held. Five corrections / additions the exploration did not have are marked **[verified]** below.

---

## 0. Verified corrections to the exploration

1. **`.board-main` has nine direct children, not four.** `#board-notes` (`web/template.html:152`), `#board-anomalies` (`:167`), `#board-findings` (`:193`), then six `.view-panel` sections (`#view-board` `:226`, `#view-calendar` `:302`, `#view-durations` `:314`, `#view-timeline` `:371`, `#view-activity` `:485`, `#view-testing` `:523`). Three of them can precede the view panel, not one. This widens R3 and is why D-05 below is keyed on a condition rather than on `#board-findings`. **[verified]**
2. **The Kanban panel's id is `view-board`, not `view-kanban`.** The Activity probe clicks `[data-view-target="board"]`. **[verified]**
3. **`#view-timeline` has no CSS rule of its own.** `.view-panel` is only `display:none` / `.is-active{display:block}` (`web/board.css:708-713`). The view's own top inset is `.timeline-intro { padding: 16px 24px 0 }` (`web/board.css:2622-2624`). So REQ-585's second half — "zero the view's own padding and fold it into the first element" — has nothing to fold here. Only the board's 24px moves. **[verified]**
4. **`renderTimelineView` re-runs on a filter change** (`web/board-filters.js:187`), not only once. It calls `releaseTimelineListeners()` first (`web/board-timeline.js:1624`), and `addTimelineListener` registers a teardown for every listener (`:156-161`), so moving the scroll listener onto `#board-main` — a node that outlives the render — leaks nothing. **[verified]**
5. **No existing browser probe writes `scrollTop` on `#timeline-scroll`.** `grep -n scrollTop timeline_browser_probe_test.go` returns nothing. So the change breaks no existing probe by making that element unscrollable. **[verified]**
6. **The Node stub's `hostSize()` helper only resizes `timeline-scroll`** (`javascript_behavior_a_test.go:1988-1994`). Its `clientHeight` currently drives the visible-row range. After the change the range reads `boardScrollHost.clientHeight`, which the stub sets to `400` by default — the same number `hostSize(900, 400)` sets today. Row counts in that probe are therefore unchanged, and `hostSize(0, 0)` still exercises the refuse-to-render guard because `plotIsMeasurable()` keeps reading the **width** off `#timeline-scroll`. No Go assertion changes. **[verified]**

---

## 1. Files to modify, in order

### T1 — `web/board.css` (the layout change)

Four edits, all inside the timeline block.

**T1a. Rewrite the block comment (`web/board.css:2614-2620`).** It currently states `.timeline-scroll` is the scroll container by design. That sentence becomes false. Replace with: the board is the scroll surface (REQ-587), the axis is sticky against it, the rows SVG is still full height and still virtualized, and the CSS-pixel/no-viewBox rule is unchanged.

**T1b. `.timeline-scroll` (`web/board.css:2883-2890`) → two declarations.**

```css
.timeline-scroll {
  position: relative;
  cursor: grab;
}
```

`max-height: 58vh`, `overflow-y: auto`, `overflow-x: hidden` and `border-top` all come off. `position: relative` stays (it is the positioned ancestor the rows SVG sits in, and it costs nothing). `cursor: grab` stays — the drag-pan is still this element's.

**T1c. New `.timeline-axis` rule.** It has none today.

```css
.timeline-axis {
  position: sticky;
  top: 0;
  z-index: 1;
  background: var(--bg-base);
  border-bottom: 1px solid var(--line-soft);
}
```

**T1d. The padding move and the focus ring** — see D-05 and D-06 for the exact rules.

### T2 — `web/board-timeline.js` (the coordinate-space change)

One new host lookup, one memo pair, six call sites. Detailed in § 3.

### T3 — `web/board-controls.js` (R2)

One guarded line in `applyView`, detailed in D-07.

### T4 — the Node lane: `generate_test.go`, `javascript_behavior_a_test.go`, `javascript_behavior_b_test.go`, `javascript_behavior_c_test.go`

Add `"board-main"` to every stub id list. Detailed in § 5.

### T5 — new file `timeline_scroll_browser_probe_test.go`

Copied from `activity_scroll_browser_probe_test.go`. Detailed in § 6.

### T6 — rebuild, measure, record

`web/` is embedded, not read from disk — `serve` re-serves the embedded assets, so the CSS/JS edits are invisible until the binary is rebuilt. Rebuild, serve, run the REQ's RED expression by hand, screenshot the Timeline at scroll 0 and mid-scroll (REQ-285: for a rendering change the screenshot is the test), and paste the numbers plus the browser build into the hand-back.

---

## 2. Architectural decisions

### D-02 — `.timeline-scroll` loses `max-height` **and both** `overflow` declarations — DECIDE & STATE

Leaving `overflow-x: hidden` beside an `overflow-y: visible` makes the `visible` compute to `auto` (CSS Overflow: a `visible` paired with a non-`visible` on the other axis becomes `auto`), so the element stays a scroll container. With `max-height` also gone it would measure `scrollHeight == clientHeight` and *look* fixed while still being a scroll box. Deleting both declarations outright is the only version that is both fixed and provable. The probe asserts the computed `overflow-y` is `visible`, not just the height comparison, precisely because the height comparison alone cannot see this mistake (D-12).

*Reason:* the one form of the change that cannot silently half-apply.

### D-03 — the sticky axis paints `var(--bg-base)`, not `var(--surface-1)` — DECIDE & STATE

`_dev/primes/prime-kanban-board.md` states it directly: the surface behind this board's SVG is `<body>`, and `#timeline-scroll` and its `<svg>` are transparent. `body { background-color: var(--bg-base) }` (`web/board.css:222`). Copying the Activity header's `--surface-1` would paint a visibly tinted strip over an untinted chart — the exact defect REQ-321 and REQ-346 each paid for once.

`z-index: 1` matches the Activity header. `.timeline-scroll` is `position: relative` with `z-index: auto` and comes *after* the axis in DOM order, so without a z-index the rows would paint over the stuck axis. `.board-main` carries `container-type: inline-size`, which applies layout containment and therefore establishes a stacking context — so `z-index: 1` is scoped inside the board and can never paint over the top bar.

*Reason:* the prime already answered this question twice; the token is not a matter of taste.

### D-04 — the 1px separator moves from `.timeline-scroll { border-top }` to `.timeline-axis { border-bottom }` — DECIDE & STATE

Today that line renders as the rule under the axis. Once the axis is sticky and the rows scroll under it, a border on the rows box scrolls away with them and the axis loses its bottom edge exactly when it is doing its job. Total height is unchanged (1px leaves one box and joins the one above it), so `rowsOffsetPx` — which is measured, not computed — needs nothing.

*Reason:* same pixel, on the element that stays on screen.

### D-05 — the board's top padding goes to **the first child of `.board-main` that is not hidden**, keyed on the condition rather than on a list — DECIDE & STATE

R3, widened by correction (1): `#board-notes`, `#board-anomalies` and `#board-findings` all sit outside the view panel, all can precede it, and none has a top margin (`.board-notes` and `.board-anomalies` are both `margin: 0 0 20px`, `web/board.css:479-484` and `:563-570`). REQ-585's recipe — move the padding onto the view's first element — reaches none of them, so zeroing the board's top padding pulls whichever strip is showing flush against the top bar.

```css
/* Since REQ-587 the Timeline's rows scroll against the board, and the axis is
   sticky. A sticky element pins to its scroll container's CONTENT box, so the
   board's own 24px top padding would be a band the rows scroll through above
   the axis in plain sight (REQ-585 measured two rows in it on Activity). The
   padding therefore moves off the board and onto whatever the reader actually
   sees first — the notes details, the anomalies strip, the verify-findings
   strip, or the view panel itself. Keyed on "the first child that is not
   hidden" rather than on a list of those four, because which ones exist is a
   property of the data and a hand-maintained list would go stale. Below the
   760px breakpoint the board's own padding is 18px, so the first element gains
   6px there; one rule is worth more than tracking six pixels per breakpoint,
   which is the same trade REQ-585 made. */
.board-main:has(> #view-timeline.is-active) {
  padding-top: 0;
}
.board-main:has(> #view-timeline.is-active) > :not([hidden]) {
  margin-top: 24px;
}
.board-main:has(> #view-timeline.is-active) > :not([hidden]) ~ :not([hidden]) {
  margin-top: 0;
}
```

Third rule zeroes it again on anything preceded by a visible sibling, so exactly one element carries it. `[hidden]` is a sound proxy: `applyView` sets `panel.hidden` alongside `is-active` (`web/board-controls.js:29`), and the three strips are toggled by the `hidden` property too.

`margin-top`, not `padding-top`: the strips have a border and a tint, so padding would open a gap *inside* the coloured box. `.board-main` has `overflow-y: auto`, which establishes a BFC, so a child's top margin cannot collapse through into the board.

Arithmetic is a no-op today. Strip visible: 24 (margin) → strip → 20 (strip's own bottom) → view → 16 (`.timeline-intro`), identical to 24 (board) → strip → 20 → view → 16. No strip: 24 → view → 16 = the same 40px above the heading it has now.

*Reason:* smallest rule that covers every arrangement of three optional siblings, and it cannot be wrong later by omission.

*Alternative rejected:* naming `#board-findings` alone. It fixes the screenshot and breaks the day a reader has notes or an anomaly and no finding.

### D-06 — the chart's focus ring moves onto the sticky axis — **ESCALATE**, proceeding with this default

R5. `.timeline-scroll:focus-visible` draws a 2px inward ring (`web/board.css:2910-2913`). On a box that is now thousands of pixels tall, the ring's top and bottom edges are off screen and the reader sees two vertical lines at the chart's left and right edges. "Tab to the chart" is shipped text (`web/template.html:452`) and the panel's accessible name says the same (`:371`), so the landing signal has to stay visible.

```css
/* The chart is one tab stop and the hint tells the reader to Tab to it, but
   since REQ-587 #timeline-scroll is a full-height element on the board's scroll
   surface: a ring around it puts its top and bottom edges thousands of pixels
   apart, both off screen. The axis is the part of the chart that is always in
   view — pinned to the board's top edge — so it carries the ring instead.
   Inward for the same reason the row ring is: when stuck, the axis is flush
   with the scrollport's top edge and an outward ring loses that edge to it. */
.timeline-chart:has(> #timeline-scroll:focus-visible) .timeline-axis {
  outline: 2px solid var(--accent-claimed);
  outline-offset: -2px;
}
```

`:has()` is required — the axis precedes the scroll host in the DOM, so no sibling combinator reaches it. It already ships in this file (D-05, and REQ-585's rule).

`.timeline-row:focus-visible` (`web/board.css:2937-2940`) is untouched; per-row keyboard focus keeps its own ring.

**Value:** the keyboard affordance the hint text promises stays visible at any scroll position. **Risk:** a reader may read the ring as "the axis has focus" rather than "the chart has focus"; reversal is deleting one rule and restoring the old one. **Escalated because** it is a visible accessibility affordance and reasonable people would pick differently (ring a new wrapper; accept two vertical lines). Record it in the REQ's Open Questions as `- [~]` and continue — per CLAUDE.md, do not stall on it.

### D-07 — entering the Timeline resets `.board-main`'s scroll position to 0 — **ESCALATE**, proceeding with this default

R2, verified: nothing resets `.board-main.scrollTop`. `applyView` only toggles `is-active`/`hidden` (`web/board-controls.js:19-31`), and the only other `scrollTop` writers in `web/` are `board-detail.js:657-658` (the drawer's own boxes, not the board) and `board-timeline.js`.

Today the Timeline page is short — the chart is capped at 58vh — so arriving from a Kanban board scrolled 800px down clamps the board to a small number and the reader lands near the top. After this change the chart is the tall thing on the page, nothing clamps, and the reader lands 800px into the chart with correct rows and no reason on screen for where they are.

```js
// The board is the Timeline's scroll surface since REQ-587, and .board-main's
// scroll position is shared by every view with nothing resetting it between
// them. Arriving from a Kanban board scrolled 800px down would drop the reader
// into the middle of the chart — rows correct, position arbitrary. Before this
// view the page was short enough that the board clamped to near the top, so
// resetting is what keeps arrival looking the way it does today.
if (viewState.view === "timeline") {
  document.getElementById("board-main").scrollTop = 0;
}
```

Placement: in `applyView`, immediately after the `viewPanels` toggle loop and **before** the `renderedOnce.timeline` block at `web/board-controls.js:67-70`, so the first `topVisibleRowAnchor()` of the first render reads 0. `applyView` has exactly two callers — a view-button click (`:187`) and init (`web/board.js:79`) — so this cannot fire mid-session; the only redundant case is clicking the Timeline button while already on the Timeline.

**Value:** arrival on the Timeline is predictable and matches what the view does today. **Risk:** leaving the Timeline and coming back loses your place in the chart, which today survives because `#timeline-scroll` keeps its own `scrollTop`. Reversal is deleting three lines. **Escalated because** it is user-visible behaviour with a real cost on one path.

*Alternative if the maintainer values the preserved place:* remember `.board-main.scrollTop` on leaving the Timeline and restore it on entry (~6 lines of state in `board-controls.js`). Reproduces today's behaviour exactly, at the price of new state in the view switcher.

*Rejected:* resetting on every view switch — it would take a scroll position away from the Kanban, Durations and Testing views, which is outside this REQ.

### D-08 — `rowsOffsetPx` and `stickyAxisHeightPx` are measured once per render and memoized, invalidated in `renderAll` beside `invalidatePlotWidth()` — DECIDE & STATE

`renderVisibleRows` **is** the scroll listener (`web/board-timeline.js:2777`) and the drag-pan calls `renderAll` once per frame. It already forces one synchronous layout per frame (`scrollHost.clientHeight` plus `timelineMeasureLabelAdvance`, measured 0.036–0.046ms on Chromium 1194, documented at `:2063-2070`). Two `getBoundingClientRect()` calls inside the handler would add two more forced layouts to every frame of a drag. Neither number can change without a re-render: both move only when content above the chart changes height — the findings strip appearing, the toolbar wrapping, the summary line rewrapping — and every one of those is a render. Same memo shape as `measuredPlotWidthPx` / `invalidatePlotWidth` (`:2004-2017`), so there is one caching idiom in the file, not two.

### D-09 — the visible-range viewport stays `boardScrollHost.clientHeight`, with no axis subtraction — DECIDE & STATE

Not subtracting the sticky axis's height renders a **superset** of what is visible — at most one extra row of nodes, hidden behind the axis. Subtracting `rowsOffsetPx` from the scroll *position* is not optional: without it the range walks off the bottom and the top rows go blank. One number has to be right; the other only has to be generous.

### D-10 — the `scroll` listener moves to `#board-main`; wheel, pan, pointer capture, keyboard, focus, hover and the ResizeObserver stay on `#timeline-scroll` — DECIDE & STATE

**Must move (S11).** `#timeline-scroll` stops scrolling, so it stops emitting `scroll`. Left there, the virtualized range freezes at whatever the last render computed: the reader scrolls the board, the first screenful of rows stays drawn where it was, and everything below it is blank SVG. This is the single failure most likely to be shipped by a change that otherwise looks complete.

**Must not move (S17–S22).** Four reasons: the wheel handler's zoom anchor is `scrollHost.getBoundingClientRect().left`, a per-chart x origin whose comment (`:2957-2963`) says both the axis and the plot must take the rect from the scroll host so they share one x scale; a `ctrl`/`meta`-wheel listener on `#board-main` would call `preventDefault` over **every other view**, since the board is shared; the pan's `setPointerCapture` is what makes a drag not a click, and capturing on the board would retarget synthesized clicks everywhere (REQ-336/337 paid for that once); `cursor: grab` / `.is-panning` belong to the chart's box, not the page's.

The ResizeObserver (S10) and `liveHostWidth()` (S7) stay on `#timeline-scroll` because they measure the **plot's width**, which is still that element's. R6 is accepted, not fixed: the measured width grows by the vertical scrollbar the inner box no longer draws (~15px), the observer delivers the change and re-renders, and every recorded plot-width number in comments and probes shifts by that much.

### D-11 — the Node stub resolves `board-main` through `getElementById`, never a class selector — DECIDE & STATE

The stub's `querySelector` returns `null` unconditionally (`generate_test.go:3357`). A `document.querySelector(".board-main")` in the renderer would make `renderTimelineView` hit its guard and return on every Node probe, silently turning several of them into assertions about nothing. `document.getElementById("board-main")` plus one map entry is the whole cost, and `makeStubNode` already supplies `scrollTop: 0`, `clientHeight: 400` and a `getBoundingClientRect` with `top: 0` — so `rowsOffsetPx` evaluates to `0` and the lane's row counts are byte-identical to today (§ 0, correction 6).

### D-12 — the probe scopes "one scroll surface" to the two elements the REQ's RED expression names — DECIDE & STATE

R7: `.timeline-table-scroll` (`web/board.css:3101-3105`) is a third `max-height: 360px; overflow: auto` box, inside a `<details>` that is closed by default. It is deliberately still a scroll box and is not in scope. The probe measures `#board-main` and `#timeline-scroll` by id, never "every element on the page", and says so in a comment. It runs the REQ's RED expression verbatim and asserts it returns `[true, false]`.

### D-13 — delete the dead `timelineVisibleRowRange` — DECIDE & STATE (optional, same commit)

`web/board-timeline.js:1056-1060` has no callers in `web/`; the only other mention is a comment at `generate_test.go:3404`. CLAUDE.md's "delete before you add" applies, and this is the render path being touched. If removing it reddens anything, drop the deletion — it is not part of the REQ.

---

## 3. The five scroll-position sites, converted

New lines at the top of `renderTimelineView` (`web/board-timeline.js:1615-1623`):

```js
var boardScrollHost = document.getElementById("board-main");
...
if (!summaryNode || !axisHost || !scrollHost || !boardScrollHost || !tableBody) {
  return;
}
```

Helpers, placed beside the plot-width memo (`:2004-2017`) so the two caches read as one idiom:

```js
var measuredRowsOffsetPx = null;
var measuredStickyAxisHeightPx = null;

// Where the rows sit inside the BOARD's scroll extent: the distance from the
// board's scroll origin to the top of the rows box. scrollTop is measured from
// the padding box and .board-main has no border, so its client rect top IS that
// origin, which is why this holds whatever the board's top padding is.
function rowsOffsetPx() {
  if (measuredRowsOffsetPx === null) {
    measuredRowsOffsetPx =
      scrollHost.getBoundingClientRect().top -
      boardScrollHost.getBoundingClientRect().top +
      boardScrollHost.scrollTop;
  }
  return measuredRowsOffsetPx;
}

// Read, not restated (REQ-322): the axis is TIMELINE_AXIS_HEIGHT plus its own
// bottom border, and a scroll target that guessed would put a focused row
// under the strip by exactly the part it guessed wrong.
function stickyAxisHeightPx() {
  if (measuredStickyAxisHeightPx === null) {
    measuredStickyAxisHeightPx = axisHost.getBoundingClientRect().height;
  }
  return measuredStickyAxisHeightPx;
}

function invalidateBoardScrollGeometry() {
  measuredRowsOffsetPx = null;
  measuredStickyAxisHeightPx = null;
}

// The board's scroll position expressed in rows coordinates — the number every
// site below used to read straight off the inner box.
function rowsScrollTop() {
  return Math.max(0, boardScrollHost.scrollTop - rowsOffsetPx());
}

// Where to put the board so a row lands just below the pinned axis. Clamped
// because a row near the chart's top would otherwise ask for a negative
// scrollTop; the browser clamps anyway, and the Node stub does not.
function boardScrollTopForRowTop(rowTopPx) {
  return Math.max(0, rowsOffsetPx() + rowTopPx - stickyAxisHeightPx());
}
```

In `renderAll` (`web/board-timeline.js:2680-2683`), beside the existing invalidation and **before** the `plotIsMeasurable()` guard — the measurement is lazy, so the refusing path still measures nothing:

```js
invalidatePlotWidth();
invalidateBoardScrollGeometry();
```

**S4 — `topVisibleRowAnchor()` (`:1751-1763`).** Read the converted position once and use it for both the comparison and the offset:

```js
var anchorScrollTop = rowsScrollTop();
...
displayItem.topPx + displayItem.height > anchorScrollTop
...
return { id: displayItem.row.id, offsetPx: anchorScrollTop - displayItem.topPx };
```

Swapping in `boardScrollHost.scrollTop` naively makes every recorded offset wrong by `rowsOffsetPx` (250–400px on this board), which drops the reader several rows away from the REQ they were reading whenever the window chips move.

**S5 — `refreshWindowRows()` (`:1788-1797`).** Ordering unchanged, comment kept and extended:

```js
rowsSvg.setAttribute("height", timelineDisplay.height);
boardScrollHost.scrollTop = rowsOffsetPx() + displayItem.topPx + anchor.offsetPx;
```

REQ-319's lesson survives the move and gets stronger: `boardScrollHost.scrollHeight` depends on the rows SVG height, so writing before the resize still clamps to the old maximum. Add one sentence to the existing comment saying the clamp is now the board's, not the box's. Measuring `rowsOffsetPx()` here is safe even though the SVG has not been resized yet — the offset describes content **above** the chart, which the resize does not move.

**S8 — `renderVisibleRows()` (`:2116-2118`).**

```js
var visible = timelineVisibleDisplayRange(
  timelineDisplay.items, rowsScrollTop(), boardScrollHost.clientHeight
);
```

**S12 — `scrollTimelineRowIntoView()` (`:2796-2808`).** Both directions convert, and the *upward* comparison must use the axis-clear top or a row scrolled "to the top" lands underneath the strip and is invisible while focused — which reads as the down-arrow doing nothing:

```js
var visibleTopPx = rowsScrollTop() + stickyAxisHeightPx();
var visibleBottomPx = rowsScrollTop() + boardScrollHost.clientHeight;
if (rowTopPx < visibleTopPx) {
  boardScrollHost.scrollTop = boardScrollTopForRowTop(rowTopPx);
} else if (rowBottomPx > visibleBottomPx) {
  boardScrollHost.scrollTop =
    rowsOffsetPx() + rowBottomPx - boardScrollHost.clientHeight;
}
```

**S23 — jump to open work, `#timeline-zoom-now` (`:3186`).**

```js
boardScrollHost.scrollTop = openDisplayItem
  ? boardScrollTopForRowTop(openDisplayItem.topPx)
  : boardScrollTopForRowTop(0);
```

The `: 0` fallback stops meaning "scroll the board to the very top" and starts meaning "bring the chart's first row just under the axis", which is what it meant when the inner box was the surface.

**S11 — the listener (`:2777`).**

```js
addTimelineListener(boardScrollHost, "scroll", renderVisibleRows);
```

**Unchanged, deliberately:** S1/S3/S6/S9/S13–S22 (`scrollHost` as the DOM host, focus target, hover, keyboard, wheel, pan, capture) and S7/S10 (width).

---

## 4. What must not move — one-line summary for the diff review

The wheel-zoom listeners (`:2985` on the chart, `:2986` on the axis), the whole pointer-pan block (`:3007-3106`) and `setPointerCapture`/`releasePointerCapture` stay bound to `#timeline-scroll`. See D-10 for the four reasons.

---

## 5. The Node lane — every site to update

Add `"board-main"` to the id list so `getElementById("board-main")` resolves. Five list literals across four files:

| File | Line | What it is |
|---|---|---|
| `generate_test.go` | 3344-3346 | `timelineRenderDomStubPreamble`'s id list — the shared one |
| `javascript_behavior_a_test.go` | 2148-2150 | re-creates every host between renders (no-match case) |
| `javascript_behavior_a_test.go` | 2156-2158 | re-creates them again (match-again case) |
| `javascript_behavior_b_test.go` | 1524-1526 | local host list |
| `javascript_behavior_c_test.go` | 2191-2193 | local host list |

Two more sites were checked and need **no** change: `javascript_behavior_a_test.go:2132-2135`, the local `getElementById` override, already falls through to `timelineStubHosts[nodeId] || null`; and `javascript_behavior_a_test.go:2002-2003` / `generate_test.go:3394`, which only read or clear `timeline-scroll`.

Miss any of the five and `renderTimelineView` hits its guard and returns for that probe — the probe still passes, measuring nothing. Cheap check after the edit: `grep -rn '"timeline-period-state"' *_test.go` must return the same five list literals, each now carrying `"board-main"`.

`generate_test.go:2759-2767` (`TestTimelinePanelStatesItsKeyboardInteraction`) asserts `tabindex="0"` on the `id="timeline-scroll"` opening tag. The template attribute is untouched, so it stays green — re-run it as a regression anyway.

No Go test asserts `58vh`, `max-height`, or the `border-top` (only a comment in the Activity probe), so the CSS change has no string lock-in to update.

---

## 6. Testing approach

### The new browser probe — `timeline_scroll_browser_probe_test.go`

One test, `TestBrowserBehaviorTimelineViewHasOneScrollSurface`, copied from `activity_scroll_browser_probe_test.go` (265 lines) and keeping its seven shape rules: two named constants with reasons; a synthetic fixture rather than the live queue; one measurement struct carrying real heights rather than booleans; REQ-291's guard applied before any comparison; two `requestAnimationFrame`s after every `scrollTop` write; scoping asserted in both directions; a `t.Logf` with every number so a green run still records the measurement. First line is `lookupBrowserForBehaviorProbe(t)`; the session is `startTrustedInputBrowserSession(..., "--window-size=1600,900")`.

**Fixture** (`writeVerifyFixture` + `buildBoard` + `generateStaticSite`, the Activity probe's exact spine):

- `timelineScrollProbeRequestCount = 120` archived REQs with `created_at`/`claimed_at`/`completed_at` — at `TIMELINE_ROW_HEIGHT = 18` plus group headers that is well over 2000px of rows against a ~830px board viewport, so the chart is taller than the fold by more than two screenfuls. The constant's comment states that arithmetic, because a fixture that does not overflow makes every assertion below pass while measuring nothing.
- **one extra REQ that lands in the anomalies strip** — `status: completed` with no `completed_at` and no commit hash, which is exactly the "terminal REQ whose completion instant cannot be resolved" the strip exists for. This is what puts a visible sibling *above* the view panel and makes D-05 testable. (`#board-findings` carries the same `.board-anomalies` class and the same position, so pinning the anomalies strip pins the same rule; producing a real verify finding costs more and proves nothing extra.)
- `timelineScrollProbeRowsScrollPixels = 1500`, with the reason in the comment: greater than one board viewport (~830px) plus two overscans (`TIMELINE_OVERSCAN_ROWS = 4` × 18px × 2 = 144px), which is what makes the before/after rendered-row sets provably disjoint.

**Steps:** click `[data-view-target="timeline"]`; click the All-days period chip (the REQ's RED case names that window); wait for `document.querySelectorAll('#timeline-scroll [data-detail-id]').length > 0`; measure; scroll; measure again; click back to `[data-view-target="board"]`; measure the padding.

**Guards first (REQ-291), all `t.Fatalf`:**

- G1 rows actually drew: `renderedRowCountBefore > 0`, and the rows SVG's height exceeds the board's client height (the chart really is taller than the fold).
- G2 every measured box has positive height: `boardClientHeight`, `timelineScrollClientHeight`, `axisHeight`.
- G3 the board actually scrolls: `boardScrollHeight > boardClientHeight` — otherwise "one scroll surface" is indistinguishable from "nothing scrolls".
- G4 the scroll landed: `boardScrollTop >= rowsOffsetPx + 1500` (read back, not assumed).

**Assertions:**

- A1 **one scroll surface.** `timelineScrollScrollHeight <= timelineScrollClientHeight + 1`, **and** the computed `overflow-y` on `#timeline-scroll` is `"visible"`. The second is not redundant — it is the only assertion that sees D-02's trap, because a leftover `overflow-x: hidden` with the `max-height` gone still reports equal heights while the element remains a scroll container.
- A2 **the REQ's own RED expression, verbatim**, asserted to return `[true, false]`. It is the sentence the REQ will be judged by; running anything else is running a different test.
- A3 **the axis sticks.** After the scroll, `axisTopOffset` (axis rect top minus board rect top) is within ±1 of 0, and `rowsVisibleAboveAxis == 0` — rows painted in the band between the board's top inner edge and the axis top, computed exactly as the Activity probe's `rowsVisibleAboveHeader`.
- A4 **the virtualized range follows the scroll** — the REQ's GREEN condition, and the assertion the Activity probe had no need for. Collect `data-detail-id` values before and after the 1500px scroll; assert both sets are non-empty and their intersection is empty, and that at least one rendered row's client rect intersects the board's visible rect after scrolling. Non-empty-and-disjoint proves the window moved; the intersection check proves the rows moved *with* it rather than being drawn somewhere nobody can see.
- A5 **the padding move, both directions.** `getComputedStyle(boardMain).paddingTop == "0px"` while on the Timeline; after clicking back to `[data-view-target="board"]`, it is not `"0px"`.
- A6 **R3 — nothing sits flush against the top bar.** The anomalies strip's rect top minus the board's inner top edge is within ±1 of 24 before any scrolling. Failure message says the strip is flush against the top bar, which is what the reader would see.
- A7 **the view's own inset survived.** `.timeline-intro`'s text top offset is > 0 (the Activity probe's `summaryTextTopOffset` shape, adapted).

**Not asserted, deliberately:** the focus ring. `:focus-visible` after a programmatic `focus()` is engine- and last-input-dependent, so a probe for it would pin the engine's heuristic rather than the rule. Evidence for D-06 is a screenshot, per REQ-285.

### What the Node lane covers

Nothing new. It sees the DOM the client builds and never a layout, so it cannot answer "how many things scroll" — that is why T5 exists. Its job here is to prove the renderer still renders after the host split: with `"board-main"` in the stub map, `rowsOffsetPx()` is 0 and `boardScrollHost.clientHeight` is 400, the same numbers `#timeline-scroll` supplied before, so every existing row-count and row-id assertion must come back byte-identical. A diff in those numbers means a conversion is wrong, not that an expectation needs updating.

### What only a real engine can settle

Whether `#timeline-scroll` is still a scroll container (D-02's trap is a *computed*-value question); whether the axis actually pins and at what offset; whether any row paints above it; the real value of `rowsOffsetPx` and `stickyAxisHeightPx`; whether the virtualized range tracks the board's scroll; the plot-width shift from R6; and how the strip and the heading sit once the padding has moved.

### Existing probes to re-run as regression

Heavy lane (`QUEUE_KANBAN_BROWSER_PROBES=on`, and set `QUEUE_KANBAN_BROWSER` — the lane reports **skipped** without it, and a skip is not a pass):

1. `TestBrowserBehaviorTimelineBarsSurviveTheDetailDrawerOpening` (`timeline_browser_probe_test.go:1319`) — the ResizeObserver/width invariant, the one R6 moves.
2. `TestBrowserBehaviorTimelineRowListIsOneTabStop` (`:3330`) — roving tabindex over the virtualized list, which now depends on S12's conversion.
3. `TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving` (`:627`) — drag-pan.
4. `TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` (`:3124`) — pointer capture.
5. `TestBrowserBehaviorActivityViewHasOneScrollSurface` — the sibling view's padding rule must stay scoped now that a second `:has()` rule targets the same property.

Plus the Node/string lane: the four `javascript_behavior_*` timeline probes and `TestTimelinePanelStatesItsKeyboardInteraction` (`generate_test.go:2737`).

---

## 7. Task count

**Six tasks: T1–T6. That is over the five-task threshold, and the orchestrator should know it.**

The recommendation is still **not to split the REQ**, because no subset ships on its own:

- T1 (CSS) alone freezes the virtualized rows — the view is worse than before.
- T2 (JS) alone reads a board that is not the scroll surface — `rowsOffsetPx` is real but the inner box still scrolls, so both surfaces stay.
- T3 is three lines that only make sense once T1+T2 land.
- T4 is forced by T2 and is not optional — skip it and four Go files quietly stop testing the renderer.
- T5+T6 are the REQ's own GREEN condition ("measured in a real engine and recorded in the hand-back; a browser probe in the existing lane pins it").

The count reflects the number of *files with different failure modes*, not independent deliverables. If the orchestrator wants a split anyway, the only clean cut is **T1+T2+T3+T4 as the change** and **T5+T6 as a follow-up probe REQ** — at the cost of merging a layout change with no engine-measured proof, which is what the prime's "a chart's correctness is partly a claim about pixels" entry exists to prevent.

---

## 8. Requirement coverage check

| Requirement (source) | Covered by |
|---|---|
| Leave exactly one scroll surface on the Timeline (What) | T1b (D-02); probe A1, A2 |
| Same style as the Activity fix — board scrolls, rows are content (What) | T1b, T1d (D-05 reuses REQ-585's `:has()` scoping) |
| The time axis sticks to the top edge, as the column header does on Activity (What) | T1c (D-03, D-04); probe A3 |
| Chart becomes full height as ordinary content (Context, M1) | T1b |
| Visible-row math reads the board's position minus the chart's offset (Context) | T2 / S4, S8 (D-08, D-09); probe A4 |
| The axis becomes sticky under the board's top edge (Context) | T1c; probe A3 |
| Wheel-to-zoom (Ctrl/Cmd + wheel) keeps working (Context) | D-10 (listeners stay); regression probes 3–4 and a manual check in T6 |
| Drag-to-pan keeps working (Context) | D-10; regression probes 3–4 |
| Five scroll-position sites converted (Context, ~1618 / 1757-1797 / 2077-2117 / 2803-2806 / 3186) | T2 / S4, S5, S8, S12, S23 |
| Reuse REQ-585's padding-move mechanism, not a second one (Context) | D-05 (same `:has()` scoping, extended to cover the three strip siblings) |
| Keyboard: Tab to the chart, arrows move/pan, +/− zoom (Context) | D-06 (the ring stays visible); D-10 (keydown stays on the chart); S12 conversion; regression probe 2 |
| Block comment at board.css ~2579 stops being true (Context) | T1a |
| RED expression returns `[true, false]` (Red-Green Proof) | probe A2; T6 by hand on the served board |
| Rows still render correctly at every scroll position (GREEN) | probe A4 |
| A row 400px below the fold appears after scrolling there (GREEN) | probe A4 (1500px, disjoint id sets) |
| The row under the pointer stays put when the window chips change (GREEN) | T2 / S4+S5 conversion — **see gap below** |
| Measured in a real engine and recorded in the hand-back (GREEN) | T6 |
| A browser probe in the existing `*_browser_probe_test.go` lane pins it (GREEN) | T5 |

**One requirement no task pins automatically: "the row under the pointer stays put when the window chips change."** It is the S4+S5 anchor behaviour, and it needs two window states and a comparison of one row's screen position across them — the Activity probe has no equivalent and neither does any existing timeline probe. Two ways to close it, and the plan does **not** silently drop it:

- **Cheapest honest option (recommended):** verify it by hand in T6 — note a REQ id at the top of the viewport, click a different period chip, confirm the same id is still at the top — and record the two screenshots plus the id in the hand-back. The REQ's GREEN says "measured in a real engine and recorded in the hand-back", which this satisfies.
- **If it should be pinned:** add an eighth assertion to T5 — record `getBoundingClientRect().top` of the top visible row's `[data-detail-id]`, click a wider period chip, wait, and assert the same id is within a few pixels of the same offset. Roughly 25 lines. It is the assertion most likely to be flaky, because it depends on the anchor row surviving the window change, so it needs its own guard that the id is still in the display list at all.

The builder should pick one and state which. Everything else in the REQ maps to a task above.
