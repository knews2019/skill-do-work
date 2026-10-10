# Integrator brief: REQ-696 (stale-wording sweep from the last run's wave-end review)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-696.

- REQ id: REQ-696. Working REQ path: `do-work/working/REQ-696-stale-wording-sweep-wave-end-review.md`. UR: UR-154 (`do-work/user-requests/UR-154/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-192201/` (run id `work-2026-10-10-192201`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-192201/REQ-696-handback.md`. Read it in full: file manifest, P-A-U text, proof record, Decisions, Discovered Tasks, proposed CHANGELOG entry, proposed lesson bullet.
- Operative name (branch = worktree basename): `worktree-agent-REQ-696-stale-wording-sweep`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-696-stale-wording-sweep`. Base `bd56c4b0`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T19:26:48Z`. The hand-back landed a while ago, so per fan-out-reference "Landed hand-back" skip the builder-work event and say so in `## Timing` notes; stamp `builder_handback_at` from the builder commit's committer date (`git log -1 --format=%cI worktree-agent-REQ-696-stale-wording-sweep`, converted to UTC Z).
- Route A, tdd: false, impact-rule-change, effort-mechanical. Test-gate probe: `do-work/runs/work-2026-10-10-192201/REQ-696-probe.sh`.
- Integration order: first integrator of this run (hand-back order). Re-read VERSION before the payloads (0.305.113 at run start, expect 0.305.114).
- Not the wave's last successful integration. Do not pass the wave-end fact. Your diff only corrects restatements; the reviewer checks the six sites read true against current code, and the builder's two report-only finds (capture-reference.md:196, forensics-guide.md:3).
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, worktree cleanup of other REQs); the coordinator does after you report.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-696-stale-wording-sweep-wave-end-review.md`.
- Changelog: the builder's proposed entry is in the hand-back; rewrite its title so it says what shipped in plain words, verify it is unused in `CHANGELOG.md`.
- UR-154 stays open: archive flat. Builder proposed no lesson.
- Heavy: staged-skills must run (it covers the edited staged-skills-contract.sh message and fixture header); a skip is not a pass.
- Prose only: a Quick-scan review is enough (Route A).
