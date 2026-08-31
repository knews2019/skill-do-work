---
id: REQ-449
title: '[impact-rule-change] Add warning-level consumer-field check to Route C plan validation'
status: pending
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
related: [REQ-448]
---

# Add Warning-Level Consumer-Field Check to Route C Plan Validation

## What

Add a fourth check to Route C plan validation (`skills/do-work/actions/work.md` § plan validation): for each planned command whose output drives an action-owned mutation, the plan names the per-record fields its consumer reads — identity, provenance, state, and outcome, as applicable. Same footing as the existing three checks: a warning, not a blocker.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
