---
id: REQ-180
title: Fix contract-regressions.sh Justfile case mismatch aborting late checks
status: completed
status_changed_at: 2026-08-15T07:41:57Z
claimed_at: 2026-08-15T09:26:47Z
completed_at: 2026-08-15T09:40:13Z
commit: 2e5a5c4
created_at: 2026-08-14T10:33:56Z
user_request: UR-040
addendum_to: REQ-179
domain: general
review_generated: true
tdd: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
write_set: [_dev/tests/contract-regressions.sh]
route: A
kb_status: pending
kb_entry:
---

# Fix contract-regressions.sh Justfile Case Mismatch Aborting Late Checks

## What

`_dev/tests/contract-regressions.sh:1797` and `:1804` reference `Justfile` (capital J), but the tracked file is lowercase `justfile`. On a case-sensitive filesystem, `extract_kanban_shutdown_line Justfile` hits `awk: cannot open .../Justfile` and the suite aborts with exit 2 at that point — every check after line ~1797 (roughly 1,500 lines including the Common Rationalizations regrowth ratchet) silently never runs. On case-insensitive filesystems (macOS default) the mismatch is invisible, which is why it survived.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Use the captured Linux RED evidence, replace only the two approved filename literals, then prove the repaired block reads the tracked lowercase file and the suite advances beyond it.
- [x] **[APPLY]:** Changed the two `Justfile` arguments in `_dev/tests/contract-regressions.sh` to `justfile` without adding case-resolution machinery.
- [x] **[UNIFY]:** Reviewed the one-file diff; ran Bash syntax, ShellCheck, exact old-primitive search, focused block comparisons, `git diff --check`, and the full contract suite through the repaired late checks.

## Open Questions

- [x] Fix by lowercasing the two references to `justfile`, or by resolving the file case-insensitively? → Lowercase the two literals to match the tracked filename.
  Recommended: lowercase the two literals to match the tracked filename — smallest fix, matches the repo's actual file.
  Also: a glob/`find -iname` resolution (tolerant but adds machinery), rename the tracked file to `Justfile` (breaks `just` conventions).

## Red-Green Proof
**RED prompt/case:** On a case-sensitive filesystem, `bash _dev/tests/contract-regressions.sh` prints `awk: cannot open ".../Justfile"` and exits 2 at `extract_kanban_shutdown_line`, before the late-suite checks run.
**Why RED now:** Two hardcoded `Justfile` literals vs the tracked lowercase `justfile`; observed at baseline during REQ-179 (2026-08-14, Linux sandbox).
**GREEN when:** The suite runs past line ~1797 on a case-sensitive filesystem (the kanban-shutdown-line check actually executes against `justfile`), and the only remaining sandbox failure is the known environmental process-tree probe.
**Validation:** User chose the recommended lowercase-literal repair via `do-work clarify` on 2026-08-15.

## Full Context

Discovered Tasks of `REQ-179` (see its archive entry). Parent UR: `do-work/user-requests/UR-040/` (archived with UR-040).

---
*Source: REQ-179 build — discovered task*

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names one shell test file, the two incorrect literals, the approved replacement, and an exact RED/GREEN suite outcome.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Root Cause

The contract suite hardcoded a case variant that happens to resolve on default macOS filesystems but does not name the tracked file byte-for-byte. With `set -e`, the failed `awk` open aborted the entire late half of the suite on case-sensitive systems.

## Implementation Summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Replaced the two capitalized `Justfile` inputs with the tracked lowercase `justfile`, keeping the existing shutdown-line comparison and negative assertions unchanged.

## Qualification

Passed — the single declared implementation file is substantive, both approved literal sites changed and no others, the old extract primitive is absent, all requirements trace to the diff, P-A-U is complete, and data-flow checks are not applicable.

## Testing

**Tests run:** `bash -n _dev/tests/contract-regressions.sh`; `shellcheck _dev/tests/contract-regressions.sh`; focused shutdown-line extraction/comparison; `bash _dev/tests/contract-regressions.sh`; `git diff --check`
**Result:** ✓ All passing; the full contract suite exits 0 and reaches its final `Contract regression checks passed.` line

**Red-green validation:**
- Captured case-sensitive baseline from REQ-179: ✗ `awk` could not open `Justfile`, exit 2 before late checks → ✓ both inputs now name tracked `justfile`, the focused comparison exits 0, and the full suite completes
- This macOS checkout is case-insensitive, so the original RED is the durable Linux evidence recorded in the REQ; the byte-exact source correction is independently verified here

**New tests added:** None — this REQ repairs the contract suite's own executable path so its existing late checks run.

*Verified by work action*

## Review

**Overall: 99%** | 2026-08-15T09:39:35Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — both path consumers now use tracked `justfile`, focused checks pass, and the complete contract suite reaches its final PASS at exit 0.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed inline by review-work action after the delegated reviewer timed out*

## Lessons Learned

**What worked:** The captured case-sensitive Linux failure named the exact byte-level mismatch, so the repair stayed at two literals and the full suite proved late checks were reachable again.
**What didn't:** Reproducing RED on a default macOS filesystem was not meaningful because case-insensitive lookup masks the bug; a disk-image workaround added ceremony without improving the captured evidence.
**Worth knowing:** Shell test fixtures and prescribed paths must use the tracked filename's exact casing even when a developer filesystem accepts variants.

## Orientation

The repository contract suite now names the lowercase root `justfile` exactly, so case-sensitive environments execute the Kanban shutdown checks and the rest of the late regression suite.

## Knowledge Handoff

- `kb_status: pending`
- No KB source was written; unattended orchestration defaults to save-for-later without user consent.
