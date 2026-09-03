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
route: B
write_set:
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
  - skills/do-work/actions/forensics.md
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-09-03T00:59:54Z
  basis:
    - Route B
    - 3-file write set
    - 2 subsystems involved
    - 6 acceptance criteria
    - cross-route regression gates
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

---

## Triage

**Route: B** - Medium

**Reasoning:** The correction is precisely specified and `## Prior Implementation` names the exact functions, but the two evidence sources the fix must consult — worktree dirtiness and the REQ's own pipeline state — are not currently reachable from `appendWorktreeFindings`, so where they come from had to be discovered.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

`classifyWorktreeMergeState` (`verify.go:819`) runs exactly one probe — `git -C repoRoot merge-base --is-ancestor <branch> HEAD` — and collapses everything else into three states. `routeWorktreeLeftover` (`verify.go:851`) then maps `worktreeMergeStateMerged` to `Fixable: true` with the remedy "the branch is already contained in HEAD, so nothing is lost". That single ancestry bit is the whole basis for both claims, which is exactly the defect: ancestry says the *commits* are safe, never that the *worktree* is.

Two facts the repository already records, neither of them a liveness signal:

1. **Worktree dirtiness** — `appendWorktreeFindings` already holds `worktreePath` from `listWorktreeAgentWorktrees` (it uses it only for `locationDetail`). `git -C <worktreePath> status --porcelain --untracked-files=all` answers whether uncommitted work is present. This is the REQ-412 case: cleanup Pass 5's non-forced `git worktree remove` would itself refuse this worktree, so calling it mechanically fixable contradicts the very command the remedy names.
2. **Unfinished pipeline state** — a `worktree-agent-REQ-NNN-*` name carries its REQ id, and a REQ still in flight is exactly the one sitting in `do-work/working/`. That is the REQ-436 case: clean, merged, and still owned by a run that has not reached review or remediation.

Both are ordinary repository reads, so the no-liveness-signal constraint from REQ-073 holds: neither asks whether a process is alive, only what the repository already says.

`Fixable`'s doc comment (`verify.go:66`) defines it as "`do-work cleanup` can mechanically resolve it", and `routeWorktreeLeftover`'s own comment says anything landing on Pass 5's consent-gated path must not be advertised otherwise. The fix is to make the merged branch state necessary but not sufficient.

*Generated in-session (single-pass discovery)*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — make merged-ness necessary but not sufficient; add the dirtiness and in-flight evidence and route both to present-and-non-fixable
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modify) — real-Git fixtures for the REQ-412 and REQ-436 cases plus the still-fixable clean finished case
- `skills/do-work/actions/forensics.md` (modify) — keep the rendered category and remedy text consistent with the corrected classifier

**Files I will NOT touch:** `skills/do-work/actions/cleanup.md` (Pass 5's consent-gated behavior is already correct — this REQ stops verify from contradicting it), and anything adding a lock, heartbeat, PID probe, or mtime heuristic.

**Acceptance criteria (restated from REQ):**
- [ ] No heartbeat, lock, PID probe, mtime heuristic, or claim registry is introduced
- [ ] Ordinary worktree dirtiness and unfinished pipeline state are distinguished from clean merged residue using evidence the repository already records
- [ ] A dirty or unfinished-run builder worktree is reported as present and non-fixable, with no `cleanup can fix` marker and no "nothing is lost" claim
- [ ] Genuinely finished, clean, merged residue is still reported as a mechanically fixable leftover
- [ ] `verify` stays read-only and still protects developer-owned worktrees outside the `worktree-agent-*` convention
- [ ] Rendered category/remedy text and user-facing verify documentation match the corrected classifier
