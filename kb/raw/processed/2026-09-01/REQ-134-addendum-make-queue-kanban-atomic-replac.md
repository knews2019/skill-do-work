---
source_type: req_lesson
req_id: REQ-134
req_path: do-work/archive/UR-032/REQ-134-queue-kanban-atomic-replacement.md
date: 2026-08-07
domain: backend
module: tools/queue-kanban
tags: [backend, addendum, queue, kanban, atomic]
---

# Lessons from REQ-134: Addendum: make queue-kanban atomic replacement cross-platform and symlink-safe

## What the REQ was about

Correct the shared queue-kanban atomic replacement path so its crash-safety contract is valid on every supported platform and `next-version` cannot silently replace a symlink with a regular file.

## Solution summary

[MAP CHANGED] Queue-kanban complete-file writes now live in `atomic_write.go`; the final primitive is selected by the `atomic_replace_*` build-tag files and is shared by Testing frontmatter writes and `next-version`.

## What worked

- A cross-platform API name does not imply a cross-platform atomicity contract; the last filesystem mutation needs an OS-specific proof.
- Validate write-target identity with `Lstat` before temporary-file creation when replacing a pathname rather than following it.

## Back-reference

See `do-work/archive/UR-032/REQ-134-queue-kanban-atomic-replacement.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `7e0536a`.
