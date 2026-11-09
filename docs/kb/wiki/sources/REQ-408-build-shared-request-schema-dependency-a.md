---
title: "Lessons from REQ-408: Build shared request, schema, dependency, atomic-file, and repository packages"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-408-build-shared-request-schema-dependency-a.md]
related:
  - page: REQ-409-implement-safe-cleanup-passes-and-explic
    rel: complements
  - page: REQ-410-implement-doctor-deterministic-forensics
    rel: complements
  - page: REQ-413-implement-capture-file-answer-release-ve
    rel: complements
  - page: REQ-414-migrate-remaining-core-checks-publicatio
    rel: complements
  - page: REQ-419-add-flat-just-recipes-collision-validati
    rel: complements
  - page: REQ-420-replace-shell-implementations-with-shims
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-408: Build shared request, schema, dependency, atomic-file, and repository packages

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Create the reusable repository model required by every request and queue command.

## Solution summary

Added the standard-library shared repository-model layer for safe atomic publication and reservations, schema normalization, byte-preserving REQ/UR parsing and field edits, one-pass repository discovery/allocation, and deterministic dependency evidence. Added a compact prime index for later command families.

## What worked

- Byte-offset frontmatter edits and rooted file/directory handles kept preservation and containment contracts testable without external dependencies.
- Assertion-level adversarial fixtures exposed collision handoff and directory-swap defects that normal package tests missed.

## What didn't work

- Inode, size, and timestamp checks alone did not detect restored-metadata in-place edits; content evidence was required.
- The first review phrased portable atomic replacement as compare-and-swap against arbitrary writers, a guarantee the standard-library replacement primitives cannot provide; narrowing the contract made the real safety boundary reviewable.

## Worth knowing

Collision evidence is a first-class repository fact and should be consulted before any dependency winner is selected. Unknown schema values remain evidence even when they are unrecognized.

## Back-reference

See `do-work/archive/REQ-408-build-shared-repository-model.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `ac2e3acd`.
