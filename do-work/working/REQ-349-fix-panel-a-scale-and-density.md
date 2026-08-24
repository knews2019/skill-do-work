---
id: REQ-349
title: "Fix panel A's scale and density"
status: claimed
claimed_at: 2026-08-24T15:00:04Z
status_changed_at: 2026-08-24T15:00:04Z
route: B
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
related: [REQ-346, REQ-347, REQ-348, REQ-350, REQ-351, REQ-352, REQ-353, REQ-354]
batch: durations-panel-improvement
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-23T22:37:52Z
  basis:
    - Route B
    - 3-file write set
    - 6 acceptance criteria
    - browser evidence
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-durations.js
  - skills/do-work-board/tools/queue-kanban/durations_test.go
  - skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Fix Panel A's Scale and Density

## What

Panel A wastes its vertical range and overplots. Move it to a square-root y scale, jitter marks
inside their own day slot, lower their opacity, and draw a per-day median line with a p25-to-p75
ribbon behind them.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route B exploration located the single renderer path, its hover index, and the existing geometry lock-in. Implement a named square-root scale and deterministic within-day jitter, aggregate accepted samples into daily quartiles, draw the overlay behind marks, and update both structural and live-browser geometry evidence.
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

On this repository's archive the median REQ is 13.2 minutes, 55 percent are under 15 and 78 percent
under 30. Against a linear 0-to-60 scale more than half the marks compete for the bottom quarter of
the panel. Horizontally a day gets about 13 SVG units here and about 8 on the consuming project, so a
busy day (38 REQs here, 55 there) is a single column of overlapping 4-unit dots. There is no opacity,
no jitter, no density encoding and no drawn trend, so the shape of the distribution is only reachable
by hovering one mark at a time.

## Detailed Requirements

- **Square-root y scale** with ticks at 0, 5, 15, 30, 45 and 60.
- **Deterministic within-day x jitter, bounded to the day's own slot** — a mark can never cross a day
  boundary. Deterministic means the same board renders the same picture twice.
- **Lower mark opacity.**
- **A per-day median line with a p25-to-p75 ribbon** behind the marks, so "is it getting slower" is
  readable without hovering.
- **Keep unchanged:** the 60-minute ceiling, the overflow lane, and the read-time exclusion rule.
- **The hover must still name the mark under the pointer.** Either keep the jitter out of the
  nearest-mark maths, or compensate for it and state how. A reader who hovers a jittered mark and
  gets its neighbour's id is a worse view than the one this REQ replaces.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first.
- The browser probe lane measures the rendered face and the geometry tests read this file's
  constants — a new scale means the probes that pin tick positions move with it, not around it.
- Generate a board and look at it. Overplotting is exactly the class of defect a green suite hides.
- **The target is legibility at 700 or more archived REQs**, not at this repository's 305. The
  consuming project is already at 692 samples across 47 active days, at about 8 SVG units per day.

## Dependencies

None. REQ-354 wires a click to the same nearest-mark resolution this REQ perturbs; whichever lands
second reconciles the two.

## Builder Guidance

**Certainty: firm.** Every parameter is specified. The one judgment call is the jitter amplitude, and
its bound is stated: the day's own slot at 8 SVG units per day on the consuming project, not 13.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repository and open Durations. Marks below 15 minutes
— 55 percent of them — are packed into the bottom quarter of panel A, a busy day is one column of
overlapping dots, and no line states the per-day trend.

**Why RED now:** `yOfMinutes` is linear to a 60-minute ceiling, marks are drawn at full opacity with
no jitter and no aggregate overlay.

**GREEN when:** panel A uses a square-root scale with ticks at 0/5/15/30/45/60, marks are jittered
inside their day slot at reduced opacity, a per-day median line and p25-p75 ribbon are drawn, the
ceiling and overflow lane are unchanged, and the hover still names the mark the reader is pointing at.

**Validation:** User confirmed (bundled invocation).

---
*Source: prompt A4, `ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html` (finding F4).*

## Triage

**Route: B** — The visual outcome and files are explicit; exploration is needed to locate the existing scale, hover geometry, and browser-probe conventions before implementation.

## Plan

**Planning not required** — Route B: exploration-guided implementation.

## Exploration

- `renderDurationsView` owns the complete Panel A path. `yOfMinutes` is linear today; the mark loop uses raw completion x coordinates, and `describeAtPointer` resolves against the same `markIndex` geometry.
- Hover stays truthful when the jittered x is used for both the circle and `markIndex.x`. Overflow/reversed y placement and the Go-produced `excludedReason` verdict remain unchanged.
- No existing quantile or jitter helper exists. Daily p25/median/p75 can be computed client-side from sorted, non-excluded samples using the payload's canonical verdict, centred on each UTC day slot.
- `generate_test.go::TestJavaScriptBehaviorDurationsDayBucketsStayInsideThePlot` pins Panel A circles to unjittered completion x values. It must move with the intended behavior, so `generate_test.go` is added to the declared scope.
- Dense browser evidence needs at least 700 samples across roughly 47 days, with actual circle coordinates, day bounds, repeat-render determinism, hover-at-circle identity, and ribbon/line bounds measured from the rendered page.

## Scope

- `skills/do-work-board/tools/queue-kanban/web/board-durations.js`
- `skills/do-work-board/tools/queue-kanban/durations_test.go`
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`
