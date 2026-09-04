---
title: "Lessons from REQ-323: Let a timeline bar be read against a date"
type: source-summary
topic_cluster: timeline-and-metrics
sources: [raw/processed/2026-09-04/REQ-323-let-a-timeline-bar-be-read-against-a-dat.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-323: Let a timeline bar be read against a date

Part of the [[concept-duration-estimation-and-breaks]] cluster.

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
