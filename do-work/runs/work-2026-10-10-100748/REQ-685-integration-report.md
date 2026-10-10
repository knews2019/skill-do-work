# REQ-685 integration report

REQ-685 (transaction rollback reuses the open repository root; the no-root rollback path is deleted) is integrated and finalized.

- Release: 0.305.99, "Transaction Rollback Reuses the Repository Folder It Already Opened, and the Second Rollback Path Is Gone".
- Archive: do-work/archive/REQ-685-rollback-uses-open-root.md
- Commits: pre 33a0f1b1 (run artifacts), merge 96636dbb (full 96636dbb9ebd388fb648a947ed735f202d684955), finalization 26a452ab. Builder branch and worktree removed.
- Gate: `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` exit 0 at the merge, wall 134 s, load 2.90, one run. Probe REQ-685-probe.sh exit 0. contract-regressions.sh exit 0 after finalization (lesson links resolve; no lesson bullet was added).
- Review: Approve, 97%. Nil-root panic stays impossible (root opened at git_transaction.go:520, all ten rollbackFailure calls pass it). Collapsed test keeps its three assertions. Bloat: 30 non-comment added lines (19 are the same test body de-indented, 11 are the root argument at ten call sites plus the signature); no new helper, option, file or test. Restatement sweep: nothing still describes the old design.
- Finding F1 (minor, report only): the new comment at git_transaction.go:704-705 cites finalization_apply.go PrimaryCommit, which is right for exact_commit.go but wrong for this site. Its "HEAD" is read by cleanup_apply.go:144, doctor_repair.go:173 and publication_commands.go:234. The REQ prescribed this wording for both sites. A one-line comment fix would correct it; I did not make it.
- Heavy lanes (detached checkout of the merge, Chrome set): queue-kanban-javascript 9 s, queue-kanban-browser 82 s, do-work-cli-integrations 64 s, staged-skills 37 s, updater 64 s, installer 26 s. All exit 0, executed, none skipped.
- Timing events recorded: builder-work, handback-merge, verification-gate (gate), review, verification-gate (heavy drain); Timing section folded.
- Note: the SCOPE-DECLARED-NOT-TOUCHED warning (three Go identifiers in the Scope list read as files) is a parser false positive; recorded in Qualification.
- No new refusals from the merged transaction code; finalization succeeded first try (FINALIZATION-COMPLETE, cleanup_complete, no blocked paths).
- Left dirty under ROOT: only untracked run-directory files of other REQs plus this report (coordinator commits at checkpoint). do-work/working/baseline.json not dirty.
- Discovered tasks: none beyond F1.
