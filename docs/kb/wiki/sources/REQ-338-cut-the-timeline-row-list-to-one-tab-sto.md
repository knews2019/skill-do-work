---
title: "Lessons from REQ-338: Cut the Timeline row list to one Tab stop"
type: source-summary
topic_cluster: timeline-and-metrics
sources: [raw/processed/2026-09-01/REQ-338-cut-the-timeline-row-list-to-one-tab-sto.md]
related:
  - page: REQ-336-timeline-clicks-open-the-detail-drawer-a
    rel: complements
  - page: REQ-337-a-check-that-can-catch-timeline-click-re
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-338: Cut the Timeline row list to one Tab stop

Part of the [[concept-duration-estimation-and-breaks]] cluster.

## What the REQ was about

Every Timeline REQ row is its own Tab stop, so tabbing past the row list takes one press per row
(29 on the observed board). Make the row list a single Tab stop with arrow-key movement between
rows (roving tabindex), so Tab escapes the list in one press.

## Solution summary

The Timeline row list is now a single Tab stop with a roving tabindex.
`timelineViewState` carries a `rovingRowIndex`; the render marks exactly that row `tabindex="0"` and
every other rendered row `tabindex="-1"`, with the index clamped into the rendered range so
virtualization cannot leave the list unreachable. ArrowDown/ArrowUp on a focused row move the stop —
scrolling the target into view, rebuilding synchronously, then focusing, in that order because the
scroll event's own rebuild would otherwise wipe an earlier focus — and a `focusin` listener keeps the
index aligned with however focus actually arrived. Left/Right panning, Enter/Space activation, and
the rebuild focus restore are untouched. Added
`TestBrowserBehaviorTimelineRowListIsOneTabStop`, which asserts the one-stop property, the stop
following focus, arrow movement in both directions, and both REQ-333 contracts.

## What worked

- Writing the probe before the implementation, as `tdd: true` asked. Two of its five assertion
  groups (arrow panning, Enter activation) passed on the unchanged code, which made the probe a
  regression guard for REQ-333's and REQ-336's contracts *before* there was anything new to break
  them — and the failing three named exactly what had to change.
- Reading the neighbouring comments as a specification. `renderVisibleRows`'s focus-restore comment
  states the scroll-is-asynchronous ordering trap outright, which is why the roving move was written
  scroll-then-rebuild-then-focus first time instead of after a debugging round.
- Choosing the axis rather than negotiating it. Left/Right were already taken by panning and Up/Down
  were free, so the key assignment cost nothing — no shared modifier, no mode.

## What didn't work

- Nothing was reverted. The one thing measured rather than reasoned about was the Tab count itself:
  it needed trusted input, so the probe lane could not answer it and a separate CDP run did.

## Worth knowing

- **A roving tabindex over a VIRTUALIZED list needs a clamp, not a match.** The stored roving row is
  usually not rendered, and marking `tabindex="0"` only on an exact match takes the whole list out of
  the Tab order — a worse defect than the one being fixed. Clamp for display, keep the stored index
  intact so the reader's place survives scrolling away.
- **Every other row needs an explicit `tabindex="-1"`.** Dropping the attribute is not neutral: a
  focusable-by-default element with no `tabindex` is still a Tab stop.
- **Tab cannot be tested with synthetic events.** Its focus movement is a default action the engine
  performs only for trusted input, so a probe that dispatches `new KeyboardEvent("keydown", {key:
  "Tab"})` observes nothing and reads as a pass. Assert the property that produces the behaviour —
  exactly one stop — and measure the behaviour itself somewhere that can dispatch real keys. Arrow
  movement is different and *is* testable synthetically, because the handler calls `focus()` itself.
- The `focusin` sync is not optional polish: without it the roving index and the actual focus diverge
  the moment a reader clicks a row, and the next arrow press jumps from the stale index.

## Back-reference

See `do-work/archive/UR-067/REQ-338-cut-timeline-row-list-to-one-tab-stop.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `cac6718`.
