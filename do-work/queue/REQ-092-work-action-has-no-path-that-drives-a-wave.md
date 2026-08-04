---
id: REQ-092
title: actions/work.md has no wave-selection or launch-before-wait path, so documented fan-out concurrency cannot be reached by following it
status: pending
created_at: 2026-08-04T00:14:18Z
user_request: UR-016
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-091]
maintenance: false
related: [REQ-073, REQ-085]
addendum_to: REQ-085
discovered_during: REQ-085
write_set: [actions/work.md, actions/work-reference.md]
---

# `actions/work.md` Has No Path That Drives a Wave

## What

`actions/work-reference.md` → Worktree Dispatch Mode → **Fan-Out Dispatch** documents several builders
running under one queue owner, with integration serialised behind them. `actions/work.md` has no step
that produces that shape:

- **Step 1** finds *the next* pending REQ.
- **Step 2** claims *one* REQ.
- **Step 6** spawns a builder and waits for it before Step 6.25 reads its output.
- **Step 10** loops back to Step 1 *after* the commit.

Every path is one-REQ-at-a-time. An orchestrator following the action file literally builds serially no
matter how many builders the harness could run, so the documented concurrency is reachable only by
departing from the instructions.

## Why

REQ-085's live acceptance run had to drive both builders by hand — every claim, dispatch, merge and
teardown was manual, because the action file describes no other way. That run is the evidence: the
fan-out section is a description of a shape the pipeline cannot currently produce on its own.

This is a gap between two shipped documents, not a broken feature. Nothing is corrupted and serial mode
is unaffected. The cost is that `--wave N` exists in Step 1 (it filters pending REQs by dependency
depth) while nothing downstream ever runs a wave *concurrently* — so the flag reads as a concurrency
feature and delivers a filtered serial run.

## Context

`do-work/archive/REQ-085-…md` → `## Testing` → finding F-02, and *What this run did not cover*, which
records precisely which properties the hand-driven run did and did not establish.

REQ-085's own Constraints parked this deliberately: "Do not build the fan-out wave loop as part of this
REQ … adding one is a separate change whose shape depends on what this run finds." That run has now
happened, so the input this REQ was waiting on exists.

## Detailed Requirements

1. **Decide, and state, whether `actions/work.md` should drive a wave at all.** A legitimate outcome is
   *no* — that fan-out stays an owner-driven procedure a human or a harness performs, with the action
   file documenting it rather than executing it. If so, say it in `actions/work.md` where a reader
   currently infers the opposite from `--wave`, and the REQ is a documentation change.
2. **If it should:** specify claim-many, dispatch-many, then integrate-one-at-a-time, honoring every
   existing per-REQ rule — one `<pre>..<merge_hash>` per REQ, one changelog entry per REQ, queue
   transitions and version bumps serial-only.
3. **Reconcile `--wave N` with whatever is decided.** Today it filters by dependency depth and then
   runs the filtered set serially. Either it becomes the wave's selector or its documentation stops
   implying concurrency.
4. **Do not weaken the exclusive-session model or add coordination state.** No lock, heartbeat, claim
   registry, or liveness probe. Fan-out's guarantees are per REQ precisely so that raising the builder
   count adds no coordination, and that property must survive.
5. **Whatever is specified must be runnable by the floor** — an agent that can read/write files and run
   shell commands. Subagents and parallelism are nice-to-haves, so the wave path must degrade to the
   serial loop rather than requiring concurrency.

## Constraints

- **This is a design REQ before it is an implementation REQ.** Requirement 1's answer may be "document
  the boundary, change no logic" — that is a success, not a dodge.
- Integration stays serial at any builder count.
- `depends_on: [REQ-091]` — the hand-back merge's index problem is inside any wave loop this would
  specify; specifying the loop on top of a sequence that cannot merge would bake the defect in.

## Dependencies

`depends_on: [REQ-091]`. `addendum_to: REQ-085`, which observed the gap. `related: REQ-073`, which
introduced the fan-out documentation.

## Builder Guidance

**Certainty: Firm on the gap, entirely open on the response.** The four one-at-a-time paths are quoted
above and verifiable in one read of `actions/work.md`. Whether the fix is prose or machinery is the
actual work, and requirement 1 exists so that "document it" is available as a real answer rather than
something a builder has to justify departing from.

Read REQ-085's *What this run did not cover* first — it names exactly which fan-out properties are
proven and which remain unexercised, which is the honest starting inventory for this decision.

## Full Context

`do-work/archive/REQ-085-run-the-live-two-builder-acceptance-test.md` → `## Testing` → finding F-02.
