# REQ-572 Exploration — one Activity row per lifecycle stamp

All paths relative to `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`.

## 1. `activity.go` and `activity_test.go`

`skills/do-work-board/tools/queue-kanban/activity.go` (105 lines, whole file is in scope):

| Symbol | Line | What it does |
|---|---|---|
| header comment | 8-31 | States "One row per REQ: its newest lifecycle stamp"; states the window is NOT applied here; points at `isWithinRecentWindow` (model.go) |
| `ActivityRow` struct | 34-47 | `RequestId`, `StampField`, `StampTime`, `Transition` |
| `buildActivityRows` | 56-78 | Loops tickets, calls `newestLifecycleStamp`, appends one row, sorts newest-first with `RequestId >` tie-break in the comparator (72-75) |
| `newestLifecycleStamp` | 85-105 | Reads every `lifecycleTimestampFields(ticket)` entry, keeps the max parseable one |

Change shape: `buildActivityRows` loops `lifecycleTimestampFields(ticket)` directly and appends one row per parseable stamp. `newestLifecycleStamp` becomes dead unless kept for another caller — grep shows it has exactly one caller, `buildActivityRows:62`. Capacity hint at line 57 (`len(tickets)`) becomes an undercount; harmless but worth bumping.

**Tie-break defect the change introduces.** The comparator only breaks ties on `RequestId`. Once one REQ can produce several rows, two stamps of the SAME REQ at the SAME instant (very common: `completed_at` and `status_changed_at` are written together by the work loop, and `created_at`/`claimed_at` can coincide) make the comparator return false both ways, so `sort.Slice` (unstable) decides. Add `StampField` as the third key. `TestBuildActivityRowsBreaksStampTiesDeterministically` (activity_test.go:163) will not catch it — it uses two different REQs.

### Existing tests (`activity_test.go`, 251 lines)

| Test | Line | Pins | Fate |
|---|---|---|---|
| `activityWindow` const | 13 | 24h window shared by the boundary tests | keep |
| `TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition` | 19 | REQ-568's captured GREEN: 5 tickets, exact order `504,505,567,503,485`, `len(inWindow)==5`, transitions per REQ | **REWRITE.** Four of the five tickets carry two or three stamps, so the row count and the whole `wantOrder` slice change. Keep the fixture instants; assert the newest-per-REQ rows still sit in the same relative order and that the older stamps now appear too |
| `TestBuildActivityRowsPicksTheNewestStampNotTheFirstDeclaredField` | 92 | `len(rows)==1`, `rows[0].StampField=="review_at"` | **REWRITE.** The ticket has four parseable stamps, so `rows` becomes 4. Rename to something like "orders one ticket's stamps newest first" and assert the sequence `review_at, planning_at, claimed_at, created_at` — that keeps the same anti-first-match intent |
| `TestBuildActivityRowsStraddlesTheWindowBoundary` | 119 | Two single-stamp tickets, `len(rows)==2`, windowing left to caller | **Survives as written** (each ticket has exactly one stamp). Verify the index assumption at 135/138 still holds after the sort change |
| `TestBuildActivityRowsSkipsTicketsWithNoParseableStamp` | 146 | Stampless and unparseable tickets skipped; `created_at` transition is "captured" | **Survives**; each fixture ticket has at most one stamp |
| `TestBuildActivityRowsBreaksStampTiesDeterministically` | 163 | Comparator, not sort stability, breaks equal-instant ties | **Survives**, but should be extended with the same-REQ/same-instant case described above |
| `TestLifecycleTimestampFieldsIsTheOneListBothReadersUse` | 186 | Anti-drift: `detectFutureTimestampFields` and the aggregation read the same list; every declared field has a non-empty transition; the `everyStamp` literal (199-215) is a deliberate second spelling of the 14-field list | **REWRITE the tail.** Lines 237-250 assert `len(rows)==1` for `everyStamp`. Under the new rule that ticket yields 14 rows, which is a *better* pin: assert `len(rows)==len(declared)` and that the set of `StampField` values equals the declared set |

Ticket fields available on `RequestTicket` for fixtures are the 14 named at activity_test.go:199-215.

## 2. `generate.go`

- `boardJavaScriptFragmentPaths` includes `web/board-activity.js` at generate.go:50 (manifest, locked by `TestBoardJavaScriptAssemblyStructure`; also asserted in generate_test.go:45 and :72). No change needed.
- `Activity []generatedActivityEntry` field + comment: generate.go:82-86. The comment says "every REQ's newest lifecycle stamp" — restate.
- `generatedActivityEntry` type + comment: generate.go:308-318. JSON tags `id`, `stampField`, `stampAt`, `transition`. **Shape stays exactly as-is** (REQ-572 requires it, REQ-573 depends on it).
- Fill site: generate.go:860-867, a plain `for _, row := range buildActivityRows(board.AllRequests)` append. No change needed beyond the comment.

**Who else reads `boardData.activity`:** nobody. Only `web/board-activity.js:21`. `web/board-controls.js:62-65` and `web/board-filters.js:181,191` reference `renderedOnce.activity` (the render cache flag), not the data.

No Go test asserts on the Activity payload contents; `grep stampField|stampAt` over `*_test.go` hits only activity_test.go comments.

## 3. `web/board-activity.js` (112 lines, whole file in scope)

- Header comment 1-13: restates "One row per REQ: its newest lifecycle stamp".
- `activityRowsWithin(windowHours)` 15-25: `Date.now()`-anchored cutoff, filters `boardData.activity` on `Date.parse(row.stampAt)`. **Unchanged** — it already filters rows, not REQs.
- `applyActivityWindowSelection` 27-31: sets `viewState.activityWindowHours`, calls `setActiveButton("#activity-window-group", "data-activity-window", …)`, re-renders.
- `activityWindowPhrase` 33-41: "last 6 hours" / "last 24 hours" / "last 7 days".
- `renderActivity` 43-112: reads `#activity-summary`, `#activity-table-body`, `#activity-empty`; windows, then filters with `requestMatchesFilters(row.id)` (56-59); builds the summary (61-65); clears and fills the tbody (67-98); sets the empty state (100-111).

**Summary text today (line 61-64):**
```
rows.length + (rows.length === 1 ? " REQ" : " REQs") + " touched in the " + activityWindowPhrase(windowHours)
+ optional " (" + windowRows.length + " before filters)"
```
Target per the REQ: `"38 transitions across 21 REQs in the last 24 hours"`. Distinct-REQ count needs a small id set built over `rows` — the client already has `requestsById`, so a plain object used as a set matches house style (see `recentlyDoneIdSet` usage in `web/board-filters.js:88`). Keep the "(N before filters)" clause; decide whether the before-filters number is transitions (simplest, matches the current meaning) and say so in a comment.

**Empty state (108-111):** `"No REQ carries a lifecycle stamp inside the …"` and `"N REQ(s) moved in this window, but the active filters hide all of them."` The second is now a transition count, so reword.

**Row build (68-98):** cells are `[id, title, status, transition, when, stampField]`; the id cell is a `<th scope="row">` with `headers=activity-table-column-req`. Each `<tr>` gets `data-activity-request="<id>"` (line 71) — that attribute stops being unique per row once a REQ has several. REQ-573 will use it for sibling highlighting, so leave it on every row (that is what makes sibling selection possible) rather than making it unique.

**Helper definitions:**
- `requestMatchesFilters(requestId, options)` — `web/board-filters.js:65-82`.
- `makeInstantWithRelativeNode(isoText)` — `web/board-core.js:81-90` (returns a `span.instant-with-relative`, or null when unparseable).
- `setActiveButton(groupSelector, attributeName, value)` — `web/board-controls.js:10`.

## 4. `web/template.html`

- View button: line 83-85, `data-view-target="activity"`.
- Activity window button group: lines 91-100, `id="activity-window-group"`, buttons `data-activity-window` = 6 / 24 (default `is-active`, `aria-pressed="true"`) / 48 / 168.
- Section markup: lines 460-486. Comment 460-467 restates the one-row-per-REQ rule. `<section id="view-activity" class="view-panel" aria-label="Recently touched REQs" hidden>`, then `<p class="activity-summary" id="activity-summary">`, `<p class="activity-empty" id="activity-empty" hidden>`, `<div class="activity-table-scroll">` wrapping a table whose `<th scope="col">` ids are `activity-table-column-req|title|status|transition|when|stamp` (lines 475-480) and `<tbody id="activity-table-body">` (483).

Only the comment needs to change unless a per-REQ visual grouping cue is added.

## 5. `web/board.css`

Activity block at `web/board.css:3151-3220`:
- `#view-activity` 3158
- `.activity-summary` 3162
- `.activity-empty` 3168
- `.activity-table-scroll` 3175 (`max-height:70vh; overflow:auto`)
- `.activity-table-scroll table` 3182
- `.activity-table-scroll th, .activity-table-scroll td` 3187 (`white-space:nowrap`, `border-bottom: 1px solid var(--line-soft)`)
- `.activity-table-scroll thead th` 3196 (sticky)
- `.activity-table-scroll tbody th` 3205 (mono)
- `td[headers="activity-table-column-title"]` 3211 (the one wrapping column)
- `td[headers="activity-table-column-stamp"]` 3216 (mono, faint)

Reusable "active" styling elsewhere: `.control-button.is-active` (board.css:353) and `.view-panel.is-active` (board.css:654). There is **no** existing selected-row or highlighted-row rule anywhere in the file (`grep is-selected|\.selected|is-highlighted` returns nothing), so REQ-573 will have to invent one. If REQ-572 wants a repeat-REQ visual cue now, a muted id cell on repeated rows is the cheapest option, but the REQ does not ask for it.

## 6. The Node behavior test lane

Helpers live in `generate_test.go`:
- `javaScriptBehaviorProbeMode = "QUEUE_KANBAN_JAVASCRIPT_PROBES"` (generate_test.go:26); `lookupNodeForJavaScriptProbe` (264) skips unless it is `on` and `node` is on PATH.
- `runJavaScriptBehaviorProbe(t, name, probeSource)` (generate_test.go:276) pipes the probe to `node -` on **stdin** (never `-e`; the assembled client blows past Linux's 128 KiB arg limit).
- `generateLiveSite(t)` (generate_test.go:361, `sync.Once`-cached) returns the generated `index.html`.
- `sliceBalancedBlockAfter(t, html, "function renderActivity(")` (generate_test.go:1999) cuts one brace-balanced function out of the assembled client.

Tests do **not** use jsdom. They hand-stub `document`. Existing example, `TestJavaScriptBehaviorTestingDoneWindowIsViewSpecific` (`javascript_behavior_a_test.go:420-450`):

```go
functionBlocks := []string{
    sliceBalancedBlockAfter(t, indexHtml, "function createElement("),
    sliceBalancedBlockAfter(t, indexHtml, "function fillColumn("),
}
javascriptProbe := `
var filterState = { searchText: "", domain: "", status: "", doneWindow: "168" };
var viewState = { view: "board" };
var nodesBySelector = {};
function makeNode() {
  return { childNodes: [], textContent: "",
    appendChild: function (childNode) { this.childNodes.push(childNode); return childNode; } };
}
var document = {
  createElement: function () { return makeNode(); },
  querySelector: function (selector) { … }
};
` + strings.Join(functionBlocks, "\n") + `
fillColumn("board", [], null, 1);
process.stdout.write(JSON.stringify([...]));`
probeOutput := runJavaScriptBehaviorProbe(t, "testing empty-copy decision", javascriptProbe)
```

For a summary-count test, slice `activityRowsWithin`, `activityWindowPhrase` and `renderActivity`; stub `document.getElementById` (renderActivity uses `getElementById`, not `querySelector`), plus globals `boardData`, `viewState`, `requestsById`, `requestMatchesFilters` (stub returning true, or slice the real one from board-filters.js plus `filterState`), `makeInstantWithRelativeNode` (stub returning null so the cell falls back to text), and `createElement`. Emit `summaryNode.textContent` and the row count. Note there is currently **no** JavaScript behavior test for the Activity view at all (`grep renderActivity|activity-summary javascript_behavior_*_test.go` is empty), so this is a new test, and the name must start with `TestJavaScriptBehavior` to land in the lane's `-run` pattern.

Also note `TestMain` refuses a green run whose probes all skipped when `QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1`.

## 7. Test commands

Module dir: `skills/do-work-board/tools/queue-kanban` (own `go.mod`; a repo-root `go build ./...` does not reach it).

Fast Go tests, budget wrapper (fails any source test file at or above 30 s, `DO_WORK_TEST_FILE_BUDGET_SECONDS`):
```
bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...
```
The gate sets `QUEUE_KANBAN_JAVASCRIPT_PROBES=off QUEUE_KANBAN_BROWSER_PROBES=off DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior` around it (`_dev/tests/maintainer-verify.sh:605-612`).

Node behavior lane:
```
QUEUE_KANBAN_BROWSER_PROBES=off QUEUE_KANBAN_JAVASCRIPT_PROBES=on \
QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 \
bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' -v .
```
(`_dev/tests/maintainer-verify.sh:709-722`; lane id `queue-kanban-javascript` in `_dev/tests/heavy-lanes.json:5-15`.) Browser lane is `queue-kanban-browser` with `QUEUE_KANBAN_BROWSER_PROBES=on`.

Plain iteration during work: `cd skills/do-work-board/tools/queue-kanban && go test -run TestBuildActivityRows ./...` plus `go vet ./...`.

Prime test map: there is no test-map file for queue-kanban. `_dev/primes/prime-kanban-board.md` is the maintainer prime and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` is the shipped routing index (§ Read first, § Traps, § Stakes, § Lessons). Neither lists tests.

## 8. Restatements of "one row per REQ" to update

| File:line | Text |
|---|---|
| `skills/do-work-board/tools/queue-kanban/activity.go:8-13` | "One row per REQ: its newest lifecycle stamp, the field that stamp lives in, and the transition it records." |
| `activity.go:33` | "`ActivityRow` is one REQ's most recent lifecycle transition." |
| `activity.go:49-55` | "`buildActivityRows` returns one row per ticket carrying at least one parseable lifecycle stamp" |
| `activity.go:80-84` | `newestLifecycleStamp` doc, if the function goes away |
| `generate.go:82-86` | "every REQ's newest lifecycle stamp and the transition it records, newest first" |
| `web/board-activity.js:2-4` | "One row per REQ: its newest lifecycle stamp and the transition that stamp records, newest first." |
| `web/board-activity.js:53-55` | "Two counts, because they answer different questions" — still true, now three numbers |
| `web/template.html:460-467` | "every REQ whose newest lifecycle stamp falls inside the selected window" |
| `CHANGELOG.md:84` (0.276.0 — Board Activity View) | "lists every REQ whose newest lifecycle stamp falls inside the window". Historical entry describing what 0.276.0 shipped — do NOT rewrite it; add the new behavior as a fresh entry under the release this REQ lands in |

Not found anywhere: `skills/do-work-board/docs/board-guide.md` never documents the Activity view (only the calendar, line 3); `_dev/primes/lessons-kanban-board.md` and `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` contain zero "activity" hits; `prime-do-kanban.md` mentions only the "queue activity calendar". So the only user-facing doc surface is the changelog.

## Concerns

- **C1 Same-REQ tie-break (real bug if unaddressed).** Add `StampField` as the third sort key. `completed_at` and `status_changed_at` frequently share an instant on one ticket.
- **C2 Near-duplicate rows.** A terminal REQ often carries `completed_at` and `status_changed_at` at the same second, giving two rows saying nearly the same thing ("completed" and "status changed to completed"). The REQ does not ask for de-duplication and the screenshot shows the "status changed to …" phrasing is already visible on the board, so ship both rows; mention it in the commit rather than inventing a suppression rule.
- **C3 Payload growth.** `board-data.js` grows roughly 5-8x for the activity array (14 possible stamps per ticket, most tickets carrying 3-6). Every archived REQ ships all its stamps, and windowing is client-side. On a queue of ~570 REQs that is a few thousand rows of four short strings, so it is fine, but it is the one size effect worth checking against the generated file.
- **C4 `data-activity-request` is no longer unique.** Any future selector assuming one node per REQ breaks. No current test or CSS rule depends on uniqueness (verified by grep), and REQ-573 wants the non-unique form for sibling highlighting.
- **C5 `model.go` is off-limits** and does not need touching: `lifecycleTimestampFields` already returns every field with its transition, so the Go change is confined to `activity.go`.
- **C6 Boundary test index assumption.** `TestBuildActivityRowsStraddlesTheWindowBoundary` indexes `rows[0]`/`rows[1]` by position; keep the one-stamp-per-ticket fixture so those indices stay meaningful.
