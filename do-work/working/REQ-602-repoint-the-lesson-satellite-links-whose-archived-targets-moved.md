---
id: REQ-602
status: claimed
domain: general
created_at: 2026-09-06T08:10:39Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-action-files.md]
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-06T08:13:55Z
  basis:
    - Route A
    - 4-file write set
    - 4 acceptance criteria
tdd: false
maintenance: true
related: [REQ-243, REQ-238]
write_set: [_dev/primes/lessons-shell-commands.md, _dev/primes/lessons-action-files.md, _dev/primes/lessons-kanban-board.md, _dev/tests/audit-lockins.sh]
title: 'Repoint the lesson-satellite links whose archived targets moved, and check satellite links'
claimed_at: 2026-09-06T08:12:01Z
status_changed_at: 2026-09-06T12:48:50Z
---

# Repoint the Lesson-Satellite Links Whose Archived Targets Moved, and Check Satellite Links

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route A. The check first, red on the fifteen; then the repointing; then the index rows.
- [x] **[APPLY]:** Two commits `a70a04f` (the check) and `543564c` (the repointing and index rows), merged as `ef8274b`. Five files: three satellites, `audit-lockins.sh`, `lessons-index.md`.
- [x] **[UNIFY]:** Orchestrator re-ran: at `a70a04f` the lock-ins exit 1 with exactly 15 FAIL lines naming the fifteen; at head and on the merged tree `Audit lock-in regressions passed.`, exit 0. The independent verifier did not run (session limit); qualify and review are owed.

## What

Fifteen lesson bullets across the three maintainer lesson satellites link to an archived REQ at
`do-work/archive/REQ-NNN-….md`. Those files now live under a user-request directory,
`do-work/archive/UR-0xx/`, so every one of the fifteen links is dead. Found while backfilling the
satellites for the REQs this run archived (commit `a38a8c4`); the new lines were checked against the
tree and resolve, the old ones were not and do not.

The fifteen, each with the file the link names and where the target actually is:

  - `lessons-shell-commands.md`: `do-work/archive/REQ-419-add-flat-just-recipes-action-delegation.md` → `do-work/archive/UR-081/REQ-419-add-flat-just-recipes-action-delegation.md`
  - `lessons-shell-commands.md`: `do-work/archive/REQ-420-replace-shell-implementations-verify-parity.md` → `do-work/archive/UR-081/REQ-420-replace-shell-implementations-verify-parity.md`
  - `lessons-action-files.md`: `do-work/archive/REQ-508-reduce-capture-templates-to-schema-backed-examples.md` → `do-work/archive/UR-098/REQ-508-reduce-capture-templates-to-schema-backed-examples.md`
  - `lessons-action-files.md`: `do-work/archive/REQ-459-stage-command-owned-calibration-with-lifecycle-release.md` → `do-work/archive/UR-081/REQ-459-stage-command-owned-calibration-with-lifecycle-release.md`
  - `lessons-action-files.md`: `do-work/archive/REQ-417-implement-interview-memory-commands.md` → `do-work/archive/UR-081/REQ-417-implement-interview-memory-commands.md`
  - `lessons-action-files.md`: `do-work/archive/REQ-409-implement-safe-cleanup.md` → `do-work/archive/UR-081/REQ-409-implement-safe-cleanup.md`
  - `lessons-action-files.md`: `do-work/archive/REQ-410-implement-doctor-forensics.md` → `do-work/archive/UR-081/REQ-410-implement-doctor-forensics.md`
  - `lessons-action-files.md`: `do-work/archive/REQ-435-complete-doctor-forensics-delegation-contract.md` → `do-work/archive/UR-081/REQ-435-complete-doctor-forensics-delegation-contract.md`
  - `lessons-action-files.md`: `do-work/archive/REQ-411-implement-queue-selection.md` → `do-work/archive/UR-081/REQ-411-implement-queue-selection.md`
  - `lessons-action-files.md`: `do-work/archive/REQ-412-implement-request-state-transactions.md` → `do-work/archive/UR-081/REQ-412-implement-request-state-transactions.md`
  - `lessons-action-files.md`: `do-work/archive/REQ-498-make-orchestrator-finalization-resumable.md` → `do-work/archive/UR-096/REQ-498-make-orchestrator-finalization-resumable.md`
  - `lessons-action-files.md`: `do-work/archive/REQ-513-commit-the-claim-footprint-in-every-mode.md` → `do-work/archive/UR-099/REQ-513-commit-the-claim-footprint-in-every-mode.md`
  - `lessons-action-files.md`: `do-work/archive/REQ-461-require-affirmative-project-owned-release-targets.md` → `do-work/archive/UR-081/REQ-461-require-affirmative-project-owned-release-targets.md`
  - `lessons-action-files.md`: `do-work/archive/REQ-570-delete-the-pending-heavy-testing-status-held-requests-stay-claimed.md` → `do-work/archive/UR-114/REQ-570-delete-the-pending-heavy-testing-status-held-requests-stay-claimed.md`
  - `lessons-kanban-board.md`: `do-work/archive/REQ-419-add-flat-just-recipes-action-delegation.md` → `do-work/archive/UR-081/REQ-419-add-flat-just-recipes-action-delegation.md`

## Why

A lesson satellite exists to be read before changing the area, and a bullet's link is the only way to
its evidence. REQ-238's own lesson in this file's neighbour says it: every fix of this class trades a
staleness risk that has a detector for a broken-link risk that has none. REQ-243 added a pointer check
for shipped Markdown; `_dev/primes/` is not shipped and nothing reads these links.

## Detailed Requirements

- Repoint each of the fifteen links to the file that exists, keeping the bullet text and family marker
  unchanged.
- Add one check that every relative link in a maintainer lesson satellite resolves from the satellite's
  own directory, keyed on the condition (a relative link whose target is missing), not on this list; put
  it where the other lock-ins live unless a better home is argued in the record.
- Prove the check red on one of the fifteen before repointing, and green after.
- The shipped satellite `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` uses canonical
  repository URLs, which a local check cannot resolve. Say in the record whether its targets moved the
  same way (compare each URL's path against the tree) and, if any did, fix them in the same pass.

## Constraints

- Prose links and one check. No change to any lesson's wording.
- The lessons index rows (`do-work/lessons-index.md`) carry a token count per satellite; recompute the
  three rows if the byte count changes.

## Open Questions

None.

