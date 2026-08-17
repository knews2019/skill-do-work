---
session_ended: 2026-08-17T21:24:00Z
last_completed: REQ-221
queue_state: 1 pending, 0 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 2
session_depth: moderate
---

# Session Checkpoint

## Completed This Session

- REQ-220: Extend runtime-boundary ownership to the remaining publication helpers (Route B, review 95%) — commit `1fb3635`, hash recorded `a980e4f`, shipped as **0.203.1**
- REQ-221: Extract the ai-report image batch into a shipped script (Route B, review 96%) — commit `e2b45bf`, hash recorded `80238ef`, shipped as **0.204.0**

Both hashes were confirmed with `record-commit-hash.sh --verify`. `maintainer-verify.sh` and `queue-kanban verify` both exit 0 at hand-back.

## Still Queued

- REQ-225: State verified-exact-publication once as a condition in the shipped shell guide (pending) — created by REQ-220's review as a `[normal]` discovered task, then flipped to `pending` by commit `e3e5b69` ("[REQ-225] Approve the shared-condition rewrite via clarify"), which landed on `main` from **outside this session**. **Not built here — awaiting the maintainer's own confirmation.** The whole reason it was captured as `pending-answers` is that it is a cascade-depth-two taste decision about reorganizing a canonical shipped document, so a consent record this session cannot attribute to the maintainer is exactly the thing it was created to avoid relying on. Scope if approved: documentation only, no script behavior changes.

## Session Notes

- A builder's report is a claim. Both REQs' red/green evidence was re-derived from git state instead: REQ-220 by reverting both scripts and re-running the suite, REQ-221 by moving the new script aside. Both produced exactly the expected FAIL sets with no pre-existing regression — and REQ-220's RED surfaced damage the captured proof had not predicted (the old code set `publication_complete=1`, so the EXIT trap deleted the backup and destroyed the prior tree). That became an assertion.
- **An extraction cannot be validated by a green suite.** REQ-221 moved a block that twelve `contract-regressions.sh` assertions matched *by address*. Deleting those assertions and deleting the behavior they guard look identical from the exit code. Each one had to be classified individually as "moved with the code" or "belongs at the old address"; one belonged at the old address and had been dropped, and was restored and mutation-tested (REQ-221 D-04).
- The runtime-boundary defect class has now been closed in four scripts, each time locally, because `prescribed-shell-primitives.md` states the nesting rule inside one script's section rather than as a condition. It has been found four times by review sweep and zero times by someone reading the guide. That observation is what REQ-225 is.
- REQ-220's fix in `install-last30days.sh` has a trap worth remembering: setting `publication_started=0` is *not* enough to stop the rollback, because `cleanup_install_paths` also fires on a surviving backup and `restore_previous_tree` opens with an unconditional `rm -rf "$target_directory"`. The obvious fix would have satisfied "exit nonzero" while destroying the colliding tree the requirement protects.
- Another checkout committed to `main` mid-session (`e3e5b69`). The concurrent-queue model held: it surfaced as ordinary new commits, not as a conflict.
- Four reservation markers (`REQ-000222`–`REQ-000225`) remain untracked and were deliberately left alone; `cleanup-req-reservations.sh` reaps them on SessionStart.
