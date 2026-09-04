---
id: REQ-504
title: '[impact-rule-change] Collapse Step 10 and Crash Recovery prose into recovery'
status: claimed
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-503]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/run-with-recovery.md
  - skills/do-work/actions/commit.md
  - skills/do-work/actions/restart-with-parallel-handoff.md
  - skills/do-work/docs/work-guide.md
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - _dev/tests/contract-regressions.sh
  - _dev/tests/contracts/request-state.sh
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands.go
  - skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
  - skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
  - skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go
  - skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_types.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_plan_test.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go
  - skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go
claimed_at: 2026-09-04T15:44:47Z
route: C
planning_at: 2026-09-04T15:57:17Z
exploration_at: 2026-09-04T15:57:17Z
preflight_at: 2026-09-04T16:00:28Z
dispatch_at: 2026-09-04T16:01:12Z
builder_handback_at: 2026-09-04T16:31:17Z
integration_at: 2026-09-04T16:31:53Z
estimate:
  p50_active_minutes: 50
  confidence: low
  calculated_at: 2026-09-04T15:45:07Z
  basis:
    - Route C
    - 5-file write set
    - 3 subsystems involved
    - 11 acceptance criteria
    - dependency depth 1
    - cross-route regression gates
    - full-suite verification
---

# Collapse Step 10 and Crash Recovery Prose Into Recovery

## What
Replace `work.md` Step 10 (155 lines: loop, checkpoint, session start) and `work-reference.md` Crash Recovery and Session Checkpoint Template with one loop sentence and one principle, now that `recover-finalization`, `run-with-recovery` and `advance` own the mechanics. Delete the sentence-predicates that pinned that prose and add behavior tests on the commands.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** The Plan and Explore agents mapped the public recovery order, structural checkpoint evidence, atomic multi-entry removal, claim-only finalization topology, checkpoint-only advance mutation, prose collapse, and replacement behavior tests across the exact widened scope.
- [x] **[APPLY]:** Added the typed public recovery and checkpoint boundaries, structural checkpoint evidence, atomic all-entry recovery, claim-only finalization handling, replacement Go coverage, and the planned prose collapse in the exact 26-file scope.
- [x] **[UNIFY]:** Reviewed all 26 changed paths and the complete diff; focused tests, race tests, vet, full Go tests, contract regressions, diff checks, and the canonical maintainer gate all passed on the builder branch with no debug artifacts.

## Why
This is the largest single deletion in the chain and the one the recovery work makes possible: the prose described what to do when a session died, and the CLI now does it.

## Context
Analysis: `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70). Measured at 721c2fb4: `work.md` 850 lines and 20 steps; about 55% of step lines are mechanics; `_dev/tests/contract-regressions.sh` holds 220 references into the two work files and pins sentences with mutation-tested predicates, which is why earlier moves into Go left prose behind.

## Detailed Requirements
- `work.md` Step 10 becomes: after integration, run `advance`; loop while it selects; on exit, `advance` writes the checkpoint through the canonical command. One paragraph.
- Crash Recovery and the session-start note reduce to: run recovery; the result says what to do; `run-with-recovery` answers ownership "mine". The takeover ladder moves into the CLI result as typed findings.
- Every `contract-regressions.sh` lane that quotes the deleted sentences is deleted in the same commit; each trap they guarded is re-expressed as a Go test on `recover-finalization` or `advance`, or recorded in the lessons satellite if it was prose-only.
- Checkpoint writes move behind `advance` (this is the first mutation `advance` gains; keep it to the checkpoint).

## Constraints
- One step per REQ, never a rewrite of `work.md`; the four-part write set (CLI command, deleted prose, deleted predicates, new behavior test) is complete or the review refuses the move.
- Judgment stays prose; `advance` emits typed findings, never paragraphs.
- The floor agent must still complete a run with `advance` output plus the remaining prose.
- Serial chain; run in one session.

## Dependencies
Depends on REQ-503 (`advance`).

## Builder Guidance
Firm on the boundary between mechanics and judgment as classified in the report's step table; dispute a row in the REQ before moving it. Latitude on prose wording. Read `_dev/primes/prime-action-files.md` before touching any action file.

## Red-Green Proof
**RED prompt/case:** Delete `work.md` Step 10 and the Crash Recovery section and run `bash _dev/tests/contract-regressions.sh`.
**Why RED now:** The suite fails on the sentence-predicates that quote those sections, and no Go test covers the behaviors they described.
**GREEN when:** After the move, the suite passes with those lanes removed, `go test ./internal/finalization ./internal/lifecycleadvance` covers session-death-after-archive and foreign-claim takeover, and `work.md` Step 10 is at most one paragraph.
**Validation:** User confirmed the direction ("more principles for the LLMs, not exact steps; the Go script does mechanics"); the per-REQ RED case is inferred during capture from the report.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3436 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".

## Review Instances — REQ-498

- [ ] `skills/do-work/actions/work.md` overview still says the orchestrator owns moves, frontmatter updates, and archiving although finalization owns those mutations.
- [ ] `skills/do-work/actions/commit.md` opening restatement still describes direct request-state delegation, exact-path committing, and a separate serial metadata commit.
- [ ] Rename or rewrite Step 8's “Archive” heading/body so its judgment-only preparation role cannot be mistaken for the actual mutation boundary.

These are active operative restatements found by REQ-498's post-remediation re-review. Fold them into this request's action-prose collapse and replace any guarding sentence predicates with behavior-level finalization/advance coverage.

## Review Instances — REQ-501

- [ ] Exercise the complete public order, including initial `recover-finalization`, so a normal interrupted claim reaches ownership recovery instead of stopping at `FINALIZATION-DISCOVERY-AMBIGUOUS`.
- [ ] Recover every supported checkpoint evidence shape, including one REQ recorded under multiple writer labels, without leaving `ALREADY-CLAIMED` residue.
- [ ] Discover or transport checkpoint evidence structurally; never interpolate an observed writer label into shell source.
- [ ] Prove uninterrupted public `run-with-recovery` → `next` → fresh `claim` behavior rather than starting the fixture at the internal `recover-claim` primitive.

These critical public-boundary residuals survived REQ-501's one remediation attempt. Fold them into this request's canonical recovery-command boundary and behavior-test replacement; the isolated `recover-claim` transaction itself already covers authority/evidence/commit guards, exact-path rollback, and unrelated-dirt preservation.

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*

## Triage

**Route: C — Complex.** This rule-changing migration deletes two large recovery/checkpoint prose surfaces, transfers checkpoint behavior into the new lifecycle command, retires sentence-predicate tests, and must close seven inherited public-boundary review findings without weakening recovery ownership or shell-safety contracts.

**Planning:** Required.

## Plan

1. Add RED tests for the actual surviving gaps: claim-only state must not be misclassified as an ambiguous finalization tail; checkpoint discovery must retain labelled, unlabelled, absent, and multiple-writer evidence; one authorized recovery must remove every same-request entry atomically without touching unrelated bytes.
2. Add a canonical lifecycle `recover` composition that runs finalization recovery first, then structurally classifies/reclaims working requests, and exposes ownership/takeover decisions as typed findings. `recover --assume-sole-authority` is the constant public boundary used by `run-with-recovery`; raw writer text never enters shell source.
3. Extend `advance` only with a checkpoint mutation mode. Ordinary `advance REQ-NNN` remains byte-for-byte read-only. The checkpoint mode preserves every foreign or unlabelled live claim record and changes no project path.
4. Cover the uninterrupted public order—recover, select, fresh claim—plus hostile writer labels, multiple labels, missing/unlabelled evidence, takeover offer-versus-authority, timestamps, unrelated dirt, and worktree cleanup in Go. Retire the now-redundant split shell recovery lane and its aggregator registration after equivalent behavior is green.
5. Collapse `work.md` Step 10 and the reference recovery/checkpoint algorithms to short principles and constant command boundaries; repair the inherited finalization ownership, Step 8 naming, commit action, handoff action, CLI prime, and work-guide restatements in the same change.
6. Keep the dependency-chain boundary intact: this request temporarily replays typed `next` and invokes `advance --checkpoint` at exhaustion; the following selection/claim request remains the sole owner of moving those operations behind `advance`.

**Plan validation:** Every inherited review instance maps to a behavior test and an owning implementation path. The captured RED premise was stale because an earlier contract split already deleted the quoted sentence predicates; this plan targets the live behavioral gap instead. Retaining the exact `## Crash Recovery (Step 1)` heading avoids unnecessary changes to doctor, forensics, cleanup, and board references.

*Generated by Plan agent.*

## Exploration

The current request-state shell lane begins at the internal `recover-claim` primitive, uses one literal writer, and never proves initial finalization recovery, multi-label/unlabelled evidence, structural transport, or public recover-to-next-to-claim order. Repository discovery currently records writer-bearing checkpoint lines only, while request state accepts only one writer or one unlabelled evidence flag and removes one entry. Finalization discovery can classify a coherent claim-only footprint as an ambiguous tail before claim recovery runs. Existing request-state mutation helpers, finalization handlers, selector APIs, cleanup worktree repair, and typed result fields provide reusable composition seams.

`advance` currently accepts one request id and promises total read-only behavior; its archived-complete phase neither selects nor writes a checkpoint. The safe first mutation is therefore an explicit checkpoint mode, with the existing command remaining read-only. The old root contract file is now only an aggregator; the live replacement owner is `_dev/tests/contracts/request-state.sh`, which can be deleted only after equivalent Go coverage lands.

*Generated by Explore agent.*

## Scope

**Files I will touch:**
- `_dev/tests/contract-regressions.sh`
- `_dev/tests/contracts/request-state.sh`
- `skills/do-work/actions/commit.md`
- `skills/do-work/actions/restart-with-parallel-handoff.md`
- `skills/do-work/actions/run-with-recovery.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/work.md`
- `skills/do-work/docs/work-guide.md`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go`
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan_test.go`
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_types.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`

**Files I will NOT touch:** doctor, forensics, cleanup, board, hook routing, other command packages, or any queue file. The retained Crash Recovery heading and existing handler registration keep those consumers valid.

**Acceptance criteria:**
- [ ] Public recovery runs finalization recovery before claim recovery, accepts coherent claim-only topology, and reaches typed selection plus a fresh claim without residue.
- [ ] Structural checkpoint evidence safely covers labelled, unlabelled, absent, hostile, and multiple-writer cases; only explicit authority removes all same-request records atomically.
- [ ] Ordinary `advance REQ-NNN` stays read-only; checkpoint mode mutates only `do-work/CHECKPOINT.md` and preserves foreign/unlabelled live entries byte-for-byte.
- [ ] Step 10, Crash Recovery, Session Checkpoint, finalization ownership, and recovery handoff prose collapse without stealing selection/claim scope from the following request.
- [ ] Every retired shell behavior has equivalent Go coverage, the exact 26-path scope passes drift checks, and focused/module/repository gates pass.

## Pre-Flight

**Git:** Clean outside `do-work/` after the concurrent release owner committed its previously observed board/action changes.

**Tests baseline:** `go test -count=1` across finalization, request-state, repository-model, result-model, and lifecycle-advance passed before implementation.

**Repository gate:** The exact unpiped `bash _dev/tests/maintainer-verify.sh` passed at revision `275b2fd131b0cf0906e94218a994b620ea843b63`, and typed green-gate evidence was recorded for that revision.

**Dependencies:** Installed; Go 1.26.1, ShellCheck 0.11.0, and all fast repository lanes launched successfully.

*Checked by work action.*

## Implementation Summary

- `_dev/tests/contract-regressions.sh` (modified) — removes the registration for the retired request-state shell lane after equivalent public Go coverage landed.
- `_dev/tests/contracts/request-state.sh` (deleted) — retires the shell fixture that began at the internal recovery primitive.
- `skills/do-work/actions/commit.md` (modified) — assigns the resumable mutation boundary to finalization instead of restating manual request-state commits.
- `skills/do-work/actions/restart-with-parallel-handoff.md` (modified) — replaces the removed checkpoint-template reference with the command-owned principle.
- `skills/do-work/actions/run-with-recovery.md` (modified) — invokes constant recover-with-sole-authority argv and hands the original scope to the work loop.
- `skills/do-work/actions/work-reference.md` (modified) — collapses crash recovery and session checkpoint algorithms to concise command-owned principles.
- `skills/do-work/actions/work.md` (modified) — collapses Step 10, clarifies orchestration judgment, and renames Step 8 to preparation.
- `skills/do-work/docs/work-guide.md` (modified) — aligns user guidance with typed recovery and checkpoint commands.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (modified) — declines coherent unfinished claim-only topology so public recovery owns it.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go` (modified) — covers the claim-only discovery exemption.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (modified) — adds explicit checkpoint mode while preserving ordinary advance reads.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go` (modified) — keeps ordinary advance byte-for-byte read-only after checkpoint mode lands.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands.go` (new) — implements the guarded checkpoint-only mutation transaction.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/checkpoint_commands_test.go` (new) — proves exact-path checkpoint mutation and live-entry preservation.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go` (new) — composes finalization-first recovery with typed authority and claim recovery.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go` (new) — covers public recovery order, hostile labels, multiple evidence shapes, takeover, selection, and fresh claim.
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (modified) — projects labelled and unlabelled checkpoint evidence structurally with legacy fallback.
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (modified) — covers canonical-section and heading-less checkpoint evidence.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified) — removes every authorized same-request checkpoint entry atomically.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modified) — proves all-entry removal and unrelated-byte preservation.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` (modified) — plans structural all-entry recovery authority.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan_test.go` (modified) — covers multi-entry planning and evidence guards.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_types.go` (modified) — adds the all-entry checkpoint recovery option.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified) — adds typed recovery and checkpoint result projections.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified) — locks text/JSON recovery and checkpoint output shape.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — documents the new command ownership and the narrowed advance mutation boundary.

**What was done:** Recovery now runs finalization discovery first and reports claim authority as typed data; explicit authority resets each request and removes all matching checkpoint evidence without interpolating writer text into shell source. Checkpoint-mode advance is the sole checkpoint mutation while ordinary request-specific advance stays read-only. The action and reference prose now state the remaining judgment principles and preserve REQ-505's ownership of selection and claim.

## Qualification

Passed — 26 files verified, 5 requirements traced, P-A-U confirmed. The four new Go package/test files are convention-discovered package members and tests, so the qualifier's static-reference warnings are expected rather than dead-code findings.

## Decisions

- **D-01 — DECIDE & STATE:** Public recovery runs finalization discovery first and stops on its typed refusal, preserving the existing resumable-finalization authority.
- **D-02 — DECIDE & STATE:** Plain recovery observes unfinished claims; only explicit takeover or sole-authority mode authorizes their reset.
- **D-03 — DECIDE & STATE:** Writer labels remain typed data and never become generated shell source.
- **D-04 — DECIDE & STATE:** Authorized recovery removes all normalized same-request checkpoint records in one request-state plan, including multiple labels, unlabelled records, aliases, and continuations.
- **D-05 — DECIDE & STATE:** Canonical-section checkpoint parsing wins when present; heading-less whole-document parsing remains as legacy read compatibility.
- **D-06 — DECIDE & STATE:** `advance --checkpoint` performs an exact-path uncommitted session-end write; ordinary advance remains read-only.
- **D-07 — DECIDE & STATE:** Recovery commits one request-state transaction per working request and stops before later requests if one transaction fails.

## Discovered Tasks

None.
