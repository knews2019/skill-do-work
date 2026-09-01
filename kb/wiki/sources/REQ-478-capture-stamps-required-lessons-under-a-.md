---
title: "Lessons from REQ-478: Capture stamps required lessons under a token budget"
type: source-summary
topic_cluster: knowledge-and-memory
sources: [raw/processed/2026-09-01/REQ-478-capture-stamps-required-lessons-under-a-.md]
related:
  - page: concept-knowledge-and-memory-systems
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-478: Capture stamps required lessons under a token budget

Part of the [[concept-knowledge-and-memory-systems]] cluster.

## What the REQ was about

`capture-request` reads the lessons index while authoring REQ payloads and stamps the relevant lessons files as mandatory reads in a new frontmatter field, keeping the stamped set's summed token estimates within one stated budget.

## Solution summary

Captured requests can now carry only the most relevant lesson reads that fit one reproducible budget. The contract narrows eligible fully-slugged satellites before dropping candidates, records every drop, and leaves unrelated requests unstamped.

## What worked

**What worked:** Separating the one budget/cost contract from its consumers kept the number unique while letting capture and the next claim-time REQ share it; a focused request-model test proved the preservation seam directly.

**What didn't:** The shared tree's paused REQ-420 work overlapped two documentation files, so whole-file staging and dirty-tree qualification could not establish attribution.

**Worth knowing:** A budgeted context pointer needs both a single limit and a mechanical cost for every entry form; otherwise a “cheaper” targeted form makes compliance impossible to verify.

## Back-reference

See `do-work/archive/UR-088/REQ-478-capture-stamps-required-lessons-under-token-budget.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `47ff6c85`.
