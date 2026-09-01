---
title: "Lessons from REQ-153: Review fix: Sweep retired core command restatements"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-153-review-fix-sweep-retired-core-command-re.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-153: Review fix: Sweep retired core command restatements

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Remove or update every live restatement of the retired core moved-command and transition-era updater contracts, and add a shipped-surface guard so this class cannot recur.

## Solution summary

Completed the root-cause sweep across every live occurrence found by exploration and made recurrence part of the staged distribution contract. Historical evidence, generic pipeline prose, branding, and legitimate core commands remain untouched.

## What worked

- Scanning the live modular tree found stale commands in source comments and user guides beyond the original five UI/hook/prime instances.
- Reusing the existing ownership table kept canonical sibling routes and single-owner validation aligned without duplicating package mappings.

## What didn't work

- Deriving only canonical action names plus a small hand-added alias set repaired the present corpus but did not cover the full deleted router vocabulary; a class-level guard needs a complete historical trigger contract.

## Worth knowing

- Current live surfaces and prime labels are clean, and the guard is intentionally narrow around history/branding/generic pipeline prose. REQ-157 is held for consent to expand alias completeness without republishing those aliases as guidance.

## Back-reference

See `do-work/archive/UR-031/REQ-153-sweep-retired-core-command-restatements.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `157b89e`.
