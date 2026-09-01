---
id: REQ-446
title: 'Review fix: Match remediation to preflight failure kind'
status: claimed
domain: general
created_at: 2026-08-31T16:40:15Z
status_changed_at: 2026-08-31T19:24:17Z
user_request: UR-081
addendum_to: REQ-432
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
tdd: true
sweep: true
sweep_key: preflight-failure-kind-remediation
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
claimed_at: 2026-09-01T00:24:50Z
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-01T00:25:26Z
  basis:
    - trivial short-circuit
---

# Review Fix: Match Remediation to Preflight Failure Kind

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Kept the Route A fix to shared failure-kind projection plus the required production-caller audit and simultaneous-state regression.
- [x] **[APPLY]:** Delegated doctor preflight findings to the shared transaction renderer and added exact dirty-target/dirty-index assertions.
- [x] **[UNIFY]:** Reviewed both changed files; focused/full Go tests, vet, exact Go 1.25 compatibility, and diff hygiene pass on the builder branch.

## What
Make every caller that projects a shared transaction-preflight failure choose recovery and verification commands from the actual failure kind. Done means the class cannot recur: simultaneous-state regressions must reject guidance that inspects a target when the selected blocker is repository-wide index state.

Fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this preflight-failure-kind remediation root cause. REQ-445 covers pathless structural cleanup findings, not failure-kind projection in shared preflight callers.

## Context
Found during review of REQ-432 (Enforce the commit guard for consumed scratch cleanup). Shared preflight now correctly selects the nonempty-index blocker before target dirt in commit mode, but doctor still emits target-only next and verification commands for every failure kind.

## Requirements
- Project dirty-index failures to cached-index diagnostic and verification commands rather than target-only status or unstaged-diff commands.
- Preserve exact target-specific guidance for dirty-target and other path-scoped failures.
- Add a simultaneous dirty-index plus dirty-target doctor regression that asserts the finding code, reason, actionable argv, no mutation, and byte-identical target retention.
- Audit every shared preflight caller's failure-kind projection so the same mismatch cannot remain elsewhere.

## Instances
- [ ] `internal/doctor/doctor_repair.go`: a selected dirty-index failure still emits target-only status and unstaged-diff argv, whose verification can succeed while the index remains nonempty. (found by REQ-432 / UR-081)

## Red-Green Proof
**RED prompt/case:** Run doctor timestamp repair in commit mode with both a dirty timestamp target and an unrelated staged file; require `GIT-DIRTY-INDEX`, cached-index next/verification commands, no mutation, and byte-identical target bytes.
**Why RED now:** Doctor maps every preflight failure to target-only commands even when the shared result identifies the repository-wide index as the blocker.
**GREEN when:** The exact simultaneous-state regression passes and every shared preflight caller chooses actionable guidance from `FailureKind`.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions
- [x] Should I process this as a new task? Cleanup is now safe, but doctor can tell users to inspect or verify the wrong state when both staged changes and target dirt exist. → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to `pending` and align doctor remediation with the selected blocker under regression coverage).
  Also: No, discard it (doctor remains fail-closed, but its recovery guidance stays misleading in the simultaneous-state edge case).

  **Answered 2026-08-31** (UTC date per `actions/work-reference.md` → **Date-only stamps**):
  User confirmed the recommendation via `do-work clarify`: add the focused correction to the queue
  so recovery and verification commands match the actual preflight failure kind, including
  simultaneous dirty-index and dirty-target regression coverage. Nothing from the captured scope
  was put out of scope.

## Triage

**Route: A** — Simple

**Reasoning:** The review provides the exact failure-kind projection seam, target behavior, and simultaneous-state regression; this is a focused bug fix with a bounded caller audit.

**Planning:** Not required

## Plan

**Planning not required** — Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

Doctor timestamp repair now derives every preflight refusal's code, stop reason, fixability, next argv, and verification argv from the shared Git transaction failure-kind renderer. Dirty-target guidance remains exact-path scoped, while dirty-index guidance now inspects and verifies the repository-wide cached index.

The caller audit found only doctor and cleanup as production `PreflightTargets` consumers. Cleanup was already failure-kind-aware and already carries a simultaneous dirty-index regression, so no cleanup change was needed.

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_repair.go` (modified) — delegates preflight finding projection to the shared renderer.
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_repair_test.go` (modified) — pins path-scoped dirty-target and repository-wide dirty-index remediation.

**Builder commit:** `b5de97a9a9c0970c714526de08f5579e05f12f5c`

**Integration range:** `73bae0ae..75fedbe1`

*Generated by work action from the builder hand-back*

## Decisions

### D-01: Reuse the shared failure renderer

**Decision:** Use `gittransaction.BuildCommandResult` rather than adding a doctor-local dirty-index branch.

**Reasoning:** The shared failure template is already the authority for failure-kind-specific recovery and verification; consuming it prevents the same projection from drifting again.

## Qualification

Passed — the two-file manifest is exact, both files are substantive in `73bae0ae..75fedbe1`, the implementation traces directly to dirty-index projection and its required simultaneous-state regression, and the production caller audit found cleanup already compliant. Route A has no `## Scope` list, so scope-drift correctly reported its documented skip.

## Testing

- Mechanical qualification — passed for integration range `73bae0ae..75fedbe1`.
- `go test -count=1 ./internal/doctor ./internal/gittransaction ./internal/cleanup` — passed.
- `go vet ./...` — passed.
- Exact Go 1.25 compatibility suite — passed across the full CLI module.
- Builder full CLI tests and diff hygiene — passed per the durable hand-back.
