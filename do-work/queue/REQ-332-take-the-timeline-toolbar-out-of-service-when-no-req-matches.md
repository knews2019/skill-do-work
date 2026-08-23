---
id: REQ-332
title: "Take the timeline toolbar out of service when no REQ matches the filters"
status: pending
created_at: 2026-08-23T12:22:00Z
user_request: UR-066
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Take the Timeline Toolbar Out of Service When No REQ Matches the Filters

## What

When the filters match nothing, `renderTimelineView` returns early after writing "No REQ matches the current
filters." The early return happens **after** `releaseTimelineListeners()` but **before** the toolbar is
wired — and the toolbar is wired with `button.onclick =`, which is outside the teardown registry. So the
previous render's handlers survive, holding the previous render's `rows`, `filterMatchedRows`, detached
`rowsSvg` and `renderAll`. One click on `Fit all` or `Week` refills the summary, the forecast and the details
table with REQs the filter excludes, over a chart that stays empty.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The empty-filter state is the one state where the view is honest on arrival and then made to lie by a single
click. A reader who searches for something that does not exist, then presses `Fit all` to get back, is told
"317 REQs in the window …" above nothing at all.

## Context

- The early return is `web/board-timeline.js:1029-1036`; `releaseTimelineListeners()` is at `:987`.
- The surviving handlers are assigned at `:1928` (`wireToolbarButton` — zoom in/out, fit, the two period
  steps), `:1961` (`Now`) and `:1992` (the three period chips).
- The From / to fields are the mirror case: their listeners are registered at `:1699-1706`, which the early
  return never reaches, so in this state the fields are **dead** — typing a date does nothing, and the field
  keeps the typed value because `syncRangeField` is never reached either. The chart's arrow keys are dead for
  the same reason.

The fix is small and mechanical: the no-match path is a real state of the view and owes the same teardown the
listener registry gets. Every control the view owns must end that path either inert-and-visibly-disabled or
wired to the current render — not silently wired to the last one.

## Detailed Requirements

1. No handler from a previous render survives the no-match early return. Whatever mechanism the fix uses, the
   condition is the rule: a control the view wires must be released on every path that leaves the view
   without re-wiring it.
2. In the no-match state the toolbar controls are visibly out of service (disabled), so the reader is not
   invited to press something that cannot act.
3. `Clear filters` — the control that actually recovers from this state — keeps working, since it lives
   outside this view.
4. Leaving the no-match state (a filter change that matches something again) restores every control, wired to
   the new render.

## Constraints

- Do not convert the toolbar to `addEventListener` without releasing it: two registrations for one button is
  how this class of bug started elsewhere in the file (see the teardown note at `:80-86`).
- Keep the two empty states distinct in wording. "No REQ matches the current filters" and "Nothing was drawn
  between …" have different remedies and the comment at `:1026-1028` is right about why.

## Red-Green Proof

**RED prompt/case:** Open Timeline. Type a search string that matches no REQ (e.g. `zzzznope`). Read
`#timeline-summary`. Press `Fit all`. Read `#timeline-summary` and count the drawn rows.

**Why RED now:** the summary correctly reads "No REQ matches the current filters." with no axis and no rows;
after the `Fit all` press it reads "317 REQs in the window 2026-05-27 23:33 UTC → 2026-08-25 04:23 UTC …"
while the chart is still empty and the details table refills with 317 excluded REQs.

**GREEN when:**
- After the `Fit all` press the summary still reads "No REQ matches the current filters." and the table is
  still empty.
- Every toolbar control reports itself disabled in that state.
- Typing in From or to in that state changes nothing and leaves the field showing a value consistent with
  what the readout says.
- Clearing the filter restores a working toolbar on the new render.
- A Node behaviour probe drives the no-match render followed by a toolbar press and fails if the summary
  changes.

**Validation:** Inferred during capture; reproduced in a live render and traced to the wiring sites listed
above.

## Full Context

See `do-work/user-requests/UR-066/input.md`.
