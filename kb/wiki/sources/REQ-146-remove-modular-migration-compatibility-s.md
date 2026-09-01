---
title: "Lessons from REQ-146: Remove Modular-Migration Compatibility Shims"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-146-remove-modular-migration-compatibility-s.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-146: Remove Modular-Migration Compatibility Shims

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Remove transitional core routes and migration rules after one modular release and confirmed client migration.

## Solution summary

Removed the one-release core command shims and every bridge-only updater/configuration migration path. The surviving current-suite updater continues through the installed validator and all-or-recover installer, marker-owned Just reconciliation, current hook composition, and permanent pipeline-guard cleanup, with updated absence and ownership contracts.

## What worked

- Deleting the compatibility branches while keeping the installed validator/installer as the sole transaction owner made the updater smaller without weakening confirmation, recovery, or managed-marker guarantees.
- Byte-identity checks for paired tools plus focused RED/GREEN contract suites kept the root and staged packages aligned through a large subtractive change.

## What didn't work

- Treating `just` as an optional syntax check also made reserved-recipe collision detection optional; the independent no-`just` smoke test exposed the invalid-file path that the main fixtures missed.
- Removing router code alone did not retire the contract: a restatement sweep still found stale commands in shipped templates, UI hints, a hook message, and transition-era prime lessons.

## Worth knowing

- The permanent updater delegates one validated archive to the installed all-or-recover installer; future work should extend that boundary rather than recreate updater-side migration logic.
- REQ-152 owns tool-independent reserved-recipe collision rejection, and REQ-153 owns the remaining retired-command and prime-restatement sweep.

## Back-reference

See `do-work/archive/UR-031/REQ-146-remove-modular-migration-shims.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0b9bcde`.
