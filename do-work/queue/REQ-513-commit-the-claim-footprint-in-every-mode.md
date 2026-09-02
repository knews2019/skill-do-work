---
id: REQ-513
title: '[impact-rule-change] Commit the claim footprint in every mode'
status: pending
created_at: 2026-09-02T20:35:18Z
user_request: UR-099
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-514, REQ-515, REQ-516, REQ-517]
batch: recovery-never-traps
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, _dev/tests/contract-regressions.sh, skills/do-work/tools/do-work-cli/internal/requeststate/]
---

# Commit the claim footprint in every mode

## What

Make Step 2's claim commit its own footprint in every mode, serial included: the queue-to-working move plus the checkpoint entry land as one bookkeeping commit at claim time. The CLI's `claim --commit` already exists; `actions/work.md` Step 2 does not use it, and worktree dispatch Step 0 stages the same moves by hand later.

The fold-first scan found no pending or pending-answers REQ in any UR that owns this claim-commit asymmetry; REQ-505 (Move selection and claim behind advance) moves the claim later but keeps its current commit behavior.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Serial mode defers the claim commit to Step 9 while worktree mode commits it at dispatch. That asymmetry is the whole bug behind the REQ-456 trap: the journaled finalizer only ever meets a dirty checkpoint in serial mode, and its complete transition treats that dirt as foreign. Once claim always commits, complete never meets dirt it made, and no special-case acceptance is needed. This is the delete-before-add fix.

## Context

Origin: REQ-456 (Wait for theme transitions before contrast measurement) finished build and review, then `complete` prepared its journal and refused at lifecycle apply with `FINALIZATION-LIFECYCLE-APPLY`, "target path do-work/CHECKPOINT.md is already dirty". The previous Route A completion, REQ-440, avoided it with a hand-made "[do-work] Record REQ-440 claim checkpoint" commit. REQ-499 through REQ-501 were Route C and never saw it because dispatch committed the claim first. The unblock for REQ-456 was the same hand-made commit, `cd9b01b0`.

## Detailed Requirements

- `actions/work.md` Step 2 invokes `claim` with `--commit` in every mode; the commit message shape is the CLI's own.
- Worktree dispatch Step 0 (`actions/work-reference.md`) stops staging claim moves and the checkpoint by hand, since the claim commit already landed; keep its guard for unrelated staged paths.
- A dirty checkpoint at claim time is shared-target dirt: the claim's refusal stands, and its finding names a verb other than `claim` (REQ-514 owns the general invariant; this REQ only makes sure the claim path has one).
- Contract predicates that pin the deferred-commit wording are deleted with the prose; a behavior test on `claim --commit` covers the serial shape.
- `git log` after a serial run shows one claim commit per REQ, followed by the implementation commit, as worktree mode already does.

## Constraints

- Do not add a second commit surface; use the existing `claim --commit` flag.
- Serial and worktree modes end with the same commit shape for the claim.
- Do not touch the finalizer here; REQ-516 owns what recovery accepts under the sole-authority assertion.

## Batch Constraints

- Judgment stays prose; mechanics stay in the Go CLI. No new prose that walks a shell sequence.
- A guard may still refuse. What it may not do is refuse for a REQ-scoped reason in a way that stops unrelated REQs, or name itself as the fix.
- Nothing here widens recovery to secret-classified or project paths; only dirt the pipeline itself wrote earlier in the run is in scope.
- Every REQ carries a behavior test on the command or a contract predicate on the action, never a sentence pin alone.

## Dependencies

REQ-517 (Pin the serial claim-to-recovery trap) depends on this REQ. Related to REQ-505, which later moves the claim behind `advance` and must inherit the commit.

## Builder Guidance

Certainty level: Firm. The flag exists; the work is wiring, predicate cleanup, and one behavior test. Read `_dev/primes/prime-action-files.md` before touching an action file.

## Red-Green Proof

**RED prompt/case:** Run a serial Route A REQ through Step 2 and inspect `git status` before Step 3.
**Why RED now:** The queue file shows as deleted, the working file as untracked, and `do-work/CHECKPOINT.md` as modified; nothing committed the claim.
**GREEN when:** After Step 2 the tree is clean for those three paths and `git log -1` is the claim's bookkeeping commit; the same is true when dispatch runs the claim in worktree mode.
**Validation:** User confirmed (verify-requests, 2026-09-02).

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3539 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing action routing and status contracts.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2643 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on semantic recovery completeness and structured evidence projection in do-work-cli internals.

## Full Context

See `do-work/user-requests/UR-099/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-02, item A1 of "how can I update the orchestrator to not end up in a trap like this?", captured by UR-099.*
