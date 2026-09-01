---
title: "Lessons from REQ-144: Activate the Four-Skill Distribution"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-144-activate-the-four-skill-distribution.md]
related:
  - page: REQ-136-define-the-four-skill-suite-contract
    rel: complements
  - page: REQ-137-ship-the-suite-aware-bridge-updater
    rel: depends-on
  - page: REQ-138-add-managed-text-section-replacement
    rel: complements
  - page: REQ-139-stage-the-modular-core-skill
    rel: complements
  - page: REQ-140-stage-the-modular-board-skill
    rel: complements
  - page: REQ-143-build-the-full-suite-installer-and-recon
    rel: depends-on
  - page: REQ-145-remove-the-stateful-pipeline
    rel: complements
  - page: REQ-146-remove-modular-migration-compatibility-s
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-144: Activate the Four-Skill Distribution

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Switch the live distribution to the four staged skills and migrate bridge-enabled client installations through one all-or-recover suite transaction.

## Solution summary

[MAP CHANGED] The distribution now installs four sibling skills—core orchestration, board, knowledge, and toolbox—from one shared manifest and version, with the monolithic root runtime removed. Fresh installs and bridge upgrades converge on the same all-or-recover suite transaction, so a client either receives one verified modular configuration or its exact managed state is restored. The narrowed prime check found stale links in both touched prime surfaces; REQ-149 and REQ-150 keep the remaining compatibility and package-reference corrections visible before the migration window closes.

## What worked

- Shipping the configuration-aware bridge first let the live cutover reuse one already-trusted validate/review/install/verify/recover transaction for direct and managed-Just updates.
- Inverting the archive contract before deleting the monolith produced a precise RED signal for every root runtime surface the release had to retire.

## What didn't work

- The earlier capability token was too broad; the hermetic configuration fixture showed that a final bridge release needed to carry managed Just and hook migration before activation.
- Content-only shim assertions and a token-level reference sweep missed alias semantics and package-relative links; REQ-149 and REQ-150 preserve those review findings as explicit follow-up work.

## Worth knowing

- Runtime source ownership now lives in the four manifest-declared packages. The repository root retains only the installer, manifest validator, and managed-text replacement bootstrap utilities.
- The print-only compatibility layer is intentionally one-release work; REQ-146 removes it after the migration window, while REQ-145 removes the preserved stateful pipeline.

## Back-reference

See `do-work/archive/UR-031/REQ-144-activate-four-skill-distribution.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `a3c2612`.
