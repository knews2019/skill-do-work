---
title: "Timeline and Metrics"
type: concept
topic_cluster: timeline-and-metrics
sources:
  - raw/processed/2026-09-04/REQ-227-add-the-timeline-view-with-two-segment-r.md
  - raw/processed/2026-09-04/REQ-233-give-the-timeline-a-keyboard-path-to-zoo.md
  - raw/processed/2026-09-04/REQ-319-list-only-the-reqs-the-selected-window-c.md
  - raw/processed/2026-09-04/REQ-320-show-and-set-the-timeline-window-s-start.md
  - raw/processed/2026-09-04/REQ-321-colour-timeline-bars-by-req-status.md
  - raw/processed/2026-09-04/REQ-322-name-the-req-on-its-own-timeline-row.md
  - raw/processed/2026-09-04/REQ-323-let-a-timeline-bar-be-read-against-a-dat.md
  - raw/processed/2026-09-04/REQ-324-give-the-timeline-drag-a-movement-thresh.md
  - raw/processed/2026-09-04/REQ-459-review-fix-stage-command-owned-calibrati.md
  - raw/processed/2026-09-01/REQ-248-anchor-the-durations-day-buckets-to-utc-.md
  - raw/processed/2026-09-01/REQ-252-record-the-browser-with-every-measured-f.md
  - raw/processed/2026-09-01/REQ-292-move-durations-label-placement-into-the-.md
  - raw/processed/2026-09-01/REQ-304-draw-a-reversed-wait-as-a-break-not-as-a.md
  - raw/processed/2026-09-01/REQ-305-say-the-timeline-forecast-describes-the-.md
  - raw/processed/2026-09-01/REQ-313-count-the-breaks-the-timeline-actually-d.md
  - raw/processed/2026-09-01/REQ-336-timeline-clicks-open-the-detail-drawer-a.md
  - raw/processed/2026-09-01/REQ-337-a-check-that-can-catch-timeline-click-re.md
  - raw/processed/2026-09-01/REQ-338-cut-the-timeline-row-list-to-one-tab-sto.md
related:
  - page: concept-kanban-board-architecture
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: high
---

# Timeline and Metrics

Architectural overview and synthesis for the Timeline and Metrics subsystem in the do-work suite.

## Key Principles & Synthesized Lessons

This cluster synthesizes evidence from 18 source documents:

- [[REQ-248-anchor-the-durations-day-buckets-to-utc-]] — Anchor the Durations day buckets to UTC midnight so Panel B stays on canvas
- [[REQ-252-record-the-browser-with-every-measured-f]] — Record the browser with every measured-face number in the Durations tests
- [[REQ-292-move-durations-label-placement-into-the-]] — Move Durations label placement into the browser and delete the measured-face constants
- [[REQ-304-draw-a-reversed-wait-as-a-break-not-as-a]] — Draw a reversed wait as a break, not as a valid bar
- [[REQ-305-say-the-timeline-forecast-describes-the-]] — Say the timeline forecast describes the whole queue, not the filtered rows
- [[REQ-313-count-the-breaks-the-timeline-actually-d]] — Count the breaks the timeline actually draws
- [[REQ-336-timeline-clicks-open-the-detail-drawer-a]] — Timeline clicks open the detail drawer again
- [[REQ-337-a-check-that-can-catch-timeline-click-re]] — A check that can catch Timeline click retargeting
- [[REQ-338-cut-the-timeline-row-list-to-one-tab-sto]] — Cut the Timeline row list to one Tab stop
- [[REQ-227-add-the-timeline-view-with-two-segment-r]] — Add the Timeline view with two-segment REQ bars
- [[REQ-233-give-the-timeline-a-keyboard-path-to-zoo]] — Give the Timeline a keyboard path to zoom and pan
- [[REQ-319-list-only-the-reqs-the-selected-window-c]] — List only the REQs the selected window covers
- [[REQ-320-show-and-set-the-timeline-window-s-start]] — Show and set the timeline window's start and end
- [[REQ-321-colour-timeline-bars-by-req-status]] — Colour timeline bars by REQ status
- [[REQ-322-name-the-req-on-its-own-timeline-row]] — Name the REQ on its own timeline row
- [[REQ-323-let-a-timeline-bar-be-read-against-a-dat]] — Let a timeline bar be read against a date
- [[REQ-324-give-the-timeline-drag-a-movement-thresh]] — Give the timeline drag a movement threshold
- [[REQ-459-review-fix-stage-command-owned-calibrati]] — Review fix: Stage command-owned calibration with lifecycle release

## Cross-References

See related system components and verification gates.
