---
source_type: req_lesson
req_id: REQ-089
req_path: do-work/archive/UR-017/REQ-089-board-drawer-copy-omits-the-tickets-frontmatter.md
date: 2026-08-04
domain: frontend
module: tools/queue-kanban
tags: [frontend, queue-kanban, board, drawer, copy]
---

# Lessons from REQ-089: The board drawer's Copy button omits the ticket's frontmatter, so the paste carries no status, domain or timestamps

## What the REQ was about

Clicking **Copy** in the board's card drawer puts the `# REQ-NNNN: <title>` heading plus the rendered
body's source Markdown on the clipboard. None of the metadata the drawer displays directly above the
body comes with it — status, domain, user request, created, claimed, tree.

## Solution summary

`RequestTicket` and `UserRequestTicket` gained `FrontmatterMarkdown`, populated at both parse sites by removing the body as an exact suffix of the raw file text — so the fence keeps its original key order, comments, spacing and line endings. `buildGeneratedBoardMarkdownData` now projects `FrontmatterMarkdown + BodyMarkdown`, which fixes the static bundle and the live server together because both call it. The click handler returns that payload verbatim on the primary path and applies `copyTextWithHeading` only on the two fallback branches; the stale comment claiming "the stored Markdown is the file body only" was replaced with one that names both shapes and says which path owns each. The guide's Copy sentence was rewritten, including the three drawer rows that cannot carry.

## What worked

- **Testing the round-trip from parsed files, not from hand-built structs.** The struct-level test that already existed passes both before and after the change, because a zero-value fence concatenates to nothing. Only a test that reads a real file can tell "the original bytes survived" from "a fence was produced" — and the fixture's comment line and irregular spacing are what make it discriminating.
- **Diffing the actually-generated bundle.** The unit tests prove the projection; extracting the payload from `board-markdown.js` and `diff`ing it against the file on disk proves the whole chain, including the JSON encoding hop that no Go test covers.

## What didn't work

- Nothing failed outright, but the REQ's file list was one short. `serve.go` shares the projection, so its contract test broke — discovered by running the suite, not by reading the REQ's Context, which named only the static generator.

## Worth knowing

- **`buildGeneratedBoardMarkdownData` is the single projection point for both the static bundle and the live server** (`generate.go:193`, `serve.go:322`). Changing it changes `do-work board` in both modes at once. Convenient here; worth knowing before assuming a change is static-only.
- **Exact-suffix removal must stay inside the `hasFrontmatter` guard.** `splitFrontmatter` returns the whole text as `bodyText` when there is no fence, so the arithmetic yields `""` correctly — but only because the guard skips it. Moving the slice outside the `if` would still compute `""` today and become wrong the moment `splitFrontmatter`'s no-fence contract changes.
- **Three drawer rows can never round-trip**, by construction: `TreeSection`, `WriteSetOverlaps` and `Dependents` are derived at parse/bucket time and are documented as never read from frontmatter. The guide now says so, so their absence reads as designed rather than as an incomplete implementation.

## Back-reference

See `do-work/archive/UR-017/REQ-089-board-drawer-copy-omits-the-tickets-frontmatter.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `ba54b5d`.
