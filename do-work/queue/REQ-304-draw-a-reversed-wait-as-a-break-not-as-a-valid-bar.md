---
id: REQ-304
title: "Draw a reversed wait as a break, not as a valid bar"
status: pending
created_at: 2026-08-20T08:37:41Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
