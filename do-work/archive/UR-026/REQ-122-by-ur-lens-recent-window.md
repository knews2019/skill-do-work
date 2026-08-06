---
id: REQ-122
title: The By UR lens counts recently-done work as active and honors the window buttons
status: completed
created_at: 2026-08-06T15:43:07Z
claimed_at: 2026-08-06T15:44:00Z
completed_at: 2026-08-06T15:51:16Z
commit: 684f507
route: A
kb_status: pending
user_request: UR-026
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: []
write_set: [tools/queue-kanban/web/board.js, tools/queue-kanban/web/template.html, tools/queue-kanban/generate_test.go]
maintenance: false
related: []
batch: upstream-by-ur-lens-report
---

# The By UR Lens Counts Recently-Done Work as Active and Honors the Window Buttons

## What

On a fully-shipped queue — `do-work/queue/` and `do-work/working/` empty, every archived
REQ terminal — the board's **By UR** lens with `URs = Active` renders nothing, while the
**Columns** lens on the same board at the same moment shows the recently-done cards for
those very URs. The board simultaneously claims there is nothing to show and shows it.
That is the normal steady state after a run finishes, so the lens is unusable for the
reader who just wants to see what a session touched.

Four defects, all in `tools/queue-kanban/web/`:

1. **`userRequestIsActive` is unsatisfiable once a queue ships.** It passes only for a UR
   holding a non-terminal REQ; when every REQ is `completed` / `completed-with-issues` /
   `cancelled`, every UR fails and the gate in `renderUserRequestLens` drops all of them.
   Not a data problem — `generate.go` emits every UR unfiltered and `walk.go` collects
   from both `user-requests/**` and `archive/**`, which is why `URs = All` populates fine.
   Active/All is a pure render-time filter.
2. **The RECENTLY DONE window never reaches this lens.** `viewState.windowHours` is read
   only by `recentlyDoneIds()`, called only by `renderColumns()`. The `[data-window-hours]`
   click handler hardcodes `renderColumns()` and does not invalidate
   `renderedOnce.userRequestLens`, while `applyView()` hides `#recent-window-group` on
   view rather than lens — so in the By UR lens the three chips stay visible, repaint
   hidden columns, and change nothing the reader can see.
3. **The hidden-count note is unreachable exactly when it matters.** The empty branch
   early-returns before the `hiddenResolvedCount > 0` note, so when *all* URs are hidden
   the reader is never told how many are waiting behind the toggle.
4. **The empty-state copy is unconditional.** `hasActiveFilters()` does not consider
   `userRequestActivity`, so the "Switch URs to All" line could print while All is
   already selected.

## Approach

Widen **Active** to mean *open work **or** a REQ completed inside the current
recently-done window*, and make the window chips drive both lenses.

Two alternatives were considered upstream and rejected, and the rejection holds on
independent review: a third `Recent` scope button splits one question across two knobs,
and auto-falling back to `All` would leave the toggle reading `Active` while displaying
the archive — a worse lie than a blank list.

- Replace the predicate with one taking a recently-done id set, short-circuiting on it
  before the status check. Build the set from the existing `recentlyDoneIds()` — reuse is
  the point: that function anchors to the wall clock rather than `generatedAtMs`, and
  sharing it is what guarantees the two lenses can never disagree about what "recent"
  means. Build it once per render, not once per UR.
- A UR that qualifies only via the window renders **all** its REQ cards, not just the
  recent ones — the UR is the unit of this lens and the `ur-count` chip already states
  the total.
- Wire the window handler to re-render the visible lens and drop the other's cache,
  mirroring the `[data-ur-activity]` handler that already gets this right.
- Drop the early `return` so the hidden count renders under the empty text; branch the
  empty copy three ways (filters matched nothing / nothing in window / tree genuinely
  empty); derive the window phrase from `viewState.windowHours` so it tracks the selected
  chip. Keep the note silent when a filter — not the scope — emptied the list, where
  "switch to All" would be a false lead.
- Give the `Active` button a `title` stating the widened rule and correct the stale
  comment above the group in `template.html`.

## Verification Plan

`board.js` has no JS harness; the shipped convention is Go tests asserting against the
**inlined** board source in the generated page (`generate_test.go`). Both halves of this
bug are invisible to markup-only assertions and to every Go model test, and neither
reproduces unless the queue is at zero, so the tests pin:

- the new predicate name **and** its call site — a half-finished rename would leave the
  lens filtering on the old rule with no symptom until the queue next hits zero;
- the absence of the superseded `userRequestIsActive`;
- `renderUserRequestLens` appearing **inside** the `[data-window-hours]` handler body,
  asserted against a slice of the source between the handler and the next statement — a
  whole-page `strings.Contains` would pass on the buggy version;
- the empty-state copy naming the window and both escapes;
- the template `title` attribute.

Gate: `go test ./...` in `tools/queue-kanban`, plus `_dev/tests/contract-regressions.sh`.

## Out of Scope

`tools/queue-kanban/prime-do-kanban.md` was checked for a by-UR paragraph describing the
old Active semantics, as the report suggested. It has none — the file records hard-won
implementation lessons, not UI behaviour, and never mentions the lens. No edit needed.

## Implementation Summary

Files changed:

- `tools/queue-kanban/web/board.js` — replaced `userRequestIsActive` with
  `userRequestHasOpenOrRecentWork(userRequest, recentlyDoneIdSet)`; added
  `recentWindowPhrase(windowHours)`; built the recently-done id set once per render in
  `renderUserRequestLens`; reworked the empty/hidden block into three branches with the
  note rendering under the empty text; extended the `[data-window-hours]` handler to
  refresh the visible lens and drop the other's cache.
- `tools/queue-kanban/web/template.html` — `title` on the Active scope button stating the
  widened rule; corrected the stale comment above the group.
- `tools/queue-kanban/generate_test.go` — four tests plus a `sliceBalancedBlockAfter`
  helper that brace-matches a block out of the inlined source.

`hasActiveFilters()` was deliberately **not** changed. The report flagged that it ignores
`userRequestActivity`, but it is shared with the calendar summary, where the scope is
genuinely irrelevant. The three-way empty branch fixes the reported symptom structurally
instead: the "switch URs to All" copy now sits behind `hiddenResolvedCount > 0`, which can
only be non-zero in Active mode, so it cannot print while All is selected.

## Verification Results

- `go test ./...` in `tools/queue-kanban`: green. All four new tests confirmed RED first
  (superseded predicate present, handler body missing the lens call, stale empty copy,
  no `title` on the chip).
- Live, against two fixture trees with the binary rebuilt after the edit (the `go:embed`
  stale-bundle trap the report warns about):
  - **Fully-shipped tree, 25 archived URs, nothing open.** Active/24h → 5 URs; Active/48h
    → 8; Active/7d → 16; All → 25. Every shown UR rendered all its REQ cards (UR-018
    showed 13 of 13). Hidden note tracked the chip: "20 URs with no open work or activity
    in the last 24 hours hidden". Search matching nothing → filter copy, note correctly
    silent. Columns → change window → back to By UR reflected the new window.
  - **All-stale tree, 9 URs, newest completion 7.8 days old.** Empty copy read "No user
    requests with open work or activity in the last 7 days — widen the RECENTLY DONE
    window, or switch URs to All to browse the archive.", with the hidden note rendered
    underneath it in DOM order. This is the branch the first fixture could not reach.
- No JS errors. The one 404 observed is `/favicon.ico`, pre-existing and unrelated.
- `_dev/tests/contract-regressions.sh` fails on 8 update-script probes (mid-update
  failure handling, dirty install). Verified identical on a stashed clean tree —
  pre-existing, not caused by this REQ, and left alone as out of scope.

## Lessons Learned

- **A predicate that reads as a property of the data can silently be a property of the
  clock.** `userRequestIsActive` looked like a pure status question and was really "is
  this queue mid-flight", so it was correct in every state anyone tested in and
  unsatisfiable in the state the board spends most of its life in. The tell was two views
  of one dataset disagreeing on screen at the same moment — worth treating as a
  high-signal bug shape rather than a rendering glitch.
- **Assert on a slice, not the page, when the claim is "X is called from Y".** Both
  `renderUserRequestLens` and `renderColumns` appear all over the inlined bundle, so a
  whole-page `strings.Contains` would have passed against the buggy version and pinned
  nothing. Brace-matching the handler out of the source is what makes the test able to
  fail — and it did fail first, which is the only reason it is worth keeping.
- **A cache invalidated in one handler needs its siblings checked.** The `renderedOnce`
  guard was handled correctly by the `[data-ur-activity]` handler and not at all by
  `[data-window-hours]`. Whenever a render cache exists, every control that changes an
  input to that render is a call site — grep the guard, not the feature.
