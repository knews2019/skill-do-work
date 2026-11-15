# REQ-735 integrator brief — run work-2026-11-15-012821

- REQ: REQ-735 ('[impact-rule-change] Decide, rule by rule, which moved standing rules load in ordinary sessions')
- REQ path: do-work/working/REQ-735-rule-scope-for-moved-standing-rules.md
- UR: UR-179 (do-work/user-requests/UR-179/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/
- Shared guide (read it fully first): do-work/runs/work-2026-11-15-012821/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-15-012821/REQ-735-handback.md
- Operative name / branch: worktree-agent-REQ-735-moved-rules-always-on-scope
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-735-moved-rules-always-on-scope
- Dispatch instant (for the builder-work timing event): 2026-11-15T07:40:29Z
- Route: B
- tdd: false
- Closes UR-179: YES (UR_CLOSES=1); archive path do-work/archive/UR-179/REQ-735-rule-scope-for-moved-standing-rules.md
- Release: yes, next patch version; re-read VERSION right before prep-payloads
- Members already finalized in this run: the eleven first-pass members (see manifest.md) and REQ-736 (0.322.1, flat archive do-work/archive/REQ-736-help-judge-line-shows-page-option.md)
- Set aside: none. Every other claimed member is still building or waiting for its integrator; their working REQ files and worktrees are not yours.
- Last successful integration of the run (run the Restatement Sweep): YES
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-15-012821/REQ-735-probe.sh
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7ebdb4e0-bf88-4ef2-9c8a-adace1fd2dd4/scratchpad/integrate-REQ-735/
- Write your final report to do-work/runs/work-2026-11-15-012821/REQ-735-integration-report.md and add it, the hand-back, the review file and the progress log to commit_paths so the tree is clean when you finish.
- UR-179 closes with this REQ: archive under do-work/archive/UR-179/, and move REQ-736 from its flat archive path into do-work/archive/UR-179/ in the same transaction, plus the UR input.md (put both REQ-736 paths and both input.md paths in commit_paths). Then run git grep -n 'archive/REQ-736' -- ':!do-work/' and re-point every hit.
- You are the last integration of this pass: pass that fact and the sibling id REQ-736 plus the parent REQ-724 (do-work/archive/UR-175/) to your reviewer so it runs the Restatement Sweep. Known candidates: every shipped text that says where the approximate-numbers rule and the keep-existing-UI rule live (they move from coding-guardrails.md to communication-style.md), every package copy of the changed crew files (they must be identical), and docs or actions that cite coding-guardrails.md for those two rules.
- Pre-flight answered PREFLIGHT-BASELINE-RETRIED: the pre-build probe failed once and passed on the retry (load from a parallel agent). Record both exit statuses in Testing, as the finding asks.
- Version: the builder proposes minor; the REQ is impact-rule-change (two rules now load in every session). Use BUMP=minor unless the diff shows otherwise; record it.
- Instruction-file review: tell the reviewer to check each moved rule word for word against /Users/t2/.claude/CLAUDE.md (read-only), that no rule ships twice, that communication-style.md stays within its style (it is installed into consumer projects and imported by this repository's CLAUDE.md), and that decision D-03 holds: the model-per-agent rule has no copy or pointer in communication-style.md (the maintainer confirmed that twice).
- Trap from REQ-736: a CHANGELOG entry that cites a repository path the installed package cannot resolve makes contract-regressions.sh fail after finalization. Name files in plain words in the entry (for example: the always-on communication rules file), or use paths that resolve from the installed core package, and run _dev/tests/shipped-package-reference-contract.sh on the CHANGELOG.new payload text before finalization.
- The coordinator already committed the hand-back file. Leave it out of commit_paths unless you changed it.
- Never edit anything under /Users/t2/.claude/.
