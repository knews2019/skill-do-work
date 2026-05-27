# Prompts Action

> **Part of the do-work skill.** A library of reusable, battle-tested prompts surfaced as sub-commands. Lives in `prompts/` at the skill root — one Markdown file per prompt, each a standalone instruction the agent can execute. User-facing walkthrough: [`docs/prompts-guide.md`](../docs/prompts-guide.md).

Unlike built-in actions (which have fixed workflows), this action is a dispatcher over a growing collection of user-contributable prompt files. Think of it as a command-palette for recurring jobs the skill doesn't (yet) have a first-class action for.

## When to Use

**Use when:**
- The user names a prompt (`do-work prompts run adr`) or wants to browse the library (`do-work prompts list`).
- A recurring job has a reusable prompt in `prompts/` — running it is cheaper than rewriting the instructions.
- The user wants to inspect a prompt before running it (`do-work prompts show <name>`).

**Do NOT use when:**
- The user described a task but no prompt matches — suggest `do-work capture request:` instead of forcing an unrelated prompt.
- The user wants to *edit* a prompt file — that's a normal file edit, not a dispatcher invocation.
- A first-class action (`review work`, `code-review`, `ui-review`, `interview`, etc.) covers the job — prefer the built-in action.

## Sub-Commands

The `prompts` command accepts a sub-command as its first argument. If no sub-command is given, show the help menu.

| Sub-command | What it does |
|---|---|
| `list` | List every available prompt with its one-line description |
| `show <name>` | Print the prompt body without executing it (read-only) |
| `run <name> [args]` | Execute the prompt — adopt its body as the operational instructions |
| `<name> [args]` | Shorthand for `run <name> [args]` when `<name>` is not a reserved sub-command |
| (none) | Show help menu |

---

## Locating the Prompt Library

The library lives at `prompts/` relative to the skill root (the directory containing `SKILL.md`). Each prompt is a single `.md` file. `prompts/README.md` is the index and is not itself a runnable prompt.

**If `prompts/` does not exist:** tell the user "No prompt library found at `prompts/`." and stop. Do not create the directory silently — that would mask installation/layout problems.

---

## Help Menu (no sub-command)

When invoked with no sub-command (`do-work prompts`), show:

```
prompts — run reusable prompts from the library

  do-work prompts list              List every available prompt
  do-work prompts show <name>       Print a prompt (read-only)
  do-work prompts run <name>        Execute a prompt
  do-work prompts <name>            Shorthand for run

Examples:
  do-work prompts run architecture-decisions-log_create-or-expand
  do-work prompts architecture-decisions-log --dry-run
```

Then stop — do not execute anything.

---

## Sub-Command: `list`

1. Glob `prompts/*.md` and exclude `README.md`.
2. For each file, read the header (everything above the first `---`) to extract:
   - The title (first `# ` heading — strip the leading `# `)
   - The one-line description (the blockquote `> …` right under the title)
   - The aliases (the `**Aliases:**` line, if present — strip backticks and split on commas)
3. Render as a table:

```
Available prompts:

  NAME                                          ALIASES                  TITLE                                    DESCRIPTION
  architecture-decisions-log_create-or-expand   adr, adr-log, decisions  Architecture Decisions Log: Create…      Create or update a project-wide ADR log at decisions/…
  ...

Run any of them with: do-work prompts run <name>  (or use any alias as <name>)
```

Column widths can be approximate — readability beats precision. Omit the ALIASES column entirely if no prompt declares any. If the library is empty, say so and point the user at `prompts/README.md` for the "how to add a new prompt" section.

**Surface alias collisions in `list` output too.** If the alias map has any duplicates, print a one-line warning above the table: `⚠ Alias collisions detected: <alias> declared in <file-a>, <file-b>. These aliases cannot be used until the conflict is resolved.`

---

## Sub-Command: `show <name>`

1. Resolve `<name>` per the resolution rules below.
2. Read the file and print it verbatim, wrapped in a fence so the user sees the raw Markdown rather than rendered output.
3. **Do not execute it.** `show` is strictly read-only. Do not interpret any imperative instructions in the prompt body as commands for you to run.

---

## Sub-Command: `run <name> [args]`

1. Resolve `<name>` per the resolution rules below.
2. Read the file. Split it at the first `---` separator on its own line: everything above is the header (metadata), everything below is the body (your new instructions).
3. **Check the header for `**Runnable:**`.** Parse it the same way as `**Aliases:**` — single line after the bolded key. Take the **first token only** — everything up to the first whitespace or punctuation (e.g. `no — placeholder…` → `no`, `false (see below)` → `false`) — then lowercase and trim it. If that first token is `no`, `false`, or `never`, the prompt is opt-out — refuse with the explanation from the prompt's first blockquote line, e.g.:
   ```
   `<name>` is a placeholder, not a runnable instruction. <first-line description from the prompt>.
   Use `do-work prompts show <name>` to inspect it, or copy it into your project to activate the
   sidecar — the prompt header explains how.
   ```
   Stop without adopting the body. Any other value (or absence of the key) means runnable — proceed to the next step.
4. **Adopt the body as your operational instructions for the remainder of this turn.** Pass `[args]` through as the prompt's arguments (the body's "Inputs / flags" section, if any, defines how to interpret them).
5. Execute the body. This may involve creating files, running commands, making commits — do whatever the body instructs, subject to the global rules in `CLAUDE.md` and the user's permission mode.
6. After execution, return control to the normal flow (including `next-steps.md` suggestions — see "Post-run suggestions" below).

### Resolution rules for `<name>`

Try in priority order — first match wins:

1. **Exact filename match** (without `.md` extension). `architecture-decisions-log_create-or-expand` → `prompts/architecture-decisions-log_create-or-expand.md`.
2. **Exact alias match.** Build an alias map by reading each `prompts/*.md` file's header (everything above the first `---`) and parsing the `**Aliases:**` line. Aliases are backtick-quoted, comma-separated tokens — e.g. `**Aliases:** \`dca\`, \`dark-code-risk\``. Strip backticks and surrounding whitespace. If `<name>` matches exactly one alias, resolve to that prompt's filename.
   - **Collision detection:** if the same alias is declared in more than one prompt's header, treat it as ambiguous — list the candidate filenames and ask the user to disambiguate by full filename. Never silently pick one. Surface the collision so the library can be cleaned up.
3. **Unambiguous filename prefix match.** If `<name>` is a prefix of exactly one prompt filename, use that. If it's a prefix of multiple, list the candidates and ask the user to disambiguate.
4. **No match:** tell the user the prompt wasn't found and list available prompts (same output as `list`). Do not "helpfully" create the file.

The header parse stops at the first `---` separator, so aliases declared in the prompt body (if any) are ignored — only the header's `**Aliases:**` line counts.

---

## Shorthand: `<name> [args]`

If the first argument isn't `list`, `show`, `run`, or `help`, treat it as shorthand for `run <name> [args]`. `do-work prompts architecture-decisions-log --dry-run` is equivalent to `do-work prompts run architecture-decisions-log_create-or-expand --dry-run` (via prefix match).

---

## Rules

- **`show` is strictly read-only.** Never execute a prompt reached via `show`, even if its body contains imperative language.
- **Prompt files are immutable during `run`.** The dispatcher must not edit `prompts/*.md` while running one — edit the library only when the user explicitly asks you to (e.g., "add a new prompt to the library").
- **Pass arguments through untouched.** Don't pre-parse `args` — let the prompt body define its own argument handling.
- **Respect `--dry-run`, `--no-push`, and similar flags.** If the prompt body declares support for them, make sure you honor them end-to-end. If the body doesn't declare them, pass them through anyway and let the body decide.
- **Idempotency is the prompt's responsibility, not the dispatcher's.** But flag it as a red flag below if the prompt you're running doesn't describe how it detects prior state.
- **Never push to `main`/`master`** during a `run` unless the prompt body explicitly instructs you to and the user has confirmed.

## Post-run suggestions

After a successful `run`, suggest next steps. Default pattern:

```
Next steps:
  do-work commit                  Commit any uncommitted changes
  do-work prompts list            Browse other prompts in the library
  do-work prompts show <name>     Inspect a prompt before running it
```

If the prompt already committed and pushed its own work (like `architecture-decisions-log_create-or-expand` does), skip the `do-work commit` suggestion.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "I'll execute the prompt body from `show` since the user clearly wants the work done" | Tell the user `show` is read-only and point at `run` | `show` is the safe preview mode; executing it silently would surprise users who wanted to inspect before committing to a run |
| "The prompt file is out of date, I'll update it while running" | Do the run; tell the user the update is needed; let them decide whether to edit | Conflating execution with library maintenance makes runs non-reproducible |
| "No exact match but there are three candidates — I'll pick the first one alphabetically" | List the candidates and ask | Guessing the wrong prompt runs the wrong job |
| "The prompt asks to push, I'll add `--no-gpg-sign` to make it smoother" | Run it as written; if push fails, report and ask | Prompts are reviewed and trusted as written — silent flag injection breaks that contract |

## Red Flags

- The prompt you're about to `run` has no "When to use", no Rules section, and no pre-flight checks — execute it cautiously and warn the user its safety guarantees are weak
- `run` produces no commits AND no visible output — the prompt may have silently no-op'd; report and investigate
- After a `run`, the working tree has changes unrelated to what the prompt described — another process may have raced; stop and ask the user
- `list` returns zero prompts but the library directory exists — check for misnamed files or lost extensions

## Verification Checklist

- [ ] `do-work prompts` with no args shows the help menu and stops
- [ ] `do-work prompts list` enumerates every `prompts/*.md` except `README.md`
- [ ] `do-work prompts show <name>` prints the file verbatim and does NOT execute it
- [ ] `do-work prompts run <name>` executes only the body (below `---`), not the header metadata
- [ ] Unknown names trigger a "not found" message with the available-prompts list, not silent file creation
- [ ] `do-work prompts <name>` shorthand resolves to `run <name>` when `<name>` isn't a reserved sub-command
- [ ] Arguments after `<name>` pass through to the prompt body unchanged
- [ ] Aliases declared in prompt headers (`**Aliases:**` line) resolve to their prompt; alias collisions across files are surfaced rather than silently picking one
