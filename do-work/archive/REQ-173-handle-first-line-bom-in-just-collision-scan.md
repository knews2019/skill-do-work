---
id: REQ-173
title: Handle first-line BOM in Just collision scan
status: completed
claimed_at: 2026-08-11T17:09:11Z
route: A
completed_at: 2026-08-11T17:19:33Z
commit: 8092258
created_at: 2026-08-11T17:00:04Z
user_request: UR-039
domain: testing
prime_files: [skills/do-work/tools/prime-do-work-update.md]
tdd: true
kb_status: pending
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-172, REQ-174]
batch: accepted-p2-fixes
write_set: [tools/replace-text-section.sh, skills/do-work/tools/replace-text-section.sh, _dev/tests/install-suite-behavior.sh]
---

# Handle First-Line BOM in Just Collision Scan

## What

Recognize a reserved Just recipe when the first line begins with a UTF-8 BOM, including when `just` is unavailable, without changing the target file's bytes.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Adapt the existing no-Just installer collision fixture to put exactly one UTF-8 BOM before a first-line `run-kanban:` definition, then run the installer behavior suite to capture the expected successful-invalid-install RED. In the fallback scanner, remove that BOM only from the first line's classification view, mirror the helper byte-for-byte, and rerun the focused/full shell contracts plus ShellCheck and mirror comparison.
- [x] **[APPLY]:** The BOM fixture first failed because the installer accepted and mutated the project. The scanner now removes one `EF BB BF` prefix only from its first-line classification value; the same fixture passes through the no-Just installer path without confirmation or mutation. Both distributed helpers contain the identical implementation.
- [x] **[UNIFY]:** Reviewed `git diff --stat` and the scoped diff: the installer test contains the raw BOM regression and unchanged-state checks, while each helper changes only scanner classification. `install-suite-behavior.sh` and the full contract regression suite pass; Bash syntax checks and warning-level ShellCheck pass for all three declared implementation files; `cmp` confirms the helper copies are byte-identical. No debug artifacts were added.

## Why

Just accepts a BOM-prefixed first definition, but the fallback scanner currently misses it. Without `just`, installation can report success while leaving a duplicate reserved recipe.

## Detailed Requirements

- Strip exactly one leading UTF-8 BOM from the first line's classification view only.
- Preserve all target bytes.
- Keep both distributed helper copies byte-identical.
- Reject the collision before confirmation or client mutation when `just` is unavailable.
- Avoid broader Unicode or encoding normalization.

## Constraints

- Preserve the no-Just collision scanner earned by REQ-152.
- Retain the existing multiline-literal behavior from REQ-162.

## Red-Green Proof

**RED prompt/case:** Install into a project whose marker-free Justfile starts with UTF-8 BOM bytes followed by `run-kanban:`, with `just` absent from `PATH`.
**Why RED now:** The helper misses the definition, installation returns success, and the resulting Justfile contains a duplicate recipe.
**GREEN when:** The installer rejects the reserved collision before confirmation or mutation while preserving the BOM-prefixed Justfile byte-for-byte.
**Validation:** User confirmed

## Full Context

See `do-work/user-requests/UR-039/input.md` and the preceding validated-feedback report.

---
*Source: fix accepted*

---

## Triage

**Route: A** - Simple

**Reasoning:** The collision is reproduced end-to-end, the two mirrored helpers and existing no-Just fixture are named, and the accepted remedy is a first-line classification-only adjustment.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Root Cause

The fallback scanner passed the target's raw first line directly to an ASCII-anchored Just identifier matcher. A leading UTF-8 BOM therefore prevented the recipe name from matching even though Just treats that BOM as a file prefix and accepts the definition. The fix keeps raw target bytes untouched and removes exactly one BOM only from the first line value used for collision classification.

## Implementation Summary

**Files changed:**
- `tools/replace-text-section.sh` (modified)
- `skills/do-work/tools/replace-text-section.sh` (modified)
- `_dev/tests/install-suite-behavior.sh` (modified)

**What was done:** The fallback definition scanner removes one UTF-8 BOM only from its first-line classification value in both distributed helper copies. The no-Just installer fixture now replays a BOM-prefixed reserved recipe and verifies pre-confirmation rejection with byte- and state-preservation.

## Qualification

Passed — 3 implementation files verified, 5 requirements traced, P-A-U confirmed. Both helper copies are byte-identical; the implementation changes only classification data, and the installer fixture proves the original target and surrounding client state remain unchanged.

## Testing

**Tests run:** `bash _dev/tests/install-suite-behavior.sh`; `_dev/tests/contract-regressions.sh`; `bash -n` and warning-level `shellcheck` on all three implementation files; helper `cmp`; `git diff --check`
**Result:** ✓ All passing

**Red-green validation:**
- No-Just BOM collision fixture: ✗ before implementation (`installer accepted a BOM-prefixed reserved recipe when Just was unavailable`) → ✓ after (`suite installer behavior probes passed.`)

**New tests added:**
- BOM-prefixed first-line reserved recipe rejection before confirmation or mutation.
- Existing byte, mode, settings, module, and Git-state preservation assertions now cover the BOM replay.

*Verified by work action*

## Review

**Overall: 99%** | 2026-08-11T17:19:05Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — the no-Just installer rejects a one-BOM, first-line reserved recipe before confirmation or mutation; both helper mirrors and all multiline scanner regressions pass.
**Suggested testing:** 1 item — a future direct helper boundary fixture could pin that a second BOM or later-line BOM is not normalized.
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Replaying the no-Just installer path preserved the real pre-confirmation and byte-identity boundaries while isolating the scanner defect.
**What didn't:** An ASCII-anchored identifier matcher silently assumed the first physical byte belonged to the Just grammar.
**Worth knowing:** UTF-8 BOM handling belongs only in the first-line classification view; the byte-preserving target and all later lines stay untouched.

## Orientation

The suite installer now rejects BOM-prefixed reserved Just recipes even without `just`, while preserving the byte-oriented managed-section transaction.
