---
id: UR-087
title: 'Non-blocking REQ-level work orchestration'
created_at: 2026-09-01T04:29:16Z
requests: [REQ-468, REQ-469, REQ-470, REQ-471, REQ-472]
word_count: 533
---

# Non-Blocking REQ-Level Work Orchestration

## Summary

Make `do-work run` continue past any single REQ that cannot advance: route blockers by cause, replace the unrelated canonical-gate hold (which today preserves the claim and stops the session) with a blocked set-aside that releases the REQ from `working/` and continues, isolate every implementation on a per-REQ branch/worktree (serial runs included) so set-aside is safe after edits, preserve enough durable state to resume, re-probe blocked conditions, and apply the behavior consistently across all flows, readers, docs, and regression tests. Never waive the gate.

The invocation arrived as `do-work validate-feedback:` wrapping this capture. Triage verified every claim against the code: the hold exists verbatim at `skills/do-work/actions/work.md` Step 6.5 item 4 and the Error Handling table, is contract-pinned by `_dev/tests/contract-regressions.sh` (canonical-gate lane), and additionally collides with crash recovery (a held claim is destructively reset on the next session). All "retain" items (blocker taxonomy, dependency gating, specialized blocked statuses, failure flow, re-probe mechanism, one-releaser semantics) already exist and are asserted, not rebuilt.

## Extracted Requests

| REQ | Title | Depends on |
|---|---|---|
| REQ-468 | Per-REQ branch/worktree isolation for all implementation, serial included | — |
| REQ-469 | Replace the unrelated canonical-gate hold with a blocked set-aside | REQ-468 |
| REQ-470 | Fold-first routing of unrelated gate failures into pending-answers REQs | REQ-469 |
| REQ-471 | Flow and reader consistency plus documentation for gate-blocked set-aside | REQ-469 |
| REQ-472 | End-to-end regression scenarios for non-blocking orchestration | REQ-469, REQ-470, REQ-471 |

## Batch Constraints

- Update the authoritative orchestration instructions, schema/restatements, board/readers, documentation, and regression tests together — `_dev/tests/contract-regressions.sh` predicates that pin edited prose change in the same commit as that prose.
- Never waive the canonical gate: a blocked implementation cannot be archived, versioned, or released until the gate passes.
- Keep one releaser per queue; integration stays serial. Preserve explicit per-REQ commits, archive rules, changelog/version behavior, scope evidence, and UR closure semantics.
- The user decided interactively during capture: isolation is always-on per-REQ branch/worktree for implementation, serial runs included (not a set-aside-time-only mechanism).
- A gate failure caused by the current REQ still follows the existing remediation and code-failure handling — only unrelated/pre-existing failures take the set-aside path.
- Search for and update stale language saying an unrelated canonical-gate failure must preserve a claim and stop the session (known sites: `skills/do-work/actions/work.md` Step 6.5 item 4, its Step 6.5 checklist line, the Error Handling table row; the contract-regressions gate lane). Archived REQs are immutable provenance and are not edited; `CHANGELOG.md` stays owner-only.

## Full Verbatim Input

> ````text
> do-work validate-feedback: do-work capture-request: Make the work orchestration non-blocking at the REQ level.
> 
> Goal
> 
> `do-work run` must continue processing other runnable REQs whenever one REQ cannot advance. A blocker affecting one REQ must be recorded durably, released from `working/`, and set aside. Only an orchestration-wide integrity failure may stop the whole run.
> 
> Required behavior
> 
> - Route blockers by cause:
>   - User decision required → `pending-answers`.
>   - External or repository condition → `blocked` with `blocked_by` and `blocked_at`.
>   - Unmet REQ dependency → retain dependency gating.
>   - Archive collision or dependency cycle → retain the existing specialized blocked status.
>   - Intent/spec/code failure → retain the existing failure and follow-up flow.
> - Replace the unrelated/pre-existing canonical-gate hold that currently leaves a REQ claimed and stops the session:
>   - Mark the active REQ `blocked`.
>   - Record the gate command, failing tests, evidence that the failures are unrelated to the REQ, and the phase where it stopped.
>   - Remove `claimed_at`, remove the current checkout’s checkpoint entry, move the REQ back to `do-work/queue/`, and continue with the next runnable REQ.
>   - Run the fold-first scan for each unrelated gate failure. Append to a matching queued REQ when possible; otherwise create a non-critical `pending-answers` REQ so the user can approve or reject fixing it.
> - Never waive the canonical gate. A blocked implementation cannot be archived, versioned, or released until the gate passes.
> - Make setting aside safe after implementation edits. A blocked REQ’s code, tests, decisions, and evidence must remain durable without contaminating the next REQ’s diff, tests, qualification, staging, or commit. Use per-REQ branch/worktree isolation for implementation, including serial runs, or an equally safe durable mechanism.
> - Record enough durable state to resume the same implementation after its condition clears. Do not discard completed work or restart from scratch unnecessarily.
> - Re-probe blocked conditions when possible. When the condition clears, return the REQ to a runnable state and resume its qualification/release pipeline.
> - Apply the behavior consistently to default, targeted, wave, fan-out, crash-recovery, checkpoint, cleanup, roadmap, clarify, and composed-summary flows.
> - Keep one releaser per queue and preserve explicit per-REQ commits, archive rules, changelog/version behavior, scope evidence, and UR closure semantics.
> 
> Acceptance tests
> 
> - REQ-A encounters an unrelated canonical-gate failure while REQ-B is pending: A becomes blocked and leaves `working/`; B is subsequently claimed and processed.
> - The unrelated failures create or fold into `pending-answers` REQs instead of being stored only in `CHECKPOINT.md`.
> - A blocked REQ with implementation edits cannot affect another REQ’s diff, qualification, tests, staging, or commit.
> - A blocked REQ’s checkpoint entry is removed when it leaves `working/`.
> - When the gate becomes green, the blocked REQ resumes from its preserved implementation and can complete normally.
> - A gate failure caused by the current REQ still follows remediation and code-failure handling rather than being misclassified as unrelated.
> - Queue summaries and roadmap output clearly distinguish blocked work, pending user decisions, dependency-gated work, and runnable work.
> - Regression tests cover serial and fan-out execution, crash recovery, repeated runs, and UR closure with blocked members.
> 
> Update the authoritative orchestration instructions, schema/restatements, board/readers, documentation, and regression tests together. Search for stale language that says an unrelated canonical-gate failure must preserve a claim and stop the entire session.
> ````

---
*Captured: 2026-09-01T04:29:16Z*
