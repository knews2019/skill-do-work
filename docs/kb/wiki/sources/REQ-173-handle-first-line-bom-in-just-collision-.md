---
title: "Lessons from REQ-173: Handle first-line BOM in Just collision scan"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-173-handle-first-line-bom-in-just-collision-.md]
related:
  - page: REQ-172-make-screenshot-source-cleanup-best-effo
    rel: complements
  - page: REQ-174-validate-root-markdown-fence-info
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-173: Handle first-line BOM in Just collision scan

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

Recognize a reserved Just recipe when the first line begins with a UTF-8 BOM, including when `just` is unavailable, without changing the target file's bytes.

## Solution summary

The fallback definition scanner removes one UTF-8 BOM only from its first-line classification value in both distributed helper copies. The no-Just installer fixture now replays a BOM-prefixed reserved recipe and verifies pre-confirmation rejection with byte- and state-preservation.

## What worked

Replaying the no-Just installer path preserved the real pre-confirmation and byte-identity boundaries while isolating the scanner defect.

## What didn't work

An ASCII-anchored identifier matcher silently assumed the first physical byte belonged to the Just grammar.

## Worth knowing

UTF-8 BOM handling belongs only in the first-line classification view; the byte-preserving target and all later lines stay untouched.

## Back-reference

See `do-work/archive/UR-039/REQ-173-handle-first-line-bom-in-just-collision-scan.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8092258`.
