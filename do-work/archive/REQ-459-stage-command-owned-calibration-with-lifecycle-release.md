---
id: REQ-459
title: 'Review fix: Stage command-owned calibration with lifecycle release'
status: completed
domain: general
created_at: 2026-08-31T21:54:18Z
user_request: UR-081
addendum_to: REQ-412
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
tdd: true
prime_files: [_dev/primes/prime-action-files.md]
claimed_at: 2026-09-01T06:08:17Z
route: B
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-01T06:08:17Z
  basis:
    - Route B
    - 3-file write set
    - one structural contract
completed_at: 2026-09-01T06:45:28Z
commit: 1c0132399f2fbe2abe57e7280175e2565c848044
---

# Review Fix: Stage Command-Owned Calibration with Lifecycle Release

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Froze the two action documents and the aggregate contract as the exact three-file scope.
- [x] **[APPLY]:** Re-keyed all four staging surfaces to the canonical completion result and added an isolated RED/GREEN contract.
- [x] **[UNIFY]:** Reviewed the exact three-file diff; focused, aggregate, canonical, scope, and diff checks passed.

## What

Make every serial and worktree Commit Phase staging instruction include `do-work/calibration-log.tsv` when the canonical `requeststate complete` transaction changed it, and pin the current command-owned completion step rather than the removed manual calibration substep.

Fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this command-owned lifecycle-staging root cause. REQ-448 mentions calibration while adding phase timestamps, but does not own the release commit's target-set condition or its regression contract.

## Context

Found during the independent re-review of REQ-412. Requeststate now appends calibration inside canonical completion, but active `work.md` and `work-reference.md` staging restatements still key that path on removed “Step 8 substep 7.5,” so an operator can omit the mutation from the release commit and leave the repository dirty.

## Requirements

- Update both serial and worktree Commit Phase staging text to key calibration on the canonical complete transaction's reported target set.
- Remove the stale Step 8 substep 7.5 condition everywhere it governs calibration staging.
- Add a contract regression that fails if either active staging path can omit command-owned calibration or cites the removed substep.

## Red-Green Proof

**RED prompt/case:** Scan the serial and worktree release staging instructions after a canonical complete result containing `do-work/calibration-log.tsv`; require both paths to stage that reported mutation and forbid the removed Step 8 substep 7.5 condition.
**Why RED now:** Both active instructions still condition staging on a substep deleted when REQ-412 moved calibration into `requeststate complete`.
**GREEN when:** The same contract passes and every release path stages command-owned calibration from the canonical result without a stale step-number dependency.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Implementation Summary

The serial and worktree Commit Phase instructions now stage `do-work/calibration-log.tsv` only when the successful canonical `complete` result reports that path among its changed or affected targets. The obsolete Step 8 substep 7.5 staging predicate is gone, and an independent structural regression covers all four active staging surfaces.

**Files changed:**
- `_dev/tests/contract-regressions.sh`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/work.md`

**Builder commit:** `133bc25bd689a25e1946a05de9e05bb95b29523c`

**Integration merge:** `1c0132399f2fbe2abe57e7280175e2565c848044`

## Qualification and Testing

- The exact four-surface contract is RED on baseline `f64b6873` and GREEN on the integrated tree.
- The full contract-regression suite and canonical maintainer verification pass.
- Scope, merge-integrity, and `git diff --check` verification pass; the browser lane had its standard no-browser skip.

## Review

**Overall: 99%** | **Risk: Low** | **Acceptance: Pass**

Independent review found no Important findings or blockers. Its only Minor observation—the unchecked P-A-U metadata—is closed in this archived record.

## Lessons Learned

- Lifecycle staging must follow the canonical transaction's reported target set, not a remembered step number or filesystem inference.
