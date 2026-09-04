# REQ-547 hand-back — stop finalize refusing a REQ that has no checkpoint entry

Branch: `worktree-agent-REQ-547-finalize-refuses-a-req-with-no-checkpoint-entry`
Commit: `b3d25c8` — `[REQ-547] stop the checkpoint deciding whether a REQ can finalize`

## File Manifest

| File | Change |
|---|---|
| `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` | `PlannedPostimages` now emits exactly the plan's declared `TargetPaths` and errors on a declared target it cannot project, so the journal's preimage and postimage sets cannot name different paths. |
| `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` | A terminal transition (complete / fail / cancel) removes every checkpoint entry for the REQ it archives, instead of only one carrying the manifest's writer label. |
| `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` | New `recover-finalization --discard-journal REQ-NNN` verb plus its refusal projection; `parseRecoverArguments` gained the value-taking flag and its exclusivity rule. |
| `skills/do-work/tools/do-work-cli/internal/finalization/finalization_types.go` | `commandOptions.DiscardJournalRequestID`. |
| `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req547_test.go` | New: absent-entry finalize, foreign-writer-label clearing, the discard verb's four cases, and lifecycle/release proof that genuinely disagreeing image sets still refuse. |
| `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` | New `TestPlannedPostimagesNamesExactlyTheDeclaredTargets`; updates to two REQ-468/REQ-538/REQ-489-era tests (see D-04). |
| `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` | One sentence in the `internal/finalization/` row naming the discard verb and its bounds (no new line; the file's length is unchanged). |

Not touched: `CHANGELOG.md`, `VERSION`, anything under `do-work/` in the worktree, anything in `internal/heavyverification/`.

## Verification

Environment: Go 1.26.1 toolchain, module root `skills/do-work/tools/do-work-cli`.

- `gofmt -l internal cmd` — no output.
- `go vet ./...` — clean.
- `go test -count=1 ./...` — every package `ok`, no failures. Run after the final edit.
- Focused: `go test -count=1 ./internal/finalization ./internal/requeststate ./internal/lifecycleadvance` — all `ok`.

### Revert-and-show-red, per behaviour change

**C-01 — `PlannedPostimages` keyed on the declared targets.** Reverted the projection tail to the old "emit every role path" loop.

- `TestPlannedPostimagesNamesExactlyTheDeclaredTargets/a_role_path_the_plan_did_not_declare_is_not_a_postimage` RED:
  `postimages = [do-work/CHECKPOINT.md do-work/archive/REQ-311.md do-work/calibration-log.tsv do-work/working/REQ-311.md], declared targets = [do-work/archive/REQ-311.md do-work/calibration-log.tsv do-work/working/REQ-311.md]` — four postimage paths against three preimage paths, the exact shape of the refused journal from 2026-09-03.
- `.../a_declared_target_with_no_projected_image_is_reported` RED: `undeclared postimage error = <nil>`.
- Restored: both GREEN (`ok ... internal/requeststate`).

**C-02 — terminal transitions clear any entry for the REQ.** Reverted `RemoveAllCheckpointClaims` back to `checkpointWithoutClaim(..., writerLabel)`.

- `TestFinalizeClearsACheckpointEntryLabelledByAnotherWriter` RED: `the archived REQ is still listed as in flight:` followed by the surviving `- REQ-740: ... — writer: vm:/other/checkout` entry. This is your field incident reproduced: typed success, archive complete, checkpoint still claiming the REQ.
- `TestCancelFromWorkingClearsAnyEntryForTheRequestAndSucceedsWithNone/entry_written_under_another_writer_label_is_cleared` RED: `the archived REQ still holds a checkpoint entry:` with `- REQ-306 ... writer: other:/repo` still present.
- `TestFinalizeCompletesWhenTheCheckpointHoldsNoEntryForTheRequest` stayed GREEN under the revert — see F-01.
- Restored: both GREEN.

**C-03 — the discard verb.** Removed the dispatch line and the `--discard-journal` parser case.

- All four subtests of `TestRecoverFinalizationDiscardsOnlyAPreMutationJournal` RED with `FINALIZATION-USAGE ... unknown recover-finalization option "--discard-journal"`.
- Restored: all four GREEN.

**C-04 — the refusal that must survive.** `TestFinalizationStillRefusesJournalImageSetsThatDisagree` (lifecycle and release subtests) passes both before and after the change, by design: it pins the `imageSetState` count check the REQ says must keep refusing. Each subtest hand-edits a journal to add one path to a postimage set under a recomputed image digest, then asserts the replay refuses with `FINALIZATION-LIFECYCLE-CONFLICT` / `FINALIZATION-RELEASE-CONFLICT` and that the unplanned path was never written. `imageSetState` itself is unchanged.

## Decisions

- **D-01 — fix the postimage side, not the plan side. DECIDE & STATE.** The two candidate fixes were "declare the unchanged checkpoint as a target" and "stop projecting a postimage for a path the plan never declared". The declared target set is the contract everything else already rests on (dirty checks, created directories, the commit-path allowlist, rollback), and declaring an unchanged file would push a no-op write back into the transaction — the exact failure REQ-468/REQ-538 fixed. So the projection follows the declaration. Reversible; touches one function.

- **D-02 — a terminal transition removes every checkpoint entry for its REQ. DECIDE & STATE.** Once the request file leaves `do-work/working/`, an "In Progress (interrupted)" entry for that REQ id is false no matter which writer wrote it, so the writer label has nothing left to protect at that moment. Two sessions cannot both hold a live claim on one REQ, and the mechanism already exists (`RemoveAllCheckpointClaims`, the same authority `recover-claim` uses). Claiming still writes and respects the label, and other REQs' entries are untouched — pinned by both changed tests. This is the field defect's root: your manifests passed `vm:/home/user/skill-do-work` where `advance` had written a different label, so the removal matched nothing and the skip was silent.

- **D-03 — a bounded discard verb rather than journal repair. DECIDE & STATE.** Repairing an image set means writing bytes into a journal that is the sole record of a half-finished transaction. Discarding is honest and checkable instead: `--discard-journal` refuses unless the phase is `prepared` and every recorded lifecycle preimage is still on disk byte-for-byte, which is precisely the pre-mutation state the REQ describes. Release preimages are deliberately not checked: at `prepared` the archive path's release preimage does not exist on disk yet, so checking them would refuse every legitimate discard.

- **D-04 — updated two tests from earlier REQs. DECIDE & STATE.** `TestCancelFromWorkingLeavesForeignOrAbsentCheckpointEntryAndSucceeds` (REQ-468 / REQ-538) asserted a foreign entry stays byte-identical; under D-02 that entry is removed. It is renamed to `TestCancelFromWorkingClearsAnyEntryForTheRequestAndSucceedsWithNone`, now asserts the REQ's own entry is cleared while an unrelated REQ's entry survives, and its "absent entry" subtest still pins the no-op-rollback failure those REQs fixed. `TestPlannedPostimagesPreservesRealFileModes` (REQ-489-era mode coverage) mutated a plan after `BuildPlan` without keeping `TargetPaths` in step; the mutations now update `TargetPaths` as `planTargets` would. Neither test lost coverage.

- **D-05 — new tests live in `finalization_req547_test.go`, not in `finalization_recovery_test.go` as the write_set listed. DECIDE & STATE.** The package's newer convention is a per-REQ file (`finalization_req499_test.go`, `_req512_`, `_req560_`); the one change inside an existing test file is the REQ-547 unit test that belongs beside the other `PlannedPostimages` tests.

- **D-06 — the discard's success finding carries no `next_argv`. DECIDE & STATE.** Re-running `finalize` needs a manifest an action authors, so inventing an argv with a placeholder path would hand automation something unrunnable. The finding's evidence says what to do; the verification argv proves the journal is gone.

## Findings

- **F-01 — the brief's RED case no longer reproduces in this tree, and that matters for how you read this REQ.** On today's `main`, `planCheckpoint` already blanked `CheckpointPath` when the checkpoint bytes would not change, so the count mismatch was unreachable through `BuildPlan` and a REQ with no checkpoint entry already finalized. I confirmed that by fixture before changing anything. What that mitigation did instead was convert the refusal into a silent skip — which is the failure you hit three times. The brief's line numbers (`state_plan.go:398`, `state_apply.go:156-158`) do not match this tree, consistent with the report coming from the released skill. Two things follow: the sets now agree by construction rather than by two conditions that happen to match, and the silent skip is gone.

- **F-02 — the silent no-match does share this seam.** The blanking of `CheckpointPath` is the single line that made the two image sets agree *and* made a writer-label mismatch invisible. Fixing it in scope was correct; the remaining skip (no entry for the REQ anywhere) is a true no-op.

## Discovered Tasks

- `internal/finalization` drops the lifecycle plan's `SkippedWork` on the floor: `requeststate.ApplyPlan` returns it, but `advanceJournal` builds its own result and never copies it, so `CHECKPOINT-ENTRY-NOT-PRESENT` never reaches a `finalize` caller. After this REQ the only surviving skip is benign, so nothing is being hidden today — but a future coupled-file skip would be equally invisible. impact-minor → report only.
- `internal/cleanup/cleanup_plan.go:232` (`ownedCheckpointRemoval`) still clears stale checkpoint claims only when the entry's writer label matches a locally derived `hostname:root`. The same staleness argument as D-02 applies to a REQ that is already archived, so cleanup can leave entries no verb will ever remove. impact-minor → report only.
- `do-work/CHECKPOINT.md` in the main tree still lists the three REQs archived during today's silent finalizations. This change stops new ones; it does not retro-clean those rows. impact-minor → report only.

## Integration Seams

- No change needed in `internal/heavyverification` or heavy-lane code; nothing in this REQ touches REQ-564's area.
- `--discard-journal` is documented in `prime-do-work-cli.md` only. No shipped action or `docs/` page mentions `recover-finalization` at all today, so there was no operator-facing page to amend without inventing one. If you want the verb named in `actions/run-with-recovery.md` or the work guide, that is a one-line addition the orchestrator owns.
- The `FINALIZATION-LIFECYCLE-CONFLICT` refusal still does not name `--discard-journal` as its next argv. REQ-514's contract is that a refusal never names itself as the fix, and I did not want to bend that rule from inside this REQ. Worth a decision if you want the exit discoverable from the refusal itself.

## Lesson Evidence

Candidate satellite bullet for `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, plus its `do-work/lessons-index.md` row — not written here, because the index lives under `do-work/`, which this worktree may not write, and the general crew rule requires both edits together.

- `[family: alternate-writer-contract-drift]` REQ-547: two sites that must agree on the same path set will drift, and the drift surfaces far away as a refusal nobody can clear. The lifecycle journal compared `planTargets`'s declared set against `PlannedPostimages`'s independently derived set; each was individually correct, and their disagreement refused finalization on every replay with no verb to clear it. Derive the second set from the first, and make an inconsistency an error at the point of projection rather than a conflict at replay. The second half of the lesson is the mitigation that preceded this REQ: silencing the mismatch by dropping the path from both sets made a writer-label no-match archive a REQ while the checkpoint went on claiming it was in flight — typed success, no finding, three occurrences in one day. When a no-op leaves shared state stale, silence is the defect; either do the work or report it.
