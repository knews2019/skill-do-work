---
id: REQ-573
title: 'Open the detail drawer from an Activity row and highlight every row of the same REQ'
status: pending
created_at: 2026-09-04T23:16:00Z
user_request: UR-115
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-572]
related: [REQ-572]
batch: activity-history
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-activity.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
---

# Open the Detail Drawer from an Activity Row and Highlight Every Row of the Same REQ

## What

Clicking a row on the Activity view does nothing today. Make a click open the same REQ detail drawer the Board opens when a card is clicked, and mark every Activity row that carries the same REQ id as selected, so the reader can scan the whole sequence of states that REQ went through to arrive where it is.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- Every Activity row (or its REQ cell) carries `data-detail-kind="req"` and `data-detail-id="<REQ id>"`, so the existing document-level click delegation in `board-controls.js` opens the drawer with no new handler. The drawer is the one in the screenshot: title, status, domain, user request, depends on, write set, created, claimed, tree, then the REQ body.
- On click, every row whose `data-activity-request` equals the clicked id gets a selected class; rows of other REQs lose it. Re-rendering the table (window change, filter change) keeps the selection for the currently open REQ.
- The highlight must be readable in light and dark themes and must not depend on color alone (a left border or bold id beside the background is enough).
- Closing the drawer clears the highlight.
- Rows are keyboard reachable the same way Board cards are (the REQ cell is a button or link, not a bare `td` with a click handler).

## Constraints

- No new definition of what a stamp means and no re-sorting in the client; this REQ touches row behavior and styling only.
- Reuse `openDetail("req", id)` through the `data-detail-kind` delegation; do not add a second drawer opener.

## Dependencies

Depends on REQ-572 (one Activity row per lifecycle stamp): highlighting sibling rows only means something once a REQ can have several rows.

## Red-Green Proof
**RED prompt/case:** In the Node behavior lane (`javascript_behavior_*_test.go`), render the Activity view with two rows for REQ-570 and one for REQ-505, dispatch a click on the REQ-570 row, and assert the detail drawer shows `REQ-570` and both REQ-570 rows carry the selected class while the REQ-505 row does not.
**Why RED now:** `renderActivity` in `board-activity.js` builds plain `tr`/`td` elements with no `data-detail-kind` attribute and no selection state, so the click does nothing and no row is marked.
**GREEN when:** The click opens the drawer for REQ-570 and exactly the rows with `data-activity-request="REQ-570"` are marked selected; on the running board, clicking any REQ-570 row opens the same drawer the Board shows for it and every REQ-570 row is visibly highlighted.
**Validation:** User confirmed (verify-requests, 2026-09-04)

## Builder Guidance

The user wrote "left side" for where the REQ should show up; the screenshot shows the drawer on the right. Treat this as "the existing detail drawer, wherever it renders", not as a request to move it.

## Assets

- `do-work/user-requests/UR-115/assets/REQ-573-screenshot-2-board-detail-drawer-req-570.png`: the Board view with the detail drawer open on the right for REQ-570 ("[impact-rule-change] Delete the pending-heavy-testing status; held requests stay claimed"). The drawer lists Status claimed, Domain general, User request UR-114, Depends on REQ-507 (pending), Unblocks REQ-571, a long Write set, Overlapping write sets (REQ-506, 507, 510, 544, 556, 557), Route C, Impact impact-rule-change, Effort estimate effort-substantive, Created Sep 4 22:52 UTC, Claimed Sep 4 23:00 UTC with a running stopwatch, Tree working, then the REQ body. On the left the Board columns show Pending 17, Claimed 1 (REQ-570), Needs input 0, Recently done 33.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`, so bare only): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban views". Over the budget on its own.

## Full Context
See `do-work/user-requests/UR-115/input.md` for complete verbatim input.

*Source: "make sure that if I click on a req it does show up on the left side just as it shows up in the board and also highlights all similar REQ-ID entries, so it can be visually scanned all the status the REQ went through to arrive there."*
