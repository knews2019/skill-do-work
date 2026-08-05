---
id: REQ-107
title: Board assigned-badge comment over-claims — "nothing trims" sits directly above truncateBadgeText
status: pending
created_at: 2026-08-05T09:43:47Z
user_request: UR-019
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set: [tools/queue-kanban/web/board.js]
related: [REQ-097, REQ-105, REQ-106]
batch: sync-review-0174
---

# Board Assigned-Badge Comment Over-Claims

## What

In `tools/queue-kanban/web/board.js`, the assigned-badge comment says the value is "Rendered verbatim — session names have no canonical vocabulary, so nothing here folds, trims, or rewrites the value" — directly above `truncateBadgeText(request.assignedTo, 18)`, which trims the badge text for display. The behavior is fine (the tooltip and drawer row carry the full value); the comment over-claims. Reword it to distinguish value normalization (none — verbatim-read class) from display truncation (badge clipped at 18 chars, full value in the tooltip and drawer).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Cosmetic, comment-only — no behavior change, no schema change, so the `model.go`/work-reference lock-step is untouched. Found by a downstream consumer's review of the 0.170.1 → 0.174.3 sync; verified here at triage (comment at `tools/queue-kanban/web/board.js:533-538`, truncation call at `:542`, full-value tooltip at `:546`, drawer row at `:1722-1723`).

## Red-Green Proof
**RED prompt/case:** Reading `board.js:535` ("nothing here folds, trims, or rewrites the value") against `:542` (`truncateBadgeText(request.assignedTo, 18)`) — the comment contradicts the very next statement.
**Why RED now:** The comment conflates "no value normalization" (true) with "no display trimming" (false).
**GREEN when:** The comment states both halves accurately: the *value* is never normalized (verbatim-read class), while the *badge display* is truncated at 18 chars with the full value preserved in the tooltip and drawer row.
**Validation:** Inferred during capture (triage-verified against the repo)

## Full Context
See `do-work/user-requests/UR-019/input.md` for complete verbatim input.

---
*Source: downstream sync-review finding 3, verified at triage — see UR-019*
