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
updated: 2026-08-08
confidence: high
---

# Decision Index

Read this first. This ADR log captures every load-bearing, still-in-force decision mined from `CHANGELOG.md` and corroborated against the current repo state, plus one declined-decision record (ADR-014) documenting a path considered and rejected. **Counts below are derived from `records/` — read the directory, not a number here**, which is how the previous hand-maintained total came to be three behind.

## Topic Clusters

- [Skill Architecture](./topics/_index_skill-architecture.md) — 7 ADRs — How the skill is structured, standardized, behaviorally guided, and distributed. Related pages: [[adr-001-modular-action-prompts-and-companion-references]], [[adr-002-load-reusable-spec-templates-during-work]], [[adr-003-always-load-karpathy-guardrails]], [[adr-011-interview-framework-with-prescriptive-templates]], [[adr-012-interview-v2-gap-closure]], [[adr-013-harden-the-vendored-skill-distribution-model]], [[adr-016-vendor-queue-kanban-into-the-skill]].
- [Workflow Orchestration](./topics/_index_workflow-orchestration.md) — How pending work is stored, how `do-work run` coordinates queue processing, and who owns a queue. Includes superseded ADR-005 and ADR-006 as history. Related pages: [[adr-004-canonicalize-pending-reqs-under-do-work-queue]], [[adr-005-pipeline-is-stateful-and-resumable]], [[adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles]], [[adr-014-considered-declined-autonomous-loop-until-done]], [[adr-015-load-maintenance-crew-via-req-marker]], [[adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser]].
- [Retired Pipeline Deliverables](./topics/_index_pipeline-deliverables.md) — 2 superseded ADRs preserving the former completion-report history. Related pages: [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]], [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]].
- [Knowledge Base](./topics/_index_knowledge-base.md) — 2 ADRs — How the BKB is structured, linked, and operated as a persistent wiki system. Related pages: [[adr-009-build-knowledge-base-as-a-compiled-interlinked-wiki]], [[adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb]].

## Navigation Notes

- [Timeline log](./log.md) — append-only timeline of the historical decisions plus this ADR bootstrap pass.
- [Progress tracker](./_progress.md) — resumable notes, scope decisions, and next ADR number.
- ADR pages live under [`decisions/records/`](./records/).
