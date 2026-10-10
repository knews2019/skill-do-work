# REQ-680 hand-back (archive fetch test keeps an absent target absent)

- Branch: worktree-agent-REQ-680-archive-fetch-absent-target-test
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-680-archive-fetch-absent-target-test
- Base: b629e5cd. Commit: 9b8dee8b ("[REQ-680] archive fetch test keeps an absent target absent"), one commit.

## File manifest
- skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch_test.go (modified): `TestTotalFailurePreservesTheTargetAndLeavesNoScratch` is now a two-row table (`pre-existing target`, `absent target`). Same server, same failing request, same failure-report fragment checks in both rows. The absent row asserts `os.Lstat` reports not-exist and `assertNoArchiveScratch` finds nothing. Comment rewritten in plain words, no consumer commit IDs or REQ numbers.

## P-A-U
- [PLAN]: Route A. Extend the one existing test into a table; reuse the existing body, server and `assertNoArchiveScratch`. No helper, no production change.
- [APPLY]: Wrapped the body in `t.Run` per row. One field `existingBytes` (empty means no seed) instead of a separate bool, so the seeded literal appears once. Absent row returns after its two assertions.
- [UNIFY]: `git diff b629e5cd --stat`: 1 file, 63 insertions, 42 deletions (most of it re-indentation). Checks: go test -v archivefetch exit 0 (3.9s wall, `--- PASS: .../absent_target` shown); REQ-680-probe.sh exit 0 (0.46s); gofmt -l empty (exit 0); go vet exit 0 (0.13s); git diff --check exit 0. File checked: archive_fetch_test.go only; no debug artifacts.

## Mutation record (not committed, reverted)
In `archive_fetch.go` `prepareDownloadCandidate` I made the stage name equal the target name and removed the stage cleanup defer (no private staging). Result: the new case failed with "a failed fetch created the absent target: <nil>" (archive_fetch_test.go:423). But `pre-existing_target` also failed ("the pre-existing target was overwritten") and about 15 other tests failed, because this function is shared by `DownloadAtomic` and the publish path. So "only the new case fails" does not hold literally for this mutation: removing staging necessarily also overwrites a pre-existing target. The new case does fail on its own distinct assertion line. Reverted with `git checkout -- archive_fetch.go`; `git diff b629e5cd --stat` shows only the test file.

## Decisions
- D-01 DECIDE & STATE: table row, not a sibling test (REQ gave latitude; it reuses the existing body, smallest addition).
- D-02 DECIDE & STATE: `existingBytes string` field rather than `seedTarget bool` plus `expectedTarget`; one field, same coverage.
- D-03 ESCALATE (coordinator): the mutation proof wording "only the new case fails" cannot be met by removing private staging (see mutation record). Value of accepting as is: the new case provably catches the named failure. Risk: none to the code; only the REQ's proof sentence is stricter than reality. Recommend recording it as "new case fails on its own line; the pre-existing row also fails, by design".

## Discovered Tasks
- `TestMissingArchiveTargetParentReportsUnattemptedRoutesInTextAndJSON` is SKIPPED in this package run (SKIP seen in the -v output). Reason not investigated. → report only

## Lessons read
- _dev/primes/lessons-releases.md (whole file). Go satellite lessons-do-work-cli.md not read (no code it names was touched).

## Anti-bloat check
`git diff --stat`: archive_fetch_test.go, 63 insertions, 42 deletions. Added beyond what the REQ named: none (no function, constant, option, file or extra test; one anonymous struct field `existingBytes` that is part of the table).

## Proposed CHANGELOG entry
Archive fetch test now proves a failed fetch never creates a missing target. A failed update must not leave a half-fetched archive at a path that held nothing before, and no test pinned that.
- `TestTotalFailurePreservesTheTargetAndLeavesNoScratch` is a table with `pre-existing target` and `absent target` rows.
- The absent row asserts no file at the target path and no scratch file beside it.
- Test only; no production change.

## Proposed lesson bullet
none

## Integration seams and wall times
None. Wall times: go test 3.9s, probe 0.46s, vet 0.13s.
