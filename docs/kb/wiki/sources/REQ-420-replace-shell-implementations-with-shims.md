---
title: "Lessons from REQ-420: Replace shell implementations with shims and prove whole-suite parity"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-420-replace-shell-implementations-with-shims.md]
related:
  - page: REQ-408-build-shared-request-schema-dependency-a
    rel: complements
  - page: REQ-409-implement-safe-cleanup-passes-and-explic
    rel: complements
  - page: REQ-410-implement-doctor-deterministic-forensics
    rel: complements
  - page: REQ-413-implement-capture-file-answer-release-ve
    rel: complements
  - page: REQ-414-migrate-remaining-core-checks-publicatio
    rel: complements
  - page: REQ-419-add-flat-just-recipes-collision-validati
    rel: depends-on
  - page: REQ-478-capture-stamps-required-lessons-under-a-
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-420: Replace shell implementations with shims and prove whole-suite parity

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Complete the migration by making every retained shell path a thin launcher and enforcing full suite parity mechanically.

## Solution summary

Migrated all 41 retained shell entry points to thin, argv-preserving compatibility launchers; completed the corresponding typed Go behavior and exact compatibility seams; consolidated audit-metrics into `do-work-cli`; added a legacy-first 110-case differential/parity oracle plus a mechanical thinness ratchet; and reconciled SessionStart, board, BKB, recovery-command, package-reference, and maintainer-gate authority.

## What worked

Capturing the legacy behavior as independent fixture expectations before cutover made it possible to delete shell domain logic without turning parity into a Go-versus-Go comparison. A single path-to-command inventory then drove thinness, parity, and staged-package checks.

## What didn't work

Several old fixture seams intercepted shell utilities or asserted literal shell constants; preserving those shapes would have kept domain ownership in shell. Moving those cases to observable status/output/argv/filesystem behavior removed the contradiction while retaining the earned oracle.

## Worth knowing

Public compatibility belongs at the launcher/runtime boundary. Typed Go owners should carry mutations, rollback, findings, and exact next/verify argv; tests should assert those effects rather than the former shell implementation shape.

## Back-reference

See `do-work/archive/REQ-420-replace-shell-implementations-verify-parity.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8fcba42f`.
