---
title: "ADR-003: Always Load Karpathy Guardrails"
type: architecture-decision-record
status: accepted
topic_cluster: skill-architecture
decided: 2026-04-12
sources:
  - CHANGELOG.md (0.62.0 The Karpathy Nod)
  - CHANGELOG.md (0.62.1 The Senior Engineer Test)
  - CLAUDE.md
  - crew-members/karpathy.md
  - actions/work.md
  - actions/review-work.md
related:
  - page: adr-002-load-reusable-spec-templates-during-work
    rel: complements
  - page: adr-007-close-the-pipeline-with-present-and-a-technical-debrief
    rel: complements
created: 2026-04-15
updated: 2026-04-15
confidence: high
---

# ADR-003: Always Load Karpathy Guardrails

Topic cluster: [[_index_skill-architecture]] ([topic index](../topics/_index_skill-architecture.md))
See also: [[adr-002-load-reusable-spec-templates-during-work]] (complements), [[adr-007-close-the-pipeline-with-present-and-a-technical-debrief]] (complements)

## Context

The do-work skill already had workflow machinery for capture, build, and review, but it wanted an always-on behavior layer that shaped how code gets written. The changelog introduced Karpathy-derived principles as a crew member loaded for every implementation pass, then immediately tightened the model by making review-work audit those principles as an informational check.

That decision remains visible in the current repo. `CLAUDE.md` lists `karpathy.md` as always loaded in Step 6, `crew-members/karpathy.md` frames the rules as orthogonal to the workflow machinery, and `review-work.md` still includes a dedicated Karpathy Principle Check.

## Decision

Every implementation run loads the Karpathy crew member alongside the general rules, regardless of domain. Review-work then performs a lightweight Karpathy Principle Check to surface blind spots without re-penalizing issues already captured by the main rubric.

The guardrails are intentionally behavioral: think before coding, prefer simplicity, make surgical changes, and stay goal-driven. They complement the workflow instead of replacing it.

## Alternatives

1. Treat the Karpathy principles as optional reading for some domains only.
This was rejected because the team wanted a universal behavior baseline, not another domain-specific toggle.

2. Bake the principles into long prose inside `work.md` only.
This was rejected because a standalone crew member is easier to load consistently and evolve independently.

3. Score the same issue twice in review, once in the normal rubric and once in a Karpathy rubric.
This was rejected because it would create noisy, redundant feedback.

## Consequences

The project now has a stable "how we code" layer that travels with every implementation run and surfaces in review. That improves consistency across domains and agents.

The trade-off is more governance material to keep calibrated. The review pass has to stay disciplined about using the Karpathy check as a mnemonic overlay rather than a second punishment system.

## References

- [CHANGELOG.md](../../CHANGELOG.md) — `0.62.0 The Karpathy Nod`, `0.62.1 The Senior Engineer Test`
- [CLAUDE.md](../../CLAUDE.md) — always-loaded crew-member rules
- [crew-members/karpathy.md](../../crew-members/karpathy.md)
- [actions/work.md](../../actions/work.md)
- [actions/review-work.md](../../actions/review-work.md)
