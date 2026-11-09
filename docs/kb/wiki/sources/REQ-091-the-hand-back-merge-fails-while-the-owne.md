---
title: "Lessons from REQ-091: The hand-back merge fails while the owner's claim bookkeeping is staged, on any install that tracks do-work/"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-01/REQ-091-the-hand-back-merge-fails-while-the-owne.md]
related:
  - page: REQ-073-fan-out-dispatch-n-concurrent-builders-u
    rel: complements
  - page: REQ-085-run-req-073-s-live-two-builder-acceptanc
    rel: complements
  - page: REQ-092-actions-work-md-has-no-wave-selection-or
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-091: The hand-back merge fails while the owner's claim bookkeeping is staged, on any install that tracks do-work/

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

`actions/work.md` Step 2 claims a REQ by **moving** it from `do-work/queue/` to `do-work/working/` and
appending an entry to `do-work/CHECKPOINT.md`. Where the consumer **tracks `do-work/`**, that move is a
staged rename sitting in the index. Step 6's hand-back sequence
(`actions/work-reference.md` → Worktree Dispatch Mode) then says to run
`git merge --no-ff --no-commit <operative_name>` — and git refuses, because the merge would touch paths
with uncommitted local changes.

## Solution summary

The hand-back sequence gained a **step 0** — settle the index before capturing `<pre>` — in both places the sequence is written: its canonical home (`actions/work-reference.md` → Worktree Dispatch Mode, *When to merge*) and the condensed restatement an orchestrator actually follows (`actions/work.md` Step 6's Hand-back merge). Step 0 says to commit the owner's bookkeeping (claim moves, `CHECKPOINT.md`, the run directory if fan-out created one), states that the step is a no-op to be skipped where `do-work/` is untracked, and states why the ordering against `<pre>` is load-bearing.

## What worked

- Numbering the new step 0 instead of renumbering — two live cross-references cite the old numbers, and a renumber would have broken both silently while looking tidier.

## Worth knowing

- The defect only exists where the consumer commits `do-work/`, and it is invisible on the common untracked install — so any future change here must be checked against both shapes, and the prescribed step has to read sensibly to a reader for whom it is a no-op.

## Back-reference

See `do-work/archive/UR-016/REQ-091-handback-merge-collides-with-the-owners-claim.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `ecf1966`.
