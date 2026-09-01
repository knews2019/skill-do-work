---
source_type: req_lesson
req_id: REQ-182
req_path: do-work/archive/UR-041/REQ-182-public-vocabulary-parity.md
date: 2026-08-15
domain: testing
module: _dev/primes
tags: [testing, public, work, schema, vocabularies]
---

# Lessons from REQ-182: Public work and schema vocabularies drift while suites stay green

## What the REQ was about

Restore parity at the public work-guide/router and testing-schema/normalizer seams, and correct the two short workflow summaries that omit canonical states while the baseline suites remain green.

## Solution summary

**Behavior:** Public aliases and testing-status aliases now have one documented inventory each plus executable parity mirrors; any one-sided addition, removal, or testing-alias remap fails the existing contract suite. Queue summaries no longer hide dependency-cycle holds.

## Worth knowing

- When prose is intentionally authoritative but runtime must remain independently readable, a seam-local exact comparison with bilateral mutation probes is enough to prevent drift without introducing a generator.
- A duplicate public inventory is itself a third drift surface; replace it with an anchored pointer before asserting parity between the remaining owner and mirror.

## Back-reference

See `do-work/archive/UR-041/REQ-182-public-vocabulary-parity.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `9ebdd06`.
