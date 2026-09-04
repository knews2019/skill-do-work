# Activity history preparation — REQ-572 and REQ-573

Read-only source inspection at `f6c43d221447c150e7a8f07b5f023f73fedc9559`. No tests, claims, source edits, queue edits, worktrees, or background jobs were created. This artifact is preparation, not a Plan adopted by a claimed request and not RED/GREEN evidence.

REQ-572 was already claimed in `do-work/working/REQ-572-show-every-lifecycle-transition-of-a-req-as-its-own-activity-row.md` when inspected, with `claimed_at: 2026-09-04T23:21:41Z`. Preserve that other owner's work. REQ-573 remains pending and depends on REQ-572; resolve its identity, dependency, source revision and existing changes again before any execution.

Read CLAUDE, the core SKILL router, both complete REQs and UR-115/input.md; inspected both UR-115 screenshots. Read the maintainer and shipped board primes, both whole board lesson satellites, shared principles, frontend, testing and prompt-injection crew. The capture-time lesson-budget drops remain historical; these whole-file touch-conditional reads are additive and do not retroactively imply narrower family reads satisfied those drops.

## Shared behavior and owners

- `skills/do-work-board/tools/queue-kanban/model.go`: `lifecycleTimestampFields` is the sole stamp/transition enumeration, currently 14 fields. The activity and future-stamp readers share it. `parseTimestamp` already defines acceptable timestamp shapes; `statusChangeTransitionForStatus` and `completionTransitionForStatus` phrase status-dependent meanings. Do not duplicate those definitions in JavaScript or edit this file for REQ-572; REQ-571 may independently change its heavy-status reader.
- `activity.go`: `buildActivityRows` currently picks `newestLifecycleStamp` once per ticket, then sorts by descending time and descending RequestId. It deliberately ships out-of-window and future parseable stamps; the client applies the live window. New same-request/same-instant ties need a deterministic third tie-break because the existing request-id tie-break cannot distinguish them.
- `generate.go:860`: `buildGeneratedBoardData` projects every returned ActivityRow verbatim to `generatedActivityEntry` (`id`, `stampField`, `stampAt`, `transition`). The payload shape already accepts multiple rows for one id; no schema expansion is needed. Generate and serve share embedded assets; rebuild before visual verification.
- `web/board-activity.js`: `activityRowsWithin` uses Date.now(), strict `stampMs > cutoffMs`; `renderActivity` applies `requestMatchesFilters`, constructs table rows with `data-activity-request`, and currently counts rows as REQs. The current filtered empty-state also calls windowRows.length REQs, so it would become wrong under repeated ids.
- `web/board-controls.js:229`: existing document click delegation resolves the closest `[data-detail-kind]` and calls `openDetail`; `detail-close` calls `closeDrawer`. No new Activity drawer opener belongs here.
- `web/board-detail.js`: `currentDetailKind/currentDetailId` (306–307) already own the open identity. `openRequestDetail` assigns them at 598–599, `openUserRequestDetail` at 623–624; both finish in `showDrawer`. `closeDrawer` resets both and returns focus, and Escape routes through it. There is no current drawer-change event or Activity selection callback.
- `web/template.html:468`: existing Activity section, six-column semantic table, named header ids, window controls and empty-state nodes. `web/board.css:3151` owns table metrics, scrolling, sticky headers and wrapping title cell. Drawer layout is shared and remains on the right as in the screenshot; the user's “left side” means that existing drawer, not relocation.
- `activity_test.go` contains the existing Go aggregation tests, including deliberate old one-row assertions that must change with REQ-572. `javascript_behavior_c_test.go` uses the existing Go/Node generated-client probes. `browser_probe_test.go` provides actual generated-page probes and trusted input sessions; use those existing facilities, not another framework.

The implementation can expose every **currently recorded parseable stamp**. Frontmatter is not an append-only event log: overwritten repeated claims/status changes cannot be reconstructed here. Do not add git-history reads, synthesized transitions, or a new event store. Keep this distinction truthful in orientation; it does not block the captured requirements, which explicitly name the stamp projection.

## REQ-572 — proposed three-task plan

1. Establish RED at aggregation and client seams. Add the captured created/claimed two-row case before changing production, plus another ticket interleaving between them. Change existing newest-only assertions to pin all parseable stamps, invalid/empty/nil skip behavior, same-instant distinct fields without deduplication, descending RequestId ties under reversed input, and deterministic same-id field ties. Drive `buildGeneratedBoardData` from a real ticket fixture to prove all rows survive the actual projection. Add a Node render case with three transitions across two REQs and filtered/empty variants so counts cannot accidentally remain row counts labelled as REQs.
2. Replace the per-ticket newest selection with one append for each parseable entry of `lifecycleTimestampFields`. Preserve descending time then descending RequestId; choose and document a stable third field-name tie-break (no new stamp-order enumeration). Remove `newestLifecycleStamp` if it has no other caller. Keep the payload fields and Go-owned phrases intact. In the client count distinct ids from filtered rows and from windowRows as needed: “3 transitions across 2 REQs in the last 24 hours”, grammatical singular variants, and explicit filtered counts when filters narrow the window. Keep both empty-state meanings and no toggle. Update obsolete newest-only comments/accessibility prose; no cosmetic CSS change is needed merely to emit more rows.
3. Run the focused Go and Node cases, the prescribed full gate/heavy lane selection through the owning run, and inspect rebuilt static plus live board behavior. With multiple rows for REQ-570, verify two real stamps appear in timestamp order and summary totals match filtered data, including 6h/24h/48h/7d changes and a clock-advanced boundary fixture. Inspect 320/768/1280 widths and both themes for readable column headings and title wrapping; capture page URL with measurements, Chromium build, screenshot paths and console results. Do not claim tests ran during this prep.

Exact proposed implementation/test Scope (all paths under `skills/do-work-board/tools/queue-kanban/`):

| Path | Purpose |
|---|---|
| `activity.go` | Emit every parsed stamp; deterministic total order; remove newest-only helper/comments |
| `activity_test.go` | RED/GREEN aggregation, ties, skip, actual payload projection controls |
| `generate.go` | Update obsolete producer comment only if present; projection stays structurally unchanged |
| `web/board-activity.js` | Distinct request and transition counts; filtered/empty wording; obsolete comment correction |
| `web/template.html` | Activity description/accessibility prose matching transitions |
| `javascript_behavior_c_test.go` | Actual shipped Activity renderer/window/filter summary tests |

The captured write_set's `web/board.css` is unnecessary for the minimal REQ-572 change; the Node behavior test file is necessary but absent from capture. Adopt that replacement in Scope before coding and record the reason. If the active owner already chose a smaller or different tested plan, review the committed result rather than imposing this preparation.

## REQ-573 — proposed three-task plan

1. On the accepted REQ-572 base establish RED by rendering two REQ-570 rows and one REQ-505 row from a production-shaped payload. Dispatch through the real document delegation, then assert drawer identity and exactly two selected rows. Cover switching to REQ-505, changing the window/filter with the REQ still open, hiding then restoring its rows, closing by button and Escape, and navigating from a REQ to its UR drawer. Include an unknown id control so a failed open does not mark a fake selection. Test a real button/link semantics assertion in the Node lane and trusted Tab/Enter/Space behavior in actual Chromium.
2. Build the REQ header cell with a native `button type="button"`, `data-detail-kind="req"`, `data-detail-id` and visible id using textContent. Existing document delegation remains the sole opener. Derive selected rows from `currentDetailKind/currentDetailId`; do not add a second selected-id store. Add an Activity synchronization helper that toggles a selected class on existing rows, called after rendering and by `showDrawer`/`closeDrawer` after their identity changes. Using the shared drawer lifecycle also clears stale marks when navigation opens a UR and avoids repainting the entire Activity table on every open. No observer, polling, additional click handler or new generic event system is needed. Style a visible background plus a left rule on the row-header cell and a clear focus-visible button outline; preserve table semantics and horizontal scrolling.
3. Execute the focused Node tests and an actual generated-page browser probe through the existing harness, then the run's canonical verification lanes. Capture light/dark screenshots at 320/768/1280; record exact browser build and page URL in each measurement. Use trusted keyboard input, not synthetic Tab, to prove reachability and return focus. Verify the same drawer title/metadata/body as the Board, every matching visible row marked, unrelated ids unmarked, marks surviving rerender, and marks gone on both close paths. Measure actual computed colors against rendered row/page backgrounds and inspect the non-color left rule at narrow width, not only a class presence assertion.

Exact proposed implementation/test Scope (all paths under `skills/do-work-board/tools/queue-kanban/`):

| Path | Purpose |
|---|---|
| `web/board-activity.js` | Native delegated REQ trigger; derived row highlight helper and rerender synchronization |
| `web/board-detail.js` | Call highlight synchronization from existing show/close lifecycle after identity is correct |
| `web/board.css` | Selection background and non-color cue; native trigger/focus styling |
| `javascript_behavior_c_test.go` | RED/GREEN actual delegation, drawer state, sibling selection and rerender/close controls |
| `activity_browser_probe_test.go` (new) | Bounded real generated-page click/keyboard and themed visual measurement seam using existing browser harness |

The captured three-file write_set omits `board-detail.js`. Add that owner to Scope before coding: the requirement to clear on every drawer close and follow the currently open REQ is owned there, not in an Activity click handler. The new focused browser test file keeps the already large test shards from absorbing another browser runtime burden; retain each fast file's 30-second limit and select strict browser tests with explicit Chromium. If a manual visual check covers all acceptance without a durable new probe, the browser file can be omitted by an explicit Plan decision, but keyboard behavior must still be verified with trusted input.

## Preservation and acceptance notes

Both requests are display-only. Preserve exact Markdown payload bytes, clipboard behavior, queue state, all existing write-surface boundaries and live/static parity. Preserve the source order delivered to JavaScript; no client sorting or timestamp-meaning inference. Current status remains the ticket's current status; the transition cell explains the recorded stamp. Existing filter semantics stay ticket based and therefore select all of a ticket's in-window transitions. No “latest only” toggle, grouping redesign, drawer repositioning, new package or test framework is required.

A large row set is a foreseeable rendering concern under the frontend guidance. Measure the actual in-window row count and responsiveness after REQ-572; do not silently introduce a virtualization subsystem into its minimal change. If more than 100 rows are simultaneously visible, resolve the frontend virtualization condition explicitly in the claimed Plan, preserving semantic keyboard access and count-from-data semantics.
