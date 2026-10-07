# REQ-638 Review: Release Guard Reads suite/modules.tsv Once

**Overall: 96%**

| Requirements | Code Quality | Test Adequacy | Scope | Risk | Acceptance |
|---|---|---|---|---|---|
| 100% | 97% | 88% | 100% | Low | Pass |

Reviewed range: da3ab5f5..6125e58c, one file, `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go` (+8/-13).

## Equivalence check (byte-identical behaviour)

- **Consumer-project exit.** Old: `len(DeclaredMaintainerReleaseRoots(...)) == 0`. That function (`release_ownership.go:127-139`) filters the `DeclaredModuleSources` result on `trackedSet[source+"/VERSION"]`. New: `!slices.ContainsFunc(shippedSources, trackedSet[source+"/VERSION"])`. Same condition on the same slice. A nil result (no tracked `suite/modules.tsv`) gives `ContainsFunc == false`, so nil is returned, as before.
- **Error paths.** Old: the first error came from `DeclaredMaintainerReleaseRoots`, which only forwards the `DeclaredModuleSources` error unchanged. Both versions wrap it as `RELEASE-SHIPPED-CHANGE-UNVERIFIABLE: %w`. The old second read could not produce a different error (same tracked set, same HEAD image). Identical.
- **Roots list and order.** `DeclaredModuleSources` returns `sortedUniqueRoots(...)`: unique and sorted. The old roots were a subset of the sources, so the old set union equals sources ∪ {suite, tools}. New code appends `suite`/`tools` only when absent, then sorts. Same set, same sort, so the refusal message string is identical. A source literally named `suite` or `tools` is handled by the `slices.Contains` check (no duplicate).
- **Aliasing.** `roots := shippedSources` shares the backing array, and `append`/`sort.Strings` may mutate it. `shippedSources` is a fresh local slice built by `sortedUniqueRoots` and is not used after this point, so this is safe. The old code had the same kind of aliasing (`roots[:0]` reuse).
- **Reads.** One `DeclaredModuleSources` call and one `headReleaseImage(repositoryRoot)` call remain. `DeclaredMaintainerReleaseRoots` and `DeclaredModuleSources` are unchanged. `DeclaredMaintainerReleaseRoots` keeps its other caller at `release_ownership.go:65`.
- **Comments.** The function doc comment (`finalization_release_guard.go:22-27`) still describes the behaviour correctly. A grep of the finalization and releaseownership packages for `seenRoots`, merge or two-read wording finds no stale comment.

## Important findings

None.

## Minor findings

- F1 `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:44`: `roots := shippedSources` followed by in-place `append` and `sort` reuses the source slice. It is safe today because `shippedSources` is not read again, but a later edit that reads `shippedSources` after line 48 would see mutated data. A one-word name such as `roots := slices.Clone(shippedSources)` would remove the trap, though it is not required. impact-negligible
- F2 `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go:168`: no existing test covers the "modules declared, but no declared source has a tracked VERSION" exit. Only the "no declaration at all" exit is tested. This is the one branch whose condition the REQ rewrote, and the equivalence is shown only by reading the code. impact-negligible
- F3 `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go:62`: refusal tests assert only the `RELEASE-WITHOUT-SHIPPED-CHANGE` token, not the root list `(skills/..., suite, tools)`. A change to the root order or content would not be caught by tests. impact-negligible

## Acceptance evidence

- `bash do-work/runs/work-2026-10-07-150106/REQ-638-probe.sh` from the repo root at the merge: `ok .../internal/finalization 6.140s`, exit 0.
- `go test -count=1 -run '^TestRelease' ./internal/finalization/ -v`: 9 top-level `--- PASS` lines, no failures.
- `go vet ./internal/finalization/ ./internal/releaseownership/`: clean (exit 0). `gofmt -l internal/finalization`: empty.
- `git diff --stat da3ab5f5..6125e58c`: 1 file, within write_set.
- `_dev/tests/maintainer-verify.sh` not run, as instructed.

## Suggested testing

1. Add a case to the release-guard tests where `suite/modules.tsv` is tracked and declares sources, but no `<source>/VERSION` is tracked. Assert the guard returns nil (pins the rewritten consumer-project exit, F2).
2. Assert the full refusal message once, including the sorted root list ending in `suite, tools` (pins byte-identical output, F3).
3. Optional: a fixture whose `modules.tsv` declares a source literally named `suite`, asserting `suite` appears once in the refusal message.
