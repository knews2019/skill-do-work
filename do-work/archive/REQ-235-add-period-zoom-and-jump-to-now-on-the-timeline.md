---
id: REQ-235
title: "Addendum: give the Timeline period-based navigation and a jump to now"
status: completed
completed_at: 2026-08-18T11:37:10Z
commit:
claimed_at: 2026-08-18T11:13:25Z
route: C
estimate:
  p50_active_minutes: 75
  confidence: medium
  calculated_at: 2026-08-18T11:13:25Z
  basis:
    - Route C
    - 4-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
    - browser evidence
    - performance work
created_at: 2026-08-18T10:22:05Z
user_request: UR-052
addendum_to: REQ-227
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-233]
batch: board-timing-views
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-timeline.js
- skills/do-work-board/tools/queue-kanban/web/template.html
- skills/do-work-board/tools/queue-kanban/web/board.css
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Addendum: Give the Timeline Period-Based Navigation and a Jump to Now

## What

On a board with 677 REQs spanning four months, the Timeline view cannot be navigated: sideways movement exists only as a mouse drag, and there is no way to land on the work that is still open. Add two things alongside the existing zoom: a day/week/month period level with previous/next stepping, and a control that jumps both the time window and the row list to the now-line and the forecast.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`, `CLAUDE.md`, the always-on crew rules, and REQ-233's just-shipped `timelinePannedWindow` / `timelineKeyboardWindow`; settled on period arithmetic that produces a *candidate* window and hands it to the existing `timelineZoomedWindow`, with the level derived from the window rather than stored.
- [x] **[APPLY]:** Four declared files touched, nothing outside them. `timeline.go` deliberately untouched — the Now jump reads the already-shipped `projection.queueEnd`, so no payload change was needed.
- [x] **[UNIFY]:** `git diff --stat` reviewed across all four files; `go vet` and `gofmt -l` clean; full `maintainer-verify.sh` exit 0 after every edit including the `wireZoomButton` → `wireToolbarButton` rename; REQ-233's two keyboard tests re-run and passing; no debug artifacts in the diff.

<!-- Boxes ticked by the ORCHESTRATOR at Step 6.3, from evidence in the builder's
     hand-back and independently re-checked against the merged tree. The dispatch
     brief omitted the P-A-U instruction (the same omission as REQ-236); a mid-run
     correction was sent but arrived after the builder had committed. Recorded here
     rather than presented as a builder self-report. -->

## Why

The user's words: "this is not working well, can not scroll horizontally, I can not jump to the remaining work."

## Context

Addendum to REQ-227 (completed, commit `17b9422`), which built the view. On the user's board every drawn bar is crushed into the last ~12% of the plot width because `Fit all` spans the whole capture history, so reading anything means zooming in and then dragging a long way — and 97 REQs are still open, i.e. exactly the part that is hardest to reach.

The user raised drawing cost as part of the reason for the period model: with a long range, stepping by a fixed period is cheaper and more predictable than free panning across months. Capture is recording that motive, not prescribing a mechanism.

Sits next to REQ-233 (keyboard zoom and pan, `pending-answers`), which touches the same files. Whichever runs second inherits the other's controls; both must end up driving the one `timelineZoomedWindow` transform rather than growing a second window model.

## Prior Implementation

REQ-227 built the view across `timeline.go` (payload: `rangeStart`, `rangeEnd`, `now`, per-row wait/work spans, projection), `generate.go` (embed), and `web/board-timeline.js` (rendering, ~735 lines), plus `template.html`, `board.css`, `board-controls.js`, `board-filters.js`, `board.js`. What exists today in `board-timeline.js`:

- `timelineViewState = { windowStartMs, windowEndMs, fitted }` is the single window model; `timelineZoomedWindow(...)` is the only transform that resizes it, clamped to `boundStartMs`/`boundEndMs` (payload range, extended to the projected queue end, plus 2% padding).
- Rows are virtualized vertically inside `#timeline-scroll` (`timelineVisibleRowRange`, re-rendered on `scroll`).
- Zoom: `−` / `+` / `Fit all` buttons (`wireZoomButton`) and ⌘/Ctrl+wheel anchored at the pointer. A plain wheel is deliberately left alone so it scrolls rows — which is why a trackpad's sideways swipe does nothing at all.
- Pan: `pointerdown`/`pointermove` drag only.
- `Now` (`timeline-zoom-now`) recentres the window on `nowMs` at the current span. It does not change the row scroll position, so the rows on screen are usually still the old ones.
- Listeners go through `addTimelineListener` because the scroll host and `window` outlive a render; anything added here must use it too.

## Detailed Requirements

- A period level control offering **day / week / month**. Choosing one sets the window to exactly that period, calendar-aligned: day = midnight→midnight UTC, week = Mon→Sun, month = the 1st→the last day.
- **Previous / next period** controls that step the window by exactly one period at the chosen level, landing on clean boundaries every time, and clamping at the ends of the available range.
- The period controls sit **alongside** the existing zoom, not in place of it: `−`, `+`, `Fit all`, ⌘/Ctrl+wheel zoom and drag-pan all keep their current behaviour. A free zoom or drag leaves the period level showing as no longer exact rather than silently lying about it.
- The `Now` control jumps to **now and the forecast**: the window covers the now-line and the projected queue-empty time, and the row list scrolls to the still-open work so the reader lands on rows they can see, not on whatever was in view before.
- All of it drives the existing `timelineViewState` / `timelineZoomedWindow` window model — no second window model, no divergence between the pointer path and the new controls.
- Stepping by period should not get slower as the board's total range grows; the user named long-range performance as a reason for this design.

## Constraints

- Timeline stays read-only. The board's three write surfaces (CLAUDE.md § Kanban Board Write Surfaces) are unchanged — nothing here writes pipeline state.
- Period arithmetic is UTC, matching the payload's stamps and the axis labels.
- New listeners go through `addTimelineListener`; new controls follow the existing `control-button` markup in `template.html`.
- Do not touch `durations.go` / `board-durations.js` — REQ-231 is open against that panel.

## Red-Green Proof

**RED prompt/case:** A Node behaviour probe in `generate_test.go` (the `TestJavaScriptBehavior*` family) driving the new period navigation over a fixture whose range spans several months: setting the level to `week` produces a window that starts on a Monday 00:00 UTC and is exactly seven days long; `next period` advances it by exactly seven days and stays Monday-aligned; stepping past the last period clamps instead of running off the range. Plus: invoking `Now` puts both `now` and the projection's `queueEnd` inside the window **and** moves the scroll host's `scrollTop` to the first still-open row.

**Why RED now:** there is no period level and no prev/next stepping to drive, and `timeline-zoom-now` only recentres the window — it never touches `scrollTop`, so the row-list assertion fails today.

**GREEN when:** the probe passes, and a headless render of a multi-month fixture shows the day/week/month controls, prev/next landing on calendar boundaries, and `Now` bringing the open work into view in one click. `bash _dev/tests/maintainer-verify.sh` exits zero.

**Validation:** User adjusted — the period model, the day/week/month levels, and the long-range performance motive are the user's own redirection of the sideways-scroll question; "alongside the existing zoom" and "calendar-aligned" were confirmed by the user during capture.

## Assets

Screenshot supplied inline with the request; the attachment was no longer resolvable on disk at capture time, so it could not be persisted under `assets/`. The full description is preserved verbatim in `do-work/user-requests/UR-052/input.md` § Screenshot Description — 677 REQs, 97 open, all bars crushed into the right ~12% of the plot, toolbar showing only `−` `+` `Now` `Fit all`.

## Builder Guidance

Certainty level: Firm on the four answered decisions (period levels, calendar alignment, alongside-not-replacing, Now = window + row list). Latitude on the toolbar layout and on how period stepping is implemented, as long as it goes through the one window model. Keep it small — this is navigation on top of a view that already works, not a redesign.

---
*Source: "do-work capture-request [screenshot] <- this is not working well, can not scroll horizontally, I can not jump to the remaining work"*

---

## Triage

**Route: C** - Complex

**Reasoning:** Multiple new controls across four files, calendar arithmetic that has to be exactly right, and a hard architectural constraint inherited from REQ-233 — every way of moving the window must route through one transform. The design question (store the level or derive it) was worth settling before writing code.

**Planning:** Required

## Plan

Period arithmetic produces a *candidate* window by pure UTC calendar maths, then hands it to the existing `timelineZoomedWindow` for the floor, ceiling and edge clamp — so the period path gets no rules of its own. The level is read back off the window rather than stored, which makes "a free zoom marks the level inexact" true by construction instead of by remembering to clear a flag.

*Generated by Plan agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modify) — period arithmetic, level derivation, Now jump, control wiring
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify) — period control group, panel label, hint line
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — the period-state readout
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — Node behaviour probe

**Files I will NOT touch:** `timeline.go` (no payload change — `projection.queueEnd` already ships), `durations.go` / `board-durations.js` (REQ-231's territory, forbidden by this REQ's constraints), `board-controls.js` (`setActiveButton` is reused unmodified).

**Acceptance criteria (restated from REQ):**
- [ ] Day/week/month levels, calendar-aligned in UTC — day midnight→midnight, week Mon→Sun, month 1st→last
- [ ] Previous/next stepping by exactly one period, landing on clean boundaries, clamping at the range ends
- [ ] Period controls sit alongside the existing zoom; `−`, `+`, `Fit all`, ⌘/Ctrl+wheel and drag-pan all keep their behaviour
- [ ] A free zoom or drag shows the level as no longer exact rather than silently lying
- [ ] `Now` covers the now-line and the projected queue-empty time, and scrolls the row list to the still-open work
- [ ] All of it drives the existing `timelineViewState` / `timelineZoomedWindow` — no second window model

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** A `‹ Day Week Month ›` control group joined the timeline toolbar ahead of the untouched zoom group, with a `#timeline-period-state` live region beside it. `timelinePeriodWindow` computes a candidate window by pure UTC calendar arithmetic (`timelinePeriodStart` floors to the level's boundary, `timelineSteppedPeriodStart` moves it by N periods) and then hands that candidate to `timelineZoomedWindow(start, end, 1, 0, bounds)` — factor 1 at anchor 0 meaning "keep this window, apply the model's floor, ceiling and edge clamp". That single call is why the period path has no floor, ceiling or clamp of its own.

No level is stored. `timelinePeriodLevelOfWindow` reads it back off `timelineViewState` — a window *is* a level exactly when it starts on that level's boundary and ends on the next one — and `renderPeriodControls()` runs from `renderAll()`, which every mover already calls. So a wheel-zoom, a drag or `Fit all` drops the pressed state and writes `custom span` with no invalidation logic anywhere.

The `Now` handler now moves both the window and the row list: it sizes the window to `[min(now, queueEnd), max(now, queueEnd)]` plus a 10% margin, and sets `scrollTop` to `timelineFirstOpenRowIndex × TIMELINE_ROW_HEIGHT`, leaving the scroll untouched when nothing is open. The local `wireZoomButton` became `wireToolbarButton`, since it now wires period steps too.

## Qualification

Passed — 4 files verified in the merge range `07e9162..7cae7a4`, 6 acceptance criteria traced.

**Merge conflict resolved by the orchestrator**, same shape as REQ-236's: this REQ and REQ-236 both appended tests at the end of `generate_test.go`. Both sides were pure appends of separate functions; the resolution keeps both, `gofmt -w` restored spacing, `go vet` is clean, no duplicate function names, and the full suite passes on the union of all three board REQs' tests.

Judgment checks, run against the merged tree rather than taken from the builder's report:
- **Decision 1's claim — that the deleted clamp was genuinely unearned — was tested, not accepted.** Forty `›` clicks past the end of the range, then one more: the axis strings are byte-identical (`13 Aug … 20 Aug`). Stepping converges and stops with no period-index guard, because the anchor is the window's own midpoint and the bounds clamp already keeps that inside the range. Deleting the second clamp was correct.
- **Calendar alignment verified against real weekdays,** not against the code's own arithmetic: the week windows start on `6 Jul`, `13 Jul`, `29 Jun` 2026, and `Date.prototype.toUTCString` confirms all three are Mondays. Month windows run 1st→1st; day windows 00:00→00:00.
- **The `Now` jump's row target was checked by measurement.** `scrollTop` moved 0 → 720, and `REQ-041` — the first still-open REQ — measures a `top` of 265 against a scroll-host top of 264. It lands at the viewport top, as specified.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Exit 0 on the merged tree, run unpiped — the run that matters, since this is where all three board REQs' tests coexist for the first time

**Red-green validation:**
- `TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow`: ✗ → ✓. The first RED was a reference-error class, so a second was produced with the code present and one behaviour wrong — `timelinePeriodStart` computing a Sunday-origin week: `the week window starts at 2026-06-14T00:00:00.000Z (weekday 0); want a Monday at 00:00 UTC`. That is the assertion failing for the reason it exists.
- **REQ-233's `TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer` and `TestTimelinePanelStatesItsKeyboardInteraction` both still pass** — checked explicitly, since this REQ is the third driver of the window REQ-233 made single-sourced.

**New tests added:**
- `TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow` — week windows start Monday 00:00 UTC and run exactly seven days; `next` advances seven days and stays Monday-aligned; stepping past the last period clamps; `Now` puts both `now` and `projection.queueEnd` inside the window and moves `scrollTop` to the first still-open row
- `rendererDeclarationLine` — a sibling of `rendererNumericConstant`, so the probe drives the shipped `TIMELINE_PERIOD_LEVEL_NAMES` rather than a copy that would go stale beside it

**Render evidence — driven in a browser against the merged tree.** Generated from the live repo (234 REQs, 19 open, 28 May → 20 Aug):
- Controls render as `‹ Day Week Month ›` with a state readout.
- `Week` → `6 Jul … 13 Jul` (`one week`, Week pressed); `›` → `13 Jul … 20 Jul`; `‹` → exact round-trip to `6 Jul`. `Month` → `1 Jul … 1 Aug`. `Day` → `16 Jul 00:00 … 17 Jul 00:00`.
- **Clamp:** 40 `›` then one more — identical axis, `custom span`, no button pressed.
- **Inexactness:** after a free `+` zoom from an exact week, the readout reads `custom span` and no level button is pressed.
- **`Now`:** `scrollTop` 0 → 720, `REQ-041` at the viewport top, the now-line on screen at x=1042, window covering now and the forecast.

*Verified by work action*

## Decisions

- **D-01**: The period-index clamp was **deleted rather than kept**. It was written, and the tests passed without it — anchoring on the window midpoint plus `timelineZoomedWindow`'s bounds clamp already makes stepping converge. Per CLAUDE.md's *delete before you add* and coding-guardrails § Earned Defense, a second clamp is exactly the "second definition of where the window goes" this REQ exists to prevent. Independently re-verified by the orchestrator (40 steps + 1, identical window). DECIDE & STATE.
- **D-02**: The level is derived from the window, never stored. A stored `periodLevel` is a second thing that can disagree with the window, and would need explicit invalidation on every zoom, drag and keypress. Deriving it makes the inexactness requirement true by construction. DECIDE & STATE.
- **D-03**: The anchor for both level selection and stepping is the window's midpoint — one rule rather than two. **Known wrinkle, flagged rather than hidden:** at the extreme end of the range the window is clamped and its midpoint can fall in the previous period, so a `‹` immediately after hitting the end can skip one period. Fixing it needs remembered state, which D-02 rules out. One press, one edge. DECIDE & STATE.
- **D-04**: At the range ends the shared clamp wins over grid alignment; such a window honestly reads `custom span` rather than claiming a whole period. Widening the bounds for the period path only would have given it a different range from pan and drag. DECIDE & STATE.
- **D-05**: `Now` covers `[min(now, queueEnd), max(now, queueEnd)]` plus a 10% margin rather than preserving the current span — preserving it cannot guarantee both are visible when the reader is zoomed in. DECIDE & STATE.
- **D-06**: `Now` puts the first still-open row at the top of the viewport, and leaves the scroll alone when nothing is open. Context rows above would have been prettier and less exactly assertable. DECIDE & STATE.
- **D-07**: `wireZoomButton` → `wireToolbarButton`, because it now wires period steps too — cleaning up its own mess, not improving adjacent code. Local function, five call sites, no references elsewhere. DECIDE & STATE.
- **D-08** (orchestrator): the `generate_test.go` merge conflict was resolved by keeping both sides in order, as with REQ-236.

## Discovered Tasks

- **[normal]** The axis tick label prints a literal `:00` minute at sub-hour tick spacing, so several ticks inside one hour all render the same label. Measured after `Now`: 7 ticks, 2 distinct labels. Pre-existing from REQ-227, but this REQ's `Now` jump lands in exactly that window by design. **Escalated to an Important review finding and queued as REQ-240** rather than left here.
- **[low]** A week window's six evenly spaced ticks land at 1.167-day intervals, so interior labels skip a day (`7, 8, 9, 10, 11`, no 12). The two ends are exact, so alignment is still visible. A period-aware tick generator would fix this and the above together, but that is an axis redesign.
- **[low]** `Fit all` spanning the whole capture history is what crushes every bar into the right-hand fraction of the plot in the first place. The period controls make that navigable; a "fit to the last N days" default would address UR-052's original complaint more directly.

## Review

**Overall: 95%** | 2026-08-18T11:37:10Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 97% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition):**
- The Timeline axis prints a literal `:00` minute, so the window `Now` lands in shows 7 ticks with 2 distinct labels — five reading `18 Aug 11:00`. The defect is REQ-227's, but this REQ makes it the routine landing state of its most-used new control, on a UR whose complaint was precisely "I can not jump to the remaining work" — gate: **user-visible** → REQ-240 created. Not a sweep: one root cause, one site, and it wants its own before/after render.

**Minor findings:** 1 (report only)
- D-03's midpoint-anchor wrinkle: a `‹` immediately after clamping at the range end can skip a period. The builder found it, chose not to fix it because the fix needs the stored state D-02 removed, and flagged it. That is the right call and the right disclosure; noted so the next reader meets it in the record rather than in the UI.

**Restatement sweep:** this REQ adds a third driver of the window REQ-233 made single-sourced, so the sweep asked what states that contract. The Timeline's interaction is described in three places — the panel `aria-label`, the `.timeline-hint` line, and the probes — and the builder updated the first two alongside the code. `_dev/primes/prime-kanban-board.md` describes write surfaces, not interactions, and this REQ adds none. REQ-227/228's archived records and `CHANGELOG.md`'s Timeline entries are dated history. No stale restatement.

**Acceptance:** Pass — all six criteria confirmed by driving the merged build in a browser, including the clamp that survives without an explicit guard and the `Now` jump landing REQ-041 at the viewport top.

**Suggested testing:** 2 items
- Read the axis at each period level on the consuming board (677 REQs, four months), which is the board this REQ was actually written for. REQ-240 will change what those labels say; a human should see both.
- The D-03 edge wrinkle is one keypress at one boundary and was reasoned about rather than measured on a real board. Worth someone stepping to the end of the range and pressing `‹` once.

**Follow-ups created:** REQ-240; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Writing the defensive clamp, then deleting it because the test passed without it. The single-window-model constraint did not just prevent divergence — it made a whole guard unnecessary, and the only way to discover that was to build the guard and then check whether removing it broke anything. "Delete before you add" is usually advice about the first draft; here it applied to code that was already written and already passing.

**What didn't:** Two of the three RED attempts across this batch started as reference errors — the constant or function simply did not exist yet. That proves the anchor is missing, not that the behaviour is wrong. The fix that worked twice was to put the code in place and break exactly one rule inside it: a Sunday-origin week here, a missing pan clamp in REQ-233. The assertion then fails for the reason it was written, which is the only failure worth calling RED.

**Worth knowing:** Deriving state instead of storing it removed an entire class of bug from this change. There is no `periodLevel` field, so there is nothing to invalidate on zoom, drag, keypress or resize, and "a free zoom marks the level inexact" needed no code at all — it is what happens when the level is a question you ask the window rather than a fact you remember about it.

## Orientation

The Timeline can now be navigated by calendar period: `Day`, `Week` and `Month` set the window to exactly that period in UTC, `‹` and `›` step by one, and `Now` jumps the window to cover the now-line and the forecast while scrolling the row list to the first still-open REQ. A readout beside the controls says which level is exact, or `custom span` when a free zoom or drag has left the grid. Lives in the queue-kanban board subsystem (`_dev/primes/prime-kanban-board.md`).

Not `[MAP CHANGED]` — and that is the point. This is the third driver of the window REQ-233 made single-sourced, and it added no fourth rule: period stepping computes a candidate and hands it to `timelineZoomedWindow` like everything else, and the level is derived rather than stored, so `timelineViewState` is still the one window model. The shape REQ-233 established absorbed a substantial feature without changing. Staleness spot-check on `_dev/primes/prime-kanban-board.md`: every referenced path resolves and the three-write-surface count is unchanged — this REQ adds none, and the Timeline stays read-only. The prime is not stale.
