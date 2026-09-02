# Review: REQ-489

**Approve with follow-up** — the canonical request-state departure path now removes complete checkpoint entries and uses exact heading lines, but the repository’s cleanup departure path still implements the old header-only removal contract and can leave the same orphaned continuation block.

Route A | implementation range `1832538d..6e92e536`

## What’s built

- `checkpointWithoutClaim` now limits removal to the exact `## In Progress (interrupted)` section and removes a matching own-label header with its adjacent nonblank indented continuation lines.
- `appendSectionEntry` and removal share exact whole-line section discovery, so inline and backticked heading mentions no longer attract insertion or removal.
- Two focused regressions cover enriched-entry removal, foreign-entry preservation, and inline-heading preservation.

## Decision / risk

The merged request-state change is clean and should remain. Before calling the repository-wide “whenever a REQ leaves working” contract closed, fold the cleanup mover into the same root-cause remediation: `internal/cleanup.ownedCheckpointRemoval` still filters matching header lines globally and preserves their indented details. This is a low-risk residual correctness defect rather than a reason to revert the merged fix.

## Findings

### Important

1. **Cleanup still removes only the claim header when it archives a terminal REQ from `working/`.** `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go:218` implements `ownedCheckpointRemoval` as a global line filter: it skips the matching `- REQ-NNN:` line and appends every following line, including Step 10’s indented `Last known state`, `Key files being modified`, and `Known issues` details. This is the same stale checkpoint-entry semantic the diff corrects in `requeststate`, and `actions/work-reference.md` explicitly makes cleanup Pass 0 one of the movers that inherits whole-entry removal. A cleanup run can therefore archive the REQ while leaving an unattributed orphan block in `CHECKPOINT.md`. — `impact-user-visible` → **suggested disposition:** fold this instance into the existing `checkpoint-section-blind-line-editing` sweep/root cause (or author one addendum follow-up if the active REQ can no longer be amended), with a cleanup-package regression that fails on the enriched entry before the fix and preserves the foreign entry byte-for-byte.

### Minor

None.

## Requirements checklist

- [~] Remove the header and all indented continuation lines when a REQ leaves `working/` — delivered for canonical `requeststate` complete/fail/cancel transitions; the cleanup recovery mover remains header-only.
- [x] Preserve behavior for a bare one-line entry — the existing lifecycle test passes, and the new loop skips no unrelated line when no continuation exists.
- [x] Preserve foreign-label entries and their continuation lines — directly asserted by `TestCheckpointClaimRemovalIncludesIndentedContinuationLines` and preserved by the section-bounded implementation.
- [x] Locate `## In Progress (interrupted)` by a whole heading line for insertion and removal — delivered through `sectionLineBounds`; the inline/backticked-heading regression passes.
- [x] Follow the direct-remedy and Finding-Closure Ratchet constraints from UR-083 — the diff is two focused files, adds the named regressions, and introduces no broader workflow machinery.

## Acceptance testing

**Result: Partial**

- GREEN independently rerun on the merged tree: `go test -count=1 ./internal/requeststate -run 'TestCheckpointClaim(RemovalIncludesIndentedContinuationLines|UsesWholeInProgressHeadingLine)$'` — exit 0.
- Adjacent package regression independently rerun: `go test -count=1 ./internal/requeststate` — exit 0.
- RED independently reproduced in an isolated archive of `6e92e536` by reversing only the production `state_apply.go` change while retaining the new tests: both named regressions failed, with the three own-label continuation lines orphaned and the Session Notes content damaged; exit 1.
- P-A-U is complete: all three REQ checkboxes are checked, and the hand-back supplies matching plan, apply, unify, file-manifest, and RED→GREEN evidence.
- `git diff --stat 1832538d..6e92e536` contains only the two declared request-state files (75 insertions, 9 deletions); no debug artifacts or undeclared implementation files were found.
- The request-state acceptance behavior passes, but the restatement/consumer sweep found the cleanup counterpath above, so repository-wide departure acceptance is Partial.

## Restatement sweep

The change redefines checkpoint claim removal from “header line” to “whole entry within the exact section.” Searches across the checkpoint contract, actions, tests, and Go consumers found the canonical documentation aligned with whole-entry removal. The request-state callers all converge on the corrected helper. The cleanup mover is the one stale executable consumer and is recorded as the Important finding above; no stale prose-only restatement was found.

## Scope, decisions, and guardrails

- Scope is exact: the implementation range modifies only `state_apply.go` and `state_apply_test.go`, matching the Implementation Summary and hand-back manifest.
- The hand-back records `## Decisions: None`; the diff introduces no dependency, API shape, or behavior choice requiring an additional decision record.
- The implementation is small, line-oriented, standard-library-only, and consistent with the CLI prime. No speculative abstraction, debug output, TODO/FIXME, naming-for-reach issue, security concern, or unrelated refactor is present.
- Klarna check: the focused tests prove the central request-state fix rather than merely checking compilation, but a green focused package did not prove the repository-wide mover contract; the cleanup sweep finding prevents over-crediting the measurable path.

## Suggested additional testing

- Extend `TestWorkingArchiveRemovesOnlyThisCheckoutCheckpointEntry` in `internal/cleanup` with an own enriched entry plus a foreign enriched entry; assert the cleanup plan deletes the complete own entry and preserves the foreign entry and continuation bytes.

## Scores

**Overall: 75%** (85% dimension average, minus the 10-point Acceptance Partial modifier)

| Dimension | Score | Notes |
|---|---:|---|
| Requirements | 75% | Three of four behavioral requirements are fully closed; whole-entry removal remains partial across working-departure movers. |
| Code Quality | 90% | The merged helper change is clear, bounded, and correct for request-state callers. |
| Test Adequacy | 75% | RED→GREEN is genuine and focused, but no cleanup-mover regression protects the repository-wide contract. |
| Scope | 100% | Exactly the two declared files changed; no drift. |
| Risk | Low | Residual cleanup behavior leaves misleading orphan checkpoint prose but does not remove foreign claims. |
| Acceptance | Partial | Canonical request-state flows pass; cleanup recovery remains stale. |

## Follow-up disposition

- Important finding: `impact-user-visible` — fold into the existing same-root checkpoint-section sweep if still writable; otherwise create one addendum follow-up owned by `internal/cleanup`. Do not split the orphaned fields into separate tasks.
- No follow-up files or do-work state were created or edited by this reviewer.

## Self-validation

Rechecked every REQ and UR-083 constraint against the exact merge range, verified the builder’s claims against the diff, reproduced both RED failures independently, reran GREEN and the full focused package, checked P-A-U, applied the restatement sweep, and recalculated the score mechanically. No additional finding emerged.

*Reviewed independently using the review-work action; review artifact only.*
