---
id: UR-070
title: Fix the tdd-vs-write-set contradiction at capture and in the queue
created_at: 2026-08-24T23:17:16Z
requests: [REQ-366]
word_count: 118
---

# Fix the tdd-vs-write-set Contradiction at Capture and in the Queue

## Summary

The user asked whether a summary of REQ-365 was accurate: `tdd: true` REQs whose `write_set` names
no test file, and the inference that GREEN was therefore being claimed on behavioural assertions
rather than code-level tests. The first half checks out (REQ-346 carried the shape); the second half
does not — `actions/work.md` Step 6.5 gates `tdd: true` on red-green evidence, and REQ-346's builder
wrote `generate_test.go` outside its declared set rather than skipping the test.

The user then asked what to fix, and endorsed the resulting recommendation by asking for it to be
captured. Three parts came out of the analysis:

- **A1** — state the rule once in `actions/capture-reference.md` § Populating `write_set`, keyed on
  the condition (a declared write set must include every class of file the REQ's own GREEN requires
  writing; tests under `tdd: true` are today's instance). Already owned by queued REQ-365.
- **A2** — drop `skills/do-work/actions/capture.md` from REQ-365's write set, and cite
  `actions/work.md`'s "Write only inside the declared scope" bullet instead of restating builder
  behaviour that already exists there. A narrowing of REQ-365 → folded into it as an addendum.
- **A3** — repair the queued REQs that already carry the shape. Five are in `do-work/queue/`
  (REQ-348, REQ-350, REQ-352, REQ-353, REQ-362); REQ-347 also carries it but is claimed and
  immutable, so it is out of scope and its builder handles it through the existing stop-and-report.
  New work → REQ-366.

A fourth idea — a lint over queued REQs' `write_set` — was considered and rejected during the
analysis: it needs a heuristic for "is this a test path", and A3 empties the backlog it would police.

## Extracted Requests

| Request | Destination |
|---|---|
| State the write-set-must-include-required-file-class rule at capture | REQ-365 (already queued) |
| Narrow REQ-365: one file, cite the existing builder rule rather than restating it | Folded into REQ-365 as an addendum |
| Repair the five queued `tdd: true` REQs whose write set names no test file | REQ-366 |

## Folded Requests

- REQ-365 (a-tdd-req-must-name-a-test-file-in-its-write-set) — the scope narrowing: state the rule once in `capture-reference.md` only, drop `capture.md` from the write set, and cite `actions/work.md`'s declared-scope bullet instead of restating what a builder does when it meets the contradiction

## Full Verbatim Input

> can this happen?
>
> """
> **TDD vs. Scope Declaration**
> REQ-365 explicitly identifies a contradiction in the project's workflow: some requirements are marked as `tdd: true` (Test Driven Development), but their declared `write_set` contains no test files. This implies that "Green" status was previously being claimed based on behavioral assertions rather than actual code-level tests.
> """

> what should we fix (first tell me, no implementation just yet)

> use do-work capture-request to capture the intent so it can be implemented

---
*Captured: 2026-08-24T23:17:16Z*
