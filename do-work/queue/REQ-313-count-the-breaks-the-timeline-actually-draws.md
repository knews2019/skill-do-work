---
id: REQ-313
title: "Count the breaks the timeline actually draws"
status: pending-answers
created_at: 2026-08-21T08:25:57Z
status_changed_at: 2026-08-21T08:25:57Z
user_request: UR-062
addendum_to: REQ-304
domain: frontend
review_generated: true
impact: impact-user-visible
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
depends_on: []
maintenance: false
---

# Count the Breaks the Timeline Actually Draws

## What

The Timeline view's summary line renders "N with broken stamps, drawn as breaks"
(`web/board-timeline.js:546-557`), and N counts rows where `row.anomaly` is true. That flag
carries `detectCompletionAnomaly`'s verdict, which `model.go:246-252` scopes to **completion
bookkeeping** — no completion instant at all on a terminal REQ. It is false for both spans the view
draws as breaks:

- a reversed **work** span, whose `completed_at` parses fine but precedes `claimed_at`
- a reversed **wait** span, which REQ-304 just started drawing as a break

So the view draws break markers and the sentence beneath them says zero.

## Why

The count is the only place a reader learns how many rows to distrust without scrolling the whole
Gantt. A count that says zero while breaks are on screen teaches the reader to ignore the sentence.

## Context

Discovered during REQ-304, which fixed the reversed-wait rendering and was explicitly forbidden by
its own Context from broadening the anomaly flag to fix the count. That constraint stands and is the
reason this is a separate REQ rather than one more line in that one:

- `timeline.go:22-24` states that `detectCompletionAnomaly` decides what is broken and this file
  consumes that verdict rather than restating it. Recomputing an anomaly inside the renderer would
  create the second definition that comment exists to prevent.
- REQ-280 closed the **detection** side — the `created_at <= claimed_at <= completed_at` ordering
  probe in `queue-kanban verify` and `forensics.md` Check 12. It did not touch this count, and
  the count is not detection.

The shape that satisfies both: count what the renderer **drew**, not what a flag says. The draw pass
already decides, per row, whether each segment became a break; the summary can total those decisions
instead of asking a flag that answers a different question.

## Detailed Requirements

- The summary line's count equals the number of rows the current render drew at least one break
  marker for — reversed wait, reversed work, or the existing completion anomaly.
- No new anomaly verdict, flag, or reason is added to the timeline row model, and
  `detectCompletionAnomaly` is not broadened.
- The count follows the filters, like every other number on that line.
- The sentence still reads correctly when the count is zero.

## Red-Green Proof

**RED prompt/case:** Render a timeline whose rows include one reversed wait and one reversed work
span, neither carrying `anomaly: true`, and read the summary text. Today it omits the break clause
entirely while two break markers are drawn.
**GREEN when:** The same render reports 2, and a render with no reversed span and no completion
anomaly still omits the clause.
**Validation:** Discovered task from REQ-304; apply `actions/work-reference.md` →
**Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-304: the Timeline view prints a line
  saying how many REQs have broken timestamps and are drawn as break markers, but it counts a flag
  that only covers one kind of breakage. Two other kinds now draw break markers and are not counted,
  so the board can show you breaks while the sentence beneath says none. Fixing it means counting
  what the view drew rather than asking the flag — the flag itself must not grow, because a second
  definition of "broken" inside the renderer is exactly what the file's own comment forbids. Should
  I process this as a new task?

  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

  Why this is yours rather than mine: the fix is small, but it decides that the renderer starts
  keeping a tally of its own drawing decisions, which is a new kind of state in a file that has so
  far only consumed verdicts. That is a design call about where "what is broken" is allowed to be
  answered, and REQ-304's constraint was written to stop me making it unilaterally.
