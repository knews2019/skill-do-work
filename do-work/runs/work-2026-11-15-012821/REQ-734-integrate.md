# REQ-734 integrator brief — run work-2026-11-15-012821

- REQ: REQ-734 ('Build the proposal as a supports-and-limits table and compare it with the map')
- REQ path: do-work/working/REQ-734-proposal-table-and-map-comparison.md
- UR: UR-178 (do-work/user-requests/UR-178/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/
- Shared guide (read it fully first): do-work/runs/work-2026-11-15-012821/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-15-012821/REQ-734-handback.md
- Operative name / branch: worktree-agent-REQ-734-supports-limits-table-page
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-734-supports-limits-table-page
- Dispatch instant (for the builder-work timing event): 2026-11-15T01:06:43Z
- Route: B
- tdd: false
- Closes UR-178: YES (UR_CLOSES=1); archive path do-work/archive/UR-178/REQ-734-proposal-table-and-map-comparison.md
- Release: NO_RELEASE=1 (no shipped file changes expected; believe the finalizer if it disagrees)
- Members already finalized in this run: all ten: REQ-729 (0.318.3), REQ-733 (no release, do-work/archive/UR-177/), REQ-728 (0.319.0), REQ-730 (0.320.0), REQ-731 (0.321.0), REQ-723 (0.321.1, do-work/archive/UR-174/), REQ-724 (0.322.0, do-work/archive/UR-175/), REQ-732 (no release), REQ-726 (no release), REQ-725 (no release, do-work/archive/UR-175/). REQ-729, REQ-728, REQ-730, REQ-731, REQ-732 and REQ-726 are archived flat because UR-176 stays open.
- Set aside: none. Every other claimed member is still building or waiting for its integrator; their working REQ files and worktrees are not yours.
- Last successful integration of the run (run the Restatement Sweep): YES
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-15-012821/REQ-734-probe.sh
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7ebdb4e0-bf88-4ef2-9c8a-adace1fd2dd4/scratchpad/integrate-REQ-734/
- Write your final report to do-work/runs/work-2026-11-15-012821/REQ-734-integration-report.md and add it, the hand-back, the review file and the progress log to commit_paths so the tree is clean when you finish.
- UR-178 holds only this REQ, so it closes: archive under do-work/archive/UR-178/ and move the UR input.md in the same transaction.
- You are the run's LAST successful integration. Pass that fact and the ten sibling ids above to your reviewer so it runs the Restatement Sweep (skills/do-work/actions/review-work.md Step 6) across the whole run: shipped text that restates something a sibling redefined. Known candidates to check: every shipped mention of the ai-report-judge usage and its flags after REQ-728 and REQ-729 (for example skills/do-work/actions/help.md and skills/do-work-toolbox/actions/stakeholder-report.md); every mention of catalog.json fields after REQ-730; every mention of the catalog rebuild after REQ-731; every shipped mention of _dev/run-prompts/ after REQ-732; the toolbox and knowledge copies of the crew files after REQ-724. A stale shipped restatement is a fix: if the sweep finds one in a shipped file, this REQ would need a release, which it does not have. In that case do NOT widen this REQ: list each stale restatement with file and line under Discovered Tasks and in your final message, and the coordinator captures one follow-up REQ.
- The builder left four comparison pictures in the run directory (REQ-734-compare-*.png), already committed by the coordinator.
- The page was built on a base without REQ-730's catalog field. After the merge, run the ai-report index command on the merged tree and commit the regenerated ai-reports/catalog.json and ai-reports/index.html pair on main as the merge follow-up (as the REQ-730 and REQ-725 integrators did): the Tinka bundle's row must list v1c04.
- After the merge, judge the new page alone with ai-report-judge --page (output to your scratch folder, nothing written into the bundle) and delete the judge's temporary capture folder afterwards; put the verdict in Testing.
- The maintainer compares this page with the REQ-733 map page to choose a standard format (the blocked REQ-727 waits on that). Tell the reviewer to open both pages at 1440 and 390 width and report in plain words which one shows more of the proposal's relations on the first screen and what each one hides.
- Never edit anything under /Users/t2/.claude/.
- The coordinator already committed every hand-back file, including yours. Leave it out of commit_paths unless you changed it.
