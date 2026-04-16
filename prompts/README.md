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
| `adr-log` | Create or update a project-wide Architecture Decision Record log at `decisions/`, modeled on the BKB wiki pattern. Layered source mining (`implementation-history.md` → `lessons-learned/` → code, with `CHANGELOG.md` as fallback). Idempotent via REQ/UR keys. Resumable, supersession-aware. |
| `weekly-signal-diff` | Weekly structural diff of AI-industry news, personalized via BKB. Ships with a 10-lane core starter universe. At run time it searches the user's project for a `weekly-signal-diff-personal.md` (at project root, `.claude/`, `do-work/`, or anywhere via glob) and loads those lanes as full members of the scan. Produces an inline digest plus a durable deliverable ingested back into BKB so next week's run can diff against it. Every loaded lane gets full coverage every week. Idempotent per week-ending date. |
| `weekly-signal-diff-personal` | Placeholder template for the personal sidecar. Ships with no real lanes. Copy it anywhere in your project (project root, `.claude/`, `do-work/`, etc.) and fill in real lanes; the main prompt auto-discovers your project-local copy. Not run directly. |
