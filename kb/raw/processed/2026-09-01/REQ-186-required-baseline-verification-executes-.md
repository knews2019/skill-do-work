---
source_type: req_lesson
req_id: REQ-186
req_path: do-work/archive/UR-041/REQ-186-baseline-suite-single-ownership.md
date: 2026-08-15
domain: testing
module: _dev/primes
tags: [testing, required, baseline, verification, executes]
---

# Lessons from REQ-186: Required baseline verification executes two child suites twice

## What the REQ was about

Give each duplicated baseline child suite one owner: remove the aggregate's redundant direct prescribed-shell edge and stop requiring a second standalone shipped-reference invocation when it has no distinct mode or fixture.

## Solution summary

**Verified unchanged:** `_dev/tests/staged-skills-contract.sh` still invokes `prescribed-shell-scripts-behavior.sh` for standalone-semantics coverage.

## Worth knowing

- A required aggregate should give each identical child invocation one owner. If a nested suite already preserves the needed standalone semantics and failure propagation, a second direct edge adds runtime without adding evidence.
- Maintainer hand-back instructions should name the aggregate baseline once and reserve standalone child runs for genuinely distinct modes, fixtures, or touched-area focus.

## Back-reference

See `do-work/archive/UR-041/REQ-186-baseline-suite-single-ownership.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0ab2b79`.
