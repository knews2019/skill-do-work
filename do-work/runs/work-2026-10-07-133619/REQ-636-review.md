## Review

**Overall: 94%** | 2026-10-07T13:50:03Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 93% |
| Test Adequacy | 85% |
| Scope | 97% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- M1 `skills/do-work-board/docs/board-guide.md:53` still says a card's last activity includes "its live `worktree-agent-REQ-NNN-*` branch tip"; after this REQ only a tip the branch owns (not reachable from the integration branch) counts, so the shipped guide restates the old meaning — impact-negligible → report only
- M2 `skills/do-work-board/tools/queue-kanban/verify.go:179-183` the `collectVerifyFindings` doc comment still says "since serve calls this outside its mtime cache"; serve now reaches `collectVerifyFindingsFromGitState` through `attachVerifyFindingsAndRequestActivity`, and only the `verify` subcommand calls `collectVerifyFindings` — impact-negligible → report only
- M3 no test pins D-02 (unresolvable integration ref counts no tips) or the detached-HEAD path (integration ref is a commit id): `cannedGitRunner` answers every `rev-parse` with `main`; the reviewer confirmed both paths against real git only — impact-negligible → report only
- M4 `TestServedResponseListsWorktreesAndBranchesOnce` counts only calls through the injected runner, so a second worktree listing reintroduced through `runGitCommand` passes (reviewer mutation 3: green); the shared state makes that regression unlikely but the test does not catch it — impact-negligible → report only
- M5 the REQ cites "D-06" (removal of the test-only wrappers) three times, but no decision text with that id exists in the REQ or the hand-back; the reasoning lives only in the 925955aa commit message — impact-negligible → report only
- Nit N1 the `verify` subcommand now runs one `for-each-ref --no-merged` whose owned tips it never uses (one extra git spawn per CLI verify) — impact-negligible → report only

**Acceptance:** Pass — full queue-kanban module green at 2a0b6c90 (67s), both new tests turn RED under mutation, and `readWorktreeAgentGitState` against a real scratch repo returns only the branch with its own commit (at-tip, behind, and merged branches excluded) on a branch HEAD, a detached HEAD, and no tips when the integration ref cannot resolve.
**Suggested testing:** 3 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

---

### Review detail (orchestrated report)

## Review: REQ-636

**Approve** — the board now counts a builder branch's tip only when the branch owns it, and one response reads worktrees and branches once, with a fixed command count, through the injected runner.
Route B | merge 2a0b6c90 (range c834fbd1..2a0b6c90; builder 977f811c, orchestrator 925955aa)

### What's built
- A `worktree-agent-REQ-NNN-*` branch still at the integration commit no longer adds a fake commit event to its REQ's card.
- `readWorktreeAgentGitState` (verify.go) makes one read: `worktree list --porcelain`, `branch --list`, `rev-parse --abbrev-ref HEAD` (+ `rev-parse HEAD` when detached), `for-each-ref --no-merged=<integration>`. Verify, the disk-space probe and activity share it. Four commands, five detached, regardless of branch count.

### Decisions / risks for you
- D-02 changes behavior when the integration ref cannot resolve: no branch tips count (before: every tip counted). Stamps still show. Reasonable; it is the only state where ownership cannot be shown.
- Verify's per-branch probes (`classifyWorktreeLeftover`, `worktreeCommittedQueueState`, `worktreeDirtyQueueState`) still spawn git through `exec` and still scale with branch count. Recorded by the builder as a discovered task (D-05); out of this REQ's scope.

### Requirements Checklist
- [x] R1 F3: tip counted only when not reachable from the integration ref that `resolveIntegrationBranchRef` resolves — delivered (`verify.go` `readWorktreeAgentGitState`, `--no-merged="+gitState.integrationRef`)
- [x] R2 F3 RED test (tip equals integration tip → claim stamp, kind `stamp`, no gap); distinct-own-tips fixture green — delivered (`TestServedActivityIgnoresABranchTipTheBranchDoesNotOwn`; `TestRequestActivityCountsALiveWorktreeAgentBranchTip` passes)
- [x] R3 F7 worktrees: one listing feeds `appendWorktreeFindings` and the disk-space probe; skipped-probe messages unchanged — delivered
- [x] R4 F7 branches: constant git commands, shared by verify and activity, injectable — delivered (D-01 uses one `for-each-ref --no-merged` in place of the suggested `branch --merged` pair; same count)
- [x] R5 `since` bound unchanged — delivered (diff does not touch lines computing `since`; only the `collectRequestActivity` call gains the tips argument)
- [x] R6 full module suite — delivered (reviewer run, below)
- [ ] R7 release / R8 lesson — N/A here, written by finalization
- [x] Constraint: finding wording and skipped-probe text unchanged — delivered
- [x] Constraint: single-word struct fields in activity_correlation.go untouched — delivered (only `cannedGitRunner` in the test file and the new `worktreeAgentGitState` gained fields)

### Check-by-check evidence

**C1 F3 correctness against real git.** Scratch repo, git 2.50.1: `main` with a `--no-ff` merge; branches `REQ-1-attip` (tip = an integration commit), `REQ-2-own` (one own commit), `REQ-3-behind` (ancestor of main), `REQ-4-mergedown` (own commit, merged into main).
- `for-each-ref --no-merged=main` → only `worktree-agent-REQ-2-own 2026-10-01T12:00:00Z`.
- Detached HEAD: `rev-parse --abbrev-ref HEAD` = `HEAD`; `--no-merged=<sha>` → same single line.
- Bogus ref: `fatal: malformed object name`, exit 128; the code returns with nil tips.
- Same repo through the production function (throwaway test, deleted): on main `ref="main" owned=map[REQ-2:[12:00Z]] calls=4`; detached `ref=<sha> owned=map[REQ-2:…] calls=5`; non-git dir `refErr=cannot resolve the integration branch…`, `owned=map[]`.
- The merged-branch exclusion is correct: a merged builder's commits reach the REQ through the merge-range attribution in the windowed log, not through the tip.
- `%(committerdate:iso-strict)` has the same format as the old `%cI`; `parseTimestamp` parsed it.

**C2 F7.** All four listing commands go through `gitCommandRunner` (`runGitCommand` = the same `git -C <root> … .Output()` call the old `exec` code made, plus a git-on-PATH guard). One state read per `attachVerifyFindingsAndRequestActivity`, the only production caller in generate.go and serve.go. Skipped-probe strings: sha1 of every quoted probe/failure literal in verify.go is identical at c834fbd1 and 2a0b6c90 (`bfd5c4a8…`); the only changed `Sprintf` lines swap `listError` for `gitState.worktreeListError`, which holds the same error value. Skip order is unchanged: worktree probes, then committed-queue-state, then disk-space. One harmless difference: `resolveIntegrationBranchRef` now runs even when the worktree listing failed. Its result is not reported in that case, so the output is unchanged.

**C3 `since` and fields.** `attachRequestActivity` diff: signature gains `gitState`, call passes `gitState.ownedTipInstantsById`; the `since` loop is untouched. No production struct field renamed except `liveBoardServer.activityGitRunner` → `liveGitRunner` (D-04, multi-word, justified).

**C4 Wrapper removal (925955aa).** Four call sites: generate_test.go ×3 and verify_test.go ×1 now call `attachVerifyReport(&data, board, collectVerifyFindings(board.RepoRoot, board, moment))`. That expression equals the old `attachVerifyFindings` body. generate_test.go ×1 and javascript_behavior_b_test.go ×1 pass `readWorktreeAgentGitState(root, runner)` with the same runner, which equals the old reading wrapper. No assertion changed. `git grep` at 2a0b6c90 outside `do-work/` for `attachVerifyFindings` (not followed by `A`), `liveBranchTipInstants`, `activityGitRunner`, `attachRequestActivityFromGitState`: no hits. Remaining hits are archived REQs and old hand-backs, which are historical records.

**C5 Restatement Sweep.** Redefined: what a live branch tip means as activity evidence, and how verify reads worktrees and branches. Swept skills/do-work-board (docs, tools, web), `_dev/primes`, and skills/do-work/docs for `log -1`, `branch tip`, `for-each-ref`, `tip date`, `last activity`, `liveBranch`. Stale: board-guide.md:53 (M1) and the `collectVerifyFindings` comment (M2). Current: the activity_correlation.go header, `collectRequestActivity`/`attachRequestActivity` docs, serve.go field and per-request comments, the generate.go attach comment. `prime-do-kanban.md` and `prime-kanban-board.md` do not describe the per-branch read or the double listing. model.go:1533 `log -1 … <hash>` is the unrelated commit-date lookup.

**C6 Test quality.**
- Mutation 1: drop `--no-merged` from the `for-each-ref` call. Result: `TestServedActivityIgnoresABranchTipTheBranchDoesNotOwn` fails ("last activity … (commit), want the claim stamp"; "largest gap = {Minutes:50…}"). The distinct-tips test stays green. The test pins F3.
- Mutation 2: add one `runner(... "log","-1" ...)` per branch. Result: `TestServedResponseListsWorktreesAndBranchesOnce` fails ("6 with 1 branch and 9 with 4"). The test pins the constant count.
- Mutation 3: a second worktree listing through `runGitCommand` in the disk-space path. Result: the test passes (M4).
- Both new tests name their failure in a doc comment. `cannedGitRunner` filters only on the exact `--no-merged=main` argument, so it also pins the argument form.
- Gaps: D-02 and the detached path (M3).

**Scope.** 8 files changed, all 8 declared in `## Scope`. Decisions D-01 to D-05 are in the hand-back. D-06 has no text (M5).

**Coding-guardrails spot check.** New names are multi-word and searchable. No debug output. `gofmt -l .` is empty and `go vet .` is clean at 2a0b6c90.

**P-A-U.** All three boxes `[x]`.

### Acceptance Testing
**Result: Pass**
- `go test -count=1 -run 'Activity|Serve|Verify|Worktree|DiskSpace|GeneratedVerify|ImplementationSpan' .`: ok, 8.7s.
- Full module `go test -count=1 .` in a detached worktree at 2a0b6c90: ok, 67.3s. QUEUE_KANBAN_BROWSER was not set, so browser lanes may skip. The orchestrator's run at 925955aa had it set.
- Real-git behavior: C1. Mutations: C6. The review worktree was removed with `git worktree remove`.

### Suggested Additional Testing
- S1 Open the served board in this repo during a fan-out. A freshly dispatched builder's claimed card should show its claim or dispatch stamp, not the branch-creation time, until the builder's first commit.
- S2 Add a canned-runner case where `rev-parse` fails (D-02: stamps only) and one where it returns `HEAD` (detached: `--no-merged=<sha>`).
- S3 Fix M1 and M2 in the same finalization commit if convenient. The board guide row is the user-facing statement of what "last activity" counts.

### Scores (on the record — not the headline)

**Overall: 94%** (mean of 100, 93, 85, 97 = 93.75; no modifiers)

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | R1-R6 delivered; R7-R8 belong to finalization |
| Code Quality | 93% | Clean shared read; one stale doc comment (M2), one unused CLI git call (N1) |
| Test Adequacy | 85% | Both new tests pin their failures under mutation; D-02/detached untested (M3), exec-path bypass (M4) |
| Scope | 97% | All 8 files declared; D-06 cited without decision text (M5) |
| Risk | Low | D-02 behavior change is deliberate and safe; no contract change outside the package |
| Acceptance | Pass | Full module green; real-git cases correct |

### Follow-ups created
None (6 findings report only)
