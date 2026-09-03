---
id: REQ-458
title: 'Addendum: classify active worktrees as present and non-fixable'
status: claimed
created_at: 2026-08-31T21:38:14Z
user_request: UR-086
addendum_to: REQ-083
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-09-03T00:59:09Z
---

# Addendum: Classify Active Worktrees as Present and Non-Fixable

## What

Correct REQ-083 (Verify reports every builder worktree as a fixable orphan, including active and unmerged ones) so a branch being merged into the integration branch is not, by itself, enough to call its worktree a leftover or mechanically fixable. A dirty worktree or a worktree belonging to an unfinished run must be reported as present and non-fixable; only clean merged residue from finished work may be reported as a fixable leftover.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

This corrects the incomplete merged-state classification introduced by REQ-083. The board screenshot showed two `merged-worktree-leftover` findings marked `cleanup can fix`:

- REQ-412 had a merged branch tip while its builder worktree contained uncommitted implementation changes. The current merge-base-only classifier therefore claimed “nothing is lost” even though cleanup's non-forced removal would refuse the dirty worktree.
- REQ-436 remained `claimed` before review and possible remediation had completed. Its clean merged worktree still belonged to the active pipeline; normal worktree cleanup occurs only after the REQ reaches its final path.

The accepted validation finding recorded `Surface-cost: N/A` because this is a direct classification correction, not a new guard, retry, fallback, or warning apparatus.

## Prior Implementation

REQ-083 added `classifyWorktreeMergeState` and `routeWorktreeLeftover`, split the old orphan category into merged, unmerged, and undetermined categories, and set `Fixable: true` for the merged branch state. Its implementation and fixture tests now live under `skills/do-work-board/tools/queue-kanban/verify.go` and `verify_test.go`; the accompanying forensics description lives under `skills/do-work/actions/forensics.md`. The recorded implementation commit is `f6c1514`.

## Requirements

- Preserve the existing no-liveness-signal decision: do not add a heartbeat, lock, PID probe, mtime heuristic, claim registry, or other process-liveness guess.
- Distinguish ordinary worktree dirtiness and unfinished pipeline state from clean merged residue using evidence the repository already records.
- Report dirty or unfinished-run builder worktrees as present and non-fixable. Do not describe them as leftovers, advertise `cleanup can fix`, or claim that nothing can be lost.
- Continue reporting genuinely finished, clean, merged residue as a mechanically fixable leftover.
- Keep `verify` read-only and preserve the current protection for developer-owned worktrees outside the `worktree-agent-*` naming convention.
- Keep all rendered category/remedy text and any user-facing verify documentation consistent with the corrected classifier.

## Red-Green Proof

**RED prompt/case:** Extend the real-Git worktree fixture in `skills/do-work-board/tools/queue-kanban/verify_test.go` with (1) a branch already merged into the integration branch whose worktree has an uncommitted ordinary source-file change and (2) a clean merged worktree whose matching REQ is still `claimed` before review/remediation finishes. Run the verify probes and inspect both findings.

**Why RED now:** Both cases route solely from `git merge-base --is-ancestor`, so each currently becomes `merged-worktree-leftover [fixable]` with the claim that cleanup is mechanical and nothing is lost.

**GREEN when:** Both regression cases report the worktree as present and non-fixable, with no `cleanup can fix` marker and no “nothing is lost” claim; a separate clean merged worktree from finished work remains a fixable leftover. The regression test must close both REQ-412 and REQ-436 instances without introducing a liveness signal.

**Validation:** User confirmed through the accepted `do-work validate-feedback` finding and this capture request.

## Assets

- `do-work/user-requests/UR-086/assets/REQ-458-screenshot-1-active-worktrees-labelled-leftovers.png` — queue board generated during the active run. Its Verify Findings strip shows two cards, for REQ-412 and REQ-436, both labelled `MERGED-WORKTREE-LEFTOVER`, both marked `cleanup can fix`, and both stating that the branch is contained in `HEAD` so nothing is lost.

---
*Source: user-approved `do-work validate-feedback` finding; full verbatim input in `do-work/user-requests/UR-086/input.md`.*
