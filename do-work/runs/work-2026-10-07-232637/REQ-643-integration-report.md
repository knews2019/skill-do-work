# REQ-643 integration report (the board shows an earmarked pending REQ under Pending → Earmarked instead of Ready)

REQ-643 is finalized and released as 0.305.79. It closed UR-140.

- Archive: `do-work/archive/UR-140/REQ-643-board-earmarked-pending-subgroup.md`, status `completed`, `phase: cleanup_complete`, no blocked paths, no reason codes. REQ-644 and UR-140's `input.md` moved into `do-work/archive/UR-140/` in the same commit.
- Commits: run artifacts `fd5dbd3b` (= `<pre>`), first merge `1e6d9909`, review fix `fa508a23` on the builder branch, re-merge `5b2c3e64` (= recorded `commit:`), finalization `f378e0ff`. Range `fd5dbd3b..5b2c3e64`: 13 files, 121+/81-.
- TDD: builder RED (compile failure, then `pending ready = [REQ-562 REQ-563], want [REQ-562]`) then GREEN. I saw the same assertion RED on a scratch copy without the fix.
- Review: first review Pass, 92%. The wave-end restatement sweep (REQ-644 inherited, still agrees) found 8 stale "display only / Ready+Waiting" restatements. I fixed F1-F7 on the builder branch (D-05 in the REQ): tooltip, board.md:120, prime-do-kanban.md, generate.go/board.css comments, priority lines, work-guide.md. The re-merge had one conflict in work-guide.md:119 with REQ-644's sentence; I kept both. Delta re-review: Pass, 94%.
- Report only: F8 (no test pins a nonzero earmarked count or the renderer group), N1 (ADR-018 says "display-only"), N2/F9 (a long comment line in model.go), delta F8 (the tooltip shows on Waiting cards too; I judged the wording true). Builder discovered tasks (board.css comment, no JS behavior test) are impact-negligible, report only.
- Gate: `maintainer-verify.sh` exit 0 at both merges (117 s, 121 s). Probe green. Contract regressions green after finalization.
- Heavy lanes at `5b2c3e64` (detached checkout, removed): queue-kanban-javascript 7 s, queue-kanban-browser 81 s (37 tests, `QUEUE_KANBAN_BROWSER` set, not skipped), staged-skills 46 s. All exit 0, all executed.
- Timing: 7 events (builder-work, handback-merge, 2 review, 3 verification-gate) folded into `## Timing` (35m 24s observed).
- Lessons: one `paired-predicate-drift` bullet in `lessons-do-kanban.md`. Its link is the canonical GitHub URL, because the file's header says it ships and relative archive links break in installs. REQ-644's bullet in `lessons-action-files.md` now points at `archive/UR-140/`.
- Cleanup: builder worktree removed and branch deleted with `-d`.
- Left dirty, not mine: `do-work/working/baseline.json`, untracked `REQ-644-integration-report.md`, and this report. I did not change `manifest.md`.
