---
id: REQ-445
title: 'Review fix: Keep pathless refusal commands runnable'
status: completed
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
completed_at: 2026-08-31T22:38:01Z
commit: 6d14a046
---

# Review Fix: Keep Pathless Refusal Commands Runnable

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Exploration localized the defect to `refusedGroupFinding`: emit repository-wide Git status only for blank structural targets, preserve exact-path diagnostics otherwise, and ratchet the complete prerequisite-refusal result matrix.
- [x] **[APPLY]:** Added the complete result-level RED matrix first, reproduced both empty duplicate pathspecs, then changed only the shared refusal constructor.
- [x] **[UNIFY]:** Reviewed both scoped files; merged-state focused/full Go, vet, exact Go 1.25, contracts, qualification, scope, diff hygiene, and the canonical maintainer gate pass.

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
- [x] `internal/cleanup/cleanup_apply.go`: duplicate group-code refusal calls the path-bearing finding helper with an empty path, producing `git status --short -- ''`. (found by REQ-430 / UR-081)

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
- [x] Every duplicate-identity and prerequisite structural refusal returns nonempty, runnable recovery/verification argv with no empty argument.
- [x] Duplicate-group identity uses the valid repository-wide Git diagnostic while named/empty missing, duplicate, repeated, transitive, and cyclic prerequisites retain runnable cleanup dry-run evidence.
- [x] Path-bearing dirty/collision refusals preserve exact repository-relative Git diagnostics.

## Implementation Summary

`refusedGroupFinding` now emits repository-wide `git status --short` when a structural group refusal has no concrete target, and retains `git status --short -- <exact path>` for path-bearing dirt and collision cases. The result-level prerequisite fixture covers duplicate identities plus direct, transitive, named/empty missing, duplicate, repeated, and cyclic prerequisite failures; it rejects empty argv elements, executes the duplicate diagnostic, and pins exact path-bearing evidence.

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go` (modified) — conditional pathless versus exact-path Git diagnostic construction.
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go` (modified) — complete structural-refusal result and runnable-command ratchet.

**Integration range:** `a6724bee..6d14a046`

*Generated by work action from the builder hand-back*

## Decisions

### D-01: Keep both refusal shapes in the existing helper

**Decision:** DECIDE & STATE — construct repository-wide or exact-path Git status argv conditionally in `refusedGroupFinding` rather than adding a helper or result-model validator.

**Reasoning:** The cleanup package authors this command contract; moving validation into the passive shared result carrier would broaden the fix without improving ownership.

### D-02: Execute only the emitted Git diagnostic in-package

**Decision:** DECIDE & STATE — execute the duplicate-identity Git argv in the fixture repository and assert prerequisite commands as the exact registered cleanup form.

**Reasoning:** Package tests should not assume a globally installed `do-work-cli`, while the Git diagnostic is directly runnable in the fixture.

## Testing

**Red-green validation:** Before the production change, both duplicate-identity findings returned `git status --short -- ''`; the result-level fixture failed on the empty element and exact diagnostic. The same fixture passes after the conditional constructor change.

- Focused cleanup result and dirty-group tests — PASS.
- Full cleanup package and full do-work-cli module tests — PASS.
- `go vet ./...` — PASS.
- Exact Go 1.25 compatibility — PASS.
- Builder scope and diff hygiene — PASS; exactly the two frozen files changed.

## Qualification

Passed — the exact two-file integration range `a6724bee..6d14a046` matches the frozen Scope, both changed files are substantive, the complete structural-refusal matrix traces directly to all acceptance criteria, and no tracked `do-work/` file appears in the implementation range.

**Merged-state checks:** focused and full cleanup tests, full do-work-cli tests, `go vet ./...`, exact Go 1.25 compatibility, contract regressions, mechanical qualification, scope drift, diff hygiene, and `bash _dev/tests/maintainer-verify.sh` all pass. The optional browser lane skipped because no browser was available.

## Review

**Overall: 99%** | 2026-08-31T22:32:38Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 98% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None.
**Minor findings:** None.
**Acceptance:** Pass — the inherited REQ-430 empty-pathspec defect reproduced against the pre-fix constructor and the merged result-level matrix passes with runnable pathless diagnostics and preserved exact-path evidence.
**Suggested testing:** None.
**Follow-ups created:** None; **sweeps appended to:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** A result-level matrix over every structural prerequisite refusal closed the class rather than only the duplicate-group example, while executing the emitted Git diagnostic proved it was actually runnable.
**What didn't:** Reusing a path-bearing recovery helper for a structural finding without a path produced an argv value that looked populated but failed at execution.
**Worth knowing:** Shared recovery-command constructors must make the presence of a concrete target explicit; blank structural targets need repository-wide diagnostics, while real targets retain exact pathspec evidence.

## Orientation

Cleanup structural refusals now emit runnable recovery commands whether or not they carry a filesystem target. The behavior and its regression matrix remain localized to `internal/cleanup`; no result-model or action-contract change was required.
