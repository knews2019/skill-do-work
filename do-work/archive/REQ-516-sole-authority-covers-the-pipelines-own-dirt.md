---
id: REQ-516
title: '[impact-rule-change] Sole-authority assertion covers the pipeline's own dirt'
status: cancelled
created_at: 2026-09-02T20:35:18Z
user_request: UR-099
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-514]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-513, REQ-514, REQ-515, REQ-517]
batch: recovery-never-traps
write_set: [skills/do-work/actions/run-with-recovery.md, skills/do-work/tools/do-work-cli/internal/finalization/, skills/do-work/tools/do-work-cli/internal/requeststate/]
completed_at: 2026-09-02T22:44:37Z
---

# Sole-authority assertion covers the pipeline's own dirt

## What

Make the escape hatch wider than the door. Under `run-with-recovery`, the sole-authority assertion covers every class of dirt the pipeline itself wrote earlier in the run, proven by the journal's recorded expected hashes; the finalizer's complete transition accepts checkpoint and queue-move dirt whose bytes match those hashes. If a guard still refuses under the assertion, it names a third verb, never `run` or `rwr` again.

The fold-first scan found REQ-512 (Complete legacy finalization semantic ownership) hardening ownership proof for legacy tails; it does not cover the pipeline's own claim footprint, so this is a new REQ.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`rwr` reuses the same Step 0.1 gate as `run` with one flag widened, so it can refuse for the exact reason `run` refused. `--assume-sole-releaser` widens release metadata classes and never the lifecycle checkpoint, which is why three `rwr` invocations against REQ-456 refused identically. The journal already holds `expected_checkpoint_sha256` and `expected_request_sha256`, and both matched the dirty files byte for byte.

## Context

`internal/requeststate/state_plan.go` opts `ExistingDirtyTargetPaths` in only for the recover transition; the defer-gate plan in `internal/publication/defer_gate.go` shows the pattern of classifying expected preimages as accepted dirt. With REQ-513 in place this path is rarely hit, but it is the second line of defence for any dirt the pipeline wrote and did not commit.

## Detailed Requirements

- Under `--assume-sole-releaser`, lifecycle apply accepts a dirty target whose bytes equal the journal's recorded expected preimage for that path, and commits the postimage it produces; a mismatch still refuses.
- Without the assertion, strict behavior is unchanged.
- `actions/run-with-recovery.md` states the rule in one sentence: the assertion covers dirt the pipeline wrote earlier in the run; a remaining refusal names a third verb.
- The refusal path under the assertion is covered by REQ-514's invariant test; this REQ adds the acceptance test with the REQ-456 journal shape as fixture.
- Secret-classified and project paths stay out of scope, as UR-097 required.

## Constraints

- The acceptance is by recorded hash, never by path class alone.
- Do not persist the assertion; it remains a per-invocation verb choice.

## Batch Constraints

- Judgment stays prose; mechanics stay in the Go CLI. No new prose that walks a shell sequence.
- A guard may still refuse. What it may not do is refuse for a REQ-scoped reason in a way that stops unrelated REQs, or name itself as the fix.
- Nothing here widens recovery to secret-classified or project paths; only dirt the pipeline itself wrote earlier in the run is in scope.
- Every REQ carries a behavior test on the command or a contract predicate on the action, never a sentence pin alone.

## Dependencies

Depends on REQ-514. Related to REQ-499, REQ-512, and REQ-513.

## Builder Guidance

Certainty level: Firm on the rule, latitude on where the hash comparison lives. Read the CLI prime and `_dev/primes/prime-action-files.md`.

## Red-Green Proof

**RED prompt/case:** With a journal at `prepared` whose `expected_checkpoint_sha256` equals the dirty `do-work/CHECKPOINT.md`, run `recover-finalization --discover --assume-sole-releaser`.
**Why RED now:** Lifecycle apply refuses with `FINALIZATION-LIFECYCLE-APPLY` because the checkpoint is dirty, although the journal proves the dirt is its own.
**GREEN when:** The same invocation reaches `cleanup_complete`, and the same journal with a checkpoint whose bytes differ from the recorded hash still refuses with a finding naming a verb other than `recover-finalization`.
**Validation:** User confirmed (verify-requests, 2026-09-02).

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2643 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on semantic recovery completeness and structured evidence projection in do-work-cli internals.
- `_dev/primes/lessons-action-files.md` — 3539 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on changing action routing and status contracts.

## Full Context

See `do-work/user-requests/UR-099/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-02, item A4 of "how can I update the orchestrator to not end up in a trap like this?", captured by UR-099.*

## Cancelled

- **When:** 2026-09-02T22:44:37Z
- **Why:** superseded by 19256ba0: finalize accepts pipeline dirt by recorded preimage hash without a flag gate
- **Decided by:** user, via `do-work abandon`
