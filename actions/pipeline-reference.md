# Pipeline Reference

> **Companion to the pipeline action.** Contains the three Pipeline Completion Report rendering templates — Plain Markdown, Marp Slide Deck, and Standalone HTML Debrief — plus the composition rules that apply to all three authored renderings. A fourth file, `{UR-NNN}-pipeline-summary.marp.html`, is produced mechanically by running `marp-cli` on the Marp source — no template, no composition rules. Extracted from `pipeline.md` to keep the main action file readable; load this file at Step 5 Completion when rendering the report.

## Pipeline Completion Report — three renderings of one dataset (plus one mechanical export)

The same facts — Final summary, Test state, Coherence, Carry-forward, Deliverables, How to verify — are rendered three ways by the LLM (`.md`, `.marp.md`, `.single.html`). One authoring pass over the data, three files on disk. A fourth file — `.marp.html` — is produced by running `npx @marp-team/marp-cli {UR-NNN}-pipeline-summary.marp.md --html` after the Marp source is written; it inherits its content mechanically and needs no template. Never author any of the three LLM renderings from scratch if another already exists; re-render from the source data so they stay consistent.

### Composition rules (apply to all three formats)

- **Serve both audiences in every file.** Each summary opens with a "What got built" narrative for the reader who has no clue, then transitions into the audit data for the reader who wants receipts. A stakeholder landing on the `.html` should understand the feature without opening any other file; a developer scanning the `.md` should reach the commits and test deltas within seconds. Never ship a summary that only audits or only educates.
- **Reuse client-brief content verbatim.** When the `present` step ran, the "What got built" narrative and architecture diagram come from `{UR-NNN}-client-brief.md` — copy the same sentences and the same diagram. Paraphrasing across files introduces drift. If the brief doesn't exist (present skipped or produced nothing), synthesize from the REQ Implementation Summaries.
- **Cite commits, not prose.** Every audit claim should trace to a commit SHA, a REQ ID, or a file path. Tables and bullet lists with pointers beat paragraphs of explanation. (The opening narrative is the exception — it's plain language for the no-clue reader.)
- **Pull from primary sources.** Final summary rows come from REQ frontmatter; coherence notes come from the review step's actual output; test deltas come from what `run` and `review` logged. Do not invent metrics.
- **Be honest about gaps.** If the baseline test count wasn't captured before the pipeline started, write "baseline not recorded" — don't guess. If no cross-REQ coherence was analyzed (single-REQ pipeline), omit that section.
- **Carry-forward ≠ auto-capture.** List candidates clearly with the command the user would run to capture each one, but never capture them automatically.
- **No format-specific editorializing.** The Marp deck must not add facts the markdown lacks; the HTML must not soften or strengthen claims for a broader audience. Format dictates rendering; rendering does not dictate facts.

## 1. Plain Markdown Report — `{UR-NNN}-pipeline-summary.md`

Developer-facing. Read in a terminal with `cat`, grepped, or pasted into a PR description. No YAML header, no CSS, no slide breaks — just markdown.

```markdown
# Pipeline Completion Report — {UR-NNN}

**Session**: {session_id} · **Duration**: {duration} · **Branch**: {branch} ({pushed|local})
**Verdict**: {PASS | PASS with caveats | FAIL}

## What got built (for the reader who has no clue)

[2-3 plain-language sentences synthesizing the UR title, the REQs' What sections, and the client brief's "What We Built" paragraph — no jargon, no commit SHAs, no REQ IDs. The reader learns what the feature *does* before they see the audit trail. If the `present` step ran, pull this straight from `{UR-NNN}-client-brief.md` so the two files stay in sync; if it didn't run, synthesize from the REQ Implementation Summaries.]

[Optional: reuse the ASCII architecture diagram from the client brief, verbatim. Skip if the work is non-architectural (config tweak, bug fix, docs).]

**Go deeper:** [`{UR-NNN}-client-brief.md`](./{UR-NNN}-client-brief.md) · [`{UR-NNN}-interactive-explainer.single.html`](./{UR-NNN}-interactive-explainer.single.html) *(only include links that actually exist on disk)*

## Final summary

| REQ | Commit | Scope | One-line |
|-----|--------|-------|----------|
| REQ-402 | 5ab214d | docs     | 4 lessons-learned files + prime links |
| REQ-410 | 9371a68 | refactor | shared `initializeDatabaseAtPath` — prod/test converged |
| REQ-413 | 9e20bde | backend  | SHA-256 index + O(log N) lookup rewrite |
| ...     | ...     | ...      | ... |

## Test state (before → after the {N}-REQ pipeline)

| Suite         | Before    | After     | Delta |
|---------------|-----------|-----------|-------|
| Go (sa1-server) | 81 tests  | 98 tests  | +17 |
| Frontend      | 1053 tests / 62 suites | 1067 tests / 65 suites | +14 tests / +3 suites |
| `go vet`      | clean     | clean     | — |

## Cross-REQ coherence highlights (verified by the review)

- **REQ-413 ↔ REQ-406**: early-exit preserved at cache-hit + fresh-match. `effectiveLimit` threaded.
- **REQ-413 ↔ REQ-407**: metric_version filter preserved in `loadCachedEdgesForSource`.
- **REQ-411 ↔ REQ-412**: orthogonal, zero file overlap, no shared state.

## Carry-forward work (implied, not captured yet)

- [Deferred item] — capture with `do work capture request: ...`
- [TODO/FIXME introduced and left for a follow-up]
- [`pending-answers` REQs awaiting user input — run `do work clarify`]

## Deliverables

Render each bullet as a relative markdown link to the file (e.g. `[...]({UR-NNN}-client-brief.md)`) so a reader opening the `.md` in GitHub, a PR, or an editor can click through to any sibling artifact. Group by audience so the reader lands on the right surface first.

**For the clueless-reader (start here if you don't know what was built):**

- [`{UR-NNN}-client-brief.md`](./{UR-NNN}-client-brief.md) — plain-language brief with architecture diagram + value prop *(if present ran)*
- [`{UR-NNN}-interactive-explainer.single.html`](./{UR-NNN}-interactive-explainer.single.html) — interactive Before/After explainer, open in any browser *(if present ran)*
- [`{UR-NNN}-video/`](./{UR-NNN}-video/) — Remotion video walkthrough (`cd` in, `npm install`, `npm run preview`) *(if present ran)*

**For the developer / reviewer (audit the run):**

- [`{UR-NNN}-pipeline-summary.md`](./{UR-NNN}-pipeline-summary.md) — this report (markdown)
- [`{UR-NNN}-pipeline-summary.marp.md`](./{UR-NNN}-pipeline-summary.marp.md) — Marp slide source (`marp --preview`)
- [`{UR-NNN}-pipeline-summary.marp.html`](./{UR-NNN}-pipeline-summary.marp.html) — Marp deck exported to HTML (for stakeholders without marp-cli)
- [`{UR-NNN}-pipeline-summary.single.html`](./{UR-NNN}-pipeline-summary.single.html) — standalone authored HTML debrief

## How to verify

1. **Check out the branch and pull latest:**
   ```
   git checkout {branch} && git pull
   ```
2. **Inspect each commit** (ordered to show the build-up):
   ```
   git show 5ab214d   # REQ-402 — lessons-learned docs
   git show 9371a68   # REQ-410 — shared init routine
   git show 9e20bde   # REQ-413 — SHA-256 index rewrite
   ```
3. **Run the tests** (matches what the pipeline ran):
   ```
   {project test command — e.g., `go test ./...` and `npm test`}
   ```
4. **Preview the other renderings:**
   ```
   npx @marp-team/marp-cli {UR-NNN}-pipeline-summary.marp.md --preview
   open do-work/deliverables/{UR-NNN}-pipeline-summary.marp.html
   open do-work/deliverables/{UR-NNN}-pipeline-summary.single.html
   ```
5. **Read the per-REQ archive** for the full trail of intent:
   ```
   do-work/archive/{UR-NNN}/REQ-*.md
   ```
```

## 2. Marp Slide Deck — `{UR-NNN}-pipeline-summary.marp.md`

Stakeholder-facing. Viewed with `marp --preview`, and also exported to `{UR-NNN}-pipeline-summary.marp.html` via `npx @marp-team/marp-cli {UR-NNN}-pipeline-summary.marp.md --html` so stakeholders without marp-cli can view the deck by opening a file. Must start with Marp YAML frontmatter (`marp: true`). Each slide separated by `---`. Keep slides scannable — no slide should fit more than ~8 rows of content; split long Final-summary tables across domain-grouped slides. Use a Mermaid `graph LR` on the coherence slide when there are 2+ cross-REQ links.

Required slide sequence (omit a slide entirely if its section has no data — don't leave empty slides):

1. **Title slide** — UR-NNN, session ID, branch, verdict badge
2. **What got built** — 2-3 plain-language bullets pulled from the client brief's "What We Built" section. No commit SHAs, no REQ IDs. This is the slide a stakeholder who wandered in late needs to orient. Skip only if no `present` step produced a brief AND the UR itself is trivially self-explanatory from its title.
3. **How it works** (when a client brief exists with an architecture diagram) — reuse the ASCII or Mermaid diagram from the brief. Skip for non-architectural changes.
4. **At-a-glance stats** — REQ count, commit count, test delta, duration (big numbers in a 2×2 or 4-column grid)
5. **What shipped — {domain}** — one slide per domain bucket (docs / backend / refactor / frontend / tests). Each is a table of REQ / commit / one-line for that domain only.
6. **Test state (before → after)** — the table, full-width
7. **Cross-REQ coherence** — Mermaid `graph LR` diagram of interacting REQs (skip for single-REQ pipelines)
8. **Coherence assertions** — verbatim review quotes, one bullet per assertion
9. **Carry-forward work** — bullets with capture commands (skip if none)
10. **How to verify** — fenced `bash` block with checkout + git-show + test commands
11. **Deliverables + next steps** — two-column layout: left column "Start here if you want to understand what was built" lists the client brief, interactive explainer, and video (when present ran); right column "Audit the run" lists the markdown and HTML summary siblings. Render each as the bare filename — stakeholders open the deck from `do-work/deliverables/`, so relative filenames are all they need to find the sibling files in the same folder.

Use this Marp frontmatter skeleton and extend the `style:` block as needed — don't invent new themes:

```yaml
---
marp: true
theme: default
paginate: true
size: 16:9
header: '{UR-NNN} — Pipeline Debrief'
footer: 'Session {session_id} · branch {branch}'
style: |
  section { font-family: system-ui, -apple-system, sans-serif; }
  h2 { color: #1e40af; border-bottom: 2px solid #e2e8f0; padding-bottom: 0.25em; }
  code { background: #f1f5f9; padding: 0.1em 0.3em; border-radius: 3px; }
  table { font-size: 0.75em; }
  th { background: #1e40af; color: white; }
  .big { font-size: 3em; font-weight: 700; color: #1e40af; }
  .label { color: #64748b; font-size: 0.9em; }
---
```

## 3. Standalone HTML Debrief — `{UR-NNN}-pipeline-summary.single.html`

Non-technical-reader-facing. Single `.html` file, zero build steps. Same content as the markdown, rendered for a browser. The `.single.` infix marks this as an LLM-authored standalone page (Tailwind + Mermaid via CDN, cross-links to siblings) — distinct from the `.marp.html` mechanical export of the Marp deck.

**Stack (CDN only — no npm, no build):**
- Tailwind CSS via `<script src="https://cdn.tailwindcss.com"></script>`
- Mermaid via `<script type="module">` import of `mermaid@10` from jsDelivr
- Vanilla JS only (no React, no framework)

**Required sections (in order):**

1. **Hero** — UR-NNN as H1, one-paragraph description, metadata badges (branch, duration, verdict)
2. **What got built** — a prose block that educates a reader who has no clue: 2-3 sentences explaining what the feature does in plain language, pulled from the client brief's "What We Built" section. No REQ IDs, no commit SHAs. When a client brief exists, this section is the primary educational entry point for a non-technical reader arriving at the HTML. Skip only if no `present` step ran and the UR is self-explanatory from its title.
3. **How it works** (when an architecture diagram exists in the client brief) — a `<div class="mermaid">` with a `graph TD` or `graph LR` rendering of the same components/data flow from the brief's architecture section. Caption each node in plain language. Skip entirely for non-architectural changes (config tweaks, bug fixes).
4. **At-a-glance stat cards** — 4-column grid of big-number stats (REQ count, commits, tests added, suites added)
5. **What shipped** — grouped sections by domain, each with a styled table of REQ / commit / one-line
6. **Test state** — the table, styled with the accent colour for the After column and green for the Delta
7. **How the work holds together** — a `<div class="mermaid">` containing the same `graph LR` from the Marp deck (Mermaid renders on load)
8. **Coherence assertions** — responsive card grid, one card per assertion, with the REQ pair in mono accent and the claim below (skip the whole section for single-REQ pipelines)
9. **Carry-forward work** — cards with a bold title, muted explanation, and the capture command in a `<pre>` block (skip if none)
10. **How to verify** — numbered headings, each followed by a copy-pasteable `<pre><code>` block
11. **Related deliverables** — a navigation card grid **before** the final follow-ups list, splitting cross-links by audience. Left card group "Understand what was built" with real `<a href="./{UR-NNN}-client-brief.md">` / `<a href="./{UR-NNN}-interactive-explainer.single.html">` / `<a href="./{UR-NNN}-video/">` anchors (only include tiles for artifacts that actually exist on disk — if present ran and produced them). Right card group "Audit the run" linking the markdown (`<a href="./{UR-NNN}-pipeline-summary.md">`), Marp source (`<a href="./{UR-NNN}-pipeline-summary.marp.md">`), and Marp HTML export (`<a href="./{UR-NNN}-pipeline-summary.marp.html">`) siblings. The `.single.html` is the most discoverable surface for a non-technical reader — it must point them to the deeper, more educational artifacts.
12. **Footer / next steps** — ordered list with `do work present {UR-NNN}` and other follow-ups

**Design requirements:**

- Light theme default; dark theme via `@media (prefers-color-scheme: dark)` overriding CSS custom properties on `:root`
- Palette: CSS variables for `--bg`, `--surface`, `--text`, `--muted`, `--accent`, `--accent-soft`, `--border`. Light: white/slate-50 / slate-900 / blue-600. Dark: slate-900 / slate-100 / blue-400.
- Font: `system-ui, -apple-system, sans-serif`
- Max content width: `max-w-6xl` centred
- Generous spacing (`py-10` / `py-16` on sections) — readable like a long-form article, not cramped like a dashboard
- Mermaid init: `mermaid.initialize({ startOnLoad: true, theme: 'default', securityLevel: 'loose' })`

**What NOT to do:**

- Don't add charts the source data doesn't support (no fabricated time-series, no fake percentages)
- Don't embed images unless the REQs reference them
- Don't pull in additional CDN scripts beyond Tailwind + Mermaid — the file must work offline once cached
- Don't add interactivity that hides data (collapsible sections are fine; JS-gated sections that require a click to reveal facts are not)
