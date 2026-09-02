---
id: REQ-494
title: '[impact-critical] Review fix: Complete already-green repository-gate repair lifecycle'
status: claimed
route: C
domain: general
created_at: 2026-09-02T03:32:03Z
user_request: UR-095
addendum_to: REQ-492
review_generated: true
impact: impact-critical
effort_estimate: effort-substantive
tdd: true
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
depends_on: [REQ-492]
related: [REQ-491, REQ-493]
sweep: true
sweep_key: repository-gate-repair-noop-downstream-guards
claimed_at: 2026-09-02T04:25:25Z
planning_at: 2026-09-02T04:28:18Z
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/review-work.md, _dev/tests/contract-regressions.sh]
---

# Review Fix: Complete Already-Green Repository-Gate Repair Lifecycle

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route C. Reuse the canonical exact no-op predicate, add narrow TDD and orchestrated-review branches, and prove the complete lifecycle plus negative neighboring cases through the existing semantic mutation harness.
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Make the already-green repository-gate repair path executable through every downstream lifecycle authority, so a repair whose fingerprint no longer reproduces can complete honestly and release its dependency-gated parents without fake implementation work. Done means no TDD, qualification, review, release, staging, or completion guard can contradict the canonical no-op evidence shape.

## Context

Found during post-remediation review of REQ-492 (Integrate repository-gate deferral and resumption into `do-work run`). The work action names and partially specifies a reviewed no-change completion, but generated repairs remain `tdd: true` and the independent review action still exits on an empty implementation diff.

## Requirements

- Define one exact durable marker/evidence shape for an already-green `repository_gate_repair: true` request.
- Exempt only that exact shape from the ordinary runnable RED/GREEN requirement; all other `tdd: true` requests remain unchanged.
- Make orchestrated independent review validate the no-op gate evidence instead of exiting for no implementation diff.
- Preserve narrow no-diff qualification, canonical complete/archive, lifecycle-only staging/commit, and no version/changelog mutation.
- Prove terminal success makes dependency parents selectable again.
- Sweep every active action/reference/contract consumer of implementation-diff, TDD, qualification, review, release, and commit guards.

## Red-Green Proof

**RED prompt/case:** Select a generated repair whose recorded diagnostic fingerprint no longer reproduces. The current `tdd: true` guard returns it to implementation, and the orchestrated review action exits on its empty implementation diff, so it cannot complete.

**Why RED now:** The special branch is defined only in work-action prose; downstream authorities retain unconditional ordinary-work guards.

**GREEN when:** A semantic fixture drives the exact already-green repair through qualification, independent no-op review, canonical completion/archive, lifecycle-only commit validation, and parent reselection, while neighboring empty-diff requests still refuse.

**Validation:** Post-remediation review finding from REQ-492; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Instances

- [ ] `impact-critical` — The unconditional `tdd: true` RED/GREEN guard rejects the canonical already-green repair evidence.
- [ ] `impact-critical` — Orchestrated `review-work` exits on the same repair’s intentionally empty implementation diff.

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/review-work.md`
- `_dev/tests/contract-regressions.sh`

**Acceptance criteria:** Only a claimed `repository_gate_repair: true` request carrying the canonical exact no-op evidence and summary bypasses ordinary RED/GREEN and empty-diff review guards; malformed and ordinary empty work still refuses, while canonical completion, lifecycle-only commit, metadata recording, and dependency reselection remain executable without a release mutation.

## Plan

1. Add a semantic RED fixture that drives the canonical no-op shape through the active TDD and review guards and proves the current contradictions.
2. Add negative fixtures for ordinary or malformed empty-diff work, nonempty project diffs, release plans, and forbidden staged paths.
3. Exempt only the canonical exact predicate from Step 6.5 RED/GREEN and route review-work to lifecycle evidence before its ordinary empty-diff exit.
4. Preserve the existing canonical qualification, complete/archive, lifecycle-only staging, metadata, and selector contracts; no CLI schema or command changes.
5. Mutation-test both narrow branches, run the full contract suite, qualification, independent review, and canonical maintainer gate.
