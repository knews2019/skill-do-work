---
id: REQ-152
title: "Review fix: Reject reserved Just recipe collisions without Just"
status: pending
domain: general
created_at: 2026-08-08T18:32:01Z
user_request: UR-031
addendum_to: REQ-146
review_generated: true
effort_estimate: normal
---

# Review Fix: Reject Reserved Just Recipe Collisions Without Just

## What
Make the suite installer reject marker-free Justfiles that already define suite-reserved recipes even when the optional `just` executable is unavailable, before confirmation or mutation can leave an invalid file.

## Context
Found during review of REQ-146. The installer currently relies on `just --list` for collision detection; without that executable, a pre-existing `run-kanban` recipe is duplicated and the installer reports success even though the resulting Justfile will not parse when Just is later installed.

This is a standalone user-visible bootstrap defect rather than part of the retired-command restatement sweep: it changes validation behavior in the paired installer and its focused regression.

## Requirements
- Detect collisions with every suite-reserved Just recipe without requiring the `just` executable.
- Reject the install before confirmation or file mutation, with an actionable error naming the collision.
- Preserve marker-managed replacement, unrelated custom recipes/content, file modes, cancellation, idempotence, and all-or-recover behavior.
- Update both installer copies identically.
- Add a no-`just` regression that starts with a marker-free reserved recipe and proves the installer fails without changing the Justfile.
- Keep existing jq/Python pipeline-guard cleanup and all installer/update contract suites passing.
