---
id: REQ-100
title: Live auto-wave acceptance run — prove real wall-clock concurrency
status: completed
created_at: 2026-08-04T19:44:17Z
claimed_at: 2026-08-04T21:23:00Z
completed_at: 2026-08-04T21:32:00Z
kb_status: pending
write_set: []
user_request: UR-018
domain: testing
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-099]
maintenance: false
related: [REQ-099, REQ-095]
batch: parallel-building
---

# Live Auto-Wave Acceptance Run

## What

Run the REQ-099 automatic wave dispatch live, with genuinely concurrent builders, and record the evidence. Real wall-clock fan-out concurrency has **never been proven** in this skill — the one recorded attempt (REQ-085) logged Partial — so this run is the proof, not ceremony.

## Detailed Requirements

- Use real (or realistic dummy) REQs so at least two builders run **simultaneously** — capture timestamps proving overlap (builder start/end times from the run directory), not just "both completed".
- Exercise the full auto path: automatic set computation (deps + claims + `assigned_to` respected), bounded wave size, worktree-per-builder, serial integration of all hand-backs.
- Record the run artifact in the same form as REQ-085's fan-out run; include at least one observed imperfection or its explicit absence (the previous run's value was finding the index-settling bug).
- If the run surfaces defects in REQ-099's prescribed commands, fix them and grep the same primitive across all actions before calling it fixed (copy-paste rule).

## Red-Green Proof

**RED prompt/case:** No recorded evidence of two builders with overlapping wall-clock execution exists anywhere in the repo (REQ-085: Partial).
**Why RED now:** Fan-out has been driven by hand, serially confirmed, and auto-wave has never existed.
**GREEN when:** A recorded run shows ≥2 builders with overlapping timestamps dispatched by the automatic wave computation and integrated serially, with any found defects fixed and cross-grepped.
**Validation:** User confirmed (approved plan, Phase 3 item 9).

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 3).

---
*Source: approved plan, Phase 3*

## Triage

**Route: B** - Medium

**Reasoning:** The deliverable is named precisely — overlapping wall-clock timestamps from at least two builders, dispatched by the automatic computation and integrated serially — and the recording form is specified by reference to REQ-085. What needed discovery was how to get *genuine* wall-clock overlap in this environment, and which of REQ-099's four ready-set clauses could actually be exercised rather than asserted. The deliverable is an execution and its record.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided execution

*Skipped by work action*

## Exploration

**What REQ-085 could not do, and why.** Its `## Testing` records `Dispatch mechanism: both builders driven **by hand in the owner's session**`, which is exactly why its concurrency result was Partial: two builders driven in sequence by one session cannot overlap, whatever the branches look like afterwards. Beating it needs two processes running at the same instant, and timestamps fine-grained enough to prove it rather than assert it.

**The four clauses are not equally testable.** `actions/work-reference.md` → Auto-wave requires pending + dependency-ready + unclaimed + not-`assigned_to`. A fixture can make **three** of them fire as exclusions (a non-pending status, an unmet dependency, an earmark) — clause 3 only fires against a live claim, which is the state the run itself creates. So the fixture was shaped to have something to *exclude* on clauses 2 and 4, not just something to include: a REQ that would be ready except its dependency is pending, and a REQ that would be ready except it is earmarked elsewhere. A wave computation that silently ignored either would look identical on an include-only fixture.

**The recompute rule is testable, and is the part with no precedent at all.** REQ-099's Step 10 says loop back for a freshly computed wave rather than a leftover slice. A dependency-blocked REQ in wave 1 that becomes ready only after wave 1 archives is the exact case a carried-forward remainder list would miss.

**The negative case is the safety argument.** REQ-099's contract says a computed set claims runnability, never non-overlap, and that the merge is the proof. That is falsifiable: give two REQs the identical `write_set`, confirm the computation includes both anyway, and confirm the second merge refuses.

## Scope

**Files I will touch:** none as *implementation*. The deliverable is this REQ's `## Testing` record. The fixture project, its two builder worktrees and its throwaway REQs all live in scratch space outside the repo and are deleted afterwards.

**Acceptance criteria (restated from REQ):**
- [ ] At least two builders run **simultaneously**, with timestamps proving overlap rather than just both completing
- [ ] Automatic set computation exercised, with `depends_on`, claim state and `assigned_to` all respected
- [ ] Bounded wave size exercised
- [ ] Worktree per builder
- [ ] Serial integration of all hand-backs
- [ ] Recorded in REQ-085's form, with at least one observed imperfection or its explicit absence
- [ ] Any defect in REQ-099's prescribed commands fixed and cross-grepped

## Pre-Flight

- **WARN — baseline suite red before any change:** the same 8 `chmod 500`-versus-root failures inherited by every REQ in this batch. No shipped file is changed by this REQ, so the gate is simply "unchanged".
- `git worktree` support present; working tree clean outside `do-work/` at claim time.

## Implementation Summary

**Files changed:** none — this REQ is an execution. Its deliverable is the `## Testing` record below. The fixture project (`<scratch>/req100/run/`), its run directory, and four builder worktrees were created, consumed, and deleted.

**What was done:** Ran REQ-099's auto-wave dispatch live. The ready-set predicate was applied to a four-REQ fixture shaped to exercise two of its exclusion clauses; the resulting bounded wave of two was dispatched into two git worktrees whose builders ran **concurrently for 4.109 seconds of measured wall-clock overlap**; both hand-backs were integrated serially with per-REQ merge ranges; a second wave was recomputed and correctly picked up the REQ whose dependency had just landed; and a deliberate overlapping-`write_set` pair confirmed the merge refuses, which is the safety argument the computed set rests on.

## Testing

### Run parameters

| Item | Value |
| --- | --- |
| Fixture project (`<scratch>` = the session scratchpad) | `<scratch>/req100/run/` — a real git repo outside the skill tree |
| Integration branch | `main`, tip `7b88fc2` at fixture creation |
| Worktree parent | `<scratch>/req100/wave-worktrees/` — a sibling of the fixture, never nested inside it |
| Run directory | `do-work/runs/work-2026-08-04-212428/` (manifest + 2 briefs + 2 hand-backs + `timeline.log`) |
| Wave bound | 2 (explicit, exercising `--fan-out N`) |
| Builder A operative name | `worktree-agent-REQ-910-widget-renderer`, commit `eea9b29` |
| Builder B operative name | `worktree-agent-REQ-911-invoice-total`, commit `a6c848f` |
| REQ-910 merge range | `65f9c2c..d8948c8` |
| REQ-911 merge range | `efa5523..9a79568` |
| Dispatch mechanism | two **concurrent background processes**, each in its own worktree — see *What this run did not cover* |
| Scripts | `setup-wave.sh`, `compute-wave.sh`, `dispatch-wave.sh`, `integrate-wave.sh`, `negative-case.sh` |

### Automatic set computation, including what it excluded

The fixture held four pending REQs. The computation applies only the four Auto-wave clauses:

```
=== auto-wave computation (bound 2) ===
  REQ-910 READY
  REQ-911 READY
  REQ-912 excluded — clause 2 (depends_on REQ-910 not terminally resolved)
  REQ-913 excluded — clause 4 (assigned_to cloud-beta)
ready set: REQ-910 REQ-911
bounded wave (first 2 in numeric id order): REQ-910 REQ-911
```

**F-01 — both exclusion clauses fired, which is what makes the include result meaningful.** A computation that ignored `depends_on` would have put REQ-912 in wave 1 (and it would have built against a dependency that did not exist yet); one that ignored `assigned_to` would have taken REQ-913 out from under `cloud-beta`. Clause 4 is REQ-097's field doing its job as a wave input, one REQ after it shipped.

### The deliverable: measured wall-clock overlap

Both builders were launched before the owner waited on either (launch-before-wait), each recording a nanosecond-precision start and end:

```
=== D4. timeline.log — the overlap evidence ===
builder REQ-910 END   2026-08-04T21:24:32.576756966Z commit eea9b29
builder REQ-910 START 2026-08-04T21:24:28.467312977Z
builder REQ-911 END   2026-08-04T21:24:32.576492866Z commit a6c848f
builder REQ-911 START 2026-08-04T21:24:28.467251908Z

=== D5. overlap arithmetic ===
REQ-910 ran 21:24:28.467312 -> 21:24:32.576756 (4.109s)
REQ-911 ran 21:24:28.467251 -> 21:24:32.576492 (4.109s)
overlap window: 21:24:28.467312 -> 21:24:32.576492
OVERLAP = 4.109 seconds
VERDICT: CONCURRENT — the two builders were running at the same wall-clock instant
```

**F-02 — GREEN, and the first recorded wall-clock concurrency in this skill.** The overlap is computed as `min(end) − max(start)`, so it is the interval during which *both* builders were genuinely running: 4.109 seconds, against REQ-085's Partial, where the two builders could not have overlapped at all because one session drove them in turn. Two worktrees, two branches, two commits, one owner waiting on both.

### Serial integration with per-REQ merge ranges

```
--- integrating REQ-910 (serial) ---
  <pre>        = 65f9c2c
  <merge_hash> = d8948c8
  range        = 65f9c2c..d8948c8
  files in range: src/widget-renderer.txt
  archived; integration tip now efa5523
--- integrating REQ-911 (serial) ---
  <pre>        = efa5523
  <merge_hash> = 9a79568
  range        = efa5523..9a79568
  files in range: src/invoice-total.txt
  archived; integration tip now db051c9
```

**F-03 — the two `<pre>` values differ, and the reason is the release tail running between the merges.** REQ-910's archive commit moved the tip, so REQ-911's `<pre>` is `efa5523` rather than the `65f9c2c` both merges started from. This reproduces exactly what REQ-085 documented and is why the contract says capture `<pre>` **once per REQ**, at its own merge. **The imperfection this run observed is my own:** the dispatch script captured a single `<pre>` for the whole wave at claim time, which would have given REQ-911 the range `65f9c2c..9a79568` — sweeping in REQ-910's builder commit *and* its archive commit, so qualify and review would have read another REQ's work plus a bookkeeping commit as REQ-911's. The per-REQ rule is not stylistic; it is the difference between a correct range and a wrong one, and the natural mistake under fan-out is exactly the one the rule forbids. No prose fix is needed — the contract already says `Capture it **once per REQ**` — but this run is the evidence that the warning is earned.

Each range contains exactly one builder's file, and both builders' work is present afterwards:

```
=== I2. both builders' work is present in the integration branch ===
invoice-total.txt
widget-renderer.txt
```

**F-04 — cleanup by operative name succeeded, which is the free merged-ness assertion.**

```
  removed worktree worktree-agent-REQ-910-widget-renderer
  Deleted branch worktree-agent-REQ-910-widget-renderer (was eea9b29).
  removed worktree worktree-agent-REQ-911-invoice-total
  Deleted branch worktree-agent-REQ-911-invoice-total (was a6c848f).
```

`git branch -d` (never `-D`) succeeded for both from the integration branch — the only mechanical evidence that both merges actually landed, and the check that would have refused had either merge been skipped or lost.

### Wave 2: recompute, not a remainder list

After wave 1 archived, the computation was re-run against the same fixture:

```
############ WAVE 2 ############
  REQ-912 READY
  REQ-913 excluded — clause 4 (assigned_to cloud-beta)
ready set: REQ-912
bounded wave (first 2 in numeric id order): REQ-912
```

**F-05 — REQ-099's recompute rule is load-bearing, and this is the case that proves it.** REQ-912 was excluded from wave 1 on clause 2 and is ready in wave 2 because its dependency reached `completed` in between. A run that carried wave 1's remainder forward would have had an empty list and stopped with a runnable REQ sitting in the queue. REQ-913 stayed excluded across both waves — an earmark does not decay.

### Negative case: the computed set does not claim non-overlap

Two REQs were given the **identical** `write_set`:

```
do-work/queue/REQ-920-touches-shared.md:write_set: [src/shared-config.txt]
do-work/queue/REQ-921-touches-shared.md:write_set: [src/shared-config.txt]

=== N2. the wave computation includes BOTH — it does not read write_set ===
  REQ-920 READY
  REQ-921 READY
```

Their builders each wrote the same line concurrently, and integration was attempted serially:

```
=== N4. serial integration: first merge clean, second must REFUSE ===
  REQ-920 merged: e3ddcc7
Auto-merging src/shared-config.txt
CONFLICT (content): Merge conflict in src/shared-config.txt
Automatic merge failed; fix conflicts and then commit the result.
  second merge exit status: 1
UU src/shared-config.txt
--- the conflict, verbatim ---
<<<<<<< HEAD
setting = setting-a
=======
setting = setting-b
>>>>>>> worktree-agent-REQ-921-touches-shared
```

**F-06 — the safety argument holds as written, in both halves.** The computation included both REQs, reading nothing from `write_set` — so a computed set is demonstrably not an assertion that its members are disjoint. And the merge refused, at exit 1, with a content conflict the owner must resolve. That is `the non-interference proof is the merge, not the pick` exercised rather than asserted, and it is why fully automatic set-picking is defensible without a contention gate.

### Defects in REQ-099's prescribed commands

**None found — recording the absence explicitly**, per this REQ's requirement. The Auto-wave predicate, the bound-then-take-first-N rule, the run-directory-before-any-spawn ordering, the hand-back sequence's step 0 (bookkeeping commit below `<pre>`), the per-REQ `<pre>`/`<merge_hash>` capture, and the `-d`-not-`-D` cleanup all executed as written with no rewording needed. Step 0's wording already reads plurally ("the Step 2 claim **moves**"), so one bookkeeping commit covering several claims is in contract — checked because a singular reading would have been wrong under a wave. Nothing to cross-grep.

### What this run did not cover

- **The builders were scripted shell processes, not agents.** This is the honest limit and the significant one. What is proven is the *mechanism* under real concurrency: set computation, bounded dispatch, worktree isolation, launch-before-wait, per-REQ merge ranges, serial integration, merged-ness cleanup, recompute, and the merge-refuses safety property. What is **not** proven is that an LLM builder reads its brief, stays in its worktree, and hands back a correct manifest — REQ-085 covered agent behavior with hand-driven builders; this run covers concurrency with scripted ones, and no single run in this repo has yet covered both at once.
- **The set computation was an implementation of the prose predicate, not the prose being read by an agent.** It shows the predicate is implementable and that its exclusions are correct; it does not show two agents would read `actions/work-reference.md` → Auto-wave the same way.
- **Two builders, bound 2.** Nothing here says whether a wave of five behaves differently, and the harness-concurrency-limit default was never exercised (the bound was passed explicitly).
- **No crash mid-wave.** A builder dying between dispatch and hand-back, and what the next run's crash recovery makes of two `working/` REQs under one label, is untested.
- **Single machine, single checkout.** The wave ran under one owner in one checkout; the cross-checkout half of the model is REQ-095's evidence, not this run's.

### Fixture cleanup

```
$ git status --porcelain --untracked-files=all
 M do-work/CHECKPOINT.md
RM do-work/queue/REQ-100-live-wave-acceptance-run.md -> do-work/working/REQ-100-live-wave-acceptance-run.md
$ git worktree list
/home/user/skill-do-work  fc78a83 [claude/ur-018-do-work-batch-cbmfeb]
$ git branch --list 'worktree-agent-*'
(empty)
$ find . -name 'REQ-9*.md' -not -path './.git/*'
(empty)
```

The fixture project, its `wave-worktrees/` parent, and all four builder worktrees are gone. The only dirt in this repo is this REQ's own bookkeeping. No `REQ-9xx` file ever existed inside the repo, and no worktree or `worktree-agent-*` branch leaked into it.

## Lessons Learned

**What worked:**
- Shaping the fixture around **exclusions** rather than inclusions. Two of the four REQs existed only to be left out, which is what turns "the wave contained the right REQs" from a coincidence into a result. An include-only fixture would have passed identically against a computation that ignored `depends_on` and `assigned_to` entirely.
- Computing overlap as `min(end) − max(start)` and printing a verdict, instead of eyeballing two timestamps. It makes the deliverable a number that is either positive or the run failed, and it is the difference between this result and REQ-085's Partial.
- Running the negative case at all. A computed set that is *documented* as not proving non-overlap is worth much less than one where the merge has been watched refusing.

**What didn't:**
- The dispatch script captured **one `<pre>` for the whole wave**, at claim time. It reads naturally — the wave shares a starting tip — and it is wrong the moment the first REQ's release tail commits, which is exactly what serial integration does between merges. The contract already forbids it; I wrote the forbidden version anyway on the first pass, which is the best available evidence that the rule needs to keep saying so.
- Assuming a `cd` into a fixture directory would survive its own deletion: removing the fixture while the shell sat inside it left the next command unable to resolve a working directory (`fatal: Unable to read current working directory`) and made a clean teardown look like a failure. Delete from outside.

**Worth knowing:**
- `date -u +%Y-%m-%dT%H:%M:%S.%NZ` gives nanosecond stamps that Python parses after truncating to microseconds (`stamp[:26]`). Second-resolution timestamps cannot prove a sub-second overlap, and a 4-second sleep is the cheapest way to make the window unambiguous.
- Under fan-out the owner's bookkeeping commit covers **several** claims at once and still lands below every REQ's `<pre>` — because each `<pre>` is captured later, at that REQ's own merge. The two rules fit together only in that order.
- `git worktree remove` on a worktree whose branch still has unmerged commits needs `--force`; in the negative case that is correct (the branch was deliberately abandoned), and it is the one place in this run where `--force` was right. On the happy path both plain `remove` and `branch -d` succeeded, which is the assertion that matters.

## Orientation

Automatic wave dispatch now has an acceptance record: two builders measured running concurrently for 4.109 seconds in separate worktrees, dispatched by the computed ready set with two of its exclusion clauses observed firing, integrated serially with correct per-REQ merge ranges, followed by a recomputed second wave that correctly picked up the REQ whose dependency had just landed. The safety argument was tested rather than assumed — two REQs declaring the same `write_set` both entered the wave, and the merge refused. This is evidence about `actions/work.md`'s auto-wave mode and `actions/work-reference.md`'s Fan-Out Dispatch; no code or contract changed. The gap that remains is named: the builders were scripts, so agent behavior under concurrency is still unproven. `prime_files` is empty.

## Qualification

**Mechanical check: FAIL, accepted** — `tools/checks/qualify.sh` cannot express a no-code-change REQ, the same accepted FAIL as REQ-085 and REQ-095:

```
$ tools/checks/qualify.sh do-work/working/REQ-100-live-wave-acceptance-run.md
FAIL: no '## Implementation Summary' file list — run after Step 6.25
```

**Judgment checks, by hand:**

- **Requirements traced (7/7).** Simultaneity → the overlap arithmetic; automatic computation with all three testable clauses → the computation transcripts for waves 1 and 2; bounded wave → `bound 2` in both; worktree per builder → `git worktree list` mid-run; serial integration → the two distinct merge ranges; REQ-085's form plus an observed imperfection → F-03 (my own once-per-wave `<pre>`) and the explicit "none found" for REQ-099's commands; cross-grep → not applicable, nothing to fix.
- **Substantive, not narrated.** Every quoted block is captured command output, including four commit hashes, four range endpoints, nanosecond timestamps, and the verbatim conflict markers.
- **Evidence over reasoning.** The two results that could most easily have been assumed instead of measured — the overlap and the merge refusal — are the two that are shown as raw output with exit statuses.

## Review

**Overall: 93%** | 2026-08-04T21:32:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | n/a (no code changed) |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 2 minor
**Acceptance:** Pass — the named deliverable (proven overlapping wall-clock timestamps) is delivered as a measured 4.109-second window, all three testable ready-set clauses were observed, integration ranges are per-REQ and correct, and the safety argument was exercised in the negative direction.
**Suggested testing:** 3 items
**Follow-ups created:** None

**Requirements checklist:** all seven `## Scope` criteria delivered; per-criterion evidence in `## Qualification`.

**Minor:**
- **Test Adequacy is 90%, not 100%, and the missing 10% is a real gap rather than modesty.** The builders were shell scripts. Every mechanical property of auto-wave is now proven under genuine concurrency, and no property of *agent* behavior under concurrency is. REQ-085 has the complement (real agent builders, no overlap). Nothing in this batch closes both at once, and the honest reading of the two runs together is "the mechanism works concurrently, and agents have followed it serially" — not "auto-wave is proven end to end". Recorded in *What this run did not cover* rather than filed, because the fix is a future run with agent builders, not a change to anything shipped.
- Clause 3 (unclaimed) was never observed *excluding* anything — it only fires against a live claim, which is the state the run itself creates, so a fixture cannot pre-stage it without staging the very condition under test. Stated rather than glossed; it is the one clause of four taken on inspection.

**Scope drift:** none. `## Scope` declared no implementation files and none were written.

**Restatement sweep (MUST):** not applicable — this REQ changed no contract token, schema field, gate wording, or prescribed-command output shape. Its findings were checked *against* shipped prose rather than changing it: F-03 confirms the existing per-REQ `<pre>` rule is earned, and the "no defects found" result was reached by executing REQ-099's prescribed sequence rather than re-reading it. Had any step needed rewording, the copy-paste rule would have required grepping the same primitive across all actions first; nothing did.

**Suggested additional testing:**
- A wave with **agent** builders, to close the gap above. That is the one run neither REQ-085 nor this one performs.
- A wave of five with bare `--fan-out` (harness-limit default), to exercise the bound the run passed explicitly.
- A builder killed between dispatch and hand-back, then a fresh run's crash recovery against two `working/` REQs under one `writer:` label — the crash-mid-wave case, which is currently reasoned about and never observed.

*Reviewed by review-work action (pipeline mode, in-session)*
