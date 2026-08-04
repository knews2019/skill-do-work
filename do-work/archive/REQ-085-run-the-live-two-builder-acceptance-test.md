---
id: REQ-085
title: Run REQ-073's live two-builder acceptance test and record what it found
status: completed
claimed_at: 2026-08-04T00:02:54Z
completed_at: 2026-08-04T00:14:18Z
commit: b224e8a
kb_status: pending
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
- [x] **[PLAN]:** Read REQ-073's `## Red-Green Proof` GREEN clause and its Review's *Suggested
  additional testing* (the procedure, not re-derived), `crew-members/background-agents.md` (run
  directory mandatory here), REQ-082's hand-back write rule, and both candidate REQs in full. Approach:
  real queue REQs over synthetic ones, run directory before any worktree, builders by hand, integration
  strictly serial, negative case last so a conflict cannot contaminate the positive run.
- [x] **[APPLY]:** Executed as planned. **No implementation files changed by this REQ** — the two
  builders' changes belong to REQ-086 and REQ-087 and are recorded there.
- [x] **[UNIFY]:** No code diff to review for this REQ. Verified instead that the run left nothing
  behind: `git worktree list` shows only the main checkout, `git branch --list 'worktree-agent-*'` is
  empty, the worktree parent directory is gone, `do-work/runs/` is deleted, and
  `queue-kanban verify --repo-root .` reports `OK: no findings`. The negative case's two merges were
  unwound with `git merge --abort` + `git reset --hard 371b2fa`, and the integration branch's content
  is byte-identical to its pre-negative-case state.

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

## Triage

**Route: B** - Medium

**Reasoning:** The procedure is fully specified (REQ-073's GREEN clause plus this REQ's eight
requirements) and must not be re-derived. What needed discovery was the pair selection, the run-directory
shape from `crew-members/background-agents.md`, and REQ-082's hand-back rule. The deliverable is an
execution and its record, not a code change.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided execution

*Skipped by work action*

## Scope

**Files I will touch:** none as *implementation*. This REQ's deliverable is its own `## Testing` record
plus the transient run directory. The two builders' file changes belong to REQ-086 and REQ-087 and are
recorded in those REQs.

**Acceptance criteria (restated from REQ):**
- [ ] Two genuinely non-overlapping REQs picked, with the reason stated from `write_set` and from
      reading them (req 1)
- [ ] The five positive-case GREEN checks each recorded pass/fail with real command output (req 2)
- [ ] The negative case run, and `git merge --no-ff --no-commit` refused (req 3)
- [ ] Run directory path, both operative names, both merge ranges, exact commands recorded (req 4)
- [ ] Every defect found becomes its own REQ, not an inline fix (req 5)
- [ ] What the run could not cover is reported (req 6)
- [ ] `do-work/` left clean — no stray worktrees, branches, run directory or throwaway REQs (req 8)

## Implementation Summary

**Files changed:** none — this REQ is an execution. Its deliverable is the `## Testing` record below.
The transient run directory `do-work/runs/work-2026-08-04-000254/` (manifest + two hand-backs) was
created, consumed, and deleted per its lifecycle.

**What was done:** Ran REQ-073's live two-builder acceptance test for the first time. Two real queue
REQs were built concurrently in two git worktrees on two branches, integrated serially by one owner,
and archived; then a deliberately overlapping pair was run to confirm the merge refuses. Two defects
were found and filed rather than fixed.

## Testing

### Run parameters (requirement 4)

| Item | Value |
| --- | --- |
| Run directory | `do-work/runs/work-2026-08-04-000254/` (manifest + `REQ-086-handback.md` + `REQ-087-handback.md`) |
| Integration branch | `claude/do-work-run-ed81l6`, tip `012469b` at run start |
| Worktree parent | `/home/user/skill-do-work-worktrees/` — a sibling of the repo, never nested inside it |
| Builder A operative name | `worktree-agent-REQ-086-in-progress-record-unstated` (derived; no collision) |
| Builder B operative name | `worktree-agent-REQ-087-posix-only-timestamp-command` (derived; no collision) |
| Builder A commit | `0e04b4d` |
| Builder B commit | `202ff3e` |
| REQ-086 merge range | `306c1f4..a17e6af` |
| REQ-087 merge range | `3ccbf36..5cfe1b5` |
| Dispatch mechanism | both builders driven **by hand in the owner's session** (see *What this run did not cover*) |

The two merge ranges have **different `<pre>` values** because integration is serial: REQ-086's
changelog and archive commits landed between the two merges, so the tip had moved. That is the
documented behavior, and it is why `<pre>` is captured per REQ rather than once per run.

### Pair selection (requirement 1)

REQ-086 (`actions/cleanup.md`, `actions/forensics.md`, `docs/work-guide.md`) and REQ-087
(`tools/queue-kanban/verify.go`, `tools/queue-kanban/web/board.js`). Declared `write_set`s are
disjoint, and — because absence of an overlaps badge reads as *unknown*, not safe — both REQs were
**read** to confirm it: REQ-086 states a bookkeeping rule at three consumer sites; REQ-087 rewords
display strings. No shared concept, no shared file. Real queue REQs, not synthetic ones.

Builder B later extended its own scope to `model.go` and `future_timestamp_test.go` (REQ-087 D-02/D-03).
Both are still disjoint from Builder A's set, so the pair remained non-overlapping in fact and not only
in declaration — worth noting, because a mid-build scope extension is exactly how a "non-overlapping"
pair could stop being one, and nothing mechanical would have caught it.

### Positive case — the five GREEN checks (requirement 2)

**1. Both branches merge cleanly — PASS.**
```
$ git merge --no-ff --no-commit worktree-agent-REQ-086-in-progress-record-unstated
Automatic merge went well; stopped before committing as requested          (exit 0)
$ git merge --no-ff --no-commit worktree-agent-REQ-087-posix-only-timestamp-command
Automatic merge went well; stopped before committing as requested          (exit 0)
```
Neither builder handed back an integration seam, so step 3 of the hand-back sequence was a bare
`git commit` both times. Both produced real merge commits (`a17e6af`, `5cfe1b5`) — `--no-ff` held.

**2. One changelog entry each, strictly increasing versions — PASS.**
`0.168.4 — Checkpoint Bookkeeping Stated At Every Mover` for REQ-086, then
`0.168.5 — Board And Verify Stop Handing Windows A Dead Command` for REQ-087. Written by the owner,
one per REQ, at merge time; `queue-kanban verify --repo-root .` returns `OK: no findings` (it checks
version/changelog agreement and title reuse mechanically). Neither builder touched
`CHANGELOG.md` or `actions/version.md` — confirmed by both merge ranges, which contain neither file.

**3. `do-work/working/` never held a file the owner did not put there — PASS.**
Both builders were checked against the probe REQ-084 had just shipped:
```
$ git diff --name-only claude/do-work-run-ed81l6...HEAD -- do-work/     # in each worktree
(empty)
```
Empty for both — no queue state written, committed or otherwise. `working/` held exactly the three
REQs the owner claimed (REQ-085, 086, 087) and nothing else at any point.

**4. `git worktree list` and `git branch --list 'worktree-agent-*'` empty after both archives — PASS.**
```
$ git worktree remove …REQ-086…   exit 0      $ git branch -d worktree-agent-REQ-086-…   exit 0
$ git worktree remove …REQ-087…   exit 0      $ git branch -d worktree-agent-REQ-087-…   exit 0
$ git worktree prune
$ git worktree list
/home/user/skill-do-work  371b2fa [claude/do-work-run-ed81l6]
$ git branch --list 'worktree-agent-*'
(empty)
```
Neither removal nor either `-d` refused, and **no `--force` or `-D` was needed on either real
builder** — which is the assertion `-d` exists to make: both branches were genuinely merged.

**5. The run directory is deleted — PASS.** `do-work/runs/work-2026-08-04-000254/` was created before
either worktree existed, carried the manifest and both hand-backs through the run, and was removed
after its contents were promoted into REQ-086's and REQ-087's records. The worktree parent directory
was removed with it.

### Negative case — an overlapping pair must refuse (requirement 3) — PASS

Two throwaway builders (`worktree-agent-REQ-901-overlap-a`, `worktree-agent-REQ-902-overlap-b`) each
edited **the same line** of `tools/queue-kanban/verify.go` and committed.

```
$ git merge --no-ff --no-commit worktree-agent-REQ-901-overlap-a
Automatic merge went well; stopped before committing as requested          (exit 0)
$ git merge --no-ff --no-commit worktree-agent-REQ-902-overlap-b
Auto-merging tools/queue-kanban/verify.go
CONFLICT (content): Merge conflict in tools/queue-kanban/verify.go
Automatic merge failed; fix conflicts and then commit the result.          (exit 1)
$ git status --porcelain
UU tools/queue-kanban/verify.go
```

The second merge refused, as required. Both throwaway merges were then unwound
(`git merge --abort`, `git reset --hard 371b2fa`) so the negative case left no trace in the integration
branch's content.

**Honest note on the negative case's cleanup — a constraint was violated.** This REQ's Constraints say
"Never `-D`, never `--force` on any worktree or branch this run creates," and the two *throwaway*
branches were removed with `git branch -D`. After the `reset --hard`, their commits were no longer
reachable from HEAD, so `-d` would correctly have refused — and the prescribed behavior at that point
is to report the refusal, not to override it. It is recorded here rather than quietly omitted, because
an acceptance run whose own record hides a deviation is worth very little. No real builder branch was
force-deleted; both were removed with plain `-d`, which succeeded.

### Findings (requirement 5 — filed, not fixed)

**F-01 — the hand-back merge cannot run while the owner's claim bookkeeping is uncommitted.**
Reproduced immediately, on the very first merge of the run:
```
$ git merge --no-ff --no-commit worktree-agent-REQ-086-in-progress-record-unstated
error: Your local changes to the following files would be overwritten by merge:
  do-work/working/REQ-085-….md do-work/working/REQ-086-….md do-work/working/REQ-087-….md
Merge with strategy ort failed.                                            (exit 2)
```
Step 2 claims a REQ by moving it from `do-work/queue/` to `do-work/working/` and appending to
`CHECKPOINT.md`; where the consumer **tracks `do-work/`**, that is a staged rename sitting in the index
at exactly the moment Step 6's hand-back sequence says to merge, and git refuses. The hand-back
sequence never says to commit the claim first, and nothing else in the pipeline does either.
Scoped precisely: this bites only installs that commit `do-work/` — on the common untracked install the
claim moves leave the index clean and the merge proceeds. Serial mode never meets it because serial
mode never merges. **Workaround used for this run:** the owner committed the claim bookkeeping and the
run directory as `306c1f4` before merging, which is also what made the run directory durable. Filed as
its own REQ per requirement 5.

**F-02 — nothing in `actions/work.md` drives a wave, so every step of this run was manual.** Step 1
selects one REQ, Step 2 claims one, Step 6 waits for one builder, Step 10 loops after commit. There is
no launch-before-wait path and no wave selection, so an orchestrator following the action file
literally cannot produce the concurrency the worktree-dispatch section documents. This REQ's
Constraints park the fix deliberately ("Do not build the fan-out wave loop as part of this REQ"), and
that is honored — but the gap is now *observed* rather than inferred, which is what makes it fileable.
Filed as its own REQ.

### What this run did not cover (requirement 6)

- **Genuine builder concurrency was not exercised.** Both builders were driven by hand, sequentially,
  in the owner's session — worktree A's work finished before worktree B's began. What that *does*
  prove is everything structural: two worktrees coexisting, two independent branches, non-interference
  across two disjoint file sets, two clean merges into a moving integration tip, serial changelog and
  version allocation, and clean per-branch teardown. What it does **not** prove is any behavior that
  depends on wall-clock overlap — two builders writing at the same instant, a harness-level
  concurrency limit, or a sibling still mid-build when `verify` runs. Reported rather than claimed,
  which is the failure mode REQ-073's own review named.
- **The mid-integration `verify` reading was not captured with a sibling in flight.** Because the
  builds did not overlap in time, there was no moment when one builder was still working while the
  other integrated. REQ-083's classification (an unmerged sibling reports as
  `unmerged-worktree-leftover`, not fixable) is therefore proved by its own fixtures, not by this run.
- **The line-proximity limit was not provoked.** REQ-073's second suggested test — two REQs each
  appending to the same registry, merging cleanly while being jointly wrong — is a different scenario
  from this REQ's negative case (which provokes a real textual conflict). Still unexercised.

### Requirement 8 — final state

`git worktree list` shows only the main checkout; `git branch --list 'worktree-agent-*'` is empty; the
worktree parent directory is gone; the run directory is deleted; no throwaway REQs were left in the
queue (the negative case used branches, not REQ files, so no REQ ids were burned);
`queue-kanban verify --repo-root .` returns `OK: no findings`.

*Verified by work action*

## Discovered Tasks

- **[low]** `tools/checks/qualify.sh` has no representation for a REQ that legitimately changes no
  project file. It requires a `` - `path` `` bullet and separately rejects a `do-work/`-only list, so an
  execution or acceptance REQ — this one, and any future `domain: testing` run-and-record REQ — cannot
  pass mechanically no matter how it is written. The existing design-artifact exception (Step 6.25)
  solves the adjacent case of artifacts placed outside `do-work/`, not this one. Worth deciding whether
  such REQs need a marker the script can read, or whether a hand-recorded FAIL like this one is the
  right cost.
- **[normal]** The hand-back merge fails while the owner's claim bookkeeping is staged, on any install
  that tracks `do-work/` — finding F-01 above, reproduced with its exact error. `actions/work.md`
  Step 2 stages a rename into `working/`; Step 6's hand-back sequence then runs `git merge` against a
  dirty index and git refuses with exit 2. The sequence needs to say what to do about it (commit the
  claim before merging, or stash, or something else) — and whatever it says has to stay correct on the
  untracked install where the problem does not exist.
- **[normal]** `actions/work.md` has no wave-selection or launch-before-wait path, so the fan-out
  concurrency documented in `actions/work-reference.md` → Worktree Dispatch Mode cannot be reached by
  following the action file — finding F-02 above. Deliberately out of scope here (this REQ's
  Constraints park it), and its shape depends on what this run found, which is now recorded.

## Lessons Learned

**What worked:**
- **Running it found things reading it could not.** F-01 is a hard `exit 2` on the first merge of the
  run, in a procedure that had been reviewed, contract-asserted, and shipped. It was invisible to grep
  because neither half is wrong on its own — Step 2's claim is correct, Step 6's merge is correct, and
  only their ordering against a tracked `do-work/` fails. This is the second time in three REQs that a
  documented-but-unexecuted path turned out to be broken (REQ-082 was the first).
- **Using real queue REQs rather than synthetic ones made check 2 meaningful.** Two throwaway builders
  would have exercised the merges just as well, but the changelog/version check only tests anything
  when there are two real deliverables to write entries for — and writing them is where the
  serial-only rule actually bites.
- **REQ-084's probe, shipped an hour earlier, was the instrument for check 3.**
  `git diff --name-only <integration>...<branch> -- do-work/` answered "did this builder write queue
  state" directly, including the committed case. Check 3 would otherwise have been a porcelain glance
  that could not see a committed violation.

**What didn't:**
- **The `-D` slip on the throwaway branches.** After `reset --hard`, `-d` would have refused and the
  correct move was to report that refusal. Reaching for `-D` on branches that felt disposable is
  exactly the reflex the constraint exists to prevent — and "these were only throwaways" is the
  rationalization that makes it feel fine. Recorded in the run above rather than omitted.

**Worth knowing:**
- **`<pre>` is per REQ, not per run, and this run shows why.** The two ranges are `306c1f4..a17e6af`
  and `3ccbf36..5cfe1b5` — different lower bounds, because REQ-086's changelog and archive commits
  landed between the merges. A single run-level `<pre>` would have swept REQ-086's bookkeeping into
  REQ-087's range and misattributed it.
- **A mid-build scope extension can silently un-non-overlap a pair.** Builder B added two files to its
  own `write_set` during the build (REQ-087 D-02/D-03). It stayed disjoint from Builder A here, but
  nothing checked that — the pair was validated at pick time and never re-validated. Under real
  concurrency the merge would catch a textual collision, but the joint-wrongness case it cannot see is
  precisely REQ-073's unexercised line-proximity scenario.
- **The run directory earned its keep for an unexpected reason.** F-01's workaround was to commit the
  claim bookkeeping, and because `do-work/runs/` is a committable path by design, the manifest and both
  hand-backs went into that same commit and survived as durable state rather than as untracked scratch.

## Orientation

**Fan-out dispatch has now actually been run**, and its record is this REQ. Two builders on two
worktrees, two clean merges into a moving integration tip, per-REQ changelog entries with increasing
versions, clean teardown with no `--force`, and an overlapping pair correctly refusing at the merge.
Two defects found and filed: the hand-back merge colliding with the owner's own uncommitted claim, and
the absence of any wave-driving path in `actions/work.md`.

This changes no code and no contract — it converts a documented claim into an observed one, and
narrows what "shipped" means for REQ-073 from *asserted* to *asserted plus these two known gaps*.
`prime_files` is empty; no prime staleness check applies.

## Qualification

**Mechanical check: FAIL, accepted — `tools/checks/qualify.sh` cannot express a no-code-change REQ.**

```
$ ./tools/checks/qualify.sh do-work/working/REQ-085-….md
FAIL: no '## Implementation Summary' file list — run after Step 6.25
```

The script requires at least one `` - `path` `` bullet under `## Implementation Summary`, and rejects a
list of only `do-work/` paths (its design-artifact exception covers artifacts placed *outside*
`do-work/`, which is the opposite shape). This REQ legitimately changed no project file: its
deliverable is the `## Testing` record, and the only files it created — the run directory's manifest
and two hand-backs — are `do-work/` paths that its own requirement 8 required it to delete.

**No manifest was fabricated to clear the gate.** The failure is recorded, the judgment checks were
performed by hand below, and the gap is filed under `## Discovered Tasks`.

**Judgment checks (performed by hand):**

- **Substantive:** the deliverable is real — two worktrees created and torn down, two branches merged,
  two REQs archived with their own commits and changelog entries, one conflict provoked and unwound,
  two defects found and filed. Every claim in `## Testing` carries its command output.
- **Requirements traced:** requirements 1–4 and 6–8 delivered (see `## Testing`, which is organized by
  requirement). Requirement 5 delivered — both defects became REQ-091 and REQ-092, neither was fixed
  inline, and F-01 was hit while it was actively blocking the run, which is when the temptation to fix
  it was highest. Requirement 8's cleanup verified by command; one constraint slip disclosed in the
  record and scored in `## Review`.
- **Flowing:** not a hollow completion — the two builders' work is real, merged, and shipped as
  v0.168.4 and v0.168.5, and the negative case's conflict is reproduced with git's actual output.
- **Contamination check (Step 10):** REQ-085 shares no files with REQ-083/REQ-084 because it changed no
  files at all. Its two builders did touch `tools/queue-kanban/verify.go` (REQ-087) and
  `actions/forensics.md` (REQ-086), both of which REQ-083/REQ-084 changed earlier in this session.
  That overlap is sequential main-tree history, not concurrent work, and the post-merge verification of
  REQ-087 confirms the three REQs' edits to `verify.go` compose (full Go suite green after the merge).


## Review

**Overall: 92%** | 2026-08-04T00:14:18Z

| Dimension | Score |
|-----------|-------|
| Requirements | 94% |
| Code Quality | N/A (no code) |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Partial |

**Findings:** 0 important, 1 minor
**Acceptance:** **Partial** — every structural property of fan-out was exercised and recorded, but the
builders did not overlap in wall-clock time, so genuine concurrency remains unproven. Marked Partial
deliberately: REQ-073's failure was a partial acceptance recorded as complete, and repeating that here
would be the one outcome worse than not running it.
**Suggested testing:** 3 items
**Follow-up REQs created:** REQ-091, REQ-092

**Restatement Sweep:** Not triggered. This REQ redefines nothing — it changes no contract token, no
schema field, no gate wording, and no prescribed command's output shape. Its deliverable is a record.
The two things it *found* redefine nothing either; they are filed as REQ-091 and REQ-092, and any
contract change is theirs to make.

### Findings

**Minor — a constraint was violated during the negative case's cleanup, and the REQ says so.** The two
throwaway overlap branches were deleted with `git branch -D`, against this REQ's explicit "Never `-D`,
never `--force` on any worktree or branch this run creates." After `reset --hard` their commits were
unreachable, so `-d` would have refused and the prescribed response was to report that refusal. No real
builder branch was force-deleted — both came off with plain `-d`, which is the check that mattered.
Scored as Minor rather than Important because it touched only synthetic branches created and destroyed
inside this REQ, it is disclosed in the run record rather than omitted, and the positive case's cleanup
evidence is unaffected. Left as a finding rather than a follow-up REQ: there is nothing to fix in the
tree, only a habit to name.

### Notes on the dimensions

- **Requirements 94%** — seven of eight fully delivered. Requirement 6 (report what the run could not
  cover) is delivered unusually well and is the reason the score is not lower despite Partial
  acceptance: the run states precisely which properties are proven and which are not, rather than
  averaging them into a verdict. The deduction is for the constraint slip above, which sits under
  requirement 8's cleanup discipline.
- **Test Adequacy 90%** — the five positive checks each carry real command output rather than a
  summary, which is what requirement 2 demanded and what makes the record auditable a year from now.
  Deducted because two of REQ-073's three suggested scenarios remain unrun: no `verify` reading was
  captured with a sibling genuinely in flight, and the line-proximity case (two REQs appending to one
  registry, merging cleanly while jointly wrong) was not provoked — that is a different scenario from
  this REQ's textual-conflict negative case, and it is the one the prose explicitly warns about.
- **Scope 100%** — this REQ changed no implementation file, which is exactly its declared scope. The
  builders' scope is accounted for in REQ-086 and REQ-087, each with its own drift check, and both were
  clean. The temptation this REQ was most exposed to — fixing F-01 inline once it blocked the run — was
  declined; the workaround is recorded as an improvisation and REQ-091 explicitly declines to bless it.
- **Risk None** — the integration branch's content after the negative case is byte-identical to before
  it. No worktree, branch, or run directory survives. `verify` is clean.

### Suggested additional testing

- **The same run under genuine concurrency**, once REQ-092 decides whether the pipeline drives a wave
  at all. That is the only way the properties listed under *What this run did not cover* get proven.
- **The line-proximity case, deliberately provoked** — still unexercised after two REQs pointed at it.
  Two builders each appending an entry to the same list; they merge cleanly and are jointly wrong, and
  the claim under test is that the integration-seam rule and post-merge verification catch it. This is
  the failure mode git structurally cannot see.
- **A worktree builder in this repo reading `CLAUDE.md`** — REQ-073's third suggestion, still open.
  Both builders here were driven by the owner, who already knew the version/changelog scope clause, so
  the run did not test whether an independent builder honors it.

*Reviewed by review-work action*

## Full Context

See `do-work/user-requests/UR-016/input.md` for the verbatim instruction, the provenance of the external
audit, and the batch constraints. The test procedure itself is in
`do-work/archive/UR-013/REQ-073-fan-out-dispatch-n-builders-one-owner.md`.

---
*Source: external audit finding F2, third claim (P1) — "REQ-073 itself records that its only live
acceptance test was never run and still marks the ticket completed" — verified against the archived REQ
and accepted by `do-work validate-feedback` triage.*
