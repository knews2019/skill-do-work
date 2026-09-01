---
title: "Lessons from REQ-138: Add Managed Text-Section Replacement"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-138-add-managed-text-section-replacement.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-138: Add Managed Text-Section Replacement

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Add a deterministic managed text-section utility and use it to own only do-work's generated Just recipes.

## Solution summary

Established the repository-owned `do-work:recipes` section and one reusable reconciler that can create or update it without reformatting the client's Justfile. The old per-recipe extraction instructions are gone; malformed or ambiguous ownership always fails before target mutation.

## What worked

- Treating the client file as bytes and ownership as one inclusive marker span made the preservation guarantee directly testable, including NUL exterior data and missing filename assumptions.
- A candidate copy provides the consent diff without requiring a dry-run mode in the mutation utility.

## What didn't work

- A first-pass legacy span would have absorbed any custom top-level recipe interleaved between the five old recipes. Rejecting ambiguous legacy structure is the only safe behavior when the old block has no end marker.

## Worth knowing

- The section file itself must be exactly one newline-terminated managed span; absent-target templates must embed those exact bytes once.
- Repeated current-content reconciliation performs no rename, preserving byte identity without unnecessary inode churn.

## Back-reference

See `do-work/archive/UR-031/REQ-138-add-managed-text-section-replacement.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `aabc363`.
