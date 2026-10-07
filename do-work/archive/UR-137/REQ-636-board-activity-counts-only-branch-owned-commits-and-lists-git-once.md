---
id: REQ-636
title: 'Board activity counts only commits the builder branch owns, and each served board response lists worktrees and agent branches once'
status: completed
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  basis:
  - Route B
  - 10-file write set
  - 6 acceptance criteria
  calculated_at: 2026-10-07T13:36:48Z
created_at: 2026-10-06T22:45:18Z
user_request: UR-137
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-635, REQ-637, REQ-638]
batch: review-0305-69-findings
required_lessons: [_dev/primes/lessons-releases.md]
write_set: [skills/do-work-board/tools/queue-kanban/activity_correlation.go, skills/do-work-board/tools/queue-kanban/activity_correlation_test.go, skills/do-work-board/tools/queue-kanban/verify.go, skills/do-work-board/tools/queue-kanban/verify_test.go, skills/do-work-board/tools/queue-kanban/serve.go, skills/do-work-board/tools/queue-kanban/generate.go, skills/do-work-board/tools/queue-kanban/generate_test.go, skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go, skills/do-work-board/docs/board-guide.md]
claimed_at: 2026-10-07T13:36:25Z
kb_status: pending
dispatch_at: 2026-10-07T13:16:25Z
builder_handback_at: 2026-10-07T13:22:50Z
review_at: 2026-10-07T13:54:49Z
integration_at: 2026-10-07T13:42:46Z
commit: 81524052bafc3baa2dae2a89adc983c39e05b8bd
heavy_verified_at: 2026-10-07T13:56:51Z
heavy_verified_revision: 81524052bafc3baa2dae2a89adc983c39e05b8bd
completed_at: 2026-10-07T13:57:15Z
release_at: 2026-10-07T13:57:15Z
---
# Board Activity Counts Only Commits the Builder Branch Owns, and Each Served Board Response Lists Worktrees and Agent Branches Once

## What
Two related fixes in the queue-kanban board's per-response git reads. Review finding F3: `liveBranchTipInstants` (`skills/do-work-board/tools/queue-kanban/activity_correlation.go:186-207`) runs `git log -1 --format=%cI <branch>` per `worktree-agent-*` branch, so a branch with no commits of its own reports the integration commit it was cut from as a commit event for that REQ. Count only commits the branch owns. Review finding F7: each served board response lists worktrees twice (`verify.go:1281` via `appendWorktreeFindings`, and `verify.go:209` for the disk-space probe) and lists agent branches twice (`verify.go:1287` and `activity_correlation.go:187`) followed by one `git log -1` per branch. List each once and share the result.

## Why
The fake commit event splits a real idle gap and can set "last activity" on a claimed card, so the card lies about when work last happened. The duplicated listing is cleanup, not a performance bug: the page reloads manually and the whole probe set measured about 40ms (`serve.go:160-174`). Both touch the same branch listing, and the review notes that if F7 is done with F3 the combined read must still exclude tips the branch does not own, so they are one slice.

## Detailed Requirements
1. F3: a `worktree-agent-REQ-NNN-*` branch contributes a commit instant only when that commit is unique to the branch, that is, not reachable from the integration branch. The integration ref is what `resolveIntegrationBranchRef` (`verify.go:1388`) already resolves for the worktree probes.
2. F3 test, written first and confirmed failing: a branch whose tip equals the integration tip (no own commits) contributes no commit event, so a claimed REQ with only a `claimed_at` stamp shows that stamp as its last activity and no idle gap is split. The existing fixture at `activity_correlation_test.go:177-180` (branches with distinct own tips) stays green.
3. F7 worktrees: run `git worktree list --porcelain` once per `collectVerifyFindings` call and pass the map to both `appendWorktreeFindings` and the disk-space probe (`appendDiskSpaceFindings`). The skipped-probe messages for an enumeration failure stay as they are.
4. F7 branches: replace the two branch listings plus N `git log -1` calls with a constant number of git commands. The suggested shape: one `git for-each-ref --format='%(refname:short) %(committerdate:iso-strict)' refs/heads/worktree-agent-*` for names and tip dates, plus one `git branch --merged <integration> --list 'worktree-agent-*'` to drop branches whose tip is already reachable from integration. Two commands total instead of two plus one per branch. Verify and activity share the result; verify currently calls `exec` directly while activity uses the injected `gitCommandRunner`, so the shared listing must still be test-injectable.
5. Leave the `since` bound in `attachRequestActivity` (`activity_correlation.go:339-349`) exactly as it is.
6. Run the full Go test suite for the queue-kanban module.
7. Release per `_dev/primes/prime-releases.md` with a CHANGELOG entry.
8. Lessons: one entry in `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (a branch tip is evidence for the REQ only when the branch owns the commit; the `git-history-evidence` family is the fit) and the matching token refresh in `do-work/lessons-index.md`, in the same commit.

## Constraints
- Do not change what counts as a verify finding or its wording; this REQ changes how the inputs are read, not what is reported.
- Do not touch the single-word private struct fields in `activity_correlation.go`; the review rejected that as a style preference.
- Batch constraint: each REQ in this batch is its own release and its own commit.

## Dependencies
None. Independent of REQ-635, REQ-637 and REQ-638, which edit different files.

## Builder Guidance
Certainty is high for F3 and the worktree half of F7. For the branch half, the two-command shape above is a suggestion; any shape with a constant number of git commands that still excludes unowned tips and remains injectable in tests is acceptable. Plumbing the shared listing from `collectVerifyFindings` into `attachRequestActivity` through `serve.go` and `generate.go` is the main design decision; keep it small.

## Red-Green Proof
**RED prompt/case:** Through the injected `gitCommandRunner`, a claimed REQ-701 with `claimed_at` one hour ago, a branch `worktree-agent-REQ-701-x` whose tip is the integration commit made ten minutes ago, and no correlated commits in the windowed log.
**Why RED now:** The card reports a commit event ten minutes ago as REQ-701's last activity, taken from a commit the branch does not own.
**GREEN when:** REQ-701's last activity is its `claimed_at` stamp one hour ago, kind `stamp`, with no idle gap; and a single served board response issues one worktree listing and a constant number of branch-related git commands regardless of branch count (asserted by counting runner calls in the test).
**Validation:** Inferred during capture. Confirmed by reading `activity_correlation.go:198` and `verify.go:202-215, 1281-1287` during triage; not executed against a live repository.

## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (7253 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matched: family `git-history-evidence` governs commit evidence attributed to a REQ, which is this REQ's subject.
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (5912 tokens, over the remaining 1334 after `lessons-releases.md`; `slugged: partial`, so no targeted form). Matched at claim: its row covers queue-kanban static output, which `generate.go` produces; this REQ changes how inputs are read, not the output.

## Full Context
See `do-work/user-requests/UR-137/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** One `worktreeAgentGitState` read in verify.go through the injectable `gitCommandRunner`: `worktree list --porcelain`, `rev-parse --abbrev-ref HEAD` (plus `rev-parse HEAD` only when detached), `branch --list worktree-agent-*`, and `for-each-ref --no-merged=<integration>` for the owned tips: four commands (five detached) whatever the branch count. generate and serve call one combined attach that reads the state once. Tests first through the served response. (Builder, from the hand-back.)
- [x] **[APPLY]:** Implemented as planned in five files at 977f811c; the orchestrator then removed the test-only wrappers in 925955aa (D-06), touching three more test files. All eight files are in the Scope.
- [x] **[UNIFY]:** `git diff --stat` 8 files, +218/-109 over c834fbd1..2a0b6c90. `gofmt -l .` empty and `go vet .` clean after each commit; no debug output or TODOs in the diff. Checked per file: verify.go (skipped-probe texts byte-identical, finding wording untouched, skip-message order kept), activity_correlation.go (since bound unchanged, struct fields unchanged, header comment updated), generate.go (the suppression list and path reduction have one home, attachVerifyReport), serve.go (field comment updated for the rename), the four test files (the canned runner no longer answers `log -1`; call sites of the removed wrappers).
*Source: review findings F3 and F7 from the pasted consumer code review of 0.305.58 to 0.305.69, triaged with `do-work validate-feedback` on 2026-10-06 and accepted.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The REQ names both defects, their lines and the intended shape (owned tips only, one listing per response), so no Plan agent is needed. The open design choice is how one shared listing reaches both verify and activity through `generate.go` and `serve.go` while staying injectable in tests, so exploration locates every caller first.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The build already existed when this run claimed the REQ (a parallel builder worked it on its branch; hand-back `do-work/runs/work-2026-10-07-133619/REQ-636-handback.md`), so exploration was taken from that hand-back and its real file list, checked against the source:

- **Required lessons consult:** `_dev/primes/lessons-releases.md` read whole (666 tokens). `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` stays dropped for budget, but its one `git-history-evidence` bullet (0.305.67, merges and pathspecs) was read because it governs commit evidence for a REQ. `_dev/primes/lessons-kanban-board.md` is a new budget drop (above). No listed file was missing.
- **Where the change goes:** `activity_correlation.go` `liveBranchTipInstants` was the second branch listing plus one `git log -1` per branch (F3 and F7). `verify.go` listed worktrees in `appendWorktreeFindings` and again for the disk-space probe, and listed branches and resolved the integration ref with `exec` directly, so none of it was injectable. `generate.go` and `serve.go` are the only production callers of the verify attach and the activity attach.
- **Pattern to follow:** the injected `gitCommandRunner` from REQ-632 (`activity_correlation.go`), which serve's tests already replace with `cannedGitRunner`.
- **Test-only callers:** after the shared read, `attachVerifyFindings` is called only by `generate_test.go` (3 sites) and `verify_test.go` (1 site), and the reading wrapper `attachRequestActivity` only by `generate_test.go` (1) and `javascript_behavior_b_test.go` (1). `collectVerifyFindings` keeps a production caller (`runVerifyProbes`, the `verify` subcommand).
- **Live evidence of F3:** on this repo the REQ-635 builder branch and this REQ's branch both pointed at main's tip with no commits of their own; `for-each-ref --no-merged=main` dropped exactly those two.

*Generated from the builder hand-back by the work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go` (modify) — the per-branch tip read is deleted; collectRequestActivity takes the owned tip instants; attachRequestActivity takes the shared read
- `skills/do-work-board/tools/queue-kanban/activity_correlation_test.go` (modify) — the canned runner models the shared read and counts calls; the RED pair (an unowned tip, git commands per response)
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — worktreeAgentGitState read once through the injected runner; the worktree probes and the disk-space probe read it
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modify) — one call site of the removed attachVerifyFindings wrapper (D-06)
- `skills/do-work-board/tools/queue-kanban/serve.go` (modify) — the combined attach per response; runner field renamed
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — the combined attach; attachVerifyReport keeps the suppression list
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — call sites of the removed wrappers (D-06)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go` (modify) — one call site of the removed wrapper (D-06)
- `skills/do-work-board/docs/board-guide.md` (modify) — the last-activity row names the owned-tip rule (review finding M1, D-07)

**Files I will NOT touch:** the verify finding wording and the skipped-probe messages, the since bound in the activity collector, the single-word struct fields in activity_correlation.go, the per-branch verify probes that still call exec (merge-base, diff, status; discovered task). Written by finalization in the main tree, not on the branch: every release path, skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md (the mandatory Traps promotion) and do-work/lessons-index.md.

**Acceptance criteria (restated from REQ):**
- [ ] A worktree-agent branch whose tip is reachable from the integration branch adds no commit event: a claimed REQ with only its claimed_at stamp shows that stamp as last activity, kind stamp, and no idle gap is split; the distinct-own-tips fixture stays green
- [ ] One `git worktree list --porcelain` per served response feeds both the worktree probes and the disk-space probe; the skipped-probe messages are unchanged
- [ ] The branch-related git commands per served response are a constant number whatever the branch count, shared by verify and activity, and injectable in tests (asserted by counting runner calls)
- [ ] The `since` bound in the activity collector is unchanged, and no verify finding or its wording changes
- [ ] The full queue-kanban Go test suite passes
- [ ] Release with a CHANGELOG entry, and one git-history-evidence lesson entry with its index token refresh

## Pre-Flight

**Git:** ✓ Integration tip 5039bb96 on `main` (this REQ's claim commit, after REQ-635's finalization e04509e4); the only dirt is this REQ's own trail (working REQ, run directory), no third-party paths
**Tests baseline:** ✓ focused board tests green (`do-work/runs/work-2026-10-07-133619/REQ-636-probe.sh`, launched by advance)
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exit 0 at 5039bb96 from the detached checkout `.git/work-run-2026-10-07-133619/drain-head`, both Go stages EXECUTING, gate wall 140s; green-gate record satisfied
**Dependencies:** ✓ Go toolchain present; no new module dependency

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/docs/board-guide.md` (modified)
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go` (modified)
- `skills/do-work-board/tools/queue-kanban/activity_correlation_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/serve.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified)

**What was done:** A served board response, and the static generate, now read the worktree-agent worktrees and branches once (`readWorktreeAgentGitState` in verify.go): the worktree list, the integration ref, the branch list, and one `git for-each-ref --no-merged=<integration>` that gives the tip dates of only the branches that own their tip. The verify worktree probes, the disk-space probe and the request-activity collector all read that one state through the injected `gitCommandRunner`, so the per-branch `git log -1` loop (`liveBranchTipInstants`) is deleted and the git command count no longer grows with the branch count. A builder branch still at the commit it was cut from no longer adds a commit event to its REQ. When the integration ref cannot be resolved, no tips count (D-02). The test-only wrappers `attachVerifyFindings` and the old reading `attachRequestActivity` are removed; tests call `attachVerifyReport` over `collectVerifyFindings`, and `attachRequestActivity` takes the shared state. After review, the board guide's last-activity row and the `collectVerifyFindings` comment were brought in line (review M1, M2). Merge range c834fbd1..81524052 (builder commit 977f811c, orchestrator commits 925955aa and c01a9f1f, first merge 2a0b6c90, second merge 81524052).

## Qualification

**Diff range:** c834fbd1..2a0b6c90 (builder commit 977f811c, orchestrator commit 925955aa, merge 2a0b6c90)
**Gate records:** qualify satisfied; scope-drift satisfied (the eight changed files equal the declared Scope).
**Warnings judged:** none raised.
**Orchestrator read of the diff:** every Detailed Requirement traces to the diff. 1: owned tips come from `for-each-ref --no-merged=<integration ref>`, and the integration ref is the one `resolveIntegrationBranchRef` resolves for the worktree probes. 2: `TestServedActivityIgnoresABranchTipTheBranchDoesNotOwn` is the captured RED case through the served response, and `TestRequestActivityCountsALiveWorktreeAgentBranchTip` (distinct own tips) stays green. 3: one worktree listing feeds `appendWorktreeFindings` and the disk-space probe, and both skipped-probe messages keep their text. 4: four branch-related commands (five when detached) whatever the branch count, all through the injected runner, shared by verify and activity. 5: the since bound is untouched. 6: the full queue-kanban suite passed at 925955aa (65s, QUEUE_KANBAN_BROWSER set). 7 and 8 are finalization's (release, lesson, index refresh).
**P-A-U honesty:** the three boxes were ticked by the orchestrator from the builder's hand-back and the orchestrator's own commit; APPLY cross-checked against git diff --stat c834fbd1..2a0b6c90 (8 files, all in Scope).
**After review:** the post-review fix c01a9f1f (board-guide.md row and one verify.go comment, no code) was merged with the same pre as 81524052; qualify's record stays on the first range, and the cumulative range c834fbd1..81524052 has 9 files, all in the updated Scope.
**Live data flow:** generate.go and serve.go call `attachVerifyFindingsAndRequestActivity`, which makes the one read and passes it to `collectVerifyFindingsFromGitState` and `attachRequestActivity`; the `verify` subcommand still reaches the same read through `collectVerifyFindings` with the real runner.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree at 2a0b6c90 (detached checkout `.git/work-run-2026-10-07-133619/drain-head`)
**Result:** ✓ All passing — exit 0 on the first run, gate wall 148s; stage queue-kanban-fast-tests (423 tests, wall 48s, slowest file strict_behavior_regression_test.go 21.73s < 30s); stage do-work-cli-fast-tests (876 tests, wall 70s, slowest file finalization_recovery_test.go 21.18s < 30s). Green-gate record satisfied by advance.

**Repository gate after the review fix:** `bash _dev/tests/maintainer-verify.sh` at 81524052 (same detached checkout), run directly because advance was past the test gate: exit 0, gate wall 154s; queue-kanban-fast-tests 423 tests, wall 46s, slowest file 21.53s; do-work-cli-fast-tests 876 tests, wall 75s, slowest file finalization_recovery_test.go 24.19s < 30s.

**Focused tests:** `do-work/runs/work-2026-10-07-133619/REQ-636-probe.sh` (board Activity, Serve, Verify, Worktree and DiskSpace tests) → exit 0, advance probe record satisfied. Full queue-kanban module at 925955aa with QUEUE_KANBAN_BROWSER set: ok, 65s.

**Red-green validation:** traced to `## Red-Green Proof`; RED taken by the builder before any production change (only the serve field rename applied so the tests compiled), GREEN at 977f811c and again at the merge:
- TestServedActivityIgnoresABranchTipTheBranchDoesNotOwn (the captured RED case through the served response): ✗ last activity was the branch tip 50 minutes after the claim, kind commit, and it split a 50-minute gap → ✓ the claim stamp, kind stamp, no gap
- TestServedResponseListsWorktreesAndBranchesOnce: ✗ 0 worktree listings through the runner (verify called exec directly) and 3 git commands with 1 branch against 6 with 4 → ✓ one worktree listing each, the same command count for 1 and 4 branches
- TestRequestActivityCountsALiveWorktreeAgentBranchTip (distinct own tips, REQ-632): ✓ before and after

**New tests added:**
- the two served-response tests above, in `activity_correlation_test.go`

**Existing tests updated (cross-REQ impact):**
- activity_correlation_test.go (REQ-632): the canned runner models the shared read and no longer answers `log -1`; collectRequestActivity call sites pass the owned tips — intentional
- generate_test.go, verify_test.go, javascript_behavior_b_test.go: call sites of the removed test-only wrappers (D-06) — intentional, no assertion changed

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: c834fbd1..2a0b6c90, replanned for c834fbd1..81524052 after the review fix (same three lanes; board-guide.md adds a staged-skills reason)
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — files changed under skills/do-work-board/tools/queue-kanban
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — files changed under skills/do-work-board/tools/queue-kanban
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed

*Verified by work action*

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
- M1 `skills/do-work-board/docs/board-guide.md:53` still says a card's last activity includes "its live `worktree-agent-REQ-NNN-*` branch tip"; after this REQ only a tip the branch owns (not reachable from the integration branch) counts, so the shipped guide restates the old meaning — impact-negligible → report only (fixed before finalization in c01a9f1f, merged as 81524052)
- M2 `skills/do-work-board/tools/queue-kanban/verify.go:179-183` the `collectVerifyFindings` doc comment still says "since serve calls this outside its mtime cache"; serve now reaches `collectVerifyFindingsFromGitState` through `attachVerifyFindingsAndRequestActivity`, and only the `verify` subcommand calls `collectVerifyFindings` — impact-negligible → report only (fixed before finalization in c01a9f1f)
- M3 no test pins D-02 (unresolvable integration ref counts no tips) or the detached-HEAD path (integration ref is a commit id): `cannedGitRunner` answers every `rev-parse` with `main`; the reviewer confirmed both paths against real git only — impact-negligible → report only
- M4 `TestServedResponseListsWorktreesAndBranchesOnce` counts only calls through the injected runner, so a second worktree listing reintroduced through `runGitCommand` passes (reviewer mutation 3: green); the shared state makes that regression unlikely but the test does not catch it — impact-negligible → report only
- M5 the REQ cites "D-06" (removal of the test-only wrappers) three times, but no decision text with that id exists in the REQ or the hand-back; the reasoning lives only in the 925955aa commit message — impact-negligible → report only (resolved: D-06 is now written under ## Decisions)
- Nit N1 the `verify` subcommand now runs one `for-each-ref --no-merged` whose owned tips it never uses (one extra git spawn per CLI verify) — impact-negligible → report only

**Acceptance:** Pass — full queue-kanban module green at 2a0b6c90 (67s), both new tests turn RED under mutation, and `readWorktreeAgentGitState` against a real scratch repo returns only the branch with its own commit (at-tip, behind, and merged branches excluded) on a branch HEAD, a detached HEAD, and no tips when the integration ref cannot resolve.
**Suggested testing:** 3 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

Full reviewer report: `do-work/runs/work-2026-10-07-133619/REQ-636-review.md`.

## Decisions

Builder decisions, copied from the hand-back (`do-work/runs/work-2026-10-07-133619/REQ-636-handback.md`).

- D-01 (DECIDE & STATE): owned tips come from one `for-each-ref --no-merged=<integration>` (names and dates of the branches that own their tip) instead of the suggested for-each-ref plus `branch --merged` pair. Same command count, and git does the filtering instead of a set difference in Go.
- D-02 (DECIDE & STATE): when the integration ref cannot be resolved, no branch tips count. Ownership cannot be shown, and the stamps remain true evidence. The previous code counted every tip.
- D-03 (ESCALATE, resolved by the orchestrator in D-06): the builder kept collectVerifyFindings, attachVerifyFindings and attachRequestActivity as one-line wrappers so two test files outside its write boundary needed no edits. Value of collapsing: fewer functions and no test-only production code. Risk: none functional; a mechanical edit of six call sites, fully reversible.
- D-04 (DECIDE & STATE): `liveBoardServer.activityGitRunner` is renamed `liveGitRunner`, because it now feeds the verify listings too.
- D-05 (DECIDE & STATE): the count test asserts git commands made through the runner. The per-branch verify probes that still call exec directly (merge-base, the committed queue-state diff, status per worktree) are findings work, not listing, and stay out of scope.

Orchestrator decisions:
- D-06 (DECIDE & STATE, scope extension): the test-only wrappers attachVerifyFindings and the reading attachRequestActivity are deleted in 925955aa, and attachRequestActivityFromGitState takes the name attachRequestActivity. Reason: after the shared read they had no production caller, and the maintainer's rule is delete before add and YAGNI. collectVerifyFindings stays because the verify subcommand (runVerifyProbes) calls it. Scope and write_set gained verify_test.go (already captured), generate_test.go and javascript_behavior_b_test.go. No assertion changed.
- D-07 (DECIDE & STATE, scope extension): review findings M1 and M2 were fixed before finalization in c01a9f1f, because the board guide is the user-facing statement of what "last activity" counts, and a shipped doc that says any branch tip counts would be wrong after this release. Scope and write_set gained `skills/do-work-board/docs/board-guide.md`. The gate and the three heavy lanes re-ran at the second merge 81524052.
- D-08 (DECIDE & STATE): dispatch_at (13:16:25Z, the branch creation in the reflog) and builder_handback_at (13:22:50Z, the builder commit) come from the branch, because the coordinator dispatched the builder before this run claimed the REQ at 13:36:25Z. They therefore precede claimed_at and the estimate. No builder-work timing event was recorded, because the recorder times from a start instant to now and would charge this run's own work to the builder.
- D-09 (DECIDE & STATE): the remaining review findings (M3, M4, N1) stay report only per their impact tokens; none is critical.

## Discovered Tasks

From the builder's hand-back and the review, impact-stamped per review-work Step 10; none is impact-critical, so nothing was queued.

- Verify's per-branch probes (classifyWorktreeLeftover, worktreeCommittedQueueState, worktreeDirtyQueueState) still spawn git through exec directly, so they are not test-injectable and still scale with the branch count. — impact-negligible → report only
- No test pins D-02 (an unresolvable integration ref counts no tips) or the detached-HEAD path (`--no-merged=<commit id>`); the reviewer confirmed both against real git. — impact-negligible → report only
- The call-count test sees only calls through the injected runner, so a listing reintroduced through runGitCommand would pass. — impact-negligible → report only
- The verify subcommand runs one `for-each-ref --no-merged` whose owned tips it never uses. — impact-negligible → report only

## Lessons Learned

**What worked:** Asking git for ownership directly (`for-each-ref --no-merged=<integration>`) instead of listing branches and subtracting a merged set in Go: one command gives names, dates and the filter, and it behaves the same on a detached checkout. Testing through the served response with a counting canned runner pinned both defects at the seam the user sees.
**What didn't:** Keeping thin wrappers so out-of-scope test files would compile left two production functions with only test callers; they were removed in a follow-up commit (D-06). The first merge also shipped a board guide that still described the old rule, caught only by the review's restatement sweep (D-07).
**Worth knowing:** A freshly cut builder branch points at the integration commit until its first commit, so any per-branch evidence read has to ask whether the branch owns its tip. The new lesson is the second `git-history-evidence` bullet in `lessons-do-kanban.md`, so the generalized trap was promoted into `prime-do-kanban.md` § Traps. Verify's per-branch probes still call exec directly and are not injectable (discovered task).

## Orientation

Now a claimed card's last activity ignores a builder branch that has no commits of its own, and each board response reads the worktree-agent worktrees and branches once, with a fixed number of git commands, shared by the VERIFY band and the activity line; lives in the queue-kanban verify probes and request-activity collector (`_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`). No map change. Prime spot-check: no prime names the removed functions or the per-branch tip read; `prime-do-kanban.md` § Traps gained the git-history-evidence line.

## Heavy Verification Plan

- Base revision: c834fbd1e1e523d0a6cba0fe105fe3d3e2cd84f6
- Target revision: 81524052bafc3baa2dae2a89adc983c39e05b8bd (landed in `commit:`)
- queue-kanban-javascript — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — board files changed
- queue-kanban-browser — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — board files changed
- staged-skills — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed (board code and board-guide.md)

## Heavy Verification Result

- Target revision: 81524052bafc3baa2dae2a89adc983c39e05b8bd
- Execution revision: 81524052bafc3baa2dae2a89adc983c39e05b8bd (detached drain checkout `.git/work-run-2026-10-07-133619/drain-head`, run 13:54:16Z to 13:56:26Z, QUEUE_KANBAN_BROWSER set to Google Chrome)
- queue-kanban-javascript: exit 0, executed, 8s
- queue-kanban-browser: exit 0, executed (not skipped), 88s
- staged-skills: exit 0, executed, 33s

Green: every selected lane present, exit 0, none skipped, none reused. The same three lanes were also green at the first merge 2a0b6c90 (8s, 89s, 38s), before the review fix.

## Timing

Observed 2026-10-07T13:36:25Z to 2026-10-07T13:56:58Z: 20m 33s total, 25m 29s attributed across 6 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 8m 50s | 1 |
| exploration-preflight | 5m 49s | 1 |
| verification-gate | 5m 49s | 2 |
| remediation | 4m 35s | 1 |
| handback-merge | 26s | 1 |

Slowest stage: review / independent review agent, report, 8m 50s, outcome success.
