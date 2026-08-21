---
id: REQ-304
title: "Draw a reversed wait as a break, not as a valid bar"
status: completed
created_at: 2026-08-20T08:37:41Z
status_changed_at: 2026-08-21T08:16:04Z
claimed_at: 2026-08-21T08:16:04Z
completed_at: 2026-08-21T08:26:15Z
kb_status: pending
commit: 5e08a31
user_request: UR-062
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
related: [REQ-280, REQ-303, REQ-305]
batch: upstream-consumer-report-2026-08-20
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-21T08:21:11Z
  basis:
    - Route B
    - 2-file write set
    - new DOM probe harness
    - 5 acceptance criteria
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-timeline.js
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Draw a Reversed Wait as a Break, Not as a Valid Bar

## What

The Timeline view draws the wait segment unconditionally (`web/board-timeline.js:639-641`), and
`drawSegment` sorts its endpoints with `Math.min`/`Math.max` (`:591-593`), so a wait whose
`claimed_at` precedes its `created_at` paints as an ordinary positive-width waiting bar. The work
segment already handles this: `:655-666` draws a break marker at the claim instant when
`workMinutes < 0`, with the comment "A reversed span has no width to draw honestly". Give the wait
the same treatment.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Mirror the existing `workMinutes < 0` branch for the wait, anchoring the break at the wait's own start instant (`createdMs`); build a timeline DOM stub modelled on `durationsRenderDomStubPreamble` and assert all three wait shapes in one render.
- [x] **[APPLY]:** `web/board-timeline.js` and `generate_test.go` only. `timeline.go` and `model.go` untouched — no second anomaly definition.
- [x] **[UNIFY]:** `git diff --stat` reviewed (2 files); `gofmt -l .` clean; `go vet ./...` clean; `go test -count=1 ./...` green in 70s; scratch probe file deleted; both changed files re-read at every edit site.

## Why

The bar and the numbers disagree in the same row. `:1007` prints the signed value into the table, so
a reversed wait reads `−12 min` in the table beside a bar drawn as a valid 12-minute wait. Broken
bookkeeping that renders as healthy is worse than broken bookkeeping that renders as broken — the
view's own break marker exists to say so.

## Context

From the 2026-08-20 consumer review, filed P2. Verified against the code; one half of the requested
remedy is deliberately **not** in scope, and the reason is a documented decision:

- `timeline.go:22-24` states that `detectCompletionAnomaly` (`model.go`) decides what is broken and
  this file consumes that verdict rather than restating it. `model.go:237-244` scopes
  `CompletionAnomaly` to completion bookkeeping only, which is why `row.anomaly` is false for a
  reversed wait. Recomputing an anomaly verdict inside the renderer would create the second
  definition that comment exists to prevent.
- Queued **REQ-280** already owns the detection side: adding the
  `created_at <= claimed_at <= completed_at` ordering probe to `queue-kanban verify` and to
  `forensics.md` Check 12. The anomaly-reporting half of this finding belongs there, not here.
- The reviewer's framing understates one thing: the summary line's anomaly count (`:541-552`,
  "N with broken stamps, drawn as breaks") already under-counts reversed **work** spans for the same
  reason, since a `completed_at` that parses but precedes `claimed_at` sets no anomaly flag. Do not
  fix that here by broadening the flag; note it against REQ-280 if it is not already covered.

## Detailed Requirements

- When a row's `waitMinutes` is negative, draw the break marker at the row's own reference instant
  instead of calling `drawSegment` — mirroring the existing `workMinutes < 0` branch, including its
  reasoning comment.
- The open-wait case (`waitOpen`, measured to the now-line) must keep its current rendering; only a
  genuinely negative span changes.
- Do not touch `drawSegment`'s min/max sorting. It is correct for every caller that reaches it.
- Do not add an anomaly verdict, flag, or reason to the timeline row model.

## Constraints

- Renderer-only. `timeline.go` and `model.go` are out of the write set on purpose.
- `web/` is embedded, not read at runtime (`prime-do-kanban.md` § Traps) — the change needs a
  rebuild to be visible, and the test must run against the generated client.

## Red-Green Proof

**RED prompt/case:** Render a timeline row whose `claimedTime` precedes its `createdTime` (negative
`waitMinutes`) and inspect the emitted SVG. Today it contains a `timeline-segment-wait` rect with
positive width, and no `timeline-segment-broken` rect.
**Why RED now:** The wait segment is drawn with no sign check, and `drawSegment` sorts the endpoints
so a reversed span still produces a left-to-right rect.
**GREEN when:** That same row emits the break marker (`timeline-segment-broken`) and no
`timeline-segment-wait` rect, while a row with an ordinary positive wait still emits its wait rect
unchanged.
**Validation:** Inferred during capture — read from the two branches side by side in the source.

## Builder Guidance

Certainty: Firm on the behavior, open on the test harness. There is no existing DOM probe for
`renderVisibleRows` — `generate_test.go`'s `durationsRenderDomStubPreamble` (~line 2305) is the
pattern to copy: stub the smallest DOM the real renderer touches, record every SVG node's tag and
attributes, and assert against that rather than re-implementing layout. Budget the harness, not the
fix; the fix itself is a handful of lines.

## Full Context
See `do-work/user-requests/UR-062/input.md` for complete verbatim input.

---
*Source: "[P2] Render negative waits as broken spans — … this unconditional call passes both endpoints to drawSegment, which sorts them with min/max and paints a normal waiting bar. Only negative work spans receive the broken marker below … so reversed waits are visually presented as valid; add equivalent negative-wait handling and anomaly reporting."*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modify) — the negative-wait branch
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — the timeline DOM stub and the three-shape behavior probe

**Files I will NOT touch:** `timeline.go` and `model.go` — the REQ forbids a second anomaly verdict in the renderer; `drawSegment` — its min/max sort is correct for every caller that should reach it.

**Acceptance criteria (restated from REQ):**
- [x] A negative `waitMinutes` draws the break marker instead of calling `drawSegment`
- [x] The open-wait case keeps its current rendering
- [x] `drawSegment`'s min/max sorting is untouched
- [x] No anomaly verdict, flag, or reason added to the timeline row model
- [x] `bash _dev/tests/maintainer-verify.sh` exits 0

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified) — the wait segment now branches on `row.waitMinutes < 0`, drawing the same fixed break marker the work segment draws, anchored at the wait's own start instant
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — `timelineRenderDomStubPreamble`, `timelineRenderProbeDriver`, `TestJavaScriptBehaviorReversedWaitDrawsAsABreak`

**What was done:** Fifteen lines in the renderer, mirroring the branch six lines below it including its
reasoning. The test runs the real `renderTimelineView` over a stub DOM and asserts three wait shapes
in one render pass — reversed, ordinary closed, and open — because a fix that turned every wait into
a break would satisfy a reversed-only test.

## Decisions

- **D-01 — The break marker anchors at `createdMs`, the wait's own start, not at the claim instant.**
  DECIDE & STATE. The work branch anchors at `workStartMs`, its own start; anchoring the wait at the
  claim would stack two markers on the same x whenever both spans are reversed, and would put the
  wait's break where the work segment begins. Reversible: it is one expression.
- **D-02 — All three wait shapes assert in one render pass over one payload.** DECIDE & STATE. The
  defect was a missing branch, so the failure mode a new test most needs to exclude is a fix that
  over-applies. A reversed-only fixture would pass against `if (true)`.

## Discovered Tasks

- **impact-user-visible** The Timeline summary line's break count is now further out of step with what
  is drawn. `web/board-timeline.js:546-557` counts `row.anomaly` and renders "N with broken stamps,
  drawn as breaks", but `row.anomaly` carries `detectCompletionAnomaly`'s verdict, which
  `model.go:246-252` scopes to completion bookkeeping only. It was already false for a reversed
  **work** span whose `completed_at` parses but precedes `claimed_at` — the REQ's own Context notes
  this — and it is false for the reversed **wait** this REQ just started drawing as a break. So the
  view draws breaks the sentence beneath it does not count. REQ-280 closed the detection side
  (`queue-kanban verify`, `forensics.md` Check 12) and did not touch this count. The fix is not
  "broaden the anomaly flag" — `timeline.go:22-24` exists to stop the renderer growing a second
  verdict — it is to count what was drawn rather than what the flag says.

## Testing

**Tests run:** `QUEUE_KANBAN_BROWSER=<chromium> bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ exit 0 — `Maintainer verification passed.` Board module: `go test -count=1 ./...` ok in 70s.

**Red-green validation:**
- **RED, reproduced on the untouched renderer.** A scratch probe rendered a row with
  `createdTime` 11:00 and `claimedTime` 10:00 and dumped every emitted rect:
  `class="timeline-segment timeline-segment-wait" x=307.5 width=188.5` — a positive-width waiting
  bar, and no `timeline-segment-broken` anywhere. Exactly the RED the REQ captured.
- **GREEN.** The same row now emits `class="timeline-segment-broken" x=493.0 width=6.0` and no wait
  rect, while the row's ordinary work bar is unchanged.

**Mutation-tested — five reverts, each caught by a different assertion:**
- M1 revert the branch (`if (false)`) → "must draw the break marker", caught
- M2 break every wait (`if (true)`) → "an ordinary positive wait must still draw its bar", caught
- M3 draw the break AND the bar → "a reversed wait must draw NO wait bar", caught
- M4 make the marker's width track the reversed magnitude → "must be the same fixed 6-unit mark", caught
- M5 route open waits into the break branch → "an unclaimed REQ must still draw its open wait bar", caught

**New tests added:**
- `TestJavaScriptBehaviorReversedWaitDrawsAsABreak` in `generate_test.go`, with
  `timelineRenderDomStubPreamble` and `timelineRenderProbeDriver` — the first DOM probe for
  `renderVisibleRows`. It runs the real renderer over a stub DOM and reads the emitted attributes;
  it re-implements no layout.

**Existing tests updated (cross-REQ impact):** none.

*Verified by work action*

## Review — 2026-08-21T08:25:28Z

**Overall: 95%**

| Dimension | Score |
|---|---|
| Requirements Compliance | 100% |
| Code Quality | 95% |
| Test Adequacy | 100% |
| Scope Discipline | 100% |
| Risk | None |

**Acceptance: Pass.**

**Requirements Compliance.** All five hold, each evidenced by a mutation. The renderer gained no
anomaly verdict — `row.anomaly` is read exactly where it already was — and `drawSegment` is
byte-identical.

**Findings**

- **F1 — Minor, accepted.** The break marker's geometry (`-3`, `width: 6`, `rowTopY + 2`,
  `TIMELINE_ROW_HEIGHT - 4`) is now written twice in the same function. Extracting a
  `drawBreakMarker` helper would remove the duplication, but the REQ asked for a mirror of the
  existing branch and widening to a refactor is scope the REQ excluded. The test pins the width, so
  the two cannot silently diverge on the number that matters.
- **F2 — Minor, deferred to Discovered Tasks.** The summary line's break count reads
  `row.anomaly`, which is false for both reversed shapes, so the view now draws breaks it does not
  count. Recorded rather than fixed: the REQ's Context explicitly forbids broadening the flag here.

**Constraint check.** Renderer-only holds — `git diff --stat` touches `web/board-timeline.js` and
`generate_test.go` and nothing else. The `web/` embed means the change needs a rebuild to be
visible, and the test reads the embedded asset, so it runs against the generated client.

## Lessons Learned

- **A missing-branch fix needs a fixture that can fail in both directions.** The obvious test is the
  reversed row alone, and `if (true)` passes it. Rendering the reversed, ordinary and open shapes
  in one pass is what makes over-application fail.
- **Anchor a mirrored branch to its own span, not to the one it mirrors.** Copying the work branch's
  `workStartMs` anchor would have put the wait's break at the claim instant, stacking two markers
  whenever both spans reverse. The mirror is the shape, not the coordinate.
- **A stub DOM is cheaper than it looks, and iterating on it is the fast path.** Four Node failures —
  `classList`, `querySelectorAll`, `setActiveButton`, the scroll geometry — each named exactly
  what to add. Guessing the whole surface up front would have taken longer than letting it fail.

## Orientation

**What changed in the map.** The Timeline view's honesty rule now covers both segments instead of
one: a span whose end precedes its start is drawn as a break marker, whether it is the wait or the
work. Nothing about detection moved — the renderer still consumes the board's verdict and never
computes one.

**What this makes true.** A REQ whose `claimed_at` precedes its `created_at` can no longer show a
healthy-looking waiting bar beside a table cell reading a negative number. The board also gained its
first DOM probe for the timeline renderer, so the next change to `renderVisibleRows` has somewhere
to assert.

**Subsystem:** the queue-kanban board tool's Timeline view. Prime: `_dev/primes/prime-kanban-board.md`.

**Discovered-task follow-up:** the summary-line break count was queued as REQ-313 (`pending-answers`, `impact-user-visible`).
