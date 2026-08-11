---
id: REQ-166
title: Simplify session-start hook and fix dead fail-soft fallback
status: completed
created_at: 2026-08-11T11:46:50Z
user_request: UR-036
domain: general
prime_files: []
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-165, REQ-167, REQ-168]
batch: stabilization-audit
write_set: [skills/do-work/hooks/session-start.sh]
claimed_at: 2026-08-11T12:23:14Z
route: A
completed_at: 2026-08-11T12:28:07Z
commit: 6538bdd
kb_status: pending
kb_entry:
---

# Simplify Session-Start Hook and Fix Dead Fail-Soft Fallback

## What

Fix the demonstrated bug in `skills/do-work/hooks/session-start.sh` — under `set -euo pipefail`, a failed `grep` in the `VERSION=$(grep … | sed …)` pipeline aborts the script before the `[ -z "$VERSION" ]` fallback can run — and gut the script to its minimal form. The banner's job is two lines of logic (read a version string, count queue files); the current 46-line script's defensive apparatus is what produced the bug.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** First add a fixture-driven behavioral probe and confirm the current hook fails the missing/reformatted version cases; then remove the `errexit`/`pipefail` trap, collapse version/count fallbacks into minimal assignments, preserve the anchored hook comment/command, wire the probe into the aggregate suite, and rerun lint/contracts.
- [x] **[APPLY]:** Added the four-case fixture probe, simplified the hook's runtime logic, and wired the probe through the existing aggregate test seam while preserving the anchored command/comment verbatim.
- [x] **[UNIFY]:** Reviewed the full diff and stat for `skills/do-work/hooks/session-start.sh`, `_dev/tests/session-start-hook-behavior.sh`, and `_dev/tests/contract-regressions.sh`; ran Bash syntax checks, ShellCheck, the four fixtures, the shipped-shell harness, similar-pattern grep, `git diff --check`, and the full aggregate contract suite.

## Why (if provided)

User: the hook feels unnecessarily complex; exemplar of "many things got more complex than needed." The fail-soft intent (banner degrades to "unknown"/0 rather than dying) is right — the implementation defeats it.

## Context

- Bug reproduced this session: with a version file lacking the `**Current version**:` line, the script exits 1 with no output; the stderr warning and `VERSION="unknown"` fallback are dead code for the missing-file and reformatted-line cases. Only the "label present but value empty" case reaches the guard.
- Keep, verbatim: the path-anchoring hook command (`${CLAUDE_PROJECT_DIR:-.}/…`) and its warning comment — it traces to a real regression (relative path resolved from project root, "No such file or directory"). Do not resurrect the relative path.
- Keep the observable contract: single stdout line `do-work v{VERSION} loaded. {N} pending REQ(s). Say 'do-work help' for commands.`; `unknown` / `0` on failure; exit 0 in every failure mode. `hooks/hooks.json` stays byte-compatible with what installers merge.
- Simplification latitude: the whole script should shrink toward the minimal form (e.g. `|| echo unknown`-style defaults inline). Preserve macOS/BSD `wc` padding handling in whatever form (`tr -d ' '` or arithmetic).
- Independent of REQ-165, but once the harness exists this script falls under its lint.

## Red-Green Proof

**RED prompt/case:** Run the hook against a skill root whose `actions/version.md` is absent or lacks the `**Current version**:` line — e.g. the scratchpad repro: `grep -m1 '^\*\*Current version\*\*:' missing.md | sed …` inside `VERSION=$(…)` under `set -euo pipefail`.
**Why RED now:** `pipefail` propagates grep's failure through the pipeline; `set -e` aborts on the failed assignment; the script exits 1 printing nothing — no banner, no stderr diagnostic, despite an explicit fallback written for exactly this case.
**GREEN when:** Same fixture prints `do-work vunknown loaded. 0 pending REQ(s). Say 'do-work help' for commands.` and exits 0; a `_dev/tests/` case covers missing-file, reformatted-line, and missing-queue-dir fixtures; happy path still prints real version and count.
**Validation:** User confirmed (bug demonstrated to the user in-session; "keep the hook slot, gut the script" direction endorsed via this capture).

## Full Context

See `do-work/user-requests/UR-036/input.md` for complete verbatim input.

---
*Source: "do-work capture-request for audit and fix to simplify and make it robust" (UR-036)*

---

## Triage

**Route: A** - Simple

**Reasoning:** The bug has a confirmed reproduction, names the exact runtime script and output contract, and prescribes a small fail-soft simplification plus focused regression fixtures.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation using the captured RED/GREEN proof and `bug-fix` specification.

*Skipped by work action*

## Root Cause

`set -euo pipefail` made the `VERSION=$(grep ... | sed ...)` assignment a fatal command: a missing file or unmatched label propagated grep's non-zero status through the pipeline, and `errexit` terminated the hook before the later empty-value fallback could execute. The warning and `VERSION="unknown"` branch therefore covered only a parsed-but-empty value, not the two failures it claimed to handle.

## Implementation Summary

**Files changed:**
- `skills/do-work/hooks/session-start.sh` (modified)
- `_dev/tests/session-start-hook-behavior.sh` (new)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Reduced the SessionStart hook's runtime logic to fail-soft version and queue-count assignments under `set -u`, preserving the anchored hook command and exact banner contract. Added fixture coverage for happy path, missing version file, reformatted version label, and missing queue directory, then wired it into the aggregate contracts.

## Qualification

Passed — 3 files verified, all runtime/output/test requirements traced, P-A-U confirmed. The hook diff removes defensive surface rather than adding a workaround, the behavioral probe exercises the copied real script at its environment boundary, and the aggregate runner reaches the probe explicitly.

## Testing

**Tests run:**
- `_dev/tests/session-start-hook-behavior.sh`
- `bash -n skills/do-work/hooks/session-start.sh _dev/tests/session-start-hook-behavior.sh _dev/tests/contract-regressions.sh`
- `shellcheck --severity=warning skills/do-work/hooks/session-start.sh _dev/tests/session-start-hook-behavior.sh _dev/tests/contract-regressions.sh`
- `_dev/tests/action-shell-blocks.sh`
- `bash _dev/tests/contract-regressions.sh`
- `git diff --check`

**Result:** ✓ All passing; aggregate contract suite exited 0.

**Red-green validation:**
- Missing `actions/version.md`: ✗ before implementation exited 2 with no banner → ✓ after implementation prints `do-work vunknown loaded. 0 pending REQ(s). Say 'do-work help' for commands.` and exits 0.
- Reformatted version label: ✗ before implementation exited 1 with no banner → ✓ after implementation prints the same `unknown`/0 fallback and exits 0.
- Happy path and missing queue directory both pass with the real version/count and real-version/0 contracts respectively.

**New tests added:**
- `_dev/tests/session-start-hook-behavior.sh` — four fixture-driven behavior cases using the real copied hook

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/contract-regressions.sh` (aggregate seam established by REQ-165) — invokes the REQ-166 hook probe; no prior assertions were weakened or changed

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-11T12:27:44Z

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
**Acceptance:** Pass — the real hook prints the exact banner with the current version/count, and all missing/reformatted version and missing-queue fixtures now fail soft with exit 0.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Writing the fixture probe before touching the hook made the failure mode undeniable: missing and reformatted version inputs exited with two distinct non-zero codes but the same missing-banner symptom.
- Testing a copied real hook inside a synthetic skill/project tree verified path derivation, environment anchoring, version parsing, queue counting, output, and exit status through the actual caller seam.

**What didn't:**
- A later empty-value guard cannot make a command substitution fail-soft when `set -e` has already terminated the script; `pipefail` made the supposedly defensive `grep | sed` pipeline the abort trigger.

**Worth knowing:**
- Keep the `${CLAUDE_PROJECT_DIR:-.}` command anchor and its warning comment intact. The simplification belongs in runtime logic, not in the installation path that previously regressed.
- For tiny status hooks, defaults that naturally produce an empty value are safer than strict-shell defensive branches whose control flow is harder than the output contract.

## Orientation

[MAP CHANGED] The core SessionStart status hook is now a minimal fail-soft banner backed by its own fixture probe in the aggregate contract suite; missing runtime metadata produces `unknown`/0 without losing the startup signal.
