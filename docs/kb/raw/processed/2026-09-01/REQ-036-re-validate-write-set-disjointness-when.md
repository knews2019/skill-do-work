---
source_type: req_lesson
req_id: REQ-036
req_path: do-work/archive/UR-007/REQ-036-dispatch-revalidation-firmed-write-sets.md
date: 2026-07-29
domain: general
module: actions
tags: [actions, re-validate, write-set, disjointness, firms]
---

# Lessons from REQ-036: Re-validate write-set disjointness when Step 5.5 firms the sets

## What the REQ was about

REQ-032's dispatch gate decides co-dispatch on capture-seeded `write_set`s, which `actions/capture-reference.md` itself calls "a hint, never a commitment" — and Step 5.5 then overwrites the field from the `## Scope` list with no conflict re-check. Add the missing re-validation: on a concurrent dispatch, the Step 5.5 mirror runs the same no-concurrent-claimant check Step 6 already requires for mid-build extensions, and serializes/partitions when it fails.

## Solution summary

Closed the unguarded second write path to `write_set`. `actions/work.md` Step 5.5 gained a co-dispatch re-validation paragraph — when more than one REQ is in flight, the mirror re-checks the firmed `## Scope` list for pairwise disjointness against every other in-flight REQ's current `write_set` before replacing the field (the same no-concurrent-claimant check Step 6 runs for mid-build extensions), serializing or partitioning the loser before its builder starts, and stating that a dispatch-time partition directive survives the mirror. The Step 1 gate gained a pointer naming Step 5.5 as the dispatch-time enforcement point (the gate itself being the initial scheduling decision on capture-seeded hints); the Step 6 write-boundary bullet now glosses the absent-set case (serial dispatch ⇒ no declared boundary, not "write nothing"); and Step 4's plan-validation file-conflict item now reconciles its flag as a warning whose real enforcement is Step 1 + Step 5.5. `actions/work-reference.md`'s Scope Declaration Template gained the matching partition-survives-mirror note. Every added clause is gated to the parallel-dispatch path; the serial/floor path is behaviorally unchanged (the absent-set gloss only widens/clarifies a serial builder's freedom). No ratchet added, no schema shape change.

## What worked

- Gating every new clause to the parallel-dispatch path (explicit "a serial run skips this entirely" hedges) made the serial/floor invariant checkable by grep, not just by argument. The adversarial hunter's ordering analysis — proving the *last* REQ to firm its Scope always catches a firming-introduced overlap — is the kind of soundness check a requirements-walk alone would skip.

## What didn't work

- Writing two cross-referencing sentences in *different* steps (the Step 1 gate pointer and the Step 4 reconciliation) produced a micro-incoherence about whether Step 1 "enforces" — a coherence REQ briefly re-introducing the exact defect class it targets. Sentences that answer the same reader question ("which step enforces?") must be drafted side-by-side, not one-per-step.

## Worth knowing

- After REQ-035, "claimant" is a loaded word — it now means a lock `claimed_reqs` holder. Write-set/file-overlap checks should say "disjointness" or "no other REQ claims this file," never "claimant," to avoid conflation.

## Back-reference

See `do-work/archive/UR-007/REQ-036-dispatch-revalidation-firmed-write-sets.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `4296e11`.
