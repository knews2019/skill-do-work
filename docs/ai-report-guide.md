# AI Report

Generates a self-contained HTML report — live screenshots + SVG callout annotations, before/after toggles, optional AI-generated diagrams when an image-gen CLI is available, and Mermaid/SVG diagrams as the always-available fallback. Output is a folder per report (`index.html` + a `screenshots/` folder, plus a `generated/` folder when AI images are used) under `ai-reports/` at the project root. Pixel-anchored proof-of-work, not a brief.

> **Not to be confused with present-work or pipeline's completion report.** `do-work present work` writes a client-facing brief, an Interactive Explainer (`.single.html`), and optionally a video — explanation of value, not literal pixels. `do-work pipeline`'s completion report is a multi-REQ developer/PM debrief (test deltas, REQ coherence graph, carry-forward work). AI Report is the *visual* artifact: "here are the screenshots, here are the annotated changes, here is the verify-it-yourself link."

## What it produces

One folder per report, with `index.html` and its images beside it:

```
ai-reports/
  2026-06-02_1430_ur-007-add-user-avatar-component/
    index.html
    screenshots/
      before-settings.png
      after-settings.png
    generated/                 (only when AI images were produced)
      01-architecture.png
```

The HTML references images by relative `src`, so moving the whole folder keeps the report working anywhere. Deleting a report is `rm -rf ai-reports/<report-slug>/`. Tailwind and Mermaid load from a CDN.

## How it adapts to what's available

| Available | Report includes |
|-----------|-----------------|
| Live dev server + `playwright-cli` (bowser) | Live screenshots before/after, with SVG callout overlays |
| Saved before/after assets in `do-work/archive/UR-NNN/assets/` (the common case after cleanup), `do-work/user-requests/UR-NNN/assets/`, `do-work/working/`, or images in the feature commit's diff | Side-by-side comparison from the saved assets |
| An image-gen CLI on PATH (e.g. `codex`, `gemini`) | Optional AI-generated architecture/concept/hero visuals, each disclosed with an "AI-generated" caption, in `generated/` |
| Nothing | Falls back to SVG architecture + Mermaid data-flow diagrams. The report still ships. |

The action picks the highest-fidelity option it can run. A bowser-less environment still produces a usable report — just with diagrams instead of screenshots.

## AI-generated visuals (optional)

Claude can't draw raster images itself, so when an image-gen CLI is on PATH (`command -v` finds it) the action delegates a few architecture/concept/hero visuals to it, keeping them in their own `generated/` folder so provenance is physical — `screenshots/` is real, `generated/` is synthetic. Every generated image carries a visible "AI-generated" caption and is verified non-empty before the HTML references it. This is strictly opportunistic: no CLI on PATH means the same SVG/Mermaid diagrams stand in, and the report is no worse off. Screenshots always outrank generated images — they're proof; generated visuals only explain structure and flow.

## SVG callout annotations

Screenshots are inert; callouts make them actionable. Each callout is an SVG overlay anchored to a pixel region with a one-line label naming what changed and why it matters. Anti-slop rules apply: lead with the conclusion (the change), justify in the prose below. The point of a callout is "your eye should land here first."

## Before/after toggle

When both before and after images exist, the report renders a single image element with a toggle (button or hover). The toggle lives in the HTML — no build step. The user sees the change as a flip, not as two adjacent images they have to scan for differences.

## Anti-slop applied inline

Every section passes the seven `crew-members/anti-slop.md` principles, applied as you write — there is no separate `slop-check` step inside `ai-report`. If a section can't justify its existence, it doesn't ship. Pixels first, prose second.

## Input

```
do-work ai-report UR-NNN          Report on every completed REQ under that UR
do-work ai-report REQ-NNN         Report on that single REQ
do-work ai-report                 Most recently completed UR in do-work/archive/
do-work ai-report most recent     Same — explicit form
```

If nothing is `status: completed` for the target, the action stops and says so — there's nothing to report on.

## Output

`ai-reports/yyyy-mm-dd_hhmm_<slug>/index.html` plus the sibling `screenshots/` folder (and `generated/` when AI images were produced). The timestamp prefix is mandatory — never just the UR/REQ ID — so reports sort chronologically. The HTML opens directly in a browser. Stakeholders can read it without running anything.

## Key rules

- **Output goes to `ai-reports/` at the project root**, not `do-work/`. The report is a project deliverable, not a do-work bookkeeping file.
- **One folder per report.** `index.html` + `screenshots/` (+ `generated/`). Moving the whole folder must keep the report working.
- **Anti-slop is inline.** Loaded in Step 1 (`crew-members/anti-slop.md`), applied to every section through Step 6. Never declare your own work clean by leaving the self-check table blank.
- **Bowser is nice-to-have, not required.** If screenshots fail, fall back to SVG/Mermaid and ship anyway.

## When NOT to use

- The work has no user-visible output (infra-only, refactor, tooling) — the report is empty by construction. Use `do-work present work` instead.
- You want a value-prop / explainer artifact for a stakeholder — use `do-work present work` (it writes a `.single.html` Interactive Explainer to `do-work/deliverables/`).
- You want a multi-REQ developer/PM debrief of a pipeline run — use `do-work pipeline`'s completion report.
- The work is still in progress — there's nothing shipped to report on.
