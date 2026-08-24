---
id: REQ-369
title: "Wait for the Timeline table rebuild condition"
status: pending
created_at: 2026-08-24T16:41:18Z
user_request: UR-069
addendum_to: REQ-348
domain: testing
review_generated: true
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec: bug-fix
depends_on: [REQ-348]
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# Wait for the Timeline Table Rebuild Condition

## What

Replace the fixed 50ms delays in `TestBrowserBehaviorTimelineListsRowsBeneathUserRequestHeaders`
with an observable-condition wait for the Timeline table/group rebuild.

## Why

REQ-348 re-review ran the probe repeatedly. One `-count=5` iteration read stale Fit-all table
contents after switching to Day, while five separately launched runs and the complete focused group
passed. The pure model tests independently prove the product transition, so this is synchronization
flakiness in the live browser probe rather than a product defect.

## Detailed Requirements

- Wait for a state that proves the requested window's table/group rebuild completed; do not increase
  a fixed delay.
- Keep the header/member ordering, window-change rebuild, virtualization, roving focus, and explicit
  table-header association assertions intact.
- Run the probe repeatedly enough to falsify the stale-table failure observed in review.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs.
- Change only the test synchronization boundary; do not weaken product assertions or change runtime
  Timeline behavior.

## Red-Green Proof

**RED:** On REQ-348's reviewed merge, a `-count=5` run failed once because the fixed 50ms delay read
the previous Fit-all table after the Day selection.

**GREEN:** The probe waits on the rebuild condition and repeated executions pass without increasing
an arbitrary timeout.

---
*Source: REQ-348 independent re-review Minor finding.*
