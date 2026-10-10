# REQ-685 hand-back

- Branch: worktree-agent-REQ-685-rollback-uses-open-root
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-685-rollback-uses-open-root
- Base: b629e5cd. Commit: 62306f91 (one commit).

## File manifest
- skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go (modified): rollbackFailure takes `root *os.Root` and calls rollbackWithRoot only; deleted rollbackWithoutRoot, openRollbackRoot, the open-and-fallback branch; fixed the 3 comments naming them; added the "HEAD is deliberate" comment (2 lines).
- .../gittransaction/git_transaction_test.go (modified): deleted TestRollbackWithoutRootHandleUnstagesRestoresFromHeadAndReportsTheRest and the "root unavailable" case; the loop over []bool and `rootAvailable` collapsed (body de-indented, assertions unchanged); removed the now-unused `syscall` import.
- .../gittransaction/exact_commit.go (modified): the same "HEAD is deliberate" comment (2 lines, wrapped; the brief said "one comment line", a single line would exceed the file's line width).
- _dev/tests/audit-lockins.sh (modified, comment only): nil-root ratchet prose now says the handle is passed in and a rooted call never tests it. Ratchet code untouched.

## P-A-U
- [PLAN]: Pass the root ExecuteTransaction opens (and defer-closes) into rollbackFailure as a required argument; use rollbackWithRoot only; delete the no-root half, its hook and its two tests; comment-only edits elsewhere. The nil-handle panic REQ-598 fixed stays impossible: rollbackFailure receives the root as a required argument and ExecuteTransaction returns before any change when the open fails, so there is no nil path. A held root follows a renamed repository folder while `git -C <path>` follows the path; nothing renames the folder today.
- [APPLY]: Done as planned, 4 files, no new function, constant, option, file or test.
- [UNIFY]: see checks below. Files checked: all 4 for diff content, gofmt, vet, stale references.

### Checks
- go test -count=1 ./internal/gittransaction/ ./internal/finalization/: exit 0 (10.9 s + 64.2 s; 1:05 total)
- grep -rn "rollbackWithoutRoot\|openRollbackRoot" internal _dev/tests: no hits (grep exit 1)
- REQ-685-probe.sh: exit 0 (about 8 s)
- gofmt -l gittransaction: empty; go vet ./...: exit 0, clean
- bash -n audit-lockins.sh: 0; bash _dev/tests/audit-lockins.sh: exit 0 ("Audit lock-in regressions passed.")
- git diff --check: exit 0
- No reruns were needed.

### git diff b629e5cd --stat
```
 _dev/tests/audit-lockins.sh                        |  22 +--
 .../internal/gittransaction/exact_commit.go        |   2 +
 .../internal/gittransaction/git_transaction.go     | 110 +++-----------
 .../gittransaction/git_transaction_test.go         | 163 +++------------------
 4 files changed, 45 insertions(+), 252 deletions(-)
```
Net -207 lines. (The REQ estimated about -55 for source only; the two deleted tests account for most of the rest.)

## Decisions
- D-01 (decided): "HEAD" comment is two wrapped lines at each site, not one, for line width. Same content the REQ names.
- D-02 (decided): the collapsed creation-intent test keeps its name and drops the t.Run wrapper; assertions unchanged.
- D-03 (decided): no new test. Rollback reporting a failing root operation as incomplete stays pinned by the retained tests (for example TestIncompleteRollbackReportsRiskWithoutRecursiveDeletion, and the creation-intent test's RollbackIncomplete assertion).

## Discovered Tasks
- None.

## Lessons read
- _dev/primes/lessons-releases.md (whole file). lessons-do-work-cli.md not read (no code its bullets name was touched beyond what the REQ dropped).

## Anti-bloat check
Diff stat above. Added beyond the REQ: nothing (no function, constant, option, file or test). Only added lines are the 4 comment lines and the de-indented test body.

## Proposed CHANGELOG entry
Title: Git transaction rollback now reuses the transaction's open repository root
Why: rollback reopened the root by path although the transaction already held one, and a near-duplicate function existed only to survive a failed reopen that cannot happen.
- rollbackFailure takes the open root as a required argument and uses rollbackWithRoot only.
- Deleted rollbackWithoutRoot, the openRollbackRoot test hook and the two tests that forced the no-root path (about 200 lines).
- The nil-handle panic fixed earlier stays impossible: the argument is required and the transaction returns before any change when its open fails.
- Comments at the two "HEAD" commit-ID sites say the literal is deliberate (a non-empty PrimaryCommit blocks rollback of a commit that landed). No behavior change.

## Proposed lesson bullet
none

## Integration seams and times
No seam expected; finalization package tests ran green (64 s). gittransaction tests 10.9 s.
