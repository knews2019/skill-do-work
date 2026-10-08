REQ-646 (board activity no longer counts main's commits through a builder's merge of main) is integrated, released as 0.305.80, and archived.

- Archive: do-work/archive/REQ-646-board-activity-snapshots-direct-merge-matches.md (flat, because UR-141 still has REQ-645 open).
- Commits: run artifacts 84f91190 (pre), merge e3ef0109 (e3ef010941eba0d8fd22dd9d25db9ef2b978427c), finalization 016bd88b. The finalization record is cleanup_complete with empty blocked_paths and reason_codes.
- Release: the board package has no CHANGELOG or VERSION of its own, so the entry went to the root CHANGELOG.md and its mirrors, the same way 0.305.79 shipped. The REQ's write_set names skills/do-work-board/CHANGELOG.md, which does not exist.
- Qualification: the qualify gate was satisfied with no findings. Touched files are exactly the two Go files in write_set. I re-ran the new test against the pre-fix source in a throwaway checkout. It failed with "wrongly attributed main-side commit main1" and "5 instants, want 4", and it passes at the merge.
- Repository gate: exit 0 on the first run, 127 s. advance recorded green-gate satisfied, and the focused probe exited 0.
- Review: Approve, 96%, Acceptance Pass. One minor finding, F1, is marked impact-negligible → report only. It is pre-existing: a merge of main whose subject has a hand-written [REQ-NNN] token is still a direct match.
- Heavy lanes ran in a detached checkout of the merge with QUEUE_KANBAN_BROWSER set. All three executed and exited 0, with no HEAVY-RUN-LANE-SKIPPED: queue-kanban-javascript 8 s, queue-kanban-browser 79 s, staged-skills 39 s.
- Timing events: builder-work 274 s, handback-merge 19 s, the gate plus probe 152 s, review 110 s, heavy drain 143 s. All are folded into ## Timing. The gate and drain events include a few seconds of bookkeeping, because the recorder has no end flag.
- Lesson: a git-history-evidence bullet was added at 0.305.80 to lessons-do-kanban.md. Its index row went from 7896 to 8110 tokens. contract-regressions.sh passed after finalization.
- Cleanup: the builder worktree, its branch (deleted with -d) and the drain checkout are gone.
- Left dirty under ROOT, not mine: do-work/working/REQ-645-recovery-removal-requires-column-zero-heading.md, plus the untracked REQ-645-handback.md and REQ-645-integrate.md. This report file is also untracked, for the checkpoint. baseline.json is clean.
- Discovered tasks: none from the builder. The only report-only item is review finding F1.
