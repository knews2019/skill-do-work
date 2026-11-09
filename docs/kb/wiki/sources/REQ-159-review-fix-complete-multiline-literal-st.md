---
title: "Lessons from REQ-159: Review fix: Complete multiline literal state in Just collision scanning"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-159-review-fix-complete-multiline-literal-st.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-159: Review fix: Complete multiline literal state in Just collision scanning

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

Make the reserved-recipe collision scanner retain lexical state for every current Just multiline literal form that can contain column-zero recipe- or alias-shaped payload. Done means valid ordinary multiline single/double strings and triple-backtick command literals cannot recur as false collisions, while real definitions around every form remain detected exactly and pre-mutation preservation remains unchanged.

## Solution summary

The collision scanner now ignores reserved-looking payload inside the three captured multiline literal families while retaining exact real-definition detection, deterministic diagnostics, and atomic pre-mutation rejection around them.

## What worked

- Just-parseable positives paired with exact sorted diagnostics and byte snapshots proved both safe acceptance and unchanged rejection behavior through the production helper.
- Persisting delimiter identity while leaving the existing definition grammar untouched kept the paired change small and reviewable.

## What didn't work

- The captured three-family boundary called ordinary backticks one-line, but Just 1.46.0 accepts them across physical lines too; implementing the broad “every current form” claim without probing adjacent grammar would have silently overstated completion.

## Worth knowing

- A safety scanner that approximates a lexer needs accepted-syntax probes for every delimiter family. Explicit Requirements should still bound the builder when a broad summary conflicts; record the adjacent accepted form for consent instead of expanding silently.

## Back-reference

See `do-work/archive/UR-031/REQ-159-complete-multiline-literal-state-in-just-collision-scanning.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `6ba3a27`.
