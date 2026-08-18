# REQ-236 hand-back — Add a URs-only lens to the Board view

**Branch:** `worktree-agent-REQ-236-urs-only-lens`
**Commit:** `266ec09` — `[REQ-236] Add a URs-only lens to the Board view`

## Integration seams

**None.** Nothing outside the write set needs a line. That was the main design
constraint and it drove the central decision (D1 below): `viewState.lens` still
takes exactly `flat` and `user-request`, so the three call sites in files I do
not own — `board-filters.js` `onFiltersChanged`, `board-testing.js`
`applyConfirmedTestingTransition`, plus `board-cards.js`
`applyRecentWindowSelection` — keep deciding correctly for the new reading with
no edit. I verified the filter path live in a browser (see Rendered, below):
typing in the search box while URs only is on screen re-renders the condensed
list, which is exactly what a third `viewState.lens` value would have broken
until someone patched two other files.

## File manifest

| Action | File | What |
|---|---|---|
| Modified | `skills/do-work-board/tools/queue-kanban/web/template.html` | Third Lens button, `URs only`, carrying `data-lens-target="user-request"` plus `data-ur-cards="folded"`. |
| Modified | `skills/do-work-board/tools/queue-kanban/web/board-controls.js` | New `userRequestCardsFolded` fold flag, `applyLensSelection()` (the one entry point a lens click goes through), `setActiveLensButton()` (keys the pressed button on lens+fold, since two buttons share a lens), and one line in `applyLens()` toggling `is-folded` on the lens host. |
| Modified | `skills/do-work-board/tools/queue-kanban/web/board-cards.js` | `renderUserRequestLens()` learns the fold: same `ur-group` markup, cards built by a shared inner `makeUserRequestCards()`, built late and dropped again on collapse in the folded reading. |
| Modified | `skills/do-work-board/tools/queue-kanban/web/board.css` | `.ur-group-row`, `.ur-group-detail`, `.ur-fold-marker`, and the tighter `.user-request-lens.is-folded` metrics. |
| Modified | `skills/do-work-board/tools/queue-kanban/generate_test.go` | Two new tests appended at the end (markup + Node behaviour probe), and one line added to the pre-existing `TestJavaScriptBehaviorByUserRequestLensUsesRecentWindowAtCaller` probe preamble (see Cross-REQ break). |

No file was created or deleted. `VERSION`, `skills/do-work/VERSION`,
`actions/version.md`, `CHANGELOG.md` untouched.

## Red-Green Proof

### RED (verbatim)

```
$ go test -count=1 -run 'TestGenerateOffersThreeLensButtons|TestJavaScriptBehaviorUserRequestsOnlyLensFoldsCardsUntilARowIsOpened' .
--- FAIL: TestGenerateOffersThreeLensButtons (0.22s)
    generate_test.go:2321: lens group is missing "data-ur-cards=\"folded\"" in the generated page
--- FAIL: TestJavaScriptBehaviorUserRequestsOnlyLensFoldsCardsUntilARowIsOpened (0.26s)
    generate_test.go:2494: URs only row UR-401 rendered cards []string{"REQ-601", "REQ-602"} before it was opened, want none
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.812s
FAIL
EXIT=1
```

The behaviour failure is the assertion, not a missing anchor: the probe declares
`userRequestCardsFolded = true` itself and drives the *existing*
`renderUserRequestLens`, which ignored it and appended every card. That is the
REQ's stated reason for red. I deliberately kept every new helper as an inner
function of `renderUserRequestLens` so the probe slices one anchor that already
existed and gets a real behavioural red rather than "anchor not found".

### GREEN (verbatim)

```
$ go test -count=1 -run 'TestGenerateOffersThreeLensButtons|TestJavaScriptBehaviorUserRequestsOnlyLensFoldsCardsUntilARowIsOpened' -v .
=== RUN   TestGenerateOffersThreeLensButtons
--- PASS: TestGenerateOffersThreeLensButtons (0.23s)
=== RUN   TestJavaScriptBehaviorUserRequestsOnlyLensFoldsCardsUntilARowIsOpened
--- PASS: TestJavaScriptBehaviorUserRequestsOnlyLensFoldsCardsUntilARowIsOpened (0.29s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	0.728s
EXIT=0
```

What the probe pins, over a three-UR fixture (`UR-401` with a pending and a
recently-completed REQ, `UR-402` with a claimed REQ, `UR-403` with only an
aged-out completed REQ):

- folded: three rows, **zero** REQ cards, every head `aria-expanded="false"`, and
  exactly one drawer trigger per row pointing at that row's own UR;
- activating row one: `aria-expanded="true"` and exactly `REQ-601, REQ-602`,
  with the sibling rows still card-free;
- activating it again: `aria-expanded="false"` and no cards;
- `Active` scope + `status: pending`: URs only and By UR render the identical
  list (`UR-401` alone — `UR-402` filtered out, `UR-403` scope-hidden), the
  opened row shows only the filter-matching `REQ-601`, and the By UR head
  carries no `aria-expanded` at all, so the fold cannot leak into the
  always-open reading.

### Cross-REQ test break (general.md § Cross-REQ Test-Break Rules)

`TestJavaScriptBehaviorByUserRequestLensUsesRecentWindowAtCaller` (REQ-122's
lock-in) failed with `ReferenceError: userRequestCardsFolded is not defined`
once the renderer started consulting the flag. The behaviour change is
intentional, so its probe preamble gained one line —
`var userRequestCardsFolded = false;` with a comment saying it pins the
always-open reading. Its assertions are unchanged and still pass.

## maintainer-verify.sh

```
Maintainer verification passed.
EXIT=0
```

Run unpiped from the worktree root, twice: once after the implementation, once
after the final CSS tweak. Exit code 0 both times.

## What I rendered and confirmed by looking at it

Built `go build -o /tmp/qk-236 .` and generated a static board against this
worktree — **232 REQs, 54 URs, 226 calendar entries** — then served it on
loopback and drove it in a real browser (Playwright). Confirmed by DOM query and
by looking at two screenshots:

- **Three lens buttons** in `#lens-group`: `Columns` (`flat`), `By UR`
  (`user-request`), `URs only` (`user-request` + `data-ur-cards="folded"`), with
  `aria-pressed` landing on exactly one of them.
- **Condensed list**: clicking `URs only` hid the columns, showed the lens host
  with class `user-request-lens is-folded`, and rendered **54 rows and 0
  `.req-card` nodes** (`URs: All`). Rows measure 44px each — id chip, title,
  `N REQ`, `Details`.
- **Expand in place**: focused a row's header and pressed **Enter** — it grew to
  8 REQ cards, `aria-expanded="true"`, the fold marker rotated, no other row
  gained cards, and the detail drawer stayed **closed**.
- **Drawer still reachable**: **Tab** from the expanded header landed on
  `.ur-group-detail` (`aria-label="Open details for UR-003"`); **Enter** opened
  the UR drawer showing `UR / UR-003` while the row stayed expanded with its 8
  cards. Both actions are keyboard-operable and neither disturbs the other.
- **Filters reach it**: typing `durations` into the search box while URs only
  was on screen re-rendered live to 2 rows (`UR-050`, `UR-051`), still folded,
  still 0 cards. This is the check that would have failed under a third
  `viewState.lens` value.
- Console: one 404 for `/favicon.ico` from my throwaway static server. No
  JavaScript errors.

Temporary artifacts (`/tmp/qk-236`, `/tmp/board-236`, the two screenshots
Playwright wrote into the main tree root) were deleted; `timeline-focus-ring.png`
in the main tree root belongs to the sibling builder and was left alone.

## How I resolved the expand-vs-drawer collision on `ur-group-head`

In the folded reading only, the header stops being the drawer trigger and
becomes the fold toggle: it drops `data-detail-kind` / `data-detail-id`, gains
`aria-expanded` and a `▸` marker, and gets a direct click listener. The drawer
moves to a sibling `button.ur-group-detail` inside a new `div.ur-group-row`
flex wrapper, carrying `data-detail-kind="ur"` so the existing delegated
`[data-detail-kind]` handler picks it up unchanged.

Why this shape rather than a disclosure triangle beside a drawer-opening header:

- The REQ says clicking the **row** expands it. A big click target that opens a
  drawer, with a 12px triangle for the primary action, inverts that.
- The delegated drawer handler matches on `closest("[data-detail-kind]")`. As
  long as the header carries that attribute, *any* click inside it — triangle
  included — also opens the drawer. Nested buttons are invalid HTML, so a
  triangle inside the header cannot cleanly opt out. Removing the attribute
  from the header is the only way the two stop fighting.
- Both controls are real `<button>`s in DOM order (fold, then Details), so
  keyboard operation and tab order come for free, and `aria-expanded` sits on
  the element that owns the state.

The always-open By UR reading is byte-for-byte unchanged: same single
`button.ur-group-head` with `data-detail-kind="ur"`, no wrapper, no
`aria-expanded`. The probe asserts that last point explicitly.

## Decisions

**D1 — URs only is a card-fold on the by-UR lens, not a third `viewState.lens`
value.** A third value (`"user-request-condensed"`) is the obvious shape, and I
started there. It breaks three call sites that ask "is the UR lens on screen?"
by naming the one value: `board-filters.js` `onFiltersChanged`,
`board-testing.js` `applyConfirmedTestingTransition`, and `board-cards.js`
`applyRecentWindowSelection`. Two of those are outside my write set, so the
feature would have shipped with the condensed list going stale on every filter
change until someone applied a seam. Modelling it as a fold keeps
`viewState.lens === "user-request"` true for both readings and every one of
those sites correct, unedited — and it is the more honest model anyway: the REQ
itself describes this as "the same group markup with the cards folded away".
The visible cost is that two Lens buttons share `data-lens-target`, which is why
`setActiveLensButton()` exists to key the pressed state on the pair. Reversible:
splitting it into a third lens value later is a mechanical change plus the two
seams.

**D2 — the fold flag lives in `board-controls.js`, not in `viewState`.**
`viewState` is declared in `web/board.js`, outside my write set. A module-level
`var userRequestCardsFolded` in the file that owns the lens controls is one
declaration next to the code that reads and writes it, and every fragment is
inlined into the same IIFE so `board-cards.js` sees it. If a later REQ moves it
into `viewState` alongside `lens` and `windowHours`, that is a rename.

**D3 — expanded/collapsed state lives in the DOM, not in a persisted set.** The
REQ calls it view state. A re-render (filter change, window change, scope
toggle) rebuilds the rows collapsed, which is the honest outcome when the set of
cards under each row has just changed. Keeping a set of open ids across
re-renders is speculative until someone asks for it.

**D4 — the condensed list tightens up.** `.user-request-lens.is-folded` cuts the
inter-group gap 18px → 8px and the row padding 14px/18px → 9px/14px, giving 44px
rows. Without this "condensed" is only condensed by the absence of cards, and on
54 URs the list still runs long. One line in `applyLens()` toggles the class.

**D5 — helpers stay inner functions of `renderUserRequestLens`.** They close
over `shownRequestIds`, `group`, and `head`, which is exactly the state they
need, and it keeps the Node probe slicing one long-standing anchor rather than
new ones — which is what made the RED behavioural instead of "anchor not found".

## Discovered Tasks

- **`.ur-synthetic` ("no input.md") is unexercised by any test and unrenderable
  on this tree.** Every one of this repo's 54 URs has an `input.md`, so neither
  lens renders the chip live, and no probe asserts it. It sits on the shared
  path and is unchanged by this REQ, but it is a marker nothing currently pins.
- **The three "is the UR lens on screen?" call sites each spell the condition
  out longhand** (`viewState.view === "board" && viewState.lens ===
  "user-request"`, in `board-filters.js`, `board-testing.js`, and
  `board-cards.js`, each followed by the same render-and-mark-fresh pair). D1
  means they are all correct today, but they are three copies of one rule across
  three files. A shared predicate would make the next lens change a one-file
  edit.
- **`board-testing.js`'s copy may be dead.** `applyConfirmedTestingTransition`
  guards on `viewState.view === "board"`, but the transition is driven from the
  Testing view's own controls, where `viewState.view` is `"testing"`. I did not
  chase whether any path reaches it from the board view; worth a look before
  someone maintains it.
