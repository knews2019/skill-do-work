---
id: REQ-346
title: "Narrow the Durations axis to a chosen window"
status: pending
created_at: 2026-08-23T22:37:52Z
user_request: UR-068
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-342, REQ-343, REQ-344, REQ-345, REQ-347, REQ-348, REQ-349, REQ-350]
batch: durations-panel-improvement
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-durations.js
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
---

# Narrow the Durations Axis to a Chosen Window

## What

The Durations axis spends most of its width on idle calendar and there is no way to narrow it. Add a
window control offering the last 30 days, the last 90 days and all history, applied to all three
panels, with the active window stated in the summary line.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

On this repository's archive 66 percent of the 86-day axis has no completions and 92 percent of the
samples fall in the final 30 days. On the consuming project a single April outlier (REQ-407) holds
the left half of the plot open by itself. The axis is linear real time on purpose — compressing idle
gaps would destroy the cadence answer — so the fix is a narrower domain, not a different scale.

## Detailed Requirements

- A window control on the Durations view: **last 30 days**, **last 90 days**, **all history**.
- **The summary line states the active window and the sample count it covers.**
- **The axis stays linear real time inside the window.** Do not compress gaps: the idle-gap cadence
  reading is why it is linear in the first place.
- **Panels B and C follow the same window** as panel A.
- **Reuse the topbar's existing recent-window control pattern** rather than inventing a second one.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first.
- **Reuse the pattern, not the value.** This window is the Durations axis domain; the board's
  recent-window control means "recently done". Sharing one value would move the Durations axis
  whenever someone adjusts a board filter, which is a different statistic wearing the same axes. Keep
  the Durations window as its own state. (Capture decision on the report's open question Q2.)
- Generate a board and look at it at each of the three settings.

## Dependencies

REQ-348's stat tiles are stated for the window actually plotted, so they read whatever state this REQ
establishes.

## Builder Guidance

**Certainty: firm.** The three window options, the linear-inside-the-window rule and the separate
state are all decided. Whether the default is 30 days or all history is yours — say which you chose
and why in the hand-back.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repository and open Durations. The axis spans all 86
days of the archive, 57 of them with no completions, and no control anywhere narrows it.

**Why RED now:** `timeStart`/`timeEnd` are fixed to the full sample range with no window state.

**GREEN when:** a window control offers 30 days, 90 days and all history; selecting one redraws
panels A, B and C over that domain with the axis still linear inside it; and the summary line names
the active window and its sample count.

**Validation:** User confirmed (bundled invocation).

---
*Source: prompt A3, `ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html` (finding F3).*
