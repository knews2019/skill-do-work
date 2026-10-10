## Review: REQ-685

**Approve**: the rollback now uses the root that the transaction already holds, the no-root path is gone, and no nil handle can reach the rooted code.
Route B | 96636dbb (range 33a0f1b1..96636dbb, 4 files, +45/-252)

### What's built
- `rollbackFailure` takes `root *os.Root` from `ExecuteTransaction` and calls only `rollbackWithRoot`. `rollbackWithoutRoot`, the `openRollbackRoot` test hook, the open-and-fallback branch and the two no-root tests are deleted.
- Comments that named the deleted code are fixed in `git_transaction.go` and `_dev/tests/audit-lockins.sh`. A two-line "HEAD is deliberate" comment is added at both commit-ID sites.

### Decisions / risks for you
- None. The held root follows a renamed folder while `git -C <path>` follows the path. Nothing renames the folder today, and the PLAN records this.

### Findings

**Important:**
- None.

**Minor:**
- F1. `git_transaction.go:704-705`: the new comment says "HEAD" is deliberate because a non-empty `PrimaryCommit` blocks rollback (`finalization_apply.go`, near line 30). That is true for `exact_commit.go`, because finalization calls `CommitExactPaths` and stores its `CommitSHA` in `journal.PrimaryCommit` (`finalization_apply.go:136-140`). No `ExecuteTransaction` result ever reaches `PrimaryCommit`. The readers that treat this site's non-empty "HEAD" as "a commit exists" are `cleanup/cleanup_apply.go:144`, `doctor/doctor_repair.go:173` and `publication/publication_commands.go:234`. The REQ text prescribed this wording for both sites, so the builder followed the capture. The comment points readers to the wrong file. Behavior is unchanged. impact-negligible → report only

**Nit:**
- None.

### Review focus answers
1. **Nil-handle panic (REQ-598) stays impossible.** `root, rootErr := os.OpenRoot(repositoryRoot)` is at `git_transaction.go:520`. A failed open returns through `failTransaction` at `:521-523`, before any mutation, and then `defer root.Close()` runs at `:524`. All ten `rollbackFailure` calls (`:639-700`) pass that same `root`. Inside `ExecuteTransaction` nothing shadows or reassigns `root` and nothing closes it early (grep for `root :=`, `root =` and `root.Close` in `:448-735`). No other file calls `rollbackFailure`. The package has eight `os.OpenRoot` sites (`:140, :223, :339, :377, :391, :423, :520, :974`), and every one returns on error before it uses the handle. Every other `*os.Root` parameter receives a handle from one of those sites. `rg` finds no `root == nil` or `nil == root` test in the package, so the audit ratchet still holds at zero.
2. **Collapsed test.** `git diff -w` on `TestCreationIntentPreservesForeignIndexBeforeTransactionStaging` removes only the loop, the `rootAvailable` switch, the `t.Run` wrapper and the hook override. Its three assertions (foreign staging kept, foreign contents kept, `RollbackIncomplete`) match the old "root available" subtest exactly. **No coverage gap:** the deleted no-root test pinned behavior that only existed on the deleted path. The rooted equivalents are still pinned: restore from HEAD (`TestPreCommitFailureRestoresTrackedAndRemovesOnlyCreatedTargets`, `TestRefusedCommitRollsBackAndReportsCommitFailed`), unstaging after `git add` (the same refused-commit test asserts `status --porcelain` is empty), existing-untracked restore (`TestExecuteTransactionExistingUntrackedTargetsRequireOptInAndRestoreBytes`), the private quarantine where REQ-598 panicked (the `TestPrivate*` rollback tests), and the typed incomplete rollback (`TestIncompleteRollbackReportsRiskWithoutRecursiveDeletion` plus seven other `RollbackIncomplete` assertions).
3. **Bloat.** 30 added lines are not comments. 19 of them are the surviving test body moved out one indent level. The other 11 are real code changes: the 10 call sites gain the `root` argument and the `rollbackFailure` signature gains the parameter. The diff adds no helper, option, file or test. The `syscall` import is removed. The helpers that the deleted function shared (`restoreExistingUntracked`, `restoreTrackedFromHead`, `restoredByTargetLoop`, `deepestFirst`, `mapKeys`) all still have callers in `rollbackWithRoot`.
4. **Restatement sweep.** See the line in `## Review` below.
5. **"HEAD is deliberate" comments.** In `finalization_apply.go:30`, the deferred pre-primary rollback returns early when `journal.PrimaryCommit != ""`. `finalization_apply.go:139` sets `PrimaryCommit` from the `CommitExactPaths` result, so the comment in `exact_commit.go:78-79` is accurate. The comment in `git_transaction.go:704-705` names the wrong reader (F1).

### Requirements Checklist
- [x] `rollbackFailure` takes the open root as a parameter and uses `rollbackWithRoot` only: delivered
- [x] `rollbackWithoutRoot`, `openRollbackRoot`, the fallback branch and both no-root tests are deleted, and the PLAN states why there is no nil path: delivered
- [x] Existing rollback-incomplete tests are kept and no new test is added: delivered
- [x] One comment at each "HEAD" site with no behavior change: delivered (two wrapped lines each, D-01). The `git_transaction.go` copy names the wrong reader (F1).
- [x] Held-root versus `git -C` note in the PLAN and commit message: delivered
- [x] Tests, vet, gofmt and the probe are clean: delivered (reviewer reran `go vet` on gittransaction and finalization, clean, and `gofmt -l`, empty)

### Acceptance Testing
**Result: Pass**
- The integrator's gate run at `96636dbb` exited 0 and the REQ probe exited 0. The reviewer's `go vet ./internal/gittransaction/ ./internal/finalization/` and `gofmt -l internal/gittransaction` were both clean. The full test suite was not rerun, as the brief instructed.

### Suggested Additional Testing
- None.

### Scores (on the record, not the headline)

**Overall: 97%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All six criteria met |
| Code Quality | 92% | One comment names the wrong reader (F1) |
| Test Adequacy | 95% | Deletion only; the rooted paths stay pinned |
| Scope | 100% | Declared and touched files match |
| Risk | None | No nil path; call-site argument change only |
| Acceptance | Pass | Gate and probe green |

### Follow-ups created
- None (1 findings report only)

## Review

**Overall: 97%** | (timestamp added by the integrator)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** F1: the comment at `git_transaction.go:704-705` cites `finalization_apply.go` `PrimaryCommit`, but only `CommitExactPaths` results reach `PrimaryCommit`. This site's "HEAD" is read by `cleanup_apply.go:144`, `doctor_repair.go:173` and `publication_commands.go:234`. impact-negligible → report only
**Acceptance:** Pass. The gate and probe are green at 96636dbb, and the reviewer's vet and gofmt runs were clean.
**Restatement sweep:** redefined the rollback design (rollback reuses the transaction's held root, and the no-root rollback path is deleted). `rg -i` searched for `rollbackWithoutRoot`, `openRollbackRoot`, `no-root rollback`, `rollback root is unavailable`, `left in place; rollback root`, `reopens the root`, `no-handle rollback`, `open rollback root`, `root unavailable`, `decides once`, `TestRollbackWithoutRoot`, `REQ-598` and nil-root/rooted-handle wording. It covered `skills/` (excluding CHANGELOG history), `_dev/`, `do-work/lessons-index.md` and `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`. `docs/` does not exist. Nothing still restates the old design. The `audit-lockins.sh:628` failure message ("REQ-598 decided the rollback handle once at its open") still agrees with the new design.
**Suggested testing:** 0 items
**Follow-ups created:** None (1 findings report only)

*Reviewed by review-work action*
