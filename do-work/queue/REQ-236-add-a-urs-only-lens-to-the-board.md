---
id: REQ-236
title: Add a URs-only lens to the Board view
status: pending
created_at: 2026-08-18T10:30:00Z
user_request: UR-053
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
