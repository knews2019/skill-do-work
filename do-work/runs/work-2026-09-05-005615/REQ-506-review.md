# Review: REQ-506

**Request changes — a failed or timed-out focused test can be reported as satisfied and make `advance` exit successfully.** The gate composition works on its ordinary paths, but its baseline override is not safe enough to authorize completion evidence.

Route C | saved implementation `06367337dd82d97416e0d9d37872cc35b56ae7bc` | reviewed range `24ed2fdda549a0759cdc571562c9b782bfeb6251..06367337dd82d97416e0d9d37872cc35b56ae7bc`.

## What's built

`advance` composes estimate, preflight, qualification/scope, focused-test comparison, and green-gate handlers into ordered request/path-bound records. Four procedural sections were replaced by a common consumer loop; qualification skepticism, finding closure, TDD, retries, and repository-gate attribution remain judgment. All 16 substantive changed paths match the Implementation Summary and declared scope. The two declared `corehelpers/checks*` paths were left unchanged because their existing handlers are reused, as Qualification already explains.

## Findings

**Important:**

- **F01 — Failed execution can clear the focused-test boundary.** `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/evidence_gates.go:169` unconditionally replaces a failed gate with `satisfied` when comparison reports `matching_red`. The comparison in `internal/corehelpers/commands.go` checks the saved baseline's launch state but does not require the current execution to have launched successfully or completed without timeout. `internal/nextselection/blocked_probe.go:107` additionally derives `launched`/`timed_out` from reserved numeric exit values, so an ordinary child exit 125/124 is indistinguishable from infrastructure failure/timeout. Public replay produced aggregate `outcome: success`, process exit 0, and a focused gate simultaneously containing `state: satisfied`, `outcome: failure`, `launched: false`, and `BLOCKED-PROBE-LAUNCH-FAILED`. A real timeout likewise exited successfully with `timed_out: true` and `matching_red`. A completion consumer is explicitly told that satisfied gates authorize durable evidence; this is a false success at that production boundary. Preserve actual process lifecycle facts and refuse baseline exclusion for unsuccessful launch, interruption, or timeout; never overwrite a failed subordinate result merely because diagnostic identity matches. — **impact-critical** → returned to the orchestrator for REQ-506's permitted remediation and, only if still unresolved afterward, critical Fold-First routing.
- **F02 — Some emitted missing-input continuations cannot be followed.** `internal/lifecycleadvance/evidence_gates.go:239` puts every missing input after `--`. At qualification, the emitted argv ends `--`, `<exact --diff-range <pre>..<merge>>`; substituting the requested range and flag yields `ADVANCE-GATE-INPUT-IRRELEVANT`, because qualification rejects a separator. The same template places missing `--gate-arg` input in the subordinate channel. The initial classifier's phase-specific argv is correct, but the missing-input record clears that continuation and substitutes the unusable one. Subordinate findings also retain direct helper remedies despite the action forbidding separate helper calls. Generate continuation tokens per phase/channel and preserve the advance boundary. — impact-user-visible → report only
- **F03 — A live reference still prescribes the retired qualification entry point.** `skills/do-work/actions/work-reference.md:502` (line 496 at the saved revision) says the merge range goes to `tools/checks/qualify.sh` through `DO_WORK_DIFF_RANGE`; the new consumer loop says not to call that helper separately and the composed gate requires `--diff-range`. This is an operative hand-back instruction, not a historical mention. Narrow it to the canonical advance input. — impact-rule-change → report only
- **F04 — The estimate short-circuit lost its operative condition during subtraction.** The removed Step 3.6 applied the floor estimate whenever effort was mechanical, or Route A had no heavy-evidence indicators. The replacement at `skills/do-work/actions/work.md:170` tells readers to extract nontrivial signals but omits that condition. `internal/lifecycleadvance/advance_commands.go:154` only emits `--trivial` when effort is mechanical **and** route is A. That narrower classifier condition predates this REQ; deleting the correct action instruction now leaves it as the only direct continuation, while `actions/estimate-reference.md:63` and `work-reference.md:164` still promise mechanical effort alone short-circuits. Preserve the original judgment condition through the new command boundary. — impact-rule-change → report only

**Minor / Nit:** None.

## Requirements checklist

- [x] Current-phase handlers run through `advance`, preserving exact discovered identity and ordered typed subordinate evidence.
- [ ] Missing-input findings provide usable advance continuations for every gate: F02.
- [ ] Focused-test evidence safely distinguishes launch failure, timeout, and baseline red: F01.
- [x] Qualification validates an exact resolvable range and scope comparison; hostile/mismatched/phase-irrelevant argv is rejected without interpolation.
- [x] Green-gate misses request a direct run; exact reported status is recorded through the existing owner.
- [x] Four old procedural headings and their work-action recipes are absent; Go public-command tests cover each replacement gate. At this saved base, `core-checks.sh` already contained no sentence predicates for those four gate procedures, so its actual delta adds removal guards rather than deleting nonexistent predicates. Retained scope checks exercise behavior and are not obsolete sentence predicates.
- [x] Qualification Anti-Rationalization Table and Finding-Closure Ratchet remain unchanged in this range; warning interpretation, TDD, retries, attribution, deferral, heavy testing, and substantive review remain prose.
- [ ] The floor-agent continuation and retained estimate-judgment contract are complete: F02–F04.
- [x] Scope, P-A-U completion, significant Decisions, four-part migration surfaces, and original UR-098 intent were checked against the diff. No unrelated finalization or queue-selection change was attributed to this REQ.

## Acceptance testing

**Result: Fail.** Independent tests ran in a detached checkout of the exact saved implementation. Existing focused packages passed uncached: lifecycleadvance 12.480s, corehelpers 14.765s, nextselection 3.541s, resultmodel 0.615s. These include the built public CLI's estimate, preflight, qualification, baseline-state, hostile-input, launch-failure, timeout, and green-record cases. The existing exact-revision record separately establishes six heavy lanes passing without skips; none was rerun or borrowed from another revision during this review.

Additional real-binary checks exposed the missing combinations:

1. A claimed Route A fixture at `test-gate` had focused probe `exit 125`, a usable saved baseline `{"test_command":"exit 125","exit_status":125,"launched":true}`, and empty baseline diagnostics. Restricting PATH to a private directory containing only `git` caused the actual `sh` launch to fail. `/usr/bin/true` was run directly and its exact zero exit supplied as canonical-gate evidence. `advance REQ-713 --request-path do-work/working/REQ-713-fixture.md --gate-arg /usr/bin/true --gate-exit-status 0 -- --probe-file focused.sh --timeout-seconds 2` exited **0**, with the contradictory focused record described in F01. No fake success was supplied for the focused probe.
2. Direct baseline `/bin/sleep 2; exit 124` completed with exit 124 and empty diagnostics. The identical focused probe with a one-second timeout was terminated by the runner; using the existing matching green-gate record, aggregate advance again exited **0**, focused state `satisfied`, `timed_out: true`, baseline state `matching_red`.
3. Ordinary immediate `exit 124` and `exit 125` probes reported timeout and unlaunched state respectively, showing the lifecycle fact/status collision without any missing executable.
4. Following the qualification missing-input continuation after replacing its placeholder with the requested `--diff-range` tokens returned `ADVANCE-GATE-INPUT-IRRELEVANT` before range validity was considered.

Current-presence check: the three F01 owner files and the F02 continuation implementation are byte-identical between the saved commit and main at review time. The cited F03/F04 prose is also still present. Later REQ-507 finalization extensions and concurrent release work do not resolve these findings. Existing REQ-503 classifier issues were not re-reported as new REQ-506 findings.

## Suggested additional testing

Add public regression combinations for a usable matching baseline plus real current launch failure/timeout, and retain negative controls for a normally exiting child whose status happens to be 124 or 125. Replay every emitted missing-input continuation after substituting only judgment-owned placeholder values. Re-review the remediation at its exact integrated range and renew whatever test evidence its changed source requires.

## Scores

**Overall: 50%** — percentage average 80.42%; Critical risk caps at 60%, then Acceptance Fail caps at 50%.

| Dimension | Score | Basis |
|---|---:|---|
| Requirements | 66.67% | Six of nine checklist items fully delivered; three partial or failed |
| Code Quality | 75% | Existing owners reused; execution-state override breaks their failure authority |
| Test Adequacy | 80% | Strong ordinary public matrix, but missing baseline × infrastructure combinations |
| Scope | 100% | All changed paths declared; two untouched declarations justified |
| Risk | Critical | Failed test execution can authorize completion evidence |
| Acceptance | Fail | Real public command reports success for failed launch and timeout |

## Review provenance and hand-off

Orchestrated, read-only preparation while REQ-506 remained queued. Read the full REQ and UR-098, review-work rubric, shared principles, relevant primes and lessons, general/testing/coding/maintenance/security guidance, and anti-slop reporting guidance. Self-validation checked both the ordinary red control and actual infrastructure failure, rather than inferring failure from reserved exit codes alone. No queue records, source, lifecycle, or commits were changed by this review. The orchestrator must re-resolve canonical claim identity, saved commit, and `resume_phase` evidence before consuming this report.

**Follow-ups created:** None. F01 is handed back for the current request's one remediation attempt; three noncritical findings remain report only. All review-owned fixtures, binaries, and detached checkout were removed after evidence was captured; no child process remains pending.
