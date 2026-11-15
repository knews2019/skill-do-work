---
id: UR-179
title: 'Follow-ups after run work-2026-11-15-012821: scope of the moved rules, and the judge line in help'
created_at: 2026-11-15T07:25:23Z
requests: [REQ-735, REQ-736]
word_count: 231
---
## Summary
After the run that finished REQ-723 to REQ-734, the maintainer made two decisions through a question dialog on 2026-11-15. Both become follow-up REQs. The input is the coordinator's summary of those two decisions.

## Extracted Requests

| REQ | Request |
| --- | --- |
| REQ-735 | Decide, rule by rule, which of the three rules moved by REQ-724 (moving the maintainer's global standing rules into the skill) get a home that loads in ordinary sessions: approximate numeric targets, never remove existing UI while redesigning, and pick the model per agent. Respect the maintainer's earlier answer that the model rule lives only in `background-agents.md`, or say plainly why a pointer is now needed. Fold in the REQ-724 review findings F3 and F5 if they fit. |
| REQ-736 | Add the `--page` option to the `ai-report judge <dir>` line in `skills/do-work-toolbox/actions/help.md` (line 18). |

## Batch Constraints
- The two REQs are independent. Neither depends on the other.
- Both change shipped files under `skills/`, so each is a release.
- Deleting `~/.claude/CLAUDE.md` is the maintainer's own step after REQ-735. No REQ reads, edits or deletes that file.

## Full Verbatim Input
> ```
> Two follow-ups after the run that finished REQ-723 to REQ-734.
> 
> 1. Rule scope after REQ-724 (moving the maintainer's global standing rules into the skill, released 0.322.0). Three moved rules load only during implementation or an agent fan-out: approximate numeric targets, never remove existing UI while redesigning, and pick the model per agent. If the maintainer deletes ~/.claude/CLAUDE.md they stop applying in ordinary sessions. The maintainer chose: capture a follow-up REQ that decides, rule by rule, which of the three move to an always-on file, so that ~/.claude/CLAUDE.md can then be deleted. Note the maintainer's earlier answer in the same run: the model-per-agent rule itself lives in background-agents.md, with no copy in communication-style.md; the follow-up must respect that or say plainly why a pointer is now needed. Related open review findings from REQ-724 to fold in if they fit: the new ask-tool rule and model rule do not say that an action's own rule wins over them (F3); the rule group about several sessions in one checkout is only partly shipped, 'never revert files you did not change' has no home in the skill (F5). Source: do-work/runs/work-2026-11-15-012821/REQ-724-integration-report.md and do-work/archive/UR-175/REQ-724-move-global-standing-rules-into-skill.md.
> 
> 2. help.md fix after REQ-728 (ai-report-judge renders one named page, released 0.319.0). The run-end restatement sweep found skills/do-work-toolbox/actions/help.md line 18: 'ai-report judge <dir>' does not show the new --page option. The maintainer chose: capture a follow-up REQ for this one line. Source: do-work/runs/work-2026-11-15-012821/REQ-734-integration-report.md.
> ```
