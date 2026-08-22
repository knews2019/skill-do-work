---
id: REQ-323
title: "Let a timeline bar be read against a date"
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
related: [REQ-318, REQ-319, REQ-320, REQ-321, REQ-322, REQ-324]
batch: timeline-ux-audit
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Let a Timeline Bar Be Read Against a Date

## What

Three additions that make a bar's position mean something: gridlines through the plot, a
drawn queue-end rule, and a minimum bar that stays readable when the window is wide.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The axis is a 26px strip above a 58vh scroll container. A bar two hundred rows down has no
vertical reference at all, and at *Fit all* over three months every completed REQ is a
1.5px sliver in which the wait and the work segment occupy the same pixel. Meanwhile the
forecast paragraph names a queue-empty instant that the chart never marks, though the
*Now* button already knows where it is.

## Context

Three separate causes, one reading problem, so one REQ:

- `renderAxis` draws ticks into `axisSvg` only; the rows SVG gets `drawNowRule` and nothing
  else.
- `drawSegment` floors width at `Math.max(1.5, …)`, and both segments are floored
  independently, so a short REQ at a wide zoom is two adjacent 1.5px marks of different
  hues.
- `timelineNowJump` already computes `queueEndMs` from `projection.queueEnd` and frames it
  in the window; nothing draws it.

`drawNowRule` is inside the rows SVG on purpose — that is what guarantees it shares the
bars' x scale. Anything added here follows the same rule.

## Detailed Requirements

- **Gridlines.** A faint vertical line at each axis tick, running the full height of the
  rows area, drawn from the same tick instants `renderAxis` computes so the two can never
  disagree. Faint enough to sit behind the bars, not compete with them.
- **Queue-end rule.** When the projection is confident and its queue-empty instant falls
  inside the window, draw a labelled vertical rule at it, visually distinct from the
  now-line, and name it in the legend. When the projection declined, or the instant is
  outside the window, draw nothing.
- **Minimum bar.** Raise the floor so a bar is visible rather than technically present.
  When a whole row's span is narrower than a readable two-segment bar, draw one marker for
  the row instead of two adjacent slivers claiming a wait/work split the pixels cannot
  show. Pick the threshold from the rendered result, and state the number and where it came
  from.

## Constraints

- Gridlines must not cost the virtualization: they are a handful of nodes per render, not
  per row.
- The now-line stays the most prominent vertical mark. A queue-end rule that outshouts it
  moves the reader's eye to a forecast instead of to the present.
- Both themes.
- Serial with the rest of the `timeline-ux-audit` batch.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repo's archive, press *Fit all*, and scroll
two hundred rows down. Ask: what date does this bar start on? There is no vertical reference
anywhere in the plot — the axis is off-screen above. Scroll back up and read the forecast
sentence: it names a queue-empty instant, and no mark on the chart corresponds to it. Look
at any completed REQ at this zoom: a 1.5px grey tick with no distinguishable wait or work.

**Why RED now:** Ticks are drawn only into the axis SVG; nothing draws `projection.queueEnd`;
`drawSegment` floors each segment independently at 1.5px.

**GREEN when:** gridlines run from each axis tick through the rows area at every zoom
level; a labelled queue-end rule appears when the projection is confident and its instant is
in the window, and is absent otherwise; a sub-threshold row draws one legible marker rather
than two 1.5px slivers; and the legend names every vertical rule and marker the chart can
draw.

**Validation:** Inferred during capture — audit findings, not one of the user's four items.

## Assets

Screenshot described in `do-work/user-requests/UR-065/input.md` — thirty-odd rows of
one-to-three-pixel ticks under an axis they cannot be measured against.

---
*Source: audit finding, UR-065 — "audit the timeline view, and make it more useful UIUX."*
