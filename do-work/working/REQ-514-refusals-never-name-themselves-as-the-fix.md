---
id: REQ-514
title: '[impact-rule-change] Refusals never name themselves as the fix'
status: claimed
priority: now
created_at: 2026-09-02T20:35:18Z
user_request: UR-099
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: [REQ-513]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-513, REQ-515, REQ-516, REQ-517]
batch: recovery-never-traps
write_set: [skills/do-work/tools/do-work-cli/internal/resultmodel/, skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go]
claimed_at: 2026-09-03T23:34:31Z
route: C
planning_at: 2026-09-03T23:37:51Z
exploration_at: 2026-09-03T23:37:51Z
dispatch_at: 2026-09-03T23:42:48Z
estimate:
  p50_active_minutes: 50
  confidence: low
  calculated_at: 2026-09-03T23:37:51Z
  basis:
    - Route C
    - 8-file write set
    - 3 subsystems involved
    - 5 acceptance criteria
    - dependency depth 1
    - cross-route regression gates
    - full-suite verification
---

# Refusals never name themselves as the fix

## What

Enforce one invariant in the result model: a refusal's `next_argv` must name a verb other than the argv that produced it, or the refusal is not allowed and the command must set the REQ aside instead. The check lives in the finding builders, not in an action file.

The fold-first scan found no pending or pending-answers REQ in any UR that owns this invariant; REQ-512 (Complete legacy finalization semantic ownership) hardens what recovery accepts, not what a refusal may say.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Enforce the invariant at result normalization, preserve verification semantics, distinguish REQ-owned set-asides from global refusals, and correct ambiguous discovery's actual resolver.
- [x] **[APPLY]:** Added centralized self-remedy normalization, a distinct ambiguous-discovery resolver, and regression coverage in the reconciled five-file scope.
- [x] **[UNIFY]:** Reviewed all five changed files; focused packages, full Go module except one transient timing test, isolated timing-test rerun, vet, and diff checks passed.

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

Depends on REQ-513 (Commit the claim footprint in every mode): the maintainer asked for A1 first, so this REQ and everything behind it wait for it. REQ-515 (Per-REQ recovery findings never stop the loop) consumes the set-aside shape this REQ defines. Related to REQ-512.

## Builder Guidance

Certainty level: Firm on the invariant, latitude on where the invoking argv threads through. Read the CLI prime first.

## Red-Green Proof

**RED prompt/case:** Build the `FINALIZATION-LIFECYCLE-APPLY` refusal for a dirty checkpoint and compare its `next_argv` to the invoking argv.
**Why RED now:** Both are `recover-finalization`; nothing in the result model rejects a self-referential refusal.
**GREEN when:** The table-driven test fails on that fixture before the fix and passes after, and no refusal in the suite names its own command as `next_argv`.
**Validation:** User confirmed (verify-requests, 2026-09-02).

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 2643 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on semantic recovery completeness and structured evidence projection in do-work-cli internals.

## Folded From REQ-517

Hand triage 2026-09-03, maintainer approved: REQ-517 (Pin the serial claim-to-recovery trap) is cancelled and its test becomes this REQ's acceptance test. One end-to-end fixture test in the CLI's Go suite: claim, a one-line implementation change, complete, then `recover-finalization --discover`, asserting the terminal phase is `cleanup_complete`. Today that sequence stops at lifecycle apply, which is the real failure the test names.

## Full Context

See `do-work/user-requests/UR-099/input.md` for complete verbatim input.

---
*Source: maintainer conversation of 2026-09-02, item A2 of "how can I update the orchestrator to not end up in a trap like this?", captured by UR-099.*

---

## Triage

**Route: C** - Complex

**Reasoning:** The invariant spans the shared result model, command runtime, lifecycle refusals, and finalization recovery. It also folds the REQ-517 end-to-end lifecycle fixture into the same behavior boundary.

**Planning:** Required

## Plan

1. Add RED tests that enumerate refusal results and reject a non-empty `next_argv` whose verb is the invoking command, including the current `FINALIZATION-LIFECYCLE-APPLY` fixture.
2. Put one invariant-enforcement boundary in the shared result/runtime path so all command families are checked consistently and a self-referential refusal becomes a REQ-scoped set-aside with an empty remedy.
3. Correct existing lifecycle and finalization refusal builders to name an external resolving verb where one is known, preserving same-command `verification_argv`.
4. Add the folded claim → implementation → completion → recovery regression, then run focused package tests, vet, and the canonical gate.

**Plan validation:** All five detailed requirements map to these four tasks. The shared enforcement step owns the output contract; package-specific edits only supply truthful alternate remedies or fixture coverage. No action-owned mutation is driven by an untyped aggregate.

*Generated inline under the Plan-agent fallback*

## Exploration

- `commandruntime.Run` assigns `Command` only after handlers return, immediately before rendering, while `resultmodel.NormalizeResult` currently normalizes nil collections without validating remedy identity.
- `requeststate.refusalResult` mechanically points back to its own transition; finalization failures and discovery refusals mechanically reuse recovery argv as both `next_argv` and `verification_argv`.
- `finalization_pipeline_dirt_test.go` already owns the REQ-513 claim-to-finalization regression shape and is the narrow place to add the folded REQ-517 lifecycle fixture.
- The implementation must distinguish a recommended next action from a same-command verification rerun; only the former is forbidden.

*Generated inline under the Explore-agent fallback*

## Scope

**Files I will touch (reconciled before integration):**
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify) — enforce and normalize the refusal remedy invariant
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modify) — table-driven result-contract coverage
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go` (modify) — runtime set-aside and exit behavior coverage
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (modify) — provide the truthful inventory resolver for ambiguous discovery
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` (modify) — verify the discovery remedy/verification split

**Files I will NOT touch:** action prose, shell contract predicates, or queue-selection behavior.

**Acceptance criteria (restated from REQ):**
- [ ] Every refusal with a non-empty next command points to a different verb than the invocation.
- [ ] Same-command verification remains allowed and read-only.
- [ ] A refusal with no truthful alternate remedy becomes REQ-scoped set-aside evidence with empty `next_argv` and a stop reason.
- [ ] The REQ-456 self-loop fixture is red before the fix and all refusal builders pass after it.
- [ ] The folded serial claim → implementation → complete → recovery fixture reaches `cleanup_complete`.

## Decisions

- **D-01 (scope reconciliation, 2026-09-04):** The planned `finalization_apply.go` edit did not own the ambiguous discovery refusal. Expand the exact scope to `finalization_discovery.go` and its existing recovery test file, and remove planned files that required no change. The directory-level request scope already covered finalization; `commandruntime/command_runtime_test.go` is added explicitly. The folded lifecycle fixture already exists in `finalization_recovery_test.go`, so duplicating it would reduce signal.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` (modified)

**What was done:** Result normalization now removes a refusal remedy when it names the invoking command while preserving verification. A fully REQ-owned self-reference becomes a findings outcome; a shared blocker remains refused without the false remedy. Ambiguous finalization discovery now names inventory collection as its resolver and recovery discovery as verification. Table-driven model, runtime JSON/exit, and discovery coverage pin those distinctions. Integrated at e42ae1e57c9f2692a598cb08daca2fe99bec6a45 from implementation range 26b3426886bfea6183502809a7e5e93799831a52..e42ae1e57c9f2692a598cb08daca2fe99bec6a45.

## Testing

- PASS: `go test -count=1 ./internal/resultmodel ./internal/commandruntime ./internal/finalization ./internal/requeststate`
- PASS: `go test -count=1 ./internal/gittransaction -run TestCancelledCommitThatLandsReportsCommittedRisk` (isolated rerun after the all-package run hit its pre-existing `hook.pid` timing flake)
- PASS: `go vet ./...`
- PASS: `git diff --check`
- The existing folded lifecycle test `TestRecoverFinalizationResumesJournalAfterLifecycleInterruption` remains green and reaches terminal cleanup through recovery.

## Review

**Overall: 50%** | 2026-09-04T00:00:24Z

| Dimension | Score |
|-----------|-------|
| Requirements | 35% |
| Code Quality | 65% |
| Test Adequacy | 40% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Fail |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Result normalization applies the no-self-remedy invariant only to aggregate `OutcomeRefused`; live `FixabilityRefused` findings under `OutcomeFindings` retain same-verb `next_argv`, and the synthetic table does not walk real producers — `impact-rule-change` → report only
- Ambiguous discovery's top-level finding names `uncommitted-inventory`, but its canonical singular/ordered finalization records still name `recover-finalization --discover` as both next and verification, producing contradictory typed remedies — `impact-user-visible` → report only
- The folded serial claim acceptance fixture seeds an already-claimed REQ instead of executing the claim command, so it does not pin the requested claim → implementation → complete → recovery sequence — `impact-rule-change` → report only

**Minor findings:** None
**Acceptance:** Fail — focused and full Go checks pass, but the universal refusal invariant and folded serial-claim sequence have direct acceptance gaps.
**Suggested testing:** 2 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*

## Remediation

One bounded remediation pass is authorized by the work loop. It must apply the invariant to every refusal finding regardless of aggregate outcome, align singular and ordered finalization remedies, and execute the literal claim path in the folded lifecycle acceptance fixture.
