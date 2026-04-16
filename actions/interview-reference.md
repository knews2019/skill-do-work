# Interview Action — Reference

> **Companion to `actions/interview.md`.** Heavy content the action file references at runtime: template authoring format, canonical entry contract, `session.json` schema, checkpoint file format, export schemas for the `work-operating-model` template, re-run mode specifications, and versioning scheme. Applies per ADR-001.

Read this file when authoring a new template, implementing the interview action, or debugging session state. The action file stays short; this file holds the specifications.

---

## Template File Format

A template declares the layers, per-layer prompts, canonical entry contract extensions, and export schemas. Templates live at `<skill-root>/interviews/<template-name>.md` (the `interviews/` directory inside the skill bundle) — one file per template.

**Structure:** YAML frontmatter + markdown body.

**Frontmatter fields (required):**

| Field | Type | What it declares |
|---|---|---|
| `name` | string | Template id — must match the filename (minus `.md`) |
| `description` | string (block scalar) | What the interview elicits and produces |
| `version` | semver string | Template version — bump when layer shape changes |
| `topic_cluster` | string | Copied verbatim into the `topic_cluster:` frontmatter of every file the `ingest` sub-command writes to `kb/raw/inbox/` |
| `layers` | list of `{id, title, order}` | Layers in declared order — determines interview sequence |
| `exports` | list of `{path, kind}` | Artifacts produced by the `export` sub-command |

**Markdown body (required sections):**

- One `## Layer N: <title>` section per layer, in declared order. Each layer section contains:
  - A one-paragraph layer purpose.
  - `### Prompt patterns` — 3+ concrete, recent-example questions. Layer 1 questions should anchor on the last one-to-two weeks.
  - `### Details shape` — layer-specific fields that extend the canonical entry contract's `details` object.
- `## Cross-layer contradiction checks` — named tensions the `review` sub-command surfaces between layers.
- `## Tone` — one-paragraph stylistic brief for the Interviewer crew persona during this template.

**Example skeleton:**

```markdown
---
name: my-template
description: |
  What this template elicits and what it produces.
version: 1.0.0
topic_cluster: my-cluster
layers:
  - id: layer_one
    title: Layer One
    order: 1
  - id: layer_two
    title: Layer Two
    order: 2
exports:
  - path: SUMMARY.md
    kind: narrative
  - path: data.json
    kind: machine-readable
---

# <Template Name> Template

Opening paragraph about intent.

## Layer 1: Layer One

One-paragraph purpose.

### Prompt patterns
- "Concrete recent-example question."
- "Another concrete question anchored on last week."
- ...

### Details shape
- `field_one` — what it captures
- `field_two` — list of `{sub, fields}`

## Layer 2: Layer Two
...

## Cross-layer contradiction checks
- **Layer One vs Layer Two** — the specific tension to surface.

## Tone
Stylistic brief for this template.
```

---

## Canonical Entry Contract

Every saved entry, regardless of template, must include all of these fields. Templates extend the `details` field with layer-specific sub-fields — they do not remove or rename the canonical fields.

| Field | Type | Description |
|---|---|---|
| `title` | string | Short name for the entry |
| `summary` | string | 1–2 sentence description |
| `cadence` | string | How often this pattern applies — e.g., `daily`, `weekly`, `per-project` |
| `trigger` | string | What event or condition activates this pattern |
| `inputs` | list of string | What this pattern needs as input |
| `stakeholders` | list of string | Who is involved |
| `constraints` | list of string | Limitations or guardrails |
| `details` | object | Layer-specific object, shape defined per-layer by the template |
| `source_confidence` | enum | `confirmed` (user stated or approved as written) or `synthesized` (abstracted from multiple examples, user approved synthesis) |
| `status` | enum | `active`, `stale`, or `aspirational` |
| `last_validated_at` | ISO 8601 timestamp | When the user last confirmed this entry — set on approval, refreshed on `update` re-run |

The Interviewer never invents these fields. If the user did not provide `constraints`, the Interviewer asks — it does not leave an empty list unchallenged unless the user explicitly says there are none.

---

## Status Vocabulary

Four independent status fields live in session state. They do not share a lifecycle — each one answers a different question. This table is the single source of truth for their transitions.

| Scope | Field | Values | Transitions |
|---|---|---|---|
| Session | `status` | `in_progress`, `complete` | `in_progress` on session start. Flips to `complete` on the final layer's approval write. Only terminal state. Cleared back to `in_progress` only by `reset` (which starts a new session). |
| Layer | `approved` | `false`, `true` | `false` on session start. Flips to `true` when the user explicitly approves the layer's checkpoint. Re-approval on `update` re-runs refreshes `approved_at` and each entry's `last_validated_at`. |
| Entry | `status` | `active`, `stale`, `aspirational` | User-set during the interview. `active` = currently true. `stale` = was true, may no longer hold. `aspirational` = user wants this pattern but it's not yet real. Updated on `update` re-runs when the user reconfirms or relabels an entry. |
| Export freshness | `.exported_at` file | ISO 8601 timestamp | Written by the `export` sub-command on successful run. Read by the next `export` run to surface staleness (session modified after last export). Cleared by `reset`. |

**There is no `superseded` state anywhere.** Prior runs are not marked superseded — they are archived as immutable `versions/v<N>-<date>/` directories. Comparison against a prior version uses the archived `session.json`, not a flag on the current one.

---

## `session.json` Schema

Session state lives at `./do-work/interview/<template>/session.json`. Full shape:

```json
{
  "template": "work-operating-model",
  "session_id": "<uuid>",
  "started_at": "<iso>",
  "last_activity_at": "<iso>",
  "status": "in_progress | complete",
  "pending_layer": "<layer-id> | null",
  "previous_version": "<version-id> | null",
  "review_completed_at": "<iso> | null",
  "review_runs": 0,
  "layers": {
    "<layer-id>": {
      "approved": true,
      "approved_at": "<iso>",
      "entries": [ /* canonical entries */ ]
    }
  }
}
```

**Field semantics:**

- `status` — `in_progress` until every declared layer is `approved: true`; flips to `complete` on the final layer's approval write.
- `pending_layer` — the id of the next layer to interview. `null` when `status: complete`.
- `previous_version` — set when the session was started via `version` re-run mode; carries `v<N>` as a back-reference for comparison queries. Otherwise `null`.
- `review_completed_at` — ISO timestamp of the most recent `review` sub-command completion. `null` until `review` runs to the end of the contradiction list at least once. Cleared only by `reset`.
- `review_runs` — monotonically increasing count of completed `review` passes. Starts at `0`. The `export` sub-command requires `review_completed_at != null && review_runs >= 1`.
- `layers.<layer-id>.entries[]` — each entry matches the canonical entry contract. Every entry is persisted only after the layer's checkpoint was explicitly approved by the user.

**Gate summary:** the `export` sub-command refuses to run unless every layer in the template is `approved: true` AND `review_completed_at != null` AND `review_runs >= 1`.

---

## Checkpoint File Format

After the Interviewer finishes asking a layer's questions and drafts canonical entries, it writes `./do-work/interview/<template>/checkpoints/<layer-id>.md` and presents it to the user in-chat for explicit approval. One file per layer. Checkpoints are transient approval artifacts — they are overwritten on re-run or revision; the authoritative record is `session.json`.

```markdown
# Checkpoint: <layer title>

## Summary
<1–2 paragraph layer summary in the user's own language>

## Entries
### <entry title>
- cadence: <value>
- trigger: <value>
- inputs: <list>
- stakeholders: <list>
- constraints: <list>
- details:
  <layer-specific fields per template>
- source_confidence: confirmed | synthesized

### <next entry title>
...

## Unresolved
<bullet list of unresolved items or "None">

## Approval
This is the <layer-id> model I'd save right now. Correct anything that's off. If this looks right, I'll save it and move to <next-layer-id>.
```

On the final layer of a template, replace the last sentence with: "If this looks right, I'll save it and wrap up the session — you can then run `do work interview <template> review` to surface cross-layer contradictions."

---

## Export Schemas — `work-operating-model`

The `export` sub-command writes these five files to `./do-work/interview/<template>/exports/`. Schemas are referenced from the template's `exports:` declaration; do not duplicate inside the template body.

### `USER.md` — narrative profile

Narrative profile of the person at work. Written in third person, present tense. Required sections, in order:

1. **Name and role** — if known, stated directly. If not captured during interview, omit (do not invent).
2. **Operating rhythm summary** — one paragraph synthesizing Layer 1 entries: when the user works, where the calendar lies, energy patterns.
3. **Key recurring decisions** — 3–5 bullets drawn from Layer 2 entries. Each bullet names the decision and the core judgment heuristic.
4. **Top dependencies** — 3–5 bullets drawn from Layer 3 entries, ordered by frequency or failure impact.
5. **Institutional knowledge** — 2–4 bullets drawn from Layer 4 entries: what the user carries that isn't written down, why it matters.
6. **Active friction points** — 3–5 bullets drawn from Layer 5 entries, ordered by priority.

Quote the user's specific language where possible. Do not add generalizations the interview did not produce.

### `SOUL.md` — decision framework

The decision framework an agent would follow when acting on behalf of this person. Derived primarily from Layer 2 entries, cross-referenced with Layer 1 (rhythm) and Layer 3 (dependencies). Required sections:

1. **When to escalate** — conditions under which the agent must pause and ask the user, not decide autonomously. Sourced from each Layer 2 entry's `escalation_rule`.
2. **When to decide autonomously** — decisions the user has marked `reversible: true` and whose `decision_inputs` the agent can fully observe.
3. **Tone rules for different audiences** — derived from Layer 3 stakeholders plus explicit tone notes the user provided; organized by audience (team, leadership, external, customers, etc.).
4. **Data sources to trust** — union of `decision_inputs` across Layer 2 entries, de-duplicated, labeled with the decision each source feeds.
5. **"Good enough" thresholds** — derived from Layer 2 `thresholds`. For each decision type, the numeric cutoff or qualitative bar the user actually uses.

### `HEARTBEAT.md` — checklist

A checklist an agent reviews on a cadence (default 30 minutes) to decide whether there's work to do for this person. Required sections:

1. **Cadence** — default `every 30 minutes`; override if the user specified a different rhythm.
2. **What to check** — bullet list drawn from Layer 3 `dependency_owner` + `deliverable` pairs, plus Layer 1 `interruptions`. Each bullet names the source/signal and the polling frequency.
3. **Signals to act on** — concrete conditions that should trigger agent action (e.g., "approval from X not received by Y time"); derived from Layer 3 `failure_impact` + `needed_by`.
4. **What to ignore** — noise sources the user explicitly does not want the agent to surface. Derived from Layer 1 `interruptions` the user marked as low-priority.

### `operating-model.json` — machine-readable dump

Full machine-readable dump of the approved session's canonical entries, grouped by layer. Shape:

```json
{
  "template": "work-operating-model",
  "template_version": "<from interviews/work-operating-model.md frontmatter>",
  "exported_at": "<iso>",
  "session_id": "<uuid>",
  "layers": {
    "operating_rhythms": { "entries": [ /* canonical entries */ ] },
    "recurring_decisions": { "entries": [ /* canonical entries */ ] },
    "dependencies": { "entries": [ /* canonical entries */ ] },
    "institutional_knowledge": { "entries": [ /* canonical entries */ ] },
    "friction": { "entries": [ /* canonical entries */ ] }
  }
}
```

### `schedule-recommendations.json` — derived schedule

Derived from Layer 1 `operating_rhythms.time_windows` + `energy_pattern` and Layer 3 `dependencies.needed_by`. Shape:

```json
{
  "generated_at": "<iso>",
  "suggested_time_blocks": [
    { "label": "Deep work", "start": "HH:MM", "end": "HH:MM", "days": ["Mon","Tue"], "rationale": "..." }
  ],
  "standing_slots": [
    { "label": "Dependency handoff — X", "start": "HH:MM", "end": "HH:MM", "days": [...], "source_dependency": "<entry title>" }
  ],
  "avoid_windows": [
    { "label": "Reactive inbox peak", "start": "HH:MM", "end": "HH:MM", "days": [...], "reason": "..." }
  ]
}
```

Every `suggested_time_blocks` entry must trace its rationale to a specific Layer 1 or Layer 3 entry — the agent does not invent blocks that aren't supported by the approved session.

---

## Re-Run Mode Specifications

When the user invokes `do work interview <template>` and the existing `session.json` has `status: complete`, the action prompts the user to choose a re-run mode. The three modes differ in what gets archived, what gets written, and what the CHANGELOG entry looks like.

### `fresh`

**Intent:** start over from scratch, preserve the old run for the record.

**Steps:**
1. Determine next version number `<N>` by scanning existing `versions/` directory (monotonically increasing, starts at `v1`).
2. Create `./do-work/interview/<template>/versions/v<N>-<YYYY-MM-DD>/`.
3. Copy the current `session.json`, the `checkpoints/` directory, and the `exports/` directory into that version folder.
4. Delete the working `checkpoints/` and `exports/` contents (the versioned copy is now the only archive).
5. Write a new empty `session.json` — fresh `session_id`, `started_at: <now>`, `status: in_progress`, `pending_layer: <first-layer-id>`, `previous_version: null`, `review_completed_at: null`, `review_runs: 0`, `layers: {}`.
6. Append to `CHANGELOG.md`:
   ```
   ## <YYYY-MM-DD HH:MM> — fresh start: archived as v<N>
   Previous run archived; new empty session started from scratch.
   ```
7. Begin Layer 1 interview.

### `update`

**Intent:** walk through the prior run and revalidate in place, keep the same session.

**Steps:**
1. Leave `session.json`, `checkpoints/`, `exports/`, and `versions/` untouched.
2. Flip `status` back to `in_progress` and set `pending_layer` to the first layer id.
3. For each layer in declared order:
   - Show the stored canonical entries in a compact form (title + cadence + one-line summary).
   - Ask: "Is this still accurate? Confirm / edit / add / remove entries."
   - Apply the user's changes: edits update entries in place; additions append; removals splice. Update each touched entry's `last_validated_at`.
   - **Empty a layer.** If the user says "remove all," "none of these apply anymore," or similar, the Interviewer may propose an empty layer. It writes a checkpoint with `## Entries` section empty and a layer summary explaining why the layer is empty now (e.g., "no standing dependencies this quarter"). The user must still explicitly approve the empty checkpoint — the gate does not bend. On approval, `layers.<layer-id>.entries` is set to `[]` and `approved_at` is refreshed. An empty layer still counts as approved and does not block `review` or `export`.
   - Write a fresh checkpoint and require explicit approval before committing the edits (same approval gate as a new interview). **Per-entry edit friction is intentional:** fixing one typo regenerates the whole layer's checkpoint. The cost is real; the approval gate is why this recipe is trustworthy. Do not invent a per-entry patch path.
4. When the final layer is confirmed, set `status: complete`, `pending_layer: null`. Do **not** reset `review_completed_at` or `review_runs` — the prior review may still stand, but the user is expected to re-run `review` if the updates changed enough to warrant it.
5. Append to `CHANGELOG.md`:
   ```
   ## <YYYY-MM-DD HH:MM> — layer updated: <layer-id>
   <added N, removed M, edited K> entries. <one-sentence summary of what shifted>
   ```
   One entry per touched layer. Layers with no changes emit no entry.
6. No new version folder is created; this is an in-place update.

### `version`

**Intent:** preserve the old run and start fresh, but with a pointer back for comparison.

**Steps:**
1. Determine next version number `<N>` and archive as in `fresh` (copy session + checkpoints + exports into `versions/v<N>-<YYYY-MM-DD>/`).
2. Clear the working `checkpoints/` and `exports/`.
3. Write a new empty `session.json` — same shape as `fresh`, except `previous_version: "v<N>"`.
4. Append to `CHANGELOG.md`:
   ```
   ## <YYYY-MM-DD HH:MM> — versioned: archived as v<N>, new session seeded
   New session references v<N> as previous_version. Use `do work interview <template> versions` to compare.
   ```
5. Begin Layer 1 interview.

The three modes are mutually exclusive per invocation. The user picks one; the action does not combine them.

---

## Versioning Scheme

- Archive directories are named `v<N>-<YYYY-MM-DD>/` under `./do-work/interview/<template>/versions/`.
- `<N>` is monotonically increasing per template; determine the next number by scanning existing version directory names and taking `max(N) + 1`. The first archive is `v1`.
- A session can reference a prior version via `previous_version: "v<N>"` (set only by the `version` re-run mode).
- **Versions are immutable.** The action never edits files inside `versions/`. If the user wants to amend a prior version, they re-interview from scratch and reference it via `previous_version`.
- The `<template> versions` sub-command enumerates the directory with one line per version: `v<N>-<date>   <layer-count> layers   <entry-count> entries`. Counts come from the archived `session.json`.

---

## Ingest Frontmatter

The `ingest` sub-command copies exports into `<repo-root>/kb/raw/inbox/` with filenames of the shape `interview-<template>-<export-basename>.md`. Each ingested file gets YAML frontmatter with these fields:

```yaml
---
title: <template-display-name> — <export-title>
source: ./do-work/interview/<template>/exports/<export-filename>
type: source-summary
topic_cluster: <value from template frontmatter>
confidence: high
created: <YYYY-MM-DD>
---
```

`topic_cluster` is copied verbatim from the template's frontmatter. `confidence: high` reflects that the source is the user's own approved operating model, not a third-party claim. The file body is the export content (for `.md` exports) or a short pointer describing the export shape + location (for `.json` exports — BKB does not ingest raw JSON).

If `kb/` does not exist when `ingest` is invoked, the action tells the user to run `do work bkb init` first and stops without writing.
