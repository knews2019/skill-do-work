---
id: UR-026
title: The board's By UR lens goes permanently blank once a queue is fully shipped
created_at: 2026-08-06T15:43:07Z
requests: [REQ-122]
word_count: 1180
---

# The Board's By UR Lens Goes Permanently Blank Once a Queue Is Fully Shipped

## Full Verbatim Input

> # By UR lens goes permanently blank on a cleared queue
>
> The board's **By UR** lens renders an empty list the moment a queue is fully shipped, and the RECENTLY DONE window buttons are a dead knob in that lens. Both are in `tools/queue-kanban/`. The investigation below was done downstream against a real 390-UR / 1052-REQ tree and the fix was prototyped end-to-end and verified live, so everything here is confirmed, not guessed. Line numbers are from the vendored copy at **v0.178.1** — anchor on the quoted code, not the numbers.
>
> ## Symptom
>
> Lens = **By UR**, URS = **Active**, RECENTLY DONE = **7d**. The list shows only:
>
> > No active user requests — every UR is fully resolved. Switch URs to All to browse the archive.
>
> Meanwhile the **Columns** lens, same board, same moment, shows 17 recently-done REQ cards every one of which is tagged `UR-489` — work that finished four hours earlier. So the board simultaneously claims there is nothing to show and shows it.
>
> Reproduces on any tree where `do-work/queue/` and `do-work/working/` are empty and every archived REQ is `completed` / `completed-with-issues` / `cancelled`. That is the *normal steady state* after a run finishes, not an exotic one — which is why the lens is effectively unusable for the reader who just wants to see what a session touched.
>
> ## Root cause
>
> Two independent defects that only bite together.
>
> **1. `userRequestIsActive` is always false once a queue ships** — `web/board.js:367`:
>
> ```js
> function userRequestIsActive(userRequest) {
>   return (userRequest.requestIds || []).some(function (requestId) {
>     var request = requestsById[requestId];
>     return request && !isTerminalResolvedStatus(request.status);
>   });
> }
> ```
>
> `isTerminalResolvedStatus` (`board.js:294`, mirroring `model.go:806`) covers `completed | completed-with-issues | cancelled`. Every UR therefore fails the test, and the gate at `board.js:835` drops all 390. Not a data problem: `generate.go:303-315` emits **every** UR unfiltered, `walk.go:160-163` collects them from both `user-requests/**` and `archive/**`, so `URS = All` populates fine. Active/All is a pure render-time filter.
>
> **2. The RECENTLY DONE window never reaches this lens.** `viewState.windowHours` is read only by `recentlyDoneIds()` (`board.js:700`), which is called only by `renderColumns()` (`board.js:734`). `renderUserRequestLens()` never mentions either. The click handler (`board.js:2053-2058`) hardcodes `renderColumns()` and does not even invalidate `renderedOnce.userRequestLens`. And `applyView()` (`board.js:1993`) hides `#recent-window-group` only when the *view* isn't the board — so in the By UR lens those three buttons stay visible, repaint hidden columns, and change nothing the reader can see.
>
> Two smaller defects sit in the same block and should go with it:
>
> - **The hidden-count note is unreachable exactly when it matters.** `board.js:912-918` early-returns on the empty branch, so the `hiddenResolvedCount > 0` note at `board.js:919` never fires when *all* URs are hidden. The reader is never told that 386 URs are waiting behind the toggle — precisely the case where they'd want to know.
> - **The empty-state copy is unconditional.** `hasActiveFilters()` (`board.js:316`) checks only `searchText / domain / status / doneWindow` — never `userRequestActivity` — so the "Switch URs to All" line would still print while All is already selected.
>
> ## The fix (direction is decided — don't re-litigate it)
>
> Widen **Active** to mean *open work **or** a REQ completed inside the current recently-done window*, and make the window buttons drive both lenses. The alternatives (a third `Recent` scope button; auto-falling back to `All`) were considered and rejected: the first adds a knob for a distinction nobody makes, the second desyncs the toggle from what's on screen.
>
> 1. **Replace the predicate** with one that takes a recently-done id set and short-circuits on it before the status check. Build the set from the existing `recentlyDoneIds(viewState.windowHours)` — reusing that function is the point, since it already anchors to the wall clock rather than `generatedAtMs` and it guarantees the two lenses can never disagree about what "recent" means. Build it **once** per render in `renderUserRequestLens`, not per UR.
> 2. **Wire the window handler** to re-render whichever lens is up and drop the other's cache — copy the shape of the `[data-ur-activity]` handler at `board.js:2097-2107`, which already gets this right.
> 3. **Fix the empty/hidden block**: drop the early `return` so the hidden count still renders under the empty text; branch the empty copy three ways (filters matched nothing / nothing in window / tree genuinely empty); derive the window phrase from `viewState.windowHours` so it reads "the last 7 days" and tracks the selected chip rather than baking in a span. Keep the note **silent** when a search matched nothing — there the scope isn't why the list is empty, so "switch to All" is a false lead.
> 4. **`web/template.html`**: the `data-ur-activity="active"` button needs a `title` spelling out the widened rule, and the stale comment above the group ("hide user requests whose REQs are all resolved") needs to match the new behaviour.
>
> A UR that qualifies only via the window should render **all** its REQ cards, not just the recent ones — the UR is the unit of this lens and the `ur-count` chip already states the total.
>
> ## Tests
>
> `board.js` has no JS harness; the shipped convention is Go tests asserting against the **inlined** board source in the generated page — see `generate_test.go:74` ("inlined board.js behaviour is missing from the page") and the `data-window-hours="168"` / `aria-pressed` assertions at `generate_test.go:504-510`. Follow it. Both halves of this bug are invisible to markup-only assertions and to every Go model test, and neither reproduces unless the queue happens to be at zero — so pin:
>
> - the new predicate name **and** its call site (a half-finished rename would leave the lens filtering on the old rule with no symptom until the queue next hits zero);
> - the absence of the superseded `userRequestIsActive`;
> - `renderUserRequestLens` appearing **inside the `[data-window-hours]` handler body** (slice the source between the handler and the next statement and assert on that slice — a whole-page `strings.Contains` would pass on the buggy version);
> - the empty-state copy naming the window and both escapes;
> - the template `title` attribute.
>
> Gate: `go test ./...` in `tools/queue-kanban`.
>
> ## Verification to reproduce before shipping
>
> Build the tool and `serve --repo-root <a repo with a fully-shipped do-work tree> --port 8099`. Note `web/*` is `go:embed`-ed (`generate.go:20`) — an already-running server keeps serving the old JS, so rebuild and restart or you will "verify" a stale bundle.
>
> With the prototype applied, driving the real page:
>
> | State | Before | After |
> | --- | --- | --- |
> | Active · 7d | empty | UR-486, UR-487, UR-488, **UR-489 (16 REQ, 16 cards)** |
> | Active · 24h | empty | UR-488, UR-489 |
> | Active · 48h | empty | UR-488, UR-489 |
> | All | 390 URs | 390 URs (unchanged) |
> | Hidden note, Active · 7d | never rendered | "386 URs with no open work or activity in the last 7 days hidden — switch URs to All to see them." |
> | Search matching nothing | "No user requests match the current filters." | unchanged, and the hidden note correctly stays silent |
> | Columns → change window → back to By UR | stale | reflects the new window |
>
> ## Housekeeping
>
> Add a CHANGELOG entry in the established voice (state what was broken and what it does now, not "improved the lens"), bump the version, and update the by-UR paragraph in `tools/queue-kanban/prime-do-kanban.md` if it describes the old Active semantics. If this repo captures its own work through the do-work flow, capture it as a REQ first rather than patching straight in.

## Context

An upstream bug report from a downstream consumer of the shipped `tools/queue-kanban/`
board, filed against the vendored copy at **v0.178.1** — which is this repo's current
version, so the quoted line numbers land exactly and no anchoring drift applies.

The report is unusually complete: the reporter read the actual source, traced both
defects to specific functions, ruled out the data layer (`generate.go` / `walk.go` emit
every UR unfiltered), prototyped the fix, and drove the real page to produce the
before/after table. The claims were re-verified against this tree before implementation
rather than taken on trust.

The design call is stated as decided, with two alternatives explicitly considered and
rejected. Reviewing them independently: widening `Active` is the right call. A third
`Recent` scope button splits one question ("what should I be looking at?") across two
knobs, and an automatic fallback to `All` would leave the toggle claiming `Active` while
displaying the archive — a worse lie than the current blank list. No pushback.

## Batch Constraints

- The `web/` assets are `go:embed`-ed, so a running `serve` keeps serving the old JS —
  any live verification must rebuild and restart, or it verifies a stale bundle.
- `board.js` is ES5-style throughout (`var`, `function`, no `Set`, no arrow functions)
  and is inlined verbatim into the generated page. New code must match that floor.
- There is no JS test harness. The shipped convention is Go tests asserting against the
  inlined board source in the generated page, and this bug is invisible to markup-only
  assertions — a whole-page `strings.Contains` would pass on the buggy version, so the
  handler-body assertion has to slice the source.
