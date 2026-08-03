---
source_type: req_lesson
req_id: REQ-075
req_path: do-work/archive/UR-013/REQ-075-write-set-display-only-reason-went-stale.md
date: 2026-08-03
domain: general
module: tools/queue-kanban
tags: [restatement-sweep, stale-premise, grep-assertions, write-set, fan-out-dispatch]
---

# Lessons from REQ-075: Five files still explain write_set's display-only status with a reason fan-out made false

## What the REQ was about

The conclusion was still right — nothing schedules, gates, or dispatches on a ticket's declared `write_set` — but several files gave a reason that a prior REQ had falsified: "under the exclusive-session model one REQ runs at a time, so the badge schedules nothing." Once several builders could run at once under a single queue owner, a reader following that reasoning concludes the opposite of the contract: if the premise no longer holds, the conclusion looks like it should fall too, and `write_set` becomes a scheduling input — which is explicitly forbidden. The correct reason after fan-out dispatch is that `write_set` is advisory input to the human's pick, never a gate, and the merge refusing is the non-interference proof. That holds at any builder count.

## Solution summary

Replaced the reason at every site while keeping every conclusion and every "absence reads as unknown, not safe" clause intact, with each site pointing at the canonical Fan-Out Dispatch section instead of restating it. The REQ named five sites; the sweep found eleven, including the `write_set` schema line itself, four Go comments, and the badge tooltip a user reads in the browser. Added a two-part contract assertion: a line sweep for prose, plus file-level negatives for the two comment-carrying source files. No behavior change — no Go logic, no schema field, no board column touched.

## What worked

- Writing the assertion first and letting it enumerate the sites, instead of trusting the REQ's list. The list said five; the check found eight, and a second grep shape found three more. On a sweep REQ the check is the inventory — the prose inventory is a starting hypothesis.

## What didn't work

- Two dead ends, both instructive. (1) Line-granularity matching — the REQ's own suggestion — is blind to wrapped Go/JS comments, so the first green run was green for the wrong reason. (2) Widening to a 3-line proximity window then false-positived on the canonical Fan-Out section itself, where "integration runs one REQ at a time" is true. The shape that works is per-class: a line sweep for prose, a file-level negative for files that have no business naming a builder count at all.

## Worth knowing

- A falsified premise leaves two fingerprints, and grepping for one of them reads as clean. Here the strong form named the count ("one REQ at a time") and the weak form named the model ("under the exclusive-session model"). The weak form is the more dangerous of the two, because the model it names is *still true* — only its relevance died — so the sentence reads correct on inspection. When a premise is retired, sweep for the thing it was *called* as well as the thing it *said*.

## Back-reference

See `do-work/archive/UR-013/REQ-075-write-set-display-only-reason-went-stale.md` for the full REQ — exploration, the four numbered decisions including the requirement-4-versus-requirement-5 conflict, review, and lessons. Shipped as v0.166.2, commit `738e9fe`.
