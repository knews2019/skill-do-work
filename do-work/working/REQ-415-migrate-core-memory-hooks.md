---
id: REQ-415
title: 'Migrate the core SessionStart and memory hooks into Go subcommands'
status: claimed
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-414]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
claimed_at: 2026-09-01T02:34:48Z
---

# Migrate the Core SessionStart and Memory Hooks into Go Subcommands

## What
Replace hook domain logic with canonical `do-work-cli` hook subcommands.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Migrate the core SessionStart hook plus memory SessionStart and Stop-capture hooks into Go subcommands.
- Preserve exact hook stdin/stdout protocols, event behavior, redaction, deduplication, timestamp/reservation repair, and failure handling.
- Keep each existing hook `.sh` path as a thin build-and-exec shim.
- Add fixture tests for malformed input, secrets, duplicate content, repair cases, and exact output.

## Constraints
- Hook launchers must remain safe in installed projects and must stop actionably when the canonical binary cannot build or run.

## Dependencies
Depends on REQ-414 (remaining core utility primitives used by SessionStart).

## Builder Guidance
Certainty level: Firm. Freeze byte-level hook protocols before moving any logic.

## Red-Green Proof
**RED prompt/case:** Replay captured valid, malformed, redacted, duplicate, and repair hook events against missing Go hook subcommands.
**Why RED now:** Hook behavior currently resides in three shipped shell implementations.
**GREEN when:** Go subcommands produce equivalent status/stdout/stderr/effects for every fixture and the original hook paths are launch-only shims.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
