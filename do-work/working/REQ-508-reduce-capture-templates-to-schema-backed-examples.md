---
id: REQ-508
title: '[impact-rule-change] Reduce capture templates to minimal examples backed by the schema layer'
status: claimed
priority: now
created_at: 2026-09-02T14:37:54Z
user_request: UR-098
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-507]
batch: orchestrator-simplification
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
write_set: [skills/do-work/actions/capture-reference.md, skills/do-work/actions/capture.md, skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go, skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go, skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go, skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go, skills/do-work/tools/do-work-cli/internal/publication/capture_files.go, skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go, skills/do-work/tools/do-work-cli/internal/publication/answer_test.go]
claimed_at: 2026-09-04T19:13:19Z
route: C
planning_at: 2026-09-04T19:21:18Z
exploration_at: 2026-09-04T19:28:16Z
dispatch_at: 2026-09-04T19:33:24Z
builder_handback_at: 2026-09-04T19:59:55Z
estimate:
  p50_active_minutes: 55
  confidence: low
  calculated_at: 2026-09-04T19:13:54Z
  basis:
    - Route C
    - 8-file write set
    - 3 subsystems involved
    - 3 acceptance criteria
    - dependency depth 4
    - cross-route regression gates
    - full-suite verification
---

# Reduce Capture Templates to Minimal Examples Backed by the Schema Layer

## What
`capture-reference.md` keeps one minimal example per record (Simple REQ, Complex REQ, UR input, Addendum REQ) and points at the schema normalizer for every field rule it enforces; per-field rule comments leave the templates.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Centralize canonical authoring and legacy read aliases in the schema layer, make capture-files validate newly authored UR/REQ shape before planning mutations, prove the missing rules through public Go tests, and reduce the four examples without moving capture judgment into code.
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
The Go schema layer already normalizes and validates frontmatter and `capture-files` refuses malformed records; the template comments restate those rules and drift from them.

## Context
Analysis: `ai-reports/2026-09-02_1651_orchestrator-simplification-analysis/index.html` (commit 1ddd7c70). Measured at 721c2fb4: `work.md` 850 lines and 20 steps; about 55% of step lines are mechanics; `_dev/tests/contract-regressions.sh` holds 220 references into the two work files and pins sentences with mutation-tested predicates, which is why earlier moves into Go left prose behind.

## Detailed Requirements
- Every rule deleted from a template comment is either already enforced by `schemanormalization` or `capture-files`, or gains a test there in the same commit.
- Templates stay a hard contract for shape; only restated rules leave.
- Delete the predicates that pinned the comments.

## Constraints
- One step per REQ, never a rewrite of `work.md`; the four-part write set (CLI command, deleted prose, deleted predicates, new behavior test) is complete or the review refuses the move.
- Judgment stays prose; `advance` emits typed findings, never paragraphs.
- The floor agent must still complete a run with `advance` output plus the remaining prose.
- Serial chain; run in one session.

## Dependencies
Depends on REQ-507 (shared writer on the contract suite).

## Builder Guidance
Firm on the boundary between mechanics and judgment as classified in the report's step table; dispute a row in the REQ before moving it. Latitude on prose wording. Read `_dev/primes/prime-action-files.md` before touching any action file.

## Red-Green Proof
**RED prompt/case:** Strip the per-field comments from the Simple REQ template and run the contract suite and `go test ./internal/schemanormalization ./internal/publication`.
**Why RED now:** Predicates that quote the comments fail; at least one rule exists only in a comment with no Go enforcement.
**GREEN when:** Suite passes without those lanes; every former comment rule has a Go test; templates are under half their current length.
**Validation:** User confirmed the direction ("more principles for the LLMs, not exact steps; the Go script does mechanics"); the per-REQ RED case is inferred during capture from the report.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3436 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-098/input.md` for complete verbatim input.

---
*Source: capture of the orchestrator simplification request (UR-098).*

## Triage

**Route: C — Complex.** The change removes instructional contracts from four capture templates only after tracing each rule to schema or capture-file enforcement, adds missing behavior coverage, and updates the active predicate owner without weakening record shape.

**Planning:** Required.

## Plan

1. Add public RED cases for alias keys/values, malformed record shape, timestamps/title scalars, blocked metadata, and phantom UR membership, plus positive fixtures for all four examples.
2. Centralize the five read-only key-alias families and canonical-write evidence in schema normalization, then make requestmodel consume that registry without changing ordinary reads.
3. Validate canonical new-record identity, linkage, field shapes, timestamps, user-text encoding, blocked-field pairing, and exact membership at the capture-files boundary before any mutation is planned.
4. Replace the Simple, Complex, UR, and Addendum fences with copyable minimal examples while retaining the named judgment contracts and updating capture's authority pointer.
5. Verify schema/read compatibility, publication behavior, template size/residue, race/module/contracts, and the direct repository gate.

## Exploration

The lossless request parser already exposes source keys, raw and decoded scalar forms, list shape, duplicate counts, and source lines. That is sufficient for `capture-files` to enforce canonical authored keys and values, safe scalar spelling, whole-second UTC timestamps, list shapes, blocked-field pairing, and exact UR membership while ordinary queue/archive reads retain legacy alias compatibility.

There is no live shell predicate pinning the four annotated templates: `_dev/tests/contract-regressions.sh` is now a dispatcher, and the only remaining `capture-reference.md` mention in the contract checks is unrelated fixture data. The predicate-deletion leg is therefore an evidence-backed no-op. Strict publication validation will require the existing answer override fixture to use canonical UR/REQ records, so that test is included in the declared boundary.

The machine must not infer Simple versus Complex, field applicability, dependency or lesson relevance, optional body sections, or whether an authored timestamp is “current”; those remain capture judgments in prose.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/actions/capture-reference.md` (modify) — replace the four annotated fences with copyable minimal examples and schema/publication pointers.
- `skills/do-work/actions/capture.md` (modify) — identify the examples as shape contracts and retain capture-time judgment ownership.
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go` (modify) — centralize key aliases and expose exact canonical-authoring evidence.
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go` (modify) — prove the key registry and canonical-versus-legacy value behavior.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (modify) — consume schema-owned aliases and project the fields publication validates.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modify) — prove alias projection and canonical-key precedence.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go` (modify) — validate canonical new-record shape before publication mutations can run.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go` (modify) — add RED/GREEN authoring-contract cases and canonicalize existing capture fixtures.
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (modify) — canonicalize the structured override-capture fixture for the shared publication validator.

**Files I will NOT touch:** `_dev/tests/contract-regressions.sh`, `_dev/tests/contracts/core-checks.sh`, publication manifest/type production files, `skills/do-work/actions/work-reference.md`, `skills/do-work/docs/capture-guide.md`, and board code; exploration found no required predicate deletion or contract wiring there.

**Acceptance criteria (restated from REQ):**
- [ ] `capture-reference.md` contains exactly one copyable minimal example for Simple REQ, Complex REQ, UR input, and Addendum REQ; the four fenced spans total at most 68 lines and contain no per-field explanatory comments.
- [ ] Every removed mechanical rule is enforced at the schema/publication boundary with public behavior tests, while legacy aliases remain accepted by ordinary readers.
- [ ] Capture judgment remains in prose, all existing publication safeguards retain their specific findings, and no replacement sentence predicate is introduced.

## Pre-Flight

**Git:** ✓ Working tree clean outside `do-work/`.
**Tests baseline:** ✓ `go test -count=1 ./internal/schemanormalization ./internal/requestmodel ./internal/publication` passed from the CLI module (17.7s on the recorded run).
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` passed directly; 375 board tests and 677 CLI tests passed, with the slowest files below 30 seconds.
**Dependencies:** ✓ Existing Go and shell toolchains satisfy the maintainer gate.

*Checked by work action*
