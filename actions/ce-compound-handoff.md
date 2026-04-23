# CE Compound Handoff

> **Part of the do-work skill.** Reusable handoff instructions called by the work action (Step 7.5) and the review-work action (Step 9.5). Offers to promote a REQ's `## Lessons Learned` into the compound-engineering knowledge base at `docs/solutions/`.

This file is not a standalone action — it is loaded by other actions as a reference. If you reached this file directly, you probably want the review-work action instead.

## Philosophy

- **Augmentation, not dependency.** do-work completes the REQ cycle by itself. If compound-engineering (CE) is not installed, this handoff prints a prepared prompt the user can run later. It never blocks and never errors.
- **One-way and terminal.** The handoff runs once per REQ, after lessons are captured. Promoted solutions become CE's problem, not do-work's. No sync, no back-reference maintenance beyond the REQ's `ce_solution_path`.
- **User pilots the dispatch.** do-work prepares structured input. The user decides whether to run `/ce-compound`, save the prompt for later, or skip. That keeps the integration harness-agnostic: any agent that understands "run this slash command if you have it" can handle this.

## When to run this

Run after lesson capture completes in the calling action. Trigger conditions (all required):

- The REQ has a `## Lessons Learned` section with at least one non-placeholder bullet (real content, not the `[bullet]` template stub).
- Review completed and did not fail (no blocking `Important` findings). Soft findings are fine — CE captures *how the REQ landed*, not whether it was perfect.
- The REQ's `ce_compound_status` frontmatter is either absent or set to `pending`. If it is already `promoted`, `declined`, or `skipped`, the handoff already ran — do not re-offer.

Skip silently (set `ce_compound_status: skipped`) when:

- The Lessons Learned section is empty, placeholder, or trivially narrow ("fixed a typo").
- The REQ was a Route A change with no exploration, no failed approaches, and no gotchas — there is nothing to compound.

## Steps

### Step 1: Assemble the handoff payload

Pull these fields from the REQ (file at `do-work/working/REQ-*.md` in pipeline mode, or `do-work/archive/...` in standalone mode):

| Payload field | Source in the REQ |
|---|---|
| `title` | REQ `title` frontmatter, trimmed |
| `date` | ISO date from `completed_at` frontmatter (YYYY-MM-DD) |
| `category` | Infer from `domain` frontmatter: `backend` → `backend`, `frontend` → `frontend`, `security` → `security`, `performance` → `performance`, `ui-design` → `frontend`, `testing` → `testing`, else → `skill-design` for do-work internal work or the REQ's own best category label |
| `module` | If `prime_files` frontmatter lists one or more paths, use the shared parent directory. Otherwise the first path in the Implementation Summary's "Files changed" list |
| `problem_type` | Derive from the REQ shape: bug-style REQ (review caught issues, test failures, production errors) → `bug`. Feature-style REQ with Lessons about friction or gaps → `workflow_issue`, `developer_experience`, or `best_practice`. See CE's `schema.yaml` for the full enum if the user has CE installed |
| `symptoms` | Bullets from the REQ's `## What` section describing observable behavior |
| `what_didnt_work` | Bullets from `## Lessons Learned` → "What didn't" |
| `solution` | Summary of `## Implementation Summary` — 2-3 sentences, not the full diff |
| `why_this_works` | Inferred from the REQ's root-cause narrative (Plan section for Route C, or the Lessons "Worth knowing" bullet) |
| `prevention` | Bullets from Lessons Learned → "Worth knowing" + "What didn't" reframed as guardrails |
| `tags` | 2-4 keywords drawn from `domain`, the REQ's title, and recurring nouns in the Lessons |

If a field has no reasonable source, leave it as `[TBD]` in the payload — CE's compound skill will prompt for missing fields if the user runs it.

### Step 2: Present the handoff to the user

Print a three-block message:

1. **What:** one-line statement of what is being offered ("Promote this REQ's lessons into docs/solutions/ via compound-engineering?").
2. **Prepared payload:** a fenced code block containing the assembled payload from Step 1, formatted as YAML keys (so a user can eyeball it quickly).
3. **Ask:** three options, phrased generically so the harness's ask-user prompt can surface them:
   - **a) Run `/ce-compound` now** — you will invoke the CE compound skill with the payload above. Only offer this if your harness recognizes that slash command or equivalent skill invocation; otherwise surface it as "copy the prompt below and run it in a CE-enabled session."
   - **b) Save for later** — record `ce_compound_status: pending` and move on. The user can run `/ce-compound REQ-NNN` later.
   - **c) Skip** — record `ce_compound_status: skipped` and do not offer again.

Use your environment's ask-user prompt. If the harness has no blocking prompt available and the action is running in unattended mode (pipeline without human in the loop), default to (b) `Save for later` — never auto-dispatch.

### Step 3: Execute the chosen path

**(a) Run now:** Invoke the CE compound skill with the payload. The canonical dispatch is the slash command `/ce-compound` (Claude Code, Cursor, Codex with CE installed). If your harness requires a different form, follow the CE skill's invocation contract — the skill's `SKILL.md` is at `.claude/skills/ce-compound/SKILL.md`, `~/.codex/skills/compound-engineering/ce-compound/SKILL.md`, or wherever your plugin installer placed it. After CE returns, capture the path of the file it created under `docs/solutions/...` and proceed to Step 4.

**(b) Save for later:** Skip to Step 4 with `status = pending` and `solution_path = null`.

**(c) Skip:** Skip to Step 4 with `status = skipped` and `solution_path = null`.

### Step 4: Update the REQ frontmatter

Append or update two frontmatter fields on the REQ (do not touch any other frontmatter):

```yaml
ce_compound_status: <promoted | pending | declined | skipped>
ce_solution_path: <docs/solutions/category/slug.md or empty if not promoted>
```

If the REQ predates CE integration and has neither field, add both. Leave legacy REQs without these fields alone — don't retroactively backfill on REQs you didn't just touch.

Use `declined` when the user actively refused; use `skipped` when you auto-skipped because the handoff conditions weren't met; use `pending` when the user chose "save for later"; use `promoted` only after CE wrote a solution file and you captured the path.

### Step 5: Return control to the caller

Print a one-line confirmation:
- `Promoted to <solution_path>.` for `promoted`
- `Handoff pending — run do-work ce-compound REQ-NNN or /ce-compound later.` for `pending`
- `Skipped compound handoff.` for `skipped` or `declined`

Then return. The caller (work Step 7.5 or review-work Step 9.5) resumes its own flow.

## Rules

- Never auto-dispatch `/ce-compound` without user consent, even if the harness allows unattended tool calls. The compound knowledge base is a persistent shared artifact; wrong entries cost more than missed ones.
- Never retry on CE invocation failure. If CE errors, fall back to `pending` and surface the error to the user — do not loop.
- Never write to `docs/solutions/` yourself. The CE compound skill owns that path. This handoff only reads back the path CE wrote.
- Never modify fields on the REQ other than `ce_compound_status` and `ce_solution_path`.
- Never block the rest of the work action on this handoff. If the user ignores the prompt or the interaction times out, default to `pending` and continue.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "CE isn't installed, so I'll just write to `docs/solutions/` myself" | Print the prepared payload as text, set `ce_compound_status: pending`, move on | The CE schema has invariants this handoff doesn't enforce — bad entries pollute the KB |
| "The Lessons section is empty but I'll invent content from the diff" | Set `ce_compound_status: skipped` and skip | Fabricated lessons are worse than none; compound value comes from real friction |
| "Let me auto-dispatch `/ce-compound` in pipeline mode since there's no human around" | Default to `pending` in unattended mode | Unattended promotion means no one reviews the category, tags, or root-cause framing |
| "The REQ already has `ce_compound_status: declined` but this time the lessons are better" | Leave it alone | User already chose. If they want to revisit, they can invoke `/ce-compound REQ-NNN` directly |

## Red Flags

- The handoff ran but `ce_compound_status` is missing from the REQ — Step 4 was skipped.
- `ce_compound_status: promoted` is set but no file exists at `ce_solution_path` — CE's write silently failed or the path was captured wrong.
- The payload has `[TBD]` in `title`, `category`, or `solution` — upstream extraction failed; do not dispatch with those stubs, surface the gap to the user first.
- Multiple REQs from the same session all show `ce_compound_status: pending` and the user clearly intended to promote — remind the user they can batch-dispatch later.

## Verification Checklist

- [ ] Payload block printed with real values in `title`, `category`, `symptoms`, `solution`, `prevention` (no `[TBD]` placeholders in core fields).
- [ ] User was asked before `/ce-compound` dispatched; unattended pipeline defaulted to `pending`.
- [ ] REQ frontmatter now has `ce_compound_status` set to exactly one of `promoted | pending | declined | skipped`.
- [ ] If `promoted`, `ce_solution_path` points to a file that exists under `docs/solutions/`.
- [ ] The calling action's flow resumed — this handoff did not early-return the parent workflow.
