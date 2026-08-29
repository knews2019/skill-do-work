---
id: REQ-421
title: '[impact-rule-change] Keep board corpus floors suite-only'
status: pending
created_at: 2026-08-29T20:26:10Z
user_request: UR-082
domain: testing
prime_files: ['_dev/primes/prime-kanban-board.md']
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-422, REQ-423, REQ-424]
batch: accepted-review-fixes
write_set: ['skills/do-work-board/tools/queue-kanban/citations_test.go']
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
---

# Keep Board Corpus Floors Suite-Only

## What
Keep citation, fence, and shipped-payload invariants active in every checkout, while applying their suite-calibrated numeric corpus floors only when `suiteCheckoutSkipReason` confirms a suite checkout.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Guard only the three suite-sized fixed floors in the citation, fence, and shipped-payload tests with the existing `suiteCheckoutSkipReason` mechanism.
- Continue running the semantic invariants in consumer checkouts.
- Add a consumer-shaped subprocess regression that runs those exact tests against a small queue and proves they pass.

## Constraints
- Do not weaken citation, fence, or shipped-payload correctness checks.
- Reuse the existing REQ-303 consumer-corpus lesson rather than adding a duplicate lesson.
- No public interface or schema change.

## Dependencies
None. It ships in the same release batch as REQ-422, REQ-423, and REQ-424.

## Builder Guidance
Certainty: Firm. The accepted review identified the exact assertions and the existing suite-checkout guard to reuse.

## Context
No pending or unassigned queue candidate shares this root cause. Provenance: accepted review finding `[P1] Skip corpus-size assertions outside suite checkouts` against `skills/do-work-board/tools/queue-kanban/citations_test.go:1124-1127`. The review observed consumer counts of 17 and 14 and requested `suiteCheckoutSkipReason` for corpus-calibrated assertions.

## Red-Green Proof
**RED prompt/case:** Run the three named tests in a subprocess whose repository contains a small consumer-shaped `do-work/` queue.
**Why RED now:** Their fixed suite-sized floors fail even when mention generation is correct.
**GREEN when:** All semantic checks still run and the consumer-shaped subprocess passes because only the numeric floors are suite-gated.
**Validation:** User accepted the review finding and supplied the implementation plan.

## Full Context
See `do-work/user-requests/UR-082/input.md` for the approved plan and batch constraints.

---
*Source: accepted review finding [P1], followed by the user-approved “Fix Accepted Review Findings” plan.*
