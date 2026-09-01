---
title: "Lessons from REQ-165: Shell-block lint harness for shipped action files"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-165-shell-block-lint-harness-for-shipped-act.md]
related:
  - page: REQ-166-simplify-session-start-hook-and-fix-dead
    rel: complements
  - page: REQ-167-deduplicate-copy-pasted-shell-primitives
    rel: complements
  - page: REQ-168-delete-or-test-audit-of-defensive-code-i
    rel: complements
  - page: REQ-170-finding-closure-ratchet-and-canonical-ea
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-165: Shell-block lint harness for shipped action files

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

Add a `_dev/tests/` check that extracts every fenced shell block (` ```bash ` / ` ```sh `) from the shipped `skills/` tree — action files, crew files, and shipped hook scripts — and lints them: `bash -n` (syntax) always, `shellcheck` when available. This makes the suite's largest defect generator (prose-prescribed shell that nothing executes until an agent runs it in a consumer repo) testable in CI.

## Solution summary

Added an attributable Bash/ShellCheck harness for all shipped Markdown shell fences and shipped shell source files, including valid fence indentation, narrow placeholder normalization, a deliberately broken indented self-test fixture, an explicit no-ShellCheck degradation path, and aggregate contract-suite invocation. Fixed the single real warning exposed by the complete scan by making the board action's prescribed repository-root directory change fail fast.

## What worked

- Pairing the checker with a deliberately broken fixture proved the failure path, while an independent count of Markdown-valid fences tested the discovery boundary the checker could not validate about itself.
- Keeping the new probe behind `_dev/tests/contract-regressions.sh` matched the repository's explicit child-probe convention and made the ratchet part of the normal suite.

## What didn't work

- The first extractor regex matched only column-zero fences. Its own clean run looked convincing while silently skipping 10 indented fences; comparing 49 reported blocks with an independent 59-block enumeration exposed the blind spot.

## Worth knowing

- SC2034 and SC2154 are excluded only for isolated Markdown fences because assignments and uses may live in separate prescribed blocks. Complete shipped `.sh` files receive the full warning set.
- A newly complete scan can reveal genuine pre-existing findings. Fix a narrow blocking defect, as with the board action's unguarded `cd`, instead of weakening the checker to preserve a green baseline.

## Back-reference

See `do-work/archive/UR-036/REQ-165-shell-block-lint-harness.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `a45d5c4`.
