---
title: "Topic Index: Knowledge Base"
type: topic-index
status: reference
topic_cluster: knowledge-base
sources:
  - CHANGELOG.md
  - CHANGELOG.md (0.43.0 The Knowledge Forge)
  - CHANGELOG.md (0.43.1 The Clean Handoff)
  - CHANGELOG.md (0.43.2 The Gap Closer)
  - CHANGELOG.md (0.43.4 The Guard Dog)
  - CHANGELOG.md (0.43.5 The Final Polish)
  - CHANGELOG.md (0.45.0 The Second Brain)
  - CHANGELOG.md (0.46.0 The Agent Crew)
  - CHANGELOG.md (0.47.0 The Full Crew)
  - SKILL.md
  - actions/bkb-reference.md
  - actions/build-knowledge-base.md
  - docs/bkb-guide.md
related:
  - page: adr-009-build-knowledge-base-as-a-compiled-interlinked-wiki
    rel: complements
  - page: adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb
    rel: complements
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# Knowledge Base

How the BKB is structured, linked, and operated as a persistent wiki system.

## ADRs

- [[adr-009-build-knowledge-base-as-a-compiled-interlinked-wiki]] — [ADR-009](../records/adr-009-build-knowledge-base-as-a-compiled-interlinked-wiki.md): Model the BKB as a persistent compiled wiki with a master index, topic indexes, logs, and a disciplined source-processing lifecycle.
- [[adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb]] — [ADR-010](../records/adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb.md): Layer typed relationships, a self-improving retrieval agent, and a built-in/custom crew model on top of the compiled BKB wiki.

## Cross-Cluster Links

- [[adr-009-build-knowledge-base-as-a-compiled-interlinked-wiki]] complements [[adr-001-modular-action-prompts-and-companion-references]] in [[_index_skill-architecture]].
- [[adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb]] complements [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]] in [[_index_pipeline-deliverables]].
