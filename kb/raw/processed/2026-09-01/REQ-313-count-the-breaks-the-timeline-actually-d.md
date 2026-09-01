---
source_type: req_lesson
req_id: REQ-313
req_path: do-work/archive/UR-062/REQ-313-count-the-breaks-the-timeline-actually-draws.md
date: 2026-08-21
domain: frontend
module: _dev/primes
tags: [frontend, count, breaks, timeline, actually]
---

# Lessons from REQ-313: Count the breaks the timeline actually draws

## What the REQ was about

The Timeline view's summary line renders "N with broken stamps, drawn as breaks"
(`web/board-timeline.js:546-557`), and N counts rows where `row.anomaly` is true. That flag
carries `detectCompletionAnomaly`'s verdict, which `model.go:246-252` scopes to **completion
bookkeeping** — no completion instant at all on a terminal REQ. It is false for both spans the view
draws as breaks:

- a reversed **work** span, whose `completed_at` parses fine but precedes `claimed_at`
- a reversed **wait** span, which REQ-304 just started drawing as a break

## Solution summary

The Timeline summary now counts unique rows in the already-filtered population
when any of the renderer's existing break causes applies: completion anomaly, reversed wait, or
drawable reversed work. A real-renderer DOM test covers every cause, multi-cause deduplication,
filtered counts, and the healthy zero-clause behavior without changing the row model or SVG geometry.

## What worked

- Counting with one boolean predicate over the already-filtered rows makes deduplication automatic
  and makes the summary inherit the same population as every neighboring number.
- A real-renderer caller-seam test caught anomaly-only, global-population, and cause-count mutations
  without introducing a second implementation of timeline filtering.

## What didn't work

- The test fixture proves summary semantics and the adjacent REQ-304 test proves marker emission, but
  it does not join those two observations by row ID. That leaves the `row.hasWork` guard defended by
  code review and production payload shape rather than a single end-to-end assertion.

## Back-reference

See `do-work/archive/UR-062/REQ-313-count-the-breaks-the-timeline-actually-draws.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0761a10`.
