---
id: REQ-116
title: Normalize route at the board's read site and correct 0.174.15's board-wide claim
status: completed
created_at: 2026-08-06T10:53:03Z
claimed_at: 2026-08-06T10:56:55Z
completed_at: 2026-08-06T11:03:00Z
commit: 2a2cd59
route: A
kb_status: promoted
kb_entry: REQ-116-normalize-route-at-the-board-s-read-site.md
user_request: UR-024
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: []
write_set: [tools/queue-kanban/model.go, tools/queue-kanban/model_test.go, CLAUDE.md]
maintenance: false
related: [REQ-117, REQ-118]
batch: schema-contract-board-fixes
---

# Normalize `route` at the Board's Read Site and Correct 0.174.15's Board-Wide Claim

## What

The board parses `route` verbatim (`tools/queue-kanban/model.go` — `Route: coerceScalarToString(fields["route"])`), so a REQ written `route: a` reaches the card as lowercase `a`. Wire that read through the normalizer REQ-111 already added, and correct the 0.174.15 changelog entry's claim that the board honors the Schema Read Contract for all nine fields — it does not, and only `domain` was ever wired.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `tools/queue-kanban/prime-do-kanban.md`, `crew-members/general.md`, `coding-guardrails.md`, `testing.md`. Mirror the `domain` wiring one field over, with two deliberate differences: call `normalizeSchemaField` rather than `resolveSchemaField` (route's documented default is `""`, so resolving would blank an unrecognized letter), and keep the present-value-only guard so an absent route stays absent for `board.js`'s `if (request.route)` gate. Test first, at parse level — the normalizer's own table test already passed while the board was wrong.
- [x] **[APPLY]:** Four test cases written first and confirmed RED for the right reason (`Route = "a", want "A"`), then one read-site change in `model.go`. Two files, both in `write_set`.
- [x] **[UNIFY]:** `git diff --stat` → 2 implementation files (`model.go` +14/−1, `model_test.go` +105). `gofmt -l .` clean; `go vet ./...` clean; `go test ./...` passes (4.0s). No debug artifacts in the diff (no `fmt.Print`, `console.log`, `TODO`, `FIXME`). Consumer sweep per REQ-111's lesson — `grep '\.Route\b'` found one consumer (`generate.go` JSON projection) and `grep route web/board.js` found two (badge text, drawer row); none compares a route literal, none filters or buckets on it, so uppercasing cannot move a card. Built the binary and ran `summary` against the live tree: 118 tickets, 0 completion anomalies.

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

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/model.go` (modified)
- `tools/queue-kanban/model_test.go` (modified)
- `CLAUDE.md` (modified)

**What was done:** The board's `route` read in `parseRequestTicket` now goes through `normalizeSchemaField("route", …)` behind a present-value-only guard, replacing the bare `coerceScalarToString`. Four parse-level test cases pin the behaviour: lowercase uppercases, a canonical letter is unchanged, a padded lowercase letter uppercases, an absent route stays empty, and an unrecognized letter is reported case-folded rather than blanked. `CLAUDE.md`'s display-parsed field enumeration gained `route` (see D-01). The changelog correction for 0.174.15's board-wide claim is written in the Commit Phase entry, not as a rewrite of the shipped entry.

## Decisions

- **D-01**: Added `route` to `CLAUDE.md` → Shipped Tooling's display-parsed field enumeration (`depends_on`, `domain`, `write_set`, `assigned_to`, blocked fields), and extended this REQ's `write_set` to cover the file. **DECIDE & STATE** — reversible, one word, intent inferable from the REQ's subject.

  Reasoning: the restatement sweep asked who else states how the board handles `route`, and the answer was the lock-step sentence that names every field the board parses for display *and* obliges any contract change to be mirrored in `model.go`. `route` was missing from that list while being parsed, badged and drawer-rowed — so the obligation was never stated for the one field that then drifted from the contract. That is not a cosmetic omission; it is the mechanism whose absence produced this REQ. Fixing the enumeration is cheaper than re-deriving why route was exempt the next time someone reads the list. Out of scope and deliberately not touched: `batch` and `related` are also parsed and also absent from that list — the same shape, but neither is a Schema Read Contract field, so neither carries a normalize-and-warn obligation to mirror.

## Testing

**Tests run:** `go test ./...` (in `tools/queue-kanban/`), plus `go test -run 'TestParseRequestTicket.*Route' .`
**Result:** ✓ All passing

**Red-green validation:**
- `TestParseRequestTicketNormalizesRoute/lowercase_letter_uppercases`: ✗ `Route = "a", want "A"` before → ✓ after
- `TestParseRequestTicketNormalizesRoute/padded_lowercase_letter_uppercases`: ✗ `Route = "c", want "C"` before → ✓ after
- `TestParseRequestTicketReportsUnrecognizedRouteUnchanged`: ✗ `Route = "z", want "Z"` before → ✓ after
- `TestParseRequestTicketNormalizesRoute/canonical_letter_is_unchanged` and `TestParseRequestTicketPreservesAbsentRoute`: green before and after by design — they are the guards against over-correcting, not the RED pair.

**New tests added:**
- `TestParseRequestTicketNormalizesRoute` (3 sub-cases)
- `TestParseRequestTicketPreservesAbsentRoute`
- `TestParseRequestTicketReportsUnrecognizedRouteUnchanged`

*Verified by work action*

## Lessons Learned

**What worked:** Writing the test at parse level rather than at the normalizer. The normalizer's own table test (`{"route", "a", "A"}`) was green the entire time the board was wrong — the test existed, passed, and proved nothing about the field the user sees. Choosing the altitude of the assertion mattered more than the number of assertions.

**What didn't:** Nothing failed, but one instinct was wrong: `resolveSchemaField` is the natural sibling call and would have been a defect here. Its contract substitutes the field's documented default, and route's default is the empty string, so `route: z` would have arrived as absent — indistinguishable from a REQ with no route at all, in the exact field re-triage reads to find the problem. The two helpers are one word apart and differ in whether the caller may invent a value.

**Worth knowing:** The reason route drifted is recorded in the wrong place to prevent it. `CLAUDE.md`'s lock-step sentence names the fields the board parses for display and obliges any contract change to be mirrored in `model.go` — and `route` was never in that list, so a field that was parsed, badged and drawer-rowed carried no mirroring obligation. When a "keep these in sync" rule is expressed as a field enumeration, a field's absence from the list is silent permission to drift. Five of the contract's nine fields remain unread by the board on purpose; the guard against re-drift is the enumeration, not the code.

## Orientation

The board now spells route letters the way the schema does: a REQ written `route: a` shows an `A` badge instead of an `a`. Lives in the queue-kanban parser (`tools/queue-kanban/prime-do-kanban.md`), one field over from where REQ-111 wired `domain`. Not `[MAP CHANGED]` — no new module, no data-flow change, no renamed concept; the Schema Read Contract already said this is what route means, and this closes the last board read site that ignored it. Prime staleness spot-check: `prime-do-kanban.md`'s referenced paths (`main.go`, `allocate.go`, `release.go`, `verify.go`, `timestamp.go`, `walk.go`, `model.go`, `generate.go`) all still exist; nothing in it went stale.

## Review

**Overall: 97%** | 2026-08-06T11:02:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 98% |
| Scope | 95% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — built the binary and ran `generate` against a synthetic tree: `route: a` reaches the payload as `"route":"A"`, an absent route as `""`.
**Suggested testing:** 2 items
**Follow-ups created:** None

**Minor:** the same enumeration D-01 fixed also omits `batch` and `related`, which the board parses too. Left alone deliberately — neither is a Schema Read Contract field, so neither carries a normalize-and-warn obligation to mirror, and widening the list to every parsed field would change what the sentence is *for*.

**Suggested additional testing:** (1) open `do-work board serve` and confirm the route badge reads `A` on a card whose file says `a` — the payload proves the value, not the rendering; (2) the `--in-set` path is untouched here but shares `normalizeSchemaField`, so REQ-118's changes should re-run these route cases.

**Restatement sweep:** ran. `grep '\.Route\b'` → one Go consumer (`generate.go`'s JSON projection); `grep route web/board.js` → badge text and drawer row, neither comparing a literal; `_dev/tests/contract-regressions.sh` carries no route-parse assertion; `CLAUDE.md`'s display-parsed enumeration was stale and is fixed under D-01. The Schema Read Contract's `route` read-site column was deliberately left naming only work-pipeline steps — the contract's own rule is that the condition, not the caller list, is the trigger, and REQ-111 set the same precedent for `domain`.

*Reviewed by review-work action*

## Triage

**Route: A** - Simple

**Reasoning:** The read site, the function to call, the test file and the enum are all named, and the contract fixes every value decision. One call-site change plus one test.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

---
*Source: `do-work/user-requests/UR-024/input.md` — Finding 1 of the 0.174.15-series feedback triage: "The board reads route verbatim at model.go:655 — wire it through normalizeSchemaField (uppercase-only, no default) with a parse-level test, and correct 0.174.15's overstated board-wide claim in the next changelog entry"*
