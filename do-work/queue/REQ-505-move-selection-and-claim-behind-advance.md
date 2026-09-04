---
id: REQ-505
title: '[impact-rule-change] Move selection and claim behind advance'
status: pending-heavy-testing
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-504]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/clarify.md
  - skills/do-work/docs/work-guide.md
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/queue_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/queue_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go
  - skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_types.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_plan_test.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
claimed_at: 2026-09-04T16:39:30Z
route: C
planning_at: 2026-09-04T16:49:45Z
exploration_at: 2026-09-04T16:49:45Z
preflight_at: 2026-09-04T16:52:22Z
dispatch_at: 2026-09-04T16:52:56Z
builder_handback_at: 2026-09-04T17:24:31Z
integration_at: 2026-09-04T17:26:13Z
status_changed_at: 2026-09-04T17:31:18Z
commit: 716187b847d1de0402b69587a2fe5cf7e7bd8516
estimate:
  p50_active_minutes: 50
  confidence: low
  calculated_at: 2026-09-04T16:40:03Z
  basis:
    - Route C
    - 8-file write set
    - 2 new files
    - 3 subsystems involved
    - 3 acceptance criteria
    - dependency depth 1
    - cross-route regression gates
    - full-suite verification
---

# Move Selection and Claim Behind advance

## What
Steps 1 (95 lines), 2.0 and 2 become `advance` phases; `work.md` keeps one selection principle (what makes a REQ claimable, keyed on the condition) and nothing procedural.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** The Plan and Explore agents mapped a typed queue-mode advance composition, explicit frozen-ledger continuation state, nested archive-collision and cycle holds, successful-probe unblock-to-claim behavior, per-request committed claim transactions, prose collapse, and replacement public tests across the widened scope.
- [x] **[APPLY]:** Implemented queue-mode selection and claim, stateless targeted continuation, guarded hold/unblock transactions, typed results, and the planned prose collapse within the declared scope.
- [x] **[UNIFY]:** Reviewed the 18-file diff, ran gofmt, focused and race tests, module tests, vet, contract regressions, and the direct maintainer gate; no debug artifacts or out-of-scope paths remained.

## Why
Selection already runs through `next` and claim through `claim`; the 114 lines of prose restate their contracts.

## Context
Analysis: `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70). Measured at 721c2fb4: `work.md` 850 lines and 20 steps; about 55% of step lines are mechanics; `_dev/tests/contract-regressions.sh` holds 220 references into the two work files and pins sentences with mutation-tested predicates, which is why earlier moves into Go left prose behind.

## Detailed Requirements
- `advance` on an unclaimed queue runs `next` with the session's targeting tokens and fan-out bound, then `archive-collision` and `claim`, reporting each as a phase result.
- The targeted-run ledger contract (REQ-453) moves from prose into `advance` state; its lanes become Go tests.
- Delete the Step 1, 2.0 and 2 prose and their predicates in the same commit.

## Constraints
- One step per REQ, never a rewrite of `work.md`; the four-part write set (CLI command, deleted prose, deleted predicates, new behavior test) is complete or the review refuses the move.
- Judgment stays prose; `advance` emits typed findings, never paragraphs.
- The floor agent must still complete a run with `advance` output plus the remaining prose.
- Serial chain; run in one session.

## Dependencies
Depends on REQ-504.

## Builder Guidance
Firm on the boundary between mechanics and judgment as classified in the report's step table; dispute a row in the REQ before moving it. Latitude on prose wording. Read `_dev/primes/prime-action-files.md` before touching any action file.

## Red-Green Proof
**RED prompt/case:** Delete Steps 1, 2.0 and 2 from `work.md` and run the contract suite.
**Why RED now:** Predicates naming Step 1 (16) and Step 2 (10) fail; the ledger contract has no Go test.
**GREEN when:** Suite passes without those lanes; `advance` tests cover default, targeted, fan-out and collision selection; `work.md` selection prose is one paragraph.
**Validation:** User confirmed the direction ("more principles for the LLMs, not exact steps; the Go script does mechanics"); the per-REQ RED case is inferred during capture from the report.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 4163 tokens, over the shared budget; `slugged: partial` so no targeted form. Matched on action routing, status contracts, and downstream readers.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the shared budget; `slugged: partial` so no targeted form. Matched on prescribed argv and migration-parity behavior.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 6265 tokens, over the shared budget; `slugged: partial` so no targeted form. Matched on frozen-ledger projection, structured evidence, collision identity, and commit preflight ordering.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*

## Triage

**Route: C — Complex.** This rule-changing migration moves default, targeted, fan-out, collision, and claim sequencing into the lifecycle command while deleting a large active action surface and replacing its prose predicates with public command behavior tests.

**Planning:** Required.

## Plan

1. Add public RED fixtures for no-argument/default, explicit-REQ, UR, fan-out, frozen-ledger replay, nested archive collision, dependency cycle, successful blocked-probe unblock, partial claim refusal, and working/archive read-only behavior.
2. Extend `advance` with a queue mode that invokes the canonical selector in-process, returns its complete typed selection evidence, and carries an explicit stateless targeted ledger with original tokens, provenance, consumed membership, non-fan-out flags, and effective dispatch bound.
3. Execute queue phases through typed transactions: classify archive collisions from the recursive repository snapshot, hold collision/cycle records with exact preimage validation, unblock only from selector-supplied successful probe evidence, rediscover after each mutation, and claim each selected REQ in its own canonical commit.
4. Add typed per-phase and claimed-request projections with tokenized continuation/verification argv, then collapse Step 1 and Steps 2.0/2 plus the targeted-ledger restatement to one command-owned selection/claim principle across work, reference, clarify, guide, and CLI-prime readers.
5. Run focused/race/module tests, behavior contracts, scope/qualification checks, and the direct canonical repository gate; preserve REQ-506's ownership of evidence-gate execution and REQ-507's ownership of finalization.

**Plan validation:** All three detailed requirements map to the tasks above. The five tasks are broader than the preferred three-task planning shape, but splitting this serial one-step migration would recreate an unsafe half-state where deleted claim prose has no complete command owner. The captured sentence-predicate RED and four-entry scope are stale after REQ-504; the live REDs are public behavioral gaps and the exact implementation scope is widened below.

## Exploration

Current `advance` accepts checkpoint mode or one REQ id; a pending queue REQ only returns a separate claim argv, so no default/targeted selection session, mutation, ledger, or phase sequence exists. `nextselection.Select` owns canonical readiness while its argument parser and command composition are private. Request state already exposes plan/apply seams for claim and unblock but has no exact collision/cycle hold transition. Result projection is singular rather than queue-session shaped.

Repository discovery already walks nested archive trees and carries complete filename/frontmatter collision evidence, so the new transaction can correctly detect `archive/UR-NNN/` twins without another parser. The legacy archive-collision helper is shallow and misses that normal shape; the new advance path will retire it from work orchestration and pin nested behavior at the public lifecycle seam. Claim plans must be rebuilt from a fresh snapshot after every prior claim because each commit changes the checkpoint preimage. Targeted replay must observe the unbounded canonical result, project frozen membership, and only then apply the saved bound so later UR members cannot consume a slot.

*Generated by Plan and Explore agents.*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/clarify.md`
- `skills/do-work/docs/work-guide.md`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/queue_commands.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/queue_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_types.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan_test.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`

**Files I will NOT touch:** command registration, repository-model production code, the shallow compatibility archive-collision helper, finalization/evidence-gate packages, board code, or the now-aggregator-only root contract file.

**Acceptance criteria:**
- [ ] Fresh default, explicit, UR, wave, and fan-out queue advances emit ordered typed phases and commit only successfully claimed records.
- [ ] Targeted continuation freezes membership/provenance, projects before bounding, ignores later members, and reports genuine blockers without shell interpolation.
- [ ] Nested archive collisions, dependency cycles, and successful blocked probes reach exact guarded hold/unblock/claim transactions while unrelated work continues.
- [ ] Working/archive advance stays read-only, checkpoint mode remains checkpoint-only, and later evidence/finalization phases remain outside this request.
- [ ] Step 1 and Steps 2.0/2 collapse to a command-owned selection/claim principle with all live readers aligned and public behavior tests replacing stale prose predicates.

## Pre-Flight

**Git:** Clean outside the owner-written run brief under `do-work/`.

**Tests baseline:** Uncached lifecycle-advance, next-selection, request-state, and result-model package tests passed before implementation.

**Repository gate:** Direct unpiped `bash _dev/tests/maintainer-verify.sh` passed at revision `f7afaccb8822c35749962bdc47a914d9a700140f`; 375 board tests and 655 CLI tests passed, and typed green-gate evidence was recorded for the exact argv.

**Dependencies:** Installed; Go 1.26.1, ShellCheck 0.11.0, and all repository lanes launched successfully.

*Checked by work action.*

## Decisions

- **D-01 — DECIDE & STATE:** Expanded the scope by one existing checkpoint test file because REQ-504 deliberately pinned pending-queue advance as read-only, while REQ-505 intentionally makes that exact phase mutating. The replacement fixture must pin working/archive read-only behavior and preserve checkpoint-only mutation; leaving the old assertion would make the stated behavior and the inherited test mutually exclusive.
- **D-02 — DECIDE & STATE:** Encoded targeted continuation as explicit tokenized CLI state using a dispatch bound plus repeated frozen members. This survives process boundaries and preserves frozen membership without shell interpolation; strict validation and hostile-token tests contain the added public evidence flags.
- **D-03 — DECIDE & STATE:** Kept ordinary default selection bounded while targeted, wave, fan-out, simple, and continuation observations use the unbounded selector seam before projection. Both paths delegate to the same selector, preserving the one-request floor without starving frozen targeted members.
- **D-04 — DECIDE & STATE:** Hold every dependency-cycle member before dispatch and count successful holds as consumed session work. This prevents cycle livelock while guarded per-record commits preserve durable and attributable state changes.

## Implementation Summary

**Files changed:**

- `skills/do-work/actions/clarify.md` (modified) — aligned clarification guidance with advance-owned queue mechanics.
- `skills/do-work/actions/work-reference.md` (modified) — replaced procedural selection and ledger instructions with typed advance evidence contracts.
- `skills/do-work/actions/work.md` (modified) — collapsed selection and claim mechanics to one command-owned principle.
- `skills/do-work/docs/work-guide.md` (modified) — updated the user-facing queue execution description.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (modified) — routed pending queue inputs through the composed advance transaction.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go` (modified) — updated public advance command coverage.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands_test.go` (modified) — preserved working/archive read-only and checkpoint-only behavior.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/queue_commands.go` (new) — implemented selection, holds, unblock, claim, and continuation orchestration.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/queue_commands_test.go` (new) — covered default, targeted, wave, fan-out, collision, cycle, probe, refusal, and hostile-token paths.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_commands.go` (modified) — exposed canonical argument parsing for lifecycle composition.
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modified) — added unbounded observation while retaining ordinary bounded selection.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified) — applied guarded hold and unblock transitions.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` (modified) — planned exact queue-state mutations from fresh snapshots.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan_test.go` (modified) — verified hold and unblock state plans.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_types.go` (modified) — represented the new mutation evidence and state transitions.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified) — projected queue phases, claims, continuation, and verification evidence.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified) — verified structured and rendered queue-advance results.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — documented the new lifecycle-advance authority boundary.

**What was done:** Queue selection and claim now execute as one public lifecycle-advance composition. It returns committed claim members, ordered typed phase findings, and tokenized continuation evidence while preserving default one-request dispatch, frozen targeted membership, guarded collision/cycle holds, successful-probe unblock, truthful partial refusal, and checkpoint-only behavior.

## Discovered Tasks

None.

## Qualification

**Passed.** The exact merge range contains 18 substantive implementation files and no undeclared path. All five acceptance criteria trace to public queue-advance behavior and aligned action readers, and the checked P-A-U state matches the merged diff. The scope ceiling intentionally reserved five selector/request-state test seams that the implementation did not need; the scope checker reported only those declared-but-untouched paths. Its two static-reference warnings are expected for convention-discovered Go package files.

## Testing

**Red-green validation:** The captured prose-predicate RED had become stale after REQ-504 removed the root sentence-predicate ownership, so the builder used the equivalent live public seam. The focused queue-mode test first failed because no-argument advance returned usage and an explicit pending target only projected a claim command; it then passed after the composed queue transaction was implemented. The final public queue matrix passed in 4.55s.

**Merged-state verification:** Focused lifecycle-advance, selector, request-state, and result-model packages passed in 16.38s; the same packages under the race detector passed in 20.73s; static analysis passed in 0.66s; the full Go module passed in 46.12s; and contract regressions passed in 16.34s. The direct, unpiped canonical maintainer gate passed with 375 board tests and 664 CLI tests. Green-gate evidence was recorded at revision `70c780c33e5411a6da8f36c5902e0eb5f19be7da`.

## Heavy Verification Plan

**Base revision:** `eb01a94f2dad78bf30f334e0614393d571ae362e`

**Target revision:** `716187b847d1de0402b69587a2fe5cf7e7bd8516`

- **do-work-cli-integrations** — argv: `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`; selected because the implementation changes the CLI lifecycle, selection, request-state, result-model, tests, and prime subtree.
- **staged-skills** — argv: `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`; selected because shipped action, guide, prime, and CLI files changed under the skills subtree.
- **updater** — argv: `bash _dev/tests/maintainer-verify.sh --heavy-lane updater`; selected because shipped CLI subtree changes must survive updater integration.
- **installer** — argv: `bash _dev/tests/maintainer-verify.sh --heavy-lane installer`; selected because shipped CLI subtree changes must survive installer integration.

## Open Questions

- [ ] Heavy lanes at `716187b847d1de0402b69587a2fe5cf7e7bd8516`: the work loop runs them at queue exhaustion and records the result here

## Orientation

[MAP CHANGED] Queue selection and claim now live behind the lifecycle-advance subsystem: `advance` returns committed claims and durable stateless continuation evidence, while the work action retains only the human-readable claimability principle and orchestration judgment.
