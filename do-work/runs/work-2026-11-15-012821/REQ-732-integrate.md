# REQ-732 integrator brief — run work-2026-11-15-012821

- REQ: REQ-732 ('[impact-negligible] Allow docs/run-prompts/ in the staged-skills contract and move the run prompts there')
- REQ path: do-work/working/REQ-732-allow-docs-run-prompts-folder.md
- UR: UR-176 (do-work/user-requests/UR-176/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/
- Shared guide (read it fully first): do-work/runs/work-2026-11-15-012821/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-15-012821/REQ-732-handback.md
- Operative name / branch: worktree-agent-REQ-732-allow-docs-run-prompts
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-732-allow-docs-run-prompts
- Dispatch instant (for the builder-work timing event): 2026-11-14T23:33:30Z
- Route: A
- tdd: false
- Closes UR-176: NO (UR_CLOSES=0); archive path do-work/archive/REQ-732-allow-docs-run-prompts-folder.md
- Release: NO_RELEASE=1 (no shipped file changes expected; believe the finalizer if it disagrees)
- Members already finalized in this run: REQ-729 (0.318.3), REQ-733 (no release), REQ-728 (0.319.0), REQ-730 (0.320.0), REQ-731 (0.321.0), REQ-723 (0.321.1), REQ-724 (0.322.0, flat archive)
- Set aside: none. Every other claimed member is still building or waiting for its integrator; their working REQ files and worktrees are not yours.
- Last successful integration of the run (run the Restatement Sweep): NO
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-15-012821/REQ-732-probe.sh
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7ebdb4e0-bf88-4ef2-9c8a-adace1fd2dd4/scratchpad/integrate-REQ-732/
- Write your final report to do-work/runs/work-2026-11-15-012821/REQ-732-integration-report.md and add it, the hand-back, the review file and the progress log to commit_paths so the tree is clean when you finish.
- UR-176 stays open (REQ-727 is blocked on the maintainer), so the archive path is flat.
- The builder moved six run prompts from _dev/run-prompts/ to docs/run-prompts/ with git mv. After the merge an empty _dev/run-prompts folder may stay on disk in the main checkout and make the probe fail: remove the empty folder with rmdir (never rm -r) and say so.
- Run the staged-skills heavy lane even if the planner does not select it: this REQ changes that lane's contract script, and the lane is not triggered by root docs/ changes (the builder's discovered task D-03).
- The builder says no shipped file changed (only _dev/, .gitattributes and docs/). Check that against the finalizer: if it answers that a release is required, believe it and switch to a patch release.
- The coordinator already committed every hand-back file, including yours. Leave it out of commit_paths unless you changed it.
