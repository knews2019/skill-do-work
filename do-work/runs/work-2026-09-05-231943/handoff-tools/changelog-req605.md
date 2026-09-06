## 0.305.16 — List Merge Commit Paths with Diff-Tree -m in Finalization Range Check (2026-09-06)

`matchingHeadCommit` in the finalization subsystem now passes `-m` to `git diff-tree` when verifying candidate commits in the `PreparedHead..HEAD` range, preventing merge commits from emitting empty path listings and bypassing exactness path constraints.

- `internal/finalization/finalization_apply.go` invokes `git diff-tree` with `-m` so candidate merge commits enumerate paths modified across each parent against `EffectiveCommitPaths`.
- `internal/finalization/finalization_apply_test.go` verifies rejection of merge commits containing foreign paths and acceptance of clean merges.
