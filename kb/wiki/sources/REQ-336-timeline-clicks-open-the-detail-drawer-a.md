---
title: "Lessons from REQ-336: Timeline clicks open the detail drawer again"
type: source-summary
topic_cluster: timeline-and-metrics
sources: [raw/processed/2026-09-01/REQ-336-timeline-clicks-open-the-detail-drawer-a.md]
related:
  - page: concept-duration-estimation-and-breaks
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-336: Timeline clicks open the detail drawer again

Part of the [[concept-duration-estimation-and-breaks]] cluster.

## What the REQ was about

Mouse clicks on Timeline (Gantt) view REQ bars no longer open the detail drawer — every mouse
click inside `#timeline-scroll` is swallowed. Restore clicking while keeping the pan behaviour
REQ-333 and REQ-324 established.

## Solution summary

Moved the Timeline pan's `setPointerCapture` out of the `pointerdown` handler
into a named `capturePanPointer` that the `pointermove` handler calls once, on the false→true pan
engage. A press that never travels 4px now takes no capture, so its synthesized click keeps its
`<rect>` target and the delegated `[data-detail-kind]` handler opens the drawer; a drag still
captures the moment it engages, so its release is still guaranteed wherever it lands. Retargeted
REQ-333's structural capture assertion to the new call site, in two halves so neither can be
deleted silently.

## What worked

- Building the CDP probe before touching the code. It cost more than reading the diff would have,
  and it is the only reason three separate claims in this area are now measured rather than
  believed: that the click is retargeted, that a 150px drag clamps without moving the window, and
  that `pointerleave` reaches the host on a buttoned exit.
- Writing the gesture that probes the *fix's own* new risk (`hover-after-subthreshold-exit`)
  before deciding whether to guard against it. It said no guard was needed, which is a cheaper
  answer than the guard.

## What didn't work

- The first `drag-outside` gesture. It moved diagonally and left the host before the 4px threshold
  tripped, so it never engaged a pan — and it therefore proved nothing about REQ-333's
  release guarantee while looking like it did. A gesture aimed at a release path has to engage the
  thing being released first.
- Measuring "it panned" from the axis label text. A 150px drag pans and leaves the labels
  identical, because the window clamps at the bound. This is REQ-324's lesson arriving a second
  time, in a different disguise, in the same view.
- The first click attempt measured a closed drawer for a press that never reached the chart: the
  bar sat at y≈1538 in a 900px viewport. An element's `getBoundingClientRect` is not a clickable
  coordinate until the element is on screen.

## Worth knowing

- **Pointer capture retargets the synthesized click, not just the pointer events.** Any handler
  that captures on `pointerdown` inside a container with delegated click handling breaks every
  click in that container, and nothing about the pan logic makes that visible. `board-detail.js`'s
  resizer captures on `pointerdown` and is fine only because nothing delegates clicks through it.
- **This view has now eaten a click for two independent reasons** (a re-render that rebuilt the
  node, and capture retargeting). The movement threshold protects against the first only. A third
  cause would look identical from the outside, so a behavioural check on the click path is worth
  more here than any amount of reasoning about the pan state machine — which is REQ-337.
- REQ-333's stated premise, that Chromium suppresses boundary events while a button is held, did
  not reproduce in Chromium 1194: `pointerleave` reached the host four times on a buttoned exit.
  The capture is still worth taking at the engage, but a future reader should not treat that
  premise as established.
- The repo's browser probe lane runs under `--dump-dom` with no CDP channel, so it cannot dispatch
  trusted input at all. Anything needing real pointer semantics has to drive
  `Input.dispatchMouseEvent` itself; Node 22's global `WebSocket` is enough to speak the protocol
  with no dependencies.

## Back-reference

See `do-work/archive/UR-067/REQ-336-timeline-clicks-open-the-detail-drawer-again.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `4527a50`.
