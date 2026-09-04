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
write_set: [skills/do-work/actions/capture-reference.md, skills/do-work/actions/capture.md, _dev/tests/contract-regressions.sh, skills/do-work/tools/do-work-cli/internal/schemanormalization/]
claimed_at: 2026-09-04T19:13:19Z
route: C
planning_at: 2026-09-04T19:21:18Z
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
