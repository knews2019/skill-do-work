---
title: "Lessons from REQ-167: Deduplicate copy-pasted shell primitives across action files"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-167-deduplicate-copy-pasted-shell-primitives.md]
related:
  - page: concept-prescribed-shell-commands
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-167: Deduplicate copy-pasted shell primitives across action files

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

Sweep the shipped `skills/` tree for shell primitives that are restated in multiple action files (the CLAUDE.md trap-list primitives are the starting inventory: untracked-file enumeration, merge-commit-safe `git show`, root-anchored ignore patterns, curl download-and-rename, `git diff-tree` file listing, etc.) and give each one exactly one canonical home; other sites reference it instead of restating it.

## Solution summary

Established a single shipped core guide for eight prescribed shell primitives and a durable primitive-to-home/former-site audit. Replaced duplicated cross-package rationale with explicit pointers while preserving executable commands and caller-specific policy, then added a regression probe that enforces the guide headings, all former-site pointers, and absence of known stale rationale phrases.

## What worked

- Separating “command executions” from “rationale copies” kept standalone actions runnable while still removing the maintenance multiplier that caused prior partial fixes.
- A phrase ratchet plus required pointers catches both obvious copy-back and silent removal of the canonical trail; the shell-fence harness independently protects the commands that remain.

## What didn't work

- Grep count alone over-reported duplication because legitimate caller commands and shared explanation use the same tokens. The inventory needed a disposition column before deletion decisions were safe.

## Worth knowing

- Cross-package shell policy belongs in core because every sibling can depend on core, while moving a shared primitive into toolbox/board/knowledge would reverse the allowed dependency direction.
- Durable audit artifacts are evidence entry points, not runtime dependencies; a no-static-reference qualifier warning is expected for this requested document class.

## Back-reference

See `do-work/archive/UR-036/REQ-167-dedupe-prescribed-shell-primitives.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `1a27c07`.
