---
source_type: req_lesson
req_id: REQ-143
req_path: do-work/archive/UR-031/REQ-143-build-full-suite-installer-reconciler.md
date: 2026-08-07
domain: general
module: tools
tags: [general, build, full, suite, installer]
---

# Lessons from REQ-143: Build the Full-Suite Installer and Reconciler

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
