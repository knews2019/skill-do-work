---
id: REQ-759
title: 'Drop routine capture verification and self-scoring'
status: completed
dispatch_at: 2026-11-25T23:12:41Z
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-11-25T23:12:22Z
  basis:
  - Route A
  - 4-file write set
  - 2 subsystems involved
  - 5 acceptance criteria
route: A
created_at: 2026-11-25T22:57:02Z
user_request: UR-190
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
maintenance: true
impact: impact-user-visible
effort_estimate: effort-substantive
batch: skill-instruction-maintenance
related: [REQ-758, REQ-760]
required_lessons: ['_dev/primes/lessons-action-files.md#cross-action-exception-closure']
write_set: [skills/do-work/next-steps.md, skills/do-work/actions/capture.md, skills/do-work/actions/verify-requests.md, skills/do-work/docs/verify-requests-guide.md]
claimed_at: 2026-11-25T23:11:32Z
builder_handback_at: 2026-11-25T23:36:20Z
review_at: 2026-11-25T23:43:04Z
commit: e7c00807cb36392703381ebf408e0aa9d04fadf3
heavy_verified_at: 2026-11-25T23:44:16Z
heavy_verified_revision: e7c00807cb36392703381ebf408e0aa9d04fadf3
completed_at: 2026-11-25T23:44:39Z
release_at: 2026-11-25T23:44:39Z
---
# Drop routine capture verification and self-scoring
## What
Remove the default verification suggestion for a new capture and the extra nudge for long or complex captures. An explicit verification reports the original clause, which REQ changed or missed it, and the correction or remaining ambiguity. Delete the confidence percentages, score bands, averaging, and post-repair re-score. Align the user guide with that output.

## Why
What matters to the maintainer is how often they have to restore an already-stated requirement. A discrepancy the verifier fixes itself, a user correction, a question already answered in the request, and a later change of preference have different meanings. This change concerns verification and its report. Building a metrics system is outside its scope.

## Constraints
Keep explicit verification and repair, comparison with folded destinations, and the gap severities. Preserve the verification that the user opts into through `validate-feedback --capture` and its gap-based run decision at `skills/do-work-toolbox/actions/validate-feedback.md:121`. Report corrections as applied only when a repair actually occurred.

## Red-Green Proof
**RED prompt/case:** Follow a new capture through `skills/do-work/next-steps.md:27` and `skills/do-work/actions/capture.md:223`. Then replay `do-work/archive/UR-011/input.md:75` against the verification report at `skills/do-work/actions/verify-requests.md:118` and the re-score at line 171.
**GREEN when:** Ordinary capture can go straight to build on the user's instruction. Explicit verification identifies that the capture described a hard stop while the original said the agent warned and recovered. The report shows the clause and disposition without a numerical score. UR-065, by contrast, records a confirmed ambiguity about blocked work in the Timeline forecast (`do-work/archive/UR-065/input.md:108`).
**Validation:** Inferred from the accepted finding and its archived cases.

## Full Context
See `do-work/user-requests/UR-190/input.md` for the verbatim capture instruction.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Remove routine verification advice and make explicit verification a clause-level report with repair dispositions. Keep the feedback action's opt-in gap gate and REQ-758's capture checks.
- [x] **[APPLY]:** Merged the four builder paths with the tutorial and feedback wording seams. Both added paths are callers of the changed capture guidance or verifier output.
- [x] **[UNIFY]:** Reviewed the six paths below against the first-parent diff and checked the range with `git diff --check`. The builder's fast contract tests and defensive audit passed. The integration gates are recorded below.
*Source: Finding 2, pasted audit, rank 2, no severity supplied. Accepted at efcd3618. Surface-cost: N/A. Remove the suggestion to run `verify-requests` after every new capture. When verification is requested, report the original clause, what the capture changed or missed, and the correction. Delete the percentage table and re-scoring step.*

## Triage

Route A. The request names the capture, verification and handoff entry points, the report change and its opt-in exception.

## Plan

**Planning not required** — Route A, direct implementation.

## Implementation Summary

- `skills/do-work/next-steps.md` (modified) — removes routine capture verification advice and keeps the mechanical-batch route.
- `skills/do-work/actions/capture.md` (modified) — ends ordinary and complex captures with the same brief summary and keeps the separate publication and linkage checks.
- `skills/do-work/actions/verify-requests.md` (modified) — reports severity, source clause, destination, discrepancy, and proposed correction or ambiguity. Follow-up records the actual disposition. Preserves fold comparison and decision revalidation.
- `skills/do-work/docs/verify-requests-guide.md` (modified) — documents the optional check, findings, and repair outcomes.
- `skills/do-work-toolbox/actions/tutorial.md` (modified) — teaches capture followed by an explicit run for feature, bug, and complex inputs. Retains the recipe for an explicit verification request.
- `skills/do-work-toolbox/actions/validate-feedback.md` (modified) — removes the stale score reference from the existing opt-in gap gate.

Implementation range `ff00a6dc7356555f34b85dcc8ea45c58e6d3d4ad..e7c00807cb36392703381ebf408e0aa9d04fadf3`. The builder SHA `9a3eda8a718a1ae263aa1ac6a971774ae3b379e7` is the merge's second parent.

## Decisions

- **D-01 (DECIDE & STATE).** Keep Step 6 as the report boundary. `validate-feedback --capture` consumes the gap list without entering Offer Fixes. Its explicit `--run` gate uses severity and owned REQs.
- **D-02 (DECIDE & STATE).** Refresh estimates after a material repair. An estimate tracks changed scope independently of capture verification.
- **D-03 (DECIDE & STATE).** Define Nit by its cosmetic effect. The placement or frequency of a source clause does not determine whether it is a requirement.

## Qualification

Passed. Request-bound `advance` accepted the merged range with no findings. All six paths serve capture, explicit verification, or a caller of those instructions.

## Testing

- The read-only replay found the routine nudge, tutorial sequence, confidence table, and re-score on the baseline. On the merge, simple and complex captures end with a summary, a user-directed `do-work run` can follow directly, and the tutorial retains an explicit verification recipe. The mechanical-batch suggestion still routes to `run-simple-reqs`.
- UR-011's “I resolved it rather than erroring” conflicts with the original hard-stop RED in REQ-067. The Important finding proposes the warning and reader-dependent recovery that the archived correction records as applied. UR-065's “Reqs that are not in the period selected should not be listed” was Ambiguous for REQ-319's separate forecast exclusions. The user confirmed that list stays whole-queue and the archived REQ records the answer.
- The public `answer` command in `verify-repair` mode inserted and resolved a question in a clean pending fixture. It returned `ANSWER-ARCHIVED-READ-ONLY` for the archived variant and preserved its bytes. The verifier marks applied corrections after success and reports recorded answers, deferrals, and open questions.
- `validate-feedback --capture` takes the Step 6 gap report and stops before Offer Fixes. Without `--run` it stops after reporting. With `--run`, it skips when all accepted findings folded, runs an owned batch with a clean or Nit-only report, and stops for an Important, Minor, or Ambiguous gap. Fold replays covered named REQ clauses, actionable sweep instances, open and resolved prose items, and a missing fold classified as Important.
- The merged UR-069 walkthrough retains a mandatory overflow count, a discretionary row count, and independent verbatim-input, publication, REQ and fold linkage checks. The `--against` workflow matches its baseline bytes.
- `bash _dev/tests/shipped-package-reference-contract.sh` passed through the request-bound focused gate. `bash _dev/tests/maintainer-verify.sh` passed directly on the first run, exit 0 in 184 seconds. The request-bound green record names `e7c00807cb36392703381ebf408e0aa9d04fadf3`.
- The typed heavy planner selected `staged-skills` for the six shipped paths. The plan and execution record follow below.

The natural-language actions were walked through with their ordinary entry points and archived examples. The answer transaction ran on disposable Git fixtures.

## Review

**Overall: 97%** | 2026-11-25T23:43:04Z

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Acceptance:** Pass. Reviewed the first-parent diff and UR-190 against the capture, verify, and feedback entry points. The independent tutorial finding and the Nit wording carried from REQ-758 are resolved in the merge. Read-only replays and direct answer fixture results are recorded under Testing.

**Restatement sweep:** Checked the default next step, optional verification, gap report, repair disposition, and opt-in feedback gate across the routers, tutorial, guides, copy prompts and package tests. REQ-758 was archived and REQ-757 was still in `working/`, so this was the current change's sweep.

**Minor findings:**

- `skills/do-work/actions/capture-reference.md` still uses “grades” for checking a fold. Its operative instructions continue to check presence and fidelity — impact-negligible → report only
- The README and help full-cycle copy prompts explicitly request verification and leave “verification fails” to the prompt author to interpret — impact-rule-change → report only

**Follow-ups created:** None.

## Orientation

Capture can lead directly to a user-directed run. An explicit verification points to the original clause and records what was repaired or still needs a decision.

## Heavy Verification Plan

- Base revision `ff00a6dc7356555f34b85dcc8ea45c58e6d3d4ad`
- Target revision `e7c00807cb36392703381ebf408e0aa9d04fadf3`
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`. All six changed paths matched the `skills` subtree.

## Heavy Verification Result

- Target revision `e7c00807cb36392703381ebf408e0aa9d04fadf3`
- Execution revision `e7c00807cb36392703381ebf408e0aa9d04fadf3`
- `staged-skills` — executed in 37 seconds, exit 0. Command `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`.
