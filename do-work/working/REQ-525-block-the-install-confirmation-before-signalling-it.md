---
id: REQ-525
title: 'Block the install confirmation before signalling it'
status: claimed
created_at: 2026-09-03T00:35:00Z
user_request: UR-085
domain: testing
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: false
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-457, REQ-524]
claimed_at: 2026-09-03T11:44:30Z
route: B
write_set:
  - skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-09-03T11:45:48Z
  basis:
    - Route B
    - 2-file write set
    - 1 subsystem involved
    - 4 acceptance criteria
    - async lifecycle behavior
---

# Block the Install Confirmation Before Signalling It

## What

Make `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation` wait until the installer is actually blocked on its confirmation prompt before delivering the signal. Today it signals on a timer, so under parallel load the installer can finish rendering its managed-install diff and exit before the signal is handled, and the assertion reads `exit = <nil>, want 130`.

## Instances

- `internal/suiteinstall` → `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation/install-suite/INT` and `/install-suite/HUP` (`suite_commands_test.go:292`). Observed failing twice under `go test ./...` parallel load during the REQ-457 run, and passing 5/5 with `-count=5` in isolation immediately afterwards. It did not fail in the pre-REQ-457 baseline, and REQ-457's diff touches no file in this package.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Diagnosis pass first (read-only, no edits), which refuted the captured premise and my own hypothesis and confirmed a product defect. Implementation pass then took PF-A's property — one exit owner on the main goroutine, deferred after recovery — rather than PF-A's literal `ExitCodeOverride` plumbing, which is unreachable inside the write set.
- [x] **[APPLY]:** Two files, both inside the widened Scope. `commandruntime`, `resultmodel` and `main.go` were left untouched; the builder stopped and reported rather than widening again.
- [x] **[UNIFY]:** Orchestrator independently re-ran `go build ./...`, `go vet ./...` (clean), `gofmt -l .` (silent), `go test -count=1 ./internal/suiteinstall/` and `go test -race -count=1 ./internal/suiteinstall/` (both green — no data race on the new atomic flag), and read the whole `install_transaction.go` diff. No debug artifacts; every added comment states an invariant, including why the record is stored before `cancelWork`.

## Finding Provenance

- **Discovered by:** REQ-457's Step 6.5 gate attribution, as the one gate failure that was neither in the recorded red baseline nor in a package this REQ's diff touches.
- **Evidence:** the captured stderr shows the installer still emitting `--- managed destination: .claude/skills/do-work ---` diff output when the signal arrives, i.e. it had not yet reached the blocked confirmation read the test name asserts against.

## Detailed Requirements

- Deliver the signal only after the installer has demonstrably reached its blocked confirmation read, rather than after a fixed delay or an unsynchronized amount of output.
- Keep asserting exit 130 and the existing managed-path non-effects; this is a synchronization fix, not a weaker assertion.
- Do not lengthen a sleep to paper over the race, and do not mark the test as skipped or flaky-tolerant.
- The fix must hold under `go test ./...` parallel load, not only under `-run` in isolation.

## Constraints

- Test-side only unless the product genuinely cannot signal readiness; if a product change is required, say so rather than widening silently.

## Dependencies

No request prerequisite.

## Red-Green Proof

**RED prompt/case:** `go test ./... ` for the module under parallel load, repeated; the subtest fails intermittently with `signal interrupt exit = <nil>, want 130`.
**Why RED now:** the test races the installer's own diff rendering instead of synchronizing on the confirmation prompt.
**GREEN when:** the subtest passes across repeated full-module parallel runs with the exit-130 and non-effect assertions unchanged.

---
*Source: REQ-457 gate attribution, captured during the work run.*

---

## Triage

**Route: B** - Medium

**Reasoning:** Captured as mechanical on the assumption that the test signalled on a timer. It does not — it polls for the prompt text — so the actual race is still unlocated and needs discovery before any edit.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Two things the capture assumed that are wrong, established before dispatch:

1. **The test does not signal on a fixed delay.** `runBuiltCLIAtBlockedConfirmation` (`suite_commands_test.go:322`) starts the built CLI, then calls `waitForPromptInFile` (`:374`), which polls the stderr file every 10 ms for up to 10 s and returns as soon as the prompt text appears. Only then does it signal. So "deliver the signal after a fixed delay" is not the defect.
2. **The prompt text is not install-only.** `grep` for `[y/N]` across `internal/suiteinstall` finds exactly one prompt, `install_transaction.go:747` — `update-suite` reaches the same one, because updating runs a full-suite install. So the failing `update-suite/*` subtests are not waiting on a string their command never prints.

What remains is the seam at `install_transaction.go:747-748`: the prompt is written with `narrateWithoutNewline`, and only then does `readConfirmation(ctx)` run. The test's poll observes the *written bytes*, which is strictly earlier than the process being blocked in a cancellable read. Anything the process does between those two points is a window in which the signal can land somewhere that does not produce exit 130 — and the observed failure is exactly `exit = <nil>`, a clean exit rather than a missed signal.

That window is the hypothesis to confirm or refute, not a conclusion. Four sightings this session, on three distinct subtests (`install-suite/INT`, `update-suite/INT`, `update-suite/TERM`, `update-suite/HUP`), always green in isolation and only failing under full-suite parallel load — which fits a window that widens when the machine is contended.

*Generated in-session (single-pass discovery)*

## Confirmed Root Cause — Scope Widened to the Product

The exploration hypothesis above is **refuted**, and so is this REQ's original title and premise. The defect is not the moment the test observes, and it is not test-side at all.

**Two goroutines race to call `os.Exit`, and nothing orders them.**

1. `install_transaction.go:200` arms signals for the whole run, long before the prompt. Its handler goroutine does `cancelWork()` → `<-transaction.recoveryFinished` → `os.Exit(130)` at `:207`.
2. `install_transaction.go:809` — `readConfirmation`'s `case <-ctx.Done():` fires and returns `false`. Cancellation reaches the blocked read correctly, so REQ-451's `[family: interruptible-blocking-io]` contract is intact.
3. `install_transaction.go:158-168` — **`!confirmed` is indistinguishable from a typed "N".** Main narrates "Installation cancelled; no files were changed." and returns `InstallResult{Outcome: OutcomeSuccess, Cancelled: true}`.
4. Deferred, LIFO: `cleanup()` → `close(recoveryFinished)` → `stopSignals()`. The close releases the handler goroutine, now committed to `os.Exit(130)`.
5. Main meanwhile marshals JSON and reaches `cmd/do-work-cli/main.go:63`, `os.Exit(ExitCode(OutcomeSuccess))` = **0**.

Whichever goroutine reaches `os.Exit` first wins. The handler usually does, because main still has marshalling and a write ahead of it; under CPU contention main wins and the process exits 0 — which the test reads as `exit = <nil>`.

**The evidence is in every failure's own captured stderr:** `Install this complete four-skill suite? [y/N] Installation cancelled; no files were changed.` Main took the declined branch and won the race.

**Deterministic reproduction:** running the built child under `GOMAXPROCS=1` fails **6 of 6** subtests every run, with or without a settle delay after the prompt — starving the handler goroutine of a scheduler slot makes main win every time. A 500 ms delay after the prompt bytes appear, with the read unquestionably parked in its `select`, still fails: that is what refutes the write-then-read window hypothesis.

**Baseline, 7 full-module runs on the unmodified tree:** this test failed in 3 of 7.

**This is user-visible.** A person pressing Ctrl-C at the install prompt gets exit 0 and a JSON result claiming `outcome: success`. No partial install is possible — `writeStarted` is false and `readConfirmation:815` re-checks `ctx.Err()` — so the filesystem is safe and the *reported status* is the defect.

**Scope decision (D-01, recorded below):** this REQ's constraint said "test-side only unless the product genuinely cannot signal readiness; if a product change is required, say so rather than widening silently." The builder said so, and the widening is deliberate rather than silent. The user-facing symptom is one thing, so it stays one REQ; the write set gains the production file.

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (modify) — make the interrupted install have exactly one exit owner, so the signal's status cannot lose a race to the ordinary return path
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (modify) — add the mid-write signal case that has no coverage, and make the existing signal test deterministic

**Files I will NOT touch:** the recovery-ordering property at `install_transaction.go:190-197` is preserved, not rewritten; and no assertion in the existing signal test is weakened — the exit-130 and managed-path non-effect checks stand exactly as written.

*Scope widened from test-only after the root cause was confirmed to be in the product — see* **Confirmed Root Cause** *above and D-01.*

**Acceptance criteria (restated from REQ):**
- [~] ~~The signal is delivered only after the installer has demonstrably reached its blocked confirmation read~~ — **premise refuted.** A 500 ms delay past that point still fails; the read is parked in its `select` and cancellation reaches it correctly. Replaced by: an interrupted confirmation reports exit 130 regardless of goroutine scheduling.
- [x] Exit 130 and the existing managed-path non-effect assertions are unchanged — this is a synchronization fix, not a weaker assertion
- [x] No sleep is lengthened to paper over the race, and the test is neither skipped nor made flaky-tolerant
- [x] The fix holds under `go test ./...` parallel load, not only under `-run` in isolation

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (modified)

**What was done:** The `recoveryFinished` channel handshake is replaced by an `interruptWasObserved atomic.Bool`. The signal handler now only records the interruption and cancels the work context — its `os.Exit(130)` and its channel wait are gone. A new `exitIfInterrupted`, deferred immediately *before* `cleanup()` so it runs immediately *after* it, is the transaction's single owner of the interrupted exit status. Recovery-before-exit is therefore defer order on one goroutine rather than a race two goroutines have to win in the right order. On the test side both built-child helpers are pinned to `GOMAXPROCS=1`, and the previously uncovered mid-write interrupt case (F4) is added.

## Decisions

- **D-01** — Scope widened from test-only to the product after the root cause was confirmed. Recorded in **Confirmed Root Cause** above; the REQ's own constraint required saying so rather than widening silently.
- **D-02** — Took PF-A's *property* rather than its literal mechanism. **ESCALATE.** PF-A as specified is unreachable inside the write set: `RunInstall` returns `InstallResult`, not `resultmodel.CommandResult`, so an `ExitCodeOverride` would have to be copied in `installResultToCommandResult` (`suite_commands.go:112`) **and** threaded through `UpdateResult` and `handleUpdateSuite` — two files outside Scope. Rather than widen again or fall back to PF-B, the builder kept the property that matters: exactly one goroutine can reach `os.Exit`, it is the main goroutine, and it does so from a deferred step after `cleanup()`. **Value:** strictly better than PF-B, which keeps two exit owners and adds a blocking wait. **Risk:** the interrupted status is *taken* rather than *returned*, so the rendered JSON result is never written on the interrupted path — a consumer capturing stdout sees nothing and must read the exit code. One line in each of those two files would make it returnable; recorded as a discovered task rather than done unilaterally.
- **D-03** — The flag is stored **before** `cancelWork()`, and that ordering is the whole correctness argument. **DECIDE & STATE.** Cancellation is the only thing that releases a main goroutine parked in `readConfirmation`'s select or in an `exec.CommandContext` subprocess, so any path that can observe the cancellation necessarily observes the record. The interrupted status is deterministic rather than scheduling-dependent. A signal landing after `exitIfInterrupted` has read the flag exits on the ordinary status, which is correct — it arrived after the work concluded.
- **D-04** — `recoveryFinished` deleted rather than left in place. **DECIDE & STATE.** With the handler no longer waiting, the channel had no reader.
- **D-05** — An interrupted mid-write install exits **130, not the rollback's 3**. **DECIDE & STATE.** The process was killed by a signal, and 128+signo is what shells and CI read; `interruptedInstallExitStatus`'s existing doc already states 130 unconditionally for HUP/INT/TERM; and the rendered result still carries `outcome: rolled_back` with the full rollback record, so the number says "interrupted" and the payload says what happened. It is also the non-behaviour-changing choice — 130 is what this case already produced whenever the handler won the race.
- **D-06** — The F4 fixture blocks in `just`, not in `cp`. **DECIDE & STATE.** `just` is invoked exactly twice per install (`:562` pre-confirmation, `:900` post-write validation), so a stub that succeeds once then parks pins the child with `writeStarted` set and every managed path already replaced. It reuses the proven shape of `installFlakyJust` and makes the test independent of whether the host has `just`.

## Open Questions

- [~] Should an interrupted confirmation get its own outcome, or keep reporting `OutcomeSuccess` / `Cancelled: true` and narrating "Installation cancelled; no files were changed."? → **D-07**: Builder chose the minimal fix — narration and outcome unchanged, exit status only. Reasoning: no second defect flows from the conflation. No filesystem effect differs, `writeStarted` is false, and the exit status now separates the two cases for anything that reads it; widening would change a `SkippedWork` code and an outcome other assertions read. **A JSON consumer still cannot distinguish "user declined" from "user interrupted"** — both render `outcome: success`, `Cancelled: true`, `INSTALL-CANCELLED` — and on the interrupted path the rendered result is never written at all, because `exitIfInterrupted` takes the process before `installResultToCommandResult`. **Value:** a distinct interrupted outcome would let an automated caller retry a genuine interruption while honouring a deliberate decline, and would make the interrupted result renderable instead of swallowed by the exit. **Risk:** a new outcome needs a number in `resultmodel.ExitCode` and would move `install-suite`'s cancelled path off exit 0, which the public shell path depends on; a new skip code alone is cheap but changes bytes existing consumers may match on.

<!-- D-XX counter: last used D-07. Next decision: D-08. -->

## Discovered Tasks

- `installResultToCommandResult` and the `UpdateResult` path have no way to carry a status override, which is why the interrupted exit must be taken inside `RunInstall` rather than returned. Plumbing `ExitCodeOverride` through those two files would let the interrupted result actually reach stdout before the exit — which is what the Open Question above needs in order to be actionable.
- `runBuiltCLIAtBlockedConfirmation` and `runBuiltCLIInterruptedAtMarker` now share most of their body. Not worth merging while their synchronization points differ, but a third built-child interruption test should collapse them.

## Testing

**Tests run:** `go build ./... && go vet ./... && gofmt -l .`; `go test -count=1 ./internal/suiteinstall/`; `go test -race -count=1 ./internal/suiteinstall/`; `go test -count=1 ./...` twelve times; canonical repository gate `bash _dev/tests/maintainer-verify.sh`.
**Result:** ✓ `internal/suiteinstall` green, including under `-race` (no data race on the new atomic flag). **The gate now fails on exactly one test** — `internal/toolboxcommands` → `TestRemediationCancellationReachesMediaGitCommitAndRollback`, the REQ-524 baseline. This REQ's own test no longer appears in it, and the gate is deterministic again for the first time this session.

**Red-green validation:** traces the REQ's replacement acceptance criterion — an interrupted confirmation reports exit 130 regardless of goroutine scheduling.

Neutered (production file restored from `HEAD`, tests kept), with the new `GOMAXPROCS=1` child pin — **9 of 9 subtests red**:
- `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation`, all six: `install-suite` and `update-suite` × HUP/INT/TERM, each `signal … exit = <nil>, want 130`
- `TestBuiltInstallExits130AfterRecoveringASignalInterruptedMidWriteInstall`, all three: `signal … exit = exit status 3, want 130`

That second failure text is the F4 conflation showing itself directly — `exit status 3` is `ExitCode(OutcomeRolledBack)`, so recovery ran correctly and the ordinary return path then reported the rollback instead of the interrupt.

Post-fix: all 9 green.

**The control, run in both directions — this is what makes the pin a lock-in rather than a preference.** On the *same neutered tree*, with `singleProcessorEnvironment()` temporarily returning bare `os.Environ()`, **0 of 9 subtests failed across 3 runs**. The defect is invisible at default `GOMAXPROCS`. So the one-processor pin is precisely what converts this flake class into a deterministic failure, at no cost.

**Deterministic reproduction of the original bug** (from the diagnosis pass, before any edit): running the built child under `GOMAXPROCS=1` failed 6 of 6 confirmation subtests on every run, with or without a settle delay after the prompt. A 500 ms delay past the prompt bytes, with the read unquestionably parked in its `select`, still failed — which is what refuted the write-then-read window hypothesis rather than merely leaving it unconfirmed.

**Baseline flakiness, for the record:** on the unmodified tree, `go test -count=1 ./...` failed this test in **3 of 7** runs. After the fix, **12 consecutive full-module runs** show exactly one failure each, always the REQ-524 baseline.

**New tests added:**
- `TestBuiltInstallExits130AfterRecoveringASignalInterruptedMidWriteInstall` — the mid-write interrupt case (F4), which had no coverage at all. Helpers `installPostWriteBlockingJust`, `runBuiltCLIInterruptedAtMarker`, `waitForMarkerFile`.
- `singleProcessorEnvironment` and the extracted `interruptingSignalCases`, shared by both built-child helpers.

**Existing tests updated:** `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation` gained the `GOMAXPROCS=1` child pin. **No assertion was weakened** — the exit-130 check and the managed-path non-effect checks are byte-unchanged, which the REQ's constraints required.

**`internal/nextselection`**, which the diagnosis pass saw fail once in seven runs, did **not** fail once across the twelve runs here. Nothing to capture.

*Verified by work action*

## Review

**Overall: 91%** | 2026-09-03T12:50:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 82% |
| Test Adequacy | 88% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 — `exitIfInterrupted` had no `installVerified` guard, so a signal landing after the last context-aware subprocess (`just --list`, `:928`) exited **130 on a complete, verified install with zero bytes on stdout** while stderr reported success. Nothing after that point reads `ctx`, so main never observes the cancellation; `cleanup` skips recovery because `installVerified` is already true. The class pre-existed as a race — the same probe on the pre-fix tree gave `exit=0` once and `exit=3` twice — but this outcome was new and certain — `impact-user-visible` → fixed in remediation (D-08)
- F2 — `os.Exit` in `RunInstall`'s defer skips `RunUpdate`'s `defer os.RemoveAll(updateTmp)`, so an interrupted `update-suite` leaks its extracted upstream tree. Measured twice: `/tmp/do-work-update.*` 258 → 261 pre-remediation, 306 → 309 after. Pre-existing, but previously the defer sometimes ran when main won the race and now it never does — `impact-negligible` (scratch under `TMPDIR`, no repository effect) → not fixed; `update_transaction.go` is outside Scope and there is no in-scope fix (D-10). Covered by the existing `ExitCodeOverride` discovered task.

**Minor findings:** 4 (report only) — the `write_set` frontmatter mirror was not refreshed when `## Scope` widened; D-06 cited pre-diff line numbers (`:562`/`:900`, now `:590`/`:928`) though its two-call-site claim is correct; the built-child signal path is now exercised only at `GOMAXPROCS=1`, so the default-scheduling direction lost the coverage it had by accident; and the new helper's 60s/30s timeouts are 6× the existing helper's with no stated reason.

**Acceptance:** Pass — independently re-run in a scratch copy, only the pre-existing REQ-524 failure remains. The reviewer worked through **eleven distinct concurrency windows** in a table and confirmed D-03's argument holds for every path that observes cancellation, identifying the single window it does not cover (F1). It also verified F4 asserts *recovery*, not just the exit status: with `runRecoveryIfNeeded` short-circuited, all three subtests fail on the narration assertion **and** all four managed-path assertions.
**Suggested testing:** 5 items — the F1 lock-in (implemented), a manual Ctrl-C at a real prompt, one built-child case left at default `GOMAXPROCS`, a real `update-suite` interrupt to size F2, and a second signal during recovery to confirm "recovery is not abortable" is intended rather than an oversight. The last four are recorded, not queued.
**Follow-ups created:** None — F1 was fixed here and F2 has no in-scope fix and is already covered; **sweeps appended to:** None

**The most valuable thing this review found was not a defect.** D-05 (an interrupted mid-write install exits 130, not the rollback's 3) was recorded as a judgment call. It is **already the repository's written contract**: `_dev/tests/install-suite-behavior.sh:650` has been asserting exit 130 after a TERM during module installation all along, inside `contract-regressions.sh` inside the canonical gate. The reviewer then showed that assertion was *latently flaky in the same way* — pre-fix code under `GOMAXPROCS=1` fails it with `installer exited 3 after a TERM during module installation (want 130)`, fixed code passes. So this REQ turned a second, previously hidden gate assertion deterministically green. Nothing anywhere in the repository reads exit 3 from this path.

*Reviewed by review-work action*

## Remediation

Verdict was Approve at 91% with F1 explicitly non-blocking. F1 was fixed anyway: it is a user-visible wrong exit status on a *successful* install, its fix is one guard, and `exitIfInterrupted`'s own doc comment already contained the argument for it.

**Remediation commit:** see `commit:` below.

- **D-08** — The guard is `installVerified` alone, not a broader "concluded successfully" condition. **DECIDE & STATE.** `cleanup` skips recovery in three states (`runRecoveryIfNeeded`, `:1012`): `!writeStarted`, `installVerified`, `recoveryRan`. Only the middle one may keep its ordinary status. **`!writeStarted` is the interrupted-confirmation case this exit owner exists for — and it also returns `OutcomeSuccess`, so any predicate phrased on outcome or on "success" would silently re-break REQ-525's own fix.** `recoveryRan` means the work was undone, so the interrupt is the informative status. **Value:** the condition is stated where it is decided and cannot be widened by accident into the cancelled-confirmation path. **Risk:** the flag is set one line before the narration, so anything inserted between it and the exit owner that can still fail would exit 0 on a broken install; the doc comment now names that invariant so a future writer has to contradict it in prose to break it.
- **D-09** — The lock-in parks the child *inside* the window with stderr back-pressure rather than timing the signal. **ESCALATE.** The reviewer's proposed shape — a `just` stub that signals its parent and exits 0 — **does not reproduce F1**: built first, the neutered tree passed 1/1 and 3/3, because the window between the reaped subprocess and `exitIfInterrupted` is sub-millisecond. Signalling from inside the stub before it exits is worse: `exec.CommandContext`'s `watchCtx` injects `context.Canceled` whenever `Cancel()` succeeds, so a cancel landing before the parent reaps turns post-write validation into a failure and produces the *mid-write* case instead. The test therefore proves the reap with a detached grandchild spinning on `kill -0`, then hands the child a narration pipe filled to exactly one buffer so the success line — the first write after `installVerified = true` — blocks until the test drains it. **Value:** no sleep, no retry, no tolerance; 3/3 red when neutered on every run, 30/30 green with the guard. **Risk:** ~120 lines of helper machinery encoding two platform facts (`exec` restores blocking mode on an `*os.File` child descriptor via `Fd()`, so the buffer size is measured on a throwaway pipe; a sub-`PIPE_BUF` write cannot interleave with the filler, which is what makes the drained remainder exactly the child's own narration). Recorded as a discovered task: there are now three near-copy built-child helpers.
- **D-10** — F2 not fixed; `update_transaction.go` untouched. **DECIDE & STATE.** `RunInstall` receives only `ExtractedSourceRoot` (`updateTmp/fresh`), so removing `updateTmp` from inside the install transaction means deleting a directory it did not create by guessing the caller's layout. There is no in-scope fix, and the Scope boundary holds rather than being widened a second time. The guard does narrow the class: an update interrupted *after* the install verifies now returns through `RunUpdate`, whose defer removes the tree.

**F1 before and after, same test:** neutered → 3 of 3 red, `HUP/INT/TERM exit = 130, want 0; stdout carried 0 bytes`, with captured narration `Install this complete four-skill suite? [y/N] Installed do-work suite v0.200.0 with four verified modules.` — a verified install, exit 130, empty stdout. With the guard → 3 of 3 green: exit 0, stdout parses as one `CommandResult` with `outcome: success` and `rollback: not_needed`, installed `VERSION` reads `0.200.0`, no recovery line in the narration.

**Neuter-and-confirm:** whole package on the neutered tree reddens **only** the new test; all nine pre-existing interrupt subtests and the `GOMAXPROCS=1` pin stay green.

**The live shell consumer still gets 130.** `bash _dev/tests/install-suite-behavior.sh` exits 0 — its `:650` stub signals during `cp` of the second module, so `installVerified` is false and the guard correctly does not apply. Orchestrator re-ran this independently.

## Lessons Learned

**What worked:**
- Making the diagnosis a separate, write-nothing pass with explicit permission to stop and report a product defect. The builder used it: it refuted the captured premise *and* my own hypothesis, then declined to touch the test. Had it been told simply to fix the flake, the honest move — reporting that the test was right — would have looked like failing the task.
- `GOMAXPROCS=1` on the child, and running the control in both directions. The pinned neutered tree fails 9 of 9 every run; the unpinned one fails 0 of 9 across three. Without that second half, the pin would have looked like a preference instead of the thing that makes these lock-ins.
- Letting the builder reject the review's proposed test. That shape provably does not reproduce the bug, and finding out cost one build rather than a shipped decorative test.

**What didn't:**
- My capture. I filed this as `effort-mechanical` test synchronization and wrote that the test "signals on a timer"; it polls for the prompt text, and the defect was a production signal-handling bug that made `Ctrl-C` report success. The triage was wrong in kind, not in size.
- My hypothesis. I proposed the window between writing the prompt bytes and entering the read. A 500 ms delay past that point still fails, with the read parked in its select — refuted, not merely unconfirmed.
- The first fix, which took 130 unconditionally and so reported exit 130 with empty stdout for a complete verified install. One exit owner was necessary and not sufficient: the owner also has to know what the work concluded.
- Four sightings across three subtests before I pulled this forward. The evidence was in the captured stderr of every one of them — `Install this complete four-skill suite? [y/N] Installation cancelled; no files were changed.` — sitting in my own gate logs, saying plainly that main had taken the declined branch.

**Worth knowing:**
- `_dev/tests/install-suite-behavior.sh:650` already asserted exit 130 from this path, inside the canonical gate, and was latently flaky the same way. A decision that looks like a judgment call is worth grepping for first: the contract may already be written down.
- `exec.CommandContext`'s `watchCtx` injects `context.Canceled` whenever `Cancel()` succeeds, so cancelling before the parent reaps a subprocess converts a successful validation into a failure. That is why signalling from inside a stub reproduces the mid-write case rather than the post-verify one.
- `exec` puts an `*os.File` child descriptor back into blocking mode via `Fd()`, so a narration pipe's buffer size has to be measured on a throwaway pipe rather than assumed.
- An interrupted `update-suite` still leaks its extracted upstream tree, because `os.Exit` in a defer skips the caller's `RemoveAll`. Scratch under `TMPDIR`, no repository effect, and it resolves with the `ExitCodeOverride` plumbing.

## Orientation

Pressing `Ctrl-C` at the suite installer's confirmation prompt now reports exit 130 instead of sometimes reporting success, and an interrupted mid-write install reports 130 after its recovery completes rather than reporting the rollback. A signal that arrives after the install has already verified keeps the ordinary exit 0 and its rendered result. Lives in the CLI's `suiteinstall` transaction.

`prime_files`: `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` — spot-checked, referenced paths all still exist. Its `[family: interruptible-blocking-io]` trap (REQ-451) governs this code and still describes it accurately: `readConfirmation` is byte-unchanged and still cancellable, which the review verified rather than assumed.

**[MAP CHANGED]** — the interruption contract now has a single exit owner on the goroutine that owns the work, and the previously hidden second consumer of that contract (`install-suite-behavior.sh`'s TERM probe) is deterministically green for the first time. New lesson family `single-exit-owner`.
