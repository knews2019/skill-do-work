---
title: "Lessons from REQ-460: Make outside-text delimiter containment condition-complete"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-04/REQ-460-make-outside-text-delimiter-containment.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-460: Make outside-text delimiter containment condition-complete

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

Replace the finite prefix list used for answer-summary containment with a condition-complete classifier for lines that Markdown can interpret as document-owned delimiters or structure.

## Solution summary

**Supplied implementation commit:** `7e16f05c4e95ebf50fcf2d065e4f0145246d46ad`

## What worked

- State the structural condition first and make examples test fixtures. Leading whitespace, ASCII punctuation, and an ordered-list digit marker cover the requested class without another brittle prefix catalog.
- Compare old and new predicates as a superset, then exercise the real `BuildAnswerPlan` seam; this proves both no weakening and byte-preserving containment.

## What didn't work

- The old implementation trimmed leading whitespace before classification and thereby erased the evidence for indented structure.
- Reasoning only from a few passing delimiters made a partially correct enumeration look complete; review also found that the original hazard explanation ignored the summary's actual mid-line write position.

## Worth knowing

- punctuation-led prose and number-like text are intentionally over-contained as the safer reversible tradeoff. Adjacent disposition and cancellation consumers are owned by REQ-528 and REQ-529.

## Back-reference

See `do-work/archive/REQ-460-make-outside-text-delimiter-containment-condition-complete.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `7e16f05c4e95ebf50fcf2d065e4f0145246d46ad`.
