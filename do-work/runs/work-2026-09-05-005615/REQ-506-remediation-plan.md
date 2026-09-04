# REQ-506 — Single-remediation plan

Repair the focused-test false success first. Include F02's continuation repair because it lives in the same `lifecycleadvance/evidence_gates.go` owner and is required by the original promise that an agent can follow `advance` evidence. Keep F03's stale reference and F04's estimate condition as report-only findings; neither is necessary to correct this execution boundary, and they would expand the concrete remediation owner set.

This is read-only preparation. No source or lifecycle changes and no new tests were run. The original review is `REQ-506-review.md` in this run directory. At dispatch, the orchestrator must re-resolve the exact claimed REQ/path and revision, settle its new baseline, record this one attempt, and give the builder an isolated checkout.

## Exact builder boundary

Eight paths, all already in REQ-506's original declared scope:

1. `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go`
2. `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go`
3. `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go`
4. `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go`
5. `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go`
6. `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go`
7. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates.go`
8. `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go`

No new result-model field or enum is needed: existing `exit_status`, `launched`, `timed_out`, `baseline_state`, record `state`/`outcome`, findings, and argv can express the corrected result. Preserve the raw-status compatibility wrappers used by queue selection. Do not change `next_selection.go`, lifecycle classification, finalization, publication/release, action prose, or the estimator. Do not refactor process-group teardown. If implementation proves another path necessary, report the exact reason before expanding this boundary.

## Task 1 — Prove the false success and continuation failure before source changes

Add focused regressions at the existing public CLI seam in `evidence_gates_test.go`, using a correctly claimed fixture with an Implementation Summary and Qualification, plus a valid directly executed canonical-gate zero so a missing green record cannot mask the focused verdict. Suggested names describe the contract rather than the implementation:

- `TestAdvanceFocusedGateNeverClearsFailedExecutionAgainstMatchingBaseline`: use a valid launched baseline with probe `exit 125`, status 125, and empty diagnostic evidence. Run the built CLI with PATH containing only a fixture-owned `git` executable/symlink, so Git evidence still works and the real `sh` launch fails. Assert current `launched:false`, the launch failure finding, non-satisfied focused gate, and nonzero aggregate advance status. The unfixed revision reports aggregate success and `satisfied` despite subordinate `failure`.
- A subcase with direct baseline `/bin/sleep 2; exit 124` actually completing at status 124, then the identical focused command with one-second timeout. Assert `timed_out:true` and that timeout cannot become `matching_red`/`satisfied` or aggregate success. Use the existing platform gating and cleanup conventions; no externally timed signal is needed for this replay.
- `TestBlockedProbeEvidencePreservesOrdinaryReservedExitValues` in the probe tests: ordinary immediate `exit 124` and `exit 125` both launched, neither timed out, and their raw exit values remain intact. At the public focused seam, valid matching baselines for these ordinary child exits must still clear; they are not infrastructure failures merely because of their integer values.
- `TestAdvanceMissingInputContinuationsPreserveArgumentChannels`: request qualification with exact request path but no range, substitute a real resolvable range into the emitted continuation without moving flags or separators, and execute it. Do the analogous round trip for missing canonical-gate tokens while preflight/focused phase tokens are already present. The unfixed qualification continuation returns `ADVANCE-GATE-INPUT-IRRELEVANT`.

Run these tests before edits and retain the real failing assertions. Keep or extend the existing green, ordinary matching-red, new-red/status-mismatch, missing/malformed/unlaunched-baseline, hostile-token, and green-record controls. The fix must not turn every baseline red into a regression or every high exit status into an infrastructure failure.

## Task 2 — Carry observed execution facts through the existing owners

The three owners currently disagree:

- The platform `runOwnedProbe` knows whether `Start` happened and which completion/timer/signal branch won, but returns only `(int, error)`.
- `RunBlockedProbeEvidenceAtRoot` reconstructs launch and timeout from statuses 125/124, losing that distinction.
- `handleBlockedCheck` again classifies those integers, `compareFocusedBaseline` validates only the saved baseline's launch state, and `composeCoreGate` then promotes any matching-red result to satisfied even when the subordinate failed.

Return observed lifecycle facts from the private platform runner, preferably through the existing evidence type rather than a second public representation. Set launch facts at the actual launch boundary and timeout only in the timer branch; keep error/interruption identity and status intact. Populate normalized diagnostic data in the existing common wrapper. The Windows implementation remains a truthful unsupported launch, and `RunBlockedProbeAtRoot`/`RunBlockedProbe` keep returning the same raw status/error surface to existing queue consumers.

In `handleBlockedCheck`, classify execution from those facts plus the actual runner error, not reserved integer equality. Only a successfully launched, normally finished execution with no runner error and no timeout may be compared for green/matching-red exclusion. The existing `not_compared` state and execution finding can express an execution that never became eligible; do not invent a baseline value that suggests the saved record itself was invalid. Preserve saved-baseline parsing and bounded diagnostic equality for executions that are eligible.

In `composeCoreGate`, preserve failed subordinate authority. Satisfy the gate only for a valid current execution whose baseline state is green or matching red; a matching-red string cannot erase a failed outcome. Ordinary matching red may retain its subordinate `findings` outcome while its gate state is satisfied, as before. Preserve exact identity, affected paths, provenance, diagnostics, and the mandatory separate canonical-gate record.

## Task 3 — Restore usable continuation argv and validate the integrated behavior

Within `evidence_gates.go`, construct missing-input continuations in their owning channel: qualification `--diff-range` and repeated `--gate-arg` flags belong before `--`; estimator/preflight/probe arguments belong after it. Reuse the known phase-specific shape and preserve already supplied valid inputs, exact REQ/path, gate argv, timeout, and paired baseline paths. Placeholder substitution must be sufficient; callers must not repair syntax or know the leaf command. A continued green-record check must not manufacture a new direct-run zero attestation.

For subordinate retry/verification commands that would re-enter an evidence helper, translate them back to the same request-bound advance phase with its accepted inputs. Preserve actual diagnostic commands such as `git diff`; these inspect a finding and do not bypass a gate. Preserve the deliberate direct canonical-gate `next_argv` plus its advance verification path. Do not blanket-replace every subordinate command or introduce a second handler registry.

Run the new focused public regressions and all four existing owner packages: `go test -count=1 ./internal/lifecycleadvance ./internal/corehelpers ./internal/nextselection ./internal/resultmodel`. Exercise the existing interruption/process cleanup coverage, and compile the changed platform signature with `GOOS=windows GOARCH=amd64 go test -c ./internal/nextselection -o <owned temporary path>`. Run gofmt and the CLI prime's applicable vet/race checks. Let the orchestrator perform the canonical post-merge gate and select renewed heavy evidence for the actual integrated range; original saved green results cannot prove changed source.

Hand back the exact eight-path manifest, RED/GREEN outputs, preserved ordinary matching-red controls, argv round trips, and unresolved report-only findings. Independent re-review should verify that real failed launch and real timeout remain non-satisfied with a valid canonical green record present, and that ordinary child exits 124/125 are neither mislabeled nor blanket-rejected. No second remediation loop is authorized by this plan.
