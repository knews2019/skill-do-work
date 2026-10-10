REQ-678 (toolbox release-check action) is integrated and released as 0.305.93. UR-151 is closed: input.md and all four REQs (675, 676, 677, 678) are under do-work/archive/UR-151/.

- Archive: do-work/archive/UR-151/REQ-678-release-check-action.md.
- Range: pre 90ef7b4d (run artifacts). First merge 335c1739. Final merge 44548b59. Finalization 0bdfee78 (cleanup_complete on the first try, nothing blocked).
- Review: 95% at 335c1739, Acceptance Pass, Approve. I fixed F1 on the builder branch (2c586918): the verdict rule had no outcome when live acceptance is "not applicable", and now that case gives unknown. Re-merged with the same pre. Re-check: 96%, Pass, no new finding. Report only: F2-F6 (below).
- Wave-end sweep: every toolbox list (router, hint, both menus, README, toolbox_actions, tutorial, next-steps, Go registry) holds all three new actions. Still stale: shared-principles.md:15 (F3, rule-change), sample-archived-req.md:99 and review-work.md:460 (F4); no sibling redirects to the new actions (F2, user-visible). F5, F6 negligible nits.
- Fixture B, run cold by the reviewer: GREEN. Not ready, deployment failed (2.0.3 served, 2.1.0 packaged), live acceptance failed, 2026-10-03 note labeled historical.
- Gate: exit 0 at 335c1739 (140 s) and at 44548b59 (127 s). Load under 5 at both starts, no load reruns. Probe exit 0 at both.
- Heavy at 44548b59: staged-skills 45 s and browser 86 s executed. javascript, cli-integrations, updater and installer were reused from my first full drain at 335c1739 (8, 63, 65, 28 s). All exit 0.
- Batch update leg: installed from e313e870 (v0.305.89) with exit 0, and none of the three actions was present. Updated to 0bdfee78 through do-work-update.sh and a local archive server: exit 0, one download, v0.305.93. All six files byte-identical, 30 citations resolve, router lists all three, review-work.md has the stages. The brief, the queue REQ and the app file kept their sha256.
- Lessons: re-pointed the REQ-675, 676 and 677 links to archive/UR-151/. New REQ-678 bullet, family added-value-skips-partition-rule. Index row 7624 -> 7756. contract-regressions exit 0 (21 s), and all four links resolve.
- Timing (7 events, folded): builder-work 7m 36s, handback-merge x2 19 s, verification-gate x3 7m 11s, review 5m 06s. Not timed: the first drain and the delta re-check.
- Discovered task (report only, negligible): the brief calls the help menu's description column 33, but it starts at index 33 (column 34).
- Left dirty under ROOT: only this report, untracked on purpose. Worktrees and branch removed. Nothing was pushed.
