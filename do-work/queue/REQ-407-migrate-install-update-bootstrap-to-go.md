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

## Folded From REQ-406 (2026-08-30)

REQ-406 built the foundation and stopped at the command seam, so these became
testable only once a real command is registered — which is this REQ. Folded here
per the Fold-First Rule rather than minted as separate REQs.

- **The `<command>` placeholder is not a runnable argv.** `usageFinding` and the
  `invalid_options` template both emit `do-work-cli --format text <command>`, which
  shows the shape but cannot be pasted. Requirement 5 asks every finding for the
  *exact* next argv. Once this REQ registers commands the runtime knows the name and
  can thread it through.
- **No test observes a successful `--commit` transaction, and `commit_failed` has no
  behavioural test.** REQ-406's fixtures cover exit 0/1/3/4 but the committing success
  path and the commit-failure kind are asserted only through their finding templates.
- **`RollbackResult.Status` has a fourth wire value, the empty string,** for results
  that never ran a Git transaction. A consumer switching on it must handle `""`
  alongside the three constants. REQ-406's D-04 explains why normalising it to
  `not_needed` is not free: every read-only command would print a rollback line
  implying a mutation was possible.
- **The success path does not consult `state.existed`.** At
  `git_transaction.go:161-166` an unrecorded change to a declared target is detected,
  but a file created without `RecordCreated` still reports `succeeded`. Harmless while
  no command is registered; worth closing when one is.
- **Text rendering of changes, skipped work and rollback errors has no direct
  assertion.** The parity test covers findings; these three sections render unasserted.

