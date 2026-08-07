---
id: REQ-135
title: "Restore the Warning-Clean ShellCheck Baseline"
status: pending
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: testing
prime_files: []
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
effort_estimate: trivial
related: [REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-143, REQ-144, REQ-145, REQ-146]
batch: do-work-four-skill-suite
write_set: [_dev/tests/contract-regressions.sh]
---

# Restore the Warning-Clean ShellCheck Baseline

## What
Fix the two warning-level `SC2034` diagnostics introduced by the foreign-listener fixture in `_dev/tests/contract-regressions.sh` before beginning modularization.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
The modularization program needs a warning-clean validation baseline so later regressions are attributable to the new work.

## Detailed Requirements
- Deliberately consume or redirect the intentionally unused `foreign_listener_output` value.
- Address the `port` variable that ShellCheck cannot see through `eval` with explicit consumption or a narrowly scoped suppression.
- Do not add a file-wide suppression.
- Do not change foreign-listener runtime behavior.

## Constraints
- This is the clean prerequisite for the full batch.
- Keep the fix local to the test fixture.

## Dependencies
None.

## Builder Guidance
Certainty level: Firm. This is a small mechanical fixture repair.

## Red-Green Proof
**RED prompt/case:** Run `shellcheck -S warning _dev/tests/contract-regressions.sh` at the current revision.
**Why RED now:** ShellCheck exits non-zero with two `SC2034` warnings for `foreign_listener_output` and `port`.
**GREEN when:** Warning-level ShellCheck exits zero while the Bash syntax check, contract regression suite, and queue-kanban tests remain green.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the complete batch decisions.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*
