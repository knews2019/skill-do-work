---
id: REQ-504
title: '[impact-rule-change] Collapse Step 10 and Crash Recovery prose into recovery'
status: claimed
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
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
preflight_at: 2026-09-04T22:36:05Z
review_at: 2026-09-04T22:28:49Z
route: C
planning_at: 2026-09-04T15:57:17Z
exploration_at: 2026-09-04T15:57:17Z
dispatch_at: 2026-09-04T16:01:12Z
builder_handback_at: 2026-09-04T22:44:47Z
integration_at: 2026-09-04T16:31:53Z
status_changed_at: 2026-09-04T21:00:44Z
commit: f412a8411057d0a833df5584657161008f315b84
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
claimed_at: 2026-09-04T22:27:42Z
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

- `_dev/primes/lessons-action-files.md` — 4362 tokens, partial coverage prevents narrowing; action/reader contract match.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, partial coverage; structural command and argv ownership match.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 6625 tokens, partial coverage; structural checkpoint evidence and alternate-writer match. All three remain additive touch-conditional builder reads.

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

## Prior Pre-Flight

**Git:** Clean outside `do-work/` after the concurrent release owner committed its previously observed board/action changes.

**Tests baseline:** `go test -count=1` across finalization, request-state, repository-model, result-model, and lifecycle-advance passed before implementation.

**Repository gate:** The exact unpiped `bash _dev/tests/maintainer-verify.sh` passed at revision `275b2fd131b0cf0906e94218a994b620ea843b63`, and typed green-gate evidence was recorded for that revision.

**Dependencies:** Installed; Go 1.26.1, ShellCheck 0.11.0, and all fast repository lanes launched successfully.

*Checked by work action.*

## Prior Implementation Summary

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

## Prior Qualification

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

## Prior Testing

- RED: before implementation, the focused lifecycle/model command failed because `recover` returned `UNKNOWN-COMMAND`, checkpoint-mode advance returned `ADVANCE-USAGE`, request-state lacked all-entry recovery types, repository discovery omitted unlabelled evidence, and finalization discovery rejected coherent claim-only topology.
- GREEN: post-merge focused tests passed for finalization, request-state, repository-model, result-model, and lifecycle-advance; request-state and lifecycle-advance race tests passed; `go vet ./...` and the uncached full CLI module passed.
- Contract replacement: `bash _dev/tests/contract-regressions.sh` passed after the redundant request-state shell lane and aggregator entry were removed.
- Canonical repository gate: direct unpiped `bash _dev/tests/maintainer-verify.sh` passed on merged `main`; 375 board tests completed in 19 seconds and 655 CLI tests in 49 seconds, with the slowest file at 22.76 seconds, below the 30-second budget.
- Green-gate evidence was recorded at revision `2069631f7b40d547a2b0218d16623e30cd717887` for the exact canonical argv.

## Heavy Verification Plan

- Base revision: `773787b74acddfdfc4c16498a89d99a5cc3ab716`
- Target revision: `f412a8411057d0a833df5584657161008f315b84`
- `queue-kanban-javascript`: `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for `_dev/tests/contract-regressions.sh` and `_dev/tests/contracts/request-state.sh`.
- `queue-kanban-browser`: `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for `_dev/tests/contract-regressions.sh` and `_dev/tests/contracts/request-state.sh`.
- `do-work-cli-integrations`: `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for `_dev/tests/contract-regressions.sh` and `_dev/tests/contracts/request-state.sh`.
- `staged-skills`: `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for `_dev/tests/contract-regressions.sh` and `_dev/tests/contracts/request-state.sh`.
- `updater`: `bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for `_dev/tests/contract-regressions.sh` and `_dev/tests/contracts/request-state.sh`.
- `installer`: `bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for `_dev/tests/contract-regressions.sh` and `_dev/tests/contracts/request-state.sh`.

## Open Questions

- [x] Heavy lanes at `f412a8411057d0a833df5584657161008f315b84`: the work loop runs them at queue exhaustion and records the result here → Confirmed: All 6 selected heavy lanes passed without skips at f412a8411057d0a833df5584657161008f315b84.


## Answer Notes

- 2026-09-04 - [ ] Heavy lanes at `f412a8411057d0a833df5584657161008f315b84`: the work loop runs them at queue exhaustion and records the result here: Confirmed: All 6 selected heavy lanes passed without skips at f412a8411057d0a833df5584657161008f315b84.
> ```
> Exact-revision heavy verification via do-work clarify. Stored base, target, selected lanes, argv and coverage reasons matched the recomputed plan. All lane results came from the detached checkout at f412a8411057d0a833df5584657161008f315b84.
> All 6 selected heavy lanes passed without skips at f412a8411057d0a833df5584657161008f315b84.
> Initial attempt: staged-skills, updater and installer each exited 1 after 0 seconds before their tests started, reporting an invalid timing-log header. Preserved the original log and initialized a fresh log using the repository test-duration-log.sh helper. Reran only those three lanes at the same revision; all passed. The earlier passing CLI integration, JavaScript and browser results remain applicable. No tracked source was changed.
> Scope: verification results only; implementation changes, fresh review and archiving remain for do-work run. Date and timestamp follow skills/do-work/actions/work-reference.md, Timestamp rule and its date-only paragraph.
> Browser lane used /Applications/Google Chrome.app/Contents/MacOS/Google Chrome through QUEUE_KANBAN_BROWSER.
> ```

## Heavy Verification Result

Target revision: `f412a8411057d0a833df5584657161008f315b84`
Execution revision: `f412a8411057d0a833df5584657161008f315b84`

- queue-kanban-javascript: exit 0, 8s — `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`
- queue-kanban-browser: exit 0, 97s — `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`
- do-work-cli-integrations: exit 0, 60s — `bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- staged-skills: exit 0, 23s — `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- updater: exit 0, 51s — `bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- installer: exit 0, 23s — `bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

## Prior Review

# Review: REQ-504

**Request changes** — the public recovery and checkpoint commands work for canonical records, but a supported legacy checkpoint can become unrecoverable and then lose its visible ownership evidence during refresh.

Route C. Reviewed `773787b74acddfdfc4c16498a89d99a5cc3ab716..f412a8411057d0a833df5584657161008f315b84`. This independent orchestrated review used the queue request as read-only evidence; the orchestrator must resolve its exact path and evidence again after claim. No request, follow-up, or lifecycle record was written.

## What's built

The migration adds finalization-first `recover`, explicit takeover authority, structural checkpoint evidence, all-entry claim removal, and checkpoint-only `advance --checkpoint`. It removes the split shell recovery fixture and large action algorithms, replacing them with command boundaries and behavioral tests. The ordinary saved-revision advance remains read-only; later requests deliberately extend it.

## Findings

**Important F-01 — supported legacy checkpoint evidence disagrees across discovery, recovery, and refresh.** `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go:347` deliberately discovers heading-less claim entries. `internal/requeststate/state_apply.go:854` (current HEAD; line 832 at the reviewed revision) requires the canonical section before removing any entry, so `recover --assume-sole-authority` refuses the same record with `RECOVER-CLAIM-CHECKPOINT-EVIDENCE`. Next, `internal/lifecycleadvance/checkpoint_commands.go:99` appends an empty canonical section below the legacy record. Discovery subsequently selects that empty section and stops returning the retained record, although refresh reported `preserved_claims: 1`. The bytes survive, but the live ownership evidence disappears from the public result. The direct public replay below confirmed both failures; the implicated functions remain unchanged at current HEAD. This fails the requirement to recover every supported checkpoint shape and preserve structurally observed live claims. — **impact-critical**; return to the orchestrator for the single remediation and fold-first decision, with no follow-up created by this reviewer.

No additional current findings. The saved revision's `work-reference.md:703` still restated automatic own-label recovery despite the new explicit-authority policy; later work removed that operative paragraph. It is historical scope context, not a remaining repair request.

## Requirements checklist

- [x] Public finalization-first recovery, coherent claim-only topology, typed ownership decisions, and constant authority argv are implemented and exercised.
- [ ] Every supported checkpoint shape recovers and remains visible after refresh: canonical labelled, unlabelled, duplicate, and hostile-label cases pass; heading-less legacy evidence fails F-01.
- [x] Authorized canonical recovery removes all same-request entries, including aliases and continuations, preserving unrelated entries and project dirt.
- [x] The saved ordinary advance is read-only; checkpoint mode changes only its checkpoint target in the passing canonical case.
- [x] Step 10 is one loop paragraph plus a context-wipe principle; Crash Recovery and session-checkpoint mechanics collapse to principles and command boundaries.
- [x] Inherited finalization-ownership, Step 8 naming, commit, handoff, and guide changes appear in the declared scope. The saved selection exception is explicitly deferred to REQ-505.
- [x] All four migration parts are present: owning commands, deleted prose, retired shell predicate/fixture registration, and replacement Go behavior tests. The original sentence-predicate RED premise was stale, and the Plan documents the replacement public-boundary RED cases.
- [x] The diff contains the declared 26 implementation paths plus its own REQ evidence record; no unrelated implementation expansion. Decisions D-01 through D-07 explain the substantive choices.

### Prior acceptance testing

**Acceptance: Fail.** Existing suites pass, but the independently exercised supported legacy flow fails.

Fresh test command at the exact target, from `skills/do-work/tools/do-work-cli`:

```text
go test -count=1 ./internal/lifecycleadvance ./internal/finalization ./internal/requeststate ./internal/repositorymodel ./internal/resultmodel
```

All five passed: lifecycleadvance 3.575s, finalization 31.723s, requeststate 4.729s, repositorymodel 1.122s, resultmodel 1.340s. These include public recover → next → fresh claim, explicit takeover versus observation, hostile labels, multiple and unlabelled entries, claim-only finalization topology, all-entry removal, and checkpoint-only mutation tests. Fresh compilation of the public CLI also passed.

The six saved heavy lanes remain valid exact-revision evidence: JavaScript 8s, browser 97s, CLI integrations 60s, staged skills 23s, updater 51s, installer 23s; every lane exit 0 with no skips. The earlier timing-header refusal and canonical-helper retry are recorded transparently in the request and preparation artifact. Heavy lanes were not rerun by this review.

### Exact independent replay

Build the target CLI with `go build -o <absolute-review-cli> ./cmd/do-work-cli`. Create a temporary Git repository, configure fixture author identity, and commit exactly these two files:

`do-work/working/REQ-713-fixture.md`:

```markdown
---
id: REQ-713
title: Legacy checkpoint replay
status: claimed
claimed_at: 2026-09-04T12:00:00Z
---

# Legacy replay
```

`do-work/CHECKPOINT.md`:

```markdown
# Session Checkpoint

- REQ-713: Legacy checkpoint replay — claimed 2026-09-04T12:00:00Z — writer: other:/checkout
  Keep detail.
```

Each following invocation used `<absolute-review-cli> --repo-root <absolute-fixture-root> --format json` followed by the shown argv:

| Argv | Observed result |
|---|---|
| `recover` | Exit 0, success; one claim with one checkpoint evidence record, writer `other:/checkout`, source line 3. |
| `recover --assume-sole-authority` | Exit 1, refused; `RECOVER-CLAIM-CHECKPOINT-EVIDENCE`; claim remains working. |
| `advance --checkpoint` | Exit 0, success; `preserved_claims: 1`; appends empty `## In Progress (interrupted)` below the old entry. |
| `recover` | Exit 0, success; the same working claim now carries `checkpoint_evidence: []`. |

The fixture ran inside the detached review checkout and was removed by its temporary-directory owner. No main queue was used as a test fixture.

## Remediation boundary and ratchet

The minimum implicated implementation paths are `internal/repositorymodel/repository_model.go`, `internal/requeststate/state_apply.go`, and `internal/lifecycleadvance/checkpoint_commands.go`, all beneath `skills/do-work/tools/do-work-cli/`. Their tests belong with those owners or the public lifecycle acceptance tests. Reuse a consistent structural interpretation of supported legacy records; avoid merely silencing the refusal or reporting the old pre-write count as preservation.

Suggested named checks for RED-before-GREEN remediation: `TestRecoverLegacyCheckpointClaimsThroughPublicCommand` and `TestAdvanceCheckpointPreservesLegacyClaimDiscovery`. The first should reproduce the refusal before the fix and prove authorized recovery removes every matching legacy entry while leaving unrelated records untouched. The second should compare public discovery before/after checkpoint refresh and prove writer, request identity, and continuation evidence remain attributable. Canonical-section, labelled/unlabelled, alias, duplicate, and no-claim controls should remain green. These names describe tests to add; this review did not install tests or fixes.

REQ-544 is an existing pending sweep about position-anchored caller-authored lifecycle evidence. It owns publication substring gates and cleanup checkpoint line selection. The checkpoint writer's substring heading check is related structural territory, but F-01 also fails on an ordinary heading-less file with no forged token. Attribution and fold belong to the orchestrator after claim; this review did not change REQ-544.

## Restatement and guardrail review

Swept core and sibling Markdown for removed `Session Checkpoint Template`, old own-label/three-hour recovery wording, and `recover-finalization --discover` consumers. Current remaining finalization references describe the still-supported primitive; changelog history is context. The current saved-to-HEAD comparison confirms that later queue/gate/finalization extensions are intentional, while the F-01 functions retain the failed logic. REQ-515's set-aside preservation is a separate later policy and is not incorrectly charged to this migration.

No speculative dependency, unrelated refactor, or naming concern was found in the new functions/files. The three new mechanics use existing typed result, repository snapshot, state plan, and transaction owners. The key test weakness is semantic composition: each canonical fixture passes while the explicitly supported legacy shape crosses inconsistent readers. Self-validation checked that this was an actual public result difference, not a complaint based only on helper code or current-version drift.

### Prior review scores

**Overall: 50%** — the acceptance-failure cap applies to the 85% arithmetic average.

| Dimension | Score |
|---|---:|
| Requirements | 80% |
| Code quality | 85% |
| Test adequacy | 75% |
| Scope | 100% |
| Risk | Critical |
| Acceptance | Fail |

**Suggested additional testing:** the two remediation checks above and their canonical controls. No browser or external-service uncertainty remains for this finding.

**Follow-ups created:** None; the orchestrator requested read-only review and owns remediation/fold routing.

**Cleanup confirmed:** all test/build/replay processes completed. The owned binary and fixture were removed; `git status --porcelain=v1 --untracked-files=all` in the detached checkout was empty. `git worktree remove .git/work-run-20260905/review-504` then completed without force. No background work remains.

## Remediation Plan

Acceptance Fail activates the work action’s one remediation attempt. The existing six-lane green record remains historical; it does not authorize completion with the reproduced legacy-checkpoint failure. Claim-time evidence was re-resolved and matched the independently reviewed exact target.

# REQ-504 Remediation Plan

Fix the disagreement between checkpoint readers and writers while retaining the already-supported legacy layout. A format migration is unnecessary: refresh can update its frontmatter and preserve the existing body, and authorized removal can consume the same claim range that discovery already recognizes.

Read-only planning at current revision `2ba5b432658853690e8e5a6d20bd2dcc147e9ada`. No source, tests, queue records, or lifecycle state changed; no tests ran in this planning pass. The exact failed public fixture and outputs are in `REQ-504-review.md` beside this file.

## Exact proposed scope

All paths below are relative to `skills/do-work/tools/do-work-cli/`:

- `internal/repositorymodel/repository_model.go`
- `internal/repositorymodel/repository_model_test.go`
- `internal/requeststate/state_apply.go`
- `internal/requeststate/state_apply_test.go`
- `internal/lifecycleadvance/checkpoint_commands.go`
- `internal/lifecycleadvance/checkpoint_commands_test.go`
- `internal/lifecycleadvance/recovery_commands_test.go`

No changes should be needed in recovery orchestration, state-plan authority rules, publication, cleanup, result schema, action prose, or the contract aggregator. If a discovered caller requires widening, record the concrete reason before doing so.

## Three coherent tasks

1. **Make the public failures RED first.** Add `TestRecoverLegacyCheckpointClaimsThroughPublicCommand` to the recovery acceptance file. Use the reviewed committed working-REQ fixture and a heading-less checkpoint, first prove observation leaves bytes intact and returns the actual writer, then require explicit-authority recovery to succeed, remove every matching legacy record and continuation, preserve unrelated records, and allow a fresh claim. Add `TestAdvanceCheckpointPreservesLegacyClaimDiscovery` to checkpoint acceptance tests: capture public recovery evidence, refresh, capture again, and compare semantic identity/writer/header evidence while allowing source-line offsets caused by frontmatter. Assert the retained entry and continuation bytes are exact and `preserved_claims` agrees with post-write discovery. Before production edits, the first test must fail on `RECOVER-CLAIM-CHECKPOINT-EVIDENCE` and the second on disappearance of checkpoint evidence. Include multiple labels, an unlabelled record, and a numeric alias as subcases of these failures rather than adding redundant smoke tests.

2. **Use one structural claim range and preserve the existing layout.** Export or narrowly expose the repository model's existing checkpoint-section helper as a descriptive two-word-or-longer identifier, such as `CheckpointClaimBounds`. Return the claim range and whether a real canonical heading exists: exact canonical heading takes precedence through the next level-two heading, otherwise the already-supported whole-document legacy range applies. Keep CRLF handling consistent with current discovery. Have repository projection, request-state removal, and the `checkpoint-absent` predicate consume that same range; this closes both the observed refusal and the false assertion of absence for a real legacy entry. Keep request identity matching and writer authorization unchanged. In `checkpointSessionBytes`, delete the branch that appends an empty canonical heading to an existing nonempty body; continue creating the normal canonical body for a genuinely empty/new checkpoint. Preserve all non-owned body bytes. Update the discovery comment that currently promises a writer-driven upgrade, since preserving a supported layout is the chosen fix.

   **Account for the fresh-claim caller in the same file.** `checkpointWithClaim` currently appends a canonical section too, which can hide another request's retained legacy evidence immediately after successful recovery. When an existing checkpoint is in legacy mode and contains observed legacy claims, append the new canonical entry line without introducing a competing section; continue using the normal section writer for canonical or empty documents. This is part of the public recovery → fresh claim acceptance control, not a separate feature. Prefer this small branch over relocating foreign records or inventing a migration. The shared range helper can distinguish structure; count/identify actual legacy claims through existing structural discovery logic rather than a loose substring check.

3. **Prove the boundaries and complete the ordinary hand-back.** Run the two named public RED→GREEN tests, then the focused owner packages. Add bounded unit controls in repository-model and request-state tests: a real canonical section ignores matching-looking entries in Completed/Notes sections; a heading token inside ordinary prose is not a heading; CRLF canonical records and LF legacy records retain their contract; absent evidence is refused when a legacy record really exists; unrelated request blocks stay byte-identical. Assert writer-specific removal retains another writer and unlabelled records, while authorized all-entry removal deletes only the requested identity. The fresh-claim control must retain the other legacy request's public evidence. Run the canonical required integration/repository verification only through the orchestrator's normal phase commands after integration; do not reuse the old target's heavy results for changed source.

## Scope and authority cautions

`RemoveAllCheckpointClaims` now also serves terminal removal in `state_plan.go`, and `RemoveOwnedCheckpointClaim` is called by cleanup. The shared range change therefore affects those callers even without editing them. Preserve their existing identity/authority predicates and cover the helpers' canonical negative controls; the accepted legacy layout should be consistently discoverable and removable, not treated as new removal authority.

REQ-544 remains a separate pending caller-authored-text anchoring sweep. Its named publication and cleanup instances are outside this seven-file remediation. Its earlier cleanup snippet has already become a call to `RemoveOwnedCheckpointClaim` at this revision, so the orchestrator should reassess that instance when REQ-544 is claimed. This plan does not fold or retire it. Exact heading recognition here is needed for the same checkpoint interpretation across the affected owners; broad substring-gate cleanup would be scope expansion.

Do not backfill the reviewed request with a passing verdict until the public reproduction is green at the integrated remediation revision. The existing five-package and six-lane green evidence remains accurate historical evidence, but it did not cover this legacy composition failure.


**Builder write boundary:** exactly the seven paths in the remediation plan, all within the original declared Scope. Preserve the original frozen estimate (50 active minutes, low confidence). Prior phase evidence is retained under Prior headings while canonical gates validate this repair.

## Remediation Decisions

- **D-08 — consistent legacy semantics:** Preserve the supported legacy layout instead of migrating foreign records. One structural bounds helper must serve discovery/removal/absence; refresh and fresh claim must retain visible evidence. This is the reviewed request’s supported-shape requirement, not a broad REQ-544 sweep.
- **D-09 — historical range attribution:** Retain original base 773787b74acddfdfc4c16498a89d99a5cc3ab716 for cumulative evidence, while recording the fresh repair merge’s parent as a supplemental attribution range. Later unrelated chain changes in the cumulative range are not this builder’s edits.

## Pre-Flight

Remediation baseline at `3501aeddd31634d0f165d193b32709c86f08cd8b`: canonical advance ran the three focused owner packages and accepted the directly executed `bash _dev/tests/maintainer-verify.sh` exit 0. Board 381 tests / 17s; CLI 723 tests / 53s; slowest CLI file 21.66s, below 30s. Both preflight and green-gate records are satisfied in `.git/work-run-20260905/REQ-504/preflight-green.json`. Original estimate remains frozen.

## Remediation Execution State

- [x] **[PLAN]:** Consume the independent failure and seven-file remediation plan; public RED first, share structural range, preserve legacy layout and authority, GREEN and boundary controls.
- [ ] **[APPLY]:** Await isolated builder.
- [ ] **[UNIFY]:** Await diff review and native checks.

Remediation builder returned at 2026-09-04T22:44:47Z; prior implementation handback was 2026-09-04T16:31:17Z. Exact seven-file branch commit: `4e351d172b14b822dd5027d3c13d12874ef5774c`. Public RED/GREEN and owner checks passed; integration and re-review remain.
