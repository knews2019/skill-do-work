---
id: REQ-514
title: '[impact-rule-change] Refusals never name themselves as the fix'
status: pending
created_at: 2026-09-02T20:35:18Z
user_request: UR-099
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-513, REQ-515, REQ-516, REQ-517]
batch: recovery-never-traps
write_set: [skills/do-work/tools/do-work-cli/internal/resultmodel/, skills/do-work/tools/do-work-cli/internal/finalization/, skills/do-work/tools/do-work-cli/internal/requeststate/]
---

# Refusals never name themselves as the fix

## What

Enforce one invariant in the result model: a refusal's `next_argv` must name a verb other than the argv that produced it, or the refusal is not allowed and the command must set the REQ aside instead. The check lives in the finding builders, not in an action file.

The fold-first scan found no pending or pending-answers REQ in any UR that owns this invariant; REQ-512 (Complete legacy finalization semantic ownership) hardens what recovery accepts, not what a refusal may say.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The REQ-456 trap had one shape: a guard refused, named itself as the fix, and a rule forbade anything else. The `FINALIZATION-LIFECYCLE-APPLY` finding carried `next_argv` and `verification_argv` both equal to `recover-finalization`, the command that had just refused. A test over the finding builders would have failed that finding the day it was written.

## Context

Findings carry `next_argv` in several result-model types (`internal/resultmodel/result_model.go`). Finalization, request-state, and publication commands all build them. The invariant is about the relation between the invoking argv and the finding's next verb, so it needs the invoking command name at build time.

## Detailed Requirements

- Every refusal finding whose `next_argv` is non-empty names a command different from the one that produced it; `verification_argv` may still be the same command, since verification is read-only.
- A guard that cannot name a different verb emits a set-aside finding with an empty `next_argv` and a stop reason, and the caller treats it as REQ-scoped, never global.
- One table-driven test walks every finding builder that refuses and asserts the invariant; the REQ-456 finding shape is the RED fixture.
- Existing self-referential findings are fixed to name their real resolving verb, or converted to set-asides.
- The invariant is a code test, not a sentence predicate in `_dev/tests/contract-regressions.sh`.

## Constraints

- No prose rule; the result model enforces it.
- Do not weaken fail-closed behavior: a guard still refuses, it only may not loop.

## Batch Constraints

- Judgment stays prose; mechanics stay in the Go CLI. No new prose that walks a shell sequence.
- A guard may still refuse. What it may not do is refuse for a REQ-scoped reason in a way that stops unrelated REQs, or name itself as the fix.
- Nothing here widens recovery to secret-classified or project paths; only dirt the pipeline itself wrote earlier in the run is in scope.
- Every REQ carries a behavior test on the command or a contract predicate on the action, never a sentence pin alone.

## Dependencies

REQ-515 (Per-REQ recovery findings never stop the loop) consumes the set-aside shape this REQ defines. Related to REQ-512.

## Builder Guidance

Certainty level: Firm on the invariant, latitude on where the invoking argv threads through. Read the CLI prime first.

## Red-Green Proof

**RED prompt/case:** Build the `FINALIZATION-LIFECYCLE-APPLY` refusal for a dirty checkpoint and compare its `next_argv` to the invoking argv.
**Why RED now:** Both are `recover-finalization`; nothing in the result model rejects a self-referential refusal.
**GREEN when:** The table-driven test fails on that fixture before the fix and passes after, and no refusal in the suite names its own command as `next_argv`.
**Validation:** Inferred during capture; the A1 through A5 list this REQ transcribes was confirmed by the maintainer in the same conversation.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2643 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on semantic recovery completeness and structured evidence projection in do-work-cli internals.

## Full Context

See `do-work/user-requests/UR-099/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-02, item A2 of "how can I update the orchestrator to not end up in a trap like this?", captured by UR-099.*
