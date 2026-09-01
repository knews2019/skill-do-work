---
title: "Lessons from REQ-387: Keep a spliced title from changing how the pasted Markdown parses"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-387-keep-a-spliced-title-from-changing-how-t.md]
related:
  - page: concept-kanban-board-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-387: Keep a spliced title from changing how the pasted Markdown parses

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

An expanded title is inserted into the document's Markdown verbatim, after a 60-character cut. Two
characters in a title can change how the pasted file parses: a pipe inside a table row, and a backtick
the cut leaves unbalanced.

## Solution summary

- `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js` (modified). Sanitizes only the existing shortened title before splicing: remove code-span backticks, then escape every ASCII punctuation character in one pass. Full appendix, offsets and drawer title helper stay unchanged. Six additi

## What worked

A Markdown source delimiter count is not a parse oracle. GFM can silently discard surplus cells, and two backticks can be an unmatched delimiter. Test the real rendered structure, preserved neighboring cell contents and author code; escape pre-existing backslashes before inserting a literal pipe escape.

## Back-reference

See `do-work/archive/UR-075/REQ-387-keep-a-spliced-title-from-changing-how-the-paste-parses.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `a0d0b350`.
