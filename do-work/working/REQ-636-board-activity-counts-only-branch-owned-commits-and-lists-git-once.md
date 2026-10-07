---
id: REQ-636
title: 'Board activity counts only commits the builder branch owns, and each served board response lists worktrees and agent branches once'
status: claimed
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
write_set: [skills/do-work-board/tools/queue-kanban/activity_correlation.go, skills/do-work-board/tools/queue-kanban/activity_correlation_test.go, skills/do-work-board/tools/queue-kanban/verify.go, skills/do-work-board/tools/queue-kanban/verify_test.go, skills/do-work-board/tools/queue-kanban/serve.go, skills/do-work-board/tools/queue-kanban/generate.go, skills/do-work/CHANGELOG.md, skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md, do-work/lessons-index.md]
claimed_at: 2026-10-07T13:36:25Z
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

## Full Context
See `do-work/user-requests/UR-137/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: review findings F3 and F7 from the pasted consumer code review of 0.305.58 to 0.305.69, triaged with `do-work validate-feedback` on 2026-10-06 and accepted.*
