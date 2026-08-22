---
id: REQ-319
title: "List only the REQs the selected window covers"
status: claimed
created_at: 2026-08-22T22:08:34Z
claimed_at: 2026-08-22T23:23:42Z
route: B
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
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-22T23:24:30Z
  basis:
    - Route B
    - 3-file write set
    - 6 acceptance criteria
    - dependency depth 1
    - browser evidence
    - cross-route regression gates
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# List Only the REQs the Selected Window Covers

## What

Drop a REQ's row from the Timeline entirely when its span falls outside the visible time
window, instead of listing it as an empty row with a clipped-to-nothing bar.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Split the one `rows` array into two: `filterMatchedRows` (fixed per render)
  and `rows` (the window subset, re-derived by `refreshWindowRows` inside `renderAll`, which
  every window-mover already funnels through). Membership as two pure functions —
  `timelineRowExtent` for a row's drawn span, `timelineRowsInWindow` for the overlap test —
  with extents computed once per render so a drag compares numbers, not ISO strings. Move
  the summary and the scroll extent into the render path; build the table on demand.
- [x] **[APPLY]:** Confined to the declared write set. One signature change outside the plan
  (`timelineNowJump`) with its probe updated in the same edit — recorded as D-01.
- [x] **[UNIFY]:** `git diff --stat` reviewed; `node --check` clean on the renderer;
  `go vet ./...` clean; `gofmt -l .` empty; no debug artifacts. Files verified —
  `web/board-timeline.js` (two new pure functions, the two-row-set split, `renderSummary`,
  the on-demand table, the Now handler, the forecast note), `generate_test.go` (one probe
  updated, one added), `web/template.html` (untouched in the end — see D-02).

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
- Generate a board and step through Day, Week and Month before hand-back. Whether the rows
  on screen match the window is a claim about pixels, and a green suite is not evidence
  about it (`_dev/primes/prime-kanban-board.md`).

## Dependencies

`depends_on: [REQ-318]` — that REQ rewrites the same row-derivation and the same subhead
sentence. Ordering, not logic: run 318 first and this one has no merge to resolve.

## Builder Guidance

**Certainty: Firm** on the behaviour, **exploratory** on the wording. The requirement is the
user's own words — "Reqs that are not in the period selected should not be listed." The
intersect-not-contain membership rule is capture's reading of "not in the period" and is
stated in the requirements above; build that reading, and if the code makes a different one
obviously right, say so in the hand-back rather than building both. The empty-window
sentence is yours to write. Scope cue: this windows what is *listed*; it does not touch the
projection, the medians, or the forecast paragraph.

## Open Questions

- [x] Does the block under the forecast — the one naming REQs that cannot be given an honest
  start time — shrink to the selected period too? → **No, it stays whole-queue.** The user
  answered at verify (2026-08-22): that block answers "what is stuck", not "what happened
  this week", and a REQ blocked since May is exactly the one that must not vanish because
  the window is on a day in June. This preserves `REQ-305`'s rule that the forecast half of
  this view describes the whole queue.

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

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is exact and the file is named, but *where* the change goes is
not: the row set is derived once per render and consumed by five things (the virtualized
scroll extent, the subhead counts, the `<details>` table, the forecast's filters-active
flag, and every window-mover's redraw path), and which of those re-derive on a window move
has to be found rather than assumed. Exploration first, no plan.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

`rows` was a single `var` computed once at the top of `renderTimelineView` and closed over by
eleven consumers. Which of them re-ran on a window move decided the whole shape of the change:

| Consumer | Ran on a window move before? | After |
|---|---|---|
| `renderVisibleRows` (the bars) | yes, via `renderAll` | yes |
| `renderAxis`, `renderPeriodControls` | yes, via `renderAll` | yes |
| the subhead counts | **no** — written once per render | moved into `renderAll` as `renderSummary` |
| `rowsSvg` height (the scroll extent) | **no** — set once at node creation | set in `renderVisibleRows` |
| the forecast's subset note | **no** | recomputed, DOM rebuilt only when the answer flips |
| the `<details>` table | **no** — built once, unvirtualized | built on open, rebuilt at most once per frame |
| `timelineNowJump`'s `scrollTop` | n/a (button only) | moved after the row refresh — see D-01 |
| the `mousemove` readout's `data-row-index` | per render | unchanged: it indexes the array the same render drew |

Three things the exploration settled that the REQ did not:

- **Bounds must not follow the window.** `boundStartMs`/`boundEndMs` come from the payload's
  range. Deriving them from the windowed set instead would drag the clamp floor up behind
  every zoom and the reader could never zoom back out. The `rows[0].createdTime` fallback
  reads `filterMatchedRows` for the same reason.
- **The empty-window state cannot be the existing early return.** That return fires before
  the axis and the toolbar exist, which is fine for "no REQ matches the filters" (there is
  nothing to build) but wrong for a window the reader zoomed into: they need the controls
  that get them back out.
- **`timelineFirstOpenRowIndex` and the projection chain are order- and set-independent**,
  so nothing else had to change.

## Decisions

- **D-01 — Split `timelineNowJump` in two.** DECIDE & STATE. It returned `{window, scrollTop}`,
  computing the scroll from the row set the caller held *before* the jump. Under window-scoped
  rows that set is exactly the one the jump replaces, so the scrollTop indexed into a stale
  list. It now returns the window alone and the handler runs the three steps in the only order
  that can be right: apply the window, refresh the rows, then ask `timelineFirstOpenRowIndex`
  where to scroll among *those*. The existing probe in `generate_test.go` was updated in the
  same edit to drive the same three steps. Reversible; no behavior lost.
- **D-02 — `web/template.html` was in the write set and ended up untouched.** The empty-window
  message is generated text, not markup, and the hint paragraph had already been corrected by
  REQ-318. Declaring a file and not needing it is the honest outcome; nothing was forced into
  it to justify the declaration.
- **D-03 — The forecast's whole-queue note names no cause.** It opened "Filters are on", which
  was true while a filter chip was the only thing that could shrink the row set. The window
  now shrinks it too — usually much harder — and both can be on at once, so the note said
  something false about the common new case. Rewritten to "This covers the whole queue, not
  the rows shown." Naming which cause is on buys the reader nothing they need and costs a
  branch that can be wrong. The parameter is renamed `showingSubset` to match.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** The Timeline's row list is now the subset of the filter-matched rows whose
drawn span overlaps the visible time window, re-derived on every window move through
`renderAll`. Membership is two pure functions: `timelineRowExtent` takes a real min/max over
every instant a row draws at — created, claimed, completed, the now-line for an open span, and
an attached forecast bar — so a reversed stamp still has a findable extent; `timelineRowsInWindow`
keeps the rows whose extent overlaps the window. Extents are computed once per render, so a
drag compares numbers. The subhead, the scroll extent, the `<details>` table and the forecast's
whole-queue note all follow that set, the reader's place in the list is anchored across the
move, and a window containing nothing says so in its own words instead of showing a blank chart.

## Qualification

Passed — 2 files verified in the diff, 8 requirements traced, P-A-U confirmed.

- Mechanical (`tools/checks/qualify.sh`): OK.
- **Substantive:** +435/−90 across two files; two new pure functions, a restructured render
  path, one probe rewritten and one added. Not whitespace.
- **Requirements traced:** membership rule → `timelineRowExtent` + `timelineRowsInWindow`;
  pure and probeable → `TestJavaScriptBehaviorTimelineRowsFollowTheWindow`; re-derived on every
  mover → `refreshWindowRows` at the top of `renderAll`, which every mover already calls;
  subhead states the window and says it is one → `renderSummary`; table follows → on-demand
  `renderTimelineTable`; empty window in its own words → `renderSummary`'s zero branch;
  forecast note still reads correctly → D-03.
- **Flowing:** not applicable — no data-fetch path; the payload is unchanged and the
  projection is untouched.
- **Contamination check (Step 10):** REQ-318 touched `timeline.go`, `timeline_test.go`,
  `board-timeline.js`, `template.html`, `generate_test.go`. This REQ touches
  `board-timeline.js` and `generate_test.go` — expected overlap, declared in both write sets,
  and the batch is deliberately serial for exactly this reason. No unexplained files.

## Testing

**Tests run:** `go vet ./...` and `go test -count=1 ./...` in
`skills/do-work-board/tools/queue-kanban`; `bash _dev/tests/maintainer-verify.sh` from the
repo root (`GOTOOLCHAIN=go1.26.1`, `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium`).
**Result:** ✓ All passing — package `ok` in 48.4s; canonical gate exit 0, strict JavaScript
and strict browser lanes included.

**Red-green validation:**
- `TestJavaScriptBehaviorTimelineRowsFollowTheWindow`: ✗ before implementation — run against
  the pre-change renderer (`git checkout HEAD -- web/board-timeline.js`) it fails with
  `anchor "function timelineRowExtent(" not found in the generated page`, because the
  membership rule did not exist → ✓ after. This is the REQ's captured RED expressed at the
  rule rather than at the DOM: the captured RED was "press Day and 309 rows are still
  listed", and the reason was that no code answered "is this row in the window".
- The probe's own fixture is the four ways a row can meet a window — before, inside,
  straddling, after — plus a reversed stamp, plus a pending REQ whose only presence in a
  future window is its forecast bar, with a no-forecast control proving the projection is
  what put it there. A containment rule passes three of those and fails three.

**Render evidence (the captured RED, reproduced end to end):** generated a board from this
repo's archive (316 REQs) and drove it in headless Chromium —
`file:///tmp/claude-0/-home-user-skill-do-work/32295d3b-538a-57cc-a4d8-2d453777559b/scratchpad/board319/probe-window.html`,
Chromium 1194 (Playwright build, `/opt/pw-browsers/chromium`), `location.href` returned with
the measurements.

| Window | Rows drawn | Scroll extent | Table | Subhead |
|---|---|---|---|---|
| Fit all | 316 | 5688px | — | "316 REQs in the window 2026-05-27 23:45 UTC → 2026-08-24 18:42 UTC" |
| Month | 52 | 936px | — | "52 REQs in the window 2026-07-01 → 2026-08-01" |
| Day | 1 | 18px | — | "1 REQ in the window … 315 outside the window, not listed." |
| Day, stepped back 40 | 0 | 0px | — | "No REQ was open between 2026-06-06 and 2026-06-07. Widen the window, step to another period, or press Fit all — 316 REQs are outside it." |
| Fit all again | 316 | 5688px | — | back to all 316 — widening never drops a row |
| Day, table open | 1 | 18px | **1 row** | table followed the window |

The scroll extent tracks exactly (52 × 18 = 936, 1 × 18 = 18), which is the assertion that
the reader is no longer scrolling past hundreds of empty rows. The Day window's single row is
`REQ-018`, a REQ whose span straddles that day — the overlap rule working on real data rather
than on a fixture.

**Forecast note checked separately** on the same board: with Fit all the paragraph carries no
subset clause; with a Day window it leads "This covers the whole queue, not the rows shown."
That check is what caught D-03 — the old wording said "Filters are on" with no filter set.

**New tests added:**
- `TestJavaScriptBehaviorTimelineRowsFollowTheWindow` — the membership rule, including the
  straddle case, the still-running case, a reversed stamp, and the forecast-only case with
  its control.

**Existing tests updated (cross-REQ impact):**
- `generate_test.go`'s `TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow`
  (from REQ-235) — drives the Now button's three steps in the new order instead of reading a
  `scrollTop` that `timelineNowJump` no longer returns. **Deliberate, per D-01**; the
  assertion it makes about where the button scrolls is unchanged.
