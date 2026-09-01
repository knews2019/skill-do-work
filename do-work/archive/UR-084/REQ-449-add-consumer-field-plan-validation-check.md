---
id: REQ-449
title: '[impact-rule-change] Add warning-level consumer-field check to Route C plan validation'
status: completed
route: A
created_at: 2026-08-31T20:38:40Z
user_request: UR-084
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-01T23:19:28Z
  basis:
    - trivial short-circuit
required_lessons: [_dev/primes/lessons-action-files.md]
related: [REQ-448]
claimed_at: 2026-09-01T23:19:07Z
dispatch_at: 2026-09-01T23:19:53Z
builder_handback_at: 2026-09-01T23:20:32Z
integration_at: 2026-09-01T23:20:32Z
review_at: 2026-09-01T23:20:32Z
completed_at: 2026-09-01T23:26:52Z
release_at: 2026-09-01T23:27:32Z
commit: bc7408f7
---

# Add Warning-Level Consumer-Field Check to Route C Plan Validation

## What

Add a fourth check to Route C plan validation (`skills/do-work/actions/work.md` § plan validation): for each planned command whose output drives an action-owned mutation, the plan names the per-record fields its consumer reads — identity, provenance, state, and outcome, as applicable. Same footing as the existing three checks: a warning, not a blocker.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route A direct edit: extend the existing validation list in `work.md`; leave `work-reference.md` unchanged because its Plan template does not enumerate the checks.
- [x] **[APPLY]:** Added the consumer-field contract as the fourth warning-only Route C plan-validation check and retained the review-time Restatement Sweep.
- [x] **[UNIFY]:** Reviewed the complete one-line implementation diff and surrounding warning semantics; verified exact categories-as-applicable wording, the Restatement Sweep relationship, and no debug artifacts.

## Detailed Requirements

- The check applies to "every command whose output drives an action-owned mutation" — the plan must "identify the exact per-record identity, provenance, state, and outcome fields required by its consumer" (categories as applicable, not a mandatory four-slot form).
- Warning-not-blocker, matching the existing plan-validation contract ("These are warnings, not blockers — the builder can adapt").
- Update the Plan Template (`skills/do-work/actions/work-reference.md` § Plan Template — Route C) if it enumerates the validation checks, in the same commit.
- State the relationship to the review-time Restatement Sweep (`skills/do-work/actions/review-work.md`): this is its plan-time counterpart, not a replacement — the sweep still runs at review.

## Constraints

- Explicitly excluded, per user decision at capture: no mandatory end-to-end consumer test, and the check is not named a "preflight" — that word is taken by Step 5.75 (`tools/checks/preflight.sh`).
- Route C only, like the rest of plan validation; Routes A and B skip planning and are untouched.

## Red-Green Proof
**RED prompt/case:** Plan validation in `skills/do-work/actions/work.md` lists three checks (requirement coverage, no orphan tasks, scope sanity). A Route C plan whose command output feeds an action-owned mutation passes validation without naming any field its consumer reads — the field-loss class of bug is first caught at review by the Restatement Sweep, after the build.
**Why RED now:** No consumer/contract notion exists in plan validation; grep for a fourth check returns nothing.
**GREEN when:** Plan validation lists the fourth warning-level check requiring per-record consumer-read fields to be named, and a Route C plan skipping it draws a warning while the pipeline proceeds.
**Validation:** User confirmed (chose the warning-level shape over the original hard requirement; test clause and "preflight" name excluded)

## Full Context
See `do-work/user-requests/UR-084/input.md` for complete verbatim input.

---
*Source: "During plan validation, require every command whose output drives an action-owned mutation to identify the exact per-record identity, provenance, state, and outcome fields required by its consumer, plus one end-to-end consumer test." (softened to warning-level at capture; consumer-test clause dropped)*

## Triage

**Route: A** - Simple

**Reasoning:** The existing Route C validation list needs one warning-only item in a single action file. The reference Plan template does not enumerate validation checks, so it requires no synchronized edit.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Root Cause

Route C plan validation checked requirement coverage, orphan work, and task count, but did not require a plan to preserve the per-record evidence an action-owned mutation consumes. That deferred identity/provenance/state/outcome field-loss detection until the review-time Restatement Sweep.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)

**What was done:** Added a fourth warning-level Route C validation check requiring each mutation-driving command plan to name the exact per-record consumer fields that apply. The text explicitly keeps the review-time Restatement Sweep as an independent later check.

## Qualification

Passed — the single-file Route A diff directly traces all requirements. The existing paragraph continues to make every listed check warning-only, the new item uses categories as applicable rather than a mandatory four-slot form, and the conditional Plan-template edit is correctly skipped because that template does not enumerate validation checks.

## Testing

**Tests run:** focused text inspection; `git diff --check`; canonical maintainer verification
**Result:** ✓ All passing. Canonical maintainer verification passed its action contracts, uncached queue-board suite, strict JavaScript lane, and full CLI suite; the optional external-browser lane was unavailable and skipped.

**Red-green validation:** The baseline validation list had exactly three checks and no consumer-field contract; the changed list has a fourth check with the required consumer evidence and warning-only continuation semantics.

**New tests added:** None — the user explicitly excluded a mandatory consumer test, and this is an instruction-only contract change.

*Verified by work action*

## Review

**Overall: 100%** | 2026-09-01T23:20:32Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — the new plan-time warning closes the field-contract omission without weakening or replacing review-time reconciliation.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by work action*

## Orientation

Route C planning now checks the structured evidence boundary between decision commands and action-owned mutations before implementation starts. The action-file prime remains current.
