---
id: REQ-166
title: Simplify session-start hook and fix dead fail-soft fallback
status: pending
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
---

# Simplify Session-Start Hook and Fix Dead Fail-Soft Fallback

## What

Fix the demonstrated bug in `skills/do-work/hooks/session-start.sh` — under `set -euo pipefail`, a failed `grep` in the `VERSION=$(grep … | sed …)` pipeline aborts the script before the `[ -z "$VERSION" ]` fallback can run — and gut the script to its minimal form. The banner's job is two lines of logic (read a version string, count queue files); the current 46-line script's defensive apparatus is what produced the bug.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
