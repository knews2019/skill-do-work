# REQ-691 integration report

- Release 0.305.112, "do-work trace Shows How Much of a Spec Is Already Captured and Built Before Anything Is Captured Again". Archive `do-work/archive/REQ-691-trace-coverage-action.md` (flat, UR-153 stays open). Finalization `0bf90ad8`, cleanup_complete first try.
- Range: pre `00116871` (run-artifacts commit). First merge `d6aacf29`. Re-merge `a270e712` after review fixes (recorded as `commit:`). Review-fix commit on the builder branch: `775e24f9`.
- Merge: one conflict, the `skills/do-work/SKILL.md` argument hint, resolved as the union (REQ-690's `status [REQ] [--watch]` plus `trace <url|image|text|UR-NNN>`). `main.go` auto-merged; REQ-690's `open-work` case and this REQ's case both kept. Every merged hunk read; none misplaced.
- Qualification: qualify and scope-drift passed; 9 declared = 9 touched. One UNWIRED warning is a false positive (Go test file).
- Review: 93%, Approve, Pass. Fixed with the reviewer's text: F1 (main.go header comment names request-commits and now), F2 (prime-do-kanban.md sentence and Read-first line), F3 (Commit Evidence gloss: the newest row above 0 can be a release commit), F5 (payload tells capture to start the UR title with the ticket id; capture has no ticket rule), F7 (impact-low → impact-negligible in the REQ's Discovered Tasks). No second reviewer (the delta is the reviewer's own text).
- Report only: F4 (trace row sits above capture only, so phrases starting with check/review/status go to those rows; Plan D-07 kept), F6 (a repeated REQ id prints rows twice).
- Gate: exit 0 at `d6aacf29` (115 s, load 1.97) and at `a270e712` (122 s, load 2.61). GREEN probe passed at both. Green record at `a270e712`.
- Heavy lanes at `a270e712`, all executed, exit 0, none skipped: queue-kanban-javascript 7 s, queue-kanban-browser 74 s, staged-skills 33 s. First drain at `d6aacf29` also green.
- Timing: 7 events (2 merge, 4 gate and drain, 1 review), folded. No builder-work event (hand-back landed long before).
- Lessons: one paired-predicate-drift bullet at the top of `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`, canonical GitHub URL to the flat archive path; REQ-689 re-points it. Index row: 9043 tokens. `contract-regressions.sh` passes after finalization.
- Cleanup: builder worktree and drain checkouts removed without force, branch deleted with -d, pruned.
- Left dirty: only this report (untracked).
- Report-only builder tasks: capture.md Step 5 never shows the manifest `operation` value (REQ-688 covers it).
