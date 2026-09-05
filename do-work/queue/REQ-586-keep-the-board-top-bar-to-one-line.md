---
id: REQ-586
title: 'Keep the board top bar to one line: single-line identity and Touched-in chips inside the Activity view'
status: pending
created_at: 2026-09-05T12:40:00Z
user_request: UR-121
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-585]
related: [REQ-585, REQ-573]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/board-activity.js
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
---

# Keep the Board Top Bar to One Line: Single-Line Identity and Touched-In Chips Inside the Activity View

## What

The top bar grows whenever its three control groups (filters, views, Touched-in) no longer fit beside the identity block: the controls wrap to two rows and the identity block (`do-work/ queue`, project, "Generated", date) wraps to four lines, so a 68 px bar becomes about 150 px. That space matters most on the Activity view, where the reader wants rows on screen to click a REQ and see every row of it highlighted (REQ-573). Two changes, both chosen by the user:

1. **O1, one-line identity.** Render the identity as one `nowrap` line, `do-work/ queue · skill-do-work2 · 12:17 UTC`, with the full "Generated 2026-09-05 12:17 UTC · 37s ago" text kept in a `title` tooltip on that line. The bar no longer grows when the controls beside it wrap.
2. **O2, Touched-in chips move into the Activity view.** Delete the `#activity-window-group` pill from the top bar and render the same four chips (6h, 24h, 48h, 7d) on the Activity view's summary line, so it reads "236 transitions across 49 REQs in the last [6h] [24h] [48h] [7d]". The top bar keeps two groups (filters, views) and fits on one row at far narrower widths. The Timeline already keeps its period controls inside its view, so this follows the existing pattern.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- The chips keep their ids, `data-activity-window` values, `aria-pressed` handling, and the `setActiveButton` call in `board-activity.js`; only their home moves. `board-controls.js` line 46 (`document.getElementById("activity-window-group").hidden = viewState.view !== "activity"`) becomes unnecessary once the group lives inside `#view-activity`, which is hidden with the view; delete it rather than keeping a no-op.
- The summary line and the chips stay on one line at desktop widths; the chips sit after the count text, styled as the same pill group, not as a second bar.
- With REQ-585 (give the Activity view one scroll surface) landed, the summary line is the natural place for the chips whether or not it is pinned (M3); if the pinned variant was chosen, the chips are pinned with it.
- The one-line identity keeps the wordmark, project name and time in that order, separated by a middle dot, and keeps the existing `id="board-project"` and `id="board-generated"` hooks that `board.js` (line 60 reads `board-generated`) and the serve-mode refresh write into; check for readers before renaming anything.
- The `@media (max-width: 760px)` rule that stacks the top bar vertically stays as it is; this REQ is about the widths above it.

## Constraints

- No new control, no new state: the window values, the default of 24h, and the persistence behavior stay exactly as today.
- Board version and changelog per `_dev/primes/prime-kanban-board.md`.

## Dependencies

Depends on REQ-585 (give the Activity view one scroll surface): both edit the Activity summary line and the same block of `board.css`, so this builds after that lands rather than beside it.

## Red-Green Proof
**RED prompt/case:** In the Node behavior lane (`javascript_behavior_*_test.go`), load the board with an Activity payload, switch to the Activity view, and assert that the element carrying the `data-activity-window="24"` button is a descendant of `#view-activity` and not of `.board-topbar`; then click the `48h` button and assert the summary text says "in the last 48 hours" and the `48h` button has `aria-pressed="true"`.
**Why RED now:** `template.html` line 94 places `#activity-window-group` inside `.board-topbar`'s `.board-controls`; the descendant assertion fails.
**GREEN when:** The assertion passes, the top bar contains exactly two `.control-group` pills on every view, and on the served board at a width where the old bar wrapped (about 1400 px), the top bar is one row of 68 px with the identity on a single line and the chips visible on the Activity summary line. The one-line identity is a layout fact the Node lane cannot see; verify it with a browser probe or a captured screenshot in the hand-back.
**Validation:** User confirmed (chose O1 and O2 from three options, 2026-09-05)

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`, so bare only): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban views". Over the budget on its own.

## Assets

- `do-work/user-requests/UR-121/assets/REQ-586-screenshot-1-top-bar-wrapped.png`: the top bar of the M1 mockup (`ai-reports/2026-09-05_1520_activity-view-double-scroll-mockups/mockups/m1-one-page-scroll.html`) in a narrow frame. Left, the identity block on four lines: "do-work/", "queue", "skill-do-work2", "Generated 2026-09-05" and "12:17 UTC" on a fifth. Right, the controls on two rows: the filter pill (Filter id or title, All domains, All statuses) above, the view pill (Board, Calendar, Durations, Timeline, Activity selected, Testing) and the Touched-in pill (6h, 24h selected, 48h, 7d) below. The mockup copies the shipped top bar rules, so the real board wraps the same way at the width where its three groups stop fitting on one row.

## Full Context
See `do-work/user-requests/UR-121/input.md` for complete verbatim input.

*Source: "this part of the header is still taking up too much vertical space, and that is precious when I want to click a req and I want to highlight all of it's occurances" / "ok, do o1 and o2 capture it first"*

## Addendum (2026-09-05)

User added:

> ````text
> while we are at it the order should be Board, Activity, Calendar, Timeline, Durations
> 
> testing can remain last
> ````

- Reorder the view buttons in the top bar to: Board, Activity, Calendar, Timeline, Durations, Testing. The order is declared once, in `template.html` lines 71 to 88 (the `data-view-target` buttons); `board-controls.js` reads the buttons from the DOM, so nothing else keeps the list.
- Asset: `do-work/user-requests/UR-122/assets/REQ-586-screenshot-2-view-tab-order.png`, the view pill today: Board, Calendar, Durations, Timeline (selected, with focus ring), Activity, Testing.
- Same proof lane as the rest of this REQ: the behavior test can assert the `data-view-target` values of the pill's buttons in document order.
