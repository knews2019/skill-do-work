---
id: REQ-493
title: '[impact-rule-change] Review fix: Complete repository-gate deferral preflight topology'
status: completed
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
dispatch_at: 2026-09-02T03:52:19Z
builder_handback_at: 2026-09-02T04:02:18Z
integration_at: 2026-09-02T04:02:19Z
review_at: 2026-09-02T04:10:16Z
release_at: 2026-09-02T04:18:52Z
kb_status: pending
write_set: [skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go, skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go, skills/do-work/actions/work-reference.md, _dev/tests/contract-regressions.sh]
completed_at: 2026-09-02T04:17:51Z
commit: 3c9caf68
---

# Review Fix: Complete Repository-Gate Deferral Preflight Topology

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route C. Add planning-time `Lstat` refusal for an occupied parent queue destination, classify the folded repair with parent/checkpoint preimages, retain the transaction's staged-input and final-boundary guards, and pin success plus every rollback position before changing production code.
- [x] **[APPLY]:** Added repair preimage classification, early parent-destination collision refusal, full transaction regressions, and the discovered downstream contract-reader update.
- [x] **[UNIFY]:** Reviewed all four implementation/contract files, verified the diff contains no debug artifacts, and passed focused publication tests plus the full contract regression suite; full CLI and canonical gates follow below.

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

- [x] `impact-user-visible` — A same-fingerprint parent can now fold into a manifest-bound tracked-dirty repair while staged bytes remain refused.
- [x] `impact-rule-change` — An occupied parent queue destination now returns a typed planning refusal with zero admitted mutations.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` — classify the folded repair and refuse occupied parent destinations during planning.
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` — prove RED/GREEN success, staged refusal, planning collision, and exact rollback at every fold mutation.
- `skills/do-work/actions/work-reference.md` — align the canonical topology contract with all classified preimages and move destinations.
- `_dev/tests/contract-regressions.sh` — update the downstream mutation contract from the superseded pre-REQ-493 topology boundary.

**Acceptance criteria:** A present-reservation manifest-bound tracked-dirty repair folds successfully, staged repair bytes remain refused, rollback restores exact bytes/modes/status at every fold mutation, and an occupied parent queue destination refuses during planning with no admitted mutations while the final exclusive apply guard remains intact.

## Plan

1. Add named failing tests for a manifest-bound tracked-dirty fold and an occupied parent queue destination; record the exact pre-fix failures.
2. Add rollback coverage after each of the three fold mutations and retain a focused staged-repair refusal test.
3. Preflight the parent queue destination before any mutation is appended, while preserving exclusive apply-time creation as the race guard.
4. Include the fold repair in shared preimage classification, deduplicating its tracked-dirty or untracked opt-in without weakening absent-reservation authority.
5. Update the owning action contract, close both review instances only after focused GREEN, and run race, vet, full CLI, and canonical maintainer verification.

## Red-Green Evidence

RED before production edits:

- `TestDeferGateFoldAcceptsManifestBoundTrackedDirtyRepair` found only parent/checkpoint in `ExistingDirtyTargetPaths`; the repair was absent.
- `TestDeferGateRefusesOccupiedParentQueueDestinationDuringPlanning` received no planning refusal.

GREEN after implementation:

- The focused four-test fold/collision/staged/rollback lane passes.
- `go test -race -count=1 ./internal/publication`, `go vet ./...`, and `go test -count=1 ./...` pass in the CLI module.
- The repository contract regression suite passes after updating its superseded REQ-492 predicate.

## Scope Drift

The builder's GREEN implementation exposed `_dev/tests/contract-regressions.sh` as a stale downstream reader of the old topology contract. It was added to `write_set` before editing; no other scope changed.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Classified folded repairs alongside parent and checkpoint preimages, authorized exact manifest-bound tracked-dirty folds without accepting staged inputs, and refused occupied parent queue destinations before admitting mutations while retaining the apply-time exclusive move guard. Added success, staged-refusal, zero-mutation collision, and exact rollback coverage for all three fold mutation positions; updated the downstream repository-gate contract reader to the completed topology semantics.

## Qualification

Canonical `qualify` and `scope-drift` both passed with no findings after the downstream contract reader was added to the declared scope.

## Testing

- Focused four-test defer-gate topology lane: passed.
- `go test -race -count=1 ./internal/publication`: passed.
- `go vet ./...`: passed.
- `go test -count=1 ./...`: passed.
- `bash _dev/tests/contract-regressions.sh`: passed.
- `gofmt -d` and `git diff --check`: clean.

## Review

**Verdict:** Approve — no findings; 99% overall, 100% acceptance.

The independent reviewer confirmed both original finding instances are behaviorally closed. Tracked-dirty repair folding, staged refusal, and exact rollback across all three mutations pass; occupied parent destinations refuse during planning with zero admitted mutations; the final exclusive move boundary remains intact. No remediation or follow-up REQ is required.

## Lessons Learned

Classify tracked fold inputs against `HEAD` so staged and unstaged changes share one complete topology observation, then leave the transaction's staged-input guard authoritative; preflight classification grants an exact dirty-input opt-in, not permission to weaken index safety.

## Orientation

- **Start here:** `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` in `BuildDeferGatePlan` and `classifyGatePreimages`.
- **Read next:** the four REQ-493 tests in `defer_gate_test.go`, then the Repository Gate Deferral and Resumption contract in `work-reference.md`.
- **Avoid:** weakening absent-reservation authority or removing `moveRootedFile`'s final exclusive destination check.
- **Verify:** publication race tests, full CLI tests, contract regressions, then the canonical maintainer gate.
