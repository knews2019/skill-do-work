---
title: "Lessons from REQ-321: Colour timeline bars by REQ status"
type: source-summary
topic_cluster: timeline-and-metrics
sources: [raw/processed/2026-09-04/REQ-321-colour-timeline-bars-by-req-status.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-321: Colour timeline bars by REQ status

Part of the [[concept-duration-estimation-and-breaks]] cluster.

## What the REQ was about

Give every timeline bar its REQ's status colour, using the same semantic tokens the board
cards and the Calendar chips already use.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (modified) — post-review, F3

## What worked

- Consuming the payload's `statusUnrecognized` verdict instead of re-deriving
- the vocabulary. And routing every swatch through the same custom property the bars read —
- which is what made the legend's contradiction (F4) a one-line fix rather than an audit.

## What didn't work

- **Opacity is one-dimensional, and it was asked to buy two separations at once.** Wait from
- work, and wait from the page. In the light palette those pull in opposite directions, so no
- single alpha exists that satisfies both — the first attempt was not a bad number, it was a
- missing channel. Two channels (an outline, and a per-theme alpha) settled it.
- **A test that asserts an ordering is not a test of a difference.** `waitOpacity <
- workOpacity` passes at 0.98 and at 0.26. The REQ's stated risk was legibility and the probe
- measured everything except legibility. Worse, the mutation I chose to "verify" it — restoring
- the pre-REQ fills — was precisely the one thing assertion 1 already greps for, which is
- REQ-293's lesson arriving one REQ later: **choose the mutation before looking at the
- assertion, or you will choose the one it already catches.**
- **And then the first fix re-created the hole.** Exempting the contrast check whenever a
- stroke was present would have let 0.98 through again. An escape hatch in an assertion needs
- its own condition to be measured — here, that the outline is visible *against the fill it
- **Manual evidence in both themes is not coverage of both themes.** The probe lane launches
- Chromium with no colour-scheme flag, so every existing browser test measures light. The dark
- palette is this board's `:root` base and nothing automated had ever looked at it.

## Worth knowing

- the surface behind a timeline bar is `<body>` — nothing between them
- paints — so a contrast measurement must read `document.body`'s background, and the real values
- are `rgb(245,247,250)` light and `rgb(12,14,18)` dark, not the `--surface-*` tokens you would
- reach for. Computing against the wrong surface gave numbers that disagreed with the review's
- by a factor of two, which is how the discrepancy was found.

## Back-reference

See `do-work/archive/UR-065/REQ-321-colour-timeline-bars-by-req-status.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `d5fc642`.
