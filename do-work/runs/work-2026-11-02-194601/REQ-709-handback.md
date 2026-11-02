# REQ-709 hand-back (recover ignores do-work notes no request record can own)

## Summary

`recover` and `recover-finalization --discover` no longer refuse on a file under `do-work/` that sits outside every request-record location, such as an untracked inbox note or a modified handoff note. `sharedFinalizationPath` now states that condition once. Discovery drops such paths from its evidence set, and the same narrowed predicate narrows the apply-time commit check at `finalization_apply.go:525`. Record-location dirt next to a note still refuses with the same code and message, and the refusal now names only the record path.

## Files

- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` (modified): the predicate is narrowed to the record locations, the new `unownedDoWorkNote` helper is added, and discovery's `dirty` set skips unowned notes (each side of a rename is judged alone).
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_unowned_notes_test.go` (created): `TestRecoverFinalizationIgnoresNotesNoRequestRecordOwns` (3 subtests) and `TestRecoverFinalizationStillRefusesRecordStateNextToANote` (9 subtests).
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go` (modified): adds `TestRecoverLeavesAnUntrackedInboxNoteForCommit` and imports `corehelpers`.

## Commits

Branch `worktree-agent-REQ-709-recover-ignores-unowned-do-work-notes`:

- `4634893c` [REQ-709] recover ignores do-work notes no request record can own

## Red-Green Evidence

The RED runs used the new tests against unchanged production code (`d5657887`). The exact failure text is saved at `/tmp/req709-builder/red-finalization.txt` and `/tmp/req709-builder/red-lifecycleadvance.txt`.

| Test | RED (before the fix) | GREEN |
|---|---|---|
| `TestRecoverLeavesAnUntrackedInboxNoteForCommit` (the REQ's RED case) | `recover refused on a note no request record owns`: `Outcome:"refused"`, `ReasonCodes:{"FINALIZATION-DISCOVERY-AMBIGUOUS"}`, `BlockedPaths:{"do-work/inbox/2026-10-28_do-work-upstream-suggestion-month-folder-report-layer.md"}` | PASS (0.12s) |
| `...IgnoresNotesNoRequestRecordOwns/no_terminal_request` | `an unowned note stopped discovery`: refused `FINALIZATION-DISCOVERY-AMBIGUOUS`, BlockedPaths `{inbox note, "do-work/working-notes.md"}` | PASS |
| `.../legacy_tail_recovers_around_an_untracked_inbox_note` | `an unowned note blocked legacy recovery`: refused `FINALIZATION-DISCOVERY-AMBIGUOUS`, BlockedPaths `{inbox note}` | PASS: `cleanup_complete`, note still `??` and not in `commit_paths` |
| `.../legacy_tail_recovers_around_a_modified_tracked_handoff_note_its_summary_names` | `a tracked handoff note blocked legacy recovery`: refused `FINALIZATION-DISCOVERY-AMBIGUOUS`, BlockedPaths `{"do-work/HANDDOWN-UR-1.md"}` | PASS: `cleanup_complete`, note still ` M` and not in `commit_paths` |
| `...StillRefusesRecordStateNextToANote/foreign_claim_line_in_the_checkpoint` | `refusal record = ...BlockedPaths:{"do-work/CHECKPOINT.md", <inbox note>}` (the refusal named the note too) | PASS: refuses, BlockedPaths exactly `{"do-work/CHECKPOINT.md"}` |
| `.../untracked_<record path>` × 8 (queue, working, archive, user-requests, `.req-reservations`, runs, CHECKPOINT.md, calibration-log.tsv) | each refusal listed `{<record path>, <inbox note>}` | PASS: each refuses `FINALIZATION-DISCOVERY-AMBIGUOUS` naming only the record path |

**Differential for each half of the fix** (production file restored afterwards and checked with `cmp`):

- A. Narrowed predicate, no `dirty` filter: only the handoff subtest fails: `recovery associated the note with REQ-750: [do-work/HANDDOWN-UR-1.md req-750.txt]`. The filter is what keeps a note out of a request's group.
- B. `dirty` filter on, old predicate restored: both legacy-tail subtests fail at the apply-time check with `FINALIZATION-AMBIGUOUS-SHARED-STATE` on the note. The narrowed predicate is what clears `finalization_apply.go:525`.

## P-A-U

- **[PLAN]:** Narrow `sharedFinalizationPath` to the condition "where request records live", keyed on the first path segment under `do-work/`. Add `unownedDoWorkNote` (under `do-work/` and not a record location) and skip such paths when building discovery's `dirty` set, judging path and rename origin separately. Every other user of the predicate (owners, unresolved-shared loop, claim-only topology, `ambiguousSharedRemainder`, apply-time check) inherits the narrowing. There is no new package, no exported API, and no change to refusal codes or messages. Lessons did not conflict with this plan.
- **[APPLY]:** Scope was exactly the 3 files above. No `do-work/` path in the worktree was written. The inventory chokepoint `corehelpers/inventory.go` is untouched (D-02).
- **[UNIFY]:** `git diff --stat main...HEAD`: `finalization_discovery.go | 39 ++++++-`, `finalization_unowned_notes_test.go | 114 +++`, `recovery_commands_test.go | 29 +++`, 3 files changed, 179 insertions(+), 3 deletions(-). `gofmt -l ./internal/finalization ./internal/lifecycleadvance` printed nothing. `go vet ./internal/finalization/ ./internal/lifecycleadvance/` was clean. Files checked:
  - `finalization_discovery.go`: read the whole diff. The predicate keeps the bare `do-work` path and every record location as evidence. The doc comment states the condition and why. There is no debug output. Every remaining caller was reviewed: the owners map, the line-163 and line-187 loops, `coherentClaimOnlyTopology`, `soleReleaserSharedPath`, `ambiguousSharedRemainder`, `associateReleaseMetadata`, and `finalization_apply.go:525`.
  - `finalization_unowned_notes_test.go`: each test has a one-line comment naming the UR-161 failure it pins. It uses only existing package helpers. Its run time is about 2.5s.
  - `recovery_commands_test.go`: the new test calls `handleRecover` with no flags, which matches the `recover` argv in `actions/commit.md` Step 1. It asserts the `corehelpers.ReadProtectedInventory` rows exactly.
  - Klarna check: the user's goal is that `do-work commit` can commit the note. The test confirms `recover` succeeds and the note is still an `A` inventory row. Association was read in code, not tested: `corehelpers.AssociateProjectPaths` always leaves `do-work/` paths unowned, so the commit action groups the note as an ordinary change. The `do-work commit` action itself (natural-language steps) was not run.

## Lessons Read

Source: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, family-targeted with a fixed-string grep. The prime `## Traps` in `prime-do-work-cli.md` was read in full.

- `[family: closed-enumeration-for-a-condition]`: lines 29 (REQ-635), 53 (0.303.7), 92 and 95 (REQ-460), 99 (REQ-461), 116 (0.283.0). Applied in three ways. The predicate states its condition in the doc comment and is an affirmative set of evidence locations, not a deny-list of note folders. The test table covers every ingredient of the condition plus a prefix-sibling boundary case (`do-work/working-notes.md`). Each half of the two-part fix was proven separately by differential.
- `[family: observed-subset-is-not-semantic-completeness]`: line 90 (REQ-498). This is unchanged in effect. Record-location paths still need the same semantic proof. Notes are now never admitted into a recovery group, so no byte is committed without ownership evidence.
- Nothing missing.

## Decisions

- **D-01** (orchestrator): The condition is "where request records live", stated once in code. That means `do-work/queue/`, `working/`, `archive/`, `user-requests/`, `.req-reservations/`, `runs/`, `do-work/CHECKPOINT.md` and `do-work/calibration-log.tsv`. No note folders are listed. Done in `sharedFinalizationPath`, whose doc comment states the condition and why.
- **D-02** (orchestrator): The rule lives in the finalization package, not in `untrackedDoWorkMetadata`, because `do-work commit` must still see the note in the inventory. The lifecycleadvance test confirms the note remains an `A` row.
- **D-03** (orchestrator): Discovery neither refuses on such a file nor associates it with a request. It is dropped from `dirty` (path and origin judged alone), and the narrowed predicate also narrows `finalization_apply.go:525`. **Consequence:** other non-record files directly under `do-work/` (for example `do-work/lessons-index.md`) also stop being evidence. A legacy tail now settles around them and leaves them dirty for `do-work commit`, where it refused before. This holds even when the REQ's Implementation Summary names the file (pinned by the handoff subtest).
- **D-04** (DECIDE & STATE): A path whose first segment under `do-work/` is a record location counts as evidence even with no further segment (for example `do-work/runs` as a file or symlink). The bare `do-work` path also stays evidence. Reason: anything that is or replaces a record location keeps today's refusal, which is the conservative reading of "exactly as today".
- **D-05** (DECIDE & STATE): The name `sharedFinalizationPath` was kept, not renamed. Reason: it keeps the diff to the predicate body and its callers unchanged, and the new doc comment carries the meaning.
- **D-06** (DECIDE & STATE): The finalization tests went into a new file, `finalization_unowned_notes_test.go`, and the public test calls `handleRecover`, not the built CLI binary. Reason: the new file's time stays easy to measure (about 2.5s, against the 30s budget). `handleRecover` is the registered handler the binary dispatches to, and the existing heavy-lanes test in the same file uses it the same way.
- **D-07** (DECIDE & STATE): Fixtures use the real UR-161 note filename instead of `note.md`. Reason: the fixture then matches the failure that was actually observed.
- **D-08** (DECIDE & STATE): The tests and the fix are in one commit, with no RED-only commit. Reason: the repository has no RED-commit convention, and a red commit breaks bisect. The RED text is recorded above and in `/tmp/req709-builder/`.
- **D-09** (DECIDE & STATE): No prose change. Reason: a grep of shipped `.md` files found no text that states the old "every do-work path is shared state" rule. `actions/commit.md` Step 1 wording is still true.

## Discovered Tasks

- impact-low: The record-location list in `sharedFinalizationPath` repeats the `do-work/` layout that `internal/repositorymodel/repository_model.go`'s walker encodes with its own literals, and there is no shared declaration. A record folder added to the walker but not to this predicate would be treated as a note: recovery would leave it dirty instead of refusing. That direction fails open. The 8-row table test in `finalization_unowned_notes_test.go` pins today's set → report only.

## Test Commands

All run from `skills/do-work/tools/do-work-cli` in the worktree:

- `go test ./internal/finalization/ ./internal/lifecycleadvance/ -count=1` before any change: `ok finalization 80.188s`, `ok lifecycleadvance 33.814s`.
- `go test ./internal/finalization/ -run 'TestRecoverFinalizationIgnoresNotesNoRequestRecordOwns|TestRecoverFinalizationStillRefusesRecordStateNextToANote' -count=1` before the fix: FAIL, with 3 + 9 subtests red as listed above.
- `go test ./internal/lifecycleadvance/ -run 'TestRecoverLeavesAnUntrackedInboxNoteForCommit' -count=1` before the fix: FAIL (refused `FINALIZATION-DISCOVERY-AMBIGUOUS`).
- The same two commands after the fix, plus the existing `TestRecoverFinalizationIgnoresFinderMetadataUnderDoWork` and `TestRecoverFinalizationStopsWhenTheRefusalOwnsNoRequest`: all PASS.
- Differentials A and B (described above): each failed only the expected subtests.
- `gofmt -l ./internal/finalization ./internal/lifecycleadvance`: empty.
- `go vet ./internal/finalization/ ./internal/lifecycleadvance/`: clean.
- `go test ./internal/finalization/ ./internal/lifecycleadvance/ -count=1` after the fix: `ok finalization 60.810s`, `ok lifecycleadvance 26.936s`.
- Not run, as the brief instructs: `_dev/tests/maintainer-verify.sh` and the heavy lanes.

## Addendum (review fix F1)

**Commit:** `6582e1d4` [REQ-709] review fix F1: pin claim-only and sole-releaser recovery beside a note. This is a test-only change to `skills/do-work/tools/do-work-cli/internal/finalization/finalization_unowned_notes_test.go` (61 insertions, 15 deletions). Production code is unchanged.

**New subtests** (inside `TestRecoverFinalizationIgnoresNotesNoRequestRecordOwns`):

- `sole-releaser legacy tail recovers around an untracked inbox note`: `seedSoleReleaserLegacyTail` plus the note, run with `--discover --assume-sole-releaser`. It asserts `cleanup_complete` and that `commit_paths` is exactly `[do-work/CHECKPOINT.md do-work/calibration-log.tsv owned.txt]`. It asserts that the `FINALIZATION-SOLE-RELEASER-ATTRIBUTED` paths are exactly `[do-work/CHECKPOINT.md do-work/calibration-log.tsv]`, and that the note is still `??`.
- `claim-only topology stays a claim beside an untracked inbox note`: the `TestRecoverFinalizationAcceptsCoherentClaimOnlyTopology` fixture plus the note. It asserts success, zero finalizations, and exactly one finding, `FINALIZATION-NONE`. It also asserts that `do-work/working/REQ-699.md` is byte-identical and the note is still `??`.

**Tightened subtests:** both existing legacy-tail subtests now compare `commit_paths` exactly through the new helper `assertSettledAroundNote`. The semantic tail is compared against `fixture.ownedPaths`, as returned by `seedSemanticLegacyTail`. The handoff tail is compared against `[req-750.txt]`. The old "note not in commit paths" loop was replaced by this exact comparison and by `assertNoteUntouched`.

**RED** was run against the pre-change `finalization_discovery.go` from `c5ae2d8e`, which is byte-identical to `d5657887`. It ran in the scratch module copy `/tmp/req709-builder/module-old/`; the worktree was not touched. Full output is in `/tmp/req709-builder/f1-red.txt`.

- sole-releaser: `a note blocked legacy recovery`: `Outcome:"refused"`, `ReasonCodes:{"FINALIZATION-DISCOVERY-AMBIGUOUS"}`, `BlockedPaths:{"do-work/inbox/2026-10-28_do-work-upstream-suggestion-month-folder-report-layer.md"}`.
- claim-only: `a note turned a claim-only topology into a refusal`: `Outcome:"refused"`, `ReasonCodes:{"FINALIZATION-DISCOVERY-AMBIGUOUS"}`, `BlockedPaths:{"do-work/CHECKPOINT.md", <inbox note>, "do-work/queue/REQ-699.md", "do-work/working/REQ-699.md"}`.
- The two tightened legacy-tail subtests are also RED on the old file (`a note blocked legacy recovery`, refused `FINALIZATION-DISCOVERY-AMBIGUOUS`). On the old code that refusal comes before the exact `commit_paths` check can run.
- To show the exact `commit_paths` check works by itself, it was also run in the scratch module against the fixed file with the `dirty` filter disabled. The result was `commit paths = [do-work/HANDDOWN-UR-1.md req-750.txt], want [req-750.txt]` (`/tmp/req709-builder/f1-variantA.txt`).

**GREEN** on the branch:

- All 5 subtests of `TestRecoverFinalizationIgnoresNotesNoRequestRecordOwns` and all 9 of `TestRecoverFinalizationStillRefusesRecordStateNextToANote` PASS.
- `gofmt -l ./internal/finalization ./internal/lifecycleadvance` printed nothing.
- `go vet ./internal/finalization/ ./internal/lifecycleadvance/` was clean.
- `go test ./internal/finalization/ -count=1`: `ok finalization 56.215s`.

**Decisions:**

- **D-10** (DECIDE & STATE): The claim-only subtest asserts exactly one `FINALIZATION-NONE` finding instead of "no findings". Reason: `FINALIZATION-NONE` is the normal success finding when nothing is left to finalize. Asserting it exactly separates "succeeded with nothing to do" from any other success shape.
- **D-11** (DECIDE & STATE): The sole-releaser subtest uses `seedSoleReleaserLegacyTail`, not the semantic tail. Reason: that tail needs `--assume-sole-releaser` to settle at all, because its checkpoint and calibration rows are attributed only under that flag. It therefore exercises the `soleReleaserSharedPath` branch the note used to reach. The expected attributed set is copied from the existing REQ-499 test, so "settles exactly as today" is checked against that test's own contract.
