---
title: "ADR-001: Modular Action Prompts and Companion References"
type: architecture-decision-record
status: accepted
topic_cluster: skill-architecture
decided: 2026-04-07
sources:
  - CHANGELOG.md (0.49.0 The Architect)
  - CHANGELOG.md (0.61.1 The Lean Cut)
  - CHANGELOG.md (0.64.1 The Companion Split)
  - CLAUDE.md
  - actions/work.md
  - actions/build-knowledge-base.md
  - actions/pipeline.md
related:
  - page: adr-005-pipeline-is-stateful-and-resumable
    rel: complements
  - page: adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats
    rel: complements
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# ADR-001: Modular Action Prompts and Companion References

Topic cluster: [[_index_skill-architecture]] ([topic index](../topics/_index_skill-architecture.md))
See also: [[adr-005-pipeline-is-stateful-and-resumable]] (complements), [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]] (complements)

## Context

The skill started as a smaller prompt bundle, then grew into a library of actions with richer routing, templates, and reporting requirements. By early April 2026, single files like `work.md`, `build-knowledge-base.md`, and later `pipeline.md` had accumulated enough embedded reference material that they became harder to load, scan, and maintain. The changelog shows a repeated pattern: extract stable scaffolding, templates, or edge-case guidance into adjacent reference files rather than keep inflating the action prompt.

This pattern is still active today. `CLAUDE.md` treats action files as standalone prompts, documents accepted variants, and lists `work-reference.md`, `bkb-reference.md`, and `pipeline-reference.md` as first-class project files. Current action prompts explicitly point readers at their companion references when they need the heavier template payload.

## Decision

The project treats each action file as the primary runnable prompt, but large reusable support material lives in companion reference files once it starts harming prompt size, scanability, or token budgets.

In practice this means:
- The action file keeps the workflow, routing, and invariants.
- Stable templates, schemas, persona packs, and error-handling tables move into adjacent `*-reference.md` files.
- The action file must explicitly tell the agent when to load the companion so the split stays discoverable.

## Alternatives

1. Keep everything in one giant action file.
This was rejected because oversized prompts become harder to load in one pass and harder to maintain safely.

2. Move details into scattered docs outside the action namespace.
This was rejected because action-specific reference material needs to stay close to the runnable prompt, not disappear into general documentation.

3. Push more behavior back into `SKILL.md`.
This was rejected because the project had already moved toward modular action ownership, not a monolithic entrypoint.

## Consequences

The positive outcome is a scalable prompt architecture: action prompts stay readable, companion files can grow independently, and new heavy-template features can follow an established extraction pattern.

The trade-off is indirection. Readers sometimes need to open two files instead of one, so each action must keep its companion link explicit and current. The repo now carries an intentional maintenance burden around cross-file consistency, but that burden is preferable to unreadable monoliths.

## References

- [CHANGELOG.md](../../CHANGELOG.md) — `0.49.0 The Architect`, `0.61.1 The Lean Cut`, `0.64.1 The Companion Split`
- [CLAUDE.md](../../CLAUDE.md) — project structure and action-file conventions
- [actions/work.md](../../actions/work.md) and [actions/work-reference.md](../../actions/work-reference.md)
- [actions/build-knowledge-base.md](../../actions/build-knowledge-base.md) and [actions/bkb-reference.md](../../actions/bkb-reference.md)
- [actions/pipeline.md](../../actions/pipeline.md) and [actions/pipeline-reference.md](../../actions/pipeline-reference.md)
