# Hand-back — REQ-581 (make the descendant-cleanup tests fail on a real process-group leak)

## Branch

- Branch: `worktree-agent-REQ-581-descendant-fixture`
- Head commit: `f5b3faf4` — `[REQ-581] make the descendant-cleanup tests fail on a real process-group leak`
- Base commit: `961697bb`
- Worktree left clean: `git status --porcelain` is empty.

## File manifest

- MODIFIED (test): `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go` — 59 insertions, 34 deletions, the only file in the write set and the only file in the commit.

No source file outside the write set was changed. `internal/nextselection/blocked_probe_unix.go` was mutated as a scratch experiment and reverted with `git checkout --` before the commit; see **Mutation revert proof** below.

## P-A-U

### [PLAN]

Read `prime-do-work-cli.md`, the four required crew-member rules plus `testing.md`, and the four lesson families the change touches. Then:

1. Apply the RED mutation from the request (both teardown functions reduced to no-op bodies) and reproduce the reviewer's reported numbers on the existing fixture, to confirm the defect is what the request says it is.
2. Keep the mutation applied and rewrite the fixture on three points:
   - the background descendant releases its inherited stdout and stderr before sleeping, so it cannot hold the runner's diagnostic pipe open and cannot delay the runner's return;
   - the descendant sleeps 30s, far past the 10s `descendantReapBudget`, so a leaked descendant is still running when the poll looks at it;
   - the probe's shell leader stops waiting on the descendant and instead holds itself for a bounded 4 seconds. The leader must still be alive when the runner's own timeout or interrupt fires (otherwise those branches never run), and with teardown broken nothing else ends it, so a bounded self-hold both keeps the branches firing and caps how long a broken runner can block.
3. Raise the interruption test's "interrupted probe did not return" bound from 5s to 15s so it outlasts the 4s leader hold. Without this the test keeps failing at that bound under the mutation instead of at its descendant assertion, which is exactly the defect the request names.
4. Rewrite the `descendantReapBudget` comment, which claimed the loop proved the descendant does not survive.
5. Revert the mutation, prove all three green with timings at or below the pre-change baseline, then run the whole module.

Approach conflicts with a lesson: none. The `reaped-by-its-own-parent` lesson is what the plan relies on — a zombie still satisfies `kill(pid, 0)`, which is why the descendant has to outlive the budget for the poll to mean anything.

### [APPLY]

Code written exactly as planned, inside `blocked_probe_test.go` only.

- Added `leakedDescendantBody = "exec 1>&- 2>&-; sleep 30"` — the shared descendant body for all three tests, with a comment explaining why releasing the inherited descriptors is what makes a leak observable at all.
- Added `probeLeaderHoldSeconds = 4` with a comment stating both jobs it does.
- Added `interruptedProbeReturnBudget = 15 * time.Second`, replacing the inline `5 * time.Second` in the interruption test.
- Rewrote the three probe scripts:
  - timeout test: `(trap '' TERM; <body>) & echo $! > <path>; sleep 4` (was `(trap '' TERM; sleep 30) & echo $! > <path>; wait`)
  - background-descendant test: `(<body>) & echo $! > <path>; exit 0` (was `sleep 30 & echo $! > <path>; exit 0`)
  - interruption test: `(<body>) & echo $! > <path>; sleep 4` (was `sleep 30 & echo $! > <path>; wait`)
  - The `trap '' TERM` in the first test is kept: it is what forces the SIGKILL escalation rather than letting the initial SIGTERM do the work.
- Added `waitForDescendantPid(t, pidPath)`, replacing three copies of pid-file reading (two unchecked `strconv.Atoi` calls and one hand-rolled 2s poll loop in the interruption test). It polls for the file, tolerates a not-yet-written or partially written file until its deadline, and registers a `t.Cleanup` SIGKILL on the pid.
- Rewrote the `descendantReapBudget` comment.

### [UNIFY]

`git diff --stat`: `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go | 93 ++++++++++++-------- ; 1 file changed, 59 insertions(+), 34 deletions(-)`.

Files verified, and what was checked on each:

- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go` — read the full diff; every changed line traces to a Detailed Requirement or to the pid-reading it replaced. Confirmed the four untouched tests in the file (`TestBlockedProbePreservesRawStatus`, `TestBlockedProbeEvidenceBoundsAndNormalizesDiagnostics`, `TestBlockedProbeEvidencePreservesOrdinaryReservedExitValues`, `TestBlockedProbeEvidenceRefusesUnrunnableInputsAsUnlaunched`) are byte-identical. Checked new names against the two-word naming rule: `leakedDescendantBody`, `probeLeaderHoldSeconds`, `interruptedProbeReturnBudget`, `waitForDescendantPid`. Grepped the added lines for `fmt.Print`, `println`, `TODO`, `FIXME`, `XXX`, `t.Skip`, `debug` — no matches.
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go` — verified byte-identical to base: `git diff` on the path is empty and it does not appear in `git status --porcelain`.

Linters, all from `skills/do-work/tools/do-work-cli`:

- `gofmt -l .` — no output (exit 0). One gofmt fix was needed and applied: the new comment block after `TestBlockedProbeEvidenceBoundsAndNormalizesDiagnostics` required a blank line the original file did not have.
- `go vet ./...` — exit 0, no output.
- `GOOS=windows GOARCH=amd64 go test -c ./internal/nextselection -o <scratch>` — exit 0. This is the prime's listed Windows blocked-probe check; the file is `//go:build unix`, so this proves the change did not break the Windows build of the package.

No `node --check` run: no client files changed.

## Test evidence

All commands run from `skills/do-work/tools/do-work-cli` unless the path says otherwise. The three tests are abbreviated below as T1 `TestBlockedProbeTimeoutKillsDescendantGroup`, T2 `TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits`, T3 `TestBlockedProbeInterruptionIsTypedAndReapsDescendants`.

### Step 1 — baseline, old fixture, unmutated tree (exit 0)

`go test -count=1 -v -run '<T1|T2|T3>' ./internal/nextselection`

```
--- PASS: TestBlockedProbeTimeoutKillsDescendantGroup (1.50s)
--- PASS: TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits (0.11s)
--- PASS: TestBlockedProbeInterruptionIsTypedAndReapsDescendants (0.63s)
ok  	.../internal/nextselection	2.624s
```

This is the number the "do not raise the unmutated runtime" constraint is measured against. It is faster than the 2.90s / 2.01s the request quotes, so the machine was quieter than the reviewer's; the comparison that matters is baseline against post-change on the same machine in the same session, below.

### Step 2 — RED reproduction, old fixture, mutated tree (exit 1)

Mutation applied to `internal/nextselection/blocked_probe_unix.go`: `cleanupReapedProcessGroup` reduced to `_ = processGroup` and `terminateOwnedProcessGroup` to `_, _, _ = processGroup, initialSignal, done`.

`go test -count=1 -v -timeout 180s -run '<T1|T2|T3>' ./internal/nextselection`

```
--- PASS: TestBlockedProbeTimeoutKillsDescendantGroup (30.01s)
--- PASS: TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits (30.01s)
    blocked_probe_test.go:114: interrupted probe did not return
--- FAIL: TestBlockedProbeInterruptionIsTypedAndReapsDescendants (5.13s)
FAIL	.../internal/nextselection	65.452s
```

This matches the request's reported numbers exactly (30.01s, 31.35s, and a T3 failure on its 5s bound). Two tests pass on a tree with a genuine process-group leak, and the one failure never reaches its descendant assertion.

### Step 3 — RED observation, new fixture, mutated tree (exit 1)

Same mutation still applied. Same command.

```
    blocked_probe_test.go:67: descendant 78883 survived 10s
--- FAIL: TestBlockedProbeTimeoutKillsDescendantGroup (14.02s)
    blocked_probe_test.go:78: descendant 82138 survived 10s
--- FAIL: TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits (10.02s)
    blocked_probe_test.go:115: descendant 83954 survived 10s
--- FAIL: TestBlockedProbeInterruptionIsTypedAndReapsDescendants (14.12s)
FAIL	.../internal/nextselection	38.454s
```

All three fail inside their own budget, on the descendant assertion, with the surviving process id named. No test in the package exits on a timeout bound. T3 reaches its descendant assertion, so its typed-interruption assertion (status 130, `InterruptionExitStatus() == 130`) still passed under the mutation — the mutation breaks teardown, not the typed result, and the test now separates the two.

### Mutation revert proof

`git checkout -- skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go`

```
$ git status --porcelain
 M skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go
$ git diff -- skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go
(empty)
```

`blocked_probe_unix.go` is absent from `git status --porcelain` and its diff against base is empty. The mutation is reverted and was never staged or committed. After the commit, `git status --porcelain` is empty entirely.

### Step 4 — GREEN, new fixture, unmutated tree (exit 0)

`go test -count=1 -v -timeout 180s -run '<T1|T2|T3>' ./internal/nextselection`

```
--- PASS: TestBlockedProbeTimeoutKillsDescendantGroup (1.50s)
--- PASS: TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits (0.11s)
--- PASS: TestBlockedProbeInterruptionIsTypedAndReapsDescendants (0.65s)
ok  	.../internal/nextselection	2.610s
```

Package wall time 2.610s against the 2.624s baseline, and per-test 1.50s / 0.11s / 0.65s against 1.50s / 0.11s / 0.63s. The runtime constraint holds: the change is timing-neutral, not a slowdown.

### Step 5 — flake check (exit 0)

`go test -count=3 -timeout 180s -run '<T1|T2|T3>' ./internal/nextselection` → `ok ... 6.942s`. Three consecutive repeats, no failure. The reap poll is the historically flaky part (the request records 1.13-1.95s under load against the old 2s budget), so it was worth repeating; the 10s budget is untouched and still generous.

### Step 6 — brief's focused commands (both exit 0)

```
$ bash _dev/tests/run-go-tests-with-budget.sh skills/do-work/tools/do-work-cli ./internal/nextselection/...
go-test budget: module=skills/do-work/tools/do-work-cli wall=3s tests=35 slowest-file=internal/nextselection/blocked_probe_test.go:2.30s limit=<30s

$ bash _dev/tests/run-go-tests-with-budget.sh skills/do-work/tools/do-work-cli ./...
go-test budget: module=skills/do-work/tools/do-work-cli wall=65s tests=762 slowest-file=internal/finalization/finalization_recovery_test.go:24.75s limit=<30s
```

Whole module green, 762 tests. `blocked_probe_test.go` measures 2.30s, well under the 30s per-file budget REQ-574 (bringing this module's test files under a 30-second per-file budget) established. The module's slowest file is `internal/finalization/finalization_recovery_test.go` at 24.75s and is untouched by this change.

The repository gate `_dev/tests/maintainer-verify.sh` was NOT run, per the brief.

### Acceptance against the request's GREEN criteria

- Each of the three fails inside its own budget on the descendant assertion under the mutation — yes, Step 3.
- The failure message names the surviving process id — yes, `descendant 78883 survived 10s` and the other two.
- No test in the package exits on a timeout bound instead — yes; the whole-package mutated run showed only these three failures, each on the descendant assertion.
- With the functions restored all three pass in their present single-digit-second range — yes, 1.50s / 0.11s / 0.65s, Step 4.
- `descendantReapBudget` comment corrected — yes.

## Lesson evidence

The request's `## Required Lessons — Dropped for Budget` names `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (7450 tokens, over the capture budget, `slugged: partial`) and the four families it matches on. The touch-conditional rule in `crew-members/general.md` applies independently of that budget drop, so I read the file's entries for all four families directly. The file exists; no listed path was missing.

- `[family: reaped-by-its-own-parent]` (REQ-543, reaping the commit hook with its own parent) — the load-bearing one. It is the source of "a zombie still satisfies `kill(pid, 0)`", which is precisely why `waitForDescendantToDisappear` could not be turned into a single immediate liveness check and why the descendant must outlive the budget instead. Also the source of "a wait on our own child ends in milliseconds, a wait on init is measured in seconds", which is what the 10s budget exists for and what the rewritten comment now says out loud.
- `[family: interruptible-blocking-io]` (REQ-451, making confirmation input interruptible) — "prove signal behavior with the real process while input remains open". The interruption test still sends a real SIGINT to the real test process and keeps its `signal.Notify` guard; only its return bound and its fixture changed.
- `[family: smoke-vs-characterization]` (REQ-414 and REQ-415, migrating core helpers and core memory hooks) — a smoke matrix stays green while semantics diverge. That is the exact defect here: three tests reported the cleanup path as proven while a leak passed through them. The fix is to make the observation happen at a moment where the two outcomes actually differ, not to add more assertions.
- `[family: silent-skip-reads-as-red]` (REQ-566, running held heavy lanes and recording lane wall time) — key a consumer on the announced condition, not on incidental facts. Applied by keying the failure on the descendant's liveness rather than on elapsed time.

Also read for context: `prime-do-work-cli.md` in full, including its `internal/nextselection` and `internal/ownedprocess` entries (the blocked probe deliberately keeps its own runner for the signal-forwarding and `128+signal` contract, so nothing here should be pushed into `ownedprocess`), and its `## Verify` line for the Windows blocked-probe compile, which I ran.

No lesson contradicted the plan. No new lesson-satellite entry is proposed: the underlying trap is already generalized in the prime as `[family: reaped-by-its-own-parent]`, and this request applied it rather than discovering something new. If the orchestrator wants the "a control that cannot fail hides its own defect" angle recorded, D-05 below is the honest content for it, but I did not write to the satellite because it is outside my write set.

## Decisions

**D-01 — Release the inherited descriptors with `exec 1>&- 2>&-` rather than `exec >/dev/null 2>&1`. DECIDE & STATE.**
Both release the runner's pipe and both make the leak observable. Closing matches the request's own wording ("closes its inherited stdout and stderr") and proves the stronger property: the descendant is not merely pointed elsewhere, it holds no copy of the parent's pipe at all. `sleep` writes nothing, so the closed descriptors are inert. Verified working on this machine's `/bin/sh`. Reversible in one line if a future platform's shell objects.

**D-02 — Bound the leader with `sleep 4` instead of `wait`. DECIDE & STATE.**
This is the change the request does not spell out but that its requirements force. Closing the descendant's descriptors alone is not enough: with `wait`, the leader itself keeps the pipe open for the descendant's full lifetime, so under the mutation the runner still blocks for 30s and the tests still cannot fail on their assertion. The leader has to end on its own. It must also still be alive when the runner's timeout (1s) or interrupt (~0.1s) fires, or the timeout and signal branches never execute, so the hold has to exceed those. 4 seconds gives roughly 4x margin over the largest of them while capping the mutated failure at about 14s per test. Named `probeLeaderHoldSeconds` so a future adjustment is one constant.

**D-03 — Raise the interruption test's return bound from 5s to 15s. DECIDE & STATE.**
Directly required: "must reach and fail its descendant assertion under the mutation instead of stopping at its 5-second bound". The bound must outlast `probeLeaderHoldSeconds`, and 15s leaves room for a loaded machine. It costs nothing on a healthy tree, where the runner returns in about 0.65s. The risk of a looser bound is that a runner which hangs forever now takes 15s to say so instead of 5s; that is the correct trade, because the 5s version was reporting a hang where the real defect was a leak.

**D-04 — Keep the polling helper and the 10s `descendantReapBudget` rather than replacing them with a single immediate liveness check. DECIDE & STATE.**
An immediate check would be the more direct reading of "assert on liveness, not on how fast it disappears", but it cannot be made correct with `kill(pid, 0)` alone: a correctly killed descendant is orphaned to init and may sit as a zombie for up to about two seconds, and a zombie answers `kill(pid, 0)` exactly like a live process. Distinguishing the two needs `ps -o state=` or `/proc`, which adds a subprocess and platform variance for no gain. Once the descendant outlives the budget, the poll is already a liveness assertion: a pid still answering at the 10s deadline has 20 seconds of `sleep` left and is genuinely running. Keeping the budget also keeps the green runtime unchanged, since the poll returns as soon as the reap lands.

**D-05 — Add `waitForDescendantPid` and a `t.Cleanup` SIGKILL on the descendant pid. DECIDE & STATE.**
Two reasons beyond removing three copies of the same code. First, the old code ignored `strconv.Atoi`'s error in two places, so a truncated pid file would have produced `pid == 0` and a confusing failure; the helper keeps polling until it has a positive integer and fails with the path if it never gets one. Second, and the reason it is worth the lines: this test is now designed to fail while a real process is still running. Without the cleanup, every failing run leaves a 30-second `sleep` behind, three at a time during a mutation experiment, on a machine the brief says other builders are sharing. The cleanup is registered at pid-read time, so it fires whether the assertion passes or fails.

**D-06 — Keep `trap '' TERM` in the timeout test. DECIDE & STATE.**
It makes the SIGKILL escalation in `terminateOwnedProcessGroup` load-bearing rather than incidental: without it the initial SIGTERM would do all the work and the escalation path would go unexercised. It costs the 500ms grace window, which is already in the 1.50s measurement.

**D-07 — No change to `blocked_probe_unix.go`. DECIDE & STATE.**
The request states the implementation is correct and this is a test-only change, and the evidence agrees: the rewritten tests pass against the unmodified implementation on the first try, with no timing change.

## Discovered Tasks

- The 500ms grace window and the second 500ms wait in `terminateOwnedProcessGroup` (`blocked_probe_unix.go:99-113`) are unnamed literals, unlike the named budgets the test file now uses. They are the largest fixed component of the timeout test's 1.50s. Naming them would make the relationship between implementation grace and test budget explicit, but nothing is wrong today. → `impact-negligible` → report only.
- `runBlockedProbeFixture` (`blocked_probe_test.go:167`) is a one-line passthrough to `RunBlockedProbeAtRoot` with a single caller and adds no behavior. It looks like a leftover from an earlier RED implementation. Deleting it is a two-line change, but it is outside this request's scope and unrelated to the defect. → `impact-negligible` → report only.
- The whole module's slowest test file is `internal/finalization/finalization_recovery_test.go` at 24.75s, which is inside the 30s per-file budget but with the least headroom of any file in the module. Nothing to do now; noting it because REQ-574 (bringing this module's test files under a 30-second per-file budget) just landed and this is the file most likely to breach next. → `impact-negligible` → report only.

## Integration seams

None. The change is confined to one test file, adds no exported name, and needs no line in any file outside the write set. No changelog or version entry is proposed from here; those paths belong to finalization.
