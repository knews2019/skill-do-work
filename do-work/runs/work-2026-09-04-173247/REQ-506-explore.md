# REQ-506 Exploration

## Package and API map

- `lifecycleadvance` may import both `corehelpers` and `gateevidence` without a cycle. `go list -deps ./internal/corehelpers` and `go list -deps ./internal/gateevidence` contain no `internal/lifecycleadvance`; both packages terminate through shared lower layers such as `commandruntime`, `resultmodel`, `repositorymodel`, `requeststate`, and `nextselection`.
- Reuse command handlers through `corehelpers.Handlers()[corehelpers.CommandEstimateP50|CommandPreflight|CommandQualify|CommandScopeDrift|CommandBlockedCheck]` and `gateevidence.Handlers()[gateevidence.CommandCheckGreenGate|CommandRecordGreenGate]`. Do not export the private `handle*` functions or add one-line delegates. This is the same composition pattern already used by `lifecycleadvance/recovery_commands.go` for finalization.
- `preflight`, `qualify`, and `scope-drift` already return typed `CommandResult` findings/changes. Green-gate check/record already returns typed `GateEvidenceResult`. `estimate-p50` is reusable but exposes its calculated block only as `ExactTextOutput`, which JSON omits; advance must project that value into its new typed gate record (or replace it with an exported typed estimator result) before returning.
- `run-blocked-check` is not presently a complete focused-test/baseline API: it accepts `--probe-file` and `--timeout-seconds`, returns status plus `BLOCKED-PROBE-*`, and discards probe stdout/stderr. `baselineRecord` is private to `corehelpers/checks.go`. Therefore status, timeout, launch failure, and `launched: false` are mechanical now, but same-test/file/failure-mode comparison cannot be automated honestly until the owned-probe API returns bounded diagnostics and corehelpers owns a typed baseline comparison.

## Exact implementation scope

Touch these files for the complete requested behavior:

- `skills/do-work/actions/work.md`
- `skills/do-work/actions/work-reference.md`
- `_dev/tests/contracts/core-checks.sh`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates.go` (new)
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go`
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_test.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_unix.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/blocked_probe_windows.go`

Do not touch `_dev/tests/contract-regressions.sh`: it is now a 76-line dispatcher, not the sentence-predicate owner. `core-checks.sh` is the correct owner for retaining the anti-rationalization table/Ratchet and rejecting reintroduced recipes. Do not modify `gateevidence`; its public handler map and result are sufficient. The four `nextselection` paths are required only because the REQ promises mechanical failure-identity comparison; dropping them requires explicitly keeping that comparison as agent judgment and narrowing the acceptance claim.

## Stale plan assumptions and RED cases

The plan is correct about handler-map composition and wrong that `run-blocked-check` can already compare failure identity. It also treats estimator output as structured although it is text-only, and it should not name `_dev/tests/contract-regressions.sh` as the predicate owner.

Add public CLI RED tests before implementation for:

1. Estimate phase: `advance REQ-NNN -- <estimator argv>` currently falls into queue-mode usage instead of executing the estimator and returning a REQ/path-bound gate record with persisted-value evidence.
2. Preflight phase: advance with tokenized test argv currently projects `preflight` instead of executing it and reporting its exact changes/findings for the selected REQ.
3. Qualification and scope phase: advance with an exact valid merge range currently returns subordinate argv instead of executing `qualify` then `scope-drift`; invalid/unresolvable ranges must produce a bound typed finding.
4. Test gate: a green probe, matching red baseline, new/different failure, timeout, launch failure, and `baseline.json` with `launched: false`; the last must refuse baseline exclusion. Assert bounded diagnostic identity, not only exit status.
5. Identity/grammar: irrelevant phase inputs, a mismatched request path, missing required inputs, and hostile tokens stay tokenized and never execute through interpolation.
6. Green evidence: exact recorded match is projected as existing evidence; a miss asks for the direct gate without claiming it ran; record/check failure stays typed and bound to the current REQ.

## Boundary: mechanics versus judgment

Move now: phase detection; exact REQ/path binding; frozen-estimate reuse; estimator execution after supplied signals; preflight execution and baseline artifact observation; qualification/scope-drift execution with exact range; owned focused-test execution; baseline launch/status/diagnostic comparison; green-evidence check/record composition; typed outcomes/findings/changes/continuation.

Keep in prose: extracting estimator signals and resolving the project test/canonical-gate argv; deciding whether qualification warnings are acceptable; substantive requirement tracing and the anti-rationalization table; interpreting whether two normalized diagnostics are semantically the same when the machine cannot prove it; TDD RED-before-GREEN validity; retry/failure classification; canonical-gate fingerprint attribution, repair/deferral choice, and isolated-base diagnosis; heavy-lane planning/hold judgment; the Finding-Closure Ratchet. `advance` may expose evidence for these decisions but must not decide them.
