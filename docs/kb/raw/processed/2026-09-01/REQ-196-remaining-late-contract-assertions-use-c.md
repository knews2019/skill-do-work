---
source_type: req_lesson
req_id: REQ-196
req_path: do-work/archive/UR-041/REQ-196-lowercase-remaining-root-justfile-contract-paths.md
date: 2026-08-15
domain: testing
module: _dev/primes
tags: [testing, remaining, late, contract, assertions]
---

# Lessons from REQ-196: Remaining late contract assertions use capitalized root Justfile

## What the REQ was about

Replace the four remaining late `assert_contains` inputs that open `Justfile` with the tracked lowercase `justfile`, and prove those assertions execute on a case-sensitive filesystem.

## Solution summary

The canonical maintainer aggregate now uses the tracked lowercase root `justfile` in every live contract input, so its late checks behave consistently on case-sensitive and case-insensitive filesystems.

## What worked

Parsing the exact live assertion source shape provides case-sensitive evidence even when the host filesystem aliases `Justfile` and `justfile`. Pinning each expected pattern separately also proves all four late assertions remain present and reachable.

## What didn't work

Fixing only the first two occurrences in REQ-180 left four later inputs hidden by macOS case-insensitive lookup. A local green aggregate was therefore not sufficient evidence that every live path used tracked casing.

## Worth knowing

Keep intentional filename variants in prose and fixture loops; the enforceable boundary is the path argument consumed by a live root-file assertion, not every textual occurrence of “Justfile.”

## Back-reference

See `do-work/archive/UR-041/REQ-196-lowercase-remaining-root-justfile-contract-paths.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `5f15929`.
