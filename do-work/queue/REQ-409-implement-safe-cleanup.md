---
id: REQ-409
title: 'Implement safe cleanup passes and explicit destructive repairs'
status: pending
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-408]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
---

# Implement Safe Cleanup Passes and Explicit Destructive Repairs

## What
Make cleanup a canonical no-LLM command that applies only provably safe repairs by default.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Implement cleanup Passes 0–4, including stranded terminal REQ archival, documentation-link repointing, and merged-worktree cleanup.
- Support text/JSON, dry-run, optional commit, and shared rollback/target guards.
- Apply only provably safe repairs by default and report conflicts with evidence and exact next actions.
- Require explicit destructive flags for blanked-record restoration and deletion of unmerged worktrees.

## Constraints
- Default `do-work-cleanup` must be safe and mechanical.
- The screenshot’s `STRANDED-FINISHED-REQ` case must be repairable directly without an LLM.

## Dependencies
Depends on REQ-408 (shared repository model).

## Builder Guidance
Certainty level: Firm. Preserve the existing cleanup pass meanings while moving deterministic decisions and mutations into Go.

## Red-Green Proof
**RED prompt/case:** Put a completed REQ in `do-work/queue/` and run `do-work-cleanup --dry-run`, then run the applying form in a clean Git fixture.
**Why RED now:** Verification reports the stranded REQ but no direct no-LLM cleanup command performs the repair.
**GREEN when:** Dry-run reports the exact archive move, apply performs it safely, JSON gives actionable evidence, and destructive cases remain refused until their explicit flags are supplied.
**Validation:** User confirmed via the supplied implementation plan and original stranded-REQ example.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
