# REQ-723 integrator brief — run work-2026-11-15-012821

- REQ: REQ-723 ('[impact-negligible] Archive-collision helper skips assets folders')
- REQ path: do-work/working/REQ-723-archive-collision-helper-skips-assets-folders.md
- UR: UR-174 (do-work/user-requests/UR-174/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/
- Shared guide (read it fully first): do-work/runs/work-2026-11-15-012821/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-15-012821/REQ-723-handback.md
- Operative name / branch: worktree-agent-REQ-723-archive-collision-skips-assets
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-723-archive-collision-skips-assets
- Dispatch instant (for the builder-work timing event): 2026-11-14T23:32:28Z
- Route: A
- tdd: true
- Closes UR-174: YES (UR_CLOSES=1); archive path do-work/archive/UR-174/REQ-723-archive-collision-helper-skips-assets-folders.md
- Release: yes, next patch version; re-read VERSION right before prep-payloads
- Members already finalized in this run: REQ-729 (0.318.3), REQ-733 (no release), REQ-728 (0.319.0), REQ-730 (0.320.0), REQ-731 (0.321.0, flat archive)
- Set aside: none. Every other claimed member is still building or waiting for its integrator; their working REQ files and worktrees are not yours.
- Last successful integration of the run (run the Restatement Sweep): NO
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-15-012821/REQ-723-probe.sh
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7ebdb4e0-bf88-4ef2-9c8a-adace1fd2dd4/scratchpad/integrate-REQ-723/
- Write your final report to do-work/runs/work-2026-11-15-012821/REQ-723-integration-report.md and add it, the hand-back, the review file and the progress log to commit_paths so the tree is clean when you finish.
- UR-174 holds only this REQ, so it closes: archive under do-work/archive/UR-174/ and move the UR input.md in the same transaction.
- No sibling touches the corehelpers package; a clean merge is expected. The builder branched before five releases, so the branch is behind main: that is normal, merge it.
- Version: a bug fix in a helper command; patch.
- The coordinator already committed every hand-back file, including yours. Leave it out of commit_paths unless you changed it.
