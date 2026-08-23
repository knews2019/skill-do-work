---
id: REQ-347
title: "Retire the Durations in-lane labels for a ranked longest-spans list"
status: pending
created_at: 2026-08-23T22:37:52Z
user_request: UR-068
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-342]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-342, REQ-343, REQ-344, REQ-345, REQ-346, REQ-348, REQ-349, REQ-350]
batch: durations-panel-improvement
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-08-23T22:37:52Z
  basis:
    - Route B
    - 6-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - dependency depth 1
    - browser evidence
    - cross-route regression gates
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-durations.js
  - skills/do-work-board/tools/queue-kanban/durations.go
  - skills/do-work-board/tools/queue-kanban/durations_test.go
  - skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board.css
---

# Retire the Durations In-Lane Labels

## What

Replace the overflow lane's direct labels with a compact ranked longest-spans list beside the chart,
and delete the label planner, its width model and the tests that exist only to place text inside the
SVG. Keep the lane, its marks, their hover and a plain remainder count.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The consuming project's board shows five labels beside "+55 more over 60 min"; this repository's
shows four beside "+9 more". That is a 7 percent labelling rate, and every sample it labels is also
one row of the table below. Behind those four labels sit browser-measured text widths, a greedy
two-row packer, a remainder reserve that runs to a fixed point (with an iteration cap that has to
live inside the function because the browser probe slices named functions out of the built client),
`durations.go`'s label-text and width model including its documented hyphen-versus-minus divergence,
and the geometry-agreement tests holding two files' constants together. It is the largest single
piece of the view.

**This deliberately reverses the separated-text-band direction chosen on 2026-08-18.** That fix
worked; the 2026-08-23 board shows clean label rows and no mark-over-label collision. The reason to
revisit it is the sample count, not a defect. The user gave an explicit yes to the reversal during
capture.

## Detailed Requirements

- **A compact ranked longest-spans list beside the chart**, longest first, carrying REQ id, UR id,
  duration, route and title. **Every over-ceiling sample is present** — this is the property the lane
  never had, and it is the reason the list is HTML beside the chart rather than text inside the SVG:
  no width measurement, no packing, no remainder sentence.
- **Keep:** the overflow lane itself, its marks, their hover, and a plain remainder count as a
  sentence.
- **Delete:** the label planner (the greedy row packer and its per-row occupancy intervals, the
  two-pass remainder reserve and its fixed-point loop with the in-function iteration cap), the
  browser text measurement that exists only to size labels, `durations.go`'s label-text and width
  model including the hyphen-versus-minus divergence it carries, and the geometry-agreement tests
  pinning the Go constants against the renderer's.
- **Name in the hand-back exactly which tests were removed and why each no longer describes a shipped
  rule.** A deleted test is a deleted claim; the hand-back is where that is accounted for.
- The list carries the UR, so REQ-342's join is a prerequisite.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first.
- **The lane is not the problem** — putting variable-width text inside it at this density is. Do not
  remove the lane or its marks while removing their labels.
- Deleting a test is only in scope when it exists to place text inside the SVG. A test that pins the
  lane, the marks, the hover or the remainder count stays.
- Generate a board and look at it, on this repository's archive and on a board with more over-ceiling
  samples if one is reachable.
- **The target is legibility at 700 or more archived REQs**, not at this repository's 305. The
  consuming project is already at 692 samples with 60 over-ceiling — a list that only reads well at
  13 rows has not solved the problem the lane had.

## Dependencies

`depends_on: REQ-342` — the list's UR column needs the join.

## Builder Guidance

**Certainty: firm, and deliberately reversing a settled decision.** Do not treat the 2026-08-18 fix
as a mistake in the hand-back: it did what it was asked to do, and this REQ retires the surface it
was fixing. Where the list lives relative to the chart (beside, below, in a details element) is
yours; that it is complete is not.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repository and open Durations. The overflow lane
labels four of thirteen over-ceiling samples and says "+9 more". The other nine are reachable only by
opening the collapsed sample table and scanning 305 rows for spans over 60 minutes.

**Why RED now:** Over-ceiling identity is carried by in-SVG labels whose count is bounded by measured
text width and lane geometry, not by the number of samples.

**GREEN when:** every over-ceiling sample appears in a ranked list beside the chart with its REQ, UR,
duration, route and title; the lane still draws its marks with their hover and a plain remainder
count; and the label planner, the width model and the geometry-agreement tests that existed only to
place text in the SVG are gone, each named in the hand-back.

**Validation:** User confirmed — explicit yes to the reversal during capture.

---
*Source: prompt A5, `ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html` (finding F5).*
