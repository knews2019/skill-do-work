---
title: "Lessons from REQ-384: Deliver the architecture report as a freeform HTML bundle"
type: source-summary
topic_cluster: presentation-and-reporting
sources: [raw/processed/2026-09-01/REQ-384-deliver-the-architecture-report-as-a-fre.md]
related:
  - page: concept-completed-work-presentation
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-384: Deliver the architecture report as a freeform HTML bundle

Part of the [[concept-completed-work-presentation]] cluster.

## What the REQ was about

Change `do-work-toolbox architecture-report` so a run publishes one beautifully rendered, self-contained HTML report — `ai-reports/<yyyy-mm-dd>_<hhmm>_architecture-report/index.html`, the same bundle home `ai-report` uses — and no markdown file at all. The HTML is deliberately freeform: the action states the quality bar and the invariants, never a fixed section-by-section layout, so a more capable future model produces a better architecture view instead of the same template filled in better.

## Solution summary

Changed the architecture-report capability, without invoking it or generating a report. Existing Markdown history remains unchanged; each future report has one self-contained index.html, an authored opening change account, and freely designed diagrams/navigation.

## What worked

A filename used as a publication marker should become visible only after the copy is complete and verified. Separating machine-readable metadata from visible structure allows deterministic history lookup without fixing the report's design.

## Back-reference

See `do-work/archive/UR-077/REQ-384-deliver-the-architecture-report-as-a-freeform-html-bundle.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `c32e1d53`.
