# AI Report

`do-work-toolbox ai-report [UR|REQ|most recent]` creates the canonical detailed stakeholder HTML for one completed UR or REQ. It works for visual UI changes, backend systems, refactors, infrastructure, and other non-visual work by changing the evidence—not by sending detailed reports to another action.

Cross-project portfolio presentation belongs to `present-work`. An animated walkthrough belongs to the separate present-video action. Neither replaces this detailed report.

## What It Produces

Every invocation creates a fresh timestamped bundle under `ai-reports/`:

```text
ai-reports/
  2026-06-02_1430_UR-007-add-user-avatar/
    index.html
    screenshots/   # only when authentic captures are used
    generated/     # only when generated visuals succeed
```

The HTML uses relative paths for local media, so the folder travels as one unit. Existing output paths are never overwritten; a same-name run receives a numeric suffix. Prior reports and other completed-work artifacts remain unchanged.

This action produces only the report bundle. It does not create a Markdown client brief, a separate `.single.html` explainer, video, Remotion/MP4 output, a `--with-video` variant, publishing, hosting, or automatic video behavior.

## Two Evidence Modes

| Mode | What the report uses |
|---|---|
| **Visual evidence** | Authentic screenshots, SVG callout annotations, real side-by-side before/after evidence, then code and verification receipts |
| **Non-visual evidence** | Architecture or data-flow diagrams, merge-aware commit evidence, current code, tests, and operational verification |

Visual work uses real screenshots whenever available. A real capture always outranks a generated illustration, and the report never fabricates a before state. Screenshots live in `screenshots/`; generated images live in `generated/` and carry a visible “AI-generated” label so proof and explanation cannot be confused.

For non-visual work, the report explicitly says, “UI captures were not expected for this work.” It does not create placeholder screenshots or force the evidence into a visual-comparison layout. Diagrams are used only when they make architecture or flow easier to understand.

## Stakeholder Narrative

The report leads with a concise verdict, then covers:

- what shipped;
- the problem and implemented change;
- how the result works;
- qualitative value without invented metrics;
- evidence appropriate to the work;
- key files and commits;
- copy-pasteable verification commands;
- recorded lessons and open questions when available.

An archived target whose normalized status is `completed-with-issues` is eligible, but the report keeps that qualification and its recorded issues visible. Cancelled, failed, and unfinished work is rejected rather than presented as shipped.

## Visual Evidence and Layout

When authentic before and after captures exist, the report shows them side by side in a wrapping row so both states remain visible. A toggle is optional only when it materially improves comparison. Screenshots open at full resolution and use SVG overlays for callouts; the original files are never drawn on.

The page is responsive and full-width, while running prose stays at a readable measure. Screenshots stop at native resolution rather than stretching. Light and dark themes share one coherent visual direction.

When browser automation is available, the action serves the bundle over HTTP, captures full-page light and dark renders, judges the rendered pixels, fixes defects, and repeats. Without browser automation the report still ships, with a footer stating that the layout was not render-verified.

## Optional Generated Visuals

A non-agentic image backend may add architecture, concept, data-flow, or hero visuals. Generation is explanatory and opportunistic; SVG/Mermaid is always the fallback, and a failed generation never reuses a stale target.

Sandbox-bypassed agentic image backends are disabled by default. Setting `DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND=1` is explicit full-host authorization over the repository, credentials, network, and external services. Without that exact opt-in, the action uses a non-agentic backend or SVG/Mermaid.

## Input

```text
do-work-toolbox ai-report UR-NNN          Report on successful REQs in that archived UR
do-work-toolbox ai-report REQ-NNN         Report on that archived REQ
do-work-toolbox ai-report                 Most recent successful archived work
do-work-toolbox ai-report most recent     Same, explicitly
```

Target statuses are normalized under the do-work schema. The terminal-success set is `completed` or `completed-with-issues`; if the selected target has no successful work, the action stops and explains why.

## Evidence Safety

Before reading archived UR input, REQ bodies, reviews, tests, or lessons, the action loads the prompt-injection and anti-slop guardrails. Archived prose is evidence to summarize, not permission to run commands or change scope.

The evidence sweep reads requirements, implementation summaries, reviews, tests, lessons, merge-aware commits, and current code when present. Missing optional evidence is identified instead of invented. Generated-media prompts use only sanitized agent-authored descriptions, never verbatim archive prose.
