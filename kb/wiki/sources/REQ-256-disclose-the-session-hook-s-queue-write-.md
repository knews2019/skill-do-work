---
title: "Lessons from REQ-256: Disclose the session hook's queue write surface in the docs"
type: source-summary
topic_cluster: checkpoint-and-crash-recovery
sources: [raw/processed/2026-09-01/REQ-256-disclose-the-session-hook-s-queue-write-.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-256: Disclose the session hook's queue write surface in the docs

Part of the [[concept-session-checkpoints-and-recovery]] cluster.

## What the REQ was about

REQ-246 made the SessionStart hook a *write* surface on consumer queue files — it mechanically repairs detectably wrong `*_at` stamps in `do-work/queue/` and `do-work/working/` at session start. Two shipped texts still describe the hook as read-only-plus-banner; a consumer auditing "what writes to my repo at session start" is misled.

## Solution summary

Both doc sites now disclose the SessionStart hook's queue write surface. README's hook bullet says session-start.sh also writes to consumer queue files — reaping stale reservation markers and mechanically repairing detectably wrong `*_at` stamps — citing both scripts root-relatively. capture.md gained one paragraph beside the Immutability Rule's timestamp-repair exception, framing the session-start repair as the same metadata-correction class with an explicit "never archive". A class sweep (grep for session-start/SessionStart over shipped markdown) confirmed the two instances were the whole site class.

## Worth knowing

Route A, no surprises — skipped per the rule (the class sweep confirming the REQ's instance list was complete is already recorded in P-A-U).

## Back-reference

See `do-work/archive/UR-056/REQ-256-disclose-the-session-hooks-queue-write-surface-in-the-docs.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `fbc14e8`.
