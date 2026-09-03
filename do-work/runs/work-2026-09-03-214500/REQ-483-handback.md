# REQ-483 builder handback

REQ-483 — Bound the architecture bundle-claim loop and restore `--commit`

## Identity

- Branch: `worktree-agent-REQ-483-bound-architecture-bundle-claim-loop`
- Base: `c27d349a`
- Commit: `f2ae022a49189acec1f9c934d0c85e83c42a52f1`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-483-bound-architecture-bundle-claim-loop`

## Outcome

The publication preflight now validates the candidate transaction without forwarding the caller's commit bit into the dry-run probe, so `--commit` reaches the real mutation and commits the generated index. Bundle claiming retries only `fs.ErrExist`; permission, read-only, and other non-collision errors stop immediately with `ARCHITECTURE-BUNDLE-CLAIM-FAILED`, the configured candidate path, and the underlying cause.

## RED → GREEN evidence

- RED: `TestArchitecturePublishStopsOnNonCollisionClaimFailure` timed out after 500 ms because the old claim loop retried a permission error forever.
- RED: `TestArchitecturePublishCommitCommitsPublishedBundle` returned `GIT-INVALID-OPTIONS` with `--dry-run and --commit cannot be combined`.
- GREEN: both focused regressions passed after the control-flow changes.
- GREEN: `go test ./internal/toolboxcommands`, `go vet ./...`, and `go test ./...` passed for the full do-work-cli module.
- Hygiene: `gofmt`, `git diff --check`, and complete two-file diff review passed.

## Manifest

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture_test.go`

## Decisions and risks

- Added one package-local claim-function seam so the non-collision error is deterministic under privileged test environments; production still delegates directly to the rooted exclusive mkdir helper.
- The typed failure reports only the repository-relative candidate path and OS cause; it never retries an error that does not prove a namespace collision.
- The initial candidate scan remains unchanged; this request governs the exclusive claim loop and commit passthrough only.

## Merge guidance

Merge commit `f2ae022a49189acec1f9c934d0c85e83c42a52f1` from the named branch. Qualify exactly the two declared files, rerun the focused regressions and full CLI module after integration, then apply the canonical gate/heavy-surface policy.

## Discovered tasks

None.
