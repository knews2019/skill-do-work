# REQ-736 integrator brief — run work-2026-11-15-012821

- REQ: REQ-736 ('[impact-negligible] Show the --page option on the ai-report judge line in help')
- REQ path: do-work/working/REQ-736-help-judge-line-shows-page-option.md
- UR: UR-179 (do-work/user-requests/UR-179/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/
- Shared guide (read it fully first): do-work/runs/work-2026-11-15-012821/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-15-012821/REQ-736-handback.md
- Operative name / branch: worktree-agent-REQ-736-help-judge-page-option
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-736-help-judge-page-option
- Dispatch instant (for the builder-work timing event): 2026-11-15T07:32:53Z
- Route: A
- tdd: false
- Closes UR-179: NO (UR_CLOSES=0); archive path do-work/archive/REQ-736-help-judge-line-shows-page-option.md
- Release: yes, next patch version; re-read VERSION right before prep-payloads
- Members already finalized in this run: the eleven first-pass members (see manifest.md); none in this pass yet
- Set aside: none. Every other claimed member is still building or waiting for its integrator; their working REQ files and worktrees are not yours.
- Last successful integration of the run (run the Restatement Sweep): NO
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-15-012821/REQ-736-probe.sh
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7ebdb4e0-bf88-4ef2-9c8a-adace1fd2dd4/scratchpad/integrate-REQ-736/
- Write your final report to do-work/runs/work-2026-11-15-012821/REQ-736-integration-report.md and add it, the hand-back, the review file and the progress log to commit_paths so the tree is clean when you finish.
- UR-179 also holds REQ-735 (still building), so UR-179 stays open and the archive path is flat.
- One changed line in skills/do-work-toolbox/actions/help.md; patch release.
- The builder could not run staged-skills-contract.sh (it runs only inside the heavy gate). Run the staged-skills heavy lane even if the planner does not select it.
- The coordinator already committed the hand-back file. Leave it out of commit_paths unless you changed it.
- Never edit anything under /Users/t2/.claude/.
