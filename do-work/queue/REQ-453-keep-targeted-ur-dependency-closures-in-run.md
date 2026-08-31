---
id: REQ-453
title: 'Keep targeted UR dependency closures in the run'
status: pending
created_at: 2026-08-31T20:49:21Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-450, REQ-451, REQ-452, REQ-454, REQ-455, REQ-456, REQ-457]
batch: accepted-validate-feedback-root-causes
---

# Keep Targeted UR Dependency Closures in the Run

## What

Keep every pending member of a targeted user-request dependency closure in the authoritative run set, then re-evaluate downstream members after their prerequisites integrate. Do not falsely declare a dependent concurrently runnable during fan-out, and do not silently leave it behind when the targeted workflow stops after the returned set.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this targeted-run dependency-closure root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding #7 — P1 — source:** `internal/nextselection/next_selection.go:210-212`

> ````text
> [P1] Keep UR dependencies in the targeted run — [prj].claude/skills/do-work/tools/do-work-cli/
> internal/nextselection/next_selection.go:210-212
> For a targeted UR-NNN containing pending B that depends on pending A, B is emitted as DEPENDENCIES-UNMET even though A is selected
> earlier in the same dependency-ordered target set. The targeted workflow stops after its returned selected set, so do-work run UR-
> NNN completes A and silently leaves B behind instead of draining the requested UR as specified in actions/work.md targeted mode.
> ````

- **Evidence:** UR members are expanded in dependency-depth order at `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go:57-86`, but every member is evaluated against the unchanged repository graph. `next_selection.go:187-213` excludes a pending dependent before an earlier member can integrate. Targeted mode stops on the returned set at `skills/do-work/actions/work.md:185-194`, while the estimate and run contract at lines `172-183` requires the loop to drain dependents.
- **Surface-cost result:** N/A — this is a direct execution-contract correction, not added defensive apparatus.

## Detailed Requirements

- Preserve the complete scoped closure or equivalent deferred run set returned for targeted UR execution.
- Re-evaluate a dependent after a prerequisite from the same targeted run integrates successfully.
- Continue until the targeted UR closure is drained or a genuine terminal blocker occurs.
- Preserve dependency-depth ordering.
- Do not classify a dependent as concurrently runnable while its prerequisite is still pending or in flight.
- Report genuine failed, cancelled, external, or unresolved dependency blockers truthfully.

## Constraints

- The authoritative result must support targeted execution without a prohibited queue rescan.
- Fan-out must retain dependency safety.

## Dependencies

No request prerequisite. This REQ changes how targeted runs represent their own internal prerequisites; it does not depend on another captured REQ.

## Builder Guidance

Certainty level: Firm. Represent deferred-in-this-run separately from permanently excluded work if needed.

## Red-Green Proof

**RED prompt/case:** Target a UR containing pending B depending on pending A, with no external blocker, and execute only from the returned authoritative result.
**Why RED now:** B is excluded against the pre-run graph and targeted mode stops after A, leaving B queued.
**GREEN when:** A runs first, B is retained and re-evaluated after A integrates, and the targeted run drains both without making B concurrently runnable too early.
**Validation:** User confirmed after validate-feedback accepted Finding #7.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding #7, captured by UR-085.*
