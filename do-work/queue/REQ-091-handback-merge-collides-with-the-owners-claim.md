---
id: REQ-091
title: The hand-back merge fails while the owner's claim bookkeeping is staged, on any install that tracks do-work/
status: pending
created_at: 2026-08-04T00:14:18Z
user_request: UR-016
domain: general
prime_files: []
tdd: false
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-073, REQ-085]
addendum_to: REQ-085
discovered_during: REQ-085
write_set: [actions/work-reference.md]
---

# The Hand-Back Merge Fails While the Owner's Claim Bookkeeping Is Staged

## What

`actions/work.md` Step 2 claims a REQ by **moving** it from `do-work/queue/` to `do-work/working/` and
appending an entry to `do-work/CHECKPOINT.md`. Where the consumer **tracks `do-work/`**, that move is a
staged rename sitting in the index. Step 6's hand-back sequence
(`actions/work-reference.md` → Worktree Dispatch Mode) then says to run
`git merge --no-ff --no-commit <operative_name>` — and git refuses, because the merge would touch paths
with uncommitted local changes.

Nothing in the pipeline commits the claim before that point, and the hand-back sequence does not
mention the state of the index at all.

## Why

This is not a hypothetical. It was the **first thing** REQ-085's live fan-out acceptance run hit, on the
very first merge — a hard `exit 2` in a procedure that had already been written, reviewed,
contract-asserted, and shipped as `completed` at v0.166.0.

It is invisible to inspection because neither half is wrong on its own. Step 2's claim is correct.
Step 6's merge is correct. Only their *ordering* against a tracked `do-work/` fails, and no single file
states both halves.

## Context

Reproduced verbatim during REQ-085 (`do-work/archive/REQ-085-…md` → `## Testing` → finding F-01):

```
$ git merge --no-ff --no-commit worktree-agent-REQ-086-in-progress-record-unstated
error: Your local changes to the following files would be overwritten by merge:
  do-work/working/REQ-085-….md do-work/working/REQ-086-….md do-work/working/REQ-087-….md
Merge with strategy ort failed.
                                                                             (exit 2)
```

**Scope is narrow and must stay that way.** The failure requires the consumer to commit `do-work/`. On
the common install `do-work/` is untracked, so the claim leaves the index clean and the merge proceeds
— any fix that assumes the tracked shape must not break the untracked one. Serial mode never meets this
at all, because serial mode never merges.

**The workaround REQ-085 used**, recorded but not blessed: the owner committed the claim bookkeeping
(and the run directory) as its own commit before merging. That happens to be convenient — it also makes
`do-work/runs/` durable, which `crew-members/background-agents.md` wants — but it was an improvisation
under a blocked merge, not a decision.

## Detailed Requirements

1. **State what the owner does about the index before the hand-back merge**, in the hand-back sequence
   itself (`actions/work-reference.md` → Worktree Dispatch Mode → *When to merge, and the range every
   evidence step reads*). The sequence is currently four steps that assume a mergeable index; whatever
   is decided becomes part of it.
2. **Whatever is prescribed must be correct on both install shapes** — tracked `do-work/` (where the
   index is dirty) and untracked (where it is not). A step that only makes sense on one of them will
   read as noise on the other and get skipped.
3. **Decide deliberately between the candidate shapes and say why.** At least: commit the claim as its
   own bookkeeping commit before dispatch or before the merge; or scope the claim's staging so it is
   never in the index at merge time. Each has consequences for what `<pre>` points at, so requirement 4.
4. **Say what the choice does to `<pre>`.** `<pre>` is captured immediately before the first merge, so a
   claim commit inserted there lands *below* `<pre>` and stays outside the merge range — which is
   correct, but only if the ordering is stated rather than left to chance.
5. **Do not introduce coordination state.** No lock, no registry, no liveness signal — the same
   constraint every worktree REQ carries.

## Constraints

- **Fan-out concurrency itself is not in scope** — that is REQ-092. This REQ is about one ordering
  problem in the hand-back sequence, and it bites a single builder just as hard as several.
- `do-work/` state remains the orchestrator's alone; nothing here changes what a builder may write.

## Dependencies

`addendum_to: REQ-085`, which found it. `related: REQ-073`, which shipped the sequence.

## Builder Guidance

**Certainty: Firm on the defect, open on the fix.** The reproduction is exact and the scope condition
(tracked `do-work/`) is established. Which of the candidate shapes is right is a genuine design
choice — log it as a D-XX.

Read REQ-085's `## Testing` before starting; the run that found this recorded the surrounding state,
including why the workaround it used should not be adopted by default without argument.

## Red-Green Proof

**RED prompt/case:** On a repo that tracks `do-work/`: claim a REQ per Step 2 (`git mv` into
`working/`, append to `CHECKPOINT.md`) without committing, then run
`git merge --no-ff --no-commit <a builder branch>`.

**Why RED now:** git exits 2 with `Your local changes to the following files would be overwritten by
merge`, and the hand-back sequence offers no instruction for that state. Reproduced during REQ-085.

**GREEN when:** an orchestrator following the hand-back sequence literally, on a repo that tracks
`do-work/`, reaches a successful merge — and the same sequence still reads correctly on an untracked
install.

**Validation:** Reproduced during REQ-085's live acceptance run; the exact error is recorded there.

## Full Context

`do-work/archive/REQ-085-run-the-live-two-builder-acceptance-test.md` → `## Testing` → finding F-01.
