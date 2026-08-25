---
id: REQ-355
title: "[impact-negligible] Review fix: bound the case-header rule at the script name"
status: completed
claimed_at: 2026-08-24T17:42:55Z
completed_at: 2026-08-24T18:12:09Z
commit: 059887b02050b8a5fa1130193aaeb3f71ba240b3
status_changed_at: 2026-08-24T18:12:09Z
route: A
created_at: 2026-08-24T09:20:00Z
user_request: UR-065
addendum_to: REQ-339
domain: testing
review_generated: true
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-24T17:42:55Z
  basis:
    - trivial short-circuit
write_set:
  - _dev/tests/prescribed-shell-case-count.sh
  - _dev/tests/contract-regressions.sh
---

# Bound the Case-Header Rule at the Script Name

## What

`count_named_case_headers` states one rule in prose and implements a wider one. The comment says a
case header "opens with the name of the script the case file covers"; the regex
`^# ${script_under_test}[^.:]*: ` has no boundary after the interpolated basename, so it accepts any
token *beginning* with that name. Add the boundary, and move the comment in the same edit.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route A applies the specified basename-boundary correction and fixture directly in the captured two-file surface.
- [x] **[APPLY]:** Added the `qualifying ...` near-miss and tightened the shared matcher without enumerating qualifier separators.
- [x] **[UNIFY]:** Reviewed both changed shell files; RED/GREEN fixture, all 17 live inventories, Bash syntax, ShellCheck, contract regressions, diff checks, and canonical maintainer verification passed.

## Why

This is the over-count mirror of the undercount REQ-339 removed, in the file REQ-339 created to be
the single home of the rule — so the same defect shape the REQ set out to close is present in its
own fix, one direction over. Nothing catches it: REQ-339's lock-in fixture carries three near-misses
covering the narrowing direction and two broadenings, but none of this shape.

`impact-negligible` because no line in the current corpus exercises it — a scan of all 17 case files
for a basename followed by a word character returns zero hits. It is worth fixing anyway because the
whole point of the file is that the stated rule and the applied rule cannot disagree.

## Detailed Requirements

- The regex requires a boundary after the interpolated basename, so a token merely *starting* with
  the script name is not a header.
- The comment's stated rule and the regex agree after the change. Both halves move together — a fix
  that corrects the regex and leaves the prose describing the old behaviour recreates the defect.
- REQ-339's lock-in fixture gains the missing near-miss so this direction is pinned like the other
  three. The reviewer's reproduction is the case to use:
  `# qualifying the untracked scan below is prose: it must not count.` in a fixture named
  `qualify.sh`, which counts 2 today against a 1-case fixture.
- The counts for all 17 real case files are unchanged by this edit (currently 101 in aggregate) —
  the corpus does not exercise the over-count, so a change in any per-file number means the fix went
  too far.

## Constraints

- `_dev/primes/prime-shell-commands.md` § *Closed Enumerations Go Stale* still governs: the
  qualifier stays open-ended. This narrows only where the qualifier may *start*, not what it may
  contain.
- Do not reintroduce the enumeration of separator spellings REQ-339 deliberately removed.

## Red-Green Proof

**RED prompt/case:** In a fixture named `qualify.sh` containing exactly one real header
(`# qualify: a real case.`) plus the line `# qualifying the untracked scan below is prose: it must
not count.`, `count_named_case_headers` returns 2. Reproduced 2026-08-24.

**GREEN when:** that fixture returns 1, the lock-in in `contract-regressions.sh` fails if the
boundary is removed again, the comment states the bounded rule, and all 17 real case files report
the same counts as before (aggregate 101).

**Validation:** Inferred during REQ-339's review — a review finding, not a user request.

---
*Source: REQ-339 review finding F1 (UR-065).*

## Triage

**Route: A** — The exact two files, failing fixture line, regex boundary, unchanged corpus count, and RED/GREEN command are specified; implement the bounded basename rule directly.

## Plan

**Planning and exploration skipped** — Route A: implement the specified regex/comment/fixture correction directly and verify every real case-file count remains unchanged.

## Decisions

- **D-01 — Preserve the live 104-case corpus, not the captured 101 count.** Three legitimate `qualify.sh` cases landed after this review finding was written, so the claimed worktree begins at 104 cases across 17 files (`qualify.sh`: 21). RED is expected 3 versus observed 4 in the isolated near-miss fixture; GREEN must keep every live per-file count and the aggregate 104 unchanged.
- **D-02 — Bound the script token without closing the qualifier language.** `[[:alnum:]_]` is the continuation class; any other first qualifier character except `.` or `:` remains allowed, so future separator spellings stay open-ended.

## Implementation Summary

- `_dev/tests/prescribed-shell-case-count.sh` (modified): requires a non-word boundary after the interpolated script basename while leaving the following period-free qualifier open-ended, with the rationale adjacent to the matcher.
- `_dev/tests/contract-regressions.sh` (modified): adds the `qualifying ...` word-prefix near-miss and updates the fixture inventory so removal of the boundary fails.

## Discovered Tasks

None.

## Testing

- RED fixture contained three valid `qualify.sh` headers plus `# qualifying a request: ...`; the original matcher returned 4 instead of 3.
- GREEN returned 3, while a complete before/after inventory kept every one of 17 real case files unchanged at aggregate 104 (`qualify.sh`: 21).
- Bash syntax, ShellCheck, focused contract regressions, diff checks, and builder canonical maintainer verification passed.
- Post-merge contract regressions passed with 104 named cases and zero failures; the full canonical maintainer gate passed, including formatting, Go vet/tests, strict JavaScript, and audit-metrics tests. Its browser lane correctly skipped because this shell-only REQ has no browser requirement.

## Qualification

- Exact merge range `cd13533..059887b` passed mechanical qualification.
- Orchestrator judgment confirmed a substantive bounded-token fix, exact two-file Route A scope, active matcher-to-fixture flow, open-ended qualifier semantics, and no generated/debug artifacts.

## Review

Independent review approved with no Important or Minor findings and one Nit: the fixture uses equivalent `qualifying ...` wording rather than the review's exact quoted sentence. Requirements 98%, Code Quality 98%, Test Adequacy 100%, Scope Discipline 100%, overall 99%, low risk, acceptance pass. Word/digit/underscore continuations fail while bare, space-, comma-, and hyphen-qualified headers remain accepted.

## Lessons Learned

When a matcher interpolates a token followed by an open-ended qualifier, define the token boundary independently from the qualifier language; that closes prefix over-counting without recreating a separator enumeration.

## Orientation

Released in 0.236.48. Named shell-case headers now require a real boundary after the script basename while preserving every live qualifier form and all 104 current cases.
