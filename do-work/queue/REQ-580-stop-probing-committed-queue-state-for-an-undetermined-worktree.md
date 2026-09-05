---
id: REQ-580
title: 'Stop probing committed queue state for a worktree whose merge state is already undetermined'
status: pending
created_at: 2026-09-05T00:19:58Z
user_request: UR-118
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
related: [REQ-579]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
---

# Stop Probing Committed Queue State for a Worktree Whose Merge State Is Already Undetermined

## What

In `appendWorktreeFindings` (verify.go), a leftover worktree whose branch is gone produces two messages: a `worktree-merge-state-undetermined` finding, and a skipped-probe line saying the committed-queue-state probe failed with "no such branch". Both come from one fact. When the merge-state classification is already undetermined for a leftover, skip the committed-queue-state probe for that leftover and say inside the undetermined finding that the queue-state check could not run either. One row per worktree, and the "not checked" fact is still reported.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The user asked whether the card and the line under it were the same information. They are two checks with one root cause, so the reader sees the same worktree name twice and has to work out that nothing new is being said. A good layout (REQ-579, render verify findings as compact rows) should not have to hide a redundant producer.

## Detailed Requirements

- The committed-queue-state probe (`worktreeCommittedQueueState`) runs only when `classifyWorktreeLeftover` did not return `worktreeLeftoverMergeStateUndetermined` for that leftover. The other dispositions keep the probe exactly as today.
- The undetermined finding's remedy gains a clause stating that, because git could not resolve the branch, whether the builder committed queue state under `do-work/` could not be checked either. "Unknown never reads as clean" holds: the fact is moved into the finding, not dropped.
- No new `SkippedProbes` entry is written for that leftover's committed-queue-state probe. Skipped entries for other reasons (git missing, integration ref unresolvable) are unchanged.
- The text renderer (`renderVerifyReport`) and the board payload need no change; they print what the report carries.

## Red-Green Proof
**RED prompt/case:** Extend `TestVerifyReportsAnUndeterminedMergeStateSeparately` (or add a sibling) on the same detached-worktree fixture: assert `report.SkippedProbes` contains no entry mentioning `committed-queue-state probe for worktree-agent-REQ-005-detached`, and that the single undetermined finding's remedy mentions that committed queue state could not be checked.
**Why RED now:** `appendWorktreeFindings` runs `worktreeCommittedQueueState` for every leftover regardless of disposition; on a detached or branchless worktree `git diff <ref>...<name>` fails and the failure is appended to `SkippedProbes`, so the fixture today yields one finding plus one skipped line for the same name.
**GREEN when:** The fixture yields exactly one undetermined finding for the name, zero skipped-probe lines for it, the finding is still not `Fixable`, and every other test in `verify_test.go` (in particular the committed-queue-state and uncommitted-queue-state cases around line 1087 and 1118) still passes.
**Validation:** Inferred during capture from the user's question ("are these finding duplicated?") and their "ok, do it" to the proposal that named this as D4.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`): matches on "Changing queue-kanban model" and its family `unknown-reads-as-clean` is the exact rule this REQ must keep. Over the 2000-token budget on its own; the rule is restated in Detailed Requirements instead.

*Source: "are these finding duplicated? the big box and the small line is the same info?" / "ok, do it"*
