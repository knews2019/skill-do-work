---
id: REQ-473
title: 'Review fix: Support documented outside-target BKB initialization'
status: pending
created_at: 2026-09-01T05:59:27Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-416]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-416, REQ-457]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-416
---

# Support Documented Outside-Target BKB Initialization

## What

Make `bkb-init` honor the shipped absolute and parent-relative target contract even when invoked from inside another Git checkout. The target must be treated as the standalone initialization root when it is outside the invocation repository, and every success/refusal must preserve the exact user target in next and verification argv.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, whose root cause is cross-repository BKB target routing. REQ-457 owns object identity during rollback, not selection of the repository/standalone authority for a documented outside target.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

- Route absolute and parent-relative new targets outside the invocation Git root to the standalone BKB initialization flow without weakening path safety.
- Preserve exact user-supplied targets in text/JSON next argv, Just recipes, and verification argv for success, dry-run, and every refusal.
- Cover invocation inside Git, target inside the same Git root, target outside it, symlink-spelled paths, spaces, and ambiguous/unavailable Git evidence.
- Differentially prove the documented `bkb init ~/research` behavior and retain all rooted publication/rollback protections.

## Red-Green Proof

**RED prompt/case:** From inside a Git checkout, run `bkb-init --kb <absolute outside path> --dry-run` and the equivalent parent-relative target.
**Why RED now:** The handler routes the outside target through the invocation repository, rejects it as escaping, and drops the target from recovery/verification argv.
**GREEN when:** Both outside target shapes use the safe standalone plan, preserve exact argv, and leave the invocation repository untouched.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-416-rereview.md`.

---
*Source: REQ-416 fresh re-review residual F2.*
