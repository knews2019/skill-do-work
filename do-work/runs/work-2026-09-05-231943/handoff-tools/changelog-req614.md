## 0.305.25 — Batch Repeated Git Reads in Finalization and Request State (2026-09-06)

Profiled recovery and state planning to batch repeated Git reads and memoize historical commits during finalization discovery.

- Batched `existingUntrackedPaths` queries in `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` using a single `git --literal-pathspecs ls-files -z -- <paths...>` command, eliminating over 1,180 per-path subprocesses.
- Batched `existingDirtyTrackedPaths` queries in `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` using a single `git status --porcelain=v1 -z --untracked-files=no -- <paths...>` command, eliminating over 660 per-path subprocesses.
- Introduced `discoverySession` in `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go` to memoize immutable HEAD file images and tracked release paths across discovery phases, keyed to the operation entry head commit to guarantee freshness across write operations.
- Reduced redundant `git show` subprocess calls from 645 to 274 (-57.5%) and overall Git subprocesses in `TestRecoverFinalization` from 2,950 to 2,446, reducing recovery wall time to 18.13s.
- Added unit tests verifying batched path classification, session memoization, negative lookups for absent files, and cache invalidation across commits.
