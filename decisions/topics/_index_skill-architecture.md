---
title: "Topic Index: Skill Architecture"
type: topic-index
status: reference
topic_cluster: skill-architecture
sources:
  - CHANGELOG.md
  - CHANGELOG.md (0.49.0 The Architect)
  - CHANGELOG.md (0.59.0 The Quality Blueprint)
  - CHANGELOG.md (0.61.1 The Lean Cut)
  - CHANGELOG.md (0.62.0 The Karpathy Nod)
  - CHANGELOG.md (0.62.1 The Senior Engineer Test)
  - CHANGELOG.md (0.64.1 The Companion Split)
  - CHANGELOG.md (0.67.0 The Open Ear)
  - CLAUDE.md
  - README.md
  - actions/build-knowledge-base.md
  - actions/capture.md
  - actions/interview.md
  - actions/interview-reference.md
  - actions/pipeline.md
  - actions/review-work.md
  - actions/work.md
  - crew-members/karpathy.md
  - interviews/work-operating-model.md
  - specs/README.md
related:
  - page: adr-001-modular-action-prompts-and-companion-references
    rel: complements
  - page: adr-002-load-reusable-spec-templates-during-work
    rel: complements
  - page: adr-003-always-load-karpathy-guardrails
    rel: complements
  - page: adr-011-interview-framework-with-prescriptive-templates
    rel: complements
created: 2026-04-15
updated: 2026-04-16
confidence: high
---

# Skill Architecture

How the skill is structured, standardized, and behaviorally guided.

## ADRs

- [[adr-001-modular-action-prompts-and-companion-references]] — [ADR-001](../records/adr-001-modular-action-prompts-and-companion-references.md): Keep action files standalone and split bulky reference material into companion files when prompt size or readability becomes a liability.
- [[adr-002-load-reusable-spec-templates-during-work]] — [ADR-002](../records/adr-002-load-reusable-spec-templates-during-work.md): Use `specs/` templates as reusable quality scaffolds that the work action loads when a REQ clearly matches a task type.
- [[adr-003-always-load-karpathy-guardrails]] — [ADR-003](../records/adr-003-always-load-karpathy-guardrails.md): Apply Karpathy-inspired behavioral guardrails in every implementation pass, then audit them during review without double-counting issues.
- [[adr-011-interview-framework-with-prescriptive-templates]] — [ADR-011](../records/adr-011-interview-framework-with-prescriptive-templates.md): Add a generalized `interview` action that runs prescriptive templates from `interviews/<name>.md`, enforces a canonical entry contract, and produces agent-ready operating artifacts. Depends on ADR-001 (modular action + companion), ADR-002 (reusable templates at runtime), and ADR-005 (stateful and resumable).

## Cross-Cluster Links

- [[adr-001-modular-action-prompts-and-companion-references]] complements [[adr-005-pipeline-is-stateful-and-resumable]] in [[_index_workflow-orchestration]].
- [[adr-001-modular-action-prompts-and-companion-references]] complements [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]] in [[_index_pipeline-deliverables]].
- [[adr-003-always-load-karpathy-guardrails]] complements [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] in [[_index_pipeline-deliverables]].
