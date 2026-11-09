---
title: "Lessons from REQ-136: Define the Four-Skill Suite Contract"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-136-define-the-four-skill-suite-contract.md]
related:
  - page: REQ-137-ship-the-suite-aware-bridge-updater
    rel: complements
  - page: REQ-138-add-managed-text-section-replacement
    rel: complements
  - page: REQ-139-stage-the-modular-core-skill
    rel: complements
  - page: REQ-140-stage-the-modular-board-skill
    rel: complements
  - page: REQ-141-stage-the-modular-knowledge-skill
    rel: complements
  - page: REQ-142-stage-the-modular-toolbox-skill
    rel: complements
  - page: REQ-143-build-the-full-suite-installer-and-recon
    rel: complements
  - page: REQ-144-activate-the-four-skill-distribution
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-136: Define the Four-Skill Suite Contract

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Define the architecture and machine-readable distribution contract for a four-skill, single-version do-work suite without activating the modular layout yet.

## Solution summary

Defined the single-version four-skill layout, added a read-only fail-closed validator and security-focused fixture suite, documented the ownership and all-or-recover transaction in ADR-019, and kept all future-layout artifacts export-ignored until the bridge rollout gate is satisfied.

## What worked

- A standalone read-only validator lets future updater and installer REQs share one security boundary while fixture tests stay hermetic.

## What didn't work

- The initial six-file plan defined the future manifest but missed the current archive boundary; the RED export check exposed that shipping an incomplete staged layout would violate the bridge-release requirement.

## Worth knowing

- `VERSION`, `suite/`, and `skills/` are intentionally export-ignored staging paths until REQ-144 removes all three guards at cutover.

## Back-reference

See `do-work/archive/UR-031/REQ-136-define-four-skill-suite-contract.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `f8aecd8`.
