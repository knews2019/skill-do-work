# v2 — close gaps in the `interview` action

This patch builds on v1. The v1 specification is stored in this repo at `decisions/imported-specs/expand-skill-do-work-interview.md` and is being (or has been) implemented into the `interview` action, companion reference, user guide, crew persona, first template, and ADR-011. Some parts of v1 were underspecified; this v2 patch closes those gaps. Do not redo v1's work. Do not rename the action or restructure the files v1 created. Patch, don't rewrite.

---

## Step 0 — Orient

Before writing anything, read:

1. `decisions/imported-specs/expand-skill-do-work-interview.md` — the v1 spec. Every decision in v1 still stands unless this v2 patch explicitly overrides it.
2. `actions/interview.md` and `actions/interview-reference.md` — as they actually exist in the repo now. V1 may have been implemented with variations; this patch adjusts the real file, not the spec.
3. `interviews/work-operating-model.md` — this is where most of your edits land.
4. `decisions/records/adr-011-interview-framework-with-prescriptive-templates.md` — you will add a sibling ADR-012 to cover v2's decisions.
5. `crew-members/` directory — you will audit whether this is the right home for `interviewer.md`. See §2.1.

If any v1 file is missing or substantively incomplete, stop and report. Do not try to complete v1's work as part of v2.

---

## Step 0.5 — Propose a plan, wait for approval

Before making any edits, produce a short written plan covering:

- Which files you will create (expected: 0–1 new ADR)
- Which files you will edit (expected: `interviews/work-operating-model.md`, `actions/interview-reference.md`, possibly `actions/interview.md`, possibly a crew file move)
- Your proposed resolution for the crew placement audit in §2.1 (move or keep)
- An estimate of the change size (lines added/modified per file)

Present the plan, then stop and wait for explicit user confirmation ("approved," "go ahead," or corrections). Do not start making edits until confirmed. This checkpoint was missing in v1 and is non-negotiable for v2.

---

## Step 1 — Close severe gaps

### 1.1 Embed concrete export templates in the Work Operating Model template

**Problem.** V1 described the five exports (`USER.md`, `SOUL.md`, `HEARTBEAT.md`, `operating-model.json`, `schedule-recommendations.json`) in prose only. An implementing agent will invent structure on first run, and different runs will produce different file shapes — defeating the "feed USER.md into any agent platform" guarantee.

**Fix.** Extend `interviews/work-operating-model.md` with an `## Export Templates` section containing explicit, substitutable templates for all five exports. Use `{{…}}` for field references against the canonical entry contract and layer-specific details. Use `{{#each …}} … {{/each}}` for per-entry iteration. An implementation of the `interview <template> export` sub-command should be able to render these templates mechanically against the approved session state.

Add this section verbatim:

~~~markdown
## Export Templates

When the `export` sub-command runs against an approved session of this template, render each artifact below using the session's canonical entries. Field references use handlebars-style `{{field}}` syntax against the canonical entry contract plus layer-specific `details`. Iteration uses `{{#each layers.<layer_id>.entries}} … {{/each}}`. Omit sections whose source layer has no qualifying entries.

### `USER.md` — narrative profile

```markdown
# Work Operating Model — {{session.role_or_name_or_repo}}

_Generated {{session.completed_at}}. Based on the work-operating-model template, version {{template.version}}._

## How the week actually runs

{{synthesis_paragraph from operating_rhythms — describe time_windows, energy_pattern, and non_calendar_reality in 2–3 sentences}}

### Deep work windows
{{#each operating_rhythms.entries where details.energy_pattern mentions "deep" or "focus"}}
- {{details.time_windows}} — {{summary}}
{{/each}}

### What the calendar hides
{{#each operating_rhythms.entries}}
- {{details.non_calendar_reality}}
{{/each}}

## Recurring decisions

{{#each recurring_decisions.entries}}
### {{details.decision_name}}
- **Cadence:** {{cadence}}. **Trigger:** {{trigger}}.
- **Inputs:** {{details.decision_inputs}}
- **Thresholds:** {{details.thresholds}}
- **Escalate when:** {{details.escalation_rule}}
- **Reversible:** {{details.reversible}}
{{/each}}

## Dependencies

{{#each dependencies.entries}}
- **{{details.dependency_owner}}** — {{details.deliverable}}, needed {{details.needed_by}}.
  Failure impact: {{details.failure_impact}}. Fallback: {{details.fallback}}.
{{/each}}

## Institutional knowledge I carry

{{#each institutional_knowledge.entries}}
- **{{details.knowledge_area}}** — {{details.why_it_matters}}.
  Currently lives: {{details.where_it_lives}}. Partial sharers: {{details.who_else_knows}}.
  Risk if missing: {{details.risk_if_missing}}.
{{/each}}

## Active friction

{{#each friction.entries sorted by details.priority desc, details.time_cost desc}}
- [{{details.priority}}] **{{title}}** — {{details.frequency}}, ~{{details.time_cost}} per occurrence.
  Workaround: {{details.current_workaround}}. Systems: {{details.systems_involved}}.
  Automation candidate: {{details.automation_candidate}}.
{{/each}}
```

### `SOUL.md` — agent decision framework

```markdown
# Agent Operating Instructions

_Use this file to decide how to act on behalf of the user described in `USER.md`. Do not override these rules with defaults inferred from general context._

## When to act autonomously

{{#each recurring_decisions.entries where details.reversible == true}}
- **{{details.decision_name}}**: apply the thresholds in `USER.md` and act. Do not escalate for this decision class.
{{/each}}

## When to escalate

{{#each recurring_decisions.entries where details.escalation_rule exists and details.escalation_rule != "never"}}
- **{{details.decision_name}}**: {{details.escalation_rule}}
{{/each}}

Additionally, always escalate when:
- A dependency from `USER.md` is late and its fallback is not defined
- An institutional_knowledge item marked `risk_if_missing: high` is needed but the owner is unreachable
- Any threshold in `USER.md` is within 10% of being crossed and the decision is irreversible

## Data sources — trust hierarchy

**Authoritative** (cite these directly, do not second-guess):
{{items appearing in 2+ recurring_decisions.details.decision_inputs}}

**Advisory** (consider, but cross-check before acting):
{{items appearing in only 1 recurring_decisions.details.decision_inputs}}

**Tacit** (do not assume present; ask the user if needed):
{{institutional_knowledge.entries where details.where_it_lives contains "head" or "undocumented"}}

## Tone rules by audience

{{#for each unique stakeholder across all layers}}
- **{{stakeholder}}**: {{derived tone — terse/formal/informal based on which layers they appear in}}
{{/for}}

## "Good enough" thresholds

{{#each recurring_decisions.entries}}
- For **{{details.decision_name}}**: proceed when {{details.thresholds}} are met. Do not hold for perfection.
{{/each}}

## What never to do

- Do not act on behalf of the user in a domain not covered by `USER.md`.
- Do not fabricate information for a decision whose `decision_inputs` are unavailable.
- Do not smooth over contradictions between `USER.md` sections — surface them.
```

### `HEARTBEAT.md` — recurring checklist

```markdown
# Heartbeat Checklist

_Review on a 30-minute cadence. For each item: act, defer, or ignore. Log the decision._

## Every heartbeat

- Scan `USER.md` dependencies. Any expected deliverable past its `needed_by` window?
  - If yes and `fallback` is defined → execute fallback per `SOUL.md`
  - If yes and no fallback → escalate
- Scan `USER.md` recurring decisions. Any whose `cadence` or `trigger` fires now?
  - If yes → pull `decision_inputs`, apply `thresholds`, act or escalate per `SOUL.md`

## First heartbeat after 08:00 local

- Load today's calendar. Compare to deep work windows in `USER.md`. Flag conflicts.
- Scan these sources for overnight changes: {{list derived from operating_rhythms.details.non_calendar_reality + recurring_decisions.details.decision_inputs}}

## First heartbeat Monday after 08:00

- Review last week's friction log from `USER.md`. Any high-priority items unchanged? Flag for user.
- For each `institutional_knowledge` entry with `risk_if_missing: high`: was this knowledge used last week? By whom? Log.

## First heartbeat on the 1st of the month

- Produce a one-page delta: what in `USER.md` no longer matches reality? Flag for the user's next quarterly interview re-run.
```

### `operating-model.json` — machine-readable dump

```json
{
  "template": "work-operating-model",
  "template_version": "{{template.version}}",
  "session_id": "{{session.session_id}}",
  "generated_at": "{{session.completed_at}}",
  "previous_version": "{{session.previous_version}}",
  "layers": {
    "operating_rhythms": { "entries": [ "{{canonical_entries}}" ] },
    "recurring_decisions": { "entries": [ "{{canonical_entries}}" ] },
    "dependencies": { "entries": [ "{{canonical_entries}}" ] },
    "institutional_knowledge": { "entries": [ "{{canonical_entries}}" ] },
    "friction": { "entries": [ "{{canonical_entries}}" ] }
  }
}
```

Canonical entries are serialized verbatim with all 11 required fields from the entry contract.

### `schedule-recommendations.json` — derived scheduling data

```json
{
  "generated_at": "{{session.completed_at}}",
  "source_template": "work-operating-model",
  "source_session": "{{session.session_id}}",
  "time_blocks": [
    {
      "label": "{{derived from operating_rhythms.details.energy_pattern}}",
      "days": ["{{day}}"],
      "start": "HH:MM",
      "end": "HH:MM",
      "type": "deep_work | admin | reactive",
      "source_entries": ["operating_rhythms.<entry_id>"]
    }
  ],
  "avoid_windows": [
    {
      "label": "{{why this window should be protected}}",
      "days": ["{{day}}"],
      "start": "HH:MM",
      "end": "HH:MM",
      "reason": "{{non_calendar_reality or friction source}}",
      "source_entries": ["<layer>.<entry_id>"]
    }
  ],
  "standing_slots": [
    {
      "label": "{{deliverable or handoff}}",
      "cadence": "{{from dependency entry}}",
      "day": "{{day or 'rolling'}}",
      "time": "HH:MM",
      "counterparty": "{{dependency_owner}}",
      "source_entries": ["dependencies.<entry_id>"]
    }
  ]
}
```

Derivation rules:
- `time_blocks` come from `operating_rhythms.entries[*].details.time_windows` joined against `details.energy_pattern` to classify type.
- `avoid_windows` come from (a) `operating_rhythms.details.non_calendar_reality` when it describes a consistent interruption pattern, and (b) high-priority `friction.entries` whose `systems_involved` implies a recurring time loss.
- `standing_slots` come from `dependencies.entries[*]` where `needed_by` has a regular cadence.
~~~

### 1.2 Plan-and-approve checkpoint is already addressed in §0.5 above

This pattern is now part of the v2 workflow itself. It does not require a change to any v1 file.

---

## Step 2 — Close material gaps

### 2.1 Audit crew placement for `interviewer.md`

**Problem.** V1 placed the interviewer persona in `crew-members/interviewer.md` by analogy to existing members (backend, frontend, debugging, etc.). On review, the existing crew members look specifically like personas for the `work` action — code-focused roles invoked during development tasks — not a generic cross-action pool.

**Fix.** Read `crew-members/*.md` in full. Determine empirically: is this a generic persona pool available to any action, or is it scoped to the `work` action?

- **If generic**: leave `crew-members/interviewer.md` where it is. No change needed.
- **If scoped to `work`**: move the interviewer persona inline into `actions/interview-reference.md` as a new `## Interviewer Persona` section. Delete `crew-members/interviewer.md`. Update any references in `actions/interview.md` (the sub-command `Crew` column) to say "see interview-reference.md §Interviewer Persona" instead of the `Interviewer` shorthand, or leave the shorthand and let the reference file be the single source of truth — whichever matches the repo's pattern for action-specific personas.
- **If unclear**: propose both options in §0.5's plan and let the user decide.

Record the resolution in ADR-012 (see §3).

### 2.2 Specify `update` re-run mode granularity

**Problem.** V1 said `update` mode loads prior layers and asks "is each layer still accurate?" but didn't specify whether this is per-layer (approve/edit the whole layer at once) or per-entry (walk through each entry individually). The granularity affects how the CHANGELOG records diffs and what `last_validated_at` means.

**Fix.** In `actions/interview-reference.md`, extend the re-run mode specification with this detail:

~~~markdown
### Update mode — entry-level granularity

When `update` mode loads a layer, walk each entry individually:

For each entry in the prior session's `layers.<layer_id>.entries`:
1. Display the entry verbatim (all fields).
2. Prompt: "Still accurate? [confirm / edit / mark-stale / delete / skip]"
3. On `confirm`: set `source_confidence: confirmed`, update `last_validated_at: <now>`, leave all other fields unchanged.
4. On `edit`: enter an interactive edit — show current values, let user override any field, produce a new checkpoint for this entry only, save after approval. Set `last_validated_at: <now>`.
5. On `mark-stale`: set `status: stale`, update `last_validated_at: <now>`. The entry remains but is flagged in exports.
6. On `delete`: remove the entry from `layers.<layer_id>.entries`. Log the deletion with full prior content in the CHANGELOG.
7. On `skip`: leave `last_validated_at` unchanged; the entry carries forward without revalidation.

After walking all existing entries, offer to add new entries by running the layer's original prompts. New entries follow the normal interview flow.

Once all entries are processed, re-approve the layer as a whole. The CHANGELOG entry for an update run summarizes: `N confirmed, N edited, N marked stale, N deleted, N added`.
~~~

### 2.3 Specify mid-layer recovery on resume

**Problem.** V1 said resume picks up from the pending layer but didn't address what happens if the user quit mid-conversation during a layer, before that layer's checkpoint was approved. The pending layer's `entries` may be empty but the user has already answered some questions in chat.

**Fix.** In `actions/interview-reference.md`, add this to the session lifecycle section:

~~~markdown
### Mid-layer recovery

Session state is only written after a layer's checkpoint is approved. If the user quits in the middle of a layer interview, on resume:

1. Detect that `pending_layer` has no approved entries in `layers.<layer_id>.entries`.
2. Check for a draft checkpoint file at `./interview/<template>/checkpoints/.draft-<layer-id>.md`. (Draft checkpoints are written opportunistically during the interview — before user approval — as a recovery aid.)
3. If a draft exists: show it to the user and ask "pick up from this draft or start the layer over?" On "pick up", load the draft entries as working state and continue from the approval step. On "start over", delete the draft and begin the layer fresh.
4. If no draft exists: begin the layer fresh.

The action should write draft checkpoints after the interview has produced candidate entries but before user approval. Drafts are deleted when the layer is approved (normal case) or explicitly discarded (start-over case).
~~~

Also amend the layer interview workflow in `actions/interview.md`: between "convert responses into canonical entries" and "write a checkpoint file," add a step "write a draft checkpoint to `./interview/<template>/checkpoints/.draft-<layer-id>.md` so the layer can be recovered on resume."

### 2.4 Specify `ingest` sub-command file mapping and frontmatter

**Problem.** V1 gave one example frontmatter block but didn't specify which files get copied to `kb/raw/inbox/` or how multi-file sources should be represented for `bkb`.

**Fix.** In `actions/interview-reference.md`, add this section:

~~~markdown
### Ingest sub-command file mapping

When `interview <template> ingest` runs, produce these files in `kb/raw/inbox/` (sibling, do not overwrite):

1. One file per export: `<template>-<export-name>.md` containing the full export content with frontmatter:

```yaml
---
title: "{{template.name}} — {{export_filename_without_ext}}"
type: source-summary
topic_cluster: operating-model
sources:
  - "interview/{{template}}/exports/{{export_filename}}"
confidence: high
created: "{{session.completed_at}}"
---
```

For JSON exports, include the JSON as a fenced code block inside the markdown body.

2. One summary file per layer: `<template>-<layer-id>.md` containing a markdown summary of that layer's entries with frontmatter:

```yaml
---
title: "{{template.name}} — {{layer.title}}"
type: concept
topic_cluster: operating-model
sources:
  - "interview/{{template}}/session.json"
related:
  - page: "{{template}}-user-md"
    rel: evidence-for
confidence: "{{majority source_confidence in layer — confirmed => high, synthesized => medium}}"
created: "{{session.completed_at}}"
---
```

The layer summary body lists each entry's `title` and `summary`, grouped under the layer heading. This gives `bkb` one wiki page per layer alongside the full exports.

3. A manifest entry appended to `kb/raw/_inbox_queue.md` for each file added, marked `ready`, with `topic_hint: operating-model` and `priority: normal`.

Total files added per ingest: 5 exports + 5 layer summaries = 10 files for the work-operating-model template.

If any file already exists in `kb/raw/inbox/` or `kb/raw/capture/` (previous ingest of the same template), ingest prefixes the new files with the current time (`HHMMSS-<filename>`) per `bkb`'s collision rule.
~~~

### 2.5 schedule-recommendations.json derivation is now covered

This was addressed in §1.1 as part of the export templates section. No separate fix needed.

---

## Step 3 — Add ADR-012

Create `decisions/records/adr-012-interview-v2-gap-closure.md` following the format of existing ADRs.

**Status.** Accepted.

**Context.** V1 (spec at `decisions/imported-specs/expand-skill-do-work-interview.md`, implementation in ADR-011) shipped the `interview` action but left five gaps: export schemas were described in prose instead of specified as templates, the crew placement for `interviewer.md` assumed a convention that may not hold, `update` re-run mode granularity was unspecified, mid-layer recovery behavior was unspecified, and `ingest` sub-command file mapping was underspecified.

**Decision.** Close the gaps through surgical patches rather than a rewrite. Embed concrete export templates in the Work Operating Model template file. Audit and resolve crew placement against the repo's actual convention. Specify entry-level granularity for update mode. Specify draft-checkpoint-based mid-layer recovery. Specify the exact file mapping and frontmatter for `ingest`. Record the plan-and-approve checkpoint as a workflow norm for future patches.

**Consequences.**

- Exports are reproducible across runs and across implementations.
- `ingest` produces a predictable shape in `kb/raw/inbox/`, making the cross-action seam reliable.
- `update` mode supports partial re-validation, preserving long-lived entries that remain accurate.
- Mid-layer quits are recoverable without re-answering.
- Crew placement reflects the repo's actual pattern rather than a misread analogy.
- Future prompt-level changes to this skill follow the plan-and-approve pattern established in §0.5.

**Alternatives considered.**

- Full rewrite of v1 to embed every gap fix — rejected as wasteful; v1's core architecture is sound.
- Per-layer granularity for update mode — rejected; coarser than users need, forces re-interviewing layers where only one entry went stale.
- Reconstruct mid-layer state from chat history — rejected as unreliable; draft checkpoints are explicit and verifiable.

Update `decisions/_master_index.md` and the relevant topic indexes to include ADR-012.

---

## Step 4 — Dry-run sanity check

Before declaring done, walk one fake session end-to-end as a mental test:

1. Assume `do-work interview work-operating-model` is invoked in a fresh repo.
2. Session starts. Layer 1 (operating_rhythms) runs. User answers 3 questions. Agent produces a draft checkpoint, then a formal checkpoint. User says "save." Entry is written.
3. User quits before Layer 2 starts.
4. Next day: `do-work interview work-operating-model` resumes. Loads session, sees `pending_layer: recurring_decisions`, starts Layer 2.
5. Layers 2, 3, 4, 5 complete normally.
6. `do-work interview work-operating-model review` runs. One tension surfaced (rhythms deep-work window collides with dependency handoff). User revises rhythms. Layer 1 resaved.
7. `do-work interview work-operating-model export` runs. All five files render using the templates in §1.1.
8. `do-work interview work-operating-model ingest` runs. 10 files land in `kb/raw/inbox/` with correct frontmatter. `_inbox_queue.md` gets 10 new rows.
9. Three months later: `do-work interview work-operating-model`. Session is `complete`. User picks `update` mode. Agent walks 14 entries across 5 layers; 8 confirmed, 3 edited, 2 marked stale, 1 deleted, 2 added. CHANGELOG logs the diff.

If any step surfaces an unspecified behavior, add a §2.x fix. If all steps round-trip cleanly, proceed.

---

## Step 5 — Verify

- [ ] `interviews/work-operating-model.md` has a new `## Export Templates` section with all five export templates.
- [ ] `actions/interview-reference.md` has new sections covering update-mode granularity, mid-layer recovery, and ingest file mapping.
- [ ] `actions/interview.md` references draft checkpoints in the layer interview workflow (§2.3).
- [ ] Crew placement for `interviewer.md` has been audited; either left in `crew-members/` with justification, or moved inline to `actions/interview-reference.md` with `crew-members/interviewer.md` deleted.
- [ ] `decisions/records/adr-012-interview-v2-gap-closure.md` exists with the content from §3.
- [ ] `decisions/_master_index.md` and relevant topic indexes list ADR-012.
- [ ] No file in this patch references OB1, Open Brain, MCP, Supabase, or any external service.
- [ ] The dry-run walkthrough in §4 round-trips without surfacing unspecified behavior.

Report what was edited and created. Stop without committing. User reviews the diff.
