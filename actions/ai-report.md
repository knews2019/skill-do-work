# AI Report Action

> **Part of the do-work skill.** Generates an HTML report of a completed feature — live screenshots with SVG callout annotations, before/after toggles, and Mermaid/SVG diagrams as a fallback when screenshots aren't available. Output is one HTML file plus a sibling `.assets/` folder of image binaries, both in `ai-reports/` in the project root.

The report exists to make a UI change **visible**: a stakeholder opens one HTML file and sees the literal pixels that changed, where they changed, and how to verify it themselves. Not a brief. Not a debrief. A pixel-anchored proof-of-work artifact.

## Philosophy

- **Pixels first, prose second.** The image is the conclusion. Text exists to explain what the eye is already seeing.
- **Anti-slop or it doesn't ship.** Every section passes the seven principles in `crew-members/anti-slop.md` (loaded in Step 1) — lead with the conclusion, verify every claim, compress, match medium to stakes.
- **HTML + sibling assets.** One `.html` file plus a `<report-stem>.assets/` folder of image binaries next to it. The HTML references the images by relative `src` — move the pair together and the report works anywhere. Tailwind + Mermaid load from a CDN.
- **Graceful when bowser is missing.** Live screenshots are nice-to-have. If `playwright-cli` isn't installed or no dev server responds, fall back to SVG architecture + Mermaid data-flow diagrams. The report still ships.

## When to Use

**Use when:**
- A UR/REQ is `status: completed` and the user wants a stakeholder-visible artifact showing *what changed visually*.
- A feature touches the UI and "show me the change" beats "describe the change."
- You have before/after assets (in `do-work/user-requests/UR-NNN/assets/`, `do-work/working/`, or root `verify-*.png`) and want a side-by-side comparison.

**Do NOT use when:**
- The work has no user-visible output (infra-only, refactor, tooling) — the report is empty by construction; use the present-work brief instead.
- The user wants an *educational* explainer (architecture + value prop + data flow) — use present-work, which writes a `.single.html` Interactive Explainer to `do-work/deliverables/`.
- The user wants a multi-REQ developer/PM **debrief** of a pipeline run (test deltas, REQ coherence graph, carry-forward work) — use the pipeline action's completion report.
- Work is still in progress — there's nothing shipped to report on yet.

## Input

`$ARGUMENTS` — one of:

| Value | Behaviour |
|-------|-----------|
| `UR-NNN` | Report covers all completed REQs under that UR |
| `REQ-NNN` | Report covers that single REQ |
| `most recent` or blank | Find the most recently completed UR in `do-work/archive/` |

## Steps

### Step 1: Load Principles

Read `crew-members/anti-slop.md`. Keep all seven principles active for every section you write below. Do not run `do-work slop-check` as a separate step — internalize and apply inline.

### Step 2: Resolve the Target

1. If a UR/REQ was specified, locate it in `do-work/archive/`.
2. If blank or "most recent": scan `do-work/archive/UR-*/` for the highest UR number whose folder contains at least one REQ file with `status: completed`.
3. For a UR target, collect all completed REQs under it. For a REQ target, collect just that file.
4. If nothing is found with `status: completed`, stop and tell the user there's nothing to report on.

Extract from each REQ:

- `title`
- `commit` (SHA from frontmatter)
- `## What` / `## Detailed Requirements` — what was requested
- `## Implementation Summary` — files changed, what was done
- `## Testing` — how to verify
- `## Review` — scores, acceptance
- `## Lessons Learned` — surprising insights (if present; Route A REQs may skip)

Also read the parent UR's `input.md` for the user's own words.

### Step 3: Collect Visual Evidence

**3a: Check for stored before/after assets.** Look in these locations, in order:

1. `do-work/archive/UR-NNN/assets/` — archived user-supplied screenshots (the common case; completed URs live here after cleanup)
2. `do-work/user-requests/UR-NNN/assets/` — live UR assets (target not yet archived)
3. `do-work/working/` — screenshots taken during development (match by UR/REQ prefix or date proximity to commit)
4. Project root — any `verify-*.png` files captured during review
5. Git diff images: `git show <commit> --name-only | grep -E '\.(png|jpg|gif)$'`

Classify found images:

- **Before** — filename contains `before`, `old`, `prior`, or was committed *before* the feature commit date
- **After** — filename contains `after`, `new`, `current`, or matches the feature commit date/SHA

**3b: Take a live screenshot if bowser is available and a server is running.**

First, detect bowser:
```
playwright-cli --help >/dev/null 2>&1 && echo "bowser: available" || echo "bowser: missing"
```

If `bowser: missing`, skip to 3c (diagram-only fallback). Do **not** prompt the user to install — note the fallback in the report's hero section and move on. If they want richer reports later, they can run `do-work install bowser`.

If `bowser: available`, probe each candidate dev server in order and capture the first that responds 200:

```
DEV_URL=""
for port in 8080 5173 3000; do
  [ "$(curl -s -o /dev/null -w '%{http_code}' "http://localhost:$port/")" = "200" ] && DEV_URL="http://localhost:$port/" && break
done
```

If `$DEV_URL` is set:

- Use the bowser skill (`playwright-cli`) to screenshot `$DEV_URL`. Save to `ai-reports/<report-stem>.assets/live.png` (the same `<report-stem>` chosen in Step 5 — `yyyy-mm-dd_hhmm_<slug>`).
- If the REQ's Implementation Summary mentions a specific route or panel (e.g., "align panel", "settings tab"), append it to `$DEV_URL` before screenshotting.

If no server responds (`$DEV_URL` empty), note it in the report and skip live screenshots.

**3c: Decide the visual strategy.**

| Situation | Strategy |
|-----------|----------|
| Both before and after images found | Side-by-side toggle + SVG callout annotations |
| Only after (or only live screenshot) | Annotated screenshot with SVG callout arrows |
| No screenshots at all | SVG architecture diagram + Mermaid data-flow chart |
| Mixed (some REQs have screenshots, some don't) | Per-REQ strategy — screenshots where available, diagrams where not |

### Step 4: Build the Visual Assets

**4a: SVG annotation overlay (when screenshots exist).**

For each screenshot, produce an inline `<svg>` positioned absolutely over the `<img>` tag. Add numbered callout circles (⬤ 1, 2, 3…) pointing at the changed UI region, with a caption legend below. Do not draw on the image itself — overlay only.

Callout anatomy:

```html
<div class="screenshot-frame" style="position:relative;display:inline-block">
  <img src="./<report-stem>.assets/after.png" alt="..." style="max-width:100%">
  <svg style="position:absolute;top:0;left:0;width:100%;height:100%"
       viewBox="0 0 [img-width] [img-height]">
    <circle cx="320" cy="140" r="14" fill="#2563eb" fill-opacity=".85"/>
    <text x="320" y="145" text-anchor="middle" fill="white" font-size="13" font-weight="bold">1</text>
    <line x1="320" y1="154" x2="320" y2="180" stroke="#2563eb" stroke-width="2" marker-end="url(#arrow)"/>
  </svg>
</div>
<ol class="callout-legend">
  <li><strong>New "Align" checkbox</strong> — added to the popover; defaults to checked.</li>
</ol>
```

Copy every collected image into `ai-reports/<report-stem>.assets/` and reference it with a relative `src` (e.g., `./<report-stem>.assets/before.png`). Use descriptive names — `before.png`, `after.png`, `live.png`, or `before-<slug>.png` when multiple before/after pairs exist. Do **not** base64-inline; the binaries live next to the HTML and travel with it.

**4b: Mermaid data-flow diagram (when no screenshots, or to supplement architecture explanation).**

```html
<script src="https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js"></script>
<script>mermaid.initialize({startOnLoad:true, theme:'base'});</script>
<pre class="mermaid">
flowchart LR
    A[User clicks Process] --> B{Align checked?}
    B -- Yes --> C[Gaussian Align]
    C --> D[Validate]
    B -- No --> D
    D --> E[Translate]
    E --> F[Grammar]
</pre>
```

Derive the diagram content from the REQ's Implementation Summary — trace the actual code path, not a generic placeholder.

**4c: SVG architecture diagram (component relationships).**

When the feature touches multiple components, produce a hand-coded SVG that shows the component graph: boxes with names, arrows for data/event flow, brief labels on edges. Keep it to one viewport (no scrolling). Use the same color palette as the report.

### Step 5: Write the Report HTML

Write the HTML file to:

```
ai-reports/yyyy-mm-dd_hhmm_<description>.html
```

**Filename rules:**

- Always use the `yyyy-mm-dd_hhmm_` prefix so reports sort chronologically. Generate with `date +%Y-%m-%d_%H%M`.
- `<description>` is a kebab-case slug: UR/REQ ID + a 2–3 word summary (e.g., `UR-246-batch-align-default`, `REQ-699-popover-checkbox`).
- Full example: `2026-05-26_1430_UR-246-batch-align-default.html`
- This handles clashes automatically and keeps the folder sortable.

Ensure both the reports directory and the per-report assets folder exist (`mkdir -p ai-reports/<report-stem>.assets`).

#### Required sections (in order)

**Hero** — Feature name + one-sentence verdict ("What shipped and whether it works"). Large type. No throat-clearing. Lead with the conclusion.

**The Change** — Two columns or a before/after toggle:
- If before+after images: sliding toggle or side-by-side with SVG callouts.
- If only current state: annotated screenshot with callout legend.
- If no screenshots: "What it looked like before" described in a styled callout box + Mermaid diagram of the new flow.

**How It Works** — Data-flow diagram (Mermaid) and/or architecture diagram (SVG). Pull from Implementation Summary. One diagram per concept — do not stack three diagrams that say the same thing.

**What Changed** — Compact table of files modified, what each does. No code snippets. Pointers only (`src/components/AlignPanel.jsx — added Align checkbox + batch flag`).

**Verify It Yourself** — Copy-pasteable shell commands from the REQ's Testing section. One `git show <sha>` block. One test-run command.

**Open Questions / Lessons** — Only if the REQ has a non-empty Lessons Learned or unresolved Open Questions. If empty, omit this section entirely.

#### Design rules

- Single `.html` file. Zero build steps. No npm installs.
- External CDN allowed only for: Tailwind CSS, Mermaid.js. Everything else inline.
- Light theme by default; dark via `@media (prefers-color-scheme: dark)`.
- CSS custom properties at `:root` for `--bg`, `--surface`, `--text`, `--accent`, `--muted`. Light: white/slate-50 bg, slate-800 text, blue-600 accent. Dark: slate-900 bg, slate-100 text, blue-400 accent.
- Large readable type: body 16px min, headings 24–40px.
- Generous whitespace: section padding ≥ 40px.
- Before/after toggle: CSS-only or minimal vanilla JS — no framework.
- Mermaid theme: `base` (works in both light and dark via CSS overrides).
- No emoji in headers or body unless the REQ itself uses them.
- No marketing language ("game-changing", "powerful", "seamless"). Factual only.
- No unearned bullet lists. If it flows as prose, write prose.

#### Before/after toggle pattern (reference implementation)

```html
<div class="toggle-group" role="group">
  <input type="radio" name="view" id="view-before" value="before" checked>
  <label for="view-before">Before</label>
  <input type="radio" name="view" id="view-after" value="after">
  <label for="view-after">After</label>
</div>
<div class="view-panels">
  <div class="panel" id="panel-before"><!-- before content --></div>
  <div class="panel" id="panel-after" hidden><!-- after content --></div>
</div>
<script>
document.querySelectorAll('input[name="view"]').forEach(radio => {
  radio.addEventListener('change', e => {
    document.querySelectorAll('.panel').forEach(p => p.hidden = true);
    document.getElementById('panel-' + e.target.value).hidden = false;
  });
});
</script>
```

### Step 6: Self-Review Against Anti-Slop

Before saving the file, run through all seven principles explicitly:

| Principle | Status | Evidence / fix applied |
|-----------|--------|------------------------|
| 1. Would I read this? | — | — |
| 2. Every claim verified | — | — |
| 3. Compressed | — | — |
| 4. Conclusion first (hero section) | — | — |
| 5. AI honesty tag present if needed | — | — |
| 6. Does this need to exist? | — | — |
| 7. Medium matches stakes | — | — |

Fix any FLAGs before writing the final file. Do not ship a Borderline or Slop report.

### Step 7: Save and Report

1. `mkdir -p ai-reports/<report-stem>.assets` (assets folder already populated in Step 3/4).
2. Write the final HTML file to `ai-reports/<report-stem>.html`.
3. The `<report-stem>.assets/` folder lives next to the HTML — keep them together when moving or sharing the report.
4. Print a short summary:

```
Report generated: ai-reports/<filename>.html

  Feature: <title>
  Evidence: <what visual assets were used>
  Diagrams: <what diagrams were generated>

Open in any browser — no build step needed.
```

Do not pad the summary. No headers or bullet lists unless there are multiple reports.

## Output Format

An HTML file at `ai-reports/yyyy-mm-dd_hhmm_<slug>.html` plus a sibling `ai-reports/yyyy-mm-dd_hhmm_<slug>.assets/` folder containing the referenced PNG/JPG binaries (each with descriptive names like `before.png`, `after.png`, `live.png`). The HTML references them via relative `src`. Plus a one-paragraph stdout summary as shown in Step 7.

## Rules

- **Output goes to `ai-reports/` in the project root** — never `do-work/deliverables/` (that's the present-work explainer's home), never `do-work/working/`, never a custom path.
- **HTML + sibling `.assets/` folder.** Images live in `<report-stem>.assets/` next to the HTML and are referenced by relative `src` — move/copy the pair together. Tailwind and Mermaid load from a CDN; without network, styling is degraded and Mermaid diagrams won't render.
- **Bowser is optional.** If `playwright-cli` is missing or no dev server responds, fall back to SVG + Mermaid diagrams. Don't error, don't block, don't prompt to install.
- **No live screenshot if the dev server isn't running** — note the absence in the hero section and use the diagram fallback. Don't fabricate a "before" state from imagination.
- **Anti-slop principles are loaded inline** (Step 1), not via a separate slop-check pass.
- **Filename uses the `yyyy-mm-dd_hhmm_` prefix** so reports sort chronologically — never just the UR/REQ ID.
- **One UR or one REQ per invocation.** Multi-target batches are not supported; the user can re-run.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "The dev server isn't running but I'll guess at the UI from the diff" | Fall back to the SVG/Mermaid diagrams — do not invent pixels | Fabricated screenshots are worse than no screenshots; they look authoritative and aren't |
| "bowser isn't installed — I'll prompt the user to install it and pause" | Skip live screenshots, use diagram fallback, mention bowser in the report's footer as a next step | The skill is designed to ship a useful report without bowser. Blocking on install defeats the optional-dep design. |
| "The REQ has no Implementation Summary but I'll write the diagram from the diff" | Stop and tell the user — the REQ wasn't actually completed properly | A missing Implementation Summary means review-work didn't run; the report would be guessing at intent |
| "I'll add a fancy intro paragraph before the hero" | Cut it — the hero IS the lead | Anti-slop principle 4: conclusion first. Throat-clearing pushes the verdict below the fold |
| "Two screenshots, one before and one after — I'll show both with no toggle, side by side" | Use the before/after toggle pattern — same screen real estate, one viewport, faster compare | Side-by-side at small screen widths squishes both; the toggle keeps each at full width |
| "I'll base64-inline the screenshots so the HTML is one file" | Save them to `<report-stem>.assets/` and reference with relative `src` | Base64 bloats the HTML ~33% per image and slows first paint; the `.assets/` folder travels next to the HTML, so the pair is just as portable |
| "This is the present-work explainer territory, I'll merge them" | Keep them separate — explainer = concept; ai-report = pixels | Two artifacts can coexist for the same UR; they answer different questions for different audiences |

## Red Flags

- The report is longer than the implementation it describes — you produced slop. Cut.
- A "Before" image was used that has no clear connection to the current feature — mislabeled. Remove it or retitle.
- The data-flow diagram is a generic template — it must match the actual code path from the REQ's Implementation Summary.
- A screenshot is base64-inlined as `src="data:image/...;base64,..."` instead of referenced from `.assets/` — contradicts the layout contract and inflates the HTML.
- The hero section buries the verdict in paragraph 2 or later — move it to sentence 1.
- Mermaid doesn't render (check CDN, check `startOnLoad:true`) — fall back to an SVG diagram instead.
- The "Verify It Yourself" section uses placeholder commands — every command must come from the REQ's Testing section or commit SHA.
- The output landed in `do-work/deliverables/` instead of `ai-reports/` — wrong action's home; move it.
- bowser was missing and you stopped instead of falling back to diagrams — the report should always ship.
- The anti-slop self-check table (Step 6) was skipped or left with `—` rows — you don't get to declare your own work clean without filling it in.

## Verification Checklist

- [ ] Anti-slop principles loaded (Step 1) and Step 6 self-check table completed with no unresolved FLAGs.
- [ ] All images saved in `<report-stem>.assets/` and referenced via relative `src` — no `src="data:image/...;base64,..."` in the HTML.
- [ ] Diagrams derived from actual REQ/code content, not generic placeholders.
- [ ] Hero section leads with the conclusion (feature name + one-sentence verdict).
- [ ] "Verify It Yourself" commands are copy-pasteable and come from the REQ.
- [ ] File saved to `ai-reports/` with the `yyyy-mm-dd_hhmm_` prefix — not `do-work/deliverables/`, not a custom path.
- [ ] No build step required to open the file — plain HTML with CDN-only externals (Tailwind + Mermaid).
- [ ] If bowser was missing, the report still shipped using SVG/Mermaid fallback — no install prompt, no block.
