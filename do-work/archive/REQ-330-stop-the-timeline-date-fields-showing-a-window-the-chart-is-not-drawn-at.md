---
id: REQ-330
title: "Stop the timeline date fields showing a window the chart is not drawn at"
status: completed
created_at: 2026-08-23T12:15:00Z
claimed_at: 2026-08-23T15:25:00Z
completed_at: 2026-08-23T15:55:00Z
commit: 22dc197
route: B
user_request: UR-066
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Stop the Timeline Date Fields Showing a Window the Chart Is Not Drawn At

## What

The From / to fields have one guard that decides whether to write a canonical value back into a field:
`syncRangeField` skips the write when the field is focused **and** its value differs from the last value this
code wrote. Three reader-visible failures come out of that one guard, including a field that goes
permanently blank in a branch whose own comment says "Restore it unconditionally".

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`. Worked the guard's truth table by hand first and confirmed F1 from it before touching code: with the attribute removed and the field focused, `"" !== null` is true, so the restore branch returns. Replaced the value comparison with the events that bound an edit; kept `timelineEndEpochToDateField` and `timelineDateFieldToEpoch` separate per REQ-320.
- [x] **[APPLY]:** `web/board-timeline.js` and `generate_test.go` as planned, plus `timeline_browser_probe_test.go` for the lifecycle probe — the field states only exist under a real engine, so a Node probe could not have reached them.
- [x] **[UNIFY]:** Verified:
  - `web/board-timeline.js` — `node --check` clean; `data-synced-value` has no remaining reader or writer; the clamp reuses `timelinePeriodStart(_, "day")` rather than adding a day-start function.
  - `generate_test.go` / `timeline_browser_probe_test.go` — `gofmt`/`go vet` clean. One trap hit and fixed: a backtick inside a JS comment terminated the Go raw string the probe lives in.
  - `bash _dev/tests/maintainer-verify.sh` exit 0.
  - Mutations: the old value-comparison guard, a change handler that does not end the edit, the removed day clamp, and a blur that does not release the edit — four distinct failures. The round-trip property (chip → fields → re-apply) still holds.

## Why

The fields are the only control that names the window in words the reader chose. A field showing a date the
chart is not drawn at is worse than no field: committing it later silently moves the window to a place
nobody asked for.

## Context

```js
function syncRangeField(field, text) {
  var readerHasEdited =
    document.activeElement === field && field.value !== field.getAttribute("data-synced-value");
  if (readerHasEdited) { return; }
  field.value = text;
  field.setAttribute("data-synced-value", text);
}
```

- **F1 — a cleared field stays blank forever.** `applyTypedRange`'s unparseable branch
  (`web/board-timeline.js:1686-1694`) calls `changedField.removeAttribute("data-synced-value")` and then
  `renderRangeControls()`. Inside `syncRangeField` the attribute is now `null` and the field's value is `""`,
  so `"" !== null` is true, the field still has focus, and the guard **returns without restoring**. The
  `removeAttribute` was meant to turn the guard off and turns it on. The comment three lines above states the
  opposite of what the code does.
- **F2 — a committed field freezes on the rejected text.** After a successful commit the render writes the
  canonical value back, but the field is still focused and still holds what the reader typed, so the guard
  skips it. Type a date that gets clamped and the field goes on displaying the date the chart is *not* drawn
  at, indefinitely.
- **F3 — a From date past the end of the range collapses the chart to an empty one-hour sliver.**
  `timelineTypedWindow` (`:381-391`) clamps each endpoint to the bounds; both collapse onto `boundEndMs`, and
  the shared settle then applies the zoom floor and slides, giving `[boundEnd − 1h, boundEnd]`. Measured:
  From `2026-09-30` gives `2026-08-25 03:23 UTC → 2026-08-25 04:23 UTC`, "Nothing was drawn…", and the field
  still reads `2026-09-30`. A date past the end should clamp to the last day the board has, not to an empty
  hour behind the frame.

## Detailed Requirements

1. Distinguish "focused" from "mid-edit" by what the reader has actually typed, not by an attribute that a
   sibling branch removes. State the invariant once, and make every branch that wants a restore able to get
   one — including the cleared-field branch, whose comment already says it must.
2. A field is never left blank. Clearing it is not a request to move the window, so the window stands and the
   field is restored to the window it is drawn at.
3. After a commit that was clamped, the field shows the window the chart is actually drawn at. The reader
   sees where they landed instead of what they asked for.
4. A typed date outside the range clamps to the nearest day the board **has**, keeping the other endpoint's
   relationship intact, rather than collapsing both endpoints onto the same bound.
5. The round trip stays exact for the windows the period chips produce: chip → fields → re-apply gives the
   identical window with the chip still lit. This already holds and must keep holding.

## Constraints

- The end field is inclusive and the window's end is exclusive. `timelineEndEpochToDateField` and
  `timelineDateFieldToEpoch` are deliberately not inverses of one function (REQ-320's lesson); do not merge
  them while fixing the guard.
- Only the field the reader changed is applied (`applyTypedRange`), so editing one endpoint never moves the
  other. Keep that.

## Red-Green Proof

**RED prompt/case:** Open Timeline. (a) Click into From, select all, delete, press Tab. (b) Type
`2026-09-30` into From and press Tab. Read the field and `#timeline-range-readout` after each.

**Why RED now:** (a) the From field is empty and stays empty through every later window move. (b) the From
field reads `2026-09-30` while the readout reads `2026-08-25 03:23 UTC → 2026-08-25 04:23 UTC` and the chart
is empty.

**GREEN when:**
- (a) leaves the window unchanged and the From field showing that window's start date.
- (b) puts the window on the last day the board has, and the From field shows that day.
- After any clamped commit, the two fields and `#timeline-range-readout` describe the same window.
- A Node behavior probe drives `syncRangeField` and `applyTypedRange`'s restore branch directly and fails if
  either regresses — the cleared-field case in both directions.

**Validation:** Inferred during capture; F1 was confirmed by reading the guard's own truth table and F2/F3
by driving a real render.

## Full Context

See `do-work/user-requests/UR-066/input.md`.
