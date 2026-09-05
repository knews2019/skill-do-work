# REQ-573 remediation hand-back — findings F1 and F4

REQ-573 is "open the detail drawer from an Activity row and highlight every row of
the same REQ". The review passed it at 79% and named eight findings. This
remediation fixes exactly two of them, F1 and F4. F2, F3, F5, F6, F7 and F8 were
out of scope and are untouched.

## Branch

- Branch: `worktree-agent-REQ-573-activity-drawer`
- Head: `9de57c92` — `[REQ-573] remediate F1 and F4: resolve the Activity highlight from the click, and mark selected rows with aria-current`
- Parent: `b169396e` (local `main`, merged in as a fast-forward before any edit)
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905-1201/worktree-agent-REQ-573-activity-drawer`

Note on the merge step: the dispatch brief's `git merge origin/main 2>/dev/null || git merge main`
did **not** put me on the integrated tree. `origin/main` is a stale remote pointer
that was already an ancestor of my branch, so the first command reported "Already
up to date" and succeeded, and the `||` fallback never ran. I ran `git merge main`
explicitly, which fast-forwarded from `05484aa6` to `b169396e`. Any other builder
given that same command line is on the pre-integration tree without knowing it.

## File manifest

| File | Action |
|---|---|
| `skills/do-work-board/tools/queue-kanban/web/board-activity.js` | modified — F1 fix in `syncActivitySelectionToClick`, F4 fix in `applyActivitySelectionHighlight` |
| `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` | modified — one new test; `removeAttribute` added to two existing stub nodes |

`skills/do-work-board/tools/queue-kanban/web/board.css` was in the permitted write
set but needed no change. The selected-row styling already keys on
`.is-activity-selected`, which the fix keeps, and `aria-current` is a programmatic
signal only. F3 (tint contrast) is the finding that would have touched this file
and it is out of scope.

Nothing under `do-work/` was staged or committed. This hand-back file is written
but deliberately left untracked.

## What F1 actually was

**The mechanism.** Both listeners involved are registered on `document` and both
are bubble-phase, so they run in registration order. `generate.go`'s
`boardJavaScriptFragmentPaths` puts `web/board-activity.js` at position 7 and
`web/board-controls.js` at position 10, so the Activity selection listener is
registered first and runs first on every click in the page. The
`[data-detail-kind]` delegation that actually opens the drawer
(`web/board-controls.js:241-247`) runs afterwards. Until it runs,
`currentDetailKind` / `currentDetailId` still name the REQ that was open **before**
the click.

The shipped code read exactly those two variables in its fallback branch:

```js
function syncActivitySelectionToClick(clickEvent) {
  var clickedRow = clickEvent.target.closest("[data-activity-request]");
  applyActivitySelectionHighlight(
    clickedRow ? clickedRow.getAttribute("data-activity-request") : selectedActivityRequestId()
  );
}
```

The row branch was correct because it never reads the drawer. The fallback branch
was the defect. Nothing recomputes the highlight afterwards: `renderActivity()`
has three callers (`board-activity.js:31` window buttons, `board-filters.js:190`
filter change, `board-controls.js:74` guarded by `renderedOnce.activity`), so a
wrong highlight survives until the reader changes a filter or the window.

**Reproduction, before the fix.** `TestJavaScriptBehaviorActivitySelectionFollowsADrawerOpenedOutsideTheTable`
replays one click in the real registration order over the shipped function bodies:
it calls `syncActivitySelectionToClick` with the click target first, and only then
moves the drawer state the way `openDetail` moves it. Sequence: open REQ-570 from
its own Activity row (drawer and table agree), then click a Board card for
REQ-505.

- Before: the table selects `[true, false, true]` — both REQ-570 rows — while the
  drawer shows REQ-505. Two surfaces naming two different REQs.
- After: the table selects `[false, true, false]` — the REQ-505 row — matching the
  drawer.

The same wrong result reproduces for a ticket link inside the drawer body and for
a Durations or Timeline mark, because all three are the same kind of trigger.

**The fix.** `syncActivitySelectionToClick` now resolves the selection from the
click event and never from a drawer that has not moved yet:

1. Click inside an Activity row → that row's `data-activity-request`. Unchanged.
2. Otherwise, the click's own `[data-detail-kind]` trigger — the same element the
   delegation is about to read. This mirrors `openDetail`'s two rules exactly: a
   `disabled` trigger opens nothing (`board-cards.js:120` renders one for a
   dependency target outside the tree, and `board-controls.js:243` checks it), and
   only the `"ur"` kind opens a drawer that is not a REQ, so a UR click clears the
   selection instead of marking it. Every other kind token — `"req"` from Board
   cards and Activity rows, `"request"` from the Timeline — is a REQ.
3. No trigger on the click at all → read the drawer, as before. This is the path
   the close button takes, and it is correct there: `#detail-close` has its own
   element-level listener, so `closeDrawer` has already run by the time this
   document-level listener sees the click.

No timing trick, no deferred callback, no dependence on fragment order in either
direction.

**F4.** `applyActivitySelectionHighlight` now sets `aria-current="true"` on a
selected row and removes the attribute from an unselected one, alongside the
existing `is-activity-selected` class. `aria-current` rather than `aria-selected`
because `aria-selected` on a row only carries meaning inside a `grid` or
`treegrid`; the Activity table is a plain table with no managed focus, and giving
it a grid role to justify `aria-selected` would change keyboard navigation far
beyond this finding. The attribute is removed rather than set to `"false"`,
because a present `aria-current` is announced whatever its value.

## Test evidence

Command (both runs):

```
QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 \
  bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban \
  -run '^TestJavaScriptBehavior' ./...
```

**RED** — new test in place, `web/board-activity.js` still as shipped:

```
=== RUN   TestJavaScriptBehaviorActivitySelectionFollowsADrawerOpenedOutsideTheTable
    javascript_behavior_c_test.go:3101: a REQ-505 drawer opened from a Board card left the table selecting []bool{true, false, true}, want []bool{false, true, false} — the drawer named REQ-505 while the table marked the REQ open before it
--- FAIL: TestJavaScriptBehaviorActivitySelectionFollowsADrawerOpenedOutsideTheTable (1.89s)
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	2.278s
go-test budget: FAIL module=skills/do-work-board/tools/queue-kanban wall=4s exit=1
```

**GREEN** — after the fix, the four Activity behavior tests:

```
--- PASS: TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests (2.03s)
--- PASS: TestJavaScriptBehaviorActivityViewHidesTheVerifyFindingsStrip (0.05s)
--- PASS: TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest (0.06s)
--- PASS: TestJavaScriptBehaviorActivitySelectionFollowsADrawerOpenedOutsideTheTable (0.05s)
ok  	github.com/knews2019/skill-do-work/queue-kanban	2.439s
```

Whole behavior lane, 60 tests, no failures:

```
go-test budget: module=skills/do-work-board/tools/queue-kanban wall=10s tests=60 slowest-file=javascript_behavior_c_test.go:2.94s limit=<30s
```

Whole package (`./...`, no `-run` filter), exit code 0:

```
go-test budget: module=skills/do-work-board/tools/queue-kanban wall=53s tests=454 slowest-file=generate_test.go:11.08s limit=<30s
exit=0
```

`gofmt -l` and `go vet ./...` both clean.

**Intermediate failure, resolved.** The first GREEN attempt failed two existing
Activity probes with `TypeError: tableRow.removeAttribute is not a function` —
their stub nodes had `setAttribute` but not `removeAttribute`, which
`applyActivitySelectionHighlight` now calls on every unselected row. Both stubs
gained the method. See Cross-REQ test edit below.

**Guard against the stubbed-`getElementById` trap.** The Node lane's
`document.getElementById` invents a stub node for any id, so an empty table body
would satisfy every `[]bool` and `[]string` comparison in this test without
proving anything. The new test asserts the three rendered rows'
`data-activity-request` values first, and fails there if the table did not render,
so no later assertion can pass vacuously.

**Cross-REQ test edit (per `crew-members/general.md`).**
`TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests` belongs to
REQ-572 (rendering the Activity table's summary counts). Its stub node gained
`removeAttribute`, with a comment in the file naming this remediation as the
cause. The change is intentional and mechanical: no assertion, expectation or
sliced function block in that test changed, and it still runs against the shipped
`renderActivity`. The same one-line addition went into REQ-573's own
`TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest`, likewise
with no assertion touched.

## Decisions

**D-R1 — F1 fix: resolve from the click's own `[data-detail-kind]` trigger.**
Chosen. It removes the ordering dependency instead of accommodating it: the
listener no longer reads any state that a later listener is about to change, so it
is correct under either fragment order and stays correct if the manifest is ever
reshuffled for an unrelated reason. It is synchronous, so the highlight and the
drawer never disagree even for one frame, and it is directly testable in the
existing Node probe lane with no timer or microtask flushing. Cost: it restates
`openDetail`'s two routing rules (the `"ur"` kind, the `disabled` check) in a
second place, about four lines. Both are pinned by the new test, so a change to
either rule in `board-detail.js` or `board-controls.js` that is not mirrored here
turns the test red.

**D-R2 — rejected: defer the drawer read by a task or a microtask.** This was the
review's first suggestion. It works in a browser, but it makes correctness depend
on event-loop timing rather than on data that is already in hand, and it leaves a
window — however short — in which the table and the drawer disagree. It is also
close to untestable in this repo's Node probe lane, which replays function bodies
synchronously and has no timer stub; proving the fix would have meant building
timer machinery in the test harness to verify a fix that did not need to be
asynchronous.

**D-R3 — rejected: have `openDetail` call `syncActivitySelectionToDrawer`.** This
is the cleanest fix on paper — one authority, the drawer, tells the table to
re-read after it changes. Two reasons against it here. It is outside the permitted
write set (`web/board-detail.js`), and it points the wrong way in the dependency
graph: `board-detail.js` is a general drawer that knows nothing about the Activity
view today, and making it call into a view-specific function adds a coupling that
would have to be repeated for the next view that wants the same behavior. Worth
reconsidering as a deliberate refactor if a second view needs drawer-open
notification; not worth introducing for one caller under a remediation.

**D-R4 — rejected: pin the fragment manifest order.** Explicitly excluded by the
dispatch brief, and the review's F5 says the same for a good reason: the hand-back's
own Discovered Task proposed a test locking `board-activity.js` ahead of
`board-controls.js`, which would have frozen the order that causes the bug. The
manifest is untouched by this commit.

**D-R5 — F4 signal: `aria-current="true"` on the `<tr>`.** Reasoning in "What F1
actually was" above. The client's existing selected-state convention is
`aria-pressed` on toggle buttons in a group (`board-controls.js:15` and `:97`,
`web/template.html`), which is the wrong ARIA for a table row: the Activity REQ
button is not a toggle, and clicking it a second time does not close the drawer.

**D-R6 — kept `is-activity-selected` alongside `aria-current`.** Styling off
`tr[aria-current="true"]` and deleting the class would have left one signal
instead of two, which is the cheaper shape. Rejected as churn: the class is
REQ-573's just-merged contract, three CSS rules and two tests read it, and the
saving is cosmetic. `crew-members/coding-guardrails.md` §3 governs.

## Integration seams

- **`board-activity.js` now depends on `openDetail`'s routing rules.** Two
  specific facts: `openDetail` (`web/board-detail.js:676`) sends only the `"ur"`
  kind to `openUserRequestDetail` and everything else to `openRequestDetail`, and
  the delegation (`web/board-controls.js:243`) skips a `disabled` trigger. If a
  third detail kind is ever added, or the disabled check moves, this listener has
  to follow. The new test fails if they drift apart.
- **No change to the fragment manifest**, so nothing about load order moved.
  `board-activity.js` may still be reordered freely with respect to
  `board-controls.js` — that is now the point.
- **Every `[data-detail-kind]` producer is a client of this behavior**, not just
  Board cards: `web/board-cards.js:118` (`dataset.detailKind`),
  `web/board-durations.js:912`, `web/board-timeline.js:1568` and `:2177` (kind
  `"request"`), `web/board-detail.js:21` (ticket links inside the drawer body),
  and the Activity REQ button itself. Any new producer of that attribute pair now
  also moves the Activity highlight, which is the intended behavior.
- **`applyActivitySelectionHighlight` now calls `removeAttribute`** on table rows.
  Any future probe that renders the Activity table with a stub node must provide
  it; the two stubs in `javascript_behavior_c_test.go` that needed it have it.
- **Untouched by this commit:** `web/board.css`, `generate.go`, `web/template.html`,
  `web/board-controls.js`, `web/board-detail.js`, and findings F2, F3, F5, F6, F7
  and F8, which remain open in the review report.

## Discovered Tasks

- The dispatch line `git merge origin/main 2>/dev/null || git merge main` silently
  no-ops when `origin/main` is a stale ancestor, leaving a worktree builder on the
  pre-integration tree with a success exit code. The `||` fallback can never fire
  in that case. A run that dispatches worktree builders should merge the local
  integration branch by name, or verify the merge landed. — impact-user-visible
  → report only
- `web/board-activity.js` builds the REQ cell button with the bare id as its
  accessible name, while `makeRequestCard` (`web/board-cards.js:123`) sets
  `aria-label` to `"REQ-570: title"`. This is review finding F8 and is left as
  reported. — impact-negligible → report only
