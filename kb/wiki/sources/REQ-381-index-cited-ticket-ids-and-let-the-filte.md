---
title: "Lessons from REQ-381: Index cited ticket ids and let the filter box match them"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-381-index-cited-ticket-ids-and-let-the-filte.md]
related:
  - page: REQ-382-expand-ticket-ids-written-as-markdown-li
    rel: complements
  - page: REQ-385-treat-an-underscore-as-a-ticket-id-bound
    rel: depends-on
  - page: REQ-386-make-the-drawer-and-the-paste-agree-abou
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-381: Index cited ticket ids and let the filter box match them

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

Typing `REQ-1679` into the board's filter box returns every card whose body cites it, not just the
card whose title happens to contain that text. The citation set is computed once on the Go side and
shipped per request, so the filter reads an index rather than re-scanning bodies in the browser.

## Solution summary

The board can now answer “which cards cite this ticket?” before Copy data loads. Citation-only hits carry a small explanation, and static generation and live refresh use the same resolved index.

## Worth knowing

- Preserve exact-record precedence before case-folded alias resolution. Sending the filter's lowercase needle directly to a case-sensitive exact resolver can incorrectly choose a compound alias.
- Search sets and annotatable occurrences have different exclusions. Derive them together, but collect citations before clipboard-only suppressions so later presentation changes cannot silently remove search results.
- A string-key lookup fed by arbitrary search text needs a null prototype or an own-property check; canonical ticket keys alone do not make the query safe.

## Back-reference

See `do-work/archive/UR-076/REQ-381-index-cited-ticket-ids-and-filter-on-them.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `961fbf84`.
