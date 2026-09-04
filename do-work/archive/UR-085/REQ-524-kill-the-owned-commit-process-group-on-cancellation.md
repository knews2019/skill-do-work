---
id: REQ-524
title: 'Kill the owned commit process group on cancellation'
status: cancelled
created_at: 2026-09-02T23:58:00Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
write_set:
  - skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go
  - skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process.go
related: [REQ-457]
completed_at: 2026-09-03T11:47:26Z
---

# Kill the Owned Commit Process Group on Cancellation

## What

Make cancellation of a media Git transaction terminate the whole owned process group, not just the direct `git` child. A `pre-commit` hook that ignores `SIGTERM` currently outlives the cancelled transaction and keeps running after the command has returned.

## Instances

- `internal/toolboxcommands` → `TestRemediationCancellationReachesMediaGitCommitAndRollback` fails with `media commit hook survived cancellation` (`report_image_process_test.go:85`). The transaction rolls back correctly and returns within the deadline; the orphaned hook process is the only unmet assertion.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Discovered by:** REQ-457's Step 5.75 pre-flight, as one of two pre-existing failures in the `bash _dev/tests/maintainer-verify.sh` baseline at `008f3d3`. Unrelated to REQ-457's created-path ownership invariant, so it is captured here rather than folded into that REQ.
- **Evidence:** the test writes a `pre-commit` hook that traps and ignores `TERM` then sleeps 30s, cancels the context once the hook is running, and asserts `syscall.Kill(pid, 0)` fails afterwards.

## Diagnosis

Measured by process-table sampling at 20 ms during a live failing run, not inferred. Three plausible causes are ruled out: `Setpgid` **is** applied (`configureCancellableProcessGroup`, `internal/gittransaction/git_transaction.go:1344`), the negative-PID group signal **is** delivered (`pidfd_open(-pid)` returns `EINVAL`, so Go falls back to `syscall.Kill(-pid, sig)`), and the hook's `trap '' TERM` sets `SIG_IGN`, which is inherited across fork and exec, so `sleep` is TERM-proof too.

The real gap is ordering plus reaping:

- **The escalation is detached and nothing waits for it.** `Cancel` returns straight after the group SIGTERM; the SIGKILL lives in an orphan goroutine (`git_transaction.go:1362`) with a 1 s timer. `runGit` has no path that proves the group is dead before returning. `runOwnedProcess` (`internal/toolboxcommands/report_image_process.go:62-64`) already has exactly that invariant — `runGit` is the one that lacks it.
- **Group-SIGTERM kills the leader first, orphaning the hook; the later SIGKILL then leaves an unreaped zombie.** `git` dies immediately, the hook is reparented to PID 1, and `Cmd.Wait` stays blocked because the hook still holds the inherited stderr pipe. Measured: the transaction returns ~20 ms after the SIGKILL, and `syscall.Kill(pid, 0)` **succeeds on a zombie**. Orphan reap latency in this environment is ~1.5-2.0 s against the test's 3 s budget.

**This is why merely making the escalation synchronous is not sufficient.** A runner that TERMs the group, waits out grace, SIGKILLs the group, then polls `ownedProcessGroupAlive` — which deliberately treats `Z` as not-alive — would return with the hook a zombie and the assertion would still read alive.

Second launch path in the same file: `indexIsEmpty` (`git_transaction.go:860`) uses a bare `exec.CommandContext` with no process group, bypassing the seam the REQ constrains to one.

## Required Shape

1. Signal **descendants only** first (every group member whose pid is not the leader), so the leader survives to reap them. `git` `waitpid()`s its own hook, which removes the orphan and zombie window entirely — and `git commit` then exits non-zero on its own, which the existing `FailureCommit` path already turns into a rollback.
2. Escalate to SIGKILL on the descendants after a grace window.
3. Only then, if the leader is still alive, signal and escalate on the leader.
4. **Block until the group is gone before returning.** No detached goroutine.
5. Keep `command.WaitDelay` (currently 2 s, `git_transaction.go:1352`) as a last-resort backstop only. The grace budget must stay well under it, or it must be raised in tandem, or `WaitDelay` will fire first and re-create a return that outruns the kill.

## Traps

- **Windows asymmetry cuts the opposite way from `toolboxcommands`.** `configureOwnedProcess` (`report_image_process_windows.go:12`) fails closed. `runGit` must **not** — every Git transaction in the module goes through it, and `configureCancellableProcessGroup` deliberately probes `Setpgid` by reflection and degrades to default cancellation where it is absent. One shared runner needs both behaviours behind one API. `internal/gittransaction` has no build-tagged files today and cross-compiles clean for `GOOS=windows`; introducing a split there is new surface and needs a matching `GOOS=windows` compile line in the prime's Verify section.
- **Three divergent group runners already exist** with three grace policies: this inline `Cancel` closure (1 s), `report_image_process.go` (1 s), and `internal/nextselection/blocked_probe_unix.go` (500 ms + 500 ms). The last one also forwards the received signal and returns `128+signal` statuses — a contract a shared runner must not flatten. Consolidating the seam and sweeping the per-caller status and fail-closed contracts belong in the same change.
- **`prctl(PR_SET_CHILD_SUBREAPER)` is not available.** It lives in `x/sys/unix`, and the prime requires standard-library-only dependencies. The Go floor is 1.25.0, enforced by `_dev/tests/do-work-cli-go125-compatibility.sh`. This is a second reason the descendants-first shape is the right one.
- `reportImageGracePeriod` (`report_image_process.go:16`) is a test seam reassigned by `report_image_process_test.go:22` and `:104`; a refactor must preserve it.

**Regression net:** `report_image_process_test.go:19,92,101` (three currently-passing owned-process tests), `git_transaction_test.go:928-1001,1146` (three `pre-commit` hook tests, all exercising a *failing* hook rather than a cancelled one), `last30days_test.go:135`, `suiteinstall/{install,update}_transaction_test.go`.

## Detailed Requirements

- Cancelling a Git-committing transaction must terminate every process the transaction launched, including grandchildren a hook spawns, not only the direct child.
- Escalate past a process that ignores the graceful signal; a hook trapping `SIGTERM` must not survive.
- Preserve the existing return-within-deadline behavior and the existing rollback-to-preimage behavior on cancellation.
- Do not leave a killed process group's exit status misreported as success.

## Constraints

- Keep the owned-process-group runner the single seam for launched subprocesses; do not add a second ad-hoc kill path.
- Preserve exact typed result and rollback contracts in `prime-do-work-cli.md`.

## Dependencies

No request prerequisite.

## Red-Green Proof

**RED prompt/case:** `go test ./internal/toolboxcommands/ -run TestRemediationCancellationReachesMediaGitCommitAndRollback` — currently fails with `media commit hook survived cancellation`.
**Why RED now:** cancellation signals only the direct `git` child, so a hook that ignores `SIGTERM` keeps running as an orphan.
**GREEN when:** that test passes without modification, the transaction still returns inside its deadline, and the target still rolls back to its preimage bytes.

---
*Source: REQ-457 pre-flight baseline, captured during the work run.*

## Cancelled

- **When:** 2026-09-03T11:47:26Z
- **Why:** folded into REQ-526 section Folded From REQ-497 and REQ-524 as an acceptance criterion (maintainer decision, 2026-09-03 roadmap triage)
- **Decided by:** user, via `do-work abandon`
