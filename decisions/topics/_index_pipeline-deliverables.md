---
title: "Topic Index: Retired Pipeline Deliverables"
type: topic-index
status: reference
topic_cluster: pipeline-deliverables
sources:
  - CHANGELOG.md
  - CHANGELOG.md (0.63.0 The Closing Act)
  - CHANGELOG.md (0.63.1 The Debrief)
  - CHANGELOG.md (0.63.2 The Triple Render)
  - CHANGELOG.md (0.64.0 The Cross-Linked Set)
  - CHANGELOG.md (0.64.1 The Companion Split)
  - README.md
  - skills/do-work-toolbox/actions/present-work.md
related:
  - page: adr-007-close-the-pipeline-with-present-and-a-technical-debrief
    rel: complements
  - page: adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats
    rel: complements
created: 2026-04-15
updated: 2026-08-08
confidence: high
---

# Retired Pipeline Deliverables

Historical record of how the retired stateful pipeline presented, summarized, and linked completed work. The live successor is `do-work-toolbox present-work`, invoked explicitly by the full-cycle prompt.

## ADRs

- [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] — [ADR-007](../records/adr-007-close-the-pipeline-with-present-and-a-technical-debrief.md) (**superseded**): Historical completion/debrief contract retired by ADR-019 and REQ-145.
- [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]] — [ADR-008](../records/adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats.md) (**superseded**): Historical three-rendering contract retired with the stateful action.

## Cross-Cluster Links

- [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] depends-on [[adr-005-pipeline-is-stateful-and-resumable]] in [[_index_workflow-orchestration]].
- [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] complements [[adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles]] in [[_index_workflow-orchestration]].
- [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]] depends-on [[adr-001-modular-action-prompts-and-companion-references]] in [[_index_skill-architecture]].
