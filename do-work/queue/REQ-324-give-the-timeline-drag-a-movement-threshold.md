---
id: REQ-324
title: "Give the timeline drag a movement threshold"
status: pending
created_at: 2026-08-22T22:08:34Z
user_request: UR-065
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-323]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-318, REQ-319, REQ-320, REQ-321, REQ-322, REQ-323]
batch: timeline-ux-audit
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Give the Timeline Drag a Movement Threshold

## What

Require a few pixels of movement before a press becomes a pan, so that clicking a row opens
its detail drawer and a hand tremor does not scroll the time axis.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The hint paragraph says "Click a row for its full detail." It works only for a perfectly
still press. `pointerdown` arms the pan immediately, so the first `pointermove` — one pixel
is enough — shifts the window and calls `renderAll`, which rebuilds every row node. The
delegated `[data-detail-kind]` click handler in `board-controls.js` then has no surviving
trigger under the pointer, and the drawer does not open. A one-pixel jiggle also moves the
chart, which nobody asked it to do.

## Context

`panState` is set in the `pointerdown` handler and consumed by `pointermove`, with no
distance test between them. Two consequences, one cause.

The second is cost: every `pointermove` runs a full `renderAll` (axis plus all visible
rows), and each render calls `xOfEpoch` several times per row, each of which calls
`plotWidth()`, each of which reads `scrollHost.clientWidth` — a forced layout read, hundreds
of times per frame during a drag.

## Detailed Requirements

- A movement threshold of a few pixels before the first pan render. Below it the press is a
  click: the drawer opens on release. Above it the pan runs, and the release does not open
  the drawer.
- Once the threshold is crossed, the pan is continuous from the original press point — no
  jump equal to the threshold at the moment it trips.
- Write the threshold decision as its own pure function so it can be probed without a DOM,
  the way this module's other decisions are.
- Coalesce renders during a drag to one per animation frame.
- Read the plot width once per render rather than once per `xOfEpoch` call.
- Keyboard activation (Enter / Space on a focused row) keeps working unchanged.

## Constraints

- Do not fix the click by making the row a native button — the row is a `<g>` in an SVG and
  `REQ-233`/`REQ-239` already settled the focus and activation contract for it.
- Serial with the rest of the `timeline-ux-audit` batch.

## Builder Guidance

**Certainty: Firm on the behaviour, latitude on the number.** A still-ish press must open
the drawer and a deliberate drag must not — that part is not open. The threshold itself is
yours: a few pixels, picked so an ordinary hand tremor stays under it and an intended drag
clears it immediately.

Scope cue: this is a bug fix, not a rework of the view's pointer model. The wheel, the
keyboard and the row activation contract stay as they are. The render coalescing rides along
because the same handler causes it; if it turns out to be more than a few lines, drop it and
say so rather than growing the REQ.

## Dependencies

`depends_on: [REQ-323]` — **ordering, not logic.** REQ-324 does not need anything REQ-323
produces; it needs REQ-323 not to be editing `web/board-timeline.js` at the same time. Every
REQ in the `timeline-ux-audit` batch writes that one file, and `write_set` is display-only —
`actions/work.md` computes a `--fan-out` wave from `depends_on` alone and explicitly does not
read `write_set`, `batch`, or the Constraints prose. Without this edge the batch's stated
serial requirement was a sentence nothing enforced, and a `--fan-out` run would have
dispatched four concurrent builders into one 1,100-line file.

**The cost, stated rather than hidden:** a chain gates on terminal *success*, so a `failed`
REQ upstream leaves the rest dependency-blocked until someone edits the chain or resolves the
failure. That is the trade for making the metadata say what the prose says.

## Red-Green Proof

**RED prompt/case:** In a browser probe against a generated board: dispatch `pointerdown` on
a timeline row, then `pointermove` 2px, then `pointerup` on the same row. Assert the detail
drawer opened. It does not — the 2px move panned the window and rebuilt the row nodes
mid-click. The same probe with no `pointermove` at all passes today, which is the point:
the behaviour depends on whether the hand was steady.

**Why RED now:** `pointerdown` arms `panState` unconditionally; the first `pointermove`
pans and re-renders.

**GREEN when:** the 2px probe opens the drawer; a probe that moves well past the threshold
pans the window and does *not* open the drawer on release; the pan itself has no jump when
the threshold trips; and a drag issues at most one render per frame.

**Validation:** Inferred during capture — audit finding, not one of the user's four items.
The diagnosis of *why* the click is lost is the audit's reading of the code plus the click
dispatch rule for a removed target; the probe above is what settles it either way, and the
threshold is the right fix regardless of which browser does what with the orphaned click.

## Assets

None. Reproduced by pointer interaction, not visible in the screenshot.

---
*Source: audit finding, UR-065 — "audit the timeline view, and make it more useful UIUX."*
