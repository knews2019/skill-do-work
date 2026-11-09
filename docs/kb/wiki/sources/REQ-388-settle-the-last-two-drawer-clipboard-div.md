---
title: "Lessons from REQ-388: Settle the last two drawer/clipboard divergences: fence info strings and ids inside paths"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-388-settle-the-last-two-drawer-clipboard-div.md]
related:
  - page: REQ-382-expand-ticket-ids-written-as-markdown-li
    rel: complements
  - page: REQ-386-make-the-drawer-and-the-paste-agree-abou
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-388: Settle the last two drawer/clipboard divergences: fence info strings and ids inside paths

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

Two places where the drawer's glossary and the paste's appendix list different ids for the same body.
Decide which surface is right in each case and make them agree.

## Solution summary

The drawer and copied appendix now list the same external references for fence metadata and file-path cases. Static boards no longer turn part of a file path into a ticket link, while live file navigation remains available.

## Worth knowing

When two projections differ, compare their final reference lists over the same source and exercise both static and live production paths. Merely finding the same regex candidates misses a drawer retry that exists only after a path fails to become a link. Keep annotation suppression independent from source-search citation collection.

## Back-reference

See `do-work/archive/UR-075/REQ-388-settle-the-last-two-drawer-clipboard-divergences.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `3ed11c17`.
