---
id: REQ-494
title: '[impact-critical] Review fix: Complete already-green repository-gate repair lifecycle'
status: completed-with-issues
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
dispatch_at: 2026-09-02T04:28:43Z
builder_handback_at: 2026-09-02T04:40:04Z
integration_at: 2026-09-02T04:40:05Z
review_at: 2026-09-02T04:43:53Z
remediation_at: 2026-09-02T04:50:18Z
re_review_at: 2026-09-02T04:53:21Z
kb_status: pending
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/review-work.md, _dev/tests/contract-regressions.sh]
completed_at: 2026-09-02T04:54:36Z
---

# Review Fix: Complete Already-Green Repository-Gate Repair Lifecycle

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route C. Reuse the canonical exact no-op predicate, add narrow TDD and orchestrated-review branches, and prove the complete lifecycle plus negative neighboring cases through the existing semantic mutation harness.
- [x] **[APPLY]:** Added the exact no-op exception to the TDD and orchestrated-review guards plus a full lifecycle semantic mutation matrix.
- [x] **[UNIFY]:** Reviewed all three scoped files, verified the two action changes are narrow, ran the full contract suite, and confirmed `git diff --check` is clean with no debug artifacts.

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

- [x] `impact-critical` — The canonical exact already-green repair evidence now precedes and narrowly bypasses ordinary `tdd: true` RED/GREEN enforcement.
- [x] `impact-critical` — Orchestrated `review-work` now validates that evidence before the ordinary empty-diff exit.

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

## Red-Green Evidence

RED before action edits: the new semantic fixture reported all 16 missing TDD/review predicates, including exact marker/evidence, green-before-edit, empty diff, intake identity, direct replay, release exclusion, stage allowlist, and negative fallbacks.

GREEN after action edits: `bash _dev/tests/contract-regressions.sh` passes, including canonical lifecycle, seven negative fixtures, and 16 widening/removal mutations.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/review-work.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Added a narrow exact-predicate exception before ordinary TDD enforcement and before orchestrated review's empty-diff exit. The exception requires the canonical claimed repair marker, durable no-op history/summary/qualification, direct green-before-edit gate evidence, empty project diff, no release mutation, and lifecycle-only staging; ordinary, malformed, nonempty, or over-staged work retains the original refusal path. The semantic fixture proves completion/archive, release exclusion, lifecycle-only commit, metadata flow, and parent reselection without changing CLI runtime interfaces.

## Review

**Verdict:** Request changes — 50%, Acceptance Fail, critical risk.

- **Critical — impact-critical:** The semantic fixture derives completion, release exclusion, lifecycle commit, and parent reselection from its own synthetic truth table rather than consuming actual REQ/Git state and invoking canonical completion, staging, metadata, and selector authorities. It can remain green while the downstream lifecycle is non-executable, so neither original finding satisfies the Finding-Closure Ratchet yet.

One remediation attempt is authorized: replace the parallel truth table with behavior-level proof through real repository state and canonical authorities while retaining the narrow negative cases.

## Remediation

Removed the synthetic lifecycle truth table. The replacement fixture creates real temporary Git repositories, parses exact claimed repair evidence, executes the recorded gate argv, derives TDD/review eligibility from real worktree and index state, rejects seven neighboring invalid cases, runs canonical `complete --commit` plus hash metadata, verifies exact commit paths and untouched release files, and proves canonical `next` selects the parent only after repair completion. The full contract suite and diff check pass.

## Re-Review

**Verdict:** Request changes — 50%, Acceptance Fail, critical risk.

- **Critical — impact-critical:** The later lifecycle now uses real Git repositories and canonical commands, but `action_decisions()` still reimplements the TDD/review decision. The shipped action guards remain phrase/mutation checked, fingerprint identity is self-asserted instead of sourced from repair intake, and the stage allowlist accepts any archive path instead of exact canonical-completion result paths. Neither original instance closes under the Finding-Closure Ratchet.

The remediation allowance is exhausted. REQ-496 carries the mandatory shared-validator closure.

## Qualification

Canonical qualification and scope-drift checks pass for the declared files; review acceptance remains failed for the behavioral authority gap.

## Testing

- Full contract regression suite: passed after remediation.
- Scoped `git diff --check`: passed.
- Independent review and fresh re-review: both failed acceptance at 50% because the decision oracle remains disconnected from the shipped guards.

## Lessons Learned

Replacing a synthetic lifecycle tail with real CLI commands is insufficient when the entry decision is still reimplemented in the test. Cross-action exceptions need one executable validator consumed by every action and exercised with actual intake evidence and the exact canonical-result path allowlist.

## Orientation

- **Start here:** `_dev/tests/contract-regressions.sh` at the REQ-494 fixture and `action_decisions()`.
- **Read next:** the no-op contract in `work-reference.md`, then the two guards in `work.md` and `review-work.md`.
- **Avoid:** another parallel predicate or an archive-directory prefix allowlist.
- **Continue with:** REQ-496.
