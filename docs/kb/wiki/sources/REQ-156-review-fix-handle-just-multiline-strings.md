---
title: "Lessons from REQ-156: Review fix: Handle Just multiline strings in collision scanning"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-156-review-fix-handle-just-multiline-strings.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-156: Review fix: Handle Just multiline strings in collision scanning

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

Make reserved-recipe collision scanning ignore recipe- and alias-shaped text inside valid triple-quoted Just variable strings.

## Solution summary

Reserved-looking recipe and alias payload inside valid Just triple-single/triple-double multiline values is ignored, while real reserved definitions before or after those values retain the existing exact sorted error and pre-mutation byte preservation.

## What worked

- Keeping the existing header classifier intact made the requested triple-string repair small and preserved exact collision behavior around the new lexical state.
- Just-parseable positive fixtures plus byte-preserving negative controls proved both acceptance and safety through the production helper.

## What didn't work

- The initial durable summary generalized from the two requested triple-quoted forms to all multiline Just values. Independent review reproduced the same false rejection for three adjacent literal forms; the final summary wording is narrowed and REQ-159 records the remaining class.

## Worth knowing

- Once a handwritten safety scanner owns a lexer boundary, adjacent valid syntax needs explicit adversarial coverage even when the requested defect names only one token family.

## Back-reference

See `do-work/archive/UR-031/REQ-156-handle-just-multiline-strings-in-collision-scan.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `db9cd11`.
