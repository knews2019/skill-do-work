# Interview Action

> **Part of the do-work skill.** Generalized interview framework. Runs prescriptive templates that elicit tacit knowledge through structured multi-layer conversations and produce agent-ready operating artifacts.

The action loads templates from `interviews/<name>.md`, runs a checkpoint-gated interview layer by layer, and produces artifacts the user can hand to an agent or feed into `bkb` as queryable knowledge. Session state persists at `./interview/<template>/session.json` and resumes across sessions per ADR-005. Heavy content — template format, canonical entry contract, session schema, export schemas, re-run mode specs — lives in the companion `actions/interview-reference.md` per ADR-001. The `ingest` sub-command produces files that land in `kb/raw/inbox/` in the format `bkb triage && bkb ingest` expects, per ADR-002.

## When to Use

**Use when:**
- The user wants to produce agent-ready operating artifacts (`USER.md` / `SOUL.md` / `HEARTBEAT.md`) and needs the five-layer Work Operating Model interview to get there.
- The user wants to onboard a delegate (human or agent) by making their tacit work patterns explicit.
- The user's operating context has shifted (new role, new team, new responsibilities) and the prior operating model needs a refresh — pick `update` or `version` re-run mode depending on whether the old run should be preserved.
- A new template has been authored in `interviews/` and the user wants to run that interview.

**Do NOT use when:**
- The user wants the agent to act on their behalf right now — that's the `work` or `pipeline` action. Interview produces the instructions; it does not execute them.
- The user wants a list of ideas or a brainstorm — use `scan-ideas` or `deep-explore`. Interview elicits structure, not possibilities.
- The user wants to review tacit knowledge that already exists as code or documentation — use `code-review`, `prime`, or `bkb query`. Interview is for knowledge that only lives in the user's head.

## Input

`$ARGUMENTS` contains the sub-command and optional template name plus modifiers.

| Invocation | Behavior |
|---|---|
| `do work interview` (no args) | Show help menu |
| `do work interview list` | List available templates |
| `do work interview <template>` | Start or resume interview for `<template>` |
| `do work interview <template> status` | Show session progress |
| `do work interview <template> review` | Run contradiction pass |
| `do work interview <template> export` | Write export artifacts |
| `do work interview <template> ingest` | Copy exports into BKB inbox |
| `do work interview <template> reset` | Archive as a version and start fresh (requires confirmation) |
| `do work interview <template> versions` | List archived runs |

## Sub-Commands

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

**Crew dispatch.** The named roles (Architect, Reviewer, Editor, Librarian) are narrative labels describing the stance the Interviewer adopts during each sub-command; the only persona file loaded is `crew-members/interviewer.md`. Load it at the start of any sub-command other than `list` and `(none)`.

---

## Locating the Template

When a sub-command references `<template>`, resolve as follows:

1. Read `interviews/<template>.md` from the repo root.
2. If the file does not exist, list available templates (`ls interviews/*.md`) and stop with: `Template '<template>' not found. Run 'do work interview list' to see available templates.`

The template's frontmatter is the contract: layers, per-layer prompts, and export declarations come from this file. The action enforces the contract; it does not improvise around missing fields.

## Locating the Session

Session state lives at `./interview/<template>/session.json` in the current working directory.

- If it does not exist and the sub-command is bare `<template>`: create the directory structure and start a fresh session (see Step 1 below).
- If it does not exist and the sub-command is anything else (`status`, `review`, `export`, `ingest`, `reset`, `versions`): stop with `No session found for '<template>'. Run 'do work interview <template>' to start one.`

---

## Sub-Command: `list`

List every template available in `interviews/`.

### Steps

1. Scan `interviews/*.md` in the repo root.
2. For each template, read the frontmatter `name`, `description`, and `version`.
3. Print a single-line summary per template plus the description on a subsequent indented line.
4. If `interviews/` does not exist or has no `.md` files, print `No templates found in interviews/. Add a template file to get started — see actions/interview-reference.md for the template format.`

### Output

```
Available templates:

  work-operating-model  (v1.0.0)
    Elicits the operating model of a person at work. Produces agent-ready
    artifacts (USER.md, SOUL.md, HEARTBEAT.md) plus machine-readable exports.

Start an interview:  do work interview <template>
```

---

## Sub-Command: `<template>` — Session Lifecycle

The core sub-command. Behavior branches on whether a session exists and its status.

### Step 1: New session (no `session.json`)

Create the directory structure:

```
./interview/<template>/
├── session.json
├── checkpoints/
├── exports/
├── versions/
└── CHANGELOG.md
```

Initial `session.json`:

```json
{
  "template": "<template>",
  "session_id": "<uuid>",
  "started_at": "<iso>",
  "last_activity_at": "<iso>",
  "status": "in_progress",
  "pending_layer": "<first-layer-id-from-template>",
  "previous_version": null,
  "review_completed_at": null,
  "review_runs": 0,
  "layers": {}
}
```

Write a header line to `CHANGELOG.md`: `# Interview CHANGELOG — <template>` followed by a blank line.

Proceed to Step 3 (layer interview workflow) starting at the first layer.

### Step 2: Existing session

Read `session.json`. Branch on `status`:

- **`status: "in_progress"`** — resume. Read `pending_layer` and proceed to Step 3 for that layer. Announce resumption briefly: "Resuming `<template>` at layer <pending_layer>. <N> of <total> layers approved so far." Do not re-show approved layers unless the user asks.

- **`status: "complete"`** — prompt for re-run mode. Present the three options verbatim:
  - `fresh` — archive current state as `versions/v<N>-<date>/`, start a new session with empty state.
  - `update` — load prior layers; for each layer, show stored entries and ask "is this still accurate? confirm / edit / add / remove." Updates in place.
  - `version` — archive current state as `versions/v<N>-<date>/`, start a new session seeded with `previous_version: v<N>`.

  Wait for the user's choice. Execute the chosen mode per the re-run specifications in `actions/interview-reference.md`.

### Step 3: Layer Interview Workflow

For each layer in the template's declared order (starting from `pending_layer`):

1. **Open with a concrete question.** Read the layer's `Prompt patterns` from the template. Pick one concrete, recent-example question and ask it. Never open a layer with an abstract question like "what do you do all day?" — that violates the Interviewer's core standards.

2. **Converse and draft.** Convert the user's responses into canonical entries matching the template's entry contract (see `actions/interview-reference.md` — Canonical Entry Contract). One question at a time. Capture the user's specific language. If a canonical field was not mentioned (e.g., `constraints`), ask — do not invent.

3. **Write the checkpoint.** When the layer has 2–5 canonical entries drafted, write `./interview/<template>/checkpoints/<layer-id>.md` using the Checkpoint File Format from `actions/interview-reference.md`. Include a 1–2 paragraph layer summary, the canonical entries, any unresolved items, and the explicit approval ask.

4. **Present and wait for approval.** Show the checkpoint contents to the user in-chat. Accepted confirmations: "save," "looks right," "confirmed," "approve," or semantic equivalents. Corrections: the user edits specific entries — regenerate the checkpoint and ask again. Never persist unconfirmed content.

5. **Persist on approval.** Write approved entries to `session.json` under `layers.<layer-id>.entries[]`. Set `layers.<layer-id>.approved = true`, `layers.<layer-id>.approved_at = <now>`, `last_activity_at = <now>`. Update each entry's `last_validated_at`. Advance `pending_layer` to the next layer id (or `null` if this was the last layer — and flip `status` to `complete`).

6. **Append to CHANGELOG.** Add one line:
   ```
   ## <YYYY-MM-DD HH:MM> — layer approved: <layer-id>
   <one-sentence summary of the real pattern surfaced>
   ```

7. **Advance.** Move to the next layer. When the final layer is approved, suggest: "All layers approved. Run `do work interview <template> review` to surface cross-layer contradictions, then `export` to write deliverables."

---

## Sub-Command: `<template> status`

Read `session.json` and report.

### Output

```
Interview status — <template>

  Started:       <started_at>
  Last activity: <last_activity_at>
  Status:        in_progress | complete
  Progress:      <approved-count> of <total-layers> layers approved

  Layers:
    [x] operating_rhythms   approved <approved_at>  (<entry-count> entries)
    [x] recurring_decisions approved <approved_at>  (<entry-count> entries)
    [ ] dependencies        pending
    [ ] institutional_knowledge  pending
    [ ] friction            pending

  Review: <N> pass(es), last completed <review_completed_at | never>
  Previous version: <previous_version | none>
```

If `session.json` doesn't exist, print the "no session" message per Locating the Session.

---

## Sub-Command: `<template> review`

Runs the cross-layer contradiction pass. Requires all layers approved.

### Preconditions

- `session.json` exists and every declared layer has `approved: true`. If any layer is unapproved, list the pending layers and stop with: "Review requires all layers approved. Missing: <list>. Run `do work interview <template>` to finish the interview."

### Steps

1. Read the template's `Cross-layer contradiction checks` section. Each named tension is a check to run (e.g., `Rhythm vs Dependencies`).

2. For each check, walk the relevant layer entries and identify pairs that instantiate the tension. Examples for `work-operating-model`:
   - **Rhythm vs Dependencies** — a Layer 1 entry claims a deep-work window that overlaps a Layer 3 dependency handoff.
   - **Decisions vs Knowledge** — a Layer 2 entry's `decision_inputs` references data the user said in Layer 4 isn't written down anywhere.
   - **Friction vs Rhythm** — a Layer 5 friction pattern implies the stated Layer 1 rhythm isn't real.
   - **Dependencies vs Knowledge** — a Layer 3 `dependency_owner` matches a Layer 4 `who_else_knows` entry naming a single point of failure.

3. For each tension found, present it explicitly, naming both entries by title. Ask the user to pick:
   - `revise-A` — rewrite the first entry.
   - `revise-B` — rewrite the second entry.
   - `both-are-true` — note the tension without rewriting.
   - `skip` — move on without recording anything.

4. **Revisions regenerate a checkpoint.** If the user picks `revise-A` or `revise-B`, write a fresh checkpoint for that layer with the revised entry and require explicit approval before overwriting the stored entry in `session.json`. Same approval gate as new interviews.

5. **Both-are-true.** If the user picks `both-are-true`, append a short note to both affected entries' `constraints` list: e.g., `known tension with <other-entry-title>`. No re-approval needed.

6. When every surfaced tension has a resolution (or was skipped), update `session.json`: set `review_completed_at = <now>`, increment `review_runs += 1`, update `last_activity_at`. Append one line to CHANGELOG:
   ```
   ## <YYYY-MM-DD HH:MM> — review pass completed (run <review_runs>)
   <N> tensions surfaced, <M> resolved, <K> skipped.
   ```

---

## Sub-Command: `<template> export`

Writes the template's declared export artifacts to `./interview/<template>/exports/`.

### Preconditions

- Every declared layer has `approved: true`. If not, list missing layers and stop.
- `review_completed_at != null` **AND** `review_runs >= 1`. If not, stop with: "Export requires the review pass to have run at least once. Run `do work interview <template> review` first."

### Steps

1. Read the template's `exports:` frontmatter. For each declared export:
   - Look up the export's schema in `actions/interview-reference.md` (Export Schemas section — one per export kind and template).
   - Compose the artifact from the approved session entries. Pull content from `session.json`; do not invent.
   - Write the file to `./interview/<template>/exports/<path>`. Overwrite any prior export.

2. Append one synthesis line to CHANGELOG:
   ```
   ## <YYYY-MM-DD HH:MM> — exports written
   <list of filenames>
   ```

3. Report to the user:
   ```
   Exports written to ./interview/<template>/exports/:
     USER.md                        narrative profile
     SOUL.md                        decision framework
     HEARTBEAT.md                   checklist
     operating-model.json           full session dump
     schedule-recommendations.json  derived schedule

   Next: do work interview <template> ingest   to feed the operating model into BKB.
   ```

---

## Sub-Command: `<template> ingest`

Copies exports into `<repo-root>/kb/raw/inbox/` with BKB-compatible frontmatter.

### Preconditions

- `./interview/<template>/exports/` exists and is non-empty. If not, stop with: "No exports found. Run `do work interview <template> export` first."
- `<repo-root>/kb/` exists. If not, stop with: "No knowledge base found. Run `do work bkb init` first."

### Steps

1. For each file in `./interview/<template>/exports/`, write a companion file to `<repo-root>/kb/raw/inbox/` named `interview-<template>-<export-basename>.md` (e.g., `interview-work-operating-model-USER.md`).

2. Prepend YAML frontmatter per the Ingest Frontmatter section in `actions/interview-reference.md`:
   ```yaml
   ---
   title: <template-display-name> — <export-title>
   source: ./interview/<template>/exports/<export-filename>
   type: source-summary
   topic_cluster: <from template frontmatter>
   confidence: high
   created: <YYYY-MM-DD>
   ---
   ```
   `topic_cluster` is copied verbatim from the template's frontmatter.

3. For markdown exports, the body is the export content. For JSON exports, the body is a short pointer describing the shape and pointing at the source path — BKB ingests the descriptive summary, not raw JSON.

4. Append a per-layer markdown summary at `kb/raw/inbox/interview-<template>-layer-<layer-id>.md` for every layer in the session. Each file's body lists that layer's canonical entries in human-readable form. Same frontmatter shape as above.

5. Report:
   ```
   Ingested <N> files into kb/raw/inbox/:
     interview-<template>-USER.md
     interview-<template>-SOUL.md
     ...

   Next: do work bkb triage && do work bkb ingest
   ```

---

## Sub-Command: `<template> reset`

Archives the current run as a version and starts fresh. Destructive — requires confirmation.

### Steps

1. Verify `session.json` exists. If not, stop with the "no session" message.

2. **Require confirmation.** `$ARGUMENTS` must include `--confirm`, OR the user must respond affirmatively ("yes", "reset", "confirm") to an interactive prompt: "Reset will archive the current run as a version and start fresh. This cannot be undone. Proceed? (yes/no)".

3. Execute the `version` re-run mode from `actions/interview-reference.md`: archive current state as `versions/v<N>-<YYYY-MM-DD>/`, clear working `checkpoints/` and `exports/`, write a new empty `session.json` with `previous_version: "v<N>"`.

4. Append to CHANGELOG:
   ```
   ## <YYYY-MM-DD HH:MM> — reset (archived as v<N>)
   Fresh session started; v<N> retained for reference.
   ```

5. Report: "Reset complete. Archived as v<N>. Run `do work interview <template>` to start the new session."

---

## Sub-Command: `<template> versions`

Enumerates `./interview/<template>/versions/`.

### Output

```
Archived versions — <template>:

  v1-2026-03-12   5 layers   18 entries
  v2-2026-04-02   5 layers   21 entries
  v3-2026-04-16   5 layers   19 entries

Read an archive:  open ./interview/<template>/versions/<version-id>/
```

Counts come from each version's archived `session.json`. If `versions/` is empty, print "No archived versions yet."

---

## Output Format

Every sub-command returns terminal output (never writes silently). In-chat, the Interviewer also drafts checkpoints and waits for user approval before writing to `session.json`.

**What gets written:**

- `./interview/<template>/session.json` — authoritative session state (written only on explicit layer approval, review completion, or export).
- `./interview/<template>/checkpoints/<layer-id>.md` — transient approval drafts (overwritten on revision).
- `./interview/<template>/exports/<filename>` — export artifacts (overwritten on re-export).
- `./interview/<template>/versions/v<N>-<date>/` — immutable archives (written by `fresh`, `version`, `reset`).
- `./interview/<template>/CHANGELOG.md` — append-only activity log (one entry per approval, review, export, archive).
- `<repo-root>/kb/raw/inbox/interview-<template>-*.md` — BKB-ready files (written by `ingest`).

---

## Rules

- **Single instance per template per repo.** Context separation comes from installing the skill in multiple repos, not from profiles or workspaces within one repo.
- **Never persist content the user did not approve.** A checkpoint must be explicitly confirmed before entries move into `session.json`. Silence is not confirmation.
- **Never invent canonical fields.** If the user did not provide `constraints`, `inputs`, or `stakeholders`, ask. Do not default to empty lists without asking.
- **Templates are the contract.** Layer order, prompts, and export shapes come from the template file. The action does not improvise a layer or skip one.
- **Versions are immutable.** Once written to `versions/`, a directory is never edited by the action. The user's `previous_version` reference is the only back-link.
- **Exports gate on review.** `export` refuses to run unless all layers are approved and at least one review pass has completed. The gate exists to catch cross-layer tensions before they propagate into agent instructions.
- **Local files only.** No MCP dependencies. No external services. Session state, templates, and exports are plain files the user can diff, grep, and commit.

---

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "The user seemed to agree — I'll save this checkpoint without a clear 'yes'" | Ask explicitly. Wait for a concrete confirmation word or equivalent. | Soft agreement becomes "I never said that" three weeks later when the operating model feels wrong. The checkpoint is the contract. |
| "They didn't mention constraints for this entry, but I can infer one from context" | Ask. "Are there constraints on this — guardrails, approvals, limits?" | Inferred fields become load-bearing in exports. An invented constraint shows up in SOUL.md as a rule the user never wrote. |
| "The review pass surfaced a tension, but it's minor — I'll note it in my head and skip" | Present it. Let the user decide `both-are-true`, `revise`, or `skip` explicitly. | Tensions you skipped silently don't exist. Tensions the user saw and chose to accept are a real decision that shows up in the final artifact. |
| "Export before review — they're just markdown files, we can always regenerate" | Require review first. Refuse to export until `review_runs >= 1`. | Unreviewed exports embed contradictions into USER.md/SOUL.md. An agent consuming those artifacts will follow the contradiction. |
| "The user wants to start over — I'll just overwrite session.json" | Archive first via `fresh` or `reset`. Never destroy a prior run. | The prior run may contain entries the user wants to reference. Versions are cheap; regret is not. |
| "The template is missing a layer the user wants to add — I'll improvise" | Stop. Tell the user to edit the template file. Then resume. | The template is the contract between runs. Improvising breaks re-run comparison and makes exports inconsistent. |
| "The user is pushing for faster pace — I'll batch three questions" | One at a time. The pace pressure is a signal to ask sharper questions, not more of them. | Batched questions produce batched non-answers. The user picks one and the other two evaporate. |

---

## Red Flags

- A checkpoint file exists for a layer that isn't marked `approved: true` in `session.json`. Means an approval step was skipped or lost; re-present the checkpoint before advancing.
- An entry's `source_confidence: confirmed` but the `title` and `summary` are the Interviewer's phrasing, not the user's. Likely inflated — downgrade to `synthesized` or re-ask the user.
- Exports exist but `review_completed_at` is `null`. The gate was bypassed; delete exports and re-run review + export.
- A CHANGELOG entry claims a layer was approved but the layer's `entries: []` is empty in `session.json`. The approval wrote nothing — the session is inconsistent. Stop and surface to the user.
- `versions/` contains a directory with a newer date than `session.json.last_activity_at`. Someone edited a version folder manually; the invariant is broken. Flag and ask the user.
- Two entries in the same layer share an identical `title`. Duplicates bypass de-duplication during export and show up twice in USER.md. Merge or rename before approving.

---

## Verification Checklist

- [ ] Every layer declared in the template has `approved: true` before export runs.
- [ ] `review_completed_at != null` AND `review_runs >= 1` before export runs.
- [ ] Every entry in `session.json` has all 11 canonical fields (title, summary, cadence, trigger, inputs, stakeholders, constraints, details, source_confidence, status, last_validated_at).
- [ ] Every `source_confidence` value is `confirmed` or `synthesized` — no other strings.
- [ ] Every `status` value is `active`, `stale`, or `aspirational`.
- [ ] CHANGELOG has one `layer approved:` line per approved layer, in the order approvals occurred.
- [ ] Export files exist in `./interview/<template>/exports/` for every declared export, and their content matches the schema in `actions/interview-reference.md`.
- [ ] Ingest output lands in `kb/raw/inbox/` with filenames of the form `interview-<template>-<export-basename>.md`.
- [ ] Versions directories follow the `v<N>-<YYYY-MM-DD>/` naming convention and `<N>` is monotonically increasing.
- [ ] No checkpoint file was written to `session.json` without an explicit user approval recorded in the CHANGELOG.

---

## Error Handling

- **Template file missing** → list available templates, suggest `do work interview list`.
- **Session corrupt (invalid JSON)** → do not attempt repair. Tell the user the file path and stop: "`./interview/<template>/session.json` has invalid JSON. Inspect and fix manually, or archive and start fresh with `do work interview <template> reset`."
- **`export` invoked with unapproved layers** → list which layers are missing approval.
- **`export` invoked before review** → tell the user to run `review` first.
- **`ingest` invoked without completed exports** → tell the user to run `export` first.
- **`ingest` invoked without `kb/`** → tell the user to run `do work bkb init` first.
- **`reset` without confirmation** → require explicit `--confirm` flag or an interactive "yes" before archiving.
- **Checkpoint revision cycle exceeds 5 rounds on one layer** → pause and ask the user directly: "We've gone five rounds on this layer. Do you want to skip ahead, take a break, or keep refining?" Do not loop indefinitely.

---

## Help Menu

When invoked with no sub-command or with `help`:

```
do work interview — Run a structured elicitation interview

  Discover:
    do work interview list                    List available templates

  Run an interview:
    do work interview <template>              Start or resume an interview
    do work interview <template> status       Show session progress
    do work interview <template> review       Run the cross-layer contradiction pass
    do work interview <template> export       Write export artifacts
    do work interview <template> ingest       Feed exports into the knowledge base

  Session lifecycle:
    do work interview <template> versions     List archived runs
    do work interview <template> reset        Archive current run, start fresh (requires --confirm)

  Typical flow:
    1. do work interview list                            See what templates are available
    2. do work interview work-operating-model            Walk the five layers, ~45 minutes
    3. do work interview work-operating-model review     Resolve cross-layer tensions
    4. do work interview work-operating-model export     Produce USER.md, SOUL.md, HEARTBEAT.md + JSON
    5. do work interview work-operating-model ingest     Feed into BKB for querying
    6. do work bkb triage && do work bkb ingest          Compile into the knowledge wiki

  Re-run cadence:
    Quarterly, or after a major role/responsibility change. Pick:
      update   — revalidate prior entries in place
      version  — archive the old run and start a fresh session referencing it
      fresh    — archive and start completely over

  Each template lives at interviews/<name>.md and declares its own layers, prompts,
  and exports. See docs/interview-guide.md for onboarding and actions/interview-reference.md
  for the template authoring spec.
```
