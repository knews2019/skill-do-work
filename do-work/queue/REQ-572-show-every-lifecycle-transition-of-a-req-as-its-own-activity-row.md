---
id: REQ-572
title: 'Show every lifecycle transition of a REQ as its own Activity row'
status: pending
created_at: 2026-09-04T23:16:00Z
user_request: UR-115
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
related: [REQ-573]
batch: activity-history
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/activity.go
  - skills/do-work-board/tools/queue-kanban/activity_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/web/board-activity.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board.css
---

# Show Every Lifecycle Transition of a REQ as Its Own Activity Row

## What

The Activity view lists one row per REQ: its newest lifecycle stamp and nothing else. Change the aggregation so every parseable lifecycle stamp inside the window becomes a row, so a REQ that was captured, claimed, dispatched, merged, reviewed, released and completed in the last 24 hours shows all of those transitions, newest first, on the same surface. The Board's detail drawer only prints Created, Claimed and Completed, and the Timeline only draws two spans, so today the full path of a REQ is readable only in its frontmatter or with `git log`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- `buildActivityRows` (`activity.go`) emits one `ActivityRow` per parseable stamp on each ticket, not one per ticket. `newestLifecycleStamp` goes away or becomes the single-row special case; do not add a second stamp list, `lifecycleTimestampFields` in `model.go` stays the one enumeration.
- Ordering stays newest first with the same deterministic tie-break, so several rows of one REQ interleave with other REQs by time.
- The client (`board-activity.js`) keeps windowing against the wall clock and keeps the shared filter chips. The summary line must say what it now counts: rows are transitions, REQs are distinct ids, and both numbers are useful ("38 transitions across 21 REQs in the last 24 hours").
- A ticket with no parseable stamp is still skipped, never dated from the zero time.
- Prefer showing all transitions by default over adding a "latest only" toggle. If a toggle is kept, "every transition" is the default state and the toggle is one extra button in the existing Activity window group, not a new control family.
- Click behavior and the Board's `data-detail-kind` attribute are REQ-573's concern (opening the drawer and highlighting sibling rows); this REQ only changes the row set and the counts.

## Constraints

- Go decides which stamps exist and what each records; the client draws and filters. No second definition of a stamp's meaning in JavaScript.
- The existing window, filter chips, and empty-state messages keep working; only the row set and the counts change.
- REQ-571 (removing the board's `pending-heavy-testing` reader case) may touch `model.go` in the same period; this REQ does not edit that file.

## Dependencies

REQ-573 (click a row to open the drawer and highlight sibling rows) depends on this REQ.

## Red-Green Proof
**RED prompt/case:** A ticket with `created_at: 2026-09-04T22:52:00Z` and `claimed_at: 2026-09-04T23:00:17Z` passed to `buildActivityRows` returns one row ("claimed"). On the running board, REQ-570 (deleting the pending-heavy-testing status) in the 24h Activity window shows a single "claimed" row; its capture eight minutes earlier is not visible.
**Why RED now:** `newestLifecycleStamp` keeps only the latest stamp per ticket by design (REQ-568, showing recently touched REQs regardless of status).
**GREEN when:** The same ticket returns two rows, "claimed" at 23:00:17 then "captured" at 22:52:00, in that order; on the board REQ-570 appears twice in the 24h window and the summary line reports both the transition count and the distinct REQ count.
**Validation:** User confirmed (verify-requests, 2026-09-04)

## Builder Guidance

The user is certain about the outcome (see every state a REQ went through) and left the shape to the builder. The Go change is small; most of the judgment is in how the table reads when one REQ has six rows among others. Keep the row shape (`id`, `stampField`, `stampAt`, `transition`) so the payload contract and REQ-573 stay stable.

## Assets

- `do-work/user-requests/UR-115/assets/REQ-572-screenshot-1-activity-view-one-row-per-req.png`: the Activity view at 24h, "38 REQs touched in the last 24 hours", columns REQ / Title / Status / What happened / When / Stamp. REQ-570 has one row: status claimed, "claimed", Sep 4 23:00 UTC, stamp `claimed_at`. REQ-505 above it shows "completed" with stamp `completed_at`; REQ-571 shows "captured" with stamp `created_at`; REQ-507 and REQ-506 show "status changed to pending" with stamp `status_changed_at`. A verify finding card sits above the table. The browser find box has "570" highlighted, 1 of 4 matches, all in the finding card and the one row.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`, so bare only): matches on "Changing queue-kanban model, parser, UI, timeline, testing, or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban parsing, views, static output". Over the budget on its own.

## Full Context
See `do-work/user-requests/UR-115/input.md` for complete verbatim input.

*Source: "is this only showing the last status of a REQ? how about if I want to see when it went through all the states of it?" / "ok, do-work capture-request"*
