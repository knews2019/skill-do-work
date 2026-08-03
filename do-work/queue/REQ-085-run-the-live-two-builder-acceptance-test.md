---
id: REQ-085
title: Run REQ-073's live two-builder acceptance test and record what it found
status: pending
created_at: 2026-08-03T17:09:21Z
user_request: UR-016
domain: testing
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-073, REQ-082]
batch: audit-remediation-external
addendum_to: REQ-073
---

# Run REQ-073's Live Two-Builder Acceptance Test and Record What It Found

## What

REQ-073 raised worktree dispatch from one builder to N and shipped as `completed` at v0.166.0. Its
`## Red-Green Proof` GREEN condition includes a live run of two concurrent builders; that run has never
happened. Two consecutive session checkpoints now carry it as deferred. Everything built since has been
serial, so grep proves the prose and nothing proves two builders compose.

This REQ's deliverable is **the run and its recorded outcome** — not a code change. Anything it breaks
becomes its own REQ.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

A capability nobody has exercised is a claim, not a feature — and this one is documented as shipped in
a user-facing changelog. REQ-073's own review was honest about it ("Acceptance: **Partial** … the REQ's
live two-builder GREEN condition was not run") and filed it as the first item under Suggested
additional testing, where it has now sat through two batches.

The cost of leaving it is concrete and already visible: REQ-082 exists because the fan-out hand-back
file has no legal write location, and that contradiction survived a full REQ, a review, and a
contract-assertion suite. It would have surfaced in the first five minutes of a live run. Every
remaining fan-out defect is in the same category — reachable only by execution.

Deferral has also stopped being tracked. `do-work/CHECKPOINT.md:38-40` is the only record, and the
queue is empty; a checkpoint bullet is not a queue entry, and the batch it was written for is closed.

## Context

**The procedure already exists.** `do-work/archive/UR-013/REQ-073-fan-out-dispatch-n-builders-one-owner.md`
→ `## Review` → *Suggested additional testing*, first item, plus the GREEN clause in that REQ's
`## Red-Green Proof`. Read both before starting; do not re-derive the check list.

**The positive case**, from REQ-073's GREEN condition — two non-overlapping REQs, two worktrees, two
branches:

- both branches merge cleanly;
- each REQ gets its own changelog entry, with strictly increasing versions;
- `do-work/working/` never holds a file the owner did not put there;
- `git worktree list` and `git branch --list 'worktree-agent-*'` are empty after both archives;
- the run directory is deleted.

**The negative case:** a deliberately overlapping pair must **fail** at
`git merge --no-ff --no-commit` rather than merging silently.

**A prior audit deliberately declined to make this a REQ.** `do-work/user-requests/UR-015/input.md`
records: "It needs a human running a real fan-out, not a REQ." That reasoning is preserved here and was
overridden by the user's instruction to capture all seven accepted findings. The override is honored by
shaping the REQ so a session *can* execute it — the pipeline dispatches builders, so an orchestrator
running this REQ is the human-equivalent — and by requirement 7, which makes an unrunnable environment
a visible failure instead of a quiet close.

**What this run is expected to collide with.** Named so a failure reads as information rather than
surprise:

- The hand-back file's write location (REQ-082). If that REQ has not landed, expect the builders to
  have nowhere legal to write and record exactly what they did instead.
- `queue-kanban verify` reporting the sibling builder as a fixable orphan mid-integration (REQ-083).
- Nothing in `actions/work.md` selects or claims a wave — Step 2 claims one REQ, Step 6 waits for one
  builder, Step 10 loops after commit. The wave has to be driven by hand for this run; that gap is
  deliberately not this REQ's fix (see Constraints).

## Detailed Requirements

1. **Pick two genuinely non-overlapping REQs** and say why they don't overlap, from their declared
   `write_set` **and** from reading them — the overlaps badge misses glob-vs-glob, `**`, and directory
   entries, and absence reads as unknown, not safe (`actions/board.md`). Real queue REQs are preferable
   to synthetic ones; if the queue offers no clean pair, synthesize two throwaway REQs and say so.
2. **Run the positive case** and record each of the five GREEN checks above as pass/fail with its
   evidence (the actual command output, not a summary).
3. **Run the negative case** — a deliberately overlapping pair — and confirm
   `git merge --no-ff --no-commit` refuses. A clean merge here is a finding, not a pass.
4. **Record the run in the REQ**, including the run directory's path, both `<operative_name>` values,
   both `<pre>..<merge_hash>` ranges, and the exact commands used. The point of this REQ is the
   evidence; a "ran it, worked" note is a failed deliverable.
5. **Every defect found becomes its own REQ**, not an inline fix. This REQ is an execution, and a
   builder that starts repairing the pipeline mid-run destroys the evidence it was dispatched to
   collect. Use `## Discovered Tasks` and let Step 8 queue them.
6. **Report what the run could not cover.** If the harness cannot dispatch two concurrent builders, or
   the negative case could not be constructed, say which check did not run — REQ-073's failure was
   exactly a partial acceptance recorded as complete, and repeating that here would be worse than not
   running it at all.
7. **If the run genuinely cannot be performed, fail the REQ** with `error_type: environment` rather
   than closing it. That keeps the gap in the pipeline where the next run can see it, which is the
   outcome UR-015's note was protecting against.
8. **Leave `do-work/` clean afterwards.** No stray worktrees, no `worktree-agent-*` branches, no
   surviving run directory, no throwaway REQs left in the queue. If a leftover cannot be removed
   without `--force`, report it and stop — that is itself a finding about the cleanup path.

## Constraints

- **Do not build the fan-out wave loop as part of this REQ.** `actions/work.md` has no wave-selection
  or launch-before-wait path, and adding one is a separate change whose shape depends on what this run
  finds. Drive the two builders by hand for this run and record that you did. (Parked deliberately —
  see UR-016's Batch Constraints.)
- **Do not fix anything the run breaks.** Requirement 5 is the whole discipline of this REQ.
- **Every worktree lives outside the repo working tree** — a nested second checkout is a documented
  corruption path (`actions/work-reference.md` → Worktree Dispatch Mode, *Where worktrees live*).
- **Never `-D`, never `--force`** on any worktree or branch this run creates. If `git worktree remove`
  or `git branch -d` refuses, that refusal is signal and gets reported (same rule as *Cleanup — happy
  path*).
- **Serial-only stays serial.** Queue transitions, REQ id allocation, `actions/version.md`, and
  `CHANGELOG.md` are the owner's and are not parallelised, whatever the run does with builders.
- **No new coordination state.** No lock, heartbeat, claim registry, or liveness probe may be
  introduced to make the run work; the forbidden-token sweep must stay green.
- `tdd: false` **is not "no proof needed."** The Red-Green Proof below is the deliverable's shape; it is
  simply not a unit test, because the thing under test is a harness capability.

## Dependencies

`addendum_to: REQ-073`, whose GREEN condition this completes. `related: REQ-082` — that contradiction is
the most likely early blocker, and running this **after** REQ-082 lands will produce a cleaner result,
but no `depends_on` is declared on purpose: a run that hits the contradiction and documents it is a
useful outcome, and gating this REQ would park it a third time.

## Builder Guidance

**Certainty: Firm on what to run; open on how to dispatch.** REQ-073 requirement 7 leaves the dispatch
mechanism deliberately unspecified, and that stands — spawned subagents and separate sessions are
indistinguishable to the owner because it synthesizes from files. Pick whichever your harness supports
and record which you used, since that is part of the result.

**Be willing to report a failure.** The valuable outcome of this REQ is an honest record, and a
half-working fan-out documented precisely is worth more than a green tick. REQ-073's review already
modelled this by marking its own acceptance Partial.

Read `crew-members/background-agents.md` before dispatching — the run directory, per-builder input and
output files, and the manifest are mandatory here, not optional, and its ceiling note holds: the pattern
makes fan-out failures survivable, not prevented.

## Red-Green Proof

**RED prompt/case:** Today, the question "have two builders ever run concurrently under one queue owner
in this repo?" has the answer **no**, recorded in `do-work/CHECKPOINT.md:38-40` and in REQ-073's own
review ("Acceptance: **Partial**"). There is no artifact anywhere in `do-work/` showing a completed
two-builder run.

**Why RED now:** `git log` shows every REQ since REQ-073 committed serially, one REQ per commit, and no
`do-work/runs/` directory has ever existed in this tree.

**GREEN when:** This REQ's body carries a `## Testing` record containing: the five positive-case checks
with their command output; the negative case refusing at `git merge --no-ff --no-commit`; the run
directory path and both merge ranges; and either a clean bill or a `## Discovered Tasks` list of what
broke. GREEN is "the run happened and its outcome is written down" — **not** "the run passed."

**Validation:** User confirmed — the user instructed that this be captured as a REQ after the triage
recommended it, overriding a prior audit's judgment that it belonged to a human rather than the queue.

## Full Context

See `do-work/user-requests/UR-016/input.md` for the verbatim instruction, the provenance of the external
audit, and the batch constraints. The test procedure itself is in
`do-work/archive/UR-013/REQ-073-fan-out-dispatch-n-builders-one-owner.md`.

---
*Source: external audit finding F2, third claim (P1) — "REQ-073 itself records that its only live
acceptance test was never run and still marks the ticket completed" — verified against the archived REQ
and accepted by `do-work validate-feedback` triage.*
