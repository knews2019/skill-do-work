---
source_type: req_lesson
req_id: REQ-493
req_path: do-work/archive/REQ-493-complete-repository-gate-deferral-preflight-topology.md
date: 2026-09-02
domain: backend
module: skills/do-work/tools/do-work-cli
tags: [backend, review, complete, repository]
---

# Lessons from REQ-493: Review fix: Complete repository-gate deferral preflight topology

## What the REQ was about

Make `defer-gate` classify every existing publication target and collision before mutation, so valid folds and invalid move topologies are decided entirely during planning. Done means the topology class cannot recur for parent, checkpoint, repair, reservation, or destination paths.

## Solution summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

## What worked

- Classify tracked fold inputs against `HEAD` so staged and unstaged changes share one complete topology observation, then leave the transaction's staged-input guard authoritative; preflight classification grants an exact dirty-input opt-in, not permission to weaken index safety.

## Back-reference

See `do-work/archive/REQ-493-complete-repository-gate-deferral-preflight-topology.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `3c9caf68`.
