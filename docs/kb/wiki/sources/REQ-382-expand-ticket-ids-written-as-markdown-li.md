---
title: "Lessons from REQ-382: Expand ticket ids written as Markdown links"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-382-expand-ticket-ids-written-as-markdown-li.md]
related:
  - page: REQ-381-index-cited-ticket-ids-and-let-the-filte
    rel: complements
  - page: REQ-387-keep-a-spliced-title-from-changing-how-t
    rel: complements
  - page: REQ-388-settle-the-last-two-drawer-clipboard-div
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-382: Expand ticket ids written as Markdown links

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

A REQ body that writes a ticket as an explicit Markdown link — `[REQ-123](https://…)` — gets neither
a title nor a glossary entry. `linkifyDetailBody` skips any text node already inside an `<a>`, so the
one place an author has gone out of their way to mark a reference is the one place REQ-378 does not
reach.

## Solution summary

- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modified). Reuses existing ticket resolution, title shortening and glossary accounting for text inside authored anchors. Adds inert title-bearing spans with non-navigation identity metadata, skips renderer-shaped autolinks and unknown

## Worth knowing

- Skipping generated DOM on a second pass prevents nesting but does not preserve first-mention/glossary memory. Reconstruct state in document order from the original mention identity while refusing to scan inserted title text; a drawer-root cache would become stale when the next ticket replaces its HTML.

## Back-reference

See `do-work/archive/UR-075/REQ-382-expand-ticket-ids-written-as-markdown-links.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `59caf025`.
