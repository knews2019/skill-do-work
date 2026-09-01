---
title: "Lessons from REQ-122: The By UR lens counts recently-done work as active and honors the window buttons"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-122-the-by-ur-lens-counts-recently-done-work.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-122: The By UR lens counts recently-done work as active and honors the window buttons

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

On a fully-shipped queue — `do-work/queue/` and `do-work/working/` empty, every archived
REQ terminal — the board's **By UR** lens with `URs = Active` renders nothing, while the
**Columns** lens on the same board at the same moment shows the recently-done cards for
those very URs. The board simultaneously claims there is nothing to show and shows it.
That is the normal steady state after a run finishes, so the lens is unusable for the
reader who just wants to see what a session touched.

Four defects, all in `tools/queue-kanban/web/`:

## Solution summary

Files changed:

## Worth knowing

- **A predicate that reads as a property of the data can silently be a property of the
  clock.** `userRequestIsActive` looked like a pure status question and was really "is
  this queue mid-flight", so it was correct in every state anyone tested in and
  unsatisfiable in the state the board spends most of its life in. The tell was two views
  of one dataset disagreeing on screen at the same moment — worth treating as a
  high-signal bug shape rather than a rendering glitch.
- **Assert on a slice, not the page, when the claim is "X is called from Y".** Both
  `renderUserRequestLens` and `renderColumns` appear all over the inlined bundle, so a
  whole-page `strings.Contains` would have passed against the buggy version and pinned
  nothing. Brace-matching the handler out of the source is what makes the test able to
  fail — and it did fail first, which is the only reason it is worth keeping.
- **A cache invalidated in one handler needs its siblings checked.** The `renderedOnce`
  guard was handled correctly by the `[data-ur-activity]` handler and not at all by
  `[data-window-hours]`. Whenever a render cache exists, every control that changes an
  input to that render is a call site — grep the guard, not the feature.

## Back-reference

See `do-work/archive/UR-026/REQ-122-by-ur-lens-recent-window.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `684f507`.
