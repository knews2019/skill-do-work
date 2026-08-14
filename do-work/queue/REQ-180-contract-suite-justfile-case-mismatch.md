---
id: REQ-180
title: Fix contract-regressions.sh Justfile case mismatch aborting late checks
status: pending-answers
created_at: 2026-08-14T10:33:56Z
user_request: UR-040
addendum_to: REQ-179
domain: general
review_generated: true
tdd: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
write_set: [_dev/tests/contract-regressions.sh]
---

# Fix contract-regressions.sh Justfile Case Mismatch Aborting Late Checks

## What

`_dev/tests/contract-regressions.sh:1797` and `:1804` reference `Justfile` (capital J), but the tracked file is lowercase `justfile`. On a case-sensitive filesystem, `extract_kanban_shutdown_line Justfile` hits `awk: cannot open .../Justfile` and the suite aborts with exit 2 at that point — every check after line ~1797 (roughly 1,500 lines including the Common Rationalizations regrowth ratchet) silently never runs. On case-insensitive filesystems (macOS default) the mismatch is invisible, which is why it survived.

## Open Questions

- [ ] Fix by lowercasing the two references to `justfile`, or by resolving the file case-insensitively?
  Recommended: lowercase the two literals to match the tracked filename — smallest fix, matches the repo's actual file.
  Also: a glob/`find -iname` resolution (tolerant but adds machinery), rename the tracked file to `Justfile` (breaks `just` conventions).

## Red-Green Proof
**RED prompt/case:** On a case-sensitive filesystem, `bash _dev/tests/contract-regressions.sh` prints `awk: cannot open ".../Justfile"` and exits 2 at `extract_kanban_shutdown_line`, before the late-suite checks run.
**Why RED now:** Two hardcoded `Justfile` literals vs the tracked lowercase `justfile`; observed at baseline during REQ-179 (2026-08-14, Linux sandbox).
**GREEN when:** The suite runs past line ~1797 on a case-sensitive filesystem (the kanban-shutdown-line check actually executes against `justfile`), and the only remaining sandbox failure is the known environmental process-tree probe.
**Validation:** Discovered during REQ-179's build; classification [normal] → pending-answers per the discovered-tasks consent flow.

## Full Context

Discovered Tasks of `REQ-179` (see its archive entry). Parent UR: `do-work/user-requests/UR-040/` (archived with UR-040).

---
*Source: REQ-179 build — discovered task*
