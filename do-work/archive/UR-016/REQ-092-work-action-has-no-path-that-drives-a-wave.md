---
id: REQ-092
title: actions/work.md has no wave-selection or launch-before-wait path, so documented fan-out concurrency cannot be reached by following it
status: completed
claimed_at: 2026-08-04T00:27:00Z
completed_at: 2026-08-04T00:27:44Z
commit: 92bebe0
kb_status: promoted
kb_entry: REQ-092-actions-work-md-has-no-wave-selection-or.md
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

## Triage

**Route: A** - Simple

**Reasoning:** The gap is quoted and verifiable in one read. The work is a single design decision
(requirement 1) plus stating it in two files. No discovery, no plan.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `actions/work.md` (modified)
- `actions/work-reference.md` (modified)

**What was done:** Answered requirement 1 with **no** — `actions/work.md` should not drive a wave — and
stated that boundary where readers currently infer the opposite. `actions/work.md` gained a paragraph
next to the sub-agent note saying it processes one REQ at a time deliberately, naming the four steps
that assume it, and pointing at `actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch
as the owner-driven procedure that does the other thing. The `--wave N` section gained a paragraph
saying the flag selects a batch and does not run one concurrently. `actions/work-reference.md`'s
Fan-Out Dispatch section gained the matching statement from the other side: this is a procedure a human
or advanced harness follows, not something the action performs.

## Decisions

- **D-01 (ESCALATE-tier reasoning, decided)** — *`actions/work.md` should not drive a wave; the boundary
  is documented instead.* This is requirement 1's answer, and the REQ explicitly allows it as a success
  rather than a dodge. Four reasons, in order of weight:
  1. **The floor comes first.** The skill's compatibility rule is that the simplest agent which can read
     and write files and run shell commands must be able to follow an action file end to end. A
     claim-many/dispatch-many/integrate-one path would sit in the main loop in front of every reader who
     cannot use it, to serve an explicitly "optional, advanced harnesses only" capability.
  2. **The dispatch mechanism is deliberately unspecified** (REQ-073 requirement 7) so that a spawned
     subagent and a separately-opened session are indistinguishable to an owner that synthesizes from
     files. A wave loop written into `actions/work.md` would have to name one and break that property.
  3. **The payoff is bounded and already known.** Integration is serial at any builder count, so
     fan-out buys wall-clock in the build phase only — REQ-085 measured exactly this. That is a poor
     trade against complexity in the path everyone reads.
  4. **The capability is not unreachable, only un-automated.** It is fully specified in
     `work-reference.md` and was executed end to end in REQ-085. The actual defect was that nothing said
     which of the two documents owns it, so a reader met `--wave` and inferred a concurrency feature.
  **Value:** `actions/work.md` stays followable by the floor, and the two documents stop disagreeing.
  **Risk:** fan-out remains manual, so it will be used rarely and its rough edges will surface slowly —
  REQ-085 found two on first contact. Fully reversible: the decision is three paragraphs, and REQ-085's
  record is the input any future "actually, build it" REQ would start from.
- **D-02 (DECIDE & STATE)** — *`--wave N` is kept, with its scope clarified, rather than renamed or
  removed.* Requirement 3 offered "becomes the wave's selector" or "stops implying concurrency". With
  D-01 answering no, only the second is available — but the flag does something genuinely useful on its
  own: a depth-N set is mutually independent, so the run cannot be derailed mid-batch by a dependency.
  Kept and described as a scoping flag.

## Qualification

Passed — 2 files verified, 5 requirements traced.

- **Requirements traced:** 1 → decided (no) and stated in both files, with the reasoning in D-01. 2 →
  not applicable, since requirement 1 resolved to no; nothing was specified for a wave loop to get
  wrong. 3 → `--wave N`'s paragraph now says it selects and does not parallelise; D-02 records why the
  flag survives. 4 → nothing was added: no lock, heartbeat, claim registry or liveness probe, and the
  per-REQ guarantee structure is untouched. 5 → satisfied by construction — the decision is precisely
  that the floor's serial loop stays the only path in `actions/work.md`.
- **Restatement Sweep:** the two documents were the only ones describing this, and they now agree
  explicitly rather than by omission — which was the defect. Checked `docs/work-guide.md`'s "One at a
  time" bullet (line 120): it already tells the user the action processes one REQ per loop iteration and
  is now corroborated rather than contradicted. `CLAUDE.md`'s agent-compatibility rule ("design for the
  floor") is the premise this decision rests on, unchanged.
- **Substantive:** three added paragraphs, all load-bearing statements of ownership; no filler.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ No new failures — 8 FAIL lines, the same pre-existing update-script probe failures
recorded under REQ-083's Discovered Tasks.

Prose-only change to action-file instructions; no behavioral code, so red-green does not apply. The
observation this REQ responds to is REQ-085's finding F-02, recorded with the four one-at-a-time step
citations that make the gap verifiable in one read.

*Verified by work action*

## Lessons Learned

**What worked:** Taking "document the boundary" as a real answer. The REQ was written to make that
outcome available without a builder having to justify departing from an implied build, and it was the
right one — the capability was never unreachable, only unowned by either document.

**Worth knowing:** The tell that this was a documentation defect rather than a missing feature: the
capability had been *executed successfully* (REQ-085) before anyone tried to automate it. A feature you
can perform by hand from the spec is specified; what was missing was a sentence saying which document
performs it.

## Orientation

**Now the two documents agree about who runs a fan-out.** `actions/work.md` states that it processes one
REQ at a time by design and points at the owner-driven procedure for the other thing; `--wave N` says it
selects a batch rather than running one concurrently; and Worktree Dispatch Mode says from its side that
it describes something a human or advanced harness performs. A reader who met `--wave` and inferred a
concurrency feature no longer can.

No logic changed and no contract moved — this closes the last of the two defects REQ-085's live run
found, by deciding ownership rather than by building machinery.

## Full Context

`do-work/archive/REQ-085-run-the-live-two-builder-acceptance-test.md` → `## Testing` → finding F-02.
