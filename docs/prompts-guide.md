# Prompts — Reusable Prompt Library

A dispatcher over `prompts/` — a growing library of standalone, runnable Markdown prompts for recurring jobs the skill doesn't have a first-class action for (ADR logs, retrospectives, audit passes, analytical frameworks, etc.).

Think of it as a command-palette for prompts you'd otherwise paste into a new session every time.

Aliases: `prompts`, `prompt`.

## Sub-commands

```
do-work prompts                    # help menu
do-work prompts list               # every available prompt + one-line descriptions
do-work prompts show <name>        # print the prompt body (read-only — does NOT execute)
do-work prompts run <name> [args]  # execute the prompt body as instructions
do-work prompts <name> [args]      # shorthand for `run <name> [args]`
```

## Name resolution

- **Exact match** on the filename (minus `.md`) wins.
- **Unambiguous prefix match** is accepted. `adr` resolves to `architecture-decisions-log_create-or-expand` if that's the only prefix match.
- **Ambiguous prefix**: candidates are listed, user disambiguates. The dispatcher never guesses.
- **No match**: not-found message + the `list` output. No silent file creation.

## How a prompt is shaped

```markdown
# <Prompt Name>

> <One-line description>

**Aliases:** <optional>
**When to use:** <2-3 bullets>
**Inputs / flags:** <optional arguments the prompt accepts>

---

<prompt body — the actual instructions the agent executes>
```

The dispatcher reads the header for `list`/`show` output. On `run`, it adopts everything below the `---` separator as its operational instructions for that turn.

## Safety model

- `show` is strictly read-only — the dispatcher refuses to execute a prompt reached via `show`, even if the body contains imperative language.
- Prompt files are immutable during `run` — the dispatcher must not edit `prompts/*.md` while running one.
- Arguments pass through untouched. The dispatcher does not pre-parse; the prompt body defines its own argument handling.
- `--dry-run`, `--no-push`, and similar flags are passed through. If the prompt body declares support, it's honored end-to-end.
- Prompts never push to `main`/`master` unless the body explicitly instructs to, and the user has confirmed.

## Adding a new prompt

1. Create `prompts/<kebab-name>.md` with the header + `---` + body shown above.
2. Keep it **idempotent** — re-running should detect existing state, not duplicate work.
3. Make it **resumable** — if execution can reasonably take multiple sessions, persist progress in a dedicated file the prompt reads on re-entry.
4. Add one line under **Available prompts** in `prompts/README.md`.

## Common patterns

- **Mine → build → log:** prompts that read prior artifacts, produce new ones, and append a timestamped entry (weekly diffs, ADR logs).
- **Interview → export:** prompts that elicit structured input from the user and emit a durable document (context docs, constraint architectures).
- **Scan → audit → rank:** prompts that read the codebase or external inputs, apply a rubric, and produce a ranked findings list.

See `prompts/README.md` for the authoritative list of prompts shipped with the skill.

## Troubleshooting

- **"No prompt library found at `prompts/`."** — Installation issue. Re-extract the skill bundle; do not create the directory manually.
- **`list` returns zero prompts** — Check for misnamed files (missing `.md`, leading underscore) in `prompts/`.
- **`run` produced no output and no commits** — The prompt may have silently no-op'd. Use `show` to inspect the body and verify it actually does work when the inputs/flags you passed are present.
- **Working tree has unexpected changes after `run`** — Stop. Another process may have raced, or the prompt touched files its description didn't mention. Investigate before committing.

## See also

- `actions/prompts.md` — the dispatcher action file (what the skill actually runs).
- `prompts/README.md` — authoritative prompt index and template.
