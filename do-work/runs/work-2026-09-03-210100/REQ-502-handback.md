# REQ-502 hand-back

Branch: `worktree-agent-REQ-502-remove-enriched-checkpoint-entries-in-cleanup-mover`

Commit: `1f03f047f8ede51d8dd38c20ba1b3be797a018c8`

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan.go` (modified) — cleanup now delegates checkpoint removal to the canonical request-state helper.
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_plan_test.go` (modified) — the cleanup seam now covers enriched own and foreign entries plus an inline heading mention.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified) — exposes the exact-section, whole-entry writer-labelled removal primitive to alternate lifecycle writers.

**What was done:** The cleanup mover no longer filters matching header lines globally. It uses the request-state implementation that targets the real In Progress section, removes the owned header and indented continuation bytes together, and preserves foreign and out-of-section bytes.

## P-A-U evidence

- **PLAN:** Reproduce at the cleanup seam, then expose the existing canonical request-state removal primitive instead of copying its parsing rules.
- **APPLY:** Added the focused regression first, observed the orphaned own continuation lines, then replaced cleanup's local line filter with the shared primitive.
- **UNIFY:** Reviewed all three changed files, ran gofmt, focused and package tests, `go vet`, `git diff --check`, and a debug-artifact scan. The builder worktree is clean after commit.

## Red-Green evidence

- RED: `go test -count=1 ./internal/cleanup -run '^TestWorkingArchiveRemovesOnlyThisCheckoutCheckpointEntry$'` failed because `Last known state: implementing` and `cleanup_plan.go` remained after the own claim header was removed.
- GREEN: the same focused command passed after the fix. `go test -count=1 ./internal/cleanup ./internal/requeststate` and `go vet ./internal/cleanup ./internal/requeststate` also passed.

## Lesson evidence

- Read `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` and its whole lesson satellite.
- Applied `alternate-writer-contract-drift`: the alternate cleanup writer now calls the canonical whole-entry implementation rather than keeping a second definition.
- Read the original REQ-489 archive and the bug-fix spec.

## Integration seams

None.

## Decisions

- **D-01 — DECIDE & STATE:** Export the narrow writer-labelled removal primitive from `requeststate` and let `cleanup` import it. This keeps the stored-format rule in its existing owner and avoids a new package for one reuse site.

## Discovered Tasks

None.
