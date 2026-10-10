# REQ-690 integration report

- Release 0.305.111, "do-work status Shows Each In-Flight REQ's Class, ETA and Remedy in One Table". Archive `do-work/archive/REQ-690-run-status-action.md` (flat, UR-153 stays open). Finalization commit `8ab3d327`, cleanup_complete on the first try.
- Range: pre `bd8f3c9e` (run-artifacts commit). First merge `ada0ccab`. Re-merge `65820f45` after the review fixes (recorded as `commit:`). My commits on the builder branch: `d6945bef` (coordinator D-11 ruling, D-19) and `b5dd78f5` (review fixes, D-20).
- D-11 ruling: before the first merge I removed both unearned refusals (zero stale-claim threshold, and a `--run` path that is not a directory). The tests for the REQ's named failures are unchanged and pass.
- Merge: one conflict, in `resultmodel/result_model.go`. I kept REQ-660's Worktree block and this REQ's RunStatus block. I read every hunk of SKILL.md, fan-out-reference, lessons and main.go: none landed in the wrong place.
- Qualification: qualify and scope-drift passed. I added `run_status_render.go` to Scope and write_set (D-18). The scope check read a backticked `--watch` in my Scope line as a path, so I removed those backticks. The 2 UNWIRED warnings are false positives (a same-package file and a test file).
- Review: 87%, Approve. Fixed with the reviewer's text: F1 (C4 sent failed and collision statuses to clarify; they now go to `do-work forensics`, with a new test that fails before the fix), F2 (the status.md shell block now keeps the run-status exit status), F6 (the C6 line quotes the board's reason), F10 (deleted the dead `verificationArgv`). No second reviewer, because the delta is the reviewer's own text.
- Report only: F3 (a mistyped `--run` silently turns C3 rows into C7 rows; I did not apply the optional fix, per the D-19 ruling), F4 (C3 says `do-work run` while an integrator is running), F5 (no tests for C2, C8 or the class order), F7 (the builder branch shows a time, not an age), F8 (prime-do-kanban.md does not name `--format json`), F9 (forensics.md:10 and README.md:176 still send "stuck" to forensics), F11 nits.
- Gate: exit 0 at `ada0ccab` (113 s, load 1.91) and at `65820f45` (121 s, load 3.53). GREEN probe passed at both. Green record at `65820f45`.
- Heavy lanes at `65820f45`, all exit 0, none skipped: queue-kanban-browser 75 s, do-work-cli-integrations 57 s, staged-skills 35 s, updater 64 s, installer 25 s, queue-kanban-javascript reused (8 s at `ada0ccab`). The first drain at `ada0ccab` was also green.
- Timing: 7 events (2 merge, 4 gate and drain, 1 review), folded. No builder-work event, because the hand-back landed long before.
- Lessons: one `cross-module-fact-handoff` bullet in `lessons-do-work-cli.md`, linked to the flat archive path; REQ-689 re-points it. Index row: 20801 tokens. This also resolves the builder's index-size task (D-15). `contract-regressions.sh` passes.
- Cleanup: builder worktree removed, branch deleted with `-d`, worktrees pruned. Drain checkouts removed without force.
- Left dirty: only this report (untracked).
- Report-only builder tasks: the prime-do-kanban.md lines (same as F8), and the builder-branch `--no-merged` rule (D-06).
