---
title: "Lessons from REQ-141: Stage the Modular Knowledge Skill"
type: source-summary
topic_cluster: knowledge-and-memory
sources: [raw/processed/2026-09-01/REQ-141-stage-the-modular-knowledge-skill.md]
related:
  - page: REQ-136-define-the-four-skill-suite-contract
    rel: depends-on
  - page: REQ-137-ship-the-suite-aware-bridge-updater
    rel: complements
  - page: REQ-138-add-managed-text-section-replacement
    rel: complements
  - page: REQ-139-stage-the-modular-core-skill
    rel: complements
  - page: REQ-140-stage-the-modular-board-skill
    rel: complements
  - page: REQ-142-stage-the-modular-toolbox-skill
    rel: complements
  - page: REQ-143-build-the-full-suite-installer-and-recon
    rel: complements
  - page: REQ-145-remove-the-stateful-pipeline
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-141: Stage the Modular Knowledge Skill

Part of the [[concept-knowledge-and-memory-systems]] cluster.

## What the REQ was about

Create a staged `skills/do-work-knowledge` package for BKB, dream, memory, interview, prompts, knowledge assets, and knowledge hooks.

## Solution summary

[MAP CHANGED] Retained knowledge now has its own staged context boundary at `skills/do-work-knowledge`. BKB synthesis, lightweight memory, consolidation, interviews, prompts, and their privacy-sensitive optional hooks live together; core can hand off consented lessons without loading or enabling the knowledge engine.

## What worked

- Keeping memory and BKB together preserves their shared ledger/value audit while removing both from core context.
- Separating the optional fragment (`memory-hooks.json`) from core's default `hooks.json` makes “installed but disabled” mechanically testable.

## What didn't work

- Copying the actions verbatim retained an obsolete rationale that memory belonged inside the monolith and left setup references pointing at the all-purpose installer; both had to become explicit package ownership.

## Worth knowing

- Exact legacy hook strings are data-migration keys. Broader path substitution risks rewriting unrelated commands in client settings.
- The raw store remains plaintext even though it is machine-local; redaction-before-truncation and never-track checks are independent defenses, not substitutes.

## Back-reference

See `do-work/archive/UR-031/REQ-141-stage-modular-knowledge-skill.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `ecd6831`.
