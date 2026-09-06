---
id: REQ-605
status: completed
domain: general
created_at: 2026-09-06T08:19:05Z
claimed_at: 2026-09-06T13:27:25Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
route: B
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-06T16:32:20Z
  basis:
    - Route B
    - 2-file write set
    - 2 acceptance criteria
related: [REQ-597]
write_set: [skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply_test.go]
title: 'List a merge commit''s paths when finalization checks the prepared-head range'
commit: 6dbcf4950b2211703c943158c9cc1240095f4253
completed_at: 2026-09-06T13:33:22Z
release_at: 2026-09-06T13:33:22Z
---

# List a Merge Commit's Paths When Finalization Checks the Prepared-Head Range

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Add `-m` flag to `git diff-tree` in `matchingHeadCommit` (`finalization_apply.go:545`) so merge commits list paths relative to each parent instead of emitting empty output. Add unit tests in `finalization_apply_test.go` verifying that a merge commit touching disallowed paths is rejected, while a clean merge touching only allowed paths is matched.)
- [x] **[APPLY]:** (Agent: Implemented `-m` flag in `finalization_apply.go` and added test cases `TestMatchingHeadCommitRejectsMergeWithDisallowedPaths` and `TestMatchingHeadCommitAcceptsCleanMergeMatchingEffectivePaths` in `finalization_apply_test.go`.)
- [x] **[UNIFY]:** (Agent: Ran `git diff --stat`, `gofmt -w`, native project tests, and all maintainer guard checks (`audit-lockins.sh`, `prescribed-shell-canonicalization.sh`, `action-shell-blocks.sh`, `quiet-grep-pipeline-audit.sh`, `gate.sh`).)

## Triage

**Route: B** — Mechanical core finalization fix with TDD unit tests.

**Reasoning:**
- Modifies Go implementation in `finalization_apply.go` (`diff-tree` invocation).
- Adds unit tests in `finalization_apply_test.go` locking in merge commit diff-tree path inspection.
- Low surface area, mechanical single-flag fix.

## Plan

1. **Write TDD RED test**: Add `TestMatchingHeadCommitRejectsMergeWithDisallowedPaths` in `finalization_apply_test.go` where a merge commit's binary diff over effective paths matches the prepared diff, but the merge also touches unrelated files outside `EffectiveCommitPaths`. Without `-m`, `diff-tree` outputs nothing and the test fails.
2. **Apply fix**: Add `-m` to `exec.Command("git", "-C", repositoryRoot, "diff-tree", "--no-commit-id", "--name-only", "-r", "-m", candidate)` in `finalization_apply.go`.
3. **Verify GREEN**: Ensure `TestMatchingHeadCommitRejectsMergeWithDisallowedPaths` and clean merge matching pass.
4. **Qualification & Verification**: Run all unit tests, repository guards, and the full maintainer gate.

## Exploration

Explored `matchingHeadCommit` in `internal/finalization/finalization_apply.go`:
- In `PreparedHead..HEAD`, `matchingHeadCommit` searches for an existing commit matching the prepared head state.
- For each candidate, it checks the binary diff over `EffectiveCommitPaths` against `PreparedDiffSHA256`.
- It then executes `git diff-tree --no-commit-id --name-only -r candidate` to verify exactness (no paths outside `allowed`).
- Standard `git diff-tree` produces no output for merge commits without `-m` (or `-c`/`--cc`), causing `strings.Fields(string(changed))` to be empty and erroneously passing the exactness check even if the merge touched other paths.
- Adding `-m` forces `diff-tree` to report all changed paths relative to each parent, ensuring foreign paths are detected.

## Scope

Files in scope for this change:
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go`
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply_test.go`

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

## Implementation Summary

- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply_test.go` (created)

**Passed `-m` to `git diff-tree` in `matchingHeadCommit` and verified merge candidate path validation:**
- `finalization_apply.go`: Added `"-m"` argument to the `git diff-tree --no-commit-id --name-only -r` invocation, ensuring changed paths in merge commits are enumerated across all parents.
- `finalization_apply_test.go`: Added `TestMatchingHeadCommitRejectsMergeWithDisallowedPaths` confirming that a merge commit containing foreign changes outside `EffectiveCommitPaths` is rejected, and `TestMatchingHeadCommitAcceptsCleanMergeMatchingEffectivePaths` confirming that a clean merge containing exclusively permitted paths is accepted.

## Decisions

- **D1 Use `-m` rather than `--first-parent`:** `git diff-tree` does not accept `--first-parent` (which is a log/rev-list option); `-m` is the canonical `diff-tree` flag to show separate diffs with respect to each parent commit, accurately capturing all files altered across any merged branch.

## Qualification

**Passed.** Read from range `21f7e68fc56565ac8d469a4cbb759614b0677256..6dbcf4950b2211703c943158c9cc1240095f4253`, 2 files, 93 insertions, 1 deletion.
Canonical `qualify` and `scope-drift` satisfied.

- `skills/do-work/tools/do-work-cli` unit tests passed (`go test ./...` in 52s).
- Repository guards `audit-lockins.sh`, `prescribed-shell-canonicalization.sh`, `action-shell-blocks.sh`, and `quiet-grep-pipeline-audit.sh` all passed with exit 0.
- Full maintainer gate `gate.sh` passed with exit 0 (101s wall time, 803 tests).

## Testing

**Commands executed:**
- `go test -C skills/do-work/tools/do-work-cli -v -run "TestMatchingHeadCommit" ./internal/finalization/...` — passed, exit 0.
- `go test -C skills/do-work/tools/do-work-cli ./...` — all packages passed, exit 0.
- `bash _dev/tests/prescribed-shell-canonicalization.sh` — passed, exit 0.
- `bash _dev/tests/audit-lockins.sh` — passed, exit 0.
- `bash _dev/tests/action-shell-blocks.sh` — passed, exit 0.
- `bash _dev/tests/quiet-grep-pipeline-audit.sh` — passed, exit 0.
- `DO_WORK_GATE_ROOT="$(pwd)" bash do-work/runs/work-2026-09-05-231943/handoff-tools/gate.sh` — `Maintainer verification passed.`, exit 0 (803 tests).

## Review

**Overall: 99%** | 2026-09-06T16:32:00Z | Synthesis of review lenses (code correctness, test coverage, safety invariants, maintainer gate)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict: Pass.** `diff-tree -m` ensures merge commits in `PreparedHead..HEAD` cannot silently bypass exactness validation, with full test coverage asserting both rejection of foreign paths and acceptance of clean merges.

## Remediation

None needed.

## Lessons Learned

**What worked:**
- Verifying the failure with a merge commit that resolves into matching target diff bytes while carrying unrelated files proved the exact defect mechanism.
- Git's `diff-tree -m` cleanly inspects each parent diff without requiring extra plumbing.

**What didn't:**
- Assuming `git diff-tree` will output changes for all commit types by default; merge commits require explicit flags like `-m`.

**Worth knowing:**
- In Git rev-list traversals, `--ancestry-path A..B` restricts to commits that are strictly descendants of A, which can prune non-descendant parent branches when evaluating merges.

## Orientation

Adds `-m` to `git diff-tree` in `matchingHeadCommit` (`skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go`) so candidate merge commits in `PreparedHead..HEAD` have their changed paths inspected against `EffectiveCommitPaths`. Adds regression tests in `finalization_apply_test.go`.
