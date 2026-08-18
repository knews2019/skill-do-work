---
id: REQ-240
title: Stop the Timeline axis printing a fake minute
status: claimed
claimed_at: 2026-08-18T11:42:03Z
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T11:42:03Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - browser evidence
created_at: 2026-08-18T11:37:10Z
user_request: UR-052
addendum_to: REQ-235
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-timeline.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Stop the Timeline Axis Printing a Fake Minute

## What

`timelineFormatAxisTick` in `web/board-timeline.js` formats any window of three days or less as `"<day> <Mon> <HH>:00"` — where `:00` is a **string literal**, not the tick's actual minute. So a tick at 11:55 prints `11:00`. Once the window is short enough that several ticks fall inside one hour, they all print the same label.

Measured on the merged tree, immediately after clicking `Now`: **seven ticks, two distinct labels** — `18 Aug 11:00` five times, then `18 Aug 12:00` twice.

## Context

Found by REQ-235's review, by rendering the board and looking at the axis rather than by any assertion — the ticks are correctly *positioned*, so nothing in the suite notices that they are incorrectly *labelled*.

The defect is REQ-227's, not REQ-235's: the literal has been there since the view was built. What REQ-235 changed is how often you meet it. Before, reaching a sub-hour window took deliberate zooming. Now the `Now` button lands you in one by design — it sizes the window to cover the now-line and the forecast, which on a healthy queue is well under an hour — and the new `Day` level is one step away from it. The single most-used new control in REQ-235 lands the reader on an axis that reads as five identical labels.

That is also why it matters more than its size suggests: UR-052's complaint was "I can not jump to the remaining work", REQ-235 built the jump, and the jump's landing state is the one place the axis is least readable.

## Requirements

- An axis tick's label states the tick's real instant. No component of the label may be a literal that the instant does not carry.
- At any window span, two ticks at different instants must not render identical labels — that is the property the current code violates, and it is what makes the axis unreadable rather than merely imprecise.
- Tick labels stay legible at the existing tick count and font size; a longer format that overlaps its neighbour trades one unreadable axis for another.
- No change to tick *positions* — this is a formatting defect, not a layout one.
- The day/week/month period windows and the `Fit all` window keep the labels they render today, where those are already correct.

## Red-Green Proof

**RED prompt/case:** a Node behaviour probe driving `timelineFormatAxisTick` over the tick instants of a sub-hour window (the shape `Now` produces), asserting that the rendered labels are pairwise distinct and that a tick at a non-zero minute does not render `:00`.
**Why RED now:** the `:00` is a literal, so a 1-hour window's ticks collapse to at most two distinct labels; measured at 7 ticks → 2 distinct on the live board.
**GREEN when:** the probe passes and a rendered board, after clicking `Now`, shows an axis whose labels are all different.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Notes for the builder

Two adjacent observations from the same review, both recorded on REQ-235 and neither in this REQ's scope — do not fix them here, but know they exist so you do not design against them:

- A week window's six evenly-spaced ticks land at 1.167-day intervals, so interior labels skip a day (`7, 8, 9, 10, 11` with 12 missing). A period-aware tick *generator* would fix that and this together, but it is an axis redesign and wants its own decision.
- `Fit all` spanning the whole capture history is what crushes the bars into the right-hand fraction of the plot in the first place.
