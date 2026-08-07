---
id: UR-028
title: Fix secret rename inventory handling
created_at: 2026-08-07T08:05:50Z
requests: [REQ-128]
word_count: 549
---

# Fix Secret Rename Inventory Handling

## Summary

Implement the two accepted review findings for secret renames, retain the do-work intent trail, release the fix as version 0.183.4, and commit the implementation, hash record, and checkpoint separately.

## Full Verbatim Input

do-work validate-feedback: 
Full review comments:

- [P1] Fail closed after a secret rename degrades to D/A — /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/tools/checks/uncommitted-inventory.sh:131-137
  With git mv .env visible-config.txt, this initially emits XD .env and X visible-config.txt, but commit Step 1 requires resetting the pre-staged X destination. Git then reports a deleted .env plus an untracked visible-config.txt; on the required re-inventory this
  branch emits XD and A, causing Steps 2/3 to read the destination's secret contents. The same occurs with rename detection disabled, so quarantine or reject ambiguous secret-deletion/addition pairs across reruns (and mirror that rule in the manual fallback). See
  actions/commit.md:63-69 and CLAUDE.md:124-138 (CLAUDE.md#L124-L138).

- [P2] Accept an already-staged secret deletion — /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/tools/checks/uncommitted-inventory.sh:139-142
  For the normal staged git mv .env visible-config.txt case, resetting the X destination leaves the source deletion already staged. The new XD row then sends commit Step 5 to run git add -u -- .env, but Git rejects that pathspec because the index no longer contains
  .env, blocking the prescribed deletion-only commit. Detect and accept a verified already-staged XD deletion (or unstage it before restaging). See actions/commit.md:143 and CLAUDE.md:124-138 (CLAUDE.md#L124-L138).

fix accepted and commit, including do-work files

PLEASE IMPLEMENT THIS PLAN:
# Fix Secret Rename Handling and Commit REQ-128

## Summary

Implement both accepted findings on the clean `main` branch without rewriting the existing separate commits. Track the work as UR-028 and REQ-128, an addendum to archived REQ-121, and release it as version 0.183.4.

## Implementation

- Make inventory explicitly enable Git rename detection regardless of `status.renames`.
- Preserve secret-derived excluded destinations across every re-inventory in commit and inspect workflows.
- If an inventory degrades to a secret deletion plus untraceable additions, fail closed before reading, associating, or staging any additions.
- Update both manual fallbacks to enforce the same quarantine and ambiguity rules.
- For an `XD` secret deletion, accept an exact already-staged deletion; otherwise stage it with `git add -u`, then verify that the cached change is deletion-only.
- Keep existing inventory tags unchanged and narrow the inaccurate copy-protection claim rather than expanding this fix into content-based copy detection.

## Work Records and Release

- Create UR-028 containing the accepted P1 and P2 feedback.
- Create REQ-128 as a Route B, maintenance, test-driven addendum to REQ-121, recording the reproduction, decisions, changed files, and verification evidence.
- Update the inventory script, commit and inspect actions, contract tests, changelog, version file, and checkpoint.
- Archive UR-028 and REQ-128 together after successful validation.

## Verification

- Add regressions for rename detection disabled, reset-and-reinventory degradation, ambiguous `XD` plus `A`, already-staged secret deletion, ordinary unstaged deletion, and unaffected ordinary additions.
- Demonstrate the new tests failing before the implementation and passing afterward.
- Run shell syntax checks, ShellCheck when available, the contract regression suite, commit-hash guards, update-script tests, queue-kanban Go tests and verification, and `git diff --check`.
- Review the complete diff and confirm no quarantined destination can reach content inspection, REQ association, or staging.

## Commits

- Commit implementation, tests, release metadata, and archived work records as `[REQ-128] Secret rename quarantine survives re-inventory (Route B)`.
- Record that implementation SHA in the archived requirement and commit it as `[REQ-128] Record implementation commit hash`.
- Commit the completed checkpoint as `do-work: session checkpoint — REQ-128 complete, UR-028 closed`.
- Include every resulting do-work file; do not amend, squash, push, or alter the already-created commits.

---
*Captured: 2026-08-07T08:05:50Z*
