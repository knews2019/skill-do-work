---
title: "ADR-008: Render Pipeline Debriefs in Three Cross-Linked Audience-Specific Formats"
type: architecture-decision-record
status: accepted
topic_cluster: pipeline-deliverables
decided: 2026-04-13
sources:
  - CHANGELOG.md (0.63.2 The Triple Render)
  - CHANGELOG.md (0.64.0 The Cross-Linked Set)
  - CHANGELOG.md (0.64.1 The Companion Split)
  - actions/pipeline.md
  - actions/pipeline-reference.md
  - actions/present-work.md
related:
  - page: adr-001-modular-action-prompts-and-companion-references
    rel: depends-on
  - page: adr-007-close-the-pipeline-with-present-and-a-technical-debrief
    rel: depends-on
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# ADR-008: Render Pipeline Debriefs in Three Cross-Linked Audience-Specific Formats

Topic cluster: [[_index_pipeline-deliverables]] ([topic index](../topics/_index_pipeline-deliverables.md))
See also: [[adr-001-modular-action-prompts-and-companion-references]] (depends-on), [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] (depends-on)

## Context

Once the pipeline produced a debrief at all, the next design question was who that debrief was for. The changelog captures a rapid sequence of decisions: generate three renderings from one dataset, open every summary with a plain-language "What got built" narrative, cross-link sibling artifacts for both stakeholder and developer audiences, and eventually split the bulky rendering templates into `pipeline-reference.md` once the prompt grew too large.

The present files preserve that entire shape. `pipeline.md` requires all three formats, `pipeline-reference.md` holds the templates and composition rules, and `present-work.md` cross-links the interactive explainer, client brief, and pipeline summaries as related reading.

## Decision

Pipeline completion data is rendered three ways from a single source dataset:
- plain markdown for developers,
- Marp slides for walkthroughs and stakeholder reviews,
- standalone HTML for non-technical readers.

All three formats must carry the same facts, start with a plain-language "What got built" entry point, and link readers to sibling artifacts that deepen either understanding or auditability. The rendering templates and cross-format rules live in `pipeline-reference.md` rather than inside `pipeline.md`.

## Alternatives

1. Generate only markdown.
This was rejected because different audiences consume the same work in different surfaces.

2. Allow each format to editorialize independently.
This was rejected because drift across formats undermines trust in the debrief.

3. Keep all rendering templates inline in `pipeline.md`.
This was rejected once the prompt crossed practical read limits.

## Consequences

The project now produces audience-appropriate deliverables without splitting the underlying facts. Readers can enter from any artifact and navigate to deeper context through deliberate cross-links.

The trade-off is a larger reporting surface area. Template parity, sibling-link hygiene, and companion-file discoverability all require ongoing maintenance attention.

## References

- [CHANGELOG.md](../../CHANGELOG.md) — `0.63.2 The Triple Render`, `0.64.0 The Cross-Linked Set`, `0.64.1 The Companion Split`
- [actions/pipeline.md](../../actions/pipeline.md)
- [actions/pipeline-reference.md](../../actions/pipeline-reference.md)
- [actions/present-work.md](../../actions/present-work.md)
