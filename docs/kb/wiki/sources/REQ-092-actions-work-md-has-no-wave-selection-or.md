---
title: "Lessons from REQ-092: actions/work.md has no wave-selection or launch-before-wait path, so documented fan-out concurrency cannot be reached by following it"
type: source-summary
topic_cluster: worktree-and-parallel-dispatch
sources: [raw/processed/2026-09-01/REQ-092-actions-work-md-has-no-wave-selection-or.md]
related:
  - page: REQ-073-fan-out-dispatch-n-concurrent-builders-u
    rel: complements
  - page: REQ-085-run-req-073-s-live-two-builder-acceptanc
    rel: complements
  - page: REQ-091-the-hand-back-merge-fails-while-the-owne
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-092: actions/work.md has no wave-selection or launch-before-wait path, so documented fan-out concurrency cannot be reached by following it

Part of the [[concept-worktree-isolation-and-parallelism]] cluster.

## What the REQ was about

`actions/work-reference.md` → Worktree Dispatch Mode → **Fan-Out Dispatch** documents several builders
running under one queue owner, with integration serialised behind them. `actions/work.md` has no step
that produces that shape:

- **Step 1** finds *the next* pending REQ.
- **Step 2** claims *one* REQ.
- **Step 6** spawns a builder and waits for it before Step 6.25 reads its output.
- **Step 10** loops back to Step 1 *after* the commit.

## Solution summary

Answered requirement 1 with **no** — `actions/work.md` should not drive a wave — and stated that boundary where readers currently infer the opposite. `actions/work.md` gained a paragraph next to the sub-agent note saying it processes one REQ at a time deliberately, naming the four steps that assume it, and pointing at `actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch as the owner-driven procedure that does the other thing. The `--wave N` section gained a paragraph saying the flag selects a batch and does not run one concurrently. `actions/work-reference.md`'s Fan-Out Dispatch section gained the matching statement from the other side: this is a procedure a human or advanced harness follows, not something the action performs.

## What worked

- Taking "document the boundary" as a real answer. The REQ was written to make that outcome available without a builder having to justify departing from an implied build, and it was the right one — the capability was never unreachable, only unowned by either document.

## Worth knowing

- The tell that this was a documentation defect rather than a missing feature: the capability had been *executed successfully* (REQ-085) before anyone tried to automate it. A feature you can perform by hand from the spec is specified; what was missing was a sentence saying which document performs it.

## Back-reference

See `do-work/archive/UR-016/REQ-092-work-action-has-no-path-that-drives-a-wave.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `92bebe0`.
