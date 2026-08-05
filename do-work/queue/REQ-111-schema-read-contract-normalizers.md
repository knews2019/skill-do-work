---
id: REQ-111
title: Implement the seven missing Schema Read Contract field normalizers
status: pending
created_at: 2026-08-05T15:53:39Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
