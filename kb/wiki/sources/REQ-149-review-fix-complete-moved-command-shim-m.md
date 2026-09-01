---
title: "Lessons from REQ-149: Review fix: Complete moved-command shim mapping"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-149-review-fix-complete-moved-command-shim-m.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-149: Review fix: Complete moved-command shim mapping

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Make every legacy command routed by the modular core print one concrete, correct sibling invocation and stop. Close the whole router-to-shim mapping class, not only the aliases observed by this review.

## Solution summary

Closed the original ambiguous-routing test gap at the permanent post-migration boundary. The compatibility shim remains deleted; current board, knowledge, and toolbox routes now have one complete ownership contract, and the updater recipe must call only the core updater.

## What worked

- Translating an expired compatibility requirement to its durable ownership invariant closed the review finding without restoring aliases that the approved migration window had retired.
- Parsing the managed update recipe as a bounded body catches wrong-command mappings that a file-wide path grep cannot distinguish.

## What didn't work

- The original placeholder-string checks proved only that some owner/path text existed; they could not detect duplicate routes, missing sibling rows, or a board call hidden in the updater recipe.

## Worth knowing

- The permanent routing contract is the 23-action sibling table. Retired core spellings are intentionally absent and should not be reintroduced as test data; REQ-153 owns remaining stale prose.

## Back-reference

See `do-work/archive/UR-031/REQ-149-complete-moved-command-shim-mapping.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `dd509cd`.
