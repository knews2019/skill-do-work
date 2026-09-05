---
id: REQ-576
title: 'Start the board card wall time at the earliest lifecycle stamp, not only claimed_at'
status: claimed
created_at: 2026-09-04T23:52:00Z
user_request: UR-116
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
related: [REQ-575, REQ-572, REQ-448]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/durations.go
  - skills/do-work-board/tools/queue-kanban/durations_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
claimed_at: 2026-09-05T00:38:08Z
---

# Start the Board Card Wall Time at the Earliest Lifecycle Stamp, Not Only `claimed_at`

## What

Change the completed card's "wall time" so its origin is the earliest parseable lifecycle stamp the REQ carries other than `created_at`, and its end stays `completed_at`. Today `measureImplementationSpan` in `durations.go` reads only `claimed_at`, so a request whose claim stamp was rewritten late reports a span that excludes all of its phase work.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

REQ-505 (moving selection and claim behind `advance`) carries `planning_at` 16:49, `dispatch_at` 16:52, `builder_handback_at` 17:24, `integration_at` 17:26, a rewritten `claimed_at` 23:00 and `completed_at` 23:01 (all 2026-09-04 UTC). The card says "wall time 1m 23s"; the drawer shows the Planning row at "-6h 10m wall since Claimed". Measuring from the earliest stamp would have shown about 6h 12m, which is the real span, using data the board already parses. REQ-575 (keeping every lifecycle stamp) stops the damage at the writer; this REQ makes the card read the record it has.

## Context

- `durations.go` `measureImplementationSpan` returns `WallMinutes = completed_at - claimed_at`, a `StampsParsed` flag, and an exclusion reason (`paused` over the ceiling, `reversed` when negative). `generate.go` copies `WallMinutes` into `implementationSpanMinutes`; `web/board-cards.js` prints it as "wall time".
- `model.go` `lifecycleTimestampFields` is the one enumeration of stamp fields. Read the origin from it; do not inline a second list.
- The Durations view's calibration span (estimate against `claimed_at` to `completed_at`) is a separate reading with its own comment saying so. Leave it unless the same helper change is a strict improvement there too; if kept separate, the card comment must say the two spans differ and why.

## Detailed Requirements

- Origin = the minimum parseable instant among the REQ's lifecycle stamps, excluding `created_at` (queue wait is not work) and `completed_at`/`release_at` (they are ends). Key this on the enumeration in `model.go`, filtered by role, not on a hand-written field list.
- A REQ with `claimed_at` as its earliest stamp behaves exactly as today, so historical cards do not change.
- A REQ with only `claimed_at` and `completed_at` still measures; a REQ whose stamps do not parse still reports no span (`StampsParsed` false), never zero.
- The `reversed` exclusion keeps firing when the origin is after `completed_at`.
- The card comment in `board-cards.js` and the Go comment on the helper say what the span now measures: earliest recorded lifecycle stamp to completion.

## Constraints

- No new frontmatter field and no timing file. Read only what the REQ already carries.
- Version, changelog and mirror handling follow `_dev/primes/prime-kanban-board.md`.

## Builder Guidance

Certainty: firm on the origin rule and on excluding `created_at` (user confirmed at verify). Exploratory on whether the Durations view calibration span should adopt the same origin; the default is to leave it and document the difference.

## Red-Green Proof

**RED prompt/case:** A ticket fixture with `planning_at: 2026-09-04T16:49:45Z`, `claimed_at: 2026-09-04T23:00:06Z`, `completed_at: 2026-09-04T23:01:29Z` passed to `measureImplementationSpan` returns `WallMinutes` of about 1.38 (1m 23s). On the running board, the REQ-505 card reads "wall time 1m 23s".
**Why RED now:** The helper reads `claimed_at` only.
**GREEN when:** The same fixture returns about 371.7 minutes (6h 11m 44s), a fixture with `claimed_at` earlier than every phase stamp returns the same value as before the change, and the rebuilt board shows REQ-505 at about "wall time 6h 11m".
**Validation:** User confirmed (verify-requests, 2026-09-05)

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5744 tokens, over the 2000-token budget and `slugged: partial`; matched because the change edits the queue-kanban duration model and card UI.
- `_dev/primes/lessons-kanban-board.md` — 4820 tokens, over budget and `slugged: partial`; matched because the change edits queue-kanban parsing consumers and static output.

## Full Context

See `do-work/user-requests/UR-116/input.md` for the verbatim input and the REQ-505 trace.

---
*Source: "capture a req for append-only stamps and the board wall time change"*
