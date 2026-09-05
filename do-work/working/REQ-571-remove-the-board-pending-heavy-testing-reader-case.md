---
id: REQ-571
title: '[impact-negligible] Remove the board''s pending-heavy-testing reader case'
status: claimed
created_at: 2026-09-04T22:52:00Z
user_request: UR-114
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-570]
related: [REQ-570]
batch: orchestrator-simplification
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/actions/board.md
  - skills/do-work-board/docs/board-guide.md
  - skills/do-work-board/tools/queue-kanban/model.go
  - skills/do-work-board/tools/queue-kanban/model_test.go
  - skills/do-work-board/tools/queue-kanban/timeline.go
  - skills/do-work-board/tools/queue-kanban/timeline_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-calendar.js
  - skills/do-work-board/tools/queue-kanban/prime-do-kanban.md
claimed_at: 2026-09-05T00:38:07Z
route: A
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-09-05T00:50:53Z
  basis:
    - Route A
    - 8-file write set
    - 2 subsystems involved
    - 3 acceptance criteria
---

# Remove the Board's pending-heavy-testing Reader Case

## What

Once REQ-570 stops writing `pending-heavy-testing`, delete the board's handling of that value: the model's status classification, the timeline's held-state events, the calendar script's label, and the sentences in the board action, guide and prime that describe it. A held request is `claimed` and renders as in progress.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The board is a separate package with its own version and a parser that must stay in lock-step with the queue's schema (`_dev/primes/prime-kanban-board.md`). After REQ-570 the value is dead in every writer, so its reader case is dead code that would keep the "held" vocabulary alive in the board's help text and tests.

## Context

- Readers found by search on 2026-09-04: `model.go`, `timeline.go`, `web/board-calendar.js`, `actions/board.md`, `docs/board-guide.md`, `prime-do-kanban.md`, `lessons-do-kanban.md`, and the tests `activity_test.go`, `board_live_test.go`, `frontmatter_cli_test.go`, `generate_test.go`, `javascript_behavior_c_test.go`, `model_test.go`, `timeline_browser_probe_test.go`, `timeline_test.go`. Confirm by search at claim time.
- Historical timeline events for archived requests may still contain the value in their recorded history; the timeline must keep rendering those records without a special label. Deleting the case must not make an old record read as an error.
- Release 0.275.3 made a held request wait in the Pending column instead of the operator's inbox; that placement rule goes with the status.

## Detailed Requirements

- Delete the status from the board's status enum, column placement, timeline event classification, and the calendar script's label map.
- Delete the sentences describing the held state from `actions/board.md`, `docs/board-guide.md` and `prime-do-kanban.md`; leave `lessons-do-kanban.md` untouched, as lessons are history.
- Update the affected tests to assert that an unknown or historical value renders without error and that a claimed request with a `## Heavy Verification Plan` section renders as in progress.
- Bump the board version and record the parser change per `_dev/primes/prime-kanban-board.md`.

## Constraints

- Depends on REQ-570; do not start while any queued or working request still carries the value.
- Tolerant reading of historical records stays; only the named case and its vocabulary go.

## Dependencies

- Depends on REQ-570 (core skill deletes the status).

## Builder Guidance

- Mechanical deletion. If a test only exists to pin the held state, delete the test rather than rewrite it.

## Red-Green Proof
**RED prompt/case:** Generate the board against a fixture queue containing a `claimed` request with a `## Heavy Verification Plan` section and an archived request whose history mentions `pending-heavy-testing`.
**Why RED now:** The model still classifies and labels the held status, and the board's help text still describes it.
**GREEN when:** The claimed request renders as in progress; the archived record renders without a held label or error; `grep -r pending-heavy-testing skills/do-work-board` matches only `lessons-do-kanban.md`; the board test package passes.
**Validation:** Inferred during capture.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5744 tokens, over the 2000-token budget and `slugged: partial`; matched on "changing queue-kanban model, parser, UI, timeline".
- `_dev/primes/lessons-kanban-board.md` — 4820 tokens, over budget and `slugged: partial`; matched on "changing queue-kanban parsing, views".

## Full Context
See `do-work/user-requests/UR-114/input.md` for complete verbatim input.

---

## Triage

**Route: A** - Simple

**Reasoning:** The request is a cleanup after REQ-570 removed a status, with its own eight-file write set declared. The reader case to remove and the tests that move with it are both named.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
