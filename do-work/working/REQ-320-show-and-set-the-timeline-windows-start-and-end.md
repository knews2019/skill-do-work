---
id: REQ-320
title: "Show and set the timeline window's start and end"
status: claimed
created_at: 2026-08-22T22:08:34Z
claimed_at: 2026-08-23T00:44:25Z
route: B
user_request: UR-065
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-319]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-318, REQ-319, REQ-321, REQ-322, REQ-323, REQ-324]
batch: timeline-ux-audit
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-23T00:45:10Z
  basis:
    - Route B
    - 4-file write set
    - 6 acceptance criteria
    - dependency depth 2
    - browser evidence
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Show and Set the Timeline Window's Start and End

## What

State the visible window's start and end instants in text, and add two date fields that set
them.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Three pure functions — parse a date field, render one, and resolve a typed
  pair into a window — with the resolution routed through `timelineZoomedWindow` like every
  other mover. A `renderRangeControls` called from `renderAll` so the readout and both fields
  track the pointer, the keyboard, the period chips, Now and Fit all. Markup in the toolbar
  as a fourth control group; the readout wraps rather than squeezing the fields.
- [x] **[APPLY]:** Confined to the declared write set. One design correction after the first
  render (D-01) — the settle preserves a span, which is wrong for a typed position.
- [x] **[UNIFY]:** `git diff --stat` reviewed; `node --check` clean; `go vet ./...` clean;
  `gofmt -l .` empty; no debug artifacts. Files verified — `web/board-timeline.js` (three
  pure functions, `renderRangeControls`, `applyTypedRange`, the `renderAll` call),
  `web/template.html` (the control group), `web/board.css` (three rules),
  `generate_test.go` (one probe added).

## Why

The window's own dates appear nowhere. Reading them means eyeballing seven axis ticks, and
the only textual state is `one week` or `custom span`. There is no way to say "show me
1 June to 15 June" — only to zoom and drag until it roughly is.

## Context

`timelineViewState` holds `windowStartMs` / `windowEndMs`, and every mover already funnels
through `timelineZoomedWindow`, which applies the one-hour floor, the range ceiling, and the
edge clamp. A typed date is a fourth mover and takes the same route: build a candidate
window, hand it to `timelineZoomedWindow` at factor 1, anchor 0. Do not give the new control
a clamp of its own.

`timelinePeriodLevelOfWindow` decides whether the period chips light up; a typed window that
happens to be exactly one calendar day should light *Day*, and the existing derive-don't-store
approach gives that for free (`REQ-235`'s lesson).

## Detailed Requirements

- A readout beside the period controls stating both endpoints as exact UTC instants to the
  minute, in the same format `timelineFormatStamp` already uses elsewhere in this view.
- The readout updates on every window move: buttons, wheel, drag, keyboard, period chips,
  *Now*, *Fit all*.
- Two date fields (`<input type="date">`, UTC) that set the window. A typed or picked start
  sets the window to begin at that UTC midnight; a typed end sets it to end at the end of
  that UTC day. Fine zoom stays where it is — the ± buttons, the wheel, the keyboard.
- Entries outside the board's range clamp to the range. An end before the start is not
  accepted silently — clamp it and let the readout show what the window actually became.
- The fields show the current window's dates, so opening the view and reading the fields
  answers "what am I looking at".
- Keep `custom span` honest: a typed window that is not exactly one calendar period still
  reads as a custom span.

## Constraints

- UTC throughout. Every other instant in this view and in the board header is UTC; a local
  date field here would be the only one and would silently disagree with the axis.
- Serial with the rest of the `timeline-ux-audit` batch.
- The toolbar already carries a legend, a five-button period group and a four-button zoom
  group, and it wraps. Placing two date fields there is a layout question — check the
  rendered result at a narrow width, do not assume.

## Dependencies

`depends_on: [REQ-319]` — the readout must state the window that REQ-319's row filter reads.
Landing this first would state a window nothing else consults.

## Open Questions

None. The picker precision was settled during capture: date fields snapping to whole UTC
days, with an exact-instant readout beside them. Date-and-time fields were rejected as
fiddly for the everyday case, and dropping time from the readout was rejected because it
would state a window that is not the one drawn once zoomed past a day.

## Builder Guidance

**Certainty: Firm.** Both halves are the user's own words — "start - end date should be
displayed, and should be selectable" — and the precision question was put to them at
capture and answered. Latitude on placement and markup: the toolbar already wraps three
groups, so where the readout and the two fields sit is a layout judgment, made against the
rendered result rather than in the abstract. Scope cue: this adds a fourth way to move the
window. It does not get its own clamp, its own floor, or its own idea of what a period is.

## Red-Green Proof

**RED prompt/case:** Open the Timeline tab and try to answer "what date does the left edge
of this chart sit on, exactly?" and then "show me 1 June to 15 June". The first is
answerable only by reading axis tick labels; the second is not answerable at all — the only
textual state is `custom span`, and no control accepts a date.

**Why RED now:** `timelineViewState` is never rendered as text, and the only movers are
relative (zoom, pan, step, fit).

**GREEN when:** the view states its window as two exact UTC instants that track every move;
typing 1 June in the start field and 15 June in the end field draws exactly that window;
a date outside the board's range clamps to the range and the readout says what it clamped
to; and pressing *Week* still lights the Week chip and updates both fields.

**Validation:** User confirmed (stated as item 3 of the request, precision settled by the
capture question).

## Assets

Screenshot described in `do-work/user-requests/UR-065/input.md` — the period group reading
`custom span` with no dates anywhere in the toolbar.

---
*Source: "3. start - end date should be displayed, and should be selectable"*

---

## Triage

**Route: B** - Medium

**Reasoning:** The behaviour is exact and the files are named, but the toolbar already carries
a legend, a five-button period group and a four-button zoom group and it wraps — where two
date fields and a readout go is a layout question the code has to answer, and the new control
has to reach `timelineZoomedWindow` by the same route every other mover uses without
acquiring a clamp of its own. Exploration first, no plan.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

`timelineViewState` already holds the window and every mover already settles through
`timelineZoomedWindow`, so the new control needed one call and no new state. Three things the
exploration decided that the REQ left open:

- **Where the readout goes.** The existing `#timeline-period-state` says "one week" or "custom
  span" and sits inside the period group. Extending it would have made one live region carry
  two unrelated facts. A fourth control group keeps the period state answering "which level"
  and the new readout answering "which instants".
- **Days in the fields, minutes in the readout.** Zoom reaches one hour and a date input cannot
  express that. Rather than hide the mismatch, the fields carry days and the readout carries
  the true instants, so a reader who has zoomed past a day sees the window they are actually
  looking at rather than the nearest date pair.
- **A field must not be overwritten mid-edit.** `renderRangeControls` runs on every window move,
  including moves the reader causes while a field has focus, so it skips the focused field.

## Decisions

- **D-01 — Clamp each typed endpoint before settling, rather than handing the raw pair to
  `timelineZoomedWindow`.** DECIDE & STATE, and the first render is what forced it. That
  function preserves the SPAN and slides, which is right for a zoom — the reader asked for a
  width — and wrong for a typed date, which is a position. Handing it an out-of-range pair
  made it pin the end to the bound and drag the start backwards to keep the width: typing
  1 July while the end field still read the board's last day produced a window starting
  30 June, and the field then redrew itself with a date nobody had typed. Each endpoint is now
  clamped into the range on its own and the shared settle still applies the floor and the edge
  rules, so the control gains no clamp of its own — it just stops asking the wrong question.
- **D-02 — A reversed pair clamps forward from the typed start; it never swaps.** An end before
  a start is a half-typed range as often as a mistake, and swapping would move the window to
  somewhere the reader did not ask for. The end is pushed to the start's own day and the
  readout then says what the window became.
- **D-03 — `timelineDateFieldToEpoch` parses the exact `YYYY-MM-DD` shape rather than calling
  `Date.parse`.** `Date.parse` accepts far more than a date field emits, and it round-trips
  31 February into 3 March rather than rejecting it. A rolled date is not the one that was
  typed, so the parse checks the round trip and returns NaN when it disagrees.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** The Timeline toolbar carries a fourth control group: two UTC date fields and
a readout stating the window's exact start and end instants to the minute. `renderRangeControls`
runs inside `renderAll`, so the readout and both fields follow the pointer, the keyboard, the
period chips, Now and Fit all, and never overwrite a field the reader is typing into. Typed
dates resolve through `timelineTypedWindow`: whole UTC days, each endpoint clamped into the
board's range on its own, the pair then settled by the same `timelineZoomedWindow` every other
mover uses so the floor and the edge rules stay in one place.

## Qualification

Passed — 4 files verified in the diff, 8 requirements traced, P-A-U confirmed.

- Mechanical (`tools/checks/qualify.sh`): OK.
- **Substantive:** three new pure functions, a new render step wired into `renderAll`, a new
  toolbar group and three CSS rules. Not whitespace.
- **Requirements traced:** readout of both endpoints to the minute → `renderRangeControls` +
  `timelineFormatStamp`; updates on every mover → the `renderAll` call, verified against six
  movers in a browser; two UTC date fields → the `timeline-range` group; start at midnight and
  end at the day's last instant → `timelineTypedWindow`, pinned by the probe; out-of-range
  clamps → D-01; reversed pair clamps rather than being accepted → D-02; fields show the
  current window → verified at first paint; `custom span` still honest → the period state is
  untouched and reads `custom span` for every typed window that is not exactly a period.
- **Flowing:** not applicable — no data-fetch path; the payload is unchanged.
- **Contamination check:** REQ-319 touched `board-timeline.js`, `generate_test.go`,
  `template.html`, `board.css`. This REQ touches the same four — expected, declared in both
  write sets, and the reason the batch is serial. No unexplained files.

## Testing

**Tests run:** `go vet ./...` and `go test -count=1 ./...` in
`skills/do-work-board/tools/queue-kanban`; `bash _dev/tests/maintainer-verify.sh` from the
repo root (`GOTOOLCHAIN=go1.26.1`, `QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium`).
**Result:** ✓ All passing — package `ok` in 50.6s; canonical gate exit 0, strict JavaScript
and strict browser lanes included.

**Red-green validation:**
- `TestJavaScriptBehaviorTimelineTypedDatesMoveTheWindow`: ✗ before implementation — against
  the pre-change renderer it fails with `anchor "function timelineDateFieldToEpoch(" not found
  in the generated page`, because no code answered "what window do these two dates mean" →
  ✓ after. This is the REQ's captured RED at the rule: the captured RED was "no control
  accepts a date", and this is why.
- The probe covers both fields, each field alone, the same date twice, a reversed pair, a date
  outside the range, two empty fields, unparseable text, `2026-02-31` (which `Date.parse` would
  roll into March), and the round trip from a mid-day instant back to its date.

**Render evidence — and it found a defect the unit probe had missed.** Generated a board from
this repo's archive (316 REQs) and drove it in headless Chromium at 1400x900 and 760x900 —
`file:///tmp/claude-0/-home-user-skill-do-work/32295d3b-538a-57cc-a4d8-2d453777559b/scratchpad/board320/probe-range.html`,
Chromium 1194 (Playwright build, `/opt/pw-browsers/chromium`), `location.href` returned with
the measurements.

| Action | Readout | Fields | Level | Rows |
|---|---|---|---|---|
| Fit all | `2026-05-27 23:44 → 2026-08-24 19:40` | 05-27 .. 08-24 | custom span | 35 |
| Type 1–15 July | `2026-07-01 00:00 → 2026-07-15 23:59` | 07-01 .. 07-15 | custom span | 10 |
| End before start | `2026-07-01 00:00 → 2026-07-01 23:59` | 07-01 .. 07-01 | custom span | 4 |
| Start 2020-01-01 | `2026-05-27 23:44 → 2026-07-01 23:59` | 05-27 .. 07-01 | custom span | 18 |
| Week chip | `2026-06-08 00:00 → 2026-06-15 00:00` | 06-08 .. 06-15 | one week | 0 |
| Zoom in twice | `2026-06-10 03:11 → 2026-06-12 20:48` | 06-10 .. 06-12 | custom span | 0 |

The first run of that table is what produced D-01: typing 1 July gave a window starting
**30 June**, because the settle preserved the span and slid the start back to keep it. The
fields then redrew with a date nobody typed. The same run showed an out-of-range start blowing
the window out to the full range instead of clamping. Both are fixed above and both now have a
case in the probe — `startAgainstCeiling` is there specifically because the unit fixture could
not see it and a browser could.

**Layout checked at both widths:** the two date fields never overlap each other
(`getBoundingClientRect()` intersection, false at 1400px and at 760px), and the toolbar's
bottom edge meets the chart's top edge exactly at both (1489/1489 and 2590/2590), so the
wrapped group pushes the chart down rather than overlapping it.

**New tests added:**
- `TestJavaScriptBehaviorTimelineTypedDatesMoveTheWindow` — the typed-window rule, thirteen
  cases including the two the render found.

**Existing tests updated (cross-REQ impact):** none.
