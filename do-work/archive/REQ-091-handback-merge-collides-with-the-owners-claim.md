---
id: REQ-091
title: The hand-back merge fails while the owner's claim bookkeeping is staged, on any install that tracks do-work/
status: completed
claimed_at: 2026-08-04T00:26:00Z
completed_at: 2026-08-04T00:25:54Z
commit: ecf1966
kb_status: pending
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
write_set: [actions/work-reference.md, actions/work.md]
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

## Triage

**Route: A** - Simple

**Reasoning:** One defect, exact reproduction, named location. The open part is which candidate shape
to prescribe (requirement 3) — a decision, not a discovery.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified)
- `actions/work.md` (modified)

**What was done:** The hand-back sequence gained a **step 0** — settle the index before capturing
`<pre>` — in both places the sequence is written: its canonical home
(`actions/work-reference.md` → Worktree Dispatch Mode, *When to merge*) and the condensed restatement
an orchestrator actually follows (`actions/work.md` Step 6's Hand-back merge). Step 0 says to commit
the owner's bookkeeping (claim moves, `CHECKPOINT.md`, the run directory if fan-out created one), states
that the step is a no-op to be skipped where `do-work/` is untracked, and states why the ordering
against `<pre>` is load-bearing.

## Decisions

- **D-01 (DECIDE & STATE)** — *Prescribe "commit the bookkeeping", placed inside the hand-back sequence,
  rather than at Step 2 or by avoiding the index.* Three candidates were considered (requirement 3):
  - *Avoid staging altogether* (plain `mv` instead of `git mv`, so nothing enters the index) —
    rejected. The refusal is about local changes to tracked paths, not about the index specifically, so
    an unstaged rename is not reliably safe; it would trade a loud, reproducible failure for a subtle
    one that depends on git's overwrite heuristics.
  - *Commit the claim at Step 2, for every run* — rejected. Serial mode has no merge and therefore no
    problem, and this would cost it its one-commit-per-request property for nothing.
  - *Commit inside the hand-back sequence* — chosen. It is scoped exactly to the mode where the defect
    exists, changes serial mode not at all, and sits where the reader is already being told what to run.
- **D-02 (DECIDE & STATE)** — *Numbered it **step 0** rather than renumbering the existing four.* Two
  other passages cite the sequence by number — "step 3 of the hand-back sequence below"
  (`work-reference.md`, Sole integrator) and "step 2 of the hand-back sequence below" (Merge, never
  rebase). Renumbering would have silently falsified both. A step 0 also reads correctly as what it is:
  a precondition, not a new stage of the merge.
- **D-03 (DECIDE & STATE, scope extension)** — *`actions/work.md` added to scope.* The Restatement
  Sweep found its Step 6 carrying a condensed three-step version of the same sequence — and that is the
  file an orchestrator actually follows, so fixing only the reference would have left the primary reader
  broken and the REQ's own purpose unmet. Declared before the edit.

## Qualification

Passed — 2 files verified, 5 requirements traced.

- **Requirements traced:** 1 → step 0 added to the sequence itself, in both homes. 2 → both copies state
  the untracked case explicitly as a skip, not an error, so neither reads as noise on the common
  install. 3 → D-01 names all three candidates and why two were rejected. 4 → both copies state the
  ordering against `<pre>` and *why* (a bookkeeping commit inside the merge range reads as an
  undeclared `do-work/` touch to qualify and review). 5 → no lock, registry, heartbeat or liveness
  signal added; the change is one prose step.
- **Restatement Sweep:** the two number-citing passages (`work-reference.md` Sole integrator → "step 3",
  Merge-never-rebase → "step 2") were checked after the edit and both remain accurate, which is D-02's
  whole point. `actions/work.md:425` and `:583` reference the hand-back merge but not its step numbers.
  Grepped `hand-back sequence|hand-back merge` across all `.md` outside `do-work/`; the CHANGELOG hits
  are historical entries, not restatements.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ No new failures — 8 FAIL lines, the same pre-existing update-script probe failures
recorded under REQ-083's Discovered Tasks.

Prose change to action-file instructions; no behavioral code, so red-green does not apply. The
regression evidence is the contract suite plus the sweep above. **The defect itself was already
reproduced** — REQ-085's run hit it live, with git's exact error recorded there — and this REQ's GREEN
condition (an orchestrator following the sequence literally reaches a successful merge on a repo that
tracks `do-work/`) is exactly what this session did after committing its bookkeeping as `306c1f4`: the
next merge succeeded, and the two merge ranges it produced are both outside that commit.

*Verified by work action*

## Lessons Learned

**What worked:** Numbering the new step 0 instead of renumbering — two live cross-references cite the
old numbers, and a renumber would have broken both silently while looking tidier.

**Worth knowing:** The defect only exists where the consumer commits `do-work/`, and it is invisible on
the common untracked install — so any future change here must be checked against both shapes, and the
prescribed step has to read sensibly to a reader for whom it is a no-op.

## Orientation

**Now the hand-back sequence tells the orchestrator what to do about its own uncommitted claim before
merging a builder branch** — previously the sequence assumed a mergeable index and git refused outright
on any install that tracks `do-work/`. Lives in the worktree-dispatch procedure
(`actions/work-reference.md` → Worktree Dispatch Mode) and its condensed twin in `actions/work.md`
Step 6.

Leaf change to a procedure — no new contract, and the merge range's definition is unchanged. What did
change is a stated ordering constraint: the bookkeeping commit must precede the `<pre>` capture, which
is the kind of thing a later edit could reorder without noticing, so the reason is written next to it.

## Full Context

`do-work/archive/REQ-085-run-the-live-two-builder-acceptance-test.md` → `## Testing` → finding F-01.
