---
title: "Lessons from REQ-163: Review fix: Complete remaining inline-link and list-fence classification"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-163-review-fix-complete-remaining-inline-lin.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-163: Review fix: Complete remaining inline-link and list-fence classification

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Finish the test-only rendered-region classifier so delimiter parity, label context, and list-item fences cannot disagree with target discovery. Done means relative and first-party targets classify consistently and fenced examples inside list items remain ignored, without changing publication, topology, raw/blob, path-containment, runtime, or documentation policy.

## Solution summary

The shipped-reference release guard now classifies all three remaining rendered-region variants consistently while preserving target normalization, publication policy, source length, newline offsets, list-paragraph continuations, and adjacent distribution contracts.

## What worked

- Relative targets beside first-party controls exposed structural extraction gaps that the bare-URL fallback had previously hidden.
- Treating a complete live link as one region and tracking list fences at their content column fixed the context errors without altering downstream target policy.
- Differential minimal reproductions against the CommonMark parser separated rendered references from link-shaped source text and exposed both directions of classifier drift: a live post-list link was hidden, while a literal backslash-separated pseudo-link was published.

## What didn't work

- The even-parity destination-opener premise in D-02 was incorrect. Backslash parity governs whether punctuation is escaped, but it does not remove source characters before structural parsing; any backslash between `]` and `(` breaks the adjacency required for an inline link.
- List-fence state tracked delimiter and indentation but not the containing list item's lifetime. Closed-only fixtures missed the valid CommonMark path where an unclosed fence ends at container dedent, causing the checker to mask subsequent top-level links through EOF.

## Worth knowing

- Markdown classifiers need both lexical and container state. Inline targets require direct `](` adjacency; backslash-separated pseudo-links should remain masked from bare-URL fallback without being extracted structurally. A list fence ends on a compatible closer or a nonblank container-ending dedent, and the dedented line must be reprocessed as live Markdown.
- Pair structural fixtures with an authoritative renderer, and cover zero through several escape characters plus explicit, missing, and over-indented fence closers. Self-consistent helper tests can otherwise ratify the same mistaken grammar assumption as the implementation.

## Back-reference

See `do-work/archive/UR-031/REQ-163-complete-remaining-inline-link-and-list-fence-classification.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `c9d1acd`.
