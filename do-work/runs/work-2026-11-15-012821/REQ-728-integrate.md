# REQ-728 integrator brief — run work-2026-11-15-012821

- REQ: REQ-728 ('ai-report-judge renders one named page of a bundle')
- REQ path: do-work/working/REQ-728-judge-renders-one-named-page.md
- UR: UR-176 (do-work/user-requests/UR-176/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/
- Shared guide (read it fully first): do-work/runs/work-2026-11-15-012821/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-15-012821/REQ-728-handback.md
- Operative name / branch: worktree-agent-REQ-728-judge-renders-one-named-page
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-728-judge-renders-one-named-page
- Dispatch instant (for the builder-work timing event): 2026-11-14T23:35:10Z
- Route: B
- tdd: true
- Closes UR-176: NO (UR_CLOSES=0); archive path do-work/archive/REQ-728-judge-renders-one-named-page.md
- Release: yes, next patch version; re-read VERSION right before prep-payloads
- Members already finalized in this run: REQ-729 (released 0.318.3, archived flat at do-work/archive/REQ-729-judge-relative-out-rule.md); REQ-733 (no release, archived at do-work/archive/UR-177/)
- Set aside: none. Every other claimed member is still building or waiting for its integrator; their working REQ files and worktrees are not yours.
- Last successful integration of the run (run the Restatement Sweep): NO
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-15-012821/REQ-728-probe.sh
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7ebdb4e0-bf88-4ef2-9c8a-adace1fd2dd4/scratchpad/integrate-REQ-728/
- Write your final report to do-work/runs/work-2026-11-15-012821/REQ-728-integration-report.md and add it, the hand-back, the review file and the progress log to commit_paths so the tree is clean when you finish.
- Expected conflict: REQ-729 already added a sentence at the end of the judgeUsage line of report_judge.go, and one paragraph in Step 7 of skills/do-work-toolbox/actions/ai-report.md. Keep the builder's [--page <file.html>] after <bundle-dir> AND REQ-729's sentence at the end of the line. Then read both merged files for a sentence that now contradicts the other member.
- Version: the builder proposes a minor bump (new --page flag). Use BUMP=minor with prep-payloads.py.
- The coordinator already committed every hand-back file, including yours. Leave it out of commit_paths unless you changed it.
- The probes run browser tests with DO_WORK_HEAVY_TESTS=1 and need Google Chrome (installed).
- After finalization, one extra check for the report: run the new ai-report-judge --page on ai-reports/2026-11-14_1547_proposal-tinka-review-surface with the REQ-733 page index-v1c03c01-*.html, write no file into the bundle, and put the verdict in your final message (REQ-733's integrator could not judge that page).
