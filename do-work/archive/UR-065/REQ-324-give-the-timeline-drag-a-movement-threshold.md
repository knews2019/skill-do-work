---
id: REQ-324
title: "Give the timeline drag a movement threshold"
status: completed
created_at: 2026-08-22T22:08:34Z
claimed_at: 2026-08-23T04:47:00Z
completed_at: 2026-08-23T05:12:00Z
commit: 3486ab2
route: B
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
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-23T04:47:20Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - dependency depth 6
    - RED reproduced before implementing
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
  - skills/do-work-board/tools/queue-kanban/browser_probe_test.go
kb_status: promoted
kb_entry: REQ-324-give-the-timeline-drag-a-movement-thresh.md
---

# Give the Timeline Drag a Movement Threshold

## What

Require a few pixels of movement before a press becomes a pan, so that clicking a row opens
its detail drawer and a hand tremor does not scroll the time axis.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reproduce RED first — the REQ says the probe settles the diagnosis either
  way. Then: `timelinePanEngaged(alreadyEngaged, pressX, pointerX)` as a pure latching
  predicate beside the module's other decisions; `panState` gains an `engaged` flag and the
  pointermove handler returns early until it trips; the shift keeps measuring from the press
  point. Coalesce through one `requestAnimationFrame` handle, flushed on release. Memoize the
  plot width, invalidated at the top of `renderAll`.
- [x] **[APPLY]:** Three files. `web/board-timeline.js`, `generate_test.go` (the pure
  threshold probe), `timeline_browser_probe_test.go` (two browser probes) plus a small
  harness addition in `browser_probe_test.go`.
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file. `node --check` clean; `go vet
  ./...` clean; `bash _dev/tests/maintainer-verify.sh` exit 0. No debug artifacts. The
  scratch reproduction pages live in the session scratchpad, not the repo.

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

---

## RED, Reproduced

Before writing any fix, the REQ's own probe was run against a board generated from the current
tree. Chromium 1194, real `PointerEvent`s on a real timeline row:

| press | drawer opened |
|---|---|
| still (no `pointermove`) | **yes** |
| 2px jiggle | **no** |
| 90px drag | no |

The diagnosis holds, and with one detail the REQ did not predict: the window was at *Fit all*,
so the pan **clamped to the same window** and the chart did not visibly move — yet the click
was still lost. The re-render alone is sufficient; the window does not have to move. That is
why the fix is a threshold on the *render*, not on the window shift.

## Decisions

**D-01 — `TIMELINE_PAN_THRESHOLD_PX = 4`, latching.** Four clears an ordinary trackpad tremor
and is reached in the first frame of any intended drag. Latching, because a slow drag that
wanders back inside four pixels would otherwise flicker the grab cursor and — worse — end as a
click that opens a drawer nobody asked for. The pure probe pins the boundary (3.99 does not
engage, 4.00 does, in both directions) and holds the constant to a 2–8px band rather than to
the literal 4.

**D-02 — The shift keeps measuring from the PRESS point, not the trip point.** The REQ's "no
jump equal to the threshold" reads two ways, and this is the one that leaves no permanent
offset: the chart tracks the pointer one-for-one from the press. Re-anchoring at the trip point
has no visible step but parks the chart four pixels behind the pointer for the rest of the
drag. The browser probe distinguishes them by dragging the same 120 pixels in one move and in
two and requiring the same end window; under trip-anchoring the two-move drag lands four pixels
short. If the other reading was meant, that test is the one to change.

**D-03 — The grab cursor moved to the trip too.** It was applied on `pointerdown`, which told
the reader they were dragging while they were clicking. Not in the REQ's list; it is the same
defect wearing a cursor, and leaving it would have made the fix invisible to the person using
it.

**D-04 — The coalescing frame handle lives at module scope.** Inside the render closure it
would have been the exact hazard the file already documents for the table-rebuild frame: a
filter change re-enters the module, and a frame scheduled by the previous render would fire
against that render's stale `rows`. It is released beside the table frame, at the same call.

**D-05 — One extra file beyond the write set.** `browser_probe_test.go` gained
`runBrowserBehaviorProbeInDirectory`. The existing runner writes the probe page into an empty
temp directory, and the generated `index.html` loads `board-data.js` from beside itself — so
the real board rendered empty and every assertion measured nothing. Six lines, and the
alternative was testing a copy of the page rather than the page (REQ-305's lesson).

## Evidence

- `bash _dev/tests/maintainer-verify.sh` → **exit 0** (`GOTOOLCHAIN=go1.26.1`,
  `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium`).
- GREEN on the same reproduction that produced RED: still press opens the drawer, 2px press
  opens the drawer, 120px drag pans and does not.
- **Measured, not asserted:** one render read `scrollHost.clientWidth` **171 times** before the
  memo and **1** after. Five pointermoves in one task produced 5 renders before coalescing and
  0-then-1-on-release after.
- Fifteen mutations across the three new probes, each failing its intended assertion and only
  that one: threshold 0 / 20, exclusive comparison, signed distance, latching removed (twice —
  pure and behavioural), grab cursor on press, grab cursor never applied, trip-point anchor,
  memo removed, memo never invalidated, coalescing removed, release dropping the frame.

## Review Fixes

Three, all accuracy rather than behaviour — the diff's logic survived the pass.

**R-01 — The threshold constant's comment blamed the pan, not the render.** It said a still-ish
press "used to pan on its first pixel, which re-rendered the rows and lost the click", which
reads as though the window moving is what did it. The RED run showed otherwise: at *Fit all* the
pan clamps to the same window and the click was lost anyway. The comment now names the render as
the cause and says why the threshold gates the render rather than the shift — the same
distinction that would have made a narrower fix pass a narrower test.

**R-02 — "hundreds of layout reads" replaced with the measured number.** 171, on this repo's
board, Chromium 1194 (`_dev/primes/prime-kanban-board.md`: record the browser build beside every
measured number). The comment also now says why invalidating at the top of `renderAll` covers
the resize case — the resize listener is `renderAll`.

**R-03 — A pointless local rebinding** in the new probe runner (`probeDirectory := siteDirectory`).

Gate re-run after the three fixes: `bash _dev/tests/maintainer-verify.sh` → **exit 0**.

## Discovered Tasks

None. REQ-325 remains `pending-answers`.

## Lessons Learned

- **Reproduce the failure before believing the explanation, because the reproduction carries
  detail the explanation does not.** The REQ's diagnosis was right, but the RED run showed the
  click dying even when the pan clamped and the window did not move at all. Had the fix been
  keyed on "the window moved" rather than on "a render happened", it would have passed a
  narrower test and left the defect at *Fit all* — the zoom level the board opens on.
- **An assertion taken after the state it describes has been cleared measures nothing.** The
  grab-cursor check read `is-panning` after `pointerup`, which removes it, so it passed with
  the cursor restored to the pressed-immediately behaviour. The mutation caught it; the
  assertion had to move inside the press.
