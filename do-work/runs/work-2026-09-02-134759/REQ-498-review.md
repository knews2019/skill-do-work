# Review: REQ-498

**Request changes** — planned finalization and one narrow legacy fixture work, but the legacy release-tail recovery that motivated the request is not implemented safely enough to let startup continue.

Route C | merge range `e8e5a79d..75648a49` | reviewed 2026-09-02T15:11:25Z

## What's built

- `finalize --manifest` and `recover-finalization --discover` share the Git-private journal engine, with explicit primary/supplied provenance, rollback before the primary commit, ordered result records, and startup/action delegation.
- The supplied no-journal fixture recovers an archived REQ plus checkpoint and one project file, records provenance, preserves one unrelated ordinary file, and permits the next claim.
- The implementation does not recover the requested semantic release/calibration/UR/follow-up tail, and its project-path association can absorb foreign hunks without proving whole-diff ownership.

## Decisions / risks for the orchestrator

- Run the one permitted remediation attempt before completion. The core discovery defect is production-path and data-integrity significant; archiving as clean would overstate acceptance.
- Consolidate the code/test findings under one root cause: incomplete semantic legacy-finalization ownership. REQ-499's sole-releaser override is adjacent but does not by itself close multi-owner semantic proof or whole-diff ownership; widen it explicitly only if that is the chosen fold.
- Route the stale action restatements to existing prose-reduction work (REQ-504) or the prose backlog after code acceptance is restored.

## Findings

### Important

1. **Legacy discovery does not implement the lifecycle-aware semantic association required by the REQ.** `finalization_discovery.go:97-116` treats every non-`do-work/` path as ordinary Implementation Summary ownership, while `:117-139` recognizes only request/archive identity and checkpoint removal for shared `do-work/` state. There is no calibration-row, UR membership/move, originating-follow-up, changelog/version/release-entry, or `release_at` proof. Consequently, dirty release metadata not listed in the summary is silently left outside every group and startup can return success; if a shared release or project path is listed, its whole current file is claimed without proving that every hunk belongs to the REQ. This misses the REQ-494 release-tail shape and can commit foreign changes. — `impact-critical` → remediate under one semantic-legacy-finalization sweep; do not treat the current REQ-499 shape as sufficient without explicit widening.

2. **Multi-group commit recognition fails in the crash window it is meant to recover.** `matchingHeadCommit` selects only the first descendant of `PreparedHead` (`finalization_apply.go:466-482`). All discovered groups are prepared against the same original HEAD. After group A commits, a group-B retry whose commit exists but whose phase write was interrupted sees A as the first descendant, rejects it, and never searches later commits for B's exact prepared diff. The next exact commit attempt then has no B diff to commit, leaving recovery stuck. The test at `finalization_recovery_test.go:132-153` covers a later commit after the target, not an earlier independent-group commit before the target. — `impact-critical` → fold into the same semantic-legacy-finalization remediation.

3. **Protected-path handling contradicts the preservation contract.** `finalization_discovery.go:36-45` refuses on every unstaged `X` or `XD` inventory row before checking the index. The request says unrelated unstaged work is preserved and only a pre-staged secret must block; it also asks for distinct protected-staged versus foreign-staged reasons. An unrelated unstaged secret currently prevents recovery and receives only `FINALIZATION-DISCOVERY-PROTECTED`. — `impact-user-visible` → fold into the same discovery remediation.

4. **The typed result omits required terminal and refusal evidence.** `FinalizationResult` exposes one `phase` and settled primary/metadata hashes (`result_model.go:247-266`), but not created-this-invocation hashes or verified/cleanup-complete progress. Successful cleanup still reports `metadata_committed` (`finalization_apply.go:556-563`). Global discovery refusals return findings without a finalization record (`finalization_commands.go:65-74`), so blocked paths/reason codes are not present in the ordered record contract that `work.md` tells actions to consume. — `impact-rule-change` → fold into the finalization contract remediation.

5. **The TDD/acceptance matrix is largely absent, and the green fixture avoids the failing surface.** `finalization_recovery_test.go` adds four tests: one narrow archive/checkpoint/project fixture, provenance validation, a generic staged refusal, and commit recognition. It does not test release/version/calibration/UR/follow-up semantics, foreign hunks, protected staged versus unstaged state, two safe groups, corrupt images, every journal phase, failures/already-green/worktree provenance, duplicate effects, or hook recovery across the required modes. The focused package passes because those requirements are not exercised. — `impact-user-visible` → tests belong in the same code remediation, not a separate test-only REQ.

6. **Action restatements still instruct the retired tail.** The canonical Step 1 and Step 9 text delegates to finalization, but `work.md:626` says already-green work invokes `complete`, `:658` says to read the canonical `complete` result, and the checklist at `:746-762` says checkpoint is read first, archive mutates before Step 9, and Step 9 stages/records through `complete`. `work-reference.md:453` repeats the direct complete/release/stage/metadata-commit path. These active restatements can send an orchestrator around the new journal. — `impact-user-visible` → append to REQ-504 if its scope remains the action-prose collapse; otherwise route as prose backlog, not a new standalone REQ.

### Minor

None beyond the Important findings.

### Nit

None.

## Requirements Checklist

- [x] Strict `finalize --manifest` and `recover-finalization --discover` commands exist and remain registered.
- [~] Journaled planned-tail phases are resumable — happy path and one lifecycle interruption work, but verification/cleanup are not represented and the requested phase matrix is absent.
- [x] Canonical request-state, publication, protected-inventory, and exact Git transaction authorities are composed rather than duplicated.
- [~] Idempotence evidence exists for one narrow discovered tail; duplicate release/version/calibration/UR effects and multi-group crash recovery are unproved and one multi-group crash window is broken.
- [x] Primary versus supplied provenance is explicit and validated for format, resolution, and ancestry.
- [x] `work` and `commit` invoke recovery before their ordinary startup/association operations.
- [ ] Legacy tails use REQ-specific semantic evidence for all shared metadata and never generic ownership.
- [ ] Unrelated unstaged/protected changes are preserved while ambiguous shared or foreign state blocks with exact typed reasons.
- [~] Ordered typed records exist, but terminal phase, created-versus-settled commits, and global discovery refusal evidence are incomplete.
- [~] Work/commit canonical sections delegate finalization, but stale operative restatements retain the direct tail.
- [x] Existing `complete`, `fail`, and `release` command interfaces remain present; the single-releaser model is retained.
- [ ] The original interruption/legacy-association acceptance matrix is implemented and demonstrated.

## Acceptance Testing

**Result: Fail**

- `go test -count=1 ./internal/finalization` passed independently.
- `bash _dev/tests/contract-regressions.sh` passed independently.
- The captured GREEN `TestRecoverFinalizationDiscoversLegacyNoJournalTail` passes, but its fixture includes only archive/working, checkpoint, one implementation file, and one unrelated ordinary file.
- A source-level acceptance trace for the requested REQ-494 release tail fails: `finalization_discovery.go` contains no semantic handler for calibration, changelog/version, UR membership/moves, follow-up provenance, or release timestamps. The generic non-`do-work/` loop either ignores those dirty paths or assigns them solely by Implementation Summary path, which is not the required evidence boundary.
- RED history was credibly recorded in the hand-back (`--discover` rejected before implementation), but the GREEN case does not match the original release-tail acceptance fixture closely enough to close it.

## Suggested Additional Testing

- Reproduce the full REQ-494 shape with archive move, enriched checkpoint removal, calibration row, UR closure/move, changelog/version mirrors, release timestamp, implementation bytes, and an unrelated ordinary edit; assert exact ownership, commit contents, provenance, and successful next claim.
- Add a foreign hunk to an otherwise owned project/shared path and prove byte-identical refusal.
- Recover two safe groups, interrupt group B after its commit but before its phase write, and prove exact commit recognition despite group A's earlier commits.
- Distinguish unrelated unstaged `X`/`XD`, pre-staged secret, and ordinary pre-staged entries with exact stable reason codes and no secret reads.
- Table-drive every phase across success, failure, already-green/no-release, primary provenance, and supplied-worktree provenance; assert one archive/calibration/release/version/provenance effect after retries.
- Assert success, per-REQ refusal, and global refusal JSON all carry non-null ordered records/slices, terminal verification/cleanup state, created hashes, settled hashes, exact blockers, reasons, and rooted next/verification argv.

## Scores (on the record — not the headline)

**Overall: 50%**

| Dimension | Score | Notes |
|---|---:|---|
| Requirements | 45% | Planned-tail foundation is present; core semantic legacy recovery and the original acceptance matrix are not. |
| Code Quality | 45% | Reuses good authorities, but ownership and multi-group recovery have critical correctness gaps. |
| Test Adequacy | 30% | Focused tests pass but omit most required safety and interruption cases. |
| Scope | 100% | All 14 changed files match the updated declared scope; the command-test expansion was justified and documented. |
| Risk | Critical | Automatic recovery can ignore shared release dirt or commit foreign hunks under insufficient ownership evidence. |
| Acceptance | Fail | The motivating legacy release-tail shape is not safely recovered. |

Raw percentage average: 55%. Critical risk caps at 60%; Acceptance Fail caps the final score at 50%.

## Follow-up Disposition

No follow-up REQ was created or queue state modified by this reviewer. The orchestrator should apply its single remediation attempt first. If re-review still finds the class open, create or fold one consolidated semantic-finalization recovery sweep for findings 1-5; route finding 6 to REQ-504 or the prose backlog.

## Self-validation

- Re-read the request, original UR-096 input, plan, exploration, builder hand-back, exact merge diff, all new tests, result projection, action consumers, and documented restatements.
- Checked P-A-U: all three boxes are complete and the hand-back's Decisions section is readable and consistent with the code.
- Re-ran focused finalization tests and the action contract suite. Their green status does not contradict the findings because the missing behaviors are not covered by those fixtures.
- Rechecked the main risk from both directions: unowned release paths are omitted from groups, while listed non-`do-work/` paths are accepted without semantic or hunk-level proof.

*Reviewed by review-work action*
