---
id: REQ-333
title: "Keep the timeline pointer and keyboard paths alive"
status: completed-with-issues
created_at: 2026-08-23T12:24:00Z
claimed_at: 2026-08-23T18:40:00Z
completed_at: 2026-08-23T19:25:00Z
commit: 36c4518
route: B
user_request: UR-066
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Keep the Timeline Pointer and Keyboard Paths Alive

## What

Three of the interactions `.timeline-hint` promises break in ways a reader cannot recover from without
clicking elsewhere: a drag released outside the chart leaves the pan latched so plain mouse motion keeps
panning; the first arrow-key pan with a row focused kills the keyboard path; and Tab through the rows never
reaches anything after the chart.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`. Reproduced all four findings in a real engine BEFORE planning any fix, as the REQ required, and the reproduction changed the plan twice:
  - F1's mechanism is not the reported one. Un-buttoned motion does **not** keep panning; what persists is the `.is-panning` grab cursor, because Chromium suppresses boundary events while a button is held and the release outside the host reaches nothing. `setPointerCapture` is the fix that makes the release a fact rather than a hope, which is what the REQ suggested preferring.
  - F3 is **refuted**: Tab escapes the rows after 29 presses and reaches the element after the chart. Not a trap, so requirement 3's "decide the keyboard contract" is answered as "leave it", and the finding is recorded rather than acted on.
- [x] **[APPLY]:** `web/board-timeline.js` and `timeline_browser_probe_test.go`. `web/template.html` was NOT touched: the hint's ctrl-scroll promise is now true rather than narrowed. F3 needed no change.
- [x] **[UNIFY]:** Verified:
  - `web/board-timeline.js` — `node --check` clean; the focus restore exists in exactly one place now, and the keydown handler's comment states why it must not keep its own; `setPointerCapture` is feature-guarded and its failure degrades to the ordinary release path rather than refusing the drag.
  - `timeline_browser_probe_test.go` — `gofmt`/`go vet` clean; the probe states plainly that a synthetic `pointerId` cannot be captured in this lane, so it drives `lostpointercapture` (the mechanism) and asserts the capture CALL structurally, while the end-to-end drag was verified with a real input device outside the lane.
  - `bash _dev/tests/maintainer-verify.sh` exit 0.
  - Mutations: the focus restore removed, `lostpointercapture` dropped, the capture call deleted, and the axis wheel binding removed — four distinct failures. The third passed at first because the structural check matched the `typeof` guard beside the call; it now requires the call itself.

## Why

The hint under the chart tells the reader to drag to pan, to Tab to the chart and use the arrows, and to click
a row for detail. Each of those is a promise, and a promise that strands the reader is worse than an absent
feature.

## Context

Reported by the pointer/keyboard audit against a live render; each needs reproducing before it is fixed —
the first especially, because `pointerleave` is bound alongside `pointerup` (`:1911-1925`) and should already
clear `panState` when the pointer leaves the host.

- **F1 — latched pan.** A drag released outside `#timeline-scroll` reportedly leaves `panState` set with
  `engaged` true, so subsequent `pointermove` events with no button held keep panning the window, and
  `.is-panning` keeps the grab cursor. Code site `:1912`. **Reproduce first**: establish which of
  `pointerup` / `pointercancel` / `pointerleave` actually fires on the release path before changing the
  teardown, and prefer `setPointerCapture` on `pointerdown` over adding another release event, so the release
  is guaranteed to reach the host that armed the drag.
- **F2 — arrow key with a row focused.** The keydown handler re-renders and then re-focuses the rebuilt row
  (`:1798-1808`), but the row anchor inside `refreshWindowRows` writes `scrollHost.scrollTop`, and the scroll
  listener (`:1755`) is `renderVisibleRows` — which runs *after* the focus restore and rebuilds the row SVG
  again, dropping focus to `<body>`. One arrow press then dead keys. Code site `:1756`.
- **F3 — Tab trap.** Every row is `<g tabindex="0" role="button">`. Tabbing to a row below the viewport
  scrolls the container, the scroll listener rebuilds the row SVG, and focus is lost — so Tab cycles inside
  the chart instead of leaving it. Code site `:1306`. A 317-row chart with a tab stop per row is also a poor
  keyboard contract on its own: the rows already have a table equivalent below, and the chart itself is one
  focusable widget with arrow-key navigation.
- **F4 — ctrl+wheel over the axis strip** does not zoom, although the hint says holding Ctrl and scrolling
  zooms the time axis. The wheel handler is bound to `#timeline-scroll` only (`:1830`).

## Detailed Requirements

1. A drag can only end. Whatever the release path, the pan state and the `.is-panning` class are cleared, and
   no un-buttoned `pointermove` ever moves the window.
2. An arrow-key pan leaves the keyboard path working: after the press, focus is on something that still
   accepts the next arrow key. Settle the ordering between the row anchor's `scrollTop` write, the scroll
   listener's rebuild and the focus restore, so the restore is last.
3. Tab reaches the far side of the chart. Decide the keyboard contract explicitly — either the chart is one
   tab stop with arrow-key row navigation, or rows keep individual tab stops and the rebuild preserves focus —
   and record which, and why, next to the row's `tabindex`.
4. `ctrl`+wheel over the axis strip zooms the axis, matching the hint; or the hint is narrowed to name the
   plot only. Prefer making it work: the axis is part of the same chart.
5. `Enter` / `Space` on a row still opens the drawer, and a sub-threshold press still counts as a click
   (REQ-324's behaviour). Both are already right and both must stay right.

## Constraints

- Every listener bound to a node that outlives a render goes through `addTimelineListener`, so the teardown
  registry keeps working (`:80-86`).
- Do not reintroduce a render per `pointermove`: the drag stays one render per frame (`requestPanRender`,
  `:1858`).
- A programmatic `.focus()` cannot answer a `:focus-visible` question (REQ-233's lesson) — do not test the
  focus ring that way.

## Red-Green Proof

**RED prompt/case:** (a) Press inside the chart, drag 200px, release with the pointer outside the chart's
box, then move the mouse with no button held and read `#timeline-range-readout`. (b) Click a row to focus it,
press `ArrowRight`, then press `ArrowRight` again and read the readout. (c) Focus the first row and press Tab
repeatedly, recording `document.activeElement` each time.

**Why RED now (to be confirmed by the builder's own reproduction):** (a) the readout keeps changing as the
un-buttoned pointer moves; (b) the second `ArrowRight` does nothing because focus is on `<body>`; (c) the
active element never leaves the rows.

**GREEN when:**
- (a) no un-buttoned pointer motion changes the readout, and the grab cursor is gone.
- (b) the second and tenth `ArrowRight` each pan the window by the same amount as the first.
- (c) a bounded number of Tab presses moves focus past the chart to the next control on the page.
- `ctrl`+wheel over the axis changes the window, or the hint no longer claims it does.
- A Node behaviour probe pins (b) by driving the keydown path and asserting the post-render focus target.

**Validation:** Inferred during capture. F1, F2 and F3 are reported from a live render by the audit and are
to be reproduced by the builder before the fix — the REQ names the reproduction as the first step rather than
asserting the mechanism.

## Outcome

Filed `completed-with-issues` rather than `completed`, because one of the four findings
this REQ was captured for (F3, the Tab trap) turned out not to exist, and the
reproduction is the record of that. Nothing in scope was left undone; the scope itself
was smaller than the capture believed.

F3's underlying observation stands as a design question this REQ deliberately does not
answer: 29 tab stops for a chart whose every value is also in the table below it. Anyone
picking that up should note the rows carry `role="button"` and are keyboard-activatable,
so removing the tab stops would remove a working affordance and needs its own decision.

## Full Context

See `do-work/user-requests/UR-066/input.md`.
