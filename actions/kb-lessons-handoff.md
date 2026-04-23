# KB Lessons Handoff

> **Part of the do-work skill.** Reusable handoff instructions called by the work action (Step 7.5) and the review-work action (Step 9.5). Offers to promote a REQ's `## Lessons Learned` into the project's knowledge base by dropping a source file into `kb/raw/inbox/` for the next `bkb triage` + `bkb ingest` cycle.

This file is not a standalone action — it is loaded by other actions as a reference. If you reached this file directly, you probably want the review-work action instead.

## Philosophy

- **Zero external dependency.** The handoff writes into do-work's own KB system (see the bkb action). Nothing outside this skill is required.
- **One-way and terminal per REQ.** The handoff runs once per REQ, after lessons are captured. Downstream processing (triage, ingest, wiki compilation) is the bkb action's job, not this handoff's.
- **User pilots the drop.** do-work prepares a structured source document and asks before writing. That keeps the handoff consistent across harnesses — any agent that can read/write files can run it.
- **Graceful degradation.** If the project has no `kb/` directory yet, the handoff records the lessons as `pending` on the REQ and points the user at `do-work bkb init`. It never blocks archival.

## When to run this

Run after lesson capture completes in the calling action. Trigger conditions (all required):

- The REQ has a `## Lessons Learned` section with at least one non-placeholder bullet (real content, not the `[bullet]` template stub).
- Review completed and did not fail (no blocking `Important` findings). Soft findings are fine — the KB records *how the REQ landed*, not whether it was perfect.
- The REQ's `kb_status` frontmatter is either absent or set to `pending`. If it is already `promoted`, `declined`, or `skipped`, the handoff already ran — do not re-offer.

Skip silently (set `kb_status: skipped`) when:

- The Lessons Learned section is empty, placeholder, or trivially narrow ("fixed a typo").
- The REQ was a Route A change with no exploration, no failed approaches, and no gotchas — there is nothing to compound.

## Steps

### Step 1: Locate the KB

Search, in order:

1. `kb/` in the current working directory (the common case).
2. `knowledge-base/` in the current working directory.
3. Any ancestor directory (up to 3 levels up) containing both `raw/` and `wiki/` subdirectories.

If none is found, set `kb_status: pending` on the REQ, print:

```
No knowledge base found. The REQ's Lessons Learned are captured in the REQ but not promoted — kb_status: pending is recorded.
To promote: run `do-work bkb init`, then re-run this handoff (e.g. `do-work review REQ-NNN`) to drop the lesson into kb/raw/inbox/, then `do-work bkb triage` to sort it and `do-work bkb ingest` to compile it into the wiki.
```

…and return. Do not auto-init the KB — that is the user's call.

If the KB is found, note its root path as `<kb>` for the next steps.

### Step 2: Assemble the source document

Pull these fields from the REQ (file at `do-work/working/REQ-*.md` in pipeline mode, or `do-work/archive/...` in standalone mode):

| Field | Source in the REQ |
|---|---|
| `title` | REQ `title` frontmatter, trimmed |
| `date` | `completed_at` frontmatter if present (standalone mode on an archived REQ); otherwise today's date in `YYYY-MM-DD` (pipeline mode runs this handoff at Step 7.5, before Step 8 writes `completed_at`, so the fallback is required). Same calendar day either way. |
| `req_id` | REQ `id` frontmatter (e.g., `REQ-042`) |
| `req_path` | Absolute path to the REQ file, so the KB entry can back-reference it |
| `domain` | REQ `domain` frontmatter (e.g., `backend`, `frontend`, `security`) |
| `module` | If `prime_files` frontmatter lists one or more paths, use the shared parent directory. Otherwise the first path in the Implementation Summary's "Files changed" list |
| `what` | One-paragraph summary from the REQ's `## What` section |
| `what_worked` | Bullets from `## Lessons Learned` → "What worked" |
| `what_didnt_work` | Bullets from `## Lessons Learned` → "What didn't" |
| `worth_knowing` | Bullets from `## Lessons Learned` → "Worth knowing" |
| `solution_summary` | 2-3 sentence summary of `## Implementation Summary` — not the full diff |
| `tags` | 2-4 keywords drawn from `domain`, the REQ's title, and recurring nouns in the Lessons |

If a field has no reasonable source, omit it from the source document. Do not invent content to fill gaps.

Render the fields as a Markdown source document with YAML frontmatter — this is the shape bkb's Compiler agent expects. Template:

```markdown
---
source_type: req_lesson
req_id: REQ-042
req_path: do-work/archive/UR-005/REQ-042-auth-fix.md
date: 2026-04-23
domain: backend
module: src/auth
tags: [auth, jwt, session]
---

# Lessons from REQ-042: Session token validation fix

## What the REQ was about

<one-paragraph summary from ## What>

## Solution summary

<2-3 sentences from ## Implementation Summary>

## What worked

- <bullet>
- <bullet>

## What didn't work

- <bullet>

## Worth knowing

- <bullet>

## Back-reference

See `do-work/archive/UR-005/REQ-042-auth-fix.md` for the full REQ — plan, exploration, implementation, review, and lessons.
```

### Step 3: Present the handoff to the user

Print a three-block message:

1. **What:** one-line statement of what is being offered ("Promote this REQ's lessons into the knowledge base at `<kb>/raw/inbox/`?").
2. **Prepared source document:** a fenced code block containing the assembled Markdown from Step 2.
3. **Ask:** three options, phrased generically so the harness's ask-user prompt can surface them:
   - **a) Drop now** — write the file to `<kb>/raw/inbox/REQ-NNN-<slug>.md`. The user can run `do-work bkb triage` next to sort it into `raw/capture/notes/` and `do-work bkb ingest` to compile it into the wiki.
   - **b) Save for later** — record `kb_status: pending` and move on. Nothing is written to the KB. The user can re-run the handoff later.
   - **c) Skip** — record `kb_status: declined` (active user refusal, terminal) and do not offer again. `skipped` is reserved for the silent auto-skip path at the top of this file, which never prompts the user.

Use your environment's ask-user prompt. If the harness has no blocking prompt available and the action is running in unattended mode (pipeline without human in the loop), default to (b) `Save for later` — never auto-write to the KB without consent.

### Step 4: Execute the chosen path

**(a) Drop now:**

1. Compute the slug: lowercased REQ title, non-alphanumerics replaced with `-`, collapsed runs of `-`, trimmed to 40 chars max. Example: `REQ-042: Session token validation fix` → `session-token-validation-fix`.
2. Compose the filename: `<req_id>-<slug>.md`. Example: `REQ-042-session-token-validation-fix.md`.
3. If `<kb>/raw/inbox/<filename>` already exists, append a numeric suffix (`-2`, `-3`, …) until the path is free. Never overwrite an existing inbox file.
4. Write the source document (from Step 2) to that path.
5. Proceed to Step 5 with `status = promoted` and `entry = <filename>`.

**(b) Save for later:** skip to Step 5 with `status = pending` and `entry = null`.

**(c) Skip:** skip to Step 5 with `status = declined` and `entry = null`.

### Step 5: Update the REQ frontmatter

Append or update two frontmatter fields on the REQ (do not touch any other frontmatter):

```yaml
kb_status: <promoted | pending | declined | skipped>
kb_entry: <filename written to raw/inbox/, or empty if not promoted>
```

`kb_entry` is the *filename only*, not a full path, so it survives bkb's later moves from `inbox/` → `capture/notes/` → `processed/YYYY-MM-DD/`. You can always `find <kb> -name "<kb_entry>"` to locate it.

If the REQ predates this handoff and has neither field, add both. Leave legacy REQs alone — do not retroactively backfill on REQs you didn't just touch.

Use `declined` when the user actively refused; use `skipped` when you auto-skipped because the handoff conditions weren't met; use `pending` when the user chose "save for later" or no KB existed; use `promoted` only after you wrote the file to `raw/inbox/`.

### Step 6: Return control to the caller

Print a one-line confirmation:
- `Promoted to <kb>/raw/inbox/<filename>. Run \`do-work bkb triage\` to sort it, then \`do-work bkb ingest\` to compile it into the wiki.` for `promoted`
- `Handoff pending — run the handoff again later or drop the content manually into <kb>/raw/inbox/.` for `pending`
- `Skipped KB handoff.` for `skipped` or `declined`

Then return. The caller (work Step 7.5 or review-work Step 9.5) resumes its own flow.

## Rules

- Never auto-drop to `raw/inbox/` without user consent, even if the harness allows unattended tool calls. The KB is a persistent shared artifact; wrong entries cost more than missed ones.
- Never overwrite an existing file in `raw/inbox/`. Use a numeric suffix on collision.
- Never write anywhere in the KB other than `raw/inbox/`. Triage and ingest are the bkb action's responsibilities — this handoff stops at the drop.
- Never modify fields on the REQ other than `kb_status` and `kb_entry`.
- Never block the rest of the work action on this handoff. If the user ignores the prompt or the interaction times out, default to `pending` and continue.
- Never auto-init the KB. If no `kb/` exists, point the user at `do-work bkb init` and set `kb_status: pending`.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "No `kb/` exists, so I'll just create one and drop the file" | Set `kb_status: pending` and point the user at `do-work bkb init` | Silently initializing a KB the user didn't ask for pollutes their project. `bkb init` is a deliberate step. |
| "The Lessons section is empty but I'll invent content from the diff" | Set `kb_status: skipped` and skip | Fabricated lessons pollute the KB. Compound value comes from real friction, not filler. |
| "Let me also run triage and ingest for the user" | Stop at the drop | `bkb triage` and `bkb ingest` are their own actions with their own agent crews. Chaining them from this handoff hides important Reviewer confidence checks. |
| "The REQ already has `kb_status: declined` but this time the lessons are better" | Leave it alone | User already chose. They can drop a file into `raw/inbox/` manually if they change their mind. |
| "Multiple REQs this session all promoted — let me batch them into one KB entry" | One file per REQ | The KB's back-references and confidence tracking depend on one source per file. |

## Red Flags

- The handoff ran but `kb_status` is missing from the REQ — Step 5 was skipped.
- `kb_status: promoted` is set but no file matching `kb_entry` exists under `<kb>/raw/` — the write silently failed or the filename was captured wrong.
- The source document has empty `what_worked` + `what_didnt_work` + `worth_knowing` sections — the extraction failed; do not drop the file, surface the gap to the user first.
- Multiple REQs in the same session all show `kb_status: pending` and the user clearly intended to promote — remind the user they can re-run the handoff per REQ, or drop files directly into `raw/inbox/`.
- `raw/inbox/` fills up faster than `bkb triage` drains it — the pipeline is stalled upstream, not a handoff problem, but surface the backlog in the final confirmation message.

## Verification Checklist

- [ ] KB location resolved (or `kb_status: pending` set with a pointer to `bkb init`).
- [ ] Source document assembled with real values in `title`, `what`, `solution_summary`, plus at least one Lessons bullet.
- [ ] User was asked before writing; unattended pipeline defaulted to `pending`.
- [ ] REQ frontmatter now has `kb_status` set to exactly one of `promoted | pending | declined | skipped`.
- [ ] If `promoted`, `kb_entry` points to a file that exists under `<kb>/raw/inbox/` (or has already moved to `capture/` or `processed/`).
- [ ] The calling action's flow resumed — this handoff did not early-return the parent workflow.
