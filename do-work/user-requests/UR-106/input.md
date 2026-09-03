---
id: UR-106
title: 'Lifecycle overhead around REQ-531: retry the gate once before a repair REQ, finalize on the REQ''s own paths'
created_at: 2026-09-03T20:05:46Z
requests: [REQ-559, REQ-560]
word_count: 196
---

# Lifecycle Overhead Around REQ-531: Retry the Gate Once Before a Repair REQ, Finalize on the REQ's Own Paths

## Summary

REQ-531 (review findings below critical stay in the report) took 87 minutes from first claim to finalize for a change its builder wrote in about 20. The Codex account of the delay names four lifecycle costs; two of them have no fix in the queue and are captured here. A transient gate failure (a broken pipe) cost 28 minutes because the pipeline deferred the REQ and minted repair REQ-548, which then found the gate already green. Other sessions' dirt in the working tree cost two "preserve" commits authored by the pipeline (061b7dbf, 83594c5e) because hand-back and finalize want a clean tree. The remaining two costs are already queued: the contract file's merge conflicts (REQ-539) and the phase count (REQ-502 to REQ-510).

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-559 | Retry a red repository gate once, immediately, before deferring the REQ or minting a repair REQ |
| REQ-560 | Hand-back and finalize check cleanliness only on the REQ's own paths; foreign dirt is left alone, never committed or surfaced as a stop |

## Batch Constraints

- Both are changes to the skill's own operating instructions plus the script or command that launches the gate or validates the tree; mechanics go in programs, judgment stays in prose (CLAUDE.md).
- Neither waits on the other or on the REQ-502 chain; both are independent of UR-104.
- Evidence for the RED state is the REQ-531 commit timeline of 2026-09-03: claim 21:19, defer 21:37, REQ-548 no-op repair 21:37 to 21:47, re-claim 21:47, preserve commits 21:47 and 22:24, finalize 22:45 (local, UTC+3).

## Full Verbatim Input

> ```
> both landed + maintainance reqs
> 
> Also see here what codex said when I asked it what took so long:
> 
> """
> The code change itself was relatively small. Most of the time went into the `do-work` lifecycle around it:
> 
> - The maintainer gate initially failed nondeterministically with a broken pipe, so the workflow required deferring REQ-531, creating and completing repair REQ-548, then resuming REQ-531.
> - Concurrent work changed the same regression-test file, causing a merge conflict that needed semantic resolution and full re-verification.
> - The skill required isolated builder work, independent review, regression checks, release/version updates, archival, worktree cleanup, and session checkpointing.
> - Finalization was briefly blocked by an unrelated existing calibration change, which had to be preserved separately.
> 
> So the delay was mainly gate recovery, concurrency handling, and mandatory lifecycle bookkeeping—not writing the 13-file implementation. In hindsight, I could have communicated those delays more clearly while they were happening.
> """
> 
> <- that is the lifecycle around it that we need to improve a lot
> 
> [assistant recommended two captures: a gate retry-once rule before a repair REQ is minted, and path-scoped cleanliness at finalize; the maintainer answered:]
> 
> do 1, 2 and 3, I'll release the cloud claims
> ```

---
*Captured: 2026-09-03T20:05:46Z*
