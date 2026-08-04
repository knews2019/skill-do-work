---
source_type: req_lesson
req_id: REQ-016
req_path: do-work/archive/UR-003/REQ-016-remove-severity-dead-field.md
date: 2026-07-01
domain: backend
module: tools/queue-kanban
tags: [backend, queue-kanban, producer-less, severity, frontmatter]
---

# Lessons from REQ-016: Remove the producer-less `severity` frontmatter field from queue-kanban

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

See `do-work/archive/UR-003/REQ-016-remove-severity-dead-field.md` for the full REQ — triage, implementation, review, and lessons. Commit `023aa50`.
