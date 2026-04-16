# Expand skill-do-work: add the `interview` action

You are working inside a local checkout of the `knews2019/skill-do-work` repository. Your task is to add a new action called `interview` — a generalized interview framework with prescriptive templates. The first template is `work-operating-model`, based on the Work Operating Model concept published by Nate B. Jones and Jonathan Edwards.

Produce a single coherent change that could be reviewed as one pull request.

---

## Step 0 — Orient yourself in the repo

Before writing anything, read these files in order and internalize the conventions they establish. Do not skim. Do not skip.

1. `CLAUDE.md` — master conventions for this repo.
2. `SKILL.md` — root action registry; you will register `interview` here.
3. `README.md` — user-facing surface; you will cross-reference it.
4. `AGENTS.md` — any agent-level guardrails.
5. `decisions/records/adr-001-modular-action-prompts-and-companion-references.md` — actions may have a companion `-reference.md` file for heavy content. Apply this: `interview` gets both `actions/interview.md` and `actions/interview-reference.md`.
6. `decisions/records/adr-002-load-reusable-spec-templates-during-work.md` — actions load reusable templates at runtime. The `interview` action loads templates from a new top-level `interviews/` directory.
7. `decisions/records/adr-005-pipeline-is-stateful-and-resumable.md` — stateful, resumable actions persist session state to disk. The `interview` action follows this pattern with `./interview/<template>/session.json`.
8. `decisions/records/adr-009-build-knowledge-base-as-a-compiled-interlinked-wiki.md` and `adr-010` — these define how `bkb` ingests sources. The `interview ingest` sub-command produces markdown that fits `bkb`'s expected inbox format.
9. `actions/build-knowledge-base.md` — read the whole thing. Your `actions/interview.md` should match its detail level, sub-command dispatch pattern, and prose register.
10. `crew-members/general.md` and one or two others — confirm the crew persona format.
11. Existing ADRs under `decisions/records/` — note the numbering and frontmatter shape; your new ADR is `adr-011`.
12. `decisions/_master_index.md`, `decisions/_progress.md`, and the `decisions/topics/_index_*.md` files — you will update the master index and the relevant topic index.

**Non-negotiable constraints:**

- Local files only. No MCP dependencies. No references to Open Brain, OB1, Supabase, Anthropic services beyond whatever client is running the action, or any other external system.
- Prescriptive template shape: templates define layers + canonical entry contract + export schemas + checkpoint format. The action enforces contracts.
- Single instance per template per repo. Context separation comes from installing the skill in multiple repos, not from profiles or workspaces.
- Fidelity to Nate's prescriptions for the `work-operating-model` template (the five layers, the canonical entry contract, the layer-specific detail shapes, the export filenames, the tone). Improvements are allowed only where this prompt explicitly names them.

---

## Step 1 — Create new files

Create these files with the exact content specified. Where content is given verbatim in fenced blocks below, use it verbatim. Where content is described as a specification, produce it in the style of the existing repo files you read in Step 0.

### 1.1 `actions/interview.md`

This is the main action prompt. Model it on `actions/build-knowledge-base.md` — same header style, same sub-command table format, same level of specification. Contents must include:

**Header section.**

A short description: "Generalized interview framework. Runs prescriptive templates that elicit tacit knowledge through structured multi-layer conversations and produce agent-ready operating artifacts."

Note that this action is part of the do-work skill, references ADR-001 (modular action + companion reference), ADR-002 (reusable template loading), and ADR-005 (stateful and resumable).

**Sub-commands table.** Use this exact set:

| Sub-command | What it does | Crew |
|---|---|---|
| `list` | List available templates in `interviews/` | Architect |
| `<template>` | Start or continue the interview for the named template | Interviewer |
| `<template> status` | Show session progress — layers approved, layer pending, last activity | Interviewer |
| `<template> review` | Run the contradiction pass across all approved layers | Interviewer + Reviewer |
| `<template> export` | Generate the template's declared export artifacts | Interviewer + Editor |
| `<template> ingest` | Copy exports into `kb/raw/inbox/` with frontmatter suitable for `bkb triage && bkb ingest` | Librarian |
| `<template> reset` | Archive the current run as a version and start fresh (requires confirmation) | Architect |
| `<template> versions` | List archived runs for the template | Architect |
| (none) | Show help menu | (any) |

**Locating the template.** When a sub-command references `<template>`, resolve as follows: the action reads `interviews/<template>.md` from the repo root. If the file does not exist, list available templates and stop.

**Locating the session.** Session state lives at `./interview/<template>/session.json` in the current working directory. If it does not exist for sub-commands other than the bare `<template>` invocation, tell the user to run `do work interview <template>` to start a session.

**Session lifecycle — the `<template>` sub-command.**

1. If `./interview/<template>/session.json` does not exist, create the directory structure and start a fresh session:
   ```
   ./interview/<template>/
   ├── session.json
   ├── checkpoints/
   ├── exports/
   ├── versions/
   └── CHANGELOG.md
   ```
2. If `session.json` exists and has `status: "in_progress"`, resume from the pending layer without asking.
3. If `session.json` exists and has `status: "complete"`, prompt the user to choose a re-run mode:
   - **fresh** — archive current state as `versions/v<N>-YYYY-MM-DD/` (full copy of session.json + checkpoints/ + exports/), then start a new session with empty state. Append a CHANGELOG entry.
   - **update** — load prior layers; for each layer, show the stored canonical entries and ask "is this still accurate? confirm / edit / add / remove entries." In-place updates only; does not create a new version folder. Append a CHANGELOG entry summarizing diffs.
   - **version** — archive current state as `versions/v<N>-YYYY-MM-DD/`, then start a new session seeded with a reference to `v<N>` in `session.json` (`previous_version: "v<N>"`). The new session begins empty; the user re-interviews from scratch but can query the prior version for comparison. Append a CHANGELOG entry.

**Layer interview workflow.** For each layer in the template's declared order:

1. Begin with concrete, recent-example questions drawn from the template's `prompts` field for that layer. Never open a layer with an abstract question.
2. Convert the user's responses into canonical entries matching the template's entry contract.
3. Write a human-readable checkpoint file to `./interview/<template>/checkpoints/<layer-id>.md` with:
   - Short layer summary
   - 2–5 canonical entries
   - Explicit unresolved items
   - The exact approval ask
4. Present the checkpoint to the user in-chat and wait for explicit confirmation. Accepted confirmations: "save," "looks right," "confirmed," "approve," or equivalent. Corrections: the user edits specific entries; regenerate the checkpoint and ask again.
5. Only after confirmation: write the approved entries into `session.json` under `layers.<layer_id>.entries[]`, mark the layer `approved: true`, update `last_validated_at`, and advance `pending_layer` to the next.
6. Append one summary line to `./interview/<template>/CHANGELOG.md`:
   ```
   ## YYYY-MM-DD HH:MM — layer approved: <layer-id>
   <one-sentence summary of the real pattern surfaced>
   ```
7. Never persist content the user did not confirm or explicitly approve as a synthesized pattern.

**`<template> review` sub-command.** After all layers are approved, the review pass walks the session and surfaces contradictions:

- rhythms vs dependencies (e.g., claimed deep-work window collides with standing meeting in dependencies)
- recurring decisions vs institutional knowledge (e.g., decision rule cites data the user also said isn't written down anywhere)
- friction vs operating rhythm (e.g., a recurring friction implies a broken rhythm that wasn't named)

Present each tension explicitly, naming the two entries that conflict. For each, the user picks: revise-A, revise-B, both-are-true-note-the-tension, or skip. Revisions rewrite the affected layer's entries and require a fresh checkpoint approval before re-saving.

**`<template> export` sub-command.** Only runs when all five layers are approved and the review pass has been completed at least once. Reads the template's `exports` declarations and writes each file to `./interview/<template>/exports/`. Appends one final synthesis line to `CHANGELOG.md`.

**`<template> ingest` sub-command.** Reads from `./interview/<template>/exports/`. For each export file plus one markdown-per-layer summary, writes a file to `<repo-root>/kb/raw/inbox/` with YAML frontmatter matching `bkb`'s expected format:

```yaml
---
title: Work Operating Model — Operating Rhythms
source: ./interview/work-operating-model/exports/USER.md
type: source-summary
topic_cluster: operating-model
confidence: high
created: YYYY-MM-DD
---
```

If `kb/` does not exist, tell the user to run `do work bkb init` first and stop.

**Error handling.**

- Template file missing → list available templates, suggest `do work interview list`.
- Session corrupt (invalid JSON) → do not attempt repair. Tell the user where the file is and stop.
- `export` invoked with unapproved layers → list which layers are missing approval.
- `ingest` invoked without completed exports → tell the user to run `export` first.
- `reset` without confirmation → require explicit `reset --confirm` or an interactive "yes" before archiving.

**Help menu.** Standard format matching `actions/build-knowledge-base.md`'s help menu, with the sub-commands listed above and a short typical-flow block at the bottom.

---

### 1.2 `actions/interview-reference.md`

Companion reference per ADR-001. This file holds heavy content that would bloat the action file. It must include, in separate clearly-headed sections:

**Template file format.** The specification for what a template author must declare. Embed this as a concrete example showing the YAML-front-matter-plus-markdown-body format templates use.

**Canonical entry contract.** Every saved entry must include these fields, regardless of template:

- `title` — short name for the entry
- `summary` — 1–2 sentence description
- `cadence` — how often this pattern applies (e.g., "daily," "weekly," "per-project")
- `trigger` — what event or condition activates this pattern
- `inputs` — list of strings: what this pattern needs as input
- `stakeholders` — list of strings: who is involved
- `constraints` — list of strings: limitations or guardrails
- `details` — layer-specific object, shape defined per-layer by the template
- `source_confidence` — `confirmed` (user stated or approved as written) or `synthesized` (abstracted from multiple examples, user approved synthesis)
- `status` — `active`, `stale`, or `aspirational`
- `last_validated_at` — ISO timestamp

**Session.json schema.** Specify the full shape:

```json
{
  "template": "work-operating-model",
  "session_id": "<uuid>",
  "started_at": "<iso>",
  "last_activity_at": "<iso>",
  "status": "in_progress | complete",
  "pending_layer": "<layer-id> | null",
  "previous_version": "<version-id> | null",
  "layers": {
    "<layer-id>": {
      "approved": true,
      "approved_at": "<iso>",
      "entries": [ /* canonical entries */ ]
    }
  }
}
```

**Checkpoint file format.** Show the template structure:

```markdown
# Checkpoint: <layer title>

## Summary
<1–2 paragraph layer summary>

## Entries
### <entry title>
- cadence: <value>
- trigger: <value>
- inputs: <list>
- stakeholders: <list>
- constraints: <list>
- details:
  <layer-specific fields>
- source_confidence: confirmed | synthesized

## Unresolved
<bullet list of unresolved items or "None">

## Approval
This is the <layer-id> model I'd save right now. Correct anything that's off. If this looks right, I'll save it and move to <next-layer-id>.
```

**Export schemas.** For each of the five work-operating-model exports, provide the expected structure. Reference these from the template file; do not duplicate.

- `USER.md` — narrative profile of the person at work (name if known, role if known, operating rhythm summary, key recurring decisions, top dependencies, institutional knowledge they carry, active friction points). Written in third person, present tense.
- `SOUL.md` — decision framework an agent would follow when acting on behalf of this person: when to escalate, when to decide autonomously, tone rules for different audiences, data sources to trust, "good enough" thresholds for each task type.
- `HEARTBEAT.md` — checklist an agent reviews on a cadence (default 30 min) to decide if there's work to do for this person: what to check, what signals to act on, what to ignore.
- `operating-model.json` — full machine-readable dump of the approved session's canonical entries, grouped by layer.
- `schedule-recommendations.json` — derived from operating_rhythms + dependencies: suggested time blocks, standing slots, avoid-windows.

**Re-run mode specifications.** Detailed behavior for fresh / update / version, including what gets archived, what gets written, and what the CHANGELOG entry looks like for each.

**Versioning scheme.** Versions are directory-named `v<N>-YYYY-MM-DD` where `<N>` is monotonically increasing per template. A session can reference a prior version via `previous_version`. Versions are immutable; the action never edits files inside `versions/`.

---

### 1.3 `docs/interview-guide.md`

User-facing guide. Match the style of `docs/bkb-guide.md`. Cover:

- What the action is and when to run it ("sit down for 45 focused minutes when you need to hand work off, onboard an agent, or make your judgment patterns explicit").
- Expected time: ~45 minutes for work-operating-model, possibly more on first run.
- Output files: list the five exports and what each is for in one line each.
- Re-run cadence: quarterly, or after a major role change.
- Integration with bkb: after export, run `do work interview <template> ingest` to turn the operating model into queryable knowledge.
- How to install the skill in multiple repos for context separation (one operating model per repo).
- Troubleshooting: what to do if a checkpoint feels wrong (answer: edit in chat before confirming; never approve something you'll regret), what to do if session.json is corrupt (don't repair, contact maintainer or start fresh), what to do if the agent starts asking abstract questions (push back: "ask about last week").

Keep it under 5 KB. The reference lives in `actions/interview-reference.md`; this doc is the onboarding read.

---

### 1.4 `crew-members/interviewer.md`

New crew persona. Match the format of existing crew members (read `crew-members/debugging.md` and `crew-members/security.md` for structure). Contents:

**Role.** You are the Interviewer. You run structured elicitation interviews that turn tacit work knowledge into explicit, delegatable structure. You are not a consultant, not a coach, not a therapist. You are an interviewer whose job is to ask the right questions in the right order and capture the answers faithfully.

**Focus.**
- Concrete before abstract. Always anchor on last week, last month, recent examples, specific misses.
- Checkpoint ritual. Every layer ends with a summary and an explicit approval ask. No saves without confirmation.
- Honest confidence. Tag every entry `confirmed` or `synthesized`. Never inflate certainty.
- Surface contradictions. Don't smooth over tensions between what the user said in different layers.
- Momentum without bulldozing. Keep the interview moving but never push past a "wait, let me rethink that."

**Standards.**
- Do not open a layer with "what do you do all day?" or any equivalent abstraction. Start with "walk me through a real Monday" or "tell me about last week's blockers."
- Do not batch questions. One question at a time; wait for the answer.
- Do not infer fields the user did not provide. If the canonical entry contract asks for `constraints` and the user didn't mention any, ask. Do not invent.
- Do not paraphrase aggressively. When the user says something specific, capture their language, not a smoothed-over version.
- Do not produce generic productivity advice at any point during the interview. You are not here to coach.

**When active.** `interview` action and all its sub-commands. Adopts tone and focus for the duration of the interview.

---

### 1.5 `interviews/work-operating-model.md`

The first template. Prescriptive shape: declares layers, prompts, entry contract extensions, and export definitions. Structure as YAML frontmatter plus markdown body.

```markdown
---
name: work-operating-model
description: |
  Elicits the operating model of a person at work. Produces agent-ready
  artifacts (USER.md, SOUL.md, HEARTBEAT.md) plus machine-readable exports.
  Based on the five-layer Work Operating Model by Nate B. Jones and
  Jonathan Edwards.
version: 1.0.0
layers:
  - id: operating_rhythms
    title: Operating Rhythms
    order: 1
  - id: recurring_decisions
    title: Recurring Decisions
    order: 2
  - id: dependencies
    title: Dependencies
    order: 3
  - id: institutional_knowledge
    title: Institutional Knowledge
    order: 4
  - id: friction
    title: Friction
    order: 5
exports:
  - path: USER.md
    kind: narrative
  - path: SOUL.md
    kind: decision-framework
  - path: HEARTBEAT.md
    kind: checklist
  - path: operating-model.json
    kind: machine-readable
  - path: schedule-recommendations.json
    kind: machine-readable
---

# Work Operating Model Template

The first job is not to automate the user. It is to help them see and describe how their work actually runs.

## Layer 1: Operating Rhythms

Map how the user's days, weeks, and months actually unfold — not the calendar version, the real one.

### Prompt patterns
- "Walk me through a real Monday from the last two weeks."
- "Where does your calendar lie to you?"
- "When are you actually good for deep work versus admin or reactive work?"
- "What repeats weekly or monthly even when it isn't formally scheduled?"

### Details shape
Every entry's `details` field must include:
- `time_windows` — list of `{start, end, label}` objects describing recurring time blocks
- `energy_pattern` — string describing when the user has energy for which kind of work
- `interruptions` — list of recurring interruptions and their sources
- `non_calendar_reality` — string describing what actually happens that isn't on the calendar

## Layer 2: Recurring Decisions

Capture the judgment calls the user makes over and over, especially ones where the answer depends on context rather than a checklist.

### Prompt patterns
- "What decisions do you make over and over where the answer depends on context, not a checklist?"
- "What do you look at before you decide?"
- "When do you escalate versus handle it yourself?"
- "Which decisions are reversible if you get them wrong?"

### Details shape
- `decision_name` — short name for the decision
- `decision_inputs` — list of data sources or signals checked
- `thresholds` — list of `{metric, value, direction}` — the numbers that matter
- `escalation_rule` — when the user passes this up or brings someone else in
- `reversible` — boolean: can this decision be undone cheaply

## Layer 3: Dependencies

Map who and what the user waits on, and what breaks when those inputs are late or wrong.

### Prompt patterns
- "What part of your week depends on someone else sending, approving, or clarifying something?"
- "What breaks when that doesn't happen on time?"
- "What's your fallback when you're blocked?"

### Details shape
- `dependency_owner` — person or system the user waits on
- `deliverable` — what they send / approve / provide
- `needed_by` — timing window
- `failure_impact` — what breaks if it's late or wrong
- `fallback` — what the user does when blocked

## Layer 4: Institutional Knowledge

Surface what the user knows that isn't written down anywhere — the context only they carry.

### Prompt patterns
- "What do you know that your team relies on but nobody has really documented?"
- "What mistakes would a smart new hire make because the real context is still in your head?"
- "What would break if you disappeared for two weeks?"

### Details shape
- `knowledge_area` — short name for the domain
- `why_it_matters` — why this context is load-bearing
- `where_it_lives` — "in my head," or a specific partial source
- `who_else_knows` — list of people who partially share this
- `risk_if_missing` — what goes wrong without this knowledge

## Layer 5: Friction

Name the recurring annoyances that eat time — the tooling gaps, the duplicate work, the waits.

### Prompt patterns
- "What keeps eating 10-20 minutes at a time?"
- "Where do you keep doing work the hard way because the systems never quite line up?"
- "What's the same broken handoff you've been complaining about for months?"

### Details shape
- `frequency` — how often this friction hits
- `time_cost` — rough minutes or hours lost per occurrence
- `current_workaround` — what the user does today
- `systems_involved` — tools, services, or people in the friction loop
- `automation_candidate` — boolean: could this reasonably be automated
- `priority` — `low`, `medium`, or `high` when the user is willing to rank

## Cross-layer contradiction checks

During the `review` sub-command, surface these specific tensions:

- **Rhythm vs Dependencies** — A claimed deep-work window that collides with a standing dependency handoff.
- **Decisions vs Knowledge** — A decision rule that cites data the user also said isn't written down anywhere.
- **Friction vs Rhythm** — A recurring friction pattern that implies the stated rhythm isn't real.
- **Dependencies vs Knowledge** — A dependency owner who is the same person the user said carries undocumented context (single point of failure).

## Tone

Direct, practical, specific. No generic productivity advice. No fake certainty. Keep momentum moving without bulldozing confirmation.
```

---

### 1.6 `decisions/records/adr-011-interview-framework-with-prescriptive-templates.md`

New ADR. Follow the format of existing ADRs (read `adr-009` and `adr-010` for structure). Contents:

**Status.** Accepted.

**Context.** The do-work skill includes actions for knowledge compilation (`bkb`), code review, commits, and pipeline orchestration, but has no action for extracting the operator's own tacit knowledge into a form agents can act on. Without this, users who want agents to take work off their plate must write their own `USER.md` / `SOUL.md` / `HEARTBEAT.md` by hand, which is the step most people fail at. The Work Operating Model by Nate B. Jones and Jonathan Edwards defines a five-layer elicitation interview that produces these files; it was published with an Open Brain / MCP dependency we don't want to take.

**Decision.** Add an `interview` action implementing a generalized elicitation framework with prescriptive templates. Templates live in `interviews/<name>.md` and declare layers, per-layer prompts, canonical entry contract, and export schemas. The action enforces contracts, handles session state, runs checkpoint approval gates, and produces exports. The first template is `work-operating-model`, using the five layers and canonical entry contract from the source prescription verbatim. Persistence is local-file only under `./interview/<template>/`. Templates are prescriptive (not minimal, not fully executable) to balance consistency of output against authoring cost for future templates.

**Consequences.**

- Users can produce agent-ready operating artifacts by running a ~45-minute interview.
- Exports can flow directly into `bkb` via a new sub-command, making the operating model queryable alongside other knowledge.
- Adding new templates (post-mortem, new-hire-onboarding, project-kickoff) requires writing a template file against the prescriptive schema; no changes to the action.
- Session state is per-CWD, which means multi-context users install the skill in multiple repos. No profile concept is introduced.
- Re-runs support `fresh`, `update`, and `version` modes; versions are immutable.
- No external service dependencies.

**Alternatives considered.**

- *Minimal template shape* — too permissive, hard to guarantee consistent output across templates.
- *Executable template shape* — too much per-template work, risks every template diverging into its own workflow.
- *Profiles within the same repo* — rejected in favor of multi-repo installation for simpler paths and cleaner mental model.
- *Direct Open Brain / MCP integration* — rejected; adds infrastructure the skill doesn't otherwise require and locks users into a specific backend.

---

## Step 2 — Edit existing files

### 2.1 `SKILL.md`

Register the new action. Locate the action registry section (it will enumerate `work`, `bkb`, `capture`, `commit`, etc.) and add an `interview` entry in alphabetical position. Style and prose register must match the surrounding entries.

### 2.2 `README.md`

In the actions list section, add an entry for `interview`. One-sentence description: "Run a structured elicitation interview and generate agent-ready operating artifacts." Link to `docs/interview-guide.md`.

### 2.3 `CHANGELOG.md`

Add a new entry at the top following the existing format. Title: "Add `interview` action — generalized elicitation framework with prescriptive templates, first template `work-operating-model`." Body: one short paragraph summarizing the scope (action, companion reference, guide, crew member, first template, ADR-011). Do not embed the full prompt content; reference the ADR.

### 2.4 `decisions/_master_index.md`

Add ADR-011 to the index following the existing format.

### 2.5 `decisions/_progress.md`

If this file tracks ADR status, add ADR-011 with status Accepted.

### 2.6 `decisions/topics/_index_skill-architecture.md` and `decisions/topics/_index_workflow-orchestration.md`

ADR-011 belongs under at least one of these topics (probably both — it's architectural and it introduces a new workflow). Add entries to both if the topic file structure supports it. If only one is appropriate based on the existing topic definitions, pick that one.

### 2.7 `.gitignore`

Add a line that excludes session artifacts from being accidentally committed:
```
# interview action — session state per repo install
interview/
```

(The `interview/` output directory is per-repo session state; it should not be versioned alongside the code the repo is actually about.)

---

## Step 3 — Verify

Before declaring done, verify all of the following:

- [ ] All six new files exist with the specified content.
- [ ] Seven existing files have been edited (SKILL.md, README.md, CHANGELOG.md, decisions/_master_index.md, decisions/_progress.md, at least one topic index, .gitignore).
- [ ] No file in the change references OB1, Open Brain, Supabase, MCP, or any external service.
- [ ] The `work-operating-model` template's five layers are in the exact order: operating_rhythms, recurring_decisions, dependencies, institutional_knowledge, friction.
- [ ] The canonical entry contract in `actions/interview-reference.md` includes all 11 required fields: title, summary, cadence, trigger, inputs, stakeholders, constraints, details, source_confidence, status, last_validated_at.
- [ ] Each of the five layers in the template has its specified `details` shape.
- [ ] The action supports all nine sub-commands listed in the table.
- [ ] The re-run modes fresh / update / version are all specified with clear archival and CHANGELOG behavior.
- [ ] ADR-011 is numbered correctly relative to existing ADRs.
- [ ] The action's help menu matches the style of `actions/build-knowledge-base.md`'s help menu.
- [ ] Prose across new files matches the register of existing repo files (read `actions/build-knowledge-base.md` and `docs/bkb-guide.md` if unsure).

When all checks pass, report what was created and edited, and stop. Do not commit — the user will review the diff first.
