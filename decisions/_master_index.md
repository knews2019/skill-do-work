---
title: "Decision Index"
type: master-index
status: reference
sources:
  - CHANGELOG.md
  - decisions/topics/
  - decisions/records/
related:
  - page: _index_skill-architecture
    rel: complements
  - page: _index_workflow-orchestration
    rel: complements
  - page: _index_pipeline-deliverables
    rel: complements
  - page: _index_knowledge-base
    rel: complements
created: 2026-04-15
updated: 2026-04-16
confidence: high
---

# Decision Index

Read this first. This ADR log captures 11 load-bearing, still-in-force decisions mined from `CHANGELOG.md` and corroborated against the current repo state.

## Topic Clusters

- [Skill Architecture](./topics/_index_skill-architecture.md) — 4 ADRs — How the skill is structured, standardized, and behaviorally guided. Related pages: [[adr-001-modular-action-prompts-and-companion-references]], [[adr-002-load-reusable-spec-templates-during-work]], [[adr-003-always-load-karpathy-guardrails]], [[adr-011-interview-framework-with-prescriptive-templates]].
- [Workflow Orchestration](./topics/_index_workflow-orchestration.md) — 3 ADRs — How pending work is stored and how the pipeline coordinates queue processing. Related pages: [[adr-004-canonicalize-pending-reqs-under-do-work-queue]], [[adr-005-pipeline-is-stateful-and-resumable]], [[adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles]].
- [Pipeline Deliverables](./topics/_index_pipeline-deliverables.md) — 2 ADRs — How completed pipeline work is presented, summarized, and linked for different audiences. Related pages: [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]], [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]].
- [Knowledge Base](./topics/_index_knowledge-base.md) — 2 ADRs — How the BKB is structured, linked, and operated as a persistent wiki system. Related pages: [[adr-009-build-knowledge-base-as-a-compiled-interlinked-wiki]], [[adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb]].

## Navigation Notes

- [Timeline log](./log.md) — append-only timeline of the historical decisions plus this ADR bootstrap pass.
- [Progress tracker](./_progress.md) — resumable notes, scope decisions, and next ADR number.
- ADR pages live under [`decisions/records/`](./records/).
