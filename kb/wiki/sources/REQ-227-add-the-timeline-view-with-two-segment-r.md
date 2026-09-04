---
title: "Lessons from REQ-227: Add the Timeline view with two-segment REQ bars"
type: source-summary
topic_cluster: timeline-and-metrics
sources: [raw/processed/2026-09-04/REQ-227-add-the-timeline-view-with-two-segment-r.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-227: Add the Timeline view with two-segment REQ bars

Part of the [[concept-duration-estimation-and-breaks]] cluster.

## What the REQ was about

Add a fifth board view, **Timeline** — a zoomable, scrollable Gantt with one horizontal bar per REQ.
Each bar carries two coloured segments: `created_at`→`claimed_at` (the wait) and
`claimed_at`→`completed_at` (the work). REQs that are claimed but not finished draw as open bars
running to now.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/timeline.go` (new)
- `skills/do-work-board/tools/queue-kanban/timeline_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (new)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-filters.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)

## What worked

- Extracting the two rules that could be silently wrong — the zoom transform and the visible-row slice — as pure functions before writing any DOM code. Both are things no screenshot catches: a zoom that drifts the anchor looks fine in a still, and a slice that is subtly wrong shows blank strips only while scrolling fast. As pure functions they became two Node probes that fail loudly under mutation. Choosing pixel coordinates over a `viewBox` was the other one: the REQ warned that Durations' fixed conversion is exactly what a zoom invalidates, and removing the conversion is a stronger answer than maintaining it.

## What didn't work

- Assuming `:root` was the light palette. It is the dark default here, with light under `@media (prefers-color-scheme: light)`, so the first render had the wait segment at full navy weight on white — the opposite of the "wait is the quieter hue" intent stated in its own comment.
- An absolutely-positioned now-line. It sat inside the scroll container, so the render's own `textContent = ""` erased it, and it would still have needed the container's padding folded into its `left` by hand. One node inside the rows SVG has neither problem.
- Copying Durations' listener habit. Durations attaches to nodes it rebuilds every render, so its handlers die with the old DOM; this view attaches to the scroll container and to `window`, which both outlive a render — and a filter change re-renders it. Five filter changes left six live scroll handlers, and every later scroll re-rendered the rows six times.

## Worth knowing

- The pattern behind the last two is the same. A view that re-renders into fresh nodes can be careless about cleanup; a view that binds to anything persistent cannot, and "the neighbour does it this way" is not evidence that the neighbour's constraints are yours. The board now has one of each, so the next view added here should decide which it is before copying either.

## Back-reference

See `do-work/archive/UR-051/REQ-227-timeline-view-with-two-segment-req-bars.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `17b9422`.
