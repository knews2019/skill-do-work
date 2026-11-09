---
title: "Lessons from REQ-140: Stage the Modular Board Skill"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-140-stage-the-modular-board-skill.md]
related:
  - page: REQ-136-define-the-four-skill-suite-contract
    rel: depends-on
  - page: REQ-137-ship-the-suite-aware-bridge-updater
    rel: complements
  - page: REQ-138-add-managed-text-section-replacement
    rel: depends-on
  - page: REQ-139-stage-the-modular-core-skill
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

# Lessons from REQ-140: Stage the Modular Board Skill

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

Create a staged `skills/do-work-board` package that owns the board action, board documentation, Just template, and complete queue-kanban Go module.

## Solution summary

[MAP CHANGED] Queue visualization is now staged as its own compiled application under `skills/do-work-board`, including its launcher, guide, managed Just template, Go server/CLI, and embedded UI. Core supplies queue semantics and version ownership; the board package supplies every visual/testing mode without enlarging core context.

## What worked

- Exporting the committed Git tree instead of copying the dirty working directory kept unrelated REQ-134 changes out of this package while preserving every accepted board fix.
- Treating the Just template as a board-owned artifact made listener shutdown, foreign-process refusal, browser opening, and install paths testable together.

## What didn't work

- A direct copy left core-schema and version references looking local to the board package; the runtime-reference contract exposed them and the core next-version seam needed an explicit `--version-file`.

## Worth knowing

- The active board module currently has unrelated uncommitted REQ-134 changes. Before cutover, REQ-144 must synchronize any board changes committed after this staging snapshot into the package.
- The queue tool supports `--version-file`; modular core must always pass it because the board package's repo-root default is only correct in the legacy source checkout.

## Back-reference

See `do-work/archive/UR-031/REQ-140-stage-modular-board-skill.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `5e9996f`.
