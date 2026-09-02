---
id: UR-096
title: 'Make do-work finalization resumable'
created_at: 2026-09-02T13:07:19Z
requests: [REQ-498]
word_count: 910
---

# Make Do-Work Finalization Resumable

## Summary

Make the release tail a journaled, resumable CLI transaction; recover journaled and safely attributable legacy tails before selecting new work.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-498 | Implement resumable orchestrator finalization and startup recovery |

## Batch Constraints

- Preserve the existing single-releaser model.
- Keep journals Git-private.
- Recover only unambiguous legacy state; never guess shared metadata ownership.
- Capture the intent, finish the current implementation at a safe vertical slice, commit it, and stop.

## Full Verbatim Input

> ```
> PLEASE IMPLEMENT THIS PLAN:
> # Resumable Orchestrator Finalization
> 
> ## Summary
> 
> The failure occurs in a gap between archiving a completed REQ and committing its release metadata. Archiving removes the REQ from `working/` and its checkpoint claim, while current recovery only scans working REQs. The abandoned tail therefore becomes invisible, but its dirty `CHECKPOINT.md` blocks every subsequent claim. Trying another ready REQ would not help because all claims mutate that shared checkpoint.
> 
> Fix this by making finalization a CLI-owned, journaled transaction and recovering unfinished finalizations before selecting another REQ.
> 
> ## Implementation Changes
> 
> ### Transactional finalization
> 
> - Add `do-work-cli finalize --manifest <path> --format json`.
> - The manifest identifies the REQ, writer, terminal result, expected lifecycle preimages, exact commit allowlist, optional release manifest, commit message, and implementation-provenance mode.
> - Before changing tracked files, validate the complete plan and atomically persist it under `git rev-parse --git-path do-work-finalization/REQ-NNN.json`.
> - Compose the existing completion, release, and Git transaction implementations rather than duplicating their rules.
> - Replace the separate archive/release/manual-staging sequence with these resumable phases:
>   1. Prepared and validated
>   2. Completion/archive applied
>   3. Release applied, when applicable
>   4. Primary commit created
>   5. Commit provenance recorded in the metadata commit
>   6. Final state verified and journal removed
> - Store sufficient preimages and target digests to identify whether each phase is unapplied, applied, or conflicting.
> - Make every phase idempotent. In particular, retries must not duplicate changelog entries, version increments, calibration rows, archive moves, or commits.
> - Before the primary commit, recoverable failures restore CLI-owned lifecycle/release/index changes when possible without discarding implementation edits. After the primary commit, recovery always rolls forward to the provenance commit.
> - Serial runs use the primary finalization commit as implementation provenance. Worktree runs accept the already-known merge/implementation hash and commit only lifecycle and release state.
> 
> ### Automatic startup recovery
> 
> - Add `do-work-cli recover-finalization --discover --format json`.
> - Run it at the start of `do-work run`, before ordinary working-REQ crash recovery, selection, or claim.
> - Resume journaled finalizations deterministically, oldest first. A journal whose final commits already exist is verified and cleaned up without creating another commit.
> - After journal recovery, discover legacy unjournaled release tails such as REQ-494.
> - Continue automatically into normal REQ selection after recovery leaves shared lifecycle state and the index safe.
> - Do not respond to a dirty shared checkpoint by skipping to another ready REQ; the checkpoint is a common claim target, so such retries cannot succeed.
> - Leave existing working-REQ recovery behavior unchanged; this change targets archive/release/commit/metadata tails.
> 
> ### Legacy-tail association
> 
> - Extend association from generic path matching to lifecycle-aware recovery groups.
> - Project files may be assigned automatically only when an archived terminal REQ lacks provenance, explicitly lists the path in its implementation summary, and no competing candidate owns that path.
> - Shared metadata must use semantic evidence:
>   - Archive and working paths contain the same REQ identity.
>   - Checkpoint changes remove that REQ’s writer entry.
>   - Calibration and user-request moves reference the same REQ.
>   - Follow-up REQs explicitly reference their originating request.
>   - Release changes contain the matching REQ/version changelog entry.
> - Never authorize `CHECKPOINT.md`, release files, or other shared metadata through the existing generic “latest completed owner” rule.
> - Run protected-inventory classification first. Never read, stage, or commit secret-classified paths; refuse recovery if such a path is already staged.
> - Stage exact paths only. A path is recoverable only when its entire current diff belongs to one recovery group; do not silently include foreign hunks.
> - Commit all unambiguous groups and record their provenance canonically.
> - Preserve unrelated unstaged changes. Stop before selection when the remainder contains ambiguous queue/lifecycle/release state, foreign staged entries, or a dirty checkpoint; report typed reasons and exact paths.
> 
> ## Interfaces and Documentation
> 
> - Return typed JSON from both commands containing the REQ, journal phase, whether recovery was resumed or discovered, created commit hashes, blocked paths, reason codes, and verification/next commands.
> - Keep existing `complete`, `release`, and state commands backward compatible.
> - Rewrite the work action’s completion steps to create one exact finalization manifest and invoke `finalize`; remove prose-directed manual release-tail staging.
> - Update crash-recovery documentation to distinguish active working-REQ recovery from finalization-tail recovery.
> - Make `do-work commit` delegate detected finalization tails to the same recovery engine before handling ordinary associated changes.
> 
> ## Test Plan
> 
> - Inject interruption after every journal phase and verify the next `do-work run` resumes exactly once.
> - Reproduce the REQ-494 state—archived REQ, dirty checkpoint/release/code changes, no journal—and verify automatic association commits it and continues to claim REQ-452 in the same run.
> - Verify duplicate release entries, version bumps, archive moves, calibration rows, and provenance commits cannot occur.
> - Verify ambiguous checkpoint hunks, competing REQ ownership, foreign staged paths, and corrupt journal preimages stop safely without changing those bytes.
> - Verify unambiguous groups are recovered even when unrelated unstaged project changes remain.
> - Verify secret-classified files are never staged and pre-staged secrets block recovery.
> - Cover successful, failed, already-green, serial, and worktree finalization paths.
> - Simulate commit-hook failure and confirm the journal remains resumable without bypassing hooks.
> - Run the focused CLI/action contract tests followed by `bash _dev/tests/maintainer-verify.sh`.
> 
> ## Assumptions
> 
> - Finalization remains single-releaser; concurrent release publication is outside this change.
> - Journals are Git-private local state, not tracked repository files.
> - Automatic recovery favors progress only when ownership is unambiguous. It recovers safe groups but never guesses ownership of shared lifecycle state.
> 
> make sure to also capture the intention via do-work capture-request
> 
> after capture, finish the current changes in a safe spot then commit and stop
> ```

---
*Captured: 2026-09-02T13:07:19Z*
