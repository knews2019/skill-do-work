---
source_type: req_lesson
req_id: REQ-147
req_path: do-work/archive/UR-033/REQ-147-reserve-request-numbers-during-allocation.md
date: 2026-08-07
domain: backend
module: tools/queue-kanban
tags: [backend, addendum, reserve, request, numbers]
---

# Lessons from REQ-147: Addendum: reserve request numbers during allocation

## What the REQ was about

Change `queue-kanban next-req` from a read-only maximum scan into an atomic reservation operation so every successful invocation owns a distinct REQ number before it prints that number.

## Solution summary

[MAP CHANGED] `queue-kanban next-req` is now the third write surface. Its only mutation is an empty marker under `do-work/.req-reservations/`, which capture retains and stages with the matching REQ/UR records.

## What worked

- A read-only max scan cannot allocate an identifier; allocation needs a durable ownership event before output.
- Per-id exclusive markers avoid stale-lock recovery entirely: a stopped caller creates a safe gap, not a lock that another process must guess how to break.
- Path containment should survive the interval after validation, so rooted filesystem handles are stronger than joining an already-checked absolute path.

## Back-reference

See `do-work/archive/UR-033/REQ-147-reserve-request-numbers-during-allocation.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `4851438`.
