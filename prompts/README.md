# Prompt Library

Reusable, battle-tested prompts for recurring jobs — ADR logs, retrospectives, audit passes, and so on. Each prompt is a standalone Markdown file an agent can execute directly.

**How to use:**

```
do-work prompts                    # short help menu
do-work prompts list               # list every available prompt
do-work prompts show <name>        # print the prompt body (read-only)
do-work prompts run <name> [args]  # execute the prompt
do-work prompts <name> [args]      # shorthand for run
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
| `architecture-decisions-log_create-or-expand` | Create or update a project-wide Architecture Decision Record log at `decisions/`, modeled on the BKB wiki pattern. Layered source mining (`implementation-history.md` → `lessons-learned/` → code, with `CHANGELOG.md` as fallback). Idempotent via REQ/UR keys. Resumable, supersession-aware. Aliases: `adr`, `adr-log`, `decisions`. |
| `weekly-structural-diff-original` | Filter signal from noise in AI news. Sort items into signal/noise, diagnose shifts across five altitudes (physics, monetization, geography, business models, geopolitics), and produce prioritized takeaways with a "what didn't change" calibration. |
| `economics-inference-stress-test` | Run any AI product through a Sora-style economics stress test — sustainability ratio, three scenarios, emoji verdict (🟢/🟡/🟠/🔴), and a concrete "what would fix it" plan. |
| `tech-infrastructure-compute-geography-risk` | Map physical-layer risks (power/grid, permitting/politics, geopolitics, data residency) across AI compute locations. Produces a risk matrix, deployment strategy, and per-location contingency playbook. |
| `economics-saas-repricing-exposure` | Estimate seat compression, compute "The Clock" until it shows in reported numbers, assess pricing-model transition readiness, and benchmark against Atlassian. |
| `business-vendor-strategic-sort` | Evaluate 2–5 AI vendors across five structural-sustainability dimensions, attach a tripwire event to each, and score portfolio-level concentration risk. |
| `tech-inference-architecture-decision` | Design an inference architecture with economics as a first-class constraint — API vs. self-hosted vs. hybrid, model selection, the Sora test, and a Now / 3× / 10× migration path. |
| `weekly-signal-diff` | Weekly structural diff of AI-industry news, personalized via BKB. Ships with a 10-lane core starter universe. At run time it searches the user's project for a `weekly-signal-diff-personal.md` (at project root, `.claude/`, `do-work/`, or anywhere via glob) and loads those lanes as full members of the scan. Produces an inline digest plus a durable deliverable ingested back into BKB so next week's run can diff against it. Every loaded lane gets full coverage every week. Idempotent per week-ending date. |
| `weekly-signal-diff-personal` | Placeholder template for the personal sidecar. Ships with no real lanes. Copy it anywhere in your project (project root, `.claude/`, `do-work/`, etc.) and fill in real lanes; the main prompt auto-discovers your project-local copy. Not run directly. |
| `prompt-kit-step0-pen-and-paper-exercises-to-prepare-prompt` | Pre-flight exercise the user does **away from the screen**. Seven questions, pen and paper, 10 minutes. The agent hands off and gets out of the way, then structures the answers into a PRE-FLIGHT BRIEF when the user returns. |
| `prompt-kit-step1-four-discipline-diagnostic` | Thorough audit across Prompt Craft, Context Engineering, Intent Engineering, and Specification Engineering. Produces a scored table, 10x gap analysis, and a personalized 4-month roadmap. |
| `prompt-kit-step2-personal-context-doc` | Structured seven-domain interview that produces a copy-paste-ready personal context document — the user's "CLAUDE.md for everything." Role, audiences, quality standards, institutional knowledge, constraints, AI patterns. |
| `prompt-kit-step3-spec-engineer` | Collaboratively build a full specification document for a real project — acceptance criteria, constraint architecture, task decomposition, evaluation criteria, definition of done. Spec an autonomous agent can execute against. |
| `prompt-kit-step4-intent-and-delegation-framework` | Extract the implicit decision-making rules the user's team operates by. Encodes priority hierarchy, decision authority map, quality thresholds, failure modes, and a Klarna Test self-check. Works for teams or as a personal framework. |
| `prompt-kit-step5-eval-harness` | Lütke-pattern eval suite over the user's actual recurring tasks. 3 test cases with refined inputs, observable quality criteria, known failure modes, and a scoring rubric fast enough to run after every model release. |
| `prompt-kit-step6-constraint-architecture` | Pre-delegation exercise that produces a four-quadrant constraint document (Must Do / Must Not / Prefer / Escalate) tied to the user's specific failure modes for a given task. Stops the "smart-but-wrong" pattern before it happens. |
