---
source_type: req_lesson
req_id: REQ-148
req_path: do-work/archive/UR-034/REQ-148-preserve-empty-quarantine-association-candidates.md
date: 2026-08-07
domain: general
module: actions
tags: [general, addendum, preserve, association, candidates]
---

# Lessons from REQ-148: Addendum: preserve association candidates with empty quarantine

## What the REQ was about

Correct REQ-128's commit and unscoped-inspect quarantine merge so an empty run-level secret quarantine preserves every safe inventory candidate for REQ association.

## Solution summary

[MAP UNCHANGED] Commit and inspect still own their candidate filtering inline; only the first-input discriminator changed, with no new tool or interface.

## Worth knowing

- `NR == FNR` is not a safe first-file test when an input may be empty; identify that input explicitly through `FILENAME` and `ARGV`.
- A regression for a two-input merge must cover the zero-record first input, not only populated joins.

## Back-reference

See `do-work/archive/UR-034/REQ-148-preserve-empty-quarantine-association-candidates.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0ed7786`.
