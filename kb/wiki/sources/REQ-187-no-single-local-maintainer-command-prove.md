---
title: "Lessons from REQ-187: No single local maintainer command proves shell plus both Go modules"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-187-no-single-local-maintainer-command-prove.md]
related:
  - page: concept-prescribed-shell-commands
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-187: No single local maintainer command proves shell plus both Go modules

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

Add one export-ignored local maintainer verification command as the source of truth for strict shell checks, the aggregate once, and vet/test in both Go modules; make documentation and any root recipe delegate to it.

## Solution summary

**[MAP CHANGED]** Repository hand-back health now has one local entrypoint: `_dev/tests/maintainer-verify.sh`. It owns strict shell verification, the aggregate, both Go modules, and available Node-backed board behavior; CLAUDE and Just only delegate to it. The prescribed-shell prime remains current and gains this ownership/recursion lesson after archive.

## What worked

- A canonical verification wrapper prevents command-list drift only when documentation and recipes are thin delegates and its behavioral contract counts/fails semantic stages rather than copying a second production command inventory.
- When a canonical gate must invoke an aggregate that also tests the gate, give the gate a fixture-only self-test mode that cannot reach the real aggregate. That closes the recursive ownership loop without duplicating a normal execution edge.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.

## Back-reference

See `do-work/archive/UR-041/REQ-187-canonical-local-maintainer-gate.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `c20110d`.
