---
id: REQ-319
title: "List only the REQs the selected window covers"
status: pending
created_at: 2026-08-22T22:08:34Z
user_request: UR-065
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-318]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-318, REQ-320, REQ-321, REQ-322, REQ-323, REQ-324]
batch: timeline-ux-audit
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# List Only the REQs the Selected Window Covers

## What

Drop a REQ's row from the Timeline entirely when its span falls outside the visible time
window, instead of listing it as an empty row with a clipped-to-nothing bar.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Pick *Day* on a 309-REQ board and you get 309 rows, four of which have a bar. The other 305
are labels above empty space, and finding the four means scrolling past them.

## Context

`renderTimelineView` derives `rows` once per render from the payload and the shared filters,
and every later window move (`renderAll`) only redraws. `drawSegment` returns early for a
span outside the window, which is what leaves the empty row behind.

The row set is read by more than the chart: the virtualized scroll extent
(`rows.length * TIMELINE_ROW_HEIGHT`), the subhead's counts, the `<details>` table, and the
`filtersActive` flag handed to `renderTimelineForecast`. All of them have to describe the
same set, and the set now changes on every window move rather than once per render.

## Detailed Requirements

- **Membership rule.** A REQ is in the window when its drawn span overlaps the window at
  all: `[createdTime, endInstant]` intersects `[windowStart, windowEnd]`, where `endInstant`
  is `completedTime`, or the now-line for an open span, or the end of an attached projected
  segment where one exists. A REQ that started before the window and is still running is
  therefore IN — the reader is asking what was happening during this period, and it was.
- Write the membership rule as its own pure function so it can be probed without a DOM, the
  way `timelineZoomedWindow` and `timelineVisibleRowRange` already are.
- Re-derive the set on every path that moves the window: the zoom buttons, the wheel, the
  drag, the keyboard, Day/Week/Month, ‹ and ›, *Now*, and *Fit all*.
- The subhead states the windowed count and says it is a window, not the whole queue —
  it currently claims all 309 with no qualifier.
- The `<details>` table lists the same set as the chart.
- A window containing no REQ says so in its own words. It must not read like the
  "no REQ matches the current filters" empty state, and it must not read like the
  "no readable `created_at`" one; the fix for an empty window is to widen it.
- The forecast paragraph's whole-queue note already fires when the drawn rows are fewer than
  the payload's — check that it still reads correctly when the reason is the window rather
  than a filter chip, and reword it if it does not (`REQ-305`'s subject).

## Constraints

- The projection is never re-derived client-side. Windowing hides rows; it does not change
  the forecast.
- Virtualization must keep working: the scroll extent follows the windowed set, and the
  scroll position needs to survive a window move without throwing the reader somewhere
  arbitrary.
- Serial with the rest of the `timeline-ux-audit` batch.

## Dependencies

`depends_on: [REQ-318]` — that REQ rewrites the same row-derivation and the same subhead
sentence. Ordering, not logic: run 318 first and this one has no merge to resolve.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repo's archive, open the Timeline tab, and
press *Day*. The subhead says 309 REQs and the scroll container holds 309 rows; all but a
handful are a label with no bar, because their spans lie outside the day on screen.

**Why RED now:** `rows` is filtered by the shared filters only and never by the window;
`drawSegment` silently skips a span outside the window and leaves the row.

**GREEN when:** the same *Day* press leaves only the REQs whose spans touch that day —
every listed row has something drawn on it; the subhead's count equals the number of rows;
stepping to the next day with › re-derives the list; a day with nothing in it says so; and
the table under "Every REQ, as a table" lists exactly the rows on the chart.

**Validation:** User confirmed (stated as item 2 of the request). The intersect-not-contain
membership rule is inferred during capture and stated above.

## Assets

Screenshot described in `do-work/user-requests/UR-065/input.md` — forty rows visible,
roughly six with a bar wide enough to see.

---
*Source: "2. Reqs that are not in the period selected should not be listed."*
