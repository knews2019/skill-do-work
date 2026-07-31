---
id: REQ-066
title: Clear two shellcheck warnings in the commit-hash guard fixture
status: completed
claimed_at: 2026-07-31T09:32:00Z
route: A
completed_at: 2026-07-31T09:34:31Z
commit: b0bd8c8
kb_status: pending
created_at: 2026-07-31T09:28:23Z
user_request: UR-010
addendum_to: REQ-064
domain: general
discovered_during: REQ-064
prime_files: []
tdd: false
maintenance: false
write_set: [_dev/tests/record-commit-hash-guards.sh, CHANGELOG.md, actions/version.md]
---

# Clear two shellcheck warnings in the commit-hash guard fixture

## What

`_dev/tests/record-commit-hash-guards.sh` trips two shellcheck warnings at default severity. Both
predate REQ-064 (confirmed against `HEAD`) and both sit in the scan-probe block REQ-063 added:

- **SC2120** at the `run_scan_script()` definition — the function references `"$@"` but no call site
  ever passes arguments.
- **SC2034** at `intact_bytes="$(wc -c < "$scan_target" | tr -d '[:space:]')"` — assigned, never read.

## Why

`shellcheck tools/checks/*.sh` is clean, but the test fixture that guards those scripts is not, so
anyone widening the lint to `_dev/tests/` gets noise instead of a signal. Neither warning affects
behavior — the suite passes and every assertion is meaningful.

## Detailed Requirements

Resolve each warning at its cause rather than with a blanket `# shellcheck disable` at the top of
the file:

- `run_scan_script` — either drop the unused argument forwarding, or add the `# shellcheck
  disable=SC2120` directive at the function if a later probe is expected to pass arguments (the
  `--restore` probes added by REQ-064 use their own `run_restore_script`, which *does* take args).
- `intact_bytes` — either assert on it (the intent was plainly to check the healthy file's size
  is unchanged, which would be a real added assertion) or remove it.

`bash _dev/tests/contract-regressions.sh` must stay green, and no existing assertion may change
meaning.

## Open Questions

- [x] Auto-approved: test-only mechanical hygiene (low). → Added to queue.

## Constraints

- Test file only — zero production-source changes.
- Do not add a file-level `shellcheck disable`; that hides future warnings too.

---
*Discovered during REQ-064 (cleanup Pass 6). Auto-queued under the test-hygiene carve-out.*

## Triage

**Route: A** (Direct to Builder)

Names one specific file, two specific warnings with line-level locations, and states the acceptable
fixes. Nothing to explore. Planning not required.

## Implementation Summary

**Files changed:**

- `_dev/tests/record-commit-hash-guards.sh` (modified) — dropped the unused `"$@"` forwarding from
  `run_scan_script` (SC2120) and turned the dead `intact_bytes` into a live assertion (SC2034).
- `CHANGELOG.md` (modified) — 0.153.1 entry.
- `actions/version.md` (modified) — 0.153.0 → 0.153.1.

**What was done:** `run_scan_script` forwarded arguments that neither of its two call sites ever
passed; the forwarding is gone and a comment records that the bare detector is the point (the
restore probes use `run_restore_script`, which genuinely does take arguments). `intact_bytes` was
computed from the pre-blanking file and never read — it is now asserted against the scanner's
`Recoverable: N bytes` line, which is a real added assertion rather than a silenced warning.

## Decisions

**D-01 — Assert `intact_bytes` rather than delete it.** DECIDE & STATE. The REQ offered both. The
recoverable byte count is the number an operator actually decides on when triaging a blanked file,
and it was the one thing in the scan output no probe covered — so the variable's existence was
pointing at a real coverage gap, not at dead code. Deleting it would have cleared the warning and
kept the gap.

**D-02 — `CHANGELOG.md` and `actions/version.md` added to `write_set`.** DECIDE & STATE. This REQ
is Route A, so it never ran Step 5.5 and its `write_set` was the capture-seeded single test file.
The repo requires a version bump and a changelog entry on every commit, so those two files were
always going to be touched; the field was extended before the write rather than after. No other
REQ was in flight.

## Qualification

Passed — 3 files verified in `git diff`, both requirements traced, no P-A-U section on this
follow-up REQ to audit.

- Both warnings traced to a specific change; `shellcheck _dev/tests/record-commit-hash-guards.sh`
  is clean at default severity, and no file-level `disable` directive was added (the REQ's explicit
  constraint — verified by grep: zero `shellcheck disable` lines in the file).
- The new assertion is substantive, not decorative: mutating it to expect the on-disk `0 bytes`
  instead of the pre-blanking size makes the suite fail with exactly that probe named, and
  reverting restores green.
- No existing assertion changed meaning — the diff adds one and removes none.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`, `bash
_dev/tests/record-commit-hash-guards.sh`, `shellcheck _dev/tests/record-commit-hash-guards.sh`
**Result:** ✓ All passing — 58 assertions (was 57), contract regressions clean, shellcheck clean at
info+.

Non-behavioral change to production code (test-file only), so red-green does not apply to the fix
itself. Evidence is regression plus the mutation test above: the added assertion fails when the
expected value is wrong and passes when it is right.

## Review

**Overall: 95%** | 2026-07-31T09:35:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 0 minor

**Restatement sweep (MUST — performed):** this diff redefines nothing that other text restates. It
changes one test helper's signature (two call sites, both in the same file, both updated by
construction since the helper now takes no arguments) and adds one assertion. No contract token,
schema field, gate wording, or prescribed command output shape is touched. `assert_output_matches`
and the `Recoverable: N bytes` output line are both consumed only here.

**Acceptance:** Pass — both warnings resolved at their cause, no blanket disable, suite green, and
the SC2034 fix landed as an added assertion rather than a deletion.
**Suggested testing:** None.
**Follow-ups created:** None.

*Reviewed by review-work action (pipeline mode)*

## Lessons Learned

**What worked:** Treating an "unused variable" warning as a coverage question rather than a lint
chore — `intact_bytes` existed because someone meant to assert the recoverable size, and the fix
that silences the warning by adding the missing assertion is strictly better than the one that
deletes the variable.

**What didn't:** The first mutation test ran a mutated *copy* of the fixture from the scratchpad.
The fixture resolves `repo_root` from its own location, so the copy failed with `FAIL:
tools/checks/record-commit-hash.sh must exist and be executable` — a wrong-reason failure that
would read as a passing mutation test if taken at face value. Mutate in place, back up first,
restore after.

**Worth knowing:** `assert_output_matches` feeds its first argument to `grep -Eq`, so an expected
value containing regex metacharacters needs escaping. Byte counts are safe; paths and hashes with
dots are not.

## Orientation

The commit-hash guard fixture now lints clean and asserts one more thing than it did: that the
scanner reports the *pre-blanking* recoverable byte count, not the zero on disk. Leaf change inside
the existing test suite — no new surface, no contract moved. No `prime_files` on this REQ and no
prime covers `_dev/tests/`.
