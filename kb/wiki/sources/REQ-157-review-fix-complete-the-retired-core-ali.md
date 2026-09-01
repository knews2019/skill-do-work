---
title: "Lessons from REQ-157: Review fix: Complete the retired core alias guard"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-157-review-fix-complete-the-retired-core-ali.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-157: Review fix: Complete the retired core alias guard

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Make the shipped-surface recurrence guard cover every former moved-command trigger, not only canonical sibling action names and a small sample of natural-language aliases.

## Solution summary

Replaced the partial canonical-action-plus-samples recurrence guard with a fixture-driven historical contract covering 186 concrete trigger rows: 117 direct heads, 22 install targets across three former families, and three bare install heads. Added fixture-integrity checks, complete table-driven root/module mutations, branding/prose/current-command boundary controls, and historical/archive/fixture collector exclusions without changing runtime routes or shipped guidance.

## What worked

- Reconstructing the contract from the deleted router/shim and install-normalization history produced an auditable complete inventory instead of another sample list.
- Full-row root/module mutations plus a separate qualification pass caught both vocabulary gaps and an over-broad branding exemption before archival.

## What didn't work

- Line-wide exemptions are too coarse for a command-occurrence guard: one legitimate test reference can hide a second real invocation on the same line.
- Selecting only the longest prefix fails when that candidate misses its boundary but a shorter historical head or prefix route still applies.

## Worth knowing

- The 186-row historical vocabulary is complete and remains test-only. REQ-160 records the two remaining occurrence-completeness edge classes and is consent-gated rather than auto-runnable.

## Back-reference

See `do-work/archive/UR-031/REQ-157-complete-retired-core-alias-guard.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `1f7a245`.
