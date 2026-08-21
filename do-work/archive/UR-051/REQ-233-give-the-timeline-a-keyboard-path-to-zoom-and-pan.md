---
id: REQ-233
title: Give the Timeline a keyboard path to zoom and pan
status: completed
completed_at: 2026-08-18T11:10:30Z
commit: 9b2578b
claimed_at: 2026-08-18T11:00:00Z
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-18T11:00:00Z
  basis:
    - Route B
    - 3-file write set
    - browser evidence
    - cross-route regression gates
status_changed_at: 2026-08-18T10:26:34Z
domain: general
created_at: 2026-08-18T01:18:57Z
user_request: UR-051
addendum_to: REQ-227
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-timeline.js
- skills/do-work-board/tools/queue-kanban/web/template.html
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Discovered Task: Give the Timeline a Keyboard Path to Zoom and Pan

## What

The Timeline view's zoom and pan are pointer-only. The three zoom buttons are reachable by keyboard, but panning the time axis is not, and both affordances are described only in a hint line of prose beside the chart rather than anywhere assistive technology reads.

## Context

Found while implementing REQ-227, which built the view. Every value the chart draws is also in the table below it, so no *data* is unreachable — what is unreachable is the navigation. A keyboard user can read every row but cannot move the window they are read in.

**Narrowed since capture.** This REQ originally recorded that rows "already take focus and open the detail drawer with Enter". That was wrong: the rows carried `role="button"` and `tabindex="0"` but the drawer opened only on a delegated `click`, which a non-native `<g>` never synthesizes from Enter or Space. PR #144's review caught it and it was fixed there — row activation now works from the keyboard and is pinned by `TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard`. What remains for this REQ is only what its title says: arrow-key panning and `+`/`-` zoom, plus stating the interaction in the panel's accessible name.

Not a regression: no other board view has zoom or pan at all, so REQ-227 added an un-keyboarded capability rather than removing a keyboarded one. That is why it is a follow-up rather than a fix inside REQ-227.

## Requirements

- With the chart focused, arrow keys pan the time axis and `+`/`-` zoom it, using the same `timelineZoomedWindow` transform the pointer path uses, so the two cannot diverge.
- The panel's accessible name states the interaction, rather than leaving it to the visual hint line beside the chart.
- Focus is visible on whatever element takes the keyboard interaction — a focus ring that exists only in the default user-agent style is not enough on a dark surface.
- No change to the pointer path's behavior.

## Red-Green Proof

**RED prompt/case:** A Node behavior probe driving the keyboard handler through the same `timelineZoomedWindow` transform, asserting that an arrow-key pan shifts the window by a bounded step and clamps at the range edges, and that `+`/`-` reach the same floor and ceiling the pointer path reaches.
**Why RED now:** there is no keyboard handler, so there is nothing to drive.
**GREEN when:** the probe passes and a headless run confirms the chart takes focus and responds to the keys.
**Validation:** Discovered during REQ-227; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-227: the new Timeline chart can be zoomed and dragged with a mouse, but there is no way to do either from the keyboard — the three zoom buttons can be tabbed to, dragging cannot be done at all, and the instructions for both sit in a line of text next to the chart rather than anywhere a screen reader announces. Nothing is unreadable: every row can be focused, opens its detail panel with Enter, and every number is repeated in the table underneath. What a keyboard user cannot do is move the time window they are reading in. Adding arrow-key panning and `+`/`-` zoom is a small, self-contained change to the one file. It is your call rather than mine because no other board view has zoom or pan at all, so this is a new capability to finish rather than a regression to repair, and you may prefer it batched with a wider accessibility pass over the board instead of done piecemeal here. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  *Answered 2026-08-18 via `do-work clarify` — user approved building the keyboard path now rather than batching it into a wider accessibility pass.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome was clear and the file was named, but the change had to reuse an existing listener registry and an existing window transform rather than adding either, so the surrounding code had to be read before anything was written.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modify) — pan transform, keyboard decision function, handler extension
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify) — accessible name, `tabindex`, hint line
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — Node behaviour probe and template-attribute test
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify, **integration seam applied by the orchestrator**) — focus ring on the chart

**Files I will NOT touch:** `timeline.go` (no payload change — the keyboard drives the client-side window model only), `board-durations.js` / `durations.go` (adjacent view, REQ-231's territory).

**Acceptance criteria (restated from REQ):**
- [ ] Arrow keys pan and `+`/`-` zoom, both through the existing `timelineZoomedWindow` transform, so the pointer and keyboard paths cannot diverge
- [ ] The panel's accessible name states the interaction
- [ ] Focus is visible on the element taking the interaction — not the user-agent default
- [ ] No change to the pointer path's behavior

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified — integration seam, applied by the orchestrator inside the merge commit)

**What was done:** Two pure functions were added to `board-timeline.js`. `timelinePannedWindow` slides the window without resizing it and clamps it to the same bounds the drag path clamps to, stepping by `TIMELINE_PAN_FRACTION` (0.15) of the *visible span* rather than a fixed duration. `timelineKeyboardWindow` is the whole keyboard decision as one pure function: Left/Right pan, `+`/`=` and `-`/`_` route through the existing `timelineZoomedWindow` at a centred 0.5 anchor, and every other key returns `null`. Routing zoom through the pointer path's own transform is what makes the two unable to acquire different floors, ceilings, or clamps.

No new listener was added. The chart already had exactly one `keydown` listener on the scroll host, registered through `addTimelineListener` and therefore already in the teardown registry; it now asks row activation first and falls through to the window transform. The handler restores focus to the same row after the render, because `renderVisibleRows` rebuilds every row node and a focused row would otherwise be a dead element by the next keypress — one arrow press, then dead keys.

`template.html` gained `tabindex="0"` on `#timeline-scroll` (without it the whole path is unreachable), an accessible name that states the interaction, and a hint line that names the keys for sighted keyboard users. `board.css` gained a `:focus-visible` ring on the chart using `--accent-claimed`, the token every other focus ring on the board already uses, at `outline-offset: -2px` so the ring is not clipped by the scroll container's own edge.

## Qualification

Passed — 4 files verified in the merge range `4bc04e6..9b2578b`, 4 acceptance criteria traced. No P-A-U section on this REQ.

Judgment checks: `timelinePannedWindow`'s clamp was read rather than assumed — it bounds the next start to `[boundStart, max(boundEnd − span, boundStart)]`, which is correct in both directions and degrades sanely when the span exceeds the bound span. `timelineKeyboardWindow` routes both zoom directions through `timelineZoomedWindow`, so the requirement that the two paths cannot diverge is structural rather than asserted. Nothing hollow: every branch returns a real window or `null`, and `null` is what leaves Enter/Space to row activation and Up/Down to native scrolling.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Exit 0, run unpiped — includes `go vet`, the uncached queue-kanban suite, and the strict JavaScript behavior lane

**Red-green validation:**
- `TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer`: ✗ → ✓. The first RED was a reference-error class (`web/board-timeline.js declares no numeric constant TIMELINE_PAN_FRACTION`), which proves nothing on its own, so a second RED was produced with the handler present but the clamp removed: `panning right settled with the window ending at 9720000000 ms, want the range edge 2592000000 ms` — 40 held presses walking the window off the data entirely. That is the assertion failing for the reason the test exists.
- `TestTimelinePanelStatesItsKeyboardInteraction`: ✗ `the timeline panel's accessible name does not state which keys pan` → ✓.

**New tests added:**
- `TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer` — Node probe asserting the pan step is a bounded fraction of the visible span, that held presses settle exactly on each range edge with the span intact, that `+`/`-` reach the same floor and ceiling the *pointer* path reaches (driven through `timelineZoomedWindow` at the wheel's off-centre 0.25 anchor, not the keyboard's 0.5), and that Enter, Space, Tab, Up, Down and an ordinary letter all move nothing. Constants are read out of the shipped `board-timeline.js`, so the probe cannot pass against numbers the view does not use.
- `TestTimelinePanelStatesItsKeyboardInteraction` — pins the panel's accessible name and the chart's `tabindex`. Both are single attributes, silently droppable in a template edit, and both load-bearing: without `tabindex` the entire keyboard path is unreachable.

**Render evidence — driven in a real browser against the merged tree, not the builder's branch.** Generated a static board from the live repo (233 REQs, 54 URs), served it on loopback, and drove it:
- The chart takes focus and reports `tabIndex: 0`.
- Two `+` presses narrowed the axis from `28 May … 20 Aug` to `22 Jun … 25 Jul`, centred.
- One ArrowRight moved it to `27 Jun … 30 Jul` — about five days on a ~33-day window, i.e. the 15% fraction.
- 60 further ArrowRight presses pinned the right edge at `20 Aug`, the range end; 200 ArrowLeft presses pinned the left edge at `28 May`, the range start. Both clamped rather than running off.
- 60 `-` presses restored the fit-all window **byte-identically** to the axis string before any zoom.
- An ArrowDown keydown came back with `defaultPrevented === false`, so the container keeps its native vertical scrolling.
- **Focus ring, checked with a real Tab press rather than a programmatic `.focus()`** — the latter does not trigger `:focus-visible` and would have produced a false negative. Tabbing from `Fit all` lands on `#timeline-scroll` with `:focus-visible` matching and a computed `outline: solid 2px rgb(58, 107, 196)` at `-2px` offset. Zoomed into the chart's top-left corner to confirm the ring is drawn unclipped on all edges.

*Verified by work action*

## Decisions

- **D-01**: Left/Right pan the time axis; Up/Down are deliberately not claimed. The time axis is horizontal, and the scroll container's Up/Down already scroll a 233-row queue — taking that away to duplicate what Left/Right provide would be a net loss. `preventDefault()` is called only on keys the view owns. DECIDE & STATE.
- **D-02**: `=` and `_` zoom alongside `+` and `-`, because `+` requires Shift on a US layout and browsers themselves accept the unshifted face. One clause. DECIDE & STATE.
- **D-03**: The pan step is 15% of the *visible span*, not a fixed duration — a fixed step is imperceptible zoomed out and a jump zoomed in. The probe asserts the fraction, reading `TIMELINE_PAN_FRACTION` from the shipped file rather than restating it. DECIDE & STATE.
- **D-04**: Keyboard zoom anchors at 0.5. There is no pointer to anchor to, and this is what the existing zoom buttons already do, so button and keyboard zoom stay identical. DECIDE & STATE.
- **D-05**: Focus is restored to the same row after a keyboard-driven render. `renderVisibleRows` rebuilds every row node, so a focused row is a dead element by the next keypress and focus falls to `<body>` — a keyboard user would get one arrow press and then nothing. Scoped to the keyboard handler; the wheel and drag paths imply a pointer, where focus is not on a row. DECIDE & STATE.
- **D-06**: The `board.css` focus ring was handed back as an integration seam rather than written by the builder, because a sibling builder held that file. Applied by the orchestrator inside the merge commit, at `outline-offset: -2px` rather than the `+2px` the other rings use — `.timeline-scroll` is an `overflow-y: auto` container flush under the axis, so a positive offset draws the ring where the surrounding card clips it. Token is `--accent-claimed`, the same one every other focus ring on the board uses. DECIDE & STATE.

## Discovered Tasks

- [normal] **`.timeline-row { outline: none; }` suppresses the row focus ring** (`web/board.css`). Rows fall back to a one-step background tint on `:focus`, which is a much weaker signal than the 2px accent ring every other focusable thing on the board gets, and it does not distinguish `:focus` from `:focus-visible`. Same requirement-3 concern as this REQ, applied to the element next door.
- [low] The Timeline's rows `<svg>` carries `role="img"` with a long chart-level `aria-label`, while the rows inside carry their own labels. A screen reader treating the SVG as an atomic image may never reach the row labels. Wants a deliberate decision about which layer owns the description.
- [low] A filter that matches nothing returns before re-fitting, so clearing it leaves the previous zoom rather than refitting. Probably correct, possibly surprising.

## Review

**Overall: 96%** | 2026-08-18T11:09:44Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 96% |
| Scope | 92% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition):** None

**Minor findings:** 2 (report only)
- Scope: the builder's declared write set could not satisfy requirement 3, because the focus ring lives in a file a sibling builder held. It stopped and handed back the exact rule rather than writing it — the correct behaviour, and the reason the score is 92 rather than 100 is the planning miss, not the conduct. A `board.css` entry belonged in the write set from the start; the REQ's own write set omitted it too.
- The row-focus suppression above is arguably in the spirit of requirement 3 ("focus is visible on whatever element takes the keyboard interaction"). It is out of the letter, since rows are not what this REQ adds, so it is routed to REQ-239 rather than folded in.

**Restatement sweep:** the diff adds a keyboard contract that is stated in three places — the panel's `aria-label`, the visual hint line, and the probe. All three were written together and agree. The sweep also asked whether any doc describes the Timeline as pointer-only: `_dev/primes/prime-kanban-board.md` does not, and the archived REQ-227/REQ-228 records are dated history. `CHANGELOG.md`'s Timeline entries describe what shipped then and stay as they are. No stale restatement.

**Acceptance:** Pass — every requirement confirmed by driving the merged build in a real browser, including the clamp at both range edges and the focus ring under a genuine Tab press.

**Suggested testing:** 2 items
- Screen-reader pass. Nothing here proves the accessible name is *announced usefully*, only that it is present and states the interaction; that is a human judgment.
- Non-US keyboard layouts. `+`/`=`/`-`/`_` are the US faces; a layout where those characters sit elsewhere was not tested.

**Follow-ups created:** REQ-239; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Making the keyboard path a *pure function* that returns a window or `null`, and routing its zoom through the transform the pointer already used. That is what turned "the two paths must not diverge" from a promise into a structural fact — there is no second floor, ceiling, or clamp to drift. It also made the probe able to assert against the pointer path's own anchor rather than a copy of it.

**What didn't:**
- Checking the focus ring with a programmatic `.focus()`. It reported `matchesFocusVisible: false` and a computed `outline: none` — which reads exactly like a missing rule, and is actually Chrome correctly refusing `:focus-visible` for non-keyboard focus. The rule was fine. Only a real Tab keypress answers this question; a scripted focus call will tell you the ring is broken when it is not.
- The write set. Requirement 3 was a CSS requirement from the moment it was written, and neither the REQ's captured write set nor the builder's scope named a stylesheet. The builder handled it correctly at the boundary, but the boundary should not have been there.

**Worth knowing:** `renderVisibleRows` rebuilds every row node, so anything holding a reference to a row across a render is holding a dead element. This bit the keyboard path (focus fell to `<body>` after one press) and will bit anything else that keeps row state — the fix is to capture the row's `data-detail-id` before the render and re-query after it.

## Orientation

The board's Timeline can now be navigated from the keyboard: Tab reaches the chart, ← and → pan the time axis by 15% of what is on screen, and `+`/`-` (or `=`/`_`) zoom it, all clamped to the same bounds the drag and wheel paths clamp to. The panel's accessible name states this, and the chart carries a visible focus ring. Lives in the queue-kanban board subsystem (`_dev/primes/prime-kanban-board.md`).

**[MAP CHANGED]** — not because of the keys, but because the Timeline now has a *keyboard contract* that the pointer path and the keyboard path both derive from one transform. Anything that later adds a third way to move the window — REQ-235's period stepping is next, and is already written against this constraint — must go through `timelineZoomedWindow` too, or the guarantee this REQ established stops holding. Staleness spot-check on `_dev/primes/prime-kanban-board.md`: every referenced path resolves and the three-write-surface count is unchanged (this REQ adds none). The prime is not stale.
