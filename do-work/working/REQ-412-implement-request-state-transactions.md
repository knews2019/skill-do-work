---
id: REQ-412
title: 'Implement request-state, checkpoint, archival, and calibration transactions'
status: claimed
claimed_at: 2026-08-31T19:48:00Z
route: C
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-411, REQ-433]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
write_set:
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_types.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_types_test.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_targets.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_targets_test.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_plan_test.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_commands.go
  - skills/do-work/tools/do-work-cli/internal/requeststate/state_commands_test.go
  - skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go
  - skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go
  - skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/abandon.md
  - skills/do-work/actions/clarify.md
  - skills/do-work/actions/commit.md
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - _dev/tests/contract-regressions.sh
estimate:
  p50_active_minutes: 95
  confidence: low
  calculated_at: 2026-08-31T19:48:00Z
  basis:
    - Route C
    - 18-file write set
    - 10 new files
    - 7 subsystems involved
    - 4 acceptance criteria
    - dependency depth 1
    - persistence changes
    - cross-route regression gates
    - full-suite verification
---

# Implement Request-State, Checkpoint, Archival, and Calibration Transactions

## What
Implement canonical Go transactions for request lifecycle state changes.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Accepted a three-task plan, explored the existing model/Git/action seams, preserved explicit-target provenance semantics, and froze the justified 20-file scope before coding.
- [x] **[APPLY]:** Added RED-first registration and existing-untracked transaction fixtures, then implemented the five lifecycle commands, shared rollback opt-in, and action delegation inside the frozen 20-file scope.
- [x] **[UNIFY]:** Reviewed the exact 20-file diff; passed gofmt, focused/full CLI tests, vet, exact Go 1.25, diff/scope hygiene, and debug-artifact checks. Main-tree contracts and the canonical gate are rerun after integration.

## Detailed Requirements
- Implement `do-work-claim`, `do-work-unblock`, `do-work-complete`, `do-work-fail`, and `do-work-cancel`.
- Synchronize checkpoints, timestamps, request movement, archival, UR closure, commit hashes, and calibration logging as atomic transactions.
- Support dry-run where meaningful, optional exact-path commit, and all shared refusal/rollback guarantees.
- Emit actionable text/JSON for successful changes, skipped work, dependency failure, and rollback state.

## Constraints
- Existing natural-language lifecycle aliases must ultimately delegate every deterministic phase to these commands.

## Dependencies
Depends on REQ-411 (canonical selection and dependency semantics).

## Builder Guidance
Certainty level: Firm. Model each lifecycle transition as an explicit transaction with fixture-proven preconditions and touched paths.

## Red-Green Proof
**RED prompt/case:** Exercise every lifecycle transition through clean, dirty-target, dependency-failed, pre-commit-failure, commit, and post-commit-verification fixtures.
**Why RED now:** Lifecycle mutations are performed through LLM action steps and multiple helpers rather than one transaction engine.
**GREEN when:** Each command makes only its exact permitted state transition, synchronizes all coupled records, and reports or rolls back failures under the shared contract.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

---

## Triage

**Route: C** - Complex

**Reasoning:** This request creates a shared lifecycle transaction engine spanning request moves, timestamps, checkpoints, UR closure, archival, calibration, optional Git commits, rollback, and text/JSON results. The persistence and cross-route regression surface requires explicit planning and exploration.

**Planning:** Required

## Plan

1. Create the pure planning half of `internal/requeststate`: typed lifecycle transitions and refusal codes, exact target resolution, and one deterministic plan of every request, checkpoint, archive/UR, calibration, and provenance path. Start with command-level and planner RED fixtures for `claim`, `unblock`, `complete`, `fail`, and `cancel`, including stale, ambiguous, dependency, collision, dry-run, and exact-path cases.
2. Implement the apply/command half around one `gittransaction` execution per pre-commit lifecycle phase. Revalidate expected bytes, preserve unknown request content, update coupled state atomically, register all five commands, and prove clean success, rollback, committed-risk, exact-index commit, renderer parity, and post-commit verification behavior. Use one injected clock/writer identity and derive calibration from the archived bytes.
3. Delegate deterministic lifecycle mechanics from `work.md`, `work-reference.md`, `abandon.md`, `clarify.md`, and `commit.md`; keep confirmation, failure classification, review/release judgment, follow-up creation, and dependent disposition in prose. Update the CLI prime and contract regressions to pin the single-authority boundary and forbid free-form fallback mutations.

**Architecture decisions:** Plan every coupled path and expected preimage before applying once; reuse the existing repository, dependency, request, Git-transaction, and result authorities; validate selector-provided exact paths instead of rescanning; derive UR closure from the projected post-transaction snapshot; and treat a successful Git commit as the boundary after which recovery is an explicit revert rather than automatic history rewriting.

**Testing approach:** Use RED-first real-command and package fixtures for all transitions, then focused request-state tests, full CLI tests and vet, exact Go 1.25 compatibility, action/commit-hash contract regressions, scope/diff hygiene, and the canonical repository gate.

**Plan validation:** All Detailed Requirements and the alias-delegation constraint map to the three tasks; every task traces to a requirement. The three-task shape stays within the Route C quality ceiling. Exploration must verify or shrink the candidate 18-file boundary before scope freezes.

*Generated by Plan agent*

## Exploration

The proposed `internal/requeststate` split is viable against existing repository, dependency, request-document, and result interfaces. One unavoidable shared-layer change expands the candidate boundary from 18 to 20 files: `gittransaction` currently refuses every existing untracked target, while the documented consumer topology may leave `do-work/` untracked. A narrow caller-owned opt-in must snapshot and restore exact untracked targets without weakening the default guard, index checks, or cleanup's separate exception.

Claim must preserve REQ-411 provenance: an explicitly named `REQ-NNN` bypasses dependency gating and clears assignment, while default and UR-expanded selections remain gated. The other command boundaries remain action-judgment inputs followed by deterministic mechanics: exact successful-probe/user-confirmed unblock evidence, action-classified failure fields, chosen terminal-success status and known implementation hash, and confirmed cancellation/dependent disposition.

Existing patterns cover the rest: one repository snapshot and canonical graph; lossless `RequestDocument` edits; complete target planning plus `MutationRecorder`; projected UR closure using the canonical terminal predicate; exact writer-labelled checkpoint removal; and calibration re-read from archived bytes. No repositorymodel, dependencygraph, requestmodel, or resultmodel change is required.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_types.go` (new) — lifecycle types, transition evidence, plans, and refusal codes.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_types_test.go` (new) — transition and option contracts.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_targets.go` (new) — exact contained target/provenance resolution.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_targets_test.go` (new) — ambiguity, stale, source-tree, and explicit-provenance cases.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` (new) — pure coupled-path and projected-state plans.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan_test.go` (new) — five-transition plan matrix and invariants.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (new) — exact revalidation, lossless writes/moves, checkpoint/UR/calibration/provenance apply.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (new) — filesystem, rollback, and commit fixtures.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_commands.go` (new) — public command parsing and typed results.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_commands_test.go` (new) — real-runtime text/JSON and exit behavior.
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modify) — explicit rollback-safe existing-untracked target opt-in.
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modify) — opt-in isolation, restoration, and unchanged default guards.
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modify) — register five lifecycle handlers.
- `skills/do-work/actions/work.md` (modify) — delegate claim, unblock, completion, and failure mechanics.
- `skills/do-work/actions/work-reference.md` (modify) — define canonical lifecycle command inputs and retained action ownership.
- `skills/do-work/actions/abandon.md` (modify) — delegate confirmed cancellation mechanics.
- `skills/do-work/actions/clarify.md` (modify) — delegate user-confirmed unblock mechanics.
- `skills/do-work/actions/commit.md` (modify) — deconflict lifecycle commit/provenance ownership.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modify) — index the new lifecycle authority and judgment boundary.
- `_dev/tests/contract-regressions.sh` (modify) — pin registration, active delegation, no-fallback, and retained-judgment contracts.

**Files I will NOT touch:** repositorymodel, dependencygraph, requestmodel, resultmodel, cleanup, queue-kanban, unrelated commands/tests, queue state from the builder, or release metadata.

**Acceptance criteria (restated from REQ):**
- [ ] `claim`, `unblock`, `complete`, `fail`, and `cancel` implement only their exact permitted transitions while preserving explicit-target provenance semantics.
- [ ] Each transition synchronizes its complete coupled checkpoint, timestamp, request move/archive, projected UR closure, commit provenance, and calibration state atomically or reports rollback/risk.
- [ ] Dry-run is byte-identical; optional exact-path commit preserves the two-commit provenance rule; dirty, stale, dependency, collision, rollback, and post-commit failures are actionable and deterministic.
- [ ] Text and JSON expose matching changes, skipped work, dependency/refusal evidence, exact next/verification commands, and rollback state.
- [ ] Natural-language aliases delegate every deterministic lifecycle phase while retaining confirmation, classification, review/release judgment, follow-up, and dependent decisions.

## Pre-Flight

**Git:** ⚠ Ten unrelated concurrent main-tree edits are preserved and excluded (the five previously observed implementation/test files plus both changelogs, both version files, and `skills/do-work/actions/version.md`); builder worktrees start from clean committed owner bookkeeping.
**Tests baseline:** ✓ The canonical gate passed at the immediately preceding implementation state; subsequent owner commits contain only queue/planning bookkeeping.
**Dependencies:** ✓ Go module dependencies and exact Go 1.25 compatibility tooling are available.

*Checked by work action*

## Implementation Summary

**What was done:** Added canonical `claim`, `unblock`, `complete`, `fail`, and `cancel` transactions. `internal/requeststate` validates exact targets and selection provenance, plans every coupled path before mutation, supports byte-identical dry-run, synchronizes request/checkpoint/archive/UR/calibration/provenance state through the shared Git transaction, and returns the existing typed text/JSON contract. The Git layer now has a narrow explicit opt-in for exact existing-untracked targets with byte/mode rollback; action prose delegates mechanics while retaining judgment.

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified) — registrations, active delegation, retained judgment, and no-fallback contracts.
- `skills/do-work/actions/abandon.md` (modified) — confirmed cancellation delegates to the canonical cancel command.
- `skills/do-work/actions/clarify.md` (modified) — confirmed unblock delegates to the canonical unblock command.
- `skills/do-work/actions/commit.md` (modified) — lifecycle provenance ownership deconflicted.
- `skills/do-work/actions/work-reference.md` (modified) — lifecycle command inputs and action boundary.
- `skills/do-work/actions/work.md` (modified) — claim, probe-unblock, complete, and fail delegation.
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified) — registers five handlers.
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified) — explicit existing-untracked target snapshot/rollback and exact committable paths.
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modified) — opt-in refusal, rollback, and commit-path coverage.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (new) — exact revalidation and lifecycle apply.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (new) — filesystem, rollback, and two-commit fixtures.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_commands.go` (new) — public parsing and typed results.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_commands_test.go` (new) — real command, renderer, and exit behavior.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` (new) — deterministic coupled-path/projected-state plans.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan_test.go` (new) — five-transition planning matrix.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_targets.go` (new) — exact contained path, ID, status, and provenance validation.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_targets_test.go` (new) — stale, ambiguous, source-tree, and dependency/provenance cases.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_types.go` (new) — transition types, evidence, options, and refusal codes.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_types_test.go` (new) — type and transition contracts.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — indexes requeststate as lifecycle authority.

**Integration range:** `01d8ba12..7fc958be`

*Generated by work action from the builder hand-back*

## Decisions

### D-01: Preserve selection provenance through claim

**Decision:** DECIDE & STATE — explicit-REQ provenance bypasses dependency gating and clears assignment; default and UR-expanded provenance remain gated.

**Reasoning:** Delegation must not narrow the established targeting override. Value: canonical mechanics preserve user intent. Risk: callers must pass and validate provenance with the exact selected path.

### D-02: Opt in exact existing-untracked transaction targets

**Decision:** DECIDE & STATE — the Git transaction accepts an explicit set of existing-untracked targets, snapshots their bytes/modes, and keeps the default refusal for every other untracked path.

**Reasoning:** Consumer repositories may leave `do-work/` untracked, so lifecycle commands otherwise cannot operate. Value: ordinary topology works with rollback guarantees. Risk: broadening the opt-in would weaken dirty-target safety; focused tests pin exactness and default refusal.

### D-03: Separate rollback paths from committable paths

**Decision:** DECIDE & STATE — a moved-away never-tracked source remains in rollback/report evidence but is omitted from Git's exact add pathspec.

**Reasoning:** Git rejects an absent path that never existed in the index. Value: untracked-source moves can commit the exact destination while still restoring the source on pre-commit failure. Risk: path classification must stay synchronized with move planning.

### D-04: Keep judgment in actions and mechanics in Go

**Decision:** DECIDE & STATE — commands consume explicit confirmation, failure classification, terminal status/hash, and dependent disposition; they never infer those choices.

**Reasoning:** The REQ consolidates deterministic lifecycle writes, not review or user judgment. Value: one mutation authority without turning prose decisions into hidden policy. Risk: incomplete caller evidence must fail closed rather than trigger a free-form fallback.

## Integration Correction

The first post-merge contract run showed that rewriting Step 8's terminal-stamp paragraph had removed `work.md`'s only resolvable sibling reference to the queue-kanban completion-time consumer. The scoped correction commit `eaef32e7` restored that semantic citation without changing ownership: requeststate still writes `completed_at` and the implementation hash; queue-kanban only consumes them. The correction merged as `a21c10ae`, and both staged-skills and full contract regressions pass on the corrected branch.

## Qualification

Passed — the 20 declared implementation files are substantive in `01d8ba12..7fc958be`; all lifecycle commands are registered and flow through typed plan/apply/results; the shared Git opt-in remains explicit; action directives delegate without fallback; requirements and P-A-U trace to code/tests. Static-reference warnings for same-package Go source/test files are expected compiled package membership. The cumulative range also contains concurrent UR-085 queue-only capture state and owner review bookkeeping, which are outside builder ownership and excluded from implementation scope judgment.

## Testing

**Red-green validation:** All five command-registration cases initially failed with `UNKNOWN-COMMAND`; the existing-untracked transaction fixture initially failed to compile without its explicit option; and the serial move fixture exposed Git's refusal of an absent never-tracked source pathspec. The final fixtures pass with exact transition, rollback, and two-commit behavior.

**Merged-state checks:**
- Focused `go test -count=1 ./internal/requeststate ./internal/gittransaction` — PASS.
- Full do-work-cli `go test -count=1 ./...` and `go vet ./...` — PASS.
- Exact Go 1.25 compatibility — PASS.
- Qualification, scope drift, and diff hygiene over `01d8ba12..7fc958be` — PASS.
- Staged-skills and full contract regressions — PASS after the one-file integration correction.
- A direct standalone invocation of `_dev/tests/record-commit-hash-guards.sh` was not a valid lane because its path is rewritten by the owning contract harness; the owning full contract ran those probes successfully.

## Initial Review

**Overall: 50%** | 2026-09-01T00:16:51+03:00

| Dimension | Score |
|-----------|-------|
| Requirements | 65% |
| Code Quality | 82% |
| Test Adequacy | 78% |
| Scope | 100% |
| Risk | Critical |
| Acceptance | Fail |

**Important findings (durable impact record):**
- `work.md` retains executable archive/UR/calibration writes after declaring canonical `complete` the sole owner, allowing double moves and duplicate calibration — impact-critical.
- `abandon.md` retains manual archive/collision mechanics after `cancel`, while the command rejects the established contained multiline cancellation-reason form — impact-user-visible.
- `state_plan.go` accepts arbitrary nonblank `error_type` values instead of the canonical intent/spec/code/environment set — impact-rule-change.
- Failed-to-cancelled body history omits the required prior `error_type` classification detail — impact-user-visible.

**Minor findings:** Missing direct ratchets for duplicate action mutations, invalid failure classifications, contained multiline reasons, and exact failed-cancellation history — impact-rule-change.
**Acceptance:** Fail — the engine is substantial and tested, but active aliases do not yet delegate every deterministic phase to one authority and two durable lifecycle semantics are incomplete.
**Follow-ups created:** None pending the one allowed remediation attempt; **sweeps appended to:** None.

*Reviewed by review-work action*

## Remediation

The single allowed remediation pass resolved all four Important findings and added direct regression coverage. `work.md` and `abandon.md` now delegate each deterministic terminal mutation to one canonical transaction; failure classification accepts only `intent`, `spec`, `code`, or `environment`; multiline cancellation reasons are represented by a safe summary plus byte-identical contained body text; and failed-to-cancelled history preserves the exact failure instant and optional prior `error_type`.

Remediation commit `8b1dc746` merged as `7fc958be`. Focused and full Go tests, vet, exact Go 1.25 compatibility, staged-skills/full contract regressions, diff hygiene, and frozen-scope checks all passed before re-review.

*Generated by work action from the remediation hand-back*

## Review

**Overall: 83%** | 2026-08-31T21:54:18Z

| Dimension | Score |
|-----------|-------|
| Requirements | 92% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Existing-untracked rollback narrows the saved mode through `Mode().Perm()` and can strip setuid, setgid, or sticky bits — impact-rule-change → appended to and promoted REQ-447.
- Commit-phase instructions still condition calibration staging on removed Step 8 substep 7.5, which can omit command-owned calibration from the lifecycle release commit — impact-user-visible → REQ-459 created.

**Minor findings:** 1 — cancellation examples still say later cleanup closes an eligible UR; routed to `do-work/prose-backlog.md`.
**Acceptance:** Partial — all original review findings are closed with direct RED/GREEN evidence; the residual mode and staging-contract gaps are isolated in queued follow-up work.
**Suggested testing:** Add complete-mode existing-untracked rollback coverage and a contract ratchet for command-owned calibration staging.
**Follow-ups created:** REQ-459; **sweeps appended to:** REQ-447.

*Re-reviewed by review-work action after the one allowed remediation attempt*

## Lessons Learned

**What worked:** Replaying each initial finding against pre-remediation production code established exact RED evidence, and separating action judgment from typed request-state mutation made the corrected ownership boundary easy to ratchet.
**What didn't:** The first pass removed duplicate lifecycle writers without sweeping later commit-phase restatements, and the shared rollback opt-in copied ordinary permission bits instead of the complete restorable mode.
**Worth knowing:** A transaction is only atomic at repository level when its action, commit staging instructions, rollback metadata, and consumer prose all agree on the same target set and ownership boundary.

## Orientation

[MAP CHANGED] Request lifecycle mutations now live in `do-work-cli`'s `internal/requeststate` transaction layer. The natural-language claim, unblock, complete, fail, and cancel actions supply judgment and consume typed results; they no longer reimplement archive, checkpoint, UR-closure, calibration, or provenance writes.
