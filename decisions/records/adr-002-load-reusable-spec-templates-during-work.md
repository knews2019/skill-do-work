---
title: "ADR-002: Load Reusable Spec Templates During Work"
type: architecture-decision-record
status: accepted
topic_cluster: skill-architecture
decided: 2026-04-10
sources:
  - CHANGELOG.md (0.59.0 The Quality Blueprint)
  - CLAUDE.md
  - README.md
  - actions/capture.md
  - actions/work.md
  - specs/README.md
related:
  - page: adr-001-modular-action-prompts-and-companion-references
    rel: extends
  - page: adr-003-always-load-karpathy-guardrails
    rel: complements
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# ADR-002: Load Reusable Spec Templates During Work

Topic cluster: [[_index_skill-architecture]] ([topic index](../topics/_index_skill-architecture.md))
See also: [[adr-001-modular-action-prompts-and-companion-references]] (extends), [[adr-003-always-load-karpathy-guardrails]] (complements)

## Context

The work action needed a way to encode recurring implementation standards without bloating every REQ or forcing the builder to rediscover expectations for common task shapes. The changelog entry that introduced `specs/` framed them as reusable blueprints for common work such as API endpoints, UI components, refactors, and bug fixes.

The current repo still uses that mechanism. `CLAUDE.md` lists `specs/` as a first-class directory, `capture.md` can hint a `suggested_spec`, `README.md` explains that specs are auto-loaded, and `work.md` still checks both task type and frontmatter hints before handing guidance to builder and reviewer.

## Decision

The project keeps reusable task-type guidance in `specs/` and loads those templates during work when either:
- the REQ clearly matches a known task type, or
- capture already recorded a `suggested_spec` hint.

Specs are advisory structure, not a replacement for reading the actual REQ. They exist to standardize output shape, quality bars, checklists, and common pitfalls for repeated classes of work.

## Alternatives

1. Put all task-type guidance directly into `work.md`.
This was rejected because it would keep inflating the busiest action file and mix reusable guidance with orchestration logic.

2. Encode task-type guidance informally in prose across docs and examples.
This was rejected because informal guidance is harder to discover and easier to apply inconsistently.

3. Require every REQ author to restate the full specification pattern manually.
This was rejected because it duplicates effort and weakens consistency across the queue.

## Consequences

The result is more consistent implementation and review behavior for recurring task types, with lighter REQ authoring overhead and a clearer extension point for future task templates.

The cost is template upkeep. Specs must stay aligned with the real workflow, and the team has to decide carefully when a task truly matches a template versus when the hint should be ignored.

## References

- [CHANGELOG.md](../../CHANGELOG.md) — `0.59.0 The Quality Blueprint`
- [CLAUDE.md](../../CLAUDE.md) — `specs/` project-structure entry
- [README.md](../../README.md) — spec-loading note
- [actions/capture.md](../../actions/capture.md) — `suggested_spec` hinting
- [actions/work.md](../../actions/work.md) — spec-loading steps
- [specs/README.md](../../specs/README.md)
