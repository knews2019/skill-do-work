---
title: "Lessons from REQ-242: Stop Panel B's slowest-day annotation colliding with its own title"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-04/REQ-242-stop-panel-b-s-slowest-day-annotation-co.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-242: Stop Panel B's slowest-day annotation colliding with its own title

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

In the Durations view, Panel B's slowest-day annotation is drawn at `y = 355` while Panel B's own title sits at `y = 350` (`DURATIONS_MEDIAN_TITLE_Y`). The two overlap: on a synthetic fixture the annotation `209 min` renders directly through the words "paused and broken spans excluded".

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

## What worked

- **A chart's empty space is only empty at the x you looked at.** The first fix passed a test suite that included a fresh, purpose-built geometric assertion — and the render showed it printing through panel B's `0` axis tick. The tick lives in the y-axis gutter, so it is invisible as a neighbour at every x except the extreme left, which is precisely the luck-of-x this REQ existed to remove. Reproducing that failure twice in one REQ is the strongest case yet for `prime-kanban-board.md`'s "generate a board and look at it": the second collision was found the same way as the first, and neither was reachable by reasoning over the constants the fix was about.
- **When a defect is "invisible because of where the data happens to fall", the fix is to remove the dependence, not to widen the margin.** A larger gap above the title would still have been a clearance that held for some slowest days and not others.
- **A measured face is per-Chromium.** The 11px label face measured 10.4278 ascent for REQ-241 and 10.1853 here on a different build. Both constants round up and away from the model, so a test written that way survives the difference — but a test asserting an exact measured number would not.

## Back-reference

See `do-work/archive/UR-051/REQ-242-stop-panel-b-annotation-colliding-with-its-title.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `48263dd`.
