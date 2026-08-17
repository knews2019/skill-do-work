---
session_ended: 2026-08-17T20:23:49Z
last_completed: REQ-219
queue_state: 2 pending, 0 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 9
session_depth: heavy
---

# Session Checkpoint

## Completed This Session

- REQ-203: Harden presentation target-ID source-seam tests (Route B, review 97%) — commit `2c9e07e`
- REQ-204: Harden ai-report generated-batch lifecycle (Route B, review 93%) — commit `7f7a5ca`; sweep created REQ-220 and REQ-221
- REQ-205: Make portfolio publication independent and exact (Route B, review 95%) — commit `4a0b90e`
- REQ-206: Finish active publication delegation (Route B, review 96%) — commit `af2a4d7`
- (correction) shell-primitives guide repointed after REQ-205's sweep missed it — commit `9822385`
- REQ-216: Teach atomic-download retry and optional credentials (Route A, review 96%) — commit `27e8ff9`
- REQ-217: Add the upstream archive fetcher with a git fallback (Route B, review 94%) — commit `0e8cf0d`
- REQ-218: Ratchet the tools download and correct the gitattributes claim (Route A, review 95%) — commit `052d384`
- UR-049 consolidated and closed — commit `c0331e8`
- REQ-219: Durations view on the Kanban board (Route C) — commit `cfffd90`; UR-050 consolidated

## Merge With the Other Checkout

Resolved and committed at `16927c9` ("Merge sync-incoming-20260817-223749 (VM share) into main"). The merge brought the other checkout's `0.202.0` release (Mechanical REQ Reservation Cleanup); REQ-219 was renumbered on top of it and shipped as **`0.203.0`**.

**The merge resurrected three queue files, and they are now deleted.** The sync branch was based on a commit predating the REQ-216/217/218 build, so the merge restored all three as they looked *before* being worked. A first attempt at cleanup renumbered them into three fresh ids instead of removing them — which silenced the three `duplicate-req-id` findings that were the only guard catching them, leaving three clean-looking `pending` REQs describing already-shipped work. Each was proven byte-identical (ids normalized) to its originally-captured pre-build file from commit `315c55c`, so the renumbering was reverted and the files removed. The archived copies in `do-work/archive/UR-049/` are strict supersets — same specification text plus Triage, Plan, Implementation Summary, Review, Lessons and verified hashes (`27e8ff9`, `0e8cf0d`, `052d384`). The pre-build snapshots remain recoverable from `16927c9`. `queue-kanban verify` is now clean.

The three reservation markers the renumbering allocated are now orphaned and were deliberately left in place — `scripts/cleanup-req-reservations.sh` reaps markers mechanically on SessionStart; an agent never deletes them.

## Still Queued

Both consented at `2026-08-17T19:58:23Z` (`pending-answers` → `pending`); no REQ is waiting on an answer.

- REQ-220: Extend runtime-boundary ownership to the remaining publication helpers (pending) — the same defect class REQ-204 fixed is alive in two shipped scripts: `generate-report-image.sh` interrupted directly leaves its backend running, and `install-last30days.sh` can nest a staging tree inside a reappearing target and still report success
- REQ-221: Extract the ai-report image batch into a shipped script (pending)

## Session Notes

- Every completed REQ was verified with `bash _dev/tests/maintainer-verify.sh` at exit 0 with zero FAIL lines on this machine. The 41 container-environment FAILs the earlier checkpoint recorded do not reproduce here.
- Three REQs landed the same lesson from different directions: an assertion over a whole document tests vocabulary, not behavior. Word-bound the verb, anchor it to the step, replay a mutation matrix against the real file.
- REQ-205 changed what "published from the same bytes" guarantees and its sweep missed a shipped guide. Sweep with a repo-wide grep, not a hand-typed path list — a mistyped path fails silently.
- REQ-216 and REQ-217 are a pair: retry alone does not close a sustained 429, so the git fallback is the load-bearing half. `git archive` is mandatory there — it is the only fetch mechanism that honors `export-ignore`.
- Two checkouts worked this queue concurrently today. The collision surfaced exactly where the model says it should — as ordinary merge conflicts plus duplicate ids — and cost one interrupted REQ commit, not any work.
- `work.md:213` offers "remove the duplicate **or rename if this is a re-do**", and the rename branch is the dangerous one: renaming a resurrected duplicate clears the `duplicate-req-id` finding without changing what the file asks for, converting a caught collision into three fresh `pending` REQs that would rebuild shipped work unguarded. Decide re-do vs stale by diffing the queue copy against the archived one with ids normalized — identical bodies mean stale, and the only surviving signal is `ur-archived-with-live-member`.
