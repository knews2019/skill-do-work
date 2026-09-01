---
title: "Lessons from REQ-148: Addendum: preserve association candidates with empty quarantine"
type: source-summary
topic_cluster: metadata-and-timestamps
sources: [raw/processed/2026-09-01/REQ-148-addendum-preserve-association-candidates.md]
related:
  - page: REQ-128-secret-rename-quarantine-survives-re-inv
    rel: extends
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-148: Addendum: preserve association candidates with empty quarantine

Part of the [[concept-timestamp-and-metadata-governance]] cluster.

## What the REQ was about

Correct REQ-128's commit and unscoped-inspect quarantine merge so an empty run-level secret quarantine preserves every safe inventory candidate for REQ association.

## Solution summary

[MAP UNCHANGED] Commit and inspect still own their candidate filtering inline; only the first-input discriminator changed, with no new tool or interface.

## Worth knowing

- `NR == FNR` is not a safe first-file test when an input may be empty; identify that input explicitly through `FILENAME` and `ARGV`.
- A regression for a two-input merge must cover the zero-record first input, not only populated joins.

## Back-reference

See `do-work/archive/UR-034/REQ-148-preserve-empty-quarantine-association-candidates.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0ed7786`.
