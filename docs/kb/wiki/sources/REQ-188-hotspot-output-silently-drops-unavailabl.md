---
title: "Lessons from REQ-188: Hotspot output silently drops unavailable tracked paths"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-188-hotspot-output-silently-drops-unavailabl.md]
related:
  - page: REQ-182-public-work-and-schema-vocabularies-drif
    rel: complements
  - page: REQ-184-live-board-origin-checks-have-no-trusted
    rel: complements
  - page: REQ-185-javascript-behavior-probes-can-all-skip-
    rel: complements
  - page: REQ-186-required-baseline-verification-executes-
    rel: complements
  - page: REQ-187-no-single-local-maintainer-command-prove
    rel: complements
created: 2026-09-01
updated: 2026-09-02
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

## Worth knowing

`topCount` applies only to numeric hotspot entries. Every churn-bearing unavailable path remains evidence, keeps its known commit count, and renders uncapped in deterministic path order.

## Back-reference

See `do-work/archive/UR-041/REQ-188-hotspot-unavailable-evidence.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8d63070`.
