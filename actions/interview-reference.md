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
| Export freshness | `session.json.last_exported_at` field | ISO 8601 timestamp or `null` | Written by the `export` sub-command on successful run. Read by the next `export` run to surface staleness (session modified after last export). `null` on a fresh session. Reset back to `null` by `fresh`, `version`, or `reset` (all of which write a new empty `session.json`). The stamp deliberately lives on `session.json` and **not** inside `exports/` — the `ingest` sub-command iterates every file in `exports/`, so a sidecar stamp there would land in `kb/raw/inbox/` as a bogus document. |

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
  "last_exported_at": "<iso> | null",
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
- `last_exported_at` — ISO timestamp of the most recent successful `export` sub-command run. `null` until the first export. The `export` sub-command compares this against `last_activity_at` in its freshness preflight to surface staleness (session modified after last export). Reset to `null` by `fresh`, `version`, and `reset`. Not a gate — preflight only announces, never blocks.
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

## Export Schemas

The `export` sub-command writes files to `./do-work/interview/<template>/exports/` using render templates defined in the template file itself. For `work-operating-model`, the five mechanical render templates live in `interviews/work-operating-model.md` under `## Export Templates` and use handlebars-style `{{field}}` plus `{{#each …}}` syntax against the canonical entry contract and layer-specific `details`. An implementation should render those templates mechanically against the approved session state.

This section documents the framework-level guarantees every rendered export must satisfy, regardless of template:

- **Narrative exports** (`USER.md` and equivalents) are written in third person, present tense. Quote the user's specific language where possible. Do not add generalizations the interview did not produce. If a field was not captured, omit it — do not invent.
- **Decision-framework exports** (`SOUL.md` and equivalents) derive from entries with `source_confidence` in (`confirmed`, `synthesized`). They must distinguish autonomous-action decisions from escalation decisions, and must list data sources by trust hierarchy.
- **Checklist exports** (`HEARTBEAT.md` and equivalents) must declare a review cadence. Each check item must name its source signal and the layer entry it derives from.
- **Machine-readable exports** (`.json`) serialize canonical entries verbatim with all 11 required fields from the entry contract. Shape is template-specific; traceability is not optional.
- **Derived scheduling exports** (`schedule-recommendations.json` and equivalents) must trace every emitted block, window, or slot to a specific layer entry via `source_entries` — the agent does not invent data unsupported by the approved session.

To add a new template, define its render templates in the template file's `## Export Templates` section using the same handlebars syntax. Do not duplicate per-template rendering here.

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
5. Write a new empty `session.json` — fresh `session_id`, `started_at: <now>`, `last_activity_at: <now>`, `status: in_progress`, `pending_layer: <first-layer-id>`, `previous_version: null`, `review_completed_at: null`, `review_runs: 0`, `last_exported_at: null`, `layers: {}`.
6. Append to `CHANGELOG.md`:
   ```
   ## <YYYY-MM-DD HH:MM> — fresh start: archived as v<N>
   Previous run archived; new empty session started from scratch.
   ```
7. Begin Layer 1 interview.

### `update`

**Intent:** walk through the prior run and revalidate in place at entry-level granularity, keep the same session.

**Steps:**
1. Leave `session.json`, `checkpoints/`, `exports/`, and `versions/` untouched. Initialize an in-memory `any_edits = false` flag for this run.
2. Flip `status` back to `in_progress` and set `pending_layer` to the first layer id.
3. For each layer in declared order, walk each entry individually. For each entry in the prior session's `layers.<layer-id>.entries`:
   1. Display the entry verbatim (all fields).
   2. Prompt: `Still accurate? [confirm / edit / mark-stale / delete / skip]`
   3. On `confirm`: set `source_confidence: confirmed`, update `last_validated_at: <now>`, leave all other fields unchanged.
   4. On `edit`: enter an interactive edit — show current values, let the user override any field, produce a new checkpoint for this entry only, save after approval. Set `last_validated_at: <now>`. Sets `any_edits = true`.
   5. On `mark-stale`: set `status: stale`, update `last_validated_at: <now>`. The entry remains but is flagged in exports. Sets `any_edits = true`.
   6. On `delete`: remove the entry from `layers.<layer-id>.entries`. Log the deletion with full prior content in the CHANGELOG. Sets `any_edits = true`.
   7. On `skip`: leave `last_validated_at` unchanged; the entry carries forward without revalidation.
4. After walking all existing entries in a layer, offer to add new entries by running the layer's original prompts. New entries follow the normal interview flow (canonical entry contract, per-entry approval). Each addition sets `any_edits = true`.
5. **Empty a layer.** If the user deletes every entry in a layer, the Interviewer proposes an empty layer with a short summary explaining why (e.g., "no standing dependencies this quarter"). The user explicitly approves the empty state. On approval, `layers.<layer-id>.entries` is set to `[]` and `approved_at` is refreshed. An empty layer still counts as approved and does not block `review` or `export`.
6. Once all entries are processed in a layer, re-approve the layer as a whole — the layer-level approval gate still applies, now recording that the entry-level walk completed. Set `layers.<layer-id>.approved_at = <now>`.
7. When the final layer is confirmed, set `status: complete`, `pending_layer: null`. **If `any_edits` is true**, also reset `review_completed_at = null` and `review_runs = 0` — the prior review covered a superseded version of the model, and the export gate must force the user back through the cross-layer contradiction pass before the next `export`. If every layer was re-confirmed without edits (`any_edits` stayed `false`), leave `review_completed_at` and `review_runs` untouched.
8. Append to `CHANGELOG.md`, one entry per touched layer:
   ```
   ## <YYYY-MM-DD HH:MM> — layer updated: <layer-id>
   <N confirmed, N edited, N marked stale, N deleted, N added>. <one-sentence summary of what shifted>
   ```
   Layers with no changes (all `confirm` or `skip`) emit no entry.
9. No new version folder is created; this is an in-place update.

### `version`

**Intent:** preserve the old run and start fresh, but with a pointer back for comparison.

**Steps:**
1. Determine next version number `<N>` and archive as in `fresh` (copy session + checkpoints + exports into `versions/v<N>-<YYYY-MM-DD>/`).
2. Clear the working `checkpoints/` and `exports/`.
3. Write a new empty `session.json` — same shape as `fresh` (including `last_activity_at: <now>`), except `previous_version: "v<N>"`.
4. Append to `CHANGELOG.md`:
   ```
   ## <YYYY-MM-DD HH:MM> — versioned: archived as v<N>, new session seeded
   New session references v<N> as previous_version. Use `do work interview <template> versions` to compare.
   ```
5. Begin Layer 1 interview.

The three modes are mutually exclusive per invocation. The user picks one; the action does not combine them.

### Mid-layer recovery

Session state is only written after a layer's checkpoint is approved. If the user quits in the middle of a layer interview — before that layer's checkpoint was approved — on resume:

1. Detect that `pending_layer` has no approved entries in `layers.<layer-id>.entries`.
2. Check for a draft checkpoint file at `./do-work/interview/<template>/checkpoints/.draft-<layer-id>.md`. (Draft checkpoints are written opportunistically during the interview — before user approval — as a recovery aid.)
3. If a draft exists: show it to the user and ask "pick up from this draft or start the layer over?" On `pick up`, load the draft entries as working state and continue from the approval step. On `start over`, delete the draft and begin the layer fresh.
4. If no draft exists: begin the layer fresh.

The action writes draft checkpoints after the interview has produced candidate entries but before user approval. Drafts are deleted when the layer is approved (normal case) or explicitly discarded (start-over case).

---

## Versioning Scheme

- Archive directories are named `v<N>-<YYYY-MM-DD>/` under `./do-work/interview/<template>/versions/`.
- `<N>` is monotonically increasing per template; determine the next number by scanning existing version directory names and taking `max(N) + 1`. The first archive is `v1`.
- A session can reference a prior version via `previous_version: "v<N>"` (set only by the `version` re-run mode).
- **Versions are immutable.** The action never edits files inside `versions/`. If the user wants to amend a prior version, they re-interview from scratch and reference it via `previous_version`.
- The `<template> versions` sub-command enumerates the directory with one line per version: `v<N>-<date>   <layer-count> layers   <entry-count> entries`. Counts come from the archived `session.json`.

---

## Ingest File Mapping

When `do work interview <template> ingest` runs, it produces files in `<repo-root>/kb/raw/inbox/` (sibling writes — do not overwrite). Two file classes are written per run: one file per export and one summary file per layer.

### 1. One file per export

Filename: `<template>-<export-name>.md`. Body is the full export content. Frontmatter:

```yaml
---
title: "{{template.name}} — {{export_filename_without_ext}}"
type: source-summary
topic_cluster: "{{template.topic_cluster}}"
sources:
  - "interview/{{template}}/exports/{{export_filename}}"
confidence: high
created: "{{session.last_exported_at}}"
---
```

For `.json` exports, include the JSON as a fenced code block inside the markdown body (BKB does not ingest raw JSON).

### 2. One summary file per layer

Filename: `<template>-<layer-id>.md`. Body is a markdown summary of that layer's entries (list each entry's `title` and `summary` under the layer heading). Frontmatter:

```yaml
---
title: "{{template.name}} — {{layer.title}}"
type: concept
topic_cluster: "{{template.topic_cluster}}"
sources:
  - "interview/{{template}}/session.json"
related:
  - page: "{{template}}-user-md"
    rel: evidence-for
confidence: "{{majority source_confidence in layer — confirmed => high, synthesized => medium}}"
created: "{{session.last_exported_at}}"
---
```

This gives BKB one wiki page per layer alongside the full exports.

### 3. Inbox manifest

Append one row to `kb/raw/_inbox_queue.md` for each file added. Each row is marked `ready`, with `topic_hint: {{template.topic_cluster}}` and `priority: normal`.

### Totals and collisions

For the `work-operating-model` template: 5 exports + 5 layer summaries = **10 files** per `ingest` run. If any target filename already exists in `kb/raw/inbox/` or `kb/raw/capture/` (previous ingest of the same template), prefix the new file with the current time (`HHMMSS-<filename>`) per BKB's collision rule.

### Preconditions

If `kb/` does not exist when `ingest` is invoked, the action tells the user to run `do work bkb init` first and stops without writing. `confidence: high` on export files reflects that the source is the user's own approved operating model, not a third-party claim.
