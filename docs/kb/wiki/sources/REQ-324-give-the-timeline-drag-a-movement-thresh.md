---
title: "Lessons from REQ-324: Give the timeline drag a movement threshold"
type: source-summary
topic_cluster: timeline-and-metrics
sources: [raw/processed/2026-09-04/REQ-324-give-the-timeline-drag-a-movement-thresh.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-324: Give the timeline drag a movement threshold

Part of the [[concept-duration-estimation-and-breaks]] cluster.

## What the REQ was about

Require a few pixels of movement before a press becomes a pan, so that clicking a row opens
its detail drawer and a hand tremor does not scroll the time axis.

## Solution summary

Three files. `web/board-timeline.js`, `generate_test.go` (the pure

## What worked

- **Reproduce the failure before believing the explanation, because the reproduction carries
- detail the explanation does not.** The REQ's diagnosis was right, but the RED run showed the
- click dying even when the pan clamped and the window did not move at all. Had the fix been
- keyed on "the window moved" rather than on "a render happened", it would have passed a
- narrower test and left the defect at *Fit all* — the zoom level the board opens on.
- **An assertion taken after the state it describes has been cleared measures nothing.** The
- grab-cursor check read `is-panning` after `pointerup`, which removes it, so it passed with
- the cursor restored to the pressed-immediately behaviour. The mutation caught it; the
- assertion had to move inside the press.

## Back-reference

See `do-work/archive/UR-065/REQ-324-give-the-timeline-drag-a-movement-threshold.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `3486ab2`.
