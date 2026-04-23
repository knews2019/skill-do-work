---
title: "ADR-007: Close the Pipeline with Present and a Technical Debrief"
type: architecture-decision-record
status: accepted
topic_cluster: pipeline-deliverables
decided: 2026-04-13
sources:
  - CHANGELOG.md (0.63.0 The Closing Act)
  - CHANGELOG.md (0.63.1 The Debrief)
  - actions/pipeline.md
  - actions/present-work.md
  - README.md
related:
  - page: adr-005-pipeline-is-stateful-and-resumable
    rel: depends-on
  - page: adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles
    rel: complements
  - page: adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats
    rel: complements
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# ADR-007: Close the Pipeline with Present and a Technical Debrief

Topic cluster: [[_index_pipeline-deliverables]] ([topic index](../topics/_index_pipeline-deliverables.md))
See also: [[adr-005-pipeline-is-stateful-and-resumable]] (depends-on), [[adr-006-pipeline-drains-follow-up-work-in-bounded-reviewed-cycles]] (complements), [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]] (complements)

## Context

The pipeline originally stopped after review, which left a manual handoff gap between shipping the work and communicating it. The changelog then added `present` as a sixth step and immediately followed that with a rule that completion must generate a technical debrief containing final summary, test state, coherence, carry-forward work, deliverables, and how-to-verify guidance.

Those expectations are now embedded in the live pipeline action. `pipeline.md` still frames completion as education, not a checkmark, and `present-work.md` still carries the sibling client-facing brief, explainer, and verification recipe that the pipeline is expected to surface.

## Decision

A pipeline run is not complete until it has:
- executed `present` as the sixth step when artifacts exist, and
- assembled a persisted Pipeline Completion Report that teaches the user what shipped and how to verify it.

The completion artifact is part of the product, not optional ceremony. Long pipelines must leave behind readable evidence, not just terminal output.

## Alternatives

1. Stop after review and tell the user to run `do-work present` separately.
This was rejected because it leaves a predictable final handoff undone.

2. Treat pipeline completion as a one-line status message.
This was rejected because multi-REQ work needs a durable summary, not ephemeral terminal text.

3. Generate only client-facing materials and skip the technical debrief.
This was rejected because developers and reviewers still need the audit trail.

## Consequences

The project now guarantees that a successful pipeline leaves behind both communication artifacts and verification artifacts. That improves handoff quality for technical and non-technical readers alike.

The cost is more completion work per pipeline run. The orchestration step must collect evidence carefully and keep report depth proportional to the scope of the pipeline.

## References

- [CHANGELOG.md](../../CHANGELOG.md) — `0.63.0 The Closing Act`, `0.63.1 The Debrief`
- [actions/pipeline.md](../../actions/pipeline.md)
- [actions/present-work.md](../../actions/present-work.md)
- [README.md](../../README.md) — pipeline section
