---
title: "Lessons from REQ-016: Remove the producer-less `severity` frontmatter field from queue-kanban"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-016-remove-the-producer-less-severity-frontm.md]
related:
  - page: concept-kanban-board-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-016: Remove the producer-less `severity` frontmatter field from queue-kanban

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

Remove the `severity` frontmatter pipeline from the queue-kanban tool. No REQ schema in this repo ever emits a top-level `severity:` key (the Schema Read Contract in `actions/work-reference.md` doesn't define one; discovered-task severity lives as an inline `[critical]`/`[normal]`/`[low]` bullet prefix inside `## Discovered Tasks`, never as frontmatter), yet the tool carries a full parse → JSON → badge pipeline for it.

## Solution summary

Removed the entire producer-less `severity` frontmatter vertical (Go struct/parse → JSON export → JS badge/drawer render → CSS) from the queue-kanban tool. Sweep confirmed no other sites existed (tests included); the shared `makeBadge` helper stays, since domain/ur/route badges use it.

## What worked

- Enumerating the full vertical (parse → export → render → style) in the REQ up front made the deletion mechanical; verifying the render against the repo's real `do-work/` tree (not just unit tests) proved the board still works end to end.

## What didn't work

- Nothing — no dead ends.

## Worth knowing

- The neighboring `batch` frontmatter field looks similar but is NOT a dead vertical — it has real producers in archived REQs (UR-002's REQ-013/REQ-014 frontmatter), so don't sweep it up in a future "same shape" cleanup without re-checking. When grepping the generated `board-data.js` for leftover fields, match the JSON key (`"severity":`) — the bare word appears legitimately in rendered REQ body prose.

## Back-reference

See `do-work/archive/UR-003/REQ-016-remove-severity-dead-field.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `023aa50`.
