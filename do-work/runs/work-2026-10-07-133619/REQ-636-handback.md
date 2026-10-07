# REQ-636 hand-back (board activity counts only branch-owned commits; one git read per response)

- Branch: `worktree-agent-REQ-636-board-activity-counts-only-branch-owned-commits-and-lists-git-once` (no name collision)
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-636-board-activity-counts-only-branch-owned-commits-and-lists-git-once`
- Commit: `977f811c` ([REQ-636] board activity counts only branch-owned commits; one git read per response), parent `13932fc2`

## File manifest (all modified, none new)

- `skills/do-work-board/tools/queue-kanban/verify.go` — new `worktreeAgentGitState` + `readWorktreeAgentGitState(repoRoot, runner)`; `collectVerifyFindings` is now a one-line wrapper over new `collectVerifyFindingsFromGitState`; `appendWorktreeFindings` and the disk-space probe read the shared state (worktree list runs once); `listWorktreeAgentWorktrees`, `listWorktreeAgentBranches`, `resolveIntegrationBranchRef` take a `gitCommandRunner` instead of calling `exec` directly. Skipped-probe messages unchanged.
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go` — deleted `liveBranchTipInstants` (branch list + one `git log -1` per branch); `collectRequestActivity` takes `ownedTipInstantsById`; `attachRequestActivity` is a wrapper over new `attachRequestActivityFromGitState`. `since` bound untouched; single-word struct fields untouched.
- `skills/do-work-board/tools/queue-kanban/generate.go` — new `attachVerifyFindingsAndRequestActivity(data, board, now, runner)` (one state read feeding both); old `attachVerifyFindings` body moved to `attachVerifyReport`; generate calls the combined function.
- `skills/do-work-board/tools/queue-kanban/serve.go` — calls the combined function; field `activityGitRunner` renamed `liveGitRunner` (it now feeds verify too).
- `skills/do-work-board/tools/queue-kanban/activity_correlation_test.go` — `cannedGitRunner` now models worktree list, rev-parse, branch list, `for-each-ref --no-merged=main` with a `mergedBranches` set, and counts calls; two new tests; existing call sites updated for the new `collectRequestActivity` parameter.

`verify_test.go` was not changed (no edit was needed).

## Red-green evidence

RED (run before any production change; only the serve field rename was applied so the tests compiled):

```
--- FAIL: TestServedActivityIgnoresABranchTipTheBranchDoesNotOwn
    last activity = "2026-10-07T13:09:28Z" (commit), want the claim stamp "2026-10-07T12:19:28Z" (stamp) — the branch tip is the integration commit
    largest gap = {Minutes:50 FromPhase:claimed ToPhase:claimed}, want none — an unowned tip must not split the idle time
--- FAIL: TestServedResponseListsWorktreesAndBranchesOnce
    worktree listings per response = 0 and 0, want 1 each
    git commands per response = 3 with 1 branch and 6 with 4, want a constant
```

(0 worktree listings through the runner because verify called `exec` directly, which is the injectability gap the REQ names.)

GREEN: both PASS, plus every existing `TestRequestActivity*` and `TestServeReadsRequestActivityOnEveryResponse` (the distinct-own-tips fixture stays green).

Live-git confirmation in this repo: `git for-each-ref 'refs/heads/worktree-agent-*'` listed the REQ-635 builder branch and this REQ-636 branch (before my commit) both at main's tip time 16:16:11; the `--no-merged=main` form dropped exactly those two. That is the F3 bug occurring on the real queue.

## P-A-U

- [x] **[PLAN]:** Introduce one `worktreeAgentGitState` read in verify.go through the injectable `gitCommandRunner`: `worktree list --porcelain`, `rev-parse --abbrev-ref HEAD` (plus `rev-parse HEAD` only when detached), `branch --list worktree-agent-*`, and `for-each-ref --no-merged=<integration> --format='%(refname:short) %(committerdate:iso-strict)' refs/heads/worktree-agent-*`. Four commands (five detached), independent of branch count. The `--no-merged` listing is the owned-tip set (F3); verify uses the names and worktree map (F7). generate and serve call one combined attach that reads the state once. Keep `collectVerifyFindings`, `attachVerifyFindings`, `attachRequestActivity` signatures as thin wrappers, because their callers in `generate_test.go` and `javascript_behavior_b_test.go` are outside the write boundary. Tests first through the served response.
- [x] **[APPLY]:** Implemented as planned in the five files above. No scope added.
- [x] **[UNIFY]:** `git diff --stat`: 5 files, +221/-100. `gofmt -l .` empty, `go vet .` clean, no debug output or TODOs in the diff. Verified per file: verify.go (skipped-probe texts byte-identical, finding wording untouched, order of skip messages preserved), activity_correlation.go (since bound unchanged, struct fields unchanged, header comment updated to the new read), generate.go (suppression list and path reduction still have one home, now `attachVerifyReport`), serve.go (comment updated for the renamed field), test file (canned runner no longer answers `log -1`, which production no longer calls).

## Decisions

- D-01 DECIDE & STATE: owned tips come from `for-each-ref --no-merged=<integration>` (one command giving names and dates of unowned-filtered branches) instead of the suggested for-each-ref + `branch --merged` pair. Same command count, and the filter is done by git, not by set difference in Go.
- D-02 DECIDE & STATE: when the integration ref cannot be resolved, no branch tips are counted. Ownership cannot be shown, and the stamps remain true evidence. The previous code counted every tip.
- D-03 ESCALATE (low): kept `collectVerifyFindings`, `attachVerifyFindings`, `attachRequestActivity` as one-line wrappers so `generate_test.go` and `javascript_behavior_b_test.go` (outside my write boundary) need no edits. After this change `attachVerifyFindings` and `attachRequestActivity` are called only by tests. Value of collapsing: three fewer functions. Risk: none functional; it is a six-call-site mechanical edit in two test files, fully reversible. Recommendation: collapse in a follow-up only if the maintainer dislikes test-only wrappers.
- D-04 DECIDE & STATE: renamed `liveBoardServer.activityGitRunner` to `liveGitRunner`, because it now feeds the verify listings too and the old name would mislead.
- D-05 DECIDE & STATE: the count test asserts git commands made through the runner. Per-branch verify commands that still use `exec` directly (`merge-base`/`classifyWorktreeLeftover`, `git diff` for committed queue state, `git status` per worktree) are findings work, not listing, and are out of this REQ's scope.

## Discovered Tasks

- Verify's per-branch probes (`classifyWorktreeLeftover`, `worktreeCommittedQueueState`, `worktreeDirtyQueueState`) still spawn `git` via `exec` directly, so they are not test-injectable and still scale with branch count. impact-low → report only

## Lessons read

- `_dev/primes/lessons-releases.md` (whole)
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`, only the `[family: git-history-evidence]` bullet (0.305.67, merge commits and pathspec)
- Primes: `_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`, `_dev/primes/prime-releases.md`
- No listed file was missing.

## Proposed lesson entry

Target: `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`, family `git-history-evidence` (second bullet of this family, so the general.md rule makes a prime `## Traps` promotion mandatory; suggested trap line below).

> - [family: git-history-evidence] 0.305.NN: **a branch tip is evidence for a REQ only when the branch owns that commit.** REQ-632 added each live `worktree-agent-REQ-NNN-*` branch's `git log -1` date to that REQ's activity. A builder branch is cut from the integration tip, so before its first commit the tip IS the integration commit, and the board reported that commit as the REQ's last activity and split its idle gap with it. On this repo during the REQ-635 to REQ-638 batch, two of four live builder branches pointed at main's tip with no commits of their own. Read tips with `git for-each-ref --no-merged=<integration ref>`: a tip reachable from integration is either unowned or already merged, and a merged builder's commits already reach the REQ through the merge's `^1..^2` range.

Suggested prime trap (`prime-do-kanban.md` § Traps): `- [family: git-history-evidence] Commit evidence for a REQ must be commits that REQ's work owns: never narrow the log with a pathspec (merges vanish), and never count a branch tip reachable from the integration branch.`

## Proposed CHANGELOG entry

Title: **Board Activity Ignores Builder Branches With No Commits of Their Own**

> A claimed card's "last activity" no longer jumps to the moment a builder branch was created. A `worktree-agent-*` branch now counts as activity for its REQ only when its tip is a commit the branch owns, not the integration commit it was cut from, so the idle gap on a fresh claim stays honest. Each served board response also reads the worktree list and the builder branches once, shared by the verify band and the activity line, with a fixed number of git commands however many branches exist.

## Integration seams

None. No signature used outside `skills/do-work-board/tools/queue-kanban/` changed. `generate_test.go` and `javascript_behavior_b_test.go` compile unchanged.

## Test commands and results

- `go test -run 'TestServedActivityIgnoresABranchTipTheBranchDoesNotOwn|TestServedResponseListsWorktreesAndBranchesOnce|TestRequestActivity|TestServeReadsRequestActivity' .` — RED as above before the change, PASS after.
- `go vet .` — clean. `gofmt -l .` — empty.
- Full module: `go test -count=1 -json .` in `skills/do-work-board/tools/queue-kanban` — PASS, 0 failures, package 83.5s wall (1:23.85 real). Per-file summed test time, all under 30s: strict_behavior_regression_test.go 24.65s, generate_test.go 11.76s, verify_test.go 6.80s, citations_test.go 4.57s, timeline_browser_probe_test.go 2.23s, browser_probe_test.go 0.84s, durations_browser_probe_test.go 0.68s, board_live_test.go 0.57s, user_request_progress_browser_probe_test.go 0.57s, serve_test.go 0.28s, activity_correlation_test.go 0.15s, testing_test.go 0.15s, every other file under 0.11s. QUEUE_KANBAN_BROWSER was not set, so browser lanes may have skipped; this REQ changes no `web/` file.
- Post-commit re-run of `-run 'Activity|Serve|Verify|Worktree|DiskSpace'` — ok, 15.3s.
