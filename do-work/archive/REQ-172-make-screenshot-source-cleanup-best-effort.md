---
id: REQ-172
title: Make screenshot source cleanup best-effort
status: completed
claimed_at: 2026-08-11T17:02:09Z
route: A
completed_at: 2026-08-11T17:08:02Z
created_at: 2026-08-11T17:00:04Z
user_request: UR-039
domain: testing
prime_files: [skills/do-work/tools/prime-do-work-update.md]
tdd: true
kb_status: pending
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
- [x] **[PLAN]:** Extend the existing executable screenshot lifecycle harness so `rm` fails only for the staged source, prove that the current block returns nonzero, then remove that fatal return while preserving strict copy verification/no-clobber installation and the existing best-effort temporary-file/directory cleanup. Verify the focused contract, the full staged-skills contract, and ShellCheck for the changed shell surface.
- [x] **[APPLY]:** Added a staged-source-only `rm` failure regression, captured the expected RED result, then removed the post-install fatal return and aligned the capture contract with warned best-effort cleanup. Copy verification, no-clobber installation, and later collision failure remain strict.
- [x] **[UNIFY]:** Reviewed the complete diff for `skills/do-work/actions/capture.md` (only post-install cleanup semantics/prose) and `_dev/tests/staged-skills-contract.sh` (focused source-`rm` failure plus later collision assertions), and reviewed this REQ's phase/root-cause notes. `git diff --check`, Bash syntax, ShellCheck at warning severity, embedded-shell ShellCheck, and the full staged-skills contract all pass; no debug artifacts or similar fatal staged-source cleanup pattern remains.

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

---

## Triage

**Route: A** - Simple

**Reasoning:** The failure is reproduced, the affected action and regression file are named, and the accepted remedy is a narrow behavior correction.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Root Cause

Staged-source deletion was grouped with pre-install copy, verification, and no-clobber failures even though it runs after the permanent asset is already verified and installed. Returning nonzero at that point created a half-success whose unchanged destination made the prescribed retry collide.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/capture.md` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)

**What was done:** Staged-source removal now warns without invalidating an already verified permanent screenshot. The executable lifecycle regression forces only that removal to fail and proves the source, destination, temporary-copy cleanup, and later no-clobber collision behavior.

## Qualification

Passed — 2 implementation files verified, 5 requirements traced, P-A-U confirmed. The action change is substantive and limited to the post-install failure classification; the regression exercises the production-prescribed block through the existing executable harness.

## Testing

**Tests run:** `bash _dev/tests/staged-skills-contract.sh`; `bash -n _dev/tests/staged-skills-contract.sh`; `shellcheck -S warning _dev/tests/staged-skills-contract.sh`; `git diff --check`
**Result:** ✓ All passing

**Red-green validation:**
- Staged-source cleanup lifecycle fixture: ✗ before implementation (`staged-source cleanup failure invalidated a verified capture`) → ✓ after (`staged skills contract: PASS`)

**New tests added:**
- Staged-source-only removal failure preserves the verified install and warning semantics.
- A subsequent destination collision remains strict and no-clobber-safe.

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-11T17:07:42Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — staged-source removal failure now warns and returns success while preserving the source and verified destination; copy verification and later no-clobber collisions remain strict.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Replaying only the staged-source `rm` separated post-install cleanup from the strict installation boundary.
**What didn't:** Treating every cleanup failure as transactional failure created a state that the normal retry could not repair.
**Worth knowing:** Once the permanent asset is byte-verified and no-clobber installed, cleanup warnings must not revoke its validity.

## Orientation

Screenshot capture now keeps a verified permanent asset valid when only staged-source cleanup fails; the no-clobber installation boundary is unchanged.
