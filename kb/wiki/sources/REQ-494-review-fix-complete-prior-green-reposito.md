---
title: "Lessons from REQ-494: Review fix: Complete already-green repository-gate repair lifecycle"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-04/REQ-494-review-fix-complete-prior-green-reposito.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-494: Review fix: Complete already-green repository-gate repair lifecycle

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

Make the already-green repository-gate repair path executable through every downstream lifecycle authority, so a repair whose fingerprint no longer reproduces can complete honestly and release its dependency-gated parents without fake implementation work. Done means no TDD, qualification, review, release, staging, or completion guard can contradict the canonical no-op evidence shape.

## Solution summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/review-work.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

## What worked

- Replacing a synthetic lifecycle tail with real CLI commands is insufficient when the entry decision is still reimplemented in the test. Cross-action exceptions need one executable validator consumed by every action and exercised with actual intake evidence and the exact canonical-result path allowlist.

## Back-reference

See `do-work/archive/REQ-494-complete-already-green-repository-gate-repair-lifecycle.md` for the full REQ — plan, exploration, implementation, review, and lessons.
