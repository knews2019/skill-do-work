---
title: "Lessons from REQ-174: Validate root Markdown fence info"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-174-validate-root-markdown-fence-info.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-174: Validate root Markdown fence info

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Make root and list fence classification share the CommonMark rule that a backtick-fence info string cannot itself contain a backtick.

## Solution summary

Centralized marker-aware fence info-string validation across root and list openings plus paragraph/container state, and added root/list/tilde differential fixtures that preserve the existing classifier contracts.

## What worked

- Markdown fence recognition is a compound contract: marker kind, info-string validity, and paragraph/container state must change together or a locally correct opener check can still mask rendered content.
- Differential fixtures against the pinned renderer catch classifier drift more reliably than asserting regex structure.

## Back-reference

See `do-work/archive/UR-039/REQ-174-validate-root-markdown-fence-info.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `bd5ecf6`.
