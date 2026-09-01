---
title: "Lessons from REQ-188: Hotspot output silently drops unavailable tracked paths"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-188-hotspot-output-silently-drops-unavailabl.md]
related:
  - page: concept-contract-verification-gates
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-188: Hotspot output silently drops unavailable tracked paths

Part of the [[concept-contract-verification-gates]] cluster.

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
