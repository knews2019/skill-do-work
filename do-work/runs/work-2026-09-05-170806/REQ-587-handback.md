# Hand-back — REQ-587 (give the Timeline view one scroll surface)

**Branch:** `worktree-agent-REQ-587-give-the-timeline-view-one-scroll-surface`
**Branch head:** `564c81a3` — "REQ-587: give the Timeline view one scroll surface"
**Worktree:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-587-give-the-timeline-view-one-scroll-surface`
**Cut from:** `d6b6adb1`. Nothing under `do-work/` was read, written, staged or committed. This file is the only main-tree write and is not staged.

Status: done. The plan's T1 through T6 all landed. Two defects that only a screenshot and a real engine could find were fixed on top of the plan (D-14, D-15). One file outside `## Scope` had to change; it is named in full under **Integration seams** and in the commit message.

## File manifest

All paths relative to `skills/do-work-board/tools/queue-kanban/`.

| Verb | File | What changed |
|---|---|---|
| modify | `web/board.css` | The block comment that called `.timeline-scroll` the scroll container; `.timeline-scroll` stripped to `position: relative` + `cursor: grab`; the first `.timeline-axis` rule (sticky, `--bg-base`, `border-bottom`); the scoped padding move (three rules); the focus ring moved to the axis plus the UA-ring suppression |
| modify | `web/board-timeline.js` | One board host lookup and guard; a memoized geometry pair with its invalidation; five scroll-position conversions; the scroll listener moved to `#board-main`; the dead `timelineVisibleRowRange` deleted |
| modify | `web/board-controls.js` | The guarded board scroll reset on entering the Timeline (D-07) |
| modify | `generate_test.go` | `"board-main"` in the shared Node-stub id list, with the reason; one stale function name in a historical comment |
| modify | `javascript_behavior_a_test.go` | `"board-main"` in two local stub id lists |
| modify | `javascript_behavior_b_test.go` | `"board-main"` in one local stub id list |
| modify | `javascript_behavior_c_test.go` | `"board-main"` in one local stub id list |
| **new** | `timeline_scroll_browser_probe_test.go` | `TestBrowserBehaviorTimelineViewHasOneScrollSurface` — four guards, eight assertions, 201-REQ fixture |
| modify | `timeline_browser_probe_test.go` | **Outside the declared Scope.** Three synchronous virtualizer settles after `scrollIntoView`; one vacuity guard re-pointed at the board's client height. See **Integration seams**. |

`git diff --stat` on the commit: 9 files, 320 insertions, 51 deletions. I checked the file list explicitly against `## Scope`: eight of the nine are the declared write set, exactly; the ninth is `timeline_browser_probe_test.go`. No path under `do-work/`, no build output (`queue-kanban` is gitignored in place), nothing else.

## Measured Evidence

Every number below was measured in **Chrome 152.0.0.0** (`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) … Chrome/152.0.0.0 Safari/537.36`; the headless probe lane reports the same build as `HeadlessChrome/152.0.0.0`), at a 1600×900 window, on the board served by a **rebuilt** binary — `web/` is embedded, so nothing below would have been visible without the rebuild.

### The REQ's own RED expression, on the served board

The two builds were served side by side from the same `do-work/` tree, on ports 8098 (pre-change, built from `d6b6adb1`) and 8097 (this branch), and the expression was run verbatim in the console on the Timeline with the All days window.

| | pre-REQ-587 (`d6b6adb1`) | this branch (`564c81a3`) |
|---|---|---|
| **RED expression** | `[true, true]` | **`[true, false]`** |
| `#board-main` client / scroll | 774 / 1721 | 774 / 15518 |
| `#timeline-scroll` client / scroll | 521 / 14318 | **14318 / 14318** |
| `#timeline-scroll` computed `overflow-y` | `auto` | **`visible`** |
| `#timeline-scroll` computed `max-height` | `522px` | `none` |
| `#board-main` computed `padding-top` | `24px` | `0px` |
| `#timeline-axis` computed `position` | `static` | `sticky` |

The `overflow-y` row is the one the height comparison cannot see: with `max-height` gone but `overflow-x: hidden` left in place, the chart would report `14318 / 14318` and still be a scroll container. Both declarations came off, and the probe asserts the computed value, not only the heights.

### The axis, pinned

Scrolled 1500px of rows into the chart, light theme: axis top offset from the board's top inner edge **0.0px**, rows painted above the axis **0**, axis height 27px. Rows offset (board scroll origin to the top of the rows box) **721.6px** in the probe fixture and **913.9px** on this repository's own board with a seven-finding verify strip on screen.

### The virtualized range follows the board

Before the scroll, 39 rendered row ids; after 1500px, 44; **intersection empty**; 34 of the rendered rows intersect the board's visible rect afterwards. Measured by the new probe, asserted in both directions.

### The row under the pointer, across a window change

The GREEN clause with no existing equivalent. Narrow window (Last 90 days), scrolled 500px of rows in, anchor row `REQ-7003` recorded at `-3.9px` from the board's inner top edge. Widened to All days — the display list went **62 → 201 rows**, so rows really were inserted above the anchor, and the board's own scroll moved **1275.0 → 1960.0** to follow it. The anchor came back at **-4.4px**: **0.5px of drift**. Pinned in the probe, not checked by hand.

### The moved padding

Timeline: `#board-main` `padding-top` = `0px`, the first unhidden child carries `margin-top: 24px` and starts 24.0px below the board's top inner edge, `#view-timeline` carries `margin-top: 0px` because a strip precedes it, and the view's own 16px inset survives (`.timeline-heading` 16.0px below the view panel's top edge). Switching to the Kanban view restores `padding-top: 24px`. Verified at **320 / 768 / 1280 / 1600** px: `padding-top: 0px` and `margin-top: 24px` on the first visible child at every one.

In the probe fixture the first visible child is the **data-warnings banner** — an element with no `id` that is not in `template.html` at all, inserted by `board-cards.js` as the board's first child. That is the case a rule naming the strips by id would have missed, and it is why the rule is keyed on "the first child that is not hidden". On this repository's own board the first visible child is `#board-findings`.

### Both themes

| | light | dark |
|---|---|---|
| `getComputedStyle(document.body).backgroundColor` | `rgb(245, 247, 250)` | `rgb(12, 14, 18)` |
| `getComputedStyle(#timeline-axis).backgroundColor` | `rgb(245, 247, 250)` | `rgb(12, 14, 18)` |
| RED expression | `[true, false]` | `[true, false]` |
| axis top offset | 0 | 0 |

Identical by construction (`--bg-base` is the `<body>` token) and confirmed by measurement in both. No tinted band over the untinted chart — the defect REQ-321 and REQ-346 each paid for.

### Interaction, by hand in a real engine

- **Tab to the chart:** two Tab presses from the toolbar land on `#timeline-scroll`, `:focus-visible` true, ring drawn on the pinned axis (`rgb(58, 107, 196) solid 2px`, offset `-2px`), axis at offset 0 so the ring is on screen at any scroll position.
- **Arrows on the chart:** ArrowDown scrolls the board (300 → 340). At the board's maximum scroll it correctly does nothing.
- **Arrows in the row list:** Tab again enters the rows (landed on `g.timeline-row`, `REQ-588`). ArrowUp roved to `REQ-456` and the board scrolled **2914 → 2400** to follow it; the focused row landed at y=468 against an axis bottom of 153, so it cleared the pinned strip. This is the S12 conversion, which has no probe assertion.
- **Ctrl+wheel zoom:** the window readout went `2026-05-27 17:03 → 2026-09-08 00:04` to `2026-06-13 12:01 → 2026-08-17 01:24`, and `preventDefault` was called, so the page does not scroll under it.
- **`+` key zoom:** `2026-06-13 12:01 → 2026-08-17 01:24` to `2026-06-25 14:31 → 2026-08-04 22:53`.
- **Drag-pan and pointer capture:** covered by the two existing probes, both green.

### Screenshots

Gitignored (`.gitignore:1`), in the main tree, absolute paths:

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.playwright-mcp/REQ-587-timeline-scroll-0-light.png` — arrival, board at scroll 0. The verify-findings strip starts 24px below the top bar; the chart is below the fold, exactly as it was before this change (the content above it is unchanged).
- `…/REQ-587-timeline-chart-top-light.png` — mid-scroll, chart's top at the board's top edge. One scrollbar. The axis pinned under the top bar with rows passing under it.
- `…/REQ-587-timeline-mid-scroll-dark.png` — the same state under `preferredColorScheme=0`. No tinted band.
- `…/REQ-587-timeline-focus-ring-DEFECT-ua-ring-on-full-height-chart.png` — **the defect D-14 fixed**, kept as evidence: Chrome's own focus ring drawn around the full-height chart, a blue box from the axis down past the last row.
- `…/REQ-587-timeline-focus-ring-FIXED-axis-only.png` — after the fix: a clean 2px ring on the pinned axis alone.

### Test runs

- `go -C skills/do-work-board/tools/queue-kanban test -count=1 ./...` — **exit 0**.
  - This worktree's own clean baseline at `d6b6adb1`, measured before any edit: **54.3s**.
  - After the change, two runs: **62.4s** and **48.9s**.
  - The brief's 42.4s reference was taken when the machine was quieter; my own baseline in this tree, taken minutes before the first edit, is the honest comparison, and 48.9s is below it. A sibling session was running browser work throughout, which is the spread. There is no measurable regression: the new probe is gated behind `QUEUE_KANBAN_BROWSER_PROBES`, so it contributes 0s to this lane.
  - The 30s per-file budget does not decompose here — `queue-kanban` is a single Go package and the number above is the whole package. The new probe costs **1.6s** in the browser lane.
- `go vet ./...` — clean. `gofmt -l .` — empty.
- **Browser lane, full:** `QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" QUEUE_KANBAN_BROWSER_PROBES=on go test -count=1 -run 'TestBrowserBehavior' .` — **exit 0, 33 probes, every one PASS, no SKIP line in the output.** I grepped for skips rather than trusting the exit code, because a skipped browser probe is not a pass.
  - The five named regression probes are among them: `…TimelineBarsSurviveTheDetailDrawerOpening`, `…TimelineRowListIsOneTabStop`, `…TimelinePressBecomesAPanOnlyAfterMoving`, `…TimelinePointerCaptureWaitsForThePanEngage`, `…ActivityViewHasOneScrollSurface`.
  - So is `TestBrowserBehaviorTimelineListsRowsBeneathUserRequestHeaders`, which the plan did not name and which this change broke; see D-16.
- **Node/string lane:** the four `javascript_behavior_*` timeline probes and `TestTimelinePanelStatesItsKeyboardInteraction` all pass inside the full run above. Every row-count and row-id expectation in that lane is **byte-identical** — no expectation was edited, which is what the plan predicted and what says the coordinate conversions are right.
- `grep -rn '"timeline-period-state"' *_test.go` returns the same five list literals, each now carrying `"board-main"`.

## P-A-U

### [PLAN]

Followed `do-work/runs/work-2026-09-05-170806/REQ-587-plan.md`, tasks T1–T6, with D-01 fixed by the REQ and D-02 through D-13 taken as written. The plan was verified against source and held everywhere I checked it, including all five Node-stub sites, the `renderAll` invalidation point, and the claim that no existing probe writes `scrollTop` on `#timeline-scroll`. Two of its predictions turned out to be incomplete once measured — see D-14 and D-15 — and one file it did not consider had to change (D-16).

The orchestrator's override was taken: the anchor requirement is pinned in the probe, with its own display-list guard, not checked by hand. It is assertion A8 in `assertTimelineAnchorSurvivesAWindowChange`.

### [APPLY]

T1 `web/board.css` → T2 `web/board-timeline.js` → T3 `web/board-controls.js` → T4 the four Node-stub files → T5 the new probe → T6 rebuild, serve, measure, screenshot. Scope stayed inside the declared list until the browser lane forced the one file named below.

### [UNIFY]

`git diff --stat` reviewed file by file; `go vet ./...` and `gofmt -l .` clean; no `console.log`, `debugger`, `TODO`, `FIXME` or `XXX` anywhere in the diff.

| File | What I checked |
|---|---|
| `web/board.css` | The old block comment no longer claims `.timeline-scroll` is the scroll container. `.timeline-scroll` has exactly two declarations and neither is an `overflow`. `.timeline-axis` paints `--bg-base` and carries the 1px rule the rows box used to own, so total height is unchanged. The three padding rules cover every arrangement of the four possible preceding siblings and leave exactly one carrying the 24px. Both focus rules read as one block. No other view's rules touched — the Activity block is byte-identical. |
| `web/board-timeline.js` | All five scroll-position sites converted; `grep -n 'scrollHost\.scrollTop\|scrollHost\.clientHeight'` returns nothing but the one comment I updated. The seventeen host/width/focus/hover/keyboard/wheel/pointer sites still name `scrollHost` — I read each of the 24 remaining `scrollHost` lines. Wheel listeners still on the chart and the axis; the whole pan block and `setPointerCapture`/`releasePointerCapture` still on the chart. The memo pair sits beside `measuredPlotWidthPx` and is invalidated on the line after `invalidatePlotWidth()`, before the `plotIsMeasurable()` guard, so the refusing path still measures nothing. The deleted `timelineVisibleRowRange` had no callers; the one historical comment naming it was updated in the same commit. |
| `web/board-controls.js` | Three lines, inside `applyView`, after the panel toggle loop and before the `renderedOnce.timeline` render, so the first anchor read sees 0. Unguarded `getElementById` deref matches the eight lines above it, which do the same. |
| `generate_test.go` | The shared stub id list carries `"board-main"` at the right indent, with a comment stating why and what it means if a row count moves. The `timelineVisibleRowRange` mention in the resize-guard comment became "the visible-row range" so the narrative survives without naming a function that does not exist. No assertion changed. |
| `javascript_behavior_{a,b,c}_test.go` | Four id lists, indentation matching each site, nothing else. No assertion changed; every row count and row id in that lane came back identical. |
| `timeline_scroll_browser_probe_test.go` | Reads as the Activity probe's seven shape rules: named constants each carrying the arithmetic that justifies them, a synthetic fixture, one measurement struct of real heights, guards before comparisons, two animation frames after every `scrollTop` write, scoping asserted in both directions, and a `t.Logf` carrying every number so a green run still records the measurement. Scope stated in the header comment (D-12). |
| `timeline_browser_probe_test.go` | Four hunks, all mechanical, none touching an assertion's meaning. The three settle lines each carry the same six-line reason; the guard change is a one-word re-point with its reason beside it. |

## Decisions

Continuing from D-13.

- **D-14 — the chart's own focus ring must be explicitly suppressed, not merely relocated. DECIDE & STATE.** D-06 moved the ring to the axis by deleting `.timeline-scroll:focus-visible`. That put Chrome's *default* focus ring back on the same full-height box, in exactly the shape D-06 exists to prevent: measured `outline: rgb(0, 95, 204) auto 1px` around a 716px element, drawn as a blue box from the axis down past the last row. Every assertion passed while it was there — `:focus-visible` after a programmatic focus is deliberately not asserted (REQ-233), and the plan's own note says the evidence for D-06 is a screenshot. It was the screenshot that found it, which is REQ-285's lesson arriving on schedule. The fix is `.timeline-scroll:focus-visible { outline: none }` in the same block as the axis rule, with a comment saying the two are one affordance and that the signal is moved rather than removed. Both screenshots are kept as before/after evidence.
- **D-15 — `rowsScrollTop()` stays signed; the floor moves to the one call site that wants it. DECIDE & STATE.** The plan wrote `Math.max(0, boardScrollTop − rowsOffset)`. That clamp makes the anchor's write disagree with its read: `topVisibleRowAnchor` records `(position − rowTop)` and `refreshWindowRows` adds it back, so a swallowed negative comes back as a positive scroll. Measured consequence — pressing any window chip while the board sat above the chart scrolled it down by `rowsOffset` (721px in the fixture, 914px on the real board), slamming the chart's first row against the top bar and hiding the heading, legend and toolbar the reader had just used. The helper is now honest and the visible-range call applies `Math.max(0, …)` itself, which is a superset in every case (a negative position can only widen the range upward) and reproduces what this view drew before the board became its scroll surface, so the chart is already populated when it scrolls into view rather than one scroll event later. Both halves carry the reasoning in comments.
- **D-16 — `timeline_browser_probe_test.go` was updated rather than reported and left red. DECIDE & STATE, with the scope widening flagged.** Four existing browser probes went red on this change, and all four were green at `d6b6adb1` — I confirmed that by stashing the whole change and re-running them, rather than assuming. `crew-members/general.md` § Cross-REQ Test-Break Rules governs: the behaviour change is intentional, so the failing tests are updated to match it and the change is documented. Two distinct causes, both mechanical:
  - Three probes call `scrollIntoView` on the chart and then aim, press and measure **in the same task**. The rows only re-render when `#board-main` emits its scroll event, one frame later, so those probes were measuring the rows the previous board position drew. From a deep starting position those rows sit thousands of pixels below the fold and `elementFromPoint` returns null for every one, which the probes correctly reported as "no press point" for a chart that was plainly on screen. Reproduced by hand before fixing: stale first row at y=2671 with `elementFromPoint` → undefined; after one `dispatchEvent(new Event("scroll"))` on the board, y=195 and a real `rect` hit. The fix is that one line at each of the three sites — it drives the shipped listener with the shipped handler and does the work the engine is about to do anyway.
  - `TestBrowserBehaviorTimelineListsRowsBeneathUserRequestHeaders` guards vacuity with "the SVG is taller than the viewport", reading the viewport as `#timeline-scroll.clientHeight`. That is now the SVG's own height, so the guard compared a number with itself and could never fire again. Re-pointed at `#board-main.clientHeight`, which is what "the viewport" now means.
  I judged this DECIDE rather than ESCALATE because leaving the heavy lane red hands back a change whose own GREEN condition is a browser measurement, and because neither edit changes what any assertion means. The orchestrator should still see it as a scope widening, so it is in the commit message and under **Integration seams**.
- **D-17 — the probe widens the window rather than narrowing it for the anchor assertion. DECIDE & STATE.** The orchestrator asked for a guard that the anchor row is still in the display list, so a legitimately-absent anchor reads as a skipped comparison rather than a false pass. That guard is present and reads the "REQs in this window, as a table" body, which lists the whole display list whatever the virtualizer drew. Going from Last-90-days to All-days means every row in the narrow window is also in the wide one, so the guard cannot fire in practice and the comparison always runs — its `t.Skipf` message says exactly that, so a skip reads as a fixture bug rather than a pass. Narrowing instead would have made a real regression indistinguishable from a legitimately dropped row.
- **D-18 — the anchor probe's fixture cycles ten user requests over 300 days. DECIDE & STATE.** A naive fixture cannot catch the conversion bug this assertion exists for. Rows are ordered newest first, so a wider window only appends older rows *below* a mid-list anchor; its `topPx` never moves, and a "forgot to convert both ways" bug preserves the position by accident. Cycling the user request means every group holds both recent and old members, so widening inserts old members into every group above the anchor. Measured: the anchor's group grew such that the board had to move 685px to hold its screen position.
- **D-19 — the fixture's anomaly REQ is kept even though it also raises a data warning. DECIDE & STATE.** The warning was unplanned, and it makes the fixture strictly better: the warnings banner becomes the board's first visible child, and it is an element with no `id` that `template.html` does not contain. The padding rule is therefore tested against the arrangement a hand-maintained list of strip ids would have got wrong. Recorded in the fixture comment.
- **D-20 — the dead `timelineVisibleRowRange` was deleted (the plan's optional D-13), and the one historical comment naming it was reworded. DECIDE & STATE.** It has no callers; its only other mention is a narrative comment in `generate_test.go` about a past defect. Deleting the function while leaving the comment naming it would strand a grep target on nothing, so the comment now says "the visible-row range". Nothing reddened.
- **D-21 — the probe measures the rows-box offset with the same formula the renderer uses. DECIDE & STATE.** REQ-322 says a constant a decision turns on must be read, never restated. This is a borderline case and I chose to restate it: the number is a fact about the laid-out page rather than a decision, and it is used only to pick a scroll target and to report. No assertion depends on it — the axis offset, the disjoint id sets, the anchor position and the padding offsets are all read back from what the page produced. Said so in a comment at the measurement.

## Discovered Tasks

- On a 1600×900 window the Timeline's chart begins 914px below the board's top edge when the verify-findings strip is showing, so it is entirely below the fold on arrival and the reader must scroll before seeing a bar. This is **unchanged by this REQ** — the same content sat above the same chart before it, and the pre-change build measures the same offset — but it is now more visible, because the chart is the tall thing on the page rather than a 58vh box. Worth a REQ on its own if the maintainer wants the Timeline to open on the chart. → report only, impact-user-visible.
- The board scrolls sideways at 320px and 768px. Measured identically on both builds (`board-main` scrollWidth 324 against clientWidth 305 at 320px), and the overflowing elements are the top bar's `nav.board-controls`, its filter input and its two selects — not the timeline. Pre-existing, unrelated to this change; I checked it only because removing `overflow-x: hidden` could plausibly have caused it, and it did not (the chart's scrollWidth equals its clientWidth at every width tested). → report only, impact-low.
- `TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` and `…PressBecomesAPanOnlyAfterMoving` both aim by scrolling and then measuring in the same task. The settle line added in D-16 fixes them, but the pattern will bite again the next time anything makes rendering depend on an asynchronous event. A shared `settleVirtualizedRows()` helper in that file would be the durable fix; three copies of one line is not it. → report only, impact-low.

## Lesson evidence

Read in full, all present:

- `_dev/primes/prime-kanban-board.md` — declared by the REQ. Its "the surface behind this board's SVG is `<body>`" entry decided D-03's token; its "a measured face is per-browser, and a constant that does not name its build cannot be argued with" is why every number above names Chrome 152.0.0.0; its "a chart's correctness is partly a claim about pixels — generate a board and look at it" is what produced the screenshot that found D-14.
- `_dev/primes/lessons-kanban-board.md` — read whole. `[family: sticky-pins-to-content-box]` (REQ-585) is the rule the padding move implements. REQ-319 (grow the SVG extent before writing scrollTop) survived the move to the board's clamp and its comment was extended rather than replaced. REQ-338 (roving tabindex over a virtualized list), REQ-336/337 (pointer capture retargets the synthesized click) and REQ-291 (a probe's default failure is a successful-looking measurement of nothing) shaped D-10 and the probe's four guards. REQ-322 (read a constant, never restate it) shaped `stickyAxisHeightPx` and is the reason D-21 states its one exception out loud. REQ-285 (for a rendering change the screenshot is the test) is what caught D-14. REQ-233 (a programmatic `.focus()` cannot answer a `:focus-visible` question) is why the ring is evidenced by a screenshot and not asserted.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — read whole. Nothing in it constrains a layout change directly; the closest is 0.295.1's "a surface that reads another's state has to be told when it changes", which is the same shape as D-16's stale-render finding.
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` — read for **Read first** and **Traps**. "`web/` is embedded, not read at runtime" is why every measurement above followed a rebuild.
- Crew: `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `frontend.md`, `ui-design.md`. `frontend.md`'s quality table drove the 320/768/1280 checks and the keyboard pass.

No listed lesson file was missing.

## Integration seams

One file outside the REQ's `## Scope`. Everything in it is already committed on my branch — this section is the disclosure, not a request for someone else to apply it.

`skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go`, four hunks:

1. `TestBrowserBehaviorTimelinePressBecomesAPanOnlyAfterMoving`, inside `trial()`, immediately after the `scrollIntoView` — the exact line added is
   `document.getElementById("board-main").dispatchEvent(new Event("scroll"));`
2. `TestBrowserBehaviorTimelinePointerAndKeyboardPathsStayAlive`'s setup, after its `scrollIntoView` — the same line.
3. `probe.aimAtARow`, after its `scrollIntoView` — the same line.
4. `timelineGroupingSnapshot()` — `viewportHeight: host.clientHeight` becomes `viewportHeight: document.getElementById("board-main").clientHeight`.

Each carries a six-line comment naming REQ-587 and the reason. If the maintainer would rather this file stayed untouched by this REQ, the alternative is to revert those four hunks and accept four red probes in the heavy lane until a follow-up REQ lands them; I did not think that was the better hand-back, and D-16 says why.

## Open question carried forward

D-06 (the focus ring on the axis) and D-07 (resetting the board's scroll on entering the Timeline) were escalated by the plan and I proceeded with both defaults, as instructed. Both are confirm-or-override items, not blockers:

- **D-06** — reversal is deleting two CSS rules and restoring `.timeline-scroll:focus-visible { outline: 2px solid var(--accent-claimed); outline-offset: -2px }`. But note D-14: whatever replaces it must still suppress the engine's own ring, or the full-height box gets one anyway.
- **D-07** — reversal is deleting three lines in `board-controls.js`. The alternative the plan named, remembering and restoring the board's position on leave and entry, is still about six lines of new state in the view switcher.
