---
id: REQ-390
title: 'Replace the timeline''s Day/Week/Month periods with trailing windows: last day, last 7/30/90/all days'
status: pending
created_at: 2026-08-27T14:15:08Z
user_request: UR-079
domain: frontend
prime_files: ['_dev/primes/prime-kanban-board.md']
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
---

# Replace the Timeline's Day/Week/Month Periods with Trailing Windows

## What
On the board's Timeline view, replace the period toolbar's calendar-period buttons (Day, Week, Month) with trailing windows ending at now: Last day, Last 7 days, Last 30 days, Last 90 days, and All days (the full recorded data range).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context
- The current toolbar is the `timeline-periods` control group in `skills/do-work-board/tools/queue-kanban/web/template.html` (buttons with `data-timeline-period="day|week|month"`), wired by `applyPeriodWindow` / `applyPeriodStep` in `web/board-timeline.js`. Those set the window to one calendar day/week/month around now (or around the panned-to point) and ‹ / › step one calendar period.
- The filter dropdown already uses trailing-window vocabulary ("Last 7 days", "Last 30 days") — the new button labels match that existing wording.
- Shipped prose and ARIA text describe the calendar-period behavior (the view panel's `aria-label`, the toolbar comment, and the `timeline-hint` paragraph in `template.html`) and must be updated to describe the trailing windows.
- Builder latitude: what ‹ / › step by under a trailing window (the window's span is the natural choice), and how "around wherever you have panned" interacts with windows now anchored at now, are the builder's call — the existing custom-span screenful stepping is precedent. "Last day" reads as the trailing 24 hours ending at now.
- The `data-timeline-period` values feed `setActiveButton` state restore and the browser probe tests (`timeline_browser_probe_test.go` clicks `day`/`week`/`month`); those tests change with the control set.

## Red-Green Proof
**RED prompt/case:** Generate a board and open the Timeline view: the period toolbar offers Day, Week and Month, each setting the window to one calendar period — there is no one-click way to see the trailing last 7/30/90 days or the whole recorded range. A probe in the `timeline_browser_probe_test.go` harness asserting a `[data-timeline-period]` control set of last-day/7/30/90/all-days buttons, and that clicking "Last 30 days" yields a window ending at the board's now and spanning 30 days, fails today.
**Why RED now:** `template.html` ships only `data-timeline-period="day|week|month"` and `applyPeriodWindow` computes calendar periods, not trailing spans.
**GREEN when:** The toolbar offers Last day, Last 7 days, Last 30 days, Last 90 days and All days; clicking one sets the window to that trailing span ending at now (All days spans the full recorded range), the probe above passes, and no Day/Week/Month calendar-period button remains.
**Validation:** Inferred during capture

## Assets
`do-work/user-requests/UR-079/assets/REQ-390-screenshot-1-timeline-period-toolbar.png` — the Timeline view of this repo's own queue (generated 2026-08-27 14:12 UTC, window 2026-08-24 → 2026-08-29). Top-right toolbar shows the current controls: ‹ Day [Week] Month › with Week active and "part of one week" as the period state, alongside the From/to date fields and the − + Now Fit-all zoom group. The request replaces that Day/Week/Month group.

---
*Source: instead of day/week/month let's have last day, last 7/30/90/all days*
