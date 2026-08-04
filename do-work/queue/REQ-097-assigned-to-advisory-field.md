---
id: REQ-097
title: assigned_to advisory field — schema line, scan skip-and-report, board parse (lock-step)
status: pending
created_at: 2026-08-04T19:44:17Z
user_request: UR-018
domain: backend
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-096]
maintenance: false
related: [REQ-096, REQ-098]
batch: parallel-building
write_set: [actions/work-reference.md, actions/work.md, tools/queue-kanban/model.go, tools/queue-kanban/model_test.go, tools/queue-kanban/web/*]
---

# `assigned_to` Advisory Field — Schema, Scan Skip, Board Parse

## What

Add the cooperative claim marker the whole model rests on: a single advisory frontmatter field, `assigned_to: "<session-name>"`, on a pending REQ. The default work-loop scan skips-and-reports such REQs; explicitly targeting one (`do-work run REQ-NNN`) overrides and clears the field. The board parses it display-only. **No verb, no status, no staleness clock, no release command** — the 0.163.0 forbidden-token ratchet stays fully intact.

## Detailed Requirements

- **Schema:** add `assigned_to` to `actions/work-reference.md`'s frontmatter block and Schema Read Contract — in the **verbatim-read class** alongside `write_set` (`:206`): no alias map, no normalization, no canonical vocabulary. Optional on every REQ; capture may seed it when the user earmarks work for a named session.
- **Scan behavior:** one skip-and-report sentence in `actions/work.md` Step 1's default selection: a pending REQ with `assigned_to` is skipped and listed ("assigned to <name>") exactly like the existing dependency-skip reporting. Explicit targeting (`do-work run REQ-NNN`) overrides the skip and **clears the field** as part of the claim.
- **Board (same commit — lock-step rule):** `tools/queue-kanban/model.go` parses `assigned_to` display-only: badge on the card + drawer metadata row, **no column logic, no scheduling**, same class as `write_set`. ~15 lines + test in `model_test.go`.
- **No `assigned_at`, no staleness threshold, no auto-release** — an assignment persists until cleared by an explicit run or hand-edit.
- **Ratchet check:** confirm `assigned_to` does not trip `_dev/tests/contract-regressions.sh`'s reservation token patterns (`reserved_for`, `reserved_at`, underscore-token forms); run the suite.
- Mirror the schema note in the co-location sentence style the Testing placeholders use (CLAUDE.md "keep the parser in lock-step" — restate inline in shipped files, never cite CLAUDE.md).

## Constraints

- Display-only at any builder count — nothing schedules, gates, or dispatches on `assigned_to` except the Step 1 skip (which is a *courtesy read*, not a gate: explicit targeting overrides).
- The forbidden-token ratchet must stay green with no test weakening.

## Red-Green Proof

**RED prompt/case:** `go test ./...` in `tools/queue-kanban/` with a new test asserting a REQ carrying `assigned_to: "cloud-alpha"` surfaces the value in the parsed model — fails before the parser change. Prose half: today no shipped file defines `assigned_to`, so a second session has no sanctioned way to see "this is earmarked".
**Why RED now:** The field does not exist; reserve (its predecessor) was deleted at 0.163.0.
**GREEN when:** The Go test passes; `actions/work.md` Step 1 documents skip-and-report + override-and-clear; contract-regressions suite green.
**Validation:** User confirmed (ask-tool answer: "assigned_to field only").

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 2, items 5–6).

---
*Source: approved plan, Phase 2*
