---
id: REQ-235
title: "Addendum: give the Timeline period-based navigation and a jump to now"
status: pending
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
