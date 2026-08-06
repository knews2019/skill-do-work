---
id: REQ-116
title: Normalize route at the board's read site and correct 0.174.15's board-wide claim
status: pending
created_at: 2026-08-06T10:53:03Z
user_request: UR-024
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: []
write_set: [tools/queue-kanban/model.go, tools/queue-kanban/model_test.go]
maintenance: false
related: [REQ-117, REQ-118]
batch: schema-contract-board-fixes
---

# Normalize `route` at the Board's Read Site and Correct 0.174.15's Board-Wide Claim

## What

The board parses `route` verbatim (`tools/queue-kanban/model.go` — `Route: coerceScalarToString(fields["route"])`), so a REQ written `route: a` reaches the card as lowercase `a`. Wire that read through the normalizer REQ-111 already added, and correct the 0.174.15 changelog entry's claim that the board honors the Schema Read Contract for all nine fields — it does not, and only `domain` was ever wired.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

`route` is not an invisible field: it is serialized into the board payload (`tools/queue-kanban/generate.go` — the `route` JSON tag), rendered as a card badge (`tools/queue-kanban/web/board.js` — `makeBadge("badge-route", …)`), and shown as a drawer metadata row. So the un-normalized value is on screen. It is also the one field of the seven added in REQ-111 that the board actually reads — the other five (`caveman`, `maintenance`, `tdd`, `error_type`, `kb_status`) have no board read site at all, which is why the changelog's claim is vacuous for them and false for `route`.

## Detailed Requirements

- Route the board's `route` read through `normalizeSchemaField("route", …)`. Uppercase-only: the contract's row is `A | B | C` with lowercase `a`/`b`/`c` → uppercase and **no documented default** ("treat as needing re-triage in Step 3"), so nothing may be substituted for an unrecognized value.
- **Do not use `resolveSchemaField` here.** It substitutes the field's default, which for `route` is the empty string — so an unrecognized `route: Z` would reach the card as blank and lose the evidence a re-triage needs. Normalize, then pass through whatever survives.
- Apply the same present-value-only guard `domain` uses: an absent `route` must stay absent, because `web/board.js` gates the badge and the drawer row on `if (request.route)`.
- Add a **parse-level** test in `tools/queue-kanban/model_test.go` — one that goes through `parseRequestTicket` and asserts the ticket's `Route` field, not one that calls `normalizeSchemaField` directly. The existing cases at `model_test.go` (the `{"route", "a", "A"}` table rows) already cover the library function; they passed while the board stayed broken, which is precisely the gap.
- Correct the changelog record. Per the user's instruction this goes in **the next entry**, not by rewriting 0.174.15 in place — the same shape 0.174.13 used to correct 0.174.12's recovered-trap evidence. The new entry must say plainly that 0.174.15 wired only `domain` to the board, and that this release adds `route`.
- Leave the five unread fields unread. Adding board read sites for `caveman`, `maintenance`, `tdd`, `error_type` or `kb_status` is **out of scope** — they have no display or column role, and inventing one to make an old changelog title true is backwards.

## Constraints

- Read-side only; no write path, no new subcommand, no third write surface (see the UR's Batch Constraints).
- Pure Go inside `tools/queue-kanban/`, plus the `CHANGELOG.md` entry that Step 9's Before-Every-Commit ritual writes anyway. No `actions/` prose change is expected: `actions/work-reference.md`'s `route` row already prescribes this behavior. Its read-site column lists only work-pipeline steps and not the board — leaving that column alone is acceptable here (the contract's own rule is that the *condition*, not the caller list, is the trigger); widening it is a judgment call for the builder, and if made, the prose and `model.go` ship in the same commit.
- Shares `parseRequestTicket` and `model_test.go` with REQ-117 — serial execution, one checkout.

## Dependencies

None. REQ-111 (archived) already added `normalizeSchemaField` and the `route` contract row; this REQ only wires the existing function to the existing read site.

## Builder Guidance

**Firm.** The enum, the alias rule and the no-default decision are fixed by the contract, and the call site is named. The only judgment is whether to widen the contract table's read-site column for `route`, and how to word the changelog correction.

## Red-Green Proof

**RED prompt/case:** In `tools/queue-kanban/model_test.go`, parse a REQ whose frontmatter says `route: a` through `parseRequestTicket` and assert `ticket.Route == "A"`. It fails today, returning `a`.
**Why RED now:** `model.go`'s read site is `Route: coerceScalarToString(fields["route"])`, and `coerceScalarToString` only trims whitespace — it never case-folds. The normalizer exists but has no caller on this path.
**GREEN when:** `route: a` parses to `A`; `route: B` stays `B`; an absent `route` still parses to the empty string (no badge, no drawer row); and an unrecognized `route: Z` is reported as `Z` rather than blanked.
**Validation:** User confirmed — the remedy was stated in the user's own capture text, which followed a triage report they accepted.

## Assets

None.

---
*Source: `do-work/user-requests/UR-024/input.md` — Finding 1 of the 0.174.15-series feedback triage: "The board reads route verbatim at model.go:655 — wire it through normalizeSchemaField (uppercase-only, no default) with a parse-level test, and correct 0.174.15's overstated board-wide claim in the next changelog entry"*
