---
id: REQ-407
title: 'Migrate bootstrap, install, update, reconciliation, validation, and fetching into Go'
status: pending
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, skills/do-work/tools/prime-do-work-update.md]
tdd: true
suggested_spec:
depends_on: [REQ-406]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
---

# Migrate Bootstrap, Install, Update, Reconciliation, Validation, and Fetching into Go

## What
Move installation and update domain logic into `do-work-cli` and remove Python/jq implementation branches.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Migrate bootstrap/install/update, byte-safe managed-section replacement, settings reconciliation, suite validation, and archive fetching to Go.
- Preserve fresh/existing-project behavior, exact rollback, custom hooks, reserved recipe collision handling, and cancellation.
- Handle CRLF, BOM, NUL, symlinks, file modes, malformed markers/JSON, and existing Just/CLAUDE/settings content.
- Eliminate do-work Python and jq branches and document the Go 1.26.1+ prerequisite for installation, update, and runtime.
- Keep Python checks only when probing a Python target capability.

## Constraints
- Preserve the public scripts-and-Just installation shape through compatibility launchers.
- Installer/update tests must pass with Python and jq absent.

## Dependencies
Depends on REQ-406 (shared CLI and transaction foundation).

## Builder Guidance
Certainty level: Firm. Characterize current byte and filesystem behavior before replacing implementation branches.

## Red-Green Proof
**RED prompt/case:** Run installer/update fixtures with Python and jq removed from PATH, including an existing managed Justfile and custom settings hooks.
**Why RED now:** Existing reconciliation and managed-section branches conditionally depend on embedded Python or jq.
**GREEN when:** All fresh/existing fixtures succeed with Go alone, preserve exact unrelated bytes/state, and roll back failures according to the shared transaction contract.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*
