# Hand-back: REQ-693 (worktree status/cleanup survive a missing worktree folder; dirt readers ignore the command's own links)

- Branch: `worktree-agent-REQ-693-worktree-missing-folder-owned-links`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-693-worktree-missing-folder-owned-links`
- Base commit: `fd9a2378`
- Commits: `22d7773e` `[REQ-693] worktree status and cleanup survive a missing folder; dirt readers ignore owned links` (one commit)
- Write boundary: exactly the eight files the brief names. Nothing under `do-work/` was touched in the worktree. This hand-back is the only main-tree write, and it is not staged.

## File manifest

- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle.go` (modified): the link rule now lives in package-level `readWorktreeLinkConfig`, `ownedWorktreeLinks`, exported `OwnedWorktreeLinks` and `worktreeDirtyPaths`. The `readLinkConfig` method is a thin caller that keeps `WORKTREE-LINK-CONFIG-UNREADABLE` / `-INVALID`. `reportStatus` marks a missing folder `Missing` and skips the dirt read. `removeBuilderWorktree` skips the dirt check, link removal and `worktree remove` for a missing folder and runs `worktree prune` before `branch -d`.
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go` (modified): Pass 5 derives the main root from the first porcelain record and reads the link config once. `worktreeClean` ignores owned links. The non-forced removal path deletes the owned links right before `git worktree remove`.
- `skills/do-work/tools/do-work-cli/internal/cleanup/worktree_lifecycle_test.go` (modified): the two new tests, plus the `context` import.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified): `WorktreeStatusRow.Missing` (`json:"missing"`). The text renderer prints `missing` where it prints `clean`/`dirty`.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/handoff.go` (modified): drops porcelain rows whose path is an owned link of that worktree, using `cleanup.OwnedWorktreeLinks`. The main root is the first `worktree list --porcelain` record.
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified): `ownedWorktreeLinkPaths` is the board's copy of the rule. `worktreeHasUncommittedWork` reads `status --porcelain -z` (still `--no-optional-locks`) and ignores owned links. The main tree path comes from `listWorktreeAgentWorktrees` and is carried in `worktreeAgentGitState.mainWorktreePath` to `classifyWorktreeLeftover`.
- `skills/do-work/actions/cleanup.md` (modified): Pass 5 step 3 says owned `do-work/worktree-links` symlinks do not count as dirt, and that Pass 5 removes them before `git worktree remove`.
- `skills/do-work/actions/fan-out-reference.md` (modified): the `worktree status` sentence adds "missing (its folder is gone)". It also says that `worktree cleanup` still deletes a missing worktree's branch.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** I read the brief, the REQ (D-01 to D-06), the crew rules, the do-work-cli prime and the kanban conventions, `lessons-releases.md`, and the three lesson families. Order: (1) add the data-only `Missing` field so the RED test compiles, (2) write both tests and run RED, (3) lift the link helpers and fix status and cleanup, (4) fix only `worktreeClean` and record the intermediate `WORKTREE-REMOVE-FAILED`, (5) add Pass 5 link removal, which turns GREEN, (6) update the handoff and board readers, (7) make the two prose edits, (8) verify.
- [x] **[APPLY]:** Done as planned. The two deviations are D-07 and D-08 below.
- [x] **[UNIFY]:** I read each of the eight files in `git diff fd9a2378` and checked for debug output, stray files and scope. Checks (run from the worktree root):

| Check | Exit | Wall |
|---|---|---|
| `REQ-693-probe.sh` | 0 | 2s |
| `REQ-693-preflight-probe.sh` | 0 | 6s |
| `gofmt -l skills/do-work/tools/do-work-cli/internal skills/do-work-board/tools/queue-kanban` (prints nothing) | 0 | 0s |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/cleanup/ ./internal/corehelpers/ ./internal/resultmodel/` | 0 | 0s |
| `go vet -C skills/do-work-board/tools/queue-kanban .` | 0 | 1s |
| `git diff --check` | 0 | 0s |
| Extra: `go test -count=1` on the whole `cleanup`, `corehelpers` and `resultmodel` packages | 0 | 15s |
| Extra: `go test -count=1 -run 'Worktree\|Verify'` in queue-kanban | 0 | 10s |

No wall-time budget failed, so nothing was rerun.

`git diff fd9a2378 --stat`:

```
 skills/do-work-board/tools/queue-kanban/verify.go  | 90 ++++++++++++++++++---
 skills/do-work/actions/cleanup.md                  |  2 +-
 skills/do-work/actions/fan-out-reference.md        |  2 +-
 .../do-work-cli/internal/cleanup/cleanup_git.go    | 36 +++++++--
 .../internal/cleanup/worktree_lifecycle.go         | 91 ++++++++++++++++------
 .../internal/cleanup/worktree_lifecycle_test.go    | 66 ++++++++++++++++
 .../do-work-cli/internal/corehelpers/handoff.go    | 18 ++++-
 .../internal/resultmodel/result_model.go           |  7 +-
 8 files changed, 267 insertions(+), 45 deletions(-)
```

## Proof record

Command: `go test -C skills/do-work/tools/do-work-cli -count=1 -v -run '^(TestWorktreeStatusAndCleanupSurviveMissingWorktreeFolder|TestPass5RemovesMergedWorktreeWhoseOnlyDirtIsAnOwnedLink)$' ./internal/cleanup/`

1. **RED (unchanged code; only the data-only `Missing` field was added so the test compiles):** exit 1, both `--- FAIL`.
   - Missing folder: `status: exit 2 rows []resultmodel.WorktreeStatusRow(nil)`, finding `WORKTREE-GIT-FAILED` with evidence `git status --porcelain=v1 -z --untracked-files=all: exit status 128: fatal: cannot change to '<tmp>/repo-worktrees/worktree-agent-REQ-41-lifecycle-demo': No such file or directory`.
   - Owned link: finding `WORKTREE-REQUIRES-CONSENT`, evidence `worktree-agent-REQ-141-owned-link is dirty=true merged=true request_state=settled`.
2. **Intermediate (only `worktreeClean` fixed, no link removal yet):** the missing-folder test `--- PASS`. The owned-link test `--- FAIL` with `WORKTREE-REMOVE-FAILED`: `fatal: '<tmp>/worktree-agent-REQ-141-owned-link' contains modified or untracked files, use --force to delete it`. This shows D-06 is needed.
3. **GREEN:** exit 0. `--- PASS: TestWorktreeStatusAndCleanupSurviveMissingWorktreeFolder (0.50s)` and `--- PASS: TestPass5RemovesMergedWorktreeWhoseOnlyDirtIsAnOwnedLink (0.34s)`.
4. **Regression:** the preflight probe exits 0 and prints "all named reader tests pass".
5. **Ad-hoc end-to-end check of the two readers that have no test (not committed; scratch fixture deleted afterwards):** I built both binaries, ran `worktree new REQ-141` in a fixture with `shared-deps` linked and `.gitignore` `shared-deps/`, and committed.
   - Valid config: `handoff-state-survey` gave `HANDOFF-WORKTREE-CLEAN` for both trees and no `-DIRTY`. `queue-kanban verify` gave `merged-worktree-leftover [fixable]`.
   - After adding a `/abs` line to the config: the survey reported `HANDOFF-WORKTREE-DIRTY` and the board reported `worktree-present-uncommitted-work`. So an invalid config filters nothing (D-04).

## Decisions

- **D-07 (DECIDE & STATE): the Pass 5 test uses REQ-141, not REQ-41.** With REQ-41, Pass 5 reports `request_state=malformed`, because `FilenameID` is zero-padded (`REQ-041`) while `requestIDFromWorktree` returns `REQ-41`. So a two-digit fixture can never reach "settled". The test still uses `lifecycleRepository`: it adds and commits `do-work/queue/REQ-141-owned-link.md`, runs `worktree new REQ-141`, then archives the REQ as `completed`. A one-line comment in the test explains why. The RED still fires for the intended reason (`dirty=true merged=true request_state=settled`).
- **D-08 (DECIDE & STATE): helper names and shapes.** The brief allowed keeping thin method callers. I kept the `readLinkConfig` method because it maps errors to the two refusal codes. I deleted the `ownedLinks` and `dirtyPaths` methods instead of keeping one-line wrappers, and the call sites use the package functions directly (delete before add). Names:
  - `readWorktreeLinkConfig(mainRoot) (linkPaths, invalidLine, readError)`: no new error type. The caller tells "unreadable" from "invalid" by which return value is set.
  - `ownedWorktreeLinks(mainRoot, worktreePath, linkPaths)`
  - `OwnedWorktreeLinks(mainRoot, worktreePath)`: the one export. It reads the config itself and returns nothing on an unreadable or invalid config.
  - `worktreeDirtyPaths(ctx, worktreePath, ownedLinks)`: the shared dirt reader, used by `worktreeClean`, `reportStatus` and `removeBuilderWorktree`.
- **D-09 (DECIDE & STATE): the status row shape is `Missing bool` (`json:"missing"`) next to `Dirty`.** A missing row has `Dirty: false` and still fills ahead/behind and the last-commit time. The text renderer prints `missing` in the clean/dirty slot.
- **D-10 (DECIDE & STATE): the board status read switched to `-z`.** The board now has to compare each path exactly against the config paths. Without `-z`, git quotes some paths (for example non-ASCII ones), and an owned link would then read as dirt. A non-empty entry shorter than 4 bytes counts as dirty, so an unreadable entry never reads as clean (`unknown-reads-as-clean`).
- **D-11 (DECIDE & STATE): the main-root derivation in Pass 5 is inlined.** It is the same one-line expression as `handleWorktree`, and I did not add a shared helper for one line.
- **D-12 (DECIDE & STATE): if removing an owned link fails in Pass 5,** it is reported as `WORKTREE-REMOVE-FAILED` (an existing code) and Pass 5 continues with the next leftover. No new finding code.

## Discovered Tasks

- impact-user-visible: Cleanup Pass 5 can never settle a builder worktree whose name has a REQ number under three digits (for example `worktree-agent-REQ-41-*`). `worktreeRequestState` compares `requestIDFromWorktree` (`REQ-41`) with the zero-padded `FilenameID` (`REQ-041`) and reports `malformed`, so such a leftover always asks for consent (`internal/cleanup/cleanup_git.go` `worktreeRequestState`). Real REQ ids are three digits or more, so the impact is low → report only
- impact-internal: `listWorktreeAgentWorktrees` in queue-kanban still parses `git worktree list --porcelain` without `-z`, so a worktree path containing a newline would be misread. This existed before this change → report only
- Review F5, F7 (`splitZero` sort) and F8 in REQ-660's review were not touched, as the brief says → report only

## Lessons read

- `_dev/primes/lessons-releases.md`: the whole file. Not relevant here, because no release payload changed.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`: family `ignored-path-is-expendable` (the REQ-660 link lesson). It explains why the link is untracked under a `dir/` ignore line, which is the RED.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`: families `paired-predicate-drift` (it is why both copies of the rule name each other in comments) and `unknown-reads-as-clean` (it is behind D-10's short-entry rule and D-04's empty filter).
- Crew rules: general, backend, coding-guardrails, shared-principles, communication-style, testing. Primes: prime-do-work-cli (Read first, Traps), prime-kanban-board (Conventions), prime-action-files, prime-releases.

## Anti-bloat check

The diff stat is above. Here is everything I added that the brief did not name, each with its reason:

- `readWorktreeLinkConfig`, `ownedWorktreeLinks`, `OwnedWorktreeLinks`, `worktreeDirtyPaths`: the brief asked for these and let me choose the names (D-08).
- `ownedWorktreeLinkPaths` (verify.go): the board's copy that D-03 requires.
- `worktreeAgentGitState.mainWorktreePath`: the brief asked to carry the main tree path through this struct.
- New `mainWorktreePath` parameter on `classifyWorktreeLeftover` and `worktreeHasUncommittedWork`, and a second return value on `listWorktreeAgentWorktrees`: the brief asked for this carry.
- The new `ownedLinks` parameter on `worktreeClean`: needed so Pass 5 reads the config once and uses the same owned set for the clean check and for removal.
- `WorktreeStatusRow.Missing`: named by the brief.
- `REQ-141-owned-link.md` written inside the test fixture: D-07.

No new flags, config, files, finding codes or remedy text. Deleted: the `lifecycleRun.ownedLinks` and `lifecycleRun.dirtyPaths` methods.

## Proposed CHANGELOG entry (the integrator adds the version)

**Builder worktree status and cleanup survive a missing folder, and dirt checks ignore the worktree command's own links**

If someone deleted a builder worktree folder by hand, `worktree status` and `worktree cleanup` stopped. Every dirt check except the worktree command itself read the `do-work/worktree-links` symlinks as uncommitted work, so finished worktrees always asked for consent.

- `do-work-cli worktree status` lists a worktree whose folder is gone as `missing` (JSON `missing: true`) and exits 0. `worktree cleanup REQ-N` prunes the record and deletes the merged branch.
- Cleanup Pass 5, the handoff survey (`handoff-state-survey`) and the board's verify check ignore owned links: configured paths that are symlinks to the same main-tree path. Pass 5 removes those links before its non-forced `git worktree remove`.
- An unreadable or invalid `do-work/worktree-links` filters nothing, so those readers still report dirt and ask for consent.
- `actions/cleanup.md` Pass 5 and `actions/fan-out-reference.md` describe both behaviors.

## Proposed lesson bullet

Satellite: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`

- [family: ignored-path-is-expendable] [REQ-693: a reader that ignores tool-created links must also remove them before a non-forced `git worktree remove`, because git refuses a worktree whose only untracked path is a symlink. Filtering the dirt check alone turned a consent finding into WORKTREE-REMOVE-FAILED. Also, after `rm -rf` of a worktree folder, `git branch -d` refuses while the prunable record exists, so prune first.](../../../../do-work/archive/UR-154/REQ-693-worktree-missing-folder-and-owned-link-dirt.md#lessons-learned)

## Integration seams and test wall times

- REQ-694 (`req append-section` and `frontmatter set` guards): same `corehelpers` package. My only new identifier there is the `cleanup` import plus the locals `mainRoot`, `statusRows`, `ownedLinks` and `dirtyRows` inside `handleHandoffSurvey`. No new package-level names, so there is no symbol clash.
- REQ-695 (run-status missing `--run`): my `result_model.go` edit stays inside `WorktreeStatusRow` (plus 2 comment lines above it) and the renderer's status-row loop. `RunStatusResult` starts 2 lines below the struct, so a merge conflict on nearby lines is possible.
- REQ-696 (stale-wording sweep): no overlap.
- Test wall times: the two new tests take 0.50s and 0.34s. The whole cleanup package takes 15s, corehelpers 10s, and the queue-kanban `Worktree|Verify` tests 9s.
