# Moved Command Shim

> One-release compatibility action for commands that moved out of core at the modular cutover. It prints the exact replacement invocation and stops; it never loads, forwards to, or reimplements a sibling action.

## Input

The complete legacy command and its arguments, with the leading `do-work` removed.

## Steps

1. Match the first legacy trigger below. Preserve the remaining arguments verbatim, except for the two renamed install targets.
2. Print exactly one line in this form and stop:

   ```text
   This command moved. Run: <replacement>
   ```

   | Legacy trigger | Replacement |
   |---|---|
   | `board`, `kanban`, `queue board`, `show the board` | `do-work-board board <arguments>` |
   | `bkb`, `build knowledge base`, `knowledge base`, `kb` | `do-work-knowledge bkb <arguments>` |
   | `memory <arguments>` | `do-work-knowledge memory <arguments>` |
   | `remember <text>`, `forget <text>`, `recall <query>`, `what do you remember <query>` | `do-work-knowledge memory <canonical subcommand and arguments>` |
   | `dream`, `consolidate memory`, `clean up wiki`, `memory cleanup` | `do-work-knowledge dream <arguments>` |
   | `interview`, `elicit`, `operating model` | `do-work-knowledge interview <arguments>` |
   | `prompts`, `prompt` | `do-work-knowledge prompts <arguments>` |
   | `install memory-module`, `install memory module`, `install memory hooks` | `do-work-knowledge setup-memory` |
   | Any legacy alias for `validate-feedback`, `code-review`, `ui-review`, `present-work`, `ai-report`, `slop-check`, `quick-wins`, `scan-ideas`, `deep-explore`, `prime`, `inspect`, `note`, `stray-check`, `tidy-repo`, or `tutorial` | `do-work-toolbox <canonical trigger> <arguments>` |
   | `install ui-design`, `install bowser`, `install last30days`, `install ideation-adhd`, and their legacy aliases | `do-work-toolbox install <canonical target>` |
   | `install just-kanban`, `install justfile`, `install run-kanban`, `install run-do-work-update` | `just run-kanban` (the suite installer already manages the Just recipes) |

3. If the trigger is not in the table, print `do-work help` guidance instead of guessing.

## Rules

- Do not dispatch another skill or read its action file.
- Do not mutate the queue, project, settings, or Git state.
- Omit the trailing space and `<arguments>` when no arguments remain.
