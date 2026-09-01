---
title: "Lessons from REQ-142: Stage the Modular Toolbox Skill"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-142-stage-the-modular-toolbox-skill.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-142: Stage the Modular Toolbox Skill

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Create a staged `skills/do-work-toolbox` package preserving the current optional reviews, discovery, reporting, repository utilities, and companion installers.

## Solution summary

[MAP CHANGED] Optional reviews, reporting, exploration, repository hygiene, and companion installation now have an independently loadable staged context boundary at `skills/do-work-toolbox`. Core queue actions remain the authority for capture and execution; board and knowledge stay separate siblings.

## What worked

- Treating copied prose references as executable dependencies exposed the exact places where an action still assumed a monolithic skill root.
- Keeping command examples canonical to `do-work-toolbox` makes the new ownership boundary visible without deleting any legacy route before cutover.

## What didn't work

- A verbatim copy left the companion installer owning board recipes, memory setup, and core updating; those workflows needed an explicit ownership prune.

## Worth knowing

- Some `docs/...` strings in tidy-repo are intended project destinations, not packaged runtime documents; the reference contract must distinguish outputs from dependencies.
- Toolbox actions still read core queue records, but they do so through explicit sibling references rather than duplicating lifecycle logic.

## Back-reference

See `do-work/archive/UR-031/REQ-142-stage-modular-toolbox-skill.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `df35345`.
