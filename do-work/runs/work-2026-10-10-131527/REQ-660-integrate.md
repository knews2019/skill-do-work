# Integrator brief: REQ-660 (worktree new, status, merge and cleanup command)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-660.

- REQ id: REQ-660. Working REQ path: `do-work/working/REQ-660-worktree-lifecycle-command.md`. UR: UR-145 (`do-work/user-requests/UR-145/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-131527/` (run id `work-2026-10-10-131527`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-131527/REQ-660-handback.md`. Read it in full: file manifest, P-A-U text, proof record, Decisions, Discovered Tasks, proposed CHANGELOG entry, proposed lesson bullet.
- Operative name (branch = worktree basename): `worktree-agent-REQ-660-worktree-lifecycle-command`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-660-worktree-lifecycle-command`. Base `bd56c4b0`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T13:23:49Z`. The hand-back landed a while ago, so per fan-out-reference "Landed hand-back" skip the builder-work event and say so in `## Timing` notes; stamp `builder_handback_at` from the builder commit's committer date (`git log -1 --format=%cI worktree-agent-REQ-660-worktree-lifecycle-command`, converted to UTC Z).
- Route B, tdd: true, impact-user-visible, effort-substantive. Test-gate probe: `do-work/runs/work-2026-10-10-131527/REQ-660-probe.sh`.
- Integration order: you are the first integrator of this run. Re-read VERSION before the payloads (0.305.101 at run start, so expect 0.305.102). Later integrators merge on top of you in hand-back order.
- Not the wave's last successful integration. Do not pass the wave-end fact. Your reviewer runs the Restatement Sweep only if your diff redefines something other text restates (fan-out-reference.md now points at the command: check every other place that restates the worktree add/merge/cleanup argv).
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, worktree cleanup of other REQs); the coordinator does after you report.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-660-worktree-lifecycle-command.md`.
- Changelog: the builder's proposed entry is in the hand-back; rewrite its title so it says what shipped in plain words, verify it is unused in `CHANGELOG.md`.
- Lessons: the builder proposed a bullet for `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`; edit it for truth, link the flat archive path. UR-145 stays open (REQ-659 and REQ-658 follow); REQ-658 closes it and re-points your bullet.
- Discovered task from the builder (cleanup Pass 5 counts the symlinks `worktree new` creates as uncommitted work): judge it; it is report only unless your review finds the new command unusable without it.
