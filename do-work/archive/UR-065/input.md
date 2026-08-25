---
id: UR-065
title: Audit the timeline view and make it more useful
created_at: 2026-08-22T22:08:34Z
requests: [REQ-318, REQ-319, REQ-320, REQ-321, REQ-322, REQ-323, REQ-324]
word_count: 52
---

# Audit the Timeline View and Make It More Useful

## Summary

The user opened the Kanban board's Timeline view on a 309-REQ queue, found it hard to
use, and asked for a full UI/UX audit. They named four changes themselves and explicitly
said their list is a starting point, not the whole scope: "let me start, but make sure to
make a full audit."

The capture is the four named items plus three findings from the audit that follow. Every
REQ targets `skills/do-work-board/tools/queue-kanban/`.

## Extracted Requests

| REQ | Source | What |
|---|---|---|
| REQ-318 | user item 1 | Newest REQ at the top instead of oldest |
| REQ-319 | user item 2 | List only the REQs the selected window covers |
| REQ-320 | user item 3 | Show the window's start and end dates, and let them be set |
| REQ-321 | user item 4 | Colour the bars by REQ status, as the Calendar view does |
| REQ-322 | audit | The row shows only a bare id; the title is 700px away at the page foot |
| REQ-323 | audit | No gridlines, no queue-end mark, and bars that collapse to 1.5px |
| REQ-324 | audit | A one-pixel pointer jiggle pans the view and eats the click |

## Audit Findings

The view was read end to end: `timeline.go` (the aggregate and the forward projection),
`web/board-timeline.js` (the whole client), the Timeline block of `web/template.html`, and
the timeline section of `web/board.css`. What follows is everything the audit turned up,
including the items deliberately **not** captured.

**Captured as REQs**

- **A row's identity is a bare id.** The 104px label column holds `REQ-012` and nothing
  else. The title exists in two places: a one-line readout at the very foot of the panel,
  below a five-sentence hint paragraph, and the collapsed `<details>` table. Neither is
  near the pointer. Scanning 300 rows tells the reader nothing about what any of them are.
  → REQ-322.
- **A bar cannot be read against a date.** Axis ticks live in a separate 26px SVG above the
  scroll container; a bar 200 rows down has no vertical reference at all. At *Fit all* over
  three months, every completed REQ is `Math.max(1.5, …)` — a 1.5px sliver in which the
  wait and the work segment are the same pixel. And the forecast paragraph names a
  queue-empty instant that the chart never marks, though the *Now* button already knows
  where it is. → REQ-323.
- **A pointer jiggle costs the click.** `pointerdown` on a row arms a pan with no movement
  threshold, so the first `pointermove` — one pixel is enough — pans the window and
  rebuilds every row node mid-click. The delegated `[data-detail-kind]` click handler then
  has no surviving trigger to fire on, so "Click a row for its full detail" works only for
  a perfectly still press. Each of those renders also re-reads `scrollHost.clientWidth`
  once per `xOfEpoch` call. → REQ-324.

**Checked and not captured**

- *Axis and bars misaligned by the scrollbar.* Not a defect. The axis SVG is wider than the
  scroll container's content box, but both draw from `plotWidth()`, which reads
  `scrollHost.clientWidth`, so tick and bar x-coordinates agree. The axis simply leaves the
  scrollbar strip blank.
- *The persisted zoom window going stale against new data.* Not reachable. `timelineViewState`
  survives a tab switch by design, and the payload's `rangeStart`/`rangeEnd` never change
  inside one page life — the board has no in-place refresh, so a new payload always arrives
  with a fresh page and a fresh window.
- *Every visible row is a tab stop.* Real, but rows are focusable by deliberate design
  (`do-work/archive/UR-051/REQ-239-give-timeline-rows-a-real-focus-ring.md`,
  `REQ-233-give-the-timeline-a-keyboard-path-to-zoom-and-pan.md`), and a roving tabindex is a
  refinement nobody has complained about. Left alone.
- *The `<details>` table is built in full on every render even while closed.* Real but
  negligible — it runs once per view switch or filter change, not per window move. REQ-319
  already rewrites which rows that table lists; building it on open can ride along there.

## Batch Constraints

- **One surface, so one at a time.** All seven REQs write
  `skills/do-work-board/tools/queue-kanban/web/board-timeline.js`. The `write_set` on each is
  honest about that overlap; the board's overlaps badge will light up and that is correct.
  Do not fan these out to parallel builders.
- **Prose that a REQ invalidates is that REQ's to fix.** The subhead sentence
  ("… in capture order, oldest at the top …"), the hint paragraph, and the legend all
  describe behaviour these REQs change. Whichever REQ falsifies a sentence updates it.
- **The screenshot is the test.** `_dev/primes/prime-kanban-board.md` — generate a board and
  look at it; a passing suite is not evidence about pixels. Record the browser and build
  beside any measured number.
- **A removed rule takes its tests with it.** REQ-318 inverts an ordering that
  `timeline_test.go` currently pins. Say so loudly in the hand-back; a quietly edited test
  looks identical in a diff (`REQ-237`'s lesson).

## Answered During Capture

- **Status colour encoding** — the user chose: the whole bar takes the REQ's status colour,
  and wait vs work is told apart by lightness rather than by hue. Rejected: colouring only
  the work segment (leaves every unclaimed REQ colourless), and a separate status stripe
  per row (leaves the bar itself status-blind).
- **Date picker precision** — the user chose: plain date fields that snap the window to
  whole UTC days, alongside a readout that always states the exact start and end instant to
  the minute. Rejected: date-and-time fields (fiddly for the common case), and day
  granularity everywhere (would state a window that is not the one drawn once zoomed
  past a day).

## Answered During Verify (2026-08-22)

- **Does the excluded list follow the window?** No. `do-work verify-requests` raised item 2's
  reach as ambiguous: REQ-319 windows the chart and the table, and the block under the
  forecast naming REQs that cannot be given an honest start time was left whole-queue. The
  user confirmed whole-queue — that block answers "what is stuck", and a REQ blocked since
  May must not vanish because the window sits on a day in June. Recorded on REQ-319 as an
  answered Open Question; preserves `REQ-305`'s rule that the forecast half of this view
  describes the whole queue.

## Screenshot

The user attached one screenshot of the Timeline view. It was not persisted as an asset —
the image was available only in the conversation, not as a file this capture could copy.
Description of what it shows:

A browser at `http://127.0.0.1:8090` on the board for a project labelled `skill-do-work2`,
generated `2026-08-21 21:49 UTC`, with the Timeline tab active among Board / Calendar /
Durations / Timeline / Testing. Heading: "How long each REQ waited, and how long it took".
Subhead: "309 REQs in capture order, oldest at the top. 25 still open, measured to the
now-line at 2026-08-21 21:49 UTC." The legend shows the five swatches (waiting, being
worked on, still open, projected, broken stamps); the period group reads "custom span"; the
zoom group shows − + Now "Fit all". The axis spans 28 May to 23 Aug. Roughly forty rows are
visible, labelled `REQ-001` through `REQ-042` in ascending order, each carrying one or two
grey slivers one to three pixels wide; a handful of longer grey bars sit at REQ-001, 006,
013, 014, 018 and 041–042. A dotted red now-line stands near the right edge at 23 Aug. The
forecast line reads "Nothing left to schedule — every remaining REQ is listed below." The
hovered row's readout at the foot reads "REQ-012 · Add do-work note command for lightweight
roadmap notes · Route B · completed · waited 17h 39m · worked 0 min". Below it, the
collapsed "Every REQ, as a table" disclosure.

The screenshot is the evidence for the four named items: oldest work fills the screen, the
rows outnumber the visible bars by an order of magnitude, the window's own dates appear
nowhere in text, and every bar is the same grey.

## Full Verbatim Input

```
do-work capture-request

audit the timeline view, and make it more useful UIUX.
let me start, but make sure to make a full audit.

Here we go:

1.  most recent REQ's should be on top.

2. Reqs that are not in the period selected should not be listed.
3. start - end date should be displayed, and should be selectable
4. req status should be color coded, like on the calendar view.
```

(Plus one attached screenshot of the Timeline view, described above.)

---
*Captured: 2026-08-22T22:08:34Z*
