---
id: REQ-573
title: 'Open the detail drawer from an Activity row and highlight every row of the same REQ'
status: claimed
created_at: 2026-09-04T23:16:00Z
user_request: UR-115
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-05T12:22:40Z
  basis:
    - trivial short-circuit
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
claimed_at: 2026-09-05T12:21:55Z
route: A
dispatch_at: 2026-09-05T12:23:39Z
builder_handback_at: 2026-09-05T12:39:24Z
integration_at: 2026-09-05T12:39:24Z
---

# Open the Detail Drawer from an Activity Row and Highlight Every Row of the Same REQ

## What

Clicking a row on the Activity view does nothing today. Make a click open the same REQ detail drawer the Board opens when a card is clicked, and mark every Activity row that carries the same REQ id as selected, so the reader can scan the whole sequence of states that REQ went through to arrive where it is.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Builder read the crew rules, the kanban prime and both lesson satellites, then the five client fragments and the fragment execution-order manifest, and settled a six-step approach: a real button carrying the detail attributes, the drawer's own state as the single selection source, one highlight function, a click path that reads the clicked row rather than the drawer, three non-colour signals, and a mutation check. Recorded under `## P-A-U` in `do-work/runs/work-2026-09-05-120117/REQ-573-handback.md`.
- [x] **[APPLY]:** One commit on the builder branch (`05484aa6`) touching exactly the three Scope files; 350 insertions, 2 deletions.
- [x] **[UNIFY]:** `git diff --stat` reviewed (3 files); linters clean; debug-artifact scan over added lines empty; two draft defects were caught in this phase and fixed before the commit. Per-file checks listed in the hand-back.

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

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the delegation mechanism to reuse (`data-detail-kind`), the three files to touch, and a captured RED/GREEN pair driving the Node behavior lane. `effort_estimate: effort-mechanical`. Nothing about the location or the pattern needs discovery.

**Planning:** Not required

## Plan

**Planning not required** - Route A: direct to builder

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-activity.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modified)

**What was done:** The Activity row's REQ cell is now a real button carrying the drawer's own `data-detail-kind` and `data-detail-id` attributes, so the existing document-level delegation opens the drawer with no second opener and the cell is reachable by keyboard the same way a Board card is. Selection is read live from the drawer's own state rather than stored a second time, which is what makes a re-render restore the highlight, the close button clear it, and a UR drawer unable to select REQ rows — none of those needed code of their own. Every row of the clicked request gets a selected class carrying three signals: a background tint, an inset bar down the REQ column, and the bold underlined id. Merge range `45a9010d..2d3981f4`; builder branch head `05484aa6`. Builder-authored `## Decisions` (D-01 to D-07) and `## Discovered Tasks` live in `do-work/runs/work-2026-09-05-120117/REQ-573-handback.md`.

## Qualification

**Passed.** Read from the merge range `45a9010d..2d3981f4`.

- One selection source. `selectedActivityRequestId()` reads the drawer's kind and id, so there is no module-level copy that can disagree with the drawer. Three of the request's own requirements fall out of that single read instead of being implemented separately.
- The click path reads the clicked row rather than the drawer, and the builder proved why from the fragment execution-order manifest rather than assuming: this fragment's document listener is registered before the delegation that opens the drawer, so at click time the drawer still names the previous request. Every non-row click falls back to the drawer read, which is what makes the close button clear the highlight.
- The keyboard path is a real `<button>`, not a click handler on a table cell, which is what the request asked for.
- The highlight does not depend on colour alone: tint, an inset bar and a bold underlined id, with both tokens defined in the light and dark palettes. The bar is an inset shadow rather than a border because a real border would shift the row under collapsed table borders.
- One cross-REQ test break, handled correctly: REQ-572's activity-summary test needed its slice list and stub node extended once the renderer started building a button. Not one assertion was changed, weakened or removed — the transition and request counts, the repeated-id contract and both empty states are pinned exactly as REQ-572 left them.

Requirements traced: detail attributes on every row, the whole matching set selected and others cleared, selection surviving a re-render, the highlight readable without colour, closing clearing it, and keyboard reachability. No new definition of a stamp and no client-side re-sorting.

*Checked by work action*

## Testing

**Focused tests (post-merge, main tree at `2d3981f4`):**
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' ./...` — exit 0, 58 tests, wall 9s, slowest file `citations_test.go` at 2.42s against the 30s per-file budget.

**Red-green validation** (traced to `## Red-Green Proof`): RED and GREEN observations, plus the builder's mutation check proving the new assertions bite, are recorded in the hand-back's `## Test evidence`.

**Not covered — browser render (builder decision D-04, escalated rather than claimed):** no browser was driven, so the tint, the inset bar and the focus ring were not seen against either palette. Every token used is defined in both palettes and every rule copies an existing pattern in the same stylesheet. The orchestrator routed this to the `queue-kanban-browser` heavy lane rather than a one-off screenshot, since that lane is already selected for this request and runs at the queue-exhaustion drain.

