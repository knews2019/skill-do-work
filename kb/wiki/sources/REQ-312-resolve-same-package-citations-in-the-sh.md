---
title: "Lessons from REQ-312: Resolve same-package citations in the shipped reference contract"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-312-resolve-same-package-citations-in-the-sh.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-312: Resolve same-package citations in the shipped reference contract

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

`_dev/tests/shipped-package-reference-contract.sh` is the guard that keeps a shipped
instruction from pointing at a file that is not there once the suite is installed under
`.claude/skills/`. It resolves **cross-package** citations only. A dangling citation to a
file in the *same* package ships silently.

Measured, not inferred. Four dangling citations were planted one at a time in
`skills/do-work/actions/work-reference.md` and the contract run after each:

## Solution summary

The shipped-reference contract now recognizes same-package citations from the
suite's live manifest topology, verifies paired source and installed targets stay inside the owning
module and exist, and reports dangling citations with their original token. Regression fixtures pin
the new classifier, containment, cross-package isolation, changelog exception, and punctuation
behavior; consumer-owned and deliberately absent paths are explicitly rooted in shipped prose.

## What worked

- Deriving citation grammar from the manifest's live content directories closed the class without a
  closed enum, while paired mutation probes proved both source and installed behavior.
- Replaying the three original silent passes in a clean clone gave the Finding-Closure Ratchet an
  exact before/after oracle instead of relying on the untouched corpus alone.

## What didn't work

- The original eleven-file scope missed an existing contract pin for the last30days prose. The
  canonical gate caught the stale expectation; D-02 expanded the scope without weakening the pin.

## Back-reference

See `do-work/archive/UR-055/REQ-312-resolve-same-package-citations-in-the-shipped-reference-contract.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `99ea028`.
