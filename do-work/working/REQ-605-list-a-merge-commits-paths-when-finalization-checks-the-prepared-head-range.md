---
id: REQ-605
status: claimed
domain: general
created_at: 2026-09-06T08:19:05Z
claimed_at: 2026-09-06T13:27:25Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
related: [REQ-597]
write_set: [skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply_test.go]
title: 'List a merge commit''s paths when finalization checks the prepared-head range'
---

# List a Merge Commit's Paths When Finalization Checks the Prepared-Head Range

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

`internal/finalization/finalization_apply.go:545` runs `git diff-tree` without `-m` over the
candidates in `PreparedHead..HEAD`. For a merge commit `diff-tree` without `-m` lists no paths, so the
`exact` loop stays true for it. Today only the preceding `git diff --binary` digest match at
`:541-543` keeps that branch unreachable. Read by REQ-597's guide builder while checking the guide's
merge-aware commit diff sentence; not a guide claim, a latent code hazard.

## Why

A check that is correct only because an earlier check happens to run first is one refactor from wrong.
The fix is one flag and one test that puts a merge commit in the range.

## Detailed Requirements

- `diff-tree` lists a merge commit's paths (`-m`, or `--first-parent` if the record argues the
  first-parent diff is what the exactness check means; say which and why).
- A test with a merge commit among the candidates that is red on the current code with the digest
  pre-check bypassed or with a merge whose paths differ from the prepared set, and green after.

## Constraints

- Shipped Go: a release. Change only the argv and the test.

## Open Questions

None.
