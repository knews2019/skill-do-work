---
id: REQ-075
title: Five files still explain write_set's display-only status with a reason fan-out made false
status: pending
created_at: 2026-08-03T15:20:00Z
user_request: UR-013
addendum_to: REQ-073
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: []
maintenance: true
review_generated: true
discovered_during: REQ-073
write_set: [actions/board.md, actions/capture-reference.md, docs/board-guide.md, tools/queue-kanban/prime-do-kanban.md, _dev/tests/contract-regressions.sh]
---

# Five Files Still Explain `write_set`'s Display-Only Status With a Reason Fan-Out Made False

## What

The conclusion is still right — nothing schedules, gates, or dispatches on `write_set` — but five files
give a reason for it that REQ-073 falsified. They say some version of:

> Under the exclusive-session model one REQ runs at a time, so the badge schedules nothing.

Since REQ-073, several builders **can** run at once under a single queue owner. So a reader following
that reasoning concludes the opposite of the contract: if the premise ("one REQ at a time") no longer
holds, the conclusion ("nothing schedules on it") looks like it should no longer hold either — and
`write_set` becomes a scheduling input. That is explicitly forbidden.

Replace the reason, keep the conclusion. The correct reason after REQ-073: `write_set` is **advisory
input to the human's pick, never a gate**, and the **merge** — `git merge --no-ff --no-commit`
refusing — is the non-interference proof. That holds at any builder count.

## Why This Is Worth a REQ Rather Than a Note

The stale text does not merely read oddly — it argues for the wrong behavior, in the files an agent is
most likely to read while touching the board or the scheduler. `actions/board.md` is loaded whenever
the board runs; `tools/queue-kanban/prime-do-kanban.md` is the prime an agent reads before editing the
tool. Both currently hand that agent a premise that has become false and a conclusion that depends on
it.

## Context

Each site, with the phrasing to replace:

- `actions/board.md:92` — "Under the exclusive-session model one REQ runs at a time, so the badge schedules nothing"
- `actions/board.md:117` — "display-only under the exclusive-session model, since one REQ runs at a time"
- `docs/board-guide.md:39` — "Under the exclusive-session model `do-work run` builds one REQ at a time, so the badge schedules nothing"
- `tools/queue-kanban/prime-do-kanban.md:57` — "Nothing schedules on `write_set` under the exclusive-session model (the work pipeline runs one REQ at a time)"
- `actions/capture-reference.md:113` — "nothing schedules on it under the exclusive-session model" (weakest of the five — it names the model without asserting one-REQ-at-a-time, so it may only need the pointer updated)

The corrected wording already exists in three places to copy from, all landed by REQ-073:
`actions/work.md` § Rules, `actions/work.md` Step 5.5, and `CLAUDE.md` § Shipped Tooling. The canonical
statement is `actions/work-reference.md` → Worktree Dispatch Mode → **Fan-Out Dispatch**.

Nothing about the board's *behavior* changes: `annotateWriteSetOverlap` still runs after bucketing,
still feeds only the badge and drawer row, and `tools/queue-kanban/` column logic stays untouched. This
is a prose-only correction.

## Detailed Requirements

1. **Replace the reason at all five sites**, keeping each file's existing voice and length. The
   conclusion ("nothing schedules / gates / dispatches on it") does not change.
2. **Do not restate the Fan-Out Dispatch contract** — point at it. `actions/work-reference.md` is its
   canonical home, and a sixth copy of the reasoning is what created this REQ.
3. **Keep "absence reads as unknown, not safe"** wherever it already appears. That property is
   independent of builder count and is load-bearing for anyone using the badge to pick a fan-out set.
4. **Add a contract-suite assertion** that no shipped file justifies `write_set`'s display-only status
   with a one-REQ-at-a-time premise — e.g. no file under `actions/`, `docs/`, or `tools/queue-kanban/`
   matches both a `write_set`/`overlaps` mention and a "one REQ at a time"-shaped clause in the same
   line. Without it, the sixth copy arrives with the next edit.
5. **Leave `tools/queue-kanban/*.go` alone.** No behavior change, no schema change.

## Discovered During

REQ-073 (`do-work/archive/` — fan-out dispatch), by that REQ's restatement sweep. The three sites
inside REQ-073's declared Scope were fixed there; these five sat outside it, and the sweep rule routes
an out-of-scope stale restatement to a follow-up rather than widening the original REQ's diff.

## Red-Green Proof

**RED prompt/case:** The assertion from requirement 4, run against the current tree — it must fail,
naming the five files.

**Why RED now:** All five still carry the stale premise; only the three REQ-073 touched were corrected.

**GREEN when:** The assertion passes, `grep -rn "one REQ at a time" actions/ docs/ tools/queue-kanban/`
returns nothing that is offered as a reason for `write_set`'s status, and the full contract suite stays
green (including REQ-073's own exactly-once invariant check).

## Full Context

See `do-work/user-requests/UR-013/input.md` for the batch input, and REQ-073's `## Review` for the
sweep that surfaced this.
