---
id: REQ-349
title: "Hide the dead filter knobs while Durations is on screen"
status: pending
created_at: 2026-08-23T22:37:52Z
user_request: UR-068
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-342, REQ-343, REQ-344, REQ-345, REQ-346, REQ-347, REQ-348, REQ-350]
batch: durations-panel-improvement
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
---

# Hide the Dead Filter Knobs While Durations Is On Screen

## What

The topbar's text filter, domain select and status select stay visible and live while the Durations
view is on screen and change nothing there. Hide them on that view, the same way the lens group is
hidden, and restore them on view change.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
