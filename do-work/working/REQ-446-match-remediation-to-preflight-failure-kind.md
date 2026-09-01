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
---

# Review Fix: Match Remediation to Preflight Failure Kind

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
