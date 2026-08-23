---
id: REQ-346
title: "Name the UR on every Durations sample"
status: pending
created_at: 2026-08-23T22:37:52Z
user_request: UR-069
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-347, REQ-348, REQ-349, REQ-350, REQ-351, REQ-352, REQ-353, REQ-354]
batch: durations-panel-improvement
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-23T22:37:52Z
  basis:
    - Route B
    - 4-file write set
    - 2 subsystems involved
    - 7 acceptance criteria
    - browser evidence
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-durations.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/generate.go
---

# Name the UR on Every Durations Sample

## What

The Durations view is the only board view that drops the REQ to UR link, so panel A's dot cloud
cannot say which user request a mark belongs to. Give the view UR identity: in the hover readout, in
the sample table, and as a grouping lane under panel A.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The Durations payload sample is five fields wide — `id`, `route`, `completionTime`, `dayKey`,
`wallMinutes` — and none of them is the UR, even though every parsed ticket carries `UserRequestId`
and the board builds a UR node with its REQ list. On this repository's archive that costs the view
66 URs of context across 305 samples, with a median of 2 and a maximum of 10 URs active on the same
day: the dots of three or four different requests are interleaved in every dense column with nothing
to tell them apart.

## Detailed Requirements

- **Prefer a client-side join.** The client already holds it at `boardData.requests[id].userRequestId`.
  Add the field to `generatedDurationSample` in `generate.go` only if the join proves awkward — and
  say in the hand-back which way it went and why.
- **Hover readout** names the UR beside the REQ id. Keep the rest of the line as it is.
- **Sample table** ("Every sample, as a table") gains a UR column.
- **UR grouping lane** under panel A: one horizontal bracket per UR spanning its samples' completion
  times, packed into a small fixed number of rows.
- **Every UR that finds no row is counted in an explicit remainder** — never silently dropped. This
  is not the builder's call. The row count is: three rows placed 54 of 61 URs on 30 days of this
  repository's data, and the consuming project has more URs in the same window.
- **Leave route colour as it is.** REQ-347 adds the colour-by channel; this REQ does not touch fill.
- **A sample whose REQ carries no `user_request` gets an explicit unknown-UR treatment**, named on
  screen rather than left blank. Twelve of this repository's 305 samples are in that state — REQ-001
  through REQ-011 and REQ-060, all pre-dating the UR system, all with parseable claim and completion
  stamps, and `buildDurationAggregate` includes every one of them. Decide the treatment (a single
  "no UR" bucket in the lane, an explicit cell value in the table, or another rule you can defend),
  apply the same rule in all three surfaces, and keep it distinct from the remainder count above:
  a UR that found no lane row and a sample that has no UR at all are different facts.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first.
- Generate a board and look at it. The lane is a new text-bearing SVG structure at a density that
  already defeated the overflow lane's labels once — a passing suite is not evidence here.
- Seven REQs in this batch write `web/board-durations.js`. Keep the diff to the identity surfaces.
- **The target is legibility at 700 or more archived REQs**, not at this repository's 305. The
  consuming project is already at 692 samples across 47 active days — check the lane against a board
  that size, or say you could not.

## Dependencies

REQ-347 and REQ-351 both read the UR join this REQ establishes.

## Builder Guidance

**Certainty: firm.** The finding, the join path and the remainder rule are all settled. Open: whether
the join is client-side or payload-side, and how many lane rows. The overflow lane's own history is
the lesson to carry — variable-width text inside a dense lane is what stopped paying for itself, so
prefer positional brackets over per-bracket labels wherever the two compete.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repository, open Durations, hover any mark in a dense
column. The readout names a REQ and no UR; the sample table has no UR column; nothing under panel A
says which marks belong to one request.

**Why RED now:** The Durations payload carries no UR field and the renderer never consults
`boardData.requests[id].userRequestId`.

**GREEN when:** the hover readout names the UR beside the REQ id; the sample table has a UR column
whose every row is either a UR id or the stated unknown-UR value — never blank, including for the
twelve pre-UR samples; and a lane under panel A brackets each UR across its samples with a stated
count of the URs that found no row.

**Validation:** User confirmed (bundled invocation, capture answer Q4).

---
*Source: prompt A1, `ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html` (finding F1).*
