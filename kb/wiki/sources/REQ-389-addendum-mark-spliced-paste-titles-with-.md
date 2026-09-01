---
title: "Lessons from REQ-389: Addendum: mark spliced paste titles with a leading arrow"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-01/REQ-389-addendum-mark-spliced-paste-titles-with-.md]
related:
  - page: REQ-387-keep-a-spliced-title-from-changing-how-t
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-389: Addendum: mark spliced paste titles with a leading arrow

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

The Copy buttons splice a ticket's title after the first body mention of its id, as
`REQ-374 (Show how long each done card took)`. Change the spliced form to
`REQ-374 (-> Show how long each done card took)` so a reader of the paste can tell the
parenthetical was inserted by the board, not written by the ticket's author.

## Solution summary

- `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js` (modified). Changes the single in-body insertion prefix to the exact ASCII arrow and space. Safe-title escaping, offsets and the appendix are unchanged. One line replaced.
- `skills/do-work-board/tools/queue-kanban/generate_test.go`

## Worth knowing

No new reusable lesson beyond the existing clipboard representation and actual-renderer lessons. This change deliberately distinguishes inserted text from author prose; do not later remove that marker to make clipboard text identical to drawer styling.

## Back-reference

See `do-work/archive/UR-078/REQ-389-mark-spliced-paste-titles-with-a-leading-arrow.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `4ed31496`.
