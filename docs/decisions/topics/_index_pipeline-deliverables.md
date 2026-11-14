---
title: "Topic Index: Pipeline Deliverables"
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
  - CHANGELOG.md (0.317.0 Tinka Canvas Offer Removed)
  - do-work/archive/UR-172/REQ-721-ai-report-closing-suggests-tinka-canvas.md
related:
  - page: adr-007-close-the-pipeline-with-present-and-a-technical-debrief
    rel: complements
  - page: adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats
    rel: complements
  - page: adr-023-adopt-no-external-review-canvas
    rel: complements
created: 2026-04-15
updated: 2026-11-14
confidence: high
---

# Pipeline Deliverables

How finished work is presented to the reviewer. This covers the live ai-report deliverable and the history of the retired stateful pipeline's completion reports. The live successor of those reports is `do-work-toolbox present-work`, invoked explicitly by the full-cycle prompt.

## ADRs

- [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] — [ADR-007](../records/adr-007-close-the-pipeline-with-present-and-a-technical-debrief.md) (**superseded**): Historical completion/debrief contract retired by ADR-019 and REQ-145.
- [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]] — [ADR-008](../records/adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats.md) (**superseded**): Historical three-rendering contract retired with the stateful action.
- [[adr-023-adopt-no-external-review-canvas]] — [ADR-023](../records/adr-023-adopt-no-external-review-canvas.md): ai-report bundles are reviewed as browser pages with feedback in chat, not on an external canvas such as Tinka.

## Cross-Cluster Links

- [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] depends-on [[adr-005-pipeline-is-stateful-and-resumable]] in [[_index_workflow-orchestration]].
- [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] complements [[adr-006-pipeline-processes-follow-up-work-in-bounded-reviewed-cycles]] in [[_index_workflow-orchestration]].
- [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]] depends-on [[adr-001-modular-action-prompts-and-companion-references]] in [[_index_skill-architecture]].
