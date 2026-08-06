---
id: REQ-117
title: An unrecognized domain must leave a footprint on the board, not become general in silence
status: completed
created_at: 2026-08-06T10:53:03Z
claimed_at: 2026-08-06T11:05:05Z
completed_at: 2026-08-06T11:09:00Z
route: A
kb_status: pending
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
- [x] **[PLAN]:** Mirror `testing_status` exactly, since it is the same contract leg one field over: two ticket fields (`OriginalDomain`, `DomainUnrecognized`), a `collectDomainWarnings` twin of `collectTestingWarnings`, appended to `board.Warnings` at the same point in `buildBoard`. Reuse `schemaFieldWarningText` rather than writing a second copy of the contract's phrasing. Keep the `general` fallback and the absent-domain guard; keep the recognized flag instead of discarding it. Test at `buildBoard` level, copying `TestUnrecognizedTestingStatusFlagsAndWarns`'s shape.
- [x] **[APPLY]:** RED first — the new tests failed at build (`DomainUnrecognized undefined`), which is the honest RED for a field that does not exist. Then `model.go`: two struct fields, the read site keeping the flag, one collector, one append. No frontend change (see D-01).
- [x] **[UNIFY]:** `git diff --stat` → 2 implementation files (`model.go` +43/−9, `model_test.go` +99). `gofmt -l .` clean; `go vet ./...` clean; `go test ./...` passes (4.2s). No debug artifacts. Consumer sweep: `board.Warnings` is printed by `main.go`'s summary and rendered by `web/board.js`'s data-warnings banner (`boardData.warnings`) — so the warning reaches the UI with no JS change. Checked for a ratchet pinning the comment being replaced (`grep 'no warning channel'` across `_dev/`, `actions/`, `decisions/`, `kb/`): none. Inventoried every `domain:` value in the live tree — 110 REQs, all canonical (94 general, 9 backend, 6 testing, 1 frontend), so this adds zero warnings to this repo's own board.

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

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/model.go` (modified)
- `tools/queue-kanban/model_test.go` (modified)

**What was done:** `RequestTicket` gained `OriginalDomain` and `DomainUnrecognized`; the domain read site keeps `resolveSchemaField`'s recognized flag instead of discarding it; `collectDomainWarnings` mirrors `collectTestingWarnings` and is appended to `board.Warnings` in `buildBoard` right after it. The `general` fallback and the absent-domain guard are unchanged. The comment claiming the board has no warning channel for domain is replaced with what is actually true, including why the old reasoning was wrong.

## Decisions

- **D-01**: The unrecognized flag stays server-side — it raises a `board.Warnings` entry and is not projected into the card payload or given a card marker. **DECIDE & STATE** — the REQ left this open to the builder and the required half is the warning.

  Reasoning: two reasons, one of them discovered mid-build. First, YAGNI — the finding was "renders with no warning of any kind", and the warnings channel answers exactly that; a per-card visual marker is a new UI affordance nobody asked for. Second, and decisive: `web/board.js` already renders `boardData.warnings` in a data-warnings banner, so the warning reaches the board's UI *without* a frontend change. That makes the projection redundant rather than merely unrequested. Verified by generating the static payload from a synthetic `domain: quantum` fixture and finding the warning string in `board-data.js`.

## Testing

**Tests run:** `go test ./...` (in `tools/queue-kanban/`), plus `go test -run Domain .`
**Result:** ✓ All passing

**Red-green validation:**
- `TestUnrecognizedDomainFlagsAndWarns`: ✗ before (build failure — `ticket.DomainUnrecognized undefined`, `ticket.OriginalDomain undefined`) → ✓ after
- `TestRecognizedDomainRaisesNoWarning`: ✗ before (same build failure) → ✓ after. This is the guard against over-correcting: a documented alias and an absent field must both stay silent.
- End-to-end: built the binary, ran `summary` against a synthetic tree with `domain: quantum` → `warnings: 1`, `! REQ-997 ⚠ domain: 'quantum' not recognized — expected one of [frontend, backend, ui-design, general, security, testing]. Treating as 'general'.` Ran `generate` and confirmed the same string in `board-data.js`.

**New tests added:**
- `TestUnrecognizedDomainFlagsAndWarns`
- `TestRecognizedDomainRaisesNoWarning`

*Verified by work action*

## Lessons Learned

**What worked:** Reading the sibling field before designing anything. `testing_status` had already solved this exact contract leg — flag on the ticket, collector over all tickets, appended in `buildBoard` — so there was no design decision left to make, just a pattern to copy. The tests were a copy too, which is the strongest signal the shapes really match.

**What didn't:** The premise in the code comment was false, and it had been reviewed at 98% when it shipped. "The board has no warning channel for it" was three greps from being disproven: `board.Warnings` is declared with the words "surfaced, never silently dropped" in the same file, the sibling field feeds it, and the frontend renders it in a banner. A stated *reason* for a design choice is a factual claim and a review has to check it like any other — a plausible-sounding justification in a comment is exactly where a wrong premise survives longest.

**Worth knowing:** `board.Warnings` is a free UI channel. Anything appended to it prints in `summary` and renders in the board's data-warnings banner with no frontend work, because `web/board.js` reads `boardData.warnings` generically. Two consequences: a new warning class costs one `append` line, and noise is cheap to introduce — which is why the recognized-alias and absent-field cases have their own test. The contract's absent-field carve-out (`resolveSchemaField` returns recognized=true for an empty value) is what keeps a real queue from warning on nearly every REQ.

## Orientation

A typo'd `domain` no longer disappears into `general` silently — the card still reads `general`, and the board now says why in its data-warnings banner and in `do-work board summary`. Lives in the queue-kanban parser (`tools/queue-kanban/prime-do-kanban.md`), sharing `collectTestingWarnings`' pattern. Not `[MAP CHANGED]` — no new module or data flow; one existing display field gained the feedback leg its contract always specified. Prime staleness spot-check: `prime-do-kanban.md`'s referenced paths all still exist; nothing went stale.

## Review

**Overall: 96%** | 2026-08-06T11:08:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — the warning prints in `summary` and appears in the generated payload the frontend banner reads; a recognized alias and an absent field stay silent.
**Suggested testing:** 2 items
**Follow-up REQs created:** None

**Minor:** the warning is prefixed as `REQ-997 ⚠ domain: …` — REQ id, then the contract's own line — whereas `collectTestingWarnings` writes a single sentence in its own words. Both read fine in the banner; the composed form was chosen so the contract's phrasing stays in one place (`schemaFieldWarningText`), at the cost of two warning classes not being worded identically. Not worth unifying by hand-copying the phrasing back out.

**Suggested additional testing:** (1) open `do-work board serve` with a typo'd domain and confirm the banner renders the warning legibly — the payload proves the string, not the rendering; (2) a REQ carrying *both* a typo'd domain and a typo'd testing_status should produce two warnings, not one swallowing the other.

**Restatement sweep:** ran. The diff changes what a stored `domain` value means to the board (it now carries a recognized/unrecognized distinction). Consumers checked: `generate.go`'s `Domain` projection and `web/board.js`'s badge + filter — both read the resolved value, which is unchanged for every recognized input, so neither needed a change. `CLAUDE.md`'s display-parsed enumeration already names `domain` and its lock-step obligation is satisfied by this commit carrying both. No prose restates the discarded-flag reasoning — the false claim lived only in the code comment, which is fixed here.

*Reviewed by review-work action*

## Triage

**Route: A** - Simple

**Reasoning:** The pattern to mirror (`testing_status`) is in the same package, the fields and the channel are named in the REQ, and the contract fixes the fallback value. No exploration needed.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

---
*Source: `do-work/user-requests/UR-024/input.md` — Finding 2 of the 0.174.15-series feedback triage: "An unrecognized domain is silently remapped to general with no footprint — add a DomainUnrecognized flag and a board.Warnings entry mirroring testing_status, and fix the false 'no warning channel' comment at model.go:619"*
