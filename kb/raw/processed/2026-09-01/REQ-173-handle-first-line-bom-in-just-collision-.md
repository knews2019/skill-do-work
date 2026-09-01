---
source_type: req_lesson
req_id: REQ-173
req_path: do-work/archive/UR-039/REQ-173-handle-first-line-bom-in-just-collision-scan.md
date: 2026-08-11
domain: testing
module: skills/do-work/tools
tags: [testing, handle, first, line, just]
---

# Lessons from REQ-173: Handle first-line BOM in Just collision scan

## What the REQ was about

Recognize a reserved Just recipe when the first line begins with a UTF-8 BOM, including when `just` is unavailable, without changing the target file's bytes.

## Solution summary

The fallback definition scanner removes one UTF-8 BOM only from its first-line classification value in both distributed helper copies. The no-Just installer fixture now replays a BOM-prefixed reserved recipe and verifies pre-confirmation rejection with byte- and state-preservation.

## What worked

**What worked:** Replaying the no-Just installer path preserved the real pre-confirmation and byte-identity boundaries while isolating the scanner defect.
**What didn't:** An ASCII-anchored identifier matcher silently assumed the first physical byte belonged to the Just grammar.
**Worth knowing:** UTF-8 BOM handling belongs only in the first-line classification view; the byte-preserving target and all later lines stay untouched.

## Back-reference

See `do-work/archive/UR-039/REQ-173-handle-first-line-bom-in-just-collision-scan.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8092258`.
