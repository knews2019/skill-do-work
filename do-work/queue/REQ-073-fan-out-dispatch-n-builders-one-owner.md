---
id: REQ-073
title: "Fan-out dispatch: N concurrent builders under one queue owner"
status: pending
created_at: 2026-08-03T11:41:15Z
user_request: UR-013
domain: general
prime_files: []
tdd: true
suggested_spec:
depends_on: []
addendum_to: REQ-069
maintenance: false
related: [REQ-071, REQ-072]
batch: parallel-builds
write_set: [actions/work-reference.md, actions/work.md, CLAUDE.md, docs/work-guide.md, _dev/tests/contract-regressions.sh]
---

# Fan-Out Dispatch: N Concurrent Builders Under One Queue Owner

## What

Raise Worktree Dispatch Mode from one builder to several concurrent builders under a single queue
owner. The exclusive-session contract is about **who owns `do-work/`**, not about how many builds run,
so a builder that owns nothing shared needs no coordination at all.

## AI Execution State (P-A-U Loop)

- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The user wants several REQs built at once without the checking overhead. Almost all of the machinery
for it already ships — `actions/work-reference.md:235` already specifies sole-integrator, "state stays
home", the four-step hand-back, the per-REQ `<pre>..<merge_hash>` range, name-based crash sweeps, and
cleanup, and the pipeline is worktree-aware throughout (`actions/work.md:377, 388, 424, 536, 572, 581,
694`). **Only the builder count is capped**, by two sentences. This REQ removes that cap and states
the boundary that makes it safe without adding any coordination.

## Context

This REQ is an addendum to **REQ-069** (`do-work/archive/UR-012/`), which adopted the exclusive-session
model at v0.161.0 and deleted ~6,500 words of orchestrator-lock, heartbeat, `claimed_reqs` and
co-dispatch re-validation machinery. That machinery all existed to police **two queue owners**, which
stays banned. This REQ reopens only the cheap half: several builders under one owner. None of the
deleted machinery returns.

- `actions/work-reference.md:237` — "The single active builder…" and "…only one builder is ever in
  flight." The two sentences to change.
- `actions/work-reference.md:53` — `## Execution Model — Exclusive Session`, asserting "one active
  session, one active REQ, one coder context."
- `_dev/tests/contract-regressions.sh:184` — asserts the exact string `one active REQ, one coder
  context` appears **exactly once** across `actions/`.
- `_dev/tests/contract-regressions.sh:241` — `router_word_budget=2650`; `SKILL.md` is at 2,588.
- `crew-members/background-agents.md:3` — loads for "any action that fans work out to background or
  parallel sub-agents", with `work (multi-REQ)` a named caller. Its run-directory pattern therefore
  becomes mandatory here, not optional.
- `actions/board.md:92` — the `overlaps` badge is advisory; **absence reads as unknown, not safe**, and
  it misses glob-vs-glob, `**`, and directory entries.
- `actions/work.md:694` — the ritual already belongs to the integrating commit, not the builder's.

## Detailed Requirements

1. **Lift the builder cap** at `actions/work-reference.md:235-237`: one integrator, several concurrent
   builders. Keep every existing guarantee — sole integrator, state stays home, merge never rebase,
   the four-step hand-back, post-merge verification, and the by-operative-name cleanup.
2. **Add a Fan-Out Dispatch subsection** covering: the human picks which REQs run together; declared
   `write_set` overlaps are surfaced as **advisory input to that pick, never a gate**; and the real
   non-interference proof is `git merge --no-ff --no-commit` refusing to merge.
3. **State the merge gate's honest limit.** Git detects conflicts by line proximity, not meaning. Two
   REQs each adding an entry to a shared registry merge cleanly and can still be wrong. The existing
   integration-seam rule is what covers that, and it only works with a single integrator.
4. **Add a Serial-Only list** naming what never parallelises: queue transitions (claim, status flips,
   archive moves), REQ id allocation, and `actions/version.md` + `CHANGELOG.md` — one changelog entry
   per REQ, written by the owner at merge time.
5. **Reword the invariant** at `actions/work-reference.md:53` from "one active REQ, one coder context"
   to a one-queue-owner formulation, and keep the two-queue-owners ban explicit so none of the deleted
   machinery has a reason to come back. Update the exactly-once assertion at
   `_dev/tests/contract-regressions.sh:184` to the new wording.
6. **Adopt the fan-out durability pattern** from `crew-members/background-agents.md` — a run directory
   created before any spawn, one input and one output file per builder, a manifest, bounded waves, and
   synthesis from files rather than from conversation:

   | Guardrail slot | Fan-out use |
   | --- | --- |
   | run directory (`:26`) | `do-work/runs/work-<YYYY-MM-DD-HHMMSS>/` |
   | per-builder input | `REQ-NNN-brief.md` — REQ body, worktree path, branch name, never-touch list, hand-back format |
   | per-builder output (`:42`) | `REQ-NNN-handback.md` — branch, file manifest, integration seams |
   | `manifest.md` (`:48`) | REQ id → builder, `worktree-agent-REQ-NNN-<suffix>`, handback file, landed status |
   | bounded waves (`:48`) | builders per wave, sized to the harness concurrency limit |

7. **Dispatch mechanism is deliberately unspecified.** Because the owner synthesizes from files and
   never from conversation, a spawned subagent and a human-driven session are indistinguishable to it.
   Do **not** document two separate routes.
8. **Brief delivery path trap.** The brief must reach the builder as prompt content or an **absolute
   main-tree path**. A repo-relative path resolves against the worktree's own stale tracked copy of
   `do-work/`.
9. **State that integration is serial.** Implementation parallelises; merge → qualify → review →
   changelog → archive runs one REQ at a time, and each merge invalidates the previous post-merge
   verification (`actions/work.md:536`). Set the expectation rather than pretending otherwise.
10. **Reconcile `CLAUDE.md`'s "Before Every Commit" section.** It mandates a version bump plus changelog
    entry for *every* commit and auto-loads in any session — including a builder in a worktree —
    pushing builders straight into the serial-only files. `actions/work.md:694` already scopes the
    ritual to the integrating commit. **Fix the instruction, not each brief:** a rule that must be
    overridden by every brief will meet a brief that forgets.
11. **Update the user-facing statement** at `docs/work-guide.md:89` ("One REQ at a time"), which becomes
    false for builders while staying true for queue ownership.

## Constraints

- **No new durable coordination state**, and none of REQ-069's deleted machinery: no lock, no
  heartbeat, no claim registry, no liveness probe, no takeover logic. The forbidden-token sweep at
  `_dev/tests/contract-regressions.sh:132-137` must stay green.
- **No new action file and no new SKILL.md routing row.** Budget headroom is 62 words.
- **`write_set` stays display-only.** Nothing may schedule, gate, or dispatch on it. `tools/queue-kanban/`
  column logic is untouched.
- Carry over `crew-members/background-agents.md:11-14`'s own ceiling note: the durability pattern makes
  fan-out failures **survivable, not prevented**. Do not describe it as a fix.
- Worktrees live outside the repo working tree (`actions/work-reference.md:235`) — a nested second
  checkout is a documented corruption path.
- **A worktree per builder is mandatory, not optional.** The original request offered "with our
  without a new workspace, is a valid variation"; the shared-tree variation was **ruled out** during
  capture and the user did not contest it. Reason: sharing one working tree means every test run,
  qualification check, and review diff in both sessions reads a tree containing the other builder's
  unfinished edits, so the evidence steps stop meaning anything and nothing downstream can tell. The
  staging race is the lesser problem. Keep this rationale in the shipped prose — without it, a future
  reader will re-offer the shared tree as a simplification.
- **"Without letting the other session know that it is running in parallel"** — the user's phrasing,
  and it is satisfied structurally rather than by suppressing information: a builder owns no queue
  state, so there is nothing for it to know. Never implement this as hiding state from a builder that
  could otherwise read it.

## Dependencies

Functionally independent of REQ-071 and REQ-072, but **run it last**: it is the largest contract change,
REQ-071 reduces its blast radius, and REQ-072's `verify` covers the worktree state it creates.

## Builder Guidance

**Certainty level: Firm** on the one-owner boundary, the serial-only list, the advisory-only status of
`write_set`, the mandatory fan-out run directory, and the deliberate silence on dispatch mechanism.

**Exploratory** on the invariant's new wording — it must survive the exactly-once assertion, stay
inside the ≤200-word Execution Model cap REQ-069 established, and read as a boundary rather than a
description. Expect to iterate on that sentence.

This is an instruction-editing REQ on the skill's own files. Prefer **rewording and deleting** over
adding: if the fan-out section grows past what the existing Worktree Dispatch Mode already says, most
of it is redundant and should be cut rather than written. That guidance governs the **shipped prose**,
not this REQ's requirement count — see `## Open Questions`.

## Open Questions

- [x] Eleven detailed requirements against a one-sentence ask carrying a lightness cue ("without the
  anxiety and checking overhead") — keep them all, or trim to the load-bearing ones? → **Keep all
  eleven.** Resolved by the user at verify time (`do-work verify-requests`, 2026-08-03). Each
  requirement traces to a specific failure surfaced during the conversation, and the builder should
  inherit that reasoning rather than rediscover it. The lightness cue was about **runtime** overhead —
  no locks, no heartbeats, no coordination checks — not about how thoroughly the request is specified.

## Red-Green Proof

**RED prompt/case:** Contract-suite assertions, which is the harness REQ-069 itself used:

1. `grep -c "The single active builder" actions/work-reference.md` is 1 today; must be 0.
2. `grep -c "only one builder is ever in flight" actions/work-reference.md` is 1 today; must be 0.
3. A grep for the Fan-Out Dispatch heading and for the Serial-Only list naming `CHANGELOG.md` returns
   nothing today; must match after.
4. The exactly-once invariant assertion at `_dev/tests/contract-regressions.sh:184` still passes
   against the **new** wording (and the old string is gone from `actions/`).
5. The forbidden-token sweep at `:132-137` and the router budget at `:241` both stay green.

**Why RED now:** `actions/work-reference.md:237` states the cap in both sentences, no fan-out section
exists, and the invariant still reads "one active REQ, one coder context."

**GREEN when:** All five assertions pass, and a live two-builder run on two non-overlapping REQs shows:
both branches merge cleanly; each REQ gets its own changelog entry with a strictly increasing version;
`do-work/working/` never holds a file the owner did not put there; `git worktree list` and
`git branch --list 'worktree-agent-*'` are empty after both archives; and the run directory is deleted.
A deliberately overlapping pair must **fail** at `git merge --no-ff --no-commit` rather than merging
silently.

**Validation:** User confirmed — the user asked for "parallel build, but without the anxiety and
checking overhead" and chose "Keep builders out of it" from an explicit option prompt, fixing the shape
as builders that never enter Step 1 and own nothing.

## Full Context

See `do-work/user-requests/UR-013/input.md` for complete verbatim input.
