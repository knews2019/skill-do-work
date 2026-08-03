---
id: REQ-071
title: Crash recovery must respect a live claim before stripping and re-queueing
status: pending
created_at: 2026-08-03T11:41:15Z
user_request: UR-013
domain: general
prime_files: []
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-072, REQ-073]
batch: parallel-builds
write_set: [actions/work.md, actions/work-reference.md, _dev/tests/contract-regressions.sh]
---

# Crash Recovery Must Respect a Live Claim Before Stripping and Re-Queueing

## What

Crash recovery is currently unconditional and destructive. `actions/work-reference.md:220-223`
tells the pipeline that every `REQ-*.md` in `do-work/working/` is its own interrupted leftover, so it
resets the frontmatter, **strips thirteen generated sections**, and moves the file back to
`do-work/queue/`. Make it respect a live claim instead: recover only what this session can show is
its own, and ask a human before taking over anything else.

## AI Execution State (P-A-U Loop)

- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Two independent reasons, and the first one alone justifies the change:

1. **It is wrong for the ordinary single-session crash.** Crash and restart thirty seconds later and
   recovery discards a finished `## Exploration`, a finished `## Plan`, and a declared `## Scope`.
   The pipeline only commits at Step 9, so those sections are almost always uncommitted and
   unrecoverable from git. The trail they form is what `SKILL.md:13` calls the skill's primary value.
2. It is also the single most destructive thing that happens when a second session is started against
   the same checkout — but that is a beneficiary, not the justification. **Write and review this REQ as
   "recovery is too aggressive."**

## Context

- `actions/work.md:116-118` — Step 1 opens with the exclusive-session assumption and then runs
  recovery **before** the CHECKPOINT read and before the queue glob.
- `actions/work-reference.md:220` — the rationale sentence: "every `working/` file is this session's
  own leftover to recover; there is no other live session whose in-flight claim a recovery could
  disturb, so recovery no longer consults any lock."
- `actions/work-reference.md:223` — substep 1 **removes** `claimed_at`. It must be read before it is
  discarded.
- `actions/work.md:225` — `claimed_at` is already written at claim time as a UTC ISO-8601 instant via
  `date -u +%Y-%m-%dT%H:%M:%SZ`, and already notes that a future-dated stamp "freezes the board's
  claim stopwatch."
- `actions/forensics.md:32,39` — Check 1 (Stuck Work) already reads `claimed_at` and already
  prescribes the manual reset. Reuse that definition; do not re-derive it.
- `tools/queue-kanban/future_timestamp_test.go:14-17` — existing skew precedent: 2-minute allowance,
  and "unparseable is not future."
- Nothing here trips the removed-machinery guard: `_dev/tests/contract-regressions.sh:132-137` bans
  exactly `Concurrent-Orchestrator Lock Guard`, `coexisting_sessions`, `claimed_reqs`, `heartbeat_at`,
  and `orchestrator-lock\.json`.

## Detailed Requirements

1. **Discriminate own-crash from foreign claim using the checkpoint file.** A claimed REQ named by
   `do-work/CHECKPOINT.md` is this session's own crash and recovers exactly as it does today. A
   claimed REQ *not* named there is left untouched — no frontmatter reset, no section stripping, no
   move — reported, and offered for takeover only when stale.
2. **Flip the order of the two things Step 1 does.** `actions/work.md:116-118` must read
   `do-work/CHECKPOINT.md` *before* Crash Recovery, since the checkpoint is now recovery's input.
3. **No checkpoint at all is ambiguous — ask, never strip.** Absence of a checkpoint must not be read
   as permission to recover.
4. **Three hours reports; it never authorizes.** A large REQ with review loops can legitimately exceed
   three hours, so the threshold is not a liveness test — it only bounds how long a dead claim goes
   unnoticed. The decision to take over always comes from a human. State this rationale where the
   threshold is defined, so a later edit cannot "simplify" it into an automatic takeover.
5. **Guard the timestamp toward asking.** An unparseable or future-dated `claimed_at` yields a
   negative or meaningless age; treat it as *immediately* eligible for the takeover prompt rather than
   never eligible, or a REQ becomes permanently protected. Follow the existing 2-minute skew
   allowance.
6. **Unattended runs must not block.** With no human to answer, a foreign claim is left alone and
   reported, and the run continues to the next queue item. Never stall the loop on the prompt, and
   never resolve it by stripping.
7. **Read `claimed_at` before substep 1 discards it** (`actions/work-reference.md:223`), the same
   ordering trap that already applies to the `## Scope` / `write_set` decision in that substep.
8. **Update the rationale sentence** at `actions/work-reference.md:220`. Its "there is no other live
   session whose in-flight claim a recovery could disturb" clause is exactly the premise this REQ
   removes.

## Constraints

- **No new durable state.** No lock file, no marker, no liveness counter, no new frontmatter field.
  This REQ reads `claimed_at` and `do-work/CHECKPOINT.md`, both of which already exist.
- The exclusive-session invariant at `actions/work-reference.md:53` is **REQ-073's** to reword. Do not
  touch it here; if this REQ needs the invariant to read differently, say so in `## Lessons Learned`
  rather than editing it.
- Interactive prompts go through `crew-members/clear-questions.md` — one decision per question, and
  the options must state their consequence.

## Dependencies

None. Deliberately first of the batch: smallest, independently justified, and it reduces the blast
radius of REQ-073.

## Builder Guidance

**Certainty level: Firm** on all six, but the provenance differs and that matters if one has to be
traded off. **User-stated, verbatim:** leave a claimed ticket claimed, the three-hour threshold, and
asking before takeover. **Analysis-derived and approved via the plan:** the checkpoint discriminator,
the timestamp guard, and the non-blocking unattended path. If any requirement has to give, it comes
from the second group — never the first.

Latitude on: where the threshold constant is stated and how it is worded; the exact prompt text; and
how the report line is formatted. Keep it simple — this is prose in two action files plus assertions,
not machinery.

## Red-Green Proof

**RED prompt/case:** Two probes, both runnable today.

1. Contract suite: an assertion that recovery is conditional on a checkpoint match — e.g. that
   `actions/work-reference.md` no longer claims "there is no other live session whose in-flight claim
   a recovery could disturb", and that `actions/work.md` Step 1 reads CHECKPOINT before recovery.
   Fails on the current tree.
2. Behavioural: hand-create `do-work/working/REQ-999-probe.md` with `status: claimed`, a
   `claimed_at` ten minutes in the past, and a `## Plan` section, alongside a `do-work/CHECKPOINT.md`
   naming a *different* REQ. Start the pipeline.

**Why RED now:** Probe 2 today strips `## Plan` and moves the file to `do-work/queue/`, because
`actions/work-reference.md:220` instructs recovery to treat every `working/` file as its own.

**GREEN when:** Probe 1's assertions pass. Probe 2 leaves `do-work/working/REQ-999-probe.md`
byte-identical, reports it as a foreign claim, and — being under three hours old — does not even offer
takeover. Re-stamping `claimed_at` to four hours ago offers takeover and still strips nothing until a
human says yes. Re-stamping it to a future instant also offers takeover rather than protecting the
file forever.

**Validation:** User adjusted — the user proposed the claim-respecting behaviour and the three-hour
ask ("basically when a ticket is claimed, leave it claimed, if it's older then 3h then ask if it
should be taken over"). The checkpoint discriminator, the timestamp guard, and the non-blocking
unattended path were added during capture and approved in the plan.

## Full Context

See `do-work/user-requests/UR-013/input.md` for complete verbatim input.
