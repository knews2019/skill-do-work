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
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
