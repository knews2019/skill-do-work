# REQ-587 pre-exploration — give the Timeline view one scroll surface

Read-only exploration. Nothing edited.

Repo: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
REQ: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/queue/REQ-587-give-the-timeline-view-one-scroll-surface.md`

All paths below are absolute unless they start with `web/`, in which case the root is
`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban/`.

---

## 1. What REQ-585 actually did

Archived record: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/UR-120/REQ-585-give-the-activity-view-one-scroll-surface.md`
Builder commit `94532c45`, merge `c08ac2b4`. Two files, 298+/8-:

- `skills/do-work-board/tools/queue-kanban/web/board.css` (41 changed lines, activity block only)
- `skills/do-work-board/tools/queue-kanban/activity_scroll_browser_probe_test.go` (265 new lines)

### The mechanism, in three declarations

**(a) Delete the inner scroll box.** `.activity-table-scroll` lost exactly two declarations:

```css
 .activity-table-scroll {
-  max-height: 70vh;
-  overflow: auto;
   margin-top: 8px;
   font-size: 12px;
 }
```

**(b) Sticky header — nothing was added.** The `thead th` rule already carried
`position: sticky; top: 0; z-index: 1; background: var(--surface-1)`
(`web/board.css:3278-3284`). Removing the inner scroll container is what re-pointed
that sticky at `.board-main` instead of at the table box. No JS involved.

**(c) The padding move, scoped with `:has()`.** This is the part REQ-587 is told to
reuse. Two paddings zeroed, one restored on the view's first element:

```css
.board-main:has(> #view-activity.is-active) {
  padding-top: 0;
}

#view-activity {
  padding: 0 24px 24px;   /* was: 16px 24px 24px */
}

.activity-summary {
  margin: 0 0 4px;
  padding-top: 40px;      /* new: 24 (board) + 16 (view) */
  font-size: 13px;
  color: var(--ink-soft);
}
```

`web/board.css:3219-3221`, `:3223-3225`, `:3227-3232`.

**Why the padding move is needed at all** (measured, not reasoned — decision D-02 in the
record): a sticky header inside a scroll container pins to that container's **content
box**, not its padding box, so `.board-main`'s 24 px top padding is a band the rows scroll
through in plain sight. Measured on Chrome 152.0.7977.76: two rows visible above the stuck
header at 700 px of scroll.

**Why `:has()` and not a class toggle** (D-03): one rule in the file that owns layout, no
second place in `board-controls.js` deciding "this is the Activity view". `board-controls.js`
was in the write boundary and left untouched.

**RED → GREEN measured**: board 665/677 → 665/3657; table 530/3509 → 3509/3509; header
55.5 px from the board's top edge → 0.0 px; 0 rows above it.

### Lessons the record leaves for REQ-587 (verbatim, "Worth knowing")

> When a view removes an inner scroll box in favour of the board's own scroll, the board's
> top padding has to move onto the first element of that view, scoped to that view.
> `.board-main:has(> #view-activity.is-active)` is the pattern; REQ-587 (the Timeline's same
> fix) should reuse it rather than invent a second scoping mechanism.

The record's Discovered Tasks also names REQ-587 explicitly as the same shape.

---

## 2. The Timeline's current scroll model

### CSS

Block comment, `web/board.css:2604-2608`:

```css
/* ---- timeline view ------------------------------------------------------
   One bar per REQ the visible time window covers. Rows are virtualized, so .timeline-scroll
   is the scroll container and the SVG inside it is full height while only the
   visible rows carry nodes. Everything is drawn in CSS pixels with no viewBox,
   which is what keeps a pointer's x equal to a plot x at every zoom level. */
```

The rule itself, `web/board.css:2883-2900`:

```css
.timeline-scroll {
  position: relative;
  max-height: 58vh;
  overflow-y: auto;
  overflow-x: hidden;
  border-top: 1px solid var(--line-soft);
  cursor: grab;
}

.timeline-scroll.is-panning {
  cursor: grabbing;
}

/* Drawn inward: .timeline-scroll is an overflow container flush under the axis,
   so a positive offset would put the ring where the surrounding card clips it. */
.timeline-scroll:focus-visible {
  outline: 2px solid var(--accent-claimed);
  outline-offset: -2px;
}
```

Its parent, `web/board.css:2812-2815`:

```css
.timeline-chart {
  position: relative;
  padding: 8px 24px 0;
}
```

`.timeline-axis` has **no CSS rule of its own** — only `.timeline-axis-svg` (`:2817`),
`.timeline-axis-tick` (`:2823`), `.timeline-axis-label` (`:2828`). It is a bare `<div>`
inside `.timeline-chart`.

The board's own scroll container, `web/board.css:443-459`: `overflow-y: auto;
padding: 24px 28px 56px; container-type: inline-size`. At the 760 px breakpoint it drops to
`padding: 18px 16px 48px` (`web/board.css:2081-2083`).

### Markup

`web/template.html:444-448`:

```html
<div class="timeline-chart">
  <div class="timeline-axis" id="timeline-axis"></div>
  <!-- tabindex="0" is what makes the chart itself reachable by Tab, so the
       arrow-key pan and +/− zoom below have somewhere to land. -->
  <div class="timeline-scroll" id="timeline-scroll" tabindex="0"></div>
</div>
```

`#board-main` is `<main id="board-main" class="board-main" tabindex="-1">`
(`web/template.html:148`). Its direct children include `#board-findings`
(`:193`) and each `#view-*` section, so a `.board-main:has(> #view-timeline.is-active)`
selector is structurally valid, exactly like the Activity one.

### Every scroll-related site in `web/board-timeline.js`

`scrollHost` is one variable serving **two distinct roles**: (a) the DOM host that owns the
rows SVG, the listeners, focus and pointer capture; (b) the scroll surface whose
`scrollTop`/`clientHeight` drive virtualization. Re-pointing splits these.

| # | Line | What it does | Role |
|---|------|--------------|------|
| S1 | 1618 | `var scrollHost = document.getElementById("timeline-scroll");` — the only lookup | both |
| S2 | 1621 | early-return guard `if (!summaryNode \|\| !axisHost \|\| !scrollHost \|\| !tableBody)` | host |
| S3 | 1665 | `scrollHost.textContent = "";` clears the rows host before a rebuild | host |
| S4 | 1757, 1759 | `topVisibleRowAnchor()` — finds the first request row whose bottom is below `scrollHost.scrollTop`, and records `offsetPx = scrollHost.scrollTop - displayItem.topPx`. This is "keep the REQ under the pointer when the window changes" | **scroll** |
| S5 | 1788-1797 | `refreshWindowRows()` — grows the rows SVG (`rowsSvg.setAttribute("height", …)`) **first**, then writes `scrollHost.scrollTop = displayItem.topPx + anchor.offsetPx`. The comment records REQ-319's lesson: writing scrollTop before growing the extent clamps it to the old maximum | **scroll** |
| S6 | 1956 | `makeTimelineSvgNode(scrollHost, "svg", {class:"timeline-rows-svg", height: timelineDisplay.height, …})` — the rows SVG is appended to it | host |
| S7 | 2002 | `liveHostWidth()` = `scrollHost.clientWidth \|\| scrollHost.getBoundingClientRect().width \|\| 0`; feeds `plotWidth()` and `plotIsMeasurable()` | host (width) |
| S8 | 2077-2085 (comment), 2116-2118 | `timelineVisibleDisplayRange(timelineDisplay.items, scrollHost.scrollTop, scrollHost.clientHeight)` — the virtualized visible-row math, inside `renderVisibleRows`, which **is** the scroll listener | **scroll** |
| S9 | 2296 | `var focusTarget = rebuiltRow \|\| scrollHost;` — focus falls back to the chart when the rebuilt row is gone | host (focus) |
| S10 | 2726 | `plotResizeObserver.observe(scrollHost)` — width-change delivery | host (width) |
| S11 | 2777 | `addTimelineListener(scrollHost, "scroll", renderVisibleRows);` | **scroll** |
| S12 | 2803-2806 | `scrollTimelineRowIntoView()` — `if (rowTopPx < scrollHost.scrollTop) scrollHost.scrollTop = rowTopPx; else if (rowBottomPx > scrollHost.scrollTop + scrollHost.clientHeight) scrollHost.scrollTop = rowBottomPx - scrollHost.clientHeight;` | **scroll** |
| S13 | 2861 | `focusin` listener (roving tabindex follows focus) | host |
| S14 | 2882 | `keydown` listener — Enter/Space open, ↑/↓ rove rows, ←/→ pan, +/− zoom | host |
| S15 | 2932 | `mousemove` listener — hover readout | host |
| S16 | 2952 | `mouseleave` listener — clears the readout | host |
| S17 | 2969 | `var bounds = scrollHost.getBoundingClientRect();` inside the wheel handler; only `bounds.left` is used | host (x only) |
| S18 | 2985, 2986 | `wheel` listener on `scrollHost` **and** on `axisHost`, both `{passive: false}` | host |
| S19 | 3007 | `pointerdown` — arms the pan | host |
| S20 | 3038-3042 | `setPointerCapture` on engage | host |
| S21 | 3050, 3062 | `pointermove` — measures `clientX` deltas, adds `is-panning` class | host (x only) |
| S22 | 3077-3094 | `pointerup`/`pointercancel`/`pointerleave`/`lostpointercapture` — `releasePointerCapture`, drops `is-panning` | host |
| S23 | 3186 | jump-to-open-work: `scrollHost.scrollTop = openDisplayItem ? openDisplayItem.topPx : 0;` inside `#timeline-zoom-now`'s handler | **scroll** |

There is **no `scrollIntoView` call anywhere** in `board-timeline.js` — S12 is a manual
`scrollTop` write, which is why it can be re-pointed at all.

Dead code note: `timelineVisibleRowRange` (`web/board-timeline.js:1056-1060`) has **no
callers** in `web/`; only test prose mentions it. If the visible-row math is touched, delete
it rather than porting it.

---

## 3. What re-pointing the scroll host to `.board-main` requires

Define once per render, not per event:

- `boardMain = document.getElementById("board-main")`
- `rowsOffsetPx` = distance from `.board-main`'s content-box top to the rows SVG's top, i.e.
  `scrollHost.getBoundingClientRect().top - boardMain.getBoundingClientRect().top + boardMain.scrollTop`
- `stickyAxisHeightPx` ≈ `TIMELINE_AXIS_HEIGHT` (26, `web/board-timeline.js:31`) once the
  axis is sticky.

Then `rowsScrollTop = Math.max(0, boardMain.scrollTop - rowsOffsetPx)` is the number every
"scroll" row in the table above wants.

**Just swap the element reference (no offset math):**
- S2 guard — add `boardMain` to it. `renderTimelineView` must bail if the board host is missing.
- S11 the `scroll` listener moves from `#timeline-scroll` to `#board-main`. **Do not leave it
  on `#timeline-scroll`** — that element stops emitting scroll events once it stops scrolling,
  and every virtualized render silently stops updating.

**Need the offset subtracted or added:**
- S4 `topVisibleRowAnchor()` — read `rowsScrollTop`, not `boardMain.scrollTop`. Naively swapping
  it makes every anchor offset wrong by `rowsOffsetPx` (~250-400 px on this board), so the
  keep-the-REQ-under-the-pointer behaviour drops the reader several rows off.
- S5 `refreshWindowRows()` write — `boardMain.scrollTop = rowsOffsetPx + displayItem.topPx + anchor.offsetPx`.
  The "grow the extent first" ordering **still matters**: `boardMain.scrollHeight` depends on the
  rows SVG height, so writing before the resize still clamps. Keep the comment; it stays true.
- S8 visible-range — `timelineVisibleDisplayRange(items, rowsScrollTop, boardMain.clientHeight)`.
  Not subtracting the axis height from the viewport is **safe**: it renders a superset. Not
  subtracting `rowsOffsetPx` from the scroll position is **not** safe — the range walks off the
  bottom and the top rows go blank.
- S12 `scrollTimelineRowIntoView` — must convert both directions **and** account for the sticky
  axis, or a row scrolled "to the top" lands underneath the axis strip and is invisible while
  focused. Target `boardMain.scrollTop = rowsOffsetPx + rowTopPx - stickyAxisHeightPx` for the
  upward case.
- S23 jump-to-open-work — same: `boardMain.scrollTop = rowsOffsetPx + topPx - stickyAxisHeightPx`.
  The `: 0` fallback becomes "scroll the chart's top into view", not "scroll the board to 0".

**Subtly wrong if swapped naively (the real traps):**

- **R1 — `overflow-x: hidden` keeps it a scroll container.** Per CSS Overflow, when one axis
  is `visible` and the other is not, the `visible` computes to `auto`. So
  `overflow-y: visible; overflow-x: hidden` leaves `.timeline-scroll` **still scrolling**, the
  probe still measures two surfaces, and it will look like the fix did nothing. Both axes must
  become `visible` (or the whole `overflow` shorthand dropped), and `max-height: 58vh` removed.
- **R2 — `.board-main.scrollTop` is shared across every view and nothing resets it.**
  `applyView` (`web/board-controls.js:19-31`) only toggles `is-active`/`hidden`; grep confirms
  no `scrollTop` write anywhere outside `board-detail.js:657-658` (the drawer) and
  `board-timeline.js`. Arriving on Timeline from a Kanban board scrolled 800 px down leaves
  the virtualization computing rows at a position that belongs to another view. On Activity
  (REQ-585) this was harmless — no virtualization. Here it needs a decision: reset
  `boardMain.scrollTop = 0` on entering the view, or clamp. Flag it as a D-xx.
- **R3 — the verify-findings strip is visible on Timeline, unlike Activity.**
  `board-controls.js:59-60`: `findingsStrip.hidden = viewState.view === "activity" || !findingsStripHasContent;`.
  `#board-findings` is a sibling of `#view-timeline` inside `#board-main`, and `.board-anomalies`
  (`web/board.css:563-570`) has `margin: 0 0 20px` — **no top margin**. So zeroing
  `.board-main`'s top padding on this view pulls the strip flush against the top bar. The REQ-585
  recipe (move the padding onto the view's first element) does not reach a sibling that lives
  outside the view. Needs a second scoped rule, e.g. giving `#board-findings` a top margin
  under the same `:has()` scope, or picking a different first element. The REQ-587 screenshot
  shows the strip present with five worktree cards, so this is visible on the real board.
- **R4 — the sticky axis needs an opaque background, and it must be `--bg-base`, not a
  `--surface-*` token.** `_dev/primes/prime-kanban-board.md` states outright: "The surface
  behind this board's SVG is `<body>`, not any `--surface-*` token. `#durations-chart`,
  `#timeline-scroll` and their `<svg>` elements are all transparent". Copying the Activity
  header's `background: var(--surface-1)` paints a visibly tinted strip over a `--bg-base`
  chart. Also needs a `z-index` above the rows (Activity used `z-index: 1`).
- **R5 — `.timeline-scroll:focus-visible` ring becomes useless.** With the box full height
  (thousands of px), the inward 2 px ring bounds the whole element; on screen the reader sees
  only two vertical lines and no top or bottom edge. The comment at `web/board.css:2896-2897`
  ("`.timeline-scroll` is an overflow container flush under the axis") stops being true and the
  keyboard affordance "Tab to the chart" loses its visible landing signal. Needs a rethink —
  ring the axis, ring a wrapper, or accept it and say so.
- **R6 — `liveHostWidth()` changes value.** `#timeline-scroll`'s `clientWidth` currently excludes
  its own vertical scrollbar (~15 px). After the change that scrollbar is on `.board-main`
  instead, so the measured plot gains those pixels. Harmless in itself (the ResizeObserver at
  S10 delivers the change and re-renders), but it means every recorded plot-width number in
  comments and probes shifts.
- **R7 — the `<details>` table below the chart is a *third* scroll box.**
  `.timeline-table-scroll` (`web/board.css:3101-3105`) is `max-height: 360px; overflow: auto`.
  It only scrolls when the reader opens the `<details>`, and it is not in the REQ's RED
  expression — but a probe that asserts "exactly one element on the page scrolls" will trip on
  it if the details is opened. Scope the probe to the two elements the REQ names, or state the
  exemption.

**Forced-reflow risk during scroll.** `renderVisibleRows` is the scroll listener (S11) and
already forces one synchronous layout per frame — `scrollHost.clientHeight` at S8 plus
`timelineMeasureLabelAdvance` (documented at `web/board-timeline.js:2063-2070`, measured
0.036-0.046 ms on Chromium 1194). Computing `rowsOffsetPx` with two
`getBoundingClientRect()` calls **inside** the handler adds two more forced layouts per
scroll event, on every frame of a drag. Measure it once per render (in `renderAll`, beside
`invalidatePlotWidth`) and cache it; it can only change when content above the chart changes,
which is always a re-render (findings strip shown/hidden, toolbar wrap, summary line length),
never a scroll.

---

## 4. Wheel-zoom, drag-pan, keyboard

**Wheel-zoom** (`handleTimelineWheel`, `web/board-timeline.js:2964-2984`):
returns immediately unless `ctrlKey || metaKey`; then `wheelEvent.preventDefault()`. Attached
twice, both `{passive: false}` (required for `preventDefault`):
`addTimelineListener(scrollHost, "wheel", handleTimelineWheel, { passive: false });` (`:2985`)
and the same handler on `axisHost` (`:2986`). The anchor is
`(wheelEvent.clientX - scrollHost.getBoundingClientRect().left - TIMELINE_LABEL_WIDTH) / plotWidth()`
— **x only**. The comment at `:2957-2963` says the rect is deliberately taken from the scroll
host either way, so the axis and the plot share one x scale.

*What breaks:* nothing, if the listeners stay on `#timeline-scroll`. `bounds.left` is unaffected
by the element becoming tall. A plain (unmodified) wheel currently scrolls `.timeline-scroll`
and afterwards scrolls `.board-main` — same visible result, one surface instead of two, which
is the point of the REQ. **Do not move the wheel listener to `.board-main`**: over the axis it
would still need its own binding, and a listener on the board would swallow ctrl-wheel over
every other view.

**Drag-pan** (`web/board-timeline.js:3006-3106`): `pointerdown` arms
(`{pointerX, pointerId, windowStartMs, engaged:false}`), `pointermove` engages past
`TIMELINE_PAN_THRESHOLD_PX` and only then calls `setPointerCapture`, `pointerup`/`cancel`/
`leave`/`lostpointercapture` tear down. It **never calls `preventDefault`** — the capture is
what makes a drag not a click (REQ-336's lesson: capturing on pointerdown retargets the
synthesized click and broke every click in the chart). All motion is `clientX` deltas divided
by `plotWidth()`.

*What breaks:* nothing from re-pointing the scroll host, as long as the listeners and the
capture stay on `#timeline-scroll`. One thing to watch: with no `preventDefault` on
`pointermove`, a vertical drag inside the chart now scrolls `.board-main` where it previously
scrolled the inner box — visually identical. `cursor: grab` / `.is-panning { cursor: grabbing }`
still apply over the element's (now full-height) box.

**Keyboard** — the hint paragraph is `web/template.html:452` and the panel's accessible name
is `:371`; both are pinned by `TestTimelinePanelStatesItsKeyboardInteraction`
(`generate_test.go:2737-2767`), which also asserts `tabindex="0"` on the `id="timeline-scroll"`
opening tag. The keydown handler is `web/board-timeline.js:2882-2921`:

- Enter/Space → `openDetail`, `preventDefault`.
- ↑/↓ → only when the event target is inside a `[data-row-index]`; `preventDefault`, then
  `moveTimelineRowFocus` (which calls `scrollTimelineRowIntoView` = S12). With focus on the
  chart itself they deliberately fall through to the browser's own scroll — **that browser
  scroll silently changes surface** from the inner box to the board. Behaviourally the same;
  worth one line in the hand-back.
- ←/→ and +/− → `timelineKeyboardWindow`, `preventDefault`, `renderAll()`.

*What breaks:* "Tab to the chart" still works (the `tabindex="0"` element is unchanged), but
its focus ring becomes near-invisible — see R5. Arrow-key row movement depends entirely on S12
being converted correctly, including the sticky-axis subtraction; get that wrong and ↓ appears
to do nothing because the focused row sits under the axis.

---

## 5. What sits below the chart

All four are **siblings of `.timeline-chart`, outside `#timeline-scroll`**, inside
`#view-timeline` (`web/template.html:450-471`):

| Element | Markup | CSS |
|---|---|---|
| forecast line | `<p class="timeline-forecast" id="timeline-forecast">` `:450` | `web/board.css:3147-3153` (`padding: 10px 24px 0`) |
| excluded list | `<div class="timeline-excluded" id="timeline-excluded">` `:451` | `web/board.css:3160+` |
| hint paragraph | `<p class="timeline-hint">…` `:452` | `web/board.css:3077-3084` |
| readout (aria-live) | `<p class="timeline-readout" id="timeline-readout">` `:453` | `web/board.css:3086-3089` |
| "The REQs in this window, as a table" | `<details class="timeline-table"><summary>…` `:454-471` | `web/board.css:3091-3105`; the inner `.timeline-table-scroll` is its own 360 px scroll box — see R7 |

This is what makes the REQ's second shape (M2, chart fills the viewport) expensive: five
elements would need a new home. The recommended shape leaves them exactly where they are, as
ordinary content that scrolls with the board.

---

## 6. Test surface

### How a browser probe is structured

Harness: `skills/do-work-board/tools/queue-kanban/browser_probe_test.go`.

- `lookupBrowserForBehaviorProbe(t)` (`:69-91`) — **first line of every probe**. It skips
  unless `QUEUE_KANBAN_BROWSER_PROBES=on` ("browser behavior probes are heavy-only"), then
  consults `QUEUE_KANBAN_BROWSER`, then a well-known-name list
  (`google-chrome`, `google-chrome-stable`, `chromium`, `chromium-browser`, `chrome`), then
  `t.Skipf`. `QUEUE_KANBAN_BROWSER` names a browser binary explicitly (a path or a PATH name);
  if it names nothing runnable the probe **fails** rather than falling back — the caller asked
  for a specific engine. Constants at `:36-48`.
- `startTrustedInputBrowserSession(t, name, siteDir, pageHTML, flags...)` (`:224-296`) — writes
  `probe.html` into the site directory, launches the binary with
  `--headless --disable-gpu --no-sandbox --disable-dev-shm-usage --user-data-dir=… --remote-debugging-pipe`,
  speaks CDP over fd 3/4, attaches to the page by URL. Callers add `--window-size=1600,900`.
- Two styles exist. **Session-driven** (Go drives step by step: `waitForPageCondition`,
  `evaluateInPage`, `decodeResult`) and **script-injected** (a `<script>` spliced before
  `</body>` writes JSON into `#queue-kanban-probe-result`, read back by
  `runBrowserBehaviorProbeInDirectory`). Both count toward the strict lane's zero-probe guard
  (`strictBrowserBehaviorDiagnostic`, `:36`).
- `generateLiveSiteInDir(t)` builds the real board; `writeVerifyFixture(t, files)`
  (`verify_test.go:23`) builds a synthetic repo root.

### REQ-585 did add such a probe

`skills/do-work-board/tools/queue-kanban/activity_scroll_browser_probe_test.go`, one function:
`TestBrowserBehaviorActivityViewHasOneScrollSurface`. **This is the file to copy.** Its shape:

1. Two named constants with reasons: `activityScrollProbeRequestCount = 40` (→ 120 rows, enough
   to overflow both boxes at 900 px) and `activityScrollProbeBoardScrollPixels = 700`.
2. A synthetic fixture (`activityScrollProbeFixtureRequest`) rather than the live queue —
   decision D-06: "the live queue's last 24 hours might not overflow anything on a quiet day
   and the probe would pass while measuring nothing".
3. A single measurement struct carrying **both boxes' heights**, not a pair of booleans, so a
   failure message says how far off the layout was.
4. **REQ-291's guard, applied three times before any comparison**: assert the expected row
   count was rendered, assert every measured box has positive height, assert the board actually
   scrolls (otherwise "the fixture is too short to tell one scroll surface from two").
5. Two `requestAnimationFrame`s after writing `scrollTop` — "not a sleep": the compositor has
   not repositioned the sticky element in the same task.
6. Scoping asserted **in both directions**: `boardMainPaddingTop == "0px"` on Activity, then
   click back to the Kanban view and assert it is *not* `0px`.
7. A `t.Logf` line with every number, so a green run still records the measurement.

### Two more probe shapes worth copying for the Timeline half

- `TestBrowserBehaviorTimelineRowListIsOneTabStop` (`timeline_browser_probe_test.go:3330`) —
  script-injected, drives the real generated board, clicks `[data-view-target="timeline"]`,
  waits, presses synthetic keys, reports a state object per step. The right shape for "the
  virtualized range follows the scroll".
- `TestBrowserBehaviorTimelineBarsSurviveTheDetailDrawerOpening` (`:1319`) — the probe that
  pins the ResizeObserver/width-change invariant, i.e. the one most likely to break if the
  host's width measurement changes (R6).
- `TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving` (`:627`) and
  `TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` (`:3124`) — the drag-pan
  regression pins. Run both; they are the ones that would catch a broken pan.

### Node-lane blast radius — the least obvious cost

`renderTimelineView` is driven directly by the Node behavior lane through a stub DOM whose
`document.getElementById` is a **fixed id map**:

`generate_test.go:3323-3364`, `timelineRenderDomStubPreamble`:

```js
var timelineStubHosts = {};
[
  "timeline-summary", "timeline-axis", "timeline-scroll", "timeline-readout",
  "timeline-table-body", "timeline-forecast", "timeline-excluded", "timeline-period-state"
].forEach(function (hostId) { timelineStubHosts[hostId] = makeStubNode("div"); });
var document = {
  getElementById: function (nodeId) { return timelineStubHosts[nodeId] || null; },
  …
  querySelector: function () { return null; }
};
```

If `board-timeline.js` starts reading `board-main`, **every one of these has to gain it** or
`renderTimelineView` hits its guard (S2) and returns, silently turning several probes into
assertions about nothing:

- `generate_test.go:3344` (the preamble), `:3394` (`walkRowGroups(timelineStubHosts["timeline-scroll"])`)
- `javascript_behavior_a_test.go:2002, 2148, 2156` (and the local `getElementById` override at `:2132-2135`)
- `javascript_behavior_b_test.go:1524`
- `javascript_behavior_c_test.go:2191`

Note `querySelector` returns `null` in the stub, so a `.board-main` **class** lookup would
break the lane outright; `document.getElementById("board-main")` plus a stub entry is the
cheap path. The stub node already carries `clientHeight: 400` and `scrollTop: 0`, so one added
entry covers it. Also `generate_test.go:2759-2767` asserts `tabindex="0"` on the
`id="timeline-scroll"` tag — safe as long as the template attribute stays.

No Go test asserts `58vh` or `max-height` anywhere (grep: only a comment in the Activity probe),
so the CSS change itself has no lock-in to update.

---

## 7. Primes and lessons

### `_dev/primes/prime-kanban-board.md` — the entries that bear on this REQ

- **"The surface behind this board's SVG is `<body>`, not any `--surface-*` token."** Names
  `#timeline-scroll` explicitly: it and its `<svg>` are transparent. REQ-321 found it under the
  timeline's bars and REQ-346 found it again under the Durations lane — "the second builder paid
  for it a second time". Directly governs the sticky axis's background (R4). Also warns that
  `--ink-faint` vs `--ink-soft` measure 1.29:1 light / 1.82:1 dark **as fills**.
- **"A chart's correctness is partly a claim about pixels — generate a board and look at it."**
  REQ-226/231/237/240 each shipped a defect every assertion passed over. Measure
  `getBoundingClientRect()` intersections in the live DOM for overlap questions.
- **"Render evidence must name the page it measured, in the same call that measures it."**
  Return `location.href` alongside every measurement — a sibling session can navigate a shared
  browser between your navigate and your evaluate. The existing timeline probes all carry
  `href: location.href` in their state objects.
- **"A measured face is per-browser, and a constant that does not name its build cannot be
  argued with."** Record the browser and build beside every measured number.
- **Browser support:** the strict lane targets current stable Chromium; Chrome 141 is not a
  compatibility target (REQ-375).
- **"`web/` is embedded, not read at runtime"** (shipped prime `prime-do-kanban.md` § Traps) —
  CSS/JS edits need a rebuild; `serve` re-serves the embedded assets, never files on disk.
  Relevant when eyeballing the fix on port 8090.
- Its only **Traps** entry is `[family: exact-basename-authority]`, unrelated to this work.

`skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` has **nothing** on scroll,
sticky, or virtualization (the one "sticky" hit is REQ-447 on file modes).

### `_dev/primes/lessons-kanban-board.md` — the family lesson, verbatim

> - [family: sticky-pins-to-content-box] [REQ-585](../../do-work/archive/UR-120/REQ-585-give-the-activity-view-one-scroll-surface.md#lessons-learned): a sticky header inside a padded scroll container pins to the container's CONTENT box, so the padding is a band the rows scroll through above it; when a view gives up its inner scroll box for the board's, move the board's top padding onto the view's first element, scoped to that view (`.board-main:has(> #view-activity.is-active)`), and measure it in an engine rather than reasoning from the specification.

Other entries in that file that bear directly on this change:

- **REQ-319** — "writing scrollTop before growing the SVG extent clamps it silently to the old
  maximum". This is S5; the conversion must preserve the ordering.
- **REQ-338** — a roving tabindex over a **virtualized** list needs a clamp into the rendered
  range; Tab itself cannot be tested with synthetic events (trusted-input default action).
- **REQ-336 / REQ-337** — pointer capture retargets the synthesized click; a text assertion
  about `setPointerCapture` must include the opening paren or it matches the `typeof` guard.
- **REQ-291** — `getBBox()` returns zeros for an unrendered element, so a browser probe's
  default failure is a successful-looking measurement of nothing. Write the result node last,
  assert positive-and-finite. The Activity probe applies this three times.
- **REQ-322** — a constant a decision turns on must be **read** by the test, never restated
  beside it. Relevant if the probe hard-codes 58vh, the axis height, or a row height.
- **REQ-285** — for a rendering change the screenshot **is** the test.

---

## Summary of the risks worth pricing before planning

- **R1** `overflow-x: hidden` alone keeps it a scroll container — the fix can look applied and
  measure RED.
- **R2** `.board-main.scrollTop` is shared across views and nothing resets it; virtualization
  makes that visible here in a way it was not on Activity.
- **R3** the verify-findings strip is visible on Timeline and sits **outside** the view panel,
  so REQ-585's padding-move recipe does not reach it.
- **R4** the sticky axis needs `--bg-base`, not `--surface-1`; the prime says so explicitly.
- **R5** `.timeline-scroll:focus-visible` loses its meaning on a full-height element, and
  "Tab to the chart" is in the shipped hint text.
- **R6/R7** plot width shifts by a scrollbar's worth; the `<details>` table is a third scroll box.
- **Node-lane cost**: six stub sites across four Go test files need a `board-main` entry the
  moment the renderer looks one up.
