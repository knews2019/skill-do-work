# Re-review: REQ-504

**Approve** — the single remediation closes both legacy-checkpoint failures from F-01, including the fresh-claim path that could otherwise hide another request's surviving evidence.

Route C. Cumulative evidence remains `773787b74acddfdfc4c16498a89d99a5cc3ab716..6a11b60c83615791769d57b082580f0b69323984`. The original 26-path migration at `f412a8411057d0a833df5584657161008f315b84` and its first independent review remain part of this review. Repair attribution is exactly `f13bdcc0d1a83c0d24e5af6262f6863ee31ef8d7..6a11b60c83615791769d57b082580f0b69323984`: seven declared paths. Intervening sibling and owner commits in the cumulative range are not attributed to this builder.

## What's built

Public recovery still settles finalization before classifying working claims and requires explicit takeover authority. Checkpoint discovery, removal, absence checks, and fresh insertion now share the existing canonical-or-legacy range. Refresh preserves the supported legacy body instead of adding an empty section that hides it. The original command ownership and prose reduction remain intact.

## Findings and F-01 closure

**Important:** None remaining. F-01 is closed by direct public-command acceptance, not merely by observing a shared helper.

**Minor:** None. **Nit:** None.

Before repair, the final public regression tests independently reproduced both reported failures: refresh reported five preserved claims while subsequent discovery returned zero, and explicit takeover refused with `RECOVER-CLAIM-CHECKPOINT-EVIDENCE`. With the integrated repair, refresh retains all five semantic claim records and exact body bytes; takeover removes all four matching labelled, duplicate, numeric-alias, and unlabelled records plus continuations. Public `next` selects the recovered request, a committed fresh claim succeeds, and public recovery still sees the unrelated writer's original claim.

## Requirements checklist

- [x] Finalization-first recovery, typed takeover decisions, and constant authority argv remain implemented and executable. Observation preserves bytes; explicit targeted and sole-authority paths pass.
- [x] Supported legacy and canonical checkpoint evidence remains recoverable and visible. Canonical section precedence, prose-only heading tokens, numeric aliases, writer-specific removal, unlabelled records, and CRLF controls pass.
- [x] False legacy absence is refused without a commit or filesystem change. All-entry removal stays restricted to the requested identity; writer-specific removal retains other writers and unlabelled records.
- [x] Checkpoint-mode advance changes only its checkpoint target and retains unrelated project dirt. Ordinary working-request advance remains read-only.
- [x] The original migration retains its four required parts: owning commands, reduced action/reference prose, retired shell fixture/registration, and public Go behavior tests. The documented stale captured RED premise was replaced with actual public-boundary failures.
- [x] The original inherited ownership, Step 8, commit, handoff, and guide repairs remain represented. Later selection, evidence-gate, and finalization extensions are separate chain work.
- [x] The seven remediation files are within the original 26-path Scope. Both execution-state checklists are complete. Decisions D-08 through D-12 explain compatibility, range attribution, and the fresh-insertion branch; the readable handback matches the diff.

## Acceptance testing

**Result: Pass.** Bounded independent commands ran from `skills/do-work/tools/do-work-cli`. The tested seven files match integration `6a11b60c83615791769d57b082580f0b69323984`; execution checkout HEAD was observed as `eabc29842dc537eb2cbc43adcfd2ae294bfbc92b`, whose only additional committed change is an unrelated report HTML edit.

- Public regression and authority group: `go test -count=1 ./internal/lifecycleadvance -run 'Test(RecoverLegacyCheckpointClaimsThroughPublicCommand|AdvanceCheckpointPreservesLegacyClaimDiscovery|RecoverPublicCommandRunsFinalizationThenRecoversEveryClaim|RecoverWithoutAuthorityOffersTypedTakeoverAndDoesNotMutateClaim|AdvanceCheckpointWritesOnlyCheckpointAndPreservesClaims|WorkingAdvanceRemainsReadOnlyAfterCheckpointMode)$' -v` — exit 0, package 3.630s, five actual tests passed. The obsolete checkpoint-only name in this selector matched no test; the correct test was then run explicitly below.
- Checkpoint-only public control: `go test -count=1 ./internal/lifecycleadvance -run '^TestAdvanceCheckpointChangesOnlyCheckpointAndPreservesLiveEntries$' -v` — exit 0, package 1.226s; exact changed-path and foreign-byte assertions passed.
- Structural and authority controls: `go test -count=1 ./internal/requeststate ./internal/repositorymodel -run 'Test(CheckpointRemovalPreservesRangeAndWriterAuthority|RecoveryRefusesFalseLegacyCheckpointAbsence|CheckpointDiscoveryUsesCanonicalOrLegacyClaimRange)$' -v` — exit 0; request-state 0.450s and repository-model 0.162s. All selected subcases passed without skips.
- Independent RED replay: a detached checkout at `f13bdcc0d1a83c0d24e5af6262f6863ee31ef8d7` received only the final two public test files. Running `go test -count=1 ./internal/lifecycleadvance -run 'Test(RecoverLegacyCheckpointClaimsThroughPublicCommand|AdvanceCheckpointPreservesLegacyClaimDiscovery)$' -v` exited 1 in 4.965s on the two original semantic failures, not compilation or fixture errors. The same tests pass above at the integrated source.

The builder's recorded RED-before-production ordering and owner-package, vet, formatting, and diff checks agree with this independent replay. This reviewer did not run a full or heavy lane. The original six-lane green record belongs to the old implementation revision and cannot verify changed source; the orchestrator owns current integrated gate and heavy evidence before completion.

## Restatement, domain, and self-validation

Swept `CheckpointClaimBounds`, the replaced private helper, checkpoint writers/removers/absence predicates, and checkpoint-heading/legacy prose across the shipped tree. Cleanup reaches the same writer-specific helper; terminal transitions reach the same all-entry helper. Existing authority predicates and request identity matching were not widened. Canonical-heading prose still describes the normal emitted format and does not authorize a legacy migration or an alternate action-owned writer. No stale operative restatement introduced by the repair was found; REQ-544's separate caller-authored publication/cleanup sweep remains separate.

Security review focused on explicit authority and exact-byte preservation. No new shell interpolation, public schema, dependency, or removal authority was introduced. The shared helper is earned by the reproduced disagreement and is used by multiple existing consumers. Self-validation checked actual test execution, the fresh-claim alternate writer, refusal non-mutation, source attribution, and reviewer cleanup.

## Scores

**Overall: 100%** — arithmetic average of the four percentage dimensions; no qualitative penalty applies.

| Dimension | Score |
|---|---:|
| Requirements | 100% |
| Code quality | 100% |
| Test adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Suggested additional testing:** finish the orchestrator-owned current-revision repository gate and selected heavy lanes. No additional manual acceptance check is required for this bounded repair.

**Follow-ups created:** None (0 findings report only).

**Cleanup confirmed:** all reviewer test processes completed. The only temporary tracked changes were the two test files in the reviewer-owned detached checkout; those were restored, its full porcelain status was empty, and `.git/work-run-20260905/re-review-504` was removed without force. No background process, fixture, or reviewer worktree remains. This report is the sole main-tree write; no source, request, queue, status, or commit was changed by this reviewer.
