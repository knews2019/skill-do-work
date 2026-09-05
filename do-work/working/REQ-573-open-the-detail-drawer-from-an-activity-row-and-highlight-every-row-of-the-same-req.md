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
review_at: 2026-09-05T14:22:50Z
heavy_verified_at: 2026-09-05T14:22:50Z
heavy_verified_revision: 7b2673b690a671ccb360c26b0c19c56ecc7356b5
commit: 2d3981f4
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


## Review

**Overall: 79%** | 2026-09-05T12:45:07Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 80% |
| Test Adequacy | 80% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 — Opening a REQ drawer from anywhere outside the Activity table leaves the highlight on the PREVIOUSLY open REQ, so the table and the drawer name different REQs. `web/board-activity.js:172-177`: when `closest("[data-activity-request]")` misses, `syncActivitySelectionToClick` falls back to `selectedActivityRequestId()` — read in the same instant that D-02 correctly identifies as stale, because `board-activity.js` (manifest position 7) runs before the `[data-detail-kind]` delegation in `board-controls.js` (position 10). D-02 fixes this for the row branch and leaves it in the fallback branch. Nothing recomputes afterwards: `renderActivity()` has exactly three callers (`board-activity.js:31` window buttons, `board-filters.js:190` filter change, `board-controls.js:74` guarded by `renderedOnce.activity`), so the wrong highlight survives until a filter or window change. Reproduced with a probe replaying both listeners in registration order over the shipped function bodies: opening REQ-570 from a drawer ticket link highlights nothing, then opening REQ-505 from the drawer highlights REQ-570's two rows while the drawer shows REQ-505. Same for a Board card click followed by a switch back to an already-rendered Activity view. Cheapest fix: give the fallback branch the delegation's own timing (defer the drawer read by a task, or have `openDetail` call `syncActivitySelectionToDrawer`); moving the fragment after `board-controls.js` also works and breaks nothing, since the row branch never reads the drawer. — impact-user-visible → report only

**Minor findings:**
- F2 — Clicking any non-REQ cell of a row (title, status, transition, when, stamp) highlights that REQ's whole set but never opens the drawer, because only the REQ cell carries `data-detail-kind`. The highlight and the drawer then disagree until the next re-render snaps the highlight back. The REQ's Detailed Requirements permit the cell-only opener ("Every Activity row (or its REQ cell)"), and the row-wide highlight matches the "clicked id" wording, so this is on-spec — but it contradicts D-01's claim that the drawer is the only selection state, and the UR asked for a click on a REQ to open the detail the way a Board card does, where the whole card is the target. — impact-user-visible → report only
- F3 — The tint alone is close to invisible in both palettes. Composited over the panel surface, `--tint-claimed` measures about 1.26:1 (dark, rgba(111,156,230,0.14) over #171c26) and about 1.18:1 (light, rgba(58,107,196,0.12) over #ffffff) — far under the 3:1 that a non-text state indicator wants. The state is carried by the 3px inset `--accent-claimed` bar (6.2:1 dark, 5.2:1 light against the same surface) and the bold underlined id, both of which sit only in the leftmost column, so a reader scanning the "when"/"stamp" columns of a wide table has essentially no signal. The REQ's own bar is met (not colour-alone, both palettes defined), so this is a strengthening suggestion, not a miss: roughly double the tint alpha, or extend a signal across the row. — impact-user-visible → report only
- F4 — The selection has no programmatic signal. No `aria-selected`, `aria-current`, or `aria-pressed` lands on the selected `tr` or its button, so a screen-reader user gets the drawer content but nothing telling them which rows belong to the open REQ. The colour-independence requirement is satisfied visually only. `web/board-activity.js:158-168`. — impact-user-visible → report only
- F5 — The hand-back's second Discovered Task proposes a test pinning `board-activity.js` ahead of `board-controls.js` in `boardJavaScriptFragmentPaths`. That would lock in the worse of the two orders. Under the reverse order the row branch is unchanged (it reads the clicked row and never the drawer) and the fallback branch becomes correct, which is F1's fix. Fix the fallback rather than pinning the order. — impact-rule-change → report only
- F6 — The new test asserts a path production never takes, which is what hid F1. `javascript_behavior_c_test.go:2812-2816` sets the drawer state and calls `syncActivitySelectionToDrawer()` directly, then asserts "a drawer opened from somewhere other than the table did not mark its rows". No production code calls `syncActivitySelectionToDrawer` on a drawer open — its only caller is the Escape `keydown` listener. A real click on a `[data-detail-kind]` element outside the table goes through `syncActivitySelectionToClick`'s fallback and produces the wrong set. The case is not vacuous (it does bite on the kind check), but it is mis-labelled and covers a hole it does not test. — impact-negligible → report only

**Nit findings:**
- F7 — Restatement Sweep result. `ai-reports/2026-09-05_1520_activity-view-double-scroll-mockups/mockups/mockup.js:40-48` builds the Activity REQ cell as plain text on a bare `th`, with `data-activity-request` but no detail button. That mockup set was authored after this REQ's dispatch and is what a later double-scroll change would be built from, so it would silently drop this REQ's opener. Out of this REQ's declared scope; recorded here per Restatement Sweep step 3. — impact-negligible → report only
- F8 — The Activity button's accessible name is the bare id. `makeRequestCard` (`web/board-cards.js:123`), the pattern the diff says it mirrors, sets `aria-label` to "REQ-570: title". In table-reading mode the adjacent title cell supplies it, so nothing is lost in practice. — impact-negligible → report only

**Sweep record (Step 6 Restatement Sweep):** The diff redefines the Activity REQ cell's content shape (text node → button element) and introduces `activity-req-button` / `is-activity-selected`; `data-activity-request` semantics are unchanged. Swept all five tokens repo-wide. Consumers found and verified: `web/template.html:480` (column header only, unaffected), `javascript_behavior_c_test.go:2464-2516` (REQ-572's test, updated correctly — see below), `_dev/primes/prime-kanban-board.md` (no restatement of the activity row shape or of listener order). One stale restatement found — F7.

**Cross-REQ test edit (D-05) verified clean.** `git diff --stat` records 2 deletions repo-wide and both are in `web/board-activity.js`; the test file has zero deleted lines, so no assertion could have been removed or rewritten. The edit to `TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests` is three added slice targets (`createElement`, `selectedActivityRequestId`, `applyActivitySelectionHighlight`), a no-op `classList.toggle` and a real `getAttribute` on the stub node, and two empty drawer-state variables. Every original assertion — transition count, distinct-REQ count, the repeated `data-activity-request` contract, both empty-state strings — is byte-identical and still runs against the shipped `renderActivity`. The no-op `toggle` cannot make anything vacuous because that test asserts no class state. Traceability to REQ-572 is stated in a comment in the file and in D-05.

**Escape and close paths (D-03) verified from source.** `closeDrawer` has exactly two callers: `board-controls.js:249` (`#detail-close`, an element-level listener that runs before any document bubble listener) and `board-detail.js:672` (`onDetailPanelKeydown`, registered with `capture = true` at `board-detail.js:650`, so it runs before the fragment's bubble `keydown`). Both leave `currentDetailKind`/`currentDetailId` empty by the time `syncActivitySelectionToClick`/`syncActivitySelectionToDrawer` reads them, so no close path leaves a stale highlight. The one leftover-highlight case is F2's — a highlight created without a drawer at all.

**Table semantics.** A `<button>` inside `<th scope="row">` does not break the row-header role; the header's accessible name is computed from its contents, and the `headers` attribute association is untouched. No finding.

**Acceptance:** Partial — Node behavior lane re-run post-merge (`go test -run '^TestJavaScriptBehaviorActivity'`, 3 tests, PASS, 2.8s); drawer opening, row-set selection, re-render restore, both close paths and keyboard reachability verified from source; opening the drawer from outside the table produces a wrong highlight (F1, reproduced); no render evidence in either palette.
**Suggested testing:** 5 items
**Follow-ups created:** None (8 findings report only)

### Suggested Additional Testing

- Browser render in both colour schemes: the tint against the panel surface, the 3px inset bar on a `th` under `border-collapse: collapse` (the case D-04 named as most likely to surprise), and the focus ring at `outline-offset: 2px` inside the scroll container.
- Keyboard walk: Tab into the table, confirm the ring is visible and not clipped by `.activity-table-scroll`'s `overflow: auto`, and that Enter and Space both open the drawer.
- The F1 repro by hand: open the Activity view, open a REQ drawer from a Board card or a ticket link inside the drawer body, and check which rows are highlighted.
- Screen-reader pass (VoiceOver table mode) on a selected row — confirms F4's gap and that the button does not disturb row-header announcement.
- A REQ with many rows (REQ-572's own follow-up suggests nine) over a 7-day window, to see whether a left-edge-only signal is enough to track a set down a long table.

*Reviewed by review-work action*

## Remediation

The review returned Acceptance **Partial** at 79 percent on one defect this request introduced (F1): a request drawer opened from anywhere outside the Activity table left the table highlighting the previously open request, because the fallback branch read the drawer's identity from a listener registered ahead of the delegation that sets it. The orchestrator judged that too poor to ship and dispatched a narrow remediation for F1 and F4 rather than archiving the finding as report-only.

**That remediation was superseded before it could be integrated.** While it was building, another session working this same checkout fixed the same defect on main as commit `7443fe11` (released as 0.295.1), by a better route: the drawer's own `setDetailTarget` is now the single writer of the selected identity and tells the Activity view directly, and *both* document click listeners are gone — so the listener-ordering hazard that caused F1 no longer exists rather than being worked around. That commit also carries its own probe, which drives the writer instead of assigning the variables in the order the old listener assumed, plus a source pin that fails if a second writer appears.

The remediation branch `worktree-agent-REQ-573-activity-drawer` (head `9de57c92`) keeps a click listener and would have reintroduced the mechanism main deliberately removed. It is abandoned at the maintainer's decision, unmerged and undeleted — `git branch -d` would refuse on unmerged work, and forcing it would destroy the only evidence that the integration did not happen.

**What ships:** this request's delivered behaviour — the Activity REQ cell as a real keyboard-reachable button carrying the drawer's `data-detail-kind`/`data-detail-id` pair, and every row of the clicked request marked selected with three non-colour signals — plus main's fix for how the selected identity is resolved. Verified green in the heavy queue-kanban lanes at `7b2673b6`, which includes `7443fe11`.

**Still open, report only:** F4 (no `aria-current` or `aria-selected` on the selected rows, so a screen-reader user has no programmatic signal), and F2, F3, F5, F6, F7, F8. F4 is a real accessibility gap; it stays a report-only finding because the rule reserves automatic follow-ups for impact-critical work, and because two other sessions are actively editing those files right now.

## Lessons Learned

**What worked:** Making the REQ cell a real `<button>` carrying the existing delegation's attributes bought drawer-opening and keyboard reachability with no second opener. Reading the selection from the drawer instead of storing a copy made three of the request's requirements fall out with no code of their own.

**What didn't:** Deriving correctness from listener registration order. The builder read the fragment execution manifest, correctly worked out that its listener runs before the delegation, and built the row branch around that fact — then left the fallback branch reading the stale value the same argument had just identified. An argument that a value is stale is an argument not to read it anywhere, not a licence to read it in the other branch. Main's fix removed the ordering dependency entirely, which is why it is the better one.

**Worth knowing:** Two sessions fixed this same defect concurrently in the same checkout, from the same review report, within about twenty minutes of each other. Nothing in the queue made that visible: a report-only finding is not a claim, so nothing stops a second reader from picking it up. That is the cost of the report-only intake brake, and it is worth knowing before the brake is tuned.

## Orientation

Clicking a request on the board's Activity view now opens the same detail drawer the Board opens, and every row belonging to that request lights up, so the whole path a request took can be scanned at a glance. Lives in the queue-kanban board subsystem (`_dev/primes/prime-kanban-board.md`), in the Activity client. The REQ cell is a real button, so it is reachable by keyboard the way a Board card is. No prime was made stale.

## Heavy Verification Result

- **Target revision:** 2d3981f4
- **Execution revision:** 7b2673b690a671ccb360c26b0c19c56ecc7356b5
- **Run at:** 2026-09-05T14:22:50Z, from a detached worktree

| Lane | Exit | Wall | Disposition |
| --- | --- | --- | --- |
| `queue-kanban-javascript` | 0 | 9s | executed |
| `queue-kanban-browser` | 0 | 141s | executed |
| `staged-skills` | 0 | 44s | executed |

Every lane this request selected was present in the run, exited 0, and none was skipped. The execution revision includes `7443fe11`, so the lanes verified the shipped combination of this request's markup and main's identity-resolution fix, not this request's mechanism in isolation.

