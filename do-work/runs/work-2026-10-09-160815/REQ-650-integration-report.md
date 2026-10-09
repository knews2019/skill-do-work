REQ-650 (the disk probe measures the repo root only) is integrated and released as 0.305.86.

- Archive: do-work/archive/REQ-650-the-disk-probe-measures-the-repo-root-only.md (flat; UR-143 stays open).
- Commits: run artifacts 54056e00 (pre), merge 674fa943ed0ce7ce37421ebbd73ffc7f03168963, finalization 1c4b5f8c4a8c491cc6340ac6805439250fa49b12. Builder worktree removed, branch deleted with -d, worktree list pruned.
- Qualify: satisfied, no findings. D-06 records the coordinator's ruling on generate_test.go, and write_set now names it.
- Gate: DO_WORK_FAST_STAGE_REUSE=off maintainer-verify exit 0 in 127 s on the first run. The probe passed in 3 s.
- Review: one reviewer agent gave 98% with Acceptance Pass. Two impact-negligible nits stay report only: verify.go builds repoRootReading before an error switch that does not use it, and generate_test.go seeds a second root finding. Restatement sweep found no stale text. The only hits are the historical 0.305.60 entries in both changelogs.
- Heavy drain at 674fa943 from a detached checkout with QUEUE_KANBAN_BROWSER set. Every lane was executed, none skipped: queue-kanban-javascript 8 s, queue-kanban-browser 89 s, staged-skills 45 s. All exited 0.
- Timing events recorded: handback-merge, verification-gate (repository gate), review, verification-gate (heavy drain). They are folded into ## Timing (8m 57s observed).
- Builder-work timing event: skipped on purpose. The hand-back landed at 16:16:41Z and I arrived at 16:35Z. The recorder has no end flag, so recording would have charged the builder about 24 minutes for a 4.7-minute build (work.md's landed hand-back rule).
- Lessons: a new disk-space-blind-spot bullet in lessons-do-kanban.md, plus the second-in-family Traps line in prime-do-kanban.md. The lessons-index row moved to 8338 tokens. The bullet links with the canonical GitHub URL, not the relative path the brief gave. The satellite's own header and work-reference require a URL because the file ships to consumers who never get do-work/archive/. contract-regressions.sh passed after finalization.
- Discovered task, report only: the builder tagged the 0.305.60 changelog note as impact-cosmetic, which is not an impact token. I reclassified it as impact-negligible. The entry stays as history.
- Dirty under ROOT: only do-work/runs/work-2026-10-09-160815/REQ-653-handback.md (untracked). It belongs to the sibling REQ-653 builder, so I left it alone. This report is also untracked, for the coordinator to commit.
