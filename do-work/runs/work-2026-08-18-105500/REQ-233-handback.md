# REQ-233 hand-back — Give the Timeline a keyboard path to zoom and pan

**Branch:** `worktree-agent-REQ-233-timeline-keyboard-path`
**Commit:** `d5b96ae` — `[REQ-233] Give the Timeline a keyboard path to zoom and pan`
**Worktree:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-233-timeline-keyboard-path`

## File manifest

| Action | File | What |
|---|---|---|
| Modified | `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` | Added `TIMELINE_PAN_FRACTION`, the pure `timelinePannedWindow` and `timelineKeyboardWindow` transforms, and extended the existing keydown listener to move the window and restore row focus after the re-render. |
| Modified | `skills/do-work-board/tools/queue-kanban/web/template.html` | Timeline panel's `aria-label` now states the interaction; `#timeline-scroll` gained `tabindex="0"`; the visual hint line gained the keyboard sentence. |
| Modified | `skills/do-work-board/tools/queue-kanban/generate_test.go` | Added `TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer` (Node probe) and `TestTimelinePanelStatesItsKeyboardInteraction` (generated-HTML assertion); added the `math` import. |

Nothing outside the declared write set was written. `git status` in the worktree shows exactly those three files.

## Red-Green Proof

### RED — run 1, before any handler existed

```
$ go test -count=1 -run 'TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer|TestTimelinePanelStatesItsKeyboardInteraction' .
--- FAIL: TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer (0.22s)
    generate_test.go:2315: web/board-timeline.js declares no numeric constant TIMELINE_PAN_FRACTION
--- FAIL: TestTimelinePanelStatesItsKeyboardInteraction (0.22s)
    generate_test.go:2494: the timeline panel's accessible name does not state which keys pan (wanted "arrow keys" in "<section id=\"view-timeline\" class=\"view-panel\" aria-label=\"Timeline\" hidden>")
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.760s
FAIL
EXIT=1
```

That first failure is the reference-error class the brief allows only with a second run, so there is one:

### RED — run 2, handler present but wrong (pan written without the bound clamp)

Intermediate `timelinePannedWindow` — the shipped body minus the clamp:

```js
function timelinePannedWindow(windowStartMs, windowEndMs, panFraction, boundStartMs, boundEndMs) {
  var windowSpanMs = windowEndMs - windowStartMs;
  var nextStartMs = windowStartMs + windowSpanMs * panFraction;
  return { windowStartMs: nextStartMs, windowEndMs: nextStartMs + windowSpanMs };
}
```

```
$ go test -count=1 -run 'TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer' .
--- FAIL: TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer (0.27s)
    generate_test.go:2430: panning right settled with the window ending at 9720000000 ms, want the range edge 2592000000 ms
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.575s
FAIL
EXIT=1
```

A real assertion failure: 40 held ArrowRight presses walked the window to 9.72e9 ms on a board whose range ends at 2.592e9 ms — off the data entirely.

### GREEN

```
$ go test -count=1 -run 'TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer|TestTimelinePanelStatesItsKeyboardInteraction' -v .
=== RUN   TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer
--- PASS: TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer (0.26s)
=== RUN   TestTimelinePanelStatesItsKeyboardInteraction
--- PASS: TestTimelinePanelStatesItsKeyboardInteraction (0.20s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	0.934s
EXIT=0
```

What the probe asserts, each separately so one regression reports one failure:

- ArrowRight/ArrowLeft shift the window by exactly `windowSpan × TIMELINE_PAN_FRACTION`, in each direction, without changing the span.
- The step is bounded — greater than zero and less than the visible span.
- 40 held presses in each direction settle exactly on the range edge with the span intact.
- `+` bottoms out at the same span the pointer path bottoms out at (driven through `timelineZoomedWindow` at the wheel's off-centre 0.25 anchor, not the keyboard's centred 0.5), and that span is the renderer's own `TIMELINE_MIN_SPAN_MS`.
- `-` tops out at the same span the pointer path tops out at, and that span is the full bound span.
- `Enter`, `" "`, `Spacebar`, `Tab`, `ArrowUp`, `ArrowDown`, `a` all return null, so none of them moves the window.

Constants come from `timelineProbePreamble`, which reads them out of the shipped `web/board-timeline.js` — the probe cannot pass against numbers the view does not use.

## maintainer-verify.sh

```
$ bash _dev/tests/maintainer-verify.sh
...
Maintainer verification passed.
EXIT=0
```

Exit code **0**. Not piped. Includes `queue-kanban go vet`, the uncached ordinary tests, and the strict JavaScript behavior lane.

## Rendered and looked at

Built `go build -o /tmp/qk-233 .` and generated a static board against the worktree (232 REQs, 54 URs). In `/tmp/board-233/index.html`:

- `id="view-timeline" class="view-panel" aria-label="Timeline. With the chart focused, arrow keys pan the time axis and the plus and minus keys zoom it."`
- `<div class="timeline-scroll" id="timeline-scroll" tabindex="0">`
- `var TIMELINE_PAN_FRACTION = 0.15;` and both transforms present in the inlined client.

Then served it on loopback and drove it in a real browser:

- `#timeline-scroll` takes focus; `document.activeElement.id === "timeline-scroll"`, `tabIndex === 0`.
- Two `+` presses narrowed the axis from `28 May … 20 Aug` to `22 Jun … 25 Jul`, centred.
- One ArrowRight moved it to `27 Jun … 30 Jul` — about 5 days, 15% of a ~33-day window.
- 40 ArrowRight presses pinned the window's right edge at `20 Aug`, the range end. 80 ArrowLeft presses pinned its left edge at `28 May`, the range start. 40 `-` presses restored the exact fit-all window.
- An ArrowDown keydown was **not** prevented and did not move the window, so the scroll container keeps its native vertical scrolling.
- Only console error on the page was a `favicon.ico` 404 from my ad-hoc static server; nothing from the board.

I did **not** verify the pointer path changed behavior, because it did not: the wheel, drag, and zoom-button handlers are untouched, and the keyboard path calls the same `timelineZoomedWindow`.

## Integration seams — one CSS rule I did not write

`web/board.css` is outside my write set (you said a sibling builder is in it), so requirement 3 is **not satisfied on my branch**. The chart's only focus indication today is the user-agent default, which I measured in the browser as `outline: rgb(0, 95, 204) auto 1px` — exactly what the REQ rules out. There is no `.timeline-scroll` focus rule in `board.css` at all.

**Exact rule to add:**

```css
.timeline-scroll:focus-visible {
  outline: 2px solid var(--accent-claimed);
  outline-offset: -2px;
}
```

**Exact location:** `skills/do-work-board/tools/queue-kanban/web/board.css`, immediately after the existing `.timeline-scroll.is-panning` block (currently line 2109-2111 on my branch), so the container's own rules stay together:

```css
.timeline-scroll.is-panning {
  cursor: grabbing;
}
                          <-- insert here
.timeline-row-hit {
```

Verified by injecting exactly that rule into the live page: it computes to `outline: rgb(58, 107, 196) solid 2px` in light and draws a clear, unclipped ring; I screenshotted it and it reads well against the chart. `--accent-claimed` is the same token every other focus ring on the board already uses (`.control-button:focus-visible`, `.req-card:focus-visible`, `.calendar-chip:focus-visible`), and in the dark palette it is `#6f9ce6` against `--surface-1: #131720` — about 7:1, well past the 3:1 non-text minimum.

`outline-offset` is **-2px**, not the `+2px` the other rules use. That is deliberate: `.timeline-scroll` is an `overflow-y: auto` container sitting flush under the axis, so a positive offset would draw the ring outside the scroll box where the axis and the surrounding card clip it. Drawn inward, the whole ring is visible on all four edges.

No other seam. Nothing else outside the write set is needed.

## Listener teardown — how I hooked in

The registry is `timelineListenerTeardowns` in `web/board-timeline.js` (declared at line 47), with `addTimelineListener(target, eventName, handler, options)` pushing a remover and `releaseTimelineListeners()` draining it at the top of `renderTimelineView()`.

I added **no new listener**. The chart already had exactly one `keydown` listener on `scrollHost`, registered through `addTimelineListener`, for row activation. I extended that same handler: it asks `timelineKeyboardActivationTarget` first and returns early on Enter/Space, then falls through to `timelineKeyboardWindow`. One listener, already in the registry, so nothing new can leak across the re-render a filter change triggers.

## Decisions

**Left/Right pan the time axis; Up/Down are left alone.** The REQ says "arrow keys pan the time axis". The time axis is horizontal, and `#timeline-scroll` is a 58vh scroll container whose Up/Down keys already scroll a 232-row queue. Claiming Up/Down would take a working interaction away to duplicate one the Left/Right pair already provides. `timelineKeyboardWindow` returns null for them, and the handler only calls `preventDefault()` on keys it actually owns — verified in the browser that an ArrowDown keydown comes back with `defaultPrevented === false`.

**`=` and `_` zoom too, alongside `+` and `-`.** On a US layout `+` needs Shift. Browsers themselves accept the unshifted `=` for zoom-in; matching that costs one clause.

**Pan step is 15% of the visible span, not a fixed duration.** A fixed number of milliseconds is imperceptible zoomed all the way out and a jump zoomed all the way in. The probe asserts the fraction rather than a magic number, reading `TIMELINE_PAN_FRACTION` out of the shipped file.

**Zoom anchors at 0.5 for the keyboard.** There is no pointer position to anchor to. This matches what the existing zoom buttons already do, so keyboard zoom and button zoom behave identically.

**Focus is restored to the same row after a keyboard-driven render.** `renderVisibleRows()` clears and rebuilds every row node, so a row that had focus is a dead element by the time the next keypress arrives and focus falls to `<body>` — a keyboard user would get one arrow press and then nothing. Moving the window never changes the vertical scroll, so the row is still inside the virtualized slice; the handler captures its `data-detail-id` before the render and focuses the replacement after. This is scoped to the keyboard handler only — the wheel and drag paths imply a pointer, where focus is not on a row.

**Added a second, non-Node test for the two template attributes.** `TestTimelinePanelStatesItsKeyboardInteraction` pins the panel's accessible name and the chart's `tabindex`. Both are single attributes, silently droppable in a template edit, and both are load-bearing: without `tabindex` the whole keyboard path is unreachable, and without the accessible name the interaction is discoverable only by sighted readers of the hint line — which is the defect this REQ exists to fix. It names those failures rather than restating the markup.

**Kept the visual hint line and extended it.** The REQ asks that the accessible name state the interaction "rather than leaving it to the visual hint line" — not that the hint line lose it. A sighted keyboard user reads the hint, not the `aria-label`, so the hint now names the keys too.

## Discovered Tasks

- **`.timeline-row { outline: none; }` suppresses the row focus ring** (`web/board.css` line ~2122-2124). Rows fall back to `.timeline-row:focus .timeline-row-hit { fill: var(--surface-2); }` — a background-tint change of one step, which is a much weaker focus signal than the 2px accent ring every other focusable thing on the board gets, and it does not distinguish `:focus` from `:focus-visible`. Out of scope here (rows are not what this REQ adds, and the file is outside my write set), but it is the same requirement-3 concern applied to the element next door.
- **The Timeline's `role="img"` rows SVG carries a long `aria-label` describing the chart, while the individual rows inside it carry their own `aria-label`s.** A screen reader treating the `<svg role="img">` as an atomic image may not reach the row labels at all. Worth a deliberate decision about which layer owns the description; I did not touch it.
- **Zoom state persists across tab switches but not across a filter change that empties the view.** `timelineViewState` is held outside `renderedOnce` on purpose, but the `rows.length === 0` path returns before re-fitting, so a filter that matches nothing and is then cleared leaves the previous zoom. Probably correct, possibly surprising; noting it rather than changing it.
