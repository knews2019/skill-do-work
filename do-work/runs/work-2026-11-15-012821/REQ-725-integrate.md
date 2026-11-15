# REQ-725 integrator brief — run work-2026-11-15-012821

- REQ: REQ-725 ('Sort the local memory notes by destination')
- REQ path: do-work/working/REQ-725-sort-local-memory-notes-by-destination.md
- UR: UR-175 (do-work/user-requests/UR-175/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/
- Shared guide (read it fully first): do-work/runs/work-2026-11-15-012821/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-15-012821/REQ-725-handback.md
- Operative name / branch: worktree-agent-REQ-725-memory-note-destination-report
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-725-memory-note-destination-report
- Dispatch instant (for the builder-work timing event): 2026-11-14T23:32:30Z
- Route: A
- tdd: false
- Closes UR-175: YES (UR_CLOSES=1); archive path do-work/archive/UR-175/REQ-725-sort-local-memory-notes-by-destination.md
- Release: NO_RELEASE=1 (no shipped file changes expected; believe the finalizer if it disagrees)
- Members already finalized in this run: REQ-729 (0.318.3), REQ-733 (no release), REQ-728 (0.319.0), REQ-730 (0.320.0), REQ-731 (0.321.0), REQ-723 (0.321.1), REQ-724 (0.322.0, flat archive do-work/archive/REQ-724-move-global-standing-rules-into-skill.md), REQ-732 (no release, flat archive), REQ-726 (no release, flat archive)
- Set aside: none. Every other claimed member is still building or waiting for its integrator; their working REQ files and worktrees are not yours.
- Last successful integration of the run (run the Restatement Sweep): NO
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-15-012821/REQ-725-probe.sh
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7ebdb4e0-bf88-4ef2-9c8a-adace1fd2dd4/scratchpad/integrate-REQ-725/
- Write your final report to do-work/runs/work-2026-11-15-012821/REQ-725-integration-report.md and add it, the hand-back, the review file and the progress log to commit_paths so the tree is clean when you finish.
- UR-175 closes with this REQ: archive under do-work/archive/UR-175/, and move REQ-724 from its flat archive path into do-work/archive/UR-175/ in the same transaction, plus the UR input.md (put both REQ-724 paths and both input.md paths in commit_paths). Then run git grep -n 'archive/REQ-724' -- ':!do-work/' and re-point every hit to do-work/archive/UR-175/ (add those files to commit_paths).
- The memory folder changed after the build: the REQ-732 integrator deleted the note run-prompts-not-in-root-docs.md and its MEMORY.md line, because REQ-732 made that note wrong. The folder now holds 46 notes plus MEMORY.md; the report has 47 rows, and the probe counts the live folder, so it fails until the report is corrected. Fix on the builder branch: remove that note's row (or follow what the probe requires), correct the counts (46 notes, 28 delete), add one sentence saying the note was deleted on 2026-11-15 when REQ-732 landed, then re-merge. Record it as a decision.
- The builder's hand-back names four rows that assumed REQ-724 and REQ-732 would merge as scoped. Both are merged: re-check those four rows against main (for REQ-724 note that the toolbox and knowledge packages now carry copies of the new lines too) and correct any row that is now wrong, on the builder branch.
- The report is a decision brief the maintainer approves row by row: tell the reviewer to sample at least ten rows across the three destinations and check each named target file or 'already recorded at' path exists on main.
- The memory folder /Users/t2/.claude/projects/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/memory/ is read-only for you and your reviewer: never edit, add or delete a file there.
- REQ-730 regenerated ai-reports/catalog.json and ai-reports/index.html; this REQ adds a new bundle. After the merge, run the ai-report index command on the merged tree and commit the regenerated pair on main as part of the merge follow-up (as REQ-730's integrator did), so the catalog lists the new report; REQ-731's new --check flag should then print nothing for this bundle. Put that output in Testing.
- The coordinator already committed every hand-back file, including yours. Leave it out of commit_paths unless you changed it.
