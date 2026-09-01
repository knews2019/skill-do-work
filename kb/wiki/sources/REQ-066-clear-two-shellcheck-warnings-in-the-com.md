---
title: "Lessons from REQ-066: Clear two shellcheck warnings in the commit-hash guard fixture"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-066-clear-two-shellcheck-warnings-in-the-com.md]
related:
  - page: concept-prescribed-shell-commands
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-066: Clear two shellcheck warnings in the commit-hash guard fixture

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

`_dev/tests/record-commit-hash-guards.sh` trips two shellcheck warnings at default severity. Both
predate REQ-064 (confirmed against `HEAD`) and both sit in the scan-probe block REQ-063 added:

- **SC2120** at the `run_scan_script()` definition — the function references `"$@"` but no call site
  ever passes arguments.
- **SC2034** at `intact_bytes="$(wc -c < "$scan_target" | tr -d '[:space:]')"` — assigned, never read.

## Solution summary

`run_scan_script` forwarded arguments that neither of its two call sites ever passed; the forwarding is gone and a comment records that the bare detector is the point (the restore probes use `run_restore_script`, which genuinely does take arguments). `intact_bytes` was computed from the pre-blanking file and never read — it is now asserted against the scanner's `Recoverable: N bytes` line, which is a real added assertion rather than a silenced warning.

## What worked

- Treating an "unused variable" warning as a coverage question rather than a lint chore — `intact_bytes` existed because someone meant to assert the recoverable size, and the fix that silences the warning by adding the missing assertion is strictly better than the one that deletes the variable.

## What didn't work

- The first mutation test ran a mutated *copy* of the fixture from the scratchpad. The fixture resolves `repo_root` from its own location, so the copy failed with `FAIL: tools/checks/record-commit-hash.sh must exist and be executable` — a wrong-reason failure that would read as a passing mutation test if taken at face value. Mutate in place, back up first, restore after.

## Worth knowing

- `assert_output_matches` feeds its first argument to `grep -Eq`, so an expected value containing regex metacharacters needs escaping. Byte counts are safe; paths and hashes with dots are not.

## Back-reference

See `do-work/archive/UR-010/REQ-066-shellcheck-warnings-in-guard-fixture.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `b0bd8c8`.
