REQ-677 (toolbox journey-qa action) is integrated and released as 0.305.92.

- Archive: do-work/archive/REQ-677-journey-qa-action.md.
- Range: pre 68bd8bd7 (run artifacts). First merge 1cb75f9f. Final merge ec1746ee. Finalization 22d053bd (cleanup_complete, nothing blocked).
- Merge seam with REQ-676: two conflicts, both resolved by keeping both names: the toolbox argument-hint and the core help row. Both probes and shipped-package-reference passed before the merge commit.
- Review: 94% at 1cb75f9f, Acceptance Pass, Approve. I fixed F1 (guide detection line) and F3 (Step 6 cleanup could revert a tracked file; now only new untracked tool output, status with --untracked-files=all) on the builder branch (84738e35) and re-merged with the same pre. Re-check: 95%, Pass, no new finding. Browser-tool clause (D-02) fits the Agent Compatibility rule, but ui-review does not count a session browser. Report only: F2 (impact-rule-change, widen ui-review or drop the clause), F4 (impact-user-visible, ui-review has no pointer to journey-qa).
- Gate: exit 0 at 1cb75f9f (124 s) and at ec1746ee (126 s), each after load fell under 5; no load reruns. Probe exit 0. contract-regressions exit 0 after finalization (22 s), so the lesson link resolves.
- Heavy: six lanes, all executed at ec1746ee, exit 0: javascript 8 s, browser 85 s, cli-integrations 63 s, staged-skills 44 s, updater 68 s, installer 28 s.
- Update leg (after 22d053bd): install from e313e870 (v0.305.89) exit 0, journey-qa and source-audit absent. Update to 22d053bd through do-work-update.sh and a local archive server: exit 0, one download, v0.305.92. Both actions and guides arrive byte-identical, links and citations resolve. Brief, queue REQ and app file keep their sha256 after install and update.
- Timing (5 events, folded): handback-merge, verification-gate x3 (two gates, heavy drain), review 5m 48s. No builder-work event: the hand-back landed long before I started. Re-check not timed.
- Lesson: one bullet in _dev/primes/lessons-action-files.md (new family read-only-action-tool-side-effects); index row 7469 -> 7624 tokens. REQ-678's integrator must re-point this link to do-work/archive/UR-151/.
- Discovered tasks (report only): run-probe quiet-grep failure (fixed by c4dda51e); duplicated usage lines in the ui-review and stray-check guides.
- Left dirty under ROOT: REQ-675 and REQ-676 integration reports (not mine) and this report (untracked on purpose). Builder worktree and branch removed. Nothing pushed.
