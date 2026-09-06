## 0.305.12 — Rollback Decides the Handle Once and Closes Nil-Handle Panic (2026-09-06)

Transaction rollback now settles the rooted filesystem handle once upon entering `rollbackFailure`. If `os.OpenRoot` fails, the transaction immediately executes `rollbackWithoutRoot` to perform Git unstaging, tracked restoration from HEAD, and preimage restoration by pathname while safely recording unavailable root errors for rooted mutations. The missing guard in `quarantineAndRollbackPrivate` that caused a nil-pointer dereference panic when rolling back identity-recorded private untracked files is eliminated, and all eight downstream nil checks across rollback loops 1–3 were removed and pinned at zero.

- `rollbackFailure` cleanly partitions into `rollbackWithRoot` and `rollbackWithoutRoot`, ensuring that rooted operations only execute under a valid handle.
- All eight defensive nil-root checks in `git_transaction.go` have been removed, and `_dev/tests/audit-lockins.sh` Finding 3 now pins the count at zero.
- Added `TestRollbackWithoutRootHandleUnstagesRestoresFromHeadAndReportsTheRest`, providing the package's first unit test verifying rollback behavior when the root handle cannot be opened.
