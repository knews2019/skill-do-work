---
id: REQ-756
title: 'Addendum: Detect committed and ignored outside-bundle changes'
status: claimed
created_at: 2026-11-20T01:30:31Z
user_request: UR-188
addendum_to: REQ-754
related: [REQ-755]
domain: backend
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-releases.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-11-20T16:14:14Z
---
# Addendum: Detect committed and ignored outside-bundle changes

## What
Repair subprocess second-opinion checks so an outside edit cannot disappear from evidence merely because the backend committed it or Git ignores its path. Existing git status checks report no outside_paths in both cases and accept the run. These are inherited gaps that the worktree implementation retained.

## Requirements
- Detect and name tracked outside-bundle edits even after the backend commits them inside the detached worktree.
- Detect and name new ignored paths outside the isolated bundle.
- Compare against evidence captured before probes/generation, not mutable post-backend HEAD or ignore rules alone.
- Apply the same condition after model probes and generation, while preserving allowed bundle mutations and unrelated real-checkout dirt behavior.
- Keep existing-output preservation, bounded cancellation/process cleanup, and refusal-page recovery.
- Cover unchanged ignored baseline inputs and both staged/unstaged ordinary outside edits as controls where relevant to the chosen implementation.

## Red-Green Proof
**RED prompt/case:** `TestSecondOpinionWorktreeOutsideSnapshot` in the toolboxcommands Go harness, starting with the preserved `TestREQ754ReviewOutsideChanges` cases for committed README edits and ignored scratch-output.txt creation.
**Why RED now:** Both cases return state=ran, judge=pass, banner=inserted and outside_paths=[] against 4e48c2ef and the pre-change source.
**GREEN when:** Both fail with the offending relative path, leave real outside files unchanged, retain permitted finished output, and remove the owned worktree. Positive normal/dirty-real-checkout cases still succeed.
**Validation:** User authorized both reproduced fixes and execution after the review and retrospective.

## Prior Implementation
REQ-754 implemented detached temporary worktrees for subprocess second opinions, current-bundle copying, exclusive import and cleanup after process reaping. It removed real-checkout outside-status checks. Commits 60186a38 and e7bd242b were merged as d1156bc6, released as 0.331.2. The helper is `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_worktree.go`; callers and checks are in `report_second_opinion.go`.

## Constraints
Keep the current uncommitted bundle contents and unrelated real-checkout changes intact. Keep backend permission flags and session-mode behavior. Retain permitted finished refusal output in the real bundle. A worktree is not a security boundary. No new dependencies or paid live-model calls are required.

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 25337 tokens, slugged partial, exceeds the 2000-token capture budget. Relevant families include canonical-worktree-cleanup and final-boundary-identity. Read touch-required sections during implementation.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** Read listed prime files and agent rules, then record the focused approach.
- [ ] **[APPLY]:** Implement the captured repair within its declared scope.
- [ ] **[UNIFY]:** Review each changed file, run formatting and focused checks, and record evidence.

## Context
No pending or pending-answers request in any UR shares this root cause. This is explicit user promotion of a report-only review finding, not an automatic critical follow-up.

## Assets
Original probes and observed results are preserved under `do-work/user-requests/UR-188/assets/`.
