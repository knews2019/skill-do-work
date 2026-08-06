---
id: REQ-119
title: An off-vocabulary route warns on the board like domain does
status: pending
created_at: 2026-08-06T11:26:17Z
user_request: UR-025
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: []
write_set: [tools/queue-kanban/model.go, tools/queue-kanban/model_test.go]
maintenance: false
related: [REQ-116, REQ-117, REQ-120]
batch: codex-pr133-findings
---

# An Off-Vocabulary `route` Warns on the Board Like `domain` Does

## What

REQ-116 made the board normalize `route` and REQ-117 gave `domain` an unrecognized flag plus a `board.Warnings` entry. `route` got the normalization and not the warning, so `route: z` now reaches the card as `Z` with no footprint anywhere — the same silence REQ-117 was written to remove, one field over.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why (if provided)

Two reasons it is worth closing now rather than living as a known asymmetry. The Schema Read Contract's clause 3 (`actions/work-reference.md` → Schema Read Contract, "Never silently drop") applies to every enum read, not to whichever field a REQ happened to name. And the board is now internally inconsistent in a way a reader would reasonably call a bug: `domain: quantum` warns, `route: z` does not, for no stated reason.

## Detailed Requirements

- Add an unrecognized flag and the raw value to the ticket (`RouteUnrecognized`, `OriginalRoute`), mirroring the `Domain`/`OriginalDomain`/`DomainUnrecognized` trio REQ-117 added.
- Raise the warning through the existing channel, using `schemaFieldWarningText("route", raw)` so the contract's phrasing stays in one place. Either extend `collectDomainWarnings` into one collector over both fields or add a sibling — the builder's call, but there must not be two hand-typed copies of the warning wording.
- **The resolved value must not change.** `route: z` still reaches the card as `Z`, not blanked and not defaulted — REQ-116 chose that deliberately (route's documented default is the empty string, so substituting it would make an unrecognized letter indistinguishable from an absent field, destroying the evidence re-triage reads). This REQ adds the footprint only.
- An absent `route` stays absent and unflagged, exactly as for `domain`; a canonical or aliased letter (`a` → `A`) resolves silently. Both need a test — a channel that fires on ordinary REQs is one readers learn to ignore.
- `schemaFieldWarningText` already handles route's no-default case ("No default is defined; reporting it unchanged."), so the warning text should need no new branch. Verify that rather than assuming it.

## Constraints

- Read-side only; no write path, no new write surface.
- Do not widen the change to the five contract fields the board does not read (`caveman`, `maintenance`, `tdd`, `error_type`, `kb_status`) — REQ-116 established they have no display role, and giving them one to satisfy a symmetry argument is backwards.

## Dependencies

None. Builds on REQ-116 (route normalization) and REQ-117 (the domain warning pattern), both archived.

## Builder Guidance

**Firm.** The pattern to mirror is in the same function and shipped one commit ago. The only judgment is one collector versus two.

## Red-Green Proof

**RED prompt/case:** Parse a REQ whose frontmatter says `route: z` through `buildBoard` and assert the board's `Warnings` list contains an entry naming `route` and `z`, and that the ticket is flagged unrecognized. Both fail today — the warnings list is empty and no flag exists.
**Why RED now:** `model.go`'s route read site calls `normalizeSchemaField` and keeps no recognition result, because REQ-116 wired normalization before REQ-117 established the warning channel for this class of field.
**GREEN when:** `route: z` still resolves to `Z` for display *and* raises the contract's warning; `route: a` resolves to `A` in silence; an absent `route` produces no badge, no drawer row, and no warning.
**Validation:** Inferred during capture — derived from Codex's P2 finding on PR #133 and verified against the code, not confirmed field-by-field with the user.

## Assets

None.

---
*Source: `do-work/user-requests/UR-025/input.md` — Codex P2 finding on PR #133: "Warn when a route remains outside the enum"*
