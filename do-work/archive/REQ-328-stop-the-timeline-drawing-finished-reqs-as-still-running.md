---
id: REQ-328
title: "[impact-critical] Stop the timeline drawing finished REQs as still running"
status: pending
created_at: 2026-08-23T12:08:00Z
user_request: UR-066
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/timeline.go
  - skills/do-work-board/tools/queue-kanban/timeline_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
---

# Stop the Timeline Drawing Finished REQs as Still Running

## What

`buildTimelineAggregate` decides whether a span is still running by asking whether a **stamp string parses**,
never by asking what state the REQ is in — and it re-parses `ticket.CompletedAt` instead of reading
`ticket.CompletionTime`, the completion instant `buildBoard` already resolved for that ticket. On this
repo's own board, **25 of the 26 rows the chart calls "still open" are finished**: 18 `completed` and 7
`cancelled`. Each is drawn as a dashed open bar running to the now-line.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Judged `impact-critical`: this is not a cosmetic mislabel, it is the chart asserting that two dozen finished
REQs are in flight right now. It is the artifact people screenshot. It also inflates every number derived
from it — the summary's open count, the payload's `RangeEnd` (pulled up to `now` by any open row), and
therefore the fitted window every other control clamps against.

## Context

Measured against `board-data.js` from `queue-kanban generate` on this repo, now `2026-08-23T11:13:09Z`:

| Row | Status | What the board resolved | What the Timeline draws |
|---|---|---|---|
| REQ-311 | `completed` | `CompletionTime` = `2026-08-21T15:20:33Z`, source `frontmatter`; the calendar places it on 2026-08-21 | `waitOpen`, a dashed bar of 3546 min running to the now-line |
| REQ-307 | `completed` | `2026-08-20T22:37:38Z`, `frontmatter` | `waitOpen`, 4192 min to the now-line |
| REQ-302 | `cancelled` | `2026-08-20T11:38:27Z`, `frontmatter` | `waitOpen`, 4476 min to the now-line |
| REQ-059 | `completed` | `CompletionUnresolved` — commit `d414590` unknown to git; already flagged `CompletionAnomaly` | `workOpen`, a dashed bar of 35764 min (24.8 days) to the now-line |
| REQ-058, 057, 056 … 051 | `completed` | same, all nine already anomalous | same |

Two distinct causes, in the same two lines:

1. **The resolved instant is ignored.** `timeline.go:107` and `:112` call `parseTimestamp(ticket.CompletedAt)`.
   `RequestTicket` already carries `CompletionTime` (resolved via `resolveCompletionTime`, `model.go:1356`,
   frontmatter then git-log) and `CompletionTimeSource`. For REQ-311/307/302 and the eight `workOpen`
   completed rows the instant exists and the Timeline simply does not read it.
2. **"Open" is inferred from a missing stamp, not from a live status.** A `cancelled` REQ with no
   `claimed_at` gets `WaitOpen = true` and a bar to now. Nothing is waiting to claim a cancelled REQ.

The nine genuinely unresolvable rows (REQ-051 … REQ-059) are the ones with no honest end at all. They are
already `CompletionAnomaly`, and `#timeline-summary` already tells the reader "9 with broken stamps, **drawn
as breaks**" — while the renderer draws them as ordinary bars, because it branches on
`row.waitMinutes < 0 || row.workMinutes < 0` (`web/board-timeline.js:1387`, `:1424`) and no row on this board
has a negative span. Routing an unresolvable terminal row to the break marker makes that sentence true and
gives the legend's "Broken stamps" swatch something to name.

## Detailed Requirements

1. `buildTimelineAggregate` reads the ticket's already-resolved completion instant rather than re-parsing
   `completed_at`. The board decides what a REQ's completion instant is in exactly one place; this consumes
   that verdict, as the file's own header already says it does for the anomaly verdict.
2. `WaitOpen` and `WorkOpen` mean **still running**. A REQ in a terminal state is never open; derive the
   terminal test from the board's existing status vocabulary (`isTerminalResolvedStatus` and the
   `cancelled` / `failed` statuses), never from a new list in this file.
3. A terminal REQ whose end instant cannot be resolved carries no open span. It is drawn as a **break**,
   which is what the summary and the legend already promise, and it keeps its anomaly reason in the row
   tooltip and the table.
4. A terminal REQ with a resolved end instant but no `claimed_at` draws its wait from `created_at` to that
   instant, with no work segment — `HasWork` stays false, and no work is invented.
5. `timelineRange` no longer pulls `RangeEnd` up to `now` for rows that are not actually open. A board whose
   newest work finished last week fits to last week plus the projection, not to now.
6. `#timeline-summary`'s "N still open, measured to the now-line" then counts only genuinely running spans.
   On this board that is 1 (REQ-325), not 26.

## Constraints

- Never fabricate an instant. `CompletionUnresolved` stays unresolved; the break marker is the honest
  drawing, not a guessed end. File mtime remains prohibited as a source (`model.go:28-31`).
- Do not restate the anomaly rule or the terminal-status rule in `timeline.go` — consume `model.go`'s
  verdicts. A second reader here becomes a second definition (REQ-219's lesson, in this module's prime).
- The renderer's break-marker branch must key on the same field the payload sets, so the two cannot drift
  into disagreeing about which rows are breaks (the rule `timelineRowSegments` already follows for
  reversed spans).

## Red-Green Proof

**RED prompt/case:** `queue-kanban generate --out DIR --repo-root .`, open Timeline, read
`#timeline-summary`; then count `rect.timeline-segment.is-open` and `rect.timeline-segment-broken` in the
rows SVG at Fit all.

**Why RED now:** the summary reads "317 REQs in the window … **26 still open**, measured to the now-line at
2026-08-23 11:13 UTC. 9 with broken stamps, drawn as breaks." 25 of those 26 are `completed` or `cancelled`;
the `timeline-segment-broken` count is **0**.

**GREEN when:**
- The summary reads "1 still open" on this board, and the one is REQ-325.
- No row whose status is `completed`, `cancelled` or `failed` carries `waitOpen` or `workOpen` in the payload.
- REQ-311's wait ends at `2026-08-21T15:20:33Z` — the instant the calendar already places it on — not at the
  now-line.
- REQ-051 … REQ-059 each draw a break marker; the `timeline-segment-broken` count is 9, matching the number
  the summary states.
- `rangeEnd` in the payload is the latest real instant plus the projection, not `now`, on a board with no
  running work.
- A Go test in `timeline_test.go` pins each of the four row shapes (running, terminal-with-instant,
  terminal-without-instant, never-claimed-and-still-pending) against a fixture, and fails in both directions.

**Validation:** Inferred during capture from the generated payload; every number above was read out of
`board-data.js` rather than estimated.

## Full Context

See `do-work/user-requests/UR-066/input.md`.
