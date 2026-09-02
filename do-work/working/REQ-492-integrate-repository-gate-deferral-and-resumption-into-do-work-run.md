---
id: REQ-492
title: 'Integrate repository-gate deferral and resumption into do-work run'
status: completed-with-issues
route: C
created_at: 2026-09-01T19:56:26Z
user_request: UR-095
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
required_lessons: [skills/do-work/tools/do-work-cli/lessons-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: [REQ-491]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 105
  confidence: low
  calculated_at: 2026-09-02T02:52:43Z
  basis:
    - Route C
    - 7-file implementation scope
    - orchestration rule change
    - arbitrary repository-gate execution and attribution
    - worktree ancestry and path-drift evidence
    - mutation-tested contract lane
    - full canonical verification
related: [REQ-469, REQ-470, REQ-471, REQ-472, REQ-491]
batch: repository-gate-dependency-recovery
claimed_at: 2026-09-02T02:43:12Z
planning_at: 2026-09-02T02:52:43Z
dispatch_at: 2026-09-02T02:53:26Z
builder_handback_at: 2026-09-02T03:06:44Z
integration_at: 2026-09-02T03:06:44Z
review_at: 2026-09-02T03:14:00Z
remediation_at: 2026-09-02T03:24:10Z
re_review_at: 2026-09-02T03:32:03Z
completed_at: 2026-09-02T03:32:03Z
kb_status: pending
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go
  - skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go
  - skills/do-work/tools/do-work-cli/prime-do-work-cli.md
  - skills/do-work/docs/command-line-guide.md
  - _dev/tests/contract-regressions.sh
---

# Integrate Repository-Gate Deferral and Resumption Into do-work run

## What

Use the canonical lifecycle from REQ-491 inside `do-work run`: establish the repository baseline before implementation, classify late failures against the saved pre-merge revision, run repairs without recursion, and safely resume deferred implementations when their dependency completes.

The duplicate scan found REQ-469 through REQ-472, but those pending REQs prescribe the superseded `blocked`/`pending-answers` lifecycle. They are related context, not fold targets, because this user-confirmed request deliberately replaces those incompatible semantics with `status: pending` plus `depends_on`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Designed the pre-build gate, typed repair closure, late attribution, saved-range resumption, repair no-op, and reservation-cleaned fold contracts with explicit failure-safe branches and mutation-tested consumers.
- [x] **[APPLY]:** Implemented the planned gate baseline, typed repair closure, late attribution/resume, no-op repair, cleaned-reservation fold, and contract mutations in the seven declared files.
- [x] **[UNIFY]:** Audited the exact seven-file diff, preserved REQ-493’s dirty/untracked fold boundary, and passed publication race, CLI vet/tests, contract regressions, gofmt, and diff checks.

## Detailed Requirements

- Run the project-declared canonical repository gate before implementation begins and again during final testing.
- When the pre-build baseline fails, invoke REQ-491's canonical deferral before any source changes exist, select the ready repair next, and resume the parent after the repair dependency completes.
- A REQ marked `repository_gate_repair: true` may run against its expected red baseline and must not recursively defer itself for the same diagnostic fingerprint.
- If a repair fails or is cancelled, leave every parent dependency-gated and continue processing unrelated ready REQs.
- For a gate failure after implementation, rerun the same gate at the saved pre-merge revision.
- When the pre-merge revision passes, classify the failure as caused by the current implementation and use the existing bounded remediation loop; never defer it as unrelated.
- When the pre-merge revision also fails with the matching unrelated fingerprint, defer through REQ-491 and persist `deferred_implementation_base` plus `deferred_implementation_merge`.
- On late resume, prove the saved merge remains in the current ancestry and its implementation paths have not drifted since deferral.
- When ancestry and paths remain valid, reuse the implementation and rerun qualification, focused tests, the canonical repository gate, and independent review before completion.
- When implementation paths drift, return the REQ to implementation and discard stale qualification/test/review trust rather than applying old evidence.
- Missing, stale, malformed, non-ancestor, or otherwise unverifiable merge evidence fails safely; do not archive or trust the deferred implementation.
- Update the work-action instructions, work-reference/schema restatements, CLI prime, selector-result consumption, composed run summaries, and contract tests that currently require an unconditional hold.

## Constraints

- Depends on REQ-491; use its transaction and selector outputs rather than restating a second mutation protocol.
- Gate attribution uses the direct command exit status and a stable diagnostic fingerprint at both current and pre-merge revisions.
- The gate remains mandatory. Deferral changes scheduling and evidence ownership, never the pass requirement for completion.
- Worktree/merge evidence is required for every post-implementation deferral.
- `pending-answers` is never used, and no user confirmation pause is introduced.
- Existing serial dirty claims are not migrated. The current REQ-440 and its `SC2034` blocker remain on the manual path and are not unblocked by this feature.
- Related REQ-469 through REQ-472 are not prerequisites and must not reintroduce `blocked` or `pending-answers` semantics.

## Dependencies

- REQ-491 (Add Canonical Repository-Gate Deferral Lifecycle) — supplies the atomic transition, repair folding, schema fields, evidence model, and selector priority this run integration consumes.

## Builder Guidance

Keep the baseline, attribution, deferral, and resume decisions explicit in the work pipeline. Pin each branch with an isolated semantic contract test so nearby vocabulary cannot satisfy it. Reuse the canonical gate command and REQ-491 result evidence, and treat path-drift detection as an authority check that invalidates stale downstream proof.

## Red-Green Proof

**RED prompt/case:** Run a queue where parent A encounters an unrelated pre-build gate failure and parent B is independently ready; today's unconditional-hold assertions stop the session, no repair dependency is selected, and A cannot resume through a verified lifecycle.
**Why RED now:** The work instructions do not own a pre-build gate/repair path, cannot distinguish a late current-diff regression from a pre-existing failure at the saved base, and have no safe late-resume ancestry/path-drift proof.
**GREEN when:** Contract and behavioral lanes prove: pre-build deferral continues the run; two parents share one repair; repair work does not recurse; successful repair resumes the deferred parent next; failed/cancelled repair leaves parents gated while unrelated work runs; current-diff failures remediate; valid late merge evidence is reused; invalid or drifted evidence fails safely; and the updated full canonical repository gate exits zero.
**Validation:** User confirmed the supplied lifecycle and test plan in this session.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3337 tokens; highest semantic match for work-pipeline and downstream-reader changes, but the satellite exceeds the 2000-token per-REQ budget and is only partially family-slugged.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens; relevant to canonical-gate command ownership and attribution, but the satellite exceeds the budget and is only partially family-slugged.

## Full Context

See `do-work/user-requests/UR-095/input.md` for complete verbatim input.

---
*Source: UR-095 — "Integrate deferral and resumption into do-work run"*

## Triage

**Route: C** - Complex

**Reasoning:** This request changes the canonical work orchestrator’s gate, deferral, failure-attribution, resumption, queue-continuation, and evidence contracts across several documented consumers. It requires an explicit plan, repository exploration, isolated semantic contract tests, and independent review.

**Planning:** Required

## Plan

1. Insert a named repository-gate baseline phase after preflight and before any Step 6 source edit. Ordinary red baselines defer through the canonical `defer-gate` transaction; repair REQs may implement only their exact recorded red fingerprints and never recursively defer; an already-green repair records a reviewed no-change completion so its parents can resume.
2. Allocate create candidates with capture’s read-only max-request/reservation scan, let `defer-gate` create the reservation atomically, and rescan/retry only fully rolled-back ID collisions. Make fold authority rely on the exact existing repair REQ preimage rather than a reservation marker cleanup legitimately removes.
3. After deferral, consume only typed `gate_deferral` result fields, add the returned repair ID to the active session closure, suppress the deferred parent even for explicit targeting, and recompute selection. Failed/cancelled/gated repairs leave parents gated while unrelated scoped/default work continues.
4. On a late worktree gate failure, run the identical structured gate argv in a detached diagnostic worktree at saved `<pre>`. A passing base routes to bounded current-diff remediation; the exact matching red fingerprint routes to `defer-gate` with both base and merge commits; mismatch, launch failure, missing range, or serial dirty implementation fails safely. Successful late deferral performs the skipped non-force builder-worktree cleanup.
5. For a late deferred parent, resolve and validate the non-empty base-to-merge range, require merge ancestry in current `HEAD`, derive rename-aware implementation paths, and reject any path-history or working-state drift. No drift reuses the range only after fresh qualification, focused tests, canonical gate, and independent review; drift clears stale pointers/evidence and returns to implementation.
6. Rewrite the old unconditional-hold mutation lane and all active restatements. Prove pre-build deferral, non-recursion, cross-UR repair closure, explicit-parent suppression, repair failure continuation, no-op repair completion, late attribution, valid resume, drift/malformed evidence, non-force cleanup, composed summaries, and reservation-cleaned folds.

**Consumer contract:** The work action executes and fingerprints the project-owned gate, proposes caller-authored manifest evidence, consumes the typed selector priority and `gate_deferral` result, and owns scheduling/attribution/path-drift judgment. The CLI remains the sole atomic mutation owner; no text-render scraping or free-form queue edit is permitted.

**Plan validation:** Every Detailed Requirement maps to tasks 1–6 and every task maps to baseline, scheduling, attribution, resume, fold authority, or verification requirements. REQ-493 remains the separate owner of tracked-dirty repair folds and occupied parent-destination preflight.

*Generated by Plan agent*

## Exploration

The executable seams are `work.md` selection, preflight, dispatch, qualification, gate, cleanup, loop, and error-summary branches; `work-reference.md` owns the full algorithm and schema/consumer restatements. REQ-491 already provides typed priority and deferral results, so no selector/result/schema changes are needed.

One blocking integration defect exists in REQ-491’s fold path: SessionStart legitimately removes a committed repair’s reservation marker, but fold validation still requires it. Creation should continue to scan request/reservation evidence, propose max+1 without writing, and let the transaction exclusively create the unpadded reservation; `queue-kanban next-req` is incompatible because it pre-creates a differently shaped marker. Fold mode instead trusts the exact repair REQ preimage and accepts an absent cleaned marker.

Targeted runs need session-local parent suppression because explicit REQs bypass dependencies, and folded repairs may belong to another UR, so only the typed returned repair ID can extend the closure. Late deferral skips normal Step 8 cleanup, requiring immediate non-force builder worktree cleanup. Serial post-implementation deferral remains fail-safe because it lacks isolated base/merge authority.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go`
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/docs/command-line-guide.md`
- `_dev/tests/contract-regressions.sh`

**Acceptance criteria:** The canonical work instructions can operationally defer a pre-build or attributable late unrelated gate failure, schedule and run the exact repair without recursion or explicit-target livelock, continue unrelated work after repair failure, validate or invalidate saved implementation evidence, and resume only after fresh proof. A committed pending repair remains foldable after reservation cleanup. All superseded unconditional-hold/blocked/pending-answers restatements are removed or explicitly limited to the serial/manual fallback.

## Pre-Flight

**Git:** Clean outside the canonical REQ-492 claim/checkpoint and the two preserved untracked `.DS_Store` files.

**Tests:** Direct unpiped `bash _dev/tests/maintainer-verify.sh` passed after claim and before source implementation, including contracts, installer behavior, board vet/tests/strict JavaScript, CLI vet, and uncached CLI tests. The strict browser lane was skipped because no browser is configured, as allowed by the gate.

**Dependencies:** REQ-491 is completed and its implementation/metadata commits are recorded. REQ-493 remains pending and is deliberately excluded from this scope.

## Root Cause

The work action treated every unrelated repository-gate failure as a session-wide hold. Although REQ-491 now supplies an atomic deferral primitive and priority evidence, the orchestrator did not execute a pre-build baseline, schedule typed repair closure, attribute late failures at an isolated base, validate saved work before resume, or tolerate cleaned fold reservations.

## Decisions

- **D-01 — KEEP:** Gate execution, semantic fingerprinting, revision attribution, path-drift judgment, and run scheduling remain action-owned; the CLI owns only typed atomic mutations and evidence projection.
- **D-02 — KEEP:** Explicit targeting never overrides the session’s just-created deferral closure; the parent is locally suppressed until repair terminal success, while the returned repair ID may cross the original UR boundary.
- **D-03 — KEEP:** A repair whose expected gate is already green completes through an explicitly reviewed no-change path, not cancellation or hollow implementation.
- **D-04 — KEEP:** Serial post-implementation failures without an isolated base/merge pair keep the fail-safe hold; only pre-build deferral is universal.
- **D-05 — KEEP:** Fold authority is the exact existing repair REQ preimage, not a reservation marker whose documented lifecycle has already ended.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)
- `skills/do-work/docs/command-line-guide.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Integrated repository-gate baselining, typed repair scheduling, explicit-target suppression, cross-UR repair closure, non-recursive repair failure, late base attribution, non-force cleanup, ancestry/path-history drift validation, and evidence-revalidation into `do-work run`. Replaced the retired unconditional-hold contract lane. Updated `defer-gate` so an absent cleaned reservation is accepted only for an exact pending repair whose bytes are clean and committed at `HEAD`; present-reservation untracked folds remain supported, while absent-reservation untracked or tracked-dirty folds remain REQ-493. The documented already-green no-op and collision-retry branches remain incomplete at downstream/runtime seams and are preserved in REQ-494/REQ-495 rather than claimed complete.

## Review

**Verdict:** Request changes (50%, acceptance fail; critical orchestration risk).

- **Critical — impact-critical:** Final-gate attribution lacks the pre-build repair exception. An ineffective repair can reproduce its own fingerprint at current and base revisions and recursively invoke `defer-gate`; a different fingerprint stops the session instead of terminally failing the repair and continuing unrelated work.
- **Important — impact-user-visible:** The already-green repair branch names a reviewed no-change completion but does not define exceptions to empty-implementation, no-diff review, release, or successful-staging guards, so it cannot actually reach terminal success and release its parents.
- **Important — impact-rule-change:** Collision retry requires “complete rollback,” but a normal planning collision is a pre-mutation refusal with no rollback result. The accepted typed states must be pre-mutation refused/zero-mutation or rolled-back/succeeded; incomplete or committed-risk outcomes must stop.
- **Minor:** The implementation summary must distinguish absent-reservation untracked/dirty folds (REQ-493) from present-reservation untracked folds, which remain supported.

Completion is on hold for the one permitted remediation pass. Required tests must prove terminal non-recursive repair failure plus unrelated continuation, executable reviewed no-change completion and parent resume, and both safe collision retry result shapes.

## Remediation

The one permitted remediation pass closed final-gate repair recursion and terminal failure/continuation, added documented no-op evidence and lifecycle exceptions, tightened collision-retry predicates, and corrected the fold boundary. Publication race, full CLI, and mutation contracts passed.

Fresh re-review found that the no-op exceptions stop short of two downstream authorities: generated repairs remain `tdd: true` under an unconditional RED/GREEN guard, and `review-work.md` still exits on an empty implementation diff. It also found that the documented collision retry states do not match actual publication wire results: planning refusals carry an empty rollback status, while generic post-mutation failures do not preserve collision identity. Because the remediation allowance is exhausted, these are mandatory follow-ups rather than a second inline patch.

## Re-Review

**Verdict:** Request changes (50%, acceptance fail; critical risk).

- **Critical — impact-critical:** Already-green repair completion remains unreachable through the unconditional TDD and no-diff review guards. REQ-494 owns the multi-action guard closure.
- **Important — impact-user-visible:** Collision retry predicates do not match emitted planning-refusal or post-mutation result shapes. REQ-495 owns the typed collision-result contract.

Closed and verified: final-gate repair failures fail/archive without recursion and continue unrelated work; reservation-cleaned clean committed folds work; present-reservation untracked folds remain supported; absent-reservation untracked/dirty folds remain REQ-493; late attribution, ancestry/path drift, suppression, cross-UR closure, cleanup, and summaries are documented and mutation-tested.

## Qualification

**Result:** Mechanical qualification and scope drift pass for the declared seven files. Acceptance remains fail because two downstream/runtime consumers were outside the declared diff and contradict the new prose path.

## Testing

**Red-green validation:** `TestDeferGateFoldAcceptsCommittedRepairAfterReservationCleanup` failed RED with `DEFER-GATE-RESERVATION-STALE` and passes GREEN. Focused publication race, `go vet ./...`, uncached full CLI tests, and full contract regressions pass. Mutation tests prove final repair non-recursion/terminal failure, unrelated continuation, gate ordering, typed scheduling, late attribution/resume, and documented no-op/collision branches, but do not execute the two contradictory downstream/runtime seams captured by REQ-494 and REQ-495.

## Lessons Learned

**What worked:** Fresh review across downstream actions caught orchestration contradictions that same-file lexical mutation tests could not; keeping REQ-493’s topology boundary explicit prevented accidental scope absorption.

**What didn't:** Adding a special path only to `work.md` and its reference is insufficient when qualification, review, result projection, and commit validation are separate authorities. Prose predicates cannot substitute for executable typed-result fixtures.

**Worth knowing:** Any lifecycle exception must be swept through every downstream gate that can reject it, and retry contracts must be defined from actual wire values rather than desired semantic labels.

## Orientation

[MAP CHANGED] `do-work run` now has an explicit repository-gate deferral/resumption lifecycle spanning baseline, repair scheduling, late attribution, and saved-range validation. The core flow and cleaned-reservation fold are present, but UR-095 is not ship-ready until REQ-494 and REQ-495 close the two remaining executable seams; REQ-493 separately owns the outstanding fold/destination topology class.
