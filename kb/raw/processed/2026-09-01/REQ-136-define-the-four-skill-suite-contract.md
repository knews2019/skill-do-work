---
source_type: req_lesson
req_id: REQ-136
req_path: do-work/archive/UR-031/REQ-136-define-four-skill-suite-contract.md
date: 2026-08-07
domain: general
module: tools
tags: [general, define, four, skill, suite]
---

# Lessons from REQ-136: Define the Four-Skill Suite Contract

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
