---
id: REQ-145
title: "Remove the Stateful Pipeline"
status: pending
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: []
tdd: true
suggested_spec: refactor
depends_on: [REQ-144]
maintenance: true
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-143, REQ-144, REQ-146]
batch: do-work-four-skill-suite
---

# Remove the Stateful Pipeline

## What
Remove the separate resumable pipeline state machine after modular cutover and replace it with a copyable full-cycle prompt.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Pipeline duplicates composition and persistent state while the existing `do-work run` orchestrator already owns testing and review.

## Detailed Requirements
- Delete pipeline action/reference, routing, `pipeline.json` lifecycle, pipeline guard hook, session-start reporting, tests, and stale docs.
- Do not weaken `do-work run`; retain triage, implementation, testing, review, lessons, archival, and commit behavior.
- Add the approved copyable prompt to core help and documentation.
- The prompt must capture the request, record its UR, verify it, run its REQs with built-in tests/review, invoke `do-work-toolbox present-work`, and stop/report on failure.
- Explain that testing is already inside `do-work run`; do not create another stateful testing stage.
- Do not retain a pipeline alias or recreate pipeline state under another name.

## Constraints
- This happens only after REQ-144 activates the suite.
- Preserve the exact approved full-cycle prompt in the UR context.

## Dependencies
Requires REQ-144.

## Builder Guidance
Certainty level: Firm. This is a deletion/narrowing pass; remove pipeline-specific machinery rather than adding another orchestration layer.

## Red-Green Proof
**RED prompt/case:** Inspect core routing/help and the runtime tree for pipeline routes, actions, hooks, or state-file handling.
**Why RED now:** Stateful pipeline machinery is present and advertised as a core workflow.
**GREEN when:** No pipeline runtime surface remains, the full-cycle prompt is copyable and accurate, and `do-work run` still passes its full orchestrator tests.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the exact replacement prompt.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*
