# REQ-723 integration report (the archive-collision helper skips assets folders)

- Release: 0.321.1 "Archive-Collision Check Ignores Evidence Files in Assets Folders" (patch: a bug fix in a helper command).
- Archive: do-work/archive/UR-174/REQ-723-archive-collision-helper-skips-assets-folders.md. UR-174 closed; its input.md moved to do-work/archive/UR-174/input.md in the same transaction.
- Range: pre 2e5b6f70 (run artifacts), merge 858c070867e9f4929cf0306f88bf2c1cf9aef3f3 (no conflicts). Finalization commit bb43e1af428c16acddfd9a31f80d1fa44a65454d (phase cleanup_complete, no blocked paths, no reason codes).
- Gate: green at the merge, 130 s, start load 3.47. GREEN probe passed through advance; green gate recorded. Contract regressions passed again after finalization.
- Review: Approve, 98% (opus reviewer, Route A quick scan). N1 nit: the test's local outputOf closure differs from the siblings' %+v format but earns its place → report only. Anti-bloat: 0.
- Heavy lanes at the merge, all executed, exit 0: do-work-cli-integrations 71 s, staged-skills 39 s, updater 70 s, installer 29 s (start load 4.87).
- Timing events: builder-work, handback-merge, verification-gate (gate), review, verification-gate (heavy drain); folded into ## Timing.
- Lessons: builder proposed none (second instance of the guard-on-one-entry-path family REQ-722 recorded); no satellite or lessons-index change.
- Cleanup: worktree removed, branch deleted (was 1be00b67), worktree prune done.
- Report-only discovered task: none.
- Left dirty in the main checkout: nothing (this report and the progress log are committed as "[REQ-723] integration report").
