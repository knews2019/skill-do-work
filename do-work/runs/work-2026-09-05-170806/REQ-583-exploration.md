# REQ-583 test-surface exploration (read-only)

All paths absolute under `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/`.
CLI root: `skills/do-work/tools/do-work-cli/`.

Everything below marked "verified" was checked by building the CLI into the scratchpad
(`go build -o $SCRATCH/do-work-cli ./cmd/do-work-cli`) and running it against throwaway
fixture repos in the scratchpad. No repository file was modified.

---

## 1. Existing test file shape

`skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go`,
537 lines, `package lifecycleadvance` (line 1) — same package as the product code, so every
unexported identifier in `evidence_gates.go` is directly callable from tests.

Style: **every** test drives a real built CLI binary as a subprocess and decodes JSON.
No test in this package calls product-code unexported functions directly; the unexported
things they call (`advanceTreeDigest`, `writeAdvanceRequest`) are test helpers only.

Test functions (evidence_gates_test.go):
- `:14 TestAdvanceExecutesEstimateGateAtPublicCLISeam` — plain, hand-rolled anonymous struct decode.
- `:52 TestAdvanceEvidenceGatesReturnTypedMissingInputs` — table-driven.
- `:78 TestAdvanceExecutesPreflightAndProjectsGreenEvidence`
- `:113 TestAdvanceQualificationUsesExactRangeAndRunsScopeDrift`
- `:142 TestAdvanceGreenGateMissRequiresDirectRunThenRecordsIt`
- `:165 TestAdvanceFocusedTestGateClassifiesBaselineStates` — table-driven, 7 rows, the focused-gate workhorse.
- `:207 TestAdvanceGateInputsFailClosedAndNeverInterpolateHostileTokens`
- `:225 TestAdvanceFocusedTestGateDistinguishesTimeoutAndLaunchFailure` — table, timeout vs `PATH=` launch failure.
- `:324 TestAdvanceFocusedGateNeverClearsFailedExecutionAgainstMatchingBaseline` — REQ-506 pin, two subtests.
- `:386 TestAdvanceMissingInputContinuationsPreserveArgumentChannels` — two subtests, follows emitted continuations.

Helpers in this file:
- `:260 runAdvanceGateJSON(t, repositoryRoot, arguments...) (resultmodel.CommandResult, int)` — prepends `--repo-root R --format json advance`.
- `:281 findAdvanceGate(result, gateID) *resultmodel.AdvanceGateRecord`
- `:293 gateHasFinding(gate, code) bool`
- `:302 hasAdvanceResultFinding(result, code) bool`
- `:442 initAdvanceGitFixture`, `:451 recordAdvanceGreenGate`, `:461 gitOnlyPathDirectory`,
  `:474 runAdvanceGateJSONWithPath` (sets `PATH=`), `:485 runAdvanceContinuation`,
  `:494 decodeAdvanceGateRun`, `:514 substituteAdvancePlaceholder`, `:530 advanceArgvIndex`.
- Constants: `:318 focusedGateRouteABody` (shortest claimed Route A body reaching the test-gate phase),
  `:322 canonicalGateFixtureBinary = "/usr/bin/true"`.

Helpers from the sibling file `advance_commands_test.go`:
- `:315 routeCBodyThrough(lastSection)`
- `:341 writeAdvanceRequest(t, root, treeSection, requestID, status, frontmatter, body) string` → returns the slash path `do-work/<section>/<ID>-fixture.md`
- `:349 writeAdvanceFile(t, root, relativePath, contents)`
- `:379 advanceCLIBinary(t)` — `sync.Once`, `go build -o <tmp>/do-work-cli ../../cmd/do-work-cli`
- `:439 runAdvanceGit(t, root, args...)`

There is **no** fixture helper that calls the gate functions in-process. A direct unit test
on `focusedGateState` would be the first of its kind in this package. That is legal (same
package clause) and is the only way to reach the M2 branches — see §3.

---

## 2. M1 — `redirectHelperRemedies`

Code: `evidence_gates.go:339-353` (`redirectHelperRemedies`), `:355-370` (`advanceArgvCommandVerb`).
Call sites: `:171` in `composeCoreGate`, `:211` in `composeGreenGate`. Both pass
`advanceGatePhaseContinuation(advance, inputs)` (`:332`), which is
`advanceGateContinuation(advance, inputs, nil, inputs.phaseArgv, inputs.separatorSeen)`.

What it changes: only two fields **on each finding** — `CommandFinding.NextArgv` (`next_argv`)
and `CommandFinding.VerificationArgv` (`verification_argv`). Nothing else. `NextJustRecipe`
is untouched (it stays `""`). The record-level `AdvanceGateRecord.NextArgv` is set elsewhere
(`missingAdvanceGateInput` `:262`, and the green-gate direct-run branch `:216`), so the
existing tests that assert on `gate.NextArgv` prove nothing about this function.

The rewrite fires when `advanceArgvCommandVerb(argv) == subordinateCommand`, i.e. the argv
starts with the literal `do-work-cli` and its first non-flag token equals the gate's own
subordinate command name.

**The producer that reliably hits this path is `run-blocked-check`**
(`internal/corehelpers/commands.go:577`): every blocked-check finding carries

```go
[]string{"do-work-cli", CommandBlockedCheck, "--probe-file", probeFile},
[]string{"do-work-cli", "--format", "json", CommandBlockedCheck, "--probe-file", probeFile}
```

so BEFORE the rewrite a finding's remedy is literally
`do-work-cli run-blocked-check --probe-file focused.sh`
(verification: `do-work-cli --format json run-blocked-check --probe-file focused.sh`),
which sends the action back into the helper the gate just ran.

AFTER the rewrite — verified by running the CLI, probe `exit 0`, no baseline:

```
finding BLOCKED-PROBE-SUCCEEDED
  next_argv        = [do-work-cli --format json advance REQ-713 --request-path do-work/working/REQ-713-fixture.md -- --probe-file focused.sh --timeout-seconds 2]
  verification_argv= [same]
```

Note the continuation carries `inputs.phaseArgv` verbatim (`--probe-file focused.sh
--timeout-seconds 2`), NOT the baseline flags `composeCoreGate` appends at `:96-99`.

Left alone (the negative control, also verified in the same run): the second finding
`FOCUSED-BASELINE-MISSING` has `next_argv: null` and stays null; `gateEvidenceFailure`'s
git remedies (`internal/gateevidence/gate_commands.go:151-152`, `git status --short` /
`git rev-parse --verify HEAD`) never match `advanceArgvCommandVerb` because argv[0] != `do-work-cli`.

**Smallest M1 test:** reuse the `TestAdvanceFocusedTestGateClassifiesBaselineStates` fixture
shape (Route A request, `focused.sh` = `exit 0`, no git init needed), run
`runAdvanceGateJSON(t, root, "REQ-xxx", "--", "--probe-file", "focused.sh", "--timeout-seconds", "2")`,
find gate `run-blocked-check`, pick the finding whose `Code == "BLOCKED-PROBE-SUCCEEDED"`, and
assert `advanceArgvCommandVerb`-equivalent: its `NextArgv` (and `VerificationArgv`) name the
`advance` verb and carry `--request-path`, not `run-blocked-check`. Deleting
`redirectHelperRemedies` flips both fields back to the helper argv, so the assertion fails.
Asserting the exact rewritten slice with `reflect.DeepEqual` is also fine — the value is
deterministic (shown above).

A second, weaker path exists via `usageResult` (`corehelpers/commands.go:427`,
`[]string{"do-work-cli", commandName}`) and `gateEvidenceFailure`'s `GATE-EVIDENCE-USAGE`
branch, but advance's own parser fails closed before most of those, so `run-blocked-check`
is the clean seam.

---

## 3. M2 — `focusedGateState` layered guard

`evidence_gates.go:183-193`:

```go
func focusedGateState(subordinateState resultmodel.AdvanceGateState, focusedTest *resultmodel.FocusedTestResult) resultmodel.AdvanceGateState {
	if subordinateState == resultmodel.AdvanceGateFailed || !focusedTest.Launched || focusedTest.TimedOut {
		return subordinateState
	}
	switch focusedTest.BaselineState { ... }
}
```

Eligibility guard, `internal/corehelpers/commands.go:545-561`:

```go
finishedOnItsOwn := probeEvidence.Launched && !probeEvidence.TimedOut && runError == nil
...
case !probeEvidence.Launched:  BLOCKED-PROBE-LAUNCH-FAILED, OutcomeFailure, error
case probeEvidence.TimedOut:   BLOCKED-PROBE-TIMED-OUT,     OutcomeFindings, warning
case runError != nil:          BLOCKED-PROBE-FAILED,        OutcomeFailure, error
case status != 0:              BLOCKED-PROBE-FAILED,        OutcomeFindings, warning
...
focusedTest := &FocusedTestResult{... BaselineState: FocusedBaselineNotCompared ...}
if baselineJSONPath != "" && finishedOnItsOwn { baselineFinding = compareFocusedBaseline(...) }
```

Why the two halves are dead through the public path:
- `BaselineState` is only moved off `not_compared` inside `compareFocusedBaseline`, which is
  called only when `finishedOnItsOwn` — launched, not timed out, no runner error.
- Every outcome that maps to `AdvanceGateFailed` (`boundAdvanceGateRecord:230`: Failure,
  Refused, RolledBack, Risk) comes from `!Launched` or `runError != nil` — both `!finishedOnItsOwn`.
  `compareFocusedBaseline` can only raise Success→Findings (`commands.go:566-568`), never to Failure.
- So `Failed`, `!Launched` and `TimedOut` all arrive with `BaselineState == not_compared`,
  where the switch has no case and falls through to `return subordinateState` — the same value
  the deleted `if` would have returned. Deleting the whole first `if` is behaviourally invisible
  at the CLI seam, which is why `TestAdvanceFocusedGateNeverClearsFailedExecutionAgainstMatchingBaseline`
  (which tries to pin exactly this) passes either way.

**How to reach the branches:** call `focusedGateState` directly. It is package-private and
`evidence_gates_test.go` is `package lifecycleadvance`, so no export and no test-only hook is
needed. Build the input by hand:

```go
// (a) failed subordinate against a matching baseline
got := focusedGateState(resultmodel.AdvanceGateFailed, &resultmodel.FocusedTestResult{
    Launched: true, TimedOut: false, BaselineState: resultmodel.FocusedBaselineMatchingRed,
})  // want AdvanceGateFailed, not AdvanceGateSatisfied
// (b) timeout against a green baseline
got = focusedGateState(resultmodel.AdvanceGateFindings, &resultmodel.FocusedTestResult{
    Launched: true, TimedOut: true, BaselineState: resultmodel.FocusedBaselineGreen,
})  // want AdvanceGateFindings
// (c) never launched against a green baseline
got = focusedGateState(resultmodel.AdvanceGateFailed, &resultmodel.FocusedTestResult{
    Launched: false, BaselineState: resultmodel.FocusedBaselineGreen,
})  // want AdvanceGateFailed
```

Each row fails if the corresponding clause of the guard is deleted (with the guard gone,
(a) and (b)/(c) all return `AdvanceGateSatisfied`). A table with the eight combinations of
{Failed/Findings} x {Launched} x {TimedOut} x {Green/MatchingRed/NewRed} is cheap and runs in
microseconds — no subprocess, no fixture repo.

Worth saying in the REQ: this is a defence-in-depth guard whose only current caller cannot
produce the input. The direct unit test is honest about that (it pins the function contract,
not an end-to-end behaviour); a public test would need a second, hypothetical producer of
`FocusedTestResult`.

---

## 4. Interrupted focused test

`corehelpers/commands.go:549-557` — exact dispatch, in order:
1. `!probeEvidence.Launched` → `BLOCKED-PROBE-LAUNCH-FAILED`, OutcomeFailure, error.
2. `probeEvidence.TimedOut` → `BLOCKED-PROBE-TIMED-OUT`, OutcomeFindings, warning.
3. `runError != nil` (launched, not timed out, runner error) → `BLOCKED-PROBE-FAILED`, **OutcomeFailure, error**.
4. `status != 0` → `BLOCKED-PROBE-FAILED`, OutcomeFindings, warning.

Case 3 is the interrupted case: `internal/nextselection/blocked_probe_unix.go:76-82` returns
`BlockedProbeEvidence{ExitStatus: 128+signal, Launched: true}` plus a typed
`BlockedProbeInterruption` when the **CLI process itself** receives SIGHUP/SIGINT/SIGTERM
while the probe runs (`signal.Notify` at `blocked_probe_unix.go:21`).

Existing coverage: `evidence_gates_test.go:347` asserts
`!gateHasFinding(*focusedGate, "BLOCKED-PROBE-LAUNCH-FAILED") || status == 0` inside the
"current launch failure" subtest of `TestAdvanceFocusedGateNeverClearsFailedExecutionAgainstMatchingBaseline`
— that is case 1 only (`PATH=` / `gitOnlyPathDirectory`, so `sh` cannot be found). Case 3
(`BLOCKED-PROBE-FAILED` at error severity, the interrupted run) has **no** assertion anywhere
in this package; the only other repo hits are in `internal/nextselection/next_selection.go:211,231`
and `next_selection_test.go:252`, which are the unrelated `next`-selection exclusion code of the
same name. So today, collapsing case 3 into case 4 (or dropping it) fails nothing in
`lifecycleadvance`.

**Public test for it, verified working.** Probe file contents:

```
kill -TERM $PPID; sleep 5
```

The probe's `$PPID` is the `do-work-cli` process (the probe runs as `sh -c` started by the CLI
with `Setpgid: true`, so the CLI stays its parent). The CLI's `signal.Notify` is already armed
when the child starts, so the signal is caught, the probe group is torn down, and the CLI still
renders JSON and exits normally. Observed (Route A REQ-713 fixture, `--timeout-seconds 20`):

```
exit status 2, outcome "failure", 0.5s wall
gate run-blocked-check: state "failed", outcome "failure"
focused_test: {exit_status:143, launched:true, timed_out:false, baseline_state:"not_compared"}
finding: BLOCKED-PROBE-FAILED, severity "error",
  observed_evidence: ["raw probe status 143; diagnostic sha256 e3b0c442...: blocked probe interrupted with exit 143"]
```

So the assertions are: `Launched == true`, `TimedOut == false`, `ExitStatus == 143`,
gate state `AdvanceGateFailed`, `gateHasFinding(..., "BLOCKED-PROBE-FAILED")`, non-zero CLI
status. That row distinguishes itself from case 4 (ordinary red) only through the gate state
`failed` vs `findings` — worth asserting both, and worth pairing with an ordinary-red row
(`exit 143` without the kill, which is launched+`runError == nil` → `findings`) as the negative
control that makes the mutation "drop case 3" red.

Do not copy the in-process technique from `internal/nextselection/blocked_probe_test.go:81-116`
(`syscall.Kill(os.Getpid(), SIGINT)` with a guard `signal.Notify` in the test) — that works
there because the runner is called in-process. Here the CLI is a subprocess, and `kill -TERM
$PPID` from inside the probe is simpler and needs no goroutine, no timing sleep, and no signal
guard in the test process. The alternative (start `exec.Cmd`, poll for a marker file, then
`command.Process.Signal(syscall.SIGTERM)`) also works but adds a race.

Note the existing `runAdvanceGateJSON` helper handles a non-zero exit fine, so no new helper
is needed. Runtime cost of the row is ~0.5s, which matters because REQ-574 put this package
under a 30s-per-file budget.

---

## 5. How tests invoke the gate publicly

Canonical shape (`evidence_gates_test.go:260`):

```go
repositoryRoot := t.TempDir()
requestPath := writeAdvanceRequest(t, repositoryRoot, "working", "REQ-713", "claimed",
    "route: A\nestimate:\n  p50_active_minutes: 5\n", focusedGateRouteABody)
writeAdvanceFile(t, repositoryRoot, "focused.sh", "exit 0")
result, status := runAdvanceGateJSON(t, repositoryRoot,
    "REQ-713", "--", "--probe-file", "focused.sh", "--timeout-seconds", "2")
focusedGate := findAdvanceGate(result, "run-blocked-check")
```

which becomes `do-work-cli --repo-root <tmp> --format json advance REQ-713 -- --probe-file
focused.sh --timeout-seconds 2`. Pre-separator inputs are `--request-path`, `--diff-range`,
repeated `--gate-arg`, `--gate-exit-status`; everything after `--` is the phase argv
(`parseAdvanceGateInputs`, `evidence_gates.go:105`).

Copy the shape from, in order of usefulness:
- `:165 TestAdvanceFocusedTestGateClassifiesBaselineStates` — the minimal focused-gate row (no git init).
- `:324 TestAdvanceFocusedGateNeverClearsFailedExecutionAgainstMatchingBaseline` — when you need git +
  a recorded green gate (`initAdvanceGitFixture` + `recordAdvanceGreenGate(t, root, canonicalGateFixtureBinary)`)
  so the green-gate half does not turn the run into `needs_input`.
- `:225 TestAdvanceFocusedTestGateDistinguishesTimeoutAndLaunchFailure` — when you need a custom `PATH`
  or to build the `exec.Cmd` by hand.

Caveat seen in the verified run: without `--gate-arg`, the test-gate phase always adds a second
record `check-green-gate` in state `needs_input`, so overall status is non-zero and
`result.Outcome` is `findings` even when the probe passes. Assert on the `run-blocked-check`
record, not on the aggregate, unless you also record a green gate.

---

## 6. Prior art

`do-work/archive/REQ-581-make-the-descendant-cleanup-tests-fail-on-a-real-process-group-leak.md`
(completed 2026-09-05, commit `92339213`, test-only, one file):
1. It named an explicit RED mutation up front — reduce `terminateOwnedProcessGroup` and
   `cleanupReapedProcessGroup` to no-op bodies — and required the builder to keep that mutation
   applied while writing the fixture, reverting only at the end.
2. The structural diagnosis was that the assertion could not fail: a surviving descendant
   inherited the parent-owned diagnostic pipe, so the runner blocked until it exited and the
   descendant was always gone by the time the poll ran.
3. The fix was to the fixture, not the budget: the descendant closes its inherited stdout/stderr
   before sleeping, and outlives the assertion budget, so a leak is still alive when the poll looks.
4. The assertion moved from "how fast it disappears" to "is it alive", and the failure message
   names the surviving pid.
5. One test was failing on an unrelated 5s "did not return" bound instead of its own assertion;
   the bound was raised past the new leader hold so the test reaches the assertion it exists for.

Same family as REQ-583: state the mutation, prove the test is red under it, and make the
assertion name the real failure rather than a proxy.

`skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (112 lines), relevant entries:
- Read first: `internal/lifecycleadvance/` owns request-bound execution of each mechanical
  evidence gate; `internal/nextselection/` owns the process-tree-owned blocked probe. The Advance
  phase table states "Queue-mode `advance` and the focused-test gate consume `nextselection`'s
  bounded blocked-probe evidence; `run-blocked-check` is not a second authority" — which is the
  contract `redirectHelperRemedies` enforces at the remedy level.
- Trap `[family: smoke-vs-characterization]` — registration/smoke checks stay green while
  semantics diverge; compare status, ordered evidence, actions, paths and effects at authority
  boundaries. This is exactly the M1/M2 situation.
- Trap `[family: reaped-by-its-own-parent]` — relevant to the interrupted-probe row: a zombie
  still satisfies `kill(pid, 0)`; signal descendants before parents.
- Trap `[family: closed-enumeration-for-a-condition]` — test the condition's ingredients, not
  the spellings today's cases use. Applies directly to the focused-gate table: the M2 rows must
  key on Launched/TimedOut/state, not on exit-status values.
- Trap `[family: silent-skip-reads-as-red]` — a lane that skips silently reads as red.
- Verify section: focused package run is `go test ./internal/lifecycleadvance`; static analysis
  `go vet ./...`; module regression `go test -count=1 ./...`.
- Lessons pointer: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (~7450 tokens,
  over the 2000-token Required Lessons budget, `slugged: partial` — REQ-581 dropped it for budget
  and recorded the four matching families instead).
