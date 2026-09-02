---
id: REQ-493
title: '[impact-rule-change] Review fix: Complete repository-gate deferral preflight topology'
status: claimed
route: C
domain: backend
created_at: 2026-09-02T02:24:22Z
user_request: UR-095
addendum_to: REQ-491
review_generated: true
impact: impact-rule-change
effort_estimate: effort-substantive
tdd: true
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-action-files.md]
depends_on: [REQ-491]
related: [REQ-492]
sweep: true
sweep_key: defer-gate-preflight-topology-incomplete
claimed_at: 2026-09-02T03:48:12Z
planning_at: 2026-09-02T03:51:55Z
write_set: [skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go, skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go, skills/do-work/actions/work-reference.md]
---

# Review Fix: Complete Repository-Gate Deferral Preflight Topology

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route C. Add planning-time `Lstat` refusal for an occupied parent queue destination, classify the folded repair with parent/checkpoint preimages, retain the transaction's staged-input and final-boundary guards, and pin success plus every rollback position before changing production code.
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Make `defer-gate` classify every existing publication target and collision before mutation, so valid folds and invalid move topologies are decided entirely during planning. Done means the topology class cannot recur for parent, checkpoint, repair, reservation, or destination paths.

## Context

Found during the post-remediation review of REQ-491 (Add canonical repository-gate deferral lifecycle). The core transaction works and its original review blockers are closed, but two topology facets remained outside the preflight classifier.

## Requirements

- Classify an existing repair path independently as tracked-dirty, tracked-clean, or untracked before a fold transaction is planned.
- Permit a manifest-bound tracked-dirty repair fold through the same exact-preimage opt-in used for other dirty targets; continue to refuse staged repair inputs.
- Cover successful fold and full rollback at every injected mutation position for the tracked-dirty repair state.
- Refuse an occupied parent queue destination while building the plan, before reservation, repair, parent, or checkpoint mutation begins.
- Preserve the atomic rollback and exact identity guarantees delivered by REQ-491.
- Keep REQ-492 (Integrate repository-gate deferral and resumption into `do-work run`) compatible with the completed topology contract.

## Red-Green Proof

**RED prompt/case:** (1) Fold a same-fingerprint deferral into a manifest-bound tracked-dirty repair file; the current planner refuses it because only parent/checkpoint paths receive dirty classification. (2) Supply an already occupied parent queue destination; the current planner reaches apply before the move collision refuses.

**Why RED now:** The shared target classifier is not applied to the fold repair, and destination existence is not validated during plan construction.

**GREEN when:** Named regression tests prove tracked-dirty repair fold success plus exact rollback at every mutation position, and prove an occupied parent destination returns a typed refusal with zero mutations.

**Validation:** Post-remediation review finding from REQ-491; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Instances

- [ ] `impact-user-visible` — A same-fingerprint parent cannot fold into a tracked-dirty repair even though the manifest binds the repair’s exact preimage.
- [ ] `impact-rule-change` — An occupied parent queue destination is detected during apply instead of at the no-mutation planning boundary.

## Scope

- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` — classify the folded repair and refuse occupied parent destinations during planning.
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` — prove RED/GREEN success, staged refusal, planning collision, and exact rollback at every fold mutation.
- `skills/do-work/actions/work-reference.md` — align the canonical topology contract with all classified preimages and move destinations.

## Plan

1. Add named failing tests for a manifest-bound tracked-dirty fold and an occupied parent queue destination; record the exact pre-fix failures.
2. Add rollback coverage after each of the three fold mutations and retain a focused staged-repair refusal test.
3. Preflight the parent queue destination before any mutation is appended, while preserving exclusive apply-time creation as the race guard.
4. Include the fold repair in shared preimage classification, deduplicating its tracked-dirty or untracked opt-in without weakening absent-reservation authority.
5. Update the owning action contract, close both review instances only after focused GREEN, and run race, vet, full CLI, and canonical maintainer verification.
