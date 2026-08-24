---
id: REQ-351
title: "Retire the Durations in-lane labels for a ranked longest-spans list"
status: completed
claimed_at: 2026-08-24T16:43:37Z
completed_at: 2026-08-24T17:36:50Z
commit: f58e1ae97e79386677509932329be8cfcbecb923
status_changed_at: 2026-08-24T17:36:50Z
route: B
created_at: 2026-08-23T22:37:52Z
user_request: UR-069
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-346]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-346, REQ-347, REQ-348, REQ-349, REQ-350, REQ-352, REQ-353, REQ-354]
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
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board.css
---

# Retire the Durations In-Lane Labels

## What

Replace the overflow lane's direct labels with a compact ranked longest-spans list beside the chart,
and delete the label planner, its width model and the tests that exist only to place text inside the
SVG. Keep the lane, its marks, their hover and a plain remainder count.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route B exploration mapped the complete seven-file planner/width-model deletion boundary and the lane/mark/hover behavior that must survive. Replace SVG labels with a complete, current-window, descending HTML list beside the chart and prove it on a 705+ sample board.
- [x] **[APPLY]:** Removed the direct SVG label planner/width model and added a complete, deterministic, current-window longest-spans list while preserving the lane, marks, hover, UR grouping, headline statistics, rolling median, cadence ticks, and table.
- [x] **[UNIFY]:** Reviewed all seven scoped files and the REQ-352 integration seam; `gofmt`, `node --check`, focused Go/browser tests, canonical maintainer verification, and diff checks passed with no debug artifacts.

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
- The list carries the UR, so REQ-346's join is a prerequisite.

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

`depends_on: REQ-346` — the list's UR column needs the join.

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

## Triage

**Route: B** — The reversal, retained/deleted behavior, six-file seed scope, and rendered outcome are explicit. Exploration must trace the full label-planner/width-model/test deletion boundary before implementation.

## Plan

**Planning not required** — Route B: exploration-guided implementation.

## Exploration

- `generate_test.go` contains three live behavior tests solely for the SVG label planner, so it is added to scope. No payload/generator change is needed: `boardData.requests[sample.id]` supplies title, route, and UR while the Durations sample supplies id and wall time.
- Remove the JavaScript label-row/leader/packer/remainder-reservation pipeline and browser text measurement, plus the Go label text/width model and tests whose only shipped claim is retired placement. Retain day-domain helpers, Panel B annotation formatting/styling, lane marks, overflow/reversed hover, UR lane, window projection, and the full sample table.
- Replace the deleted longest-first/count claims with one behavior-first dense-board proof: every current-window over-ceiling sample appears in a deterministic descending HTML list outside the SVG, with REQ, UR, duration, route, and title; SVG label/leader nodes are absent while lane marks and exact hover identity remain.
- Use a chart/list wrapper: desktop grid with a 280–360px scrollable aside aligned to the chart, stacked below about 1000px. Every item remains in the DOM; 320/768/1280 layouts must avoid new horizontal clipping.
- A plain sentence states the complete over-ceiling count; it is not a `+N more` remainder because nothing is omitted.

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/web/board-durations.js`
- `skills/do-work-board/tools/queue-kanban/durations.go`
- `skills/do-work-board/tools/queue-kanban/durations_test.go`
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`
- `skills/do-work-board/tools/queue-kanban/web/template.html`
- `skills/do-work-board/tools/queue-kanban/web/board.css`

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified): deletes direct SVG duration labels, leader/packing/reservation machinery, and measured-face sizing; renders every positive current-window span over 60 minutes into a deterministic descending semantic list with REQ, duration, UR, route, and title.
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified): removes the retired label text/width payload model while preserving the underlying samples and domain geometry.
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified): removes Go tests that specified only the retired label formatter/width model.
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go` (modified): replaces label-placement-only probes with a 705-sample, 94-overflow responsive matrix that proves complete list content/order, absent SVG labels, preserved circles/hover, and bounded layout.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified): removes renderer tests for the retired two-row planner and pins the complete adjacent list contract without disturbing REQ-352's headline/rolling/cadence lock-ins.
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified): wraps the chart with an adjacent semantic longest-spans aside and complete-count sentence.
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified): removes label-only SVG styling and adds a bounded 340px desktop aside that stacks below the chart under 1000px.

## Decisions

- **D-01 — Keep every qualifying row in semantic HTML.** An ordered list makes the complete set accessible and inspectable without recreating an SVG packing limit.
- **D-02 — Interpret “over 60 minutes” strictly.** The list contains positive `wallMinutes > 60` samples; the exact count sentence says every one is listed.
- **D-03 — Make equal durations deterministic.** REQ id and completion time break ties after descending wall minutes.
- **D-04 — Preserve independent geometry helpers.** Day-bucket projection and Panel B annotation helpers remain even where their historical names mention labels, because current non-retired behavior and tests still consume them.
- **D-05 — Reconcile the concurrent headline release at merge time.** The shared renderer retains REQ-352's quantile/stat/rolling/tick paths while removing only the direct-label planner and adding the complete list.

## Discovered Tasks

- A full browser-enabled run reproduced two pre-existing Timeline failures: `TestBrowserBehaviorTimelineBarsSurviveTheDetailDrawerOpening` and `TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage`. They reproduce outside this seven-file write set and require a separate Timeline investigation.

## Testing

- RED proved the generated page lacked the complete adjacent longest-spans list and its count contract.
- Focused renderer tests passed for the complete list and preserved Panel A/day-bucket behavior; `node --check`, `gofmt`, `go vet`, diff checks, and canonical maintainer verification passed.
- The merged combined browser matrix passed all six dense 320/768/1280 cases and all four REQ-352 headline cases. Every dense case asserted all 94 qualifying rows, absent SVG labels/leaders/remainder text, preserved circle/hover identity, deterministic content, and bounded responsive geometry.
- The post-merge `GOTOOLCHAIN=go1.26.1 go test ./... -count=1 -timeout=10m` suite passed in 43.849s.
- The orchestrator generated this repository's real queue: the current 30-day window listed all 23 over-ceiling spans in descending order, retained headline/rolling/cadence output, removed direct SVG REQ labels, and stacked the list below the chart at 320px. The only console item was the unrelated missing favicon.
- Twelve deleted tests are accounted for in the builder hand-back: each specified only the retired direct-label formatter, measured-width model, two-row packer, fixed-point reservation, leader placement, or incomplete remainder sentence. Generic span formatting, lane, marks, hover, day-domain, Panel B annotation, and REQ-352 tests remain.

## Qualification

- Exact merge range `ae2368c..f58e1ae` passed mechanical qualification.
- Scope drift passed: the seven-file Implementation Summary exactly matches the declared Scope.
- Orchestrator judgment confirmed substantive implementation, requirement coverage, complete request-metadata flow, deliberate deletion boundaries, preserved REQ-352 integration seams, and no generated/debug artifacts in the merge.

## Review

**Overall: 98%** | 2026-08-24T17:35:53Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 94% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings:** None.

**Minor findings:** 2 (report only) — the dense probe's `--force-dark-mode` cases report the light body colour even though separately inspected Playwright screenshots prove dark rendering, and two renderer comments still describe retired label-planning rules.

**Acceptance:** Pass — independent focused behavior and ten-case Chromium runs passed; all 94 dense rows, preserved lane/hover, exact seven-file scope, twelve deleted-test explanations, REQ-352 integration seams, and real/dense screenshots were verified.

**Suggested testing:** 1 item — make the dense automated theme cases assert an actually emulated dark surface if that coverage is revised later.

**Follow-ups created:** REQ-370 for the distinct pre-existing, non-falsifiable Timeline pointer-capture mutation; **sweeps appended to:** None.

*Reviewed by review-work action*

## Lessons Learned

Deleting a visual packing system safely means accounting for every retired test claim while proving the information, interaction, and neighboring releases survive. Completeness belongs in semantic HTML; dense SVG should keep only the marks and geometry it can state honestly.

## Orientation

Released in 0.236.47. Durations keeps its lane and hoverable marks but moves every current-window span over 60 minutes into one complete, deterministic, responsive longest-spans list beside the chart.
