---
title: "ADR-008: Render Pipeline Debriefs in Three Cross-Linked Audience-Specific Formats"
type: architecture-decision-record
status: superseded
topic_cluster: pipeline-deliverables
decided: 2026-04-13
sources:
  - CHANGELOG.md (0.63.2 The Triple Render)
  - CHANGELOG.md (0.64.0 The Cross-Linked Set)
  - CHANGELOG.md (0.64.1 The Companion Split)
  - actions/pipeline.md
  - actions/present-work.md
related:
  - page: adr-001-modular-action-prompts-and-companion-references
    rel: depends-on
  - page: adr-007-close-the-pipeline-with-present-and-a-technical-debrief
    rel: depends-on
  - page: adr-019-four-skill-suite-contract
    rel: superseded-by
created: 2026-04-15
updated: 2026-08-08
confidence: high
---

# ADR-008: Render Pipeline Debriefs in Three Cross-Linked Audience-Specific Formats

> **Superseded 2026-08-08.** ADR-019 and REQ-145 retired the pipeline-owned three-format debrief surface. This page remains the historical record of that output contract.

Topic cluster: [[_index_pipeline-deliverables]] ([topic index](../topics/_index_pipeline-deliverables.md))
See also: [[adr-001-modular-action-prompts-and-companion-references]] (depends-on), [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] (depends-on)

## Context

Once the pipeline produced a debrief at all, the next design question was who that debrief was for. The changelog captures a rapid sequence of decisions: generate three renderings from one dataset, open every summary with a plain-language "What got built" narrative, cross-link sibling artifacts for both stakeholder and developer audiences, and split the bulky rendering templates into a `pipeline-reference.md` companion once the prompt grew too large (later re-inlined back into `pipeline.md` when trimming brought the combined size back under the token budget).

That shape remained live until the stateful action and companion reference were removed at modular cutover. `do-work-toolbox present-work` continues to own its own client brief, explainer, and optional walkthrough without recreating pipeline summary state or templates.

## Decision

Pipeline completion data is rendered three ways from a single source dataset:
- plain markdown for developers,
- Marp slides for walkthroughs and stakeholder reviews,
- standalone HTML for non-technical readers.

All three formats must carry the same facts, start with a plain-language "What got built" entry point, and link readers to sibling artifacts that deepen either understanding or auditability. The rendering templates and cross-format rules live in `actions/pipeline-reference.md`, pointed at by name from `pipeline.md`'s Output Format section. (They were inline in `pipeline.md` from 0.76.0 until REQ-030 re-split them out on 2026-07-27 — see adr-001 for the split/re-inline history of this pair. The decision recorded here is about the three formats and their cross-links, not about which file the templates sit in.)

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
- [skills/do-work-toolbox/actions/present-work.md](../../skills/do-work-toolbox/actions/present-work.md)
- [[adr-019-four-skill-suite-contract]] — superseding modular-suite contract
