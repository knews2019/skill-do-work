---
id: REQ-329
title: "Land the timeline Now and Fit all buttons somewhere readable"
status: pending
created_at: 2026-08-23T12:14:00Z
user_request: UR-066
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-328]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Land the Timeline Now and Fit All Buttons Somewhere Readable

## What

`Now` sizes its window from the span between the now-line and the forecast's queue-empty instant. When the
queue is empty or nearly so that span is zero, the margin falls back to half the **zoom floor**, and the
button lands the reader on a one-hour window — the tightest zoom the view has. Zoom-in is then completely
dead, and one `›` press walks into the cosmetic bound padding where nothing is drawn at all. `Fit all`
has the mirror problem: it fits the payload's range rather than the rows on screen, so under a filter most
of the plot is blank.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

`Now` is the view's "take me to what matters" control and the one the hint sends readers to. Landing on a
one-hour window means a reader who presses it sees a handful of bars clipped flush at both frame edges with
no surrounding context, and the obvious next move — zoom in — does nothing.

## Context

`timelineNowJump` (`web/board-timeline.js:459`) computes
`marginMs = max((latest - earliest) * 0.1, TIMELINE_MIN_SPAN_MS / 2)`. With `queueEnd` at or near `now` the
first term is ~0 and the floor is 30 minutes, giving a 1-hour window. Measured on this repo's board
(now `2026-08-23 11:13`, nothing left to schedule): `Now` gives `2026-08-23 10:43 UTC → 2026-08-23 11:43 UTC`.
That is exactly `TIMELINE_MIN_SPAN_MS`, so `timelineZoomedWindow`'s floor (`:185`) refuses every zoom-in —
the `+` button, `ctrl`+wheel and the `+` key are all silent no-ops with no disabled state and no feedback.

On the user's own board the forecast span was 33 minutes (`11:08 → 11:41`), which gives the same shape.

Two adjacent states from the same family:

- **`Now` then `›`** steps a day into the render's 2% bound padding (`:1053-1056`), a stretch that exists
  only so a bar at the range edge is not flush against the frame. Nothing is drawn there, the Day chip is
  lit, and `›` is dead. The padding is cosmetic and should not be a place a period control can land.
- **`Fit all`** assigns the payload bounds verbatim (`:1943`). Under a filter matching a handful of rows this
  leaves the great majority of the plot empty. Bounds must stay the payload's — narrowing them would trap a
  zoomed-in reader, which the code comment at `:1038` correctly explains — but the *button* is asking "show
  me everything on screen", and everything on screen is the filtered set.
- **Zoom-out then zoom-in** does not return the window at either end of the range: the centre anchor plus the
  edge clamp slides it inward by about 30% of its span. Minor, same function, worth settling here.

## Detailed Requirements

1. `Now` lands on a window that shows the now-line **in context**. Floor its margin on a span a reader can
   read a day's work against rather than on half the zoom floor, so the landing window is never the zoom
   floor itself and zoom-in always has somewhere to go.
2. Both lines `Now` is responsible for — the now-line and the queue-empty rule — stay inside the window it
   lands on, which is the property the current margin exists to guarantee. Keep it.
3. `Fit all` fits the **filter-matched rows plus the projection**, not the raw payload range. The clamp
   bounds stay the payload's, so the reader can still pan and zoom outside the filtered extent.
4. A period step cannot land wholly inside the cosmetic bound padding. Either keep the padding out of the
   period controls' reach or make the step stop at the last period that contains data — state which, and why,
   where the padding is computed.
5. A control that cannot move the window says so rather than silently doing nothing: at minimum the zoom and
   step buttons carry a disabled state when the window is already at the floor, the ceiling, or the last
   period.

## Constraints

- `depends_on: [REQ-328]`. Today's `queueEnd`/`RangeEnd` are distorted by rows wrongly marked open, so
  measuring this behaviour before that lands would calibrate against the wrong numbers.
- Do not give `Now`, `Fit all` or the zoom buttons their own floor or clamp: they keep routing through
  `timelineZoomedWindow`, which stays the single place the floor, ceiling and edge rules live.
- Keep `TIMELINE_MIN_SPAN_MS` as the zoom floor. This REQ changes where `Now` lands, not how far zoom reaches.

## Red-Green Proof

**RED prompt/case:** Generate a board, open Timeline, press `Now`, read `#timeline-range-readout`, then press
`+` three times and read it again. Separately: filter to one domain, press `Fit all`, and measure what
fraction of the plot width the drawn bars span.

**Why RED now:** `Now` gives `2026-08-23 10:43 UTC → 2026-08-23 11:43 UTC` (exactly one hour); the three `+`
presses leave that string byte-identical. Filtered to `ui-design` (1 row), `Fit all` gives the full
89-day range for a row whose bar occupies a few pixels.

**GREEN when:**
- `Now` lands on a window strictly wider than `TIMELINE_MIN_SPAN_MS`, with the now-line inside it, and one
  `+` press measurably narrows it.
- The queue-empty rule is still inside the `Now` window when the forecast is confident.
- `Fit all` under a filter produces a window whose span is the filtered rows' own extent (plus the
  projection), and pressing `−` from there still widens past it.
- No `›` press from the `Now` window lands on a window containing no drawn row while data remains ahead of it.
- The zoom-in button reports itself disabled in the state where it cannot act.

**Validation:** Inferred during capture from a reproduced render; the numbers above were read from the live
DOM, not estimated.

## Full Context

See `do-work/user-requests/UR-066/input.md`.
