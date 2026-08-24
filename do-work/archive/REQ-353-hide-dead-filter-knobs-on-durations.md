---
id: REQ-353
title: "Hide the dead filter knobs while Durations is on screen"
status: completed
claimed_at: 2026-08-24T16:43:37Z
completed_at: 2026-08-24T16:58:09Z
commit: 2dd9f3c
status_changed_at: 2026-08-24T16:58:09Z
route: A
created_at: 2026-08-23T22:37:52Z
user_request: UR-069
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-346, REQ-347, REQ-348, REQ-349, REQ-350, REQ-351, REQ-352, REQ-354]
batch: durations-panel-improvement
estimate:
  p50_active_minutes: 10
  confidence: high
  calculated_at: 2026-08-24T16:43:37Z
  basis:
    - Route A
    - 1-file write set
    - existing visibility mechanism
    - 4 acceptance criteria
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
---

# Hide the Dead Filter Knobs While Durations Is On Screen

## What

The topbar's text filter, domain select and status select stay visible and live while the Durations
view is on screen and change nothing there. Hide them on that view, the same way the lens group is
hidden, and restore them on view change.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route A reused the existing `applyView` hidden-state mechanism; no separate state or data-flow path was needed.
- [x] **[APPLY]:** The isolated builder changed only `web/board-controls.js` to hide all three shared filters on Durations and restore them on other views.
- [x] **[UNIFY]:** Node syntax, focused assembled-client tests, the full module suite, diff checks, generated-page keyboard round trip, hidden/display states, and zero-console behavior passed. Worktree and merge range contain no artifact or scope drift.

## Why

This is not an oversight in one place and a decision in another — both halves were decided and only
one was applied. `applyView` in `web/board-controls.js` hides the lens and recent-window groups on
views where they do nothing, with the comment that the topbar should never advertise dead knobs.
`onFiltersChanged` in `web/board-filters.js` deliberately excludes Durations, because a filtered
distribution is a different statistic wearing the same axes. The conclusion was reached; the controls
were left on screen.

## Detailed Requirements

- Hide the text filter, the domain select and the status select while Durations is the active view.
- Use the mechanism `applyView` already uses for the lens group, not a second one.
- Restore them on view change.
- **Do not widen `onFiltersChanged` to include Durations.** If a filtered durations distribution is
  wanted later, that is a separate request with its own caption rules, not a silent widening of this
  one.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first.
- One file. Do not touch `web/board-filters.js`.
- Generate a board, switch to Durations and back, and look at it.

## Builder Guidance

**Certainty: firm.** This applies a rule the code already states, in the place that already states it.
The smallest change in the batch — keep it that way.

## Red-Green Proof

**RED prompt/case:** Generate a board, switch to the Durations view, type into the text filter or
change the domain select. The controls are visible and accept input; the chart does not change.

**Why RED now:** `applyView` hides the lens and recent-window groups on views where they do nothing
but leaves the three filter controls visible on Durations, which `onFiltersChanged` deliberately
excludes.

**GREEN when:** switching to Durations hides the text filter, domain select and status select the
same way the lens group is hidden, switching away restores them, and `onFiltersChanged` still
excludes Durations.

**Validation:** User confirmed (bundled invocation).

---
*Source: prompt A6, `ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html` (finding F6).*

## Triage

**Route: A** — This is one existing visibility rule applied to three existing controls in one file, with no state or data-flow change.

## Plan

**Planning and exploration skipped** — Route A: implement the existing `applyView` visibility pattern directly and verify the round trip.

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified): applies the existing `hidden` visibility mechanism to search, domain, and status filters while Durations is active and restores them on every other view without changing stored filter state.

## Decisions

- **D-01 — Reuse direct `hidden` assignments in `applyView`.** This matches the existing lens, recent-window, Durations, and testing control visibility path.
- **D-02 — Preserve shared-filter values.** Durations only hides dead controls; switching away restores their existing state and does not create a filtered Durations statistic.
- **D-03 — Keep the explicit one-file constraint.** No new test file was added; assembled-client tests, full module tests, and direct generated-page RED/GREEN evidence cover the change.

## Discovered Tasks

None.

## Testing

- Generated-page RED showed search, domain, and status controls visible, displayed, and keyboard reachable on Durations.
- GREEN through Chromium 151 proved all three `hidden=true`/`display=none` on Durations, absent from sequential focus, restored on Board/Timeline, and hidden again on return; filter values and focus remained intact with zero board-script console errors.
- Node syntax, focused assembled-client tests, diff checks, and builder full module suite (91.970s) passed. Post-merge suite passed in 49.715s; independent reviewer full suite passed in 50.355s.

## Qualification

- Exact merge range `907f106..2dd9f3c` passed mechanical qualification.
- Orchestrator judgment confirmed substantive user-visible behavior, unchanged filter data flow, one-file scope, and no debug/generated artifacts.

## Review

Independent review approved with no findings: Requirements 100%, Code Quality 100%, Test Adequacy 96%, Scope 100%, overall 99%, low risk. Reviewer independently repeated the generated-board keyboard hide/restore acceptance.

## Lessons Learned

When a view deliberately ignores shared filters, the controls must disappear through the same view-state mechanism that already hides every other dead knob; preserving their stored values makes the round trip predictable.

## Orientation

Released in 0.236.45. Search, domain, and status controls now disappear while Durations is active and restore unchanged on other views; Durations remains intentionally unfiltered.
