---
title: "Lessons from REQ-158: Review fix: Complete rendered-region classification in shipped Markdown references"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-158-review-fix-complete-rendered-region-clas.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-158: Review fix: Complete rendered-region classification in shipped Markdown references

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Make one rendered-versus-ignored decision govern every target-discovery path in the shipped-package Markdown guard for the already-approved escaped-link, indented-code, and inline-code classes. Done means these classes cannot recur as either release-blocking false positives or broken-link false negatives; this does not add a documentation policy or broaden the downstream reference contract.

## Solution summary

The shipped Markdown guard now feeds structural links, reference definitions, and bare first-party URLs from the same rendered text while retaining the existing downstream publication policy and source-offset invariants.

## What worked

- Exact production-helper probes made disagreement between structural extraction and the bare-URL fallback observable without changing downstream policy.
- Length-preserving masks plus paired live/hidden controls protected offsets and caught over-masking during re-qualification.

## What didn't work

- Implementing only the reproduced opening-bracket and top-level paragraph forms left the same classification rule incomplete at other link delimiters and list-item paragraph contexts; REQ-161 records the remaining sweep for consent.

## Worth knowing

- A rendered-region guard needs delimiter-complete and block-context-complete tables, not one fixture per syntax family. Full live-corpus contracts can stay green while nearby production-helper variants still fail.

## Back-reference

See `do-work/archive/UR-031/REQ-158-complete-rendered-region-classification.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `47b71fd`.
