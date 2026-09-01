---
title: "Lessons from REQ-152: Review fix: Reject reserved Just recipe collisions without Just"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-152-review-fix-reject-reserved-just-recipe-c.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-152: Review fix: Reject reserved Just recipe collisions without Just

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

Make the suite installer reject marker-free Justfiles that already define suite-reserved recipes even when the optional `just` executable is unavailable, before confirmation or mutation can leave an invalid file.

## Solution summary

Existing marker-free Justfiles can no longer acquire duplicate suite-reserved recipes when Just is unavailable. The replacer derives the protected namespace from the actual managed section, checks only top-level definitions/aliases outside replaceable ownership, and rejects before writing its candidate.

## What worked

- Deriving reserved names from the managed section prevents the safety gate from drifting when the shipped recipe set changes.
- Running the opt-in guard inside the replacer rejects collisions before confirmation or any client mutation while reusing managed-span ownership.

## What didn't work

- A line-local Just header parser missed multiline variable-string state; valid recipe-shaped string content can false-positive even though normal variables, comments, bodies, attributes, prefixes, and aliases are covered.

## Worth knowing

- Optional `just --list` remains the full syntax validator, but collision safety no longer depends on it. REQ-156 is held for consent to close the triple-quoted string edge without weakening real collision detection.

## Back-reference

See `do-work/archive/UR-031/REQ-152-reject-reserved-just-recipe-collisions-without-just.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `e2230b8`.
