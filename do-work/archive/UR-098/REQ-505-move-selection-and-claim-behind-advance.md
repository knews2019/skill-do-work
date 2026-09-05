---
id: REQ-505
title: '[impact-rule-change] Move selection and claim behind advance'
status: completed
kb_status: pending
review_at: 2026-09-04T23:01:28Z
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
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
route: C
planning_at: 2026-09-04T16:49:45Z
exploration_at: 2026-09-04T16:49:45Z
preflight_at: 2026-09-04T16:52:22Z
dispatch_at: 2026-09-04T16:52:56Z
builder_handback_at: 2026-09-04T17:24:31Z
integration_at: 2026-09-04T17:26:13Z
status_changed_at: 2026-09-04T21:02:21Z
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
heavy_verified_at: 2026-09-04T21:02:21Z
heavy_verified_revision: 716187b847d1de0402b69587a2fe5cf7e7bd8516
claimed_at: 2026-09-04T23:00:06Z
completed_at: 2026-09-04T23:01:29Z
release_at: 2026-09-04T23:01:29Z
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
- `_dev/primes/lessons-action-files.md` — 4362 tokens, over the shared budget; `slugged: partial` so no targeted form. Matched on action routing, status contracts, and downstream readers.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the shared budget; `slugged: partial` so no targeted form. Matched on prescribed argv and migration-parity behavior.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 6929 tokens, over the shared budget; `slugged: partial` so no targeted form. Matched on frozen-ledger projection, structured evidence, collision identity, and commit preflight ordering.

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

- [x] Heavy lanes at `716187b847d1de0402b69587a2fe5cf7e7bd8516`: the work loop runs them at queue exhaustion and records the result here → Confirmed: All 4 selected heavy lanes passed without skips at 716187b847d1de0402b69587a2fe5cf7e7bd8516.

## Prior Orientation

[MAP CHANGED] Queue selection and claim now live behind the lifecycle-advance subsystem: `advance` returns committed claims and durable stateless continuation evidence, while the work action retains only the human-readable claimability principle and orchestration judgment.


## Answer Notes

- 2026-09-04 - [ ] Heavy lanes at `716187b847d1de0402b69587a2fe5cf7e7bd8516`: the work loop runs them at queue exhaustion and records the result here: Confirmed: All 4 selected heavy lanes passed without skips at 716187b847d1de0402b69587a2fe5cf7e7bd8516.
> ```
> Exact-revision heavy verification via do-work clarify. Stored base, target, selected lanes, argv and coverage reasons matched the recomputed plan. All lane results came from the detached checkout at 716187b847d1de0402b69587a2fe5cf7e7bd8516.
> All 4 selected heavy lanes passed without skips at 716187b847d1de0402b69587a2fe5cf7e7bd8516.
> Initial attempt: staged-skills, updater and installer each exited 1 after 0 seconds before their tests started, reporting an invalid timing-log header. Preserved the original log and initialized a fresh log using the repository test-duration-log.sh helper. Reran only those three lanes at the same revision; all passed. The earlier passing CLI integration result remains applicable. No tracked source was changed.
> Scope: verification results only; implementation changes, fresh review and archiving remain for do-work run. Date and timestamp follow skills/do-work/actions/work-reference.md, Timestamp rule and its date-only paragraph.
> ```

## Heavy Verification Result

Target revision: `716187b847d1de0402b69587a2fe5cf7e7bd8516`
Execution revision: `716187b847d1de0402b69587a2fe5cf7e7bd8516`

- do-work-cli-integrations: exit 0, 61s — `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- staged-skills: exit 0, 25s — `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- updater: exit 0, 52s — `bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- installer: exit 0, 23s — `bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

## Review

Canonical claim375735da preserved exact saved commit/heavy revision716187b847d1de0402b69587a2fe5cf7e7bd8516 and exposed resume_phase:review. Stored baseeb01a94f2dad78bf30f334e0614393d571ae362e resolves exactly. The prematurely recorded Orientation was retained as Prior Orientation before the phase advanced. Independent saved-revision review below remains applicable.

**Disposition:** Partial71.25% is completed under Step7; two noncritical findings remain report only. No remediation, new request, or additional heavy execution is required for this review-only resume.


**Approve with follow-ups — ordinary claims and bounded continuation work, but two queue edge cases miss their required disposition.** Acceptance is Partial, so the orchestrator should apply its remediation gate before completing this request.

This is read-only review preparation before claim. Reviewed queued `do-work/queue/REQ-505-move-selection-and-claim-behind-advance.md`, observed `status: pending`, `route: C`, `commit` and `heavy_verified_revision` both `716187b847d1de0402b69587a2fe5cf7e7bd8516`, and `heavy_verified_at: 2026-09-04T21:02:21Z`. The orchestrator must re-resolve identity and resume evidence when claiming. No queue, request, source, or lifecycle state was changed by this review.

### What's built

Queue-mode `advance` delegates readiness to the canonical selector, commits claims separately, holds detected cycles and explicit archive collisions, and returns frozen membership with tokenized continuation. The action's selection/claim procedure is reduced to one command-owned principle. Later evidence-gate and finalization behavior belongs to REQ-506/507 and was not judged as a regression against this saved implementation.

### Review

**Overall: 71.25%**

| Dimension | Score |
|---|---|
| Requirements | 60% |
| Code Quality | 85% |
| Test Adequacy | 80% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

The four percentage dimensions average 81.25%; Partial acceptance applies the specified ten-point penalty. Three of the five declared acceptance criteria are fully delivered; frozen-member blockers and collision handling are partial.

**Important findings:**

- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/queue_commands.go:286` — Default no-argument selection drops an archive-colliding pending request before the hold loop: `freezeQueueMembers` admits only `DEPENDENCY-CYCLE` exclusions on this path, while the canonical selector classifies the queue/archive duplicate as `DEPENDENCY-AMBIGUOUS`. A nested archived twin therefore leaves the queue copy `pending`, with no `archive-collision-hold` phase or commit. Explicit-target mode works because it includes the exclusion, so the existing explicit collision test misses the default-run regression. Users see an unresolved pending request instead of the required durable collision hold. — impact-user-visible → report only
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/queue_commands.go:139` — A frozen unconsumed member that disappears from a UR expansion is absent from both per-request maps, so the claim loop breaks without emitting a selection blocker. Replaying after that member is removed returns exit 0, `outcome: success`, empty phases/findings/claims, and the same nonempty continuation with that member still unconsumed. The run can repeat forever or appear drained without the request-specific blocker the transferred REQ-453 contract requires. The same branch is still present in the main source at review time. — impact-user-visible → report only

**Minor findings:** None.
**Acceptance:** Partial — focused public command tests pass, but independent CLI fixtures reproduce both findings at the exact saved revision.
**Suggested testing:** 2 items, detailed below.
**Follow-ups created:** None (2 findings report only).

### Requirements checklist

- [x] Default, explicit, UR, wave and fan-out queue modes select through the canonical readiness authority and commit successful claims independently. Existing public fixtures verify exact committed paths, clean trees, assignment/dependency provenance and partial claim refusal.
- [ ] Frozen membership, provenance, original flags and dispatch bounds survive continuation; later UR members are ignored and projection precedes bounding. These work, but a disappearing unconsumed UR member lacks its required genuine blocker (finding 2).
- [ ] Nested archive collisions, dependency cycles and successful blocked probes reach guarded state transactions while ordinary work continues. Explicit collisions, cycles and successful/failed probes work; default archive collisions never reach the hold transaction (finding 1).
- [x] Working/archive single-target advance remains read-only at this saved revision, and checkpoint mode retains its separate boundary. Later gate/finalization extension is outside this review range.
- [x] Step 1 and Steps 2.0/2 were collapsed to one selection-and-claim principle, and live guide/reference/clarify/prime readers were aligned. Public Go behavior tests replace the executable ownership of the removed procedure.

### Scope and traceability

Reviewed exact range `eb01a94f2dad78bf30f334e0614393d571ae362e..716187b847d1de0402b69587a2fe5cf7e7bd8516`. Its 18 substantive implementation paths match the Implementation Summary and declared Scope; one additional changed working REQ is owner bookkeeping. Five declared test seams were intentionally unused, as Qualification records. P-A-U is fully checked, and Decisions D-01 through D-04 explain the checkpoint-test update, stateless continuation flags, unbounded observation seam, and cycle consumption.

UR-098's four-part migration constraint was checked explicitly. This range has the CLI owner, deleted prose, and new behavior tests, but no predicate-file deletion. The REQ records that its original sentence-predicate RED was stale after REQ-504. Independently confirmed that the preceding REQ-504 range removed 88 lines from `_dev/tests/contracts/request-state.sh` plus its aggregator invocation, and there is no remaining named targeted-ledger/Step-2 selection predicate to delete at this base. Treat this as an already-retired predicate surface, not a missing live owner or a reason to manufacture a test-file edit. It does not excuse the missing public behavior cases above.

The restatement sweep covered the deleted targeting ledger, selection/claim steps, read-only advance statements, and heavy-review resume references across active actions, guide and contract surfaces at the target. The REQ-453 missing-member requirement is visible in the deleted contract and is the behavioral source of finding 2. No additional actionable stale restatement was established. No approach directive was assigned. Naming, standard-library dependency direction and surgical-change guardrails were checked; neither finding involves data loss or unsafe mutation.

### Acceptance evidence

Created a detached worktree under `.git/work-run-20260905/review-505` at exactly `716187b847d1de0402b69587a2fe5cf7e7bd8516` and built its public CLI. No heavy or full repository gate was rerun.

Fresh command: `go test -count=1 ./internal/lifecycleadvance ./internal/nextselection ./internal/requeststate ./internal/resultmodel` from the saved checkout's CLI module. All passed: lifecycleadvance 8.422s, nextselection 3.521s, requeststate 5.666s, resultmodel 1.252s. The lifecycle package runs the public executable matrix for default and explicit claims, UR chains and forks, later-member exclusion, successful/failed probes, nested explicit collisions, cycles, dirty partial refusal, wave/fan-out, hostile tokens, and working/archive/checkpoint boundaries.

Saved exact-revision heavy evidence was independently read from the request: CLI integrations exit 0 / 61s; staged skills exit 0 / 25s; updater exit 0 / 52s; installer exit 0 / 23s. All four are recorded without skips and at the same execution and target revision. The earlier log-header failures were recorded as infrastructure-only attempts followed by successful same-revision reruns; no result sharing across revisions was assumed.

### Reproduce finding 1

1. In a throwaway Git repository, commit `do-work/queue/REQ-901-fixture.md` with canonical `id: REQ-901`, `title: fixture`, `status: pending`; also commit `do-work/archive/UR-900/REQ-901-fixture.md` with the same id and `status: completed`.
2. Run the exact saved binary with `--repo-root <fixture> --format json advance`, with no target tokens.
3. Observed exit 0 and success; `excluded[0].code` is `DEPENDENCY-AMBIGUOUS`; `queue_advance.frozen_members`, `phases` and `claimed` are empty; queue file remains `pending`. Expected a guarded committed `archive-collision-hold`, exact duplicate-path evidence, and `blocked-archive-collision` while retaining the archive file.

### Reproduce finding 2

1. In a separate throwaway Git repository, commit pending `REQ-901` and `REQ-902` queue files, both with `user_request: UR-900`.
2. Invoke the saved binary with `--repo-root <fixture> --format json advance UR-900`. It commits REQ-901's claim and returns a ledger with REQ-901 consumed and REQ-902 unconsumed.
3. Remove only the unconsumed REQ-902 queue fixture and commit this external change. Replay the returned continuation argv exactly, adding the fixture repository root.
4. Observed exit 0 and success; selected/findings/phases/claims are empty; the selector's only exclusion identifies `UR-900`, not REQ-902; the returned continuation still contains unconsumed REQ-902 and is unchanged. Expected a typed non-success blocker with exact frozen REQ-902 identity/path and verification guidance. Both fixture removal and all commands were confined to the reviewer-owned throwaway repository.

### Suggested additional testing

1. Add a public no-target collision case alongside the existing explicit collision case, including an unrelated ready request to prove the hold and continued claim both occur.
2. Add a public UR continuation case whose unconsumed member disappears; require a typed exact-member blocker, non-success exit, and no claim or mutation. Preserve already-consumed members and the frozen set.

Self-validation checked failure paths after the passing matrix instead of treating existing tests as sufficient acceptance. Both independent fixtures failed their intended expectation before any changes; no repair was attempted. Reviewer-owned checkout, binary and fixture repositories were removed after recording this report. No reviewer subprocess or agent work remains pending.

## Lessons Learned

**What worked:** Public CLI fixtures and exact-revision heavy evidence establish ordinary selection/claim behavior without borrowing later changes.

**What did not:** Testing explicit collisions missed the default no-target path, and a frozen-ledger happy path missed a disappearing unconsumed member. Both remain visible in this Review as report-only limitations.

**Worth knowing:** A missing frozen member needs attributable blocker evidence; absent selector-map entries must not be mistaken for successful exhaustion. Review the public input modes as well as the helper’s happy path.

## Orientation

[MAP CHANGED] Queue selection and committed claims run through advance, with typed continuation for targeted runs. Two reviewed edge cases remain: default archive collisions can remain pending, and a disappearing frozen member can return an unchanged continuation without its blocker. These noncritical findings are report-only.
