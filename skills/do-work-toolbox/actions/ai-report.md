# AI Report Action

> **Part of the do-work-toolbox skill.** Creates the canonical detailed stakeholder HTML for one completed UR or REQ. It adapts its evidence to visual UI work, backend work, refactors, infrastructure, and other non-visual changes while preserving a timestamped, self-contained `ai-reports/` bundle. User-facing walkthrough: [`docs/ai-report-guide.md`](../docs/ai-report-guide.md).

`ai-report` is the only action that produces detailed stakeholder-facing HTML. The narrow open-questions digest for an outside stakeholder is not a detailed report — that is `stakeholder-report.md`, a non-routed file invoked by core. Cross-project portfolio presentation belongs to `present-work`; an animated walkthrough belongs to the separate present-video action.

## Philosophy

- **Evidence before decoration.** Use the strongest authentic evidence the work produced, then explain it for a stakeholder.
- **One report shape, two evidence modes.** Visual work leads with real captures; non-visual work leads with architecture, code, commit, test, and operational evidence.
- **Self-contained timestamped bundle.** The report and its local assets travel together under `ai-reports/<report-slug>/`.
- **Rendered output is the artifact.** When browser automation is available, inspect full-page light and dark renders; source review is not visual QA.

## When to Use

**Use when:**

- The user wants a detailed presentation of one completed UR or REQ.
- The work may be visual, backend, refactoring, infrastructure, or another evidence-bearing completed change.
- A stakeholder needs the verdict, shipped behavior, value, key files and commits, and verification in one HTML report.

**Do NOT use when:**

- The user wants a cross-project portfolio; use `do-work-toolbox present-work`.
- The user wants an animated walkthrough; use the separate present-video action.
- The target is unfinished or unsuccessful; report its status instead of presenting it as shipped.

## Input

`$ARGUMENTS` is `UR-NNN`, `REQ-NNN`, `most recent`, or blank. One invocation covers one UR or one REQ; blank is the explicit `most recent` form.

## Steps

### Step 1: Resolve and Read the Completed Work

Read and follow [`completed-work-presentation-reference.md`](completed-work-presentation-reference.md) in full **before opening archived user content**. It is the sole contract for safety load order, target resolution, archive fields, missing evidence, merge-aware commit and current-code inspection, evidence honesty, and no-overwrite publication. Do not recreate those rules here.

Build the reference's provenance ledger for the selected work. Its commit inspection follows the canonical [Merge-aware commit diff](../../do-work/docs/prescribed-shell-primitives.md#merge-aware-commit-diff) contract. Complete this read before drafting stakeholder claims or generating media.

### Step 2: Choose the Evidence Mode and Bundle Path

Choose the mode from the work, not from the tools installed:

| Mode | Condition | Primary evidence |
|---|---|---|
| **Visual evidence** | The shipped result has a UI or other visible state for which authentic captures are relevant | Real screenshots, SVG callouts, authentic before/after comparison, then code/tests |
| **Non-visual evidence** | UI captures were not expected for the work | Architecture or data-flow diagrams, merge-aware commit evidence, current code, tests, and operational verification |

A multi-REQ UR may use the appropriate mode per section, but it remains one report. Never force a non-visual change into screenshot-shaped cards, and never downgrade visual work to generic diagrams when authentic captures are available.

Derive `<report-slug>` as `yyyy-mm-dd_hhmm_<description>`, where the description contains the UR/REQ ID and a short kebab-case summary. Use `ai-reports/<report-slug>/` as this consumer's preferred bundle path and apply the shared reference's **Collision-Safe Publication** section before creating it. Create `screenshots/` only when authentic captures will be included and `generated/` only when current-run generated images succeed.

### Step 3: Collect Mode-Appropriate Evidence

#### Visual evidence mode

Search archived `assets/`, matching development captures under `do-work/working/`, and image paths from the feature commits. Path-only commit inspection follows the canonical [Commit file listing](../../do-work/docs/prescribed-shell-primitives.md#commit-file-listing) rule. A loose project-root image is not report evidence.

Use provenance, not a filename guess alone, to classify a capture as before or after. Prefer, in order:

1. an authentic before-and-after pair;
2. an authentic current or live capture;
3. architecture/data-flow explanation with an explicit note that captures were unavailable.

If browser automation is available and a relevant dev server is already running, capture the current shipped route into the report's `screenshots/` folder. Do not block or prompt for an install when browser automation or a server is missing. Never fabricate a screenshot or a visual before state.

Real screenshots outrank synthetic visuals. Keep them in `screenshots/`, separate from every generated image in `generated/`, and describe that provenance in captions.

#### Non-visual evidence mode

Trace the implemented architecture or data flow from the merge-aware commit plus current-code inspection, then connect it to recorded tests and operational checks. The report must say exactly: **“UI captures were not expected for this work.”** Use diagrams only when they materially clarify relationships or flow; code, commit, test, and operational receipts remain the evidence.

Do not create placeholder screenshots, an invented before panel, or empty visual comparison controls.

#### Optional generated visuals

When a generated concept, architecture, or flow image would improve comprehension, follow **Image Generation Backend** in [`ai-report-reference.md`](ai-report-reference.md). Generation is opportunistic explanation, never proof and never a screenshot substitute.

The agentic fallback remains disabled unless `DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND=1` is explicitly set. When that opt-in applies, the reference's helper must run the backend from a `mktemp -d` scratch directory protected with `chmod 700`; otherwise use the non-agentic backend or SVG/Mermaid fallback.

### Step 4: Build Evidence Visuals

For visual evidence, copy each real capture into `screenshots/` with a descriptive name and reference it by relative path. Wrap screenshots in links to their own full-resolution files. Put numbered callouts in an inline SVG overlay with the capture's real `viewBox`, and set the overlay to `pointer-events:none` so it never blocks the link.

Show an authentic before/after pair side by side in a wrapping row by default. Use the optional toggle in **Before/After Toggle Reference Implementation** only when the frames genuinely cannot fit side by side and interaction improves comparison.

For either evidence mode, use Mermaid, hand-authored SVG, or a successfully generated image for architecture and data flow. Derive diagram content from actual implementation and current-code evidence. Follow **SVG Data-Viz Rules** in `ai-report-reference.md`; keep synthetic output visibly labeled and physically separate from screenshots.

### Step 5: Write the Detailed HTML

Write `ai-reports/<report-slug>/index.html` using **Report Design Rules** in `ai-report-reference.md`. Required narrative, in this order:

1. **Verdict** — what shipped, whether verification supports it, and any recorded completed-with-issues qualification.
2. **What Shipped** — concise delivered behavior.
3. **Problem and Change** — the prior problem and the implemented change, without inventing a visual before state.
4. **How It Works** — an evidence-derived flow or architecture explanation.
5. **Value Delivered** — qualitative stakeholder value only; no fabricated metrics.
6. **Evidence** — visual or non-visual receipts appropriate to the selected mode, including the required no-UI-captures statement in non-visual mode.
7. **Key Files and Commits** — compact file roles, commit identifiers when recorded, and any current-code drift from the historical commit.
8. **Verify It Yourself** — copy-pasteable recorded test, operational, and canonical commit-inspection commands, accurately labeled as run or unrun.
9. **Lessons and Open Questions** — include when present; handle missing optional records under the shared evidence contract.

Use one coherent responsive layout with full-width wrapping bands, readable prose measures, and native-resolution image caps. Add interaction only when it improves comparison or comprehension; static HTML is the default.

### Step 6: Review the Claims and Artifact Boundaries

Apply every current principle from `../../do-work/crew-members/anti-slop.md`; do not rely on a copied principle count. Verify each claim against the provenance ledger, lead with the verdict, compress repetition, disclose synthetic media, and remove any visual that only decorates.

Confirm that the invocation creates only the report bundle. It must not create a Markdown client brief, a separate `.single.html` explainer, a video, Remotion/MP4 output, a `--with-video` path, or any automatic video behavior. It also does not publish, host, or search for distribution targets.

### Step 7: Render and Judge

When browser automation is available, serve the report folder over HTTP—never `file://`—and take full-page screenshots in both light and dark rendering contexts. Save judge captures outside the report bundle, inspect the images themselves, fix defects, and repeat until clean. Any report containing text-bearing SVG requires at least two render-and-judge passes.

Judge width usage, table shape, diagram informativeness, emphasis hierarchy, light/dark contrast, SVG label collisions, edge clipping, screenshot sharpness, and responsive stacking. If browser automation is unavailable, ship the report and state in its footer that the layout was not render-verified.

### Step 8: Verify and Report the Result

Confirm `index.html` exists, every relative asset resolves, screenshots open at full resolution, synthetic assets are disclosed, and the final bundle satisfies the shared **Collision-Safe Publication** section. Remove temporary judge captures and stop the temporary HTTP server.

Print a compact summary containing the report path, target and verdict, evidence mode, evidence used, recorded issues, and light/dark render status.

## Output Format

A fresh self-contained folder at `ai-reports/yyyy-mm-dd_hhmm_<slug>/` containing `index.html`, plus `screenshots/` when authentic captures are used and `generated/` when generated visuals succeed. All local assets use relative references. The timestamped folder is the only stakeholder artifact this action publishes.

## Rules

- Screenshots are authentic evidence; generated images and diagrams are explanation. Never blur their provenance.
- Browser and image generation tooling are optional. Missing optional tooling changes the evidence presentation, not whether a valid report can be produced.
- Tailwind CSS and Mermaid.js are the only allowed CDN dependencies; everything else is inline or co-located.

## Verification Checklist

- [ ] Shared completed-work reference loaded before archive content; target, evidence ledger, and **Collision-Safe Publication** contracts satisfy it.
- [ ] Visual evidence uses authentic captures, SVG annotations, responsive before/after where available, and distinct synthetic provenance.
- [ ] Non-visual evidence states UI captures were not expected and uses commit, current-code, architecture/data-flow, test, and operational receipts without fabricated screenshots.
- [ ] Stakeholder narrative includes verdict, shipped change, problem/change, operation, qualitative value, files/commits, verification, and available lessons/questions.
- [ ] Report is responsive, self-contained at the folder level, and render-judged in full-page light and dark when browser automation exists.
- [ ] No brief, separate explainer, video, publishing, hosting, or search artifact was created.
