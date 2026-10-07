# REQ-638 Hand-back (release guard reads suite/modules.tsv once)

- Branch: `worktree-agent-REQ-638-release-guard-reads-modules-tsv-once` (no name collision)
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-638-release-guard-reads-modules-tsv-once`
- Base: `13932fc2`. Commit: `4650565a` ([REQ-638] release guard reads suite/modules.tsv once)

## File manifest
- Modified: `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go` (+8 / -13). Adds the `slices` import, removes the second `DeclaredModuleSources` read, the second `headReleaseImage` build, and the set merge.
- New: none.

## Red-green evidence
`tdd: false`. Per the REQ's Red-Green Proof, no runnable RED exists because behaviour is unchanged. The proof is that the existing suite is green before and after the change.
- Before: `go test -count=1 ./internal/finalization/` → `ok ... 54.194s`
- After: `go test -count=1 ./...` → every package `ok` (finalization 71.357s under full-module parallel load). The `-json` rerun of finalization showed 0 fail events.
- Relevant existing tests: `finalization_release_guard_test.go` (5.95s, green).

## P-A-U
- **[PLAN]:** Read `DeclaredModuleSources` once, through a single `headReleaseImage`. The consumer-project exit becomes "no declared source has a tracked `<source>/VERSION`". This is exactly `len(DeclaredMaintainerReleaseRoots(...)) == 0`, because that function filters the same sources on the same condition. Roots = sources plus `suite` and `tools` when they are absent, then sorted. This is byte-identical: `DeclaredModuleSources` already returns sorted, unique entries, the old maintainer roots were a subset of those sources, and the old set merge only removed duplicates, which the absent-check reproduces. `DeclaredMaintainerReleaseRoots` and `DeclaredModuleSources` stay unchanged.
- **[APPLY]:** One file was edited as planned. The tracked-VERSION check is a `slices.ContainsFunc` one-liner, written inline.
- **[UNIFY]:** `git diff --stat` shows 1 file changed (8+/13-). `gofmt -l ./internal/finalization/` printed nothing. `go vet ./internal/finalization/ ./internal/releaseownership/` exit 0. The full module `go test -count=1 ./...` is all ok. The diff has no debug artifacts. `sort` is still used (the roots sort). `DeclaredMaintainerReleaseRoots` still has another caller inside releaseownership, so it is not orphaned.

## Decisions
- D-01 DECIDE & STATE: The tracked-VERSION check stays inline as `slices.ContainsFunc` and is not a named helper, because it is used in one place only (this was the Builder Guidance latitude).
- D-02 DECIDE & STATE: `suite` and `tools` are added only when absent (`slices.Contains`), not through a set. This keeps the old dedup behaviour for a declared source literally named `suite` or `tools`, so the output stays identical for every input.
- D-03 DECIDE & STATE: The REQ write_set lists `skills/do-work/CHANGELOG.md`, but the brief forbids touching it. I followed the brief, so the orchestrator writes the release and CHANGELOG.

## Discovered Tasks
- None.

## Lessons read
- `_dev/primes/lessons-releases.md` (whole file, required_lessons)
- Primes: `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, `_dev/primes/prime-releases.md`
- Crew: general.md, coding-guardrails.md, shared-principles.md, communication-style.md, backend.md
- No "Dropped for Budget" satellites are listed in the REQ.

## Proposed lesson
None warranted, because this was a behaviour-preserving cosmetic deletion with no incident. If the orchestrator wants one anyway:
- `[family: subset-merged-into-superset]` (lessons-releases.md) — REQ-638: the release guard read `DeclaredMaintainerReleaseRoots` and then `DeclaredModuleSources`, and merged them, but the first function is a filter of the second. When one reader is derived from another, call the base reader once and apply the filter as a predicate, so the declaration is parsed once.

## Proposed CHANGELOG entry
Title: **Release Guard Reads the Module Declaration Once**
Body: The shipped-change release guard now parses `suite/modules.tsv` once instead of twice and no longer merges a subset of module roots into their superset. Behaviour is unchanged.
(Patch bump. Shipped Go change under `skills/`.)

## Integration seams
None. Only `releaseShippedChangeError` changed, and its signature and outputs are the same.

## Test commands and results
- `gofmt -l ./internal/finalization/` → empty
- `go vet ./internal/finalization/ ./internal/releaseownership/` → exit 0
- `go test -count=1 ./...` (module `skills/do-work/tools/do-work-cli`) → all 32 test packages ok, wall 1:13.6. `internal/releaseownership` has no test files.
- Per-file serial elapsed sum for `internal/finalization` (budget under 30s): req499 17.33s, recovery 16.27s, release_guard 5.95s, review_regressions 5.58s, req547 3.35s, commands 1.26s, req560 1.00s, pipeline_dirt 0.79s, review_release 0.68s, apply 0.42s, req512 0.35s, req557 0.08s, req565 0.00s. Every file is under budget. The package as a whole takes 54-71s.
