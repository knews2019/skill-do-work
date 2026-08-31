---
id: REQ-451
title: 'Make confirmation input interruptible'
status: pending
created_at: 2026-08-31T20:49:21Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-450, REQ-452, REQ-453, REQ-454, REQ-455, REQ-456, REQ-457]
batch: accepted-validate-feedback-root-causes
---

# Make Confirmation Input Interruptible

## What

Make install and update confirmation input cancellation-aware so `SIGINT`, `SIGHUP`, and `SIGTERM` cannot leave the process waiting forever at the prompt. Before writes begin, signal handling must exit with the documented signal status without waiting on recovery that the blocked input path itself prevents.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this blocking confirmation-input shutdown root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding #3 — P2 — source:** `internal/suiteinstall/install_transaction.go:775`

> ````text
> [P2] Make confirmation input interruptible — [prj].claude/skills/do-work/tools/do-
> work-cli/internal/suiteinstall/install_transaction.go:775-775
> When install or update is waiting at the confirmation prompt, pressing Ctrl-C or receiving HUP/TERM cancels only the work context
> and then waits for recoveryFinished, while this blocking ReadString does not observe that context. RunInstall consequently cannot
> return to close the channel, so the process hangs instead of exiting with the documented signal status; make the confirmation read
> cancellation-aware or avoid waiting on recovery before writes begin.
> ````

- **Evidence:** `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go:157` and `201-210` cancel the context and wait for `recoveryFinished`; lines `729-730` and `770-775` reach a blocking `ReadString` that does not observe the context. A process reproduced at an open FIFO remained alive after `SIGINT` and required `SIGKILL`.
- **Surface-cost result:** Earned — the hang is concrete and reproduced. A cancellation-aware confirmation boundary plus a subprocess signal regression is cheaper than an uninterruptible installer.

## Detailed Requirements

- Make confirmation reads observe cancellation or restructure the pre-write signal path so it does not wait on blocked input.
- Preserve the documented signal exit status, including exit 130 for `SIGINT`.
- Preserve recovery guarantees after mutation begins.
- Cover install and update confirmation behavior.

## Constraints

- Do not weaken rollback or recovery for signals received after writes begin.
- Do not rely on sending another input byte to release the blocked read.

## Dependencies

No request prerequisite.

## Builder Guidance

Certainty level: Firm. Distinguish the pre-write confirmation boundary from mutation-time recovery.

## Red-Green Proof

**RED prompt/case:** Start the real command as a subprocess with confirmation input held open, wait until the prompt, send `SIGINT`, and impose a short exit deadline.
**Why RED now:** The blocking read ignores context cancellation while signal handling waits for a completion channel that `RunInstall` cannot close.
**GREEN when:** The subprocess exits with status 130 within the deadline without extra input or forced termination; equivalent HUP/TERM coverage preserves their documented statuses.
**Validation:** User confirmed after validate-feedback accepted Finding #3.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding #3, captured by UR-085.*
