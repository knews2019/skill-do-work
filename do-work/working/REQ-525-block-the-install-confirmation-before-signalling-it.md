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
