# Hand-back — REQ-573 (open the detail drawer from an Activity row and highlight every row of the same REQ)

## Branch

- Branch: `worktree-agent-REQ-573-activity-drawer`
- Head commit: `05484aa6` — `[REQ-573] open the detail drawer from an Activity row and select the REQ's whole set of rows`
- Base commit: `5f4821ab`
- One commit on the branch. The hand-back file is not staged and not committed.

## File manifest

| Verb | Path |
|---|---|
| modified | `skills/do-work-board/tools/queue-kanban/web/board-activity.js` |
| modified | `skills/do-work-board/tools/queue-kanban/web/board.css` |
| modified (tests) | `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` |

`git diff --stat` against the base: 3 files changed, 350 insertions, 2 deletions.

Nothing outside the declared write set was touched. No file under `do-work/` was written except this hand-back.

Tests touched:

- **Added** `TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest` in `javascript_behavior_c_test.go` — the REQ's captured RED/GREEN pair.
- **Updated** `TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests` (REQ-572's test, one Activity row per lifecycle stamp) — see D-05.

## P-A-U

### [PLAN]

Read before code: `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`, `_dev/primes/prime-kanban-board.md`, and both lesson satellites named in the brief. Then read `web/board-activity.js`, `web/board-controls.js`, `web/board-detail.js` (`openDetail`/`closeDrawer`/`showDrawer`), `web/board-cards.js` (`makeRequestCard`, the keyboard pattern to copy), `generate.go`'s fragment execution-order manifest, and the activity block of `web/board.css`.

Approach:

1. **Opening the drawer.** In `renderActivity`, build the REQ cell's content as a `<button type="button" class="activity-req-button">` carrying `data-detail-kind="req"` and `data-detail-id="<REQ id>"`, mirroring `makeRequestCard`. The existing document-level `[data-detail-kind]` delegation in `board-controls.js` then opens the drawer with no new handler, and a real button is Tab-reachable where a `td` with a click listener is not. The `tr` keeps its existing `data-activity-request`.
2. **Selection state.** Do not introduce a second copy of "which REQ is selected". All web fragments share one function scope (`board-clipboard.js` already reads `currentDetailId`), so the highlight reads the drawer's own `currentDetailKind` / `currentDetailId`. That single read makes three requirements fall out with no extra state: re-render restores the highlight, closing the drawer clears it, and a UR drawer never selects REQ rows.
3. **Applying it.** `applyActivitySelectionHighlight(selectedRequestId)` toggles `is-activity-selected` on every row whose `data-activity-request` matches, so one REQ's whole set of rows lights up. Called at the end of `renderActivity` from the drawer state, and from a document click listener.
4. **The click and close paths.** `generate.go` runs `board-activity.js` before `board-controls.js`, so this fragment's document listener is registered first and runs first — the drawer still names the *previous* REQ while it runs. So the click path reads the clicked row's own id; every other click re-reads the drawer, which is what makes the drawer's close button clear the highlight (its element-level listener has already run `closeDrawer`). Escape is the drawer's only other close path (`closeDrawer` has exactly two callers), so a bubble-phase `keydown` listener covers it — `board-detail.js` dismisses from a capture-phase listener, so the state is already cleared by then.
5. **Styling.** Three signals so the highlight is not colour alone: `--tint-claimed` background, an inset 3px `--accent-claimed` bar down the REQ column, and the bold underlined id. Both tokens are redefined in the light block, so both themes are covered. The bar is an inset shadow rather than a border because under `border-collapse: collapse` a real border would shift the row.
6. **Verify** with the two focused lanes, plus a mutation check to prove the new assertions bite.

### [APPLY]

Coded as planned, inside the three declared files. No other file was created, edited, or deleted. Two deviations from the first draft, both caught during UNIFY and both fixed before commit:

- The selected-row rule first carried `padding-left: 9px`, which would have shoved the row's text 9px sideways at the moment of selection. The inset is now permanent on the REQ column (`thead` header included, so the column still lines up) and only the bar appears on selection.
- REQ-572's probe stub had to be widened; see D-05.

### [UNIFY]

`git diff --stat` (against `5f4821ab`):

```
 .../queue-kanban/javascript_behavior_c_test.go     | 233 +++++++++++++++++++++
 .../tools/queue-kanban/web/board-activity.js       |  69 +++++-
 .../do-work-board/tools/queue-kanban/web/board.css |  50 +++++
 3 files changed, 350 insertions(+), 2 deletions(-)
```

Linters, each run from the worktree:

| Command | Result |
|---|---|
| `gofmt -l .` (in `skills/do-work-board/tools/queue-kanban`) | no output — clean |
| `go vet ./...` | exit 0, no findings |
| `node --check web/board-activity.js` | exit 0 — parses |

Debug-artifact sweep on added lines: `git diff -U0 \| grep -E "^\+.*(console\.log\|debugger\|TODO\|FIXME\|XXX)"` returned nothing.

Per-file review:

- **`web/board-activity.js`** — checked that the REQ cell branch is an `else if` chain that leaves the instant and plain-text branches untouched; that `createElement`, `currentDetailKind` and `currentDetailId` all resolve in the shared fragment scope; that `applyActivitySelectionHighlight` guards a missing table body the same way `renderActivity` already guards its three nodes; that nothing here calls `openDetail` or registers a second drawer opener; that no stamp meaning is re-derived and no client-side re-sorting was added (the REQ's Constraints); that the two new top-level `document.addEventListener` calls are the only statements outside a function, matching how `board-controls.js` registers its own delegation. Names are two-word and greppable (`applyActivitySelectionHighlight`, `selectedActivityRequestId`, `syncActivitySelectionToClick`, `syncActivitySelectionToDrawer`, `detailButton`, `clickedRow`, `activity-req-button`, `is-activity-selected`).
- **`web/board.css`** — checked that every token used (`--tint-claimed`, `--accent-claimed`, `--ink-strong`) is defined in the base palette *and* redefined in the `prefers-color-scheme: light` block, so nothing gets its only definition in one theme; that the new rules are scoped under `.activity-table-scroll` / `.activity-req-button` and cannot reach any other view; that the focus-visible ring copies the `.req-card` / `.control-button` convention byte for byte; that no selection-only rule changes layout metrics.
- **`javascript_behavior_c_test.go`** — checked that the new probe stubs nothing that is under test (the selection helpers are sliced from the shipped bundle, not redefined), that `closest` supports only the `[attribute]` form the client actually passes and throws on anything else rather than silently pretending, that every new block sliced by `sliceBalancedBlockAfter` contains no brace inside a string literal (the helper's stated precondition), and that the assertions are on values, not truthiness.

## Test evidence

All commands run from the worktree root `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905-1201/worktree-agent-REQ-573-activity-drawer`. The repository gate (`_dev/tests/maintainer-verify.sh`) was **not** run.

| # | Command | Exit |
|---|---|---|
| T1 | `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest$' ./...` (before implementation) | 1 — RED |
| T2 | same as T1, mutated implementation | 1 — RED, assertion level |
| T3 | `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' ./...` | 0 — 58 tests, wall 7s, `limit=<30s` |
| T4 | `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...` | 0 — 392 tests, wall 44s, slowest file `generate_test.go` 9.42s, `limit=<30s` |

Every test-file invocation finished under the 30s per-file budget the brief sets; T4's 44s is the whole-module wall, and the budget script's own per-file limit reported `<30s`.

### RED

Test: `TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest`. Before any implementation:

```
=== RUN   TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest
    javascript_behavior_c_test.go:2678: anchor "function selectedActivityRequestId(" not found in the generated page
--- FAIL: TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest (1.81s)
```

That is the shipped client not defining the behavior, which is what the REQ's "Why RED now" states: `renderActivity` built plain `tr`/`td` with no `data-detail-kind` and no selection state.

This lane extracts named function blocks out of the generated bundle, so a missing behavior surfaces first as a missing block rather than as an assertion diff. `testing.md` asks for the assertion itself to fail, so I also ran a **mutation check (T2)** after GREEN: with the row match in `applyActivitySelectionHighlight` replaced by "the first row only", the test fails on the assertion, proving it pins behavior rather than the presence of a name.

```
    javascript_behavior_c_test.go:2872: clicking one REQ-570 row selected []bool{true, false, false},
    want []bool{true, false, true} — every row of that REQ, and no other REQ's row
--- FAIL: TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest (1.62s)
```

The mutation was reverted from a saved copy immediately after that run; the committed source is the unmutated implementation, and T3/T4 were re-run after the revert.

### GREEN

```
=== RUN   TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest
--- PASS: TestJavaScriptBehaviorActivityRowClickSelectsEveryRowOfTheSameRequest (0.05s)
```

with the whole Node behavior lane at exit 0 (T3) and the whole Go module at exit 0 (T4).

The new test covers the REQ's captured proof and its four Detailed Requirements:

1. Three rows render for two REQs (two for REQ-570, one for REQ-505); each REQ cell is a `button` with `type="button"`, its label is the REQ id, and it carries `data-detail-kind="req"` plus the matching `data-detail-id` — the entire contract with the `board-controls.js` delegation, which is what opens the drawer.
2. Clicking one REQ-570 row selects `{true, false, true}` — both REQ-570 rows, not the REQ-505 row.
3. A re-render with the drawer open restores the same set with nothing remembered in between.
4. Clicking the REQ-505 row moves the selection to `{false, true, false}`.
5. Clearing the drawer and clicking elsewhere (the close-button path) clears the highlight; so does `syncActivitySelectionToDrawer` after Escape.
6. A drawer opened from outside the table marks that REQ's rows.
7. A UR drawer whose id reads like a REQ id selects nothing — the kind is checked, not just the id.

### Untested

The two `document.addEventListener` registrations themselves are top-level statements in the fragment, and this lane can only slice named function blocks, so the probe drives `syncActivitySelectionToClick` and `syncActivitySelectionToDrawer` directly rather than through a dispatched DOM event. What is unproven by the Node lane is only the wiring (that a real click reaches those functions) and the rendered pixels of the highlight in the two themes. A browser probe would close both. I did not run one: the prime warns that a sibling agent sharing the browser can navigate it between a navigate and an evaluate, and several builders are running in this session. See **Discovered Tasks** and **D-04**.

## Lesson evidence

| Satellite | Read | Missing paths |
|---|---|---|
| `_dev/primes/prime-kanban-board.md` (the prime itself) | yes, in full | none |
| `_dev/primes/lessons-kanban-board.md` | yes, in full | none |
| `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` | yes, in full | none |

Both satellites were listed in the REQ under "Required Lessons — Dropped for Budget", so they were read here rather than at claim time. Every path either file names that this change touches resolved on disk. Entries that shaped the work:

- The prime's **"a chart's correctness is partly a claim about pixels"** and **"render evidence must name the page it measured, in the same call that measures it"** are why the visual half is reported as Untested with a named reason instead of being asserted from a passing suite (D-04).
- The prime's **"the surface behind this board's SVG is `<body>`"** does not apply — this is a table with real painted backgrounds, not a transparent SVG, so the tint is measured against a surface the reader actually sees.
- `lessons-do-kanban.md` **REQ-117** ("a stated reason in a comment is a factual claim and reviews must check it like any other") is why the fragment-ordering claim in the new comment was checked against `generate.go`'s manifest rather than assumed.
- `lessons-do-kanban.md` **REQ-116** ("pick the altitude of a test, not just its count") is why the assertions run through the real `renderActivity` and the real selection helpers sliced out of the generated bundle, not against hand-built structures.
- `lessons-kanban-board.md` was read in full; nothing in it contradicts this approach.

Neither satellite records the fragment-ordering trap this REQ ran into. Both files are outside my write set, so I did not add an entry — see **Integration seams**.

## Decisions

**D-01 — The drawer is the only selection state. (DECIDE & STATE)**
`selectedActivityRequestId()` returns `currentDetailKind === "req" ? currentDetailId : ""`, read live from `board-detail.js`. All fragments share one function scope (`board-clipboard.js` already reads `currentDetailId`), so no new plumbing was needed. Three of the REQ's requirements then need no code of their own: a re-render restores the highlight because it re-reads, closing the drawer clears it because `closeDrawer` clears both variables, and a UR drawer cannot select REQ rows because the kind is part of the read. A module-level `selectedActivityRequestId` variable would have been a second copy that can disagree with the drawer. Reversible: it is one small function.

**D-02 — The click path reads the clicked row, not the drawer. (DECIDE & STATE)**
`generate.go`'s execution-order manifest runs `board-activity.js` (position 7) before `board-controls.js` (position 10), so this fragment's document listener is registered first and therefore fires first in the bubble phase — at that moment the drawer still names the previous REQ. Reading the clicked row's own `data-activity-request` makes the click path independent of listener order. Every non-row click falls back to the drawer read, which is exactly what makes the drawer's close button clear the highlight: `#detail-close` has an element-level listener, so `closeDrawer` has already run. Verified by reading the manifest, not assumed.

**D-03 — A `keydown` listener for Escape rather than a `MutationObserver` on the drawer. (DECIDE & STATE)**
`closeDrawer` has exactly two callers: the `#detail-close` click (covered by D-02's fallback) and the Escape handler in `board-detail.js`. A bubble-phase `keydown` listener covers the second, and `board-detail.js` dismisses from a *capture*-phase listener, so the state is already cleared when this one runs. A `MutationObserver` on the drawer's `hidden` attribute would cover every close path including future ones, but it cannot see a row-to-row switch while the drawer stays open, so it would not have replaced the click listener — it would have been a third mechanism for a case that does not exist yet. If a third close path is ever added, this listener list is where it has to be named.

**D-04 — No browser render evidence, reported as Untested rather than claimed. (ESCALATE)**
The prime is explicit that a passing suite is not evidence about pixels, and this REQ has a visual requirement ("readable in light and dark themes"). I did not drive a browser.
*Value:* a render check would confirm the tint, the inset bar and the focus ring against both palettes, and confirm the click wiring end to end — the two things the Node lane cannot see.
*Risk:* low but not zero. Every token used is defined in both palettes and every rule copies an existing pattern in the same stylesheet, so a broken theme is unlikely; what a render would most plausibly catch is a weak tint on the light surface, or `box-shadow` on a `th` under `border-collapse: collapse` rendering differently than expected. Fully reversible — it is three CSS rules. I skipped it because the prime also warns that a shared browser instance can be navigated by a sibling agent between the navigate and the measurement, several builders are running in this session, and the brief limited me to two focused test commands. If the orchestrator wants the check before merge, it is one board generation and a screenshot in both colour schemes.

**D-05 — Updated REQ-572's activity-summary test instead of leaving it broken. (DECIDE & STATE)**
`TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests` (REQ-572, one Activity row per lifecycle stamp) failed with `ReferenceError: createElement is not defined` once `renderActivity` started building a button. Per `general.md`'s Cross-REQ Test-Break rules the change is intentional, so the test was updated to match: three function blocks added to its slice list (`createElement`, `selectedActivityRequestId`, `applyActivitySelectionHighlight`), `classList.toggle` and `getAttribute` added to its stub node, and `currentDetailKind` / `currentDetailId` declared empty. Not one assertion was changed, weakened, or removed — the transition and REQ counts, the repeated-id contract and both empty states are pinned exactly as REQ-572 left them, and the comment in the file says which test now owns selection.

**D-06 — No `disabled` branch for a REQ the board never parsed. (DECIDE & STATE)**
`makeRequestCard` disables the card when `requestsById` has no entry, because it renders dependency targets that can point outside the tree. Activity rows come from stamps on tickets the board itself parsed, so that case does not arise here, and `renderActivity`'s existing `requestsById[row.id] || {}` already degrades to a blank title and status. Adding a disabled path would be defensive surface with no incident behind it.

**D-07 — The highlight carries three signals. (DECIDE & STATE)**
Background tint, an inset bar down the REQ column, and the bold underlined id. The REQ asks for "a left border or bold id beside the background"; both were cheap, and the bar gives the eye a left edge to track down a long table while the underline says which cell the selection is about. The bar is `box-shadow: inset` rather than `border-left` because under `border-collapse: collapse` a real border would shift the row, and the column's 9px inset is permanent (header included) so a selection never moves the text.

## Discovered Tasks

- The two `document.addEventListener` registrations in `web/board-activity.js` are top-level statements, which the Node behavior lane cannot slice, so the wiring from a real DOM click to `syncActivitySelectionToClick` is asserted only indirectly (through the `data-detail-*` attributes and by driving the functions directly). A browser probe in the strict lane, of the shape `timeline_browser_probe_test.go` already uses, would close that gap for the Activity view and would also measure the highlight's contrast in both themes. `impact-negligible` → report only.
- `web/board-controls.js` and `web/board-activity.js` now each register a document-level `click` listener, and their relative order is decided by `generate.go`'s fragment manifest rather than by anything either file states. Nothing pins that order today; reordering the manifest would silently change which id the highlight reads on a click. A single test asserting `board-activity.js` precedes `board-controls.js` in `boardJavaScriptFragmentPaths` would make the dependency visible. `impact-negligible` → report only.
- `lessons-do-kanban.md` records no entry about fragment execution order deciding document-listener order, though the manifest comment in `generate.go` calls itself "the sole execution-order manifest". The next person to add a listener in an early fragment will re-derive it. See **Integration seams** for the exact line. `impact-negligible` → report only.

## Integration seams

One line, for a file outside my write set. Optional — it is a lessons entry, not a code change.

**File:** `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`
**Where:** appended as the last bullet of the list.
**Line:**

```
- [family: fragment-order-decides-listener-order] REQ-573: `generate.go`'s `boardJavaScriptFragmentPaths` decides more than load order — it decides the firing order of every document-level listener, because each fragment registers its own at execution time. `board-activity.js` sits three entries before `board-controls.js`, so on a click its listener runs BEFORE the `[data-detail-kind]` delegation has called `openDetail`, and the drawer's `currentDetailId` still names the PREVIOUS REQ. A later fragment reading drawer state on a click is correct; an earlier one must read the clicked element instead. Nothing pins that order, so a manifest reshuffle would change the behavior silently and no assertion would notice.
```

If this line is taken, `do-work/lessons-index.md` needs its `lessons-do-kanban.md` row refreshed in the same edit (slug set gains `fragment-order-decides-listener-order`, size estimate up by roughly 130 tokens, coverage still `partial`). Both files are the orchestrator's to write. This is a first occurrence of the family, and it constrains any future change that adds a listener to this client, so `general.md`'s rule points at promoting it into the prime's `## Traps` — I leave that call to the orchestrator, since `_dev/primes/prime-kanban-board.md` is also outside my write set.

No other line belongs in a file outside the write set. The change needs no version bump, changelog entry, or template/Go-side edit: `template.html` already carries the table and the column header ids, the new attributes and classes are produced entirely by the client, and `generate.go` already lists `web/board-activity.js` and `web/board.css`.
