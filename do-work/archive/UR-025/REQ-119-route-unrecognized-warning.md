---
id: REQ-119
title: An off-vocabulary route warns on the board like domain does
status: completed
created_at: 2026-08-06T11:26:17Z
claimed_at: 2026-08-06T11:28:00Z
completed_at: 2026-08-06T11:32:00Z
commit: c327f24
route: A
kb_status: pending
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
- [x] **[PLAN]:** Mirror REQ-117's domain trio for route (`OriginalRoute`, `RouteUnrecognized`), and take the "one collector" branch of the REQ's open judgment: generalize `collectDomainWarnings` into `collectSchemaFieldWarnings` iterating a field/flag/original triple, so the warning wording stays in `schemaFieldWarningText` alone. Keep `normalizeSchemaField` at the read site (not `resolveSchemaField`) so `route: z` still resolves to `Z`. Derive the flag with `isKnownSchemaFieldValue` rather than a second `resolveSchemaField` call, which would discard the case-folded value.
- [x] **[APPLY]:** Tests first, RED at build (`ticket.RouteUnrecognized undefined`). Then two struct fields, the flag at the read site, the collector generalized, and its one call site renamed. Two files, both in `write_set`.
- [x] **[UNIFY]:** `git diff --stat` → 2 implementation files. `gofmt -l .` clean; `go vet ./...` clean; `go test ./...` green (4.4s), including main's `open_work` suite and REQ-117's domain tests, which pass unchanged through the collector rename. No debug artifacts. Verified `schemaFieldWarningText` needed no new branch, as the REQ asked: route's no-default row already produces "No default is defined; reporting it unchanged." Built the binary and ran `summary` against a synthetic tree holding both a bad route and a bad domain — `warnings: 2`, each naming its own field and written value.

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

## Triage

**Route: A** - Simple

**Reasoning:** The pattern shipped one commit ago in the same function; the REQ names the fields and the constraint. One judgment call, recorded as D-01.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

---

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/model.go` (modified)
- `tools/queue-kanban/model_test.go` (modified)

**What was done:** `RequestTicket` gained `OriginalRoute` and `RouteUnrecognized`; the route read site derives the flag with `isKnownSchemaFieldValue` while keeping the case-folded value. `collectDomainWarnings` became `collectSchemaFieldWarnings`, iterating both contract fields the board reads, so the contract's warning wording still exists in exactly one place.

## Decisions

- **D-01**: One collector over both fields, not a sibling `collectRouteWarnings`. **DECIDE & STATE** — the REQ left this open and both satisfy it.

  Reasoning: a sibling function would have been the second site formatting a contract warning, and the REQ's own constraint was that there must not be two hand-typed copies of the wording. Generalizing keeps one loop, one formatter call, and makes the next field an entry in a slice rather than a new function plus a new `append` in `buildBoard`. The cost is that the function name no longer says "domain", so its doc comment states the pairing rule explicitly: a field's flag and its `Original*` value are useless apart.

## Testing

**Tests run:** `go test ./...` (in `tools/queue-kanban/`), plus `go test -run 'Route|Domain' .`
**Result:** ✓ All passing

**Red-green validation:**
- `TestUnrecognizedRouteFlagsAndWarns`: ✗ before (build failure — `ticket.RouteUnrecognized undefined`) → ✓ after
- `TestRecognizedRouteRaisesNoWarning`: ✗ before (same build failure) → ✓ after — the silence guard for the alias and absent cases
- Regression, green before and after: REQ-117's `TestUnrecognizedDomainFlagsAndWarns` and `TestRecognizedDomainRaisesNoWarning` both pass through the collector rename, which is the evidence the generalization preserved domain's behaviour
- End-to-end: `summary` against a tree with `route: z` **and** `domain: quantum` → `warnings: 2`, route's reading `⚠ route: 'z' not recognized — expected one of [A, B, C]. No default is defined; reporting it unchanged.`

**New tests added:**
- `TestUnrecognizedRouteFlagsAndWarns`
- `TestRecognizedRouteRaisesNoWarning`

*Verified by work action*

## Lessons Learned

**What worked:** Taking the reviewer's finding at face value only after checking it. Codex's claim was that the board was out of lock-step with the contract for `route`, and reading the read site confirmed it exactly — normalization without a recognition result. The finding was also *structurally* predictable: REQ-116 and REQ-117 were captured from the same review round and split by field, so the first shipped normalization before the second had established the channel. Splitting one contract leg across two REQs leaves a window where the fields disagree.

**What didn't:** Nothing failed, but the first instinct — derive the flag with a second `resolveSchemaField` call for symmetry with domain — would have been wrong. For route that call returns the empty-string default, so the flag would have been right and the value destroyed. `isKnownSchemaFieldValue` on the already-normalized value is the correct pairing when a field has no default.

**Worth knowing:** The board's two contract-field warnings now share one collector, so adding a third field is a slice entry plus its read-site flag — but the two halves are useless apart, and nothing enforces that pairing except the doc comment. A field given a flag and no `Original*` value would warn naming an empty string; a field given `Original*` and no flag would never warn at all.

## Orientation

The board's data-warnings banner now covers `route` as well as `domain` — `route: z` still shows as `Z`, and the board says it is off-vocabulary instead of displaying it silently. Lives in the queue-kanban parser (`tools/queue-kanban/prime-do-kanban.md`), closing the asymmetry REQ-116/117 left between two fields of one contract leg. Not `[MAP CHANGED]`: one collector generalized, no new module or data flow. Prime staleness spot-check: `prime-do-kanban.md`'s referenced paths all still exist (it gained `open_work.go` from main's 0.176.0).

## Review

**Overall: 97%** | 2026-08-06T11:31:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 97% |
| Test Adequacy | 97% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — both warnings observed together in `summary`, with route's no-default phrasing correct and domain's unchanged.
**Follow-up REQs created:** None

**Minor:** the flag/`Original*` pairing in `collectSchemaFieldWarnings` is documented but not enforced — a future field wired half-way would warn on an empty value or never warn. A table-driven read site would make it structural; not worth the indirection for two fields.

**Restatement sweep:** ran. The diff renames a function with one call site (`buildBoard`) and adds two ticket fields. `grep -rn 'collectDomainWarnings'` → no remaining references anywhere, including tests, which assert through `board.Warnings` rather than the collector. No prose names the collector. `CLAUDE.md`'s parser enumeration already lists both `domain` and `route` (the merge resolution kept that), so the lock-step obligation is satisfied by this commit carrying both the code and — for route — the enumeration entry that made the obligation apply.

*Reviewed by review-work action*

---
*Source: `do-work/user-requests/UR-025/input.md` — Codex P2 finding on PR #133: "Warn when a route remains outside the enum"*
