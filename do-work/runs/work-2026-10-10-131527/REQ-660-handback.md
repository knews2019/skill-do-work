# REQ-660 hand-back (do-work-cli worktree new, status, merge and cleanup)

- Branch: `worktree-agent-REQ-660-worktree-lifecycle-command`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-660-worktree-lifecycle-command`
- Base commit (before the first edit): `bd56c4b0`
- Commits: `cd557b36` `[REQ-660] add worktree new, status, merge and cleanup lifecycle command` (one commit)
- Integration seams: none. All six paths are committed on the branch.

## File manifest

- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle.go` (new): the `worktree` handler with `new`, `status`, `merge` and `cleanup`; reuses `parseWorktrees`, `requestIDFromWorktree`, `cleanupGit`/`cleanupGitBytes`, `gitExitSuccess`, `splitZero`; never calls `ApplyWorktreeRepairs`.
- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go` (new): exactly the five named tests, through the real `commandruntime` seam with `--format json`.
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_commands.go` (modified): `"worktree": handleWorktree` added to `Handlers()`.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified): `WorktreeLifecycleResult` and `WorktreeStatusRow` types after `LifecycleTimingResult`; `Worktree` field on its own line after `LifecycleTiming`; a text render block right after the timing block; `time` import.
- `skills/do-work/actions/fan-out-reference.md` (modified): one short paragraph each in `### Naming`, `### Where worktrees live` (link file plus `worktree status`), `### When to merge` (after step 4), and one sentence appended in `### Cleanup — happy path`. No rule text rewritten; every hand sequence stays valid.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modified): the `internal/cleanup/` Package routing bullet now names the `worktree new|status|merge|cleanup` command. Nothing else.

## P-A-U

**[PLAN]** Add the typed result block first so the tests compile, write the five tests against `Handlers()` through `commandruntime.NewRuntime`, run RED (missing handler), then implement one file in `internal/cleanup/`. All git work runs against the main tree, found as the first `git worktree list --porcelain -z` record. Each subcommand checks every precondition before its first side effect: `new` checks request uniqueness, base branch, link config and the name before `git worktree add`; `merge` checks MERGE_HEAD, index, ahead count and queue guard before `git merge`; `cleanup` checks merged-ness and dirt (owned links filtered out) before removing any link. Refusal `next_argv` is always a read-only git command.

**[APPLY]** Implemented as planned. One correction during GREEN: `RequestFile.FilenameID` is the normalized id (not the raw filename prefix), so the suffix is derived by stripping `^REQ-[0-9]+-?` from the basename with a regex instead.

**[UNIFY]** `git diff bd56c4b0 --stat`:

```
 skills/do-work/actions/fan-out-reference.md        |   8 +-
 .../internal/cleanup/cleanup_commands.go           |   2 +-
 .../internal/cleanup/worktree_lifecycle.go         | 505 +++++++++++++++++++++
 .../internal/cleanup/worktree_lifecycle_test.go    | 219 +++++++++
 .../internal/resultmodel/result_model.go           |  56 +++
 .../tools/do-work-cli/lessons-do-work-cli.md       |   2 +-
 6 files changed, 789 insertions(+), 3 deletions(-)
```

Checks (from the worktree root):

| Check | Exit | Wall time |
|---|---|---|
| `gofmt -l .../internal/cleanup .../internal/resultmodel` | 0, no output | <1 s |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/cleanup/ ./internal/resultmodel/` | 0 | 0.13 s (cached build) |
| `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/resultmodel/` | 0 | 0.81 s |
| `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/cleanup/` (whole touched package, regression) | 0 | 16.3 s |
| five named tests (GREEN) | 0 | 2.65 s |
| `bash .../REQ-660-probe.sh` | 0, "REQ-660 probe: all checks passed" | 4.61 s |
| `git diff --check` | 0 | <1 s |

Extra manual check: a `mktemp -d` fixture (outside both trees, removed after) with a REQ file named `REQ-5-Odd Name's Slug!.md`: `new` derived `worktree-agent-REQ-5-oddnamesslug`, `status` showed one clean row, `merge` refused the empty hand-back with exit 1, `cleanup` removed it with exit 0, an unknown subcommand refused with exit 1. Text rendering read correctly.

Files checked by reading the final diff: all six above. No debug artifacts. `"--force"` and `"-D"` do not appear in `worktree_lifecycle.go`.

## Proof record

- **RED** (before the handler existed; result type already present so the tests compiled): exit 1. All five tests `--- FAIL`, each with `worktree [new REQ-41] carried no typed worktree result: exit 2` and the runtime's `UNKNOWN-COMMAND` finding.
- **GREEN**: exit 0.
  ```
  --- PASS: TestWorktreeLifecycleNewMergeCleanupLeavesNothingBehind (0.69s)
  --- PASS: TestWorktreeMergeRefusesBuilderCommitUnderDoWork (0.37s)
  --- PASS: TestWorktreeMergeReportsEmptyHandBackWithoutCommit (0.31s)
  --- PASS: TestWorktreeCleanupRefusesDirtyWorktreeBeforeAnySideEffect (0.30s)
  --- PASS: TestWorktreeNewAppendsNumericSuffixOnCollision (0.24s)
  ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/cleanup	2.200s
  ```
  No skips.
- **Queue-guard mutation**: the guard condition was disabled (`false &&`); `TestWorktreeMergeRefusesBuilderCommitUnderDoWork` failed with `merge: exit 0 outcome success findings []`. After restoring the file from a backup copy: `--- PASS` (0.46 s). `grep -c 'false &&'` showed 1 during the mutation and 0 after the restore.
- **Probe**: exit 0, 4.61 s.
- No load reruns were needed. Nothing failed under load.

## Decisions

- D-07: `new` returns outcome `success` (exit 0) when the worktree was created, even when a link was skipped or a collision suffix was used. Those are `warning` findings (`WORKTREE-LINK-SKIPPED`, `WORKTREE-NAME-COLLISION`). Reason: exit 1 would tell a script the worktree was not created. DECIDE & STATE.
- D-08: Refusal findings use a read-only `git ...` argv as `next_argv` (for example `git -C <main> status --short`, the queue-guard diff itself), never `do-work-cli worktree status`. Reason: `NormalizeResult` clears a refusal's `next_argv` whose verb equals the command (`worktree`), so that suggestion would silently disappear. `worktree status` stays in `verification_argv`. DECIDE & STATE.
- D-09: A merge that stops on conflicts returns outcome `findings` (exit 1, code `WORKTREE-MERGE-CONFLICT`, `conflicted_paths` set), not `refused`, because a merge is now in progress (a side effect happened). A `git merge` error with no conflicted paths returns `failure` with git's stderr. DECIDE & STATE.
- D-10: `status` uses each worktree's `HEAD` (from the porcelain record) for ahead/behind and last-commit time, not the branch name. For an attached worktree these are the same commit; a detached builder worktree still gets a row. DECIDE & STATE.
- D-11: `status` filters owned links out of the dirty check, using the same helper as `cleanup`. Without this, every worktree with a link reads as dirty, and test 1 requires "clean" after `new` created a link. So `worktreeClean` is not reused; a filtered porcelain read (`dirtyPaths`) replaces it in this file only. DECIDE & STATE.
- D-12: `status` refuses on a detached main tree, like `merge` and `cleanup`, because there is no integration branch to count against. Only `merge` and `cleanup` refuse a `worktree-agent-*` branch checked out in the main tree; `new` and `status` do not. DECIDE & STATE.
- D-13: `merge` also refuses when `MERGE_HEAD` exists (brief item: "a merge is already in progress"), checked first, before the index check. DECIDE & STATE.
- D-14: The suffix keeps only `[a-z0-9-]` after lowercasing: other characters are dropped, not replaced by `-`, and leading/trailing `-` are trimmed. Spaces therefore join words (`odd name` becomes `oddname`). If the slug is empty, the name is `worktree-agent-REQ-N`. Reversible. DECIDE & STATE.
- D-15: Link config lines are trimmed and cleaned. An absolute path, a `..` segment, or a line that cleans to `.` refuses the whole command before anything is created. There is no comment syntax (YAGNI). DECIDE & STATE.
- D-16: All git work runs in the main tree (the first porcelain record), including `merge`, even when the command is started from a linked tree. REQ discovery and the link config also read the main tree. This follows "State stays home". DECIDE & STATE.

## Discovered Tasks

- impact-minor: `cleanup` Pass 5 (`ApplyWorktreeRepairs` in `internal/cleanup/cleanup_git.go`) uses `worktreeClean`, which counts the symlinks `worktree new` creates as dirt. A merged leftover that still has its links would be reported as dirty and need consent there, even though `worktree cleanup` would remove it. Confirmed by code reading: `worktreeClean` is an unfiltered `status --porcelain`. → report only
- impact-minor: The Crash Recovery sweep (`actions/work-reference.md` → Crash Recovery (Step 1)) has the same dirt reading for linked leftovers. It is not checked in code here. → report only
- impact-negligible: `worktree cleanup` leaves the empty `<repo>-worktrees/` parent directory after the last worktree is removed. The hand path does the same. → report only

## Lessons read

- Satellites: `_dev/primes/lessons-releases.md` (whole file).
- Primes: `_dev/primes/prime-shell-commands.md`, `_dev/primes/prime-action-files.md`, `_dev/primes/prime-releases.md`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (Traps).
- Families quoted in the brief and applied: `destructive-next-argv` (D-08), `commit-preflight-before-side-effect` (preflight order in each subcommand; test 4 pins it for cleanup), `opaque-evidence-projection` (typed `worktree` block, no packed detail strings), `fixture-cost-is-subprocess-spawning` (five tests total about 1.9 s; the setup is one shared helper).
- Not read (dropped for budget, as the REQ records): `lessons-do-work-cli.md` as a whole, `lessons-shell-commands.md`, `lessons-action-files.md`.

## Anti-bloat check

`git diff --stat` is above. Things added that the brief did not name, each with its reason:

- `WorktreeStatusRow` type (resultmodel): the brief lists status rows as a field with six members, which needs a row type.
- `time` import (resultmodel): the text block shows the last commit as an age.
- `worktreeLinksConfig` const: the one path string, used by the reader and in findings.
- `lifecycleRequestPattern` and `lifecycleFilenamePrefix` vars: validate the `REQ-N` argument; strip the id prefix from the filename for the suffix.
- `lifecycleRun` type with methods `integrationBranch`, `createWorktree`, `reportStatus`, `mergeBuilderBranch`, `removeBuilderWorktree`, `resolveOperativeName`, `builderBranches`, `readLinkConfig`, `ownedLinks`, `dirtyPaths`, `finding`, `stop`, `finish`, plus the functions `handleWorktree`, `parseLifecycleArguments`, `lifecycleNameTaken`: one function per subcommand and per rule the brief names. `ownedLinks` and `dirtyPaths` are shared by `status` and `cleanup`. `finding`/`stop`/`finish` build the result from one observation set.
- Test helpers `lifecycleRepository`, `runWorktreeCommand`, `lifecycleWorktreePath`, `hasFindingCode`, `commitInWorktree`, const `lifecycleOperativeName`: shared fixture setup for the five tests, so git spawns stay low.
- Finding codes: `WORKTREE-USAGE`, `WORKTREE-GIT-FAILED`, `WORKTREE-DETACHED-HEAD`, `WORKTREE-BUILDER-BRANCH-CHECKED-OUT`, `WORKTREE-DISCOVERY-FAILED`, `WORKTREE-REQUEST-NOT-UNIQUE`, `WORKTREE-BASE-NOT-A-BRANCH`, `WORKTREE-NAME-COLLISION`, `WORKTREE-LINK-SKIPPED`, `WORKTREE-LINK-FAILED`, `WORKTREE-LINK-CONFIG-UNREADABLE`, `WORKTREE-LINK-CONFIG-INVALID`, `WORKTREE-MERGE-IN-PROGRESS`, `WORKTREE-INDEX-NOT-EMPTY`, `WORKTREE-EMPTY-HAND-BACK`, `WORKTREE-QUEUE-GUARD`, `WORKTREE-MERGE-CONFLICT`, `WORKTREE-OPERATIVE-NAME-UNRESOLVED`, `WORKTREE-NOT-MERGED`, `WORKTREE-DIRTY`: one per refusal or failure the brief lists.
- Flags: only `--from` (new) and `--name` (merge, cleanup), as the brief names. No other flags, config keys or files.

## Proposed CHANGELOG entry

**Builder Worktree Lifecycle Command**

Builders no longer need hand-typed `git worktree add`, `ln -s`, merge and cleanup commands: `do-work-cli worktree` runs the same steps `actions/fan-out-reference.md` already defines, and it refuses before any side effect when a step is unsafe.

- `worktree new REQ-NNN [--from <branch>]` derives `worktree-agent-REQ-NNN-<suffix>` from the filename, adds `-2`, `-3` on a collision without deleting anything, creates the worktree in the sibling `<repo>-worktrees/` directory, and symlinks the paths listed in an optional `do-work/worktree-links` file.
- `worktree status` lists each builder worktree with ahead/behind, dirty or clean, and last-commit age.
- `worktree merge REQ-NNN [--name <operative_name>]` runs hand-back steps 1 to 4 (clean index, `<pre>`, queue guard, `--no-ff` merge, `[REQ-NNN]` commit) and reports `<pre>` and `<merge_hash>`; an empty hand-back refuses without a commit, and a conflict stops with the merge in progress.
- `worktree cleanup REQ-NNN [--name <operative_name>]` refuses an unmerged branch or a dirty worktree, then removes its links, the worktree and the branch, and prunes. No subcommand passes `--force` or `-D`.
- The hand sequences in `fan-out-reference.md` stay valid for harnesses without the command.

## Proposed lesson bullet

Satellite: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`.

- [family: destructive-next-argv] [REQ-660: a refusal's `next_argv` whose verb is the command itself is cleared by `NormalizeResult`, so a self-referencing read-only suggestion such as `do-work-cli worktree status` silently vanishes; point refusals at a read-only `git` argv and keep the self-check in `verification_argv`](../../../../do-work/archive/REQ-660-worktree-lifecycle-command.md)

## Integration seams and test wall times

- Seams: none to apply. Expected adjacent-line merges: `result_model.go` (REQ-658, REQ-690 may add a `CommandResult` field or render block next to mine; keep both); `fan-out-reference.md` (REQ-689 edits other sections; my edits stay in Naming, Where worktrees live, When to merge and Cleanup — happy path); `lessons-do-work-cli.md` (I edited only the `internal/cleanup/` routing bullet in place).
- Test wall times: five named tests 2.65 s total (individual 0.24 to 0.69 s); whole `internal/cleanup` package 16.3 s; `internal/resultmodel` 0.81 s; probe 4.61 s.
