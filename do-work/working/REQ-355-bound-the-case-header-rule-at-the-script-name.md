---
id: REQ-355
title: "[impact-negligible] Review fix: bound the case-header rule at the script name"
status: claimed
claimed_at: 2026-08-24T17:42:55Z
status_changed_at: 2026-08-24T17:42:55Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
