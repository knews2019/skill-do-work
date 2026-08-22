---
id: REQ-321
title: "Colour timeline bars by REQ status"
status: pending
created_at: 2026-08-22T22:08:34Z
user_request: UR-065
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-318, REQ-319, REQ-320, REQ-322, REQ-323, REQ-324]
batch: timeline-ux-audit
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Colour Timeline Bars by REQ Status

## What

Give every timeline bar its REQ's status colour, using the same semantic tokens the board
cards and the Calendar chips already use.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Today a blocked REQ, a pending REQ and a completed one draw identical grey-and-blue bars.
Status — the thing the rest of the board colours everything by — is invisible here unless
you open a row.

## Context

`board.css` already defines the semantic palette every other view reads: `--accent-pending`,
`--accent-claimed`, `--accent-blocked`, `--accent-done`, `--ink-faint` for cancelled, plus
the matching `--tint-*`. The Calendar's `.calendar-chip[data-status=…]` block is the
reference mapping. The Timeline instead has its own `--timeline-wait` / `--timeline-work` /
`--timeline-projected` trio.

The client already holds each row's status: `requestsById[row.id].status`, which
`renderVisibleRows` reads for the cancelled case and `timelineRowDescription` speaks aloud.
Nothing new is needed in the payload.

`calendarDayBreakdown` spells every status out rather than prefix-matching, so a typo like
`blockd-dependency-cycle` falls through to the unrecognized group instead of being counted
as real blocked work. Same rule here: exact match, and an unrecognized status takes the same
accent the Calendar gives it.

## Detailed Requirements

- **The whole bar takes the status colour** — both segments. Wait vs work is told apart by
  lightness: a pale wash for the wait, the solid accent for the work. This was the user's
  choice at capture, over colouring the work segment alone (which leaves every unclaimed REQ
  colourless) and over a separate per-row status stripe (which leaves the bar status-blind).
- Reuse the existing `--accent-*` / `--tint-*` tokens. Do not mint a second palette for the
  same statuses; a REQ must not read as one colour on a card and another on a bar.
- Map status to class through one pure function, exact-match over the full status
  vocabulary, with an explicit fallback for an unrecognized value. Keep it in step with
  `actions/work-reference.md`'s status vocabulary and `model.go`.
- Cancelled keeps its existing dimmed treatment; broken stamps keep their break marker;
  projected segments keep their hatch — a forecast must never read as measured work.
- The open-span dashed outline still has to mean "running to the now-line" once the fill is
  a status hue.
- The legend is now describing two encodings — hue is status, lightness is wait vs work.
  Rewrite it to say both. It is the view's only colour key.
- The status stays spoken as well as coloured: the row `aria-label` and the table column
  already do this, and must survive.

## Constraints

- Both themes. The dark and light blocks in `board.css` both define the accent tokens;
  check the rendered bars in each rather than trusting the token names.
- Contrast between the pale wait wash and the surface has to survive at a 10px bar height.
  Render it and look, per the prime's rule about pixels.
- Serial with the rest of the `timeline-ux-audit` batch.

## Builder Guidance

**Certainty: Firm on the encoding, exploratory on the values.** The user asked for status
colour "like on the calendar view" and picked the whole-bar encoding from three options at
capture, so which channel carries what is settled. The two lightness values that separate
wait from work are not: pick them, render a mixed-status board in both themes, and adjust
what the render shows. Scope cue: reuse the palette that exists. A REQ must read as the
same colour on its card, its calendar chip and its bar, and a second set of tokens for the
same five statuses would be the thing that later drifts.

## Red-Green Proof

**RED prompt/case:** Generate a board holding at least one REQ in each of pending, claimed,
blocked, completed and cancelled, and open the Timeline tab. Every bar is the same
grey-and-blue pair; the only rows that look different are cancelled ones (dimmed) and
broken-stamp ones (a red break marker). Nothing on the chart distinguishes a blocked REQ
from a completed one.

**Why RED now:** `drawSegment` is called with `timeline-segment-wait` /
`timeline-segment-work` and nothing else; status never reaches a class.

**GREEN when:** each bar carries its REQ's status colour from the shared `--accent-*`
tokens, matching what the same REQ's Calendar chip shows; wait and work are still
distinguishable within one bar; an unrecognized status falls to the same accent the Calendar
gives it rather than to a real status's colour; the legend states hue-is-status and
lightness-is-phase; and the screenshot of a mixed-status board shows the difference at a
glance.

**Validation:** User adjusted — the user asked for "color coded, like on the calendar view"
and chose the whole-bar encoding from the options offered at capture.

## Assets

Screenshot described in `do-work/user-requests/UR-065/input.md` — forty uniformly grey rows.

---
*Source: "4. req status should be color coded, like on the calendar view."*
