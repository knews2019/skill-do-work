---
id: REQ-350
title: "Narrow the Durations axis to a chosen window"
status: claimed
claimed_at: 2026-08-24T15:21:20Z
status_changed_at: 2026-08-24T15:21:20Z
route: B
created_at: 2026-08-23T22:37:52Z
user_request: UR-069
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-346, REQ-347, REQ-348, REQ-349, REQ-351, REQ-352, REQ-353, REQ-354]
batch: durations-panel-improvement
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-24T15:21:20Z
  basis:
    - Route B
    - 3-file seeded write set
    - 5 acceptance criteria
    - browser evidence across three settings
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-durations.js
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Narrow the Durations Axis to a Chosen Window

## What

The Durations axis spends most of its width on idle calendar and there is no way to narrow it. Add a
window control offering the last 30 days, the last 90 days and all history, applied to all three
panels, with the active window stated in the summary line.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route B exploration traced a Durations-only state/control, a latest-completion-anchored whole-day domain, and one projected samples/days slice feeding all panels, summaries, hover indexes, and the table. Add focused production-renderer behavior tests and inspect all three windows.
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

REQ-352's stat tiles are stated for the window actually plotted, so they read whatever state this REQ
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

## Triage

**Route: B** — The three states, reusable control pattern, affected renderer, and rendered outcome are explicit. Focused exploration is required to trace the independent state, filtering/domain boundaries, and browser-probe conventions before implementation.

## Plan

**Planning not required** — Route B: exploration-guided implementation.

## Exploration

- Reuse the existing `control-group`, `control-button`, active class, and `aria-pressed` pattern with a distinct `data-durations-window` attribute. Never reuse the board's `data-window-hours` state.
- Keep a whitelisted Durations window state in `board-durations.js`, defaulted to the last 30 days. Anchor 30/90-day whole-UTC-day domains to the latest completion rather than wall-clock time so a static board does not age into an empty chart; `all` preserves the current full-history domain.
- Project samples and precomputed days once through `[timeStart, timeEnd)`, then feed that projection to all three panels, UR lane, colour context, hover indexes, summary, and table. The epoch-to-x mapping remains affine, so idle gaps retain their real width.
- Existing control CSS is sufficient. No payload, Go aggregation, `board.js`, or browser-probe change is required.
- `generate_test.go` must be in scope: renderer behavior needs sample/day/table/count and linear-domain proof, and the existing 400-day all-history geometry fixture must explicitly choose `all` under the new 30-day default.

## Scope

- `skills/do-work-board/tools/queue-kanban/web/board-durations.js`
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js`
- `skills/do-work-board/tools/queue-kanban/web/template.html`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`
