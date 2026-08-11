---
id: REQ-172
title: Make screenshot source cleanup best-effort
status: pending
created_at: 2026-08-11T17:00:04Z
user_request: UR-039
domain: testing
prime_files: [skills/do-work/tools/prime-do-work-update.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
related: [REQ-173, REQ-174]
batch: accepted-p2-fixes
write_set: [skills/do-work/actions/capture.md, _dev/tests/staged-skills-contract.sh]
---

# Make Screenshot Source Cleanup Best-Effort

## What

After a screenshot has been copied, byte-verified, and installed without clobbering, a failure to remove the staged source must warn without invalidating the permanent asset.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The current fatal cleanup branch leaves a verified destination in place but returns failure. The documented retry then fails at the no-clobber hard link because the destination already exists.

## Detailed Requirements

- Keep copy verification and no-clobber installation strict.
- Treat staged-source removal as warned best-effort cleanup after verified installation.
- Keep the installed destination valid and referenceable when source cleanup fails.
- Do not add destination rollback machinery.
- Add a regression that fails only the staged-source removal.

## Constraints

- Preserve the existing best-effort semantics for `.copying` and empty-directory cleanup.
- The regression must assert the verified destination remains intact, the staged source remains, the `.copying` file is absent, and later collisions remain no-clobber-safe.

## Red-Green Proof

**RED prompt/case:** Force only removal of the staged screenshot source to fail after a verified no-clobber install.
**Why RED now:** The capture block returns nonzero, leaves the destination present, and an exact retry fails with `File exists`.
**GREEN when:** Capture succeeds with a warning, preserves the staged source and verified destination, removes `.copying`, and retains no-clobber behavior on later collisions.
**Validation:** User confirmed

## Full Context

See `do-work/user-requests/UR-039/input.md` and the preceding validated-feedback report.

---
*Source: fix accepted*
