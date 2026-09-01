---
id: REQ-463
title: '[impact-critical] Review fix: Require committed fresh authority for reservation cleanup'
status: pending
created_at: 2026-09-01T02:23:45Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-415]
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-414, REQ-457]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-414
sweep: true
sweep_key: reservation-cleanup-committed-fresh-authority
---

# Require Committed Fresh Authority for Reservation Cleanup

## What

Delete a REQ reservation only from committed request evidence when inside any Git worktree and only after identity and eligibility are revalidated immediately before deletion. Done means uncommitted files and stale observations can never authorize removal of in-flight coordination state.

REQ-415 consumes this cleanup operation from core SessionStart and has owner-approved scope to implement this exact critical closure. This dependency prevents a parallel write collision; after REQ-415 review, disposition this sweep from the same regression evidence rather than duplicating the fix.

The fold-first scan found no eligible pending or pending-answers REQ, sweep or otherwise, in any UR that shares this reservation-authority root cause. REQ-457 governs rollback ownership of transaction-created destinations, not the committed-evidence and fresh-eligibility predicate for reservation cleanup.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] `internal/corehelpers/reservations.go`: failure of `git ls-tree HEAD` in an unborn repository falls back to uncommitted on-disk REQs, allowing an uncommitted REQ-203 to delete reservation REQ-000203. (found by REQ-414 / UR-081)
- [ ] `internal/corehelpers/reservations.go`: the second stat checks identity but does not recompute age and claimed eligibility immediately before deletion. (found by REQ-414 / UR-081)

## Requirements

- Distinguish not-a-Git-worktree from a Git worktree whose HEAD is unborn or otherwise unreadable; inside Git, only committed request evidence may authorize landed-marker cleanup.
- Revalidate marker identity, age threshold, and committed-request eligibility from the final observation immediately before deletion.
- Preserve the marker on any ambiguity, concurrent change, Git evidence failure, or stale eligibility transition and return actionable typed evidence.
- Add deterministic fixtures for unborn repositories, uncommitted matching REQs, final-stat mtime/identity changes, and a genuinely eligible committed marker.

## Red-Green Proof

**RED prompt/case:** Initialize a repository without a commit, add an uncommitted REQ file matching an old reservation, and run cleanup; separately change marker eligibility between the first and final stat.
**Why RED now:** On-disk fallback and first-observation eligibility can authorize deletion without committed or fresh evidence.
**GREEN when:** Both adversarial markers remain byte-identical while a marker backed by committed evidence and unchanged final eligibility is deleted.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-414-rereview.md`.

---
*Source: REQ-414 fresh re-review finding 2.*
