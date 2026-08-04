---
source_type: req_lesson
req_id: REQ-082
req_path: do-work/archive/UR-016/REQ-082-fanout-handback-has-no-legal-write-location.md
date: 2026-08-03
domain: general
module: actions
tags: [actions, fan-out, hand-back, legal, write]
---

# Lessons from REQ-082: The fan-out hand-back file has no legal write location

## What the REQ was about

Fan-Out Dispatch makes a per-builder output file mandatory — `REQ-NNN-handback.md` inside
`do-work/runs/work-<timestamp>/` — and `crew-members/background-agents.md` is explicit that the
sub-agent writes that file itself. Worktree Dispatch Mode is equally explicit that a builder **never
writes the main tree**, and `do-work/` exists in the main tree only. There is no location satisfying
both rules, so the mandatory hand-back has no legal execution.

## Solution summary

`Sole integrator` now reads "The builder never writes the main tree or its branch, **with exactly one exception: its own `do-work/runs/work-<timestamp>/REQ-NNN-handback.md`**" — reached by the absolute main-tree path the orchestrator hands it, **never staged, committed, or merged**, and explicitly bounded: a sibling's hand-back, `manifest.md`, anything else under `do-work/runs/`, and every other main-tree path remain violations. One sentence, as the Constraints budgeted, carrying the *why* (the file exists because the transcript is not durable) so a maintenance pass cannot read it as redundant.

## What worked

- **Fixing the enumeration with a *negative* assertion alongside the positive one.** A single "contains the condition" assertion is satisfiable by a sentence that states the condition *and* keeps the stale list beside it — which is the likeliest accidental regression, since adding is easier than replacing. The pair makes re-adding the list a failure on its own.
- **Putting the carve-out inside the sentence it modifies.** It is the cheapest possible defence against the specific failure requirement 8 describes, and it cost nothing.

## What didn't work

- nothing failed. The sibling-enumeration grep requirement 5 mandated came back empty, which is worth recording as a *result* — the previous two REQs in this batch both had their site inventories turn out to be floors, and this one genuinely did not.

## Worth knowing

- **A contradiction between two files is invisible to a per-file assertion suite, and this one was also unreachable in practice.** Nothing has fan-out dispatched since REQ-073 shipped, so no run has ever produced a hand-back file — the contract has been broken for its entire life without a single observable symptom. Contract defects in unexercised paths are found by reading, not by running, which is an argument for REQ-085 actually happening.
- **The dangerous direction of a path-resolution bug is the one where the write succeeds.** Inbound, a repo-relative brief path yields nothing or a stale snapshot — visible. Outbound, the same mistake writes a real file into the builder's branch, gets swept into the merge as committed scratch, and the orchestrator quietly reads nothing. Same mechanism, and the return direction is worse, which is why the trap sentence now says so rather than being left as an exercise.

## Back-reference

See `do-work/archive/UR-016/REQ-082-fanout-handback-has-no-legal-write-location.md` for the full REQ — triage, implementation, review, and lessons. Commit `1cff0a7`.
