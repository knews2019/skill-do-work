---
id: REQ-602
status: completed
domain: general
created_at: 2026-09-06T08:10:39Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-action-files.md]
route: A
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
completed_at: 2026-09-06T12:50:24Z
commit: ef8274bef8ea83c6961b6ad3d1d12848c011c5e8
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

## Triage

**Route: A** — Direct build.

**Reasoning:** Fifteen link targets are known and listed with where each file now lives; the check is
one condition over three files whose shape `audit-lockins.sh` already has a dozen examples of. Nothing to
explore: the dead links were found by resolving every relative link in the three satellites against the
tree, and the list is that resolution's output.

**Planning:** Skipped.

**The check goes in before the repointing so the red is real.** Fifteen dead links are the fixture; a
check written after they are fixed has never failed.

## Plan

**Planning not required** — Route A: one check, fifteen link edits, three index rows.

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `_dev/tests/audit-lockins.sh` (modified)
- `_dev/primes/lessons-action-files.md` (modified)
- `_dev/primes/lessons-kanban-board.md` (modified)
- `_dev/primes/lessons-shell-commands.md` (modified)
- `do-work/lessons-index.md` (modified)

**The check, in its own commit first.** A new block in `audit-lockins.sh` headed "Lesson-satellite links
resolve (REQ-602)": a satellite is any `_dev/primes/lessons-*.md` (a glob, with a FAIL when it matches
nothing); a relative link is any `](target)` whose target is non-empty, does not start with `#` and is
not an absolute URL; the target, fragment stripped, must resolve from `_dev/primes/`. One FAIL line per
broken link, naming the satellite, the line and the target as written. At `a70a04f`, with the links
still broken, the lock-ins exit 1 with exactly the fifteen FAIL lines the request lists and no other
(re-run by the orchestrator after the builder).

**The repointing.** Fifteen links in three satellites now name the `UR-0xx/` directory each target moved
under; bullet text and family markers unchanged. Three `lessons-index.md` rows recomputed by program
(tokens and families). Lock-ins green at head and on the merged tree.

**The shipped satellite's canonical URLs** were reported, not edited: see the hand-back's URL report
(`do-work/runs/work-2026-09-05-231943/REQ-602-handback.md`).

**Owed to the next session:** the independent verifier and the review did not run (session limit);
canonical `qualify` over the merge range, the three-lens review, and finalization (not a release: no
shipped file changed).

## Decisions

- **D1 The check lives in `audit-lockins.sh`**, beside the other lock-ins the gate already runs, rather
  than in a fifth script.
- **D2 The satellite bullets' text is untouched**, only the link targets: the lesson wording is the
  archive's, not this request's.

## Qualification

**Passed.** Merged range `66e9992f..ef8274be`, 5 files, 69 insertions and 18 deletions.
- Canonical `qualify` satisfied on the merged range.
- `_dev/tests/audit-lockins.sh` added the satellite relative link check block.
- Fifteen satellite links in `_dev/primes/lessons-action-files.md`, `_dev/primes/lessons-kanban-board.md`, and `_dev/primes/lessons-shell-commands.md` repointed to their respective `UR-0xx/` directories.
- `do-work/lessons-index.md` rows updated with recalculated token counts.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` canonical URLs audited; all 43 URLs resolve to tracked files.

## Testing

**Red-Green Evidence:**
- **RED:** At commit `a70a04f`, `bash _dev/tests/audit-lockins.sh` failed with exit 1, printing exactly 15 FAIL lines for the fifteen unresolved links across the three maintainer satellites.
- **GREEN:** At commit `543564c` and on the merged tree `ef8274b`, `bash _dev/tests/audit-lockins.sh` exited 0 (`Audit lock-in regressions passed.`).
- **Mutation:** Mutating one link in `_dev/primes/lessons-kanban-board.md` back to the un-repointed path reproduced the failure (exit 1), and restoring it returned to exit 0.
- **Guards:** `_dev/tests/quiet-grep-pipeline-audit.sh` and `_dev/tests/action-shell-blocks.sh` exit 0.
- Fast gate on the merged tree passed cleanly (exit 0).

## Review

**Overall: 95%** | 2026-09-06T13:00:00Z | Synthesis of adversarial verification lenses

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Verdict: Pass.** The implementation is clean, minimal, and fully verified.
- The satellite link check in `_dev/tests/audit-lockins.sh` is condition-based rather than an enumerated list of files, handles fragments and external URLs cleanly, and avoids SIGPIPE traps using herestrings.
- Red/green evidence reproduced: exactly 15 FAIL lines at `a70a04f`, exit 0 at `543564c` and on the merged tree `ef8274b`.
- The fifteen links in the three maintainer satellites were repointed to their existing `UR-0xx/` targets without modifying bullet text or family slugs.
- All 43 canonical URLs in the shipped satellite `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` resolve to tracked files.
- `do-work/lessons-index.md` token counts updated to match the recomputed byte lengths.
- Scope is clean: exactly the 5 declared files, no changes under `skills/` (non-release).

## Remediation

None needed; all requirements satisfied and verified green.

## Lessons Learned

**What worked:**
- Committing the check first (`a70a04f`) to establish authentic RED evidence before applying the repointing fixes (`543564c`).
- Using condition-based matching (`_dev/primes/lessons-*.md` and regex-based relative link matching) rather than hardcoded file lists, ensuring future maintainer satellites automatically inherit link integrity checks.

**What didn't:**
- Initial backfilling of lesson satellites omitted link validation against repository moves, allowing broken relative links to accumulate silently.

**Worth knowing:**
- When archiving REQs into nested UR folders, cross-references and satellite relative links must be updated or checked by an automated gate.

## Orientation

`_dev/tests/audit-lockins.sh` now enforces that all relative links in maintainer lesson satellites (`_dev/primes/lessons-*.md`) resolve to real files on disk. The 15 broken links in `lessons-action-files.md`, `lessons-kanban-board.md`, and `lessons-shell-commands.md` have been repointed to their `UR-0xx/` locations, and index token counts in `do-work/lessons-index.md` are synchronized. Subsystem: maintainer lock-in tests and prime documentation.
