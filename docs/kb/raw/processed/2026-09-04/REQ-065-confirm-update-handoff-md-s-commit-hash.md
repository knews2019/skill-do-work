---
source_type: req_lesson
req_id: REQ-065
req_path: do-work/archive/UR-010/REQ-065-handoff-writeback-guidance.md
date: 2026-07-31
domain: general
module: general
tags: [general, confirm, update, handoff]
---

# Lessons from REQ-065: Confirm: update HANDOFF.md's commit-hash write-back guidance

## What the REQ was about

A `[low]` discovery from REQ-062's restatement sweep. `do-work/HANDOFF.md:35` currently tells this repo's future sessions:

## Solution summary

Rewrote the single bullet at `do-work/HANDOFF.md:35`.

## What worked

- The `commit:` field is redundant by construction in any repo do-work has run in: Step 9 mandates commit titles of the form `[{id}] {title} (Route {route})`, so `git log --grep='\[REQ-NNN\]'` recovers the hash, and the grep survives history rewrites that invalidate a recorded hash. Its one non-redundant use is disambiguating the `--no-ff` merge commit from the builder's branch commit in worktree dispatch mode. Worth knowing before anyone invests further in that field — the guard script earns its keep by protecting the *file*, not by protecting the hash.
- When deleting a local instruction that competes with a shipped one, delete it *and* say why in place. A bare deletion is indistinguishable from an omission, and the next session restores it.

## Back-reference

See `do-work/archive/UR-010/REQ-065-handoff-writeback-guidance.md` for the full REQ — plan, exploration, implementation, review, and lessons.
