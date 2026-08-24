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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
