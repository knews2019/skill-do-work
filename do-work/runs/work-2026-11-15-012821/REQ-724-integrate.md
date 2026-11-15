# REQ-724 integrator brief — run work-2026-11-15-012821

- REQ: REQ-724 ('[impact-rule-change] Move the maintainer''s global standing rules into the skill')
- REQ path: do-work/working/REQ-724-move-global-standing-rules-into-skill.md
- UR: UR-175 (do-work/user-requests/UR-175/input.md)
- Run dir: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/
- Shared guide (read it fully first): do-work/runs/work-2026-11-15-012821/INTEGRATOR-GUIDE.md
- Hand-back: do-work/runs/work-2026-11-15-012821/REQ-724-handback.md
- Operative name / branch: worktree-agent-REQ-724-move-standing-rules-into-skill
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-724-move-standing-rules-into-skill
- Dispatch instant (for the builder-work timing event): 2026-11-15T00:01:58Z
- Route: B
- tdd: false
- Closes UR-175: NO (UR_CLOSES=0); archive path do-work/archive/REQ-724-move-global-standing-rules-into-skill.md
- Release: yes, next patch version; re-read VERSION right before prep-payloads
- Members already finalized in this run: REQ-729 (0.318.3), REQ-733 (no release), REQ-728 (0.319.0), REQ-730 (0.320.0), REQ-731 (0.321.0), REQ-723 (0.321.1, do-work/archive/UR-174/)
- Set aside: none. Every other claimed member is still building or waiting for its integrator; their working REQ files and worktrees are not yours.
- Last successful integration of the run (run the Restatement Sweep): NO
- Mid-run addendum text: none
- Do NOT run Step 10 (no loop, no `advance --checkpoint`, no `recover`).

Notes
- GREEN probe: do-work/runs/work-2026-11-15-012821/REQ-724-probe.sh
- Your private scratch folder: /private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7ebdb4e0-bf88-4ef2-9c8a-adace1fd2dd4/scratchpad/integrate-REQ-724/
- Write your final report to do-work/runs/work-2026-11-15-012821/REQ-724-integration-report.md and add it, the hand-back, the review file and the progress log to commit_paths so the tree is clean when you finish.
- UR-175 also holds REQ-725 (the memory-note report, still waiting for its integrator), so UR-175 stays open and the archive path is flat.
- Decision D-05 (the model-per-agent rule lives in background-agents.md, with no copy or pointer in communication-style.md) was confirmed by the maintainer through a question on 2026-11-15. It is no longer an escalation: write no pending-answers follow-up for it, and record the confirmation in Decisions.
- Version: minor (BUMP=minor). The REQ is impact-rule-change: it adds rules that change agent behavior in every consumer project.
- This is an instruction-file change. Tell the reviewer to check each added rule against the source /Users/t2/.claude/CLAUDE.md (read-only) for lost meaning, and against the skill's writing rules: conditions rather than closed lists, no machine-local paths or names, no rule shipped twice, no existing shipped rule reworded. Also run the Restatement Sweep for this REQ's own diff if it redefines something other shipped text restates.
- communication-style.md is installed into consumer projects and imported by this repository's CLAUDE.md: check that the contract tests naming that file still pass.
- The coordinator already committed every hand-back file, including yours. Leave it out of commit_paths unless you changed it.
