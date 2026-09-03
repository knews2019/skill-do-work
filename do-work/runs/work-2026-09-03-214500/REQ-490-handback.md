# REQ-490 Builder Handback

- Branch: `codex/REQ-490-wave-depth`
- Commit: `87b72c7789d58152f1eaf4253629ba7b12bfa65f`
- Base: `7035dc3657080918d021a6b10ede0f7e00d99235`

## Changed files

- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go`
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go`

## Implementation

`queueDependencyDepth` now treats an edge absent from the current node's authoritative `UnmetDependencies` as satisfied, so duplicate records already resolved successful by the dependency graph contribute no queue depth. Pending unmet dependencies still recurse; missing, cyclic, single-record, and genuinely ambiguous-unsatisfied behavior is unchanged.

## RED/GREEN

The new `TestWaveZeroIncludesDependentWithSatisfiedDuplicateRecords` fixture creates two terminal-success records for REQ-311 and a pending REQ-312 dependent. Before the production edit it failed because selection was empty and REQ-312 was excluded with `WAVE-MISMATCH: dependency depth is 1, not requested wave 0`. After the edit it selects only REQ-312 and reports `DependencyDepth: 0`.

## Verification

- `go test -count=1 ./internal/nextselection -run 'TestWaveZeroIncludesDependentWithSatisfiedDuplicateRecords|TestWaveDepthAndFanOutAreSeparateSelectionAxes'` — pass
- `go test -count=1 ./internal/nextselection` — pass
- `go vet ./...` — pass
- `go test -count=1 ./...` — pass
- `git diff --check` — pass

The complete two-file diff was reviewed; it contains no lifecycle, generated, release, debug, or out-of-scope changes.

## Merge guidance

Merge the single commit with a no-fast-forward integration commit. The changed test file is a declared CLI heavy surface, so after focused/full integration checks and the canonical fast gate, park the request at `pending-heavy-testing` unless exact-revision heavy permission is available.
