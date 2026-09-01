---
source_type: req_lesson
req_id: REQ-309
req_path: do-work/archive/UR-055/REQ-309-run-the-repos-canonical-gate-before-hand-back.md
date: 2026-08-21
domain: general
module: _dev/primes
tags: [general, repo, canonical, gate, hand]
---

# Lessons from REQ-309: Run the repo's canonical gate before hand-back, not only the changed area's tests

## What the REQ was about

REQ-283 archived with a Testing section listing four green checks — `go test ./...` and a fresh build in `skills/do-work-board/tools/queue-kanban`, plus `queue-kanban verify` returning `OK: no findings`. Every one of them was true. None of them was `bash _dev/tests/maintainer-verify.sh`, which the change had just turned red by adding a second `./actions/board.md` routing row that `_dev/tests/staged-skills-contract.sh` counts.

The gate stayed red across REQ-279, REQ-295 and REQ-283's own metadata commits, and REQ-262's run was the first to notice — because Step 5.75's pre-flight ran the command and reported the baseline failing.

## Solution summary

[MAP CHANGED] The work action now inherits any canonical repository-wide gate explicitly declared
by project guidance and requires its direct zero exit status on the final tree, independently of
per-request prime selection. The contract-regression suite mutation-tests that policy; REQ-317 will
align the remaining downstream error-handling reader.

## What worked

**What worked:** A semantic detector plus adversarial mutations turned a prose policy into an
executable contract, and replaying REQ-283 proved the new gate catches the exact escaped defect.

**What didn't:** Testing only the newly inserted Step 6.5 lane missed a later generic Error Handling
row that can oppose it. REQ-317 carries that downstream-reader reconciliation.

**Worth knowing:** Focused tests establish attribution; a project-declared repository gate establishes
whether the final tree is hand-backable. They are complementary verdicts, not substitutes.

## Back-reference

See `do-work/archive/UR-055/REQ-309-run-the-repos-canonical-gate-before-hand-back.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `d9bf150`.
