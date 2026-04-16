---
title: "ADR-005: Pipeline Is Stateful and Resumable"
type: architecture-decision-record
status: accepted
topic_cluster: workflow-orchestration
decided: 2026-04-08
sources:
  - CHANGELOG.md (0.51.5 The Full Send)
  - README.md
  - actions/pipeline.md
  - hooks/pipeline-guard.sh
related:
  - page: adr-004-canonicalize-pending-reqs-under-do-work-queue
    rel: depends-on
  - page: adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles
    rel: complements
  - page: adr-007-close-the-pipeline-with-present-and-a-technical-debrief
    rel: complements
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# ADR-005: Pipeline Is Stateful and Resumable

Topic cluster: [[_index_workflow-orchestration]] ([topic index](../topics/_index_workflow-orchestration.md))
See also: [[adr-004-canonicalize-pending-reqs-under-do-work-queue]] (depends-on), [[adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles]] (complements), [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] (complements)

## Context

The project needed an end-to-end path that could investigate, capture, verify, run, and review a request without re-implementing those actions inline. The first pipeline release made orchestration itself a first-class action, with explicit state tracking, resume behavior, and a stop hook to guard mid-run exits. Later refinements preserved the same core shape rather than replacing it.

That architecture is still intact. The current `pipeline.md` describes the pipeline as a six-step orchestrator, uses `do-work/pipeline.json` as source of truth, requires foreground dispatch, and reiterates that the pipeline never re-implements existing action logic.

## Decision

The pipeline is a resumable orchestration layer over existing actions. It owns macro-step state in `do-work/pipeline.json`, writes state transitions before and after dispatch, resumes from the first unfinished step, and keeps the state file out of version control.

The orchestrator does not duplicate capture, work, review, or present logic. Instead, it passes the right context and artifacts to those existing actions and records which step has completed.

## Alternatives

1. Ask users to run the actions manually in sequence.
This was rejected because it loses end-to-end continuity and puts too much sequencing burden on the operator.

2. Re-implement each action's logic inside the pipeline.
This was rejected because duplicated behavior would drift from the standalone actions and make fixes harder to apply consistently.

3. Keep pipeline progress only in chat context.
This was rejected because the workflow must survive crashes, context limits, and multi-session resumes.

## Consequences

The project gains a durable full-cycle mode that can be resumed safely and inspected through a single state file. It also keeps action ownership clear because the pipeline orchestrates rather than redefines them.

The trade-off is more state-management discipline. `pipeline.json` must reflect reality, stay out of commits, and be updated carefully around every step transition.

## References

- [CHANGELOG.md](../../CHANGELOG.md) — `0.51.5 The Full Send`
- [README.md](../../README.md) — pipeline overview
- [actions/pipeline.md](../../actions/pipeline.md)
- [hooks/pipeline-guard.sh](../../hooks/pipeline-guard.sh)
