---
id: REQ-111
title: Implement the seven missing Schema Read Contract field normalizers
status: completed
created_at: 2026-08-05T15:53:39Z
claimed_at: 2026-08-05T19:13:04Z
completed_at: 2026-08-05T19:19:07Z
commit: e77383a
route: A
user_request: UR-021
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: []
write_set: [tools/queue-kanban/model.go, tools/queue-kanban/model_test.go]
maintenance: false
related: [REQ-112]
batch: census-durable-findings
---

# Implement the Seven Missing Schema Read Contract Field Normalizers

## What

The Schema Read Contract (`actions/work-reference.md` L200–210) defines nine enum-or-boolean fields, each with an alias map, a canonical enum, and a documented default-plus-warning on an unrecognized value. Only **two** have a mechanical implementation anywhere in the repo: `normalizeStatus` (`tools/queue-kanban/model.go` L718) and `normalizeTestingStatus` (`tools/queue-kanban/testing.go` L59). Add normalizers for the other seven, following `normalizeStatus`'s existing shape.

The seven, with their contract rows:

| Field | Canonical enum | Normalization the contract requires | Default on unknown |
|---|---|---|---|
| `domain` | `frontend`, `backend`, `ui-design`, `general`, `security`, `testing` | `back-end`/`back_end` → `backend`; `front-end`/`front_end` → `frontend`; `ui_design` → `ui-design`; `sec` → `security`; `test` → `testing` | `general` |
| `route` | `A`, `B`, `C` | lowercase `a`/`b`/`c` → uppercase | needs re-triage |
| `caveman` | `false`, `true`, `lite`, `full`, `ultra` | `yes`/`on` → `true`; `light` → `lite` | `false` |
| `maintenance` | `true`, `false` | `yes`/`on`/`t` → `true`; `no`/`off`/`f` → `false` | `false` |
| `tdd` | `true`, `false` | `test_first`/`yes`/`on`/`t` → `true`; `no`/`off`/`f` → `false` | `false` |
| `error_type` | `intent`, `spec`, `code`, `environment` | (no aliases identified) | `code` |
| `kb_status` | `promoted`, `pending`, `declined`, `skipped` | `skip` → `skipped`; `rejected` → `declined` | `pending` |

`actions/work-reference.md` L200–210 is the source of truth for every row above — read it rather than trusting this table, and if the two disagree, the contract wins and this REQ's table is the stale copy.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `normalizeStatus` (`model.go`) and `normalizeTestingStatus` (`testing.go`) to find the house convention: a `normalizeX(raw) string` that never decides, paired with an `isKnownX(normalized) bool` the caller uses to warn. Mirror that, but table-driven — the seven rows are structurally identical, and REQ-112 receives the field name as runtime data from a CLI flag, so a lookup table is required rather than speculative. Verify: `go test ./...` green, `gofmt -l` and `go vet` clean.
- [x] **[APPLY]:** Added the contract table plus three functions to `model.go` and wired `domain` at its read site. Nothing else touched — `status` and `testing_status` keep their existing normalizers and are dispatched to rather than forked.
- [x] **[UNIFY]:** `git diff --stat` → 2 files (`model.go` +150/−1, `model_test.go` +166). `gofmt -l .` clean; `go vet ./...` clean; `go test ./...` passes. Grepped added lines for `console.log`/`debugger`/`TODO`/`FIXME`/`XXX`/`fmt.Print` debug artifacts: none. Built the binary and ran `queue-kanban summary` against the live tree — 112 REQs, 0 completion anomalies, REQ-112 correctly bucketed as waiting-on-deps. Ran `_dev/tests/contract-regressions.sh`: 7 failures, byte-identical to the pre-existing `main` baseline.

## Why (if provided)

`domain` is currently read verbatim via `coerceScalarToString` (`model.go` L626), so `domain: back-end` silently mis-selects the crew file the REQ meant to load, with nothing anywhere to catch it. The contract has been correct and centralized since 0.76.2; only its enforcement is missing. This is the finding the census rated most durable, because it is about the absence of a mechanism rather than about any line number.

## Detailed Requirements

- One normalizer per field, mirroring `normalizeStatus`'s signature and placement in `model.go`.
- Each returns the canonical value on a recognized input or alias, and signals the unrecognized case so the caller can emit the contract's warning — `⚠ {field}: '{value}' not recognized — expected one of [{enum}]. Treating as '{default}'.` Do **not** silently remap an unknown value; the warning is the whole point of normalize-and-warn (the contract's item 3 records why: silence is what let `dependencies:` go unnoticed pre-0.76.2).
- Wire `domain` through its new normalizer at the existing read site (`model.go` L626) so the board stops reading it verbatim.
- **Do not change any write path.** The contract states write paths are unaffected (`work-reference.md` L212) — capture and the work pipeline always emit canonical values. This REQ is read-side only.

## Constraints

- Pure Go inside `tools/queue-kanban/`. No action prose changes, no new subcommand — REQ-112 owns exposure.
- Keep the parser in lock-step with the schema: `CLAUDE.md` → Shipped Tooling requires a contract change and its `model.go` update in the same commit. Here the contract is unchanged and only the implementation lands, so no `actions/work-reference.md` edit is expected.

## Dependencies

None. REQ-112 depends on this one.

## Builder Guidance

**Firm.** The seven rows are specified by an existing contract, not open to design. The only judgment is code placement and how the unrecognized-value signal is returned (error, sentinel, or second return value) — match whatever `normalizeStatus` and `normalizeTestingStatus` already do rather than introducing a third convention.

## Red-Green Proof

**RED prompt/case:** In `tools/queue-kanban/model_test.go`, assert that a REQ whose frontmatter says `domain: back-end` parses to `Domain == "backend"`. It fails today, returning the literal `back-end`, because `model.go` L626 reads the field with `coerceScalarToString` and no normalizer exists.
**Why RED now:** Only `status` and `testing_status` have normalizers. The other seven fields are read verbatim or not read at all, so every alias the contract promises to accept is silently passed through unchanged.
**GREEN when:** `domain: back-end` parses to `backend`, each of the other six fields resolves its documented aliases, and an unrecognized value for any of the seven produces the contract's warning text plus the documented default instead of the raw input.
**Validation:** Inferred during capture — derived from the contract table at `actions/work-reference.md` L200–210, not confirmed field-by-field with the user.

## Assets

None.

---
*Source: census finding — 7 of 9 Schema Read Contract enum fields have no normalizer anywhere in the repo (`decisions/audits/2026-08-05-shell-logic-in-prose-census.md` §3 closing paragraph)*

---

## Triage

**Route: A** - Simple

**Reasoning:** The seven rows are fully specified by an existing contract table, the target file is named, and the house convention to mirror already exists in the same file. No exploration or planning needed.

## Implementation Summary

**What was done:** Added the Schema Read Contract's seven missing field normalizers to the board's parser as a lookup table plus three functions, and wired `domain` through it at its read site. `status` and `testing_status` keep their existing normalizers and are dispatched to.

**Files changed:**
- `tools/queue-kanban/model.go` (modified) — added `schemaFieldContract` type, the `schemaReadContractFields` table (7 rows), `normalizeSchemaField`, `isKnownSchemaFieldValue`, `resolveSchemaField`; wired `domain` at the `parseRequestTicket` read site
- `tools/queue-kanban/model_test.go` (modified) — 5 new tests: contract aliases across all 7 fields, default-fallback-without-silent-remap, dispatch of `status`/`testing_status` to their own normalizers, the `domain: back-end` RED case, and an absent-domain regression guard

**Not changed, deliberately:** no write path (the contract states write paths are unaffected), no action prose, no new subcommand (REQ-112 owns exposure), and `normalizeStatus`/`normalizeTestingStatus` are untouched.

## Decisions

- **D-01: A lookup table rather than seven `normalizeX` functions.** DECIDE & STATE. The REQ said "mirroring `normalizeStatus`'s signature," which reads as seven functions, but the rows are structurally identical (alias map, enum, default) so seven copies would be seven places for the contract to drift from its home — the exact failure the census was written about. Decisive factor: REQ-112 receives the field name as **runtime data** from a CLI flag, so it needs a lookup regardless; seven hardcoded functions would require a switch that re-derives this table. The table is foundation for an already-captured sibling, not speculation.
- **D-02: `status` and `testing_status` dispatched, not folded into the table.** DECIDE & STATE. Copying their alias maps into the table would create the second definition the table exists to prevent, and `coding-guardrails.md` § Surgical Changes says don't refactor what isn't broken. `normalizeSchemaField` routes those two field names to the existing functions, so callers get a complete 9-field surface without either function changing.
- **D-03: an absent field resolves to the default with `recognized=true`.** DECIDE & STATE. All seven fields are optional, so treating absence as a contract violation would emit a warning for nearly every REQ in a real queue and train readers to ignore the channel the contract added the warning for. Absence is not an unrecognized value.
- **D-04: the `domain` call site normalizes only a *present* value.** DECIDE & STATE, and it caught a regression — see Testing. `resolveSchemaField`'s absent→default behaviour is right for a reader that must pick a crew file and wrong for the renderer, so the call site is explicit rather than the function being changed.

## Qualification

Passed — 2 files verified present in the diff, both mechanical checks and judgment checks clean. Every requirement traced: 7 normalizers present and covered by tests, `domain` wired at `model.go`'s read site, no write path touched, no prose changed.

## Testing

**Tests run:** `go test ./...` (the tool's own suite), `gofmt -l .`, `go vet ./...`, `_dev/tests/contract-regressions.sh`, plus a `queue-kanban summary` smoke run against the live tree.

**Result:** Go suite passes; `gofmt` and `vet` clean. Contract suite shows its 7 pre-existing update-script failures, byte-identical to the `main` baseline verified earlier this session — no new regressions. Smoke run parsed 112 REQs with 0 completion anomalies and bucketed REQ-112 as waiting-on-deps, confirming the parser change didn't disturb bucketing.

**Red-green validation:**
- RED: `TestParseRequestTicketNormalizesDomain` and the three normalizer tests failed to compile — `undefined: normalizeSchemaField`, `undefined: resolveSchemaField`. Confirmed before writing any implementation.
- GREEN: all four pass. `domain: back-end` now parses to `backend`.

**A regression the RED/GREEN pair did not cover, caught during UNIFY:** wiring `domain` through `resolveSchemaField` directly made an *absent* domain resolve to `"general"`. `web/board.js` gates both the domain badge (L481) and the filter dropdown (L379) on `if (request.domain)`, so every domain-less card would have gained a badge and a filter entry it never had — a visible UI change nobody asked for, and one the existing suite did not catch. Wrote `TestParseRequestTicketPreservesAbsentDomain`, confirmed it failed with `Domain = "general"`, then fixed the call site to normalize only a present value. This is why the tests are 5 and not 4.

## Review

**Approve** — implements all seven contract rows against the existing convention, and the one regression it introduced was caught, guarded by a test, and fixed before commit.

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Minor:** the `schemaReadContractFields` table restates a contract whose home is `actions/work-reference.md`, so the two can now drift. Mitigated by a comment naming the contract as source of truth and requiring same-commit changes — the same lock-step rule `CLAUDE.md` already imposes on this parser — but it is a second copy, and the census's own central finding is that second copies drift.
**Acceptance:** Pass — `domain: back-end` → `backend`; unrecognized values fall back with a recognized=false signal; absent fields stay absent for the renderer.
**Follow-ups created:** None.

## Lessons Learned

**What worked:** Reading both existing normalizers before designing anything. The house convention (normalizer never decides; caller warns) was not obvious from the contract text alone, and matching it meant the new code needed no explanation for reviewers who know `normalizeStatus`.

**What didn't:** The first wiring of `domain` was wrong in a way the REQ's own RED/GREEN proof could not detect. The proof asserted that a *present* alias normalizes; it said nothing about absence, so the suite went green over a UI regression. The catch came from UNIFY's habit of asking who consumes the changed field — `grep '\.Domain'` in Go, then `grep domain web/*.js` — not from the test plan. A RED/GREEN pair proves the stated behaviour; it does not bound the change.

**Worth knowing:** `resolveSchemaField` deliberately maps an absent field to the contract's default, which is correct for a *decision-making* reader (work.md Step 6 must pick some crew file) and wrong for a *rendering* reader (the board must show nothing). Any future caller has to choose consciously. The distinction is documented at both the function and the call site, because getting it wrong is invisible in Go and visible only in the browser.

## Orientation

The board's parser now enforces the Schema Read Contract for all nine of its enum fields instead of two: `domain: back-end` resolves to `backend`, and the other six aliased fields resolve through one table. Lives in `tools/queue-kanban/model.go`, the parser the Schema Read Contract is kept in lock-step with. No user-facing capability is added — REQ-112 exposes this to prose. `[MAP CHANGED]` for the parser only: a new `schemaReadContractFields` table is now the mechanical home of a contract that previously existed only as prose plus two hand-written normalizers.
