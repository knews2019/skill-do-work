---
id: REQ-434
title: 'Review fix: Refuse unsupported timestamp ordering anchors'
status: claimed
domain: general
created_at: 2026-08-30T22:01:41Z
claimed_at: 2026-08-31T15:44:18Z
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-31T15:44:18Z
  basis:
    - trivial short-circuit
user_request: UR-081
addendum_to: REQ-410
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
tdd: true
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
---

# Review Fix: Refuse Unsupported Timestamp Ordering Anchors

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What
Prevent unsupported but parseable timestamp shapes from becoming ordering anchors or changing supported successor fields during doctor repair.

## Context
Found during terminal re-review of REQ-410 after its single remediation. An offset or fractional `created_at` is refused as a repair target but still participates as the predecessor for `claimed_at`; doctor can therefore clamp a supported successor to the unsupported instant, potentially beyond `now`, and commit it. The legacy predicate orders only fields with a supported comparison key. Fold-first scan found no pending REQ or sweep in any UR that shares this mixed-shape ordering-anchor root cause.

## Requirements
- Exclude unsupported timestamp shapes from repair ordering anchors as well as from direct repair.
- Keep every unsupported field byte-identical while still diagnosing and refusing it explicitly.
- Never derive or clamp a supported successor from an unsupported predecessor, including offsets and fractional seconds.
- Preserve the supported `created_at <= claimed_at <= completed_at` ordering behavior and current-time ceiling.

## Red-Green Proof
**RED prompt/case:** Create a record with an offset or fractional `created_at` followed by a supported whole-second `claimed_at`; build and apply the doctor timestamp plan and assert the supported successor is not clamped from the unsupported predecessor.
**Why RED now:** Comparable parse success is currently enough to make the unsupported field an ordering anchor.
**GREEN when:** Mixed unsupported/supported fixtures report the unsupported field, leave both bytes unchanged unless an independently valid repair exists, and all supported-shape ordering fixtures continue to pass.
**Validation:** Terminal review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

---

## Triage

**Route: A** - Simple

**Reasoning:** The unsupported-anchor predicate, exact mixed-shape fixtures, and expected byte-preserving outcome are explicit. This is a focused doctor timestamp-plan correction.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
