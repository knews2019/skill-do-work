---
id: UR-075
title: 'Beautifully rendered architecture report'
created_at: 2026-08-26T17:56:52Z
requests: [REQ-378]
word_count: 30
---

# Beautifully Rendered Architecture Report

## Full Verbatim Input

> ````text
> basically I want it beautifully rendered so it's easy to understand and so it has all the bells and whistles
> ````

> ````text
> ok, ask away and capture the intent via do-work capture-request
> ````

"It" is the architecture report: the conversation preceding this capture established that `do-work architecture-report` publishes Markdown only (the action mandates "Markdown with inline Mermaid only. No HTML.") and that no REQ existed asking for an ai-report-style rendering — the user had remembered the bundle-shape decision (commit b89a914) as a format decision.

## Capture-Time Decisions

Resolved interactively during capture, via two rounds of structured questions:

- **Scope: capability only.** The action gains the capability; the current 2026-08-26 report is not rendered as part of this work. A rendered report arrives the next time the action is invoked.
- **Home: the ai-report home.** User's words: "at the same place as any ai-report lives" — the report's own dated bundle under `ai-reports/`, holding `index.html`.
- **Treatment: redesigned, not faithful.** The HTML is free to reorganize and tell the story visually rather than mirroring a fixed section list.
- **No markdown file at all, and freeform.** User's words on the proof question: "yes, but note that I don't want the markdown file at all, and the HTML file needs to be freeform in the sense that as models evolve, I keep getting a better architecture view into it."
- **Delta survives as authored prose.** Byte-identical carry-forward cannot survive freeform re-authoring; the user chose an authored "changed since last report" opening section, written by reading the previous report, over dropping change tracking entirely.
- **The existing markdown bundle stays.** `ai-reports/2026-08-26_1709_architecture-report/` is kept as committed history; future runs ignore markdown bundles when looking for a prior report.

---
*Captured: 2026-08-26T17:56:52Z*
