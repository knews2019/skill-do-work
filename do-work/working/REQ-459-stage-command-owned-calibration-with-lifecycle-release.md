---
id: REQ-459
title: 'Review fix: Stage command-owned calibration with lifecycle release'
status: claimed
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
---

# Review Fix: Stage Command-Owned Calibration with Lifecycle Release

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
