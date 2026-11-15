# Hand-back — REQ-723 (the archive-collision helper skips assets folders)

Branch `worktree-agent-REQ-723-archive-collision-skips-assets`, tip `1be00b67`, base `df1bf423`.

## Commits
- `3cbdc83e` [REQ-723] RED: archive-collision test pins the assets/ skip in both directions
- `1be00b67` [REQ-723] archive-collision helper skips assets/ folders

## Files
- modified `skills/do-work/tools/do-work-cli/internal/corehelpers/archive_collision_test.go` (new `TestArchiveCollisionSkipsAssetsFolders` at the end)
- modified `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go` (`fs.SkipDir` for folders named `assets` in `handleArchiveCollision`, after the archive-root checks)

## Scope changes
None.

## Red-green evidence
- After commit 1 (RED): `go test -count=1 -run TestArchiveCollision ./internal/corehelpers/` →
  `--- FAIL: TestArchiveCollisionSkipsAssetsFolders` / `archive_collision_test.go:65: evidence file in assets/ reported as archived twin: outcome=findings output="do-work/archive/UR-001/assets/REQ-001-evidence.md\n"`
- After commit 2: `go test -count=1 ./internal/corehelpers/` → `ok ... 8.785s` (new test and both existing archive-collision tests pass, unmodified).
- GREEN probe: fails at base by design (test absent; confirmed by the pre-dispatch agent at main HEAD, and the test failed at my RED commit). At tip `1be00b67`: `REQ-723 probe PASS`.

## Verification
1. `gofmt -l ./internal/corehelpers` → no output. `go vet ./internal/corehelpers/` → clean. `go test -count=1 ./internal/corehelpers/` → ok (8.785s, single run, no rerun needed).
2. End to end: `mktemp -d` Git repo holding only `do-work/archive/UR-001/assets/REQ-001-evidence.md`; `bash <worktree>/skills/do-work/tools/do-work-cli.sh --repo-root <repo> --format text archive-collision REQ-001` → no output, exit 0.
3. `bash .../REQ-723-probe.sh` from the worktree → `REQ-723 probe PASS`.

## P-A-U
- PLAN: copy the `cleanup_git.go:115` assets skip (comment and branch) into the `handleArchiveCollision` WalkDir callback; RED test first, both directions.
- APPLY: one test function (with a small local `outputOf` closure so failure messages print the output text, not a pointer), one 4-line branch.
- UNIFY: `git diff --stat df1bf423..HEAD` → 2 files, +28: archive_collision_test.go (+24), commands.go (+4); checked no debug output, gofmt/vet clean, missing-archive / not-a-directory / symlink paths untouched (skip placed after the root checks).

## Decisions
- **D-03** DECIDE & STATE: failure messages print `outcome` and the output text via a local closure instead of `%+v` of the whole `CommandResult` (which prints ~40 fields and a pointer address). Kept local to the test; no shared helper added.

## Discovered Tasks
None.

## Integration seams
Touched only the two files in the write set. No shared files.

## Proposed CHANGELOG entry
**Archive-Collision Check Ignores Evidence in assets/ Folders**
What changes for an installed copy: `do-work-cli archive-collision REQ-NNN` (and `tools/checks/archive-collision.sh`) no longer reports evidence files inside `do-work/archive/**/assets/` as archived twins; real archived records are still reported.

## Proposed lesson
none (the REQ-722 `guard-on-one-entry-path` bullet already covers this; this REQ is a second instance of it).

## Lessons read
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` line 5 (`[family: guard-on-one-entry-path]`, REQ-722), as quoted in the brief and read in the worktree.
