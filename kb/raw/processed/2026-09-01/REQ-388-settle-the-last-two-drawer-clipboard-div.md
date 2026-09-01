---
source_type: req_lesson
req_id: REQ-388
req_path: do-work/archive/UR-075/REQ-388-settle-the-last-two-drawer-clipboard-divergences.md
date: 2026-08-27
domain: frontend
module: _dev/primes
tags: [frontend, settle, last, drawer, clipboard]
---

# Lessons from REQ-388: Settle the last two drawer/clipboard divergences: fence info strings and ids inside paths

## What the REQ was about

Two places where the drawer's glossary and the paste's appendix list different ids for the same body.
Decide which surface is right in each case and make them agree.

## Solution summary

The drawer and copied appendix now list the same external references for fence metadata and file-path cases. Static boards no longer turn part of a file path into a ticket link, while live file navigation remains available.

## What worked

When two projections differ, compare their final reference lists over the same source and exercise both static and live production paths. Merely finding the same regex candidates misses a drawer retry that exists only after a path fails to become a link. Keep annotation suppression independent from source-search citation collection.

## Back-reference

See `do-work/archive/UR-075/REQ-388-settle-the-last-two-drawer-clipboard-divergences.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `3ed11c17`.
