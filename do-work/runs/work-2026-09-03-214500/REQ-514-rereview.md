# REQ-514 Independent Re-review

## Decision Brief

### WHAT'S BUILT

- The remediation now applies same-verb refusal normalization to `FixabilityRefused` findings even when their aggregate result is `outcome=findings`.
- The cited cleanup scratch refusal no longer emits a raw `cleanup` self-remedy.
- Ambiguous discovery now aligns its top-level finding and both typed finalization projections on `uncommitted-inventory`, retaining `recover-finalization --discover` only as verification.
- The folded lifecycle fixture now begins with a pending queue REQ, executes the real committed `claim` handler, changes one implementation file, interrupts completion after lifecycle apply, and recovers to a clean terminal repository.

### DECISION

**Request changes** — F2's exact discovery defect and F3's missing real-claim sequence are closed, and F1's live `OutcomeFindings` defect is closed. Two acceptance surfaces remain: the mandated all-refusal-builders test is still a synthetic result-model table, and journal-backed recovery failures still publish the invoking recovery command as the canonical typed next action.

Route C | original integration range `26b3426886bfea6183502809a7e5e93799831a52..e42ae1e57c9f2692a598cb08daca2fe99bec6a45` | remediation range `f319cbbb..bbc1292809990afa72c5093ad29f31ff3dac7b48`

### DECISIONS / RISKS FOR YOU

- No product decision is needed. A narrow source-and-test correction is sufficient: distinguish journal-failure recovery verification from a truthful resolving command in typed records, and replace or augment the synthetic invariant table with the requested real-producer walk.

### FINDINGS

**Important:**

- R1. The remediation closes F1's cited live cleanup counterexample, but not the explicit test requirement that “one table-driven test walks every finding builder that refuses.” `TestRefusalRemediesNeverNameTheInvokingCommand` still constructs four synthetic `CommandResult` values by hand rather than invoking any production finding builder ([`result_model_test.go`](../../../skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go#L252)). The new fourth row is valuable coverage for `OutcomeFindings`, and the real cleanup regression independently pins the cited producer, but there is no enumerated producer inventory or mechanically complete walk across the refusal builders in finalization, request-state, cleanup, publication, installation, hooks, and the other command families. A newly added raw producer can therefore escape direct acceptance coverage even though the render boundary will normalize its top-level finding. — `impact-rule-change` → report only
- R2. Journal-backed failures returned by the live `recover-finalization` producer still publish a self-referential remedy in the canonical typed finalization records. `finalizationFailure` assigns `recoveryArgv(journal)` to both the finding and record ([`finalization_apply.go`](../../../skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go#L617)); `finalizationRecord` stores that argv as both `NextArgv` and `VerificationArgv` ([`finalization_apply.go`](../../../skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go#L627)). During recovery, `handleRecoverFinalization` appends that record into singular/ordered output and returns it on the first failure ([`finalization_commands.go`](../../../skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go#L50)). Runtime supplies `Command="recover-finalization"` before rendering ([`command_runtime.go`](../../../skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go#L63)); normalization removes the top-level finding's same-verb next action but `normalizeFinalization` only fills nil slices and leaves typed argv identity unchanged ([`result_model.go`](../../../skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go#L581)). The focused corrupt-journal test exercises this producer but asserts only outcome and reason code, not typed argv ([`finalization_recovery_test.go`](../../../skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go#L458)). This conflicts with the action contract that ordered `finalizations` are canonical and global/per-REQ refusals use that record contract ([`work-reference.md`](../../../skills/do-work/actions/work-reference.md#L1101)). The discovery-specific F2 fix is correct; this is the same consumer-visible contradiction on the journal path. — `impact-user-visible` → report only

**Minor:**

- R3. F3's real claim-to-recovery sequence is now present and its strong postconditions establish successful cleanup, but the folded acceptance fixture does not literally assert the required `Finalization.Phase == cleanup_complete`; its recovery assertion checks success, non-nil, and `Resumed`, then proves distinct commits, archive provenance, implementation bytes, journal removal, and a clean worktree ([`finalization_commands_test.go`](../../../skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands_test.go#L70)). Add the phase assertion so a future success-shaped response cannot satisfy the named acceptance contract without the terminal typed state. — `impact-rule-change` → report only

**Nit:** None.

### ORIGINAL-FINDING CLOSURE

| Finding | Status | Evidence |
|---|---|---|
| F1: aggregate-only normalization and live cleanup self-remedy | **Partially closed** | `NormalizeResult` now treats `FixabilityRefused` as a refusal independently of aggregate outcome ([`result_model.go`](../../../skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go#L541)); the cited cleanup builder emits no next argv ([`cleanup_apply.go`](../../../skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go#L619)); its real producer test asserts empty next argv. The required table still does not walk production builders (R1). |
| F2: discovery typed records contradict the finding | **Closed** | `discoveryRefusal` assigns inventory collection to finding, singular record, and ordered record, with recovery retained as verification ([`finalization_discovery.go`](../../../skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go#L1663)); focused source and live CLI checks agree. Journal-backed typed failures remain a separate unclosed consumer surface (R2). |
| F3: fixture seeds an already-claimed request | **Closed with minor assertion gap** | The fixture writes a pending queue REQ and invokes the real `claim --commit` handler, verifies its claim commit, then performs the implementation/finalize/interruption/recovery sequence ([`finalization_commands_test.go`](../../../skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands_test.go#L17)). The exact typed phase assertion is absent (R3). |

### REQUIREMENTS CHECKLIST

- [ ] Every refusal with non-empty `next_argv` names a different verb than the invocation — top-level findings are normalized, but journal-backed recovery's canonical typed records still name the invoking recovery verb (R2).
- [x] Same-command `verification_argv` remains allowed and preserved — result normalization and the live discovery output preserve it.
- [x] A no-alternate REQ-owned self-refusal becomes set-aside evidence with empty next argv and a stop reason — the generic owned/unowned result-model cases and runtime exit behavior pass.
- [ ] One table-driven test walks every refusal finding builder — the only table is synthetic (R1).
- [ ] Existing self-referential refusal evidence is fixed or converted — the cited cleanup and discovery shapes are fixed; journal-backed typed recovery failures remain (R2).
- [x] The invariant lives in Go, not an action prose or shell predicate.
- [ ] Folded serial claim → one-line implementation → complete → recovery reaches typed `cleanup_complete` — the real sequence and terminal repository effects pass, but the named typed phase is not asserted (R3).
- [x] Fail-closed behavior is preserved — no reviewed change widens recovery or mutates refused evidence.
- [x] Remediation scope is reconciled through D-02 and the seven-file range contains no lifecycle, queue, release, or unrelated changes.

### ACCEPTANCE TESTING

**Result: Fail**

- `go test -count=1 ./internal/resultmodel ./internal/commandruntime ./internal/cleanup ./internal/finalization ./internal/requeststate` — PASS (`resultmodel` 0.419s, `commandruntime` 1.372s, `cleanup` 17.276s, `finalization` 60.235s, `requeststate` 9.873s).
- `go test -count=1 ./internal/finalization -run 'Test(DiscoveryRefusalNamesInventoryAsTheResolvingVerb|RecoverFinalizationResumesJournalAfterLifecycleInterruption|RecoverFinalizationRefusesCorruptJournalImage)' -v` — PASS; this exercised the real claim/recovery fixture, discovery remedy split, and a live journal-failure producer.
- `go test -count=1 ./internal/resultmodel ./internal/cleanup -run 'Test(RefusalRemediesNeverNameTheInvokingCommand|ConsumedUntrackedRunCommitRefusesScratchOnlyAndMixedDeletion)' -v` — PASS, including the four result-model rows and both cleanup scratch variants.
- Live producer check: built the current CLI, staged an ordinary foreign file in a temporary Git repository, and ran `do-work-cli --repo-root <fixture> --format json recover-finalization --discover` — PASS for F2: top finding, singular `finalization`, and ordered `finalizations[0]` all named `do-work-cli --format json uncommitted-inventory`; all three retained `recover-finalization --discover` as verification.
- Static live-consumer trace — FAIL for R2: the exercised corrupt-journal producer routes `finalizationFailure` through `appendFinalizationResult`; the rendering boundary normalizes findings but does not normalize typed remedy identity.
- Restatement sweep: searched all internal Go producers for `OutcomeRefused`, `FixabilityRefused`, and `NextArgv`, then traced canonical finalization consumers. The sweep confirmed the discovery remediation and exposed R1/R2; no action prose restatement was added.
- `git diff --check 26b3426886bfea6183502809a7e5e93799831a52..e42ae1e57c9f2692a598cb08daca2fe99bec6a45` — PASS.
- `git diff --check f319cbbb..bbc1292809990afa72c5093ad29f31ff3dac7b48` — PASS.

### SUGGESTED ADDITIONAL TESTING

- Exercise the journal-backed `recover-finalization` refusal through runtime JSON rendering and assert the top finding plus singular/ordered finalization records contain no same-verb non-empty `next_argv`, while same-command verification remains intact.
- Replace or supplement the synthetic result-model table with the requested table of real refusal producers/builders and make additions to the refusal-builder inventory explicit.
- Add `recovered.Finalization.Phase == string(PhaseCleanupComplete)` to the actual claim-to-recovery fixture.

### SCORES (ON THE RECORD — NOT THE HEADLINE)

**Overall: 50%**

| Dimension | Score | Notes |
|---|---:|---|
| Requirements | 68% | F2 discovery and the real claim flow are corrected, but the universal builder test and a canonical typed recovery surface remain incomplete. |
| Code Quality | 76% | The centralized finding normalization is compact and fail-closed; typed record semantics are not governed by the same invariant. |
| Test Adequacy | 66% | Focused producer tests are green and materially improved, but the universal table remains synthetic and journal typed argv is unasserted. |
| Scope | 100% | The original and remediation ranges are coherent and D-02 reconciles the seven-file remediation expansion. |
| Risk | Low | The remaining defect does not broaden mutation, but it can send canonical typed consumers back to the command that just refused. |
| Acceptance | Fail | A production journal-backed refusal still exposes the trap shape in canonical typed output. |

Raw percentage average: 77.5%. `Acceptance: Fail` caps the recorded overall score at 50%.

### FOLLOW-UPS CREATED

None (3 findings report only). This re-review did not mutate queue, request, source, or release state.

### SELF-VALIDATION

Reviewed the complete original and remediation ranges, the original review, remediation handback, authoritative REQ, brief, UR-099 intake, output normalization and runtime render seams, cleanup producer, discovery producer, journal failure producer, finalization aggregation, canonical action consumer, and the actual claim-to-recovery fixture. Re-ran the focused packages and named acceptance tests, performed a live discovery refusal check, checked both ranges for whitespace errors, and re-read every cited line. F2's exact discovery defect and F3's real-claim gap are genuinely closed; R1/R2/R3 are bounded residuals rather than restatements of the original failures.

## Append-Ready Durable Re-review Block

```markdown
## Re-review

**Overall: 50%** | 2026-09-04T00:13:17Z

| Dimension | Score |
|---|---:|
| Requirements | 68% |
| Code Quality | 76% |
| Test Adequacy | 66% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Fail |

**Original-finding closure:**
- F1 partially closed: aggregate-independent normalization and the cited cleanup producer are fixed; the mandated table still constructs synthetic results instead of walking all refusal builders.
- F2 closed for discovery: finding, singular record, and ordered record use inventory as next action and recovery only as verification.
- F3 closed for the real claim path: the fixture runs committed claim before implementation and recovery; it omits the literal typed `cleanup_complete` assertion.

**Important findings (each with its recorded impact token — durable audit record):**
- The required table-driven walk of every real refusal finding builder remains a four-row synthetic `CommandResult` table, so producer completeness is not pinned — `impact-rule-change` → report only
- Journal-backed `recover-finalization` failures normalize the top finding but retain `recover-finalization` as `next_argv` in canonical singular/ordered finalization records — `impact-user-visible` → report only

**Minor findings:**
- The real claim-to-recovery fixture proves terminal effects but does not literally assert `Finalization.Phase == cleanup_complete` — `impact-rule-change` → report only

**Acceptance:** Fail — focused tests and the live discovery output pass, but a production journal-backed refusal still exposes a self-referential canonical typed next action.
**Suggested testing:** 3 items
**Follow-ups created:** None (3 findings report only)

*Re-reviewed by review-work action*
```
