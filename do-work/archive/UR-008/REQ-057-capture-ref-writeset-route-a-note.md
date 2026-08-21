---
id: REQ-057
title: "capture-reference.md write_set firming note omits Route A (never runs Step 5.5)"
status: completed
route: A
claimed_at: 2026-07-29T13:55:51Z
commit: a68d8b5
created_at: 2026-07-29T13:04:42Z
status_changed_at: 2026-07-29T13:24:10Z
user_request: UR-008
addendum_to: REQ-045
domain: general
prime_files: []
tdd: false
depends_on: []
write_set: [actions/capture-reference.md]
maintenance: false
review_generated: false
---

# capture-reference.md write_set note omits Route A

## What
Discovered during REQ-045. `actions/capture-reference.md` (~:113) says capture's `write_set` value is one that "the work pipeline's Scope step … firms it up and overwrites it." That is true only for Routes B and C — a **Route A** REQ never runs Step 5.5, so it keeps capture's value for the whole run (and, after REQ-045, that value is what **Step 3** re-validates for a co-dispatched Route A REQ). A one-clause qualification closes the imprecision.

## Open Questions
- [x] Fix this one-clause doc imprecision? → Confirmed: yes — add a clause noting a Route A REQ keeps the capture-seeded `write_set` (it never runs Step 5.5), consistent with REQ-045's Step 3 re-validation.

## Full Context
Discovered task from REQ-045 (dispatch re-validation). See `do-work/archive/UR-008/REQ-045-dispatch-revalidation-route-a-gap.md` → `## Discovered Tasks`.

## Triage

**Route:** A
**Reasoning:** Single file (`actions/capture-reference.md`); one-clause addition noting a Route A REQ keeps the capture-seeded `write_set` (never runs Step 5.5). Answered yes via clarify. Route A.
**Rigor:** Standard main-context review (part of the parallel disjoint-write_set batch 051/052/054/057/058; single-file, no spec-cluster overlap).

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## P-A-U

- [x] **P**lan → **A**ct → **U**nify

UNIFY: The clause lands in the one place a capture-time reader learns that their `write_set` is provisional, and it now states the exception in the same vocabulary the pipeline uses — Routes B/C firm the set at `actions/work.md` Step 5.5, Route A keeps the capture-seeded set and re-validates it at Step 3 (REQ-045). No surrounding text restructured; the "hint, never a commitment" phrase is preserved verbatim because `actions/work.md` Step 5.5 quotes it.

## Implementation Summary

**Files changed:**
- `actions/capture-reference.md` (modified) — qualified the Scope-step-overwrites claim as Routes B/C only; noted a Route A REQ keeps its capture-seeded `write_set` for Step 3 re-validation

**What was done:** In the **Populating `write_set`** paragraph, qualified the claim that the pipeline's Scope step "firms it up and overwrites it" as true for Routes B and C only, and added that a Route A REQ never runs Step 5.5 — so it keeps the capture-seeded value for the whole run, and that value is what `actions/work.md` Step 3 re-validates for disjointness when a Route A REQ was co-dispatched (per REQ-045). One sentence added inside the existing paragraph; no other text touched.
