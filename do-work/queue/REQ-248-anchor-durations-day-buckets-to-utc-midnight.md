---
id: REQ-248
title: Anchor the Durations day buckets to UTC midnight so Panel B stays on canvas
status: pending
created_at: 2026-08-18T13:54:59Z
status_changed_at: 2026-08-18T13:54:59Z
user_request: UR-051
addendum_to: REQ-242
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Anchor the Durations Day Buckets to UTC Midnight So Panel B Stays on Canvas

## What

Panel B's bars are placed with `xOfEpoch`, which maps each day bucket's **midnight**, while `timeStart` is the **first completion instant**. The two disagree by however far into its first day the earliest sample falls, so the leftmost bar is drawn to the left of the plot area — and on a board with one or two active days the disagreement dominates the whole span and the panel renders off-canvas entirely.

## Context

Found by REQ-242's builder as an unrelated pre-existing quirk, then confirmed and extended by REQ-242's independent review. This is not cosmetic at low day counts.

## Instances

- [ ] **Leftmost bar sits in the axis gutter on the real board.** `x=37.1 width=12` spans 37.1–49.1, entirely left of `DURATIONS_MARGIN_LEFT` (54), and the render shows it struck through by the "0" axis tick. Visible on this repository's own board today.
- [ ] **One active day: Panel B renders empty.** `timeSpan` collapses to the intra-day sample span (about 3 hours), so `xOfEpoch(midnight)` maps to roughly minus three plot-widths — measured annotation at `x=-3330`, bar at `x=-3342`. Both completely off-canvas.
- [ ] **Two active days: same failure, smaller magnitude.** Annotation at `x=-336.5`, bar at `x=-348.5`.

## Requirements

- Every Panel B bar renders inside the plot area at every day count, including one and two active days.
- The slowest-day annotation renders on canvas at every day count — it exists to state a value a clipped bar cannot, and cannot do that from off-screen.
- No change to `DURATIONS_MEDIAN_TITLE_Y` or `describeAtPointer`'s A/B boundary.
- REQ-241's and REQ-242's guarantees hold unchanged: 0 same-row label overlaps, 0 label/mark overlaps, the annotation clear of every neighbour in its strip.

## Builder Guidance

The suggested root fix from the review is to floor `timeStart` to its UTC midnight and ceil `timeEnd` to the following midnight before computing `timeSpan`, so the axis domain and the day buckets share one origin. Verify that against the other panels before adopting it — Panels A and C read the same domain.

**Generate a board and look at it**, at one, two and many active days. Measure in the live DOM.

## Red-Green Proof

**RED prompt/case:** a test asserting every Panel B bar's x-range and the annotation's x-range fall inside the plot area, evaluated on one-day, two-day and many-day fixtures.
**Why RED now:** measured `x=-3330` on a one-day board and `x=37.1` against a left margin of 54 on the real board.
**GREEN when:** the assertion passes at every day count and a render at one and two days shows Panel B populated.
**Validation:** Review finding on REQ-242; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
