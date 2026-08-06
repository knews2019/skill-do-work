---
id: REQ-117
title: An unrecognized domain must leave a footprint on the board, not become general in silence
status: pending
created_at: 2026-08-06T10:53:03Z
user_request: UR-024
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: []
write_set: [tools/queue-kanban/model.go, tools/queue-kanban/model_test.go]
maintenance: false
related: [REQ-116, REQ-118]
batch: schema-contract-board-fixes
---

# An Unrecognized `domain` Must Leave a Footprint on the Board, Not Become `general` in Silence

## What

Since REQ-111 wired `domain` through `resolveSchemaField`, a typo'd `domain: quantum` renders on the board as a `general` badge with no warning anywhere — the recognized flag is discarded at the call site in `tools/queue-kanban/model.go`. Before that change the typo was at least visible verbatim, so the typo is now *better hidden* than it was. Add the flag and the data warning that `testing_status` already has, keep the `general` fallback, and delete the comment claiming the board has no channel for it.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

The comment defending the silence says the board "has no warning channel for it." That is false, and the disproof is one field away: `board.Warnings` is declared in `model.go` as "data-shape warnings (e.g. duplicate ids, unrecognized statuses, future-dated stamps) — surfaced, never silently dropped", it is fed by `collectTestingWarnings`, and `testing_status` additionally carries an unrecognized flag through to the card. The Schema Read Contract's clause 3 is literally titled **Never silently drop** — the warning is described there as "the missing feedback channel that allowed `dependencies:` to go unnoticed pre-0.76.2" — and the `domain` row prescribes the default *and* the warning, not the default alone.

## Detailed Requirements

- Add a `DomainUnrecognized` (or equivalently named) flag on the request ticket, set when a **present** `domain` value survives normalization outside the canonical enum. Mirror the existing `testing_status` shape rather than inventing a second convention.
- Raise one `board.Warnings` entry per affected ticket, worded with the existing `schemaFieldWarningText("domain", raw)` formatter — do not hand-type a second copy of the contract's warning text.
- **Keep the `general` fallback.** The contract's `domain` row names `general` as the default on unknown; this REQ adds the footprint, it does not change the resolved value.
- **Keep the absent-domain guard.** An absent `domain` must stay absent and unflagged — `web/board.js` gates the badge and the filter dropdown on `if (request.domain)`, and REQ-111's `TestParseRequestTicketPreservesAbsentDomain` exists to hold that line. An absent field is not a violation (`resolveSchemaField` returns recognized=true for it), so it must raise no warning.
- Replace the false comment. Whatever remains must not claim a missing channel; if a reason for a design choice is stated, it has to be a true one.
- Decide and state whether the flag reaches the card payload (`generate.go` + `web/board.js`) or stays server-side feeding only the warnings list. Either is acceptable; the warning is the required half, the visual flag is the optional half. If the flag does reach the frontend, it is display-only — **no column logic, no filter logic, no scheduling.**

## Constraints

- Read-side only; no write path, no new write surface (see the UR's Batch Constraints).
- The `domain` badge and filter behavior for *recognized* values must not change — this REQ only adds a signal for the unrecognized case.
- Shares `parseRequestTicket` and `model_test.go` with REQ-116 — serial execution, one checkout.
- If the flag is surfaced in the payload, the parser-to-schema lock-step rule applies (`CLAUDE.md` → Shipped Tooling): the board's parsed-fields contract and `model.go` change in the same commit.

## Dependencies

None. Independent of REQ-116, though both edit the same function.

## Builder Guidance

**Firm on the requirement, open on the surface.** That an unrecognized `domain` must produce a warning is settled by the contract. Whether the unrecognized flag also reaches the card as a visual marker is the builder's call — pick one, implement it, and say which in the Implementation Summary.

## Red-Green Proof

**RED prompt/case:** Parse a REQ whose frontmatter says `domain: quantum` and assert two things: the board's `Warnings` list contains an entry naming `domain` and `quantum`, and the ticket is flagged unrecognized. Both fail today — the warnings list is empty and no flag exists, while `ticket.Domain` is silently `general`.
**Why RED now:** `model.go`'s domain read site discards `resolveSchemaField`'s second return value (`normalizedDomain, _ = …`), so the recognized/unrecognized distinction is thrown away at the moment it is computed.
**GREEN when:** `domain: quantum` still resolves to `general` for display, but the board reports the contract's warning line for it; `domain: back-end` still resolves to `backend` silently (a recognized alias must stay quiet); and a REQ with no `domain` field produces no badge, no filter entry, and no warning.
**Validation:** User confirmed — the remedy was stated in the user's own capture text, which followed a triage report they accepted.

## Assets

None.

---
*Source: `do-work/user-requests/UR-024/input.md` — Finding 2 of the 0.174.15-series feedback triage: "An unrecognized domain is silently remapped to general with no footprint — add a DomainUnrecognized flag and a board.Warnings entry mirroring testing_status, and fix the false 'no warning channel' comment at model.go:619"*
