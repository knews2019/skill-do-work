---
id: REQ-332
title: "Take the timeline toolbar out of service when no REQ matches the filters"
status: completed
created_at: 2026-08-23T12:22:00Z
claimed_at: 2026-08-23T16:00:00Z
completed_at: 2026-08-23T16:35:00Z
commit: ce0deb2
route: B
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
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`. Chose to keep `onclick` and clear it, rather than move the toolbar to `addEventListener` — two registrations for one button is the failure mode the teardown note in this file already warns about. One selector naming every owned control, so a future control inherits the rule instead of being added to a list.
- [x] **[APPLY]:** `web/board-timeline.js` and `generate_test.go` as planned, plus a `timeline-zoom` class on the zoom group in `web/template.html` so the selector can name it, and `timeline_browser_probe_test.go` for the flake fix below.
- [x] **[UNIFY]:** Verified:
  - `web/board-timeline.js` — `node --check` clean; the retire/enable pair is symmetric and the enable sits past the early return, where the render owns the controls.
  - `web/template.html` — the added class has no CSS rule, so it is a selector hook and changes nothing visually (checked `board.css`).
  - `generate_test.go` — the new probe carries a vacuity guard: if the stub wires no control, the handler half of the test cannot fail, and it says so.
  - `bash _dev/tests/maintainer-verify.sh` exit 0.
  - Mutations: no retirement, handlers cleared but controls left pressable, and never re-enabled — three distinct failures.
  - Live render: 11 controls disabled with nothing matching, summary unchanged after pressing Fit all, all 11 restored on clearing the filter.
  - **A flake I introduced, found by the gate and fixed here, not quarantined:** the drawer probe passed alone and failed inside the full suite. Its fixed 300ms waits were adequate on an idle machine and not on a loaded one. Both new probes now poll for the condition; `setTimeout` rather than `requestAnimationFrame`, because headless `--dump-dom` has no compositor to drive frames and an rAF poll never resolves. The drawer mutation was re-run after the change to confirm the probe still bites.

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
