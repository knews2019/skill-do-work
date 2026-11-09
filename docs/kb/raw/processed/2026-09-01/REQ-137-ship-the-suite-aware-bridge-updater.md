---
source_type: req_lesson
req_id: REQ-137
req_path: do-work/archive/UR-031/REQ-137-ship-suite-aware-bridge-updater.md
date: 2026-08-07
domain: general
module: tools
tags: [general, ship, suite, aware, bridge]
---

# Lessons from REQ-137: Ship the Suite-Aware Bridge Updater

## What the REQ was about

Ship a bridge release whose updater understands both the current all-in-one archive and the future four-skill suite, while leaving the live layout unchanged.

## Solution summary

Shipped the bridge updater without activating the modular archive. Both public update paths now share one engine that recognizes legacy and suite artifacts, validates a complete suite before mutation, warns and confirms once, verifies reviewed bytes, and automatically recovers only its explicit managed paths.

## What worked

- Turning every managed write into an explicit source/destination plan made review, dirty reporting, recovery, and verification share one boundary.
- A deterministic failing `cp` wrapper tested the real trap and Git recovery path without permission-dependent fixtures.

## What didn't work

- The first suite validator call trusted the downloaded archive's validator. An adversarial RED fixture exposed that the already-installed bridge must remain the authority for migration validation.

## Worth knowing

- The bridge requires the exact project Git worktree root and clears confirmed dirty managed content from both index and worktree before installing reviewed bytes.
- Future suite destinations are checked both textually by the manifest validator and physically against existing client-side symlinks.

## Back-reference

See `do-work/archive/UR-031/REQ-137-ship-suite-aware-bridge-updater.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `bef7334`.
