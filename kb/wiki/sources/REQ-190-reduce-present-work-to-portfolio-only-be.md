---
title: "Lessons from REQ-190: Reduce present-work to portfolio-only behavior"
type: source-summary
topic_cluster: presentation-and-reporting
sources: [raw/processed/2026-09-01/REQ-190-reduce-present-work-to-portfolio-only-be.md]
related:
  - page: concept-completed-work-presentation
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-190: Reduce present-work to portfolio-only behavior

Part of the [[concept-completed-work-presentation]] cluster.

## What the REQ was about

Make `do-work-toolbox present-work all|portfolio` a portfolio aggregation command only. Remove the competing per-item detail, brief, explainer, and Remotion workflows, and give bare or item-specific invocations compact non-writing guidance instead of silently delegating.

## Solution summary

Reduced `present-work` from a mixed item/portfolio artifact generator to one cross-project portfolio command. Only `all|portfolio` may read archives, prompt, or write; bare and item-specific invocations now stop after compact guidance, while portfolio output uses one evidence-backed draft and optional byte-identical no-clobber snapshot.

## What worked

- Deleting the entire detail-mode span before rebuilding exposed a small, auditable dispatcher and one portfolio workflow instead of leaving hidden per-item fallthroughs.
- Keeping snapshot publication as a prose contract avoided unsafe shell state across an interactive prompt while still making byte identity and no-clobber behavior explicit.

## What didn't work

- The first dispatcher described canonical-looking UR/REQ tokens without inheriting the suite-wide case-insensitive, numeric-value token grammar.
- Resolving a collision-safe snapshot name was not enough; the first publication wording refreshed the canonical file before exclusive snapshot success, permitting partial completion.

## Back-reference

See `do-work/archive/UR-042/REQ-190-reduce-present-work-to-portfolio-only.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `c66d11c`.
