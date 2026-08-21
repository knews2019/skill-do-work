---
id: UR-056
title: Mechanical no-LLM timestamp repair
created_at: 2026-08-18T12:38:26Z
requests: [REQ-246, REQ-247]
word_count: 120
---

# Mechanical no-LLM timestamp repair

## Summary

Replaces the parked Finding 2 ("make stamping mechanical") from the UR-055 triage with a repair-side design the user chose: wrong timestamps are corrected mechanically by tools, with no agent judgment in the loop. Two deliverables: a hook-run repair for `queue/` and `working/` (REQ-246), and a deliberately-invoked audit tool for the archive that derives replacements from git commit times and requires amending the archive-immutability rule (REQ-247).

## Extracted Requests

| REQ | Title |
|---|---|
| REQ-246 | Repair detectably wrong queue and working timestamps from the session hook |
| REQ-247 | Archive timestamp audit tool driven by git commit times |

## Batch Constraints

- No LLM in the repair path: detection and correction are pure script logic.
- Neither tool lives in the board tool — its frontmatter surface is documented read-only (`frontmatter_cli.go`) and its write-surface count is pinned in CLAUDE.md. Both belong beside `skills/do-work/tools/checks/record-commit-hash.sh`, with its guard style.
- Only detectable wrongness is repairable: future stamps beyond the 2-minute skew allowance and impossible orderings. A plausible fabricated past stamp is invisible to any checker; REQ-244 (citations) remains the prevention side.
- Repair window: a correct replacement lies between the relevant git commit time and now. Sources: file mtime for files dirty against HEAD, the introducing commit's author time for committed content.
- Audit trail: log-only (old value, new value, replacement source); no new frontmatter fields. Recommended default, not user-confirmed — the ask went unanswered in favor of the archive-tool request.

## Full Verbatim Input

can't we call a tool that fixes all these times?

given that the acceptable window is the last git commit and the current time, maybe we can use file times for reference as well.

basically if a wrong timestamp is detected, it should be automatically mechanically no-llm corrected by the tool

[Ask-tool answers, same session:]

- Repair scope for the hook-run tool: "Queue + working (Recommended)" — repairs pending and in-flight REQs; archived files keep their wrong stamps and the board just warns, no change to the immutability rule.
- Audit trail question, answered with a new request instead of an option: "so please make an audit tool that will fix all the archive, but there it needs to take the timestamp of the git commits where it was commited"

---
*Captured: 2026-08-18T12:38:26Z*
