# REQ-506 owner map

Read-only survey by four parallel readers plus a synthesis pass, produced before
builder dispatch. Line numbers are at HEAD `e4d78d81`. This is a map, not a patch.

All paths are under `ROOT = /home/user/skill-do-work/skills/do-work/tools/do-work-cli/`. Absolute file paths are given in full the first time each file appears; later references use the `ROOT`-relative form.

Files that carry the defect:
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go`
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go`
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go`
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go`
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go`
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates.go`

# 1. Defect sites

## (a) 124/125 reconstruction loses the observed launch/timeout facts

**D1 — `internal/nextselection/blocked_probe_unix.go:15`, the `runOwnedProbe` contract.**
`func runOwnedProbe(repositoryRoot string, probeBytes []byte, timeout time.Duration, diagnosticWriter io.Writer) (int, error)`. The function knows every fact first-hand: whether `command.Start()` succeeded (line 38), whether the timer fired (line 64), whether a signal arrived (line 68). It discards all of that and returns one integer. This is the origin of the defect; every other site below is a downstream attempt to rebuild what this line threw away. It is unexported and platform-split, so its signature can change without touching any public surface, but both build-tagged variants must change in lock-step.

**D2 — `blocked_probe_unix.go:64-67`, the timeout branch returns `err == nil`.**
```go
case <-timer.C:
	terminateOwnedProcessGroup(processGroup, syscall.SIGTERM, done)
	<-diagnosticDone
	return BlockedProbeTimeoutStatus, nil
```
The error channel carries no timeout information, so the integer 124 is the only carrier. Any consumer that wants "did this time out" must compare against 124, and cannot distinguish it from a probe that ran to completion and exited 124 itself.

**D3 — `blocked_probe_unix.go:21, 40, 52`, three genuine launch failures collapse to the same 125 as everything else.**
Pipe creation failure (21), `command.Start()` failure (40), and process-group isolation failure (52) each `return BlockedProbeLaunchStatus, err`. These are the only true "never ran" cases on unix, and they are the only ones that carry a non-nil error. The remediation should treat these three as the definition of `Launched == false`, not the number 125.

**D4 — `blocked_probe_unix.go:129`, `probeExitStatus` synthesizes 125 with no error.**
```go
func probeExitStatus(err error) int {
	...
	return BlockedProbeLaunchStatus
}
```
Reached when an `*exec.ExitError` carries a `Sys()` that is not a `syscall.WaitStatus`. The call site is line 60, inside the normal-completion branch, which returns `status, nil`. So a process that definitely ran reaches the reconstruction as 125 with a nil error, and is reported as `Launched: false` with nothing to explain it.

**D5 — `blocked_probe_windows.go:12`, the platform stub.**
```go
func runOwnedProbe(_ string, _ []byte, _ time.Duration, _ io.Writer) (int, error) {
	return BlockedProbeLaunchStatus, fmt.Errorf("standard-library process-tree ownership is unavailable on windows")
}
```
This is the one site where 125 is semantically right (nothing launched). It must keep returning "not launched" under whatever the new internal shape is, and it is the reason the internal shape must be expressible on both build tags.

**D6 — `blocked_probe.go:107-108`, the primary reconstruction.**
```go
	ExitStatus: status, Launched: status != BlockedProbeLaunchStatus,
	TimedOut: status == BlockedProbeTimeoutStatus, Diagnostic: diagnostic,
```
Two booleans inferred from one integer. A probe that runs and deliberately exits 125 is published as `launched: false`; a probe that runs and deliberately exits 124 is published as `timed_out: true`. Both booleans are already the right carriers — they just need to be assigned from observation instead of from arithmetic.

**D7 — `blocked_probe.go:94-102`, the three input-validation guards.**
```go
	return BlockedProbeEvidence{ExitStatus: BlockedProbeLaunchStatus}, fmt.Errorf("repository root is empty")   // 95
	return BlockedProbeEvidence{ExitStatus: BlockedProbeLaunchStatus}, fmt.Errorf("timeout must be positive")   // 98
	return BlockedProbeEvidence{ExitStatus: BlockedProbeLaunchStatus}, fmt.Errorf("probe is empty")             // 101
```
These bypass the reconstruction at 107-108 entirely and rely on the zero value of `Launched`/`TimedOut` happening to be `false`. Correct today by accident. If the booleans start being set explicitly, these three literals must set them explicitly too, or they become the new inconsistency.

**D8 — `blocked_probe.go:86-89`, `RunBlockedProbeAtRoot` is the lossy narrowing.**
```go
func RunBlockedProbeAtRoot(repositoryRoot string, probeBytes []byte, timeoutSeconds int) (int, error) {
	evidence, err := RunBlockedProbeEvidenceAtRoot(repositoryRoot, probeBytes, timeoutSeconds)
	return evidence.ExitStatus, err
}
```
`Launched`, `TimedOut`, `Diagnostic`, `DiagnosticSHA256` are computed one frame down and dropped. Its two production callers (`internal/nextselection/next_commands.go:33`, `internal/lifecycleadvance/queue_commands.go:256`) then re-derive meaning from the integer. This signature must not change (see §3); the fix is that its callers gain an evidence-preserving path, not that this function changes shape.

**D9 — `internal/nextselection/next_selection.go:207-230`, three inconsistent re-derivations in one block.**
```go
207		if probeRunner == nil {
208			evidence.ProbeStatus = resultmodel.ProbeLaunchFailed
209			evidence.ProbeAttempted = true
210			evidence.ProbeExitCode = 125
...
221		if probeError != nil || exitCode != 0 {
222			reason := fmt.Sprintf("blocked_check failed this run with exit %d", exitCode)
223			if probeError != nil {
224				evidence.ProbeStatus = resultmodel.ProbeLaunchFailed
225				reason = "blocked_check could not launch: " + probeError.Error()
226			} else if exitCode == 124 {
227				evidence.ProbeStatus = resultmodel.ProbeTimedOut
```
Three separate wrongs. Line 210 fabricates the number 125 for a probe no process ever attempted, and line 209 sets `ProbeAttempted = true` for that same never-attempted probe. Line 223 keys launch-failure on `probeError != nil` alone, which is the opposite convention from `commands.go:550` (which ORs error with 125), so the two packages disagree about what launch failure means. Line 226 re-parses 124 for timeout. Note that a signal interruption is caught earlier at line 218 by `blockedProbeInterruptionStatus`, which only accepts 129, 130, 143 (`blocked_probe.go:141`); any other `128+signal` falls through to line 223 and is labelled launch failure.

**D10 — `internal/corehelpers/commands.go:547-552`, the classification ladder ignores the booleans it already has.**
```go
547	if status == 124 {
548		code = "BLOCKED-PROBE-TIMED-OUT"
549	}
550	if status == 125 || runError != nil {
551		code, outcome, severity = "BLOCKED-PROBE-LAUNCH-FAILED", resultmodel.OutcomeFailure, resultmodel.SeverityError
552	}
```
`probeEvidence.TimedOut` and `probeEvidence.Launched` are in scope and are copied into `focusedTest` five lines later (554-555), but the decision reads the integers. Two consequences beyond the misclassification: the emitted JSON can carry `timed_out: false` next to code `BLOCKED-PROBE-TIMED-OUT` once the booleans become truthful, unless this ladder is rewritten off the same facts; and `runError != nil` at 550 swallows the typed `BlockedProbeInterruption` from `blocked_probe_unix.go:76`, so a Ctrl-C during a focused test is reported as `BLOCKED-PROBE-LAUNCH-FAILED` for a probe that demonstrably launched and ran. `handleBlockedCheck` never calls `blockedProbeInterruptionStatus`; its only caller is `next_selection.go:218`.

**D11 — `commands.go:547-549` leaves outcome and severity at the `status != 0` values.**
The timeout rung rewrites only `code`; `outcome` stays `OutcomeFindings` and `severity` stays `SeverityWarning` from line 545. Whether a timeout should stay advisory is a policy call the remediation has to make explicitly rather than inherit from statement ordering. Precedence in this ladder is textual: 125 overwrites 124.

## (b) `compareFocusedBaseline` validates only the saved baseline

**D12 — `commands.go:596`, the launch gate tests the wrong run.**
```go
596	if !baseline.Launched {
597		focusedTest.BaselineState = resultmodel.FocusedBaselineUnusable
```
`baseline.Launched` is `baselineRecord.Launched` (`internal/corehelpers/checks.go:18-22`) decoded from `do-work/working/baseline.json`, written by `handlePreflight` at `checks.go:98`. Verified by reading the whole function: between lines 579 and 616 the identifiers `focusedTest.Launched` and `focusedTest.TimedOut` never appear. The only `focusedTest` fields read are `ExitStatus` (600, 610), `CommandText` (610), `DiagnosticSHA256` (610) and `ProbeFile` (615). So the current execution's launch and timeout facts, which `handleBlockedCheck` just populated at 554-555, are written into the result and then ignored by the classifier that decides `baseline_state`.

**D13 — `commands.go:610`, the match predicate compares status numbers only.**
```go
610	if baseline.ExitStatus != 0 && baseline.ExitStatus == focusedTest.ExitStatus && strings.TrimSpace(baseline.TestCommand) == focusedTest.CommandText && baselineDiagnosticSHA256 == focusedTest.DiagnosticSHA256 {
611		focusedTest.BaselineState = resultmodel.FocusedBaselineMatchingRed
```
A current run that timed out (124), failed to launch (125) or was killed by a signal (143) matches a saved baseline carrying the same number and the same diagnostic, and is classified `matching_red`. This is reachable in practice: `runBaselineCommand` (`checks.go:29-49`) sets `launched = exitStatus != 126 && exitStatus != 127` (line 46), so a baseline whose own command timed out under an outer timeout wrapper is stored as `{"exit_status":124,"launched":true}` and sails through the line-596 gate. The launch-failure variant is also reachable: with an empty diagnostic on both sides the SHA256 values are equal by construction. `matching_red` is the state that clears the gate (see D14), so this is the load path from "the test never ran" to "advance says satisfied".

**D14 — `commands.go:600-603`, the green rung short-circuits every identity check.**
```go
600	if focusedTest.ExitStatus == 0 {
601		focusedTest.BaselineState = resultmodel.FocusedBaselineGreen
602		return resultmodel.CommandFinding{}
```
Exit 0 is `green` before command text or diagnostic identity is compared, so a probe that is a completely different command from the baseline's still reports green. Lower severity than D13 because exit 0 cannot be produced by the timeout or launch-failure paths, but it is the same class of error: state decided from one integer instead of from the observed facts.

## (c) `composeCoreGate` promotes matching-red over a failed subordinate

**D15 — `internal/lifecycleadvance/evidence_gates.go:166-176`, the unguarded promotion arm.**
```go
166	if gateID == corehelpers.CommandBlockedCheck && subordinate.FocusedTest != nil {
167		record.FocusedTest = subordinate.FocusedTest
168		switch subordinate.FocusedTest.BaselineState {
169		case resultmodel.FocusedBaselineGreen, resultmodel.FocusedBaselineMatchingRed:
170			record.State = resultmodel.AdvanceGateSatisfied
171		case resultmodel.FocusedBaselineMissing, resultmodel.FocusedBaselineUnusable, resultmodel.FocusedBaselineNewRed:
172			if record.State != resultmodel.AdvanceGateFailed {
173				record.State = resultmodel.AdvanceGateFindings
174			}
175		}
176	}
```
The demotion arm at 171-174 refuses to weaken an `AdvanceGateFailed` state. The promotion arm at 169-170 has no such guard. `boundAdvanceGateRecord` (210-236) has already mapped `OutcomeFailure` to `AdvanceGateFailed` at line 216; line 170 overwrites it with `AdvanceGateSatisfied`. Combined with D13, a `run-blocked-check` that returned `OutcomeFailure` with `BLOCKED-PROBE-LAUNCH-FAILED` becomes a satisfied gate.

**D16 — `evidence_gates.go:232` plus `evidence_gates.go:275`, the overwrite is invisible downstream.**
`boundAdvanceGateRecord` sets `Outcome: subordinate.Outcome` (232) and never revisits it, so the record ships with `State: satisfied` and `Outcome: failure` at the same time. `aggregateAdvanceGateResult` switches only on `record.State`:
```go
275		switch record.State {
276		case resultmodel.AdvanceGateFailed:
277			outcome = resultmodel.OutcomeFailure
```
So the aggregate returns `OutcomeSuccess` while line 273 still copies the `BLOCKED-PROBE-LAUNCH-FAILED` / SeverityError finding into the result. `advance` reports success and carries an error finding in the same payload. Also note `subordinate.ExitCodeOverride` (set to the raw probe status at `commands.go:576`) is never copied into the gate record, so the raw status leaves no trace in the advance exit status either.

**D17 — `evidence_gates.go:163-164`, no nil check on the handler.**
```go
163	handler := corehelpers.Handlers()[gateID]
164	subordinate := handler(executionContext, append([]string(nil), arguments...))
```
An unmapped `gateID` panics. Not part of REQ-506's stated defect, but this function is being edited and a nil-map lookup one line above a call is cheap to close.

## (d) Continuation argv emitted in the wrong channel relative to `--`

**D18 — `evidence_gates.go:239`, one template for five different argument shapes.**
```go
239	next := []string{"do-work-cli", "--format", "json", CommandAdvance, advance.RequestID, "--request-path", advance.RequestPath, "--", "<" + expected + ">"}
```
`missingAdvanceGateInput` always emits `--` and always puts the placeholder after it. Three of its five call sites need post-separator argv, two need pre-separator flags.

**D19 — `evidence_gates.go:63`, qualification asks for `--diff-range` after `--`.**
```go
63		records = append(records, missingAdvanceGateInput(advance, corehelpers.CommandQualify, "exact --diff-range <pre>..<merge>"))
```
Produces `... --request-path <path> -- <exact --diff-range <pre>..<merge>>`. `--diff-range` is parsed pre-separator (`parseAdvanceGateInputs` 118-123) and consumed at 62-66. Worse, the same branch refuses outright when a separator is present:
```go
59		if inputs.separatorSeen || len(inputs.gateArgv) > 0 || inputs.gateExitStatus != nil {
60			return irrelevantAdvanceGateInput(advance, "qualification accepts only --diff-range and the discovered request path")
```
Following the emitted continuation verbatim yields `ADVANCE-GATE-INPUT-IRRELEVANT` and `OutcomeRefused`, not progress. It also nests angle brackets. The correct template already exists at `internal/lifecycleadvance/advance_commands.go:226`: `--diff-range <pre>..<merge_hash>` with no separator.

**D20 — `evidence_gates.go:182`, the green gate asks for `--gate-arg` after `--`.**
```go
182		return missingAdvanceGateInput(advance, gateevidence.CommandCheckGreenGate, "one --gate-arg per canonical repository-gate argv token")
```
`--gate-arg` is a repeated pre-separator flag (parsed 124-126), and `parseAdvanceGateInputs` returns at the first `--` (104-108). Anything written after the separator lands in `phaseArgv` and `inputs.gateArgv` stays empty, so `composeGreenGate` emits the identical needs-input record again on the next invocation — a loop with no exit. In the `test-gate` branch the stray token is additionally handed to `handleBlockedCheck`, which rejects it at `commands.go:528` (`unknown option`) as a `usageResult`, i.e. `OutcomeFailure`. Correct templates exist at `advance_commands.go:218` and `:234`: `--gate-arg <canonical-gate-argv-token> -- <phase argv>`.

**D21 — `evidence_gates.go:53`, correct placement, nonsense token.**
`expected` is `"resolved test argv after --; use a bare -- when no test command exists"` and D18 wraps it in angle brackets, producing `<resolved test argv after --; use a bare -- when no test command exists>` as a single argv token. The channel is right; the token is a sentence pretending to be a placeholder. Same shape at line 44 (`<estimator argv after -->`) and line 79 (`<run-blocked-check argv after -->`).

**D22 — `evidence_gates.go:287-303`, the one correct builder, used once.**
`advanceGateVerificationArgv` emits diff-range (289-291), gate-exit-status (292-294), repeated `--gate-arg` (295-297) all before `--`, then `phaseArgv` after (298-301). Its only caller is line 199 in `composeGreenGate`. This is the shape D19 and D20 should be producing, and it is already in the file.

**D23 — `commands.go:570-572`, the probe finding's retry argv is weaker than the run that failed.**
```go
572		[]string{"do-work-cli", CommandBlockedCheck, "--probe-file", probeFile}, []string{"do-work-cli", "--format", "json", CommandBlockedCheck, "--probe-file", probeFile})}
```
Drops `--timeout-seconds`, `--baseline-json` and `--baseline-failures`. The prescribed retry runs at the 30-second default (`commands.go:496`) and ends at `baseline_state: not_compared`, so following it can never reproduce the failure it is offered for. It also re-enters the helper directly instead of returning to `advance`. `boundAdvanceFinding` (`evidence_gates.go:258-266`) only prepends IDs and paths, so this argv survives verbatim into the advance result.

**D24 — `commands.go:583, 587, 592, 598, 607, 615`, six error-severity findings with no argv at all.**
Every `compareFocusedBaseline` finding passes `nil, nil` for next and verify. `helperFinding` (350-360) delegates to `findingSpecificCommands` (362-422), whose table has no case for any `FOCUSED-*` code and falls through to `return nil, nil` at 421. `FOCUSED-NEW-RED`, `FOCUSED-BASELINE-UNREADABLE`, `FOCUSED-BASELINE-INVALID`, `FOCUSED-BASELINE-NOT-LAUNCHED` and `FOCUSED-BASELINE-FAILURE-EVIDENCE-MISSING` all ship as SeverityError with empty `NextArgv` and empty `VerificationArgv`. The heavy matrix test at `internal/corehelpers/commands_test.go:288-294` asserts every non-info finding has both, but its 17-command table never reaches this path.

# 2. Type surface

Everything needed to carry observed launch, timeout and error facts already exists. No new public representation is required.

**T1 — `BlockedProbeEvidence`, `ROOT/internal/nextselection/blocked_probe.go:20-26`.**
```go
type BlockedProbeEvidence struct {
	ExitStatus       int
	Launched         bool
	TimedOut         bool
	Diagnostic       string
	DiagnosticSHA256 string
}
```
No json tags; never marshalled; copied field-by-field into `resultmodel.FocusedTestResult` at `commands.go:553-557`. `Launched` and `TimedOut` are the correct destinations for the observed facts — the change at D6 is what feeds them, not the fields themselves. Because it is never marshalled, adding a field here (for example an interruption or termination-cause discriminator) changes no wire format. It is exported, so any addition is still an API addition for the one external consumer at `commands.go:541`.

**T2 — `BlockedProbeInterruption`, `blocked_probe.go:62-72`, with `Error() string` (66) and `InterruptionExitStatus() int` (70).**
The existing typed error for the signal case, constructed only at `blocked_probe_unix.go:76`, consumed structurally through the anonymous `interface{ error; InterruptionExitStatus() int }` at `blocked_probe.go:132-136`. This is the precedent to follow for any other outcome that needs to be distinguishable without a magic number: a typed error already flows end to end through the existing `(int, error)` signature.

**T3 — `resultmodel.FocusedTestResult`, `ROOT/internal/resultmodel/result_model.go:356-367`.** The public wire contract.
```go
	ProbeFile        string                   `json:"probe_file"`
	ExitStatus       int                      `json:"exit_status"`
	Launched         bool                     `json:"launched"`
	TimedOut         bool                     `json:"timed_out"`
	Diagnostic       string                   `json:"diagnostic"`
	DiagnosticSHA256 string                   `json:"diagnostic_sha256"`
	BaselineState    FocusedTestBaselineState `json:"baseline_state"`
	BaselineStatus   int                      `json:"baseline_exit_status"`
	BaselineLaunched bool                     `json:"baseline_launched"`
	CommandText      string                   `json:"command_text"`
```
`launched` and `timed_out` are already published and already reach `advance` via `AdvanceGateRecord.FocusedTest` (`result_model.go:456`, `json:"focused_test,omitempty"`). `baseline_exit_status` and `baseline_launched` describe the saved side; the current side has `exit_status`, `launched`, `timed_out`. Every fact D12 needs is on this struct already — it just is not read.

**T4 — `FocusedTestBaselineState`, `result_model.go:342-351`.** Six values, all already defined: `not_compared`, `baseline_missing`, `baseline_unusable`, `matching_red`, `new_red`, `green`. The shipped action prose at `/home/user/skill-do-work/skills/do-work/actions/work.md:339` already promises a six-way distinction including timeout and launch failure, so the remediation must express those through existing states (most naturally `baseline_unusable`, which already means "cannot exclude anything") rather than by adding enum members.

**T5 — `SelectionProbeStatus`, `result_model.go:197-206`.** `not_applicable`, `missing`, `succeeded`, `failed`, `timed_out`, `launch_failed`. Carried on `SelectionRecord` (219-221) and `SelectionExclusion` (252-254) as `ProbeStatus` / `ProbeAttempted` / `ProbeExitCode`. `timed_out` and `launch_failed` already exist; D9's job is to set them from facts, not to invent states.

**T6 — `AdvanceGateRecord`, `result_model.go:444-458`.** `State AdvanceGateState` (`needs_input` / `satisfied` / `findings` / `failed`, 423-430), `Outcome CommandOutcome`, `Findings`, `Changes`, `OutputLines`, `NextArgv`, `VerificationArgv`, `FocusedTest *FocusedTestResult`, `GreenGate *GateEvidenceResult`. The record can already express "failed with a focused test attached"; D15 is the only reason it does not.

**T7 — `baselineRecord`, `ROOT/internal/corehelpers/checks.go:18-22`.**
```go
type baselineRecord struct {
	TestCommand string `json:"test_command"`
	ExitStatus  int    `json:"exit_status"`
	Launched    bool   `json:"launched"`
}
```
Unexported, written at `checks.go:98`, read at `commands.go:589`. Its on-disk JSON (`do-work/working/baseline.json`) is a shipped format that fixtures and action prose already write by hand, so its three keys must keep their names and meanings.

**T8 — the internal seam that can change freely.** `runOwnedProbe` is unexported and defined once per build tag (`blocked_probe_unix.go:15`, `blocked_probe_windows.go:11`). Its return type is the natural place to carry observed facts — returning `BlockedProbeEvidence`-shaped data or a small unexported outcome struct costs nothing publicly, provided both variants change together and `probeExitStatus` (unix 116-130) stops being the only classifier of `command.Wait()`'s error.

# 3. Caller compatibility

Exported symbols in `nextselection` and who depends on them. Signatures marked **frozen** must not change.

**C1 — `RunBlockedProbeAtRoot(repositoryRoot string, probeBytes []byte, timeoutSeconds int) (int, error)`, `blocked_probe.go:86`. Frozen.**
Production callers: `ROOT/internal/lifecycleadvance/queue_commands.go:256` and `ROOT/internal/nextselection/next_commands.go:33`. Both wrap it as the `probeRunner` closure of type `func(probeBytes []byte, timeoutSeconds int) (int, error)`, which is the contract `Select` takes. Also `blocked_probe_test.go:145` (`runBlockedProbeFixture`). Keeping the raw-status-plus-error surface is required; giving selection richer facts means either a second evidence-shaped runner alongside it or a typed error, not a changed signature here.

**C2 — `RunBlockedProbe(probeBytes []byte, timeoutSeconds int) (int, error)`, `blocked_probe.go:77`. Frozen.**
No production caller. Tests only: `blocked_probe_test.go:19, 46, 69` and `ROOT/internal/nextselection/next_commands_test.go:185`. It is a cwd-based wrapper over C1 whose only failure path is `os.Getwd()` returning 125 at line 80.

**C3 — `RunBlockedProbeEvidenceAtRoot(repositoryRoot string, probeBytes []byte, timeoutSeconds int) (BlockedProbeEvidence, error)`, `blocked_probe.go:93`.**
One external caller: `ROOT/internal/corehelpers/commands.go:541`. One test caller: `blocked_probe_test.go:28`. This is the evidence-preserving seam and the natural place for the fix to surface truthful facts.

**C4 — `BlockedProbeDiagnosticIdentity(diagnostic, repositoryRoot string) (string, string)`, `blocked_probe.go:115`.**
Called from `blocked_probe.go:105` and from `commands.go:609` inside `compareFocusedBaseline`, where only the second return value is used. Both sides of the SHA256 comparison at `commands.go:610` go through it, so any change to its normalization silently invalidates every stored `baseline-failures.txt` identity.

**C5 — `BlockedProbeTimeoutStatus = 124` and `BlockedProbeLaunchStatus = 125`, `blocked_probe.go:15-16`.**
`BlockedProbeTimeoutStatus` is referenced outside its own file only at `blocked_probe_test.go:47`. `BlockedProbeLaunchStatus` is referenced only inside the package (`blocked_probe.go:80, 95, 98, 101, 107`; `blocked_probe_unix.go:21, 40, 52, 129`; `blocked_probe_windows.go:12`). `corehelpers/commands.go:547, 550` and `next_selection.go:210, 226` use the bare literals 124 and 125 rather than the constants — those literals are part of what has to go. The constants themselves must survive as the values the process still exits with (`commands.go:576` `ExitCodeOverride: status`, asserted at `evidence_gates_test.go:226-227`).

**C6 — `BlockedProbeInterruption` and its `InterruptionExitStatus() int`, `blocked_probe.go:62-72`.**
Constructed once (`blocked_probe_unix.go:76`), consumed structurally at `blocked_probe.go:132-136` via `blockedProbeInterruptionStatus`, whose only caller is `next_selection.go:218`. The structural interface means the concrete type is not depended on by name outside the package, but `errors.As` against `interface{ InterruptionExitStatus() int }` is pinned by `blocked_probe_test.go:130-133`.

**C7 — `blockedProbeInterruptionStatus(err error) (int, bool)`, `blocked_probe.go:131-141`.** Unexported. Line 140 accepts only 129, 130, 143 as interruptions. Any other `128+signal` falls through and is misclassified by D9.

**C8 — `corehelpers.Handlers()`, `commands.go:49`, and the 21 `Command*` constants, `commands.go:25-47`.**
`evidence_gates.go:163` looks handlers up by name. `commands_test.go:15-25` asserts the map has exactly 21 entries, all non-nil. Adding or renaming a public subcommand breaks that test by design.

**C9 — JSON field names.** `focused_test`, `exit_status`, `launched`, `timed_out`, `baseline_state`, `baseline_exit_status`, `baseline_launched`, `command_text`, `probe_status`, `probe_attempted`, `probe_exit_code`, `gate_records`, `state`, `outcome`. These are read by shipped action prose (`/home/user/skill-do-work/skills/do-work/actions/work.md:251` and `:339` both describe `launched: false` and the six-way focused-test classification) and by the heavy matrix tests. Enum string values (`matching_red`, `launch_failed`, `satisfied`, …) are equally load-bearing.

# 4. Test conventions a new regression test must follow

Two possible homes, with different rules.

## Home A — `ROOT/internal/nextselection/blocked_probe_test.go` (probe-level, in-process)

Build tag, line 1, mandatory for anything touching the unix runner:
```go
//go:build unix
```
No heavy-lane guard in this file; tests run in the default lane. Isolation is `t.TempDir()`, no `t.Cleanup`. The existing wrapper:
```go
func runBlockedProbeFixture(repositoryRoot string, probeBytes []byte, timeoutSeconds int) (int, error) {
	return RunBlockedProbeAtRoot(repositoryRoot, probeBytes, timeoutSeconds)
}
```
Assertion shape is one compound `if` and one `t.Fatalf` with `%#v` on the whole evidence value, as at line 32-34:
```go
	if evidence.ExitStatus != 29 || !evidence.Launched || evidence.TimedOut {
		t.Fatalf("evidence=%#v", evidence)
	}
```
Signal tests install their own guard first (lines 92-94) so an early delivery is safe, and poll descendants with `syscall.Kill(pid, 0)` on a deadline.

## Home B — `ROOT/internal/lifecycleadvance/evidence_gates_test.go` (end-to-end, out of process)

**No build tags anywhere in `internal/lifecycleadvance`.** `evidence_gates_test.go:1` is `package lifecycleadvance`. Do not add one.

**No heavy-lane guard.** The repo idiom `if testing.Short() || os.Getenv("DO_WORK_HEAVY_TESTS") != "1" { t.Skip(...) }` is used in `cmd/do-work-cli/gate_evidence_integration_test.go:12`, `internal/doctor/doctor_commands_test.go:21` and `internal/nextselection/next_commands_test.go:49`, and deliberately not here. Do not add one.

**No `t.Cleanup` in this package.** Isolation is `repositoryRoot := t.TempDir()` inside each `t.Run` closure (lines 67, 182, 230). The built binary is intentionally never removed.

**The CLI seam is a binary built once per package run**, `advance_commands_test.go:379-397`, backed by `advanceCLIBinaryOnce sync.Once` / `advanceCLIBinaryPath` / `advanceCLIBinaryErr` (19-23):
```go
func advanceCLIBinary(t *testing.T) string {
	t.Helper()
	advanceCLIBinaryOnce.Do(func() {
		temporaryDirectory, err := os.MkdirTemp("", "advance-cli-test-*")
		...
		command := exec.Command("go", "build", "-o", advanceCLIBinaryPath, "../../cmd/do-work-cli")
```

**The standard driver**, `evidence_gates_test.go:252-271`:
```go
func runAdvanceGateJSON(t *testing.T, repositoryRoot string, arguments ...string) (resultmodel.CommandResult, int) {
	t.Helper()
	commandArguments := []string{"--repo-root", repositoryRoot, "--format", "json", "advance"}
	commandArguments = append(commandArguments, arguments...)
	command := exec.Command(advanceCLIBinary(t), commandArguments...)
	output, runError := command.CombinedOutput()
	...
}
```
Bypass it only when the helper cannot express the need. The precedent is `evidence_gates_test.go:235-238`, which needs a modified environment:
```go
	command := exec.Command(advanceCLIBinary(t), commandArguments...)
	if test.emptyPath {
		command.Env = append(os.Environ(), "PATH=")
	}
```
`PATH=` is the established way to force a real 125 launch failure end to end.

**Fixture helpers**, all in `ROOT/internal/lifecycleadvance/advance_commands_test.go`, same package:
```go
func writeAdvanceRequest(t *testing.T, repositoryRoot, treeSection, requestID, status, frontmatter, body string) string {   // 341
	requestPath := filepath.ToSlash(filepath.Join("do-work", treeSection, requestID+"-fixture.md"))
	contents := "---\nid: " + requestID + "\ntitle: Fixture " + requestID + "\nstatus: " + status + "\n" + frontmatter + "---\n\n" + body
	writeAdvanceFile(t, repositoryRoot, requestPath, contents)
	return requestPath
}
func writeAdvanceFile(t *testing.T, repositoryRoot, relativePath, contents string)                                          // 349
func routeCBodyThrough(lastSection string) string                                                                          // 315
func runAdvanceGit(t *testing.T, repositoryRoot string, arguments ...string) []byte                                         // 439
```
The return of `writeAdvanceRequest` is the slash-form relative path, which is exactly what `--request-path` expects.

**Reaching the test-gate phase.** Route A with an inline body is the established shortcut (lines 183 and 231 use the identical literal):
```go
	body := "## Triage\n\nRoute A.\n\n## Plan\n\nPlanning not required.\n\n## Implementation Summary\n\n- `owned.go` (modified)\n\n## Qualification\n\nPassed.\n"
	writeAdvanceRequest(t, repositoryRoot, "working", "REQ-713", "claimed", "route: A\nestimate:\n  p50_active_minutes: 5\n", body)
```
Route A skips the `scope-drift` gate (`evidence_gates.go:75`), so `records[0]` is the probe gate. `routeCBodyThrough("Qualification")` is the Route C equivalent when scope-drift must also run.

**Planting a canonical green record.** Two routes, both in use. Out of band via the public verb (lines 87-91):
```go
	recordCommand := exec.Command(advanceCLIBinary(t), "--repo-root", repositoryRoot, "--format", "json", "record-green-gate", "--gate-exit-status", "0", "--", gateArgv[0], gateArgv[1], gateArgv[2])
```
In band via `advance` itself, passing `--gate-exit-status 0` alongside `--gate-arg` (line 157), which routes `composeGreenGate` to `record-green-gate` (`evidence_gates.go:187-191`). Both need a git repo with one commit: `init -q`, `config user.name "Advance Gate Test"`, `config user.email "advance@example.invalid"`, `add .`, `commit -qm "fixture"` (lines 81-85, 117-121, 145-149).

**Baseline fixture files** go at `do-work/working/baseline.json` and `do-work/working/baseline-failures.txt` (lines 186-189). `evidence_gates.go:82-87` appends those two default paths automatically when the phase argv names neither.

**Assertion helpers**, local to `evidence_gates_test.go`: `findAdvanceGate(result, gateID) *resultmodel.AdvanceGateRecord` (273), `gateHasFinding(gate, code) bool` (285), `hasAdvanceResultFinding(result, code) bool` (294).

**House assertion shape** — one compound `if`, one `t.Fatalf` dumping the whole record and result (line 192-194):
```go
	if focusedGate == nil || focusedGate.FocusedTest == nil || focusedGate.FocusedTest.BaselineState != test.wantState || focusedGate.State != test.wantGateState {
		t.Fatalf("focused gate=%#v result=%#v", focusedGate, result)
	}
```
Table tests use an anonymous struct slice with a leading `name` field and `t.Run(test.name, func(t *testing.T) { … })`. Probes stay POSIX `sh` (`exit 0`, `echo same; exit 17`, `sleep 5`), never bash-isms — that is what lets this package skip a unix build tag.

# 5. Existing controls to preserve

Each of these currently passes and pins a contract the remediation must not break.

**`ROOT/internal/nextselection/blocked_probe_test.go`**

- `TestBlockedProbePreservesRawStatus` (18) — `exit 37` returns `status == 37, err == nil`. A probe's own status is never folded into a 0-4 envelope.
- `TestBlockedProbeEvidenceBoundsAndNormalizesDiagnostics` (25) — status 29 with `Launched == true` and `TimedOut == false` (line 32); diagnostic is CRLF-normalized, root-redacted to `<repo-root>`, truncated at `blockedProbeDiagnosticLimit` with the `[diagnostic truncated]` marker, and hashed. Line 32 is the only assertion in the repo that touches `Launched`/`TimedOut`, and only on the happy path.
- `TestBlockedProbeTimeoutKillsDescendantGroup` (42) — 1s timeout on a probe whose child traps SIGTERM returns `status == BlockedProbeTimeoutStatus` **with `err == nil`** (line 47), and the descendant is dead within 2s. Two contracts: the 124 value, and that timeout is not signalled through the error channel today. If the fix starts returning a non-nil error on timeout, this assertion changes and the change must be deliberate.
- `TestBlockedProbeCleansBackgroundDescendantAfterLeaderExits` (65) — leader exits 0, orphaned `sleep 30` is reaped within 2s.
- `TestBlockedProbeInterruptionIsTypedAndReapsDescendants` (88) — SIGINT to the test process yields `status == 130` and an error satisfying `errors.As(..., &interface{ InterruptionExitStatus() int })` returning 130. The typed-error precedent (T2) must survive.

**`ROOT/internal/corehelpers/commands_test.go`**

- `TestEveryRemainingUtilityHasOneHandler` (15) — `len(Handlers()) == 21`, every `Command*` constant maps to a non-nil handler.
- `TestBlockedCheckReturnsTypedBoundedBaselineComparison` (165) — matching red (`echo 'same failure'; exit 23` against a baseline recording 23) gives `FocusedBaselineMatchingRed`, non-empty `DiagnosticSHA256`, `FOCUSED-BASELINE-MATCH` present and `FOCUSED-NEW-RED` absent; then `launched:false` gives `FocusedBaselineUnusable` + `FOCUSED-BASELINE-NOT-LAUNCHED`. This is the only helper-level test of the classification and it exercises neither 124, nor 125, nor `runError != nil`.
- `TestNonInformationalFindingsReceiveCommandSpecificActions` (43) — four codes with exact argv, and no family-wide placeholder.
- `TestDryRunSurfacesDoNotMutateBaselineDownloadOrTimestamps` (127) — `handlePreflight --dry-run` emits `PREFLIGHT-BASELINE-DRY-RUN` and writes nothing.
- `TestAllSeventeenPublicCommandsRunInTextAndJSONWithStableStatusAndNoDryRunEffects` (202, heavy-only) — for all 17 commands: exit status never above 1, JSON decodes with non-nil `Findings`/`Changes`, **every non-info finding has non-empty `NextArgv` and `VerificationArgv`** (288-294), no family-wide action strings, text and JSON exit statuses agree, every JSON finding and change appears in the text output in order, and `git status --porcelain=v1 --untracked-files=all` is byte-identical before and after. D24's empty-argv findings would fail 288-294 if that table ever reached them; whatever fills them must not reintroduce a family-wide string from the forbidden set at 295-300.

**`ROOT/internal/lifecycleadvance/evidence_gates_test.go`**

- `TestAdvanceExecutesEstimateGateAtPublicCLISeam` (14) — one record, `gate_id=estimate-p50`, `satisfied`, `advance_executed`, non-empty `output_lines`.
- `TestAdvanceEvidenceGatesReturnTypedMissingInputs` (52) — four phases, each returning `needs_input` plus `ADVANCE-GATE-INPUT-REQUIRED` and a nonzero exit status when its judgment-owned input is absent. Fixing D19 and D20 changes the argv inside those records; this test asserts the state and code, not the argv, so it should keep passing. A new assertion on argv placement belongs here or beside it.
- `TestAdvanceExecutesPreflightAndProjectsGreenEvidence` (78) — exactly two records, `preflight` satisfied with non-empty `Changes` then `green-gate` satisfied with `existing_evidence` and `GreenGate.Matches == true`; `do-work/working/baseline.json` is written.
- `TestAdvanceQualificationUsesExactRangeAndRunsScopeDrift` (113) — real `--diff-range` gives `qualify` (`merged_range`) then `scope-drift`; a bogus range gives nonzero status, `findings`, `QUALIFY-DIFF-RANGE-INVALID`.
- `TestAdvanceGreenGateMissRequiresDirectRunThenRecordsIt` (142) — `needs_input` with `len(NextArgv) == 3` and `ADVANCE-GREEN-GATE-DIRECT-RUN-REQUIRED`, then `--gate-exit-status 0` flips to `satisfied` with `GateEvidenceRecorded` and exactly one change. The `len(NextArgv) == 3` assertion pins `composeGreenGate`'s raw-gate-argv `NextArgv` (line 198), which is a different field from the `missingAdvanceGateInput` template in D18 — do not conflate them.
- `TestAdvanceFocusedTestGateClassifiesBaselineStates` (165) — five rows: green/satisfied, matching red/satisfied, different red status/new_red+findings, new red/new_red+findings, unlaunched baseline/unusable+findings. The fix must keep all five. Note that the two `AdvanceGateSatisfied` rows are the promotion arm at `evidence_gates.go:169-170` doing its intended job over an `OutcomeFindings` subordinate; D15 is only about the unguarded case where the subordinate failed.
- `TestAdvanceGateInputsFailClosedAndNeverInterpolateHostileTokens` (199) — a `--request-path` disagreeing with the discovered request gives `ADVANCE-EVIDENCE-MISMATCH`; a `$(touch …)` token makes `GateRecords[0].State == AdvanceGateFailed` and the marker file must not exist.
- `TestAdvanceFocusedTestGateDistinguishesTimeoutAndLaunchFailure` (217) — `sleep 5` with `--timeout-seconds 1` gives `FocusedTest.ExitStatus == 124` and gate `findings`; `PATH=` gives `125` and gate `failed`. This is the closest existing control to REQ-506 and the one to extend: it asserts only `ExitStatus` and gate state, never `Launched` or `TimedOut`, and both rows have no baseline planted, so both land on `FocusedBaselineMissing` and never exercise the promotion arm.

**Gaps with no control at all** (nothing in any `_test.go` greps these): `FOCUSED-BASELINE-MISSING`, `FOCUSED-BASELINE-INVALID`, `FOCUSED-BASELINE-UNREADABLE`, `FOCUSED-BASELINE-FAILURE-EVIDENCE-MISSING`; the three input-validation guards at `blocked_probe.go:94-102`; the windows stub; a probe that exits 124 or 125 under its own power; and a subordinate `OutcomeFailure` meeting `matching_red` at `evidence_gates.go:169-170`. The last one is the regression test REQ-506 most needs.
