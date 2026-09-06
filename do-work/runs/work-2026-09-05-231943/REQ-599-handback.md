# REQ-599 Hand-Back (decide in-flight-ness from the walked root, not a path substring)

Branch: `worktree-agent-REQ-599-walked-root`
Head: `8164eb979623b1ff727e92b89be7468fea015d69`
Worktree: `/home/user/skill-do-work-worktrees/worktree-agent-REQ-599-walked-root` (clean after commit)

## Defect

`AssociateProjectPaths` in `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go`
decided in-flight-ness with `active := strings.Contains(filepath.ToSlash(path), "/working/")` on the
absolute walked path. A checkout beneath a directory named `working` made every archived REQ look
active, the terminal-success status filter was skipped, and a blocked archived REQ claimed paths.

## Fix

The two roots are now a slice of `{directory, active}` pairs (`do-work/working` true,
`do-work/archive` false). The walk callback reads `walkedRoot.active` and never inspects the path.
The slice keeps the original walk order (working first, then archive), so the latest-completion
tie-break between roots is unchanged.

## Red run (test added, inventory.go untouched)

Command: `go test ./internal/corehelpers/ -run 'TestAssociationUnderWorkingDirectoryCheckoutSkipsBlockedArchivedRequest' -count=1 -v`

```
=== RUN   TestAssociationUnderWorkingDirectoryCheckoutSkipsBlockedArchivedRequest
    inventory_test.go:396: blocked archived request claimed project.txt through the checkout path: owner="REQ-905"
--- FAIL: TestAssociationUnderWorkingDirectoryCheckoutSkipsBlockedArchivedRequest (0.00s)
FAIL
FAIL	github.com/knews2019/skill-do-work/do-work-cli/internal/corehelpers	0.006s
FAIL
exit=1
```

## Green run (after the fix)

Same command:

```
=== RUN   TestAssociationUnderWorkingDirectoryCheckoutSkipsBlockedArchivedRequest
--- PASS: TestAssociationUnderWorkingDirectoryCheckoutSkipsBlockedArchivedRequest (0.00s)
PASS
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/corehelpers	0.007s
exit=0
```

## Whole package, build, vet, gofmt

- Baseline before any change: `go test ./internal/corehelpers/ -count=1` -> `ok ... 6.482s`, exit 0.
- After the fix: `go test ./internal/corehelpers/ -count=1` -> `ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/corehelpers	4.169s`, exit 0.
- `go build ./...` exit 0. `go vet ./...` exit 0.
- `gofmt -l internal/corehelpers/inventory.go internal/corehelpers/inventory_test.go` printed nothing.

## Gate

Command: `DO_WORK_GATE_ROOT=/home/user/skill-do-work-worktrees/worktree-agent-REQ-599-walked-root bash <scratchpad>/gate.sh`

- Exit status 0, wall 86s. Last line: `Maintainer verification passed.`
- do-work-cli fast tests: 797 tests, wall 28s. queue-kanban fast tests: 402 tests, wall 28s.
- The heavy tier did not run (log line 57: `SKIP: staged skills, updater, and installer probes require maintainer-verify.sh --heavy.`), so the two heavyverification tests known to fail at the branch point for environmental reasons did not appear in this run.
- Full log: `/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/req599/gate.log`.

## Normal checkout: rule unchanged

The observable rule for a checkout that is not under a `working` directory is the same as before and is covered by existing tests that ran green in the package run above:

- `TestGenericAssociationNeverOwnsSharedDoWorkMetadata` (commands_test.go:28): an `archive/` REQ with `status: completed` claims `project.txt`.
- `TestProtectedInventoryPreservesAssociationPartitionByClassification` (inventory_test.go:458) via `writeInventoryOwner`: a `working/` REQ with `status: claimed` (non-terminal) claims its paths.
- `TestTerminalSuccessAliases` (inventory_test.go:364): the terminal-success alias list, and `cancelled` rejected.

The new test also holds a `working/` REQ with `status: blocked` in the same `working`-nested checkout and asserts it still claims its path, so a fix that disabled `working/` ownership would fail too.

## Diff summary

```
8164eb979623b1ff727e92b89be7468fea015d69 [REQ-599] Decide a REQ in-flight from the walked root, not a path substring

 .../do-work-cli/internal/corehelpers/inventory.go  | 21 +++++++++++-----
 .../internal/corehelpers/inventory_test.go         | 28 ++++++++++++++++++++++
 2 files changed, 43 insertions(+), 6 deletions(-)
```

Two files only: `inventory.go` (fix) and `inventory_test.go` (the new test, placed after `TestTerminalSuccessAliases`, following the package's `foo.go` / `foo_test.go` pairing).

## Found and not fixed

- Release bookkeeping (changelog entry, version bump, `skills/do-work/CHANGELOG.md` mirror) per `_dev/primes/prime-releases.md` is not in this commit. History shows the builder commit and the `[REQ-NNN] release:` commit are separate steps, so this is left to the finalize step.
- The REQ-599 file's P-A-U checkboxes are not ticked; the file is outside the two-file rule.
- Observation only: walk order matters for the tie-break. A `working/` REQ with no `completed_at` (zero time) loses a path to an `archive/` completed REQ that names the same path with any `completed_at`. This is pre-existing behavior and was preserved by keeping the roots in a slice rather than a map; whether it is the intended precedence was not evaluated here.
