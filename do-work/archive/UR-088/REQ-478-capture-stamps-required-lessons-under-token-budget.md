---
id: REQ-478
title: '[impact-rule-change] Capture stamps required lessons under a token budget'
status: completed
created_at: 2026-09-01T10:47:44Z
user_request: UR-088
domain: general
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: false
suggested_spec:
depends_on: [REQ-477]
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-477, REQ-479]
batch: lessons-transfer-routing
write_set: [skills/do-work/actions/capture.md, skills/do-work/actions/capture-reference.md, skills/do-work/actions/work-reference.md]
claimed_at: 2026-09-01T15:59:53Z
route: C
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-09-01T16:02:12Z
  basis:
    - Route C
    - 4-file write set
    - 2 subsystems involved
    - 7 acceptance criteria
    - dependency depth 1
    - cross-route regression gates
kb_status: promoted
kb_entry: REQ-478-capture-stamps-required-lessons-under-a-.md
completed_at: 2026-09-01T16:11:59Z
commit: 47ff6c85
---

# Capture Stamps Required Lessons Under a Token Budget

## What

`capture-request` reads the lessons index while authoring REQ payloads and stamps the relevant lessons files as mandatory reads in a new frontmatter field, keeping the stamped set's summed token estimates within one stated budget.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Defined one named budget contract, the ranking/narrowing algorithm, schema classification, and focused preservation proof before implementation.
- [x] **[APPLY]:** Added capture-time lesson routing, template/schema documentation, and the normalized writer round-trip test within the declared four-file scope.
- [x] **[UNIFY]:** Reviewed the isolated four-file delta, ran focused Go tests and contract searches, and passed the canonical maintainer gate on clean HEAD plus only this REQ's patch.

## Detailed Requirements

- While authoring REQ payloads (`skills/do-work/actions/capture.md` Step 5), capture reads the index (REQ-477) and decides which lessons files are relevant to the request being captured, stamping them in a new frontmatter field (suggested `required_lessons: [paths]`) on each minted REQ.
- **Token budget.** The stamped set's summed index estimates must stay within a stated budget (suggested ~2000 tokens per REQ; this REQ decides the number and where it is stated so it is one findable constant, not scattered).
- **Over budget:** capture ranks by relevance and stamps the best-fitting subset. Because lesson bullets carry family slugs (REQ-477), capture may stamp a targeted reference (path plus family slug) so the builder greps only the relevant bullets instead of reading the whole satellite — the cheapest way to stay in budget without dropping a match.
- **What was considered and dropped is noted in the REQ body, never silently.** Empty/absent when nothing matches — never invented.
- Add the field to the Simple REQ template (`skills/do-work/actions/capture-reference.md`) and the Request File Schema (`skills/do-work/actions/work-reference.md`).
- **Verify lossless preservation:** check whether `internal/requestmodel`/`internal/schemanormalization` must learn the field (unknown-field preservation may already cover it), and whether the board needs anything (display optional; the parser lock-step rule in `_dev/primes/prime-kanban-board.md` governs if it does).

## Constraints

- Plain files only; capture on the floor agent must be able to match index hooks with read/grep alone.
- The budget constant is stated once in one findable place, never scattered.
- Never invent a stamp: no index match means no field emitted (same never-invent posture as `assigned_to`).

## Dependencies

Depends on REQ-477 (index format and family slugs).

## Builder Guidance

Certainty level: Firm on the mechanism; latitude on the field name, budget value, and relevance-ranking judgment wording.

- [~] Budget value → builder decides; ~2000 tokens recommended.

## Red-Green Proof

**RED prompt/case:** Capture a request that touches rollback/deletion paths in do-work-cli internals today: the minted REQ carries no lessons pointer of any kind, and nothing in the capture flow reads the lesson satellites.
**Why RED now:** Capture has no lessons-routing step; builders reach lessons only through the touch-conditional rule at `work.md:404`, which the 2026-08-31 run showed does not transfer (REQ-415 repeated REQ-414's recorded family).
**GREEN when:** The same capture stamps `required_lessons` (or the chosen field) referencing the matching satellite within the stated budget, notes any dropped candidates in the REQ body, and the field is documented in both the capture template and the Request File Schema with lossless round-tripping verified.
**Validation:** User confirmed (approved plan, 2026-09-01 session).

## Full Context

See `do-work/user-requests/UR-088/input.md` for complete verbatim input.

---
*Source: UR-088 (Lessons routing with token-budgeted mandatory reads and a fold-gate fix)*

## Addendum (2026-09-01)

User added (v4 revision, validate-feedback Findings 2 and 5 — Accept):

> ```
> Two entry forms: a bare `path` (whole satellite) or `path#family-slug` (targeted; only the bullets carrying that slug). Targeted entries are permitted only for satellites whose index line says `slugged: full`; a `partial` satellite is stamped bare or not at all, so an un-slugged bullet can never be grepped past. [...] Cost rule, stated next to the constant: a bare path costs its index estimate; a targeted entry costs the size of its matching bullets (grep the slug, wc the lines, same formula as the index). Relevance ranking, in order: (a) the request text or its likely touched paths match a family named in the hook; (b) the satellite's owning prime governs a path the request names; (c) most recent same-family recurrence. Over budget, capture prefers narrowing a match to a targeted entry over dropping it; anything still dropped is listed in the REQ body under a fixed heading, never silently. [...] add a round-trip predicate to _dev/tests/contract-regressions.sh proving `required_lessons` survives internal/requestmodel/internal/schemanormalization unchanged (do not assume unknown-field preservation covers it — prove it)
> ```

- Entry forms are fixed: bare `path` (whole satellite) or `path#family-slug` (only the bullets carrying that slug).
- Targeted entries only when the satellite's index line says `slugged: full`; a `partial` satellite is stamped bare or not at all.
- Cost rule sits next to the budget constant: bare path costs its index estimate; a targeted entry costs the size of its matching bullets (grep the slug, `wc` the lines, same formula as the index). This supersedes the original "summed index estimates" wording, which left a targeted entry's cost undefined.
- Relevance ranking, in order: (a) request text or likely touched paths match a family named in the hook; (b) the satellite's owning prime governs a path the request names; (c) most recent same-family recurrence.
- Over budget, prefer narrowing a match to a targeted entry over dropping it; anything still dropped is listed in the REQ body under one fixed heading.
- The schema check upgrades from "verify whether" to a required round-trip predicate proving `required_lessons` survives the `internal/requestmodel`/`internal/schemanormalization` writers unchanged — do not assume unknown-field preservation covers it; prove it. Home is `_dev/tests/contract-regressions.sh` or a focused Go test, builder's call.
- Provenance: validate-feedback 2026-09-01, Findings 2 and 5. Surface-cost: N/A (specification of the accepted mechanism plus a proof test).

## Triage

**Route: C** — this changes the capture contract, shared REQ schema, and lossless request-model behavior.

## Plan

1. Define one canonical 2000-token budget, the two accepted entry forms, their mechanical costs, and the fixed dropped-candidates heading in the capture template reference.
2. Make capture rank index-backed candidates, prefer targeted narrowing, and omit the field when nothing matches.
3. Document `required_lessons` as a verbatim path-list field and prove an ordinary normalized request-model write preserves it byte-for-byte.

**Plan validation:** The three tasks cover routing, schema/template documentation, budget/drop behavior, and the required round-trip proof without changing the optional board display.

## Exploration

- Capture's only durable REQ authoring path is Step 5's payload construction before the `capture-files` transaction.
- The index already supplies path, hook family slugs, mechanical token estimate, coverage, and owning prime.
- `requestmodel.TypedRecord()` runs schema normalization, while `SetScalar` is the common lossless writer; one focused test can exercise both without modifying the paused REQ-420 contract script.
- The board has no scheduling or display requirement for this field, so leaving it unparsed avoids a new write/read surface.

## Scope

**Files I will touch:**
- `skills/do-work/actions/capture.md` (modify) — capture-time matching, ranking, narrowing, stamping, and omission rules
- `skills/do-work/actions/capture-reference.md` (modify) — single budget constant, entry/cost/drop contract, and Simple REQ field
- `skills/do-work/actions/work-reference.md` (modify) — frontmatter schema and verbatim-read classification
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modify) — normalized writer round-trip predicate

**Acceptance criteria:** Relevant captures stamp only index-backed entries within 2000 tokens; targeted entries require full slug coverage; drops use one fixed heading; no match emits nothing; schema writers preserve the field unchanged; the board remains intentionally unchanged.

## Decisions

- **D-01 — DECIDE & STATE:** Use 2000 tokens, the accepted recommendation, and state it once in `capture-reference.md`; all consumers cite the named contract.
- **D-02 — DECIDE & STATE:** Prove preservation in a focused Go test rather than `_dev/tests/contract-regressions.sh`, which is currently part of the paused REQ-420 write set.
- **D-03 — DECIDE & STATE:** Do not add board display. The field governs builder reads, not scheduling, and display was explicitly optional.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/capture.md` (modified) — index-backed relevance ranking, targeted narrowing, budget/drop recording, and no-match omission
- `skills/do-work/actions/capture-reference.md` (modified) — single 2000-token budget, exact entry/cost forms, fixed drop heading, and Simple REQ field
- `skills/do-work/actions/work-reference.md` (modified) — optional verbatim-read schema field and preservation/display boundary
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modified) — byte-for-byte preservation through normalization plus state mutation

**What was done:** Captured requests can now carry only the most relevant lesson reads that fit one reproducible budget. The contract narrows eligible fully-slugged satellites before dropping candidates, records every drop, and leaves unrelated requests unstamped.

## Qualification

Passed — the four declared files are substantive, every original/addendum criterion maps to the diff, and the optional board surface remains intentionally unchanged.

## Testing

- Focused `go test ./internal/requestmodel ./internal/schemanormalization`: PASS.
- Contract searches: PASS — exactly one `REQUIRED_LESSONS_TOKEN_BUDGET` definition; field/template/schema/drop-heading references present.
- `git diff --check` on the REQ scope: PASS.
- Clean isolated `bash _dev/tests/maintainer-verify.sh`: PASS (strict browser lane skipped because no browser was available, per the gate's own contract).

## Review

**Verdict: Approve.** The implementation states every accepted ranking, cost, coverage, omission, and dropped-candidate rule at an executable instruction site, and the focused test proves rather than assumes unknown-field preservation. No Important findings.

**Acceptance:** Pass

## Lessons Learned

**What worked:** Separating the one budget/cost contract from its consumers kept the number unique while letting capture and the next claim-time REQ share it; a focused request-model test proved the preservation seam directly.

**What didn't:** The shared tree's paused REQ-420 work overlapped two documentation files, so whole-file staging and dirty-tree qualification could not establish attribution.

**Worth knowing:** A budgeted context pointer needs both a single limit and a mechanical cost for every entry form; otherwise a “cheaper” targeted form makes compliance impossible to verify.

## Orientation

[MAP CHANGED] Capture now projects relevant family-keyed lessons into each minted REQ through a bounded `required_lessons` field. The action and do-work-cli schema contracts share one budget definition, while the board stays outside the data path.
