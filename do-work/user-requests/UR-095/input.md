---
id: UR-095
title: 'Capture repository-gate dependency recovery'
created_at: 2026-09-01T19:56:26Z
requests: [REQ-491, REQ-492]
word_count: 590
---

# Capture Repository-Gate Dependency Recovery

## Summary

Replace the session-stopping response to unrelated repository-gate failures with a transactional dependency lifecycle. One REQ establishes the canonical deferral primitive and queue ordering; the dependent REQ integrates pre-build deferral, late deferral, and safe resumption into `do-work run`.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-491 | Add the canonical repository-gate deferral transaction, repair-REQ lifecycle, schema fields, evidence, folding, and selector priority. |
| REQ-492 | Integrate repository-gate deferral, repair execution, attribution, and safe resumption into `do-work run`. |

## Batch Constraints

- REQ-492 depends on REQ-491.
- Both REQs are `impact-user-visible`, `effort-substantive`, and `tdd: true`.
- `pending-answers` is never used because repository-gate repair needs no user decision.
- Worktree/merge evidence is mandatory for post-implementation deferral.
- Existing serial dirty claims are not migrated. This feature does not unblock REQ-440; its existing `SC2034` blocker remains on the current manual path.
- Related REQ-469 through REQ-472 describe the superseded `blocked`/`pending-answers` solution shape. This batch records the newer dependency-based design and must not inherit those incompatible semantics.
- Capture only; implementation belongs to a later `do-work run`.

## Full Verbatim Input

> ```
> PLEASE IMPLEMENT THIS PLAN:
> # Capture Repository-Gate Dependency Recovery
> 
> ## Summary
> 
> Capture one complex UR containing two dependency-ordered REQs. The change will make future `do-work run` sessions convert unrelated repository-gate failures into explicit dependency work instead of stopping the entire run.
> 
> ## Captured REQs
> 
> 1. **Add canonical repository-gate deferral lifecycle**
>    - Add a transactional `defer-gate --manifest` CLI operation.
>    - Atomically create or fold a gate-repair REQ, move the current REQ from `working` to `queue`, set it back to `pending`, append the repair to `depends_on`, remove its checkpoint claim, and record gate evidence.
>    - Add frontmatter markers:
>      - `gate_deferred: true`
>      - `repository_gate_repair: true` on generated repair REQs
>      - `deferred_implementation_base` and `deferred_implementation_merge` for late deferrals with existing worktree merge evidence.
>    - Prioritize ready gate-repair REQs first, then ready deferred parents, then ordinary stable queue order.
>    - Use the source REQ’s UR for generated repairs. Fold matching failures into one sweep REQ and let multiple parents depend on it.
> 
> 2. **Integrate deferral and resumption into `do-work run`**
>    - Depends on REQ 1.
>    - Run the canonical repository gate before implementation and again during final testing.
>    - On a failing pre-build baseline, defer the current REQ before source changes exist, run the repair next, and resume the parent after the dependency completes.
>    - A gate-repair REQ may run against its expected red baseline and must not recursively defer itself for the same failure.
>    - For a late failure, rerun the gate at the saved pre-merge revision:
>      - Pre-merge passes: treat it as caused by the current implementation and use the normal remediation loop.
>      - Pre-merge also fails: defer using the persisted implementation merge range.
>    - On late resume, verify the merge remains in ancestry and its implementation paths have not drifted, then rerun qualification, focused tests, the canonical gate, and review. Path drift returns the REQ to implementation instead of trusting stale evidence.
>    - If the repair fails or is cancelled, leave the parent dependency-gated and continue with other ready requests.
> 
> Both REQs are `impact-user-visible`, `effort-substantive`, and `tdd: true`.
> 
> ## Interfaces and Evidence
> 
> - Generated repair REQs are `status: pending`, share the original `user_request`, use `related` rather than `addendum_to`, and carry a root-cause `sweep_key`.
> - The parent receives a `## Repository Gate Deferral` record containing the command, exit status, diagnostic fingerprint, dependency REQ, and saved merge range when applicable.
> - `pending-answers` is never used because no user decision is required.
> - Update the request schema, selector result contract, work-action instructions, CLI prime, and existing contract tests that currently require an unconditional hold.
> 
> ## Test Plan
> 
> - Pre-build gate failure creates one repair dependency and allows the run to continue.
> - Two parents with the same failure fold into one repair sweep.
> - Gate-repair work does not recursively generate itself.
> - Successful repair selects and resumes the deferred parent next.
> - Failed repair leaves the parent gated while unrelated ready work continues.
> - Current-diff gate regressions remain in the remediation path.
> - Late unrelated failures preserve and reuse valid merge evidence.
> - Missing, stale, or drifted merge evidence fails safely.
> - Transaction collision or publication failure rolls back every queue and checkpoint mutation.
> - The full canonical repository gate passes after updating the old hold assertions.
> 
> ## Assumptions
> 
> - Worktree/merge evidence is required for post-implementation deferral.
> - Existing serial dirty claims are not migrated. Therefore this feature will not itself unblock the current REQ-440; its existing `SC2034` blocker must be resolved through the current manual path before these new REQs can run.
> - Capture creates the UR and REQs only; it does not execute them.
> ```

---
*Captured: 2026-09-01T19:56:26Z*
