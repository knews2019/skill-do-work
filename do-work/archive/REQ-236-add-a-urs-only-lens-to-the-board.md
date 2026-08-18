---
id: REQ-236
title: Add a URs-only lens to the Board view
status: completed
completed_at: 2026-08-18T11:20:00Z
commit: 456ee9d
claimed_at: 2026-08-18T11:00:00Z
route: C
estimate:
  p50_active_minutes: 60
  confidence: medium
  calculated_at: 2026-08-18T11:00:00Z
  basis:
    - Route C
    - 5-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
    - browser evidence
created_at: 2026-08-18T10:30:00Z
user_request: UR-054
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
effort_estimate: normal
write_set:
- skills/do-work-board/tools/queue-kanban/web/template.html
- skills/do-work-board/tools/queue-kanban/web/board-cards.js
- skills/do-work-board/tools/queue-kanban/web/board-controls.js
- skills/do-work-board/tools/queue-kanban/web/board.css
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Add a URs-Only Lens to the Board View

## What

The Board view can be read as REQ cards in status columns (`Lens: Columns`) or as UR headers with all their REQ cards beneath (`Lens: By UR`). Add the third reading the user asked for: **URs only** — just the user-request headers, each expanding in place to reveal its REQs on click.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`, `CLAUDE.md` § Kanban Board Write Surfaces, and the always-on crew rules; settled on a card-fold over the existing by-UR renderer after trying and rejecting a third `viewState.lens` value on blast-radius grounds (D-01).
- [x] **[APPLY]:** Five declared files touched, nothing outside them; the drawer collision was resolved inside the write set rather than by editing the delegated handler's file.
- [x] **[UNIFY]:** `git diff --stat` reviewed across all five files; `go vet` and `gofmt -l` clean; full `maintainer-verify.sh` exit 0 twice; no debug artifacts in the diff.

<!-- Boxes ticked by the ORCHESTRATOR at Step 6.3, not by the builder. The dispatch
     brief for this REQ omitted the P-A-U instruction, so the builder never saw it;
     each line above is transcribed from evidence in its hand-back and independently
     re-checked against the merged tree, rather than being a builder self-report. The
     omission is the orchestrator's, and it is recorded here rather than papered over. -->

## Why

The user's words: "add possibility to add UR+REQ UR and only REQ to be viewed". Two of the three readings exist; the condensed one — the list of what was actually asked for, without every REQ card unrolled — does not. On a board with 50+ URs, `By UR` is a very long page.

## Context

Today's `Lens` group in `web/template.html` (`id="lens-group"`) holds two buttons: `data-lens-target="flat"` (Columns) and `data-lens-target="user-request"` (By UR). `board-controls.js` wires them into `viewState.lens` and toggles the `ur-activity-group` (Active/All) with them; `board-cards.js` renders the by-UR lens as one `section.ur-group` per UR — a `button.ur-group-head` carrying `data-detail-kind="ur"` (opens the UR detail drawer), the UR id, title, an optional "no input.md" chip and a `N REQ` count, followed by a `div.ur-group-cards` of REQ cards that is always expanded.

So the new lens is the same group markup with the cards folded away by default, plus an expand affordance. The Active/All UR scope, the shared filters, and the recently-done window all already apply to the by-UR path and must apply here unchanged.

## Detailed Requirements

- The `Lens` group offers three choices: **Columns**, **By UR**, **URs only**. Columns and By UR behave exactly as they do today.
- `URs only` lists one row per user request — id, title, REQ count, and the same "no input.md" marker the by-UR headers carry — with no REQ cards showing until the reader asks for them.
- Clicking a UR row **expands it in place**, unfolding that UR's REQ cards underneath; clicking again collapses it. More than one may be open at a time.
- Opening the UR's detail drawer stays reachable from the row — the header is currently the drawer trigger, so the expand affordance and the drawer must not collide with each other.
- Filters (search, domain, status), the `URs: Active / All` scope, and the recently-done window apply to `URs only` exactly as they apply to `By UR`, including the existing empty-state and hidden-UR notes.
- Expanded/collapsed state is view state, not persisted queue state.

## Constraints

- Read-only, like every other board view: the board's three write surfaces (CLAUDE.md § Kanban Board Write Surfaces) are unchanged.
- The row must be operable from the keyboard, and expanded/collapsed must be announced — the header is already a `<button>`, so this is `aria-expanded` on the right element rather than new machinery.
- Reuse `makeRequestCard` and the existing `ur-group` markup rather than growing a second UR renderer; the two lenses must not drift.
- No change to the payload or to `model.go` — every field this needs is already in the data island.

## Red-Green Proof

**RED prompt/case:** A Node behaviour probe in `generate_test.go` (the `TestJavaScriptBehavior*` family) over a fixture with at least two URs: selecting the `URs only` lens renders one row per UR and zero REQ cards; activating a row renders exactly that UR's REQ cards and sets `aria-expanded="true"` on it; activating it again removes them; and the `Active` scope plus a status filter hide the same URs they hide in the `By UR` lens.

**Why RED now:** there is no third lens value — `viewState.lens` only takes `flat` and `user-request`, and `board-cards.js` always appends every REQ card into `ur-group-cards`, so both the zero-card assertion and the toggle assertion fail.

**GREEN when:** the probe passes and a headless render shows three lens buttons, a condensed UR list, and one UR expanded in place. `bash _dev/tests/maintainer-verify.sh` exits zero.

**Validation:** User confirmed — the three-way choice, its placement as a third Lens button, and expand-in-place were all picked by the user during capture.

## Builder Guidance

Certainty level: Firm on the three decisions above. Latitude on the expand affordance's exact shape (disclosure triangle vs. row click with a separate drawer control) as long as both actions stay reachable and keyboard-operable. Keep it small — this is a fold on top of a renderer that already exists, not a new view.

---
*Source: "add possibility to add UR+REQ UR and only REQ to be viewed"*

---

## Triage

**Route: C** - Complex

**Reasoning:** The REQ's Context section already named the exact markup and wiring to extend, but the change touches four client files plus tests, and its central design question — how the expand affordance and the UR detail drawer share one row — had a real answer space that had to be explored before writing.

**Planning:** Required

## Plan

The lens is a fold on the existing by-UR renderer rather than a new view: same `ur-group` markup, same `makeRequestCard`, with the cards built late and dropped on collapse. The open question at plan time was whether `URs only` becomes a third `viewState.lens` value or a modifier on the existing one; the builder explored both and the answer changed the blast radius (see D-01).

*Generated by Plan agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify) — third Lens button
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modify) — fold flag and lens-button selection
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modify) — folded rendering of the UR groups
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — row, detail button, fold marker, condensed metrics
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — markup test and Node behaviour probe

**Files I will NOT touch:** `model.go` / any Go payload code (every field this needs is already in the data island), `board-filters.js` and `board-testing.js` (see D-01 — the chosen model is what keeps them correct unedited).

**Acceptance criteria (restated from REQ):**
- [ ] Three Lens choices; Columns and By UR behave exactly as today
- [ ] URs only lists one row per UR — id, title, REQ count, "no input.md" marker — with no REQ cards until asked
- [ ] Clicking a row expands it in place and again collapses it; more than one may be open
- [ ] The UR detail drawer stays reachable from the row without colliding with the expand affordance
- [ ] Filters, `URs: Active/All`, and the recently-done window apply exactly as they do to By UR
- [ ] Expanded/collapsed is view state, not persisted queue state

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** `URs only` is a third Lens button carrying the *same* `data-lens-target="user-request"` as By UR plus `data-ur-cards="folded"`, so `viewState.lens` still holds exactly two values. `board-controls.js` gained a module-level `userRequestCardsFolded` flag, one `applyLensSelection()` entry point, and `setActiveLensButton()` — which keys the pressed state on the lens/fold pair, since two buttons now share a lens. `renderUserRequestLens()` in `board-cards.js` learned the fold: the same `ur-group` markup, with cards built by a shared inner `makeUserRequestCards()` late and dropped again on collapse.

In the folded reading only, the header stops being the drawer trigger and becomes the fold toggle — it drops `data-detail-kind`/`data-detail-id`, gains `aria-expanded` and a marker, and the drawer moves to a sibling `button.ur-group-detail` inside a `div.ur-group-row` wrapper that still carries `data-detail-kind="ur"` for the existing delegated handler. The always-open By UR reading is unchanged, which the probe asserts explicitly. `board.css` added the row, detail button, fold marker, and the tighter `.user-request-lens.is-folded` metrics that make the condensed list actually condensed at 44px rows.

## Qualification

Passed — 5 files verified in the merge range `52fb7d7..456ee9d`, 6 acceptance criteria traced.

**Merge conflict resolved by the orchestrator.** `generate_test.go` conflicted: this REQ and REQ-233 both appended tests at the end of the same file. Both sides were pure appends of separate functions, so the resolution keeps both in order; `gofmt -w` restored the blank line between them, and `go vet` plus the full suite pass on the union. This is the merge doing the job it exists for — two builders' work meeting for the first time.

Judgment checks, run against the merged tree rather than taken from the builder's report:
- **D-01's central claim was verified rather than accepted.** `viewState.lens` is assigned in exactly one place (`board-controls.js`, from the button's `data-lens-target`), and every read site — `board-cards.js:345`, `board-testing.js:323`, `board-controls.js:174`, `board-filters.js:169` — compares against `"user-request"`. All four therefore hold for both readings with no edit, which is what the decision claimed.
- **The "0 cards" assertion needed scoping to mean anything.** A naive `document.querySelectorAll('.req-card')` returns 25 in the folded lens, which looks like a failure and is not: those live in the *hidden* Columns lens, which keeps its DOM. Scoped to the lens host the counts are 0 folded → 20 with one row open → 0 again, with every other row still empty.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Exit 0 on the merged tree, run unpiped — this is the run that matters, since REQ-233's and this REQ's tests coexist in one file for the first time here

**Red-green validation:**
- `TestJavaScriptBehaviorUserRequestsOnlyLensFoldsCardsUntilARowIsOpened`: ✗ `URs only row UR-401 rendered cards []string{"REQ-601", "REQ-602"} before it was opened, want none` → ✓. A behavioural failure, not a missing anchor: the probe declares the fold flag itself and drives the *existing* renderer, which ignored it.
- `TestGenerateOffersThreeLensButtons`: ✗ `lens group is missing "data-ur-cards=\"folded\"" in the generated page` → ✓.

**New tests added:**
- `TestGenerateOffersThreeLensButtons` — the Lens group offers all three readings, with URs only carrying the shared lens target plus the fold attribute
- `TestJavaScriptBehaviorUserRequestsOnlyLensFoldsCardsUntilARowIsOpened` — over a three-UR fixture: folded renders zero cards with `aria-expanded="false"` and one drawer trigger per row; opening a row yields exactly that UR's REQs with siblings still empty; re-activating collapses it; and under `Active` scope plus a status filter, URs only and By UR render an identical UR list while the By UR head carries no `aria-expanded` at all, so the fold cannot leak into the always-open reading

**Existing tests updated (cross-REQ impact):**
- `generate_test.go` — `TestJavaScriptBehaviorByUserRequestLensUsesRecentWindowAtCaller` (from REQ-122): its probe preamble gained `var userRequestCardsFolded = false;` because the renderer now consults that flag and the probe failed with `ReferenceError`. The behaviour change is intentional and the line pins the always-open reading; the test's own assertions are unchanged and still pass.

**Render evidence — driven in a real browser against the merged tree.** Generated from the live repo (234 REQs, 54 URs):
- Three lens buttons present, `aria-pressed` on exactly one.
- Selecting URs only gives the host class `user-request-lens is-folded`, 6 rows under the default `Active` scope, **0 REQ cards inside the lens host**, `aria-expanded="false"`, and the header carrying no `data-detail-kind` — so it is no longer the drawer trigger.
- Opening the first row: `aria-expanded="true"` and 20 cards, with every sibling row still empty. Closing it: back to 0.
- The drawer is still reachable — the row's own `button.ur-group-detail` (`aria-label="Open details for UR-042"`) opened the UR drawer titled "Consolidate completed-work presentation commands" **while the row stayed collapsed**. The two controls do not fight.
- **The filter path was exercised, because D-01 rests on it.** Typing `durations` into `#filter-search` with URs only on screen re-rendered live to 2 rows (`UR-050`, `UR-051`), still folded, still 0 cards; clearing restored 6. That is `board-filters.js` — a file this REQ never touched — driving the new reading correctly.

*Verified by work action*

## Decisions

- **D-01**: `URs only` is a card-fold on the by-UR lens, not a third `viewState.lens` value. A third value breaks three call sites that ask "is the UR lens on screen?" by naming the one value, two of which are outside this REQ's write set — the feature would have shipped with the condensed list going stale on every filter change until someone applied a seam. The fold keeps `viewState.lens === "user-request"` true for both readings and every site correct unedited, and it matches the REQ's own words ("the same group markup with the cards folded away"). Cost: two Lens buttons share `data-lens-target`, which is why `setActiveLensButton()` keys the pressed state on the pair. Reversible. DECIDE & STATE.
- **D-02**: The fold flag lives in `board-controls.js`, not `viewState` — `viewState` is declared in `web/board.js`, outside the write set, and every fragment inlines into the same IIFE. Moving it into `viewState` later is a rename. DECIDE & STATE.
- **D-03**: Expanded/collapsed lives in the DOM, not a persisted set. A re-render rebuilds rows collapsed, which is honest when the set of cards under each row has just changed. DECIDE & STATE.
- **D-04**: `.user-request-lens.is-folded` tightens the inter-group gap 18px → 8px and row padding to give 44px rows. Without it "condensed" means only "no cards", and 54 URs still runs long. DECIDE & STATE.
- **D-05**: Helpers stay inner functions of `renderUserRequestLens`, closing over the state they need — which is also what let the probe slice one long-standing anchor and get a behavioural RED instead of "anchor not found". DECIDE & STATE.
- **D-06** (orchestrator): the `generate_test.go` merge conflict was resolved by keeping both sides in order. Neither side edited the other's region; the conflict was two appends at one file's end.

## Discovered Tasks

- [low] `.ur-synthetic` (the "no input.md" chip) is unexercised by any test and unrenderable on this tree — all 54 URs have an `input.md`, so neither lens renders it live and no probe asserts it. Unchanged by this REQ, but a marker nothing currently pins.
- [low] The three "is the UR lens on screen?" call sites each spell the condition out longhand across three files, followed by the same render-and-mark-fresh pair. All correct today under D-01, but three copies of one rule; a shared predicate would make the next lens change a one-file edit.
- [low] `board-testing.js`'s copy may be dead: `applyConfirmedTestingTransition` guards on `viewState.view === "board"`, but the transition is driven from the Testing view, where `viewState.view` is `"testing"`. Worth confirming before someone maintains it.

## Review

**Overall: 97%** | 2026-08-18T11:18:40Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 96% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition):** None

**Minor findings:** 2 (report only)
- The header changes role between the two readings — drawer trigger in By UR, fold toggle in URs only. That is the right resolution of the collision the REQ flagged, and the probe pins both shapes, but it is a genuine asymmetry: a future reader editing `ur-group-head` has to know which reading they are in. The inner-function structure keeps that local, which is most of why it is minor rather than important.
- The fold marker is small and low-contrast at the row's left edge. Within the latitude the REQ granted on the affordance's shape, and the whole row is clickable, so discoverability does not depend on seeing it — noted for the maintainer's eye rather than as a defect.

**Restatement sweep:** the diff changes what `viewState.lens` means in practice — two buttons now map to one value — so the sweep asked who states that mapping. The four read sites were checked individually and all compare against `"user-request"`, which is exactly why D-01 works; none restates a two-buttons-two-values assumption. `_dev/primes/prime-kanban-board.md` describes the board's write surfaces, not its lenses, and is unaffected. No stale restatement.

**Acceptance:** Pass — all six criteria confirmed by driving the merged build in a browser, including the filter path that D-01 depends on and the drawer-versus-expand separation.

**Suggested testing:** 2 items
- The "no input.md" chip cannot be exercised on this repository; a board with a UR missing its `input.md` is the only way to see it in either lens.
- Multiple rows open at once is allowed by the design and asserted only indirectly (siblings stay empty when one opens). Worth a human pass on how the condensed list reads with several rows unfolded at 54 URs.

**Follow-ups created:** None — all three discoveries are `[low]` and would each be a one-line question; they are recorded in `## Discovered Tasks` on this archived REQ, where `do-work forensics` and a future lens change will find them, rather than as three separate `pending-answers` REQs competing for the maintainer's attention.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Trying the obvious model first and letting it fail on blast radius. A third `viewState.lens` value is the shape everyone reaches for, and it is wrong here for a reason invisible from the REQ text: three files ask "is the UR lens on screen?" by naming the one value, and two of them were outside the write set. Modelling the new reading as a *modifier* on the existing lens kept all three correct with no edit — and the check that proves it is not a test but a grep for every site that reads `viewState.lens`.

**What didn't:** Counting `.req-card` document-wide to prove "no cards until asked". It returns 25 in a correctly-folded lens, because the hidden Columns lens keeps its DOM — a number that reads as a flat failure of the headline requirement. Scoping the query to the lens host is what makes the assertion mean what it says. A DOM count is only evidence when its root is the thing under test.

**Worth knowing:** Two Lens buttons now share `data-lens-target="user-request"`, so anything keying UI state off that attribute alone will light up both. `setActiveLensButton()` exists precisely because `aria-pressed` was the first casualty of that; the next thing to key off the lens target will be the second.

## Orientation

The Board view offers a third reading: **URs only**, a condensed one-row-per-user-request list that expands in place to show that UR's REQ cards, with the detail drawer moved to its own button on the row so the two actions stop competing. Filters, the Active/All scope, and the recently-done window reach it unchanged. Lives in the queue-kanban board subsystem (`_dev/primes/prime-kanban-board.md`).

**[MAP CHANGED]** — modestly, and in one specific way: a *lens* and a *lens button* are no longer one-to-one. `viewState.lens` still holds two values while the Lens group offers three choices, with the third distinguished by a fold modifier. Anything that infers the visible reading from `viewState.lens` alone now sees By UR and URs only as the same thing, which is exactly the property that let this ship without touching `board-filters.js` or `board-testing.js` — and exactly the assumption a future lens will have to check. Staleness spot-check on `_dev/primes/prime-kanban-board.md`: every referenced path resolves and the three-write-surface count is unchanged; this REQ adds none, and the lens is read-only. The prime is not stale.
