---
source_type: req_lesson
req_id: REQ-188
req_path: do-work/archive/UR-041/REQ-188-hotspot-unavailable-evidence.md
date: 2026-08-15
domain: backend
module: skills/do-work-toolbox/tools/audit-metrics
tags: [backend, hotspot, output, silently, drops]
---

# Lessons from REQ-188: Hotspot output silently drops unavailable tracked paths

## What the REQ was about

Keep unreadable or otherwise unavailable tracked paths visible in hotspot output as `NOT-MEASURED`, while preserving valid measured rows and warning that the ranking is incomplete.

## Solution summary

**[MAP CHANGED]** Hotspot reporting now has an explicit completeness channel: numeric churn × size rankings stay capped, while every churn-bearing path unavailable in the current worktree remains visible as uncapped `NOT-MEASURED` evidence.

## What worked

- Partitioning measured rankings from unavailable evidence keeps arithmetic honest while making completeness visible.
- Removing regular tracked files after commit is a portable real-Git fixture for the same read boundary as a broken symlink.

## What didn't work

- Treating a per-path measurement failure as a harmless `continue` produced a plausible but incomplete report.
- A one-measured-row fixture cannot mutation-lock numeric ordering or capping even when it proves unavailable rows bypass the cap.

## Back-reference

See `do-work/archive/UR-041/REQ-188-hotspot-unavailable-evidence.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8d63070`.
