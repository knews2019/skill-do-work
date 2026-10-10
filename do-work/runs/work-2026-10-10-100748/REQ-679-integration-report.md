# REQ-679 integration report

REQ-679 (timeline browser probes use the fixed fixture; Previous and Next assertions return) is released as 0.305.98, "Timeline Browser Probes No Longer Depend on the Live Queue's Dates".

- Archive: `do-work/archive/REQ-679-timeline-probes-fixed-data.md` (flat path; UR-152 stays open).
- Commits: run artifacts `77b1d33f`, merge `0aab0216` (full `0aab0216f7d2aac86d5d56737dbe9d509373a85c`, builder commit `bab949ea`), finalization `3d7e5d5b`. Finalization record: `cleanup_complete`, no blocked paths, no reason codes. Worktree and branch removed; `contract-regressions.sh` passes after the move.
- Gate: `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` passed first run, 120 s, start load 4.03, no other gate running. Green probe passed through `advance` (both named browser tests passed, none skipped).
- Review: 96%, approve, no important findings. Three minor, all report only: F1 the clause (3) comment claims "no forecast past now" without proof; F2 the clause (6) comment still says the filtered width "is live queue data"; F3 no drawn-segments assertion on the trailing 7 days itself. Anti-bloat count: 0 new helpers, options, files or tests; 1 local constant, 2 struct fields, 2 fixture rows.
- Heavy lanes (detached checkout of the merge, Chrome set): queue-kanban-javascript 7 s, queue-kanban-browser 72 s, staged-skills 36 s; all executed, exit 0, none skipped.
- Timing events: builder-work, handback-merge, verification-gate (gate, probe, heavy plan), review, verification-gate (heavy drain). The heavy-drain event was recorded with a start about 1 minute late (11:38:40Z instead of 11:37:42Z), so its elapsed time is short by about a minute; the CLI has no end-time option.
- Lesson: builder bullet appended to `_dev/primes/lessons-kanban-board.md` (family `probe-reads-live-queue`, flat archive link); `do-work/lessons-index.md` row refreshed (6008 tokens, new family added).
- Dirty under ROOT: only untracked sibling files (REQ-680..686 hand-backs, integrate files and integration reports) and this report. Nothing of mine is staged.
- Report-only discovered tasks: (1) the typed-week rows use absolute 2026 dates, valid while the clock is after 2026-08-02; (2) `REQ-679-probe.sh` is green-only and runs in a few seconds; the red proof is the hand-back's scratch-worktree run.
- No refusal from rebuilt tools after the merge.
