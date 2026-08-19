# Core Help Action

> Invoked by bare `do-work`, `do-work help`, or `<core-command> help`. It explains the modular command boundary and stops without executing work.

## Full menu

```text
do-work — core request and queue orchestration

  do-work capture-request: <task>   Preserve intent as a UR and linked REQs
  do-work run [REQ|UR ...]          Build, test, review, archive, and commit ready work
  do-work verify-requests [REQ|UR|--against source ...]
                                      Check capture quality or revalidate queued work after a reversal
  do-work review-work [REQ|UR]      Review completed implementation and acceptance evidence
  do-work clarify                   Resolve pending questions in a batch
  do-work abandon [REQ|UR] [why]    Cancel and archive work that should not be done
  do-work cleanup                   Consolidate completed archive records
  do-work commit                    Build focused, traceable commits
  do-work forensics                 Diagnose queue and archive integrity
  do-work roadmap [scope]           Summarize ready, blocked, stale, and completed work
  do-work handoff                   Hand off to a fresh session with a paste-ready restart prompt
  do-work version                   Show the installed version and recent releases
  do-work update                    Update the complete four-skill suite
  do-work recap                     Summarize recent URs and their REQs
  do-work help                      Show this menu

Extensions installed beside core — run <package> help for usage on any of these:

  do-work-board      board [serve|static|summary|cli]
  do-work-knowledge  bkb · memory · dream · interview · prompts · setup-memory
  do-work-toolbox    validate-feedback · code-review · maintainability-audit
                     ui-review · ai-report · present-work · present-video · slop-check
                     quick-wins · scan-ideas · deep-explore · prime · inspect
                     note · stray-check · tidy-repo · tutorial · install
```

## Full cycle without persistent state

`do-work run` already owns implementation, testing, and review for every REQ, so a full cycle composes the public commands without a separate testing stage or resumable state file. Copy this prompt and replace the final placeholder:

```text
Use the installed do-work suite to complete this request end to end:

1. Use do-work to capture the request below and record the resulting UR ID.
2. Run do-work verify-requests for that UR. Stop and report if verification fails.
3. Run the UR's REQs through do-work run. Require its built-in tests and review to pass.
4. Use do-work-toolbox ai-report for the same UR.
5. Report the implementation, tests, decisions, and deliverable paths.

Request:
<paste request here>
```

For `<core-command> help`, read that action's Input and When to Use sections and return no more than 15 lines: purpose, usage, accepted arguments, and two examples. Never execute the command while serving help.

## Rules

- Keep this menu aligned with the core router.
- Name every sibling subcommand so a user who knows a command name can find its package, but keep usage detail in that package's own `help` — never restage a sibling's menu here. When a sibling gains or drops a command, update this list in the same commit.
- Help never mutates the queue, project, or Git state.
