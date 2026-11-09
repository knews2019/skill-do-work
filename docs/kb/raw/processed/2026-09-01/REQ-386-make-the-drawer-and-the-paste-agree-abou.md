---
source_type: req_lesson
req_id: REQ-386
req_path: do-work/archive/UR-075/REQ-386-agree-on-the-restating-h1-between-drawer-and-paste.md
date: 2026-08-27
domain: frontend
module: _dev/primes
tags: [frontend, drawer, paste, agree, about]
---

# Lessons from REQ-386: Make the drawer and the paste agree about a body H1 that restates the title

## What the REQ was about

The drawer deletes a body's opening H1 when it restates the frontmatter title, then decides which
mention expands. The clipboard keeps that H1 and counts it as the first prose mention. Pick one rule
and apply it to both.

## Solution summary

Saving a copied ticket back to disk no longer breaks duplicate-heading suppression. The drawer and Copy agree about which visible prose occurrence first receives the ticket title.

## Worth knowing

Rendered heading text is not the Markdown heading source. Reuse the renderer, account for its preprocessing, and explicitly match JavaScript whitespace and full lowercase before claiming two languages perform the same comparison. When reparsing a fragment, carry reference definitions from the document or the fragment can silently change a heading's text. A copy/save/rebuild test catches heading annotation that a single-surface title test misses.

## Back-reference

See `do-work/archive/UR-075/REQ-386-agree-on-the-restating-h1-between-drawer-and-paste.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `59577def`.
