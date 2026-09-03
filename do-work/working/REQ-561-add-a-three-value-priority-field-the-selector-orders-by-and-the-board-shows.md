---
id: REQ-561
title: 'Add a three-value priority field the selector orders by and the board shows'
status: claimed
created_at: 2026-09-03T20:38:55Z
user_request: UR-107
domain: backend
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-kanban-board.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-530]
write_set:
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/capture-reference.md
  - skills/do-work/actions/capture.md
  - skills/do-work/actions/work.md
  - skills/do-work/tools/do-work-cli/internal/schemanormalization/
  - skills/do-work/tools/do-work-cli/internal/nextselection/
  - skills/do-work-board/tools/queue-kanban/model.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - do-work/queue/
claimed_at: 2026-09-03T20:58:09Z
---

# Add a Three-Value Priority Field the Selector Orders By and the Board Shows

## What

Add one optional frontmatter field to the REQ schema, `priority`, with the closed set `now`, `next`, `later`; absent reads as `next`, anything else normalizes to `next` with a warning under the Schema Read Contract. The selector orders ready work inside its ordinary class by priority (`now` before `next` before `later`) and keeps the existing queue order inside each value; the gate-repair and deferred-parent classes stay above it and `depends_on` stays the hard gate. The board sorts the pending column by priority inside its ready and waiting groups and shows a small `now` or `later` tag on the card (no tag for `next`, the default), the way impact tags are shown today. Capture may set the field from the user's words at Step 1 and addenda may change it. The landing commit stamps the current queue per the 23:20 triage table in `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/index.html`: the build-now set `now`, the deferred set `later`, the rest untouched.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

"can we make a priority list that would reflect in the kanban board as well". A targeted run keeps caller order but lives in a shell command; the board cannot show it and the next session cannot see it.

## Context

- Selection order today (`work-reference.md`, Selection Order): `repository_gate_repair: true` first, ready `gate_deferred: true` parents second, ordinary work third, existing queue order inside each class. Priority sorts inside the third class only.
- The board's pending column already splits ready from waiting on dependencies and shows impact tags on cards; priority is one more sort key and one more tag, not a new column.
- REQ-530 (order ready work by the newest REQ it unblocks) is superseded: with a priority field the maintainer says what is urgent instead of the selector inferring it from REQ numbers. Its rule is not adopted as the tie-break; existing queue order stays the tie-break so the change is one field and nothing else. Cancel REQ-530 with this REQ's landing hash.
- The triage table in the velocity report (commit a524cf16 and later) is the first priority list: 28 REQs to build now (the REQ-502 chain, REQ-559, REQ-560, REQ-542, REQ-539, REQ-475, REQ-483, REQ-496, REQ-512, REQ-544, REQ-514, REQ-515, REQ-485, REQ-527, REQ-534, REQ-547, REQ-490, REQ-535, REQ-536, REQ-545), 13 deferred (REQ-549 to REQ-558, REQ-482, REQ-486, and REQ-530 until cancelled). Read the report's current table at build time; it may have moved.

## Detailed Requirements

- Schema: `priority` documented in `work-reference.md` Request File Schema and the Schema Read Contract table (closed set, default `next`, normalization warning), and in `capture-reference.md`'s Simple REQ template as a commented optional line; `capture.md` Step 1 gets a one-line priority assessment: emit only when the user's words rank the work, never invent.
- do-work-cli: `schemanormalization` reads and normalizes the field; `nextselection` orders the ordinary class by it and projects the value on selected and excluded records so a caller never infers it from display text. Table-driven tests: default absent, each value, an invalid value, and a `later` REQ that is the only ready dependency of a `now` REQ (the dependency still runs first).
- Board: `model.go` parses the field (lock-step with the parser rule in prime-kanban-board.md), `generate.go` sorts the pending column by priority inside ready and waiting, `board-cards.js` renders the tag, `board.css` styles it in both themes. A model test pins the sort; the static snapshot and the live board agree.
- Stamp: the landing commit edits `priority:` on the pending REQ files named by the triage table; those are capture-editable pending files, staged by explicit path, never a claimed file.
- Version bump and changelog entry; the changelog title says what shipped.

## Constraints

- One integrating commit; the schema, selector and board land together so a stamped queue is never malformed to any reader.
- `depends_on` is never overridden by priority; a `now` REQ behind an unfinished dependency waits.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.

## Red-Green Proof
**RED prompt/case:** Put `priority: now` on a pending REQ, run `do-work-cli next --format json` and regenerate the board.
**Why RED now:** the field is unknown: the selector ignores it and orders by queue number, and the board shows no tag and no ordering.
**GREEN when:** `next` lists the `now` REQ before other ordinary ready work and projects `priority` on the record; a `now` REQ behind an unmet dependency is still excluded; the board's pending column shows the `now` REQ at the top of its ready group with a tag; `later` REQs sink to the bottom of their group; and after the stamp the board's pending column reads in the triage table's order.
**Validation:** Inferred during capture; the maintainer approved the capture ("capture UR-107").

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5562 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes the board model and cards.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5660 tokens, over budget and `slugged: partial`. Matched because this REQ changes schema normalization and selection.
- `_dev/primes/lessons-action-files.md` — 4050 tokens, over budget and `slugged: partial`. Matched because this REQ adds a schema field read by several actions.

## Full Context
See `do-work/user-requests/UR-107/input.md` for complete verbatim input.
