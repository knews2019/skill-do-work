---
title: "ADR-010: Use Typed Relationships, Retrieval Memory, and Agent Crew in BKB"
type: architecture-decision-record
status: accepted
topic_cluster: knowledge-base
decided: 2026-04-06
sources:
  - CHANGELOG.md (0.45.0 The Second Brain)
  - CHANGELOG.md (0.46.0 The Agent Crew)
  - CHANGELOG.md (0.47.0 The Full Crew)
  - actions/build-knowledge-base.md
  - actions/bkb-reference.md
  - docs/bkb-guide.md
  - SKILL.md
related:
  - page: adr-009-build-knowledge-base-as-a-compiled-interlinked-wiki
    rel: depends-on
  - page: adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats
    rel: complements
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# ADR-010: Use Typed Relationships, Retrieval Memory, and Agent Crew in BKB

Topic cluster: [[_index_knowledge-base]] ([topic index](../topics/_index_knowledge-base.md))
See also: [[adr-009-build-knowledge-base-as-a-compiled-interlinked-wiki]] (depends-on), [[adr-008-render-pipeline-debriefs-in-three-cross-linked-audience-specific-formats]] (complements)

## Context

After the compiled wiki architecture landed, the BKB needed better connective tissue and operational roles. The next changelog entries added typed `related:` relationships, a retrieval agent that learns from query history, three-tier query routing to avoid wiki bloat, an eight-agent operational crew, and later maintenance/extension features like defrag, garden, and custom agents.

Those features remain part of the current BKB design. The reference docs still require typed `related:` entries, `wiki/agent.md` hot-topic prioritization, crew dispatch tables, and maintenance sub-commands that depend on those richer semantics.

## Decision

The BKB uses more than static pages. It layers three reinforcing mechanisms onto the compiled wiki:
- typed relationships in frontmatter so pages can express how they connect,
- a retrieval-memory file (`wiki/agent.md`) that prioritizes future lookups based on prior useful queries,
- an agent crew model that assigns specialized roles to ingest, query, lint, resolve, defrag, garden, and crew-management operations.

Query output is also routed intentionally: only answers that synthesize across sources become new wiki pages, which keeps the wiki from filling up with low-value one-off responses.

## Alternatives

1. Use flat untyped links between pages.
This was rejected because maintenance, contradiction tracking, and multi-hop query behavior all benefit from explicit relationship semantics.

2. Treat every query answer as a new wiki page.
This was rejected because it would create wiki bloat and lower the value density of stored knowledge.

3. Run all BKB operations through one generic persona.
This was rejected because the project wanted stable role boundaries for structure, cross-linking, QA, and readability.

## Consequences

The knowledge base becomes richer and more maintainable: links carry meaning, queries learn from history, and operational roles stay explicit. It also makes advanced maintenance commands like `defrag`, `garden`, and custom crew extensions coherent instead of bolted on.

The cost is more metadata and more operational surface area. Typed relationships, agent logs, and crew dispatch rules all need to stay in sync with the underlying pages.

## References

- [CHANGELOG.md](../../CHANGELOG.md) — `0.45.0 The Second Brain`, `0.46.0 The Agent Crew`, `0.47.0 The Full Crew`
- [actions/build-knowledge-base.md](../../actions/build-knowledge-base.md)
- [actions/bkb-reference.md](../../actions/bkb-reference.md)
- [docs/bkb-guide.md](../../docs/bkb-guide.md)
- [SKILL.md](../../SKILL.md) — BKB help and sub-command surface
