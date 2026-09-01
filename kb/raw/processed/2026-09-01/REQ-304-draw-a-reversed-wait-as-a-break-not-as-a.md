---
source_type: req_lesson
req_id: REQ-304
req_path: do-work/archive/UR-062/REQ-304-draw-a-reversed-wait-as-a-break-not-as-a-valid-bar.md
date: 2026-08-21
domain: frontend
module: _dev/primes
tags: [frontend, draw, reversed, wait, break]
---

# Lessons from REQ-304: Draw a reversed wait as a break, not as a valid bar

## What the REQ was about

The Timeline view draws the wait segment unconditionally (`web/board-timeline.js:639-641`), and
`drawSegment` sorts its endpoints with `Math.min`/`Math.max` (`:591-593`), so a wait whose
`claimed_at` precedes its `created_at` paints as an ordinary positive-width waiting bar. The work
segment already handles this: `:655-666` draws a break marker at the claim instant when
`workMinutes < 0`, with the comment "A reversed span has no width to draw honestly". Give the wait
the same treatment.

## Solution summary

Fifteen lines in the renderer, mirroring the branch six lines below it including its
reasoning. The test runs the real `renderTimelineView` over a stub DOM and asserts three wait shapes
in one render pass — reversed, ordinary closed, and open — because a fix that turned every wait into
a break would satisfy a reversed-only test.

## What worked

- **A missing-branch fix needs a fixture that can fail in both directions.** The obvious test is the
  reversed row alone, and `if (true)` passes it. Rendering the reversed, ordinary and open shapes
  in one pass is what makes over-application fail.
- **Anchor a mirrored branch to its own span, not to the one it mirrors.** Copying the work branch's
  `workStartMs` anchor would have put the wait's break at the claim instant, stacking two markers
  whenever both spans reverse. The mirror is the shape, not the coordinate.
- **A stub DOM is cheaper than it looks, and iterating on it is the fast path.** Four Node failures —
  `classList`, `querySelectorAll`, `setActiveButton`, the scroll geometry — each named exactly
  what to add. Guessing the whole surface up front would have taken longer than letting it fail.

## Back-reference

See `do-work/archive/UR-062/REQ-304-draw-a-reversed-wait-as-a-break-not-as-a-valid-bar.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `5e08a31`.
