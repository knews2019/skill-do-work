---
id: REQ-495
title: 'Review fix: Emit retryable repository-gate collision evidence'
status: cancelled
domain: general
created_at: 2026-09-02T03:32:03Z
user_request: UR-095
addendum_to: REQ-492
review_generated: true
impact: impact-user-visible
effort_estimate: effort-substantive
tdd: true
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-action-files.md]
depends_on: [REQ-492]
related: [REQ-491, REQ-493]
completed_at: 2026-09-03T11:47:24Z
---

# Review Fix: Emit Retryable Repository-Gate Collision Evidence

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Align `defer-gate` collision results with the work action’s bounded max+1 retry contract, preserving whether a collision refused before mutation or occurred during a transaction that fully rolled back.

## Context

Found during post-remediation review of REQ-492 (Integrate repository-gate deferral and resumption into `do-work run`). No pending REQ or sweep shares this result-projection root cause; REQ-493 concerns target topology rather than retry evidence.

## Requirements

- Emit a typed pre-mutation collision result that unambiguously carries zero changes and the canonical no-rollback-needed state consumed by the action.
- Preserve collision identity through a post-mutation transaction failure when and only when retrying with a fresh ID is safe.
- Require successful complete rollback before classifying a post-mutation collision as retryable.
- Keep incomplete rollback, committed-risk, stale preimages, unsafe topology, and non-collision refusals non-retryable.
- Update JSON/text parity, action/reference predicates, and behavioral tests from actual emitted values rather than regex-only prose.

## Red-Green Proof

**RED prompt/case:** Exercise a planning repair-ID collision and an injected post-mutation collision. The former emits an empty rollback status instead of `not_needed`; the latter loses collision identity behind generic `mutation_failed`, so neither matches the documented safe retry predicate.

**Why RED now:** Publication refusal and Git transaction result layers do not carry the semantic collision outcome the orchestration consumes.

**GREEN when:** End-to-end typed-result fixtures prove planning collision and fully rolled-back collision each retry, while incomplete/committed-risk and non-collision outcomes stop without rescan.

**Validation:** Post-remediation review finding from REQ-492; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Cancelled

- **When:** 2026-09-03T11:47:24Z
- **Why:** edge-case hardening with no observed occurrence; recapture if a retry-contract mismatch is seen in a real run (maintainer decision, 2026-09-03 roadmap triage)
- **Decided by:** user, via `do-work abandon`
