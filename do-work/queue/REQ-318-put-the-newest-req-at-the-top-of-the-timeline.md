---
id: REQ-318
title: "Put the newest REQ at the top of the timeline"
status: pending
created_at: 2026-08-22T22:08:34Z
user_request: UR-065
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-319, REQ-320, REQ-321, REQ-322, REQ-323, REQ-324]
batch: timeline-ux-audit
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline.go
  - skills/do-work-board/tools/queue-kanban/timeline_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
---

# Put the Newest REQ at the Top of the Timeline

## What

Reverse the Timeline view's row order so the most recent REQ is the first row under the
axis and the oldest is last.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

On a 309-REQ queue the current work is 5,700 pixels below the fold. Every visit to the
view opens on REQ-001.

## Context

`buildTimelineAggregate` (`timeline.go`) sorts rows by `created_at` ascending with
`requestIdLess` as the tiebreak, and `board-timeline.js` renders them in payload order.
Reversing in Go keeps one ordering decision in one place and keeps the client a renderer.

Two things read the row order and must still mean what they say afterwards:

- `timelineFirstOpenRowIndex` / `timelineNowJump` scroll to the first still-open row. With
  newest first, that row sits near the top rather than near the bottom; the function still
  works, but confirm the *Now* button lands somewhere sensible rather than assuming it.
- The subhead sentence in `renderTimelineView` literally reads "in capture order, oldest at
  the top". It is false the moment this lands and is this REQ's to fix.

## Detailed Requirements

- Newest `created_at` first, oldest last.
- The tiebreak for two REQs captured in the same instant stays deterministic — reverse it
  with the sort so a build never swaps two rows.
- The subhead states the new order in plain words.
- `timeline_test.go` currently pins oldest-first. That assertion is in scope and must be
  rewritten to pin newest-first, not deleted. Name the change in the hand-back: a quietly
  edited test looks identical in a diff.

## Constraints

- Serial with the rest of the `timeline-ux-audit` batch — all of them write
  `web/board-timeline.js`.
- The Calendar view already reads newest-first; matching it is the point, not a coincidence.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repo's archive and open the Timeline tab.
The first row under the axis is `REQ-001`; the newest REQ is only reachable by scrolling to
the bottom of the list. In Go, `buildTimelineAggregate` returns `rows[0]` as the earliest
`created_at`, and `timeline_test.go` asserts exactly that.

**Why RED now:** The sort in `buildTimelineAggregate` is ascending by `created_at`.

**GREEN when:** `rows[0]` is the newest `created_at` and the last row is the oldest; two
rows sharing an instant keep a stable, reversed id order across repeated builds; the
generated board opens with recent work on screen; the subhead no longer says "oldest at the
top".

**Validation:** User confirmed (stated as item 1 of the request).

## Assets

Screenshot described in `do-work/user-requests/UR-065/input.md` — rows `REQ-001` through
`REQ-042` filling the visible list in ascending order.

---
*Source: "1.  most recent REQ's should be on top."*
