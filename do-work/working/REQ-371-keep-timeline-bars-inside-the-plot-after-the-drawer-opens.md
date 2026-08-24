---
id: REQ-371
title: "[impact-critical] Review fix: keep Timeline bars inside the plot after the drawer opens"
status: claimed
claimed_at: 2026-08-24T20:34:04Z
status_changed_at: 2026-08-24T20:34:04Z
route: C
created_at: 2026-08-24T18:22:57Z
user_request: UR-066
addendum_to: REQ-331
domain: frontend
review_generated: true
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-331]
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 50
  confidence: medium
  calculated_at: 2026-08-24T20:34:04Z
  basis:
    - Route C
    - 2-file write set
    - 4 acceptance criteria
    - browser evidence
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# Review Fix: Keep Timeline Bars Inside the Plot After the Drawer Opens

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Measured the retained-browser observer/frame boundary, preserved host-width change as the trigger, kept drag work frame-coalesced, and designed a compositor-independent live-width fallback plus stronger shared-scale assertions.
- [x] **[APPLY]:** Recorded the last rendered host width, rendered directly on observer-delivered width change, added a teardown-owned positive-width condition poll, and strengthened open/close browser snapshots without naming drawer callers.
- [x] **[UNIFY]:** Reviewed both files; repeated Chromium, required mutation, drag-cost, zero-size, focused Timeline, full module, vet, syntax, canonical, diff, and artifact checks passed apart from the separately queued REQ-370 baseline.

## What

Restore REQ-331's drawer-resize invariant under the retained Chromium: after a Timeline row opens the
detail drawer, rendered bars must be remeasured against and remain inside the narrowed plot host.

## Context

Independent review of REQ-354 reproduced
`TestBrowserBehaviorTimelineBarsSurviveTheDetailDrawerOpening` in isolation on Chromium 1228. The
drawer opened and 40 segment nodes remained in the DOM, but none intersected the 851px-wide host:
the leftmost segment started at x=1480 while the host ended at x=903. This is the exact user-visible
failure REQ-331 fixed, not REQ-370's separate pointer-capture mutation-sensitivity gap.

## Requirements

- Reproduce the isolated drawer-opening failure on the retained browser and identify why the
  ResizeObserver/invalidation path no longer redraws against the narrowed host.
- Restore the condition-based remeasurement contract; do not enumerate drawer callers or weaken the
  bar-intersection assertion.
- Keep the shared Timeline axis and row SVG on the same width measurement, and preserve drawer open
  and close behavior without requiring a window move.
- Make the existing trusted browser probe green repeatedly and demonstrate that suppressing the
  required remeasurement still makes it fail.

## Red-Green Proof

**RED prompt/case:** Run `TestBrowserBehaviorTimelineBarsSurviveTheDetailDrawerOpening` with Chromium
1228. After the drawer opens, 40 segments exist but zero intersect the host; leftmost x=1480 is beyond
host right x=903.

**Why RED now:** The plot retains coordinates from the wide layout even though the drawer has narrowed
its host, so DOM node count stays nonzero while the visible chart is blank.

**GREEN when:** the drawer-open and drawer-closed snapshots both contain visible intersecting segments,
their coordinates are measured against the current host width, and a targeted invalidation mutation
returns the probe to RED.

**Validation:** Independent review finding from REQ-354; auto-queued because it restores an
`impact-critical` user-visible invariant.

---
*Source: REQ-354 independent review; regression of REQ-331 (UR-066).*

## Triage

**Route: C** — This is an impact-critical regression in an asynchronous ResizeObserver/render invalidation path. The two-file target is known, but the retained-browser event boundary must be measured, the condition-based repair planned, both drawer directions proven, and the existing mutation ratchet preserved before release.

## Prior Implementation

REQ-331 introduced a `ResizeObserver` on the Timeline scroll host, one-width-per-render memoization, zero-size render refusal, and the existing drawer/open-close browser probe. Under Chromium 1228 the observer still schedules through `requestAnimationFrame`, while the probe runner's DOM-dump mode does not reliably drive compositor frames after the drawer layout change; 42 segments remain at wide-layout coordinates outside the narrowed host. Reproduce and measure the current boundary before choosing the repair.

## Plan

1. Instrument the generated page to measure host-width, ResizeObserver, and frame delivery around drawer open without altering production behavior.
2. Keep live host-width change as the condition and preserve one width measurement per render, zero-size refusal, and drag frame coalescing.
3. Add a compositor-independent fallback only if direct observer rendering remains insufficient under repeated Chromium runs.
4. Strengthen the existing trusted probe to pin bar containment, axis/row shared scale, unchanged range, and close recovery; replay a required-remeasurement mutation.

**Plan validation:** Each impact-critical invariant maps to a measured event boundary, an existing product guard, or a non-vacuous browser assertion inside the declared two-file scope.

## Exploration

- Initial Chromium 1228 RED narrowed the host from 1481px to 851px while 40 segments remained at wide-layout coordinates; none were visible and the leftmost began at x=1476 beyond host-right x=903.
- During the bounded post-click timer chain, neither ResizeObserver nor animation-frame callbacks were delivered. A longer idle budget delivered an old frame, then the observer, then another queued frame.
- Rendering directly from ResizeObserver still failed one of five runs, so observer delivery alone is not a stable retained-browser boundary. A 50ms teardown-owned poll of the same positive live-width condition supplies the missing compositor-independent trigger and never renders when width is unchanged.
- The existing render-cost and zero-size probes already pin one measurement per render, pan coalescing, and hidden-host refusal. Strengthen the drawer probe with equal nonempty axis/gridline x scales and an unchanged nonempty range readout across before/open/closed snapshots.

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js`
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go`

**Acceptance criteria:**

- A positive live host-width change triggers condition-based remeasurement without enumerating drawer or layout callers; zero-size hosts remain non-renderable and pan rendering remains frame-coalesced.
- Drawer-open and drawer-close snapshots retain visible segments inside current plot bounds without a window move, using one shared axis/row scale and an unchanged time range.
- The named retained-Chromium probe passes repeatedly, while suppressing the width-change condition returns the exact blank-chart RED.
- Focused Timeline/static, drag-cost, full module, and canonical gates pass; unrelated REQ-370 pointer-capture baseline evidence stays isolated.

## Decisions

- **D-01 — Compare live width with the last real render measurement.** The value directly represents whether the current chart scale matches its host and avoids a caller list.
- **D-02 — Keep drag and resize scheduling separate.** Pointer movement remains animation-frame-coalesced; resize recovery may render directly because the retained headless engine can park compositor frames after layout changes.
- **D-03 — Add a bounded-frequency condition fallback after observer-only evidence failed.** The 50ms check reads a width cheaply, renders only on positive change, and is torn down with the Timeline listeners.
- **D-04 — Pin shared scale and range, not only visible node count.** Equal axis/gridline coordinates and stable range text prove the redraw used one current measurement without changing the selected window.
- **D-05 — Describe the condition, not one delivery mechanism.** Review found commentary still naming ResizeObserver/rAF as the sole repair even though observer and timer now deliver one shared positive-width-change condition; correct both changed files before release.

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified): tracks the host width used by the last render, re-renders on observer-delivered positive width change, and adds a teardown-owned condition poll for engines that park observer/frame delivery while leaving pan coalescing and zero-size refusal intact.
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified): strengthens the drawer open/close probe with shared axis/gridline scale and stable-range assertions while retaining bar intersection, current-right-edge, URL, and recovery checks.

## Discovered Tasks

None. The remaining Timeline pointer-capture mutation failure is already REQ-370.

## Testing

- Initial RED measured host 1481px→851px, 40 segments/0 visible, leftmost x=1476 beyond host-right x=903, with no observer/frame delivery during the bounded poll.
- Final named Chromium 1228 probe passed 15 consecutive runs; the drawer probe plus drag render-cost probe passed three paired runs, and five focused Timeline JavaScript fixtures including zero-size passed.
- Reversing the required width-change condition returned the exact 40-node/0-visible RED and was restored. Axis ticks and row gridlines now share equal nonempty x coordinates, and the nonempty range readout remains unchanged through open and close.
- Go vet, Node syntax, uncached full queue-kanban tests, and the builder browserless canonical gate passed. Browser-focused and strict lanes leave only the independently queued REQ-370 pointer-capture mutation failure; one transient grouping probe passed isolation five times.
- Initial independent review found no behavioral defect and scored the change 98/100, but conditionally accepted one Minor: comments in both changed files still described the superseded observer/rAF-only repair and removal mutation. Comment-only remediation will name the shared condition and both delivery paths.
