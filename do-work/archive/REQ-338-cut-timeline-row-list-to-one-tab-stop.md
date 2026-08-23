---
id: REQ-338
title: "Cut the Timeline row list to one Tab stop"
status: completed
claimed_at: 2026-08-23T20:32:41Z
completed_at: 2026-08-23T20:47:14Z
commit: cac6718
kb_status: pending
created_at: 2026-08-23T18:30:26Z
user_request: UR-067
addendum_to: REQ-333
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-23T20:32:41Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - browser evidence
    - cross-route regression gates
    - full-suite verification
route: B
related: [REQ-336, REQ-337]
batch: timeline-click-regression
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
---

# Cut the Timeline Row List to One Tab Stop

## What

Every Timeline REQ row is its own Tab stop, so tabbing past the row list takes one press per row
(29 on the observed board). Make the row list a single Tab stop with arrow-key movement between
rows (roving tabindex), so Tab escapes the list in one press.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md` (REQ-233, REQ-235, REQ-239, REQ-324
  lesson links), `general.md`, `coding-guardrails.md`, `communication-style.md`, `testing.md`.
  1. Write the probe first (`tdd: true`) and confirm it fails → verify: 35 rows tabbable, no focus
     movement on ArrowDown.
  2. Store the roving stop in `timelineViewState` and render `tabindex` from it, clamped into the
     rendered range so the virtualization cannot take the list out of the Tab order → verify:
     exactly one tabbable row, all others `-1`.
  3. Take **Up/Down** for row movement and leave Left/Right panning untouched — the vertical axis
     matches the row list and REQ-333's arrow-pan contract is horizontal → verify: ArrowRight still
     pans and does not move row focus.
  4. Sync the roving index from `focusin` so an arrow press moves from where focus actually is →
     verify: the tabbable row follows a directly-focused row.
  5. Measure the real Tab count over CDP, since Tab's movement is a trusted-input default action
     the probe lane cannot dispatch → verify: RED many presses, GREEN one.
- [x] **[APPLY]:** Both declared files, nothing else. `timelineTabbableRowIndex`,
  `scrollTimelineRowIntoView`, `moveTimelineRowFocus`, a `focusin` listener, the rendered
  `tabindex`, and `rovingRowIndex` on the view state; plus the probe. One neighbouring comment
  corrected because the change made it wrong (D-03).
- [x] **[UNIFY]:** `git diff --stat` → 2 files, both declared.
  - `web/board-timeline.js` — read every hunk: one state field, one pure helper, two scoped
    functions, one listener, one rendered attribute, one corrected comment. `node --check` clean.
  - `timeline_browser_probe_test.go` — read the whole hunk: one probe with a vacuity guard and five
    assertion groups. `gofmt -l` prints nothing.
  - No debug artifacts: no `console.log`, `debugger`, `alert`, `TODO` in the diff. The CDP Tab probe
    and the temporary RED revert live in the scratchpad; `diff -q` confirms `board-timeline.js` is
    byte-identical to its pre-revert copy.
  - Native lint/tests: full board suite with the browser lane enabled → ok; canonical gate → exit 0.

## Why

Under REQ-333 the user refuted this as a keyboard trap — Tab does escape — and left it. During
this capture the user chose to fix the tedium anyway: a keyboard user pays one press per row to
get past the chart. Not urgent, but noticed in real use.

## Context

REQ-333 (commit 36c4518) established the current keyboard contract: rows are focusable, Enter
and Space on a focused row open the drawer, and arrow keys pan the chart without killing the
keyboard path. Its F3 finding recorded the 29-press escape as "not a trap, leave it" — this REQ
supersedes that "leave it" by the user's explicit choice at capture.

## Detailed Requirements

- The row list occupies one Tab stop: Tab from any focused row leaves the list in one press
  (in each direction).
- Arrow keys move focus between rows (roving tabindex).
- Enter and Space on the focused row still open the drawer — do not regress the drawer lane
  REQ-336 restores.
- Do not break REQ-333's contracts: arrow-key panning with a row focused must keep working
  (the builder decides how row-movement keys and pan keys share the keyboard — e.g. which axis
  moves rows vs pans), and focus must survive the node rebuilds REQ-333's focus restore handles.

## Builder Guidance

Certainty: Firm on the shape (one stop, roving tabindex — user picked this option at capture);
builder's latitude on key assignment and on how the roving index interacts with the existing
focus-restore and pan handlers in `web/board-timeline.js`.

## Red-Green Proof

**RED prompt/case:** In the probe lane, generate a board with multiple REQ rows, focus the first
Timeline row, and press Tab: focus lands on the next row, not past the list — escaping takes as
many presses as there are rows (29 observed).
**Why RED now:** Every row carries its own tabindex, so each is a distinct Tab stop.
**GREEN when:** A probe-lane test asserts: Tab from a focused row exits the row list in one
press, arrow keys move focus between rows, and Enter/Space on the focused row still opens the
drawer. Test written first and confirmed failing before the change (`tdd: true`).
**Validation:** User confirmed (chose the roving-tabindex option during capture)

## Dependencies

None. Overlaps REQ-336 on `web/board-timeline.js` (declared in both write sets) — serial
execution or coordination needed if dispatched in parallel.

---
*Source: UR-067 — see `do-work/user-requests/UR-067/input.md` for complete verbatim input.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The shape is settled (one Tab stop, roving tabindex, user-chosen) so nothing needs planning. What needs discovery is how a roving index shares the keyboard with REQ-333's arrow-key panning and survives REQ-333's focus restore across node rebuilds — three interacting handlers in one large module.

**Planning:** Not required

## Exploration

**Where the Tab stops come from.** One place: the row group's markup at
`web/board-timeline.js:1897`, `tabindex: "0"` on every rendered row. Rows are virtualized through
`timelineVisibleRowRange` (`:885`), so the count is the rendered window's size (33-35 on this
board), not the 332 rows in the payload — the user's 29 was the same measurement on a smaller
window.

**The three handlers a roving index has to share the keyboard with.**

- `timelineKeyboardWindow` (`:306`) owns Left/Right (pan) and `+`/`-` (zoom) and returns null for
  everything else. Its comment said that null "leaves Enter and Space to row activation and Up/Down
  to scrolling the queue" — the vertical keys were free, which is what makes them the right axis
  for row movement and leaves REQ-333's arrow-pan contract untouched.
- `timelineKeyboardActivationTarget` (`:1192`) owns Enter/Space. Asked first in the listener
  (`:2450`), so it keeps precedence.
- `renderVisibleRows` (`:1834`) owns focus restore across rebuilds, keyed on the focused row's
  `data-detail-id`. Its own comment records why the keydown handler must NOT restore focus itself:
  the scroll event a scroll write schedules is asynchronous, so a restore done in the handler is
  wiped by the rebuild that event triggers. That is the ordering constraint the roving move has to
  respect.

**What virtualization does to a roving index.** The stored roving row is frequently not rendered.
Rendering `tabindex="0"` only for an exact match would then mark nothing tabbable and take the whole
list out of the Tab order — strictly worse than the 29 stops. Clamping the roving index into the
rendered range for display keeps exactly one stop, keeps it reachable, and leaves the stored index
alone so scrolling back restores the reader's place.

**RED, measured two ways before any change:**

- Probe lane: `35 rows with tabindex=0 out of 35`, `0 carry tabindex=-1`, ArrowDown leaves focus on
  the same row. (Assertions (d) and (e) — arrow panning and Enter — passed already, so the probe
  starts as a failing test of the new behaviour and a passing test of the old contracts.)
- Real input over CDP, which is the only way to measure Tab at all (its focus movement is a default
  action on a trusted event, so a synthetic `KeyboardEvent` moves nothing): `33 rendered rows, 33
  tabbable, 6 Tab presses and still inside the row list` — the probe caps at six.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modify) — roving tabindex, arrow
  movement, focus sync.
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modify) — the probe that
  pins it.

**Files I will NOT touch:** `web/board-controls.js`, `web/board-detail.js`, `board.css` (no visual
change — the existing focus ring from REQ-239 is what a roving stop needs and it already applies to
`.timeline-row:focus-visible`).

**Acceptance criteria (restated from REQ):**
- [ ] The row list occupies one Tab stop; Tab leaves the list in one press, in each direction
- [ ] Arrow keys move focus between rows (roving tabindex)
- [ ] Enter and Space on the focused row still open the drawer
- [ ] REQ-333's contracts hold: arrow-key panning with a row focused still works, and focus survives
      the node rebuilds the focus restore handles

## Decisions

- **D-01 — DECIDE & STATE. Up/Down rove the rows; Left/Right keep panning.** The REQ left the key
  assignment to the builder. Vertical keys match a vertical list, and REQ-333's arrow-pan contract is
  horizontal, so this needs no trade at all — the two never contend. Up/Down are intercepted **only
  when focus is on a row**: with focus on the chart itself they stay the browser's own scrolling,
  which is the one way left to move through the queue without entering the list.
- **D-02 — DECIDE & STATE. The rendered tab stop is the roving index clamped into the rendered
  range.** With virtualization the roving row is usually not rendered, and marking nothing tabbable
  would take the list out of the Tab order entirely — worse than the defect. The clamp is display
  only: `timelineViewState.rovingRowIndex` still names the row the reader was on, so scrolling back
  restores their place rather than resetting to the top.
- **D-03 — DECIDE & STATE. One neighbouring comment was corrected.** `timelineKeyboardWindow`'s
  header said returning null "leaves … Up/Down to scrolling the queue". That is now only true when
  focus is not on a row, so the sentence became a half-truth about the very change this REQ makes —
  the board prime's REQ-231 shape. Rewritten in place; it is one comment in a declared file, not a
  follow-up.
- **D-04 — DECIDE & STATE. A `focusin` listener syncs the roving index, and it writes attributes
  rather than re-rendering.** Without the sync, arrowing after clicking row 20 would move from row 0,
  because the index and the actual focus had diverged. `focusin` fires on every Tab and click into the
  list, so a render per focus change would sweep every row for what is two attribute writes.
- **D-05 — DECIDE & STATE. No Home/End, no PageUp/PageDown, no wrap-around at the ends.** The REQ
  asks for one Tab stop and arrow movement; the rest is the ARIA pattern's optional half and nobody
  asked for it. `moveTimelineRowFocus` clamps at both ends and returns without a render when the
  index would not change.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified)

**What was done:** The Timeline row list is now a single Tab stop with a roving tabindex.
`timelineViewState` carries a `rovingRowIndex`; the render marks exactly that row `tabindex="0"` and
every other rendered row `tabindex="-1"`, with the index clamped into the rendered range so
virtualization cannot leave the list unreachable. ArrowDown/ArrowUp on a focused row move the stop —
scrolling the target into view, rebuilding synchronously, then focusing, in that order because the
scroll event's own rebuild would otherwise wipe an earlier focus — and a `focusin` listener keeps the
index aligned with however focus actually arrived. Left/Right panning, Enter/Space activation, and
the rebuild focus restore are untouched. Added
`TestBrowserBehaviorTimelineRowListIsOneTabStop`, which asserts the one-stop property, the stop
following focus, arrow movement in both directions, and both REQ-333 contracts.

## Qualification

Passed — 2 files verified, 4 requirements traced, P-A-U confirmed.

- Mechanical: `tools/checks/qualify.sh` exit 0; `tools/checks/scope-drift.sh` exit 0.
- Substantive (check 2): +298/-8; read every hunk. The JS side is one state field, one pure helper,
  two scoped functions, one listener, one rendered attribute and one corrected comment; the test side
  is one probe with a vacuity guard.
- Requirements traced (check 3): one Tab stop → the rendered attribute plus the probe's (a), and
  measured for real at one press; arrow movement → `moveTimelineRowFocus` plus (c); Enter/Space →
  (e), untouched code; REQ-333's contracts → (d) plus the full suite, including REQ-333's own
  keyboard probe.
- Data flows (check 6): nothing stubbed or hardcoded. `timelineTabbableRowIndex` is a pure function
  of the live roving index and the live rendered range; `moveTimelineRowFocus` reads `rows.length`
  rather than a remembered count, so a filter change cannot leave it pointing past the list.
- Contamination check: the previous REQ (REQ-337) touched only `timeline_browser_probe_test.go`, and
  this REQ appends a new probe rather than editing REQ-337's test. `web/board-timeline.js` was last
  touched by REQ-336, whose `capturePanPointer` hunk is untouched here — confirmed by reading the
  diff, and by REQ-337's check still passing.

## Testing

**Tests run:**
- `go test -count=1 ./...` in `skills/do-work-board/tools/queue-kanban` with the browser lane enabled
- `bash _dev/tests/maintainer-verify.sh` from the repo root against the final tree
- A CDP trusted-key probe counting real Tab presses, before and after

**Result:** ✓ board suite `ok … 119.742s`; canonical gate exit 0.

**Red-green validation (test written first, per `tdd: true`):**

| assertion | before | after |
|---|---|---|
| rows with `tabindex="0"` | ✗ 35 of 35 | ✓ 1 of 35 |
| rows with `tabindex="-1"` | ✗ 0 | ✓ 34 |
| tab stop follows a directly-focused row | ✗ none tabbable to name | ✓ the focused row |
| ArrowDown moves focus to the next row | ✗ stays on the same row | ✓ moves, stop moves with it |
| ArrowUp moves it back | ✗ stays | ✓ moves back |
| ArrowRight still pans, without moving row focus | ✓ already held | ✓ unchanged |
| Enter opens the drawer on the focused row | ✓ already held | ✓ unchanged |

The last two rows are the REQ-333/REQ-336 contracts: they passed before the change and after it, so
the probe is a regression guard for them as well as a lock-in for the new behaviour.

**Measured with real keyboard input** (CDP `Input.dispatchKeyEvent`, Chromium 1194, viewport
1400×900, because Tab's focus movement is a default action the engine performs only for trusted
events — a synthetic `KeyboardEvent` moves nothing, which is why the probe above asserts the property
instead):

```
RED    renderedRows=33 tabbableRows=33 -> 6 Tab presses and still inside the row list (probe caps at 6)
GREEN  renderedRows=33 tabbableRows=1  -> 1 Tab press, then outside the row list (on a BUTTON)
```

**New tests added:**
- `TestBrowserBehaviorTimelineRowListIsOneTabStop`

**Existing tests updated (cross-REQ impact):** none. REQ-333's keyboard probe, REQ-336's fix and
REQ-337's check all pass unchanged, which is the evidence that the roving index did not take a key
or a focus path they rely on.

*Verified by work action*

## Review

**Overall: 96%** | 2026-08-23T20:45:42Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 94% |
| Test Adequacy | 96% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** 2 (report only).
1. `moveTimelineRowFocus` calls `renderVisibleRows()` for a change that is, in the common case, two
   `tabindex` attributes and a focus — a full row sweep per arrow press. Held at ~33 rendered rows
   and one press per keystroke, so it is not a measured problem; the `focusin` path already takes the
   attribute-only route where it matters (every entry into the list). Optimising the arrow path would
   mean duplicating the render's own scroll/rebuild ordering, which is the part that was hard to get
   right.
2. Screen-reader announcement of the moved row is unverified. The row's accessible name comes from
   its `<title>` and focus movement should announce it, but nothing here measures that.

**Acceptance:** Pass — one Tab press measured with real input, arrow movement and both prior
contracts green in the probe, canonical gate exit 0.
**Suggested testing:** 3 items
**Follow-ups created:** None; **sweeps appended to:** None

### Requirements Checklist

- [x] The row list occupies one Tab stop; Tab leaves in one press — delivered, and measured with
      trusted input (1 press, from 6+)
- [x] Arrow keys move focus between rows (roving tabindex) — delivered (Up/Down; probe (c))
- [x] Enter and Space still open the drawer — delivered (probe (e); the activation path is untouched)
- [x] REQ-333's arrow panning still works with a row focused — delivered (probe (d), plus REQ-333's
      own keyboard probe passing)
- [x] Focus survives the node rebuilds the focus restore handles — delivered; the roving move relies
      on that restore rather than replacing it, and the ordering constraint the restore's comment
      records is honoured (scroll → synchronous rebuild → focus)
- [x] Tab leaves the list "in each direction" — Shift+Tab is the same single stop by construction:
      one `tabindex="0"` and 34 explicit `-1`s is direction-independent. Not separately measured,
      which is the honest limit of this row.

### Restatement Sweep

The diff redefines what Up/Down do inside the chart. Swept every statement of the keyboard contract:

- **`timelineKeyboardWindow`'s header comment was a half-truth and is fixed in this diff** (D-03) —
  it said returning null leaves Up/Down to scrolling the queue, which now holds only when focus is
  not on a row.
- `web/board-timeline.js:1188-1191` (the `role="button"` / Enter-and-Space comment) — still exactly
  true; that path is untouched.
- `renderVisibleRows`'s focus-restore comment — still true, and the new roving move is a *third*
  caller that depends on it rather than a change to it.
- The chart's on-screen keyboard hint: grepped the generated page and the Go templates for the hint
  text. It names ctrl+wheel zoom and arrow panning; it does not claim anything about Up/Down, so
  there is nothing stale to correct. Checked rather than assumed.
- `skills/do-work-board/actions/board.md` and `skills/do-work/docs/` — no shipped prose describes the
  Timeline's per-key behaviour.

### Code Review Notes

- **Naming for reach:** four new identifiers — `timelineTabbableRowIndex`,
  `scrollTimelineRowIntoView`, `moveTimelineRowFocus`, `rovingRowIndex` — all multi-word and
  greppable. `rovingFrom`, `previousStop`, `nextRowGroup` are locals within a screen of their use.
- **The ordering in `moveTimelineRowFocus` is the load-bearing part** and it is documented in place:
  scroll, synchronous rebuild, then focus. Written the other way round, the asynchronous scroll event
  the scroll write schedules would rebuild over the focused node — the exact regression REQ-333's
  restore comment records. Verified by the probe's ArrowDown/ArrowUp pair rather than by reading.
- **`timelineTabbableRowIndex` is pure and separately reasonable:** it takes three numbers and
  returns one, and its `Math.max(firstRenderedRow, lastRenderedRow - 1)` guard means an empty
  rendered range cannot produce an index below the range's own start.
- **`isFinite` on the parsed `data-row-index`** is a real guard, not decoration: `Number("")` is 0,
  which would silently rove to the first row, and `Number(null)` is 0 as well.
- **Code Quality 94%, not higher,** for Minor 1: the arrow path pays a full render where the focus
  path does not, and the asymmetry is a deliberate trade rather than a clean design.

### Acceptance Testing

**Result: Pass**
- Probe: RED first (35 of 35 tabbable, no arrow movement), GREEN after, with the two prior contracts
  green throughout.
- Real input: one Tab press to leave the list, against 6+ before.
- Full board suite with the browser lane on, and the canonical gate: both clean, so REQ-333's
  keyboard probe, REQ-336's click fix and REQ-337's capture check all survive.
- Adjacent path exercised deliberately: the CDP click probe from REQ-336 still opens the drawer on a
  bar click after this change, so the roving `tabindex` did not disturb the pointer path.

### Suggested Additional Testing

- Shift+Tab with real input, to measure the reverse direction rather than infer it from the
  attribute count.
- A screen reader on a real desktop: confirm an arrow-moved row announces its `<title>`, and that the
  roving stop does not make the list read as a single control.
- A filtered board where the roving row falls outside the filtered set, then arrow from the list.
  `moveTimelineRowFocus` clamps against the live `rows.length`, so this should be safe, but it is the
  interaction least covered by the probe.

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Writing the probe before the implementation, as `tdd: true` asked. Two of its five assertion
  groups (arrow panning, Enter activation) passed on the unchanged code, which made the probe a
  regression guard for REQ-333's and REQ-336's contracts *before* there was anything new to break
  them — and the failing three named exactly what had to change.
- Reading the neighbouring comments as a specification. `renderVisibleRows`'s focus-restore comment
  states the scroll-is-asynchronous ordering trap outright, which is why the roving move was written
  scroll-then-rebuild-then-focus first time instead of after a debugging round.
- Choosing the axis rather than negotiating it. Left/Right were already taken by panning and Up/Down
  were free, so the key assignment cost nothing — no shared modifier, no mode.

**What didn't:**
- Nothing was reverted. The one thing measured rather than reasoned about was the Tab count itself:
  it needed trusted input, so the probe lane could not answer it and a separate CDP run did.

**Worth knowing:**
- **A roving tabindex over a VIRTUALIZED list needs a clamp, not a match.** The stored roving row is
  usually not rendered, and marking `tabindex="0"` only on an exact match takes the whole list out of
  the Tab order — a worse defect than the one being fixed. Clamp for display, keep the stored index
  intact so the reader's place survives scrolling away.
- **Every other row needs an explicit `tabindex="-1"`.** Dropping the attribute is not neutral: a
  focusable-by-default element with no `tabindex` is still a Tab stop.
- **Tab cannot be tested with synthetic events.** Its focus movement is a default action the engine
  performs only for trusted input, so a probe that dispatches `new KeyboardEvent("keydown", {key:
  "Tab"})` observes nothing and reads as a pass. Assert the property that produces the behaviour —
  exactly one stop — and measure the behaviour itself somewhere that can dispatch real keys. Arrow
  movement is different and *is* testable synthetically, because the handler calls `focus()` itself.
- The `focusin` sync is not optional polish: without it the roving index and the actual focus diverge
  the moment a reader clicks a row, and the next arrow press jumps from the stale index.

## Orientation

Tabbing past the Timeline chart now takes one press instead of one per row, and the arrow keys walk
the rows. Lives in the board's Timeline view (`skills/do-work-board/tools/queue-kanban`) — the row
markup, the keydown handler, and a new focus-sync listener. Left/Right panning, Enter/Space
activation, and the rebuild focus restore are unchanged, so no contract moved and the system's shape
is the same. `_dev/primes/prime-kanban-board.md` gains one lesson link and its referenced paths all
still resolve.
