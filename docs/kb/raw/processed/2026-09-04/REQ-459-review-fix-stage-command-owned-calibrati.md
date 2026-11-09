---
source_type: req_lesson
req_id: REQ-459
req_path: do-work/archive/REQ-459-stage-command-owned-calibration-with-lifecycle-release.md
date: 2026-09-01
domain: general
module: _dev/primes
tags: [general, review, stage, command]
---

# Lessons from REQ-459: Review fix: Stage command-owned calibration with lifecycle release

## What the REQ was about

Make every serial and worktree Commit Phase staging instruction include `do-work/calibration-log.tsv` when the canonical `requeststate complete` transaction changed it, and pin the current command-owned completion step rather than the removed manual calibration substep.

## Solution summary

The serial and worktree Commit Phase instructions now stage `do-work/calibration-log.tsv` only when the successful canonical `complete` result reports that path among its changed or affected targets. The obsolete Step 8 substep 7.5 staging predicate is gone, and an independent structural regression covers all four active staging surfaces.

## What worked

- Lifecycle staging must follow the canonical transaction's reported target set, not a remembered step number or filesystem inference.

## Back-reference

See `do-work/archive/REQ-459-stage-command-owned-calibration-with-lifecycle-release.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `1c0132399f2fbe2abe57e7280175e2565c848044`.
