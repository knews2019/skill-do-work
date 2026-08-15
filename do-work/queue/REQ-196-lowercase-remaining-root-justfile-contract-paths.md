---
id: REQ-196
title: Remaining late contract assertions use capitalized root Justfile
status: pending
created_at: 2026-08-15T12:11:13Z
user_request: UR-041
addendum_to: REQ-180
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
effort_estimate: normal
related: [REQ-180, REQ-187]
write_set: [_dev/tests/contract-regressions.sh]
---

# Remaining Late Contract Assertions Use Capitalized Root Justfile

## What

Replace the four remaining late `assert_contains` inputs that open `Justfile` with the tracked lowercase `justfile`, and prove those assertions execute on a case-sensitive filesystem.

## Why

REQ-180 repaired the two earlier casing mismatches that aborted the aggregate, but four later assertions still open a nonexistent capital-J root path. Default macOS filesystems resolve the wrong spelling, so the aggregate can pass locally while those assertions fail on Linux or another case-sensitive filesystem.

## Context

- Discovered during REQ-187 while tracing the root Just delegate and intentionally left outside that canonical-gate Scope.
- The remaining inputs are the four `"Justfile"` arguments near the final managed-marker, board-source, and updater-fallback assertions in `_dev/tests/contract-regressions.sh`.
- This is an addendum to REQ-180's exact tracked-casing repair.

## Detailed Requirements

- Change all four remaining root-file inputs from `Justfile` to `justfile` without changing their assertion meanings.
- Keep intentional prose and fixture references to the general `Justfile` filename variant unchanged.
- Add or reuse case-sensitive evidence that the four late assertions read the tracked root file and remain reachable.
- Preserve the canonical maintainer gate and aggregate final pass marker.

## Constraints

- Do not add case-insensitive path resolution or filesystem-specific production logic.
- Do not rewrite unrelated `Justfile` prose or multi-variant fixtures.
- Lock-in limit: zero live contract inputs open a root filename whose case differs from the tracked `justfile`.

## Red-Green Proof

**RED prompt/case:** Run the late root-justfile assertions against a case-sensitive checkout; all four capital-J inputs name a nonexistent file.
**Why RED now:** The current checkout tracks only lowercase `justfile`, while the four inputs survive because the development filesystem is case-insensitive.
**GREEN when:** All live root inputs use exact tracked casing, the focused case-sensitive replay reaches the four assertions, and the full aggregate passes through its final marker.

## Open Questions

- [x] Auto-approved: test-only mechanical hygiene (normal). → Added to queue.

---
*Source: discovered during REQ-187 canonical maintainer gate implementation.*
