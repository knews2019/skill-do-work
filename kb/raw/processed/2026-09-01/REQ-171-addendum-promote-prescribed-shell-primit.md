---
source_type: req_lesson
req_id: REQ-171
req_path: do-work/archive/UR-038/REQ-171-promote-prescribed-shell-to-tested-scripts.md
date: 2026-08-11
domain: general
module: skills/do-work/general
tags: [general, addendum, promote, prescribed, shell]
---

# Lessons from REQ-171: Addendum: promote prescribed shell primitives to shipped, fixture-tested scripts

## What the REQ was about

Graduate the canonical prescribed-shell primitives from documented prose to real, shipped script files with fixture-repo execution tests. Each *multi-line* primitive in `skills/do-work/docs/prescribed-shell-primitives.md` (and the remaining multi-line blocks in action files, e.g. capture.md Step 4's screenshot copy/verify/link block) becomes a `.sh` file under a per-package `scripts/` directory (core first: `skills/do-work/scripts/`); call sites keep a one-line intent statement plus the invocation; `_dev/tests/` gains a fixture-repo scaffold (mktemp repo, git init, seeded queue/version fixtures) that *executes* each script and asserts output and exit codes. The dividing line is "does this block contain logic that can be wrong" — one-liners and illustrative fragments stay inline, covered by the existing lint harness as residue.

## Solution summary

- The executable ownership map is `skills/do-work/docs/prescribed-shell-primitives.md`; the durable 17/21/2 census lives in `decisions/audits/2026-08-11-prescribed-shell-primitives.md`.
- Start behavioral changes in `_dev/tests/prescribed-shell-scripts-behavior.sh`, then keep canonicalization, staged-package, and action-shell lint green.

## What worked

- Shell quoting and Git path literalness are different layers: quotes stop the shell, but Git pathspec magic still needs an explicit literal-path boundary.
- Promotion inventories must mirror exact changed paths in both Scope and `write_set`; grouped glob prose is useful explanation but not an auditable boundary.
- Migrating inline logic safely means moving the original adversarial fixture with it—the coordinated screenshot race was the evidence that made the conflict resolution deterministic.

## Back-reference

See `do-work/archive/UR-038/REQ-171-promote-prescribed-shell-to-tested-scripts.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `5a18faf`.
