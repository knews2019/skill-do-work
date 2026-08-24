---
id: UR-073
title: Reconcile the tdd-vs-write-set capture against what main shipped
created_at: 2026-08-24T23:37:06Z
requests: [REQ-372, REQ-373]
word_count: 46
---

# Reconcile the tdd-vs-write-set Capture Against What Main Shipped

## Summary

This branch captured two things while main was 147 commits ahead: a scope narrowing on REQ-365 and a
REQ-366 repairing the queued `tdd: true` REQs whose `write_set` named no test file. The merge made
both obsolete — main completed REQ-365 (`6265f1c`) and archived every REQ the repair targeted, and
the numbers REQ-366 and UR-070 belong to main's own work. Both were dropped in the merge commit.

Re-reading main's shipped text against the rest of the pipeline surfaced two findings that are still
live, captured here.

**F1 — the builder's out-of-scope response is now stated twice, differently.** The shipped
conditional-completeness invariant (`actions/capture-reference.md` → Populating `write_set`) tells a
builder meeting the contradiction to flag it, **proceed** with the required file class, and report in
the handback. `actions/work.md` § "Write only inside the declared scope" tells a builder discovering
a needed out-of-scope file to **stop-and-report to the orchestrator, never a silent write**, and puts
the extension of the Scope list and `write_set` on the orchestrator. A builder facing that moment now
reads two different actions, and the newer one lives in a capture-side file builders do not load.
→ REQ-372.

**F2 — the TDD-evidence gate accepted browser-probe evidence in place of a test.** REQ-353 completed
`tdd: true` with a one-file `write_set` and no test: its D-03 records "No new test file was added",
and its `## Testing` section carries generated-page RED/GREEN through Chromium instead. Step 6.5's
stated bar is test-first evidence — "write the failing test first, confirm it fails, then make it
pass" — so either the gate accepts more than it says, or that REQ should have been captured
`tdd: false`. Either way the shape the invariant exists to prevent was resolved by lowering the
evidence bar rather than by adding a test.
→ REQ-373.

## Reconciliation Decisions

Recorded during `do-work verify-requests`, because the input names them and no REQ carries them.

- **Version and changelog: deliberately not bumped.** The input asked to "update change log version".
  This branch changes nothing under `skills/` — the merge commit and this capture touch queue data
  only — and `main`'s own capture commits (`9edc122`, `dff8fb3`, `1f889e3`) bump neither `VERSION`
  nor `CHANGELOG.md`. User confirmed the no-bump reading on 2026-08-25. The shared version stays
  `0.236.61`.
- **UR/REQ renumbering: done in the merge commit.** `UR-070` and `REQ-366` belong to `main`; this
  capture took `UR-073`, `REQ-372` and `REQ-373`.
- **Merge fallout: verified, not assumed.** `bash _dev/tests/maintainer-verify.sh` exits 0 on the
  merged tree ("Maintainer verification passed"), with Go 1.26.1, ShellCheck 0.11.0 and `just`
  supplied to the container; the strict browser lane skips with no browser configured.

## Extracted Requests

| REQ | Title |
|---|---|
| REQ-372 | State the builder's out-of-scope response once |
| REQ-373 | Decide what the TDD-evidence gate accepts as test-first evidence |

## Batch Constraints

Both are instruction-consistency REQs over `skills/do-work/actions/`. Neither may turn `write_set`
into a gate, and neither ships code. They overlap on `actions/work.md`, so run them serially or let
the board's overlaps badge do its job.

## Full Verbatim Input

> do git pull/merge from the main branch and reconcile, i.e. but not limited to update change log version, update the UR/REQ numbers, fix all issues that the merge could cause, if the same thing is implemented twice, check which implementation is better and keep that one

---
*Captured: 2026-08-24T23:37:06Z*
