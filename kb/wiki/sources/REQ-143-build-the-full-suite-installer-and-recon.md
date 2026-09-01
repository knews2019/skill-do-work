---
title: "Lessons from REQ-143: Build the Full-Suite Installer and Reconciler"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-143-build-the-full-suite-installer-and-recon.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-143: Build the Full-Suite Installer and Reconciler

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Build the canonical fresh-install bootstrap and client-configuration reconciler for the complete four-skill suite, without activating the live modular distribution yet.

## Solution summary

[MAP CHANGED] `tools/install-do-work-suite.sh` is the staged suite’s canonical bootstrap/reconciliation boundary. It owns archive validation, full four-module installation, managed Just recipes, core-hook composition, known memory-path migration, confirmation, byte verification, and recovery; live distribution still remains in bridge mode until REQ-144.

## What worked

- Passing the already-downloaded bootstrap archive into the installer makes “one artifact” mechanically testable and avoids a second network trust boundary.
- Preparing Just and settings candidates before confirmation moves malformed-input failures ahead of all managed writes while keeping post-write recovery independently testable.

## What didn't work

- An initial bootstrap assertion accidentally matched only the generic `--` token; tightening it to the literal archive argument prevented a false-positive contract.
- Resetting signal traps to defaults before recovery let a second signal interrupt rollback; recovery must ignore termination signals until originals are restored.

## Worth knowing

- The managed-section utility requires Python for existing Justfiles. A Python-free fresh project can still copy the complete board template and take the explicit manual settings step; an existing Justfile fails unchanged rather than receiving an unsafe shell rewrite.
- Hook migration is limited to `command` string values containing the two known legacy memory paths. Other settings strings and every unrelated hook entry survive.

## Back-reference

See `do-work/archive/UR-031/REQ-143-build-full-suite-installer-reconciler.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8f22cbe`.
