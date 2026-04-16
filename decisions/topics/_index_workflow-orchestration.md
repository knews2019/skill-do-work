---
title: "Topic Index: Workflow Orchestration"
type: topic-index
status: reference
topic_cluster: workflow-orchestration
sources:
  - CHANGELOG.md
  - CHANGELOG.md (0.51.5 The Full Send)
  - CHANGELOG.md (0.56.0 The Clean Sweep)
  - CHANGELOG.md (0.56.1 The Safety Net)
  - CHANGELOG.md (0.56.2 The Tight Scope)
  - CHANGELOG.md (0.60.3 The Paved Path)
  - CHANGELOG.md (0.60.5 The Honest Mirror)
  - CLAUDE.md
  - README.md
  - actions/capture.md
  - actions/cleanup.md
  - actions/pipeline.md
  - actions/review-work.md
  - actions/work.md
  - docs/work-guide.md
  - hooks/pipeline-guard.sh
related:
  - page: adr-004-canonicalize-pending-reqs-under-do-work-queue
    rel: complements
  - page: adr-005-pipeline-is-stateful-and-resumable
    rel: complements
  - page: adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles
    rel: complements
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# Workflow Orchestration

How pending work is stored and how the pipeline coordinates queue processing.

## ADRs

- [[adr-004-canonicalize-pending-reqs-under-do-work-queue]] — [ADR-004](../records/adr-004-canonicalize-pending-reqs-under-do-work-queue.md): Treat `do-work/queue/` as the canonical home for pending REQ files and update every workflow around that assumption.
- [[adr-005-pipeline-is-stateful-and-resumable]] — [ADR-005](../records/adr-005-pipeline-is-stateful-and-resumable.md): Treat the pipeline as a stateful orchestrator that dispatches existing actions, records progress in `do-work/pipeline.json`, and resumes across sessions.
- [[adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles]] — [ADR-006](../records/adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles.md): After the formal pipeline completes, continue draining pending work in explicit run-review loops with iteration caps and REQ-targeted reviews.

## Cross-Cluster Links

- [[adr-005-pipeline-is-stateful-and-resumable]] complements [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] in [[_index_pipeline-deliverables]].
- [[adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles]] complements [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] in [[_index_pipeline-deliverables]].
