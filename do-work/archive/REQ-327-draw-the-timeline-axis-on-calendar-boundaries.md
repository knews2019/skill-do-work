---
id: REQ-327
title: "Draw the timeline axis on calendar boundaries"
status: completed
created_at: 2026-08-23T12:06:00Z
claimed_at: 2026-08-23T14:30:00Z
completed_at: 2026-08-23T15:20:00Z
commit: b2384cf
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

# Draw the Timeline Axis on Calendar Boundaries

## What

`timelineAxisTickInstants` divides the window into exactly `TIMELINE_AXIS_TICK_COUNT` equal parts, so on any
window that is not a whole multiple of that count the ticks land at arbitrary instants — and the label
formatter then prints a **date alone** for a tick that is not at midnight. The axis claims day boundaries it
does not have, and a calendar day inside the window gets no tick at all.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`. Prototyped the ladder standalone before touching the renderer and checked every window the view can reach for "date-only implies midnight", label distinctness and 4–8 ticks. Returned the chosen gap alongside the instants so the formatter reads the same number that positioned the ticks.
- [x] **[APPLY]:** `web/board-timeline.js` and `generate_test.go` as planned. Also added `rendererBracketDeclaration` to the test file: `sliceBalancedBlockAfter` balances braces, so pointing it at an array of object literals returns the first element only. `TIMELINE_YEAR_MS` deleted — the year threshold it existed for is gone.
- [x] **[UNIFY]:** Verified:
  - `web/board-timeline.js` — `node --check` clean; no reader left for `TIMELINE_YEAR_MS`; the tick walk lives in one function; the duplicate `plotSpanPx` local the queue-end caption had was folded into the one `renderAxis` already measures.
  - `generate_test.go` — `gofmt`/`go vet` clean; both changed lock-in tests state what REQ-327 changed and why the old shape passed; the one-source check counts calls from OUTSIDE the walk rather than a raw string count, because the walk legitimately calls the stepper twice.
  - `bash _dev/tests/maintainer-verify.sh` exit 0.
  - Mutations: equal-parts spacing, epoch-aligned week gaps, epoch-aligned fortnight gaps, the span-based year threshold, and gridlines deriving their own instants — five distinct failures. Two passed on the first attempt because no fixture reached those rungs; a three-month window and a vacuity guard counting both week rungs were added for exactly that.
  - Live render: the reported window now reads `1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 Jun`; month gives Mondays, day gives four-hour marks, Fit all gives fortnightly Mondays.

## Why

It is visible in the screenshot the user opened with. Their week axis reads
`1 Jun | 2 Jun | 3 Jun | 4 Jun | 5 Jun | 6 Jun | 8 Jun` — **7 June is missing**, and the two ticks either
side of the gap are the same pixel distance apart as every other pair. A reader lining a bar up against
"4 Jun" is reading a gridline that is actually at 4 Jun 12:00.

## Context

Six intervals across a 7-day week is 28 hours per tick. `timelineFormatAxisTick`
(`web/board-timeline.js:156`) drops the time when `spanMs / TIMELINE_AXIS_TICK_COUNT >= TIMELINE_DAY_MS`,
which is exactly when the ticks stop being at midnight — so the format switch and the defect trigger
together. Reproduced on this repo's board:

| Window | Ticks the axis prints | What is wrong |
|---|---|---|
| Week `2026-07-06 → 2026-07-13` | `6, 7, 8, 9, 10, 11, 13 Jul` | 12 Jul absent; ticks 28h apart, labelled as dates |
| Month `2026-07-01 → 2026-08-01` | `1, 6, 11, 16, 21, 26 Jul, 1 Aug` | ticks 5.17 days apart, labelled as dates |
| Fit all `2026-05-27 23:33 → 2026-08-25 04:23` | `27 May, 11 Jun, 26 Jun, …` | ticks at 23:33-ish, labelled as dates |

Two smaller defects live in the same function and are in scope because the fix decides them:

- The rightmost label names a day the window **excludes** (`13 Jul` on a window ending at 13 Jul 00:00),
  contradicting the `to` field, which calls 12 Jul the last day included.
- A window that crosses a calendar-year boundary but is shorter than a year gets no year on any label
  (`spanMs >= TIMELINE_YEAR_MS` is the wrong condition for "could this date be ambiguous").

The existing lock-in test (`generate_test.go:4677`
`TestJavaScriptBehaviorTimelineAxisLabelsNameTheirOwnInstant`) passes today because it only asks that each
label's numbers belong to its own instant — a date-only label on a 04:00 tick satisfies that. It also
**reimplements the tick spacing inline** rather than calling `timelineAxisTickInstants`, which is REQ-305's
lesson: a probe that reimplements the function under test cannot hold its call site. Both need fixing here.

## Detailed Requirements

1. Ticks sit on calendar boundaries. Pick the gap from a ladder of natural steps (minutes → hours → days →
   weeks → months → years) — the entry whose interval count is closest to `TIMELINE_AXIS_TICK_COUNT`, which
   becomes a **target** rather than an exact count — then align: sub-day and day steps to UTC midnight, week
   steps to Monday (the same week rule `timelinePeriodStart` uses), month steps to the 1st.
2. A date-only label is drawn **only** for a tick at UTC midnight. State it as a property, not a comment.
3. The label's time-vs-date-only switch keys on the **chosen gap**, and the year suffix on whether the
   window spans more than one calendar year — not on `spanMs` against a year.
4. `drawGridlines` keeps reading `timelineAxisTickInstants`, so a gridline can never mean a different instant
   from the label above it. The existing structural test for that stays green.
5. Rewrite the axis-label lock-in test to drive the real `timelineAxisTickInstants` and add the new property:
   every date-only label's instant has `getUTCHours() === 0 && getUTCMinutes() === 0`. Keep the existing
   distinct-labels and numbers-belong-to-the-instant assertions.
6. A tick very near either edge anchors its text inward so it is not clipped, the way the queue-end caption
   already does (`:1568`).

## Constraints

- Windows can be as short as `TIMELINE_MIN_SPAN_MS` (one hour) and as long as the archive; the ladder must
  cover both without a caller having to pick. Guard the generation loop against a pathological window.
- The degenerate zero-span window must still produce finite instants — `generate_test.go:3635` asserts it.
- Do not change `TIMELINE_AXIS_TICK_COUNT`'s value; change what it means and say so where it is declared.

## Red-Green Proof

**RED prompt/case:** Generate a board, open Timeline, press Week, read the axis labels and each tick's
instant.

**Why RED now:** The week of 6 Jul prints `6, 7, 8, 9, 10, 11, 13 Jul`. 12 Jul has no tick, and the ticks
labelled `7 Jul` … `11 Jul` are at 04:00, 08:00, 12:00, 16:00 and 20:00.

**GREEN when:** the same window prints `6, 7, 8, 9, 10, 11, 12, 13 Jul`, every one of those instants is a UTC
midnight, and the same property holds for every window the view can reach — a validated ladder gives:

| Window | Ticks after the fix |
|---|---|
| 1 h (the Now landing) | `10:50, 11:00, 11:10, 11:20, 11:30, 11:40` |
| 1 day | `00:00, 04:00, 08:00, 12:00, 16:00, 20:00, 00:00` |
| 1 week | 8 consecutive midnights |
| 1 month | the month's Mondays |
| Fit all (90 d) | biweekly Mondays, 6 ticks |
| 2 years | quarter starts, with years on the labels |

**Validation:** Inferred during capture; the ladder above was prototyped and every case checked for
"date-only implies midnight", label distinctness and 4–8 ticks per axis before this REQ was written.

## Full Context

See `do-work/user-requests/UR-066/input.md`.
