---
title: "Lessons from REQ-166: Simplify session-start hook and fix dead fail-soft fallback"
type: source-summary
topic_cluster: checkpoint-and-crash-recovery
sources: [raw/processed/2026-09-01/REQ-166-simplify-session-start-hook-and-fix-dead.md]
related:
  - page: concept-session-checkpoints-and-recovery
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-166: Simplify session-start hook and fix dead fail-soft fallback

Part of the [[concept-session-checkpoints-and-recovery]] cluster.

## What the REQ was about

Fix the demonstrated bug in `skills/do-work/hooks/session-start.sh` — under `set -euo pipefail`, a failed `grep` in the `VERSION=$(grep … | sed …)` pipeline aborts the script before the `[ -z "$VERSION" ]` fallback can run — and gut the script to its minimal form. The banner's job is two lines of logic (read a version string, count queue files); the current 46-line script's defensive apparatus is what produced the bug.

## Solution summary

Reduced the SessionStart hook's runtime logic to fail-soft version and queue-count assignments under `set -u`, preserving the anchored hook command and exact banner contract. Added fixture coverage for happy path, missing version file, reformatted version label, and missing queue directory, then wired it into the aggregate contracts.

## What worked

- Writing the fixture probe before touching the hook made the failure mode undeniable: missing and reformatted version inputs exited with two distinct non-zero codes but the same missing-banner symptom.
- Testing a copied real hook inside a synthetic skill/project tree verified path derivation, environment anchoring, version parsing, queue counting, output, and exit status through the actual caller seam.

## What didn't work

- A later empty-value guard cannot make a command substitution fail-soft when `set -e` has already terminated the script; `pipefail` made the supposedly defensive `grep | sed` pipeline the abort trigger.

## Worth knowing

- Keep the `${CLAUDE_PROJECT_DIR:-.}` command anchor and its warning comment intact. The simplification belongs in runtime logic, not in the installation path that previously regressed.
- For tiny status hooks, defaults that naturally produce an empty value are safer than strict-shell defensive branches whose control flow is harder than the output contract.

## Back-reference

See `do-work/archive/UR-036/REQ-166-simplify-session-start-hook.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `6538bdd`.
