---
id: REQ-346
title: "Name the UR on every Durations sample"
status: completed
completed_at: 2026-08-24T12:45:00Z
claimed_at: 2026-08-24T10:05:00Z
created_at: 2026-08-23T22:37:52Z
user_request: UR-069
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
route: B
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
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Name the UR on Every Durations Sample

## What

The Durations view is the only board view that drops the REQ to UR link, so panel A's dot cloud
cannot say which user request a mark belongs to. Give the view UR identity: in the hover readout, in
the sample table, and as a grouping lane under panel A.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `prime-kanban-board.md` and UR-069's Batch Constraints. Chose the client-side join so one fact keeps one definition, and sized the lane against the 700-sample target rather than this repository.
- [x] **[APPLY]:** Four files. `generate.go` was in the declared write set and deliberately untouched; `generate_test.go` was not and was added — see the scope note below.
- [x] **[UNIFY]:** Audited by the orchestrator against `0fa6360..f1a1ca3`, then independently re-rendered: 65 brackets, `+2 URs with no free row`, 311 table rows, 0 blank UR cells, exactly 12 reading `no UR recorded`. The reviewer re-rendered on three boards including a 692-sample synthetic and confirmed 0 overlaps at target density.

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

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome was fully specified but the lane was a new SVG structure whose density behaviour had to be discovered by rendering, not derived.

**Planning:** Not required.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modify) — join, hover, lane, unknown bucket
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify) — the UR column header
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — bracket rules
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — the behavioural lock-in

**Files I will NOT touch:** `generate.go` — declared in the capture-seeded write set for the payload-side join, which proved unnecessary.

**Scope note.** `generate_test.go` was **not** in the capture-seeded write set, and this REQ is `tdd: true`. The builder flagged it as a stop-and-report and proceeded rather than hand back an untested rendering change. The orchestrator upheld that and has extended `write_set` accordingly, per `actions/work.md` § "Write only inside the declared scope". The generalisable half — a `tdd: true` REQ whose write set names no test file is uncompletable as written — is captured as REQ-365.

**Acceptance criteria (restated from REQ):**
- [x] Hover names the UR beside the REQ id
- [x] Sample table gains a UR column, populated for every row
- [x] A grouping lane brackets each UR, packed into a fixed number of rows
- [x] Every UR that finds no row is counted in an explicit remainder
- [x] A sample with no `user_request` gets a named unknown treatment, distinct from the remainder
- [x] Route colour untouched

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** `durationSampleUserRequestId` reads `boardData.requests[id].userRequestId` and `durationUserRequestName` turns an empty result into the single constant `DURATIONS_UNKNOWN_USER_REQUEST_NAME = "no UR recorded"`, so all three surfaces share one definition and the unknown rule cannot diverge between them. Brackets pack left-edge-first into six rows; whatever finds no row is counted right-aligned on the lane's title line. Unknown-UR samples get a reserved row below the packed six, dashed, gutter-labelled `no UR`, which never enters the packer.

## Testing

**Tests run:** `GOTOOLCHAIN=go1.26.1 go test -count=1 ./...`; the full gate with the browser lane run rather than skipped
**Result:** ✓ Gate exit 0, `TestMaintainerStrictBrowserBehaviorLane` PASS

**Rendered evidence (orchestrator, independently):** 65 brackets, `+2 URs with no free row`, 311 rows, 0 blank UR cells, 12 `no UR recorded`.

**At target density (reviewer):** a 692-sample/140-UR/47-active-day board gives 110 brackets + 1 unknown + `+30 URs with no free row` = 140, **0 overlaps**, min gap exactly 2.0 units, min width exactly 3.0, extent inside the plot bounds, 692 rows, 0 blanks.

**Six behavioural mutations, all caught:** unknown name colliding with the remainder wording, unknown name emptied, remainder suppressed, unknown bucket nulled, hidden count zeroed, UR-less samples defaulted to `UR-000`.

**Two defects the render caught that no assertion would have:** a bracket running past the plot's right edge on a UR completing at the axis end, and an alternating row tone measuring **1.293:1** — verified by the reviewer to three decimal places against every figure recorded in the diff.

*Verified by work action*

## Review

**Overall: 87%** | 2026-08-24T12:37:51Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 70% |
| Scope | 85% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:**
- **The lane's geometry is held by nothing.** Five mutations of its constants pass the whole suite: removing the plot-right-edge clamp, separation 2→0, row pitch 10→6 (rows overlap 7-high brackets), min width 3→0 (single-day URs vanish), and `DURATIONS_UR_LANE_TOP` 358→430. The reviewer rendered that last one and measured **7 bracket-over-text intersections** with panel B's title against 0 on shipped code; the orchestrator reproduced it independently. The board already has this rule one panel up (`TestBrowserBehaviorDurationsLabelRowsClearTheirNeighbours`); the new lane did not inherit it. The shipped code is correct — this is coverage, not a defect. — `impact-rule-change` → **REQ-364**
- **The `write_set` mirror was not extended** when `generate_test.go` landed outside it. Bookkeeping fixed here; the generalisable rule — a `tdd: true` REQ whose write set names no test file is uncompletable as written — → **REQ-365**
- **Second sighting of the `<body>`-is-the-surface trap.** `body` is `rgb(245,247,250)` light / `rgb(12,14,18)` dark while `#durations-chart` and its `<svg>` are both transparent. REQ-321 recorded the identical discovery for the timeline's bars. Two builders have now paid to rediscover it inside a feature comment. → promoted to `_dev/primes/prime-kanban-board.md` in this REQ's Lessons-Capture Phase, not a REQ.

**Minor:** the "6 rows → 12" figure in the source comment is a property of one synthetic fixture, not a prediction for the consuming project (the reviewer's equivalent board gives 30); the lane hover has no horizontal distance cutoff, consistent with panel A's existing nearest-mark behaviour; six `"/probe.html"` literals remain beside the constant.

**Acceptance:** Pass — rendered and measured on three generated boards, all three hover surfaces driven live, 6/6 behavioural mutations caught.

**Follow-ups created:** REQ-364, REQ-365

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Rendering the board and measuring it. Two defects that every assertion passed over — a bracket over the plot's right edge, and a row tone at 1.29:1 — were only visible in the render. The repository's standing rule for chart changes earned itself again.

**What didn't:** The lock-in pins *behaviour* and not *geometry*. Six mutations of the naming and counting rules are caught; five mutations of the lane's constants ship green, including one that puts seven brackets through panel B's title. The panel one level up already has a neighbour-clearance probe; the new lane did not inherit it, because the test was written against the requirements rather than against the failure modes the render had just demonstrated.

**Worth knowing:** The surface behind this view's SVG is `<body>`, not any `--surface-*` token, so a tone chosen against a surface variable measures nothing like what a reader sees. That is now recorded in the prime rather than in a feature comment, because this is its second sighting.

## Orientation

The Durations view can say which user request a mark belongs to — in the hover, in the sample table, and as a bracket lane under panel A. Lives in the board tool's Durations view. `[MAP CHANGED]` is not warranted: no new module, no new data flow, and the payload is unchanged — the join was already in the browser and this view simply stopped ignoring it.
