---
id: REQ-135
title: "Restore the Warning-Clean ShellCheck Baseline"
status: completed
completed_at: 2026-08-07T20:33:34Z
commit: 99901ae
claimed_at: 2026-08-07T20:31:33Z
route: A
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
- [x] **[PLAN]:** Add no-op reads at the two fixture seams so ShellCheck can see the intentionally captured values without changing the simulated listener behavior; prove the warning-level baseline RED before editing and GREEN afterward.
- [x] **[APPLY]:** Added explicit no-op reads for `port` and `foreign_listener_output` inside the existing foreign-listener fixture; no fixture control flow or assertions changed.
- [x] **[UNIFY]:** Reviewed the scoped `_dev/tests/contract-regressions.sh` hunk for behavior changes and debug artifacts; warning-level ShellCheck, Bash syntax, contract regressions, and queue-kanban Go tests pass.

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

---

## Triage

**Route: A** - Simple

**Reasoning:** The regression is two known warning diagnostics in one test fixture with an explicitly constrained, behavior-preserving remedy.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Explicitly consumed the fixture's captured output and the `port` value read indirectly through `eval`, removing both warning-level `SC2034` diagnostics without changing the foreign-listener simulation or its assertions. The root cause was ShellCheck's inability to trace the dynamic `eval` read and an intentionally unused command-substitution result.

## Qualification

Passed — 1 implementation file verified, all 4 detailed requirements traced to the scoped hunk, P-A-U confirmed, and no hollow wiring or debug artifacts found.

## Testing

**Tests run:** `shellcheck -S warning _dev/tests/contract-regressions.sh`; `bash -n _dev/tests/contract-regressions.sh`; `/bin/bash _dev/tests/contract-regressions.sh`; `(cd tools/queue-kanban && go test ./...)`
**Result:** ✓ All passing

**Red-green validation:**
- Warning-level ShellCheck: ✗ before implementation with exactly two `SC2034` warnings (`foreign_listener_output`, `port`) → ✓ after implementation with exit 0

**New tests added:**
- No new test case; the warning-level ShellCheck command is the captured regression proof, backed by the existing foreign-listener contract fixture.

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-07T20:33:16Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings:** None
**Minor findings:** 0
**Acceptance:** Pass — warning-level ShellCheck is clean, Bash syntax and the full fixture suite pass, and the refusal behavior remains covered.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Orientation

The contract-regression baseline is warning-clean again, so later suite-split work can treat new ShellCheck warnings as regressions rather than inherited noise.
