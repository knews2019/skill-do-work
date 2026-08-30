---
id: REQ-412
title: 'Implement request-state, checkpoint, archival, and calibration transactions'
status: pending
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
---

# Implement Request-State, Checkpoint, Archival, and Calibration Transactions

## What
Implement canonical Go transactions for request lifecycle state changes.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
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
