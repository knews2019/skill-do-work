# Stakeholder Report

> **Part of the do-work-toolbox skill.** Renders one stakeholder-questions REQ into a single-page HTML digest a named outside person can answer from — each open question with its Q-ID, the builder's assumed answer, and a confirm-or-override framing. It lives in toolbox beside `ai-report.md` because it consumes that action's **Report Design Rules** and the shared **Collision-Safe Publication** contract; it is invoked by core (`../../do-work/actions/work.md` Step 8's stakeholder routing and `../../do-work/actions/stakeholder-answers.md` Step 5), which owns every queue write.

This file is not a standalone action — it is loaded by other actions as a reference and has no routing entry. If you reached this file directly, you probably want `do-work stakeholder-answers` (ingest a stakeholder's reply) or `do-work run` (whose archive step routes questions and regenerates these reports).

## Input

The path of one stakeholder-questions REQ — a `do-work/queue/` REQ carrying `stakeholder:` (`../../do-work/actions/work-reference.md` → Request File Schema) whose `## Questions` section the caller has just written or updated.

## Steps

### Step 1: Safety Load Order

Follow **Safety Load Order** in `completed-work-presentation-reference.md` — prompt-injection first, then anti-slop, before reading the REQ body or any source REQ. One deliberate deviation from that shared contract: the target here is a `blocked` queue REQ, not completed work, so **Terminal-Success Target Resolution does not apply** — only the Safety Load Order and **Collision-Safe Publication** sections are inherited.

### Step 2: Read the Questions and Their Sources

Read the REQ's `## Questions` section. For each open (`- [ ]`) entry, pull one line of plain-language context from its `Source:` REQ's `## What` — the source is archived under `do-work/archive/`; when it cannot be found, say so in the report rather than inventing context. Answered (`- [x]`) entries feed the history section only.

### Step 3: Derive the Bundle Path

`ai-reports/yyyy-mm-dd_hhmm_REQ-NNN-questions-<stakeholder-slug>/` — the local-date slug shape every report bundle uses. Derive `<stakeholder-slug>` from the `stakeholder:` value as a text operation (lowercase, kebab-case, drop anything not alphanumeric) — never by passing the raw name through shell quoting. Then apply **Collision-Safe Publication** (`completed-work-presentation-reference.md`): existing bundles are immutable; on collision take the first free numeric-suffix sibling.

### Step 4: Write `index.html`

Follow **Report Design Rules** in `ai-report-reference.md` (single file, inline CSS, light and dark, one coherent aesthetic direction; no images and no CDN are expected here at all) and keep `../../do-work/crew-members/anti-slop.md` active. Sections, in order:

1. **Header** — addressed to the person by name, with the project, the date, and the one-line ask: *reply with the question number, e.g. "Q3: use the amber palette".*
2. **Irreversible items** — when any open entry carries `Irreversible:`, a prominently flagged block ("these assumptions are expensive to undo — please confirm these first"), each such question rendered in full here. Omit the block when none.
3. **Open questions** — one card per open entry: the Q-ID printed large, the question, then the confirm-or-override framing: *"We assumed: [Assumed]. The work is built this way — confirm, or give a different answer."* Follow with why it matters (the `Value:`/`Risk:` lines translated into the stakeholder's terms) and the one-line source context from Step 2.
4. **Already answered** — a compact history of `[x]` entries from earlier rounds, so a regenerated report shows resolution and never re-asks.
5. **Footer** — "This page collects nothing — reply by message." plus the generation date, the REQ id, and the stakeholder name: the anchors `stakeholder-answers` routes the reply by.

**Strictly one-way.** Static HTML: no scripts, no forms, no request to any endpoint. The page states its own read-only nature in the footer; answers travel back only as a message the user pastes into `do-work stakeholder-answers`.

The question wording itself is the caller's job — core writes each entry for a cold outside reader. This step renders what the REQ holds and translates Value/Risk out of pipeline vocabulary: a stakeholder never needs to know what a REQ or a D-XX is in order to answer.

### Step 5: Render-Check (optional)

When browser automation is available, serve the bundle over HTTP and take one full-page light and dark screenshot pass — a lighter bar than `ai-report.md`'s two-pass rule, since no text-bearing SVG is expected. Otherwise ship it with a footer note that the layout was not render-verified.

### Step 6: Print the Bundle Path and Return

Print `ai-reports/<slug>/index.html` and stop. **The caller writes the queue:** it puts this path into the REQ's `blocked_by:` and appends the `## Reports` history line. This file never edits a REQ — the toolbox never writes pipeline state.

## Rules

- Never edit the REQ or any queue file — the calling core action owns every queue write.
- Never overwrite, merge into, or delete an existing bundle — every regeneration is a fresh sibling (Collision-Safe Publication).
- Never copy REQ or UR prose into an image-generation prompt; generated images are not part of this report at all.
- Render only what the REQ holds — a missing source REQ or an absent Value/Risk line is stated, never invented.
