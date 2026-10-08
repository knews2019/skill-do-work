# REQ-644 integration report (capture separates a session earmark from work that needs the user at the keyboard)

REQ-644 is finalized and released as 0.305.78.

- Archive: `do-work/archive/REQ-644-capture-earmark-vs-operator-blocked.md` (flat, because UR-140 stays open for REQ-643). Status `completed`. `phase: cleanup_complete`, no blocked paths, no reason codes.
- Commits: run artifacts `cd91f838` (= `<pre>`), first merge `58168fe7`, review fix `54005faf` on the builder branch, re-merge `3602e584` (= recorded `commit:`), finalization `771876d4`. Range `cd91f838..3602e584`: 2 files, 4+/4-.
- Review: first review Pass, 95%. F1 (impact-negligible): the new text said "Needs Input · Blocked" but the board heading says "Needs input · Blocked". I fixed it on the builder branch, re-merged, re-gated, re-drained, and the delta re-review passed at 97%. F2 (impact-rule-change → report only): "`blocked_by` naming that person" could tempt an agent to write a real name when the user only says "me". F3 (impact-negligible → report only): a `blocked` REQ with an unmet `depends_on` shows under Pending → Waiting, not Needs input · Blocked.
- Gate: `maintainer-verify.sh` exit 0 at both merges (121 s, 124 s). Probe green (about 21 s each). Contract regressions green after finalization, so the lesson link resolves.
- Heavy lane `staged-skills`: executed, exit 0, 38 s at `58168fe7` and 35 s at `3602e584` (detached drain checkout, removed afterwards).
- Timing events: builder-work, handback-merge, 4 verification-gate (2 gates, 2 drains), and review. Folded into `## Timing` (14m 50s observed).
- Lesson: one `alternate-writer-contract-drift` bullet in `_dev/primes/lessons-action-files.md`. Its `do-work/lessons-index.md` row token count is now 6586.
- Cleanup: worktree removed and branch `worktree-agent-REQ-644-...` deleted with `-d`.
- Left dirty, and not mine: `do-work/working/REQ-643-board-earmarked-pending-subgroup.md`, `do-work/working/baseline.json`, and untracked `REQ-643-handback.md` (REQ-643's builder has handed back). This report is untracked. I did not update the `manifest.md` row for REQ-644, which still says "integrator starting".
- Discovered task (impact-negligible → report only): the addendum rule at `capture.md:143` handles only an addendum that names a session. It does not say how to flip a pending REQ to `blocked` when an addendum says the user must be present.
