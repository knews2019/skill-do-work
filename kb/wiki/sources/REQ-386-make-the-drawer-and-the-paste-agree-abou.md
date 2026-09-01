---
title: "Lessons from REQ-386: Make the drawer and the paste agree about a body H1 that restates the title"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-386-make-the-drawer-and-the-paste-agree-abou.md]
related:
  - page: concept-kanban-board-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-386: Make the drawer and the paste agree about a body H1 that restates the title

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

The drawer deletes a body's opening H1 when it restates the frontmatter title, then decides which
mention expands. The clipboard keeps that H1 and counts it as the first prose mention. Pick one rule
and apply it to both.

## Solution summary

Saving a copied ticket back to disk no longer breaks duplicate-heading suppression. The drawer and Copy agree about which visible prose occurrence first receives the ticket title.

## What worked

Rendered heading text is not the Markdown heading source. Reuse the renderer, account for its preprocessing, and explicitly match JavaScript whitespace and full lowercase before claiming two languages perform the same comparison. When reparsing a fragment, carry reference definitions from the document or the fragment can silently change a heading's text. A copy/save/rebuild test catches heading annotation that a single-surface title test misses.

## Back-reference

See `do-work/archive/UR-075/REQ-386-agree-on-the-restating-h1-between-drawer-and-paste.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `59577def`.
