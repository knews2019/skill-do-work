---
title: "ADR-009: Build Knowledge Base as a Compiled Interlinked Wiki"
type: architecture-decision-record
status: accepted
topic_cluster: knowledge-base
decided: 2026-04-05
sources:
  - CHANGELOG.md (0.43.0 The Knowledge Forge)
  - CHANGELOG.md (0.43.1 The Clean Handoff)
  - CHANGELOG.md (0.43.2 The Gap Closer)
  - CHANGELOG.md (0.43.4 The Guard Dog)
  - CHANGELOG.md (0.43.5 The Final Polish)
  - actions/build-knowledge-base.md
  - actions/bkb-reference.md
  - docs/bkb-guide.md
related:
  - page: adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb
    rel: complements
  - page: adr-001-modular-action-prompts-and-companion-references
    rel: complements
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# ADR-009: Build Knowledge Base as a Compiled Interlinked Wiki

Topic cluster: [[_index_knowledge-base]] ([topic index](../topics/_index_knowledge-base.md))
See also: [[adr-010-use-typed-relationships-retrieval-memory-and-agent-crew-in-bkb]] (complements), [[adr-001-modular-action-prompts-and-companion-references]] (complements)

## Context

The repo introduced a dedicated BKB command to turn raw source material into a structured knowledge base rather than re-derive answers from scratch each time. Early releases then hardened that shape by fixing file lifecycle mistakes, removing ghost folders, clarifying stable `raw/processed/` source paths, and strengthening schema rules around logs, confidence, and source handling.

The architecture is still present in the current repo. `build-knowledge-base.md`, `bkb-reference.md`, and `docs/bkb-guide.md` all describe the same compiled-wiki structure: `wiki/_master_index.md`, topic indexes, per-page frontmatter, `wiki/log.md`, and source movement from inbox/capture into processed storage.

## Decision

The BKB is a compiled Markdown wiki, not an ad hoc note pile. Raw material flows through a lifecycle (`inbox` -> `capture` -> `processed`), while the knowledge layer is organized as:
- a top-level `_master_index.md`,
- topic-cluster indexes,
- interlinked articles,
- daily/monthly maintenance logs,
- an append-only `log.md`, and
- a schema file that defines frontmatter and workflow rules.

Ingest must move sources to their stable processed location, and wiki pages cite that stable path rather than transient staging paths.

## Alternatives

1. Keep a flat folder of notes without explicit indexes or lifecycle stages.
This was rejected because the knowledge base is meant to scale and remain queryable over time.

2. Answer BKB queries directly from raw sources every time.
This was rejected because the project explicitly wants compounding knowledge, not repeated cold reads.

3. Copy sources into processed storage while leaving capture copies behind.
This was rejected because it re-queues already-ingested material and muddies provenance.

## Consequences

The result is a durable knowledge-base architecture with explicit provenance, navigable indexes, and a clear maintenance surface. It also gives the project a reusable wiki pattern that this ADR log can now borrow.

The trade-off is structural overhead. The BKB requires schema discipline, log upkeep, and source-lifecycle correctness that a looser note system would avoid.

## References

- [CHANGELOG.md](../../CHANGELOG.md) — `0.43.0 The Knowledge Forge`, `0.43.1 The Clean Handoff`, `0.43.2 The Gap Closer`, `0.43.4 The Guard Dog`, `0.43.5 The Final Polish`
- [actions/build-knowledge-base.md](../../actions/build-knowledge-base.md)
- [actions/bkb-reference.md](../../actions/bkb-reference.md)
- [docs/bkb-guide.md](../../docs/bkb-guide.md)
