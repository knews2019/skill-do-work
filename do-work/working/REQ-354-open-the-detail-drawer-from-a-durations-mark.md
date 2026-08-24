---
id: REQ-354
title: "Open the detail drawer from a Durations mark"
status: claimed
claimed_at: 2026-08-24T17:42:55Z
status_changed_at: 2026-08-24T17:42:55Z
route: B
created_at: 2026-08-23T22:37:52Z
user_request: UR-069
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-346, REQ-347, REQ-348, REQ-349, REQ-350, REQ-351, REQ-352, REQ-353]
batch: durations-panel-improvement
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-24T17:42:55Z
  basis:
    - Route B
    - 4-file anticipated write set
    - 2 subsystems involved
    - 4 acceptance criteria
    - browser evidence
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-durations.js
  - skills/do-work-board/tools/queue-kanban/web/board-detail.js
---

# Open the Detail Drawer From a Durations Mark

## What

A Durations mark is a dead end: the nearest-mark hover writes one line into a status paragraph and
nothing else happens. Make a click open the detail drawer for the nearest REQ, and give the view a
keyboard path to the same information.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Having identified an outlier, the reader has to leave the view and search for its id by hand. The
board and timeline views both open a detail drawer; Durations has no click target and no keyboard
path except the sample table inside a collapsed `<details>`.

## Detailed Requirements

- **A click on panel A's hover surface opens the detail drawer** for the nearest REQ — the same
  drawer the board and timeline views open, not a second one.
- **A keyboard path to the same information** that does not require opening the collapsed sample
  table, and that can reach **every plotted sample** — not only the over-ceiling ones. A route that
  covers outliers alone leaves an ordinary under-ceiling REQ reachable only through the collapsed
  305-row table, which is the exact limitation this REQ exists to remove, and it would leave keyboard
  readers with strictly less reach than mouse readers.
- **Keep the hover readout exactly as it is.**
- **Click and hover resolve to the same mark, always.** A reader must never open a REQ other than the
  one the readout names. If REQ-349's jitter has landed, the click inherits whatever compensation the
  hover uses; if it has not, whichever REQ lands second reconciles the two.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first.
- The Timeline's click regression is the cautionary case: pointer capture on `pointerdown` retargeted
  the synthesized click so no mouse click reached the delegated handler (REQ-336). Verify with a real
  mouse click, not a synthesized `PointerEvent`.
- Keyboard access follows the Timeline's roving-tabindex precedent (REQ-338): one Tab stop for a list,
  arrows to move within it. Do not add one Tab stop per mark — there are 305 of them here.
- Generate a board and look at it.

## Dependencies

None declared. REQ-349 perturbs the nearest-mark maths this REQ reuses; see the requirement above.

## Builder Guidance

**Certainty: firm on the click, open on the keyboard path.** The shape is yours, but completeness is
not: whatever you choose must reach every plotted sample. A roving list over the over-ceiling samples
alone does not qualify, and neither does REQ-351's ranked list on its own — both cover outliers only.
A day-then-mark traversal (arrow between days, arrow between the marks within a day) does, and so
does a roving list over the full sample set; pick one, state why, and make sure it reaches the same
drawer the click opens.

**`tdd: false` deliberately.** The browser probe lane cannot dispatch trusted input today — that is
the whole subject of REQ-341, and it is why REQ-324's click lock-in missed the Timeline regression
and REQ-336's RED had to be reproduced outside the suite over the DevTools Protocol. A structural
check is available and is worth writing; if REQ-341 has landed by the time this is claimed, write the
behavioural probe instead and say so.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repository, open Durations, hover the 10h 55m outlier
until the readout names REQ-064, then click it. Nothing opens. Tab through the view: the only route to
the same information is a collapsed `<details>` holding 305 table rows.

**Why RED now:** The hover surface binds `mousemove` and `mouseleave` only, and there is no keyboard
affordance outside the sample table.

**GREEN when:** clicking a mark opens the same detail drawer the board and timeline views open, for
the REQ the readout names; a keyboard path reaches that drawer without expanding the sample table;
and the hover readout is unchanged.

**Validation:** User confirmed (bundled invocation).

---
*Source: prompt A8, `ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html` (finding F8).*

## Triage

**Route: B** — The click outcome and completeness rule are explicit, but exploration must choose and trace the one-tab-stop keyboard path, the shared detail-drawer entry point, nearest-mark selection, and trusted browser-probe conventions.

## Plan

**Planning not required** — Route B: exploration-guided implementation.
