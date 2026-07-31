---
id: UR-007
title: Encourage parallel REQ execution — write-set gate, worktree dispatch, decomposition nudge
created_at: 2026-07-28T20:52:05Z
requests: [REQ-032, REQ-033, REQ-034]
word_count: 210
---

# Encourage parallel REQ execution

## Summary

Across a conversation reviewing a consumer project's queue (g1w-game-find-the-difference), the user asked how to make do-work runs more parallel and whether the skill itself can be updated to encourage parallelism. The user initially proposed timed per-file lock files (~3 min TTL); discussion converged instead on declared write-sets with disjointness-based scheduling, orchestrator-managed worktree builders (a pattern the user's runner already demonstrated: `worktree-agent-*` branches merged by the orchestrator), serial-only resource classes for things like DB migrations, post-merge verification, worktree cleanup ownership (orchestrator at archive; crash recovery + cleanup for orphans), and a capture-time decomposition nudge. The user approved the three-REQ split and said "run the capture."

## Extracted Requests

| REQ | Title |
|-----|-------|
| REQ-032 | Write-set contract, parallel-dispatch gate, serial-only resource classes |
| REQ-033 | Worktree dispatch mode: orchestrator-managed builders, integration, cleanup |
| REQ-034 | Capture-time decomposition nudge + board write-set overlap badges |

## Batch Constraints

- Design for the floor: the existing serial loop must remain the default and fully functional for agents without subagent/worktree support; everything here is an advanced-harness mode that degrades gracefully.
- No new action files — changes land in existing actions (`work.md`, `work-reference.md`, `capture.md`, `capture-reference.md`) and the board tool; SKILL.md's word budget applies.
- Any new frontmatter field must update `tools/queue-kanban/model.go` in the same commit (parser lock-step rule).
- Timed per-file locks were considered and **rejected** (TTL breaks over live slow agents — same defect class as the 15s orchestrator-mutex break removed in 0.140.4). Do not reintroduce them.
- Serial-only resource lists must follow the Closed Enumerations Go Stale rule: state the condition, mark examples illustrative.

## Full Verbatim Input

> can we do something to update this skill to encourage them to run in parallel?
>
> For example editing .css files can we lock the file first (so for example writing a lock file, who locked, when, and then assume that it will write it in 3 minutes max) this way they can run in parallel especially because the playwright/verification phase can take long enough as it is. The only problem I can see if there is db migration and two workers try to do it without a correlation plan, but with dynamic workflows those can be mapped as well, and they can find ways to signals eachother and not step on eachother's toes.

[Screenshots: consumer project's Kanban board showing wave-2 REQs 1261–1264 ready with identical `NEEDS REQ-1258/1259/1260` rows; runner transcript confirming four parallel builders with single-writer discipline.]

> so we also have worktrees as well:

[Screenshot: lazygit showing four `worktree-agent-<hash>` branches and `wip: … (integrated by orchestrator)` merge commits; `do-work/working/` holding REQ-1262/1263.]

> who will delete the merged worktrees?

> run the capture

Assistant's user-approved plan (context for the three REQs): (A) `write_set` frontmatter + parallel-dispatch gate on pairwise-disjoint write-sets + serial-only resource classes; (B) worktree dispatch mode — `worktree-agent-REQ-NNN-<suffix>` naming, orchestrator as sole writer to main, post-merge verification before archive, archive-time `git worktree remove` + `git branch -d` (the `-d` doubling as a merge assertion), crash-recovery sweep for merged orphans, consent-gated `cleanup` deletion for unmerged ones, live-session guard, do-work state stays in the main tree; (C) capture-time slicing guidance preferring per-file REQ boundaries + board badges surfacing write-set overlaps between queued REQs.

---
*Captured: 2026-07-28T20:52:05Z*
