---
source_type: req_lesson
req_id: REQ-323
req_path: do-work/archive/UR-065/REQ-323-let-a-timeline-bar-be-read-against-a-date.md
date: 2026-08-23
domain: frontend
module: _dev/primes
tags: [frontend, timeline, read, against]
---

# Lessons from REQ-323: Let a timeline bar be read against a date

## What the REQ was about

Three additions that make a bar's position mean something: gridlines through the plot, a
drawn queue-end rule, and a minimum bar that stays readable when the window is wide.

## Solution summary

Five files. `web/board-timeline.js` (the three parts),

## What worked

- **A prominence metric that ignores a channel the design uses will contradict the render.**
- The first version of the vertical-rule probe measured contrast times stroke width and
- reported the light-theme forecast rule as louder than the now-rule — the opposite of what the
- screenshot showed. Dash duty was the missing channel: a `2 4` rule inks a third of its length
- and a `3 3` rule inks half. The fix was the metric, not the design, and the way to tell which
- is which was having the screenshot first.
- **A guard that no mutation can break is either dead or untested — check which before
- assuming the second.** Removing the `isFinite` branch from the collapse function changed no
- test result. The instinct is to write a test that reaches it; the truth was that nothing
- can, because the sentinel segment is only ever emitted alone.

## Back-reference

See `do-work/archive/UR-065/REQ-323-let-a-timeline-bar-be-read-against-a-date.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `7adaa2e`.
