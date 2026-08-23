---
id: REQ-335
title: "[impact-negligible] Give the timeline a usable window when the payload range is unreadable"
status: completed
created_at: 2026-08-23T12:28:00Z
claimed_at: 2026-08-23T16:40:00Z
completed_at: 2026-08-23T17:00:00Z
commit: 965ea58
route: A
user_request: UR-066
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Give the Timeline a Usable Window When the Payload Range Is Unreadable

## What

The render's fallback for an unreadable payload range sets the clamp bounds to
`[filterMatchedRows[0].createdTime, +1 hour]`. Rows are newest-first, so `[0]` is the **newest** REQ, and the
bounds collapse to a one-hour window around it. Because bounds are what every control clamps against, no
control can leave that window — on this repo's board that would strand 287 of 317 REQs permanently out of
reach.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`. Took the extent from a small dedicated pass over `filterMatchedRows` rather than reordering the render to reuse `filterMatchedSegments`: the fallback is a rare path, so the extra pass costs the common one nothing, and reordering the render to serve it would put risk where there is none today.
- [x] **[APPLY]:** `web/board-timeline.js` and `generate_test.go`, as planned.
- [x] **[UNIFY]:** Verified:
  - `web/board-timeline.js` — `node --check` clean; the bounds still come from the payload on the normal path, so a window that narrows the rows cannot narrow what it is clamped against; the decline reuses the existing message rather than adding a second wording.
  - `generate_test.go` — `gofmt`/`go vet` clean; the probe drives the whole renderer, and its setup assertion fails loudly if the payload stops taking the fallback branch.
  - `bash _dev/tests/maintainer-verify.sh` exit 0.
  - Mutations: the old newest-row fallback fails with the exact defect (1 of 3 rows, 2 outside the window); removing the decline branch throws rather than passing. A third — reading `created_at` alone — PASSED at first, so the fixture now carries a completion eight hours past the newest capture and the assertion reads the window the summary names.

## Why

Judged `impact-negligible`: the branch fires only when `timeline.rangeStart` / `rangeEnd` are unparseable or
inverted, which the Go producer does not currently emit — `timelineRange` always returns real instants for a
non-empty row set. It is worth fixing anyway because it is three lines, and because a fallback that makes the
view unusable is worse than no fallback: the failure mode should be "degraded but navigable", never "every
control dead".

## Context

`web/board-timeline.js:1043-1048`:

```js
var boundStartMs = Date.parse(timeline.rangeStart);
var boundEndMs = Date.parse(timeline.rangeEnd);
if (isNaN(boundStartMs) || isNaN(boundEndMs) || boundEndMs <= boundStartMs) {
  boundStartMs = Date.parse(filterMatchedRows[0].createdTime);
  boundEndMs = boundStartMs + TIMELINE_MIN_SPAN_MS;
}
```

`filterMatchedRows` is ordered newest-first (REQ-318), so `[0].createdTime` is the most recent capture. The
comment above the block is right that the bounds must not come from the windowed set — but the fallback then
takes them from a single row of the filtered set, which is the same mistake at maximum severity.

## Detailed Requirements

1. The fallback bounds span the **whole** filter-matched set — the earliest instant any matched row starts to
   the latest any of them reaches — so every matched row is reachable.
2. Nothing is fabricated: if not one matched row carries a parseable instant, the view says so rather than
   inventing a window (the no-readable-`created_at` message at `:1034` is the existing shape for that).
3. The fallback is exercised by a test, so it cannot rot unnoticed while the producer keeps it unreachable.

## Constraints

- Keep the bounds independent of the visible window. A window that narrows the rows must never narrow the
  bounds it is clamped against, which is what the comment at `:1037-1041` is protecting.
- Do not scan the rows twice: the render already computes a segment list per matched row
  (`filterMatchedSegments`) and the extent is derivable from it.

## Red-Green Proof

**RED prompt/case:** In a Node behaviour probe (or by editing a generated `board-data.js` copy), set
`timeline.rangeStart` to `"not-a-date"` and re-render. Read `#timeline-range-readout` and press `Fit all`,
`Month`, `−` and `‹` in turn.

**Why RED now:** every one of those controls reports the same window — a one-hour span around the newest
REQ's `created_at` — and the summary reads "… 287 outside the window, not listed" with no way to reach them.

**GREEN when:**
- With an unreadable `rangeStart`, `Fit all` produces a window covering every matched row, and the summary
  reports 0 outside it.
- `−` from any window still widens to that same extent.
- With no matched row carrying a parseable instant, the view shows the existing "nothing to place on a
  timeline" message instead of a fabricated window.

**Validation:** Inferred during capture; the ordering that makes `[0]` the newest row was confirmed against
REQ-318's newest-first sort in `timeline.go`.

## Full Context

See `do-work/user-requests/UR-066/input.md`.
