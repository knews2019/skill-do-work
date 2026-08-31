---
id: REQ-445
title: 'Review fix: Keep pathless refusal commands runnable'
status: claimed
domain: general
created_at: 2026-08-31T15:34:58Z
status_changed_at: 2026-08-31T19:24:17Z
user_request: UR-081
addendum_to: REQ-430
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
tdd: true
sweep: true
sweep_key: pathless-refusal-recovery-argv
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
claimed_at: 2026-08-31T22:07:41Z
route: B
write_set:
  - skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go
  - skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-31T22:08:35Z
  basis:
    - Route B
    - 3-file write set
    - 3 acceptance criteria
    - cross-route regression gates
    - full-suite verification
---

# Review Fix: Keep Pathless Refusal Commands Runnable

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Exploration localized the defect to `refusedGroupFinding`: emit repository-wide Git status only for blank structural targets, preserve exact-path diagnostics otherwise, and ratchet the complete prerequisite-refusal result matrix.
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What
Make every structural cleanup refusal that has no concrete target path emit a runnable recovery command. Done means the class cannot recur: result-level coverage must reject empty or otherwise invalid path arguments for every applicable pathless refusal.

Fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this pathless-refusal recovery-command root cause.

## Context
Found during review of REQ-430 (Couple UR closure to terminal member archival). Duplicate operation-group codes are correctly refused before mutation, but the shared path-bearing refusal helper emits an empty Git pathspec and its suggested command exits 128.

## Requirements
- Add a valid no-path recovery-command form for structural cleanup refusals, or otherwise ensure the duplicate-group refusal emits runnable next and verification commands.
- Preserve exact path-targeted Git diagnostics for refusals that do have a concrete target path.
- Add a result-level ratchet that covers every applicable structural dependency refusal and fails on empty invalid path arguments.

## Instances
- [ ] `internal/cleanup/cleanup_apply.go`: duplicate group-code refusal calls the path-bearing finding helper with an empty path, producing `git status --short -- ''`. (found by REQ-430 / UR-081)

## Red-Green Proof
**RED prompt/case:** Exercise a cleanup plan with duplicate operation-group codes and execute or validate every emitted next command; require each command to be non-empty and runnable without an empty pathspec.
**Why RED now:** The duplicate-group finding appends an empty path argument to `git status --short --`, which exits 128.
**GREEN when:** The same result-level test passes and every applicable pathless structural refusal carries valid actionable command evidence.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions
- [x] Should I process this as a new task? The cleanup safety behavior already works, but users who hit a duplicate internal group identity receive a recovery command that immediately fails. → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to `pending` and repair the command contract with a regression test).
  Also: No, discard it (the safe refusal remains, but its recovery command stays unusable for this edge case).

  **Answered 2026-08-31** (UTC date per `actions/work-reference.md` → **Date-only stamps**):
  User confirmed the recommendation via `do-work clarify`: add the focused fix to the queue
  so every applicable pathless structural cleanup refusal emits runnable recovery commands
  with result-level regression coverage. Nothing from the captured scope was put out of scope.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Triage

**Route: B** — Moderate

**Reasoning:** The fix is localized to cleanup finding construction and its tests, but it changes a shared actionable-command contract and must ratchet every applicable pathless structural refusal rather than patch one example.

**Planning:** Not required — requirements and target behavior are sufficiently concrete for focused exploration.

## Plan

Planning not required (Route B).

## Exploration

Duplicate operation-group identity is the only currently broken pathless case: `ApplyPlan` intentionally supplies no target, while `refusedGroupFinding` still appends an empty pathspec. Dependency refusals already emit runnable cleanup dry-run commands. The existing `TestOperationGroupPrerequisitesFailClosed` result seam covers direct/transitive, named-missing, duplicate, repeated, and cyclic prerequisites; extend it with the empty-required-code branch, argv argument validation, executable duplicate-group diagnostics, and an exact-path preservation assertion.

No `resultmodel` change is justified: it only carries the package-authored argv. The candidate three-file estimate is frozen, but exploration narrows implementation to the two cleanup files below.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go` (modify) — construct valid pathless versus exact-path group-refusal diagnostics.
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (modify) — result-level structural-refusal and exact-path regression matrix.

**Files I will NOT touch:** resultmodel, Git transaction primitives, other cleanup planners/commands, action prose, or release metadata.

**Acceptance criteria (restated from REQ):**
- [ ] Every duplicate-identity and prerequisite structural refusal returns nonempty, runnable recovery/verification argv with no empty argument.
- [ ] Duplicate-group identity uses the valid repository-wide Git diagnostic while named/empty missing, duplicate, repeated, transitive, and cyclic prerequisites retain runnable cleanup dry-run evidence.
- [ ] Path-bearing dirty/collision refusals preserve exact repository-relative Git diagnostics.
