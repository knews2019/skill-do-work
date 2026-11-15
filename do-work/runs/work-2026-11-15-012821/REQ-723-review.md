## Review: REQ-723

**Approve** — the archive-collision helper now skips `assets/` folders with the same rule as `repositorymodel` and cleanup, and the new test pins both directions.
Route A | 858c0708 (range 2e5b6f70..858c0708: RED `3cbdc83e`, fix `1be00b67`, merge `858c0708`)

### What's built
- `do-work-cli archive-collision REQ-NNN` no longer lists evidence files inside `do-work/archive/**/assets/` as archived twins. A real archived record beside an `assets/` folder is still listed.

### Decisions / risks for you
- None. Builder decisions D-01 (copy the skip rule, do not switch to the `repositorymodel` snapshot), D-02 (skip only `assets`) and D-03 (local `outputOf` closure for readable failure messages) are recorded in the REQ and match the diff.

### Findings

**Important:**
- None.

**Minor:**
- None.

**Nit:**
- N1: `archive_collision_test.go:55-60` adds a local `outputOf` closure, while the two sibling tests print failures with `%+v`. The closure is 6 lines, stays inside the one test, and makes a failure print the output text instead of a pointer address, so it earns its place. Not counted as bloat. — impact-negligible → report only

**Anti-bloat count: 0.** No helper, option, file or test beyond what the REQ named. The one new test is the test the REQ's GREEN clause asks for.

### Requirements Checklist

- [x] Skip folders named `assets` in the `handleArchiveCollision` WalkDir — delivered at `commands.go:174-177`: `entry.IsDir() && entry.Name() == "assets"` returns `fs.SkipDir`. This is the same condition as `repository_model.go:192` (any depth) and `cleanup_git.go:115`, with the same comment text as the cleanup copy.
- [x] Assets-only archive gives success with empty output — delivered, first half of `TestArchiveCollisionSkipsAssetsFolders` (asserts `OutcomeSuccess`, non-nil output, output `""`).
- [x] A real record in the same UR folder is still listed — delivered, second half asserts `OutcomeFindings` and output exactly `do-work/archive/UR-001/REQ-001-done.md\n`. The exact match also proves the evidence file stays hidden when a real record exists, so the test fails if the skip is removed or grows too wide.
- [x] The two existing archive-collision tests still pass unmodified — confirmed (run below; the diff does not touch them).
- [x] Missing-archive, not-a-directory and symlink behavior unchanged — the new branch sits after the `walkError` and root checks, so a missing root still returns nil and a file at the root still errors. The archive root itself is named `archive`, so the skip can never fire on it. A directory symlink named `assets` reports `IsDir() == false`, so the branch ignores it, and WalkDir still does not follow it, as before.

### Acceptance Testing

**Result: Pass** (stage covered: local Go unit tests only)
- `go test -count=1 -run TestArchiveCollision -v ./internal/corehelpers/` in `skills/do-work/tools/do-work-cli` on the main checkout at 858c0708: all three tests PASS (`TestArchiveCollisionFindsNestedRequestsWithoutFollowingDirectoryLinks`, `TestArchiveCollisionDistinguishesMissingArchiveFromUnreadableArchive`, `TestArchiveCollisionSkipsAssetsFolders`), `ok` in 0.187 s.
- Red-green evidence in the REQ Testing section: the new test failed at RED `3cbdc83e` with the evidence path in the output, and passes at `1be00b67` and the merge.

### Suggested Additional Testing

- None needed beyond the integrator's planned heavy lanes. The only launcher, `skills/do-work/tools/checks/archive-collision.sh`, passes arguments straight to the Go command, so the unit test covers the behavior.

### Code review notes (quick scan)

- Code Quality: one 4-line branch, idiomatic `fs.SkipDir`, comment states why. No debug output.
- Risk: None. The command is only reached through the compatibility launcher. The `advance` archive-collision hold reads the `repositorymodel` snapshot, which already skipped `assets/`, so the two now agree.
- Restatement sweep (not the wave's last integration, one-line check): `grep -rn archive-collision skills/` finds the launcher (no output contract stated), a historical CHANGELOG entry (line 925, still true) and `blocked-archive-collision` status docs that describe `advance`, not this command. Nothing restates the scan set, so nothing is stale.

### Scores (on the record — not the headline)

**Overall: 98%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All five checklist items delivered |
| Code Quality | 95% | Matches the reference rule and comment; nothing to remove |
| Test Adequacy | 95% | Both directions pinned with exact output; red-green recorded |
| Scope | 100% | 2 files, both in `write_set` |
| Risk | None | Root, missing-archive and symlink paths untouched |
| Acceptance | Pass | Local Go unit tests |

### Follow-ups created
- None (1 findings report only)

## Review

**Overall: 98%** | <TS>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** None. Nit N1 (no score weight): `archive_collision_test.go:55-60` local `outputOf` closure differs from the siblings' `%+v` failure format, but it is test-local and prints the output text instead of a pointer, so it earns its place — impact-negligible → report only. Anti-bloat count: 0.
**Acceptance:** Pass — local `go test -count=1 -run TestArchiveCollision ./internal/corehelpers/` at 858c0708: 3/3 PASS; red-green recorded at `3cbdc83e` → `1be00b67`.
**Restatement sweep:** redefined the archive-collision command's scan set (now excludes `assets/` folders); grep of `skills/` finds no restatement of it, nothing stale.
**Suggested testing:** 0 items
**Follow-ups created:** None (1 findings report only)

*Reviewed by review-work action*
