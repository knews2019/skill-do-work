---
id: REQ-451
title: 'Make confirmation input interruptible'
status: completed
route: C
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
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-09-01T23:50:40Z
  basis:
    - Route C
    - 5-file write set
    - 2 subsystems involved
    - 4 acceptance criteria
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
required_lessons: [skills/do-work/tools/do-work-cli/lessons-do-work-cli.md]
related: [REQ-450, REQ-452, REQ-453, REQ-454, REQ-455, REQ-456, REQ-457]
batch: accepted-validate-feedback-root-causes
claimed_at: 2026-09-01T23:50:34Z
planning_at: 2026-09-01T23:54:23Z
dispatch_at: 2026-09-01T23:58:26Z
builder_handback_at: 2026-09-02T00:09:20Z
integration_at: 2026-09-02T00:09:20Z
review_at: 2026-09-02T00:54:06Z
write_set:
  - skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go
  - skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go
  - skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction_test.go
  - skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go
  - _dev/tests/install-suite-behavior.sh
completed_at: 2026-09-02T00:55:13Z
release_at: 2026-09-02T00:56:47Z
commit: 21036776
kb_status: promoted
kb_entry: REQ-451-make-confirmation-input-interruptible.md
---

# Make Confirmation Input Interruptible

## What

Make install and update confirmation input cancellation-aware so `SIGINT`, `SIGHUP`, and `SIGTERM` cannot leave the process waiting forever at the prompt. Before writes begin, signal handling must exit with the documented signal status without waiting on recovery that the blocked input path itself prevents.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this blocking confirmation-input shutdown root cause.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Make confirmation selection context-aware without changing the write-started recovery owner; prove install/update prompt signals with real subprocesses and mutation-time recovery in the shell lane.
- [x] **[APPLY]:** Added buffered cancellation-aware confirmation reads, blocked-reader install/update tests, a direct-binary HUP/INT/TERM matrix, and exact mutation-time exit-130 coverage in five scoped files.
- [x] **[UNIFY]:** Reviewed all five files and verified context race checks, preserved yes/EOF behavior, no-input signal proof, direct PID signaling, process/goroutine cleanup, unchanged recovery boundaries, and no debug artifacts.

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

## Triage

**Route: C** - Complex

**Reasoning:** The fix crosses blocking input, signal cancellation, process exit codes, and recovery sequencing. It needs a precise pre-write/mutation-time boundary plus real subprocess regressions for install and update.

**Planning:** Required

## Plan

1. Make the one-line confirmation read context-aware by performing the existing blocking read in a single goroutine, sending its result through a buffered channel, and selecting against `ctx.Done()`. Re-check cancellation before accepting an affirmative line so cancellation wins close races.
2. Thread cancellation through `reviewAndConfirm` into the existing cleanup path without changing `armSignalRecovery` or the `writeStarted` boundary. Before writes, cancellation returns and closes recovery completion immediately; after writes begin, the existing single recovery owner still completes rollback before exit.
3. Cover install and delegated update confirmation with blocked-reader seam tests that cancel without sending data, plus real built-CLI HUP/INT/TERM subprocess matrices that keep stdin open, signal the binary directly, require exit 130, and prove no fresh install or update mutation landed.
4. Strengthen the existing installer behavior script's mutation-time TERM lane to assert exact exit 130 while retaining its full recovery checks.

**Files:** `internal/suiteinstall/install_transaction.go`, `install_transaction_test.go`, `update_transaction_test.go`, `suite_commands_test.go`, and `_dev/tests/install-suite-behavior.sh`.

**Architectural decision:** The installer intentionally preserves its documented uniform signal exit status 130 for HUP, INT, and TERM. A blocked generic `io.Reader` cannot be forcibly stopped safely; the buffered result channel prevents a late completion from blocking, while the outer confirmation call becomes cancellation-aware and the CLI process can exit without another input byte.

*Generated by Plan agent*

## Exploration

`RunInstall` already passes its work context through `reviewAndConfirm` before `writeAndVerify`; `writeStarted` is set only after backups and before the first mutation. The existing signal coordinator groups HUP/INT/TERM under exit 130 and waits for `recoveryFinished`, so making confirmation return on context cancellation unblocks that coordinator without weakening post-write recovery. Update delegates to `RunInstall` with the same context and confirmation reader. Existing suite fixtures, archive helpers, and mutation-time TERM shell coverage provide the required seams.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go`
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go`
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction_test.go`
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go`
- `_dev/tests/install-suite-behavior.sh`

**Acceptance criteria:** Confirmation cancellation returns without another input byte; the real install and update commands exit 130 for HUP/INT/TERM while blocked at the prompt; no project mutation lands before confirmation; mutation-time TERM still waits for and proves full recovery; generic-reader handling does not block a late sender or change existing EOF/yes semantics.

## Root Cause

Signal handling cancelled the install context and waited for `recoveryFinished`, but confirmation blocked synchronously in `ReadString` without observing that context. Before writes, the main goroutine could never unwind to close the channel the signal goroutine awaited, leaving the process hung at the prompt.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (modified)
- `_dev/tests/install-suite-behavior.sh` (modified)

**What was done:** Confirmation input now reads into a buffered result channel and selects against cancellation, checking cancellation both before launch and before accepting yes. The existing write-started recovery machinery remains unchanged; real install/update subprocesses now prove HUP/INT/TERM exit 130 without input, and mutation-time TERM still proves full recovery.

## Qualification

Passed — five declared files verified against all requirements. Pre-write cancellation can unwind without rollback; after mutation begins the original recovery owner still completes before signal exit. Tests signal the built CLI directly, hold stdin open, send no byte, and reap every process/reader path.

## Testing

**Tests run:** controlled blocking mutation RED; focused signal test; repeated blocked-reader/signal runs; `go test -race -count=1 ./internal/suiteinstall`; `go vet ./...`; `bash _dev/tests/install-suite-behavior.sh`; `git diff --check`; full uncached module; canonical maintainer verification
**Result:** ✓ All checks pass. The two unrelated UTC-date-boundary failures were repaired separately in `f7a99afe`; focused package tests, race detection, vet, installer behavior, the full CLI module, and canonical maintainer verification are now green with the REQ-451 implementation present.

**Red-green validation:** With cancellation ignored under the new signature, the real install/INT case exceeded its five-second deadline while stdin stayed open. Restoring the context select made the same test pass in 0.886s; ten repeated signal/blocked-reader runs passed.

**New tests added:** blocked-reader install/update cancellation and direct-binary install/update × HUP/INT/TERM exit-130/non-mutation matrix.

*Verified by work action*

## Review

**Overall: 99%** | 2026-09-02T00:54:06Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 99% |
| Test Adequacy | 100% |
| Scope | 100% |
| UX | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None
**Minor findings:** 0
**Acceptance:** Pass — prompt-time signals exit without input and mutation-time recovery remains intact.
**Follow-ups created:** None.

**Residual risk:** A generic blocked `io.Reader` goroutine cannot be forcibly cancelled and may remain until its reader returns. The buffered result channel prevents a late sender from blocking, and real CLI processes exit after recovery.

*Reviewed by work action*

## Orientation

Install and update confirmation are now interruptible without weakening transactional recovery after writes begin. The do-work CLI prime remains current.

## Gate Hold Resolution

The unrelated UTC-sensitive knowledge-command test failures were repaired in standalone commit `f7a99afe`. With that baseline restored, REQ-451 passed its focused interruptibility test, `go test -race -count=1 ./internal/suiteinstall`, `go vet ./...`, `_dev/tests/install-suite-behavior.sh`, the full uncached CLI module suite, `git diff --check`, scope drift, and canonical maintainer verification.

## Lessons Learned

**What worked:** Separating pre-write confirmation cancellation from the existing write-started recovery owner kept the fix small and made real-process signal tests decisive.

**What did not work:** The first canonical verification was held by date-bound fixtures outside this request. Repairing those tests in a standalone commit restored the gate without widening REQ-451.

**Worth knowing:** A buffered result channel is necessary because a generic `io.Reader` may complete after the caller returns; post-write recovery remains owned by the unchanged transaction boundary.
