# REQ-680 integration report (archive fetch test keeps an absent target absent)

REQ-680 is done and released as 0.305.95 ("Archive Fetch Test Proves a Failed Fetch Never Creates a Target That Did Not Exist").

- Archive: do-work/archive/REQ-680-archive-fetch-absent-target-test.md
- Run-artifacts commit 1dbe058e (pre). Merge 2414535d (full 2414535d0073764c550bc451b93750b1bdcf1afb). Finalization commit 37660ddd (phase cleanup_complete, blocked_paths and reason_codes empty).
- Change: `TestTotalFailurePreservesTheTargetAndLeavesNoScratch` is a two-row table (`pre-existing target`, `absent target`) in archive_fetch_test.go only. 63+/42-, no production change.
- Gate: DO_WORK_FAST_STAGE_REUSE=off maintainer-verify exit 0, wall 128 s, load 4.74, quiet machine. Green probe exit 0 (BLOCKED-PROBE-SUCCEEDED). Contract regressions after finalization: pass.
- Review (quick scan, general-purpose reviewer): 100%, Approve. Two nits, both impact-negligible, report only: N1 the failure message prints `<nil>` when the absent target was created; N2 `assertNoArchiveScratch` is called in both branches. Anti-bloat count: 0 new helpers, options, files, constants or test functions. Restatement sweep: nothing redefined.
- Heavy lanes (detached checkout of the merge, all exit 0, executed, none skipped): do-work-cli-integrations 65 s, staged-skills 39 s, updater 68 s, installer 44 s.
- Mutation proof (D-03, per the coordinator ruling): new row fails on its own assertion line; the pre-existing row and about 15 other staging tests fail too, as expected. Recorded in `## Testing`.
- Timing events: builder-work 1070 s, handback-merge 7 s, verification-gate (repository gate) 150 s, review 94 s, verification-gate (heavy drain) 226 s. Folded into `## Timing`.
- Lessons: none proposed, no satellite or lessons-index edit.
- Report-only discovered task: `TestMissingArchiveTargetParentReportsUnattemptedRoutesInTextAndJSON` shows SKIP in the normal package run because it calls `t.Skip("go-run integration is heavy-only")` (archive_fetch_test.go:512). Expected, nothing to fix.
- Cleanup: builder worktree removed, branch deleted with `git branch -d`, worktrees pruned, drain checkout removed.
- Left dirty under ROOT: only the sibling REQs' untracked hand-back, integrate and integration-report files (REQ-679, 681 to 686) and this report (untracked, for the coordinator's checkpoint commit). Nothing tracked is dirty.
- VERSION is now 0.305.95. The next integrator expects 0.305.96.
