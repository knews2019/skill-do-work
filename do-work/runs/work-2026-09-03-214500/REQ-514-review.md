# REQ-514 Independent Review

## Decision Brief

### WHAT'S BUILT

- Top-level `outcome=refused` results now remove a `next_argv` that resolves to the invoking CLI verb while preserving `verification_argv`.
- A self-referential refusal with only REQ-owned blockers renders as `outcome=findings`; an unowned blocker remains globally refused.
- Ambiguous legacy-finalization findings now name `uncommitted-inventory` as their top-level resolver.

### DECISION

**Request changes** — the original REQ-456 self-loop is closed at the rendered top-level finding, but the universal refusal invariant, the typed finalization projection, and the folded serial-claim acceptance test are incomplete.

Route C | merge range `26b3426886bfea6183502809a7e5e93799831a52..e42ae1e57c9f2692a598cb08daca2fe99bec6a45`

### DECISIONS / RISKS FOR YOU

- No product choice is needed. The current patch should receive one bounded remediation pass before finalization.

### FINDINGS

**Important:**

- F1. The invariant runs only when the whole result has `outcome=refused`, not for every refusal finding as required. [`NormalizeResult` checks `result.Outcome == OutcomeRefused`](../../../skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go#L541), so a `FixabilityRefused` finding carried by `outcome=findings` keeps a self-referential remedy. The live cleanup producer returns `do-work-cli cleanup` as `next_argv` ([`cleanup_apply.go`](../../../skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go#L626)), `ApplyPlan` promotes results with findings to `OutcomeFindings` ([`cleanup_apply.go`](../../../skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply.go#L185)), and its existing test still requires that forbidden self-remedy ([`cleanup_apply_test.go`](../../../skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_apply_test.go#L615)). The new “every finding builder” test uses three synthetic `CommandResult` values instead of walking real refusal producers ([`result_model_test.go`](../../../skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go#L252)), so it cannot catch this production counterexample. — `impact-rule-change` → report only
- F2. Ambiguous finalization discovery now emits contradictory typed remedies. The top-level finding correctly names `uncommitted-inventory`, but the same result's singular and ordered `FinalizationResult.NextArgv` still name `recover-finalization --discover`, identical to the invoking verb ([`finalization_discovery.go`](../../../skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go#L1663)). `normalizeFinalization` only replaces nil slices and does not enforce remedy identity ([`result_model.go`](../../../skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go#L580)); the added finalization test asserts only `result.Findings[0]` ([`finalization_recovery_test.go`](../../../skills/do-work/tools/do-work-cli/internal/finalization/finalization_recovery_test.go#L344)). This matters because the shipped action declares ordered `finalizations` the canonical record consumed by callers ([`work-reference.md`](../../../skills/do-work/actions/work-reference.md#L1101)). — `impact-user-visible` → report only
- F3. The folded serial claim → implementation → complete → recovery acceptance fixture was not delivered end to end. The retained test writes a REQ already in `status: claimed`, commits that seed, and then edits the implementation ([`finalization_commands_test.go`](../../../skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands_test.go#L16)); it never invokes the `claim` command. It proves finalization recovery from a seeded claimed state, but not the requested serial claim transaction that caused the REQ-456 trap. The builder handback's statement that this existing test exercises the required claim fixture therefore overstates coverage. — `impact-rule-change` → report only

**Minor:** None.

**Nit:** None.

### REQUIREMENTS CHECKLIST

- [ ] Every refusal finding with non-empty `next_argv` names a different verb — partially delivered; only top-level `OutcomeRefused` results are normalized, and F1 is a live counterexample.
- [x] Same-command `verification_argv` remains allowed and preserved — delivered by result-model and runtime tests.
- [ ] No-truthful-alternate REQ refusal becomes set-aside evidence without a remedy — delivered for the tested `OutcomeRefused`/`AffectedIDs` shape, not enforced across every refusal finding.
- [ ] One table-driven test walks every refusal builder — not delivered; the table contains synthetic results and omits actual producers.
- [ ] Existing self-referential refusal evidence is fixed or converted — partially delivered; the top-level finalization finding is fixed, but ordered finalization evidence and `OutcomeFindings` producers remain self-referential.
- [x] The invariant is implemented in Go rather than a prose/shell predicate — delivered.
- [ ] Folded serial claim → one-line implementation → complete → recovery reaches `cleanup_complete` — not delivered as an actual claim-command sequence.
- [x] Fail-closed behavior is preserved — no reviewed path broadens recovery or mutates refused state.
- [x] Scope is reconciled to the five files in the merged range — delivered; no `do-work/` branch content or unrelated source is present.

### ACCEPTANCE TESTING

**Result: Fail**

- `go test -count=1 ./internal/resultmodel ./internal/commandruntime ./internal/finalization ./internal/requeststate` — PASS in 68.49s.
- `go test -count=1 ./...` — PASS in 66.19s.
- `go vet ./...` — PASS in 0.14s.
- `git diff --check 26b3426886bfea6183502809a7e5e93799831a52..e42ae1e57c9f2692a598cb08daca2fe99bec6a45` — PASS.
- Diff and restatement sweep — FAIL on F1 and F2: a live refusal producer and the canonical typed finalization record retain self-referential `next_argv`.
- Folded finding-closure check — FAIL on F3: the named serial-claim sequence is simulated from an already-claimed file rather than executed.

### SUGGESTED ADDITIONAL TESTING

- Add the required real-producer refusal table across every result shape that can carry `FixabilityRefused`, including `OutcomeFindings` and ordered/singular typed finalization records; assert the rendered JSON contains no self-referential non-empty `next_argv`.
- Add the literal serial fixture through the public handlers/runtime: execute `claim`, make the one-line implementation change, finalize completion, interrupt at lifecycle apply, then run recovery and assert `cleanup_complete` plus a clean repository.

### SCORES (ON THE RECORD — NOT THE HEADLINE)

**Overall: 50%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 35% | Central invariant and folded acceptance sequence are only partially delivered. |
| Code Quality | 65% | Small, readable boundary, but it keys on aggregate outcome and leaves conflicting typed evidence. |
| Test Adequacy | 40% | Focused tests pass, but the universal table is synthetic and the serial claim is not executed. |
| Scope | 100% | Exact reconciled five-file range; no drift. |
| Risk | Low | Fail-closed behavior remains, but callers can still receive the trap shape the REQ aimed to forbid. |
| Acceptance | Fail | Static acceptance found production counterexamples despite green suites. |

Raw percentage average: 60%. `Acceptance: Fail` caps the recorded overall score at 50%.

### FOLLOW-UPS CREATED

None (3 findings report only). The review did not mutate queue or request state.

### SELF-VALIDATION

Rechecked the result normalization order, both finalization projections, the runtime render seam, live cleanup refusal builders, existing cleanup expectations, the folded recovery fixture, and the action consumer. No additional findings emerged. The original top-level REQ-456 case does normalize correctly; the findings above are the unclosed contract surfaces rather than a claim that the entire patch is ineffective.

## Append-Ready Durable Review Block

```markdown
## Review

**Overall: 50%** | 2026-09-04T00:00:24Z

| Dimension | Score |
|-----------|-------|
| Requirements | 35% |
| Code Quality | 65% |
| Test Adequacy | 40% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Fail |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Result normalization applies the no-self-remedy invariant only to aggregate `OutcomeRefused`; live `FixabilityRefused` findings under `OutcomeFindings` retain same-verb `next_argv`, and the synthetic table does not walk real producers — `impact-rule-change` → report only
- Ambiguous discovery's top-level finding names `uncommitted-inventory`, but its canonical singular/ordered finalization records still name `recover-finalization --discover` as both next and verification, producing contradictory typed remedies — `impact-user-visible` → report only
- The folded serial claim acceptance fixture seeds an already-claimed REQ instead of executing the claim command, so it does not pin the requested claim → implementation → complete → recovery sequence — `impact-rule-change` → report only

**Minor findings:** None
**Acceptance:** Fail — focused and full Go checks pass, but the universal refusal invariant and folded serial-claim sequence have direct acceptance gaps.
**Suggested testing:** 2 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*
```
