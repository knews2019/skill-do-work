---
source_type: req_lesson
req_id: REQ-462
req_path: do-work/archive/REQ-462-preserve-deletion-precedence-across-git-xy-inventory-states.md
date: 2026-09-01
domain: backend
module: skills/do-work/tools/do-work-cli
tags: [backend, review, preserve, deletion]
---

# Lessons from REQ-462: Review fix: Preserve deletion precedence across Git XY inventory states

## What the REQ was about

Classify every Git porcelain XY state from the path's usable filesystem condition, with deletion taking precedence whenever either side says the path is absent. Done means combined status states cannot be mistaken for readable or associable paths.

## Solution summary

The typed inventory classifier now applies deletion precedence whenever either Git porcelain column is `D`, and recognizes additions only from index-column `A` or `??`. A 45-state retained differential covers ordinary, secret, rename/copy, unmerged, ambiguity, mutation, and protected-association behavior.

## What worked

- Two-column Git status classifiers need explicit precedence over the path's usable state; checking for additions before deletions can label an absent path readable.
- *Source: REQ-414 fresh re-review finding 1.*

## Back-reference

See `do-work/archive/REQ-462-preserve-deletion-precedence-across-git-xy-inventory-states.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `24abcf96408d33440498f67dcb6a59ef4240c03a`.
