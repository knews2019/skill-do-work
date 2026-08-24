---
id: REQ-369
title: "Wait for the Timeline table rebuild condition"
status: completed
claimed_at: 2026-08-24T21:21:06Z
completed_at: 2026-08-24T21:47:55Z
commit: ed15507d3c94573c007aff7e66d40abacfd72c88
status_changed_at: 2026-08-24T21:47:55Z
route: B
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
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-24T21:21:06Z
  basis:
    - trivial short-circuit
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

## Triage

**Route: B** — The single test-only target and stale-table symptom are known, but exploration must identify an observable rebuild-completion condition that cannot self-satisfy on the previous window, preserve every existing assertion, and prove stability under repeated retained-browser runs.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Reproduced the stale Fit-all read twice in 20 Chromium runs, traced the table's replace-all publication boundary, and chose a pre-action MutationObserver with explicit non-vacuity and failure serialization.
- [x] **[APPLY]:** Replaced both fixed delays with first-row replacement waits and added result fields/fatals proving each action began from a real table and published a replacement before existing assertions ran.
- [x] **[UNIFY]:** Reviewed the exact one-file diff; 60 repeated final runs, inverted-condition RED, related grouping/focus/theme probes, full module, vet, canonical, formatting, diff, and artifact checks passed.

## Plan

**Planning not required** — Route B: explore the renderer's publication boundary, replace elapsed-time success with a non-vacuous observable condition, and replay the captured race repeatedly.

## Exploration

- The captured race reproduced 2/20 on retained Chromium 1228: after the Day action, both snapshots still contained the same 363-REQ Fit-all table.
- `renderTimelineTable` clears and recreates every row synchronously inside one scheduled frame. A different nonempty first-row node therefore proves the requested action reached a complete replacement table, while control state alone proves only that the handler ran.
- A timer-poll version still failed once in 20 because virtual timers could exhaust before Chromium delivered a compositor frame. Arm a child-list MutationObserver before the action and keep a 10-second deadline only to serialize missing publication as `rebuildObserved: false`.
- Preserve every existing grouping, order, window-difference, virtualization, roving-focus, and header-association assertion after the two observed rebuilds.

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go`

**Acceptance criteria:**

- Both Fit-all and Day actions arm an observer before invocation and continue only after a different nonempty first row proves complete table replacement.
- Existing prior-row and observed-rebuild guards make the condition non-vacuous; the deadline can only report explicit failure, never success.
- All existing grouping, order, window-change, virtualization, roving-focus, and table-header association assertions remain intact.
- The captured stale read reproduces on the base, final retained-Chromium runs pass repeatedly, and an inverted identity condition returns explicit RED without runtime changes.

## Decisions

- **D-01 — Use first-row identity as the publication boundary.** The renderer's established replace-all behavior makes identity change a complete-table signal rather than an indirect control-state proxy.
- **D-02 — Arm before triggering.** No synchronous or scheduled rebuild can slip between the baseline snapshot and observer registration.
- **D-03 — Treat the deadline as failure serialization only.** A successful path is driven exclusively by observed node replacement; deadline expiry writes an explicit false result for Go to reject.

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified): replaces both fixed 50ms waits with a pre-action MutationObserver that requires nonempty first-row replacement, reports explicit deadline failure, and leaves every prior behavior/accessibility assertion intact.

## Discovered Tasks

None.

## Testing

- Base Chromium 1228 repeated runs reproduced stale Fit-all content 2/20. Final condition wait passed 30/30 twice after fitting both possible deadline serializations inside the existing virtual-time budget.
- Inverting the node identity condition returned explicit RED: requested rebuild not observed for both Fit-all and Day. Related grouping, roving-focus, and light/dark group-header probes passed three repeated runs each.
- Go vet, uncached full queue-kanban tests, formatting, diff checks, and the builder browserless canonical gate passed; retained-browser evidence came from the focused repeated runs.
- Independent review approved at 10/10 with no Critical, Major, or Minor findings. It independently reproduced base RED 2/20, final GREEN 20/20, the mutation RED, and confirmed the deadline cannot self-satisfy stale content.
- On merged main, the named retained-Chromium probe passed another 20/20 runs. The final canonical gate passed all contracts, queue-kanban tests, strict JavaScript, and audit metrics; its optional browser lane made the standard no-browser skip after the repeated focused evidence.

## Qualification

- Exact merge range `be6f2f1..ed15507d3c94573c007aff7e66d40abacfd72c88` passed mechanical qualification.
- Scope drift passed: the one changed test file exactly matches the declared Scope and Implementation Summary.
- Orchestrator judgment confirmed substantive non-vacuous publication synchronization, complete preservation of behavioral/accessibility assertions, explicit failure-only deadline flow, and no generated/debug artifacts.

## Review

Independent review approved at 10/10 with no Critical, Major, or Minor findings and low residual risk. It confirmed arm-before-action ordering, complete replace-all publication semantics, non-vacuous prior-row guards, failure-only deadline behavior, every retained assertion, base RED 2/20, final GREEN 20/20, and inverted-condition RED.

## Lessons Learned

An asynchronous test should wait on the artifact the renderer publishes, not the control state that requested it or a duration guessed from one machine. Arm before the action, require a real prior value, and serialize timeout as explicit failure.

## Orientation

Released in 0.236.59. The Timeline grouping probe now waits for actual table replacement, eliminating stale Fit-all reads without weakening its behavior or accessibility assertions.
