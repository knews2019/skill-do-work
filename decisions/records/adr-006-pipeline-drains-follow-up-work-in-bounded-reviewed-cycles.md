---
title: "ADR-006: Pipeline Drains Follow-Up Work in Bounded Reviewed Cycles"
type: architecture-decision-record
status: accepted
topic_cluster: workflow-orchestration
decided: 2026-04-10
sources:
  - CHANGELOG.md (0.56.0 The Clean Sweep)
  - CHANGELOG.md (0.56.1 The Safety Net)
  - CHANGELOG.md (0.56.2 The Tight Scope)
  - actions/pipeline.md
  - actions/work.md
  - actions/review-work.md
related:
  - page: adr-004-canonicalize-pending-reqs-under-do-work-queue
    rel: depends-on
  - page: adr-005-pipeline-is-stateful-and-resumable
    rel: depends-on
  - page: adr-007-close-the-pipeline-with-present-and-a-technical-debrief
    rel: complements
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# ADR-006: Pipeline Drains Follow-Up Work in Bounded Reviewed Cycles

Topic cluster: [[_index_workflow-orchestration]] ([topic index](../topics/_index_workflow-orchestration.md))
See also: [[adr-004-canonicalize-pending-reqs-under-do-work-queue]] (depends-on), [[adr-005-pipeline-is-stateful-and-resumable]] (depends-on), [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] (complements)

## Context

Once the pipeline existed, the next load-bearing question was how it should behave when follow-up REQs remained in the queue. The changelog shows a sequence of releases that first added queue continuation, then tightened its safeguards: explicit error handling, a three-cycle cap, and a rule that continuation reviews must target the REQ IDs from that cycle rather than broad UR scopes.

The current `pipeline.md` still encodes the same bounded continuation design. It scans `do-work/queue/` after completion, records the REQs about to be processed, runs standard queue work, reviews each REQ individually, stops after three cycles, and does not reopen the formal pipeline state machine during that continuation.

## Decision

The pipeline may drain remaining pending REQs after its primary request is complete, but only through a bounded post-pipeline continuation loop:
- capture the pending REQ IDs for the current continuation cycle,
- run the work action on that batch,
- review each REQ individually by REQ ID,
- stop after three cycles or on error.

The continuation is intentionally distinct from the formal pipeline state machine. `pipeline.json` stays complete, and failures in continuation yield explicit recovery commands instead of silent retries.

## Alternatives

1. End the pipeline immediately after the primary request finishes.
This was rejected because review-generated or previously pending work would be left hanging even when the user clearly asked for a full-cycle automation path.

2. Let the continuation drain forever until the queue is empty.
This was rejected because review steps can create follow-ups indefinitely, which risks runaway loops.

3. Review continuation batches by UR.
This was rejected because UR-scoped review can re-review unrelated completed REQs instead of the specific batch just processed.

## Consequences

The benefit is a more complete queue-draining experience that still preserves review rigor and operator control. Follow-up work does not vanish after the first pipeline success.

The trade-off is extra orchestration complexity. The continuation loop needs careful targeting, user-facing recovery guidance, and clear separation from the main `pipeline.json` lifecycle.

## References

- [CHANGELOG.md](../../CHANGELOG.md) — `0.56.0 The Clean Sweep`, `0.56.1 The Safety Net`, `0.56.2 The Tight Scope`
- [actions/pipeline.md](../../actions/pipeline.md)
- [actions/work.md](../../actions/work.md)
- [actions/review-work.md](../../actions/review-work.md)
