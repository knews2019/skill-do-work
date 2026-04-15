# Prompt Library

Reusable, battle-tested prompts for recurring jobs — ADR logs, retrospectives, audit passes, and so on. Each prompt is a standalone Markdown file an agent can execute directly.

**How to use:**

```
do work prompts                    # short help menu
do work prompts list               # list every available prompt
do work prompts show <name>        # print the prompt body (read-only)
do work prompts run <name> [args]  # execute the prompt
do work prompts <name> [args]      # shorthand for run
```

Resolution rules: `<name>` matches the filename without the `.md` extension. Exact match wins; otherwise a single unambiguous prefix match is accepted.

**How a prompt file is shaped:**

```markdown
# <Prompt Name>

> <One-line description>

**Aliases:** <optional>
**When to use:** <2-3 bullets>
**Inputs / flags:** <optional arguments the prompt accepts>

---

<prompt body — the actual instructions the agent executes>
```

The dispatcher (`actions/prompts.md`) reads the header for `list`/`show` output and adopts the body below the `---` separator when `run` is invoked.

**How to add a new prompt:**

1. Create `prompts/<kebab-name>.md` with the header + `---` + body.
2. Keep prompts **idempotent** — re-running should detect existing state, not duplicate work.
3. Make prompts **resumable** — if execution can reasonably take multiple sessions, persist progress in a dedicated file the prompt reads on re-entry.
4. Add one line under **Available prompts** below.

**Available prompts:**

| Name | What it does |
|---|---|
| `adr-log` | Create or update a project-wide Architecture Decision Record log at `decisions/`, modeled on the BKB wiki pattern. Mines `CHANGELOG.md` for load-bearing decisions. Resumable, idempotent, handles supersession. |
