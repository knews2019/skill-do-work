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
integration_at: 2026-09-04T20:00:14Z
review_at: 2026-09-04T20:08:35Z
kb_status: pending
kb_entry:
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
- [x] **[APPLY]:** Implemented canonical capture validation and minimal examples within the nine declared source files.
- [x] **[UNIFY]:** Reviewed every file listed in Implementation Summary for schema/read compatibility, safe publication refusals, copyable examples, and fixture intent. Diff check, Go vet, focused/race/module tests, contract regressions, and the direct builder gate passed; no debug artifacts remain.

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
- [x] `capture-reference.md` contains exactly one copyable minimal example for Simple REQ, Complex REQ, UR input, and Addendum REQ; the four fenced spans total at most 68 lines and contain no per-field explanatory comments.
- [ ] Every removed mechanical rule is enforced at the schema/publication boundary with public behavior tests, while legacy aliases remain accepted by ordinary readers.
- [ ] Capture judgment remains in prose, all existing publication safeguards retain their specific findings, and no replacement sentence predicate is introduced.

## Pre-Flight

**Git:** ✓ Working tree clean outside `do-work/`.
**Tests baseline:** ✓ `go test -count=1 ./internal/schemanormalization ./internal/requestmodel ./internal/publication` passed from the CLI module (17.7s on the recorded run).
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` passed directly; 375 board tests and 677 CLI tests passed, with the slowest files below 30 seconds.
**Dependencies:** ✓ Existing Go and shell toolchains satisfy the maintainer gate.

*Checked by work action*

## Implementation Summary

**What was done:** Replaced the four annotated capture examples with 67 fenced lines (previously 137), moved canonical authored-field validation into schema/publication code, and preserved tolerant legacy reads and capture judgment.

**Files changed:**
- `skills/do-work/actions/capture-reference.md` (modified) — minimal Simple, Complex extension, UR, and Addendum examples with schema pointers and separate optional proof/question shapes.
- `skills/do-work/actions/capture.md` (modified) — publication authority pointer and corrected assignment-example reference.
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go` (modified) — shared key aliases and canonical-authoring evidence.
- `skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go` (modified) — canonical/alias/default behavior and registry tests.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (modified) — schema-backed key aliases and field projections.
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modified) — legacy aliases and canonical-key precedence tests.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go` (modified) — canonical record validation before exposing record mutations.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go` (modified) — schema refusals, optional shapes, unordered membership, actual-example tests, and canonical fixtures.
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (modified) — canonical structured-override capture fixture without changing answer behavior.

**Predicate deletion:** No active predicate still pinned the removed comments; this leg is an evidence-backed no-op. No replacement sentence predicate was added.

## Decisions

- D-01 — DECIDE & STATE: Minimal Simple and Addendum examples use tdd:false and omit optional verdicts to avoid claiming nonexistent proof. Capture still judges applicability and preserves the separate proof shape.
- D-02 — DECIDE & STATE: UR membership is a duplicate-free set, independent of order. Existing missing-member findings stay specific; phantom/duplicate members receive schema findings.
- D-03 — DECIDE & STATE: Only non-default impact is mirrored in the title, matching the existing title convention. The machine validates an authored verdict, never chooses it.
- D-04 — DECIDE & STATE: No retired-comment predicate remains to delete; actual published examples are exercised through the public capture boundary instead.
- D-05 — DECIDE & STATE: The user narrowed this run to finishing this request and stopping. Run its selected heavy lanes immediately at this scoped queue boundary, retain its claim through review/finalization, and neither claim later work nor drain other held requests.

## Discovered Tasks

- [impact-negligible] The staged-skills heavy-only guard still describes user permission even though selected lanes run automatically at exhaustion; the guard itself works. → report only

## Qualification

**Result: Pass.** Canonical advance returned satisfied qualify and scope-drift records, with no findings, for range 366e1796bc2b0ca4f5b4a344e3c511a4c680dc8c..c00227166b288b97c60377cc06e7a5bfa736a0e8. All nine declared source files changed substantively. The orchestrator read the schema registry, request projections, publication validator, and example diff: public capture planning reaches validation before record mutations, and ordinary typed reads still normalize legacy values. The examples retain record shape and capture-time judgment; no live comment-pinning predicate remains.

**Requirement trace:** Minimal examples are covered by the actual-document public test and 67-line measurement; removed mechanical comments map to schema/publication behavior tests; legacy aliases and canonical-key precedence have explicit tests; capture judgments and optional proof/question shapes remain prose.

## Testing

**Fast tier: Pass.** The direct, unpiped maintainer gate passed on merged source at d3f26cfe84bf4a4f45e9e50caf183f2bd8558fc1: 375 board tests and 686 CLI tests, slowest file 19.09s, below the 30s file budget. Canonical advance recorded satisfied scope-drift, focused-test, and green-gate evidence. The uncached focused schema/requestmodel/publication packages passed in 0.134s, 0.261s, and 18.467s respectively; no baseline exclusions were needed.

**Builder verification:** Full CLI module tests and Go vet passed; race checks passed for schema/requestmodel/publication (1.617s, 1.901s, 26.541s). Contract regressions and diff checks passed. The standalone staged-skills script's heavy-only refusal was not counted as a pass; selected heavy lanes run through the typed runner below.

**Red-green validation:** The resumed builder directly observed the new reordered-membership test fail before changing positional comparison to duplicate-free set equality, then pass in focused/race/module/gate runs. The pre-interruption builder reported schema-refusal RED-before-code then GREEN; its original console transcript was unavailable after the model switch, so that evidence remains explicitly attributed. A separate reconstructed base-source overlay made all 23 schema-regression cases fail and current source passes them; reconstruction is regression evidence, not relabeled original TDD. The captured comment-predicate RED premise was stale: no live predicate still pinned those comments. The actual four examples now pass public capture planning and total 67 lines versus 137 before.

**Existing tests updated (cross-REQ impact):** REQ-413 (capture, answer, and release transactions; cf111a50 and remediation 1b8777d2) supplied the original raw-input and F6 structured-override fixtures; REQ-485 (canonical reservation marker filenames; 271aa8ae and 1425b982) supplied reservation/collision fixtures. Only payload metadata was canonicalized. Existing path, linkage, raw-byte containment, reservation, fold, and answer assertions remain and passed; these requests' behavior was not weakened.

**Heavy verification:** Typed runner completed at execution revision 0e5e7e37bdea6af82c38b248c54c017fbb9ac7d4 against the saved target c00227166b288b97c60377cc06e7a5bfa736a0e8. All four selected lanes ran, none skipped: do-work-cli-integrations exit 0 (46s; 700 CLI tests plus 110 script cases), staged-skills exit 1 (22s), updater exit 0 (51s), installer exit 0 (25s). The runner returned HEAVY-RUN-LANE-RED for staged-skills; command-level success is not claimed as lane success.

**Inherited heavy failure:** Staged-skills reports “core runtime must resolve actions/work.md through sibling do-work-board”. Its assertion at lines 960–971 and the work.md subject are byte-identical at the pre-merge base and current HEAD (work.md blob 85f98f893ab09d38be4629dedfbec0836f84b023; test blob 384def9f8f9695654e39bd9dca968eb3500a0b48). The predicate's required do-work-board/ text is absent at the base too; neither surface belongs to this request's changed files. This establishes an inherited failure, not a REQ-508 regression. No heavy-verified marker is asserted. Keep the failure visible and finish with issues rather than expanding the user's current-only scope to an unrelated repair.

## Review

**Overall: 82%** | 2026-09-04T20:06:32Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 88% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings:**
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go:193` and `:225` accept simultaneous list and nested-map evidence: public capture dry-runs accepted malformed mixed `depends_on` and UR `requests` fields. Reject mixed shape at the canonical-authoring boundary while retaining tolerant legacy reads. — impact-user-visible → report only
- `skills/do-work/actions/capture.md:7` still describes the removed alias table and forbids fields/enums not shown in now-minimal examples, contradicting the schema-authority pointer and potentially suppressing valid capture assessments. — impact-user-visible → report only

**Minor findings:** None.
**Acceptance:** Partial — canonical capture and scalar refusal work; mixed list/mapping shapes incorrectly pass. Four selected heavy lanes ran: three passed and staged-skills failed on an unchanged, independently attributed baseline assertion. No heavy-verified claim.
**Score calculation:** Mean 92%, minus 10 for Partial acceptance = 82%.
**Suggested testing:** 2 items.
**Follow-ups created:** None (2 findings report only).
**Self-validation:** Requirements, exact diff, final Testing, cross-REQ fixtures, restatement sweep, P-A-U, positive/negative public acceptance, and finding impact routing checked.

*Reviewed by independent review-work action. Public dry-run acceptance used an exclusive temporary fixture: canonical block lists succeeded, scalar lists refused without changes, mixed REQ/UR lists incorrectly succeeded; no queue was written.*

## Lessons Learned

**What worked:**
- Testing the actual published examples through the public capture boundary removed sentence-pinning without abandoning executable shape checks.
- Separate canonical authoring evidence preserved tolerant legacy reads while adding write-side refusals.

**What didn't:**
- Model-switch recovery retained code but lost the original console transcript; inherited testimony and reconstructed regression proof had to stay separately labeled.
- The first membership comparison wrongly depended on array order; a test-first reordered case exposed and corrected it.
- Review found that a tolerant parser can expose both list and nested-map evidence, and that capture's opening pointer still described exhaustive templates. These remain reported issues.

**Worth knowing:**
- Strict writers using tolerant parsers must judge exclusive shape evidence, not just the presence of a usable projection. Tests should read the shipped examples themselves as well as deliberately malformed payloads.

## Orientation

Capture now delegates mechanical authoring checks to the schema and publication layers. Its four examples show compact record shapes, while legacy readers retain aliases and the agent still judges applicability. Known limits are the mixed-list validation edge, stale opening instruction, and inherited staged-skills assertion recorded above. The touched action/shell prime pointers still resolve; no new subsystem or release owner was introduced.

## Finalization Judgment

Finish as completed-with-issues: core acceptance works, independent review is Partial at 82%, and one selected heavy lane remains red on an unchanged baseline assertion. Preserve all findings without starting another request, as the user explicitly instructed. No critical discovery or escalated builder decision requires a new queue item.

Release judgment: 0.275.3 → 0.275.4, a patch enforcing existing capture contracts and simplifying examples; no new user workflow or intended valid-input contract. VERSION and root CHANGELOG.md are project-owned; skills/do-work/VERSION and skills/do-work/CHANGELOG.md are required suite-maintainer mirrors. Exact release payload dry-run passed. KB handoff remains pending; no KB was initialized or written.

## Heavy Verification Plan

Base revision: 366e1796bc2b0ca4f5b4a344e3c511a4c680dc8c
Target revision: c00227166b288b97c60377cc06e7a5bfa736a0e8
Manifest: _dev/tests/heavy-lanes.json
Planner result: success; no uncovered paths, uncertainty, or forced-all selection.

### do-work-cli-integrations

Argv: ["bash","_dev/tests/maintainer-verify.sh","--heavy-lane","do-work-cli-integrations"]
Reasons:
- skills/do-work/tools/do-work-cli/internal/publication/answer_test.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/publication/capture_files.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go matched subtree skills/do-work/tools/do-work-cli

### staged-skills

Argv: ["bash","_dev/tests/maintainer-verify.sh","--heavy-lane","staged-skills"]
Reasons:
- skills/do-work/actions/capture-reference.md matched subtree skills
- skills/do-work/actions/capture.md matched subtree skills
- skills/do-work/tools/do-work-cli/internal/publication/answer_test.go matched subtree skills
- skills/do-work/tools/do-work-cli/internal/publication/capture_files.go matched subtree skills
- skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go matched subtree skills
- skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go matched subtree skills
- skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go matched subtree skills
- skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go matched subtree skills
- skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go matched subtree skills

### updater

Argv: ["bash","_dev/tests/maintainer-verify.sh","--heavy-lane","updater"]
Reasons:
- skills/do-work/tools/do-work-cli/internal/publication/answer_test.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/publication/capture_files.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go matched subtree skills/do-work/tools/do-work-cli

### installer

Argv: ["bash","_dev/tests/maintainer-verify.sh","--heavy-lane","installer"]
Reasons:
- skills/do-work/tools/do-work-cli/internal/publication/answer_test.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/publication/capture_files.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization.go matched subtree skills/do-work/tools/do-work-cli
- skills/do-work/tools/do-work-cli/internal/schemanormalization/schema_normalization_test.go matched subtree skills/do-work/tools/do-work-cli
