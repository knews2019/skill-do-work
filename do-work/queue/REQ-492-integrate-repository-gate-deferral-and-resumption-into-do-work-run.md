---
id: REQ-492
title: 'Integrate repository-gate deferral and resumption into do-work run'
status: pending
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
related: [REQ-469, REQ-470, REQ-471, REQ-472, REQ-491]
batch: repository-gate-dependency-recovery
---

# Integrate Repository-Gate Deferral and Resumption Into do-work run

## What

Use the canonical lifecycle from REQ-491 inside `do-work run`: establish the repository baseline before implementation, classify late failures against the saved pre-merge revision, run repairs without recursion, and safely resume deferred implementations when their dependency completes.

The duplicate scan found REQ-469 through REQ-472, but those pending REQs prescribe the superseded `blocked`/`pending-answers` lifecycle. They are related context, not fold targets, because this user-confirmed request deliberately replaces those incompatible semantics with `status: pending` plus `depends_on`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files`, `required_lessons`, and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
