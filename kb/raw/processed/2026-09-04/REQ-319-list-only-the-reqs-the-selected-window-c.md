---
source_type: req_lesson
req_id: REQ-319
req_path: do-work/archive/UR-065/REQ-319-list-only-the-reqs-the-selected-window-covers.md
date: 2026-08-23
domain: frontend
module: _dev/primes
tags: [frontend, list, only, reqs]
---

# Lessons from REQ-319: List only the REQs the selected window covers

## What the REQ was about

Drop a REQ's row from the Timeline entirely when its span falls outside the visible time
window, instead of listing it as an empty row with a clipped-to-nothing bar.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified) — post-review, F4/F8
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified) — post-review, F8

## What worked

- `: run the sweep before narrowing a declaration, not after.
- **D-03 — The forecast's whole-queue note names no cause.** It opened "Filters are on", which
- was true while a filter chip was the only thing that could shrink the row set. The window
- now shrinks it too — usually much harder — and both can be on at once, so the note said
- something false about the common new case. Rewritten to "This covers the whole queue, not
- the rows shown." Naming which cause is on buys the reader nothing they need and costs a
- branch that can be wrong. The parameter is renamed `showingSubset` to match.

## Back-reference

See `do-work/archive/UR-065/REQ-319-list-only-the-reqs-the-selected-window-covers.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `efc3dff`.
