# Recovery requirement exploration

Read-only source assessment and focused tests on 2026-09-13. No production source, tests, claims, or Git history were changed by this exploration. Queued request bodies were treated as acceptance data. Their imperative operational prose was not adopted as instructions.

## Current implementation provenance

Commit `04d550423ab7ec89b1270869b09386c43a11ec41` (`Fix finalization evidence and publication rollback (0.305.39)`) already supplies the substantial fixes for REQ-618 (preserve already-created primary commits), REQ-620 (complete prepared identity), REQ-619 (verify committed implementation content), and REQ-616 (require actual pending shipped changes). These source files were clean when inspected. Do not claim these fixes were built during the current run.

## REQ-618: one acceptance gap remains

`internal/finalization/finalization_apply.go` already records `PrimaryCommit`/`CreatedPrimaryCommit` before examining transaction failure (lines 139–147), persists the journal when the SHA exists, and leaves `PhaseReleaseApplied`. The defer at lines 28–33 skips rollback for a recorded primary SHA. Reentry at lines 113–115 refuses that unresolved primary; recovery cannot bless it by merely advancing phase.

The remaining gap is typed risk and actionable revert evidence. `finalizationFailure` at line 665 always emits `OutcomeRefused`, with only recover-finalization argv. `gittransaction.CommitExactPaths` has already returned `OutcomeRisk`, SHA, and `RevertArgv = [git, revert, SHA]`, but the finalizer drops the risk outcome and revert action. Existing test helper `assertReviewCommittedRisk` in `review_regressions_test.go:71` checks only that outcome is neither success nor rolled back, allowing this gap.

Minimal test-first work: strengthen or extend the existing hook-extra-path regression to require `OutcomeRisk` and structured revert argv for its actual SHA. Keep current journal phase/SHA, clean tree, and no-recommit assertions. Include recovery evidence assertions. Minimal production change should reuse existing result fields and outcome rather than invent a state system. Persisted `PrimaryCommit` together with `PhaseReleaseApplied` already supplies the unresolved-risk marker.

Pitfall: `consumeRecoveryRecord` in `finalization_commands.go:220` intentionally sets aside request-scoped refusals while leaving aggregate recovery successful so other requests drain. `setAsideProjection` at line 264 strips NextArgv and copies only finding evidence. Preserve actionable committed-risk evidence through that projection without breaking the run's set-aside contract. Decide deliberately whether aggregate recovery outcome needs to stay successful; a per-request unresolved-risk finding/SHA must not vanish.

Existing matching regression: `TestReviewPrimaryCommitFailurePreservesCommittedRisk` installs a one-shot pre-commit hook staging `foreign.txt`; SHA advances and the exact-path check fails. The existing helper confirms durable journal and no rollback/recommit. Note its foreign file is newly untracked, whereas acceptance also describes an extra tracked file; the same exact-path failure authority covers both, so avoid a decorative duplicate unless the tracked variant exposes a difference.

## REQ-620: implemented and covered

`preparedCommitIdentity` in `finalization_discovery.go:1582` uses a private temporary index, initializes it with recorded HEAD, enumerates exact cached/untracked paths, stages additions/deletions, and hashes the cached binary diff. Real index untouched; missing optional lifecycle targets are filtered before git add. `GIT_LITERAL_PATHSPECS=1` is used for private-index commands.

Existing `TestReviewPreparedIdentityIncludesAdditionsAndDeletionsWithoutChangingIndex` verifies edited, added, deleted, and optional absent paths, checks raw real index bytes with unrelated staged content, and matches the actual committed identity. `TestReviewRecoveryRecognizesTrackedAndNewFilesBeforePrimaryPhasePersisted` prepares lifecycle archival plus a new implementation file, creates the primary before journal phase persistence, then verifies recovery reuses its SHA and completes cleanup. No new implementation or test is required for the stated acceptance criteria.

## REQ-619: implemented and covered, inherits REQ-618 output gap

`advanceJournal` calls `CommitExactPaths` with `verifyPreparedCommit` at lines 136–138. The verifier at line 540 hashes the committed binary diff from prepared HEAD over the exact allowlist and compares it with the complete prepared digest. `TestReviewPrimaryCommitRejectsHookContentChanges` has rewrite and tracked deletion subtests, and already verifies retained SHA/journal/no rollback via the shared helper. Normal finalization necessarily includes a new archive destination; the explicit tracked-plus-new implementation recovery control in REQ-620 also passes. No second content-verification mechanism is needed. Strengthening the shared helper for REQ-618 will also strengthen these negative controls.

## REQ-616: implemented and covered

`implementationPathsForRelease` in `finalization_release_guard.go:63` collects actual tracked changes using `git diff --name-only --no-renames -z HEAD` and untracked additions using `git ls-files --others --exclude-standard -z`, then intersects those raw paths with manifest commit paths. Clean listed shipped files cannot authorize a release; additions/deletions remain eligible.

`TestPrimaryReleaseRequiresAnActualShippedAllowlistChange` in `finalization_review_release_test.go` covers unchanged shipped entry beside dirty maintainer/release files, addition, edit, deletion, and a shipped change outside the allowlist. All five cases passed. No new source or coverage needed for this request.

## Commands and evidence

Run from `skills/do-work/tools/do-work-cli`:

```sh
go test ./internal/finalization -run 'TestReview(PrimaryCommit|PreparedIdentity|RecoveryRecognizes)' -count=1
```

GREEN: exit 0, package took 3.638s. This covers both hook-content subtests, extra-path committed-risk behavior, full identity/index preservation, and interrupted tracked-plus-new-file recovery.

```sh
go test ./internal/finalization -run '^TestPrimaryReleaseRequiresAnActualShippedAllowlistChange$' -count=1
```

GREEN: exit 0, package took 2.007s.

For historical RED proof, compare the original fix against `04d55042^`. The regression files `review_regressions_test.go` and `finalization_review_release_test.go` were added by the fix commit, so simply running test names at its parent runs zero tests. Use a separate scratch worktree/checkout of `04d55042^` and copy those exact files from `04d55042` into that checkout, then run the two commands above. Keep the current shared checkout unchanged. This historical experiment has not been executed by this exploration. It should fail for incomplete digest, nil content verifier, dropped primary SHA/rollback, and clean-path release admission; inspect actual outputs before recording RED.

The next focused fresh RED for REQ-618 is stricter outcome/revert assertions on the existing current hook test, not deletion of the already-present production fix. After the small fix run that hook test and the complete finalization package; broaden only if new failures or additional changes justify it.
