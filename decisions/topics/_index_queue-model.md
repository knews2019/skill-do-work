# Topic: Queue Model

How captured work is represented on disk, how it moves through the pipeline, and what's off-limits once a builder claims it.

## Shape of a Captured Request

Every invocation of the capture action produces **two coordinated artifacts**:

1. A **User Request (UR)** folder — the umbrella grouping related work.
2. One or more **Request (REQ)** files inside that folder — individual units of work the builder processes.

Pending REQs live at `do-work/queue/`. Once a builder claims a REQ, the UR folder moves to `do-work/working/`. Once complete, it moves to `do-work/archive/`. Files inside `working/` and `archive/` are immutable — any follow-up is a new REQ with an `addendum_to:` pointer, not an edit.

## ADRs in This Cluster

| ADR  | Title                                | Status   |
|------|--------------------------------------|----------|
| [[../adr/0002-ur-req-pairing\|0002]]              | UR + REQ pairing is mandatory on every capture | accepted |
| [[../adr/0003-immutable-inflight-archived\|0003]] | Immutable `working/` and `archive/` folders   | accepted |
| [[../adr/0008-queue-canonical-path\|0008]]        | Queue canonical path — `do-work/queue/`        | accepted |

## How They Relate

- [[../adr/0002-ur-req-pairing|0002]] sets the shape (UR + REQ).
- [[../adr/0003-immutable-inflight-archived|0003]] protects the shape after claim — `addendum_to` is the only legal way to extend an in-flight UR.
- [[../adr/0008-queue-canonical-path|0008]] nails down *where* the pending shape lives, resolving a recurring class of stale-path bugs.

## Downstream Behaviors That Depend on This Cluster

- The work action's folder moves (`queue/ → working/ → archive/`) assume both the UR/REQ shape and the immutability invariant.
- The pipeline action references URs by id when dispatching sub-steps — this only works because every capture produces a UR.
- The verify-requests and review-work actions check coherence **across addendum chains**, which would be impossible if in-flight edits were allowed.
- The cleanup action's sweep logic relies on the canonical queue path to find loose REQs.
