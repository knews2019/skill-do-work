---
id: REQ-491
title: 'Add canonical repository-gate deferral lifecycle'
status: pending
created_at: 2026-09-01T19:56:26Z
user_request: UR-095
domain: backend
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
required_lessons: [skills/do-work/tools/do-work-cli/lessons-do-work-cli.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-469, REQ-470, REQ-471, REQ-472, REQ-492]
batch: repository-gate-dependency-recovery
---

# Add Canonical Repository-Gate Deferral Lifecycle

## What

Add one transactional `defer-gate --manifest` operation that converts an unrelated repository-gate failure into explicit repair work and safely returns the parent REQ to the dependency-gated queue. Extend the request model and selector so repair work and resumed parents run in the intended order without stopping unrelated work.

The duplicate scan found REQ-469 through REQ-472, but those pending REQs prescribe the superseded `blocked`/`pending-answers` lifecycle. They are related context, not fold targets, because this user-confirmed request deliberately replaces those incompatible semantics with `status: pending` plus `depends_on`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files`, `required_lessons`, and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- Add a canonical transactional CLI operation named `defer-gate --manifest`.
- In one atomic publication boundary, create or fold a repository-gate repair REQ, move the active parent from `do-work/working/` back to `do-work/queue/`, set the parent to `status: pending`, append the repair REQ to `depends_on`, remove the parent's checkpoint claim, and append durable gate evidence.
- Add and normalize the request frontmatter markers `gate_deferred: true`, `repository_gate_repair: true`, `deferred_implementation_base`, and `deferred_implementation_merge`.
- Generated repair REQs are `status: pending`, use the source REQ's `user_request`, use `related` rather than `addendum_to`, carry `repository_gate_repair: true`, and carry a root-cause `sweep_key`.
- Matching diagnostic fingerprints fold into one repair sweep REQ, allowing multiple parent REQs to depend on the same repair.
- A parent receives a `## Repository Gate Deferral` record containing the gate command, direct exit status, diagnostic fingerprint, dependency REQ, and the saved implementation merge range when applicable.
- `pending-answers` is never part of this lifecycle because the repair requires no user decision.
- Extend selector ordering to choose ready `repository_gate_repair` REQs first, then ready `gate_deferred` parents, then ordinary ready REQs in existing stable queue order.
- Update the request schema and lossless normalized projection for all new fields, the selector result contract and evidence, the CLI prime, and every alternate reader/writer affected by the changed status and dependency contracts.
- The transaction must refuse missing or invalid inputs, collisions, stale preimages, and unsafe publication topology before mutation; any failure after mutation begins must roll back every queue, REQ, reservation, and checkpoint change.

## Constraints

- This is a dependency relationship: the parent remains `pending` and names the repair in `depends_on`; do not introduce a repository-gate-specific blocked status.
- Fold by stable root-cause identity (`sweep_key` plus diagnostic evidence), never by title similarity alone.
- Preserve the existing stable order within each selector priority class.
- Multiple parents may depend on one repair sweep, and folding must not overwrite prior parents or evidence.
- The source REQ's UR remains the repair REQ's UR. Do not create an addendum relationship.
- The parent may carry `deferred_implementation_base` and `deferred_implementation_merge` only when late deferral has valid worktree/merge evidence.
- Existing serial dirty claims are outside the migration boundary; REQ-440 remains on its manual recovery path.
- Related REQ-469 through REQ-472 are not prerequisites and must not reintroduce `blocked` or `pending-answers` semantics.

## Dependencies

None. REQ-492 depends on this canonical transaction and selector contract.

## Builder Guidance

Treat the publication operation as the single mutation owner. Model the parent transition, repair creation/fold, checkpoint removal, and evidence append as one planned transaction with explicit preimages and rollback identity. Follow the request-model lesson that every new frontmatter field needs explicit typed and normalized projection, and sweep every downstream reader rather than stopping at the schema declaration.

## Red-Green Proof

**RED prompt/case:** With a claimed parent REQ and a reproducible unrelated gate fingerprint, invoke the proposed `defer-gate --manifest`; today the CLI has no such operation, no atomic parent/repair/checkpoint transition exists, and queue selection has no repair/deferred priority classes.
**Why RED now:** The current contract either holds the claim or routes through the older blocked/pending-answers design, so one unrelated red gate can stop the run and there is no canonical dependency record to resume from.
**GREEN when:** Focused CLI tests prove one successful atomic deferral, same-fingerprint folding across two parents, repair-first then deferred-parent selector order, lossless schema projection, and complete rollback for collision and publication failure. The existing stable order remains unchanged outside the two new priority classes.
**Validation:** User confirmed the supplied lifecycle and test plan in this session.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 3337 tokens; relevant because status contracts and alternate readers change, but the satellite exceeds the 2000-token per-REQ budget and is only partially family-slugged.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens; relevant because a prescribed CLI transaction and gate commands change, but the satellite exceeds the budget and is only partially family-slugged.

## Full Context

See `do-work/user-requests/UR-095/input.md` for complete verbatim input.

---
*Source: UR-095 — "Add canonical repository-gate deferral lifecycle"*
