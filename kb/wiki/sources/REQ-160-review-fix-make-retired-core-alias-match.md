---
title: "Lessons from REQ-160: Review fix: Make retired core alias matching occurrence-complete"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-160-review-fix-make-retired-core-alias-match.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-160: Review fix: Make retired core alias matching occurrence-complete

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Make the test-only retired-command guard evaluate command occurrences and overlapping trigger candidates completely. An exemption or boundary-invalid longer candidate must not hide another valid retired command. Done means this false-negative class cannot recur while direct-alias suffix exclusions remain intact.

## Solution summary

The test-only recurrence guard now evaluates every occurrence and eligible overlapping historical candidate without weakening direct alias boundaries or changing the 186-row inventory/runtime contract.

## What worked

- Returning occurrence spans made the exemption rule exact and independently testable without weakening the retired-trigger vocabulary or runtime surface.
- Keeping candidate fallback family-scoped repaired the historical install/setup overlap while preserving direct-alias boundaries across all 186 rows.

## What didn't work

- The first lifecycle qualification omitted the mandatory P-A-U record even though implementation evidence was complete; the second attempt restored it before review.

## Worth knowing

- Longest-first command matching must validate each candidate before committing to it, and exemptions belong to source spans rather than whole lines. Independent review's repeated/different-occurrence matrix is the durable guard against both false-negative classes.

## Back-reference

See `do-work/archive/UR-031/REQ-160-make-retired-core-alias-matching-occurrence-complete.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `3d8613a`.
