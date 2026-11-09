---
title: "Lessons from REQ-061: Add a just shortcut for do-work updates"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-04/REQ-061-add-a-just-shortcut-for-do-work-updates.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-061: Add a just shortcut for do-work updates

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Extend the target repository's installed `just` recipes so users can run `just run-do-work-update` to update the project-local do-work skill programmatically, without invoking the agent-driven `do-work update` action. Install the recipe alongside the existing `just run-kanban` shortcut.

## Solution summary

**Files changed:**
- `tools/do-work-update.sh` (new) — guarded project-local updater with confirmation, rollback, runtime exclusions, and post-update audit.
- `tools/prime-do-work-update.md` (new) — maintenance map and safety invariants for the updater.
- `actions/install.md` (modified) — installs, upgrades, verifies, and documents the fourth just recipe.
- `justfile` (modified) — mirrors the shipped `run-do-work-update` recipe.
- `SKILL.md` (modified) — accepts `install run-do-work-update` as an alias for the recipe installer.
- `README.md` (modified) — document the new shortcut.
- `docs/version-guide.md` (modified) — document the new shortcut.
- `docs/board-guide.md` (modified) — keep the board shortcut recipe enumeration accurate.
- `actions/board.md` (modified) — keep the standing-shortcut recipe enumeration accurate.
- `actions/help.md` (modified) — describe the installer as including the project-local updater.
- `_dev/tests/contract-regressions.sh` (modified) — pins the recipe, executable updater, runtime exclusion, and overwrite-confirmation contract.
- `actions/version.md` (modified) — release version 0.150.13.
- `CHANGELOG.md` (modified) — record the new capability release.

## What worked

- Keeping the justfile recipe to project-relative root resolution made the same command work for both consumer installs and this repository's direct development recipe.
- A synthetic tarball update proved the safeguards that static checks cannot: prompt/interview cleanup, runtime exclusion, version verification, and rollback creation.

## What didn't work

- The first semantic-version comparison used `index` as an `awk` loop variable; macOS `awk` treats it as a built-in and rejected the script. Renaming it to `part_number` fixed the portable path and is recorded in the updater prime.

## Worth knowing

- The terminal shortcut deliberately prompts on every available update after showing the full diff. This is stricter than an unattended overwrite and preserves the existing customization-review boundary without using an agent turn.

## Back-reference

See `do-work/archive/UR-009/REQ-061-just-shortcut-for-do-work-update.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `924e668`.
