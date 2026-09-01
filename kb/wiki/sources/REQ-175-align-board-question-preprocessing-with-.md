---
title: "Lessons from REQ-175: Align board question preprocessing with valid Markdown fences"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-175-align-board-question-preprocessing-with-.md]
related:
  - page: REQ-174-validate-root-markdown-fence-info
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-175: Align board question preprocessing with valid Markdown fences

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

Make the board's Open Questions hard-break preprocessor distinguish real fenced code blocks from invalid backtick-info lookalikes, matching the Markdown renderer it feeds.

## Solution summary

The question-option preprocessor now rejects invalid backtick fence candidates using the same info-string condition as pinned Goldmark. Focused renderer coverage proves the invalid opener remains prose with two option-line breaks while valid backtick and tilde fences stay byte-identical during preprocessing.

## What worked

A production-seam RED/GREEN renderer test exposed the preprocessor/Goldmark disagreement directly, and the marker-aware info check kept the fix small.

## What didn't work

Prefix-only fence detection duplicated only part of Markdown's fence contract, so invalid prose was silently classified as code before Goldmark rendered it.

## Worth knowing

Backtick-fence info strings cannot contain backticks; tilde-fence info strings can. Any lightweight Markdown preprocessor must preserve that marker-specific distinction.

## Back-reference

See `do-work/archive/UR-039/REQ-175-align-board-question-preprocessing-with-valid-markdown-fences.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `17e1ea0`.
