---
source_type: req_lesson
req_id: REQ-128
req_path: do-work/archive/UR-028/REQ-128-secret-rename-quarantine-and-staged-deletion.md
date: 2026-08-07
domain: general
module: tools/checks
tags: [general, secret, rename, quarantine, survives]
---

# Lessons from REQ-128: Secret rename quarantine survives re-inventory

## What the REQ was about

Correct REQ-121's uncommitted-inventory workflow so a secret rename remains excluded after an index reset degrades it into a deletion plus addition, and so an already-staged secret deletion can be committed without a failing restage.

## Solution summary

Inventory now overrides disabled rename detection, buffers classifications, and converts ambiguous additions beside a secret-shaped deletion to `X`. Commit and inspect keep excluded paths quarantined across re-inventory and mirror the rule in their manual fallbacks. Commit accepts an exact already-staged secret deletion and otherwise stages then verifies deletion-only metadata. Version 0.183.4 records the fix.

## What worked

- Buffering the complete NUL-delimited inventory made the ambiguous `XD` plus `A` rule fail closed without changing the public tag or exit-code interface.
- Cached name/status metadata was enough to distinguish an exact already-staged deletion without reading secret content.

## What didn't work

- Relying on Git's current rename record loses provenance after an index reset; later inventories need both an ambiguity rule and a run-level quarantine.
- The first action draft assumed a shell variable could survive separate prescribed command blocks; re-deriving a deterministic Git-private path is required.

## Worth knowing

`git add -u -- <path>` can reject a rename source whose deletion is already present in the index, so deletion-only workflows must verify cached metadata before deciding whether to stage.

## Back-reference

See `do-work/archive/UR-028/REQ-128-secret-rename-quarantine-and-staged-deletion.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `7bb03d2`.
