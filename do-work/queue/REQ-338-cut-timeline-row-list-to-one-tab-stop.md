---
id: REQ-338
title: "Cut the Timeline row list to one Tab stop"
status: pending
created_at: 2026-08-23T18:30:26Z
user_request: UR-067
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-336, REQ-337]
batch: timeline-click-regression
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# Cut the Timeline Row List to One Tab Stop

## What

Every Timeline REQ row is its own Tab stop, so tabbing past the row list takes one press per row
(29 on the observed board). Make the row list a single Tab stop with arrow-key movement between
rows (roving tabindex), so Tab escapes the list in one press.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Under REQ-333 the user refuted this as a keyboard trap — Tab does escape — and left it. During
this capture the user chose to fix the tedium anyway: a keyboard user pays one press per row to
get past the chart. Not urgent, but noticed in real use.

## Context

REQ-333 (commit 36c4518) established the current keyboard contract: rows are focusable, Enter
and Space on a focused row open the drawer, and arrow keys pan the chart without killing the
keyboard path. Its F3 finding recorded the 29-press escape as "not a trap, leave it" — this REQ
supersedes that "leave it" by the user's explicit choice at capture.

## Detailed Requirements

- The row list occupies one Tab stop: Tab from any focused row leaves the list in one press
  (in each direction).
- Arrow keys move focus between rows (roving tabindex).
- Enter and Space on the focused row still open the drawer — do not regress the drawer lane
  REQ-336 restores.
- Do not break REQ-333's contracts: arrow-key panning with a row focused must keep working
  (the builder decides how row-movement keys and pan keys share the keyboard — e.g. which axis
  moves rows vs pans), and focus must survive the node rebuilds REQ-333's focus restore handles.

## Builder Guidance

Certainty: Firm on the shape (one stop, roving tabindex — user picked this option at capture);
builder's latitude on key assignment and on how the roving index interacts with the existing
focus-restore and pan handlers in `web/board-timeline.js`.

## Red-Green Proof

**RED prompt/case:** In the probe lane, generate a board with multiple REQ rows, focus the first
Timeline row, and press Tab: focus lands on the next row, not past the list — escaping takes as
many presses as there are rows (29 observed).
**Why RED now:** Every row carries its own tabindex, so each is a distinct Tab stop.
**GREEN when:** A probe-lane test asserts: Tab from a focused row exits the row list in one
press, arrow keys move focus between rows, and Enter/Space on the focused row still opens the
drawer. Test written first and confirmed failing before the change (`tdd: true`).
**Validation:** User confirmed (chose the roving-tabindex option during capture)

## Dependencies

None. Overlaps REQ-336 on `web/board-timeline.js` (declared in both write sets) — serial
execution or coordination needed if dispatched in parallel.

---
*Source: UR-067 — see `do-work/user-requests/UR-067/input.md` for complete verbatim input.*
