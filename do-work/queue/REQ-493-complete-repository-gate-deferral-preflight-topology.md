---
id: REQ-493
title: '[impact-rule-change] Review fix: Complete repository-gate deferral preflight topology'
status: pending
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
---

# Review Fix: Complete Repository-Gate Deferral Preflight Topology

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
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
