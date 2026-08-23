---
id: REQ-326
title: "Anchor and step the timeline period controls on a position, not a width"
status: pending
created_at: 2026-08-23T12:05:00Z
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

# Anchor and Step the Timeline Period Controls on a Position, Not a Width

## What

The Timeline's Day / Week / Month chips and its `‹` / `›` step arrows all compute a calendar period and
then hand it to `timelineZoomedWindow`, which **preserves the span and slides**. That is right for a zoom
and wrong for a period: a period is a position. Four separate reader-visible failures fall out of it, plus
one from the anchor the chips choose. Fix the family together — they share one cause and one code site.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

This is the request the user led with: "when I click on week, every entry disappears because it goes far in
the past instead of be sensible and show me the now." The chip is the control a reader reaches for first,
and on the default window it lands on an empty screen.

## Context

`timelinePeriodWindow` (`web/board-timeline.js:280`) builds the candidate `[periodStart, nextPeriodStart)`
and passes it to `timelineZoomedWindow` at `:286` with factor 1. When the candidate sits outside the
bounds, `timelineZoomedWindow:194-197` pins the offending edge and **drags the other edge to keep the
width** — so the window keeps a period's *length* and loses a period's *alignment*. `timelineTypedWindow`
already learned this exact lesson and clamps each endpoint before settling (`:383-391`, with the comment
explaining why); the period path never got the same treatment.

`applyPeriodWindow` (`:1980`) then anchors every press on `(windowStartMs + windowEndMs) / 2`. Once a
previous step has been slid by the clamp, that midpoint has drifted into the neighbouring period, so the
next step is computed off the wrong period.

### The five failures

- **F1 — a chip lands in the middle of the archive.** From the default Fit-all window the midpoint anchor
  is the arithmetic middle of the whole capture history. On this repo's board (317 REQs, now
  2026-08-23 11:13 UTC) Week gives `2026-07-06 00:00 → 2026-07-13 00:00`, 1 of 317 REQs; Day gives
  `2026-07-16`, 1 of 317. On the user's board it gave the week of 1 June with **795 REQs outside it**.
- **F2 — a step at a range edge steps by the wrong amount, then dies.** From the week of 17–24 Aug one `›`
  gives `2026-08-18 04:23 → 2026-08-25 04:23`: forward 1 day 4h23m, not 7 days, and the window is no longer
  a week, so the chip goes dark and the state reads "custom span". Two further `›` presses are byte-identical
  no-ops with no feedback. Month is the same shape and additionally *shrinks*: Fit all → `›` gives 31 days,
  a second `›` gives 30.
- **F3 — `›` then `‹` lands a period earlier than where the reader started.** 17–24 Aug → `›` → `‹` gives
  10–17 Aug. The forward-then-back pair moved the reader a whole week backwards and dropped the current week
  off screen. Month skips a period outright: a window starting 2026-05-28 → Month → `›` jumps from
  `2026-05-27 23:33 → 2026-06-27 23:33` to `2026-07-01 → 2026-08-01`; June is never shown.
- **F4 — a chip near a bound produces a window that is not that period.** The chip stays dark and
  `#timeline-period-state` reads "custom span" for a press the reader just made, and the month containing
  now cannot be reached as a month at all.
- **F5 — `‹` / `›` on a wide window zoom instead of stepping.** `timelineNearestPeriodLevel` (`:414`) picks
  the level closest in *span* to what is on screen, so on any window wider than about 19 days the arrows
  jump to month scale rather than stepping the window the reader is looking at.

### The judgement call this REQ has to make, and the recommendation

**What should a chip anchor on?** Recommendation: **the now-line when the current window contains it, and
the window's midpoint otherwise.** Rationale: if now is on screen it is part of what the reader is looking
at, so "Week" means the current week; if the reader has panned back to March, "Week" means a week of March
and the midpoint is the right anchor. The rule is a condition, not a special case for the default window,
and it is idempotent — pressing a chip twice never moves the window.

Stepping is **not** re-anchored on now. `‹` and `›` must step from the window's own alignment or repeated
presses would all step from the same place. Anchor a step on the period containing `windowStartMs` going
forward and the period containing `windowEndMs - 1` going back, so the two are inverses.

## Detailed Requirements

1. `timelinePeriodWindow` clamps the candidate period's endpoints into `[boundStartMs, boundEndMs]` before
   settling, the way `timelineTypedWindow` does. An edge period is **cut short**, never slid.
2. A chip press (`stepCount === 0`) anchors on `nowMs` when
   `windowStartMs <= nowMs <= windowEndMs`, and on the window midpoint otherwise.
3. A step (`stepCount !== 0`) anchors on the window's current alignment, not its midpoint: forward from the
   period containing `windowStartMs`, back from the period containing `windowEndMs - 1`. `‹` immediately
   after `›` returns to the window the reader came from, and no press skips a period.
4. `‹` / `›` step by the level the chips report when one is lit, and otherwise by a level the reader would
   recognise as "about what is on screen" — fix F5 by keying the choice on the window, not by adding a list.
5. The chip / `#timeline-period-state` readout stays honest in every resulting state. A window genuinely cut
   short by a bound is allowed to read "custom span"; a window that IS the period must light its chip.

## Constraints

- `nowMs` is already in scope at the call site (`renderTimelineView`'s `nowMs`), so no new payload field.
- Do not give the period path its own floor or clamp beyond the endpoint clamp above: the shared settle in
  `timelineZoomedWindow` stays the one place the floor and ceiling live.
- Do not touch the axis tick arithmetic — that is REQ-327's scope, and both files' diffs must stay separable.

## Red-Green Proof

**RED prompt/case:** Generate a board from this repo (`queue-kanban generate --out DIR`), open it, click
Timeline, click Week. Read `#timeline-range-readout` and `#timeline-summary`.

**Why RED now:** The readout says `2026-07-06 00:00 UTC → 2026-07-13 00:00 UTC` and the summary says
"1 REQ in the window … 316 outside the window, not listed" — a window in the middle of the archive with
effectively nothing drawn. Then: Now → Week → `›` → `‹` ends on `2026-08-10 00:00 → 2026-08-17 00:00`,
a week earlier than the Week press produced.

**GREEN when:**
- Week from the default Fit-all window lands on the calendar week containing the now-line
  (`2026-08-17 00:00 → 2026-08-24 00:00` on this board), with the Week chip lit.
- Day lands on the calendar day containing now; Month on the calendar month containing now.
- From a window that does NOT contain now (say From 2026-06-03 / to 2026-06-17), Week still lands on a week
  inside that window, not on the current one.
- `›` then `‹` from any exact period returns the identical window, at both range edges and mid-range.
- No `‹` / `›` press produces a window whose span is not the period's, except one deliberately cut short by
  a bound — and that one is reported as "custom span" rather than lighting a chip.

**Validation:** Inferred during capture from a reproduced render; the anchor rule is the capture's
recommendation and is stated above so review can overrule it in one place.

## Full Context

See `do-work/user-requests/UR-066/input.md`.
