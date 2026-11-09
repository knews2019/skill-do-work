---
source_type: req_lesson
req_id: REQ-141
req_path: do-work/archive/UR-031/REQ-141-stage-modular-knowledge-skill.md
date: 2026-08-07
domain: general
module: skills/do-work/general
tags: [general, stage, modular, knowledge, skill]
---

# Lessons from REQ-141: Stage the Modular Knowledge Skill

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
