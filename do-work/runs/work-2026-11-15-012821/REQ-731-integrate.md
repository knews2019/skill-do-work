# REQ-731 integrator brief — run work-2026-11-15-012821

- REQ: REQ-731 ('ai-report prints the catalog rebuild command when a new report is missing from the catalog')
- REQ path: do-work/working/REQ-731-report-prints-catalog-rebuild-reminder.md
- UR: UR-176 (do-work/user-requests/UR-176/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/
- Shared guide (read it fully first): do-work/runs/work-2026-11-15-012821/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-15-012821/REQ-731-handback.md
- Operative name / branch: worktree-agent-REQ-731-report-catalog-rebuild-reminder
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-731-report-catalog-rebuild-reminder
- Dispatch instant (for the builder-work timing event): 2026-11-14T23:35:12Z
- Route: B
- tdd: false
- Closes UR-176: NO (UR_CLOSES=0); archive path do-work/archive/REQ-731-report-prints-catalog-rebuild-reminder.md
- Release: yes, next patch version; re-read VERSION right before prep-payloads
- Members already finalized in this run: REQ-729 (0.318.3), REQ-733 (no release, do-work/archive/UR-177/), REQ-728 (0.319.0), REQ-730 (0.320.0, flat archive do-work/archive/REQ-730-catalog-lists-second-opinion-pages.md)
- Set aside: none. Every other claimed member is still building or waiting for its integrator; their working REQ files and worktrees are not yours.
- Last successful integration of the run (run the Restatement Sweep): NO
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-15-012821/REQ-731-probe.sh
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7ebdb4e0-bf88-4ef2-9c8a-adace1fd2dd4/scratchpad/integrate-REQ-731/
- Write your final report to do-work/runs/work-2026-11-15-012821/REQ-731-integration-report.md and add it, the hand-back, the review file and the progress log to commit_paths so the tree is clean when you finish.
- Expected overlap: REQ-730 already changed handleAIReportIndex, the usage string and the ends of report_index.go and report_index_test.go, and added a paragraph to skills/do-work-toolbox/actions/ai-report.md and the guide. Your builder added a --check branch, a usage change and new functions at the end of the same files. Resolve conflicts by keeping both members' code and text; then read every merged hunk and run the package tests before the gate.
- The builder's contract-regressions run failed only on the 30-second per-file budget at load 39; rerun it on a quiet machine and record the result.
- After the merge, run the new check against the merged tree for one bundle that is in the catalog and one that is not, and put both outputs in Testing.
- The coordinator already committed every hand-back file, including yours. Leave it out of commit_paths unless you changed it.
- Version: a new read-only flag; minor (BUMP=minor) matches how REQ-728 and REQ-730 were bumped in this run. Decide and record it.
