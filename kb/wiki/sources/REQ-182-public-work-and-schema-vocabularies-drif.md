---
title: "Lessons from REQ-182: Public work and schema vocabularies drift while suites stay green"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-182-public-work-and-schema-vocabularies-drif.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-182: Public work and schema vocabularies drift while suites stay green

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Restore parity at the public work-guide/router and testing-schema/normalizer seams, and correct the two short workflow summaries that omit canonical states while the baseline suites remain green.

## Solution summary

**Behavior:** Public aliases and testing-status aliases now have one documented inventory each plus executable parity mirrors; any one-sided addition, removal, or testing-alias remap fails the existing contract suite. Queue summaries no longer hide dependency-cycle holds.

## What worked

- When prose is intentionally authoritative but runtime must remain independently readable, a seam-local exact comparison with bilateral mutation probes is enough to prevent drift without introducing a generator.
- A duplicate public inventory is itself a third drift surface; replace it with an anchored pointer before asserting parity between the remaining owner and mirror.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.

## Back-reference

See `do-work/archive/UR-041/REQ-182-public-vocabulary-parity.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `9ebdd06`.
