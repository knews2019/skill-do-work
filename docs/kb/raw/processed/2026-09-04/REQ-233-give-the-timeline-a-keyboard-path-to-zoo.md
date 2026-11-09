---
source_type: req_lesson
req_id: REQ-233
req_path: do-work/archive/UR-051/REQ-233-give-the-timeline-a-keyboard-path-to-zoom-and-pan.md
date: 2026-08-18
domain: general
module: _dev/primes
tags: [general, give, timeline, keyboard]
---

# Lessons from REQ-233: Give the Timeline a keyboard path to zoom and pan

## What the REQ was about

The Timeline view's zoom and pan are pointer-only. The three zoom buttons are reachable by keyboard, but panning the time axis is not, and both affordances are described only in a hint line of prose beside the chart rather than anywhere assistive technology reads.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified — integration seam, applied by the orchestrator inside the merge commit)

## What worked

- Making the keyboard path a *pure function* that returns a window or `null`, and routing its zoom through the transform the pointer already used. That is what turned "the two paths must not diverge" from a promise into a structural fact — there is no second floor, ceiling, or clamp to drift. It also made the probe able to assert against the pointer path's own anchor rather than a copy of it.

## What didn't work

- Checking the focus ring with a programmatic `.focus()`. It reported `matchesFocusVisible: false` and a computed `outline: none` — which reads exactly like a missing rule, and is actually Chrome correctly refusing `:focus-visible` for non-keyboard focus. The rule was fine. Only a real Tab keypress answers this question; a scripted focus call will tell you the ring is broken when it is not.
- The write set. Requirement 3 was a CSS requirement from the moment it was written, and neither the REQ's captured write set nor the builder's scope named a stylesheet. The builder handled it correctly at the boundary, but the boundary should not have been there.

## Worth knowing

- `renderVisibleRows` rebuilds every row node, so anything holding a reference to a row across a render is holding a dead element. This bit the keyboard path (focus fell to `<body>` after one press) and will bit anything else that keeps row state — the fix is to capture the row's `data-detail-id` before the render and re-query after it.

## Back-reference

See `do-work/archive/UR-051/REQ-233-give-the-timeline-a-keyboard-path-to-zoom-and-pan.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `9b2578b`.
