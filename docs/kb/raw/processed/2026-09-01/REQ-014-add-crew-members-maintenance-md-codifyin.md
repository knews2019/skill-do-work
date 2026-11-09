---
source_type: req_lesson
req_id: REQ-014
req_path: do-work/archive/UR-002/REQ-014-delete-before-you-add-rule.md
date: 2026-06-28
domain: general
module: crew-members
tags: [crew-members, delete-before-you-add, maintenance, codifying]
---

# Lessons from REQ-014: Add crew-members/maintenance.md codifying delete-before-you-add

## What the REQ was about

Add a new crew-member rule file `crew-members/maintenance.md` that codifies
subtraction-first maintenance from the Agent Maintenance Loop: when fixing a drifting
agent/action during a maintenance pass, try removing or narrowing before adding any new
instruction. No such rule exists today — the closest is `karpathy.md`'s "Simplicity
First" and "don't delete adjacent dead code," which is an implementation-time rule, not
a maintenance principle.

## Solution summary

Created a new maintenance-time crew rule codifying subtraction-first maintenance (try removing/narrowing before adding an instruction; prove any addition against a replay case). Distinguished it explicitly from karpathy.md's implementation-time surgical-changes rule and from YAGNI (which it points at rather than restates). Wired the trigger into CLAUDE.md's Agent Rules as a third named contract and referenced it from quick-wins.md — refined from the REQ's "loaded by quick-wins" because quick-wins is read-only and loads no crew rules (recorded as D-01).

## What worked

- Reading `quick-wins.md` *before* wiring caught that it's read-only and loads no crew rules — turning the REQ's "loaded by quick-wins" into the honest "referenced from quick-wins" (D-01). Exploring the wiring target beat trusting the REQ's recommended phrasing.

## What didn't work

- The first pass updated only one of CLAUDE.md's two Agent-Rules bullets that enumerate triggers — ironically leaving stale exactly the kind of closed enumeration this REQ's rule warns about. The independent review caught it; fixed before commit. Lesson: when a value appears in two parallel lists, grep both (the sibling REQ-013 lesson "author one canonical source / read full context before editing" applies here too).

## Worth knowing

- maintenance.md is referenced but not yet auto-*loaded* during real subtraction work — there's no `maintain` action, and adding it to `work.md` Step 6's loader would misfire on every REQ. The load path closes when a dedicated maintenance action lands. Until then, the actor carries the quick-wins pointer across manually.

## Back-reference

See `do-work/archive/UR-002/REQ-014-delete-before-you-add-rule.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `db4d661`.
