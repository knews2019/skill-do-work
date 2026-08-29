---
id: REQ-421
title: '[impact-rule-change] Keep board corpus floors suite-only'
status: completed
claimed_at: 2026-08-29T20:31:49Z
completed_at: 2026-08-29T20:46:08Z
commit: f532dcf2
route: B
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
- [x] **[PLAN]:** Guard only the three fixed numeric floors with the existing suite detector, add a small-checkout subprocess regression that invokes the exact tests, then run the focused tests uncached.
- [x] **[APPLY]:** Added a small consumer subprocess fixture that runs the exact citation, fence, and shipped-payload tests; guarded only their three suite-calibrated numeric floors with `suiteCheckoutSkipReason`.
- [x] **[UNIFY]:** Reviewed the complete diff, ran the focused uncached tests, the full uncached queue-kanban suite, and the canonical maintainer verifier; all required lanes passed.

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

---

## Triage

**Route: B** - Medium

**Reasoning:** The fix is localized, but the subprocess regression must preserve the semantic assertions while changing only the suite-calibrated floors.

**Planning:** Not required; the user supplied an implementation plan.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/citations_test.go` (modify) — suite-gate numeric floors and add the consumer subprocess regression

**Files I will NOT touch:** citation generation or shipped client behavior.

**Acceptance criteria (restated from REQ):**
- [x] The three semantic invariants run in all checkouts.
- [x] Only their fixed corpus floors skip outside suite checkouts.
- [x] A consumer-shaped subprocess runs the exact tests and passes.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/citations_test.go` (modified)

**What was done:** The live-corpus tests now retain their semantic checks in consumer checkouts while suite-only count floors are gated by the existing checkout detector. A subprocess fixture proves the exact three tests pass against a small queue containing both prose and fenced mentions.

## Testing

**Tests run:** Focused citation/fence/payload Go tests; `go test -count=1 ./...`; `bash _dev/tests/maintainer-verify.sh`.

**Result:** All passed. The canonical verifier exited 0; its optional browser lane skipped because no browser was configured.

**Red-green validation:**
- `TestConsumerCheckoutRunsCitationFenceAndShippedPayloadInvariants`: RED with 1 REQ mention, 2 fenced REQs, and 1 shipped mention failing the old floors → GREEN after the three floors were suite-gated.

**New tests added:**
- Consumer-shaped subprocess coverage invoking all three exact live-corpus tests.

## Lessons Learned

REQ-303 (Run the pinned live archive assertions only in a suite checkout) already owns this lesson: a repo-independence test must build a consumer root and invoke the real binary from it. This regression reuses that pattern and adds no duplicate prime-satellite entry.
