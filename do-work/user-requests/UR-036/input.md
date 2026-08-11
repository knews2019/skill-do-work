---
id: UR-036
title: Audit and fix to simplify and make it robust
created_at: 2026-08-11T11:46:50Z
requests: [REQ-165, REQ-166, REQ-167, REQ-168]
word_count: 12
---

# Audit and fix to simplify and make it robust

## Summary

Stabilization pass on the suite itself. The user's goal (stated in conversation): stop reviews from finding 3–5 defects every pass. The conversation diagnosed three recurring defect generators — untested prose-prescribed shell, copy-pasted shell primitives, and untested defensive code that itself breaks — and converged on a plan: build a lint harness for fenced shell blocks first, then deduplicate primitives, then run a delete-or-test audit of defensive code, fixing the session-start hook's demonstrated pipefail bug along the way.

## Extracted Requests

| REQ | Title | Depends on |
|---|---|---|
| REQ-165 | Shell-block lint harness for shipped action files | — |
| REQ-166 | Simplify session-start hook and fix dead fail-soft fallback | — |
| REQ-167 | Deduplicate copy-pasted shell primitives across action files | REQ-165 |
| REQ-168 | Delete-or-test audit of defensive code in shipped skills | REQ-165 |

## Batch Constraints

- **Goal is convergence, not polish:** every fix should close a defect *class* (via a ratchet test) or shrink surface (via deletion). A fix that does neither is out of scope.
- Review findings become ratchet tests — a finding isn't closed until a test would catch its recurrence (existing pattern: `_dev/tests/contract-regressions.sh`).
- The harness (REQ-165) is repo-only under `_dev/tests/` — it does not ship, so the "design for the floor" / no-toolchain rule for shipped packages is not affected.
- Shipped action-file edits must preserve agent compatibility (generalized language, standalone-prompt property) and must not cite `CLAUDE.md`/`AGENTS.md`.
- Serial-only files (`CHANGELOG.md`, `actions/version.md`) are integrator-owned; per-REQ builders must not touch them.
- Prior incident context: CLAUDE.md's "Prescribed Shell Commands" trap list and "Closed Enumerations Go Stale" section are the historical record of the defect classes this batch closes.

## Full Verbatim Input

do-work capture-request for audit and fix to simplify and make it robust

### Conversation context (same session, condensed)

- User goal, stated verbatim earlier in session: "my goal is to stabilize the skill, so on review I'm not getting again 3-5 thing to fix" and "I feel this is just making it unecesarily complex, in fact many things got more complex than needed".
- Demonstrated bug: `skills/do-work/hooks/session-start.sh` runs under `set -euo pipefail`; when `actions/version.md` is missing or its `**Current version**:` line is reformatted, the `grep | sed` pipeline fails, the `VERSION=$(...)` assignment aborts the script before the `[ -z "$VERSION" ]` fallback runs — silent exit 1, no banner, the "unknown" fail-soft path is dead code for the two most likely failure modes. Reproduced in scratchpad this session.
- Agreed direction: keep the hook slot (structural cross-session signal), gut the script to minimal form; build the shell-block lint harness first (biggest generator), then dedup primitives, then audit defensive code with a delete-or-test rule.

---
*Captured: 2026-08-11T11:46:50Z*
