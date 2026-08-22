---
id: REQ-320
title: "Show and set the timeline window's start and end"
status: pending
created_at: 2026-08-22T22:08:34Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
