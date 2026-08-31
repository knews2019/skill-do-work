---
id: UR-085
title: 'Capture eight uncaptured accepted validate-feedback root causes'
created_at: 2026-08-31T20:49:21Z
requests: [REQ-450, REQ-451, REQ-452, REQ-453, REQ-454, REQ-455, REQ-456, REQ-457]
word_count: 67
---

# Capture Eight Uncaptured Accepted Validate-Feedback Root Causes

## Summary

Capture the eight accepted root causes represented by validate-feedback Findings #1, #3, #4, #7, #13, #14, #15, and #16. Preserve the verbatim claims, duplicate severities and source numbers, validation Evidence, and Surface-cost result. Consolidate duplicate comments by root cause rather than minting one request per comment.

## Extracted Requests

| Request | Root cause | Source findings |
|---|---|---|
| REQ-450 (Exclude already-claimed requests before selection) | Claim evidence is omitted from selection eligibility | #1 P1, #9 P1, #17 P2 |
| REQ-451 (Make confirmation input interruptible) | Blocking confirmation input prevents signal shutdown | #3 P2 |
| REQ-452 (Refuse ambiguous explicit request IDs) | Explicit targeting bypasses duplicate-record ambiguity | #4 P2, #8 P1, #18 P2 |
| REQ-453 (Keep targeted UR dependency closures in the run) | Targeted UR selection drops pending dependents selected for the same run | #7 P1 |
| REQ-454 (Expose UR source tokens in selection records) | Source-token provenance is discarded from result records | #13 P2 |
| REQ-455 (Summarize estimates for the complete run set) | Estimate summaries cover only immediately selected work | #14 P2 |
| REQ-456 (Wait for theme transitions before contrast measurement) | Browser contrast tests sample transitioning colors | #15 P2 |
| REQ-457 (Record cleanup move destinations after exclusive creation) | Rollback ownership is recorded before the process creates the destination | #16 P1 |

## Batch Constraints

- Preserve each finding's verbatim claim, every duplicate comment's severity and source number, the validation Evidence, and the Surface-cost result.
- Consolidate Findings #1/#9/#17 into REQ-450 and Findings #4/#8/#18 into REQ-452.
- Do not duplicate existing accepted-feedback requests.
- No dependency ordering is imposed among the eight independently accepted root causes; shared implementation files alone are not prerequisites.

## Existing Coverage and Reconciliation

- REQ-438 (Refuse mismatched Git transaction roots) already captures Finding #2.
- REQ-439 (Anchor trailing windows before display padding) already captures Findings #11/#19.
- REQ-440 (Refuse non-file static board targets) already captures Finding #5.
- REQ-441 (Validate HTTP archives before publication) already captures Finding #6.
- REQ-442 (Reserve time for untimed claimed work) already captures Finding #12.
- Finding #10 is already fixed by commit `2c82ef12`.
- REQ-443 (Keep Git fallback archive prefixes to one component) remains pending but now needs reconciliation with commit `2c82ef12`; this capture records the mismatch without duplicating or editing REQ-443.

## Full Verbatim Input

> ````text
> do-work capture-request: Capture the eight uncaptured accepted root causes from validate-feedback Findings #1, #3, #4, #7, #13, #14, #15, and #16. Preserve each finding’s verbatim claim, all duplicate comment severities/source numbers, Evidence, and Surface-cost result. Consolidate #1/#9/#17 and #4/#8/#18 into one REQ per root cause. Do not duplicate REQ-438 through REQ-442; record that #10 is already fixed by commit 2c82ef12 and that pending REQ-443 now needs reconciliation.
> ````

---
*Captured: 2026-08-31T20:49:21Z*
