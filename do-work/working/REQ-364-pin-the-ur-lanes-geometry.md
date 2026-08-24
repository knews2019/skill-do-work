---
id: REQ-364
title: "[impact-rule-change] Pin the UR lane's geometry"
status: claimed
claimed_at: 2026-08-24T18:57:05Z
status_changed_at: 2026-08-24T18:57:05Z
route: B
created_at: 2026-08-24T12:50:00Z
user_request: UR-069
addendum_to: REQ-346
domain: testing
review_generated: true
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-24T18:57:05Z
  basis:
    - Route B
    - 1-file write set
    - 4 acceptance criteria
    - browser evidence
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go
---

# Pin the UR Lane's Geometry

## What

REQ-346's UR lane is pinned by six behavioural mutations and by **no geometric one**. Five mutations
of its constants pass the entire suite, including one that puts seven brackets through panel B's
title. Give the lane the neighbour-clearance probe the panel above it already has.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The shipped code is correct. This is missing coverage, not a defect — which is exactly why it is
worth closing now, while the failure modes are known and someone has already rendered them.

Five mutations, each passing the full suite:

| mutation | what it does to the render |
|---|---|
| remove the plot-right-edge clamp | a UR completing at the axis end draws outside the plot |
| `DURATIONS_UR_BRACKET_SEPARATION` 2 → 0 | same-row brackets touch |
| `DURATIONS_UR_LANE_ROW_PITCH` 10 → 6 | rows overlap 7-high brackets |
| `DURATIONS_UR_BRACKET_MIN_WIDTH` 3 → 0 | single-day URs vanish |
| `DURATIONS_UR_LANE_TOP` 358 → 430 | **7 measured bracket-over-text intersections** with panel B's title |

The last was rendered and measured by the reviewer (bracket `y=430 h=7` against title bbox
`y=431.66 h=13.77`, 0 on shipped code) and reproduced independently by the orchestrator: the suite
returns `ok` in 12.5s with brackets crashing through the title below them.

The right-edge clamp is the fix for a defect REQ-346's own render caught — and nothing would notice
its removal. `TestBrowserBehaviorDurationsLabelRowsClearTheirNeighbours` already does this job one
panel up; the new lane did not inherit it because the test was written against the requirements
rather than against the failure modes the render had just demonstrated.

## Detailed Requirements

- A browser probe measures the rendered lane and fails on each of the five mutations above.
- **Measure the rendered geometry, do not restate the constants.** The board prime's REQ-322 lesson
  governs: a constant a decision turns on must be READ by the test, never restated beside it. Assert
  relationships — no bracket intersects a text bbox, no two same-row brackets are closer than the
  separation the code uses, every bracket is inside the plot bounds — rather than asserting the
  numbers those relationships currently produce.
- Follow `TestBrowserBehaviorDurationsLabelRowsClearTheirNeighbours` rather than inventing a shape;
  it is the same property one panel up.
- Return `location.href` with every measurement (the prime's render-evidence rule).

## Constraints

- `_dev/primes/prime-kanban-board.md` governs. Read it first.
- Do not change `web/board-durations.js`. The geometry is correct; this REQ adds the check that keeps
  it correct. If the probe reveals a real defect, that is a separate REQ.
- The probe joins the strict browser lane, so it must SKIP cleanly with no browser and must not
  weaken `TestMaintainerStrictBrowserBehaviorLaneRejectsZeroProbes`.

## Red-Green Proof

**RED prompt/case:** Set `DURATIONS_UR_LANE_TOP = 430` in `web/board-durations.js` and run
`GOTOOLCHAIN=go1.26.1 go test -count=1 -run 'Durations' .` — it returns `ok` in 12.5s while seven
brackets overlap panel B's title. Reproduced 2026-08-24.

**GREEN when:** that mutation fails with a message naming the intersection, the other four fail too,
and the unmutated board passes.

**Validation:** Inferred during REQ-346's review; the counterexample was rendered by the reviewer and
independently reproduced by the orchestrator.

---
*Source: REQ-346 review finding F1 (UR-069).*

## Triage

**Route: B** — Production geometry is intentionally unchanged, but exploration must reuse the existing neighbour-clearance probe shape, identify runtime geometry constants without restating them, and design five mutation-sensitive browser relationships in the one declared test file.

## Plan

**Planning not required** — Route B: exploration-guided implementation.

## Exploration

- Reuse the generated-site browser-probe shape and inject measurement access before the assembled Durations IIFE closes, so the probe reads the live production geometry constants instead of restating them.
- Build one long-axis fixture with crowded one-sample URs, multiple occupied rows, a final-day minimum-width/right-clamp witness, and exactly one no-UR sample in its reserved row.
- Measure every bracket and neighbouring text node from rendered SVG/client rectangles; return `location.href`, the live plot/lane/separation/pitch/height/min-width constants, compact bracket geometry, and console errors in one typed JSON result.
- Assert plot containment, live minimum width, same-row live separation, integer row/pitch placement, non-overlap, exact unknown-row placement, and clearance from lane/gutter/Panel B text. Vacuity guards require the right-edge, minimum-width, multi-row, crowded-separation, and unknown-row witnesses.
- Five mutations disable the right clamp, remove separation from both comparisons, use bracket height as row pitch, remove the minimum-width extension, and move the row-top expression into Panel B. Each must fail this one named probe.

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go`

**Acceptance criteria:**

- One generated-board browser probe measures rendered UR-lane geometry and returns same-measurement URL provenance with non-vacuous fixture witnesses.
- Every bracket remains inside live plot bounds, respects live minimum width and same-row separation, and occupies non-overlapping rows derived from the live lane top/pitch.
- The unknown row and every named bracket clear all rendered lane, gutter, and Panel B text neighbours.
- Each of the five known production mutations fails the named probe; the ordinary no-browser lane skips cleanly while the strict browser lane retains its zero-probe guard.
