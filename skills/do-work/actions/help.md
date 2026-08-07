# Core Help Action

> Invoked by bare `do-work`, `do-work help`, or `<core-command> help`. It explains the modular command boundary and stops without executing work.

## Full menu

```text
do-work — core request and queue orchestration

  do-work capture-request: <task>   Preserve intent as a UR and linked REQs
  do-work run [REQ|UR ...]          Build, test, review, archive, and commit ready work
  do-work verify-requests [REQ|UR]  Check captured requirements against their source request
  do-work review-work [REQ|UR]      Review completed implementation and acceptance evidence
  do-work clarify                   Resolve pending questions in a batch
  do-work abandon [REQ|UR] [why]    Cancel and archive work that should not be done
  do-work cleanup                   Consolidate completed archive records
  do-work commit                    Build focused, traceable commits
  do-work forensics                 Diagnose queue and archive integrity
  do-work roadmap [scope]           Summarize ready, blocked, stale, and completed work
  do-work version                   Show the installed version and recent releases
  do-work update                    Update the complete four-skill suite
  do-work recap                     Summarize recent URs and their REQs
  do-work help                      Show this menu

Extensions installed beside core:
  do-work-board       Kanban, Testing, calendar, summaries, and board CLI
  do-work-knowledge   BKB, memory, dream, interviews, and prompts
  do-work-toolbox     Reviews, reports, presentation, inspection, and repo utilities

Temporary compatibility:
  do-work pipeline <request>        Stateful end-to-end flow; removed after modular cutover
```

For `<core-command> help`, read that action's Input and When to Use sections and return no more than 15 lines: purpose, usage, accepted arguments, and two examples. Never execute the command while serving help.

## Rules

- Keep this menu aligned with the core router.
- Do not duplicate extension menus; point to the named sibling skill.
- Help never mutates the queue, project, or Git state.
