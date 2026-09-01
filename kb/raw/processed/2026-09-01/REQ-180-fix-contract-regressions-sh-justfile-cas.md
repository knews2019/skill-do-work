---
source_type: req_lesson
req_id: REQ-180
req_path: do-work/archive/UR-040/REQ-180-contract-suite-justfile-case-mismatch.md
date: 2026-08-15
domain: general
module: _dev/primes
tags: [general, contract, regressions, justfile, case]
---

# Lessons from REQ-180: Fix contract-regressions.sh Justfile case mismatch aborting late checks

## What the REQ was about

`_dev/tests/contract-regressions.sh:1797` and `:1804` reference `Justfile` (capital J), but the tracked file is lowercase `justfile`. On a case-sensitive filesystem, `extract_kanban_shutdown_line Justfile` hits `awk: cannot open .../Justfile` and the suite aborts with exit 2 at that point — every check after line ~1797 (roughly 1,500 lines including the Common Rationalizations regrowth ratchet) silently never runs. On case-insensitive filesystems (macOS default) the mismatch is invisible, which is why it survived.

## Solution summary

Replaced the two capitalized `Justfile` inputs with the tracked lowercase `justfile`, keeping the existing shutdown-line comparison and negative assertions unchanged.

## What worked

**What worked:** The captured case-sensitive Linux failure named the exact byte-level mismatch, so the repair stayed at two literals and the full suite proved late checks were reachable again.
**What didn't:** Reproducing RED on a default macOS filesystem was not meaningful because case-insensitive lookup masks the bug; a disk-image workaround added ceremony without improving the captured evidence.
**Worth knowing:** Shell test fixtures and prescribed paths must use the tracked filename's exact casing even when a developer filesystem accepts variants.

## Back-reference

See `do-work/archive/UR-040/REQ-180-contract-suite-justfile-case-mismatch.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `2e5a5c4`.
